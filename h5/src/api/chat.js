import request from './request'

export function getProductInfo(code) {
  return request({ url: `/product/${code}/info`, method: 'get' })
}

export function initSession(data) {
  return request({ url: '/chat/init', method: 'post', data })
}

export function getChatHistory(params) {
  return request({ url: '/chat/history', method: 'get', params })
}

export function queryFaq(params) {
  return request({ url: '/chat/faq', method: 'get', params })
}

export function uploadFile(formData) {
  return request({
    url: '/upload',
    method: 'post',
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

export function markMessagesAsRead(data) {
  return request({ url: '/chat/mark-read', method: 'post', data })
}

export function getUnreadCount(params) {
  return request({ url: '/chat/unread-count', method: 'get', params })
}

export function closeSession() {
  return request({ url: '/chat/close', method: 'post' })
}
