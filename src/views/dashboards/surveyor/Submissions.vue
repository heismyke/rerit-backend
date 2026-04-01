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
const selectedSubmission = ref<any>(null)
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })

const newSubmission = ref({ property: '', type: 'Land Survey' })

const submissions = ref([
  { id: 'SUB-001', surveyId: 'SURV-001', type: 'Land Survey', property: 'Plot 15, Ikoyi', status: 'Approved', date: '2024-01-10', verifiedBy: 'Admin' },
  { id: 'SUB-002', surveyId: 'SURV-002', type: 'Boundary Survey', property: 'Block 3, Lekki', status: 'Under Review', date: '2024-01-12', verifiedBy: '-' },
  { id: 'SUB-003', surveyId: 'SURV-003', type: 'Topographic Survey', property: 'Plot 88, VI', status: 'Approved', date: '2024-01-08', verifiedBy: 'Admin' },
  { id: 'SUB-004', surveyId: 'SURV-004', type: 'Subdivision Survey', property: 'Block 7, Epe', status: 'Pending', date: '2024-01-14', verifiedBy: '-' },
])

const filteredSubmissions = computed(() => submissions.value.filter(s => {
  const matchesSearch = s.property.toLowerCase().includes(searchQuery.value.toLowerCase()) || s.id.toLowerCase().includes(searchQuery.value.toLowerCase())
  const matchesStatus = filterStatus.value === 'all' || s.status === filterStatus.value
  return matchesSearch && matchesStatus
}))

const totalPages = computed(() => Math.ceil(filteredSubmissions.value.length / itemsPerPage.value))
const paginatedSubmissions = computed(() => { const start = (currentPage.value - 1) * itemsPerPage.value; return filteredSubmissions.value.slice(start, start + itemsPerPage.value) })
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }

const showToast = (message: string) => { toast.value = { show: true, message }; setTimeout(() => { toast.value.show = false }, 3000) }
const openViewModal = (s: any) => { selectedSubmission.value = s; showViewModal.value = true }
const handleAddSubmission = () => {
  const newId = 'SUB-' + String(submissions.value.length + 1).padStart(3, '0')
  submissions.value.unshift({ id: newId, surveyId: 'SURV-' + newId.split('-')[1], ...newSubmission.value, status: 'Pending', date: new Date().toISOString().split('T')[0], verifiedBy: '-' })
  showAddModal.value = false; newSubmission.value = { property: '', type: 'Land Survey' }; showToast('Submission created successfully')
}
</script>

<template>
  <div class="min-h-screen flex bg-[#f5f7fa]">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />
    <div class="flex-1 flex flex-col">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6 shrink-0">
        <div class="flex items-center gap-4"><span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span><span class="text-[#d1d5db]">/</span><span class="text-[#1f2937] text-sm font-medium">Submissions</span></div>
        <div class="flex items-center gap-4"><span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span><button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button></div>
      </header>
      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between"><h2 class="text-[13px] font-semibold text-[#1f2937]">Submissions</h2><button @click="showAddModal = true" class="btn-primary text-[11px]">New Submission</button></div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4">
            <input v-model="searchQuery" type="text" placeholder="Search submissions..." class="input-field max-w-md" />
            <select v-model="filterStatus" class="input-field w-48"><option value="all">All Status</option><option value="Approved">Approved</option><option value="Under Review">Under Review</option><option value="Pending">Pending</option></select>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead><tr><th class="table-header">Submission ID</th><th class="table-header">Survey ID</th><th class="table-header">Type</th><th class="table-header">Property</th><th class="table-header">Status</th><th class="table-header">Date</th><th class="table-header">Verified By</th><th class="table-header">Actions</th></tr></thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="sub in paginatedSubmissions" :key="sub.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ sub.id }}</td><td class="table-cell">{{ sub.surveyId }}</td><td class="table-cell">{{ sub.type }}</td><td class="table-cell text-[#6b7280]">{{ sub.property }}</td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': sub.status === 'Approved', 'bg-yellow-50 text-yellow-700': sub.status === 'Under Review', 'bg-gray-100 text-gray-600': sub.status === 'Pending'}">{{ sub.status }}</span></td>
                  <td class="table-cell text-[#9ca3af]">{{ sub.date }}</td><td class="table-cell text-[#9ca3af]">{{ sub.verifiedBy }}</td>
                  <td class="table-cell"><button @click="openViewModal(sub)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">View</button></td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="px-6 py-4 flex items-center justify-between border-t border-[#f3f4f6]">
            <p class="text-[11px] text-[#6b7280]">Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, filteredSubmissions.length) }} of {{ filteredSubmissions.length }} entries</p>
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
          <div class="bg-[#B90B0B] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">New Submission</h3><button @click="showAddModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Property Location</label><input v-model="newSubmission.property" type="text" placeholder="Enter property address" class="input-field w-full" /></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Submission Type</label><select v-model="newSubmission.type" class="input-field w-full"><option>Land Survey</option><option>Boundary Survey</option><option>Topographic Survey</option><option>Subdivision Survey</option></select></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-end"><button @click="showAddModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button><button @click="handleAddSubmission" class="px-4 py-2 text-[11px] bg-[#B90B0B] text-white rounded-lg hover:bg-[#991010]">Submit</button></div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showViewModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#B90B0B] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Submission Details</h3><button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div><p class="text-[11px] text-gray-500">Submission ID</p><p class="text-[13px] font-medium">{{ selectedSubmission?.id }}</p></div>
              <div><p class="text-[11px] text-gray-500">Status</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': selectedSubmission?.status === 'Approved', 'bg-yellow-50 text-yellow-700': selectedSubmission?.status === 'Under Review'}">{{ selectedSubmission?.status }}</span></div>
              <div><p class="text-[11px] text-gray-500">Survey ID</p><p class="text-[13px]">{{ selectedSubmission?.surveyId }}</p></div>
              <div><p class="text-[11px] text-gray-500">Date</p><p class="text-[13px]">{{ selectedSubmission?.date }}</p></div>
            </div>
            <div><p class="text-[11px] text-gray-500">Type</p><p class="text-[13px]">{{ selectedSubmission?.type }}</p></div>
            <div><p class="text-[11px] text-gray-500">Property</p><p class="text-[13px]">{{ selectedSubmission?.property }}</p></div>
            <div><p class="text-[11px] text-gray-500">Verified By</p><p class="text-[13px]">{{ selectedSubmission?.verifiedBy }}</p></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-end"><button @click="showViewModal = false" class="px-4 py-2 text-[11px] bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200">Close</button></div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
