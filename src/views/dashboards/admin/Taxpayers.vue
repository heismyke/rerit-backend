<script setup lang="ts">
import { useRoleStore } from '@/stores'
import { useRouter } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import { ref, computed } from 'vue'

const { selectedRole, user, logout } = useRoleStore()
const router = useRouter()
const handleLogout = () => { logout(); router.push('/') }

const searchQuery = ref('')
const currentPage = ref(1)
const itemsPerPage = ref(5)

const showAddModal = ref(false)
const showEditModal = ref(false)
const showViewModal = ref(false)
const selectedTaxpayer = ref<any>(null)
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })

const newTaxpayer = ref({ name: '', email: '', phone: '', riskLevel: 'Low' })
const editTaxpayer = ref({ name: '', email: '', phone: '', riskLevel: 'Low' })

const taxpayers = ref([
  { id: 'TP-001', name: 'Emeka Okonkwo', email: 'emeka@example.com', phone: '+234 801 234 5678', properties: 3, riskLevel: 'Low' },
  { id: 'TP-002', name: 'Adaobi Nnamdi', email: 'adaobi@example.com', phone: '+234 802 345 6789', properties: 1, riskLevel: 'Medium' },
  { id: 'TP-003', name: 'Chidi Okafor', email: 'chidi@example.com', phone: '+234 803 456 7890', properties: 5, riskLevel: 'Low' },
  { id: 'TP-004', name: 'Folake Adeyemi', email: 'folake@example.com', phone: '+234 804 567 8901', properties: 2, riskLevel: 'High' },
])

const filteredTaxpayers = computed(() => taxpayers.value.filter(t => t.name.toLowerCase().includes(searchQuery.value.toLowerCase()) || t.id.toLowerCase().includes(searchQuery.value.toLowerCase())))
const totalPages = computed(() => Math.ceil(filteredTaxpayers.value.length / itemsPerPage.value))
const paginatedTaxpayers = computed(() => { const start = (currentPage.value - 1) * itemsPerPage.value; return filteredTaxpayers.value.slice(start, start + itemsPerPage.value) })
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }

const showToast = (message: string) => {
  toast.value = { show: true, message }
  setTimeout(() => { toast.value.show = false }, 3000)
}

const openAddModal = () => { newTaxpayer.value = { name: '', email: '', phone: '', riskLevel: 'Low' }; showAddModal.value = true }
const openEditModal = (t: any) => { selectedTaxpayer.value = t; editTaxpayer.value = { ...t }; showEditModal.value = true }
const openViewModal = (t: any) => { selectedTaxpayer.value = t; showViewModal.value = true }

const handleAddTaxpayer = () => {
  const newId = 'TP-' + String(taxpayers.value.length + 1).padStart(3, '0')
  taxpayers.value.unshift({ id: newId, ...newTaxpayer.value, properties: 0 })
  showAddModal.value = false
  showToast('Taxpayer added successfully')
}

const handleUpdateTaxpayer = () => {
  const index = taxpayers.value.findIndex(t => t.id === selectedTaxpayer.value.id)
  if (index !== -1) { taxpayers.value[index] = { ...taxpayers.value[index], ...editTaxpayer.value }; showToast('Taxpayer updated') }
  showEditModal.value = false
}

