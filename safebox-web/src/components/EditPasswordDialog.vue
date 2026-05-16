<template>
  <el-dialog
    :model-value="visible"
    @update:model-value="(val: boolean) => emit('update:visible', val)"
    title="修改密码记录"
    width="520px"
    :close-on-click-modal="false"
    class="rounded-2xl overflow-hidden"
  >
    <el-form :model="form" :rules="rules" ref="formRef" label-position="top" class="pt-2">
      <el-form-item label="描述" prop="description">
        <input
          v-model="form.description"
          type="text"
          placeholder="例：中信银行手机App登录账号"
          class="w-full px-4 py-2.5 rounded-xl border border-slate-200 bg-white focus:border-blue-400 focus:ring-2 focus:ring-blue-100 outline-none transition-all text-sm"
        />
      </el-form-item>

      <el-form-item label="用户名 / 账号" prop="username">
        <input
          v-model="form.username"
          type="text"
          placeholder="请输入用户名或账号"
          class="w-full px-4 py-2.5 rounded-xl border border-slate-200 bg-white focus:border-blue-400 focus:ring-2 focus:ring-blue-100 outline-none transition-all text-sm"
        />
      </el-form-item>

      <el-form-item label="密码">
        <div class="w-full space-y-2">
          <!-- 当前密码（解密后显示） -->
          <div v-if="decryptedPwd && !changePwdMode" class="flex items-center gap-2 px-4 py-2.5 rounded-xl border border-slate-200 bg-slate-50 text-sm text-slate-500 font-mono">
            {{ decryptedPwd }}
          </div>
          <!-- 修改模式：输入新密码 -->
          <div v-if="changePwdMode" class="relative">
            <input
              v-model="form.password"
              :type="showPwd ? 'text' : 'password'"
              placeholder="请输入新密码"
              class="w-full px-4 py-2.5 pr-11 rounded-xl border border-slate-200 bg-white focus:border-blue-400 focus:ring-2 focus:ring-blue-100 outline-none transition-all text-sm"
            />
            <button
              type="button"
              @click="showPwd = !showPwd"
              class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition-colors cursor-pointer"
            >
              <Eye v-if="!showPwd" class="w-5 h-5" />
              <EyeOff v-else class="w-5 h-5" />
            </button>
          </div>
          <button
            type="button"
            v-if="!changePwdMode"
            @click="showDecrypt = true"
            class="text-sm text-blue-500 hover:text-blue-600 cursor-pointer flex items-center gap-1"
          >
            <Edit class="w-3.5 h-3.5" />
            修改密码
          </button>
          <button
            type="button"
            v-else
            @click="changePwdMode = false; form.password = ''"
            class="text-sm text-slate-400 hover:text-slate-600 cursor-pointer"
          >
            取消修改
          </button>
        </div>
        <p class="text-xs text-amber-500 mt-1.5 flex items-center gap-1">
          <ShieldAlert class="w-3.5 h-3.5" />
          修改密码时需输入加密密钥重新加密
        </p>
      </el-form-item>
    </el-form>

    <template #footer>
      <button
        @click="emit('update:visible', false)"
        class="px-5 py-2 rounded-xl text-slate-500 hover:bg-slate-100 transition-colors cursor-pointer"
      >
        取消
      </button>
      <button
        @click="handleSubmit"
        :disabled="submitting"
        class="px-6 py-2 rounded-xl bg-blue-500 text-white font-medium hover:bg-blue-600 disabled:opacity-50 transition-colors cursor-pointer"
      >
        {{ submitting ? '保存中...' : '保 存' }}
      </button>
    </template>
  </el-dialog>

  <!-- 解密密钥弹窗（移到外层，避免被父弹窗遮挡） -->
  <DecryptDialog
    v-model:visible="showDecrypt"
    :description="props.item.Description"
    @confirmed="handleDecrypted"
  />

  <!-- 加密密钥弹窗（移到外层，避免被父弹窗遮挡） -->
  <EncryptDialog
    v-model:visible="showEncryptDialog"
    :description="form.description"
    @confirmed="handleEncrypted"
  />
</template>

<script setup lang="ts">
import { reactive, ref, watch, computed } from 'vue'
import { Eye, EyeOff, ShieldAlert, Edit } from 'lucide-vue-next'
import { ElMessage } from 'element-plus'
import { decrypt, encrypt } from '@/utils/crypto'
import DecryptDialog from './DecryptDialog.vue'
import EncryptDialog from './EncryptDialog.vue'
import type { PasswordItem } from '@/types'

