<script setup lang="ts">
import { useRoleStore } from '@/stores'
import { useRouter } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import { ref, computed } from 'vue'

const { selectedRole, user, logout } = useRoleStore()
const router = useRouter()
const handleLogout = () => {
  logout()
  router.push('/')
}

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
  photos: [] as string[],
  measurements: { length: '', width: '', total: '' },
  condition: '',
  occupancy: '',
  ownershipVerified: false,
  gpsVerified: false,
  findings: '',
})

const submissions = ref([
  {
    id: 'SUB-001',
    surveyId: 'SURV-001',
    property: 'Plot 8, Banana Island',
    owner: 'Folake Adeyemi',
    type: 'Property Verification',
    status: 'Pending',
    submittedDate: new Date().toISOString().split('T')[0],
    verifiedBy: '-',
    photos: [] as string[],
    measurements: '',
    condition: '',
    occupancy: '',
  },
  {
    id: 'SUB-002',
    surveyId: 'SURV-002',
    property: 'Block 8, Victoria Island',
    owner: 'Ngozi Adebayo',
    type: 'Property Verification',
    status: 'Pending',
    submittedDate: new Date().toISOString().split('T')[0],
    verifiedBy: '-',
    photos: [] as string[],
    measurements: '',
    condition: '',
    occupancy: '',
  },
])

const filteredSubmissions = computed(() =>
  submissions.value.filter((s) => {
    const matchesSearch =
      s.property.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      s.owner.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      s.id.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesStatus = filterStatus.value === 'all' || s.status === filterStatus.value
    return matchesSearch && matchesStatus
  }),
)

const totalPages = computed(() => Math.ceil(filteredSubmissions.value.length / itemsPerPage.value))
const paginatedSubmissions = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  return filteredSubmissions.value.slice(start, start + itemsPerPage.value)
})
const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) currentPage.value = page
}

const showToast = (message: string) => {
  toast.value = { show: true, message }
  setTimeout(() => {
    toast.value.show = false
  }, 3000)
}

const resetForm = () => {
  newSubmission.value = {
    surveyId: '',
    property: '',
    owner: '',
    coordinates: '',
    photos: [],
    measurements: { length: '', width: '', total: '' },
    condition: '',
    occupancy: '',
    ownershipVerified: false,
    gpsVerified: false,
    findings: '',
  }
}

const openSubmitModal = () => {
  resetForm()
  showSubmitModal.value = true
}

const openViewModal = (s: any) => {
  selectedSubmission.value = s
  showViewModal.value = true
}

const handlePhotoUpload = (event: Event) => {
  const input = event.target as HTMLInputElement
  if (input.files) {
    const fileNames = Array.from(input.files).map((f) => f.name)
    newSubmission.value.photos = [...newSubmission.value.photos, ...fileNames]
  }
}

const removePhoto = (index: number) => {
  newSubmission.value.photos.splice(index, 1)
}

const calculateArea = () => {
  if (newSubmission.value.measurements.length && newSubmission.value.measurements.width) {
    const l = parseFloat(newSubmission.value.measurements.length)
    const w = parseFloat(newSubmission.value.measurements.width)
    if (!isNaN(l) && !isNaN(w)) {
      newSubmission.value.measurements.total = (l * w).toLocaleString() + ' sqm'
    }
  }
}

const isFormValid = computed(() => {
  return (
    newSubmission.value.surveyId &&
    newSubmission.value.property &&
    newSubmission.value.photos.length > 0 &&
    newSubmission.value.measurements.total &&
    newSubmission.value.condition &&
    newSubmission.value.occupancy &&
    newSubmission.value.findings.trim()
  )
})

const handleSubmit = () => {
  if (!isFormValid.value) {
    showToast('Please complete all required fields')
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
    photos: [...newSubmission.value.photos],
    measurements: newSubmission.value.measurements.total,
    condition: newSubmission.value.condition,
    occupancy: newSubmission.value.occupancy,
  })
  showSubmitModal.value = false
  showToast('Survey submission created successfully')
}
</script>

