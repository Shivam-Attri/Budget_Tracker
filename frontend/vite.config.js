import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    // This is needed for Docker container port mapping
    host: true, 
    port: 5173, // This is the port the dev server will run on
    // Optional: Enables hot-reloading in some Docker setups
    watch: {
      usePolling: true,
    },
  },
  build: {
    // This ensures the output target is modern and compatible
    target: 'esnext'
  }
})
