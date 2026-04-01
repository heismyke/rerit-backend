<script setup lang="ts">
import { useRoleStore } from '@/stores'
import { useRouter } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import { ref, computed } from 'vue'

const { selectedRole, user, logout } = useRoleStore()
const router = useRouter()
const handleLogout = () => { logout(); router.push('/') }

const searchQuery = ref('')
const filterPriority = ref('all')
const currentPage = ref(1)
const itemsPerPage = ref(5)

const showViewModal = ref(false)
const selectedFlag = ref<any>(null)
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })
const aiInsights = ref('')
const estimatedRent = ref('')

const flaggedProperties = ref([
  { id: 'FLAG-001', property: 'Plot 8, Banana Island', owner: 'Folake Adeyemi', reason: 'Undeclared renovation', priority: 'Critical', status: 'Under Investigation', flaggedDate: '2024-01-15', investigator: 'Agent A', occupantId: 'OCC-001', location: 'Banana Island, Lagos', riskLevel: 'Critical', taxType: 'Property Tax', houseType: 'Residential Mansion', declaredValue: 'N500,000,000', estimatedValue: 'N850,000,000', declaredRent: 'N2,500,000', lat: 6.4281, lng: 3.4219 },
  { id: 'FLAG-002', property: 'Block 12, Maitama', owner: 'Unknown Entity', reason: 'Suspected tax evasion', priority: 'High', status: 'Pending Review', flaggedDate: '2024-01-12', investigator: 'Agent B', occupantId: 'OCC-002', location: 'Maitama, Abuja', riskLevel: 'High', taxType: 'Land Use Charge', houseType: 'Commercial Complex', declaredValue: 'N200,000,000', estimatedValue: 'N450,000,000', declaredRent: 'N800,000', lat: 9.0579, lng: 7.4951 },
  { id: 'FLAG-003', property: 'Plot 45, VI', owner: 'Chinedu & Partners', reason: 'Value discrepancy', priority: 'Medium', status: 'Resolved', flaggedDate: '2024-01-08', investigator: 'Agent A', occupantId: 'OCC-003', location: 'Victoria Island, Lagos', riskLevel: 'Medium', taxType: 'Development Levy', houseType: 'Mixed Use', declaredValue: 'N180,000,000', estimatedValue: 'N220,000,000', declaredRent: 'N1,200,000', lat: 6.4281, lng: 3.4219 },
  { id: 'FLAG-004', property: 'Estate 7, Lekki', owner: 'Global Ventures Ltd', reason: 'Document forgery', priority: 'High', status: 'Under Investigation', flaggedDate: '2024-01-10', investigator: 'Agent C', occupantId: 'OCC-004', location: 'Lekki Phase 1, Lagos', riskLevel: 'High', taxType: 'Property Tax', houseType: 'Residential Estate', declaredValue: 'N350,000,000', estimatedValue: 'N620,000,000', declaredRent: 'N1,800,000', lat: 6.4312, lng: 3.5012 },
])

const filteredProperties = computed(() => flaggedProperties.value.filter(p => {
  const matchesSearch = p.property.toLowerCase().includes(searchQuery.value.toLowerCase()) || p.owner.toLowerCase().includes(searchQuery.value.toLowerCase()) || p.id.toLowerCase().includes(searchQuery.value.toLowerCase())
  const matchesPriority = filterPriority.value === 'all' || p.priority === filterPriority.value
  return matchesSearch && matchesPriority
}))

const totalPages = computed(() => Math.ceil(filteredProperties.value.length / itemsPerPage.value))
const paginatedProperties = computed(() => { const start = (currentPage.value - 1) * itemsPerPage.value; return filteredProperties.value.slice(start, start + itemsPerPage.value) })
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }

const showToast = (message: string) => { toast.value = { show: true, message }; setTimeout(() => { toast.value.show = false }, 3000) }
const openViewModal = (f: any) => { selectedFlag.value = f; aiInsights.value = ''; estimatedRent.value = ''; showViewModal.value = true }
const startInvestigation = (f: any) => { showToast('Investigation started for ' + f.id) }
const escalateAudit = () => { showToast('Audit escalated for ' + selectedFlag.value?.id); showViewModal.value = false }
const getComplianceNotice = () => { showToast('Compliance notice generated for ' + selectedFlag.value?.id) }

