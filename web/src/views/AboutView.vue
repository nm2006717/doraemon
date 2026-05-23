<script setup lang="ts">
import { ref, onMounted } from 'vue'

const lang = ref<'en' | 'zh'>('zh')
const content = ref<{ en: string; zh: string }>({ en: '', zh: '' })
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await fetch('/api/readme')
    content.value = await res.json()
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="max-w-4xl mx-auto">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">关于</h1>
      <div class="flex rounded-lg overflow-hidden border border-gray-300">
        <button
          @click="lang = 'zh'"
          :class="['px-3 py-1 text-sm', lang === 'zh' ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-100']"
        >中文</button>
        <button
          @click="lang = 'en'"
          :class="['px-3 py-1 text-sm', lang === 'en' ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-100']"
        >English</button>
      </div>
    </div>

    <div v-if="loading" class="text-gray-500">加载中...</div>
    <div
      v-else
      class="prose prose-sm max-w-none"
      v-html="lang === 'zh' ? content.zh : content.en"
    />
  </div>
</template>
