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
const activeTab = ref('all')

const showAddModal = ref(false)
const showEditModal = ref(false)
const showViewModal = ref(false)
const selectedUser = ref<any>(null)
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })

const newUser = ref({ name: '', email: '', role: 'Admin', department: 'IT' })
const editUser = ref({ name: '', email: '', role: 'Admin', department: 'IT', status: 'Active' })

const users = ref([
  { id: 'U-001', name: 'Aisha Mohammed', email: 'aisha.mohammed@nrs.gov.ng', role: 'Admin', department: 'IT', status: 'Active', lastLogin: '2024-01-15 09:30', created: '2023-06-01' },
  { id: 'U-002', name: 'Chukwudi Okafor', email: 'c.okafor@nrs.gov.ng', role: 'Auditor', department: 'Audit', status: 'Active', lastLogin: '2024-01-15 08:15', created: '2023-06-15' },
  // { id: 'U-003', name: 'Funke Adebayo', email: 'f.adebayo@nrs.gov.ng', role: 'Compliance', department: 'Compliance', status: 'Active', lastLogin: '2024-01-14 16:45', created: '2023-07-01' },
  { id: 'U-004', name: 'Ngozi Eze', email: 'n.eze@nrs.gov.ng', role: 'Tax Payer', department: 'Tax Payer', status: 'Active', lastLogin: '2024-01-15 10:00', created: '2023-08-01' },
  { id: 'U-006', name: 'Segun Fashola', email: 's.fashola@nrs.gov.ng', role: 'Admin', department: 'Operations', status: 'Active', lastLogin: '2024-01-14 14:30', created: '2023-08-15' },
  { id: 'U-007', name: 'Amina Yusuf', email: 'a.yusuf@nrs.gov.ng', role: 'Auditor', department: 'Audit', status: 'Pending', lastLogin: '-', created: '2024-01-10' },
])

const roles = ['Admin', 'Auditor', 'Tax Payer']
// const roles = ['Admin', 'Auditor', 'Compliance', 'Tax Payer']

const filteredUsers = computed(() => {
  let items = users.value.filter(u => u.name.toLowerCase().includes(searchQuery.value.toLowerCase()) || u.email.toLowerCase().includes(searchQuery.value.toLowerCase()) || u.id.toLowerCase().includes(searchQuery.value.toLowerCase()))
  if (activeTab.value !== 'all') { items = items.filter(u => u.status.toLowerCase() === activeTab.value) }
  return items
})

const totalPages = computed(() => Math.ceil(filteredUsers.value.length / itemsPerPage.value))
const paginatedUsers = computed(() => { const start = (currentPage.value - 1) * itemsPerPage.value; return filteredUsers.value.slice(start, start + itemsPerPage.value) })
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }
const stats = computed(() => ({ total: users.value.length, active: users.value.filter(u => u.status === 'Active').length, inactive: users.value.filter(u => u.status === 'Inactive').length, pending: users.value.filter(u => u.status === 'Pending').length }))

const showToast = (message: string) => { toast.value = { show: true, message }; setTimeout(() => { toast.value.show = false }, 3000) }
const openAddModal = () => { newUser.value = { name: '', email: '', role: 'Admin', department: 'IT' }; showAddModal.value = true }
const openEditModal = (u: any) => { selectedUser.value = u; editUser.value = { ...u }; showEditModal.value = true }
const openViewModal = (u: any) => { selectedUser.value = u; showViewModal.value = true }

const handleAddUser = () => {
  const newId = 'U-' + String(users.value.length + 1).padStart(3, '0')
  users.value.unshift({ id: newId, ...newUser.value, status: 'Pending', lastLogin: '-', created: new Date().toISOString().split('T')[0] })
  showAddModal.value = false; showToast('User added successfully')
}

const handleUpdateUser = () => {
  const index = users.value.findIndex(u => u.id === selectedUser.value.id)
  if (index !== -1) { users.value[index] = { ...users.value[index], ...editUser.value }; showToast('User updated') }
  showEditModal.value = false
}

