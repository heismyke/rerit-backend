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
const showEditModal = ref(false)
const showViewModal = ref(false)
const selectedRecord = ref<any>(null)
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })

const newRecord = ref({ plotNo: '', block: '', location: '', size: '', owner: '', titleType: 'C of O' })
const editRecord = ref({ plotNo: '', block: '', location: '', size: '', owner: '', titleType: 'C of O', status: 'Active' })

const landRecords = ref([
  { id: 'LR-001', plotNo: 'Plot 42', block: 'Block A', location: 'Victoria Island', size: '500 sqm', owner: 'Emeka Okonkwo', titleType: 'C of O', status: 'Active', lastUpdated: '2024-01-10' },
  { id: 'LR-002', plotNo: 'Plot 15', block: 'Block C', location: 'Lekki Phase 1', size: '750 sqm', owner: 'Adaobi Nnamdi', titleType: 'Survey Plan', status: 'Pending', lastUpdated: '2024-01-12' },
  { id: 'LR-003', plotNo: 'Plot 88', block: 'Block B', location: 'Ikoyi', size: '1,200 sqm', owner: 'Chidi Okafor', titleType: 'C of O', status: 'Active', lastUpdated: '2024-01-08' },
  { id: 'LR-004', plotNo: 'Plot 8', block: 'Block D', location: 'Banana Island', size: '600 sqm', owner: 'Folake Adeyemi', titleType: 'Excision', status: 'Flagged', lastUpdated: '2024-01-15' },
])

const filteredRecords = computed(() => landRecords.value.filter(r => {
  const matchesSearch = r.owner.toLowerCase().includes(searchQuery.value.toLowerCase()) || r.location.toLowerCase().includes(searchQuery.value.toLowerCase()) || r.id.toLowerCase().includes(searchQuery.value.toLowerCase())
  const matchesStatus = filterStatus.value === 'all' || r.status === filterStatus.value
  return matchesSearch && matchesStatus
}))

const totalPages = computed(() => Math.ceil(filteredRecords.value.length / itemsPerPage.value))
const paginatedRecords = computed(() => { const start = (currentPage.value - 1) * itemsPerPage.value; return filteredRecords.value.slice(start, start + itemsPerPage.value) })
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }

const showToast = (message: string) => { toast.value = { show: true, message }; setTimeout(() => { toast.value.show = false }, 3000) }
const openViewModal = (r: any) => { selectedRecord.value = r; showViewModal.value = true }
const openEditModal = (r: any) => { selectedRecord.value = r; editRecord.value = { ...r }; showEditModal.value = true }

const handleAddRecord = () => {
  const newId = 'LR-' + String(landRecords.value.length + 1).padStart(3, '0')
  landRecords.value.unshift({ id: newId, ...newRecord.value, status: 'Pending', lastUpdated: new Date().toISOString().split('T')[0] })
  showAddModal.value = false; newRecord.value = { plotNo: '', block: '', location: '', size: '', owner: '', titleType: 'C of O' }; showToast('Record added successfully')
}

const handleUpdateRecord = () => {
  const index = landRecords.value.findIndex(r => r.id === selectedRecord.value.id)
  if (index !== -1) { landRecords.value[index] = { ...landRecords.value[index], ...editRecord.value }; showToast('Record updated') }
  showEditModal.value = false
}
</script>