const handleDeleteTaxpayer = () => {
  taxpayers.value = taxpayers.value.filter(t => t.id !== selectedTaxpayer.value.id)
  showEditModal.value = false
  showToast('Taxpayer deleted')
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
          <span class="text-[#1f2937] text-sm font-medium">Taxpayers</span>
        </div>
        <div class="flex items-center gap-4">
          <span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span>
          <button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button>
        </div>
      </header>
      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between">
            <h2 class="text-[13px] font-semibold text-[#1f2937]">All Taxpayers</h2>
            <button @click="openAddModal" class="btn-primary text-[11px]">Add Taxpayer</button>
          </div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4">
            <input v-model="searchQuery" type="text" placeholder="Search..." class="input-field max-w-md" />
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr>
                  <th class="table-header">ID</th>
                  <th class="table-header">Name</th>
                  <th class="table-header">Email</th>
                  <th class="table-header">Phone</th>
                  <th class="table-header">Properties</th>
                  <th class="table-header">Risk Level</th>
                  <th class="table-header">Actions</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="t in paginatedTaxpayers" :key="t.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ t.id }}</td>
                  <td class="table-cell">{{ t.name }}</td>
                  <td class="table-cell text-[#6b7280]">{{ t.email }}</td>
                  <td class="table-cell text-[#6b7280]">{{ t.phone }}</td>
                  <td class="table-cell">{{ t.properties }}</td>
                  <td class="table-cell">
                    <span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': t.riskLevel === 'Low', 'bg-yellow-50 text-yellow-700': t.riskLevel === 'Medium', 'bg-red-50 text-red-700': t.riskLevel === 'High'}">{{ t.riskLevel }}</span>
                  </td>
                  <td class="table-cell">
                    <div class="flex gap-2">
                      <button @click="openViewModal(t)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">View</button>
                      <button @click="openEditModal(t)" class="px-3 py-1 text-[11px] bg-[#2D5A27] text-white rounded hover:bg-[#1e3d1a]">Edit</button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="px-6 py-4 flex items-center justify-between border-t border-[#f3f4f6]">
            <p class="text-[11px] text-[#6b7280]">Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, filteredTaxpayers.length) }} of {{ filteredTaxpayers.length }} entries</p>
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
      <div v-if="showAddModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Add Taxpayer</h3>
            <button @click="showAddModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-4">
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Full Name</label><input v-model="newTaxpayer.name" type="text" placeholder="Enter name" class="input-field w-full" /></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Email</label><input v-model="newTaxpayer.email" type="email" placeholder="Enter email" class="input-field w-full" /></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Phone</label><input v-model="newTaxpayer.phone" type="text" placeholder="+234 ..." class="input-field w-full" /></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Risk Level</label><select v-model="newTaxpayer.riskLevel" class="input-field w-full"><option>Low</option><option>Medium</option><option>High</option></select></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-end">
            <button @click="showAddModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button>
            <button @click="handleAddTaxpayer" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Add Taxpayer</button>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showEditModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Edit Taxpayer</h3>
            <button @click="showEditModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-4">
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Full Name</label><input v-model="editTaxpayer.name" type="text" class="input-field w-full" /></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Email</label><input v-model="editTaxpayer.email" type="email" class="input-field w-full" /></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Phone</label><input v-model="editTaxpayer.phone" type="text" class="input-field w-full" /></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Risk Level</label><select v-model="editTaxpayer.riskLevel" class="input-field w-full"><option>Low</option><option>Medium</option><option>High</option></select></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-between">
            <button @click="handleDeleteTaxpayer" class="px-4 py-2 text-[11px] text-red-600 border border-red-200 rounded-lg hover:bg-red-50">Delete</button>
            <div class="flex gap-3">
              <button @click="showEditModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button>
              <button @click="handleUpdateTaxpayer" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Save Changes</button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showViewModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Taxpayer Details</h3>
            <button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div><p class="text-[11px] text-gray-500">ID</p><p class="text-[13px] font-medium">{{ selectedTaxpayer?.id }}</p></div>
              <div><p class="text-[11px] text-gray-500">Risk Level</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': selectedTaxpayer?.riskLevel === 'Low', 'bg-yellow-50 text-yellow-700': selectedTaxpayer?.riskLevel === 'Medium', 'bg-red-50 text-red-700': selectedTaxpayer?.riskLevel === 'High'}">{{ selectedTaxpayer?.riskLevel }}</span></div>
              <div><p class="text-[11px] text-gray-500">Name</p><p class="text-[13px] font-medium">{{ selectedTaxpayer?.name }}</p></div>
              <div><p class="text-[11px] text-gray-500">Properties</p><p class="text-[13px]">{{ selectedTaxpayer?.properties }}</p></div>
            </div>
            <div><p class="text-[11px] text-gray-500">Email</p><p class="text-[13px]">{{ selectedTaxpayer?.email }}</p></div>
            <div><p class="text-[11px] text-gray-500">Phone</p><p class="text-[13px]">{{ selectedTaxpayer?.phone }}</p></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-end"><button @click="showViewModal = false" class="px-4 py-2 text-[11px] bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200">Close</button></div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
