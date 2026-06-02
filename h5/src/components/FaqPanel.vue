<template>
  <div v-if="visible" class="faq-panel">
    <div class="faq-header">
      <span>常见问题</span>
      <button class="faq-close" @click="$emit('close')">×</button>
    </div>
    <div class="faq-search">
      <input
        v-model="keyword"
        type="text"
        placeholder="搜索问题..."
        @compositionstart="composing = true"
        @compositionend="onCompositionEnd"
        @input="onInput"
      />
    </div>
    <div class="faq-list">
      <div v-if="loading" class="faq-loading">搜索中...</div>
      <div v-else-if="list.length === 0" class="faq-empty">
        {{ keyword ? '未找到相关问题' : '输入关键词搜索' }}
      </div>
      <div v-for="item in list" :key="item.id" class="faq-item" @click="$emit('select', item)">
        <div class="faq-question">{{ item.question }}</div>
        <div class="faq-answer-preview">{{ truncate(item.answer, 60) }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { queryFaq } from '@/api/chat'

defineProps({
  visible: { type: Boolean, default: false }
})

const emit = defineEmits(['close', 'select'])

const keyword = ref('')
const list = ref([])
const loading = ref(false)
const composing = ref(false)
let debounceTimer = null

function doSearch() {
  if (!keyword.value) {
    list.value = []
    return
  }
  loading.value = true
  queryFaq({ keyword: keyword.value, page: 1, pageSize: 10 })
    .then(res => {
      list.value = res.data?.list || []
    })
    .finally(() => { loading.value = false })
}

function onInput() {
  if (composing.value) return
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(doSearch, 300)
}

function onCompositionEnd(e) {
  composing.value = false
  keyword.value = e.target.value
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(doSearch, 300)
}

function truncate(text, len) {
  if (!text) return ''
  return text.length > len ? text.substring(0, len) + '...' : text
}
</script>
