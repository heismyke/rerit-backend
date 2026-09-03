import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'

window.addEventListener('vite:preloadError', () => {
  const reloadKey = 'rerit:chunk-reload'
  if (sessionStorage.getItem(reloadKey)) return
  sessionStorage.setItem(reloadKey, '1')
  window.location.reload()
})

const app = createApp(App)
app.use(router)
router.isReady().then(() => {
  sessionStorage.removeItem('rerit:chunk-reload')
  app.mount('#app')
})
