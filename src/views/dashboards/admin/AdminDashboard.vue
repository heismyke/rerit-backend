<script setup lang="ts">
import { Bar, Line, Doughnut } from 'vue-chartjs'
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
import { getRiskColor } from '@/utils/riskScoring'

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend, ArcElement, PointElement, LineElement)

const { selectedRole, user, logout } = useRoleStore()
const router = useRouter()

const handleLogout = () => {
  logout()
  router.push('/')
}

const lineData = {
  labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
  datasets: [
    { label: 'Tax Revenue', data: [45, 52, 38, 65, 72, 58, 80, 75, 68, 85, 90, 95], borderColor: '#2D5A27', backgroundColor: 'rgba(45, 90, 39, 0.1)', fill: true, tension: 0.3, pointRadius: 3 },
  ],
}

const barData = {
  labels: ['Lagos', 'Abuja', 'Rivers', 'Kano', 'Oyo', 'Delta'],
  datasets: [{
    label: 'Properties',
    data: [1245, 892, 567, 423, 356, 298],
    backgroundColor: '#2D5A27',
    borderRadius: 4,
  }],
}

const doughnutData = {
  labels: ['Verified', 'Pending', 'Flagged', 'Under Review'],
  datasets: [{
    data: [4567, 1234, 456, 789],
    backgroundColor: ['#16a34a', '#f59e0b', '#2D5A27', '#6366f1'],
    borderWidth: 0,
    cutout: '65%',
  }],
}

const stats = [
  { label: 'Total Properties', value: '12,847' },
  { label: 'Tax Collected', value: 'N4.5B' },
  { label: 'Active Cases', value: '1,234' },
  { label: 'Flagged', value: '456' },
]

const riskDistribution = [
  { level: 'Critical', count: 89, color: 'bg-red-500' },
  { level: 'High', count: 234, color: 'bg-orange-500' },
  { level: 'Medium', count: 567, color: 'bg-yellow-500' },
  { level: 'Low', count: 1892, color: 'bg-green-500' },
]

const highRiskProperties = [
  { id: 'PROP-004', owner: 'Folake Adeyemi', address: 'Plot 8, Banana Island', riskScore: 85, riskLevel: 'Critical' },
  { id: 'PROP-012', owner: 'Chukwuemeka Obi', address: 'Block 15, Victoria Island', riskScore: 78, riskLevel: 'Critical' },
  { id: 'PROP-023', owner: 'Ngozi Adebayo', address: 'Plot 5, Maitama', riskScore: 72, riskLevel: 'High' },
]

const recentActivity = [
  { id: 1, action: 'New property registered', user: 'Emeka Okonkwo', time: '2 mins ago', type: 'property' },
  { id: 2, action: 'Tax payment received', user: 'Adaobi Nnamdi', time: '15 mins ago', type: 'payment' },
  { id: 3, action: 'Case flagged for review', user: 'System', time: '32 mins ago', type: 'case' },
  { id: 4, action: 'Survey submitted', user: 'Aisha Ibrahim', time: '1 hour ago', type: 'survey' },
  { id: 5, action: 'Rule triggered', user: 'Auto-system', time: '2 hours ago', type: 'rule' },
]