<template>
  <div class="min-h-screen flex bg-[#f5f7fa]">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />
    <div class="flex-1 flex flex-col">
      <header
        class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6 shrink-0"
      >
        <div class="flex items-center gap-4">
          <span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span
          ><span class="text-[#d1d5db]">/</span
          ><span class="text-[#1f2937] text-sm font-medium">Submissions</span>
        </div>
        <div class="flex items-center gap-4">
          <span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span
          ><button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button>
        </div>
      </header>
      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between">
            <h2 class="text-[13px] font-semibold text-[#1f2937]">Survey Submissions</h2>
            <button @click="openSubmitModal" class="btn-primary text-[11px]">New Submission</button>
          </div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search by property, owner, or ID..."
              class="input-field max-w-md"
            />
            <select v-model="filterStatus" class="input-field w-40">
              <option value="all">All Status</option>
              <option value="Approved">Approved</option>
              <option value="Under Review">Under Review</option>
              <option value="Pending">Pending</option>
            </select>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr>
                  <th class="table-header">Submission ID</th>
                  <th class="table-header">Survey ID</th>
                  <th class="table-header">Property</th>
                  <th class="table-header">Owner</th>
                  <th class="table-header">Photos</th>
                  <th class="table-header">Status</th>
                  <th class="table-header">Date</th>
                  <th class="table-header">Actions</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="sub in paginatedSubmissions" :key="sub.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ sub.id }}</td>
                  <td class="table-cell text-[#6b7280]">{{ sub.surveyId }}</td>
                  <td class="table-cell">{{ sub.property }}</td>
                  <td class="table-cell text-[11px]">{{ sub.owner }}</td>
                  <td class="table-cell text-[#9ca3af]">{{ sub.photos.length }} photos</td>
                  <td class="table-cell">
                    <span
                      class="px-2 py-0.5 text-[11px] font-medium rounded-full"
                      :class="{
                        'bg-green-100 text-green-700': sub.status === 'Approved',
                        'bg-yellow-100 text-yellow-700': sub.status === 'Under Review',
                        'bg-gray-100 text-gray-600': sub.status === 'Pending',
                      }"
                      >{{ sub.status }}</span
                    >
                  </td>
                  <td class="table-cell text-[#9ca3af]">{{ sub.submittedDate }}</td>
                  <td class="table-cell">
                    <button
                      @click="openViewModal(sub)"
                      class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]"
                    >
                      View
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="px-6 py-4 flex items-center justify-between border-t border-[#f3f4f6]">
            <p class="text-[11px] text-[#6b7280]">
              Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to
              {{ Math.min(currentPage * itemsPerPage, filteredSubmissions.length) }} of
              {{ filteredSubmissions.length }} entries
            </p>
            <div class="flex items-center gap-1">
              <button
                @click="goToPage(currentPage - 1)"
                :disabled="currentPage === 1"
                class="px-3 py-1 text-[11px] border border-[#e5e7eb] rounded disabled:opacity-50"
              >
                Prev
              </button>
              <button
                v-for="p in totalPages"
                :key="p"
                @click="goToPage(p)"
                class="px-3 py-1 text-[11px] border border-[#e5e7eb] rounded"
                :class="currentPage === p ? 'bg-[#1f2937] text-white' : ''"
              >
                {{ p }}
              </button>
              <button
                @click="goToPage(currentPage + 1)"
                :disabled="currentPage === totalPages"
                class="px-3 py-1 text-[11px] border border-[#e5e7eb] rounded disabled:opacity-50"
              >
                Next
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>

    <Teleport to="body"
      ><div
        v-if="toast.show"
        class="fixed bottom-4 right-4 bg-green-600 text-white px-4 py-2 rounded-lg shadow-lg z-50"
      >
        {{ toast.message }}
      </div></Teleport
    >

    <Teleport to="body">
      <div
        v-if="showSubmitModal"
        class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4 overflow-y-auto"
      >
        <div class="bg-white rounded-xl shadow-xl w-full max-w-3xl my-8">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Submit Survey Report</h3>
            <button @click="showSubmitModal = false" class="text-white/80 hover:text-white">
              ✕
            </button>
          </div>
          <div class="p-6 space-y-6">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5"
                  >Survey ID *</label
                >
                <input
                  v-model="newSubmission.surveyId"
                  type="text"
                  placeholder="e.g., SURV-001"
                  class="input-field w-full"
                />
              </div>
              <div>
                <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5"
                  >Property Location *</label
                >
                <input
                  v-model="newSubmission.property"
                  type="text"
                  placeholder="Enter property address"
                  class="input-field w-full"
                />
              </div>
              <div>
                <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5"
                  >Owner Name</label
                >
                <input
                  v-model="newSubmission.owner"
                  type="text"
                  placeholder="Enter owner name"
                  class="input-field w-full"
                />
              </div>
              <div>
                <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5"
                  >GPS Coordinates</label
                >
                <input
                  v-model="newSubmission.coordinates"
                  type="text"
                  placeholder="e.g., 6.4281° N, 3.4219° E"
                  class="input-field w-full"
                />
              </div>
            </div>

            <div class="bg-[#EEEEEE] rounded-lg p-4">
              <h4 class="text-[11px] text-[#6b7280] mb-3 font-semibold">PHYSICAL INSPECTION</h4>
              <div class="space-y-4">
                <div class="border border-[#d1d5db] rounded-lg p-4 bg-white">
                  <div class="flex items-start gap-3">
                    <input
                      type="checkbox"
                      v-model="newSubmission.gpsVerified"
                      class="w-5 h-5 mt-0.5 rounded border-[#d1d5db]"
                    />
                    <div class="flex-1">
                      <p class="text-[13px] font-semibold text-[#1f2937] mb-2">
                        📍 GPS Verification
                      </p>
                      <div class="flex items-center gap-2">
                        <input
                          v-model="newSubmission.coordinates"
                          type="text"
                          placeholder="Enter coordinates"
                          class="flex-1 px-3 py-2 border border-[#e5e7eb] rounded text-[12px]"
                        />
                        <button
                          v-if="newSubmission.coordinates"
                          @click="newSubmission.gpsVerified = true"
                          class="px-3 py-2 text-[11px] bg-green-600 text-white rounded hover:bg-green-700"
                        >
                          Verified
                        </button>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="border border-[#d1d5db] rounded-lg p-4 bg-white">
                  <p class="text-[13px] font-semibold text-[#1f2937] mb-3">📷 Property Photos *</p>
                  <label class="block mb-3">
                    <span
                      class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded cursor-pointer hover:bg-[#1e3d1a] inline-block"
                    >
                      Upload Photos
                    </span>
                    <input
                      type="file"
                      multiple
                      accept="image/*"
                      @change="handlePhotoUpload"
                      class="hidden"
                    />
                  </label>
                  <div v-if="newSubmission.photos.length > 0" class="flex flex-wrap gap-2 mb-2">
                    <span
                      v-for="(file, idx) in newSubmission.photos"
                      :key="idx"
                      class="px-3 py-1.5 bg-green-100 text-green-700 text-[11px] rounded flex items-center gap-2"
                    >
                      {{ file }}
                      <button
                        @click="removePhoto(idx)"
                        class="text-red-500 hover:text-red-700 font-bold"
                      >
                        ×
                      </button>
                    </span>
                  </div>
                  <p v-else class="text-[11px] text-[#9ca3af]">No photos uploaded yet</p>
                </div>

                <div class="border border-[#d1d5db] rounded-lg p-4 bg-white">
                  <p class="text-[13px] font-semibold text-[#1f2937] mb-3">
                    📐 Physical Measurements *
                  </p>
                  <div class="flex items-end gap-3">
                    <div>
                      <label class="text-[10px] text-[#6b7280]">Length (m)</label>
                      <input
                        v-model="newSubmission.measurements.length"
                        type="number"
                        placeholder="0"
                        class="w-24 px-3 py-2 border border-[#e5e7eb] rounded text-[12px]"
                      />
                    </div>
                    <span class="text-[#6b7280] pb-2">×</span>
                    <div>
                      <label class="text-[10px] text-[#6b7280]">Width (m)</label>
                      <input
                        v-model="newSubmission.measurements.width"
                        type="number"
                        placeholder="0"
                        class="w-24 px-3 py-2 border border-[#e5e7eb] rounded text-[12px]"
                      />
                    </div>
                    <button
                      @click="calculateArea"
                      class="px-4 py-2 text-[11px] bg-blue-600 text-white rounded hover:bg-blue-700"
                    >
                      Calculate
                    </button>
                    <span
                      v-if="newSubmission.measurements.total"
                      class="px-3 py-2 bg-green-100 text-green-700 text-[12px] rounded font-medium"
                      >= {{ newSubmission.measurements.total }}</span
                    >
                  </div>
                </div>

                <div class="border border-[#d1d5db] rounded-lg p-4 bg-white">
                  <p class="text-[13px] font-semibold text-[#1f2937] mb-3">
                    🏠 Condition Assessment *
                  </p>
                  <div class="flex flex-wrap gap-2">
                    <label
                      v-for="rating in ['Excellent', 'Good', 'Fair', 'Poor']"
                      :key="rating"
                      class="cursor-pointer"
                    >
                      <input
                        type="radio"
                        v-model="newSubmission.condition"
                        :value="rating"
                        class="hidden peer"
                      />
                      <span
                        class="px-4 py-2 text-[11px] border border-[#d1d5db] rounded peer-checked:bg-[#2D5A27] peer-checked:text-white peer-checked:border-[#2D5A27]"
                        >{{ rating }}</span
                      >
                    </label>
                  </div>
                </div>

                <div class="border border-[#d1d5db] rounded-lg p-4 bg-white">
                  <p class="text-[13px] font-semibold text-[#1f2937] mb-3">👥 Occupancy Status *</p>
                  <div class="flex flex-wrap gap-2">
                    <label
                      v-for="status in [
                        'Occupied',
                        'Vacant',
                        'Partially Occupied',
                        'Under Construction',
                      ]"
                      :key="status"
                      class="cursor-pointer"
                    >
                      <input
                        type="radio"
                        v-model="newSubmission.occupancy"
                        :value="status"
                        class="hidden peer"
                      />
                      <span
                        class="px-4 py-2 text-[11px] border border-[#d1d5db] rounded peer-checked:bg-[#2D5A27] peer-checked:text-white peer-checked:border-[#2D5A27]"
                        >{{ status }}</span
                      >
                    </label>
                  </div>
                </div>

                <div class="border border-[#d1d5db] rounded-lg p-4 bg-white">
                  <div class="flex items-start gap-3">
                    <input
                      type="checkbox"
                      v-model="newSubmission.ownershipVerified"
                      class="w-5 h-5 mt-0.5 rounded border-[#d1d5db]"
                    />
                    <div class="flex-1">
                      <p class="text-[13px] font-semibold text-[#1f2937]">✓ Ownership Verified</p>
                      <p class="text-[11px] text-[#6b7280]">
                        Confirm that the current occupant matches the owner on record
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div>
              <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5"
                >Survey Findings & Notes *</label
              >
              <textarea
                v-model="newSubmission.findings"
                rows="4"
                placeholder="Describe your findings from the physical inspection..."
                class="w-full px-4 py-3 border border-[#e5e7eb] rounded-lg text-[13px] resize-none focus:ring-2 focus:ring-[#2D5A27] focus:border-transparent"
              ></textarea>
            </div>

            <div class="flex items-center justify-between pt-4 border-t border-[#e5e7eb]">
              <div class="text-[11px]" :class="isFormValid ? 'text-green-600' : 'text-red-500'">
                {{
                  isFormValid
                    ? '✓ All required fields completed'
                    : '○ Complete all required fields marked with *'
                }}
              </div>
              <div class="flex gap-3">
                <button
                  @click="showSubmitModal = false"
                  class="px-4 py-2 text-[11px] border border-[#e5e7eb] text-[#6b7280] rounded-lg hover:bg-gray-50"
                >
                  Cancel
                </button>
                <button
                  @click="handleSubmit"
                  class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]"
                  :disabled="!isFormValid"
                >
                  Submit Report
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div
        v-if="showViewModal"
        class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4 overflow-y-auto"
      >
        <div class="bg-white rounded-xl shadow-xl w-full max-w-2xl my-8">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Submission Details</h3>
            <button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-6">
            <div class="grid grid-cols-2 gap-4">
              <div class="bg-[#EEEEEE] rounded-lg p-4">
                <h4 class="text-[11px] text-[#6b7280] mb-3 font-semibold">SUBMISSION INFO</h4>
                <div class="space-y-2">
                  <div>
                    <p class="text-[10px] text-[#9ca3af]">Submission ID</p>
                    <p class="text-[13px] font-semibold text-[#1f2937]">
                      {{ selectedSubmission?.id }}
                    </p>
                  </div>
                  <div>
                    <p class="text-[10px] text-[#9ca3af]">Survey ID</p>
                    <p class="text-[13px] font-semibold text-[#1f2937]">
                      {{ selectedSubmission?.surveyId }}
                    </p>
                  </div>
                  <div>
                    <p class="text-[10px] text-[#9ca3af]">Status</p>
                    <span
                      class="px-2 py-0.5 text-[11px] font-medium rounded-full"
                      :class="{
                        'bg-green-100 text-green-700': selectedSubmission?.status === 'Approved',
                        'bg-yellow-100 text-yellow-700':
                          selectedSubmission?.status === 'Under Review',
                        'bg-gray-100 text-gray-600': selectedSubmission?.status === 'Pending',
                      }"
                      >{{ selectedSubmission?.status }}</span
                    >
                  </div>
                </div>
              </div>
              <div class="bg-[#EEEEEE] rounded-lg p-4">
                <h4 class="text-[11px] text-[#6b7280] mb-3 font-semibold">PROPERTY INFO</h4>
                <div class="space-y-2">
                  <div>
                    <p class="text-[10px] text-[#9ca3af]">Property</p>
                    <p class="text-[13px] font-semibold text-[#1f2937]">
                      {{ selectedSubmission?.property }}
                    </p>
                  </div>
                  <div>
                    <p class="text-[10px] text-[#9ca3af]">Owner</p>
                    <p class="text-[13px] font-semibold text-[#1f2937]">
                      {{ selectedSubmission?.owner }}
                    </p>
                  </div>
                  <div>
                    <p class="text-[10px] text-[#9ca3af]">Submitted</p>
                    <p class="text-[13px] font-semibold text-[#1f2937]">
                      {{ selectedSubmission?.submittedDate }}
                    </p>
                  </div>
                </div>
              </div>
            </div>

            <div class="bg-[#EEEEEE] rounded-lg p-4">
              <h4 class="text-[11px] text-[#6b7280] mb-3 font-semibold">INSPECTION RESULTS</h4>
              <div class="grid grid-cols-3 gap-4">
                <div>
                  <p class="text-[10px] text-[#9ca3af]">Photos</p>
                  <p class="text-[13px] font-semibold text-[#1f2937]">
                    {{ selectedSubmission?.photos?.length || 0 }} uploaded
                  </p>
                </div>
                <div>
                  <p class="text-[10px] text-[#9ca3af]">Measurements</p>
                  <p class="text-[13px] font-semibold text-[#1f2937]">
                    {{ selectedSubmission?.measurements || 'N/A' }}
                  </p>
                </div>
                <div>
                  <p class="text-[10px] text-[#9ca3af]">Condition</p>
                  <span
                    class="px-2 py-0.5 text-[11px] font-medium rounded-full bg-blue-100 text-blue-700"
                    >{{ selectedSubmission?.condition || 'N/A' }}</span
                  >
                </div>
                <div>
                  <p class="text-[10px] text-[#9ca3af]">Occupancy</p>
                  <span
                    class="px-2 py-0.5 text-[11px] font-medium rounded-full bg-blue-100 text-blue-700"
                    >{{ selectedSubmission?.occupancy || 'N/A' }}</span
                  >
                </div>
              </div>
            </div>

            <div class="flex items-center justify-between pt-4 border-t border-[#e5e7eb]">
              <div class="text-[11px] text-[#6b7280]">
                Verified by: {{ selectedSubmission?.verifiedBy }}
              </div>
              <button
                @click="showViewModal = false"
                class="px-4 py-2 text-[11px] bg-[#f3f4f6] text-[#374151] rounded-lg hover:bg-[#e5e7eb]"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
