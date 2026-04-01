<script setup lang="ts">
import { Bar, Line } from 'vue-chartjs'
import { Chart as ChartJS, CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend, PointElement, LineElement } from 'chart.js'
import { useRoleStore } from '@/stores'
import { useRouter } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import { ref } from 'vue'

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend, PointElement, LineElement)

const { selectedRole, user, logout } = useRoleStore()
const router = useRouter()
const handleLogout = () => { logout(); router.push('/') }

const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })
const showExportModal = ref(false)

const revenueData = {
  labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
  datasets: [
    { label: 'Collected', data: [450, 520, 380, 650, 720, 580, 800, 750, 680, 850, 900, 950], borderColor: '#16a34a', backgroundColor: 'rgba(22, 163, 74, 0.1)', fill: true, tension: 0.3 },
    { label: 'Outstanding', data: [120, 110, 130, 100, 90, 110, 80, 95, 85, 70, 60, 50], borderColor: '#B90B0B', backgroundColor: 'rgba(185, 11, 11, 0.1)', fill: true, tension: 0.3 },
  ],
}

const payments = ref([
  { id: 'PAY-001', taxpayer: 'Emeka Okonkwo', property: 'PROP-001', amount: 'N2,500,000', date: '2024-01-15', method: 'Bank Transfer', status: 'Completed' },
  { id: 'PAY-002', taxpayer: 'Adaobi Nnamdi', property: 'PROP-002', amount: 'N800,000', date: '2024-01-14', method: 'Card', status: 'Completed' },
  { id: 'PAY-003', taxpayer: 'Chidi Okafor', property: 'PROP-003', amount: 'N1,800,000', date: '2024-01-13', method: 'USSD', status: 'Pending' },
])

const showToast = (message: string) => {
  toast.value = { show: true, message }
  setTimeout(() => { toast.value.show = false }, 3000)
}

const handleExportReport = (format: string) => {
  showExportModal.value = false
  showToast(`Report exported as ${format.toUpperCase()}`)
}

const viewPayment = (payment: any) => {
  showToast(`Viewing payment ${payment.id} for ${payment.taxpayer}`)
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
          <span class="text-[#1f2937] text-sm font-medium">Revenue & Payments</span>
        </div>
        <div class="flex items-center gap-4">
          <button @click="showExportModal = true" class="px-3 py-1.5 text-[11px] bg-green-600 text-white rounded hover:bg-green-700">Export Report</button>
          <span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span>
          <button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button>
        </div>
      </header>
      <main class="flex-1 p-6">
        <div class="grid grid-cols-4 gap-4 mb-6">
          <div class="stat-card"><p class="metric-label">Total Collected</p><p class="metric-value">N4.5B</p></div>
          <div class="stat-card"><p class="metric-label">This Month</p><p class="metric-value">N950M</p></div>
          <div class="stat-card"><p class="metric-label">Outstanding</p><p class="metric-value">N120M</p></div>
          <div class="stat-card"><p class="metric-label">Transactions</p><p class="metric-value">12.8K</p></div>
        </div>
        <div class="card p-5 mb-6">
          <h3 class="section-title">Revenue Trend</h3>
          <div class="h-64"><Line :data="revenueData" :options="{ responsive: true, maintainAspectRatio: false, plugins: { legend: { position: 'bottom' } }, scales: { x: { grid: { display: false } }, y: { grid: { color: '#f3f4f6' } } } }" /></div>
        </div>
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between"><h2 class="text-[13px] font-semibold text-[#1f2937]">Recent Payments</h2><button @click="showExportModal = true" class="btn-primary text-[11px]">Record Payment</button></div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr>
                  <th class="table-header">ID</th>
                  <th class="table-header">Taxpayer</th>
                  <th class="table-header">Property</th>
                  <th class="table-header">Amount</th>
                  <th class="table-header">Date</th>
                  <th class="table-header">Method</th>
                  <th class="table-header">Status</th>
                  <th class="table-header">Actions</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-[#f3f4f6]">
                <tr v-for="p in payments" :key="p.id" class="hover:bg-[#f9fafb]">
                  <td class="table-cell font-medium">{{ p.id }}</td>
                  <td class="table-cell">{{ p.taxpayer }}</td>
                  <td class="table-cell text-[#6b7280]">{{ p.property }}</td>
                  <td class="table-cell">{{ p.amount }}</td>
                  <td class="table-cell text-[#9ca3af]">{{ p.date }}</td>
                  <td class="table-cell text-[#6b7280]">{{ p.method }}</td>
                  <td class="table-cell">
                    <span class="px-2 py-0.5 text-[11px] font-medium rounded-full" :class="{'bg-green-50 text-green-700': p.status === 'Completed', 'bg-yellow-50 text-yellow-700': p.status === 'Pending'}">{{ p.status }}</span>
                  </td>
                  <td class="table-cell">
                    <div class="flex gap-2">
                      <button @click="viewPayment(p)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">View</button>
                      <button class="px-3 py-1 text-[11px] bg-[#B90B0B] text-white rounded hover:bg-[#991010]">Receipt</button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </main>
    </div>

    <Teleport to="body"><div v-if="toast.show" class="fixed bottom-4 right-4 bg-green-600 text-white px-4 py-2 rounded-lg shadow-lg z-50">{{ toast.message }}</div></Teleport>

    <Teleport to="body">
      <div v-if="showExportModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-sm">
          <div class="bg-[#B90B0B] px-6 py-4 flex justify-between items-center">
            <h3 class="text-base font-semibold text-white">Export / Actions</h3>
            <button @click="showExportModal = false" class="text-white/80 hover:text-white">✕</button>
          </div>
          <div class="p-6 space-y-3">
            <button @click="handleExportReport('pdf')" class="w-full px-4 py-3 text-left text-[13px] border border-gray-200 rounded-lg hover:bg-gray-50 flex items-center gap-3">
              <span class="text-red-600">📄</span> Export as PDF
            </button>
            <button @click="handleExportReport('excel')" class="w-full px-4 py-3 text-left text-[13px] border border-gray-200 rounded-lg hover:bg-gray-50 flex items-center gap-3">
              <span class="text-green-600">📊</span> Export as Excel
            </button>
            <button @click="handleExportReport('csv')" class="w-full px-4 py-3 text-left text-[13px] border border-gray-200 rounded-lg hover:bg-gray-50 flex items-center gap-3">
              <span class="text-blue-600">📋</span> Export as CSV
            </button>
          </div>
          <div class="px-6 py-4 border-t border-gray-100"><button @click="showExportModal = false" class="w-full px-4 py-2 text-[11px] border border-gray-300 rounded-lg hover:bg-gray-50">Cancel</button></div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
