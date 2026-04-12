import axios from 'axios'
import type { ApiResponse } from '@/types'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'
import router from '@/router'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
})

// 请求拦截器：自动附加 Bearer Token
request.interceptors.request.use((config) => {
  const authStore = useAuthStore()
  if (authStore.token) {
    config.headers.Authorization = `Bearer ${authStore.token}`
  }
  return config
})

// 响应拦截器：统一处理错误和 401 跳转
request.interceptors.response.use(
  (response) => {
    const data: ApiResponse = response.data
    if (data.Code !== 0 && data.Msg) {
      ElMessage.error(data.Msg)
      return Promise.reject(new Error(data.Msg))
    }
    return data
  },
  (error) => {
    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      authStore.logout()
      ElMessage.error('登录已过期，请重新登录')
      router.push('/login')
      return Promise.reject(error)
    }
    ElMessage.error(error.message || '网络请求异常')
    return Promise.reject(error)
  }
)

export default request
