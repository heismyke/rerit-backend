export type RoleSlug = 'admin' | 'auditor' | 'taxpayer' | 'compliance'

export interface Role {
  id: RoleSlug
  name: string
  description: string
  icon: string
}
