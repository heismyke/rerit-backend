<script setup lang="ts">
import { Bar, Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
  PointElement,
  LineElement,
} from 'chart.js'
import { useRoleStore } from '@/stores'
import { useRouter } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend, PointElement, LineElement)

const { selectedRole, user, logout } = useRoleStore()
const router = useRouter()

const handleLogout = () => {
  logout()
  router.push('/')
}

const lineData = {
  labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
  datasets: [
    { label: 'Audit Cases', data: [45, 52, 38, 65, 72, 58, 80, 75, 68, 85, 90, 95], borderColor: '#B90B0B', backgroundColor: 'rgba(185, 11, 11, 0.1)', fill: true, tension: 0.3, pointRadius: 3 },
  ],
}

const histogramData = {
  labels: ['Lagos', 'Abuja', 'Rivers', 'Kano', 'Oyo', 'Delta'],
  datasets: [{
    label: 'High Risk Cases',
    data: [24, 18, 15, 12, 10, 8],
    backgroundColor: '#B90B0B',
    borderRadius: 4,
  }],
}

const stats = [
  { label: 'Total Audit Cases', value: '1,247' },
  { label: 'Pending Review', value: '156' },
  { label: 'Under Investigation', value: '89' },
  { label: 'Resolved Cases', value: '1,002' },
]

const recentAudits = [
  { id: 'AUD-001', property: 'Plot 42, Victoria Island', owner: 'Emeka Okonkwo', status: 'Under Investigation', date: 'Jan 15', priority: 'High' },
  { id: 'AUD-002', property: 'Block 7, Lekki Phase 2', owner: 'Adaobi Nnamdi', status: 'Pending Review', date: 'Jan 14', priority: 'Medium' },
  { id: 'AUD-003', property: '15 Admiralty Way, Lekki', owner: 'Chidi Okafor', status: 'Resolved', date: 'Jan 13', priority: 'Low' },
  { id: 'AUD-004', property: 'Plot 8, Banana Island', owner: 'Folake Adeyemi', status: 'Under Investigation', date: 'Jan 12', priority: 'High' },
  { id: 'AUD-005', property: 'Block 12, Maitama', owner: 'Global Ventures Ltd', status: 'Pending Review', date: 'Jan 11', priority: 'Critical' },
]
</script>

<template>
  <div class="min-h-screen flex bg-[#f5f7fa]">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />

    <div class="flex-1">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6">
        <div class="flex items-center gap-4">
          <span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span>
          <span class="text-[#d1d5db]">/</span>
          <span class="text-[#1f2937] text-sm font-medium">Overview</span>
        </div>
        <div class="flex items-center gap-4">
          <span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span>
          <button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button>
        </div>
      </header>

      <main class="p-6">
        <div class="grid grid-cols-4 gap-4 mb-6">
          <div v-for="stat in stats" :key="stat.label" class="stat-card">
            <p class="metric-label">{{ stat.label }}</p>
            <p class="metric-value">{{ stat.value }}</p>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-4 mb-6">
          <div class="card p-5">
            <h3 class="section-title">Audit Cases Over Time</h3>
            <div class="h-64">
              <Line :data="lineData" :options="{ responsive: true, maintainAspectRatio: false, plugins: { legend: { display: true, position: 'bottom', labels: { boxWidth: 10, padding: 12 } } }, scales: { x: { grid: { display: false } }, y: { grid: { color: '#f3f4f6' } } } }" />
            </div>
          </div>

          <div class="card p-5">
            <h3 class="section-title">High Risk Cases by Region</h3>
            <div class="h-64">
              <Bar :data="histogramData" :options="{ responsive: true, maintainAspectRatio: false, plugins: { legend: { display: false } }, scales: { x: { grid: { display: false } }, y: { grid: { color: '#f3f4f6' } } } }" />
            </div>
          </div>
        </div>

        <div class="card">
          <div class="px-5 py-4 border-b border-[#e5e7eb]">
            <h3 class="section-title mb-0">Recent Audit Cases</h3>
          </div>
          <table class="w-full">
            <thead>
              <tr>
                <th class="table-header">Audit ID</th>
                <th class="table-header">Property</th>
                <th class="table-header">Owner</th>
                <th class="table-header">Status</th>
                <th class="table-header">Priority</th>
                <th class="table-header">Date</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-[#f3f4f6]">
              <tr v-for="audit in recentAudits" :key="audit.id">
                <td class="table-cell font-medium">{{ audit.id }}</td>
                <td class="table-cell text-[#6b7280]">{{ audit.property }}</td>
                <td class="table-cell">{{ audit.owner }}</td>
                <td class="table-cell">
                  <span class="inline-flex items-center px-2 py-0.5 rounded text-[11px] font-medium"
                    :class="{
                      'bg-yellow-50 text-yellow-700': audit.status === 'Pending Review',
                      'bg-blue-50 text-blue-700': audit.status === 'Under Investigation',
                      'bg-green-50 text-green-700': audit.status === 'Resolved',
                    }">{{ audit.status }}</span>
                </td>
                <td class="table-cell">
                  <span class="inline-flex items-center px-2 py-0.5 rounded text-[11px] font-medium"
                    :class="{
                      'bg-red-100 text-red-700': audit.priority === 'Critical' || audit.priority === 'High',
                      'bg-yellow-100 text-yellow-700': audit.priority === 'Medium',
                      'bg-gray-100 text-gray-600': audit.priority === 'Low',
                    }">{{ audit.priority }}</span>
                </td>
                <td class="table-cell text-[#9ca3af]">{{ audit.date }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </main>
    </div>
  </div>
</template>
