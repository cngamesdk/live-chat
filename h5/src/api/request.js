import axios from 'axios'

const service = axios.create({
  baseURL: '/api/v1',
  timeout: 30000
})

// Request interceptor
service.interceptors.request.use(
  config => {
    const token = localStorage.getItem('chat_session_token')
    if (token) {
      config.headers['X-Session-Token'] = token
    }
    return config
  },
  error => Promise.reject(error)
)

// Response interceptor
service.interceptors.response.use(
  response => {
    const res = response.data
    if (res.code !== 0) {
      console.error('API Error:', res.msg)
      return Promise.reject(new Error(res.msg || 'Error'))
    }
    return res
  },
  error => {
    console.error('Request Error:', error)
    return Promise.reject(error)
  }
)

export default service
