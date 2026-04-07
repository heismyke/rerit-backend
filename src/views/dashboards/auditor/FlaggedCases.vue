<script setup lang="ts">
import { useRoleStore } from '@/stores'
import { useRouter } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import { ref, computed } from 'vue'
import { useAuditorCasesStore, type FlaggedCase } from '@/stores/auditorCasesStore'

const { selectedRole, user, logout } = useRoleStore()
const router = useRouter()
const handleLogout = () => { logout(); router.push('/') }

const searchQuery = ref('')
const filterReason = ref('all')
const currentPage = ref(1)
const itemsPerPage = ref(5)

const showAiModal = ref(false)
const showViewModal = ref(false)
const showEditModal = ref(false)
const showResultModal = ref(false)
const selectedFlag = ref<any>(null)
const aiInsights = ref('')
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })

const { flaggedCases } = useAuditorCasesStore()

const filteredFlags = computed(() => flaggedCases.value.filter(f => {
  const matchesSearch = f.id.toLowerCase().includes(searchQuery.value.toLowerCase())
    || f.filingId.toLowerCase().includes(searchQuery.value.toLowerCase())
    || f.property.toLowerCase().includes(searchQuery.value.toLowerCase())
    || f.taxpayer.toLowerCase().includes(searchQuery.value.toLowerCase())
  const matchesReason = filterReason.value === 'all' || f.reason === filterReason.value
  return matchesSearch && matchesReason
}))

const totalPages = computed(() => Math.ceil(filteredFlags.value.length / itemsPerPage.value))
const paginatedFlags = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  return filteredFlags.value.slice(start, start + itemsPerPage.value)
})
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }

const openAiModal = (f: any) => {
  selectedFlag.value = f
  aiInsights.value = ''
  showAiModal.value = true
}

const showToast = (message: string) => { toast.value = { show: true, message }; setTimeout(() => { toast.value.show = false }, 3000) }
const openViewModal = (f: any) => { selectedFlag.value = f; showViewModal.value = true }
const openEditModal = (f: any) => {
  selectedFlag.value = f
  editFlag.value = { priority: f.priority, status: f.status }
  showEditModal.value = true
}
const openResultModal = (f: any) => {
  selectedFlag.value = f
  resultForm.value = { status: f.resultStatus || 'Compliant', message: f.resultNotes || '' }
  showResultModal.value = true
}

const editFlag = ref<{ priority: FlaggedCase['priority']; status: FlaggedCase['status'] }>({ priority: 'Medium', status: 'Pending Review' })
const resultForm = ref<{ status: NonNullable<FlaggedCase['resultStatus']>; message: string }>({ status: 'Compliant', message: '' })

const handleUpdateFlag = () => {
  const index = flaggedCases.value.findIndex(f => f.id === selectedFlag.value.id)
  if (index !== -1) {
    flaggedCases.value[index] = { ...flaggedCases.value[index], ...editFlag.value }
    showToast('Flag updated')
  }
  showEditModal.value = false
}

const handleSendResult = () => {
  const index = flaggedCases.value.findIndex(f => f.id === selectedFlag.value.id)
  if (index === -1) return
  const resultStatus = resultForm.value.status
  flaggedCases.value[index] = {
    ...flaggedCases.value[index],
    resultStatus,
    resultNotes: resultForm.value.message,
    resultSentAt: new Date().toISOString().split('T')[0],
    status: resultStatus === 'Compliant' ? 'Resolved' : flaggedCases.value[index].status,
  }
  showResultModal.value = false
  showToast('Result sent to taxpayer')
}

const generateAiInsights = () => {
  if (!selectedFlag.value) return
  const f = selectedFlag.value
  aiInsights.value =
    `AI Summary for ${f.property}\n\n` +
    `• Filing ID: ${f.filingId}\n` +
    `• Reason: ${f.reason}\n` +
    `• Risk Signal: Medium–High\n\n` +
    `Suggested Actions:\n` +
    `1. Request supporting documents from ${f.taxpayer}\n` +
    `2. Cross-check land registry and ownership records\n` +
    `3. Compare declared values with market benchmarks\n`
}
</script>

