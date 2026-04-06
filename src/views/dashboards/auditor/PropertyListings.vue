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
const filterRisk = ref('all')
const currentPage = ref(1)
const itemsPerPage = ref(5)

const showViewModal = ref(false)
const showAuditModal = ref(false)
const selectedProperty = ref<any>(null)
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })

const newAudit = ref({ priority: 'Medium', dueDate: '', notes: '' })

const auditCases = ref([
  { id: 'AUD-2024-001', property: 'Plot 42, Victoria Island', owner: 'Emeka Okonkwo', auditor: 'John Smith', priority: 'High', status: 'In Progress', started: '2024-01-10', due: '2024-01-25' },
  { id: 'AUD-2024-002', property: 'Block 7, Lekki Phase 2', owner: 'Adaobi Nnamdi', auditor: 'Sarah Johnson', priority: 'Medium', status: 'Pending', started: '2024-01-12', due: '2024-01-30' },
  { id: 'AUD-2024-003', property: '15 Admiralty Way, Lekki', owner: 'Chidi Okafor', auditor: 'John Smith', priority: 'Low', status: 'Completed', started: '2024-01-05', due: '2024-01-20' },
])

const properties = ref([
  { id: 'PROP-001', owner: 'Emeka Okonkwo', address: 'Plot 42, Victoria Island', type: 'Commercial', value: 'N250,000,000', status: 'Verified', riskScore: 25, riskLevel: 'Low', lastAudit: '2024-01-10', coordinates: '6.4281° N, 3.4219° E', declaredRent: 'N12,000,000' },
  { id: 'PROP-002', owner: 'Adaobi Nnamdi', address: 'Block 7, Lekki Phase 2', type: 'Residential', value: 'N80,000,000', status: 'Pending', riskScore: 35, riskLevel: 'Medium', lastAudit: '2024-01-12', coordinates: '6.4312° N, 3.5012° E', declaredRent: 'N4,500,000' },
  { id: 'PROP-003', owner: 'Chidi Okafor', address: '15 Admiralty Way, Lekki', type: 'Mixed Use', value: 'N180,000,000', status: 'Verified', riskScore: 20, riskLevel: 'Low', lastAudit: '2024-01-08', coordinates: '6.4281° N, 3.4219° E', declaredRent: 'N9,000,000' },
  { id: 'PROP-004', owner: 'Folake Adeyemi', address: 'Plot 8, Banana Island', type: 'Residential', value: 'N500,000,000', status: 'Flagged', riskScore: 85, riskLevel: 'Critical', lastAudit: '2024-01-15', coordinates: '6.4281° N, 3.4219° E', declaredRent: 'N2,500,000' },
  { id: 'PROP-005', owner: 'Ibrahim Bello', address: 'Block 3, Ikoyi', type: 'Commercial', value: 'N350,000,000', status: 'Verified', riskScore: 30, riskLevel: 'Low', lastAudit: '2024-01-05', coordinates: '6.4536° N, 3.3958° E', declaredRent: 'N18,000,000' },
  { id: 'PROP-006', owner: 'Global Ventures Ltd', address: 'Estate 7, Lekki', type: 'Residential Estate', value: 'N620,000,000', status: 'Flagged', riskScore: 72, riskLevel: 'High', lastAudit: '2024-01-10', coordinates: '6.4312° N, 3.5012° E', declaredRent: 'N1,800,000' },
  { id: 'PROP-007', owner: 'Nigerian Holdings Ltd', address: '21 Broad Street, Lagos', type: 'Commercial', value: 'N890,000,000', status: 'Pending', riskScore: 45, riskLevel: 'Medium', lastAudit: '2024-01-18', coordinates: '6.4281° N, 3.4219° E', declaredRent: 'N45,000,000' },
])

const filteredProperties = computed(() => properties.value.filter(p => {
  const matchesSearch = p.owner.toLowerCase().includes(searchQuery.value.toLowerCase()) || p.address.toLowerCase().includes(searchQuery.value.toLowerCase()) || p.id.toLowerCase().includes(searchQuery.value.toLowerCase())
  const matchesStatus = filterStatus.value === 'all' || p.status === filterStatus.value
  const matchesRisk = filterRisk.value === 'all' || p.riskLevel === filterRisk.value
  return matchesSearch && matchesStatus && matchesRisk
}))

const totalPages = computed(() => Math.ceil(filteredProperties.value.length / itemsPerPage.value))
const paginatedProperties = computed(() => { const start = (currentPage.value - 1) * itemsPerPage.value; return filteredProperties.value.slice(start, start + itemsPerPage.value) })
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }

const showToast = (message: string) => { toast.value = { show: true, message }; setTimeout(() => { toast.value.show = false }, 3000) }

