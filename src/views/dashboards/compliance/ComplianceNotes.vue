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
const filterType = ref('all')
const currentPage = ref(1)
const itemsPerPage = ref(5)

const showAddModal = ref(false)
const showViewModal = ref(false)
const showEditModal = ref(false)
const selectedNote = ref<any>(null)
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })

const newNote = ref({ title: '', caseId: '', type: 'Investigation', content: '' })
const editNote = ref({ title: '', caseId: '', type: 'Investigation', content: '' })

const notes = ref([
  { id: 'NOTE-001', title: 'Suspicious Valuation Pattern', caseId: 'CASE-001', type: 'Investigation', author: 'Agent A', created: '2024-01-15', lastUpdated: '2024-01-15', content: 'Property values appear to be consistently under-reported in this area. Recommend further investigation.' },
  { id: 'NOTE-002', title: 'Document Verification Required', caseId: 'CASE-002', type: 'Action Item', author: 'Agent B', created: '2024-01-12', lastUpdated: '2024-01-14', content: 'Owner has been contacted to provide additional documentation for property verification.' },
  { id: 'NOTE-003', title: 'Owner Interview Scheduled', caseId: 'CASE-003', type: 'Meeting', author: 'Agent A', created: '2024-01-10', lastUpdated: '2024-01-13', content: 'Interview scheduled for January 20th, 2024 at 10:00 AM at the NRS office.' },
  { id: 'NOTE-004', title: 'Evidence Collected', caseId: 'CASE-004', type: 'Evidence', author: 'Agent C', created: '2024-01-08', lastUpdated: '2024-01-12', content: 'Photographs and video evidence of property condition collected. Documents attached.' },
])

const filteredNotes = computed(() => {
  return notes.value.filter(n => {
    const matchesSearch = n.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      n.id.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesType = filterType.value === 'all' || n.type === filterType.value
    return matchesSearch && matchesType
  })
})

const totalPages = computed(() => Math.ceil(filteredNotes.value.length / itemsPerPage.value))

const paginatedNotes = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  return filteredNotes.value.slice(start, start + itemsPerPage.value)
})

const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
  }
}

const showToast = (message: string) => {
  toast.value = { show: true, message }
  setTimeout(() => { toast.value.show = false }, 3000)
}

const openAddModal = () => {
  newNote.value = { title: '', caseId: '', type: 'Investigation', content: '' }
  showAddModal.value = true
}

const openViewModal = (note: any) => {
  selectedNote.value = note
  showViewModal.value = true
}

const openEditModal = (note: any) => {
  selectedNote.value = note
  editNote.value = { title: note.title, caseId: note.caseId, type: note.type, content: note.content }
  showEditModal.value = true
}

const handleAddNote = () => {
  const newId = 'NOTE-' + String(notes.value.length + 1).padStart(3, '0')
  notes.value.unshift({
    id: newId,
    title: newNote.value.title,
    caseId: newNote.value.caseId,
    type: newNote.value.type,
    author: user.value?.name || 'Agent',
    created: new Date().toISOString().split('T')[0],
    lastUpdated: new Date().toISOString().split('T')[0],
    content: newNote.value.content
  })
  showAddModal.value = false
  showToast('Note added successfully')
}

const handleUpdateNote = () => {
  const index = notes.value.findIndex(n => n.id === selectedNote.value.id)
  if (index !== -1) {
    notes.value[index] = {
      ...notes.value[index],
      title: editNote.value.title,
      caseId: editNote.value.caseId,
      type: editNote.value.type,
      content: editNote.value.content,
      lastUpdated: new Date().toISOString().split('T')[0]
    }
    showToast('Note updated successfully')
  }
  showEditModal.value = false
}

