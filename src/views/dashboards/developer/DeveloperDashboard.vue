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
  labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'],
  datasets: [{
    label: 'Tax Paid',
    data: [2.0, 2.5, 1.8, 3.0, 2.2, 3.5],
    backgroundColor: '#2D5A27',
    borderRadius: 2,
    borderSkipped: false,
  }],
}

const doughnutData = {
  labels: ['Properties', 'Outstanding', 'Pending'],
  datasets: [{
    data: [3, 1, 1],
    backgroundColor: ['#2D5A27', '#f59e0b', '#d1d5db'],
    borderWidth: 0,
    cutout: '70%',
  }],
}

const lineData = {
  labels: ['Q1 23', 'Q2 23', 'Q3 23', 'Q4 23', 'Q1 24'],
  datasets: [{
    label: 'Property Value',
    data: [45, 48, 52, 58, 65],
    borderColor: '#2D5A27',
    backgroundColor: 'rgba(45, 90, 39, 0.1)',
    fill: true,
    tension: 0.3,
    pointRadius: 0,
  }],
}

const stats = [
  { label: 'Total Projects', value: '12' },
  { label: 'Compliance Rate', value: '94%' },
  { label: 'Active Cases', value: '2' },
  { label: 'Amount Due', value: 'N4.2M' },
]

const notifications = [
  { title: 'Tax Assessment Notice', message: 'Q1 2024 assessment issued', time: '2h ago' },
  { title: 'Payment Confirmed', message: 'Property #PROP-001 payment received', time: '1d ago' },
  { title: 'Compliance Review', message: 'Annual review scheduled for Mar 15', time: '3d ago' },
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
            <h3 class="section-title">Payment Trend</h3>
            <div class="h-48">
              <Bar :data="barData" :options="{ responsive: true, maintainAspectRatio: false, plugins: { legend: { display: false } }, scales: { x: { grid: { display: false } }, y: { grid: { color: '#f3f4f6' } } } }" />
            </div>
          </div>

          <div class="card p-5">
            <h3 class="section-title">Portfolio Overview</h3>
            <div class="h-40 flex items-center justify-center">
              <Doughnut :data="doughnutData" :options="{ responsive: true, maintainAspectRatio: false, plugins: { legend: { position: 'bottom', labels: { boxWidth: 10, padding: 12 } } } }" />
            </div>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-4 mb-6">
          <div class="card">
            <div class="px-5 py-4 border-b border-[#e5e7eb]">
              <h3 class="section-title mb-0">Recent Notifications</h3>
            </div>
            <div class="p-4">
              <div v-for="notif in notifications" :key="notif.title" class="py-3 border-b border-[#f3f4f6] last:border-0">
                <p class="text-sm font-medium text-[#1f2937]">{{ notif.title }}</p>
                <p class="text-xs text-[#6b7280] mt-0.5">{{ notif.message }}</p>
                <p class="text-[11px] text-[#9ca3af] mt-1">{{ notif.time }}</p>
              </div>
            </div>
          </div>

          <div class="card">
            <div class="px-5 py-4 border-b border-[#e5e7eb]">
              <h3 class="section-title mb-0">Quick Actions</h3>
            </div>
            <div class="p-4 grid grid-cols-3 gap-3">
              <button @click="router.push('/developer/notices')" class="p-4 border border-[#e5e7eb] rounded text-center hover:border-[#2D5A27] hover:bg-red-50/30 transition">
                <p class="text-xs font-medium text-[#1f2937]">Respond to Notice</p>
                <p class="text-[11px] text-[#9ca3af] mt-1">View & respond</p>
              </button>
              <button @click="router.push('/developer/payments')" class="p-4 border border-[#e5e7eb] rounded text-center hover:border-[#2D5A27] hover:bg-red-50/30 transition">
                <p class="text-xs font-medium text-[#1f2937]">Make Payment</p>
                <p class="text-[11px] text-[#9ca3af] mt-1">Pay tax dues</p>
              </button>
              <button class="p-4 border border-[#e5e7eb] rounded text-center hover:border-[#2D5A27] hover:bg-red-50/30 transition">
                <p class="text-xs font-medium text-[#1f2937]">Request Refund</p>
                <p class="text-[11px] text-[#9ca3af] mt-1">If applicable</p>
              </button>
            </div>
          </div>
        </div>

        <div class="card">
          <div class="px-5 py-4 border-b border-[#e5e7eb]">
            <h3 class="section-title mb-0">Property Value Trend</h3>
          </div>
          <div class="p-5 h-48">
            <Line :data="lineData" :options="{ responsive: true, maintainAspectRatio: false, plugins: { legend: { display: false } }, scales: { x: { grid: { display: false } }, y: { grid: { color: '#f3f4f6' } } } }" />
          </div>
        </div>
      </main>
    </div>
  </div>
</template>
