import { createRouter, createWebHistory } from 'vue-router'
import publicRoutes from './public'
import { getRoleById } from '@/data'
import { useRoleStore } from '@/stores'
import type { RoleSlug } from '@/types'

const router = createRouter({
  history: createWebHistory(),
  routes: [...publicRoutes],
})

const demoUsers: Record<RoleSlug, { email: string; name: string }> = {
  admin: { email: 'admin@fctirs.gov.ng', name: 'Admin' },
  auditor: { email: 'auditor@fctirs.gov.ng', name: 'Auditor' },
  taxpayer: { email: 'taxpayer@example.com', name: 'Taxpayer' },
}

router.beforeEach((to) => {
  const access = to.meta.access
  if (access !== 'admin' && access !== 'auditor' && access !== 'taxpayer') return true

  const roleId = access as RoleSlug
  const role = getRoleById(roleId)
  if (!role) return true

  const roleStore = useRoleStore()
  roleStore.setRole(role)
  if (!roleStore.user.value) {
    const demoUser = demoUsers[roleId]
    roleStore.login(demoUser.email, demoUser.name)
  }

  return true
})

export default router
