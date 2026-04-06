<script setup lang="ts">
import { useRoleStore } from '@/stores'
import { useRouter } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import { ref } from 'vue'

const { selectedRole, user, logout } = useRoleStore()
const router = useRouter()

const handleLogout = () => {
  logout()
  router.push('/')
}

const activeTab = ref('profile')

const profile = ref({
  name: 'John Smith',
  email: 'j.smith@nrs.gov.ng',
  employeeId: 'EMP-2024-001',
  department: 'Audit Department',
  role: 'Senior Auditor',
})

const notifications = ref({ email: true, sms: false, push: true })
const security = ref({ twoFactor: false, sessionTimeout: '30' })
</script>

<template>
  <div class="min-h-screen bg-[#f5f6fa] flex">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />

    <div class="flex-1 flex flex-col">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6 shrink-0">
        <div class="flex items-center gap-4">
          <span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span>
          <span class="text-[#d1d5db]">/</span>
          <span class="text-[#1f2937] text-sm font-medium">Settings</span>
        </div>
        <div class="flex items-center gap-4">
          <span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span>
          <button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button>
        </div>
      </header>

      <main class="flex-1 p-6">
        <div class="bg-white border border-[#e5e7eb] rounded-lg">
          <div class="border-b border-[#e5e7eb]">
            <nav class="flex">
              <button @click="activeTab = 'profile'" class="px-5 py-3 text-[11px] font-medium border-b-2" :class="activeTab === 'profile' ? 'border-[#2D5A27] text-[#2D5A27]' : 'border-transparent text-[#6b7280] hover:text-[#1f2937]'">Profile</button>
              <button @click="activeTab = 'notifications'" class="px-5 py-3 text-[11px] font-medium border-b-2" :class="activeTab === 'notifications' ? 'border-[#2D5A27] text-[#2D5A27]' : 'border-transparent text-[#6b7280] hover:text-[#1f2937]'">Notifications</button>
              <button @click="activeTab = 'security'" class="px-5 py-3 text-[11px] font-medium border-b-2" :class="activeTab === 'security' ? 'border-[#2D5A27] text-[#2D5A27]' : 'border-transparent text-[#6b7280] hover:text-[#1f2937]'">Security</button>
            </nav>
          </div>

          <div class="p-6">
            <div v-if="activeTab === 'profile'" class="space-y-6">
              <div class="flex items-center gap-5 pb-6 border-b border-[#e5e7eb]">
                <div class="w-16 h-16 bg-[#2D5A27] rounded-full flex items-center justify-center text-xl font-bold text-white">JS</div>
                <button class="btn-secondary text-[11px]">Change Photo</button>
              </div>
              <div class="grid grid-cols-2 gap-5">
                <div><label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Full Name</label><input v-model="profile.name" type="text" class="input-field" /></div>
                <div><label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Email</label><input v-model="profile.email" type="email" class="input-field" /></div>
                <div><label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Employee ID</label><input v-model="profile.employeeId" type="text" disabled class="input-field bg-[#f9fafb]" /></div>
                <div><label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Department</label><input v-model="profile.department" type="text" class="input-field" /></div>
              </div>
              <div class="flex justify-end"><button class="btn-primary text-[11px]">Save Changes</button></div>
            </div>

            <div v-if="activeTab === 'notifications'" class="space-y-0">
              <div class="flex items-center justify-between py-4 border-b border-[#e5e7eb]">
                <div><p class="text-sm font-medium text-[#1f2937]">Email Notifications</p><p class="text-xs text-[#9ca3af]">Receive audit updates via email</p></div>
                <label class="relative inline-flex items-center cursor-pointer"><input v-model="notifications.email" type="checkbox" class="sr-only peer" /><div class="w-10 h-5 bg-[#d1d5db] peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-[#d1d5db] after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-[#2D5A27]"></div></label>
              </div>
              <div class="flex items-center justify-between py-4 border-b border-[#e5e7eb]">
                <div><p class="text-sm font-medium text-[#1f2937]">SMS Notifications</p><p class="text-xs text-[#9ca3af]">Urgent alerts via SMS</p></div>
                <label class="relative inline-flex items-center cursor-pointer"><input v-model="notifications.sms" type="checkbox" class="sr-only peer" /><div class="w-10 h-5 bg-[#d1d5db] peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-[#d1d5db] after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-[#2D5A27]"></div></label>
              </div>
              <div class="flex items-center justify-between py-4">
                <div><p class="text-sm font-medium text-[#1f2937]">Push Notifications</p><p class="text-xs text-[#9ca3af]">Real-time browser alerts</p></div>
                <label class="relative inline-flex items-center cursor-pointer"><input v-model="notifications.push" type="checkbox" class="sr-only peer" /><div class="w-10 h-5 bg-[#d1d5db] peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-[#d1d5db] after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-[#2D5A27]"></div></label>
              </div>
            </div>

            <div v-if="activeTab === 'security'" class="space-y-6">
              <div class="flex items-center justify-between py-4 border-b border-[#e5e7eb]">
                <div><p class="text-sm font-medium text-[#1f2937]">Two-Factor Authentication</p><p class="text-xs text-[#9ca3af]">Add extra security to your account</p></div>
                <label class="relative inline-flex items-center cursor-pointer"><input v-model="security.twoFactor" type="checkbox" class="sr-only peer" /><div class="w-10 h-5 bg-[#d1d5db] peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-[#d1d5db] after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-[#2D5A27]"></div></label>
              </div>
              <div class="py-4 border-b border-[#e5e7eb]">
                <label class="block text-[11px] font-medium text-[#6b7280] mb-2">Session Timeout</label>
                <select v-model="security.sessionTimeout" class="input-field w-48">
                  <option value="15">15 minutes</option>
                  <option value="30">30 minutes</option>
                  <option value="60">1 hour</option>
                </select>
              </div>
              <div class="flex gap-3">
                <button class="btn-secondary text-[11px]">Change Password</button>
                <button class="btn-secondary text-[11px]">View Login History</button>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>
