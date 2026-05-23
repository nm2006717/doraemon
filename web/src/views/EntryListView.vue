<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api/client'
import AppLayout from '../components/AppLayout.vue'

const router = useRouter()
const entries = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const perPage = ref(20)
const totalPages = ref(0)
const category = ref('')
const searchQuery = ref('')
const loading = ref(false)
const showDeleteConfirm = ref(false)
const deleteTarget = ref<any>(null)

const categoryLabels: Record<string, string> = {
  note: '笔记',
  link: '链接',
  file: '文件',
  image: '图片'
}

async function loadEntries() {
  loading.value = true
  try {
    const res = await api.listEntries({ category: category.value, page: page.value, per_page: perPage.value })
    entries.value = res.entries
    total.value = res.total
    totalPages.value = res.total_pages
  } finally {
    loading.value = false
  }
}

async function handleSearch() {
  if (!searchQuery.value.trim()) {
    loadEntries()
    return
  }
  loading.value = true
  try {
    entries.value = await api.searchEntries(searchQuery.value)
    total.value = entries.value.length
    totalPages.value = 1
    page.value = 1
  } finally {
    loading.value = false
  }
}

function confirmDelete(entry: any) {
  deleteTarget.value = entry
  showDeleteConfirm.value = true
}

async function doDelete() {
  if (!deleteTarget.value) return
  await api.deleteEntry(deleteTarget.value.ID)
  showDeleteConfirm.value = false
  deleteTarget.value = null
  loadEntries()
}

watch([category, page], () => {
  if (!searchQuery.value) loadEntries()
})

onMounted(loadEntries)
</script>

<template>
  <AppLayout>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">知识条目</h1>
      <router-link to="/entries/new" class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">新建</router-link>
    </div>

    <div class="flex gap-4 mb-4">
      <select v-model="category" class="px-3 py-2 border rounded-md">
        <option value="">全部分类</option>
        <option value="note">笔记</option>
        <option value="link">链接</option>
        <option value="file">文件</option>
        <option value="image">图片</option>
      </select>
      <div class="flex-1 flex gap-2">
        <input v-model="searchQuery" @keyup.enter="handleSearch" placeholder="搜索..."
          class="flex-1 px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500" />
        <button @click="handleSearch" class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700">搜索</button>
      </div>
    </div>

    <div class="bg-white rounded-lg shadow">
      <div v-if="loading" class="p-4 text-gray-500">加载中...</div>
      <div v-else-if="entries.length === 0" class="p-4 text-gray-500">暂无条目</div>
      <table v-else class="w-full">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-4 py-2 text-left text-sm font-medium text-gray-500">ID</th>
            <th class="px-4 py-2 text-left text-sm font-medium text-gray-500">标题</th>
            <th class="px-4 py-2 text-left text-sm font-medium text-gray-500">分类</th>
            <th class="px-4 py-2 text-left text-sm font-medium text-gray-500">标签</th>
            <th class="px-4 py-2 text-left text-sm font-medium text-gray-500">创建时间</th>
            <th class="px-4 py-2 text-left text-sm font-medium text-gray-500">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="entry in entries" :key="entry.ID" class="border-t hover:bg-gray-50">
            <td class="px-4 py-2 text-sm">{{ entry.ID }}</td>
            <td class="px-4 py-2">
              <router-link :to="`/entries/${entry.ID}`" class="text-blue-600 hover:underline">{{ entry.Title }}</router-link>
            </td>
            <td class="px-4 py-2 text-sm">{{ categoryLabels[entry.Category] || entry.Category }}</td>
            <td class="px-4 py-2 text-sm text-gray-500">{{ entry.Tags }}</td>
            <td class="px-4 py-2 text-sm text-gray-500">{{ new Date(entry.CreatedAt).toLocaleDateString() }}</td>
            <td class="px-4 py-2">
              <button @click="confirmDelete(entry)" class="text-red-600 hover:underline text-sm">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="totalPages > 1" class="flex justify-center gap-2 mt-4">
      <button @click="page--" :disabled="page <= 1"
        class="px-3 py-1 border rounded disabled:opacity-50">上一页</button>
      <span class="px-3 py-1">{{ page }} / {{ totalPages }}</span>
      <button @click="page++" :disabled="page >= totalPages"
        class="px-3 py-1 border rounded disabled:opacity-50">下一页</button>
    </div>

    <div v-if="showDeleteConfirm" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white p-6 rounded-lg shadow-lg max-w-sm w-full">
        <h3 class="text-lg font-medium mb-2">确认删除</h3>
        <p class="text-gray-600 mb-4">确定要删除「{{ deleteTarget?.Title }}」吗？此操作不可撤销。</p>
        <div class="flex justify-end gap-2">
          <button @click="showDeleteConfirm = false" class="px-4 py-2 border rounded hover:bg-gray-50">取消</button>
          <button @click="doDelete" class="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700">删除</button>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