const calculateComplianceGap = computed(() => {
  if (!selectedFlag.value || !estimatedRent.value) return null
  const declared = parseFloat(selectedFlag.value.declaredRent.replace(/[^0-9.]/g, ''))
  const estimated = parseFloat(estimatedRent.value.replace(/[^0-9.]/g, ''))
  const gap = estimated - declared
  return gap
})

const hasTaxIssue = computed(() => {
  if (!calculateComplianceGap.value) return false
  return calculateComplianceGap.value > 0
})

const getAIRecommendations = () => {
  if (!selectedFlag.value || !estimatedRent.value.trim()) {
    showToast('Please enter the estimated rent to get AI insights')
    return
  }
  showToast('Analyzing rent data...')
  const gap = calculateComplianceGap.value ?? 0
  const riskScore = selectedFlag.value?.riskLevel === 'Critical' ? '85/100' : selectedFlag.value?.riskLevel === 'High' ? '68/100' : '45/100'
  const declaredRentNum = parseFloat(selectedFlag.value?.declaredRent?.replace(/[^0-9.]/g, '') || '0')
  const gapPercent = declaredRentNum > 0 ? Math.round((gap / declaredRentNum) * 100) : 0
  setTimeout(() => {
    aiInsights.value = 'AI Analysis for ' + selectedFlag.value?.property + ':\n\n• Risk Score: ' + riskScore + '\n\n• Rent Analysis:\n  - Declared Rent (Tax Payer): ' + selectedFlag.value?.declaredRent + '/year\n  - Estimated Rent (Market Value): ' + estimatedRent.value + '/year\n  - Compliance Gap: ' + (gap > 0 ? '+' : '') + 'N' + Math.abs(gap).toLocaleString() + (gap > 0 ? ' (TAX ISSUE)' : ' (Compliant)') + '\n\n• Key Findings:\n  - ' + (gap > 0 ? 'Estimated rent exceeds declared rent by ' + gapPercent + '%' : 'Rent declaration appears accurate') + '\n  - Property has been flagged for ' + selectedFlag.value?.reason + '\n  - Based on location (' + selectedFlag.value?.location + ') and property type (' + selectedFlag.value?.houseType + '), market rate analysis suggests ' + (gap > 0 ? 'under-declaration' : 'accurate reporting') + '\n\n• Suggested Actions:\n  1. ' + (gap > 0 ? 'Issue compliance notice for rent under-declaration' : 'Document findings and close case') + '\n  2. Request supporting documents from taxpayer\n  3. Schedule physical property inspection\n  4. Compare with similar properties in ' + selectedFlag.value?.location?.split(',')[0] + ' area'
    showToast('AI analysis complete')
  }, 1500)
}
</script>

