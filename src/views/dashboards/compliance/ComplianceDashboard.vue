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
  labels: ['Undeclared', 'Evasion', 'Forgery', 'Discrepancy', 'Other'],
  datasets: [{
    label: 'Cases',
    data: [25, 18, 12, 20, 14],
    backgroundColor: '#2D5A27',
    borderRadius: 2,
    borderSkipped: false,
  }],
}

const doughnutData = {
  labels: ['Resolved', 'Under Review', 'Active', 'Pending'],
  datasets: [{
    data: [72, 12, 3, 2],
    backgroundColor: ['#16a34a', '#f59e0b', '#2D5A27', '#d1d5db'],
    borderWidth: 0,
    cutout: '70%',
  }],
}

const lineData = {
  labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'],
  datasets: [
    { label: 'Resolved', data: [10, 15, 12, 18, 22, 25], borderColor: '#16a34a', backgroundColor: 'rgba(22, 163, 74, 0.1)', fill: true, tension: 0.3, pointRadius: 0 },
    { label: 'New Cases', data: [8, 12, 14, 16, 18, 20], borderColor: '#2D5A27', backgroundColor: 'rgba(45, 90, 39, 0.1)', fill: true, tension: 0.3, pointRadius: 0 },
  ],
}

const stats = [
  { label: 'Total Cases', value: '89' },
  { label: 'Under Investigation', value: '12' },
  { label: 'Resolved', value: '72' },
  { label: 'Escalated', value: '5' },
]

const recentCases = [
  { id: 'CASE-001', title: 'Undeclared Property', location: 'Lekki Phase 2', priority: 'High', status: 'Active' },
  { id: 'CASE-002', title: 'Tax Evasion Suspected', location: 'Victoria Island', priority: 'Critical', status: 'Under Review' },
  { id: 'CASE-003', title: 'Document Forgery', location: 'Ikoyi', priority: 'Medium', status: 'Resolved' },
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
            <h3 class="section-title">Monthly Activity</h3>
            <div class="h-48">
              <Line :data="lineData" :options="{ responsive: true, maintainAspectRatio: false, plugins: { legend: { position: 'bottom', labels: { boxWidth: 10, padding: 12 } } }, scales: { x: { grid: { display: false } }, y: { grid: { color: '#f3f4f6' } } } }" />
            </div>
          </div>

          <div class="card p-5">
            <h3 class="section-title">Cases by Type</h3>
            <div class="h-40 flex items-center justify-center">
              <Bar :data="barData" :options="{ responsive: true, maintainAspectRatio: false, indexAxis: 'y', plugins: { legend: { display: false } }, scales: { x: { grid: { color: '#f3f4f6' } }, y: { grid: { display: false } } } }" />
            </div>
          </div>
        </div>

        <div class="grid grid-cols-3 gap-4 mb-6">
          <div class="card p-5 col-span-2">
            <h3 class="section-title">Case Status</h3>
            <div class="h-40 flex items-center justify-center">
              <Doughnut :data="doughnutData" :options="{ responsive: true, maintainAspectRatio: false, plugins: { legend: { position: 'bottom', labels: { boxWidth: 10, padding: 12 } } } }" />
            </div>
          </div>

          <div class="card">
            <div class="px-5 py-4 border-b border-[#e5e7eb]">
              <h3 class="section-title mb-0">Recent Cases</h3>
            </div>
            <div class="p-4">
              <div v-for="c in recentCases" :key="c.id" class="py-3 border-b border-[#f3f4f6] last:border-0">
                <p class="text-sm font-medium text-[#1f2937]">{{ c.title }}</p>
                <p class="text-xs text-[#6b7280] mt-0.5">{{ c.location }}</p>
                <div class="flex items-center gap-2 mt-2">
                  <span class="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-medium"
                    :class="{
                      'bg-red-50 text-red-700': c.priority === 'Critical' || c.priority === 'High',
                      'bg-yellow-50 text-yellow-700': c.priority === 'Medium',
                    }">{{ c.priority }}</span>
                  <span class="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-medium"
                    :class="{
                      'bg-red-50 text-red-700': c.status === 'Active',
                      'bg-yellow-50 text-yellow-700': c.status === 'Under Review',
                      'bg-green-50 text-green-700': c.status === 'Resolved',
                    }">{{ c.status }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>
