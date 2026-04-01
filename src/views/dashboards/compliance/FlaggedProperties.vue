<script setup lang="ts">
import { useRoleStore } from '@/stores'
import { useRouter } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import { ref, computed } from 'vue'

const { selectedRole, user, logout } = useRoleStore()
const router = useRouter()
const handleLogout = () => { logout(); router.push('/') }

const searchQuery = ref('')
const filterPriority = ref('all')
const currentPage = ref(1)
const itemsPerPage = ref(5)

const showViewModal = ref(false)
const selectedFlag = ref<any>(null)
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })

const flaggedProperties = ref([
  { id: 'FLAG-001', property: 'Plot 8, Banana Island', owner: 'Folake Adeyemi', reason: 'Undeclared renovation', priority: 'Critical', status: 'Under Investigation', flaggedDate: '2024-01-15', investigator: 'Agent A' },
  { id: 'FLAG-002', property: 'Block 12, Maitama', owner: 'Unknown Entity', reason: 'Suspected tax evasion', priority: 'High', status: 'Pending Review', flaggedDate: '2024-01-12', investigator: 'Agent B' },
  { id: 'FLAG-003', property: 'Plot 45, VI', owner: 'Chinedu & Partners', reason: 'Value discrepancy', priority: 'Medium', status: 'Resolved', flaggedDate: '2024-01-08', investigator: 'Agent A' },
  { id: 'FLAG-004', property: 'Estate 7, Lekki', owner: 'Global Ventures Ltd', reason: 'Document forgery', priority: 'High', status: 'Under Investigation', flaggedDate: '2024-01-10', investigator: 'Agent C' },
])

const filteredProperties = computed(() => flaggedProperties.value.filter(p => {
  const matchesSearch = p.property.toLowerCase().includes(searchQuery.value.toLowerCase()) || p.owner.toLowerCase().includes(searchQuery.value.toLowerCase()) || p.id.toLowerCase().includes(searchQuery.value.toLowerCase())
  const matchesPriority = filterPriority.value === 'all' || p.priority === filterPriority.value
  return matchesSearch && matchesPriority
}))

const totalPages = computed(() => Math.ceil(filteredProperties.value.length / itemsPerPage.value))
const paginatedProperties = computed(() => { const start = (currentPage.value - 1) * itemsPerPage.value; return filteredProperties.value.slice(start, start + itemsPerPage.value) })
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }

const showToast = (message: string) => { toast.value = { show: true, message }; setTimeout(() => { toast.value.show = false }, 3000) }
const openViewModal = (f: any) => { selectedFlag.value = f; showViewModal.value = true }
const startInvestigation = (f: any) => { showToast(`Investigation started for ${f.id}`) }
</script>