<template>
  <div class="min-h-screen bg-[#f5f6fa] flex">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />
    <div class="flex-1 flex flex-col">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6 shrink-0">
        <div class="flex items-center gap-4">
          <span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span>
          <span class="text-[#d1d5db]">/</span>
          <span class="text-[#1f2937] text-sm font-medium">Flagged Cases</span>
        </div>
        <div class="flex items-center gap-4">
          <span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span>
          <button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button>
        </div>
      </header>

      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between">
            <h2 class="text-[13px] font-semibold text-[#1f2937]">Flagged From FCT‑IRS</h2>
            <span class="text-[11px] text-[#6b7280]">Mismatch between FCT‑IRS filing and RERIT data</span>
          </div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4">
            <input v-model="searchQuery" type="text" placeholder="Search by ID, filing, property, or taxpayer..." class="input-field flex-1" />
            <select v-model="filterReason" class="input-field w-64">
              <option value="all">All Reasons</option>
              <option>Declared rent below benchmark</option>
            </select>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr>
                  <th class="table-header">Flag ID</th>
                  <th class="table-header">FCT‑IRS Filing ID</th>
                  <th class="table-header">Property</th>
                  <th class="table-header">Taxpayer</th>
                  <th class="table-header">Reason</th>
                  <th class="table-header">Received</th>
                  <th class="table-header">Status</th>
                  <th class="table-header">Result</th>
                  <th class="table-header">Actions</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="flag in paginatedFlags" :key="flag.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ flag.id }}</td>
                  <td class="table-cell">{{ flag.filingId }}</td>
                  <td class="table-cell">{{ flag.property }}</td>
                  <td class="table-cell">{{ flag.taxpayer }}</td>
                  <td class="table-cell">
                    <span class="px-2 py-0.5 text-[11px] bg-red-50 text-red-700 rounded">{{ flag.reason }}</span>
                  </td>
                  <td class="table-cell text-[#9ca3af]">{{ flag.receivedAt }}</td>
                  <td class="table-cell">
                    <span class="px-2 py-0.5 text-[11px] font-medium rounded-full"
                      :class="{
                        'bg-yellow-50 text-yellow-700': flag.status === 'Pending Review',
                        'bg-blue-50 text-blue-700': flag.status === 'In Review',
                        'bg-green-50 text-green-700': flag.status === 'Resolved',
                      }">{{ flag.status }}</span>
                  </td>
                  <td class="table-cell">
                    <span class="px-2 py-0.5 text-[11px] font-medium rounded-full"
                      :class="{
                        'bg-gray-100 text-gray-600': !flag.resultStatus,
                        'bg-green-50 text-green-700': flag.resultStatus === 'Compliant',
                        'bg-red-50 text-red-700': flag.resultStatus === 'Non-Compliant',
                      }">{{ flag.resultStatus || 'Pending' }}</span>
                  </td>
                  <td class="table-cell">
                    <div class="flex gap-2">
                      <button @click="openViewModal(flag)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">View</button>
                      <button @click="openEditModal(flag)" class="px-3 py-1 text-[11px] bg-green-50 text-green-700 rounded hover:bg-green-100">Edit</button>
                      <button @click="openResultModal(flag)" class="px-3 py-1 text-[11px] bg-[#2D5A27]/10 text-[#2D5A27] rounded hover:bg-[#2D5A27]/20">Send Result</button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="px-6 py-4 flex items-center justify-between border-t border-[#f3f4f6]">
            <p class="text-[11px] text-[#6b7280]">Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, filteredFlags.length) }} of {{ filteredFlags.length }} entries</p>
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
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Flagged Case Details</h3>
            <button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div><p class="text-[11px] text-gray-500">Flag ID</p><p class="text-[13px] font-medium">{{ selectedFlag?.id }}</p></div>
              <div><p class="text-[11px] text-gray-500">Status</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-yellow-50 text-yellow-700': selectedFlag?.status === 'Pending Review', 'bg-blue-50 text-blue-700': selectedFlag?.status === 'In Review', 'bg-green-50 text-green-700': selectedFlag?.status === 'Resolved'}">{{ selectedFlag?.status }}</span></div>
              <div><p class="text-[11px] text-gray-500">Priority</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-red-50 text-red-700': selectedFlag?.priority === 'Critical', 'bg-orange-50 text-orange-700': selectedFlag?.priority === 'High', 'bg-yellow-50 text-yellow-700': selectedFlag?.priority === 'Medium', 'bg-gray-100 text-gray-600': selectedFlag?.priority === 'Low'}">{{ selectedFlag?.priority }}</span></div>
              <div><p class="text-[11px] text-gray-500">FCT‑IRS Filing ID</p><p class="text-[13px]">{{ selectedFlag?.filingId }}</p></div>
            </div>
            <div><p class="text-[11px] text-gray-500">Property</p><p class="text-[13px]">{{ selectedFlag?.property }}</p></div>
            <div><p class="text-[11px] text-gray-500">Taxpayer</p><p class="text-[13px]">{{ selectedFlag?.taxpayer }}</p></div>
            <div><p class="text-[11px] text-gray-500">Reason</p><p class="text-[13px]">{{ selectedFlag?.reason }}</p></div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <p class="text-[11px] text-gray-500">Result</p>
                <span class="px-2 py-0.5 text-[11px] font-medium rounded-full"
                  :class="{
                    'bg-gray-100 text-gray-600': !selectedFlag?.resultStatus,
                    'bg-green-50 text-green-700': selectedFlag?.resultStatus === 'Compliant',
                    'bg-red-50 text-red-700': selectedFlag?.resultStatus === 'Non-Compliant',
                  }">{{ selectedFlag?.resultStatus || 'Pending' }}</span>
              </div>
              <div>
                <p class="text-[11px] text-gray-500">Result Sent</p>
                <p class="text-[13px]">{{ selectedFlag?.resultSentAt || 'Not sent' }}</p>
              </div>
            </div>
            <div v-if="selectedFlag?.resultNotes">
              <p class="text-[11px] text-gray-500">Result Notes</p>
              <p class="text-[13px]">{{ selectedFlag?.resultNotes }}</p>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-between">
            <div class="flex gap-2">
              <button @click="openResultModal(selectedFlag)" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Send Result</button>
              <button @click="openAiModal(selectedFlag)" class="px-4 py-2 text-[11px] bg-[#111827] text-white rounded-lg hover:bg-[#0b1220]">AI Review</button>
            </div>
            <button @click="showViewModal = false" class="px-4 py-2 text-[11px] bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200">Close</button>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showEditModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Edit Flag</h3>
            <button @click="showEditModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Priority</label>
                <select v-model="editFlag.priority" class="input-field w-full"><option>Low</option><option>Medium</option><option>High</option><option>Critical</option></select>
              </div>
              <div>
                <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Status</label>
                <select v-model="editFlag.status" class="input-field w-full"><option>Pending Review</option><option>In Review</option><option>Resolved</option></select>
              </div>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-end">
            <div class="flex gap-3">
              <button @click="showEditModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button>
              <button @click="handleUpdateFlag" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Save</button>
            </div>
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
              <p class="text-[11px] text-gray-500">Flag</p>
              <p class="text-[13px] font-medium">{{ selectedFlag?.id }} — {{ selectedFlag?.property }}</p>
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

    <Teleport to="body">
      <div v-if="showAiModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-lg">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">AI Review</h3>
            <button @click="showAiModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-4">
            <div class="bg-[#f9fafb] rounded-lg p-4 border border-[#e5e7eb]">
              <p class="text-[11px] text-gray-500">Flag</p>
              <p class="text-[13px] font-medium">{{ selectedFlag?.id }} — {{ selectedFlag?.property }}</p>
              <p class="text-[11px] text-gray-500 mt-2">Reason: {{ selectedFlag?.reason }}</p>
            </div>
            <div>
              <label class="block text-[11px] font-medium text-gray-600 mb-1.5">AI Insights</label>
              <textarea v-model="aiInsights" rows="8" placeholder="Generate AI insights for this flagged case..." class="w-full px-4 py-3 border border-[#e5e7eb] rounded-lg text-[13px] resize-none focus:ring-2 focus:ring-[#2D5A27] focus:border-transparent"></textarea>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-end">
            <button @click="generateAiInsights" class="px-4 py-2 text-[11px] bg-[#1f2937] text-white rounded-lg hover:bg-[#111827]">Generate</button>
            <button @click="showAiModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Close</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