const props = defineProps<{
  visible: boolean
  item: PasswordItem
}>()
const emit = defineEmits<{
  (e: 'update:visible', val: boolean): void
  (e: 'success'): void
}>()

const formRef = ref()
const submitting = ref(false)
const showPwd = ref(false)
const showDecrypt = ref(false)
const showEncryptDialog = ref(false)
const changePwdMode = ref(false)
const decryptedPwd = ref('')
const decryptingForReEncrypt = ref(false) // 是否在为"仅修改描述/用户名"而解密

const form = reactive({
  description: '',
  username: '',
  password: '',
})

const rules = {
  description: [{ required: true, message: '请输入描述信息', trigger: 'blur' }],
  username: [{ required: true, message: '请输入用户名或账号', trigger: 'blur' }],
}

// 记录原始值，判断 description/username 是否变更
const originalDesc = ref('')
const originalUser = ref('')

// description 或 username 是否有变更（不含 password 变更）
const descOrUserChanged = computed(
  () => form.description !== originalDesc.value || form.username !== originalUser.value
)

// 每次打开弹窗，用 item 填充表单
watch(
  () => props.visible,
  (val) => {
    if (val) {
      form.description = props.item.Description
      form.username = props.item.Username
      form.password = ''
      changePwdMode.value = false
      decryptedPwd.value = ''
      showPwd.value = false
      decryptingForReEncrypt.value = false
      // 保存原始值用于比较
      originalDesc.value = props.item.Description
      originalUser.value = props.item.Username
    }
  }
)

function handleDecrypted(decryptionKey: string) {
  try {
    decryptedPwd.value = decrypt(props.item.Password, decryptionKey, props.item.Description, props.item.Username)
    form.password = decryptedPwd.value
    changePwdMode.value = true
    showDecrypt.value = false
    decryptingForReEncrypt.value = false
  } catch (error: any) {
    ElMessage.error(error.message || '解密失败')
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  // 场景一：什么都没变（包括 password 未改动），直接用原密文保存
  if (!changePwdMode.value && !descOrUserChanged.value) {
    await doSaveDirect()
    return
  }

  // 场景二：改了密码（有明文），直接走加密流程
  if (changePwdMode.value && form.password) {
    showEncryptDialog.value = true
    return
  }

  // 场景三：只改了 description/username，密码未变
  // 需要：旧 key 解密 → 新 key 重新加密
  if (descOrUserChanged.value && !changePwdMode.value) {
    decryptingForReEncrypt.value = true
    showDecrypt.value = true
    return
  }
}

// 解密完成后（场景三：仅修改描述/用户名）
// handleDecrypted 中 decryptingForReEncrypt 为 true 时，复用解密结果走加密流程
watch(decryptedPwd, (val) => {
  if (decryptingForReEncrypt.value && val) {
    decryptingForReEncrypt.value = false
    // decryptedPwd 已有明文（用旧 key 解的），form.password 也在
    // 直接打开加密弹窗，用新 key 重新加密
    showEncryptDialog.value = true
  }
})

function handleEncrypted(encryptionKey: string) {
  doSaveWithEncryption(encryptionKey)
}

async function doSaveDirect() {
  submitting.value = true
  try {
    const { updatePassword } = await import('@/api/pwdManage')
    await updatePassword({
      id: props.item.Id,
      description: form.description,
      username: form.username,
      password: props.item.Password,
    })
    ElMessage.success('已更新')
    emit('update:visible', false)
    emit('success')
  } catch (error: any) {
    ElMessage.error(error.message || '保存失败')
  } finally {
    submitting.value = false
  }
}

async function doSaveWithEncryption(keyStr: string) {
  submitting.value = true
  try {
    const { updatePassword } = await import('@/api/pwdManage')
    const passwordToSave = encrypt(form.password, keyStr, form.description, form.username)
    await updatePassword({
      id: props.item.Id,
      description: form.description,
      username: form.username,
      password: passwordToSave,
    })
    ElMessage.success('密码已更新')
    emit('update:visible', false)
    emit('success')
  } catch (error: any) {
    ElMessage.error(error.message || '保存失败')
  } finally {
    submitting.value = false
    showEncryptDialog.value = false
  }
}
</script>
