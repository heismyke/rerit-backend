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
const activeTab = ref('all')

const showViewModal = ref(false)
const selectedItem = ref<any>(null)
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })

const complianceItems = ref([
  { id: 'C-001', property: 'Plot 45, Victoria Island', owner: 'Alhaji Garba & Sons', type: 'Commercial', status: 'Compliant', lastChecked: '2024-01-15', riskLevel: 'Low' },
  { id: 'C-002', property: 'Block 7, Lekki Phase 2', owner: 'Mrs. Folake Adeyemi', type: 'Residential', status: 'Pending Review', lastChecked: '2024-01-14', riskLevel: 'Medium' },
  { id: 'C-003', property: 'Suite 301, Marina Towers', owner: 'TechStart Ltd', type: 'Commercial', status: 'Non-Compliant', lastChecked: '2024-01-13', riskLevel: 'High' },
  { id: 'C-004', property: 'Plot 12, Ikoyi', owner: 'Chief Okafor', type: 'Residential', status: 'Compliant', lastChecked: '2024-01-12', riskLevel: 'Low' },
  { id: 'C-005', property: 'Block C, Victoria Island', owner: 'Global Finance Corp', type: 'Commercial', status: 'Pending Review', lastChecked: '2024-01-11', riskLevel: 'Medium' },
  { id: 'C-006', property: 'Flat 5, Banana Island', owner: 'Dr. Emeka Okonkwo', type: 'Residential', status: 'Compliant', lastChecked: '2024-01-10', riskLevel: 'Low' },
  { id: 'C-007', property: 'Plot 88, Wuse Zone 5', owner: 'Nigerian Holdings Ltd', type: 'Commercial', status: 'Non-Compliant', lastChecked: '2024-01-09', riskLevel: 'High' },
])

const filteredItems = computed(() => {
  let items = complianceItems.value.filter(i => i.property.toLowerCase().includes(searchQuery.value.toLowerCase()) || i.owner.toLowerCase().includes(searchQuery.value.toLowerCase()) || i.id.toLowerCase().includes(searchQuery.value.toLowerCase()))
  if (activeTab.value !== 'all') { items = items.filter(i => i.status.toLowerCase().replace(' ', '-') === activeTab.value) }
  return items
})

const totalPages = computed(() => Math.ceil(filteredItems.value.length / itemsPerPage.value))
const paginatedItems = computed(() => { const start = (currentPage.value - 1) * itemsPerPage.value; return filteredItems.value.slice(start, start + itemsPerPage.value) })
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }

const stats = computed(() => ({ total: complianceItems.value.length, compliant: complianceItems.value.filter(i => i.status === 'Compliant').length, pending: complianceItems.value.filter(i => i.status === 'Pending Review').length, nonCompliant: complianceItems.value.filter(i => i.status === 'Non-Compliant').length }))

const showToast = (message: string) => { toast.value = { show: true, message }; setTimeout(() => { toast.value.show = false }, 3000) }
const openViewModal = (item: any) => { selectedItem.value = item; showViewModal.value = true }
const runComplianceCheck = () => { showToast('Compliance check initiated for all properties') }
</script>