const openViewModal = (p: any) => {
  selectedProperty.value = p
  newAudit.value = { priority: 'Medium', dueDate: '', notes: '' }
  showViewModal.value = true
}

const openAuditModal = () => {
  showViewModal.value = false
  showAuditModal.value = true
}

const startAudit = () => {
  if (!selectedProperty.value) return
  const newId = 'AUD-2024-' + String(auditCases.value.length + 1).padStart(3, '0')
  auditCases.value.unshift({
    id: newId,
    property: selectedProperty.value.address,
    owner: selectedProperty.value.owner,
    auditor: user.value?.name || 'Auditor',
    priority: newAudit.value.priority,
    status: 'Pending',
    started: new Date().toISOString().split('T')[0],
    due: newAudit.value.dueDate
  })
  const propIndex = properties.value.findIndex(p => p.id === selectedProperty.value.id)
  if (propIndex !== -1) properties.value[propIndex].status = 'Under Audit'
  showAuditModal.value = false
  showToast('Audit case ' + newId + ' created')
}

const getRiskColor = (level: string) => {
  switch (level) {
    case 'Critical': return 'bg-red-100 text-red-700'
    case 'High': return 'bg-orange-100 text-orange-700'
    case 'Medium': return 'bg-yellow-100 text-yellow-700'
    case 'Low': return 'bg-green-100 text-green-700'
    default: return 'bg-gray-100 text-gray-700'
  }
}
</script>

