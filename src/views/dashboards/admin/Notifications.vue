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

const showViewModal = ref(false)
const selectedNotification = ref<any>(null)
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })

const notifications = ref([
  { id: 'N-001', type: 'System', title: 'System Maintenance Scheduled', message: 'Scheduled maintenance on Jan 20, 2024 at 2:00 AM. Expected downtime: 30 minutes.', isRead: false, createdAt: '2024-01-15 09:30' },
  { id: 'N-002', type: 'Compliance', title: 'New Compliance Alert', message: '15 properties flagged for review in Lagos State.', isRead: false, createdAt: '2024-01-15 08:15' },
  { id: 'N-003', type: 'Payment', title: 'Large Payment Received', message: 'Nigerian Holdings Ltd paid N25,000,000 for property tax.', isRead: true, createdAt: '2024-01-14 16:45' },
  { id: 'N-004', type: 'Audit', title: 'Audit Report Ready', message: 'Q4 2023 audit report has been generated and is ready for review.', isRead: true, createdAt: '2024-01-14 14:30' },
  { id: 'N-005', type: 'Survey', title: 'Survey Submission', message: 'New survey submitted by Surveyor Chukwudi Okafor for review.', isRead: false, createdAt: '2024-01-14 11:00' },
  { id: 'N-006', type: 'User', title: 'New User Registration', message: 'Amina Yusuf registered as Auditor. Awaiting approval.', isRead: false, createdAt: '2024-01-13 10:00' },
  { id: 'N-007', type: 'System', title: 'Database Backup Complete', message: 'Weekly database backup completed successfully.', isRead: true, createdAt: '2024-01-13 06:00' },
  { id: 'N-008', type: 'Compliance', title: 'Compliance Report Due', message: 'Monthly compliance report is due in 5 days.', isRead: true, createdAt: '2024-01-12 09:00' },
])

const types = ['System', 'Compliance', 'Payment', 'Audit', 'Survey', 'User']

const filteredNotifications = computed(() => {
  let items = notifications.value.filter(n => n.title.toLowerCase().includes(searchQuery.value.toLowerCase()) || n.message.toLowerCase().includes(searchQuery.value.toLowerCase()) || n.id.toLowerCase().includes(searchQuery.value.toLowerCase()))
  if (activeTab.value !== 'all') { items = items.filter(n => n.type.toLowerCase() === activeTab.value) }
  return items
})

const totalPages = computed(() => Math.ceil(filteredNotifications.value.length / itemsPerPage.value))
const paginatedNotifications = computed(() => { const start = (currentPage.value - 1) * itemsPerPage.value; return filteredNotifications.value.slice(start, start + itemsPerPage.value) })
const goToPage = (page: number) => { if (page >= 1 && page <= totalPages.value) currentPage.value = page }
const stats = computed(() => ({ total: notifications.value.length, unread: notifications.value.filter(n => !n.isRead).length, system: notifications.value.filter(n => n.type === 'System').length, alerts: notifications.value.filter(n => ['Compliance', 'Audit', 'Payment'].includes(n.type)).length }))

const showToast = (message: string) => { toast.value = { show: true, message }; setTimeout(() => { toast.value.show = false }, 3000) }
const markAsRead = (id: string) => { const n = notifications.value.find(n => n.id === id); if (n) n.isRead = true; showToast('Marked as read') }
const markAllAsRead = () => { notifications.value.forEach(n => n.isRead = true); showToast('All notifications marked as read') }
const openViewModal = (n: any) => { selectedNotification.value = n; if (!n.isRead) markAsRead(n.id); showViewModal.value = true }
</script>

