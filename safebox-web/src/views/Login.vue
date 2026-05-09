<template>
  <div class="min-h-screen gradient-bg flex items-center justify-center p-4">
    <div class="glass rounded-3xl shadow-2xl w-full max-w-md p-10">
      <!-- Logo & 标题 -->
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-blue-500 to-cyan-400 shadow-lg shadow-blue-200 mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
          </svg>
        </div>
        <h1 class="text-2xl font-bold text-slate-800">SafeBox 密码管家</h1>
        <p class="text-slate-500 mt-1 text-sm">安全可靠地管理您的密码</p>
      </div>

      <!-- 登录表单 -->
      <el-form :model="form" :rules="rules" ref="formRef" @submit.prevent="handleLogin">
        <el-form-item prop="name" class="!mb-5">
          <div class="flex items-center gap-3">
            <label class="text-sm font-medium text-slate-700 w-14 text-right shrink-0">用户名</label>
            <div class="relative flex-1">
              <input
                v-model="form.name"
                type="text"
                placeholder="请输入用户名"
                class="w-full px-4 pr-11 py-2.5 rounded-xl border border-slate-200 bg-white/60 focus:bg-white focus:border-blue-400 focus:ring-2 focus:ring-blue-100 outline-none transition-all text-sm box-border"
                autocomplete="username"
              />
              <span class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-300 pointer-events-none select-none">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z" /></svg>
              </span>
            </div>
          </div>
        </el-form-item>

        <el-form-item prop="password" class="!mb-5">
          <div class="flex items-center gap-3">
            <label class="text-sm font-medium text-slate-700 w-14 text-right shrink-0">密码</label>
            <div class="relative flex-1">
              <input
                v-model="form.password"
                :type="showPwd ? 'text' : 'password'"
                placeholder="请输入密码"
                class="w-full px-4 pr-11 py-2.5 rounded-xl border border-slate-200 bg-white/60 focus:bg-white focus:border-blue-400 focus:ring-2 focus:ring-blue-100 outline-none transition-all text-sm box-border"
                autocomplete="current-password"
                @keyup.enter="handleLogin"
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
          </div>
        </el-form-item>

        <button
          type="submit"
          :disabled="loading"
          class="w-full py-3 rounded-xl bg-gradient-to-r from-blue-500 to-cyan-400 text-white font-semibold shadow-lg shadow-blue-200 hover:shadow-xl hover:shadow-blue-300 hover:-translate-y-0.5 active:translate-y-0 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed mt-2 cursor-pointer"
        >
          {{ loading ? '登录中...' : '登 录' }}
        </button>

        <p class="text-center text-sm text-slate-500 mt-6">
          还没有账号？
          <router-link to="/register" class="text-blue-500 hover:text-blue-600 font-medium">立即注册</router-link>
        </p>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Eye, EyeOff } from 'lucide-vue-next'
import { loginUser } from '@/api/user'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const formRef = ref()
const loading = ref(false)
const showPwd = ref(false)

const form = reactive({ name: '', password: '' })

const rules = {
  name: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 3, message: '密码至少3个字符', trigger: 'blur' },
  ],
}

async function handleLogin() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  const res = await loginUser({ name: form.name, password: form.password })
  if (!res.Data?.token) {
    ElMessage.error('登录失败，未获取到令牌')
    loading.value = false
    return
  }

  authStore.setToken(res.Data.token)

  // fetchUserInfo 失败不影响登录，静默处理
  authStore.fetchUserInfo().catch(() => {})

  const displayName = authStore.userInfo?.Username || form.name
  ElMessage.success(`欢迎回来，${displayName}`)
  loading.value = false

  // 强制跳转并 reload，确保 Pinia store 和路由状态完全刷新
  window.location.href = '/'
}
</script>
