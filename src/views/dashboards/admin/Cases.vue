<script setup lang="ts">
import { useRoleStore } from '@/stores'
import { useRouter } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import { ref, computed } from 'vue'
import { getRiskColor, getRiskScoreColor } from '@/utils/riskScoring'

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
const selectedCase = ref<any>(null)
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })

const newCase = ref({ title: '', type: 'Tax Evasion', priority: 'Medium', status: 'Open', assigned: '', riskScore: 30, riskLevel: 'Low' })
const editCase = ref({ title: '', type: 'Tax Evasion', priority: 'Medium', status: 'Open', assigned: '' })

const cases = ref([
  { id: 'CASE-001', title: 'Undeclared Property', type: 'Tax Evasion', priority: 'High', status: 'Open', assigned: 'John Smith', date: '2024-01-15', riskScore: 65, riskLevel: 'High' },
  { id: 'CASE-002', title: 'Value Discrepancy', type: 'Valuation', priority: 'Medium', status: 'In Progress', assigned: 'Sarah Johnson', date: '2024-01-14', riskScore: 45, riskLevel: 'Medium' },
  { id: 'CASE-003', title: 'Document Forgery', type: 'Fraud', priority: 'Critical', status: 'Open', assigned: 'Michael Brown', date: '2024-01-13', riskScore: 88, riskLevel: 'Critical' },
  { id: 'CASE-004', title: 'Late Payment', type: 'Compliance', priority: 'Low', status: 'Resolved', assigned: 'John Smith', date: '2024-01-12', riskScore: 20, riskLevel: 'Low' },
])

const getCaseRiskLevel = (priority: string, status: string) => {
  if (priority === 'Critical' && status === 'Open') return { score: 88, level: 'Critical' }
  if (priority === 'High') return { score: 65, level: 'High' }
  if (priority === 'Medium') return { score: 45, level: 'Medium' }
  return { score: 20, level: 'Low' }
}

const filteredCases = computed(() => cases.value.filter(c => c.title.toLowerCase().includes(searchQuery.value.toLowerCase()) || c.id.toLowerCase().includes(searchQuery.value.toLowerCase())))
const totalPages = computed(() => Math.ceil(filteredCases.value.length / itemsPerPage.value))
const paginatedCases = computed(() => { const start = (currentPage.value - 1) * itemsPerPage.value; return filteredCases.value.slice(start, start + itemsPerPage.value) })
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }

const showToast = (message: string) => { toast.value = { show: true, message }; setTimeout(() => { toast.value.show = false }, 3000) }
const openAddModal = () => { newCase.value = { title: '', type: 'Tax Evasion', priority: 'Medium', status: 'Open', assigned: '', riskScore: 30, riskLevel: 'Low' }; showAddModal.value = true }
const openEditModal = (c: any) => { selectedCase.value = c; editCase.value = { ...c }; showEditModal.value = true }
const openViewModal = (c: any) => { selectedCase.value = c; showViewModal.value = true }

const handleAddCase = () => {
  const newId = 'CASE-' + String(cases.value.length + 1).padStart(3, '0')
  cases.value.unshift({ id: newId, date: new Date().toISOString().split('T')[0], ...newCase.value })
  showAddModal.value = false; showToast('Case created successfully')
}

const handleUpdateCase = () => {
  const index = cases.value.findIndex(c => c.id === selectedCase.value.id)
  if (index !== -1) { cases.value[index] = { ...cases.value[index], ...editCase.value }; showToast('Case updated') }
  showEditModal.value = false
}

const handleDeleteCase = () => {
  cases.value = cases.value.filter(c => c.id !== selectedCase.value.id)
  showEditModal.value = false; showToast('Case deleted')
}
</script>

