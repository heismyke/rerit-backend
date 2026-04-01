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
const selectedReportType = ref('all')

const showGenerateModal = ref(false)
const showViewModal = ref(false)
const selectedReport = ref<any>(null)
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })

const newReport = ref({ title: '', type: 'Revenue' })

const reports = ref([
  { id: 'RPT-001', title: 'Monthly Revenue Summary', type: 'Revenue', generatedBy: 'Admin User', generatedAt: '2024-01-15 09:30', format: 'PDF', size: '2.4 MB' },
  { id: 'RPT-002', title: 'Property Ownership Changes', type: 'Ownership', generatedBy: 'Admin User', generatedAt: '2024-01-14 14:20', format: 'Excel', size: '1.8 MB' },
  { id: 'RPT-003', title: 'Compliance Status Report', type: 'Compliance', generatedBy: 'Admin User', generatedAt: '2024-01-13 11:45', format: 'PDF', size: '3.1 MB' },
  { id: 'RPT-004', title: 'Taxpayer Activity Log', type: 'Activity', generatedBy: 'Admin User', generatedAt: '2024-01-12 16:00', format: 'CSV', size: '856 KB' },
  { id: 'RPT-005', title: 'Audit Findings Summary', type: 'Audit', generatedBy: 'Admin User', generatedAt: '2024-01-11 10:15', format: 'PDF', size: '1.5 MB' },
  { id: 'RPT-006', title: 'Annual Tax Collection Report', type: 'Revenue', generatedBy: 'Admin User', generatedAt: '2024-01-10 08:00', format: 'PDF', size: '5.2 MB' },
  { id: 'RPT-007', title: 'Survey Submissions Status', type: 'Survey', generatedBy: 'Admin User', generatedAt: '2024-01-09 13:30', format: 'Excel', size: '2.9 MB' },
])

const reportTypes = ['Revenue', 'Ownership', 'Compliance', 'Activity', 'Audit', 'Survey']

const filteredReports = computed(() => {
  let items = reports.value.filter(r => r.title.toLowerCase().includes(searchQuery.value.toLowerCase()) || r.id.toLowerCase().includes(searchQuery.value.toLowerCase()))
  if (selectedReportType.value !== 'all') { items = items.filter(r => r.type === selectedReportType.value) }
  return items
})

const totalPages = computed(() => Math.ceil(filteredReports.value.length / itemsPerPage.value))
const paginatedReports = computed(() => { const start = (currentPage.value - 1) * itemsPerPage.value; return filteredReports.value.slice(start, start + itemsPerPage.value) })
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }

const showToast = (message: string) => { toast.value = { show: true, message }; setTimeout(() => { toast.value.show = false }, 3000) }
const openViewModal = (report: any) => { selectedReport.value = report; showViewModal.value = true }
const handleGenerateReport = () => {
  const newId = 'RPT-' + String(reports.value.length + 1).padStart(3, '0')
  reports.value.unshift({ id: newId, title: newReport.value.title, type: newReport.value.type, generatedBy: 'Admin User', generatedAt: new Date().toISOString().replace('T', ' ').substring(0, 16), format: 'PDF', size: '1.0 MB' })
  showGenerateModal.value = false; newReport.value = { title: '', type: 'Revenue' }; showToast('Report generated successfully')
}
const handleDownload = (report: any) => { showToast(`Downloading ${report.title}...`) }
</script>

