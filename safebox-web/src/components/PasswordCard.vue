<template>
  <div class="card-hover bg-white rounded-2xl p-6 shadow-sm border border-slate-100">
    <!-- 卡片头部：描述 + 操作按钮 -->
    <div class="flex items-start justify-between gap-2 mb-4">
      <h3 class="font-semibold text-slate-800 text-base leading-snug line-clamp-2 flex-1 min-w-0">
        {{ props.item.Description }}
      </h3>
      <div class="flex items-center gap-1 shrink-0">
        <button
          @click="openEditDialog"
          class="p-1.5 rounded-lg text-slate-400 hover:text-blue-500 hover:bg-blue-50 transition-colors cursor-pointer"
          title="修改"
        >
          <Pencil class="w-4 h-4" />
        </button>
        <button
          @click="handleDelete"
          class="p-1.5 rounded-lg text-slate-400 hover:text-red-500 hover:bg-red-50 transition-colors cursor-pointer"
          title="删除"
        >
          <Trash2 class="w-4 h-4" />
        </button>
      </div>
    </div>

    <!-- 卡片主体信息 -->
    <div class="space-y-3">
      <!-- 用户名 -->
      <div class="flex items-start gap-2">
        <User class="w-4 h-4 text-slate-400 mt-0.5 shrink-0" />
        <div class="min-w-0">
          <span class="text-xs text-slate-400 block">用户名</span>
          <span class="text-sm text-slate-700 break-all">{{ props.item.Username }}</span>
        </div>
      </div>

      <!-- 密码（脱敏 + 解密） -->
      <div class="flex items-start gap-2">
        <Key class="w-4 h-4 text-slate-400 mt-0.5 shrink-0" />
        <div class="min-w-0 flex-1">
          <span class="text-xs text-slate-400 block">密码</span>
          <div v-if="!decryptedVisible" class="flex items-center gap-2">
            <span class="text-sm text-slate-500 tracking-widest font-mono">••••••••</span>
            <button
              @click="showDecrypt = true"
              class="shrink-0 p-1.5 rounded-lg text-slate-400 hover:text-blue-500 hover:bg-blue-50 transition-colors cursor-pointer"
              title="查看密码"
            >
              <Eye class="w-4 h-4" />
            </button>
          </div>
          <div v-else class="flex items-center gap-2 group">
            <input
              ref="pwdInputRef"
              :type="showPlainPwd ? 'text' : 'password'"
              :value="decryptedPwd"
              readonly
              class="flex-1 text-sm text-emerald-600 font-mono bg-emerald-50/50 px-2 py-1 rounded-lg border border-emerald-100 outline-none min-w-0"
            />
            <button
              @click="showPlainPwd = !showPlainPwd"
              class="shrink-0 p-1.5 rounded-lg text-slate-400 hover:text-slate-600 transition-colors cursor-pointer"
              title="切换显示"
            >
              <EyeOff v-if="showPlainPwd" class="w-4 h-4" />
              <Eye v-else class="w-4 h-4" />
            </button>
            <button
              @click="copyPassword(decryptedPwd)"
              class="shrink-0 p-1.5 rounded-lg text-slate-400 hover:text-emerald-500 hover:bg-emerald-50 transition-colors cursor-pointer"
              title="复制密码"
            >
              <ClipboardCopy class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部时间信息 -->
    <div class="mt-4 pt-4 border-t border-slate-100 flex items-center gap-1.5 text-xs text-slate-400">
      <Clock class="w-3.5 h-3.5" />
      {{ formatTime(props.item.CreateTime) }}
    </div>

    <!-- 解密密钥弹窗 -->
    <DecryptDialog
      v-model:visible="showDecrypt"
      :description="props.item.Description"
      @confirmed="handleDecrypted"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { User, Key, Eye, EyeOff, Clock, ClipboardCopy, Pencil, Trash2 } from 'lucide-vue-next'
import { decrypt } from '@/utils/crypto'
import type { PasswordItem } from '@/types'
import DecryptDialog from './DecryptDialog.vue'

const props = defineProps<{ item: PasswordItem }>()
const emit = defineEmits<{
  (e: 'deleted', id: number): void
  (e: 'edit', item: PasswordItem): void
}>()

const showDecrypt = ref(false)
const decryptedVisible = ref(false)
const decryptedPwd = ref('')
const showPlainPwd = ref(false)
const pwdInputRef = ref<HTMLInputElement>()

function handleDecrypted(decryptionKey: string) {
  try {
    const plainText = decrypt(props.item.Password, decryptionKey, props.item.Description, props.item.Username)
    decryptedPwd.value = plainText
    decryptedVisible.value = true
    showDecrypt.value = false
    ElMessage.success('解密成功')
  } catch (error: any) {
    ElMessage.error(error.message || '解密失败')
  }
}

async function copyPassword(text: string) {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch {
    const el = pwdInputRef.value
    if (el) {
      el.select()
      document.execCommand('copy')
    }
    ElMessage.success('已复制到剪贴板')
  }
}

function openEditDialog() {
  emit('edit', props.item)
}

async function handleDelete() {
  try {
    await ElMessageBox.confirm(
      `确定要删除「${props.item.Description}」吗？此操作不可恢复。`,
      '删除确认',
      { confirmButtonText: '删除', cancelButtonText: '取消', type: 'warning' }
    )
    const { deletePassword } = await import('@/api/pwdManage')
    await deletePassword(props.item.Id)
    ElMessage.success('删除成功')
    emit('deleted', props.item.Id)
  } catch {
    // 用户取消
  }
}

function formatTime(timeStr: string): string {
  if (!timeStr) return ''
  const d = new Date(timeStr)
  return d.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}
</script>
