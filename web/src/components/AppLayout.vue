<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()

function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<template>
  <div class="min-h-screen flex">
    <aside class="w-56 bg-gray-800 text-white flex flex-col">
      <div class="p-4 text-xl font-bold border-b border-gray-700">Doraemon</div>
      <nav class="flex-1 p-4 space-y-2">
        <router-link to="/" class="block px-3 py-2 rounded hover:bg-gray-700">仪表盘</router-link>
        <router-link to="/entries" class="block px-3 py-2 rounded hover:bg-gray-700">知识条目</router-link>
        <router-link to="/entries/new" class="block px-3 py-2 rounded hover:bg-gray-700">新建条目</router-link>
      </nav>
      <div class="p-4 border-t border-gray-700 space-y-2">
        <router-link to="/about" class="block px-3 py-2 text-sm text-gray-300 rounded hover:bg-gray-700">关于</router-link>
        <button @click="logout" class="w-full px-3 py-2 text-sm bg-gray-700 rounded hover:bg-gray-600">退出登录</button>
      </div>
    </aside>
    <main class="flex-1 p-6 overflow-auto">
      <slot />
    </main>
  </div>
</template>
