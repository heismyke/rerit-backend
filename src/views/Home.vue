<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { roles, getRouteByRoleId } from '@/data'
import { useRoleStore } from '@/stores'

const router = useRouter()
const { setRole } = useRoleStore()

const showLoginModal = ref(false)
const selectedRoleForLogin = ref<typeof roles[number] | null>(null)
const email = ref('')
const password = ref('')
const showForgotPassword = ref(false)

const handleRoleSelect = (role: typeof roles[number]) => {
  selectedRoleForLogin.value = role
  setRole(role)
  showLoginModal.value = true
}

const handleLogin = () => {
  const { login } = useRoleStore()
  login(email.value, 'User')
  showLoginModal.value = false
  router.push(getRouteByRoleId(selectedRoleForLogin.value!.id))
}

const closeModal = () => {
  showLoginModal.value = false
  showForgotPassword.value = false
  email.value = ''
  password.value = ''
}
</script>

<template>
  <div class="min-h-screen flex">
    <div class="hidden lg:flex lg:w-1/2 relative overflow-hidden">
      <div class="absolute inset-0">
        <img src="/tax.jpg" alt="" class="w-full h-full object-cover" />
        <div class="absolute inset-0 bg-gradient-to-r from-black/60 to-black/30"></div>
      </div>
      <div class="relative z-10 flex items-center justify-center p-12 w-full">
        <div class="max-w-lg">
          <h1 class="text-5xl font-bold text-white mb-4 tracking-tight">ReRiT</h1>
          <p class="text-xl font-light text-white/80 mb-6">Real Estate Revenue & Information System</p>
          <p class="text-white/60 leading-relaxed">
            Streamlining property tax collection, audit management, and compliance tracking for a smarter Nigeria.
          </p>
          <div class="mt-12 grid grid-cols-2 gap-4">
            <div class="bg-white/10 backdrop-blur rounded-lg p-4 border border-white/10">
              <p class="text-2xl font-bold text-white">10,000+</p>
              <p class="text-xs text-white/70 mt-1">Properties Managed</p>
            </div>
            <div class="bg-white/10 backdrop-blur rounded-lg p-4 border border-white/10">
              <p class="text-2xl font-bold text-white">N50B+</p>
              <p class="text-xs text-white/70 mt-1">Tax Revenue Collected</p>
            </div>
            <div class="bg-white/10 backdrop-blur rounded-lg p-4 border border-white/10">
              <p class="text-2xl font-bold text-white">500+</p>
              <p class="text-xs text-white/70 mt-1">Surveyors</p>
            </div>
            <div class="bg-white/10 backdrop-blur rounded-lg p-4 border border-white/10">
              <p class="text-2xl font-bold text-white">36</p>
              <p class="text-xs text-white/70 mt-1">States Covered</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="w-full lg:w-1/2 flex items-center justify-center p-8 bg-[#f5f7fa]">
      <div class="w-full max-w-md">
        <div class="lg:hidden text-center mb-8">
          <h1 class="text-3xl font-bold text-[#2D5A27]">ReRiT</h1>
          <p class="text-sm text-gray-500 mt-1">Revenue System</p>
        </div>

        <div class="mb-6">
          <img src="/fct-irs.svg" alt="FCT-IRS" class="h-16" />
        </div>

        <h2 class="text-xl font-semibold text-[#1f2937] mb-1">Welcome</h2>
        <p class="text-sm text-gray-500 mb-8">Select your role to continue</p>

        <div class="space-y-3">
          <button
            v-for="role in roles"
            :key="role.id"
            @click="handleRoleSelect(role)"
            class="w-full bg-white border border-gray-200 rounded-xl p-4 hover:border-[#2D5A27] hover:bg-red-50/30 transition-all text-left flex items-center gap-4 shadow-md"
          >
            <span class="text-xl text-[#2D5A27]">{{ role.icon }}</span>
            <div>
              <h3 class="text-sm font-medium text-[#1f2937]">{{ role.name }}</h3>
              <p class="text-xs text-gray-500">{{ role.description }}</p>
            </div>
          </button>
        </div>
      </div>
    </div>

    <Teleport to="body">
      <div v-if="showLoginModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md overflow-hidden">
          <div class="bg-[#2D5A27] px-6 py-4">
            <h3 class="text-base font-semibold text-white">Sign In</h3>
            <p class="text-xs text-white/80 mt-0.5">{{ selectedRoleForLogin?.name }}</p>
          </div>

          <div class="p-6">
            <template v-if="!showForgotPassword">
              <div class="space-y-4">
                <div>
                  <label class="block text-xs font-medium text-gray-600 mb-1.5">Email Address</label>
                  <input
                    v-model="email"
                    type="email"
                    placeholder="Enter your email"
                    class="input-field"
                  />
                </div>
                <div>
                  <label class="block text-xs font-medium text-gray-600 mb-1.5">Password</label>
                  <input
                    v-model="password"
                    type="password"
                    placeholder="Enter your password"
                    class="input-field"
                  />
                </div>
                <button
                  @click="showForgotPassword = true"
                  class="text-xs text-[#2D5A27] hover:text-[#6a0707] font-medium"
                >
                  Forgot Password?
                </button>
                <button @click="handleLogin" class="btn-primary w-full">Sign In</button>
              </div>
            </template>

            <template v-else>
              <div class="space-y-4">
                <p class="text-sm text-gray-600">Enter your email address and we'll send you a link to reset your password.</p>
                <div>
                  <label class="block text-xs font-medium text-gray-600 mb-1.5">Email Address</label>
                  <input v-model="email" type="email" placeholder="Enter your email" class="input-field" />
                </div>
                <button @click="showForgotPassword = false" class="text-xs text-[#2D5A27] font-medium">Back to Login</button>
                <button @click="showForgotPassword = false" class="btn-primary w-full">Send Reset Link</button>
              </div>
            </template>

            <button @click="closeModal" class="mt-4 w-full text-center text-xs text-gray-500 hover:text-gray-700 py-2">Cancel</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
