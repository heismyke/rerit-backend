<script setup lang="ts">
import { useRoleStore } from '@/stores'
import { useRouter } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import { ref, computed } from 'vue'
import { useAuditorCasesStore, type AuditCase } from '@/stores/auditorCasesStore'

const { selectedRole, user, logout } = useRoleStore()
const router = useRouter()
const handleLogout = () => { logout(); router.push('/') }

const searchQuery = ref('')
const filterPriority = ref('all')
const currentPage = ref(1)
const itemsPerPage = ref(5)

const showAddModal = ref(false)
const showViewModal = ref(false)
const showEditModal = ref(false)
const showResultModal = ref(false)
const selectedCase = ref<any>(null)
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })

type CaseDraft = {
  property: string
  owner: string
  priority: AuditCase['priority']
  status: AuditCase['status']
  due: string
}
const newCase = ref<CaseDraft>({ property: '', owner: '', priority: 'Medium', status: 'Pending', due: '' })
const editCase = ref<CaseDraft>({ property: '', owner: '', priority: 'Medium', status: 'Pending', due: '' })

const { auditCases, moveAuditToSuccessful, moveAuditToFlagged } = useAuditorCasesStore()

const filteredCases = computed(() => auditCases.value.filter(c => {
  const matchesSearch = c.id.toLowerCase().includes(searchQuery.value.toLowerCase()) || c.property.toLowerCase().includes(searchQuery.value.toLowerCase()) || c.owner.toLowerCase().includes(searchQuery.value.toLowerCase())
  const matchesPriority = filterPriority.value === 'all' || c.priority === filterPriority.value
  return matchesSearch && matchesPriority
}))

const totalPages = computed(() => Math.ceil(filteredCases.value.length / itemsPerPage.value))
const paginatedCases = computed(() => { const start = (currentPage.value - 1) * itemsPerPage.value; return filteredCases.value.slice(start, start + itemsPerPage.value) })
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }

const showToast = (message: string) => { toast.value = { show: true, message }; setTimeout(() => { toast.value.show = false }, 3000) }
const openViewModal = (c: any) => { selectedCase.value = c; showViewModal.value = true }
const openEditModal = (c: any) => { selectedCase.value = c; editCase.value = { ...c }; showEditModal.value = true }
const openResultModal = (c: any) => {
  selectedCase.value = c
  resultForm.value = {
    status: c.resultStatus || 'Compliant',
    message: c.resultNotes || '',
  }
  showResultModal.value = true
}

const resultForm = ref<{ status: NonNullable<AuditCase['resultStatus']>; message: string }>({
  status: 'Compliant',
  message: '',
})

const handleAddCase = () => {
  const newId = 'AUD-2024-' + String(auditCases.value.length + 1).padStart(3, '0')
  auditCases.value.unshift({
    id: newId,
    auditor: 'Current User',
    started: new Date().toISOString().split('T')[0],
    ...newCase.value,
    resultStatus: null,
    resultNotes: '',
    resultSentAt: null
  })
  showAddModal.value = false; newCase.value = { property: '', owner: '', priority: 'Medium', status: 'Pending', due: '' }; showToast('Audit case created successfully')
}

const handleUpdateCase = () => {
  const index = auditCases.value.findIndex(c => c.id === selectedCase.value.id)
  if (index !== -1) {
    const updated = { ...auditCases.value[index], ...editCase.value }
    auditCases.value[index] = updated
    if (updated.status === 'Completed') {
      moveAuditToSuccessful(updated)
      auditCases.value = auditCases.value.filter(c => c.id !== updated.id)
    }
    if (updated.status === 'Flagged') {
      moveAuditToFlagged(updated)
      auditCases.value = auditCases.value.filter(c => c.id !== updated.id)
    }
    showToast('Audit case updated')
  }
  showEditModal.value = false
}

const handleSendResult = () => {
  const index = auditCases.value.findIndex(c => c.id === selectedCase.value.id)
  if (index === -1) return
  const resultStatus = resultForm.value.status
  const nextStatus = resultStatus === 'Compliant' ? 'Completed' : 'Flagged'
  const updated = {
    ...auditCases.value[index],
    status: nextStatus,
    resultStatus,
    resultNotes: resultForm.value.message,
    resultSentAt: new Date().toISOString().split('T')[0],
  }
  auditCases.value[index] = updated
  if (updated.status === 'Completed') {
    moveAuditToSuccessful(updated)
    auditCases.value = auditCases.value.filter(c => c.id !== updated.id)
  }
  if (updated.status === 'Flagged') {
    moveAuditToFlagged(updated)
    auditCases.value = auditCases.value.filter(c => c.id !== updated.id)
  }
  showResultModal.value = false
  showToast('Result sent to taxpayer')
}

