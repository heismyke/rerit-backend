<script setup lang="ts">
import { useRoleStore } from '@/stores'
import { useRouter } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import { ref } from 'vue'

const { selectedRole, user, logout } = useRoleStore()
const router = useRouter()
const handleLogout = () => { logout(); router.push('/') }

const selectedPaymentMethod = ref('card')
const selectedProperty = ref('PROP-001')
const amount = ref('2,500,000')
const toast = ref<{ show: boolean; message: string }>({ show: false, message: '' })

const paymentMethods = [
  { id: 'card', name: 'Bank Card', icon: '💳', description: 'Pay with debit/credit card' },
  { id: 'transfer', name: 'Bank Transfer', icon: '🏦', description: 'Direct bank transfer' },
  { id: 'ussd', name: 'USSD Code', icon: '📱', description: 'Pay via USSD code' },
]

const properties = [
  { id: 'PROP-001', name: 'Commercial Complex A', amount: 'N2,500,000' },
  { id: 'PROP-002', name: 'Residential Estate B', amount: 'N1,800,000' },
  { id: 'PROP-003', name: 'Mixed Use Development C', amount: 'N3,500,000' },
]

const showToast = (message: string) => { toast.value = { show: true, message }; setTimeout(() => { toast.value.show = false }, 3000) }

const makePayment = () => {
  showToast('Payment processed successfully! Receipt sent to your email.')
}

const updateAmount = () => {
  const prop = properties.find(p => p.id === selectedProperty.value)
  if (prop) amount.value = prop.amount.replace('N', '').replace(/,/g, '')
}
</script>

<template>
  <div class="min-h-screen bg-[#f5f6fa] flex">
    <Sidebar v-if="selectedRole?.id" :role-id="selectedRole.id" />
    <div class="flex-1 flex flex-col">
      <header class="h-14 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6 shrink-0">
        <div class="flex items-center gap-4"><span class="text-[#6b7280] text-sm">{{ selectedRole?.name }}</span><span class="text-[#d1d5db]">/</span><span class="text-[#1f2937] text-sm font-medium">Make Payment</span></div>
        <div class="flex items-center gap-4"><span class="text-[11px] text-[#9ca3af]">{{ user?.email }}</span><button @click="handleLogout" class="btn-ghost text-[11px]">Logout</button></div>
      </header>
      <main class="flex-1 p-6">
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div class="lg:col-span-2 space-y-6">
            <div class="bg-white border border-[#e5e7eb] rounded-lg p-5">
              <h3 class="text-[13px] font-semibold text-[#1f2937] mb-4">Select Property</h3>
              <select v-model="selectedProperty" @change="updateAmount" class="input-field">
                <option v-for="prop in properties" :key="prop.id" :value="prop.id">{{ prop.name }} - {{ prop.amount }}</option>
              </select>
            </div>
            <div class="bg-white border border-[#e5e7eb] rounded-lg p-5">
              <h3 class="text-[13px] font-semibold text-[#1f2937] mb-4">Payment Amount</h3>
              <div class="flex items-center gap-2"><span class="text-[11px] text-[#6b7280]">N</span><input v-model="amount" type="text" class="input-field" /></div>
            </div>
            <div class="bg-white border border-[#e5e7eb] rounded-lg p-5">
              <h3 class="text-[13px] font-semibold text-[#1f2937] mb-4">Payment Method</h3>
              <div class="space-y-3">
                <label v-for="method in paymentMethods" :key="method.id" class="flex items-center gap-4 p-4 border rounded cursor-pointer transition" :class="selectedPaymentMethod === method.id ? 'border-[#B90B0B] bg-red-50/30' : 'border-[#e5e7eb] hover:border-[#d1d5db]'">
                  <input v-model="selectedPaymentMethod" type="radio" :value="method.id" class="w-4 h-4 text-[#B90B0B]" />
                  <span class="text-lg">{{ method.icon }}</span>
                  <div><p class="text-[13px] font-medium text-[#1f2937]">{{ method.name }}</p><p class="text-[11px] text-[#9ca3af]">{{ method.description }}</p></div>
                </label>
              </div>
            </div>
            <div v-if="selectedPaymentMethod === 'card'" class="bg-white border border-[#e5e7eb] rounded-lg p-5">
              <h3 class="text-[13px] font-semibold text-[#1f2937] mb-4">Card Details</h3>
              <div class="space-y-4">
                <div><label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Card Number</label><input type="text" placeholder="1234 5678 9012 3456" class="input-field" /></div>
                <div class="grid grid-cols-2 gap-4">
                  <div><label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">Expiry Date</label><input type="text" placeholder="MM/YY" class="input-field" /></div>
                  <div><label class="block text-[11px] font-medium text-[#6b7280] mb-1.5">CVV</label><input type="text" placeholder="123" class="input-field" /></div>
                </div>
              </div>
            </div>
            <div v-if="selectedPaymentMethod === 'transfer'" class="bg-white border border-[#e5e7eb] rounded-lg p-5">
              <h3 class="text-[13px] font-semibold text-[#1f2937] mb-4">Bank Transfer Details</h3>
              <div class="space-y-2 text-[13px]">
                <p><span class="text-[#6b7280]">Bank Name:</span> <span class="font-medium">First Bank of Nigeria</span></p>
                <p><span class="text-[#6b7280]">Account Name:</span> <span class="font-medium">ReRiT Tax Collection</span></p>
                <p><span class="text-[#6b7280]">Account Number:</span> <span class="font-medium">3084721934</span></p>
                <p><span class="text-[#6b7280]">Reference:</span> <span class="font-medium">{{ selectedProperty }}-2024</span></p>
              </div>
            </div>
            <div v-if="selectedPaymentMethod === 'ussd'" class="bg-white border border-[#e5e7eb] rounded-lg p-5">
              <h3 class="text-[13px] font-semibold text-[#1f2937] mb-4">USSD Payment</h3>
              <div class="text-center p-4 bg-[#f9fafb] rounded"><p class="text-[11px] text-[#9ca3af] mb-2">Dial this code on your registered phone number</p><p class="text-2xl font-bold text-[#B90B0B]">*901*000*1234#</p></div>
            </div>
            <button @click="makePayment" class="btn-primary w-full py-3">Proceed to Pay N{{ Number(amount).toLocaleString() }}</button>
          </div>
          <div class="space-y-6">
            <div class="bg-white border border-[#e5e7eb] rounded-lg p-5">
              <h3 class="text-[13px] font-semibold text-[#1f2937] mb-4">Payment Summary</h3>
              <div class="space-y-3 text-[13px]">
                <div class="flex justify-between"><span class="text-[#6b7280]">Property</span><span class="font-medium">{{ selectedProperty }}</span></div>
                <div class="flex justify-between"><span class="text-[#6b7280]">Amount</span><span class="font-bold text-lg">N{{ Number(amount).toLocaleString() }}</span></div>
                <div class="flex justify-between"><span class="text-[#6b7280]">Processing Fee</span><span class="font-medium">N0.00</span></div>
                <div class="border-t border-[#e5e7eb] pt-3 flex justify-between"><span class="text-[#1f2937] font-medium">Total</span><span class="font-bold text-[#B90B0B]">N{{ Number(amount).toLocaleString() }}</span></div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>

    <Teleport to="body"><div v-if="toast.show" class="fixed bottom-4 right-4 bg-green-600 text-white px-4 py-2 rounded-lg shadow-lg z-50">{{ toast.message }}</div></Teleport>
  </div>
</template>