<template>
  <div class="min-h-screen flex bg-[#f5f7fa]">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />
    <div class="flex-1 flex flex-col">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6 shrink-0">
        <div class="flex items-center gap-4"><span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span><span class="text-[#d1d5db]">/</span><span class="text-[#1f2937] text-sm font-medium">Flagged Properties</span></div>
        <div class="flex items-center gap-4"><span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span><button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button></div>
      </header>
      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb]"><h2 class="text-[13px] font-semibold text-[#1f2937]">Flagged Properties</h2></div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4">
            <input v-model="searchQuery" type="text" placeholder="Search..." class="input-field max-w-md" />
            <select v-model="filterPriority" class="input-field w-48"><option value="all">All Priority</option><option value="Critical">Critical</option><option value="High">High</option><option value="Medium">Medium</option></select>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead><tr><th class="table-header">Flag ID</th><th class="table-header">Property</th><th class="table-header">Owner</th><th class="table-header">Reason</th><th class="table-header">Priority</th><th class="table-header">Status</th><th class="table-header">Date</th><th class="table-header">Investigator</th><th class="table-header">Actions</th></tr></thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="flag in paginatedProperties" :key="flag.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ flag.id }}</td><td class="table-cell text-[#6b7280]">{{ flag.property }}</td><td class="table-cell">{{ flag.owner }}</td><td class="table-cell">{{ flag.reason }}</td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-red-50 text-red-700': flag.priority === 'Critical', 'bg-orange-50 text-orange-700': flag.priority === 'High', 'bg-yellow-50 text-yellow-700': flag.priority === 'Medium'}">{{ flag.priority }}</span></td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-blue-50 text-blue-700': flag.status === 'Under Investigation', 'bg-yellow-50 text-yellow-700': flag.status === 'Pending Review', 'bg-green-50 text-green-700': flag.status === 'Resolved'}">{{ flag.status }}</span></td>
                  <td class="table-cell text-[#9ca3af]">{{ flag.flaggedDate }}</td><td class="table-cell">{{ flag.investigator }}</td>
                  <td class="table-cell"><div class="flex gap-2"><button @click="openViewModal(flag)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">View</button><button @click="startInvestigation(flag)" class="px-3 py-1 text-[11px] bg-[#B90B0B] text-white rounded hover:bg-[#6a0707]">Investigate</button></div></td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="px-6 py-4 flex items-center justify-between border-t border-[#f3f4f6]">
            <p class="text-[11px] text-[#6b7280]">Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, filteredProperties.length) }} of {{ filteredProperties.length }} entries</p>
            <div class="flex items-center gap-1">
              <button @click="goToPage(currentPage - 1)" :disabled="currentPage === 1" class="px-3 py-1 text-[11px] border border-[#e5e7eb] rounded disabled:opacity-50">Prev</button>
              <button v-for="p in totalPages" :key="p" @click="goToPage(p)" class="px-3 py-1 text-[11px] border border-[#e5e7eb] rounded" :class="currentPage === p ? 'bg-[#1f2937] text-white' : ''">{{ p }}</button>
              <button @click="goToPage(currentPage + 1)" :disabled="currentPage === totalPages" class="px-3 py-1 text-[11px] border border-[#e5e7eb] rounded disabled:opacity-50">Next</button>
            </div>
          </div>
        </div>
      </main>
    </div>

    <Teleport to="body"><div v-if="toast.show" class="fixed bottom-4 right-4 bg-green-600 text-white px-4 py-2 rounded-lg shadow-lg z-50">{{ toast.message }}</div></Teleport>

    <Teleport to="body">
      <div v-if="showViewModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#B90B0B] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Flag Details</h3><button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div><p class="text-[11px] text-gray-500">Flag ID</p><p class="text-[13px] font-medium">{{ selectedFlag?.id }}</p></div>
              <div><p class="text-[11px] text-gray-500">Status</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-blue-50 text-blue-700': selectedFlag?.status === 'Under Investigation', 'bg-yellow-50 text-yellow-700': selectedFlag?.status === 'Pending Review'}">{{ selectedFlag?.status }}</span></div>
              <div><p class="text-[11px] text-gray-500">Priority</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-red-50 text-red-700': selectedFlag?.priority === 'Critical', 'bg-orange-50 text-orange-700': selectedFlag?.priority === 'High'}">{{ selectedFlag?.priority }}</span></div>
              <div><p class="text-[11px] text-gray-500">Investigator</p><p class="text-[13px]">{{ selectedFlag?.investigator }}</p></div>
            </div>
            <div><p class="text-[11px] text-gray-500">Property</p><p class="text-[13px]">{{ selectedFlag?.property }}</p></div>
            <div><p class="text-[11px] text-gray-500">Owner</p><p class="text-[13px]">{{ selectedFlag?.owner }}</p></div>
            <div><p class="text-[11px] text-gray-500">Reason</p><p class="text-[13px]">{{ selectedFlag?.reason }}</p></div>
            <div><p class="text-[11px] text-gray-500">Flagged Date</p><p class="text-[13px]">{{ selectedFlag?.flaggedDate }}</p></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-3">
            <button @click="showViewModal = false" class="px-4 py-2 text-[11px] bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200">Close</button>
            <button @click="startInvestigation(selectedFlag); showViewModal = false" class="px-4 py-2 text-[11px] bg-[#B90B0B] text-white rounded-lg hover:bg-[#991010]">Start Investigation</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
