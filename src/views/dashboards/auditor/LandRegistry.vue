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

const showAddModal = ref(false)
const showViewModal = ref(false)
const showEditModal = ref(false)
const selectedRecord = ref<any>(null)

const newRecord = ref({
  plotNo: '',
  block: '',
  location: '',
  size: '',
  owner: '',
  titleType: 'C of O',
})

const landRecords = ref([
  { id: 'LR-001', plotNo: 'Plot 42', block: 'Block A', location: 'Victoria Island', size: '500 sqm', owner: 'Emeka Okonkwo', titleType: 'C of O', status: 'Active', lastUpdated: '2024-01-10' },
  { id: 'LR-002', plotNo: 'Plot 15', block: 'Block C', location: 'Lekki Phase 1', size: '750 sqm', owner: 'Adaobi Nnamdi', titleType: 'Survey Plan', status: 'Pending', lastUpdated: '2024-01-12' },
  { id: 'LR-003', plotNo: 'Plot 88', block: 'Block B', location: 'Ikoyi', size: '1,200 sqm', owner: 'Chidi Okafor', titleType: 'C of O', status: 'Active', lastUpdated: '2024-01-08' },
  { id: 'LR-004', plotNo: 'Plot 8', block: 'Block D', location: 'Banana Island', size: '600 sqm', owner: 'Folake Adeyemi', titleType: 'Excision', status: 'Flagged', lastUpdated: '2024-01-15' },
])

const filteredRecords = computed(() => {
  return landRecords.value.filter(r => {
    const matchesSearch = r.owner.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      r.location.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      r.id.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesStatus = filterStatus.value === 'all' || r.status === filterStatus.value
    return matchesSearch && matchesStatus
  })
})

const totalPages = computed(() => Math.ceil(filteredRecords.value.length / itemsPerPage.value))

const paginatedRecords = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  return filteredRecords.value.slice(start, start + itemsPerPage.value)
})

const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
  }
}

const handleView = (record: any) => {
  selectedRecord.value = record
  showViewModal.value = true
}

const handleEdit = (record: any) => {
  selectedRecord.value = { ...record }
  showEditModal.value = true
}

const handleAddRecord = () => {
  const newId = 'LR-' + String(landRecords.value.length + 1).padStart(3, '0')
  const today = new Date().toISOString().split('T')[0]
  landRecords.value.push({
    id: newId,
    ...newRecord.value,
    status: 'Pending',
    lastUpdated: today,
  })
  showAddModal.value = false
  newRecord.value = { plotNo: '', block: '', location: '', size: '', owner: '', titleType: 'C of O' }
}

const handleUpdateRecord = () => {
  const index = landRecords.value.findIndex(r => r.id === selectedRecord.value.id)
  if (index !== -1) {
    landRecords.value[index] = {
      ...selectedRecord.value,
      lastUpdated: new Date().toISOString().split('T')[0],
    }
  }
  showEditModal.value = false
}

const closeModal = () => {
  showAddModal.value = false
  showViewModal.value = false
  showEditModal.value = false
  selectedRecord.value = null
}
</script>

