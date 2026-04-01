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
const currentPage = ref(1)
const itemsPerPage = ref(5)

const showAddModal = ref(false)
const showEditModal = ref(false)
const showViewModal = ref(false)
const selectedRule = ref<any>(null)
const toast = ref<{ show: boolean; message: string; type: string }>({ show: false, message: '', type: 'success' })

const newRule = ref({ name: '', description: '', status: 'Active', riskThreshold: 50, category: 'Valuation' })
const editRule = ref({ name: '', description: '', status: 'Active', riskThreshold: 50, category: 'Valuation' })

const rules = ref([
  { id: 'R-001', name: 'High Value Property Review', description: 'Flag properties above N500M for manual review', triggerCount: 234, status: 'Active', created: '2024-01-10', riskThreshold: 70, category: 'Valuation' },
  { id: 'R-002', name: 'Late Payment Alert', description: 'Send alert when payment is 30 days overdue', triggerCount: 156, status: 'Active', created: '2024-01-08', riskThreshold: 60, category: 'Payment' },
  { id: 'R-003', name: 'Duplicate Owner Detection', description: 'Flag multiple properties owned by same entity', triggerCount: 89, status: 'Active', created: '2024-01-05', riskThreshold: 75, category: 'Ownership' },
  { id: 'R-004', name: 'Survey Expiry Check', description: 'Alert when survey is older than 2 years', triggerCount: 67, status: 'Draft', created: '2024-01-03', riskThreshold: 40, category: 'Survey' },
  { id: 'R-005', name: 'Value Discrepancy', description: 'Flag when declared value differs by 30%', triggerCount: 45, status: 'Active', created: '2024-01-01', riskThreshold: 80, category: 'Valuation' },
])

const filteredRules = computed(() => {
  return rules.value.filter(r => r.name.toLowerCase().includes(searchQuery.value.toLowerCase()) || r.id.toLowerCase().includes(searchQuery.value.toLowerCase()))
})

const totalPages = computed(() => Math.ceil(filteredRules.value.length / itemsPerPage.value))

const paginatedRules = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  return filteredRules.value.slice(start, start + itemsPerPage.value)
})

const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) currentPage.value = page
}

const showToast = (message: string, type: string = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => { toast.value.show = false }, 3000)
}

const openAddModal = () => {
  newRule.value = { name: '', description: '', status: 'Active', riskThreshold: 50, category: 'Valuation' }
  showAddModal.value = true
}

const openEditModal = (rule: any) => {
  selectedRule.value = rule
  editRule.value = { ...rule }
  showEditModal.value = true
}

const openViewModal = (rule: any) => {
  selectedRule.value = rule
  showViewModal.value = true
}

const handleAddRule = () => {
  const newId = 'R-' + String(rules.value.length + 1).padStart(3, '0')
  rules.value.unshift({
    id: newId,
    name: newRule.value.name,
    description: newRule.value.description,
    triggerCount: 0,
    status: newRule.value.status,
    created: new Date().toISOString().split('T')[0],
    riskThreshold: newRule.value.riskThreshold,
    category: newRule.value.category
  })
  showAddModal.value = false
  showToast('Rule created successfully')
}

const handleUpdateRule = () => {
  const index = rules.value.findIndex(r => r.id === selectedRule.value.id)
  if (index !== -1) {
    rules.value[index] = { ...rules.value[index], ...editRule.value }
    showToast('Rule updated successfully')
  }
  showEditModal.value = false
}

const handleDeleteRule = () => {
  rules.value = rules.value.filter(r => r.id !== selectedRule.value.id)
  showEditModal.value = false
  showToast('Rule deleted successfully')
}

const getRiskBadgeClass = (threshold: number) => {
  if (threshold >= 75) return 'bg-red-50 text-red-700'
  if (threshold >= 50) return 'bg-orange-50 text-orange-700'
  if (threshold >= 30) return 'bg-yellow-50 text-yellow-700'
  return 'bg-green-50 text-green-700'
}
</script>