<template>
  <div class="min-h-screen flex bg-[#f5f7fa]">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />
    <div class="flex-1">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6">
        <div class="flex items-center gap-4"><span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span><span class="text-[#d1d5db]">/</span><span class="text-[#1f2937] text-sm font-medium">Compliance Overview</span></div>
        <div class="flex items-center gap-4"><span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span><button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button></div>
      </header>
      <main class="flex-1 p-6">
        <div class="grid grid-cols-4 gap-4 mb-6">
          <div class="bg-white border border-[#e5e7eb] rounded-lg p-4"><p class="text-[11px] text-[#6b7280]">Total Properties</p><p class="text-2xl font-semibold text-[#1f2937] mt-1">{{ stats.total }}</p></div>
          <div class="bg-white border border-[#e5e7eb] rounded-lg p-4"><p class="text-[11px] text-[#6b7280]">Compliant</p><p class="text-2xl font-semibold text-green-600 mt-1">{{ stats.compliant }}</p></div>
          <div class="bg-white border border-[#e5e7eb] rounded-lg p-4"><p class="text-[11px] text-[#6b7280]">Pending Review</p><p class="text-2xl font-semibold text-yellow-600 mt-1">{{ stats.pending }}</p></div>
          <div class="bg-white border border-[#e5e7eb] rounded-lg p-4"><p class="text-[11px] text-[#6b7280]">Non-Compliant</p><p class="text-2xl font-semibold text-red-600 mt-1">{{ stats.nonCompliant }}</p></div>
        </div>
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between">
            <div class="flex gap-4">
              <button @click="activeTab = 'all'" class="text-[13px] font-medium pb-1 border-b-2" :class="activeTab === 'all' ? 'border-[#B90B0B] text-[#B90B0B]' : 'border-transparent text-[#6b7280]'">All</button>
              <button @click="activeTab = 'compliant'" class="text-[13px] font-medium pb-1 border-b-2" :class="activeTab === 'compliant' ? 'border-[#B90B0B] text-[#B90B0B]' : 'border-transparent text-[#6b7280]'">Compliant</button>
              <button @click="activeTab = 'pending-review'" class="text-[13px] font-medium pb-1 border-b-2" :class="activeTab === 'pending-review' ? 'border-[#B90B0B] text-[#B90B0B]' : 'border-transparent text-[#6b7280]'">Pending Review</button>
              <button @click="activeTab = 'non-compliant'" class="text-[13px] font-medium pb-1 border-b-2" :class="activeTab === 'non-compliant' ? 'border-[#B90B0B] text-[#B90B0B]' : 'border-transparent text-[#6b7280]'">Non-Compliant</button>
            </div>
            <button @click="runComplianceCheck" class="btn-primary text-[11px]">Run Compliance Check</button>
          </div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4"><input v-model="searchQuery" type="text" placeholder="Search by property, owner, or ID..." class="input-field max-w-md" /></div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead><tr><th class="table-header">ID</th><th class="table-header">Property</th><th class="table-header">Owner</th><th class="table-header">Type</th><th class="table-header">Status</th><th class="table-header">Risk Level</th><th class="table-header">Last Checked</th><th class="table-header">Actions</th></tr></thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="item in paginatedItems" :key="item.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ item.id }}</td><td class="table-cell">{{ item.property }}</td><td class="table-cell">{{ item.owner }}</td><td class="table-cell"><span class="px-2 py-0.5 text-[11px] bg-[#f3f4f6] text-[#374151] rounded">{{ item.type }}</span></td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': item.status === 'Compliant', 'bg-yellow-50 text-yellow-700': item.status === 'Pending Review', 'bg-red-50 text-red-700': item.status === 'Non-Compliant'}">{{ item.status }}</span></td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': item.riskLevel === 'Low', 'bg-yellow-50 text-yellow-700': item.riskLevel === 'Medium', 'bg-red-50 text-red-700': item.riskLevel === 'High'}">{{ item.riskLevel }}</span></td>
                  <td class="table-cell text-[#9ca3af]">{{ item.lastChecked }}</td>
                  <td class="table-cell"><div class="flex gap-2"><button @click="openViewModal(item)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">Review</button><button class="px-3 py-1 text-[11px] bg-[#B90B0B] text-white rounded hover:bg-[#991010]">Details</button></div></td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="px-6 py-4 flex items-center justify-between border-t border-[#f3f4f6]">
            <p class="text-[11px] text-[#6b7280]">Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, filteredItems.length) }} of {{ filteredItems.length }} entries</p>
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
          <div class="bg-[#B90B0B] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Compliance Details</h3><button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div><p class="text-[11px] text-gray-500">ID</p><p class="text-[13px] font-medium">{{ selectedItem?.id }}</p></div>
              <div><p class="text-[11px] text-gray-500">Status</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': selectedItem?.status === 'Compliant', 'bg-yellow-50 text-yellow-700': selectedItem?.status === 'Pending Review', 'bg-red-50 text-red-700': selectedItem?.status === 'Non-Compliant'}">{{ selectedItem?.status }}</span></div>
              <div><p class="text-[11px] text-gray-500">Property</p><p class="text-[13px] font-medium">{{ selectedItem?.property }}</p></div>
              <div><p class="text-[11px] text-gray-500">Risk Level</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': selectedItem?.riskLevel === 'Low', 'bg-yellow-50 text-yellow-700': selectedItem?.riskLevel === 'Medium', 'bg-red-50 text-red-700': selectedItem?.riskLevel === 'High'}">{{ selectedItem?.riskLevel }}</span></div>
            </div>
            <div><p class="text-[11px] text-gray-500">Owner</p><p class="text-[13px]">{{ selectedItem?.owner }}</p></div>
            <div><p class="text-[11px] text-gray-500">Type</p><p class="text-[13px]">{{ selectedItem?.type }}</p></div>
            <div><p class="text-[11px] text-gray-500">Last Checked</p><p class="text-[13px]">{{ selectedItem?.lastChecked }}</p></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-3">
            <button @click="showViewModal = false" class="px-4 py-2 text-[11px] bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200">Close</button>
            <button @click="showToast('Compliance review submitted')" class="px-4 py-2 text-[11px] bg-[#B90B0B] text-white rounded-lg hover:bg-[#991010]">Submit Review</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
