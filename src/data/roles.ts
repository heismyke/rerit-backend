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
    id: 'taxpayer',
    name: 'Tax Payer',
    description: 'Manage your property portfolio and tax obligations',
    icon: '◫',
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
    taxpayer: '/taxpayer',
    compliance: '/compliance',
  }
  return routes[id]
}
