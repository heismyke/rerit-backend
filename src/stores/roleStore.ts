import { ref } from 'vue'
import type { Role, RoleSlug } from '@/types'

const selectedRole = ref<Role | null>(null)
const isAuthenticated = ref(false)
const user = ref<{ email: string; name: string } | null>(null)

export const useRoleStore = () => {
  const setRole = (role: Role) => {
    selectedRole.value = role
  }

  const clearRole = () => {
    selectedRole.value = null
  }

  const login = (email: string, name: string) => {
    isAuthenticated.value = true
    user.value = { email, name }
  }

  const logout = () => {
    isAuthenticated.value = false
    user.value = null
    selectedRole.value = null
  }

  return {
    selectedRole,
    isAuthenticated,
    user,
    setRole,
    clearRole,
    login,
    logout,
  }
}
