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

const showAddModal = ref(false)
const showViewModal = ref(false)
const selectedProperty = ref<any>(null)
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })

const newProperty = ref({ name: '', location: '', type: 'Commercial', value: '', declaredRent: '' })

const properties = ref([
  { id: 'PROP-001', name: 'Commercial Complex A', location: 'Victoria Island, Lagos', type: 'Commercial', value: 'N250,000,000', declaredRent: 'N12,000,000', status: 'Compliant', nextDue: 'Mar 31, 2024' },
  { id: 'PROP-002', name: 'Residential Estate B', location: 'Lekki Phase 1, Lagos', type: 'Residential', value: 'N180,000,000', declaredRent: 'N8,500,000', status: 'Pending', nextDue: 'Mar 31, 2024' },
  { id: 'PROP-003', name: 'Mixed Use Development C', location: 'Ikoyi, Lagos', type: 'Mixed Use', value: 'N350,000,000', declaredRent: 'N18,000,000', status: 'Compliant', nextDue: 'Jun 30, 2024' },
  { id: 'PROP-004', name: 'Office Tower D', location: 'Admiralty Way, Lekki', type: 'Commercial', value: 'N420,000,000', declaredRent: 'N22,000,000', status: 'Compliant', nextDue: 'Jun 30, 2024' },
  { id: 'PROP-005', name: 'Land Parcel E', location: 'Epe, Lagos', type: 'Land', value: 'N45,000,000', declaredRent: 'N1,200,000', status: 'Under Review', nextDue: 'Pending' },
])

const filteredProperties = computed(() => properties.value.filter(p => {
  const matchesSearch = p.name.toLowerCase().includes(searchQuery.value.toLowerCase()) || p.location.toLowerCase().includes(searchQuery.value.toLowerCase()) || p.id.toLowerCase().includes(searchQuery.value.toLowerCase())
  const matchesStatus = filterStatus.value === 'all' || p.status === filterStatus.value
  return matchesSearch && matchesStatus
}))

const totalPages = computed(() => Math.ceil(filteredProperties.value.length / itemsPerPage.value))
const paginatedProperties = computed(() => { const start = (currentPage.value - 1) * itemsPerPage.value; return filteredProperties.value.slice(start, start + itemsPerPage.value) })
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }

const showToast = (message: string) => { toast.value = { show: true, message }; setTimeout(() => { toast.value.show = false }, 3000) }
const openViewModal = (p: any) => { selectedProperty.value = p; showViewModal.value = true }
const handleAddProperty = () => {
  const newId = 'PROP-' + String(properties.value.length + 1).padStart(3, '0')
  properties.value.unshift({ id: newId, ...newProperty.value, status: 'Pending', nextDue: 'Pending' })
  showAddModal.value = false; newProperty.value = { name: '', location: '', type: 'Commercial', value: '', declaredRent: '' }; showToast('Property registered successfully')
}
const payTax = (p: any) => { showToast(`Redirecting to payment for ${p.id}`) }
</script>

