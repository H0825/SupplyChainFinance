import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    proxy: {
      '/api': 'http://localhost:8888',
      '/workflow': 'http://localhost:8888',
      '/health': 'http://localhost:8888'
    }
  },
  build: {
    outDir: 'dist',
    emptyOutDir: true
  }
})
