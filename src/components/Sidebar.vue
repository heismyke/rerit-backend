<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'
import type { RoleSlug } from '@/types'
import { useRoleStore } from '@/stores'

defineProps<{
  roleId: RoleSlug
}>()

const { user } = useRoleStore()
const router = useRouter()
const route = useRoute()

const auditorMenu = [
  { label: 'Overview', path: '/auditor' },
  { label: 'Properties', path: '/auditor/properties' },
  { label: 'Audit Cases', path: '/auditor/audit-cases' },
  { label: 'Land Registry', path: '/auditor/land-registry' },
  { label: 'Settings', path: '/auditor/settings' },
]

const adminMenu = [
  { label: 'Dashboard', path: '/admin' },
  { label: 'Rules & Automation', path: '/admin/rules' },
  { label: 'Properties', path: '/admin/properties' },
  { label: 'Taxpayers', path: '/admin/taxpayers' },
  { label: 'Revenue & Payments', path: '/admin/revenue' },
  { label: 'Cases', path: '/admin/cases' },
  { label: 'Audits', path: '/admin/audits' },
  { label: 'Surveys', path: '/admin/surveys' },
  { label: 'Compliance', path: '/admin/compliance' },
  { label: 'Reports', path: '/admin/reports' },
  { label: 'Users & Roles', path: '/admin/users' },
  { label: 'Notifications', path: '/admin/notifications' },
  { label: 'Settings', path: '/admin/settings' },
]

const developerMenu = [
  { label: 'Overview', path: '/developer' },
  { label: 'Properties', path: '/developer/properties' },
  { label: 'Notices', path: '/developer/notices' },
  { label: 'Payments', path: '/developer/payments' },
  { label: 'Settings', path: '/developer/settings' },
]

const surveyorMenu = [
  { label: 'Overview', path: '/surveyor' },
  { label: 'Surveys', path: '/surveyor/surveys' },
  { label: 'Submissions', path: '/surveyor/submissions' },
  { label: 'Settings', path: '/surveyor/settings' },
]

const complianceMenu = [
  { label: 'Overview', path: '/compliance' },
  { label: 'Properties', path: '/compliance/properties' },
  { label: 'Flagged', path: '/compliance/flagged' },
  { label: 'Registry', path: '/compliance/land-registry' },
  { label: 'Notes', path: '/compliance/notes' },
  { label: 'Settings', path: '/compliance/settings' },
]

const getMenu = (roleId: RoleSlug) => {
  switch (roleId) {
    case 'auditor':
      return auditorMenu
    case 'admin':
      return adminMenu
    case 'developer':
      return developerMenu
    case 'surveyor':
      return surveyorMenu
    case 'compliance':
      return complianceMenu
    default:
      return []
  }
}
</script>

<template>
  <aside class="sidebar w-56 bg-[#EEEEEE] flex flex-col">
    <div class="px-4 py-4 border-b border-gray-300">
      <div class="flex items-center gap-3 mb-4">
        <div class="w-10 h-10 rounded-full bg-[#2D5A27] flex items-center justify-center">
          <span class="text-sm font-semibold text-white">{{ user?.name?.charAt(0) || 'U' }}</span>
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-[13px] font-medium text-[#1f2937] truncate">{{ user?.name || 'User' }}</p>
          <p class="text-[10px] text-[#6b7280] truncate">{{ user?.email || 'user@nrs.gov.ng' }}</p>
        </div>
      </div>
    </div>
    <div class="px-4 py-3 border-b border-gray-300">
      <h1 class="text-sm font-semibold text-[#2D5A27] tracking-tight">FCT-IRS</h1>
      <p class="text-[10px] text-[#6b7280] mt-0.5">Revenue System</p>
    </div>
    <nav class="flex-1 p-4 overflow-y-auto">
      <ul class="space-y-0.5">
        <li v-for="item in getMenu(roleId)" :key="item.path">
          <button
            @click="router.push(item.path)"
            class="w-full sidebar-item"
            :class="route.path === item.path ? 'sidebar-item-active' : 'sidebar-item-inactive'"
          >
            <span>{{ item.label }}</span>
          </button>
        </li>
      </ul>
    </nav>
    <div class="p-4 border-t border-gray-300 shrink-0">
      <div class="flex flex-col items-center gap-2">
        <img src="/fct-irs.svg" alt="FCT-IRS" class="h-8 opacity-50" />
        <p class="text-[11px] text-[#9ca3af]">v1.0.0</p>
      </div>
    </div>
  </aside>
</template>
