<template>
  <el-dialog
    :model-value="visible"
    @update:model-value="(val: boolean) => emit('update:visible', val)"
    title=""
    width="420px"
    :close-on-click-modal="false"
    append-to-body
    class="rounded-2xl overflow-hidden"
  >
    <div class="text-center pt-2 pb-4">
      <!-- 解锁图标 -->
      <div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-blue-50 border border-blue-100 mb-5">
        <Unlock class="w-8 h-8 text-blue-500" />
      </div>
      <h3 class="text-lg font-bold text-slate-800 mb-2">解密查看密码</h3>
      <p class="text-sm text-slate-500 leading-relaxed px-4">
        正在解锁：<strong class="text-blue-600 truncate inline-block max-w-[220px] align-bottom">{{ description || '该条记录' }}</strong><br/>
        <span class="text-xs text-slate-400 mt-1 block">请输入您设置的加密密钥以解密密码</span>
      </p>
    </div>

    <div class="space-y-4 pb-2">
      <label class="block text-sm font-medium text-slate-700 ml-1">解密密钥</label>
      <div class="relative">
        <input
          v-model="decryptionKey"
          :type="showKey ? 'text' : 'password'"
          placeholder="请输入解密密钥"
          autofocus
          class="w-full px-4 py-3 pr-11 rounded-xl border border-blue-200 bg-blue-50/30 focus:border-blue-400 focus:ring-2 focus:ring-blue-100 outline-none transition-all text-sm"
          @keyup.enter="handleConfirm"
        />
        <button
          type="button"
          @click="showKey = !showKey"
          class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition-colors cursor-pointer"
        >
          <Eye v-if="!showKey" class="w-5 h-5" />
          <EyeOff v-else class="w-5 h-5" />
        </button>
      </div>
      <p v-if="errorMsg" class="text-xs text-red-500 ml-1 flex items-center gap-1">
        <AlertCircle class="w-3.5 h-3.5" />
        {{ errorMsg }}
      </p>

      <div class="flex gap-3 pt-2">
        <button
          @click="emit('update:visible', false)"
          class="flex-1 py-2.5 rounded-xl border border-slate-200 text-slate-600 hover:bg-slate-50 transition-colors cursor-pointer"
        >
          取消
        </button>
        <button
          @click="handleConfirm"
          :disabled="!decryptionKey.trim()"
          class="flex-1 py-2.5 rounded-xl bg-blue-500 text-white font-medium hover:bg-blue-600 disabled:opacity-40 disabled:cursor-not-allowed transition-colors cursor-pointer"
        >
          解密查看
        </button>
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { Unlock, Eye, EyeOff, AlertCircle } from 'lucide-vue-next'
import { decrypt } from '@/utils/crypto'

const props = defineProps<{
  visible: boolean
  description?: string
}>()
const emit = defineEmits<{
  (e: 'update:visible', val: boolean): void
  (e: 'confirmed', plainText: string): void
}>()

const decryptionKey = ref('')
const showKey = ref(false)
const errorMsg = ref('')

// 每次打开重置
watch(
  () => props.visible,
  (val) => {
    if (val) {
      decryptionKey.value = ''
      showKey.value = false
      errorMsg.value = ''
    }
  }
)

function handleConfirm() {
  const key = decryptionKey.value.trim()
  if (!key) return

  errorMsg.value = ''

  // 这里我们只传递密钥给父组件，由父组件调用 decrypt
  // 因为父组件持有 cipherText 数据
  emit('confirmed', key)
}
</script>
