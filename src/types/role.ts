export type RoleSlug = 'admin' | 'auditor' | 'developer' | 'surveyor' | 'compliance'

export interface Role {
  id: RoleSlug
  name: string
  description: string
  icon: string
}
