<script setup lang="ts">
import { useRoleStore } from '@/stores'
import { useRouter } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import { ref, computed } from 'vue'

const { selectedRole, user, logout } = useRoleStore()
const router = useRouter()
const handleLogout = () => { logout(); router.push('/') }

const searchQuery = ref('')
const filterStatus = ref('all')
const currentPage = ref(1)
const itemsPerPage = ref(5)

const showAddModal = ref(false)
const showViewModal = ref(false)
const selectedProperty = ref<any>(null)
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })

const newProperty = ref({ owner: '', address: '', type: 'Residential', value: '' })

const properties = ref([
  { id: 'PROP-001', owner: 'Emeka Okonkwo', address: 'Plot 42, Victoria Island', type: 'Commercial', value: 'N250,000,000', status: 'Verified', lastAudit: '2024-01-10' },
  { id: 'PROP-002', owner: 'Adaobi Nnamdi', address: 'Block 7, Lekki Phase 2', type: 'Residential', value: 'N80,000,000', status: 'Pending', lastAudit: '2024-01-12' },
  { id: 'PROP-003', owner: 'Chidi Okafor', address: '15 Admiralty Way, Lekki', type: 'Mixed Use', value: 'N180,000,000', status: 'Verified', lastAudit: '2024-01-08' },
  { id: 'PROP-004', owner: 'Folake Adeyemi', address: 'Plot 8, Banana Island', type: 'Residential', value: 'N500,000,000', status: 'Flagged', lastAudit: '2024-01-15' },
  { id: 'PROP-005', owner: 'Ibrahim Bello', address: 'Block 3, Ikoyi', type: 'Commercial', value: 'N350,000,000', status: 'Verified', lastAudit: '2024-01-05' },
  { id: 'PROP-006', owner: 'Oluwaseun Adeyinka', address: 'Plot 15, Maitama', type: 'Land', value: 'N120,000,000', status: 'Verified', lastAudit: '2024-01-03' },
  { id: 'PROP-007', owner: 'Nigerian Holdings Ltd', address: '21 Broad Street, Lagos', type: 'Commercial', value: 'N890,000,000', status: 'Pending', lastAudit: '2024-01-18' },
])

const filteredProperties = computed(() => properties.value.filter(p => {
  const matchesSearch = p.owner.toLowerCase().includes(searchQuery.value.toLowerCase()) || p.address.toLowerCase().includes(searchQuery.value.toLowerCase()) || p.id.toLowerCase().includes(searchQuery.value.toLowerCase())
  const matchesStatus = filterStatus.value === 'all' || p.status === filterStatus.value
  return matchesSearch && matchesStatus
}))

const totalPages = computed(() => Math.ceil(filteredProperties.value.length / itemsPerPage.value))
const paginatedProperties = computed(() => { const start = (currentPage.value - 1) * itemsPerPage.value; return filteredProperties.value.slice(start, start + itemsPerPage.value) })
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }

const showToast = (message: string) => { toast.value = { show: true, message }; setTimeout(() => { toast.value.show = false }, 3000) }
const openViewModal = (p: any) => { selectedProperty.value = p; showViewModal.value = true }
const handleAddProperty = () => {
  const newId = 'PROP-' + String(properties.value.length + 1).padStart(3, '0')
  properties.value.unshift({ id: newId, ...newProperty.value, status: 'Pending', lastAudit: new Date().toISOString().split('T')[0] })
  showAddModal.value = false; newProperty.value = { owner: '', address: '', type: 'Residential', value: '' }; showToast('Property added successfully')
}
const startAudit = (p: any) => { showToast(`Audit started for ${p.id}`) }
</script>