<template>
  <div class="min-h-screen flex bg-[#f5f7fa]">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />
    <div class="flex-1 flex flex-col">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6 shrink-0">
        <div class="flex items-center gap-4"><span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span><span class="text-[#d1d5db]">/</span><span class="text-[#1f2937] text-sm font-medium">Flagged Properties</span></div>
        <div class="flex items-center gap-4"><span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span><button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button></div>
      </header>
      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb]"><h2 class="text-[13px] font-semibold text-[#1f2937]">Flagged Properties</h2></div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4">
            <input v-model="searchQuery" type="text" placeholder="Search..." class="input-field max-w-md" />
            <select v-model="filterPriority" class="input-field w-48"><option value="all">All Priority</option><option value="Critical">Critical</option><option value="High">High</option><option value="Medium">Medium</option></select>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead><tr><th class="table-header">Flag ID</th><th class="table-header">Property</th><th class="table-header">Owner</th><th class="table-header">Reason</th><th class="table-header">Priority</th><th class="table-header">Status</th><th class="table-header">Date</th><th class="table-header">Investigator</th><th class="table-header">Actions</th></tr></thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="flag in paginatedProperties" :key="flag.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ flag.id }}</td><td class="table-cell text-[#6b7280]">{{ flag.property }}</td><td class="table-cell">{{ flag.owner }}</td><td class="table-cell">{{ flag.reason }}</td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-red-50 text-red-700': flag.priority === 'Critical', 'bg-orange-50 text-orange-700': flag.priority === 'High', 'bg-yellow-50 text-yellow-700': flag.priority === 'Medium'}">{{ flag.priority }}</span></td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-blue-50 text-blue-700': flag.status === 'Under Investigation', 'bg-yellow-50 text-yellow-700': flag.status === 'Pending Review', 'bg-green-50 text-green-700': flag.status === 'Resolved'}">{{ flag.status }}</span></td>
                  <td class="table-cell text-[#9ca3af]">{{ flag.flaggedDate }}</td><td class="table-cell">{{ flag.investigator }}</td>
                  <td class="table-cell"><div class="flex gap-2"><button @click="openViewModal(flag)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">View</button><button @click="startInvestigation(flag)" class="px-3 py-1 text-[11px] bg-[#B90B0B] text-white rounded hover:bg-[#6a0707]">Investigate</button></div></td>
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
        <div class="bg-white rounded-xl shadow-xl w-full max-w-4xl my-8">
          <div class="bg-[#B90B0B] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Property Investigation Details</h3>
            <button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-6">
            <div class="grid grid-cols-4 gap-4">
              <div class="bg-[#EEEEEE] rounded-lg p-4">
                <h4 class="text-[11px] text-[#6b7280] mb-3 font-semibold">OCCUPANT INFO</h4>
                <div class="space-y-3">
                  <div>
                    <p class="text-[10px] text-[#9ca3af]">Occupant ID</p>
                    <p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedFlag?.occupantId }}</p>
                  </div>
                  <div>
                    <p class="text-[10px] text-[#9ca3af]">Location</p>
                    <p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedFlag?.location }}</p>
                  </div>
                </div>
              </div>
              <div class="bg-[#EEEEEE] rounded-lg p-4">
                <h4 class="text-[11px] text-[#6b7280] mb-3 font-semibold">PROPERTY INFO</h4>
                <div class="space-y-3">
                  <div>
                    <p class="text-[10px] text-[#9ca3af]">Property ID</p>
                    <p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedFlag?.id }}</p>
                  </div>
                  <div>
                    <p class="text-[10px] text-[#9ca3af]">Risk Level</p>
                    <span class="inline-block px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-red-50 text-red-700': selectedFlag?.riskLevel === 'Critical', 'bg-orange-50 text-orange-700': selectedFlag?.riskLevel === 'High', 'bg-yellow-50 text-yellow-700': selectedFlag?.riskLevel === 'Medium'}">{{ selectedFlag?.riskLevel }}</span>
                  </div>
                  <div>
                    <p class="text-[10px] text-[#9ca3af]">Tax Type</p>
                    <p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedFlag?.taxType }}</p>
                  </div>
                </div>
              </div>
              <div class="bg-[#EEEEEE] rounded-lg p-4">
                <h4 class="text-[11px] text-[#6b7280] mb-3 font-semibold">VALUATION</h4>
                <div class="space-y-3">
                  <div>
                    <p class="text-[10px] text-[#9ca3af]">House Type</p>
                    <p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedFlag?.houseType }}</p>
                  </div>
                  <div>
                    <p class="text-[10px] text-[#9ca3af]">Declared Value</p>
                    <p class="text-[13px] font-semibold text-red-600">{{ selectedFlag?.declaredValue }}</p>
                  </div>
                  <div>
                    <p class="text-[10px] text-[#9ca3af]">Estimated Value</p>
                    <p class="text-[13px] font-semibold text-green-700">{{ selectedFlag?.estimatedValue }}</p>
                  </div>
                </div>
              </div>
              <div class="bg-[#EEEEEE] rounded-lg p-4">
                <h4 class="text-[11px] text-[#6b7280] mb-3 font-semibold">RENT ANALYSIS</h4>
                <div class="space-y-3">
                  <div>
                    <p class="text-[10px] text-[#9ca3af]">Declared Rent (Tax Payer)</p>
                    <p class="text-[13px] font-semibold text-red-600">{{ selectedFlag?.declaredRent }}/year</p>
                  </div>
                  <div>
                    <p class="text-[10px] text-[#9ca3af]">Estimated Rent (Market)</p>
                    <input v-model="estimatedRent" type="text" placeholder="N0" class="w-full px-2 py-1 text-[13px] font-semibold border border-[#d1d5db] rounded focus:ring-1 focus:ring-[#B90B0B] focus:border-transparent" />
                  </div>
                  <div v-if="calculateComplianceGap !== null">
                    <p class="text-[10px] text-[#9ca3af]">Compliance Gap</p>
                    <div class="flex items-center gap-2">
                      <p class="text-[13px] font-semibold" :class="hasTaxIssue ? 'text-red-600' : 'text-green-700'">
                        {{ hasTaxIssue ? '+' : '' }}N{{ Math.abs(calculateComplianceGap).toLocaleString() }}
                      </p>
                      <span v-if="hasTaxIssue" class="px-1.5 py-0.5 text-[10px] bg-red-100 text-red-700 rounded-full font-medium">TAX ISSUE</span>
                      <span v-else class="px-1.5 py-0.5 text-[10px] bg-green-100 text-green-700 rounded-full font-medium">OK</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div>
              <h4 class="text-[11px] text-[#6b7280] mb-2 font-semibold">PROPERTY LOCATION</h4>
              <div class="w-full h-64 bg-[#EEEEEE] rounded-lg overflow-hidden relative">
                <div class="absolute inset-0 flex items-center justify-center">
                  <div class="text-center">
                    <div class="w-16 h-16 bg-[#B90B0B] rounded-full flex items-center justify-center mx-auto mb-2 shadow-lg">
                      <span class="text-white text-2xl">📍</span>
                    </div>
                    <p class="text-[13px] font-semibold text-[#1f2937]">{{ selectedFlag?.location }}</p>
                    <p class="text-[11px] text-[#6b7280]">Coordinates: {{ selectedFlag?.lat }}, {{ selectedFlag?.lng }}</p>
                  </div>
                </div>
                <div class="absolute bottom-3 right-3 bg-white/90 px-3 py-1 rounded text-[10px] text-[#6b7280]">
                  Google Maps Placeholder
                </div>
              </div>
            </div>

            <div>
              <h4 class="text-[11px] text-[#6b7280] mb-2 font-semibold">AI INSIGHTS</h4>
              <textarea v-model="aiInsights" rows="5" :placeholder="`Enter the property details above to obtain insights...\n\nExample: Analyze this property for tax compliance risks based on the declared vs estimated value discrepancy.`" class="w-full px-4 py-3 border border-[#e5e7eb] rounded-lg text-[13px] resize-none focus:ring-2 focus:ring-[#B90B0B] focus:border-transparent"></textarea>
              <button @click="getAIRecommendations" class="mt-3 px-4 py-2 text-[11px] bg-[#B90B0B] text-white rounded-lg hover:bg-[#991010] flex items-center gap-2">
                <span>🤖</span> Get AI Recommendations
              </button>
            </div>

            <div class="flex items-center justify-between pt-4 border-t border-[#e5e7eb]">
              <div class="text-[11px] text-[#6b7280]">
                <span class="font-medium">Owner:</span> {{ selectedFlag?.owner }} | 
                <span class="font-medium">Reason:</span> {{ selectedFlag?.reason }} | 
                <span class="font-medium">Investigator:</span> {{ selectedFlag?.investigator }}
              </div>
              <div class="flex gap-3">
                <button @click="escalateAudit" class="px-4 py-2 text-[11px] bg-[#1f2937] text-white rounded-lg hover:bg-[#374151]">
                  Escalate Audit
                </button>
                <button @click="getComplianceNotice" class="px-4 py-2 text-[11px] bg-[#B90B0B] text-white rounded-lg hover:bg-[#991010]">
                  Get Compliance Notice
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
