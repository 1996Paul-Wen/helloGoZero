<template>
  <el-dialog
    :model-value="visible"
    @update:model-value="(val: boolean) => emit('update:visible', val)"
    title=""
    width="420px"
    :close-on-click-modal="false"
    append-to-body
    class="rounded-2xl overflow-hidden encrypt-dialog"
  >
    <div class="text-center pt-2 pb-4">
      <!-- 安全图标 -->
      <div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-amber-50 border border-amber-100 mb-5">
        <Shield class="w-8 h-8 text-amber-500" />
      </div>
      <h3 class="text-lg font-bold text-slate-800 mb-2">设置加密密钥</h3>
      <p class="text-sm text-slate-500 leading-relaxed px-4">
        请输入一个<strong class="text-amber-600">加密密钥</strong>，用于 AES 加密保护您的密码。<br/>
        <span class="text-xs text-amber-500">⚠️ 请妥善保管此密钥，忘记后将无法恢复密码。</span>
      </p>
    </div>

    <div class="space-y-4 pb-2">
      <label class="block text-sm font-medium text-slate-700 ml-1">加密密钥</label>
      <div class="relative">
        <input
          v-model="encryptionKey"
          :type="showKey ? 'text' : 'password'"
          placeholder="请输入加密密钥"
          autofocus
          class="w-full px-4 py-3 pr-11 rounded-xl border border-amber-200 bg-amber-50/30 focus:border-amber-400 focus:ring-2 focus:ring-amber-100 outline-none transition-all text-sm"
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

      <div class="flex gap-3 pt-2">
        <button
          @click="emit('update:visible', false)"
          class="flex-1 py-2.5 rounded-xl border border-slate-200 text-slate-600 hover:bg-slate-50 transition-colors cursor-pointer"
        >
          取消
        </button>
        <button
          @click="handleConfirm"
          :disabled="!encryptionKey.trim()"
          class="flex-1 py-2.5 rounded-xl bg-amber-500 text-white font-medium hover:bg-amber-600 disabled:opacity-40 disabled:cursor-not-allowed transition-colors cursor-pointer"
        >
          确认加密
        </button>
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { Shield, Eye, EyeOff } from 'lucide-vue-next'

const props = defineProps<{ visible: boolean; description?: string }>()
const emit = defineEmits<{
  (e: 'update:visible', val: boolean): void
  (e: 'confirmed', encryptionKey: string): void
}>()

const encryptionKey = ref('')
const showKey = ref(false)

// 每次打开重置
watch(
  () => props.visible,
  (val) => {
    if (val) {
      encryptionKey.value = ''
      showKey.value = false
    }
  }
)

function handleConfirm() {
  const key = encryptionKey.value.trim()
  if (!key) return
  emit('confirmed', key)
}
</script>