<template>
  <div class="min-h-screen flex bg-[#f5f7fa]">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />

    <div class="flex-1">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6">
        <div class="flex items-center gap-4">
          <span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span>
          <span class="text-[#d1d5db]">/</span>
          <span class="text-[#1f2937] text-sm font-medium">Rules & Automation</span>
        </div>
        <div class="flex items-center gap-4">
          <span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span>
          <button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button>
        </div>
      </header>

      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between">
            <h2 class="text-[13px] font-semibold text-[#1f2937]">Automation Rules & Risk Scoring</h2>
            <button @click="openAddModal" class="btn-primary text-[11px]">Create Rule</button>
          </div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4">
            <input v-model="searchQuery" type="text" placeholder="Search rules..." class="input-field max-w-md" />
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr>
                  <th class="table-header">Rule ID</th>
                  <th class="table-header">Name</th>
                  <th class="table-header">Category</th>
                  <th class="table-header">Risk Threshold</th>
                  <th class="table-header">Triggers</th>
                  <th class="table-header">Status</th>
                  <th class="table-header">Created</th>
                  <th class="table-header">Actions</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="rule in paginatedRules" :key="rule.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ rule.id }}</td>
                  <td class="table-cell">{{ rule.name }}</td>
                  <td class="table-cell">
                    <span class="px-2 py-0.5 text-[11px] bg-blue-50 text-blue-700 rounded">{{ rule.category }}</span>
                  </td>
                  <td class="table-cell">
                    <div class="flex items-center gap-2">
                      <div class="w-16 h-2 bg-gray-200 rounded-full overflow-hidden">
                        <div class="h-full rounded-full" :class="getRiskBadgeClass(rule.riskThreshold).replace('bg-', 'bg-').split(' ')[0].replace('text-', '')" :style="{ width: rule.riskThreshold + '%', backgroundColor: rule.riskThreshold >= 75 ? '#dc2626' : rule.riskThreshold >= 50 ? '#ea580c' : rule.riskThreshold >= 30 ? '#ca8a04' : '#16a34a' }"></div>
                      </div>
                      <span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="getRiskBadgeClass(rule.riskThreshold)">{{ rule.riskThreshold }}%</span>
                    </div>
                  </td>
                  <td class="table-cell">
                    <span class="px-2 py-0.5 text-[11px] bg-[#f3f4f6] text-[#374151] rounded">{{ rule.triggerCount }}</span>
                  </td>
                  <td class="table-cell">
                    <span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="rule.status === 'Active' ? 'bg-green-50 text-green-700' : 'bg-yellow-50 text-yellow-700'">{{ rule.status }}</span>
                  </td>
                  <td class="table-cell text-[#9ca3af]">{{ rule.created }}</td>
                  <td class="table-cell">
                    <div class="flex gap-2">
                      <button @click="openEditModal(rule)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">Edit</button>
                      <button @click="openViewModal(rule)" class="px-3 py-1 text-[11px] bg-[#B90B0B] text-white rounded hover:bg-[#991010]">View</button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="px-6 py-4 flex items-center justify-between border-t border-[#f3f4f6]">
            <p class="text-[11px] text-[#6b7280]">Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, filteredRules.length) }} of {{ filteredRules.length }} entries</p>
            <div class="flex items-center gap-1">
              <button @click="goToPage(currentPage - 1)" :disabled="currentPage === 1" class="px-3 py-1 text-[11px] border border-[#e5e7eb] rounded hover:bg-[#f9fafb] disabled:opacity-50">Prev</button>
              <button v-for="p in totalPages" :key="p" @click="goToPage(p)" class="px-3 py-1 text-[11px] border border-[#e5e7eb] rounded hover:bg-[#f9fafb]" :class="currentPage === p ? 'bg-[#1f2937] text-white border-[#1f2937]' : ''">{{ p }}</button>
              <button @click="goToPage(currentPage + 1)" :disabled="currentPage === totalPages" class="px-3 py-1 text-[11px] border border-[#e5e7eb] rounded hover:bg-[#f9fafb] disabled:opacity-50">Next</button>
            </div>
          </div>
        </div>
      </main>
    </div>

    <Teleport to="body">
      <div v-if="toast.show" class="fixed bottom-4 right-4 bg-green-600 text-white px-4 py-2 rounded-lg shadow-lg z-50">
        {{ toast.message }}
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showAddModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#B90B0B] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Create New Rule</h3>
            <button @click="showAddModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-4">
            <div>
              <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Rule Name</label>
              <input v-model="newRule.name" type="text" placeholder="Enter rule name" class="input-field w-full" />
            </div>
            <div>
              <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Category</label>
              <select v-model="newRule.category" class="input-field w-full">
                <option value="Valuation">Valuation</option>
                <option value="Payment">Payment</option>
                <option value="Ownership">Ownership</option>
                <option value="Survey">Survey</option>
                <option value="Compliance">Compliance</option>
              </select>
            </div>
            <div>
              <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Description</label>
              <textarea v-model="newRule.description" placeholder="Enter description" rows="3" class="input-field w-full"></textarea>
            </div>
            <div>
              <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Risk Threshold (%) - Properties above this score will be flagged</label>
              <input v-model="newRule.riskThreshold" type="number" min="0" max="100" class="input-field w-full" />
              <div class="mt-2 flex gap-2">
                <span class="px-2 py-0.5 text-[10px] bg-green-50 text-green-700 rounded">0-29 Low</span>
                <span class="px-2 py-0.5 text-[10px] bg-yellow-50 text-yellow-700 rounded">30-49 Medium</span>
                <span class="px-2 py-0.5 text-[10px] bg-orange-50 text-orange-700 rounded">50-74 High</span>
                <span class="px-2 py-0.5 text-[10px] bg-red-50 text-red-700 rounded">75-100 Critical</span>
              </div>
            </div>
            <div>
              <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Status</label>
              <select v-model="newRule.status" class="input-field w-full">
                <option value="Active">Active</option>
                <option value="Draft">Draft</option>
              </select>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-end">
            <button @click="showAddModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button>
            <button @click="handleAddRule" class="px-4 py-2 text-[11px] bg-[#B90B0B] text-white rounded-lg hover:bg-[#991010]">Create Rule</button>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showEditModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#B90B0B] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Edit Rule</h3>
            <button @click="showEditModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-4">
            <div>
              <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Rule Name</label>
              <input v-model="editRule.name" type="text" class="input-field w-full" />
            </div>
            <div>
              <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Category</label>
              <select v-model="editRule.category" class="input-field w-full">
                <option value="Valuation">Valuation</option>
                <option value="Payment">Payment</option>
                <option value="Ownership">Ownership</option>
                <option value="Survey">Survey</option>
                <option value="Compliance">Compliance</option>
              </select>
            </div>
            <div>
              <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Description</label>
              <textarea v-model="editRule.description" rows="3" class="input-field w-full"></textarea>
            </div>
            <div>
              <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Risk Threshold (%)</label>
              <input v-model="editRule.riskThreshold" type="number" min="0" max="100" class="input-field w-full" />
            </div>
            <div>
              <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Status</label>
              <select v-model="editRule.status" class="input-field w-full">
                <option value="Active">Active</option>
                <option value="Draft">Draft</option>
              </select>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-between">
            <button @click="handleDeleteRule" class="px-4 py-2 text-[11px] text-red-600 border border-red-200 rounded-lg hover:bg-red-50">Delete</button>
            <div class="flex gap-3">
              <button @click="showEditModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button>
              <button @click="handleUpdateRule" class="px-4 py-2 text-[11px] bg-[#B90B0B] text-white rounded-lg hover:bg-[#991010]">Save Changes</button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showViewModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#B90B0B] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Rule Details</h3>
            <button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <p class="text-[11px] text-gray-500">Rule ID</p>
                <p class="text-[13px] font-medium">{{ selectedRule?.id }}</p>
              </div>
              <div>
                <p class="text-[11px] text-gray-500">Status</p>
                <span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="selectedRule?.status === 'Active' ? 'bg-green-50 text-green-700' : 'bg-yellow-50 text-yellow-700'">{{ selectedRule?.status }}</span>
              </div>
              <div>
                <p class="text-[11px] text-gray-500">Rule Name</p>
                <p class="text-[13px] font-medium">{{ selectedRule?.name }}</p>
              </div>
              <div>
                <p class="text-[11px] text-gray-500">Category</p>
                <span class="px-2 py-0.5 text-[11px] bg-blue-50 text-blue-700 rounded">{{ selectedRule?.category }}</span>
              </div>
            </div>
            <div>
              <p class="text-[11px] text-gray-500">Risk Threshold</p>
              <div class="flex items-center gap-2 mt-1">
                <div class="flex-1 h-3 bg-gray-200 rounded-full overflow-hidden">
                  <div class="h-full rounded-full" :style="{ width: selectedRule?.riskThreshold + '%', backgroundColor: selectedRule?.riskThreshold >= 75 ? '#dc2626' : selectedRule?.riskThreshold >= 50 ? '#ea580c' : selectedRule?.riskThreshold >= 30 ? '#ca8a04' : '#16a34a' }"></div>
                </div>
                <span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="getRiskBadgeClass(selectedRule?.riskThreshold)">{{ selectedRule?.riskThreshold }}%</span>
              </div>
            </div>
            <div>
              <p class="text-[11px] text-gray-500">Description</p>
              <p class="text-[13px]">{{ selectedRule?.description }}</p>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <p class="text-[11px] text-gray-500">Trigger Count</p>
                <p class="text-[13px] font-medium">{{ selectedRule?.triggerCount }}</p>
              </div>
              <div>
                <p class="text-[11px] text-gray-500">Created</p>
                <p class="text-[13px]">{{ selectedRule?.created }}</p>
              </div>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-end">
            <button @click="showViewModal = false" class="px-4 py-2 text-[11px] bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200">Close</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
