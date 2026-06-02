<template>
  <div class="chat-container">
    <ChatHeader
      :logo="product.logo"
      :name="product.name"
      :connected="connected"
      :sessionClosed="sessionClosed"
      @close="handleCloseSession"
    />

    <ChatBody
      :messages="messages"
      :welcomeTitle="product.welcomeTitle"
      :welcomeMessage="product.welcomeMessage"
      :queuePosition="queuePosition"
      @preview="previewUrl = $event"
    />

    <FaqPanel
      :visible="showFaq"
      @close="showFaq = false"
      @select="selectFaq"
    />

    <ChatInput
      :disabled="sessionClosed"
      @send-message="sendMessage"
      @send-image="sendImage"
      @send-video="sendVideo"
      @toggle-faq="showFaqPanel"
    />

    <!-- Image/Video preview modal -->
    <Teleport to="body">
      <div v-if="previewUrl" class="preview-modal" @click="previewUrl = ''">
        <img v-if="isImagePreview" :src="previewUrl" class="preview-content" />
        <video v-else-if="isVideoPreview" :src="previewUrl" class="preview-content" controls autoplay />
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useChat } from '@/composables/useChat'
import ChatHeader from '@/components/ChatHeader.vue'
import ChatBody from '@/components/ChatBody.vue'
import ChatInput from '@/components/ChatInput.vue'
import FaqPanel from '@/components/FaqPanel.vue'
import { isImage, isVideo } from '@/utils/index'

const previewUrl = ref('')

const isImagePreview = computed(() => {
  return previewUrl.value && (previewUrl.value.startsWith('data:') ||
    /\.(jpg|jpeg|png|gif|webp)/i.test(previewUrl.value))
})
const isVideoPreview = computed(() => {
  return previewUrl.value && /\.(mp4|mov|webm)/i.test(previewUrl.value)
})

const {
  product,
  messages,
  loading,
  faqList,
  showFaq,
  queuePosition,
  sessionClosed,
  connected,
  init,
  sendMessage,
  sendImage,
  sendVideo,
  selectFaq,
  showFaqPanel,
  closeSession
} = useChat()

function handleCloseSession() {
  if (confirm('确定要结束会话吗？')) {
    closeSession()
  }
}

onMounted(() => {
  // Read config from URL params or script tag data attributes
  const params = new URLSearchParams(window.location.search)
  const productCode = params.get('product_code') || 'DEMO'
  const userId = params.get('user_id') || 'guest_' + Date.now()
  const userName = params.get('user_name') || userId
  const userToken = params.get('user_token') || ''

  init(productCode, userId, userName, userToken)
})
</script>
