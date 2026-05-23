<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api/client'
import AppLayout from '../components/AppLayout.vue'

const router = useRouter()
const title = ref('')
const content = ref('')
const category = ref('note')
const tags = ref('')
const error = ref('')
const loading = ref(false)

async function handleSubmit() {
  error.value = ''
  if (!title.value.trim()) {
    error.value = '请填写标题'
    return
  }

  loading.value = true
  try {
    const entry = await api.createEntry({
      title: title.value,
      content: content.value,
      category: category.value,
      tags: tags.value
    })
    router.push(`/entries/${entry.ID}`)
  } catch (e: any) {
    error.value = e.message || '创建失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <AppLayout>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">新建条目</h1>
      <router-link to="/entries" class="px-4 py-2 border rounded hover:bg-gray-50">返回列表</router-link>
    </div>

    <div class="bg-white rounded-lg shadow p-6">
      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">标题</label>
          <input v-model="title" type="text"
            class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">内容</label>
          <textarea v-model="content" rows="6"
            class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"></textarea>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">分类</label>
          <select v-model="category"
            class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
            <option value="note">笔记</option>
            <option value="link">链接</option>
            <option value="file">文件</option>
            <option value="image">图片</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">标签</label>
          <input v-model="tags" type="text" placeholder="多个标签用逗号分隔"
            class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500" />
        </div>

        <p v-if="error" class="text-red-500 text-sm">{{ error }}</p>

        <div class="flex justify-end">
          <button type="submit" :disabled="loading"
            class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50">
            {{ loading ? '创建中...' : '创建' }}
          </button>
        </div>
      </form>
    </div>
  </AppLayout>
</template>
