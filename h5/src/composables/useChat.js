import { ref, reactive, watch, nextTick } from 'vue'
import { getProductInfo, initSession, getChatHistory, queryFaq, uploadFile, markMessagesAsRead, getUnreadCount, closeSession as closeSessionApi } from '@/api/chat'
import { useWebSocket } from './useWebSocket'

// Normalize message keys from snake_case (API) to camelCase (UI)
function normalizeMessage(m) {
  return {
    id: m.id || Date.now(),
    content: m.content || '',
    msgType: m.msg_type || m.msgType || 'text',
    senderType: m.sender_type || m.senderType || 'system',
    senderName: m.sender_name || m.senderName || '系统',
    attachmentURL: m.attachment_url || m.attachmentURL || '',
    isFaqReply: m.is_faq_reply ?? m.isFaqReply ?? false,
    isRead: m.is_read ?? m.isRead ?? false,
    messageId: m.message_id || m.messageId || m.id,
    createdAt: m.created_at || m.createdAt || null
  }
}

export function useChat() {
  const product = reactive({
    code: '',
    name: '',
    logo: '',
    welcomeTitle: '',
    welcomeMessage: ''
  })

  const session = reactive({
    id: 0,
    token: '',
    userId: ''
  })

  const messages = ref([])
  const loading = ref(false)
  const faqList = ref([])
  const showFaq = ref(false)
  const queuePosition = ref(-1)
  const unreadCount = ref(0)
  const sessionClosed = ref(false)

  const ws = useWebSocket()

  watch(() => ws.message.value, (msg) => {
    if (!msg) return
    handleWsMessage(msg)
  })

  function handleWsMessage(msg) {
    switch (msg.type) {
      case 'message':
        const newMsg = normalizeMessage({
          id: msg.message_id || Date.now(),
          content: msg.content,
          msg_type: msg.msg_type || 'text',
          sender_type: msg.sender_type,
          sender_name: msg.sender_name,
          attachment_url: msg.attachment_url,
          is_faq_reply: msg.is_faq_reply,
          is_read: msg.is_read || false,
          message_id: msg.message_id,
          created_at: new Date().toISOString()
        })
        messages.value.push(newMsg)
        // Update unread count if it's from agent
        if (msg.sender_type === 'agent' && !msg.is_read) {
          unreadCount.value++
        }
        // Auto mark as read after a short delay
        if (msg.sender_type === 'agent' && msg.message_id) {
          setTimeout(() => {
            markAsRead([msg.message_id])
          }, 1000)
        }
        break
      case 'queue':
        queuePosition.value = msg.queue_position
        break
      case 'connected':
        messages.value.push(normalizeMessage({
          id: Date.now(),
          content: `已为您接通人工客服 ${msg.agent_name || ''}`,
          msg_type: 'system',
          sender_type: 'system',
          sender_name: '系统',
          created_at: new Date().toISOString()
        }))
        queuePosition.value = -1
        break
      case 'closed':
        sessionClosed.value = true
        messages.value.push(normalizeMessage({
          id: Date.now(),
          content: '会话已结束，感谢您的咨询',
          msg_type: 'system',
          sender_type: 'system',
          sender_name: '系统',
          created_at: new Date().toISOString()
        }))
        break
    }
  }

  async function init(productCode, userId, userName, userToken) {
    loading.value = true
    try {
      const info = await getProductInfo(productCode)
      if (info.data) {
        Object.assign(product, {
          code: info.data.product_code,
          name: info.data.name,
          logo: info.data.logo,
          welcomeTitle: info.data.welcome_title,
          welcomeMessage: info.data.welcome_message
        })
      }

      const res = await initSession({
        product_code: productCode,
        user_id: userId,
        user_name: userName || userId,
        user_token: userToken || '',
        source: 'h5'
      })

      session.id = res.data.session_id
      session.token = res.data.token
      session.userId = res.data.user_id
      localStorage.setItem('chat_session_token', res.data.token)

      ws.connect(res.data.token)

      const history = await getChatHistory({ page: 1, pageSize: 50 })
      if (history.data && history.data.list) {
        messages.value = history.data.list.map(normalizeMessage)
        // Auto mark unread agent messages as read
        const unreadAgentMessages = messages.value
          .filter(m => m.senderType === 'agent' && !m.isRead && m.messageId)
          .map(m => m.messageId)
        if (unreadAgentMessages.length > 0) {
          setTimeout(() => {
            markAsRead(unreadAgentMessages)
          }, 1500)
        }
      }
    } catch (e) {
      console.error('Init chat failed:', e)
    } finally {
      loading.value = false
    }
  }

  function sendMessage(content, msgType = 'text', attachmentURL = '') {
    if (sessionClosed.value) {
      console.warn('会话已关闭，无法发送消息')
      return
    }
    // Only send via WebSocket; the echo back handles the UI update
    ws.send({ type: 'message', content, msg_type: msgType, attachment_url: attachmentURL })
  }

  async function sendImage(file) {
    const formData = new FormData()
    formData.append('file', file)
    try {
      const res = await uploadFile(formData)
      sendMessage('[图片]', 'image', res.data.url)
    } catch (e) {
      console.error('Upload failed:', e)
    }
  }

  async function sendVideo(file) {
    const formData = new FormData()
    formData.append('file', file)
    try {
      const res = await uploadFile(formData)
      sendMessage('[视频]', 'video', res.data.url)
    } catch (e) {
      console.error('Upload failed:', e)
    }
  }

  async function searchFaq(keyword) {
    if (!keyword) {
      faqList.value = []
      return
    }
    try {
      const res = await queryFaq({ keyword, page: 1, pageSize: 5 })
      faqList.value = res.data?.list || []
    } catch (e) {
      console.error('FAQ search failed:', e)
    }
  }

  function selectFaq(faq) {
    showFaq.value = false
    messages.value.push(normalizeMessage({
      id: Date.now(),
      content: faq.question,
      msg_type: 'text',
      sender_type: 'user',
      sender_name: session.userId,
      created_at: new Date().toISOString()
    }))
    messages.value.push(normalizeMessage({
      id: Date.now() + 1,
      content: faq.answer,
      msg_type: 'text',
      sender_type: 'system',
      sender_name: '智能客服',
      is_faq_reply: true,
      created_at: new Date().toISOString()
    }))
  }

  function destroy() {
    ws.disconnect()
    localStorage.removeItem('chat_session_token')
  }

  async function markAsRead(messageIds) {
    try {
      await markMessagesAsRead({ message_ids: messageIds })
      // Update local state
      messageIds.forEach(id => {
        const msg = messages.value.find(m => m.id === id || m.messageId === id)
        if (msg) {
          msg.isRead = true
        }
      })
      // Refresh unread count
      await refreshUnreadCount()
    } catch (e) {
      console.error('Mark as read failed:', e)
    }
  }

  async function refreshUnreadCount() {
    try {
      const res = await getUnreadCount({ sender_type: 'agent' })
      unreadCount.value = res.data?.count || 0
    } catch (e) {
      console.error('Get unread count failed:', e)
    }
  }

  async function closeSession() {
    try {
      await closeSessionApi()
      sessionClosed.value = true
      messages.value.push(normalizeMessage({
        id: Date.now(),
        content: '您已结束会话，感谢您的咨询',
        msg_type: 'system',
        sender_type: 'system',
        sender_name: '系统',
        created_at: new Date().toISOString()
      }))
    } catch (e) {
      console.error('Close session failed:', e)
    }
  }

  return {
    product,
    session,
    messages,
    loading,
    faqList,
    showFaq,
    queuePosition,
    unreadCount,
    sessionClosed,
    connected: ws.connected,
    init,
    sendMessage,
    sendImage,
    sendVideo,
    searchFaq,
    selectFaq,
    markAsRead,
    refreshUnreadCount,
    closeSession,
    showFaqPanel: () => { showFaq.value = !showFaq.value },
    destroy
  }
}
