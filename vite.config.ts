import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
    base: '/',
    root: './web',
    publicDir: './web/public',
    server: {
        host: '127.0.0.1',
        port: 3032,
    },
    build: {
        outDir: './web/dist',
    },
    plugins: [vue()],
})
