<template>
  <div class="chat-body" ref="bodyRef">
    <!-- Welcome message when empty -->
    <div v-if="messages.length === 0 && welcomeMessage" class="welcome-area">
      <div class="welcome-title">{{ welcomeTitle || '欢迎咨询' }}</div>
      <div class="welcome-text">{{ welcomeMessage }}</div>
    </div>

    <!-- Message list -->
    <div v-for="msg in messages" :key="msg.id" class="message-wrapper">
      <MessageBubble :message="msg" @preview="url => $emit('preview', url)" />
    </div>

    <!-- Queue position indicator -->
    <div v-if="queuePosition > 0" class="queue-indicator">
      排队中，前方还有 {{ queuePosition }} 人...
    </div>

    <div ref="scrollAnchor" />
  </div>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue'
import MessageBubble from './MessageBubble.vue'

const props = defineProps({
  messages: { type: Array, default: () => [] },
  welcomeTitle: { type: String, default: '' },
  welcomeMessage: { type: String, default: '' },
  queuePosition: { type: Number, default: -1 }
})

defineEmits(['preview'])

const bodyRef = ref(null)
const scrollAnchor = ref(null)

function scrollToBottom() {
  nextTick(() => {
    scrollAnchor.value?.scrollIntoView({ behavior: 'smooth' })
  })
}

watch(() => props.messages.length, scrollToBottom)
watch(() => props.queuePosition, scrollToBottom)
</script>
