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

const showSubmitModal = ref(false)
const showViewModal = ref(false)
const selectedSubmission = ref<any>(null)
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })

const newSubmission = ref({
  surveyId: '',
  property: '',
  owner: '',
  coordinates: '',
  photos: '',
  measurements: '',
  condition: 'Good',
  occupancy: 'Occupied',
  findings: ''
})

const submissions = ref([
  { id: 'SUB-001', surveyId: 'SURV-001', property: 'Plot 15, Ikoyi', owner: 'Emeka Okonkwo', type: 'Routine Survey', status: 'Approved', submittedDate: '2024-01-10', verifiedBy: 'Admin', photos: '12 photos', coordinates: '6.4536° N, 3.3958° E', measurements: '2,500 sqm', condition: 'Excellent', occupancy: 'Occupied' },
  { id: 'SUB-002', surveyId: 'SURV-002', property: 'Block 8, Victoria Island', owner: 'Ngozi Adebayo', type: 'Property Verification', status: 'Under Review', submittedDate: '2024-01-12', verifiedBy: '-', photos: '8 photos', coordinates: '6.4281° N, 3.4219° E', measurements: '1,800 sqm', condition: 'Good', occupancy: 'Partially Occupied' },
  { id: 'SUB-003', surveyId: 'SURV-003', property: 'Plot 45, Victoria Island', owner: 'Chinedu & Partners', type: 'Value Assessment', status: 'Approved', submittedDate: '2024-01-13', verifiedBy: 'Admin', photos: '15 photos', coordinates: '6.4281° N, 3.4219° E', measurements: '3,200 sqm', condition: 'Good', occupancy: 'Occupied' },
  { id: 'SUB-004', surveyId: 'SURV-004', property: 'Estate 3, Lekki', owner: 'Aisha Bello', type: 'Ownership Verification', status: 'Pending', submittedDate: '2024-01-14', verifiedBy: '-', photos: '10 photos', coordinates: '6.4312° N, 3.5012° E', measurements: '2,100 sqm', condition: 'Fair', occupancy: 'Vacant' },
  { id: 'SUB-005', surveyId: 'SURV-005', property: 'Block 5, Maitama', owner: 'Obi Foundation', type: 'Document Verification', status: 'Approved', submittedDate: '2024-01-08', verifiedBy: 'Admin', photos: '20 photos', coordinates: '9.0579° N, 7.4951° E', measurements: '4,500 sqm', condition: 'Excellent', occupancy: 'Occupied' },
])

const filteredSubmissions = computed(() => submissions.value.filter(s => {
  const matchesSearch = s.property.toLowerCase().includes(searchQuery.value.toLowerCase()) || s.owner.toLowerCase().includes(searchQuery.value.toLowerCase()) || s.id.toLowerCase().includes(searchQuery.value.toLowerCase())
  const matchesStatus = filterStatus.value === 'all' || s.status === filterStatus.value
  return matchesSearch && matchesStatus
}))

const totalPages = computed(() => Math.ceil(filteredSubmissions.value.length / itemsPerPage.value))
const paginatedSubmissions = computed(() => { const start = (currentPage.value - 1) * itemsPerPage.value; return filteredSubmissions.value.slice(start, start + itemsPerPage.value) })
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }

