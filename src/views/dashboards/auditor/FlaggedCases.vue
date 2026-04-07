<script setup lang="ts">
import { useRoleStore } from '@/stores'
import { useRouter } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import { ref, computed } from 'vue'

const { selectedRole, user, logout } = useRoleStore()
const router = useRouter()
const handleLogout = () => { logout(); router.push('/') }

const searchQuery = ref('')
const filterMismatch = ref('all')
const currentPage = ref(1)
const itemsPerPage = ref(5)

const showAiModal = ref(false)
const selectedFlag = ref<any>(null)
const aiInsights = ref('')

const flaggedCases = ref([
  { id: 'FLG-2024-011', filingId: 'FCT-IRS-00211', property: 'Plot 18, Wuse II', taxpayer: 'Nwosu Holdings', mismatch: 'Declared rent below benchmark', receivedAt: '2024-01-18', status: 'Pending Review' },
  { id: 'FLG-2024-012', filingId: 'FCT-IRS-00225', property: 'Block 4, Maitama', taxpayer: 'Ayo Martins', mismatch: 'Property size mismatch', receivedAt: '2024-01-18', status: 'Pending Review' },
  { id: 'FLG-2024-013', filingId: 'FCT-IRS-00231', property: 'Unit 12, Gwarinpa', taxpayer: 'Saka Ventures', mismatch: 'Missing land registry reference', receivedAt: '2024-01-19', status: 'In Review' },
  { id: 'FLG-2024-014', filingId: 'FCT-IRS-00237', property: 'Plot 9, Asokoro', taxpayer: 'Dara Okafor', mismatch: 'Owner identity mismatch', receivedAt: '2024-01-20', status: 'In Review' },
])

const filteredFlags = computed(() => flaggedCases.value.filter(f => {
  const matchesSearch = f.id.toLowerCase().includes(searchQuery.value.toLowerCase())
    || f.filingId.toLowerCase().includes(searchQuery.value.toLowerCase())
    || f.property.toLowerCase().includes(searchQuery.value.toLowerCase())
    || f.taxpayer.toLowerCase().includes(searchQuery.value.toLowerCase())
  const matchesMismatch = filterMismatch.value === 'all' || f.mismatch === filterMismatch.value
  return matchesSearch && matchesMismatch
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

const generateAiInsights = () => {
  if (!selectedFlag.value) return
  const f = selectedFlag.value
  aiInsights.value =
    `AI Summary for ${f.property}\n\n` +
    `• Filing ID: ${f.filingId}\n` +
    `• Mismatch: ${f.mismatch}\n` +
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
            <select v-model="filterMismatch" class="input-field w-64">
              <option value="all">All Mismatches</option>
              <option>Declared rent below benchmark</option>
              <option>Property size mismatch</option>
              <option>Missing land registry reference</option>
              <option>Owner identity mismatch</option>
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
                  <th class="table-header">Mismatch</th>
                  <th class="table-header">Received</th>
                  <th class="table-header">Status</th>
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
                    <span class="px-2 py-0.5 text-[11px] bg-red-50 text-red-700 rounded">{{ flag.mismatch }}</span>
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
                    <button @click="openAiModal(flag)" class="px-3 py-1 text-[11px] bg-[#2D5A27]/10 text-[#2D5A27] rounded hover:bg-[#2D5A27]/20">
                      AI Review
                    </button>
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
              <p class="text-[11px] text-gray-500 mt-2">Mismatch: {{ selectedFlag?.mismatch }}</p>
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