const handleDeleteCase = () => {
  auditCases.value = auditCases.value.filter(c => c.id !== selectedCase.value.id)
  showEditModal.value = false; showToast('Audit case deleted')
}
</script>

<template>
  <div class="min-h-screen bg-[#f5f6fa] flex">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />
    <div class="flex-1 flex flex-col">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6 shrink-0">
        <div class="flex items-center gap-4"><span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span><span class="text-[#d1d5db]">/</span><span class="text-[#1f2937] text-sm font-medium">Audit Cases</span></div>
        <div class="flex items-center gap-4"><span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span><button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button></div>
      </header>
      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between"><h2 class="text-[13px] font-semibold text-[#1f2937]">Audit Cases</h2><button @click="showAddModal = true" class="btn-primary text-[11px]">Create Case</button></div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4">
            <input v-model="searchQuery" type="text" placeholder="Search by case ID, property, or owner..." class="input-field flex-1" />
            <select v-model="filterPriority" class="input-field w-48"><option value="all">All Priority</option><option value="Critical">Critical</option><option value="High">High</option><option value="Medium">Medium</option><option value="Low">Low</option></select>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead><tr><th class="table-header">Case ID</th><th class="table-header">Property</th><th class="table-header">Owner</th><th class="table-header">Auditor</th><th class="table-header">Priority</th><th class="table-header">Status</th><th class="table-header">Result</th><th class="table-header">Started</th><th class="table-header">Due Date</th><th class="table-header">Actions</th></tr></thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="audit in paginatedCases" :key="audit.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ audit.id }}</td><td class="table-cell">{{ audit.property }}</td><td class="table-cell">{{ audit.owner }}</td><td class="table-cell">{{ audit.auditor }}</td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-red-50 text-red-700': audit.priority === 'Critical', 'bg-orange-50 text-orange-700': audit.priority === 'High', 'bg-yellow-50 text-yellow-700': audit.priority === 'Medium', 'bg-gray-100 text-gray-600': audit.priority === 'Low'}">{{ audit.priority }}</span></td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-blue-50 text-blue-700': audit.status === 'In Progress', 'bg-yellow-50 text-yellow-700': audit.status === 'Pending', 'bg-green-50 text-green-700': audit.status === 'Completed', 'bg-red-50 text-red-700': audit.status === 'Flagged'}">{{ audit.status }}</span></td>
                  <td class="table-cell">
                    <span class="px-2 py-0.5 text-[11px] font-medium rounded-full"
                      :class="{
                        'bg-gray-100 text-gray-600': !audit.resultStatus,
                        'bg-green-50 text-green-700': audit.resultStatus === 'Compliant',
                        'bg-red-50 text-red-700': audit.resultStatus === 'Non-Compliant',
                      }">{{ audit.resultStatus || 'Pending' }}</span>
                  </td>
                  <td class="table-cell text-[#9ca3af]">{{ audit.started }}</td><td class="table-cell text-[#9ca3af]">{{ audit.due }}</td>
                  <td class="table-cell">
                    <div class="flex gap-2">
                      <button @click="openViewModal(audit)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">View</button>
                      <button @click="openEditModal(audit)" class="px-3 py-1 text-[11px] bg-green-50 text-green-700 rounded hover:bg-green-100">Edit</button>
                      <button @click="openResultModal(audit)" class="px-3 py-1 text-[11px] bg-[#2D5A27]/10 text-[#2D5A27] rounded hover:bg-[#2D5A27]/20">Send Result</button>
                    </div>
                  </td>
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
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Create Audit Case</h3><button @click="showAddModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Property</label><input v-model="newCase.property" type="text" placeholder="Enter property address" class="input-field w-full" /></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Owner</label><input v-model="newCase.owner" type="text" placeholder="Enter owner name" class="input-field w-full" /></div>
            <div class="grid grid-cols-2 gap-4">
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Priority</label><select v-model="newCase.priority" class="input-field w-full"><option>Low</option><option>Medium</option><option>High</option><option>Critical</option></select></div>
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Due Date</label><input v-model="newCase.due" type="date" class="input-field w-full" /></div>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-end"><button @click="showAddModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button><button @click="handleAddCase" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Create</button></div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showEditModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Edit Audit Case</h3><button @click="showEditModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Priority</label><select v-model="editCase.priority" class="input-field w-full"><option>Low</option><option>Medium</option><option>High</option><option>Critical</option></select></div>
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Status</label><select v-model="editCase.status" class="input-field w-full"><option>Completed</option><option>Flagged</option></select></div>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-between"><button @click="handleDeleteCase" class="px-4 py-2 text-[11px] text-red-600 border border-red-200 rounded-lg hover:bg-red-50">Delete</button><div class="flex gap-3"><button @click="showEditModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button><button @click="handleUpdateCase" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Save</button></div></div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showViewModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Audit Case Details</h3><button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div><p class="text-[11px] text-gray-500">Case ID</p><p class="text-[13px] font-medium">{{ selectedCase?.id }}</p></div>
              <div><p class="text-[11px] text-gray-500">Status</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-blue-50 text-blue-700': selectedCase?.status === 'In Progress', 'bg-yellow-50 text-yellow-700': selectedCase?.status === 'Pending', 'bg-green-50 text-green-700': selectedCase?.status === 'Completed'}">{{ selectedCase?.status }}</span></div>
              <div><p class="text-[11px] text-gray-500">Priority</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-red-50 text-red-700': selectedCase?.priority === 'Critical', 'bg-orange-50 text-orange-700': selectedCase?.priority === 'High'}">{{ selectedCase?.priority }}</span></div>
              <div><p class="text-[11px] text-gray-500">Auditor</p><p class="text-[13px]">{{ selectedCase?.auditor }}</p></div>
            </div>
            <div><p class="text-[11px] text-gray-500">Property</p><p class="text-[13px]">{{ selectedCase?.property }}</p></div>
            <div><p class="text-[11px] text-gray-500">Owner</p><p class="text-[13px]">{{ selectedCase?.owner }}</p></div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <p class="text-[11px] text-gray-500">Result</p>
                <span class="px-2 py-0.5 text-[11px] font-medium rounded-full"
                  :class="{
                    'bg-gray-100 text-gray-600': !selectedCase?.resultStatus,
                    'bg-green-50 text-green-700': selectedCase?.resultStatus === 'Compliant',
                    'bg-red-50 text-red-700': selectedCase?.resultStatus === 'Non-Compliant',
                  }">{{ selectedCase?.resultStatus || 'Pending' }}</span>
              </div>
              <div>
                <p class="text-[11px] text-gray-500">Result Sent</p>
                <p class="text-[13px]">{{ selectedCase?.resultSentAt || 'Not sent' }}</p>
              </div>
            </div>
            <div v-if="selectedCase?.resultNotes">
              <p class="text-[11px] text-gray-500">Result Notes</p>
              <p class="text-[13px]">{{ selectedCase?.resultNotes }}</p>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div><p class="text-[11px] text-gray-500">Started</p><p class="text-[13px]">{{ selectedCase?.started }}</p></div>
              <div><p class="text-[11px] text-gray-500">Due Date</p><p class="text-[13px]">{{ selectedCase?.due }}</p></div>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-between">
            <button @click="openResultModal(selectedCase)" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Send Result</button>
            <button @click="showViewModal = false" class="px-4 py-2 text-[11px] bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200">Close</button>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showResultModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Send Filing Result</h3>
            <button @click="showResultModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-4">
            <div>
              <p class="text-[11px] text-gray-500">Case</p>
              <p class="text-[13px] font-medium">{{ selectedCase?.id }} — {{ selectedCase?.property }}</p>
            </div>
            <div>
              <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Result</label>
              <select v-model="resultForm.status" class="input-field w-full">
                <option>Compliant</option>
                <option>Non-Compliant</option>
              </select>
            </div>
            <div>
              <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Notes to Taxpayer</label>
              <textarea v-model="resultForm.message" rows="4" placeholder="Explain the outcome, required actions, or next steps." class="w-full px-4 py-3 border border-[#e5e7eb] rounded-lg text-[13px] resize-none focus:ring-2 focus:ring-[#2D5A27] focus:border-transparent"></textarea>
            </div>
            <div class="rounded-lg bg-[#f9fafb] border border-[#e5e7eb] p-3 text-[11px] text-[#6b7280]">
              This will appear in the taxpayer portal as the official filing result.
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-end">
            <button @click="showResultModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button>
            <button @click="handleSendResult" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Send</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