<template>
  <div class="min-h-screen bg-[#f5f6fa] flex">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />
    <div class="flex-1 flex flex-col">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6 shrink-0">
        <div class="flex items-center gap-4"><span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span><span class="text-[#d1d5db]">/</span><span class="text-[#1f2937] text-sm font-medium">My Properties</span></div>
        <div class="flex items-center gap-4"><span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span><button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button></div>
      </header>
      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between"><h2 class="text-[13px] font-semibold text-[#1f2937]">My Properties</h2><button @click="showAddModal = true" class="btn-primary text-[11px]">Add Property</button></div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4">
            <input v-model="searchQuery" type="text" placeholder="Search by name, location, or ID..." class="input-field flex-1" />
            <select v-model="filterStatus" class="input-field w-48"><option value="all">All Status</option><option value="Compliant">Compliant</option><option value="Pending">Pending</option><option value="Under Review">Under Review</option></select>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead><tr><th class="table-header">Property ID</th><th class="table-header">Name</th><th class="table-header">Location</th><th class="table-header">Type</th><th class="table-header">Value</th><th class="table-header">Declared Rent</th><th class="table-header">Status</th><th class="table-header">Next Due</th><th class="table-header">Actions</th></tr></thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="property in paginatedProperties" :key="property.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ property.id }}</td><td class="table-cell">{{ property.name }}</td><td class="table-cell text-[#6b7280]">{{ property.location }}</td><td class="table-cell">{{ property.type }}</td><td class="table-cell">{{ property.value }}</td>
                  <td class="table-cell text-[#6b7280]">{{ property.declaredRent }}/yr</td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': property.status === 'Compliant', 'bg-yellow-50 text-yellow-700': property.status === 'Pending', 'bg-blue-50 text-blue-700': property.status === 'Under Review'}">{{ property.status }}</span></td>
                  <td class="table-cell text-[#9ca3af]">{{ property.nextDue }}</td>
                  <td class="table-cell"><div class="flex gap-2"><button @click="openViewModal(property)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">View</button><button @click="payTax(property)" class="px-3 py-1 text-[11px] bg-[#B90B0B] text-white rounded hover:bg-[#6a0707]">Pay</button></div></td>
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
      <div v-if="showAddModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#B90B0B] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Register Property</h3><button @click="showAddModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Property Name</label><input v-model="newProperty.name" type="text" placeholder="Enter property name" class="input-field w-full" /></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Location</label><input v-model="newProperty.location" type="text" placeholder="Enter location" class="input-field w-full" /></div>
            <div class="grid grid-cols-2 gap-4">
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Type</label><select v-model="newProperty.type" class="input-field w-full"><option>Commercial</option><option>Residential</option><option>Mixed Use</option><option>Land</option></select></div>
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Declared Value (NGN)</label><input v-model="newProperty.value" type="text" placeholder="N0" class="input-field w-full" /></div>
            </div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Declared Annual Rent (NGN) - This is your self-reported rental income</label><input v-model="newProperty.declaredRent" type="text" placeholder="N0" class="input-field w-full" /></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-end"><button @click="showAddModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button><button @click="handleAddProperty(); newProperty.declaredRent = ''" class="px-4 py-2 text-[11px] bg-[#B90B0B] text-white rounded-lg hover:bg-[#991010]">Register</button></div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showViewModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#B90B0B] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Property Details</h3><button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div><p class="text-[11px] text-gray-500">Property ID</p><p class="text-[13px] font-medium">{{ selectedProperty?.id }}</p></div>
              <div><p class="text-[11px] text-gray-500">Status</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': selectedProperty?.status === 'Compliant', 'bg-yellow-50 text-yellow-700': selectedProperty?.status === 'Pending'}">{{ selectedProperty?.status }}</span></div>
              <div><p class="text-[11px] text-gray-500">Name</p><p class="text-[13px] font-medium">{{ selectedProperty?.name }}</p></div>
              <div><p class="text-[11px] text-gray-500">Type</p><p class="text-[13px]">{{ selectedProperty?.type }}</p></div>
            </div>
            <div><p class="text-[11px] text-gray-500">Location</p><p class="text-[13px]">{{ selectedProperty?.location }}</p></div>
            <div><p class="text-[11px] text-gray-500">Declared Value</p><p class="text-[13px] font-medium text-green-700">{{ selectedProperty?.value }}</p></div>
            <div><p class="text-[11px] text-gray-500">Declared Annual Rent</p><p class="text-[13px] font-medium text-[#1f2937]">{{ selectedProperty?.declaredRent }}/year</p></div>
            <div><p class="text-[11px] text-gray-500">Next Due</p><p class="text-[13px]">{{ selectedProperty?.nextDue }}</p></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-3">
            <button @click="showViewModal = false" class="px-4 py-2 text-[11px] bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200">Close</button>
            <button @click="payTax(selectedProperty); showViewModal = false" class="px-4 py-2 text-[11px] bg-[#B90B0B] text-white rounded-lg hover:bg-[#991010]">Pay Tax</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
