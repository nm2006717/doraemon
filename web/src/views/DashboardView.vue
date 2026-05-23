<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../api/client'
import AppLayout from '../components/AppLayout.vue'

const stats = ref<{ total_entries: number; by_category: Record<string, number>; recent_entries: any[] } | null>(null)
const loading = ref(true)

onMounted(async () => {
  try {
    stats.value = await api.getDashboard()
  } finally {
    loading.value = false
  }
})

const categoryLabels: Record<string, string> = {
  note: '笔记',
  link: '链接',
  file: '文件',
  image: '图片'
}
</script>

<template>
  <AppLayout>
    <h1 class="text-2xl font-bold mb-6">仪表盘</h1>

    <div v-if="loading" class="text-gray-500">加载中...</div>

    <template v-else-if="stats">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
        <div class="bg-white p-4 rounded-lg shadow">
          <div class="text-sm text-gray-500">总条目</div>
          <div class="text-3xl font-bold">{{ stats.total_entries }}</div>
        </div>
        <div v-for="(count, cat) in stats.by_category" :key="cat" class="bg-white p-4 rounded-lg shadow">
          <div class="text-sm text-gray-500">{{ categoryLabels[cat] || cat }}</div>
          <div class="text-3xl font-bold">{{ count }}</div>
        </div>
      </div>

      <div class="bg-white rounded-lg shadow">
        <div class="px-4 py-3 border-b font-medium">最近条目</div>
        <div v-if="!stats.recent_entries || stats.recent_entries.length === 0" class="p-4 text-gray-500">暂无条目</div>
        <table v-else class="w-full">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-2 text-left text-sm font-medium text-gray-500">标题</th>
              <th class="px-4 py-2 text-left text-sm font-medium text-gray-500">分类</th>
              <th class="px-4 py-2 text-left text-sm font-medium text-gray-500">标签</th>
              <th class="px-4 py-2 text-left text-sm font-medium text-gray-500">创建时间</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="entry in stats.recent_entries" :key="entry.ID" class="border-t hover:bg-gray-50">
              <td class="px-4 py-2">
                <router-link :to="`/entries/${entry.ID}`" class="text-blue-600 hover:underline">{{ entry.Title }}</router-link>
              </td>
              <td class="px-4 py-2 text-sm">{{ categoryLabels[entry.Category] || entry.Category }}</td>
              <td class="px-4 py-2 text-sm text-gray-500">{{ entry.Tags }}</td>
              <td class="px-4 py-2 text-sm text-gray-500">{{ new Date(entry.CreatedAt).toLocaleDateString() }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>
  </AppLayout>
</template>
