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

const showAddModal = ref(false)
const showEditModal = ref(false)
const showViewModal = ref(false)
const selectedAudit = ref<any>(null)
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })

const newAudit = ref({ property: '', auditor: '', status: 'Scheduled', result: '-' })
const editAudit = ref({ property: '', auditor: '', status: 'Scheduled', result: '-' })

const audits = ref([
  { id: 'AUD-001', property: 'Plot 42, Victoria Island', auditor: 'John Smith', status: 'Completed', result: 'Verified', date: '2024-01-15' },
  { id: 'AUD-002', property: 'Block 7, Lekki Phase 2', auditor: 'Sarah Johnson', status: 'In Progress', result: '-', date: '2024-01-14' },
  { id: 'AUD-003', property: '15 Admiralty Way, Lekki', auditor: 'Michael Brown', status: 'Completed', result: 'Flagged', date: '2024-01-13' },
])

const filteredAudits = computed(() => audits.value.filter(a => a.property.toLowerCase().includes(searchQuery.value.toLowerCase()) || a.id.toLowerCase().includes(searchQuery.value.toLowerCase())))
const totalPages = computed(() => Math.ceil(filteredAudits.value.length / itemsPerPage.value))
const paginatedAudits = computed(() => { const start = (currentPage.value - 1) * itemsPerPage.value; return filteredAudits.value.slice(start, start + itemsPerPage.value) })
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }

const showToast = (message: string) => { toast.value = { show: true, message }; setTimeout(() => { toast.value.show = false }, 3000) }
const openAddModal = () => { newAudit.value = { property: '', auditor: '', status: 'Scheduled', result: '-' }; showAddModal.value = true }
const openEditModal = (a: any) => { selectedAudit.value = a; editAudit.value = { ...a }; showEditModal.value = true }
const openViewModal = (a: any) => { selectedAudit.value = a; showViewModal.value = true }

const handleAddAudit = () => {
  const newId = 'AUD-' + String(audits.value.length + 1).padStart(3, '0')
  audits.value.unshift({ id: newId, date: new Date().toISOString().split('T')[0], ...newAudit.value })
  showAddModal.value = false; showToast('Audit scheduled successfully')
}

const handleUpdateAudit = () => {
  const index = audits.value.findIndex(a => a.id === selectedAudit.value.id)
  if (index !== -1) { audits.value[index] = { ...audits.value[index], ...editAudit.value }; showToast('Audit updated') }
  showEditModal.value = false
}

const handleDeleteAudit = () => {
  audits.value = audits.value.filter(a => a.id !== selectedAudit.value.id)
  showEditModal.value = false; showToast('Audit deleted')
}
</script>

<template>
  <div class="min-h-screen flex bg-[#f5f7fa]">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />
    <div class="flex-1">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6">
        <div class="flex items-center gap-4"><span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span><span class="text-[#d1d5db]">/</span><span class="text-[#1f2937] text-sm font-medium">Audits</span></div>
        <div class="flex items-center gap-4"><span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span><button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button></div>
      </header>
      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between"><h2 class="text-[13px] font-semibold text-[#1f2937]">All Audits</h2><button @click="openAddModal" class="btn-primary text-[11px]">Create Audit</button></div>
          <div class="p-4 border-b border-[#e5e7eb]"><input v-model="searchQuery" type="text" placeholder="Search..." class="input-field max-w-md" /></div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead><tr><th class="table-header">ID</th><th class="table-header">Property</th><th class="table-header">Auditor</th><th class="table-header">Status</th><th class="table-header">Result</th><th class="table-header">Date</th><th class="table-header">Actions</th></tr></thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="a in paginatedAudits" :key="a.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ a.id }}</td><td class="table-cell text-[#6b7280]">{{ a.property }}</td><td class="table-cell">{{ a.auditor }}</td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-blue-50 text-blue-700': a.status === 'In Progress', 'bg-green-50 text-green-700': a.status === 'Completed', 'bg-yellow-50 text-yellow-700': a.status === 'Scheduled'}">{{ a.status }}</span></td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': a.result === 'Verified', 'bg-red-50 text-red-700': a.result === 'Flagged', 'bg-gray-100 text-gray-600': a.result === '-'}">{{ a.result }}</span></td>
                  <td class="table-cell text-[#9ca3af]">{{ a.date }}</td>
                  <td class="table-cell"><div class="flex gap-2"><button @click="openViewModal(a)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">View</button><button @click="openEditModal(a)" class="px-3 py-1 text-[11px] bg-[#2D5A27] text-white rounded hover:bg-[#1e3d1a]">Edit</button></div></td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="px-6 py-4 flex items-center justify-between border-t border-[#f3f4f6]">
            <p class="text-[11px] text-[#6b7280]">Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, filteredAudits.length) }} of {{ filteredAudits.length }} entries</p>
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
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Schedule Audit</h3><button @click="showAddModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Property</label><input v-model="newAudit.property" type="text" placeholder="Enter property address" class="input-field w-full" /></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Auditor</label><input v-model="newAudit.auditor" type="text" placeholder="Enter auditor name" class="input-field w-full" /></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Status</label><select v-model="newAudit.status" class="input-field w-full"><option>Scheduled</option><option>In Progress</option><option>Completed</option></select></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-end"><button @click="showAddModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button><button @click="handleAddAudit" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Schedule</button></div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showEditModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Edit Audit</h3><button @click="showEditModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Property</label><input v-model="editAudit.property" type="text" class="input-field w-full" /></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Auditor</label><input v-model="editAudit.auditor" type="text" class="input-field w-full" /></div>
            <div class="grid grid-cols-2 gap-4">
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Status</label><select v-model="editAudit.status" class="input-field w-full"><option>Scheduled</option><option>In Progress</option><option>Completed</option></select></div>
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Result</label><select v-model="editAudit.result" class="input-field w-full"><option>-</option><option>Verified</option><option>Flagged</option></select></div>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-between"><button @click="handleDeleteAudit" class="px-4 py-2 text-[11px] text-red-600 border border-red-200 rounded-lg hover:bg-red-50">Delete</button><div class="flex gap-3"><button @click="showEditModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button><button @click="handleUpdateAudit" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Save</button></div></div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showViewModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Audit Details</h3><button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div><p class="text-[11px] text-gray-500">Audit ID</p><p class="text-[13px] font-medium">{{ selectedAudit?.id }}</p></div>
              <div><p class="text-[11px] text-gray-500">Status</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-blue-50 text-blue-700': selectedAudit?.status === 'In Progress', 'bg-green-50 text-green-700': selectedAudit?.status === 'Completed'}">{{ selectedAudit?.status }}</span></div>
              <div><p class="text-[11px] text-gray-500">Auditor</p><p class="text-[13px] font-medium">{{ selectedAudit?.auditor }}</p></div>
              <div><p class="text-[11px] text-gray-500">Result</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': selectedAudit?.result === 'Verified', 'bg-red-50 text-red-700': selectedAudit?.result === 'Flagged', 'bg-gray-100 text-gray-600': selectedAudit?.result === '-'}">{{ selectedAudit?.result }}</span></div>
            </div>
            <div><p class="text-[11px] text-gray-500">Property</p><p class="text-[13px]">{{ selectedAudit?.property }}</p></div>
            <div><p class="text-[11px] text-gray-500">Scheduled Date</p><p class="text-[13px]">{{ selectedAudit?.date }}</p></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-end"><button @click="showViewModal = false" class="px-4 py-2 text-[11px] bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200">Close</button></div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
