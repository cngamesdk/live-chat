<template>
  <div class="chat-header">
    <div class="header-left">
      <img v-if="logo" :src="logo" class="header-logo" alt="logo" />
      <span class="header-name">{{ name || '在线客服' }}</span>
    </div>
    <div class="header-right">
      <span class="status-dot" :class="connected ? 'online' : 'offline'" />
      <span class="status-text">{{ connected ? '已连接' : '连接中...' }}</span>
      <button v-if="!sessionClosed && showCloseButton" class="close-btn" @click="$emit('close')" title="结束会话">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup>
defineProps({
  logo: { type: String, default: '' },
  name: { type: String, default: '在线客服' },
  connected: { type: Boolean, default: false },
  sessionClosed: { type: Boolean, default: false },
  showCloseButton: { type: Boolean, default: true }
})

defineEmits(['close'])
</script>

<style scoped>
.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #fff;
  border-bottom: 1px solid #e0e0e0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-logo {
  width: 32px;
  height: 32px;
  border-radius: 50%;
}

.header-name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-dot.online {
  background: #4CAF50;
}

.status-dot.offline {
  background: #999;
}

.status-text {
  font-size: 12px;
  color: #666;
}

.close-btn {
  padding: 4px;
  background: transparent;
  border: none;
  cursor: pointer;
  color: #999;
  transition: color 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  color: #f44336;
}
</style>
