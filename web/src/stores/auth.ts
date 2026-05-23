import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '../api/client'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const initialized = ref<boolean | null>(null)

  const isLoggedIn = () => !!token.value

  async function checkStatus() {
    const res = await api.getStatus()
    initialized.value = res.initialized
    return res.initialized
  }

  async function login(username: string, password: string) {
    const res = await api.login(username, password)
    token.value = res.token
    localStorage.setItem('token', res.token)
  }

  async function setup(username: string, password: string) {
    await api.setup(username, password)
    initialized.value = true
  }

  function logout() {
    token.value = ''
    localStorage.removeItem('token')
  }

  return { token, initialized, isLoggedIn, checkStatus, login, setup, logout }
})
