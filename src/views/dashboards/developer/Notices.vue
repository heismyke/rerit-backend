<script setup lang="ts">
import { useRoleStore } from '@/stores'
import { useRouter } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import { ref, computed, onMounted } from 'vue'
import { api } from '@/services/api'

type Notice = {
  id: string
  title: string
  property: string
  amount: string
  dueDate: string
  status: string
  type: string
  response?: string
}

const { selectedRole, user, logout } = useRoleStore()
const router = useRouter()
const handleLogout = () => { logout(); router.push('/') }

const searchQuery = ref('')
const filterType = ref('all')
const currentPage = ref(1)
const itemsPerPage = ref(5)

const showViewModal = ref(false)
const showRespondModal = ref(false)
const selectedNotice = ref<Notice | null>(null)
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })
const isLoading = ref(false)
const errorMessage = ref('')
const response = ref('')

const fallbackNotices: Notice[] = [
  { id: 'NT-2024-001', title: 'Q1 2024 Assessment Notice', property: 'Commercial Complex A', amount: 'N2,500,000', dueDate: 'Mar 31, 2024', status: 'Pending', type: 'Assessment' },
  { id: 'NT-2024-002', title: 'Property Valuation Update', property: 'Residential Estate B', amount: 'N1,800,000', dueDate: 'Feb 28, 2024', status: 'Responded', type: 'Valuation' },
  { id: 'NT-2024-003', title: 'Annual Compliance Review', property: 'Office Tower D', amount: 'N4,200,000', dueDate: 'Jun 30, 2024', status: 'Resolved', type: 'Compliance' },
  { id: 'NT-2024-004', title: 'Document Verification Request', property: 'Mixed Use Development C', amount: '-', dueDate: 'Feb 15, 2024', status: 'Pending', type: 'Documentation' },
]

const notices = ref<Notice[]>(fallbackNotices)

const filteredNotices = computed(() => notices.value.filter(n => {
  const matchesSearch = n.title.toLowerCase().includes(searchQuery.value.toLowerCase()) || n.id.toLowerCase().includes(searchQuery.value.toLowerCase())
  const matchesType = filterType.value === 'all' || n.type === filterType.value
  return matchesSearch && matchesType
}))

const totalPages = computed(() => Math.ceil(filteredNotices.value.length / itemsPerPage.value))
const paginatedNotices = computed(() => { const start = (currentPage.value - 1) * itemsPerPage.value; return filteredNotices.value.slice(start, start + itemsPerPage.value) })
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }

const showToast = (message: string) => { toast.value = { show: true, message }; setTimeout(() => { toast.value.show = false }, 3000) }
const openViewModal = (n: Notice) => { selectedNotice.value = n; showViewModal.value = true }
const openRespondModal = (n: Notice) => { selectedNotice.value = n; response.value = n.response || ''; showRespondModal.value = true }
const loadNotices = async () => {
  isLoading.value = true
  errorMessage.value = ''
  try {
    notices.value = await api.getNotices<Notice[]>()
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Unable to load notices'
    notices.value = fallbackNotices
  } finally {
    isLoading.value = false
  }
}

const submitResponse = async () => {
  if (!selectedNotice.value) return
  try {
    const updated = await api.respondNotice<Notice>(selectedNotice.value.id, response.value)
    const index = notices.value.findIndex(n => n.id === updated.id)
    if (index !== -1) notices.value[index] = updated
    showRespondModal.value = false
    showToast('Response submitted successfully')
  } catch (error) {
    showToast(error instanceof Error ? error.message : 'Unable to submit response')
  }
}

onMounted(loadNotices)
</script>

