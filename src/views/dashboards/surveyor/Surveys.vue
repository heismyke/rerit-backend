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

const checklist = ref({
  photos: { checked: false, files: [] as string[], notes: '' },
  gps: { checked: false, verified: false, notes: '' },
  measurements: { checked: false, length: '', width: '', total: '', notes: '' },
  condition: { checked: false, rating: '', notes: '' },
  occupancy: { checked: false, status: '', notes: '' },
  ownership: { checked: false, verified: false, notes: '' }
})

const surveys = ref([
  { id: 'SURV-001', property: 'Plot 8, Banana Island', owner: 'Folake Adeyemi', type: 'Property Verification', priority: 'High', status: 'Pending', assignedDate: '2024-01-15', dueDate: '2024-01-18', coordinates: '6.4281° N, 3.4219° E', description: 'Verify property exists and matches declared records. High-value property flagged for investigation.' },
  { id: 'SURV-002', property: 'Block 12, Maitama', owner: 'Unknown Entity', type: 'Ownership Verification', priority: 'Critical', status: 'In Progress', assignedDate: '2024-01-14', dueDate: '2024-01-17', coordinates: '9.0579° N, 7.4951° E', description: 'Verify ownership claims. Property flagged for suspected tax evasion.' },
  { id: 'SURV-003', property: 'Plot 45, Victoria Island', owner: 'Chinedu & Partners', type: 'Value Assessment', priority: 'Medium', status: 'Completed', assignedDate: '2024-01-10', dueDate: '2024-01-15', coordinates: '6.4281° N, 3.4219° E', description: 'Physical inspection for value discrepancy verification.', completedDate: '2024-01-13' },
  { id: 'SURV-004', property: 'Estate 7, Lekki', owner: 'Global Ventures Ltd', type: 'Document Verification', priority: 'High', status: 'Pending', assignedDate: '2024-01-16', dueDate: '2024-01-19', coordinates: '6.4312° N, 3.5012° E', description: 'Verify documents submitted. Property flagged for document forgery investigation.' },
  { id: 'SURV-005', property: 'Plot 15, Ikoyi', owner: 'Emeka Okonkwo', type: 'Routine Survey', priority: 'Low', status: 'Completed', assignedDate: '2024-01-05', dueDate: '2024-01-12', coordinates: '6.4536° N, 3.3958° E', description: 'Routine property verification survey.', completedDate: '2024-01-10' },
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

const resetChecklist = () => {
  checklist.value = {
    photos: { checked: false, files: [], notes: '' },
    gps: { checked: false, verified: false, notes: '' },
    measurements: { checked: false, length: '', width: '', total: '', notes: '' },
    condition: { checked: false, rating: '', notes: '' },
    occupancy: { checked: false, status: '', notes: '' },
    ownership: { checked: false, verified: false, notes: '' }
  }
}

const openViewModal = (s: any) => {
  resetChecklist()
  selectedSurvey.value = s
  showViewModal.value = true
}

const handlePhotoUpload = (event: Event) => {
  const input = event.target as HTMLInputElement
  if (input.files) {
    const fileNames = Array.from(input.files).map(f => f.name)
    checklist.value.photos.files = [...checklist.value.photos.files, ...fileNames]
    if (checklist.value.photos.files.length > 0) {
      checklist.value.photos.checked = true
    }
  }
}

const removePhoto = (index: number) => {
  checklist.value.photos.files.splice(index, 1)
  if (checklist.value.photos.files.length === 0) {
    checklist.value.photos.checked = false
  }
}

const verifyGPS = () => {
  checklist.value.gps.verified = true
  checklist.value.gps.checked = true
  showToast('GPS coordinates verified')
}

const calculateArea = () => {
  if (checklist.value.measurements.length && checklist.value.measurements.width) {
    const l = parseFloat(checklist.value.measurements.length)
    const w = parseFloat(checklist.value.measurements.width)
    if (!isNaN(l) && !isNaN(w)) {
      checklist.value.measurements.total = (l * w).toLocaleString() + ' sqm'
      checklist.value.measurements.checked = true
    }
  }
}

const allItemsChecked = computed(() => {
  return checklist.value.photos.checked &&
    checklist.value.gps.checked &&
    checklist.value.measurements.checked &&
    checklist.value.condition.checked &&
    checklist.value.occupancy.checked &&
    checklist.value.ownership.checked
})

const startInspection = (s: any) => {
  const index = surveys.value.findIndex(x => x.id === s.id)
  if (index !== -1) surveys.value[index].status = 'In Progress'
  showToast('Inspection started for ' + s.id)
}

const completeInspection = (s: any) => {
  if (!allItemsChecked.value) {
    showToast('Please complete all checklist items before submitting')
    return
  }
  const index = surveys.value.findIndex(x => x.id === s.id)
  if (index !== -1) {
    surveys.value[index].status = 'Completed'
    surveys.value[index].completedDate = new Date().toISOString().split('T')[0]
  }
  showToast('Inspection completed. Redirecting to submit report...')
  showViewModal.value = false
  setTimeout(() => {
    router.push('/surveyor/submissions')
  }, 1500)
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
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': survey.status === 'Completed', 'bg-blue-50 text-blue-700': survey.status === 'In Progress', 'bg-gray-100 text-gray-600': survey.status === 'Pending'}">{{ survey.status }}</span></td>
                  <td class="table-cell text-[#9ca3af]">{{ survey.dueDate }}</td>
                  <td class="table-cell"><button @click="openViewModal(survey)" class="px-3 py-1 text-[11px] bg-[#B90B0B] text-white rounded hover:bg-[#991010]">View Details</button></td>
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
        <div class="bg-white rounded-xl shadow-xl w-full max-w-3xl my-8">
          <div class="bg-[#B90B0B] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Survey Assignment - {{ selectedSurvey?.id }}</h3>
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
              <h4 class="text-[11px] text-[#6b7280] mb-3 font-semibold">INSPECTION CHECKLIST</h4>
              <div class="space-y-4">
                <div class="border border-[#d1d5db] rounded-lg p-4 bg-white">
                  <div class="flex items-start gap-3">
                    <input type="checkbox" v-model="checklist.photos.checked" :disabled="selectedSurvey?.status === 'Completed'" class="w-5 h-5 mt-0.5 rounded border-[#d1d5db]" />
                    <div class="flex-1">
                      <p class="text-[13px] font-semibold text-[#1f2937] mb-2">Property Photos</p>
                      <div v-if="selectedSurvey?.status !== 'Completed'">
                        <label class="block mb-2">
                          <span class="px-3 py-1.5 text-[11px] bg-[#B90B0B] text-white rounded cursor-pointer hover:bg-[#991010] inline-block">
                            📷 Upload Photos
                          </span>
                          <input type="file" multiple accept="image/*" @change="handlePhotoUpload" class="hidden" />
                        </label>
                        <div v-if="checklist.photos.files.length > 0" class="flex flex-wrap gap-2 mb-2">
                          <span v-for="(file, idx) in checklist.photos.files" :key="idx" class="px-2 py-1 bg-green-100 text-green-700 text-[11px] rounded flex items-center gap-1">
                            {{ file }}
                            <button @click="removePhoto(idx)" class="text-red-500 hover:text-red-700">×</button>
                          </span>
                        </div>
                        <textarea v-model="checklist.photos.notes" rows="2" placeholder="Add notes about the photos..." class="w-full px-3 py-2 border border-[#e5e7eb] rounded text-[12px] resize-none"></textarea>
                      </div>
                      <div v-else-if="checklist.photos.files.length > 0" class="flex flex-wrap gap-2">
                        <span v-for="(file, idx) in checklist.photos.files" :key="idx" class="px-2 py-1 bg-green-100 text-green-700 text-[11px] rounded">{{ file }}</span>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="border border-[#d1d5db] rounded-lg p-4 bg-white">
                  <div class="flex items-start gap-3">
                    <input type="checkbox" v-model="checklist.gps.checked" :disabled="selectedSurvey?.status === 'Completed'" class="w-5 h-5 mt-0.5 rounded border-[#d1d5db]" />
                    <div class="flex-1">
                      <p class="text-[13px] font-semibold text-[#1f2937] mb-2">GPS Verification</p>
                      <div class="flex items-center gap-2 mb-2">
                        <span class="text-[12px] text-[#6b7280]">Coordinates: {{ selectedSurvey?.coordinates }}</span>
                        <button v-if="!checklist.gps.verified && selectedSurvey?.status !== 'Completed'" @click="verifyGPS" class="px-3 py-1 text-[11px] bg-green-600 text-white rounded hover:bg-green-700">Verify</button>
                        <span v-else class="px-2 py-0.5 bg-green-100 text-green-700 text-[11px] rounded">✓ Verified</span>
                      </div>
                      <textarea v-model="checklist.gps.notes" rows="2" placeholder="Add notes about GPS verification..." class="w-full px-3 py-2 border border-[#e5e7eb] rounded text-[12px] resize-none" :disabled="selectedSurvey?.status === 'Completed'"></textarea>
                    </div>
                  </div>
                </div>

                <div class="border border-[#d1d5db] rounded-lg p-4 bg-white">
                  <div class="flex items-start gap-3">
                    <input type="checkbox" v-model="checklist.measurements.checked" :disabled="selectedSurvey?.status === 'Completed'" class="w-5 h-5 mt-0.5 rounded border-[#d1d5db]" />
                    <div class="flex-1">
                      <p class="text-[13px] font-semibold text-[#1f2937] mb-2">Physical Measurements</p>
                      <div v-if="selectedSurvey?.status !== 'Completed'" class="flex items-end gap-2 mb-2">
                        <div>
                          <label class="text-[10px] text-[#6b7280]">Length (m)</label>
                          <input v-model="checklist.measurements.length" type="number" placeholder="0" class="w-24 px-2 py-1 border border-[#e5e7eb] rounded text-[12px]" />
                        </div>
                        <span class="text-[#6b7280] pb-1">×</span>
                        <div>
                          <label class="text-[10px] text-[#6b7280]">Width (m)</label>
                          <input v-model="checklist.measurements.width" type="number" placeholder="0" class="w-24 px-2 py-1 border border-[#e5e7eb] rounded text-[12px]" />
                        </div>
                        <button @click="calculateArea" class="px-3 py-1 text-[11px] bg-blue-600 text-white rounded hover:bg-blue-700">Calculate</button>
                        <span v-if="checklist.measurements.total" class="px-2 py-1 bg-green-100 text-green-700 text-[12px] rounded font-medium">= {{ checklist.measurements.total }}</span>
                      </div>
                      <div v-else class="mb-2">
                        <span class="text-[12px] text-[#6b7280]">Total Area: {{ checklist.measurements.total || 'N/A' }}</span>
                      </div>
                      <textarea v-model="checklist.measurements.notes" rows="2" placeholder="Add notes about measurements..." class="w-full px-3 py-2 border border-[#e5e7eb] rounded text-[12px] resize-none" :disabled="selectedSurvey?.status === 'Completed'"></textarea>
                    </div>
                  </div>
                </div>

                <div class="border border-[#d1d5db] rounded-lg p-4 bg-white">
                  <div class="flex items-start gap-3">
                    <input type="checkbox" v-model="checklist.condition.checked" :disabled="selectedSurvey?.status === 'Completed'" class="w-5 h-5 mt-0.5 rounded border-[#d1d5db]" />
                    <div class="flex-1">
                      <p class="text-[13px] font-semibold text-[#1f2937] mb-2">Condition Assessment</p>
                      <div v-if="selectedSurvey?.status !== 'Completed'" class="flex gap-2 mb-2">
                        <label v-for="rating in ['Excellent', 'Good', 'Fair', 'Poor']" :key="rating" class="flex items-center gap-1 cursor-pointer">
                          <input type="radio" v-model="checklist.condition.rating" :value="rating" @change="checklist.condition.checked = true" class="hidden peer" />
                          <span class="px-3 py-1 text-[11px] border border-[#d1d5db] rounded peer-checked:bg-[#B90B0B] peer-checked:text-white peer-checked:border-[#B90B0B]">{{ rating }}</span>
                        </label>
                      </div>
                      <span v-else class="px-2 py-0.5 bg-blue-100 text-blue-700 text-[11px] rounded">{{ checklist.condition.rating || 'N/A' }}</span>
                      <textarea v-model="checklist.condition.notes" rows="2" placeholder="Add notes about property condition..." class="w-full px-3 py-2 border border-[#e5e7eb] rounded text-[12px] resize-none mt-2" :disabled="selectedSurvey?.status === 'Completed'"></textarea>
                    </div>
                  </div>
                </div>

                <div class="border border-[#d1d5db] rounded-lg p-4 bg-white">
                  <div class="flex items-start gap-3">
                    <input type="checkbox" v-model="checklist.occupancy.checked" :disabled="selectedSurvey?.status === 'Completed'" class="w-5 h-5 mt-0.5 rounded border-[#d1d5db]" />
                    <div class="flex-1">
                      <p class="text-[13px] font-semibold text-[#1f2937] mb-2">Occupancy Verification</p>
                      <div v-if="selectedSurvey?.status !== 'Completed'" class="flex flex-wrap gap-2 mb-2">
                        <label v-for="status in ['Occupied', 'Vacant', 'Partially Occupied', 'Under Construction']" :key="status" class="flex items-center gap-1 cursor-pointer">
                          <input type="radio" v-model="checklist.occupancy.status" :value="status" @change="checklist.occupancy.checked = true" class="hidden peer" />
                          <span class="px-3 py-1 text-[11px] border border-[#d1d5db] rounded peer-checked:bg-[#B90B0B] peer-checked:text-white peer-checked:border-[#B90B0B]">{{ status }}</span>
                        </label>
                      </div>
                      <span v-else class="px-2 py-0.5 bg-blue-100 text-blue-700 text-[11px] rounded">{{ checklist.occupancy.status || 'N/A' }}</span>
                      <textarea v-model="checklist.occupancy.notes" rows="2" placeholder="Add notes about occupancy..." class="w-full px-3 py-2 border border-[#e5e7eb] rounded text-[12px] resize-none mt-2" :disabled="selectedSurvey?.status === 'Completed'"></textarea>
                    </div>
                  </div>
                </div>

                <div class="border border-[#d1d5db] rounded-lg p-4 bg-white">
                  <div class="flex items-start gap-3">
                    <input type="checkbox" v-model="checklist.ownership.checked" :disabled="selectedSurvey?.status === 'Completed'" class="w-5 h-5 mt-0.5 rounded border-[#d1d5db]" />
                    <div class="flex-1">
                      <p class="text-[13px] font-semibold text-[#1f2937] mb-2">Ownership Check</p>
                      <div class="flex items-center gap-2 mb-2">
                        <span class="text-[12px] text-[#6b7280]">Owner on record: {{ selectedSurvey?.owner }}</span>
                        <button v-if="!checklist.ownership.verified && selectedSurvey?.status !== 'Completed'" @click="checklist.ownership.verified = true; checklist.ownership.checked = true" class="px-3 py-1 text-[11px] bg-green-600 text-white rounded hover:bg-green-700">Verify Ownership</button>
                        <span v-else class="px-2 py-0.5 bg-green-100 text-green-700 text-[11px] rounded">✓ Verified</span>
                      </div>
                      <textarea v-model="checklist.ownership.notes" rows="2" placeholder="Add notes about ownership verification..." class="w-full px-3 py-2 border border-[#e5e7eb] rounded text-[12px] resize-none" :disabled="selectedSurvey?.status === 'Completed'"></textarea>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div class="flex items-center justify-between pt-4 border-t border-[#e5e7eb]">
              <div class="text-[11px] text-[#6b7280]">
                <span :class="allItemsChecked ? 'text-green-600 font-medium' : 'text-red-500'">
                  {{ allItemsChecked ? '✓ All items completed' : '○ Complete all checklist items before submitting' }}
                </span>
              </div>
              <div class="flex gap-3">
                <button v-if="selectedSurvey?.status === 'Pending'" @click="startInspection(selectedSurvey); showViewModal = false" class="px-4 py-2 text-[11px] bg-[#1f2937] text-white rounded-lg hover:bg-[#374151]">Start Inspection</button>
                <button v-else-if="selectedSurvey?.status === 'In Progress'" @click="completeInspection(selectedSurvey)" class="px-4 py-2 text-[11px] bg-green-600 text-white rounded-lg hover:bg-green-700">Complete & Submit</button>
                <button v-else-if="selectedSurvey?.status === 'Completed'" class="px-4 py-2 text-[11px] bg-gray-200 text-gray-500 rounded-lg cursor-not-allowed">Completed {{ selectedSurvey?.completedDate }}</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
