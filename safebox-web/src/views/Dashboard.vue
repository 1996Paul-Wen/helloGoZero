<template>
  <div class="min-h-screen bg-slate-50">
    <!-- 顶部导航栏 -->
    <header class="fixed top-0 left-0 right-0 h-16 bg-white/80 backdrop-blur-md border-b border-slate-200/60 z-40">
      <div class="max-w-6xl mx-auto h-full flex items-center justify-between px-6">
        <div class="flex items-center gap-3">
          <div class="w-9 h-9 rounded-xl bg-gradient-to-br from-blue-500 to-cyan-400 flex items-center justify-center shadow-md">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
          </div>
          <h1 class="text-lg font-bold text-slate-800">SafeBox</h1>
        </div>

        <div class="flex items-center gap-4">
          <span class="text-sm text-slate-500">你好，<strong class="text-slate-700">{{ authStore.username }}</strong></span>
          <button
            @click="handleLogout"
            class="px-4 py-1.5 rounded-lg text-sm text-slate-500 hover:text-red-500 hover:bg-red-50 transition-all cursor-pointer"
          >
            退出登录
          </button>
        </div>
      </div>
    </header>

    <!-- 主内容区 -->
    <main class="pt-[72px] pb-12 px-4 sm:px-6 max-w-6xl mx-auto">
      <!-- 搜索与操作栏 -->
      <div class="flex flex-col sm:flex-row gap-4 mt-6 mb-8">
        <div class="relative flex-1">
          <Search class="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜索密码记录（如：中信、银行、app）..."
            class="w-full pl-12 pr-4 py-3.5 rounded-2xl border border-slate-200 bg-white focus:border-blue-400 focus:ring-2 focus:ring-blue-100 outline-none transition-all text-sm shadow-sm"
            @keyup.enter="handleSearch"
          />
        </div>
        <div class="flex gap-3">
          <button
            @click="handleSearch"
            :disabled="loading"
            class="flex items-center justify-center gap-2 px-5 py-3.5 rounded-2xl border border-slate-200 bg-white text-slate-600 font-medium hover:bg-slate-50 hover:border-slate-300 transition-all whitespace-nowrap cursor-pointer disabled:opacity-50"
          >
            <List class="w-5 h-5" />
            查询已存密码
          </button>
          <button
            @click="showAddDialog = true"
            class="flex items-center justify-center gap-2 px-6 py-3.5 rounded-2xl bg-gradient-to-r from-blue-500 to-cyan-400 text-white font-semibold shadow-lg shadow-blue-200 hover:shadow-xl hover:shadow-blue-300 hover:-translate-y-0.5 active:translate-y-0 transition-all whitespace-nowrap cursor-pointer"
          >
            <Plus class="w-5 h-5" />
            新增密码
          </button>
        </div>
      </div>

      <!-- 加载中 -->
      <div v-if="loading" class="flex justify-center py-20">
        <div class="w-8 h-8 border-2 border-blue-400 border-t-transparent rounded-full animate-spin" />
      </div>

      <!-- 空状态 -->
      <div v-else-if="passwords.length === 0" class="text-center py-20">
        <div class="w-24 h-24 mx-auto mb-6 rounded-full bg-slate-100 flex items-center justify-center">
          <Lock class="w-10 h-10 text-slate-300" />
        </div>
        <h3 class="text-lg font-semibold text-slate-600 mb-2">{{ searchQuery ? '未找到匹配结果' : '还没有密码记录' }}</h3>
        <p class="text-sm text-slate-400 mb-6">{{ searchQuery ? '试试其他关键词' : '点击「新增密码」添加您的第一条密码记录' }}</p>
        <button
          v-if="!searchQuery"
          @click="showAddDialog = true"
          class="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-blue-50 text-blue-600 font-medium hover:bg-blue-100 transition-colors cursor-pointer"
        >
          <Plus class="w-4 h-4" />
          添加第一条密码
        </button>
      </div>

      <!-- 密码卡片列表 -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
        <PasswordCard
          v-for="item in passwords"
          :key="item.Id"
          :item="item"
          @deleted="handleDeleted"
          @edit="openEditDialog"
        />
      </div>
    </main>

    <!-- 新增密码弹窗 -->
    <AddPasswordDialog
      v-model:visible="showAddDialog"
      @success="handleAdded"
    />

    <!-- 修改密码弹窗 -->
    <EditPasswordDialog
      v-if="editItem"
      v-model:visible="showEditDialog"
      :item="editItem"
      @success="handleEdited"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, Lock, List } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { queryPasswords } from '@/api/pwdManage'
import type { PasswordItem } from '@/types'
import PasswordCard from '@/components/PasswordCard.vue'
import AddPasswordDialog from '@/components/AddPasswordDialog.vue'
import EditPasswordDialog from '@/components/EditPasswordDialog.vue'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const showAddDialog = ref(false)
const showEditDialog = ref(false)
const editItem = ref<PasswordItem | null>(null)
const searchQuery = ref('')
const passwords = ref<PasswordItem[]>([])

onMounted(async () => {
  if (!authStore.isLoggedIn) return
  await Promise.all([authStore.fetchUserInfo(), handleLoadAll()])
})

async function handleSearch() {
  loading.value = true
  try {
    const res = await queryPasswords(searchQuery.value)
    passwords.value = res.Data || []
  } finally {
    loading.value = false
  }
}

async function handleLoadAll() {
  searchQuery.value = ''
  loading.value = true
  try {
    const res = await queryPasswords('')
    passwords.value = res.Data || []
  } finally {
    loading.value = false
  }
}

function handleAdded() {
  searchQuery.value = ''
  handleLoadAll()
}

function handleDeleted(id: number) {
  handleLoadAll()
}

function openEditDialog(item: PasswordItem) {
  editItem.value = item
  nextTick(() => {
    showEditDialog.value = true
  })
}

function handleEdited() {
  handleLoadAll()
}

async function handleLogout() {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    authStore.logout()
    router.push('/login')
  } catch {
    // 用户取消
  }
}
</script>