<template>
  <div class="min-h-screen flex bg-[#f5f7fa]">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />

    <div class="flex-1 flex flex-col">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6 shrink-0">
        <div class="flex items-center gap-4">
          <span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span>
          <span class="text-[#d1d5db]">/</span>
          <span class="text-[#1f2937] text-sm font-medium">Land Registry</span>
        </div>
        <div class="flex items-center gap-4">
          <span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span>
          <button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button>
        </div>
      </header>

      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between">
            <h2 class="text-[13px] font-semibold text-[#1f2937]">Land Registry</h2>
            <button @click="showAddModal = true" class="btn-primary text-[11px]">Add Record</button>
          </div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4">
            <input v-model="searchQuery" type="text" placeholder="Search..." class="input-field max-w-md" />
            <select v-model="filterStatus" class="input-field w-48">
              <option value="all">All Status</option>
              <option value="Active">Active</option>
              <option value="Pending">Pending</option>
              <option value="Flagged">Flagged</option>
            </select>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr>
                  <th class="table-header">Record ID</th>
                  <th class="table-header">Plot No</th>
                  <th class="table-header">Block</th>
                  <th class="table-header">Location</th>
                  <th class="table-header">Size</th>
                  <th class="table-header">Owner</th>
                  <th class="table-header">Title</th>
                  <th class="table-header">Status</th>
                  <th class="table-header">Updated</th>
                  <th class="table-header">Actions</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="record in paginatedRecords" :key="record.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ record.id }}</td>
                  <td class="table-cell">{{ record.plotNo }}</td>
                  <td class="table-cell">{{ record.block }}</td>
                  <td class="table-cell text-[#6b7280]">{{ record.location }}</td>
                  <td class="table-cell">{{ record.size }}</td>
                  <td class="table-cell">{{ record.owner }}</td>
                  <td class="table-cell">
                    <span class="px-2 py-0.5 text-[11px] bg-[#f3f4f6] text-[#6b7280] rounded">{{ record.titleType }}</span>
                  </td>
                  <td class="table-cell">
                    <span class="px-2 py-0.5 text-[11px] font-medium rounded-full"
                      :class="{
                        'bg-green-50 text-green-700': record.status === 'Active',
                        'bg-yellow-50 text-yellow-700': record.status === 'Pending',
                        'bg-red-50 text-red-700': record.status === 'Flagged',
                      }">{{ record.status }}</span>
                  </td>
                  <td class="table-cell text-[#9ca3af]">{{ record.lastUpdated }}</td>
                  <td class="table-cell">
                    <div class="flex gap-2">
                      <button @click="handleView(record)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">View</button>
                      <button @click="handleEdit(record)" class="px-3 py-1 text-[11px] bg-[#2D5A27] text-white rounded hover:bg-[#1e3d1a]">Edit</button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="px-6 py-4 flex items-center justify-between border-t border-[#f3f4f6]">
            <p class="text-[11px] text-[#6b7280]">Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, filteredRecords.length) }} of {{ filteredRecords.length }} entries</p>
            <div class="flex items-center gap-1">
              <button @click="goToPage(currentPage - 1)" :disabled="currentPage === 1" class="px-3 py-1 text-[11px] border border-[#e5e7eb] rounded hover:bg-[#f9fafb] disabled:opacity-50 disabled:cursor-not-allowed">Prev</button>
              <button v-for="p in totalPages" :key="p" @click="goToPage(p)" class="px-3 py-1 text-[11px] border border-[#e5e7eb] rounded hover:bg-[#f9fafb]" :class="currentPage === p ? 'bg-[#1f2937] text-white border-[#1f2937]' : ''">{{ p }}</button>
              <button @click="goToPage(currentPage + 1)" :disabled="currentPage === totalPages" class="px-3 py-1 text-[11px] border border-[#e5e7eb] rounded hover:bg-[#f9fafb] disabled:opacity-50 disabled:cursor-not-allowed">Next</button>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>

  <Teleport to="body">
    <div v-if="showAddModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-xl shadow-xl w-full max-w-lg overflow-hidden max-h-[90vh] overflow-y-auto">
        <div class="bg-[#2D5A27] px-6 py-4">
          <h3 class="text-base font-semibold text-white">Add New Record</h3>
        </div>
        <div class="p-6">
          <div class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Plot Number</label>
                <input v-model="newRecord.plotNo" type="text" class="input-field" placeholder="e.g. Plot 42" />
              </div>
              <div>
                <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Block</label>
                <input v-model="newRecord.block" type="text" class="input-field" placeholder="e.g. Block A" />
              </div>
            </div>
            <div>
              <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Location</label>
              <input v-model="newRecord.location" type="text" class="input-field" placeholder="e.g. Victoria Island" />
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Size</label>
                <input v-model="newRecord.size" type="text" class="input-field" placeholder="e.g. 500 sqm" />
              </div>
              <div>
                <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Title Type</label>
                <select v-model="newRecord.titleType" class="input-field">
                  <option value="C of O">C of O</option>
                  <option value="Survey Plan">Survey Plan</option>
                  <option value="Excision">Excision</option>
                  <option value="Lease">Lease</option>
                </select>
              </div>
            </div>
            <div>
              <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Owner Name</label>
              <input v-model="newRecord.owner" type="text" class="input-field" placeholder="e.g. Emeka Okonkwo" />
            </div>
            <div>
              <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Upload Documents</label>
              <div class="border-2 border-dashed border-[#d1d5db] rounded-lg p-6 text-center hover:border-[#2D5A27] transition cursor-pointer">
                <input type="file" multiple class="hidden" id="addFileUpload" />
                <label for="addFileUpload" class="cursor-pointer">
                  <div class="text-[#6b7280]">
                    <p class="text-sm font-medium">Click to upload or drag and drop</p>
                    <p class="text-[11px] mt-1">PDF, JPG, PNG up to 10MB</p>
                  </div>
                </label>
              </div>
            </div>
          </div>
          <div class="flex gap-3 mt-6">
            <button @click="closeModal" class="btn-secondary flex-1">Cancel</button>
            <button @click="handleAddRecord" class="btn-primary flex-1">Add Record</button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showViewModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-xl shadow-xl w-full max-w-lg overflow-hidden">
        <div class="bg-[#2D5A27] px-6 py-4">
          <h3 class="text-base font-semibold text-white">View Record Details</h3>
        </div>
        <div class="p-6">
          <div class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-[11px] font-medium text-[#9ca3af] mb-1">Record ID</label>
                <p class="text-sm text-[#1f2937] font-medium">{{ selectedRecord?.id }}</p>
              </div>
              <div>
                <label class="block text-[11px] font-medium text-[#9ca3af] mb-1">Status</label>
                <span class="px-2 py-0.5 text-[11px] font-medium rounded-full"
                  :class="{
                    'bg-green-50 text-green-700': selectedRecord?.status === 'Active',
                    'bg-yellow-50 text-yellow-700': selectedRecord?.status === 'Pending',
                    'bg-red-50 text-red-700': selectedRecord?.status === 'Flagged',
                  }">{{ selectedRecord?.status }}</span>
              </div>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-[11px] font-medium text-[#9ca3af] mb-1">Plot Number</label>
                <p class="text-sm text-[#1f2937]">{{ selectedRecord?.plotNo }}</p>
              </div>
              <div>
                <label class="block text-[11px] font-medium text-[#9ca3af] mb-1">Block</label>
                <p class="text-sm text-[#1f2937]">{{ selectedRecord?.block }}</p>
              </div>
            </div>
            <div>
              <label class="block text-[11px] font-medium text-[#9ca3af] mb-1">Location</label>
              <p class="text-sm text-[#1f2937]">{{ selectedRecord?.location }}</p>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-[11px] font-medium text-[#9ca3af] mb-1">Size</label>
                <p class="text-sm text-[#1f2937]">{{ selectedRecord?.size }}</p>
              </div>
              <div>
                <label class="block text-[11px] font-medium text-[#9ca3af] mb-1">Title Type</label>
                <p class="text-sm text-[#1f2937]">{{ selectedRecord?.titleType }}</p>
              </div>
            </div>
            <div>
              <label class="block text-[11px] font-medium text-[#9ca3af] mb-1">Owner</label>
              <p class="text-sm text-[#1f2937]">{{ selectedRecord?.owner }}</p>
            </div>
            <div>
              <label class="block text-[11px] font-medium text-[#9ca3af] mb-1">Last Updated</label>
              <p class="text-sm text-[#1f2937]">{{ selectedRecord?.lastUpdated }}</p>
            </div>
            <div>
              <label class="block text-[11px] font-medium text-[#9ca3af] mb-1">Documents</label>
              <div class="border border-[#e5e7eb] rounded-lg p-4 mt-1">
                <p class="text-sm text-[#6b7280]">No documents uploaded</p>
              </div>
            </div>
          </div>
          <div class="mt-6">
            <button @click="closeModal" class="btn-secondary w-full">Close</button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showEditModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-xl shadow-xl w-full max-w-lg overflow-hidden max-h-[90vh] overflow-y-auto">
        <div class="bg-[#2D5A27] px-6 py-4">
          <h3 class="text-base font-semibold text-white">Edit Record</h3>
        </div>
        <div class="p-6">
          <div class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Plot Number</label>
                <input v-model="selectedRecord.plotNo" type="text" class="input-field" />
              </div>
              <div>
                <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Block</label>
                <input v-model="selectedRecord.block" type="text" class="input-field" />
              </div>
            </div>
            <div>
              <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Location</label>
              <input v-model="selectedRecord.location" type="text" class="input-field" />
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Size</label>
                <input v-model="selectedRecord.size" type="text" class="input-field" />
              </div>
              <div>
                <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Title Type</label>
                <select v-model="selectedRecord.titleType" class="input-field">
                  <option value="C of O">C of O</option>
                  <option value="Survey Plan">Survey Plan</option>
                  <option value="Excision">Excision</option>
                  <option value="Lease">Lease</option>
                </select>
              </div>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Owner</label>
                <input v-model="selectedRecord.owner" type="text" class="input-field" />
              </div>
              <div>
                <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Status</label>
                <select v-model="selectedRecord.status" class="input-field">
                  <option value="Active">Active</option>
                  <option value="Pending">Pending</option>
                  <option value="Flagged">Flagged</option>
                </select>
              </div>
            </div>
            <div>
              <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Upload Documents</label>
              <div class="border-2 border-dashed border-[#d1d5db] rounded-lg p-6 text-center hover:border-[#2D5A27] transition cursor-pointer">
                <input type="file" multiple class="hidden" id="editFileUpload" />
                <label for="editFileUpload" class="cursor-pointer">
                  <div class="text-[#6b7280]">
                    <p class="text-sm font-medium">Click to upload or drag and drop</p>
                    <p class="text-[11px] mt-1">PDF, JPG, PNG up to 10MB</p>
                  </div>
                </label>
              </div>
            </div>
          </div>
          <div class="flex gap-3 mt-6">
            <button @click="closeModal" class="btn-secondary flex-1">Cancel</button>
            <button @click="handleUpdateRecord" class="btn-primary flex-1">Save Changes</button>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>