<template>
  <div class="min-h-screen flex bg-[#f5f7fa]">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />
    <div class="flex-1">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6">
        <div class="flex items-center gap-4"><span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span><span class="text-[#d1d5db]">/</span><span class="text-[#1f2937] text-sm font-medium">Reports & Analytics</span></div>
        <div class="flex items-center gap-4"><span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span><button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button></div>
      </header>
      <main class="flex-1 p-6">
        <div class="grid grid-cols-4 gap-4 mb-6">
          <div class="bg-white border border-[#e5e7eb] rounded-lg p-4"><p class="text-[11px] text-[#6b7280]">Total Reports</p><p class="text-2xl font-semibold text-[#1f2937] mt-1">{{ reports.length }}</p></div>
          <div class="bg-white border border-[#e5e7eb] rounded-lg p-4"><p class="text-[11px] text-[#6b7280]">This Month</p><p class="text-2xl font-semibold text-[#B90B0B] mt-1">{{ reports.filter(r => r.generatedAt.startsWith('2024-01')).length }}</p></div>
          <div class="bg-white border border-[#e5e7eb] rounded-lg p-4"><p class="text-[11px] text-[#6b7280]">PDF Reports</p><p class="text-2xl font-semibold text-blue-600 mt-1">{{ reports.filter(r => r.format === 'PDF').length }}</p></div>
          <div class="bg-white border border-[#e5e7eb] rounded-lg p-4"><p class="text-[11px] text-[#6b7280]">Excel Reports</p><p class="text-2xl font-semibold text-green-600 mt-1">{{ reports.filter(r => r.format === 'Excel').length }}</p></div>
        </div>
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between"><h2 class="text-[13px] font-semibold text-[#1f2937]">Generated Reports</h2><button @click="showGenerateModal = true" class="btn-primary text-[11px]">Generate New Report</button></div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4 flex-wrap"><input v-model="searchQuery" type="text" placeholder="Search reports..." class="input-field max-w-xs" /><select v-model="selectedReportType" class="input-field max-w-xs"><option value="all">All Types</option><option v-for="type in reportTypes" :key="type" :value="type">{{ type }}</option></select></div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead><tr><th class="table-header">Report ID</th><th class="table-header">Title</th><th class="table-header">Type</th><th class="table-header">Generated By</th><th class="table-header">Generated At</th><th class="table-header">Format</th><th class="table-header">Size</th><th class="table-header">Actions</th></tr></thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="report in paginatedReports" :key="report.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ report.id }}</td><td class="table-cell">{{ report.title }}</td><td class="table-cell"><span class="px-2 py-0.5 text-[11px] bg-[#f3f4f6] text-[#374151] rounded">{{ report.type }}</span></td>
                  <td class="table-cell text-[#6b7280]">{{ report.generatedBy }}</td><td class="table-cell text-[#9ca3af]">{{ report.generatedAt }}</td>
                  <td class="table-cell"><span class="px-2 py-0.5 text-[11px] font-medium rounded" :class="{'bg-red-50 text-red-700': report.format === 'PDF', 'bg-green-50 text-green-700': report.format === 'Excel', 'bg-blue-50 text-blue-700': report.format === 'CSV'}">{{ report.format }}</span></td>
                  <td class="table-cell text-[#9ca3af]">{{ report.size }}</td>
                  <td class="table-cell"><div class="flex gap-2"><button @click="handleDownload(report)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">Download</button><button @click="openViewModal(report)" class="px-3 py-1 text-[11px] bg-[#B90B0B] text-white rounded hover:bg-[#991010]">View</button></div></td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="px-6 py-4 flex items-center justify-between border-t border-[#f3f4f6]">
            <p class="text-[11px] text-[#6b7280]">Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, filteredReports.length) }} of {{ filteredReports.length }} entries</p>
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
      <div v-if="showGenerateModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#B90B0B] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Generate Report</h3><button @click="showGenerateModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Report Title</label><input v-model="newReport.title" type="text" placeholder="Enter report title" class="input-field w-full" /></div>
            <div><label class="block text-[11px] font-medium text-gray-600 mb-1.5">Report Type</label><select v-model="newReport.type" class="input-field w-full"><option v-for="type in reportTypes" :key="type" :value="type">{{ type }}</option></select></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex gap-3 justify-end"><button @click="showGenerateModal = false" class="px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button><button @click="handleGenerateReport" class="px-4 py-2 text-[11px] bg-[#B90B0B] text-white rounded-lg hover:bg-[#991010]">Generate</button></div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showViewModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#B90B0B] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Report Details</h3><button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div><p class="text-[11px] text-gray-500">Report ID</p><p class="text-[13px] font-medium">{{ selectedReport?.id }}</p></div>
              <div><p class="text-[11px] text-gray-500">Format</p><span class="px-2 py-0.5 text-[11px] font-medium rounded" :class="{'bg-red-50 text-red-700': selectedReport?.format === 'PDF', 'bg-green-50 text-green-700': selectedReport?.format === 'Excel'}">{{ selectedReport?.format }}</span></div>
              <div><p class="text-[11px] text-gray-500">Title</p><p class="text-[13px] font-medium">{{ selectedReport?.title }}</p></div>
              <div><p class="text-[11px] text-gray-500">Type</p><p class="text-[13px]">{{ selectedReport?.type }}</p></div>
            </div>
            <div><p class="text-[11px] text-gray-500">Generated By</p><p class="text-[13px]">{{ selectedReport?.generatedBy }}</p></div>
            <div><p class="text-[11px] text-gray-500">Generated At</p><p class="text-[13px]">{{ selectedReport?.generatedAt }}</p></div>
            <div><p class="text-[11px] text-gray-500">File Size</p><p class="text-[13px]">{{ selectedReport?.size }}</p></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-3">
            <button @click="showViewModal = false" class="px-4 py-2 text-[11px] bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200">Close</button>
            <button @click="handleDownload(selectedReport)" class="px-4 py-2 text-[11px] bg-[#B90B0B] text-white rounded-lg hover:bg-[#991010]">Download</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