<template>
  <div class="min-h-screen flex bg-[#f5f7fa]">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />
    <div class="flex-1 flex flex-col">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6 shrink-0">
        <div class="flex items-center gap-4"><span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span><span class="text-[#d1d5db]">/</span><span class="text-[#1f2937] text-sm font-medium">Property Listings</span></div>
        <div class="flex items-center gap-4"><span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span><button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button></div>
      </header>
      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between"><h2 class="text-[13px] font-semibold text-[#1f2937]">Properties</h2><button @click="showAddModal = true" class="btn-primary text-[11px]">Add Property</button></div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4">
            <input v-model="searchQuery" type="text" placeholder="Search by owner, address, or ID..." class="input-field max-w-md" />
            <select v-model="filterStatus" class="input-field w-48"><option value="all">All Status</option><option value="Verified">Verified</option><option value="Pending">Pending</option><option value="Flagged">Flagged</option></select>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead><tr><th class="table-header">Property ID</th><th class="table-header">Owner</th><th class="table-header">Address</th><th class="table-header">Type</th><th class="table-header">Value</th><th class="table-header">Status</th><th class="table-header">Last Audit</th><th class="table-header">Actions</th></tr></thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="property in paginatedProperties" :key="property.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ property.id }}</td><td class="table-cell">{{ property.owner }}</td><td class="table-cell text-[#6b7280]">{{ property.address }}</td><td class="table-cell">{{ property.type }}</td><td class="table-cell">{{ property.value }}</td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': property.status === 'Verified', 'bg-yellow-50 text-yellow-700': property.status === 'Pending', 'bg-red-50 text-red-700': property.status === 'Flagged'}">{{ property.status }}</span></td>
                  <td class="table-cell text-[#9ca3af]">{{ property.lastAudit }}</td>
                  <td class="table-cell"><div class="flex gap-2"><button @click="openViewModal(property)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">View</button><button @click="startAudit(property)" class="px-3 py-1 text-[11px] bg-[#B90B0B] text-white rounded hover:bg-[#6a0707]">Audit</button></div></td>
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
      <div v-if="showAddModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#B90B0B] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Add Property</h3><button @click="showAddModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Owner</label><input v-model="newProperty.owner" type="text" placeholder="Enter owner name" class="input-field w-full" /></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Address</label><input v-model="newProperty.address" type="text" placeholder="Enter address" class="input-field w-full" /></div>
            <div class="grid grid-cols-2 gap-4">
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Type</label><select v-model="newProperty.type" class="input-field w-full"><option>Residential</option><option>Commercial</option><option>Mixed Use</option><option>Land</option></select></div>
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Value (NGN)</label><input v-model="newProperty.value" type="text" placeholder="N0" class="input-field w-full" /></div>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-end"><button @click="showAddModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button><button @click="handleAddProperty" class="px-4 py-2 text-[11px] bg-[#B90B0B] text-white rounded-lg hover:bg-[#991010]">Add Property</button></div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showViewModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#B90B0B] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Property Details</h3><button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div><p class="text-[11px] text-gray-500">Property ID</p><p class="text-[13px] font-medium">{{ selectedProperty?.id }}</p></div>
              <div><p class="text-[11px] text-gray-500">Status</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': selectedProperty?.status === 'Verified', 'bg-yellow-50 text-yellow-700': selectedProperty?.status === 'Pending', 'bg-red-50 text-red-700': selectedProperty?.status === 'Flagged'}">{{ selectedProperty?.status }}</span></div>
              <div><p class="text-[11px] text-gray-500">Owner</p><p class="text-[13px] font-medium">{{ selectedProperty?.owner }}</p></div>
              <div><p class="text-[11px] text-gray-500">Type</p><p class="text-[13px]">{{ selectedProperty?.type }}</p></div>
            </div>
            <div><p class="text-[11px] text-gray-500">Address</p><p class="text-[13px]">{{ selectedProperty?.address }}</p></div>
            <div><p class="text-[11px] text-gray-500">Value</p><p class="text-[13px] font-medium text-green-700">{{ selectedProperty?.value }}</p></div>
            <div><p class="text-[11px] text-gray-500">Last Audit</p><p class="text-[13px]">{{ selectedProperty?.lastAudit }}</p></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-3">
            <button @click="showViewModal = false" class="px-4 py-2 text-[11px] bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200">Close</button>
            <button @click="startAudit(selectedProperty); showViewModal = false" class="px-4 py-2 text-[11px] bg-[#B90B0B] text-white rounded-lg hover:bg-[#991010]">Start Audit</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