const handleDeleteUser = () => {
  users.value = users.value.filter(u => u.id !== selectedUser.value.id)
  showEditModal.value = false; showToast('User deleted')
}
</script>

<template>
  <div class="min-h-screen flex bg-[#f5f7fa]">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />
    <div class="flex-1">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6">
        <div class="flex items-center gap-4"><span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span><span class="text-[#d1d5db]">/</span><span class="text-[#1f2937] text-sm font-medium">Users & Roles</span></div>
        <div class="flex items-center gap-4"><span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span><button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button></div>
      </header>
      <main class="flex-1 p-6">
        <div class="grid grid-cols-4 gap-4 mb-6">
          <div class="bg-white border border-[#e5e7eb] rounded-lg p-4"><p class="text-[11px] text-[#6b7280]">Total Users</p><p class="text-2xl font-semibold text-[#1f2937] mt-1">{{ stats.total }}</p></div>
          <div class="bg-white border border-[#e5e7eb] rounded-lg p-4"><p class="text-[11px] text-[#6b7280]">Active</p><p class="text-2xl font-semibold text-green-600 mt-1">{{ stats.active }}</p></div>
          <div class="bg-white border border-[#e5e7eb] rounded-lg p-4"><p class="text-[11px] text-[#6b7280]">Inactive</p><p class="text-2xl font-semibold text-gray-500 mt-1">{{ stats.inactive }}</p></div>
          <div class="bg-white border border-[#e5e7eb] rounded-lg p-4"><p class="text-[11px] text-[#6b7280]">Pending</p><p class="text-2xl font-semibold text-yellow-600 mt-1">{{ stats.pending }}</p></div>
        </div>
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between">
            <div class="flex gap-4">
              <button @click="activeTab = 'all'" class="text-[13px] font-medium pb-1 border-b-2" :class="activeTab === 'all' ? 'border-[#2D5A27] text-[#2D5A27]' : 'border-transparent text-[#6b7280]'">All Users</button>
              <button @click="activeTab = 'active'" class="text-[13px] font-medium pb-1 border-b-2" :class="activeTab === 'active' ? 'border-[#2D5A27] text-[#2D5A27]' : 'border-transparent text-[#6b7280]'">Active</button>
              <button @click="activeTab = 'inactive'" class="text-[13px] font-medium pb-1 border-b-2" :class="activeTab === 'inactive' ? 'border-[#2D5A27] text-[#2D5A27]' : 'border-transparent text-[#6b7280]'">Inactive</button>
              <button @click="activeTab = 'pending'" class="text-[13px] font-medium pb-1 border-b-2" :class="activeTab === 'pending' ? 'border-[#2D5A27] text-[#2D5A27]' : 'border-transparent text-[#6b7280]'">Pending</button>
            </div>
            <button @click="openAddModal" class="btn-primary text-[11px]">Add User</button>
          </div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4"><input v-model="searchQuery" type="text" placeholder="Search by name, email, or ID..." class="input-field max-w-md" /></div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead><tr><th class="table-header">User ID</th><th class="table-header">Name</th><th class="table-header">Email</th><th class="table-header">Role</th><th class="table-header">Department</th><th class="table-header">Status</th><th class="table-header">Last Login</th><th class="table-header">Actions</th></tr></thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="u in paginatedUsers" :key="u.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ u.id }}</td><td class="table-cell">{{ u.name }}</td><td class="table-cell text-[#6b7280]">{{ u.email }}</td><td class="table-cell"><span class="px-2 py-0.5 text-[11px] bg-[#2D5A27]/10 text-[#2D5A27] rounded font-medium">{{ u.role }}</span></td>
                  <td class="table-cell text-[#9ca3af]">{{ u.department }}</td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': u.status === 'Active', 'bg-gray-100 text-gray-500': u.status === 'Inactive', 'bg-yellow-50 text-yellow-700': u.status === 'Pending'}">{{ u.status }}</span></td>
                  <td class="table-cell text-[#9ca3af]">{{ u.lastLogin }}</td>
                  <td class="table-cell"><div class="flex gap-2"><button @click="openEditModal(u)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">Edit</button><button @click="openViewModal(u)" class="px-3 py-1 text-[11px] bg-[#2D5A27] text-white rounded hover:bg-[#1e3d1a]">View</button></div></td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="px-6 py-4 flex items-center justify-between border-t border-[#f3f4f6]">
            <p class="text-[11px] text-[#6b7280]">Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, filteredUsers.length) }} of {{ filteredUsers.length }} entries</p>
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
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Add User</h3><button @click="showAddModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Full Name</label><input v-model="newUser.name" type="text" placeholder="Enter name" class="input-field w-full" /></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Email</label><input v-model="newUser.email" type="email" placeholder="Enter email" class="input-field w-full" /></div>
            <div class="grid grid-cols-2 gap-4">
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Role</label><select v-model="newUser.role" class="input-field w-full"><option v-for="role in roles" :key="role" :value="role">{{ role }}</option></select></div>
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Department</label><input v-model="newUser.department" type="text" placeholder="Enter department" class="input-field w-full" /></div>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-end"><button @click="showAddModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button><button @click="handleAddUser" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Add User</button></div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showEditModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Edit User</h3><button @click="showEditModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Full Name</label><input v-model="editUser.name" type="text" class="input-field w-full" /></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Email</label><input v-model="editUser.email" type="email" class="input-field w-full" /></div>
            <div class="grid grid-cols-3 gap-4">
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Role</label><select v-model="editUser.role" class="input-field w-full"><option v-for="role in roles" :key="role" :value="role">{{ role }}</option></select></div>
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Department</label><input v-model="editUser.department" type="text" class="input-field w-full" /></div>
              <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Status</label><select v-model="editUser.status" class="input-field w-full"><option>Active</option><option>Inactive</option><option>Pending</option></select></div>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-between"><button @click="handleDeleteUser" class="px-4 py-2 text-[11px] text-red-600 border border-red-200 rounded-lg hover:bg-red-50">Delete</button><div class="flex gap-3"><button @click="showEditModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button><button @click="handleUpdateUser" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Save</button></div></div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showViewModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">User Details</h3><button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div><p class="text-[11px] text-gray-500">User ID</p><p class="text-[13px] font-medium">{{ selectedUser?.id }}</p></div>
              <div><p class="text-[11px] text-gray-500">Status</p><span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': selectedUser?.status === 'Active', 'bg-gray-100 text-gray-500': selectedUser?.status === 'Inactive', 'bg-yellow-50 text-yellow-700': selectedUser?.status === 'Pending'}">{{ selectedUser?.status }}</span></div>
              <div><p class="text-[11px] text-gray-500">Name</p><p class="text-[13px] font-medium">{{ selectedUser?.name }}</p></div>
              <div><p class="text-[11px] text-gray-500">Role</p><span class="px-2 py-0.5 text-[11px] bg-[#2D5A27]/10 text-[#2D5A27] rounded font-medium">{{ selectedUser?.role }}</span></div>
            </div>
            <div><p class="text-[11px] text-gray-500">Email</p><p class="text-[13px]">{{ selectedUser?.email }}</p></div>
            <div><p class="text-[11px] text-gray-500">Department</p><p class="text-[13px]">{{ selectedUser?.department }}</p></div>
            <div><p class="text-[11px] text-gray-500">Last Login</p><p class="text-[13px]">{{ selectedUser?.lastLogin }}</p></div>
            <div><p class="text-[11px] text-gray-500">Created</p><p class="text-[13px]">{{ selectedUser?.created }}</p></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-end"><button @click="showViewModal = false" class="px-4 py-2 text-[11px] bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200">Close</button></div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
