import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { guest: true },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { guest: true },
  },
  {
    path: '/',
    name: 'Dashboard',
    component: () => import('@/views/Dashboard.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/login',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// 全局导航守卫
router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()

  // 需要登录的页面，未登录则跳转登录页
  if (to.meta.requiresAuth && !authStore.token) {
    next({ name: 'Login' })
    return
  }
  // 访客页面（已登录则跳转首页）
  if (to.meta.guest && authStore.token) {
    next({ name: 'Dashboard' })
    return
  }
  next()
})

export default router