<template>
  <div class="min-h-screen bg-[#f5f6fa] flex">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />
    <div class="flex-1 flex flex-col">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6 shrink-0">
        <div class="flex items-center gap-4"><span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span><span class="text-[#d1d5db]">/</span><span class="text-[#1f2937] text-sm font-medium">Notices</span></div>
        <div class="flex items-center gap-4"><span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span><button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button></div>
      </header>
      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between"><h2 class="text-[13px] font-semibold text-[#1f2937]">View Notices</h2></div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4">
            <input v-model="searchQuery" type="text" placeholder="Search notices..." class="input-field flex-1" />
            <select v-model="filterType" class="input-field w-48"><option value="all">All Types</option><option value="Assessment">Assessment</option><option value="Valuation">Valuation</option><option value="Compliance">Compliance</option><option value="Documentation">Documentation</option></select>
          </div>
          <div v-if="isLoading" class="px-6 py-3 text-[12px] text-[#6b7280] border-b border-[#f3f4f6]">Loading notices...</div>
          <div v-else-if="errorMessage" class="px-6 py-3 text-[12px] text-yellow-700 bg-yellow-50 border-b border-yellow-100">Showing demo notices because the backend did not respond.</div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead><tr><th class="table-header">Notice ID</th><th class="table-header">Title</th><th class="table-header">Property</th><th class="table-header">Amount</th><th class="table-header">Due Date</th><th class="table-header">Type</th><th class="table-header">Status</th><th class="table-header">Actions</th></tr></thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="notice in paginatedNotices" :key="notice.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ notice.id }}</td><td class="table-cell">{{ notice.title }}</td><td class="table-cell text-[#6b7280]">{{ notice.property }}</td><td class="table-cell">{{ notice.amount }}</td><td class="table-cell text-[#9ca3af]">{{ notice.dueDate }}</td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] bg-[#f3f4f6] text-[#6b7280] rounded">{{ notice.type }}</span></td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-yellow-50 text-yellow-700': notice.status === 'Pending', 'bg-blue-50 text-blue-700': notice.status === 'Responded', 'bg-green-50 text-green-700': notice.status === 'Resolved'}">{{ notice.status }}</span></td>
                  <td class="table-cell"><div class="flex gap-2"><button @click="openViewModal(notice)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">View</button><button v-if="notice.status === 'Pending'" @click="openRespondModal(notice)" class="px-3 py-1 text-[11px] bg-[#2D5A27] text-white rounded hover:bg-[#6a0707]">Respond</button></div></td>
                </tr>
                <tr v-if="!isLoading && paginatedNotices.length === 0">
                  <td colspan="8" class="px-6 py-8 text-center text-[12px] text-[#6b7280]">No notices found.</td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="px-6 py-4 flex items-center justify-between border-t border-[#f3f4f6]">
            <p class="text-[11px] text-[#6b7280]">Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, filteredNotices.length) }} of {{ filteredNotices.length }} entries</p>
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
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Notice Details</h3><button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div><p class="text-[11px] text-gray-500">Notice ID</p><p class="text-[13px] font-medium">{{ selectedNotice?.id }}</p></div>
              <div><p class="text-[11px] text-gray-500">Status</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-yellow-50 text-yellow-700': selectedNotice?.status === 'Pending', 'bg-blue-50 text-blue-700': selectedNotice?.status === 'Responded'}">{{ selectedNotice?.status }}</span></div>
              <div><p class="text-[11px] text-gray-500">Type</p><p class="text-[13px]">{{ selectedNotice?.type }}</p></div>
              <div><p class="text-[11px] text-gray-500">Due Date</p><p class="text-[13px]">{{ selectedNotice?.dueDate }}</p></div>
            </div>
            <div><p class="text-[11px] text-gray-500">Title</p><p class="text-[13px] font-medium">{{ selectedNotice?.title }}</p></div>
            <div><p class="text-[11px] text-gray-500">Property</p><p class="text-[13px]">{{ selectedNotice?.property }}</p></div>
            <div><p class="text-[11px] text-gray-500">Amount</p><p class="text-[13px] font-medium text-green-700">{{ selectedNotice?.amount }}</p></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-end"><button @click="showViewModal = false" class="px-4 py-2 text-[11px] bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200">Close</button></div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showRespondModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Respond to Notice</h3><button @click="showRespondModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div class="bg-[#f9fafb] rounded-lg p-4">
              <p class="text-[11px] text-gray-500">Notice</p>
              <p class="text-[13px] font-medium">{{ selectedNotice?.title }}</p>
              <p class="text-[11px] text-gray-500 mt-2">Property: {{ selectedNotice?.property }}</p>
            </div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Your Response</label><textarea v-model="response" rows="4" placeholder="Enter your response..." class="input-field w-full"></textarea></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-end"><button @click="showRespondModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button><button @click="submitResponse" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Submit Response</button></div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
