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
const selectedSurvey = ref<any>(null)
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })

const newSurvey = ref({ property: '', surveyor: '', type: 'Land Survey', status: 'Pending' })
const editSurvey = ref({ property: '', surveyor: '', type: 'Land Survey', status: 'Pending' })

const surveys = ref([
  { id: 'SURV-001', property: 'Plot 15, Ikoyi', surveyor: 'Aisha Ibrahim', type: 'Land Survey', status: 'Approved', date: '2024-01-10' },
  { id: 'SURV-002', property: 'Block 3, Lekki', surveyor: 'Babatunde Olatunji', type: 'Boundary Survey', status: 'Under Review', date: '2024-01-12' },
  { id: 'SURV-003', property: 'Plot 88, VI', surveyor: 'Aisha Ibrahim', type: 'Topographic Survey', status: 'Approved', date: '2024-01-08' },
])

const filteredSurveys = computed(() => surveys.value.filter(s => s.property.toLowerCase().includes(searchQuery.value.toLowerCase()) || s.id.toLowerCase().includes(searchQuery.value.toLowerCase())))
const totalPages = computed(() => Math.ceil(filteredSurveys.value.length / itemsPerPage.value))
const paginatedSurveys = computed(() => { const start = (currentPage.value - 1) * itemsPerPage.value; return filteredSurveys.value.slice(start, start + itemsPerPage.value) })
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }

const showToast = (message: string) => { toast.value = { show: true, message }; setTimeout(() => { toast.value.show = false }, 3000) }
const openAddModal = () => { newSurvey.value = { property: '', surveyor: '', type: 'Land Survey', status: 'Pending' }; showAddModal.value = true }
const openEditModal = (s: any) => { selectedSurvey.value = s; editSurvey.value = { ...s }; showEditModal.value = true }
const openViewModal = (s: any) => { selectedSurvey.value = s; showViewModal.value = true }

const handleAddSurvey = () => {
  const newId = 'SURV-' + String(surveys.value.length + 1).padStart(3, '0')
  surveys.value.unshift({ id: newId, date: new Date().toISOString().split('T')[0], ...newSurvey.value })
  showAddModal.value = false; showToast('Survey created successfully')
}

const handleUpdateSurvey = () => {
  const index = surveys.value.findIndex(s => s.id === selectedSurvey.value.id)
  if (index !== -1) { surveys.value[index] = { ...surveys.value[index], ...editSurvey.value }; showToast('Survey updated') }
  showEditModal.value = false
}

const handleDeleteSurvey = () => {
  surveys.value = surveys.value.filter(s => s.id !== selectedSurvey.value.id)
  showEditModal.value = false; showToast('Survey deleted')
}
</script>

<template>
  <div class="min-h-screen flex bg-[#f5f7fa]">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />
    <div class="flex-1">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6">
        <div class="flex items-center gap-4"><span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span><span class="text-[#d1d5db]">/</span><span class="text-[#1f2937] text-sm font-medium">Surveys</span></div>
        <div class="flex items-center gap-4"><span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span><button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button></div>
      </header>
      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between"><h2 class="text-[13px] font-semibold text-[#1f2937]">All Surveys</h2><button @click="openAddModal" class="btn-primary text-[11px]">Create Survey</button></div>
          <div class="p-4 border-b border-[#e5e7eb]"><input v-model="searchQuery" type="text" placeholder="Search..." class="input-field max-w-md" /></div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead><tr><th class="table-header">ID</th><th class="table-header">Property</th><th class="table-header">Surveyor</th><th class="table-header">Type</th><th class="table-header">Status</th><th class="table-header">Date</th><th class="table-header">Actions</th></tr></thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="s in paginatedSurveys" :key="s.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ s.id }}</td><td class="table-cell text-[#6b7280]">{{ s.property }}</td><td class="table-cell">{{ s.surveyor }}</td><td class="table-cell">{{ s.type }}</td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': s.status === 'Approved', 'bg-yellow-50 text-yellow-700': s.status === 'Under Review', 'bg-gray-100 text-gray-600': s.status === 'Pending'}">{{ s.status }}</span></td>
                  <td class="table-cell text-[#9ca3af]">{{ s.date }}</td>
                  <td class="table-cell"><div class="flex gap-2"><button @click="openViewModal(s)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">View</button><button @click="openEditModal(s)" class="px-3 py-1 text-[11px] bg-[#2D5A27] text-white rounded hover:bg-[#1e3d1a]">Edit</button></div></td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="px-6 py-4 flex items-center justify-between border-t border-[#f3f4f6]">
            <p class="text-[11px] text-[#6b7280]">Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, filteredSurveys.length) }} of {{ filteredSurveys.length }} entries</p>
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
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Create Survey</h3><button @click="showAddModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Property</label><input v-model="newSurvey.property" type="text" placeholder="Enter property address" class="input-field w-full" /></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Surveyor</label><input v-model="newSurvey.surveyor" type="text" placeholder="Enter surveyor name" class="input-field w-full" /></div>
            <div class="grid grid-cols-2 gap-4">
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Type</label><select v-model="newSurvey.type" class="input-field w-full"><option>Land Survey</option><option>Boundary Survey</option><option>Topographic Survey</option></select></div>
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Status</label><select v-model="newSurvey.status" class="input-field w-full"><option>Pending</option><option>Under Review</option><option>Approved</option></select></div>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-end"><button @click="showAddModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button><button @click="handleAddSurvey" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Create</button></div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showEditModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Edit Survey</h3><button @click="showEditModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Property</label><input v-model="editSurvey.property" type="text" class="input-field w-full" /></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Surveyor</label><input v-model="editSurvey.surveyor" type="text" class="input-field w-full" /></div>
            <div class="grid grid-cols-2 gap-4">
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Type</label><select v-model="editSurvey.type" class="input-field w-full"><option>Land Survey</option><option>Boundary Survey</option><option>Topographic Survey</option></select></div>
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Status</label><select v-model="editSurvey.status" class="input-field w-full"><option>Pending</option><option>Under Review</option><option>Approved</option></select></div>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-between"><button @click="handleDeleteSurvey" class="px-4 py-2 text-[11px] text-red-600 border border-red-200 rounded-lg hover:bg-red-50">Delete</button><div class="flex gap-3"><button @click="showEditModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button><button @click="handleUpdateSurvey" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Save</button></div></div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showViewModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Survey Details</h3><button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div><p class="text-[11px] text-gray-500">Survey ID</p><p class="text-[13px] font-medium">{{ selectedSurvey?.id }}</p></div>
              <div><p class="text-[11px] text-gray-500">Status</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': selectedSurvey?.status === 'Approved', 'bg-yellow-50 text-yellow-700': selectedSurvey?.status === 'Under Review'}">{{ selectedSurvey?.status }}</span></div>
              <div><p class="text-[11px] text-gray-500">Surveyor</p><p class="text-[13px] font-medium">{{ selectedSurvey?.surveyor }}</p></div>
              <div><p class="text-[11px] text-gray-500">Type</p><p class="text-[13px]">{{ selectedSurvey?.type }}</p></div>
            </div>
            <div><p class="text-[11px] text-gray-500">Property</p><p class="text-[13px]">{{ selectedSurvey?.property }}</p></div>
            <div><p class="text-[11px] text-gray-500">Created</p><p class="text-[13px]">{{ selectedSurvey?.date }}</p></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-end"><button @click="showViewModal = false" class="px-4 py-2 text-[11px] bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200">Close</button></div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
