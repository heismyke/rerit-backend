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

const activeSection = ref('general')

const generalSettings = ref({
  systemName: 'ReRiT - Real Estate Revenue & Information System',
  organization: 'National Revenue System (NRS)',
  address: 'Federal Ministry of Finance, Abuja, Nigeria',
  phone: '+234 9 234 5678',
  email: 'support@nrs.gov.ng',
  timezone: 'Africa/Lagos',
  fiscalYearStart: 'January',
})

const taxSettings = ref({
  defaultTaxRate: '0.5',
  latePaymentPenalty: '5',
  penaltyFrequency: 'monthly',
  minimumPropertyValue: '100000',
  currency: 'NGN',
})

const notificationSettings = ref({
  emailNotifications: true,
  smsNotifications: false,
  systemAlerts: true,
  auditReminders: true,
  complianceAlerts: true,
})

const securitySettings = ref({
  twoFactorAuth: false,
  sessionTimeout: '30',
  passwordExpiry: '90',
  ipWhitelist: '',
})
</script>

<template>
  <div class="min-h-screen flex bg-[#f5f7fa]">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />

    <div class="flex-1">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6">
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
          <div class="px-6 py-4 border-b border-[#e5e7eb] flex items-center justify-between">
            <h2 class="text-[13px] font-semibold text-[#1f2937]">System Settings</h2>
            <button class="btn-primary text-[11px]">Save Changes</button>
          </div>
          
          <div class="flex border-b border-[#e5e7eb]">
            <button @click="activeSection = 'general'" class="px-6 py-3 text-[13px] font-medium border-b-2" :class="activeSection === 'general' ? 'border-[#2D5A27] text-[#2D5A27]' : 'border-transparent text-[#6b7280]'">General</button>
            <button @click="activeSection = 'tax'" class="px-6 py-3 text-[13px] font-medium border-b-2" :class="activeSection === 'tax' ? 'border-[#2D5A27] text-[#2D5A27]' : 'border-transparent text-[#6b7280]'">Tax Configuration</button>
            <button @click="activeSection = 'notifications'" class="px-6 py-3 text-[13px] font-medium border-b-2" :class="activeSection === 'notifications' ? 'border-[#2D5A27] text-[#2D5A27]' : 'border-transparent text-[#6b7280]'">Notifications</button>
            <button @click="activeSection = 'security'" class="px-6 py-3 text-[13px] font-medium border-b-2" :class="activeSection === 'security' ? 'border-[#2D5A27] text-[#2D5A27]' : 'border-transparent text-[#6b7280]'">Security</button>
          </div>

          <div class="p-6">
            <div v-if="activeSection === 'general'" class="space-y-6">
              <div class="grid grid-cols-2 gap-6">
                <div>
                  <label class="block text-[11px] font-medium text-[#374151] mb-1.5">System Name</label>
                  <input v-model="generalSettings.systemName" type="text" class="input-field w-full" />
                </div>
                <div>
                  <label class="block text-[11px] font-medium text-[#374151] mb-1.5">Organization</label>
                  <input v-model="generalSettings.organization" type="text" class="input-field w-full" />
                </div>
                <div>
                  <label class="block text-[11px] font-medium text-[#374151] mb-1.5">Address</label>
                  <input v-model="generalSettings.address" type="text" class="input-field w-full" />
                </div>
                <div>
                  <label class="block text-[11px] font-medium text-[#374151] mb-1.5">Phone</label>
                  <input v-model="generalSettings.phone" type="text" class="input-field w-full" />
                </div>
                <div>
                  <label class="block text-[11px] font-medium text-[#374151] mb-1.5">Email</label>
                  <input v-model="generalSettings.email" type="email" class="input-field w-full" />
                </div>
                <div>
                  <label class="block text-[11px] font-medium text-[#374151] mb-1.5">Timezone</label>
                  <select v-model="generalSettings.timezone" class="input-field w-full">
                    <option value="Africa/Lagos">Africa/Lagos (WAT)</option>
                    <option value="Africa/Abuja">Africa/Abuja (WAT)</option>
                  </select>
                </div>
                <div>
                  <label class="block text-[11px] font-medium text-[#374151] mb-1.5">Fiscal Year Start</label>
                  <select v-model="generalSettings.fiscalYearStart" class="input-field w-full">
                    <option>January</option>
                    <option>April</option>
                  </select>
                </div>
              </div>
            </div>

            <div v-if="activeSection === 'tax'" class="space-y-6">
              <div class="grid grid-cols-2 gap-6">
                <div>
                  <label class="block text-[11px] font-medium text-[#374151] mb-1.5">Default Tax Rate (%)</label>
                  <input v-model="taxSettings.defaultTaxRate" type="number" step="0.1" class="input-field w-full" />
                </div>
                <div>
                  <label class="block text-[11px] font-medium text-[#374151] mb-1.5">Late Payment Penalty (%)</label>
                  <input v-model="taxSettings.latePaymentPenalty" type="number" class="input-field w-full" />
                </div>
                <div>
                  <label class="block text-[11px] font-medium text-[#374151] mb-1.5">Penalty Frequency</label>
                  <select v-model="taxSettings.penaltyFrequency" class="input-field w-full">
                    <option value="monthly">Monthly</option>
                    <option value="quarterly">Quarterly</option>
                    <option value="annually">Annually</option>
                  </select>
                </div>
                <div>
                  <label class="block text-[11px] font-medium text-[#374151] mb-1.5">Minimum Property Value (NGN)</label>
                  <input v-model="taxSettings.minimumPropertyValue" type="number" class="input-field w-full" />
                </div>
                <div>
                  <label class="block text-[11px] font-medium text-[#374151] mb-1.5">Currency</label>
                  <select v-model="taxSettings.currency" class="input-field w-full">
                    <option value="NGN">Nigerian Naira (NGN)</option>
                  </select>
                </div>
              </div>
            </div>

            <div v-if="activeSection === 'notifications'" class="space-y-4">
              <div class="flex items-center justify-between py-3 border-b border-[#f3f4f6]">
                <div>
                  <p class="text-[13px] font-medium text-[#1f2937]">Email Notifications</p>
                  <p class="text-[11px] text-[#6b7280]">Receive notifications via email</p>
                </div>
                <label class="relative inline-flex items-center cursor-pointer">
                  <input v-model="notificationSettings.emailNotifications" type="checkbox" class="sr-only peer">
                  <div class="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-[#2D5A27]"></div>
                </label>
              </div>
              <div class="flex items-center justify-between py-3 border-b border-[#f3f4f6]">
                <div>
                  <p class="text-[13px] font-medium text-[#1f2937]">SMS Notifications</p>
                  <p class="text-[11px] text-[#6b7280]">Receive notifications via SMS</p>
                </div>
                <label class="relative inline-flex items-center cursor-pointer">
                  <input v-model="notificationSettings.smsNotifications" type="checkbox" class="sr-only peer">
                  <div class="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-[#2D5A27]"></div>
                </label>
              </div>
              <div class="flex items-center justify-between py-3 border-b border-[#f3f4f6]">
                <div>
                  <p class="text-[13px] font-medium text-[#1f2937]">System Alerts</p>
                  <p class="text-[11px] text-[#6b7280]">Receive system maintenance and update alerts</p>
                </div>
                <label class="relative inline-flex items-center cursor-pointer">
                  <input v-model="notificationSettings.systemAlerts" type="checkbox" class="sr-only peer">
                  <div class="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-[#2D5A27]"></div>
                </label>
              </div>
              <div class="flex items-center justify-between py-3 border-b border-[#f3f4f6]">
                <div>
                  <p class="text-[13px] font-medium text-[#1f2937]">Audit Reminders</p>
                  <p class="text-[11px] text-[#6b7280]">Receive reminders for scheduled audits</p>
                </div>
                <label class="relative inline-flex items-center cursor-pointer">
                  <input v-model="notificationSettings.auditReminders" type="checkbox" class="sr-only peer">
                  <div class="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-[#2D5A27]"></div>
                </label>
              </div>
              <div class="flex items-center justify-between py-3">
                <div>
                  <p class="text-[13px] font-medium text-[#1f2937]">Compliance Alerts</p>
                  <p class="text-[11px] text-[#6b7280]">Receive alerts for compliance issues</p>
                </div>
                <label class="relative inline-flex items-center cursor-pointer">
                  <input v-model="notificationSettings.complianceAlerts" type="checkbox" class="sr-only peer">
                  <div class="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-[#2D5A27]"></div>
                </label>
              </div>
            </div>

            <div v-if="activeSection === 'security'" class="space-y-6">
              <div class="grid grid-cols-2 gap-6">
                <div>
                  <label class="block text-[11px] font-medium text-[#374151] mb-1.5">Session Timeout (minutes)</label>
                  <input v-model="securitySettings.sessionTimeout" type="number" class="input-field w-full" />
                </div>
                <div>
                  <label class="block text-[11px] font-medium text-[#374151] mb-1.5">Password Expiry (days)</label>
                  <input v-model="securitySettings.passwordExpiry" type="number" class="input-field w-full" />
                </div>
                <div>
                  <label class="block text-[11px] font-medium text-[#374151] mb-1.5">IP Whitelist (comma-separated)</label>
                  <input v-model="securitySettings.ipWhitelist" type="text" placeholder="e.g., 192.168.1.1, 10.0.0.1" class="input-field w-full" />
                </div>
              </div>
              <div class="flex items-center justify-between py-3 border-t border-[#f3f4f6]">
                <div>
                  <p class="text-[13px] font-medium text-[#1f2937]">Two-Factor Authentication</p>
                  <p class="text-[11px] text-[#6b7280]">Require 2FA for all admin users</p>
                </div>
                <label class="relative inline-flex items-center cursor-pointer">
                  <input v-model="securitySettings.twoFactorAuth" type="checkbox" class="sr-only peer">
                  <div class="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-[#2D5A27]"></div>
                </label>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>
