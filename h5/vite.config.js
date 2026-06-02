import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd())
  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': resolve(__dirname, 'src')
      }
    },
    server: {
      port: 8081,
      host: '0.0.0.0',
      proxy: {
        '/api': {
          target: env.VITE_API_BASE || 'http://127.0.0.1:8890',
          changeOrigin: true
        },
        '/ws': {
          target: env.VITE_WS_BASE || 'ws://127.0.0.1:8890',
          ws: true
        },
        '/uploads': {
          target: env.VITE_API_BASE || 'http://127.0.0.1:8890',
          changeOrigin: true
        }
      }
    },
    build: {
      outDir: 'dist',
      assetsDir: 'assets'
    }
  }
})
