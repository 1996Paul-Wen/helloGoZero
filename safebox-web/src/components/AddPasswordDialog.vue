<template>
  <el-dialog
    :model-value="visible"
    @update:model-value="(val: boolean) => emit('update:visible', val)"
    title="新增密码记录"
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

      <el-form-item label="密码" prop="password">
        <div class="relative w-full">
          <input
            v-model="form.password"
            :type="showPwd ? 'text' : 'password'"
            placeholder="请输入要保存的明文密码"
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
        <p class="text-xs text-amber-500 mt-1.5 flex items-center gap-1">
          <ShieldAlert class="w-3.5 h-3.5" />
          保存时将使用 AES 加密保护您的密码
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

  <!-- 加密密钥二级弹窗（移到外层，避免被父弹窗遮挡） -->
  <EncryptDialog
    v-model:visible="showEncryptDialog"
    :description="form.description"
    @confirmed="handleEncrypted"
  />
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Eye, EyeOff, ShieldAlert } from 'lucide-vue-next'
import { savePassword } from '@/api/pwdManage'
import { encrypt } from '@/utils/crypto'
import EncryptDialog from './EncryptDialog.vue'

const props = defineProps<{ visible: boolean }>()
const emit = defineEmits<{
  (e: 'update:visible', val: boolean): void
  (e: 'success'): void
}>()

const formRef = ref()
const submitting = ref(false)
const showPwd = ref(false)
const showEncryptDialog = ref(false)

const form = reactive({
  description: '',
  username: '',
  password: '',
})

const rules = {
  description: [{ required: true, message: '请输入描述信息', trigger: 'blur' }],
  username: [{ required: true, message: '请输入用户名或账号', trigger: 'blur' }],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 1, message: '密码不能为空', trigger: 'blur' },
  ],
}

async function handleSubmit() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  // 打开加密密钥弹窗
  showEncryptDialog.value = true
}

/** 加密密钥确认后的回调 */
function handleEncrypted(encryptionKey: string) {
  doSaveWithEncryption(encryptionKey)
}

async function doSaveWithEncryption(keyStr: string) {
  submitting.value = true
  try {
    const encryptedPassword = encrypt(form.password, keyStr, form.description, form.username)
    await savePassword({
      description: form.description,
      username: form.username,
      password: encryptedPassword,
    })
    ElMessage.success('密码已安全保存')
    emit('update:visible', false)

    // 重置表单
    form.description = ''
    form.username = ''
    form.password = ''

    emit('success')
  } finally {
    submitting.value = false
    showEncryptDialog.value = false
  }
}
</script>
