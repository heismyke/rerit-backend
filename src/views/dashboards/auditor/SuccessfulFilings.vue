<script setup lang="ts">
import { useRoleStore } from '@/stores'
import { useRouter } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import { ref, computed, onMounted } from 'vue'
import { useAuditorCasesStore, type SuccessfulFiling } from '@/stores/auditorCasesStore'

const { selectedRole, user, logout } = useRoleStore()
const router = useRouter()
const handleLogout = () => { logout(); router.push('/') }

const searchQuery = ref('')
const currentPage = ref(1)
const itemsPerPage = ref(5)
const showEditModal = ref(false)
const selectedFiling = ref<any>(null)
const editFiling = ref<{ status: SuccessfulFiling['status'] }>({ status: 'Validated' })
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })

const { successfulFilings, loadAuditorCases, updateSuccessfulFiling } = useAuditorCasesStore()

onMounted(() => {
  loadAuditorCases().catch(() => showToast('Using offline filing data'))
})

const filteredFilings = computed(() => successfulFilings.value.filter(f => {
  const q = searchQuery.value.toLowerCase()
  return f.id.toLowerCase().includes(q)
    || f.filingId.toLowerCase().includes(q)
    || f.property.toLowerCase().includes(q)
    || f.taxpayer.toLowerCase().includes(q)
}))

const totalPages = computed(() => Math.ceil(filteredFilings.value.length / itemsPerPage.value))
const paginatedFilings = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  return filteredFilings.value.slice(start, start + itemsPerPage.value)
})
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }
const showToast = (message: string) => { toast.value = { show: true, message }; setTimeout(() => { toast.value.show = false }, 3000) }

const openEditModal = (f: any) => {
  selectedFiling.value = f
  editFiling.value = { status: f.status }
  showEditModal.value = true
}

const handleUpdateFiling = async () => {
  const index = successfulFilings.value.findIndex(f => f.id === selectedFiling.value.id)
  if (index !== -1) {
    const updated = { ...successfulFilings.value[index], ...editFiling.value }
    try {
      await updateSuccessfulFiling(updated.id, updated)
    } catch {
      successfulFilings.value[index] = updated
    }
    showToast('Filing updated')
  }
  showEditModal.value = false
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
          <span class="text-[#1f2937] text-sm font-medium">Successful Filings</span>
        </div>
        <div class="flex items-center gap-4">
          <span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span>
          <button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button>
        </div>
      </header>

      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between">
            <h2 class="text-[13px] font-semibold text-[#1f2937]">Successful Filings</h2>
            <span class="text-[11px] text-[#6b7280]">Validated by RERIT</span>
          </div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4">
            <input v-model="searchQuery" type="text" placeholder="Search by ID, filing, property, or taxpayer..." class="input-field flex-1" />
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr>
                  <th class="table-header">Record ID</th>
                  <th class="table-header">FCT‑IRS Filing ID</th>
                  <th class="table-header">Property</th>
                  <th class="table-header">Taxpayer</th>
                  <th class="table-header">Validated</th>
                  <th class="table-header">Status</th>
                  <th class="table-header">Actions</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="filing in paginatedFilings" :key="filing.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ filing.id }}</td>
                  <td class="table-cell">{{ filing.filingId }}</td>
                  <td class="table-cell">{{ filing.property }}</td>
                  <td class="table-cell">{{ filing.taxpayer }}</td>
                  <td class="table-cell text-[#9ca3af]">{{ filing.validatedAt }}</td>
                  <td class="table-cell">
                    <span class="px-2 py-0.5 text-[11px] font-medium rounded-full bg-green-50 text-green-700">{{ filing.status }}</span>
                  </td>
                  <td class="table-cell">
                    <button @click="openEditModal(filing)" class="px-3 py-1 text-[11px] bg-green-50 text-green-700 rounded hover:bg-green-100">Edit</button>
                  </td>
                </tr>
                <tr v-if="paginatedFilings.length === 0">
                  <td colspan="7" class="px-6 py-8 text-center text-[12px] text-[#6b7280]">No successful filings found.</td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="px-6 py-4 flex items-center justify-between border-t border-[#f3f4f6]">
            <p class="text-[11px] text-[#6b7280]">Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, filteredFilings.length) }} of {{ filteredFilings.length }} entries</p>
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
      <div v-if="showEditModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Edit Filing</h3>
            <button @click="showEditModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-4">
            <div class="bg-[#f9fafb] rounded-lg p-4 border border-[#e5e7eb]">
              <p class="text-[11px] text-gray-500">Filing</p>
              <p class="text-[13px] font-medium">{{ selectedFiling?.id }} — {{ selectedFiling?.property }}</p>
            </div>
            <div>
              <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Status</label>
              <select v-model="editFiling.status" class="input-field w-full">
                <option>Reviewed</option>
              </select>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-end">
            <div class="flex gap-3">
              <button @click="showEditModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button>
              <button @click="handleUpdateFiling" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Save</button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