<template>
  <div class="min-h-screen flex bg-[#f5f7fa]">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />
    <div class="flex-1">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6">
        <div class="flex items-center gap-4"><span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span><span class="text-[#d1d5db]">/</span><span class="text-[#1f2937] text-sm font-medium">Cases</span></div>
        <div class="flex items-center gap-4"><span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span><button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button></div>
      </header>
      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between"><h2 class="text-[13px] font-semibold text-[#1f2937]">All Cases</h2><button @click="openAddModal" class="btn-primary text-[11px]">Create Case</button></div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4">
            <input v-model="searchQuery" type="text" placeholder="Search..." class="input-field max-w-md" />
            <select v-model="filterStatus" class="input-field w-48"><option value="all">All Status</option><option value="Open">Open</option><option value="In Progress">In Progress</option><option value="Resolved">Resolved</option></select>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead><tr><th class="table-header">ID</th><th class="table-header">Title</th><th class="table-header">Type</th><th class="table-header">Priority</th><th class="table-header">Status</th><th class="table-header">Risk Score</th><th class="table-header">Assigned</th><th class="table-header">Date</th><th class="table-header">Actions</th></tr></thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="c in paginatedCases" :key="c.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ c.id }}</td><td class="table-cell">{{ c.title }}</td><td class="table-cell text-[#6b7280]">{{ c.type }}</td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-red-100 text-red-700': c.priority === 'Critical' || c.priority === 'High', 'bg-yellow-100 text-yellow-700': c.priority === 'Medium', 'bg-gray-100 text-gray-600': c.priority === 'Low'}">{{ c.priority }}</span></td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-blue-50 text-blue-700': c.status === 'In Progress', 'bg-yellow-50 text-yellow-700': c.status === 'Open', 'bg-green-50 text-green-700': c.status === 'Resolved'}">{{ c.status }}</span></td>
                  <td class="table-cell">
                    <div class="flex items-center gap-2">
                      <span class="font-semibold" :class="getRiskScoreColor(getCaseRiskLevel(c.priority, c.status).score)">{{ getCaseRiskLevel(c.priority, c.status).score }}</span>
                      <span class="px-1.5 py-0.5 text-[10px] font-medium rounded" :class="getRiskColor(getCaseRiskLevel(c.priority, c.status).level)">{{ getCaseRiskLevel(c.priority, c.status).level }}</span>
                    </div>
                  </td>
                  <td class="table-cell text-[#6b7280]">{{ c.assigned }}</td><td class="table-cell text-[#9ca3af]">{{ c.date }}</td>
                  <td class="table-cell"><div class="flex gap-2"><button @click="openViewModal(c)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">View</button><button @click="openEditModal(c)" class="px-3 py-1 text-[11px] bg-[#2D5A27] text-white rounded hover:bg-[#1e3d1a]">Edit</button></div></td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="px-6 py-4 flex items-center justify-between border-t border-[#f3f4f6]">
            <p class="text-[11px] text-[#6b7280]">Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, filteredCases.length) }} of {{ filteredCases.length }} entries</p>
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
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Create Case</h3><button @click="showAddModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Title</label><input v-model="newCase.title" type="text" placeholder="Enter case title" class="input-field w-full" /></div>
            <div class="grid grid-cols-2 gap-4">
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Type</label><select v-model="newCase.type" class="input-field w-full"><option>Tax Evasion</option><option>Valuation</option><option>Fraud</option><option>Compliance</option></select></div>
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Priority</label><select v-model="newCase.priority" class="input-field w-full"><option>Low</option><option>Medium</option><option>High</option><option>Critical</option></select></div>
            </div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Assigned To</label><input v-model="newCase.assigned" type="text" placeholder="Enter name" class="input-field w-full" /></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-end"><button @click="showAddModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button><button @click="handleAddCase" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Create</button></div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showEditModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Edit Case</h3><button @click="showEditModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Title</label><input v-model="editCase.title" type="text" class="input-field w-full" /></div>
            <div class="grid grid-cols-2 gap-4">
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Type</label><select v-model="editCase.type" class="input-field w-full"><option>Tax Evasion</option><option>Valuation</option><option>Fraud</option><option>Compliance</option></select></div>
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Priority</label><select v-model="editCase.priority" class="input-field w-full"><option>Low</option><option>Medium</option><option>High</option><option>Critical</option></select></div>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Status</label><select v-model="editCase.status" class="input-field w-full"><option>Open</option><option>In Progress</option><option>Resolved</option></select></div>
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Assigned To</label><input v-model="editCase.assigned" type="text" class="input-field w-full" /></div>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-between"><button @click="handleDeleteCase" class="px-4 py-2 text-[11px] text-red-600 border border-red-200 rounded-lg hover:bg-red-50">Delete</button><div class="flex gap-3"><button @click="showEditModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button><button @click="handleUpdateCase" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Save</button></div></div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showViewModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Case Details</h3><button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div><p class="text-[11px] text-gray-500">Case ID</p><p class="text-[13px] font-medium">{{ selectedCase?.id }}</p></div>
              <div><p class="text-[11px] text-gray-500">Status</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-blue-50 text-blue-700': selectedCase?.status === 'In Progress', 'bg-yellow-50 text-yellow-700': selectedCase?.status === 'Open', 'bg-green-50 text-green-700': selectedCase?.status === 'Resolved'}">{{ selectedCase?.status }}</span></div>
              <div><p class="text-[11px] text-gray-500">Title</p><p class="text-[13px] font-medium">{{ selectedCase?.title }}</p></div>
              <div><p class="text-[11px] text-gray-500">Priority</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-red-100 text-red-700': selectedCase?.priority === 'Critical' || selectedCase?.priority === 'High', 'bg-yellow-100 text-yellow-700': selectedCase?.priority === 'Medium'}">{{ selectedCase?.priority }}</span></div>
            </div>
            <div><p class="text-[11px] text-gray-500">Type</p><p class="text-[13px]">{{ selectedCase?.type }}</p></div>
            <div><p class="text-[11px] text-gray-500">Assigned To</p><p class="text-[13px]">{{ selectedCase?.assigned }}</p></div>
            <div><p class="text-[11px] text-gray-500">Created</p><p class="text-[13px]">{{ selectedCase?.date }}</p></div>
            <div class="border-t border-gray-100 pt-4">
              <p class="text-[11px] text-gray-500 mb-2">Risk Assessment</p>
              <div class="flex items-center gap-2">
                <span class="text-sm font-bold" :class="getRiskScoreColor(getCaseRiskLevel(selectedCase?.priority, selectedCase?.status).score)">{{ getCaseRiskLevel(selectedCase?.priority, selectedCase?.status).score }}</span>
                <span class="px-2 py-0.5 text-[11px] font-medium rounded" :class="getRiskColor(getCaseRiskLevel(selectedCase?.priority, selectedCase?.status).level)">{{ getCaseRiskLevel(selectedCase?.priority, selectedCase?.status).level }}</span>
              </div>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-end"><button @click="showViewModal = false" class="px-4 py-2 text-[11px] bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200">Close</button></div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
