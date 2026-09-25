import { createRouter, createWebHistory } from 'vue-router'
import WelcomeView from '../views/WelcomeView.vue'
import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import ReceivablesView from '../views/ReceivablesView.vue'
import BusinessView from '../views/BusinessView.vue'
import FinanceView from '../views/FinanceView.vue'
import AuditView from '../views/AuditView.vue'
import DocumentsView from '../views/DocumentsView.vue'
import RiskLabView from '../views/RiskLabView.vue'
import ControlCenterView from '../views/ControlCenterView.vue'
import AISupport from '../components/AISupport.vue'

const routes = [
  { path: '/', redirect: '/welcome' },
  { path: '/welcome', name: 'welcome', component: WelcomeView, meta: { public: true } },
  { path: '/login', name: 'login', component: LoginView, meta: { public: true } },

  { path: '/dashboard', name: 'dashboard', component: DashboardView },
  { path: '/receivables', name: 'receivables', component: ReceivablesView },
  { path: '/business', name: 'business', component: BusinessView, meta: { roles: ['business', 'admin'] } },
  { path: '/finance', name: 'finance', component: FinanceView, meta: { roles: ['finance', 'admin'] } },
  { path: '/audit', name: 'audit', component: AuditView, meta: { roles: ['admin'] } },
  { path: '/documents', name: 'documents', component: DocumentsView, meta: { roles: ['business', 'finance', 'admin'] } },
  { path: '/risk-lab', name: 'risk-lab', component: RiskLabView, meta: { roles: ['finance', 'admin'] } },
  { path: '/ai-assistant', name: 'ai-assistant', component: AISupport, meta: { roles: ['business', 'finance', 'admin'] } },
  { path: '/control-center', name: 'control-center', component: ControlCenterView, meta: { roles: ['admin'] } }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to) => {
  const raw = localStorage.getItem('scf_user')
  const user = raw ? JSON.parse(raw) : null
  const role = (user?.role || '').toLowerCase()

  if (to.meta.public) {
    return true
  }
  if (!user) {
    return '/login'
  }
  if (to.meta.roles && !to.meta.roles.includes(role)) {
    if (role === 'business') return '/business'
    if (role === 'finance') return '/finance'
    return '/dashboard'
  }
  return true
})

export default router
