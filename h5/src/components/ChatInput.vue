<template>
  <div class="chat-input-area">
    <!-- Upload preview -->
    <div v-if="previewUrl" class="upload-preview">
      <img v-if="previewType === 'image'" :src="previewUrl" alt="preview" />
      <video v-else-if="previewType === 'video'" :src="previewUrl" controls />
      <button class="preview-remove" @click="clearPreview">×</button>
    </div>

    <!-- Toolbar -->
    <div class="input-toolbar">
      <label class="tool-btn" title="图片">
        <input type="file" accept="image/*" hidden @change="onFileChange($event, 'image')" />
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
          <circle cx="8.5" cy="8.5" r="1.5"/>
          <polyline points="21,15 16,10 5,21"/>
        </svg>
      </label>
      <label class="tool-btn" title="视频">
        <input type="file" accept="video/*" hidden @change="onFileChange($event, 'video')" />
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polygon points="23,7 16,12 23,17 23,7"/>
          <rect x="1" y="5" width="15" height="14" rx="2" ry="2"/>
        </svg>
      </label>
      <button class="tool-btn" @click="$emit('toggleFaq')" title="常见问题">FAQ</button>
    </div>

    <!-- Input row -->
    <div class="input-row">
      <textarea
        ref="inputRef"
        v-model="text"
        class="input-textarea"
        :placeholder="disabled ? '会话已结束' : placeholder"
        :disabled="disabled"
        rows="1"
        @keydown.enter.exact.prevent="sendText"
        @input="autoResize"
      />
      <button class="send-btn" :disabled="!canSend" @click="sendText">发送</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  placeholder: { type: String, default: '输入您的问题...' },
  disabled: { type: Boolean, default: false }
})

const emit = defineEmits(['send-message', 'send-image', 'send-video', 'toggleFaq'])

const text = ref('')
const previewUrl = ref('')
const previewType = ref('')
const fileToSend = ref(null)
const inputRef = ref(null)

const canSend = computed(() => !props.disabled && (text.value.trim() || fileToSend.value))

function sendText() {
  if (props.disabled) return

  const content = text.value.trim()
  if (!content && !fileToSend.value) return

  if (fileToSend.value) {
    if (previewType.value === 'image') {
      emit('send-image', fileToSend.value)
    } else {
      emit('send-video', fileToSend.value)
    }
    clearPreview()
  } else {
    emit('send-message', content)
  }
  text.value = ''
}

function onFileChange(event, type) {
  const file = event.target.files[0]
  if (!file) return
  fileToSend.value = file
  previewType.value = type
  previewUrl.value = URL.createObjectURL(file)
}

function clearPreview() {
  previewUrl.value = ''
  previewType.value = ''
  fileToSend.value = null
}

function autoResize() {
  const el = inputRef.value
  if (el) {
    el.style.height = 'auto'
    el.style.height = Math.min(el.scrollHeight, 120) + 'px'
  }
}
</script>
