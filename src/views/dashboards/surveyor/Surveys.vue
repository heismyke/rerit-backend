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
const filterPriority = ref('all')
const currentPage = ref(1)
const itemsPerPage = ref(5)

const showViewModal = ref(false)
const selectedSurvey = ref<any>(null)
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })

const surveys = ref([
  { id: 'SURV-001', property: 'Plot 8, Banana Island', owner: 'Folake Adeyemi', type: 'Property Verification', priority: 'High', status: 'Pending', assignedDate: '2024-01-15', dueDate: '2024-01-18', coordinates: '6.4281° N, 3.4219° E', description: 'Verify property exists and matches declared records. High-value property flagged for investigation.' },
  { id: 'SURV-002', property: 'Block 12, Maitama', owner: 'Unknown Entity', type: 'Ownership Verification', priority: 'Critical', status: 'Pending', assignedDate: '2024-01-14', dueDate: '2024-01-17', coordinates: '9.0579° N, 7.4951° E', description: 'Verify ownership claims. Property flagged for suspected tax evasion.' },
  { id: 'SURV-003', property: 'Plot 45, Victoria Island', owner: 'Chinedu & Partners', type: 'Value Assessment', priority: 'Medium', status: 'Pending', assignedDate: '2024-01-10', dueDate: '2024-01-15', coordinates: '6.4281° N, 3.4219° E', description: 'Physical inspection for value discrepancy verification.' },
  { id: 'SURV-004', property: 'Estate 7, Lekki', owner: 'Global Ventures Ltd', type: 'Document Verification', priority: 'High', status: 'Pending', assignedDate: '2024-01-16', dueDate: '2024-01-19', coordinates: '6.4312° N, 3.5012° E', description: 'Verify documents submitted. Property flagged for document forgery investigation.' },
  { id: 'SURV-005', property: 'Plot 15, Ikoyi', owner: 'Emeka Okonkwo', type: 'Routine Survey', priority: 'Low', status: 'Pending', assignedDate: '2024-01-05', dueDate: '2024-01-12', coordinates: '6.4536° N, 3.3958° E', description: 'Routine property verification survey.' },
])

const filteredSurveys = computed(() => surveys.value.filter(s => {
  const matchesSearch = s.property.toLowerCase().includes(searchQuery.value.toLowerCase()) || s.owner.toLowerCase().includes(searchQuery.value.toLowerCase()) || s.id.toLowerCase().includes(searchQuery.value.toLowerCase())
  const matchesStatus = filterStatus.value === 'all' || s.status === filterStatus.value
  const matchesPriority = filterPriority.value === 'all' || s.priority === filterPriority.value
  return matchesSearch && matchesStatus && matchesPriority
}))

const totalPages = computed(() => Math.ceil(filteredSurveys.value.length / itemsPerPage.value))
const paginatedSurveys = computed(() => { const start = (currentPage.value - 1) * itemsPerPage.value; return filteredSurveys.value.slice(start, start + itemsPerPage.value) })
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }

const showToast = (message: string) => { toast.value = { show: true, message }; setTimeout(() => { toast.value.show = false }, 3000) }
const openViewModal = (s: any) => { selectedSurvey.value = s; showViewModal.value = true }

const startInspection = (s: any) => {
  const index = surveys.value.findIndex(x => x.id === s.id)
  if (index !== -1) surveys.value[index].status = 'In Progress'
  showToast('Inspection started for ' + s.id + '. Proceed to Submissions to upload findings.')
  showViewModal.value = false
}
</script>

