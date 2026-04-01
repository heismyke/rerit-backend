<script setup lang="ts">
import { Bar, Doughnut, Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
  ArcElement,
  PointElement,
  LineElement,
} from 'chart.js'
import { useRoleStore } from '@/stores'
import { useRouter } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend, ArcElement, PointElement, LineElement)

const { selectedRole, user, logout } = useRoleStore()
const router = useRouter()

const handleLogout = () => {
  logout()
  router.push('/')
}

const barData = {
  labels: ['Land', 'Boundary', 'Topo', 'Construction', 'Subdivision'],
  datasets: [{
    label: 'Surveys',
    data: [15, 12, 10, 5, 3],
    backgroundColor: '#B90B0B',
    borderRadius: 2,
    borderSkipped: false,
  }],
}

const doughnutData = {
  labels: ['Approved', 'Under Review', 'Pending'],
  datasets: [{
    data: [38, 5, 2],
    backgroundColor: ['#16a34a', '#f59e0b', '#d1d5db'],
    borderWidth: 0,
    cutout: '70%',
  }],
}

const lineData = {
  labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'],
  datasets: [{
    label: 'Submissions',
    data: [8, 12, 6, 15, 10, 12],
    borderColor: '#B90B0B',
    backgroundColor: 'rgba(220, 38, 38, 0.1)',
    fill: true,
    tension: 0.3,
    pointRadius: 0,
  }],
}

const stats = [
  { label: 'Total Surveys', value: '45' },
  { label: 'Pending Review', value: '5' },
  { label: 'Approved', value: '38' },
  { label: 'This Month', value: '12' },
]

const recentSurveys = [
  { id: 'SURV-001', property: 'Plot 15, Ikoyi', type: 'Land Survey', status: 'Approved', date: 'Jan 10' },
  { id: 'SURV-002', property: 'Block 3, Lekki', type: 'Boundary Survey', status: 'Under Review', date: 'Jan 12' },
  { id: 'SURV-003', property: 'Plot 88, VI', type: 'Topographic Survey', status: 'Approved', date: 'Jan 08' },
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

        <div class="grid grid-cols-3 gap-4 mb-6">
          <div class="card p-5 col-span-2">
            <h3 class="section-title">Monthly Submissions</h3>
            <div class="h-48">
              <Line :data="lineData" :options="{ responsive: true, maintainAspectRatio: false, plugins: { legend: { display: false } }, scales: { x: { grid: { display: false } }, y: { grid: { color: '#f3f4f6' } } } }" />
            </div>
          </div>

          <div class="card p-5">
            <h3 class="section-title">Surveys by Type</h3>
            <div class="h-40 flex items-center justify-center">
              <Bar :data="barData" :options="{ responsive: true, maintainAspectRatio: false, indexAxis: 'y', plugins: { legend: { display: false } }, scales: { x: { grid: { color: '#f3f4f6' } }, y: { grid: { display: false } } } }" />
            </div>
          </div>
        </div>

        <div class="card mb-6">
          <div class="px-5 py-4 border-b border-[#e5e7eb]">
            <h3 class="section-title mb-0">Approval Status</h3>
          </div>
          <div class="p-5 h-48 flex items-center justify-center">
            <Doughnut :data="doughnutData" :options="{ responsive: true, maintainAspectRatio: false, plugins: { legend: { position: 'bottom', labels: { boxWidth: 10, padding: 12 } } } }" />
          </div>
        </div>

        <div class="card">
          <div class="px-5 py-4 border-b border-[#e5e7eb]">
            <h3 class="section-title mb-0">Recent Surveys</h3>
          </div>
          <table class="w-full">
            <thead>
              <tr>
                <th class="table-header">Survey ID</th>
                <th class="table-header">Property</th>
                <th class="table-header">Type</th>
                <th class="table-header">Status</th>
                <th class="table-header">Date</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-[#f3f4f6]">
              <tr v-for="survey in recentSurveys" :key="survey.id">
                <td class="table-cell font-medium">{{ survey.id }}</td>
                <td class="table-cell text-[#6b7280]">{{ survey.property }}</td>
                <td class="table-cell">{{ survey.type }}</td>
                <td class="table-cell">
                  <span class="inline-flex items-center px-2 py-0.5 rounded text-[11px] font-medium"
                    :class="{
                      'bg-green-50 text-green-700': survey.status === 'Approved',
                      'bg-yellow-50 text-yellow-700': survey.status === 'Under Review',
                    }">{{ survey.status }}</span>
                </td>
                <td class="table-cell text-[#9ca3af]">{{ survey.date }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </main>
    </div>
  </div>
</template>
