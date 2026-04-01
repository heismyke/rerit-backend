import { createRouter, createWebHistory } from 'vue-router'
import publicRoutes from './public'

const router = createRouter({
  history: createWebHistory(),
  routes: [...publicRoutes],
})

export default router