const showToast = (message: string) => { toast.value = { show: true, message }; setTimeout(() => { toast.value.show = false }, 3000) }
const openSubmitModal = () => {
  newSubmission.value = { surveyId: '', property: '', owner: '', coordinates: '', photos: '', measurements: '', condition: 'Good', occupancy: 'Occupied', findings: '' }
  showSubmitModal.value = true
}
const openViewModal = (s: any) => { selectedSubmission.value = s; showViewModal.value = true }
const handleSubmit = () => {
  if (!newSubmission.value.surveyId || !newSubmission.value.property) {
    showToast('Please fill in required fields')
    return
  }
  const newId = 'SUB-' + String(submissions.value.length + 1).padStart(3, '0')
  submissions.value.unshift({
    id: newId,
    surveyId: newSubmission.value.surveyId,
    property: newSubmission.value.property,
    owner: newSubmission.value.owner,
    type: 'Property Inspection',
    status: 'Pending',
    submittedDate: new Date().toISOString().split('T')[0],
    verifiedBy: '-',
    photos: newSubmission.value.photos || '0 photos',
    coordinates: newSubmission.value.coordinates || 'Not recorded',
    measurements: newSubmission.value.measurements || 'Not measured',
    condition: newSubmission.value.condition,
    occupancy: newSubmission.value.occupancy
  })
  showSubmitModal.value = false
  showToast('Survey submission created successfully')
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
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between"><h2 class="text-[13px] font-semibold text-[#1f2937]">Survey Submissions</h2><button @click="openSubmitModal" class="btn-primary text-[11px]">New Submission</button></div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4">
            <input v-model="searchQuery" type="text" placeholder="Search by property, owner, or ID..." class="input-field max-w-md" />
            <select v-model="filterStatus" class="input-field w-40"><option value="all">All Status</option><option value="Approved">Approved</option><option value="Under Review">Under Review</option><option value="Pending">Pending</option></select>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead><tr><th class="table-header">Submission ID</th><th class="table-header">Survey ID</th><th class="table-header">Property</th><th class="table-header">Owner</th><th class="table-header">Photos</th><th class="table-header">Status</th><th class="table-header">Date</th><th class="table-header">Actions</th></tr></thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="sub in paginatedSubmissions" :key="sub.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ sub.id }}</td>
                  <td class="table-cell text-[#6b7280]">{{ sub.surveyId }}</td>
                  <td class="table-cell">{{ sub.property }}</td>
                  <td class="table-cell text-[11px]">{{ sub.owner }}</td>
                  <td class="table-cell text-[#9ca3af]">{{ sub.photos }}</td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': sub.status === 'Approved', 'bg-yellow-50 text-yellow-700': sub.status === 'Under Review', 'bg-gray-100 text-gray-600': sub.status === 'Pending'}">{{ sub.status }}</span></td>
                  <td class="table-cell text-[#9ca3af]">{{ sub.submittedDate }}</td>
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
      <div v-if="showSubmitModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4 overflow-y-auto">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-2xl my-8">
          <div class="bg-[#B90B0B] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Submit Survey Report</h3>
            <button @click="showSubmitModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-6">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Survey ID *</label>
                <input v-model="newSubmission.surveyId" type="text" placeholder="e.g., SURV-001" class="input-field w-full" />
              </div>
              <div>
                <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Property Location *</label>
                <input v-model="newSubmission.property" type="text" placeholder="Enter property address" class="input-field w-full" />
              </div>
              <div>
                <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Owner Name</label>
                <input v-model="newSubmission.owner" type="text" placeholder="Enter owner name" class="input-field w-full" />
              </div>
              <div>
                <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">GPS Coordinates</label>
                <input v-model="newSubmission.coordinates" type="text" placeholder="e.g., 6.4281° N, 3.4219° E" class="input-field w-full" />
              </div>
            </div>

            <div class="bg-[#EEEEEE] rounded-lg p-4">
              <h4 class="text-[11px] text-[#6b7280] mb-3 font-semibold">PHYSICAL INSPECTION</h4>
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Number of Photos</label>
                  <input v-model="newSubmission.photos" type="text" placeholder="e.g., 12 photos" class="input-field w-full" />
                </div>
                <div>
                  <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Measurements</label>
                  <input v-model="newSubmission.measurements" type="text" placeholder="e.g., 2,500 sqm" class="input-field w-full" />
                </div>
                <div>
                  <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Condition</label>
                  <select v-model="newSubmission.condition" class="input-field w-full">
                    <option>Excellent</option>
                    <option>Good</option>
                    <option>Fair</option>
                    <option>Poor</option>
                  </select>
                </div>
                <div>
                  <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Occupancy Status</label>
                  <select v-model="newSubmission.occupancy" class="input-field w-full">
                    <option>Occupied</option>
                    <option>Vacant</option>
                    <option>Partially Occupied</option>
                    <option>Under Construction</option>
                  </select>
                </div>
              </div>
            </div>

            <div>
              <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Survey Findings & Notes</label>
              <textarea v-model="newSubmission.findings" rows="4" placeholder="Describe your findings from the physical inspection..." class="w-full px-4 py-3 border border-[#e5e7eb] rounded-lg text-[13px] resize-none focus:ring-2 focus:ring-[#B90B0B] focus:border-transparent"></textarea>
            </div>

            <div class="flex items-center justify-between pt-4 border-t border-[#e5e7eb]">
              <div class="text-[11px] text-[#6b7280]">Ensure all inspection checklist items are completed</div>
              <div class="flex gap-3">
                <button @click="showSubmitModal = false" class="px-4 py-2 text-[11px] border border-[#e5e7eb] text-[#6b7280] rounded-lg hover:bg-gray-50">Cancel</button>
                <button @click="handleSubmit" class="px-4 py-2 text-[11px] bg-[#B90B0B] text-white rounded-lg hover:bg-[#991010]">Submit Report</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showViewModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4 overflow-y-auto">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-2xl my-8">
          <div class="bg-[#B90B0B] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Submission Details</h3>
            <button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-6">
            <div class="grid grid-cols-2 gap-4">
              <div class="bg-[#EEEEEE] rounded-lg p-4">
                <h4 class="text-[11px] text-[#6b7280] mb-3 font-semibold">SUBMISSION INFO</h4>
                <div class="space-y-3">
                  <div><p class="text-[10px] text-[#9ca3af]">Submission ID</p><p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedSubmission?.id }}</p></div>
                  <div><p class="text-[10px] text-[#9ca3af]">Survey ID</p><p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedSubmission?.surveyId }}</p></div>
                  <div><p class="text-[10px] text-[#9ca3af]">Type</p><p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedSubmission?.type }}</p></div>
                </div>
              </div>
              <div class="bg-[#EEEEEE] rounded-lg p-4">
                <h4 class="text-[11px] text-[#6b7280] mb-3 font-semibold">PROPERTY INFO</h4>
                <div class="space-y-3">
                  <div><p class="text-[10px] text-[#9ca3af]">Property</p><p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedSubmission?.property }}</p></div>
                  <div><p class="text-[10px] text-[#9ca3af]">Owner</p><p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedSubmission?.owner }}</p></div>
                  <div><p class="text-[10px] text-[#9ca3af]">Coordinates</p><p class="text-[12px] font-semibold text-[#1f2937]">{{ selectedSubmission?.coordinates }}</p></div>
                </div>
              </div>
            </div>

            <div class="bg-[#EEEEEE] rounded-lg p-4">
              <h4 class="text-[11px] text-[#6b7280] mb-3 font-semibold">PHYSICAL VERIFICATION</h4>
              <div class="grid grid-cols-3 gap-4">
                <div><p class="text-[10px] text-[#9ca3af]">Photos</p><p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedSubmission?.photos }}</p></div>
                <div><p class="text-[10px] text-[#9ca3af]">Measurements</p><p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedSubmission?.measurements }}</p></div>
                <div><p class="text-[10px] text-[#9ca3af]">Condition</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-100 text-green-700': selectedSubmission?.condition === 'Excellent', 'bg-blue-100 text-blue-700': selectedSubmission?.condition === 'Good', 'bg-yellow-100 text-yellow-700': selectedSubmission?.condition === 'Fair', 'bg-red-100 text-red-700': selectedSubmission?.condition === 'Poor'}">{{ selectedSubmission?.condition }}</span></div>
                <div><p class="text-[10px] text-[#9ca3af]">Occupancy</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full bg-blue-100 text-blue-700">{{ selectedSubmission?.occupancy }}</span></div>
                <div><p class="text-[10px] text-[#9ca3af]">Submitted</p><p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedSubmission?.submittedDate }}</p></div>
                <div><p class="text-[10px] text-[#9ca3af]">Status</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': selectedSubmission?.status === 'Approved', 'bg-yellow-50 text-yellow-700': selectedSubmission?.status === 'Under Review', 'bg-gray-100 text-gray-600': selectedSubmission?.status === 'Pending'}">{{ selectedSubmission?.status }}</span></div>
              </div>
            </div>

            <div class="flex items-center justify-between pt-4 border-t border-[#e5e7eb]">
              <div class="text-[11px] text-[#6b7280]">Verified by: {{ selectedSubmission?.verifiedBy }}</div>
              <button @click="showViewModal = false" class="px-4 py-2 text-[11px] bg-[#f3f4f6] text-[#374151] rounded-lg hover:bg-[#e5e7eb]">Close</button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
