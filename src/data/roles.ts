import type { Role, RoleSlug } from '@/types'

export const roles: Role[] = [
  {
    id: 'admin',
    name: 'Admin',
    description: 'System control center and management',
    icon: '◉',
  },
  {
    id: 'auditor',
    name: 'Auditor',
    description: 'Audit property records and tax compliance',
    icon: '◈',
  },
  {
    id: 'developer',
    name: 'Real Estate Developer (Tax Payer)',
    description: 'Manage your property portfolio and tax obligations',
    icon: '◫',
  },
  {
    id: 'surveyor',
    name: 'Surveyor',
    description: 'Submit property surveys and valuations',
    icon: '◎',
  },
  {
    id: 'compliance',
    name: 'Compliance Officer (NRS Officer)',
    description: 'Ensure regulatory compliance and enforcement',
    icon: '◉',
  },
]

export const getRoleById = (id: RoleSlug): Role | undefined => {
  return roles.find(role => role.id === id)
}

export const getRouteByRoleId = (id: RoleSlug): string => {
  const routes: Record<RoleSlug, string> = {
    admin: '/admin',
    auditor: '/auditor',
    developer: '/developer',
    surveyor: '/surveyor',
    compliance: '/compliance',
  }
  return routes[id]
}