<template>
  <div class="min-h-screen flex bg-[#f5f7fa]">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />
    <div class="flex-1 flex flex-col">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6 shrink-0">
        <div class="flex items-center gap-4"><span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span><span class="text-[#d1d5db]">/</span><span class="text-[#1f2937] text-sm font-medium">Property Listings</span></div>
        <div class="flex items-center gap-4"><span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span><button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button></div>
      </header>
      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between"><h2 class="text-[13px] font-semibold text-[#1f2937]">Properties for Review</h2><span class="text-[11px] text-[#6b7280]">Auditor: {{ user?.name || 'Auditor' }}</span></div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4 flex-wrap">
            <input v-model="searchQuery" type="text" placeholder="Search by owner, address, or ID..." class="input-field max-w-md" />
            <select v-model="filterStatus" class="input-field w-40"><option value="all">All Status</option><option value="Verified">Verified</option><option value="Pending">Pending</option><option value="Flagged">Flagged</option><option value="Under Audit">Under Audit</option></select>
            <select v-model="filterRisk" class="input-field w-40"><option value="all">All Risk</option><option value="Critical">Critical</option><option value="High">High</option><option value="Medium">Medium</option><option value="Low">Low</option></select>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead><tr><th class="table-header">Property ID</th><th class="table-header">Owner</th><th class="table-header">Address</th><th class="table-header">Type</th><th class="table-header">Value</th><th class="table-header">Risk</th><th class="table-header">Status</th><th class="table-header">Last Audit</th><th class="table-header">Actions</th></tr></thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="property in paginatedProperties" :key="property.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ property.id }}</td>
                  <td class="table-cell">{{ property.owner }}</td>
                  <td class="table-cell text-[#6b7280]">{{ property.address }}</td>
                  <td class="table-cell text-[11px]">{{ property.type }}</td>
                  <td class="table-cell">{{ property.value }}</td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="getRiskColor(property.riskLevel)">{{ property.riskLevel }} ({{ property.riskScore }})</span></td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': property.status === 'Verified', 'bg-yellow-50 text-yellow-700': property.status === 'Pending', 'bg-red-50 text-red-700': property.status === 'Flagged', 'bg-blue-50 text-blue-700': property.status === 'Under Audit'}">{{ property.status }}</span></td>
                  <td class="table-cell text-[#9ca3af]">{{ property.lastAudit }}</td>
                  <td class="table-cell"><button @click="openViewModal(property)" class="px-3 py-1 text-[11px] bg-[#2D5A27] text-white rounded hover:bg-[#1e3d1a]">Review</button></td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="px-6 py-4 flex items-center justify-between border-t border-[#f3f4f6]">
            <p class="text-[11px] text-[#6b7280]">Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, filteredProperties.length) }} of {{ filteredProperties.length }} entries</p>
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
            <h3 class="text-base font-semibold text-white">Property Review - {{ selectedProperty?.id }}</h3>
            <button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-6">
            <div class="grid grid-cols-2 gap-4">
              <div class="bg-[#EEEEEE] rounded-lg p-4">
                <h4 class="text-[11px] text-[#6b7280] mb-3 font-semibold">PROPERTY INFO</h4>
                <div class="space-y-3">
                  <div><p class="text-[10px] text-[#9ca3af]">Property ID</p><p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedProperty?.id }}</p></div>
                  <div><p class="text-[10px] text-[#9ca3af]">Owner</p><p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedProperty?.owner }}</p></div>
                  <div><p class="text-[10px] text-[#9ca3af]">Address</p><p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedProperty?.address }}</p></div>
                  <div><p class="text-[10px] text-[#9ca3af]">Type</p><p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedProperty?.type }}</p></div>
                </div>
              </div>
              <div class="bg-[#EEEEEE] rounded-lg p-4">
                <h4 class="text-[11px] text-[#6b7280] mb-3 font-semibold">RISK ASSESSMENT</h4>
                <div class="space-y-3">
                  <div><p class="text-[10px] text-[#9ca3af]">Risk Level</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="getRiskColor(selectedProperty?.riskLevel)">{{ selectedProperty?.riskLevel }}</span></div>
                  <div><p class="text-[10px] text-[#9ca3af]">Risk Score</p><p class="text-[15px] font-bold text-[#1f2937]">{{ selectedProperty?.riskScore }}/100</p></div>
                  <div><p class="text-[10px] text-[#9ca3af]">Status</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': selectedProperty?.status === 'Verified', 'bg-yellow-50 text-yellow-700': selectedProperty?.status === 'Pending', 'bg-red-50 text-red-700': selectedProperty?.status === 'Flagged'}">{{ selectedProperty?.status }}</span></div>
                </div>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div class="bg-[#EEEEEE] rounded-lg p-4">
                <h4 class="text-[11px] text-[#6b7280] mb-3 font-semibold">VALUATION</h4>
                <div class="space-y-3">
                  <div><p class="text-[10px] text-[#9ca3af]">Declared Value</p><p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedProperty?.value }}</p></div>
                  <div><p class="text-[10px] text-[#9ca3af]">Declared Rent</p><p class="text-[13px] font-semibold text-red-600">{{ selectedProperty?.declaredRent }}/year</p></div>
                </div>
              </div>
              <div class="bg-[#EEEEEE] rounded-lg p-4">
                <h4 class="text-[11px] text-[#6b7280] mb-3 font-semibold">LOCATION</h4>
                <div class="space-y-3">
                  <div><p class="text-[10px] text-[#9ca3af]">GPS Coordinates</p><p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedProperty?.coordinates }}</p></div>
                  <div><p class="text-[10px] text-[#9ca3af]">Last Audit</p><p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedProperty?.lastAudit }}</p></div>
                </div>
              </div>
            </div>

            <div class="border-t border-[#e5e7eb] pt-4">
              <h4 class="text-[11px] text-[#6b7280] mb-3 font-semibold">AUDITOR ACTIONS</h4>
              <div class="flex items-center justify-between">
                <div class="text-[11px] text-[#6b7280]">
                  {{ selectedProperty?.status === 'Flagged' ? 'This property has been flagged for investigation.' : selectedProperty?.status === 'Verified' ? 'This property appears compliant.' : 'This property requires review.' }}
                </div>
                <button v-if="selectedProperty?.status === 'Flagged' || selectedProperty?.status === 'Pending'" @click="openAuditModal" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">
                  Create Audit Case
                </button>
                <button v-else class="px-4 py-2 text-[11px] bg-gray-200 text-gray-500 rounded-lg cursor-not-allowed">
                  No Action Required
                </button>
              </div>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-[#e5e7eb] flex justify-end">
            <button @click="showViewModal = false" class="px-4 py-2 text-[11px] bg-[#f3f4f6] text-[#374151] rounded-lg hover:bg-[#e5e7eb]">Close</button>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showAuditModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Create Audit Case</h3>
            <button @click="showAuditModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-4">
            <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-3">
              <p class="text-[11px] text-yellow-800">Creating audit for: <strong>{{ selectedProperty?.address }}</strong></p>
              <p class="text-[11px] text-yellow-800">Owner: <strong>{{ selectedProperty?.owner }}</strong></p>
              <p class="text-[11px] text-yellow-800">Risk: <strong>{{ selectedProperty?.riskLevel }} ({{ selectedProperty?.riskScore }})</strong></p>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Priority</label>
                <select v-model="newAudit.priority" class="input-field w-full">
                  <option>Low</option>
                  <option>Medium</option>
                  <option>High</option>
                  <option>Critical</option>
                </select>
              </div>
              <div>
                <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Due Date</label>
                <input v-model="newAudit.dueDate" type="date" class="input-field w-full" />
              </div>
            </div>
            <div>
              <label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Notes</label>
              <textarea v-model="newAudit.notes" rows="3" placeholder="Add notes about this audit case..." class="input-field w-full resize-none"></textarea>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-end">
            <button @click="showAuditModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button>
            <button @click="startAudit" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Create Case</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
