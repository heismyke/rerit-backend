<script setup lang="ts">
import { useRoleStore } from '@/stores'
import { useRouter } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import { ref, computed } from 'vue'

const { selectedRole, user, logout } = useRoleStore()
const router = useRouter()
const handleLogout = () => { logout(); router.push('/') }

const searchQuery = ref('')
const currentPage = ref(1)
const itemsPerPage = ref(5)

const successfulFilings = ref([
  { id: 'SUC-2024-101', filingId: 'FCT-IRS-00311', property: 'Plot 12, Jabi', taxpayer: 'Kola Ibrahim', validatedAt: '2024-01-18', status: 'Validated' },
  { id: 'SUC-2024-102', filingId: 'FCT-IRS-00318', property: 'Block 5, Garki', taxpayer: 'Laila Musa', validatedAt: '2024-01-18', status: 'Validated' },
  { id: 'SUC-2024-103', filingId: 'FCT-IRS-00327', property: 'Unit 3, Wuse I', taxpayer: 'Prime Estates Ltd', validatedAt: '2024-01-19', status: 'Validated' },
  { id: 'SUC-2024-104', filingId: 'FCT-IRS-00333', property: 'Plot 6, Maitama', taxpayer: 'Uche Nwankwo', validatedAt: '2024-01-20', status: 'Validated' },
])

const filteredFilings = computed(() => successfulFilings.value.filter(f => {
  const q = searchQuery.value.toLowerCase()
  return f.id.toLowerCase().includes(q)
    || f.filingId.toLowerCase().includes(q)
    || f.property.toLowerCase().includes(q)
    || f.taxpayer.toLowerCase().includes(q)
}))

const totalPages = computed(() => Math.ceil(filteredFilings.value.length / itemsPerPage.value))
const paginatedFilings = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  return filteredFilings.value.slice(start, start + itemsPerPage.value)
})
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }
</script>

<template>
  <div class="min-h-screen bg-[#f5f6fa] flex">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />
    <div class="flex-1 flex flex-col">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6 shrink-0">
        <div class="flex items-center gap-4">
          <span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span>
          <span class="text-[#d1d5db]">/</span>
          <span class="text-[#1f2937] text-sm font-medium">Successful Filings</span>
        </div>
        <div class="flex items-center gap-4">
          <span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span>
          <button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button>
        </div>
      </header>

      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between">
            <h2 class="text-[13px] font-semibold text-[#1f2937]">Successful Filings</h2>
            <span class="text-[11px] text-[#6b7280]">Validated by RERIT</span>
          </div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4">
            <input v-model="searchQuery" type="text" placeholder="Search by ID, filing, property, or taxpayer..." class="input-field flex-1" />
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr>
                  <th class="table-header">Record ID</th>
                  <th class="table-header">FCT‑IRS Filing ID</th>
                  <th class="table-header">Property</th>
                  <th class="table-header">Taxpayer</th>
                  <th class="table-header">Validated</th>
                  <th class="table-header">Status</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="filing in paginatedFilings" :key="filing.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ filing.id }}</td>
                  <td class="table-cell">{{ filing.filingId }}</td>
                  <td class="table-cell">{{ filing.property }}</td>
                  <td class="table-cell">{{ filing.taxpayer }}</td>
                  <td class="table-cell text-[#9ca3af]">{{ filing.validatedAt }}</td>
                  <td class="table-cell">
                    <span class="px-2 py-0.5 text-[11px] font-medium rounded-full bg-green-50 text-green-700">{{ filing.status }}</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="px-6 py-4 flex items-center justify-between border-t border-[#f3f4f6]">
            <p class="text-[11px] text-[#6b7280]">Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, filteredFilings.length) }} of {{ filteredFilings.length }} entries</p>
            <div class="flex items-center gap-1">
              <button @click="goToPage(currentPage - 1)" :disabled="currentPage === 1" class="px-3 py-1 text-[11px] border border-[#e5e7eb] rounded disabled:opacity-50">Prev</button>
              <button v-for="p in totalPages" :key="p" @click="goToPage(p)" class="px-3 py-1 text-[11px] border border-[#e5e7eb] rounded" :class="currentPage === p ? 'bg-[#1f2937] text-white' : ''">{{ p }}</button>
              <button @click="goToPage(currentPage + 1)" :disabled="currentPage === totalPages" class="px-3 py-1 text-[11px] border border-[#e5e7eb] rounded disabled:opacity-50">Next</button>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>