<template>
  <div class="min-h-screen flex bg-[#f5f7fa]">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />
    <div class="flex-1">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6">
        <div class="flex items-center gap-4"><span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span><span class="text-[#d1d5db]">/</span><span class="text-[#1f2937] text-sm font-medium">Notifications</span></div>
        <div class="flex items-center gap-4"><span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span><button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button></div>
      </header>
      <main class="flex-1 p-6">
        <div class="grid grid-cols-4 gap-4 mb-6">
          <div class="bg-white border border-[#e5e7eb] rounded-lg p-4"><p class="text-[11px] text-[#6b7280]">Total Notifications</p><p class="text-2xl font-semibold text-[#1f2937] mt-1">{{ stats.total }}</p></div>
          <div class="bg-white border border-[#e5e7eb] rounded-lg p-4"><p class="text-[11px] text-[#6b7280]">Unread</p><p class="text-2xl font-semibold text-[#2D5A27] mt-1">{{ stats.unread }}</p></div>
          <div class="bg-white border border-[#e5e7eb] rounded-lg p-4"><p class="text-[11px] text-[#6b7280]">System</p><p class="text-2xl font-semibold text-blue-600 mt-1">{{ stats.system }}</p></div>
          <div class="bg-white border border-[#e5e7eb] rounded-lg p-4"><p class="text-[11px] text-[#6b7280]">Alerts</p><p class="text-2xl font-semibold text-yellow-600 mt-1">{{ stats.alerts }}</p></div>
        </div>
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between">
            <div class="flex gap-4">
              <button @click="activeTab = 'all'" class="text-[13px] font-medium pb-1 border-b-2" :class="activeTab === 'all' ? 'border-[#2D5A27] text-[#2D5A27]' : 'border-transparent text-[#6b7280]'">All</button>
              <button v-for="type in types" :key="type" @click="activeTab = type.toLowerCase()" class="text-[13px] font-medium pb-1 border-b-2" :class="activeTab === type.toLowerCase() ? 'border-[#2D5A27] text-[#2D5A27]' : 'border-transparent text-[#6b7280]'">{{ type }}</button>
            </div>
            <button @click="markAllAsRead" class="btn-secondary text-[11px]">Mark All as Read</button>
          </div>
          <div class="p-4 border-b border-[#e5e7eb] flex gap-4"><input v-model="searchQuery" type="text" placeholder="Search notifications..." class="input-field max-w-md" /></div>
          <div class="divide-y divide-[#f3f4f6]">
            <div v-for="notification in paginatedNotifications" :key="notification.id" class="p-4 hover:bg-[#f9fafb]" :class="{ 'bg-blue-50/50': !notification.isRead }">
              <div class="flex items-start justify-between">
                <div class="flex items-start gap-3">
                  <div class="mt-1"><div class="w-2 h-2 rounded-full" :class="notification.isRead ? 'bg-transparent' : 'bg-[#2D5A27]'"></div></div>
                  <div>
                    <div class="flex items-center gap-2">
                      <span class="px-2 py-0.5 text-[11px] font-medium rounded" :class="{'bg-blue-50 text-blue-700': notification.type === 'System', 'bg-yellow-50 text-yellow-700': notification.type === 'Compliance', 'bg-green-50 text-green-700': notification.type === 'Payment', 'bg-purple-50 text-purple-700': notification.type === 'Audit', 'bg-orange-50 text-orange-700': notification.type === 'Survey', 'bg-gray-100 text-gray-700': notification.type === 'User'}">{{ notification.type }}</span>
                      <h4 class="text-[13px] font-medium text-[#1f2937]">{{ notification.title }}</h4>
                    </div>
                    <p class="text-[12px] text-[#6b7280] mt-1">{{ notification.message }}</p>
                    <p class="text-[11px] text-[#9ca3af] mt-2">{{ notification.createdAt }}</p>
                  </div>
                </div>
                <div class="flex gap-2">
                  <button v-if="!notification.isRead" @click="markAsRead(notification.id)" class="px-3 py-1 text-[11px] bg-[#f3f4f6] text-[#374151] rounded hover:bg-[#e5e7eb]">Mark Read</button>
                  <button @click="openViewModal(notification)" class="px-3 py-1 text-[11px] bg-[#2D5A27] text-white rounded hover:bg-[#1e3d1a]">View</button>
                </div>
              </div>
            </div>
          </div>
          <div class="px-6 py-4 flex items-center justify-between border-t border-[#f3f4f6]">
            <p class="text-[11px] text-[#6b7280]">Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, filteredNotifications.length) }} of {{ filteredNotifications.length }} entries</p>
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
      <div v-if="showViewModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md">
          <div class="bg-[#2D5A27] px-6 py-4 flex justify-between items-center"><h3 class="text-base font-semibold text-white">Notification Details</h3><button @click="showViewModal = false" class="text-white/80 hover:text-white">✕</button></div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div><p class="text-[11px] text-gray-500">ID</p><p class="text-[13px] font-medium">{{ selectedNotification?.id }}</p></div>
              <div><p class="text-[11px] text-gray-500">Type</p><span class="px-2 py-0.5 text-[11px] font-medium rounded" :class="{'bg-blue-50 text-blue-700': selectedNotification?.type === 'System', 'bg-yellow-50 text-yellow-700': selectedNotification?.type === 'Compliance'}">{{ selectedNotification?.type }}</span></div>
            </div>
            <div><p class="text-[11px] text-gray-500">Title</p><p class="text-[13px] font-medium">{{ selectedNotification?.title }}</p></div>
            <div><p class="text-[11px] text-gray-500">Message</p><p class="text-[13px]">{{ selectedNotification?.message }}</p></div>
            <div><p class="text-[11px] text-gray-500">Created</p><p class="text-[13px]">{{ selectedNotification?.createdAt }}</p></div>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-end"><button @click="showViewModal = false" class="px-4 py-2 text-[11px] bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200">Close</button></div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
