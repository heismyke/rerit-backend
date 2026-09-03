import { ref } from 'vue'
import type { Role, RoleSlug } from '@/types'
import { api } from '@/services/api'

const selectedRole = ref<Role | null>(null)
const isAuthenticated = ref(false)
const user = ref<{ email: string; name: string } | null>(null)
const authToken = ref<string | null>(null)

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

  const loginWithApi = async (email: string, password: string, role: Role) => {
    const response = await api.login({ email, password, role: role.id })
    authToken.value = response.token
    selectedRole.value = role
    isAuthenticated.value = true
    user.value = { email: response.user.email, name: response.user.name }
  }

  const logout = () => {
    isAuthenticated.value = false
    user.value = null
    selectedRole.value = null
    authToken.value = null
  }

  return {
    selectedRole,
    isAuthenticated,
    user,
    authToken,
    setRole,
    clearRole,
    login,
    loginWithApi,
    logout,
  }
}