<template>
  <div class="min-h-screen flex bg-[#f5f7fa]">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />
    <div class="flex-1 flex flex-col">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6 shrink-0">
        <div class="flex items-center gap-4"><span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span><span class="text-[#d1d5db]">/</span><span class="text-[#1f2937] text-sm font-medium">Assigned Surveys</span></div>
        <div class="flex items-center gap-4"><span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span><button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button></div>
      </header>
      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between"><h2 class="text-[13px] font-semibold text-[#1f2937]">Assigned Surveys</h2><span class="text-[11px] text-[#6b7280]">Surveyor: {{ user?.name || 'Agent' }}</span></div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4 flex-wrap">
            <input v-model="searchQuery" type="text" placeholder="Search by property, owner, or ID..." class="input-field max-w-md" />
            <select v-model="filterStatus" class="input-field w-40"><option value="all">All Status</option><option value="Pending">Pending</option><option value="In Progress">In Progress</option><option value="Completed">Completed</option></select>
            <select v-model="filterPriority" class="input-field w-40"><option value="all">All Priority</option><option value="Critical">Critical</option><option value="High">High</option><option value="Medium">Medium</option><option value="Low">Low</option></select>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead><tr><th class="table-header">Survey ID</th><th class="table-header">Property</th><th class="table-header">Owner</th><th class="table-header">Type</th><th class="table-header">Priority</th><th class="table-header">Status</th><th class="table-header">Due Date</th><th class="table-header">Actions</th></tr></thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="survey in paginatedSurveys" :key="survey.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ survey.id }}</td>
                  <td class="table-cell text-[#6b7280]">{{ survey.property }}</td>
                  <td class="table-cell">{{ survey.owner }}</td>
                  <td class="table-cell text-[11px]">{{ survey.type }}</td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-red-100 text-red-700': survey.priority === 'Critical', 'bg-orange-100 text-orange-700': survey.priority === 'High', 'bg-yellow-100 text-yellow-700': survey.priority === 'Medium', 'bg-gray-100 text-gray-600': survey.priority === 'Low'}">{{ survey.priority }}</span></td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-100 text-green-700': survey.status === 'Completed', 'bg-blue-100 text-blue-700': survey.status === 'In Progress', 'bg-gray-100 text-gray-600': survey.status === 'Pending'}">{{ survey.status }}</span></td>
                  <td class="table-cell text-[#9ca3af]">{{ survey.dueDate }}</td>
                  <td class="table-cell"><button @click="openViewModal(survey)" class="px-3 py-1 text-[11px] bg-[#2D5A27] text-white rounded hover:bg-[#1e3d1a]">View Details</button></td>
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
      <div v-if="showViewModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4 overflow-y-auto">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-2xl my-8">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Survey Assignment Details</h3>
            <button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-6">
            <div class="grid grid-cols-2 gap-4">
              <div class="bg-[#EEEEEE] rounded-lg p-4">
                <h4 class="text-[11px] text-[#6b7280] mb-3 font-semibold">SURVEY INFO</h4>
                <div class="space-y-2">
                  <div><p class="text-[10px] text-[#9ca3af]">Survey ID</p><p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedSurvey?.id }}</p></div>
                  <div><p class="text-[10px] text-[#9ca3af]">Survey Type</p><p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedSurvey?.type }}</p></div>
                  <div><p class="text-[10px] text-[#9ca3af]">Priority</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-red-100 text-red-700': selectedSurvey?.priority === 'Critical', 'bg-orange-100 text-orange-700': selectedSurvey?.priority === 'High', 'bg-yellow-100 text-yellow-700': selectedSurvey?.priority === 'Medium'}">{{ selectedSurvey?.priority }}</span></div>
                </div>
              </div>
              <div class="bg-[#EEEEEE] rounded-lg p-4">
                <h4 class="text-[11px] text-[#6b7280] mb-3 font-semibold">PROPERTY INFO</h4>
                <div class="space-y-2">
                  <div><p class="text-[10px] text-[#9ca3af]">Property</p><p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedSurvey?.property }}</p></div>
                  <div><p class="text-[10px] text-[#9ca3af]">Owner</p><p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedSurvey?.owner }}</p></div>
                  <div><p class="text-[10px] text-[#9ca3af]">GPS Coordinates</p><p class="text-[12px] font-semibold text-[#1f2937]">{{ selectedSurvey?.coordinates }}</p></div>
                </div>
              </div>
            </div>

            <div class="bg-[#EEEEEE] rounded-lg p-4">
              <h4 class="text-[11px] text-[#6b7280] mb-3 font-semibold">ASSIGNMENT DETAILS</h4>
              <div class="grid grid-cols-3 gap-4">
                <div><p class="text-[10px] text-[#9ca3af]">Assigned Date</p><p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedSurvey?.assignedDate }}</p></div>
                <div><p class="text-[10px] text-[#9ca3af]">Due Date</p><p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedSurvey?.dueDate }}</p></div>
                <div><p class="text-[10px] text-[#9ca3af]">Status</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-100 text-green-700': selectedSurvey?.status === 'Completed', 'bg-blue-100 text-blue-700': selectedSurvey?.status === 'In Progress', 'bg-gray-100 text-gray-600': selectedSurvey?.status === 'Pending'}">{{ selectedSurvey?.status }}</span></div>
              </div>
              <div class="mt-4"><p class="text-[10px] text-[#9ca3af]">Description</p><p class="text-[13px] text-[#1f2937] mt-1">{{ selectedSurvey?.description }}</p></div>
            </div>

            <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
              <h4 class="text-[11px] text-blue-800 mb-2 font-semibold">NEXT STEPS</h4>
              <ol class="text-[12px] text-blue-700 space-y-1 list-decimal list-inside">
                <li>Review the survey assignment details above</li>
                <li>Click "Start Inspection" to begin</li>
                <li>Visit the property location ({{ selectedSurvey?.coordinates }})</li>
                <li>Go to <strong>Submissions</strong> page to upload photos, measurements, and findings</li>
              </ol>
            </div>

            <div class="flex items-center justify-between pt-4 border-t border-[#e5e7eb]">
              <div class="text-[11px] text-[#6b7280]">Complete the inspection and submit findings in Submissions</div>
              <div class="flex gap-3">
                <button @click="showViewModal = false" class="px-4 py-2 text-[11px] border border-[#e5e7eb] text-[#6b7280] rounded-lg hover:bg-gray-50">Close</button>
                <button v-if="selectedSurvey?.status === 'Pending'" @click="startInspection(selectedSurvey)" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Start Inspection</button>
                <button v-else-if="selectedSurvey?.status === 'In Progress'" @click="showViewModal = false; router.push('/surveyor/submissions')" class="px-4 py-2 text-[11px] bg-blue-600 text-white rounded-lg hover:bg-blue-700">Submit Findings</button>
                <button v-else class="px-4 py-2 text-[11px] bg-gray-200 text-gray-500 rounded-lg cursor-not-allowed">Completed</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
