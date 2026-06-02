import { ref, onUnmounted } from 'vue'

export function useWebSocket() {
  const connected = ref(false)
  const message = ref(null)
  let ws = null
  let reconnectTimer = null
  let pingTimer = null

  function connect(token) {
    const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = import.meta.env.VITE_WS_BASE || `${protocol}//${location.host}`
    const url = `${host}/ws/chat?token=${encodeURIComponent(token)}`

    ws = new WebSocket(url)

    ws.onopen = () => {
      connected.value = true
      // Heartbeat
      pingTimer = setInterval(() => {
        if (ws && ws.readyState === WebSocket.OPEN) {
          ws.send(JSON.stringify({ type: 'ping' }))
        }
      }, 30000)
    }

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        message.value = data
      } catch (e) {
        console.error('WS parse error:', e)
      }
    }

    ws.onclose = () => {
      connected.value = false
      clearInterval(pingTimer)
      reconnectTimer = setTimeout(() => connect(token), 3000)
    }

    ws.onerror = () => {
      ws?.close()
    }
  }

  function send(msg) {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(msg))
    }
  }

  function disconnect() {
    clearTimeout(reconnectTimer)
    clearInterval(pingTimer)
    if (ws) {
      ws.close()
      ws = null
    }
    connected.value = false
  }

  onUnmounted(disconnect)

  return { connected, message, connect, send, disconnect }
}