const handleDeleteNote = () => {
  notes.value = notes.value.filter(n => n.id !== selectedNote.value.id)
  showEditModal.value = false
  showToast('Note deleted')
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
          <span class="text-[#1f2937] text-sm font-medium">Compliance Notes</span>
        </div>
        <div class="flex items-center gap-4">
          <span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span>
          <button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button>
        </div>
      </header>

      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between">
            <h2 class="text-[13px] font-semibold text-[#1f2937]">Compliance Notes</h2>
            <button @click="openAddModal" class="btn-primary text-[11px]">Add Note</button>
          </div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4">
            <input v-model="searchQuery" type="text" placeholder="Search..." class="input-field max-w-md" />
            <select v-model="filterType" class="input-field w-48">
              <option value="all">All Types</option>
              <option value="Investigation">Investigation</option>
              <option value="Action Item">Action Item</option>
              <option value="Meeting">Meeting</option>
              <option value="Evidence">Evidence</option>
            </select>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr>
                  <th class="table-header">Note ID</th>
                  <th class="table-header">Title</th>
                  <th class="table-header">Case ID</th>
                  <th class="table-header">Type</th>
                  <th class="table-header">Author</th>
                  <th class="table-header">Created</th>
                  <th class="table-header">Updated</th>
                  <th class="table-header">Actions</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="note in paginatedNotes" :key="note.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ note.id }}</td>
                  <td class="table-cell">{{ note.title }}</td>
                  <td class="table-cell">
                    <span class="px-2 py-0.5 text-[11px] bg-[#f3f4f6] text-[#6b7280] rounded">{{ note.caseId }}</span>
                  </td>
                  <td class="table-cell">
                    <span class="px-2 py-0.5 text-[11px] font-medium rounded-full"
                      :class="{
                        'bg-red-50 text-red-700': note.type === 'Investigation',
                        'bg-blue-50 text-blue-700': note.type === 'Action Item',
                        'bg-purple-50 text-purple-700': note.type === 'Meeting',
                        'bg-green-50 text-green-700': note.type === 'Evidence',
                      }">{{ note.type }}</span>
                  </td>
                  <td class="table-cell">{{ note.author }}</td>
                  <td class="table-cell text-[#9ca3af]">{{ note.created }}</td>
                  <td class="table-cell text-[#9ca3af]">{{ note.lastUpdated }}</td>
                  <td class="table-cell">
                    <div class="flex gap-2">
                      <button @click="openViewModal(note)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">View</button>
                      <button @click="openEditModal(note)" class="px-3 py-1 text-[11px] bg-[#2D5A27] text-white rounded hover:bg-[#1e3d1a]">Edit</button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="px-6 py-4 flex items-center justify-between border-t border-[#f3f4f6]">
            <p class="text-[11px] text-[#6b7280]">Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, filteredNotes.length) }} of {{ filteredNotes.length }} entries</p>
            <div class="flex items-center gap-1">
              <button @click="goToPage(currentPage - 1)" :disabled="currentPage === 1" class="px-3 py-1 text-[11px] border border-[#e5e7eb] rounded hover:bg-[#f9fafb] disabled:opacity-50 disabled:cursor-not-allowed">Prev</button>
              <button v-for="p in totalPages" :key="p" @click="goToPage(p)" class="px-3 py-1 text-[11px] border border-[#e5e7eb] rounded hover:bg-[#f9fafb]" :class="currentPage === p ? 'bg-[#1f2937] text-white border-[#1f2937]' : ''">{{ p }}</button>
              <button @click="goToPage(currentPage + 1)" :disabled="currentPage === totalPages" class="px-3 py-1 text-[11px] border border-[#e5e7eb] rounded hover:bg-[#f9fafb] disabled:opacity-50 disabled:cursor-not-allowed">Next</button>
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
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Add Note</h3>
            <button @click="showAddModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-4">
            <div>
              <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Title</label>
              <input v-model="newNote.title" type="text" placeholder="Enter note title" class="input-field w-full" />
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Case ID</label>
                <input v-model="newNote.caseId" type="text" placeholder="CASE-XXX" class="input-field w-full" />
              </div>
              <div>
                <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Type</label>
                <select v-model="newNote.type" class="input-field w-full">
                  <option>Investigation</option>
                  <option>Action Item</option>
                  <option>Meeting</option>
                  <option>Evidence</option>
                </select>
              </div>
            </div>
            <div>
              <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Content</label>
              <textarea v-model="newNote.content" rows="4" placeholder="Enter note content..." class="input-field w-full"></textarea>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-end">
            <button @click="showAddModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button>
            <button @click="handleAddNote" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Add Note</button>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showViewModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Note Details</h3>
            <button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <p class="text-[11px] text-gray-500">Note ID</p>
                <p class="text-[13px] font-medium">{{ selectedNote?.id }}</p>
              </div>
              <div>
                <p class="text-[11px] text-gray-500">Case ID</p>
                <span class="px-2 py-0.5 text-[11px] bg-[#f3f4f6] text-[#6b7280] rounded">{{ selectedNote?.caseId }}</span>
              </div>
              <div>
                <p class="text-[11px] text-gray-500">Author</p>
                <p class="text-[13px] font-medium">{{ selectedNote?.author }}</p>
              </div>
              <div>
                <p class="text-[11px] text-gray-500">Type</p>
                <span class="px-2 py-0.5 text-[11px] font-medium rounded-full"
                  :class="{
                    'bg-red-50 text-red-700': selectedNote?.type === 'Investigation',
                    'bg-blue-50 text-blue-700': selectedNote?.type === 'Action Item',
                    'bg-purple-50 text-purple-700': selectedNote?.type === 'Meeting',
                    'bg-green-50 text-green-700': selectedNote?.type === 'Evidence',
                  }">{{ selectedNote?.type }}</span>
              </div>
            </div>
            <div>
              <p class="text-[11px] text-gray-500">Title</p>
              <p class="text-[13px] font-medium">{{ selectedNote?.title }}</p>
            </div>
            <div>
              <p class="text-[11px] text-gray-500">Content</p>
              <p class="text-[13px] text-gray-700">{{ selectedNote?.content }}</p>
            </div>
            <div class="grid grid-cols-2 gap-4 pt-2 border-t border-gray-100">
              <div>
                <p class="text-[11px] text-gray-500">Created</p>
                <p class="text-[12px]">{{ selectedNote?.created }}</p>
              </div>
              <div>
                <p class="text-[11px] text-gray-500">Last Updated</p>
                <p class="text-[12px]">{{ selectedNote?.lastUpdated }}</p>
              </div>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-end">
            <button @click="showViewModal = false" class="px-4 py-2 text-[11px] bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200">Close</button>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showEditModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Edit Note</h3>
            <button @click="showEditModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-4">
            <div>
              <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Title</label>
              <input v-model="editNote.title" type="text" class="input-field w-full" />
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Case ID</label>
                <input v-model="editNote.caseId" type="text" class="input-field w-full" />
              </div>
              <div>
                <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Type</label>
                <select v-model="editNote.type" class="input-field w-full">
                  <option>Investigation</option>
                  <option>Action Item</option>
                  <option>Meeting</option>
                  <option>Evidence</option>
                </select>
              </div>
            </div>
            <div>
              <label class="block text-[11px] font-medium text-gray-600 mb-1.5">Content</label>
              <textarea v-model="editNote.content" rows="4" class="input-field w-full"></textarea>
            </div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-between">
            <button @click="handleDeleteNote" class="px-4 py-2 text-[11px] text-red-600 border border-red-200 rounded-lg hover:bg-red-50">Delete</button>
            <div class="flex gap-3">
              <button @click="showEditModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button>
              <button @click="handleUpdateNote" class="px-4 py-2 text-[11px] bg-[#2D5A27] text-white rounded-lg hover:bg-[#1e3d1a]">Save</button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
