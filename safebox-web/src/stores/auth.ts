import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UserInfo } from '@/types'
import { describeUser } from '@/api/user'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const userInfo = ref<UserInfo | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const username = computed(() => userInfo.value?.Username || '')

  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  async function fetchUserInfo() {
    try {
      const res = await describeUser()
      if (res.Data) {
        userInfo.value = res.Data
      }
    } catch {
      console.error('获取用户信息失败')
    }
  }

  async function login(newToken: string) {
    setToken(newToken)
    await fetchUserInfo()
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    username,
    setToken,
    login,
    logout,
    fetchUserInfo,
  }
})
