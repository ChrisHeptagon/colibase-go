import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  base: '/',
  publicDir: 'public',
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:6700',
        changeOrigin: true,
      },
    },
  },
})
