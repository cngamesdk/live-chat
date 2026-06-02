<template>
  <div class="message-bubble" :class="bubbleClass">
    <template v-if="!isSystem">
      <div class="bubble-avatar">{{ avatarText }}</div>
      <div class="bubble-content">
        <div class="bubble-sender">{{ message.senderName }}</div>
        <div class="bubble-body" :class="bodyClass">
          <!-- Text message -->
          <div v-if="message.msgType === 'text'" class="bubble-text">
            {{ message.content }}
          </div>

          <!-- Image message -->
          <div v-else-if="message.msgType === 'image'" class="bubble-image">
            <img :src="message.attachmentURL || message.content" @click="previewImage" alt="image" />
          </div>

          <!-- Video message -->
          <div v-else-if="message.msgType === 'video'" class="bubble-video">
            <video :src="message.attachmentURL || message.content" controls preload="metadata" />
          </div>

          <!-- FAQ reply tag -->
          <div v-if="message.isFaqReply" class="faq-tag">自动回复</div>

          <!-- Unread indicator -->
          <div v-if="!isUser && !message.isRead" class="unread-dot"></div>
        </div>
        <div class="bubble-time">{{ timeText }}</div>
      </div>
    </template>

    <!-- System message centered -->
    <template v-else>
      <div class="system-message">
        <div class="system-content">{{ message.content }}</div>
        <div class="system-time">{{ timeText }}</div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { formatTime } from '@/utils/index'

const props = defineProps({
  message: { type: Object, required: true }
})

const emit = defineEmits(['preview'])

const isUser = computed(() => props.message.senderType === 'user')
const isSystem = computed(() => props.message.senderType === 'system')

const bubbleClass = computed(() => ({
  'bubble-user': isUser.value,
  'bubble-agent': !isUser.value && !isSystem.value,
  'bubble-system': isSystem.value
}))

const bodyClass = computed(() => ({
  'bg-user': isUser.value,
  'bg-agent': !isUser.value && !isSystem.value,
  'bg-system': isSystem.value
}))

const avatarText = computed(() => {
  if (isUser.value) return '我'
  if (isSystem.value) return '系'
  return (props.message.senderName || '客').charAt(0)
})

const timeText = computed(() => {
  const t = props.message.created_at || props.message.createdAt
  return t ? formatTime(t) : ''
})

function previewImage() {
  emit('preview', props.message.attachmentURL || props.message.content)
}
</script>

<style scoped>
.message-bubble {
  display: flex;
  margin-bottom: 16px;
  gap: 8px;
}

.bubble-user {
  flex-direction: row-reverse;
}

.bubble-agent {
  flex-direction: row;
}

.bubble-system {
  justify-content: center;
}

.bubble-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #e0e0e0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  color: #666;
  flex-shrink: 0;
}

.bubble-user .bubble-avatar {
  background: #4CAF50;
  color: white;
}

.bubble-content {
  max-width: 70%;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.bubble-user .bubble-content {
  align-items: flex-end;
}

.bubble-sender {
  font-size: 12px;
  color: #999;
  padding: 0 8px;
}

.bubble-body {
  padding: 10px 14px;
  border-radius: 8px;
  word-break: break-word;
}

.bg-user {
  background: #4CAF50;
  color: white;
}

.bg-agent {
  background: #f5f5f5;
  color: #333;
}

.bg-system {
  background: #fff3cd;
  color: #856404;
}

.bubble-text {
  line-height: 1.5;
}

.bubble-image img {
  max-width: 200px;
  max-height: 200px;
  border-radius: 4px;
  cursor: pointer;
}

.bubble-video video {
  max-width: 250px;
  border-radius: 4px;
}

.faq-tag {
  margin-top: 4px;
  font-size: 11px;
  color: #ff9800;
  background: rgba(255, 152, 0, 0.1);
  padding: 2px 6px;
  border-radius: 3px;
  display: inline-block;
}

.bubble-time {
  font-size: 11px;
  color: #999;
  padding: 0 8px;
}

/* System message centered style */
.system-message {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.system-content {
  background: #f0f0f0;
  color: #666;
  padding: 6px 12px;
  border-radius: 12px;
  font-size: 12px;
  text-align: center;
  max-width: 80%;
}

.system-time {
  font-size: 11px;
  color: #999;
}

.unread-dot {
  position: absolute;
  top: -4px;
  right: -4px;
  width: 8px;
  height: 8px;
  background: #ff4444;
  border-radius: 50%;
  border: 2px solid white;
}

.bubble-body {
  position: relative;
}
</style>
