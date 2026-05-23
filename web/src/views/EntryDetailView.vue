<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { api } from '../api/client'
import AppLayout from '../components/AppLayout.vue'

const router = useRouter()
const route = useRoute()
const entry = ref<any>(null)
const loading = ref(true)
const error = ref('')

const editing = ref(false)
const editTitle = ref('')
const editContent = ref('')
const editCategory = ref('')
const editTags = ref('')
const saving = ref(false)
const saveError = ref('')

const showDeleteConfirm = ref(false)

const categoryLabels: Record<string, string> = {
  note: '笔记',
  link: '链接',
  file: '文件',
  image: '图片'
}

async function loadEntry() {
  loading.value = true
  try {
    entry.value = await api.getEntry(Number(route.params.id))
  } catch (e: any) {
    error.value = e.message || '加载失败'
  } finally {
    loading.value = false
  }
}

function startEdit() {
  editTitle.value = entry.value.Title
  editContent.value = entry.value.Content
  editCategory.value = entry.value.Category
  editTags.value = entry.value.Tags || ''
  editing.value = true
  saveError.value = ''
}

function cancelEdit() {
  editing.value = false
  saveError.value = ''
}

async function handleSave() {
  saveError.value = ''
  if (!editTitle.value.trim()) {
    saveError.value = '请填写标题'
    return
  }
  saving.value = true
  try {
    entry.value = await api.updateEntry(Number(route.params.id), {
      title: editTitle.value,
      content: editContent.value,
      category: editCategory.value,
      tags: editTags.value
    })
    editing.value = false
  } catch (e: any) {
    saveError.value = e.message || '保存失败'
  } finally {
    saving.value = false
  }
}

async function doDelete() {
  await api.deleteEntry(Number(route.params.id))
  router.push('/entries')
}

onMounted(loadEntry)
</script>

<template>
  <AppLayout>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">条目详情</h1>
      <div class="flex gap-2">
        <router-link to="/entries" class="px-4 py-2 border rounded hover:bg-gray-50">返回列表</router-link>
        <button v-if="entry && !editing" @click="startEdit"
          class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">编辑</button>
        <button v-if="entry && !editing" @click="showDeleteConfirm = true"
          class="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700">删除</button>
      </div>
    </div>

    <div v-if="loading" class="text-gray-500">加载中...</div>
    <div v-else-if="error" class="text-red-500">{{ error }}</div>

    <div v-else-if="!editing" class="bg-white rounded-lg shadow p-6 space-y-4">
      <div>
        <h2 class="text-xl font-medium">{{ entry.Title }}</h2>
      </div>
      <div class="flex items-center gap-2">
        <span class="inline-block px-2 py-1 bg-blue-100 text-blue-800 rounded text-sm">
          {{ categoryLabels[entry.Category] || entry.Category }}
        </span>
        <span v-if="entry.Tags" class="text-sm text-gray-500">{{ entry.Tags }}</span>
      </div>
      <div class="whitespace-pre-wrap text-gray-800">{{ entry.Content }}</div>
      <div class="text-sm text-gray-500">
        创建时间：{{ new Date(entry.CreatedAt).toLocaleString() }}
      </div>
    </div>

    <div v-else class="bg-white rounded-lg shadow p-6">
      <form @submit.prevent="handleSave" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">标题</label>
          <input v-model="editTitle" type="text"
            class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">内容</label>
          <textarea v-model="editContent" rows="6"
            class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"></textarea>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">分类</label>
          <select v-model="editCategory"
            class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
            <option value="note">笔记</option>
            <option value="link">链接</option>
            <option value="file">文件</option>
            <option value="image">图片</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">标签</label>
          <input v-model="editTags" type="text" placeholder="多个标签用逗号分隔"
            class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500" />
        </div>

        <p v-if="saveError" class="text-red-500 text-sm">{{ saveError }}</p>

        <div class="flex justify-end gap-2">
          <button type="button" @click="cancelEdit" class="px-4 py-2 border rounded hover:bg-gray-50">取消</button>
          <button type="submit" :disabled="saving"
            class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50">
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </div>
      </form>
    </div>

    <div v-if="showDeleteConfirm" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white p-6 rounded-lg shadow-lg max-w-sm w-full">
        <h3 class="text-lg font-medium mb-2">确认删除</h3>
        <p class="text-gray-600 mb-4">确定要删除「{{ entry?.Title }}」吗？此操作不可撤销。</p>
        <div class="flex justify-end gap-2">
          <button @click="showDeleteConfirm = false" class="px-4 py-2 border rounded hover:bg-gray-50">取消</button>
          <button @click="doDelete" class="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700">删除</button>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