<template>
  <div class="min-h-screen flex bg-[#f5f7fa]">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />
    <div class="flex-1 flex flex-col">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6 shrink-0">
        <div class="flex items-center gap-4"><span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span><span class="text-[#d1d5db]">/</span><span class="text-[#1f2937] text-sm font-medium">Land Registry</span></div>
        <div class="flex items-center gap-4"><span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span><button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button></div>
      </header>
      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between"><h2 class="text-[13px] font-semibold text-[#1f2937]">Land Registry</h2><button @click="showAddModal = true" class="btn-primary text-[11px]">Add Record</button></div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4">
            <input v-model="searchQuery" type="text" placeholder="Search..." class="input-field max-w-md" />
            <select v-model="filterStatus" class="input-field w-48"><option value="all">All Status</option><option value="Active">Active</option><option value="Pending">Pending</option><option value="Flagged">Flagged</option></select>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead><tr><th class="table-header">Record ID</th><th class="table-header">Plot No</th><th class="table-header">Block</th><th class="table-header">Location</th><th class="table-header">Size</th><th class="table-header">Owner</th><th class="table-header">Title</th><th class="table-header">Status</th><th class="table-header">Updated</th><th class="table-header">Actions</th></tr></thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="record in paginatedRecords" :key="record.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ record.id }}</td><td class="table-cell">{{ record.plotNo }}</td><td class="table-cell">{{ record.block }}</td><td class="table-cell text-[#6b7280]">{{ record.location }}</td><td class="table-cell">{{ record.size }}</td><td class="table-cell">{{ record.owner }}</td><td class="table-cell"><span class="px-2 py-0.5 text-[11px] bg-[#f3f4f6] text-[#6b7280] rounded">{{ record.titleType }}</span></td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': record.status === 'Active', 'bg-yellow-50 text-yellow-700': record.status === 'Pending', 'bg-red-50 text-red-700': record.status === 'Flagged'}">{{ record.status }}</span></td>
                  <td class="table-cell text-[#9ca3af]">{{ record.lastUpdated }}</td>
                  <td class="table-cell"><div class="flex gap-2"><button @click="openViewModal(record)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">View</button><button @click="openEditModal(record)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">Edit</button></div></td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="px-6 py-4 flex items-center justify-between border-t border-[#f3f4f6]">
            <p class="text-[11px] text-[#6b7280]">Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, filteredRecords.length) }} of {{ filteredRecords.length }} entries</p>
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
          <div class="bg-[#B90B0B] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Add Record</h3><button @click="showAddModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4"><div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Plot Number</label><input v-model="newRecord.plotNo" type="text" placeholder="e.g. Plot 42" class="input-field w-full" /></div><div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Block</label><input v-model="newRecord.block" type="text" class="input-field w-full" /></div></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Location</label><input v-model="newRecord.location" type="text" class="input-field w-full" /></div>
            <div class="grid grid-cols-2 gap-4"><div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Size</label><input v-model="newRecord.size" type="text" class="input-field w-full" /></div><div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Title Type</label><select v-model="newRecord.titleType" class="input-field w-full"><option>C of O</option><option>Survey Plan</option><option>Excision</option></select></div></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Owner</label><input v-model="newRecord.owner" type="text" class="input-field w-full" /></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-end"><button @click="showAddModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button><button @click="handleAddRecord" class="px-4 py-2 text-[11px] bg-[#B90B0B] text-white rounded-lg hover:bg-[#991010]">Add</button></div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showEditModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#B90B0B] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Edit Record</h3><button @click="showEditModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4"><div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Plot Number</label><input v-model="editRecord.plotNo" type="text" class="input-field w-full" /></div><div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Block</label><input v-model="editRecord.block" type="text" class="input-field w-full" /></div></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Location</label><input v-model="editRecord.location" type="text" class="input-field w-full" /></div>
            <div class="grid grid-cols-2 gap-4"><div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Size</label><input v-model="editRecord.size" type="text" class="input-field w-full" /></div><div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Status</label><select v-model="editRecord.status" class="input-field w-full"><option>Active</option><option>Pending</option><option>Flagged</option></select></div></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Owner</label><input v-model="editRecord.owner" type="text" class="input-field w-full" /></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-end"><button @click="showEditModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button><button @click="handleUpdateRecord" class="px-4 py-2 text-[11px] bg-[#B90B0B] text-white rounded-lg hover:bg-[#991010]">Save</button></div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showViewModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#B90B0B] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Record Details</h3><button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4"><div><p class="text-[11px] text-gray-500">Record ID</p><p class="text-[13px] font-medium">{{ selectedRecord?.id }}</p></div><div><p class="text-[11px] text-gray-500">Status</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': selectedRecord?.status === 'Active', 'bg-yellow-50 text-yellow-700': selectedRecord?.status === 'Pending'}">{{ selectedRecord?.status }}</span></div></div>
            <div class="grid grid-cols-2 gap-4"><div><p class="text-[11px] text-gray-500">Plot</p><p class="text-[13px]">{{ selectedRecord?.plotNo }}</p></div><div><p class="text-[11px] text-gray-500">Block</p><p class="text-[13px]">{{ selectedRecord?.block }}</p></div></div>
            <div><p class="text-[11px] text-gray-500">Location</p><p class="text-[13px]">{{ selectedRecord?.location }}</p></div>
            <div class="grid grid-cols-2 gap-4"><div><p class="text-[11px] text-gray-500">Size</p><p class="text-[13px]">{{ selectedRecord?.size }}</p></div><div><p class="text-[11px] text-gray-500">Title</p><p class="text-[13px]">{{ selectedRecord?.titleType }}</p></div></div>
            <div><p class="text-[11px] text-gray-500">Owner</p><p class="text-[13px]">{{ selectedRecord?.owner }}</p></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-end"><button @click="showViewModal = false" class="px-4 py-2 text-[11px] bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200">Close</button></div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
