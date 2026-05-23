import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', name: 'login', component: () => import('../views/LoginView.vue'), meta: { public: true } },
    { path: '/setup', name: 'setup', component: () => import('../views/SetupView.vue'), meta: { public: true } },
    { path: '/', name: 'dashboard', component: () => import('../views/DashboardView.vue') },
    { path: '/entries', name: 'entries', component: () => import('../views/EntryListView.vue') },
    { path: '/entries/new', name: 'entry-create', component: () => import('../views/EntryCreateView.vue') },
    { path: '/entries/:id', name: 'entry-detail', component: () => import('../views/EntryDetailView.vue') }
  ]
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()

  if (auth.initialized === null) {
    await auth.checkStatus()
  }

  if (!auth.initialized && to.name !== 'setup') {
    return { name: 'setup' }
  }

  if (auth.initialized && to.name === 'setup') {
    return { name: 'login' }
  }

  if (!to.meta.public && !auth.isLoggedIn()) {
    return { name: 'login' }
  }
})

export default router
