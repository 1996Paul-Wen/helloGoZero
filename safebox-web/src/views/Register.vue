<template>
  <div class="min-h-screen gradient-bg flex items-center justify-center p-4">
    <div class="glass rounded-3xl shadow-2xl w-full max-w-md p-10">
      <!-- Logo & 标题 -->
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-emerald-500 to-teal-400 shadow-lg shadow-emerald-200 mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
          </svg>
        </div>
        <h1 class="text-2xl font-bold text-slate-800">创建账号</h1>
        <p class="text-slate-500 mt-1 text-sm">开始安全地管理您的密码</p>
      </div>

      <!-- 注册表单 -->
      <el-form :model="form" :rules="rules" ref="formRef" @submit.prevent="handleRegister">
        <el-form-item prop="name">
          <label class="block text-sm font-medium text-slate-700 mb-1.5">用户名</label>
          <input
            v-model="form.name"
            type="text"
            placeholder="请设置用户名"
            class="w-full px-4 py-2.5 rounded-xl border border-slate-200 bg-white/60 focus:bg-white focus:border-emerald-400 focus:ring-2 focus:ring-emerald-100 outline-none transition-all text-sm"
            autocomplete="username"
          />
        </el-form-item>

        <el-form-item prop="password">
          <label class="block text-sm font-medium text-slate-700 mb-1.5">密码</label>
          <div class="relative">
            <input
              v-model="form.password"
              :type="showPwd ? 'text' : 'password'"
              placeholder="请设置密码（至少3位）"
              class="w-full px-4 py-2.5 pr-11 rounded-xl border border-slate-200 bg-white/60 focus:bg-white focus:border-emerald-400 focus:ring-2 focus:ring-emerald-100 outline-none transition-all text-sm"
              autocomplete="new-password"
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
        </el-form-item>

        <el-form-item prop="confirmPassword">
          <label class="block text-sm font-medium text-slate-700 mb-1.5">确认密码</label>
          <input
            v-model="form.confirmPassword"
            type="password"
            placeholder="再次输入密码"
            class="w-full px-4 py-2.5 rounded-xl border border-slate-200 bg-white/60 focus:bg-white focus:border-emerald-400 focus:ring-2 focus:ring-emerald-100 outline-none transition-all text-sm"
            autocomplete="new-password"
          />
        </el-form-item>

        <button
          type="submit"
          :disabled="loading"
          class="w-full py-3 rounded-xl bg-gradient-to-r from-emerald-500 to-teal-400 text-white font-semibold shadow-lg shadow-emerald-200 hover:shadow-xl hover:shadow-emerald-300 hover:-translate-y-0.5 active:translate-y-0 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed mt-2 cursor-pointer"
        >
          {{ loading ? '注册中...' : '注 册' }}
        </button>

        <p class="text-center text-sm text-slate-500 mt-6">
          已有账号？
          <router-link to="/login" class="text-emerald-500 hover:text-emerald-600 font-medium">去登录</router-link>
        </p>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Eye, EyeOff } from 'lucide-vue-next'
import { createUser } from '@/api/user'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)
const showPwd = ref(false)

const form = reactive({
  name: '',
  password: '',
  confirmPassword: '',
})

const validateConfirmPass = (_rule: any, value: string, callback: any) => {
  if (value !== form.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules: FormRules = {
  name: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, message: '用户名至少2个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请设置密码', trigger: 'blur' },
    { min: 3, message: '密码至少3个字符', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validateConfirmPass, trigger: 'blur' },
  ],
}

async function handleRegister() {
  if (!formRef.value) return
  await formRef.value.validate().catch(() => {})

  loading.value = true
  try {
    const res = await createUser({ name: form.name, password: form.password })
    if (res.Data?.UserID) {
      ElMessage.success('注册成功，请登录')
      router.push('/login')
    }
  } finally {
    loading.value = false
  }
}
</script>