const getActivityColor = (type: string) => {
  switch (type) {
    case 'property': return 'bg-blue-50 text-blue-700'
    case 'payment': return 'bg-green-50 text-green-700'
    case 'case': return 'bg-red-50 text-red-700'
    case 'survey': return 'bg-purple-50 text-purple-700'
    case 'rule': return 'bg-yellow-50 text-yellow-700'
    default: return 'bg-gray-50 text-gray-700'
  }
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
          <span class="text-[#1f2937] text-sm font-medium">Dashboard</span>
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
            <h3 class="section-title">Revenue Over Time</h3>
            <div class="h-64">
              <Line :data="lineData" :options="{ responsive: true, maintainAspectRatio: false, plugins: { legend: { display: true, position: 'bottom', labels: { boxWidth: 10, padding: 12 } } }, scales: { x: { grid: { display: false } }, y: { grid: { color: '#f3f4f6' } } } }" />
            </div>
          </div>

          <div class="card p-5">
            <h3 class="section-title">Property Status</h3>
            <div class="h-48 flex items-center justify-center">
              <Doughnut :data="doughnutData" :options="{ responsive: true, maintainAspectRatio: false, plugins: { legend: { position: 'bottom', labels: { boxWidth: 10, padding: 12 } } } }" />
            </div>
          </div>
        </div>

        <div class="grid grid-cols-3 gap-4 mb-6">
          <div class="card p-5">
            <h3 class="section-title">Properties by Region</h3>
            <div class="h-48">
              <Bar :data="barData" :options="{ responsive: true, maintainAspectRatio: false, plugins: { legend: { display: false } }, scales: { x: { grid: { display: false } }, y: { grid: { color: '#f3f4f6' } } } }" />
            </div>
          </div>

          <div class="card p-5">
            <div class="flex items-center justify-between mb-4">
              <h3 class="section-title mb-0">Risk Distribution</h3>
              <span class="text-[11px] text-[#2D5A27] font-medium">2,782 Assessed</span>
            </div>
            <div class="space-y-3">
              <div v-for="risk in riskDistribution" :key="risk.level" class="flex items-center gap-3">
                <div class="w-16 text-[11px] font-medium" :class="risk.level === 'Critical' ? 'text-red-600' : risk.level === 'High' ? 'text-orange-600' : risk.level === 'Medium' ? 'text-yellow-600' : 'text-green-600'">{{ risk.level }}</div>
                <div class="flex-1 bg-gray-100 rounded-full h-2">
                  <div class="h-2 rounded-full" :class="risk.color" :style="{ width: (risk.count / 1892 * 100) + '%' }"></div>
                </div>
                <div class="w-10 text-[11px] text-right font-medium text-gray-600">{{ risk.count }}</div>
              </div>
            </div>
          </div>

          <div class="card p-5">
            <div class="flex items-center justify-between mb-4">
              <h3 class="section-title mb-0">High Risk Properties</h3>
              <button @click="router.push('/admin/properties')" class="text-[11px] text-[#2D5A27] hover:underline">View All</button>
            </div>
            <div class="space-y-2">
              <div v-for="prop in highRiskProperties" :key="prop.id" class="flex items-center justify-between py-2 border-b border-[#f3f4f6] last:border-0">
                <div class="flex-1">
                  <p class="text-[12px] font-medium text-gray-800">{{ prop.id }}</p>
                  <p class="text-[10px] text-gray-500">{{ prop.owner }}</p>
                </div>
                <div class="flex items-center gap-2">
                  <span class="text-[11px] font-bold" :class="prop.riskScore >= 75 ? 'text-red-600' : 'text-orange-600'">{{ prop.riskScore }}</span>
                  <span class="px-1.5 py-0.5 text-[10px] font-medium rounded" :class="getRiskColor(prop.riskLevel)">{{ prop.riskLevel }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-4 gap-4">
          <div class="card p-4 text-center hover:shadow-md transition cursor-pointer" @click="router.push('/admin/properties')">
            <div class="w-10 h-10 bg-blue-50 rounded-full flex items-center justify-center mx-auto mb-2">
              <span class="text-lg text-blue-600">◈</span>
            </div>
            <p class="text-[11px] text-[#6b7280]">Properties</p>
            <p class="text-lg font-semibold text-[#1f2937]">12.8K</p>
          </div>
          <div class="card p-4 text-center hover:shadow-md transition cursor-pointer" @click="router.push('/admin/taxpayers')">
            <div class="w-10 h-10 bg-purple-50 rounded-full flex items-center justify-center mx-auto mb-2">
              <span class="text-lg text-purple-600">◎</span>
            </div>
            <p class="text-[11px] text-[#6b7280]">Taxpayers</p>
            <p class="text-lg font-semibold text-[#1f2937]">8.4K</p>
          </div>
          <div class="card p-4 text-center hover:shadow-md transition cursor-pointer" @click="router.push('/admin/revenue')">
            <div class="w-10 h-10 bg-green-50 rounded-full flex items-center justify-center mx-auto mb-2">
              <span class="text-lg text-green-600">◫</span>
            </div>
            <p class="text-[11px] text-[#6b7280]">Revenue</p>
            <p class="text-lg font-semibold text-[#1f2937]">N4.5B</p>
          </div>
          <div class="card p-4 text-center hover:shadow-md transition cursor-pointer" @click="router.push('/admin/cases')">
            <div class="w-10 h-10 bg-red-50 rounded-full flex items-center justify-center mx-auto mb-2">
              <span class="text-lg text-red-600">◉</span>
            </div>
            <p class="text-[11px] text-[#6b7280]">Active Cases</p>
            <p class="text-lg font-semibold text-[#1f2937]">1.2K</p>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>
