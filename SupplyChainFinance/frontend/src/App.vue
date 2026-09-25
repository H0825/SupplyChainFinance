<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AppSidebar from './components/AppSidebar.vue'
import { useAuthStore } from './stores/auth'
import { requestPreferDB } from './api/http'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()
const sidebarCollapsed = ref(localStorage.getItem('scf_sidebar_collapsed') === '1')
const themeMode = ref(localStorage.getItem('scf_theme_mode') || 'light')
const showOnboarding = ref(false)
const workflowSnapshot = ref({ enterprises: [], receivables: [], documents: [] })
const workflowLoading = ref(false)

const onboardingMap = {
  business: {
    title: '企业端新手引导',
    intro: '按业务顺序完成企业登记、应收创建、确权与结算，不需要记忆复杂ID。',
    steps: [
      { title: '1. 企业登记', desc: '先进入企业业务页，登记企业身份与属性。' },
      { title: '2. 创建应收', desc: '系统可自动生成应收ID，并关联付款方与凭证。' },
      { title: '3. 确权与结算', desc: '后续步骤会从已有应收中自动筛选候选记录。' }
    ],
    startRoute: '/business'
  },
  finance: {
    title: '金融端新手引导',
    intro: '先看总览，再做融资处理与风控核验，系统会自动联动可融资应收。',
    steps: [
      { title: '1. 金融总览', desc: '优先查看链上状态、交易趋势和业务规模。' },
      { title: '2. 融资处理', desc: '单笔融资可直接从可融资应收中选择。' },
      { title: '3. 风控与复核', desc: '去风控模型和应收池做贷前复核。' }
    ],
    startRoute: '/finance'
  },
  admin: {
    title: '管理端新手引导',
    intro: '管理端负责总览、链上监控、用户映射和审计追踪。',
    steps: [
      { title: '1. 平台总览', desc: '先看全局经营与链上运行指标。' },
      { title: '2. 管理控制台', desc: '检查用户、角色与企业映射关系。' },
      { title: '3. 审计中心', desc: '排查异常交易和权限使用记录。' }
    ],
    startRoute: '/dashboard'
  }
}

const baseWorkflowMap = {
  business: [
    { route: '/dashboard', title: '看运营总览', desc: '先确认当前链上状态和业务规模。', nextRoute: '/business', nextLabel: '去企业业务' },
    { route: '/business', title: '完成企业登记', desc: '先登记企业身份，后续流程才更顺。', nextRoute: '/business', nextLabel: '去登记企业' },
    { route: '/business', title: '创建应收与凭证', desc: '系统会自动生成应收ID，并联动凭证上传。', nextRoute: '/business', nextLabel: '去创建应收' },
    { route: '/business', title: '完成确权与结算', desc: '待确认和可结算候选会自动筛出。', nextRoute: '/business', nextLabel: '去处理业务' }
  ],
  finance: [
    { route: '/dashboard', title: '看金融总览', desc: '优先确认链状态、应收规模和近期交易。', nextRoute: '/finance', nextLabel: '去融资处理' },
    { route: '/finance', title: '完成融资登记', desc: '从可融资应收中选择并录入授信。', nextRoute: '/finance', nextLabel: '去融资处理' },
    { route: '/risk-lab', title: '做风控核查', desc: '结合风控模型与应收池做贷前筛查。', nextRoute: '/risk-lab', nextLabel: '去风控模型' },
    { route: '/receivables', title: '复核应收资产', desc: '回看主信息和融资记录，确认流程闭环。', nextRoute: '/receivables', nextLabel: '查看应收池' }
  ],
  admin: [
    { route: '/dashboard', title: '看平台总览', desc: '先掌握系统业务分布与链上运行。', nextRoute: '/dashboard', nextLabel: '查看总览' },
    { route: '/control-center', title: '检查映射关系', desc: '核对用户、企业和链账户关系。', nextRoute: '/control-center', nextLabel: '去管理控制台' },
    { route: '/audit', title: '查看审计数据', desc: '追踪异常操作与近期关键交易。', nextRoute: '/audit', nextLabel: '去审计中心' },
    { route: '/documents', title: '核对文档归档', desc: '检查凭证文档上传与CID是否完整。', nextRoute: '/documents', nextLabel: '看文档中心' }
  ]
}

const isPublicPage = computed(() => Boolean(route.meta?.public))
const appClasses = computed(() => [
  'app-shell',
  themeMode.value === 'dark' ? 'theme-dark' : 'theme-light',
  sidebarCollapsed.value ? 'sidebar-collapsed' : ''
])
const onboardingContent = computed(() => onboardingMap[auth.role] || onboardingMap.admin)
const workflowSteps = computed(() => buildWorkflowSteps(auth.role, workflowSnapshot.value, route.path, auth.user))
const workflowGuide = computed(() => workflowSteps.value.find((item) => item.state === 'current') || workflowSteps.value.find((item) => item.state === 'next') || workflowSteps.value[0])
const completedStepCount = computed(() => workflowSteps.value.filter((item) => item.state === 'done').length)

function go(to) {
  router.push(to)
}

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

function toggleTheme() {
  themeMode.value = themeMode.value === 'dark' ? 'light' : 'dark'
}

function onboardingKey(role) {
  return `scf_onboarding_seen_${role}`
}

function visitKey(role) {
  return `scf_route_visits_${role}`
}

function getVisitedRoutes(role) {
  try {
    return JSON.parse(localStorage.getItem(visitKey(role)) || '[]')
  } catch (_) {
    return []
  }
}

function markRouteVisited(role, path) {
  if (!role || role === 'guest' || !path) return
  const visited = new Set(getVisitedRoutes(role))
  visited.add(path)
  localStorage.setItem(visitKey(role), JSON.stringify([...visited]))
}

function dismissOnboarding() {
  localStorage.setItem(onboardingKey(auth.role), '1')
  showOnboarding.value = false
}

function startOnboardingFlow() {
  dismissOnboarding()
  router.push(onboardingContent.value.startRoute)
}

function buildWorkflowSteps(role, snapshot, currentPath, user) {
  const templates = baseWorkflowMap[role] || baseWorkflowMap.admin
  const visited = new Set(getVisitedRoutes(role))
  const receivables = snapshot.receivables || []
  const enterprises = snapshot.enterprises || []
  const documents = snapshot.documents || []
  const address = String(user?.address || '').toLowerCase()

  let completedFlags = []
  if (role === 'business') {
    const ownEnterprise = enterprises.some((item) => String(item.wallet || '').toLowerCase() === address)
    const ownReceivables = receivables.filter((item) => String(item.issuer || '').toLowerCase() === address || !address)
    completedFlags = [
      visited.has('/dashboard'),
      ownEnterprise,
      ownReceivables.length > 0 || documents.length > 0,
      ownReceivables.some((item) => ['CONFIRMED', 'FINANCED', 'SETTLED'].includes(item.statusText))
    ]
  } else if (role === 'finance') {
    completedFlags = [
      visited.has('/dashboard'),
      receivables.some((item) => item.statusText === 'FINANCED' || item.statusText === 'SETTLED'),
      visited.has('/risk-lab'),
      visited.has('/receivables')
    ]
  } else {
    completedFlags = [
      visited.has('/dashboard'),
      visited.has('/control-center'),
      visited.has('/audit'),
      visited.has('/documents')
    ]
  }

  const currentIndex = completedFlags.findIndex((flag) => !flag)
  return templates.map((item, index) => {
    let state = 'todo'
    if (completedFlags[index]) state = 'done'
    else if (currentIndex === index || currentPath === item.route) state = 'current'
    else if (currentIndex !== -1 && index === currentIndex + 1) state = 'next'
    return { ...item, state }
  })
}

async function loadWorkflowSnapshot() {
  if (!auth.user || isPublicPage.value) return
  workflowLoading.value = true
  try {
    const [enterprises, receivables, documents] = await Promise.all([
      requestPreferDB('/api/db/enterprises', '/api/enterprise').catch(() => ({ data: [] })),
      requestPreferDB('/api/db/receivables', '/api/receivable', { withRole: true }).catch(() => ({ data: [] })),
      requestPreferDB('/api/db/documents', '/api/db/documents', { withRole: true }).catch(() => ({ data: [] }))
    ])
    workflowSnapshot.value = {
      enterprises: enterprises?.data || [],
      receivables: receivables?.data || [],
      documents: documents?.data || []
    }
  } finally {
    workflowLoading.value = false
  }
}

watch(sidebarCollapsed, (value) => {
  localStorage.setItem('scf_sidebar_collapsed', value ? '1' : '0')
})

watch(themeMode, (value) => {
  localStorage.setItem('scf_theme_mode', value)
  document.documentElement.setAttribute('data-theme', value)
})

watch(
  () => auth.role,
  (role) => {
    if (!role || role === 'guest') return
    showOnboarding.value = localStorage.getItem(onboardingKey(role)) !== '1'
  },
  { immediate: true }
)

watch(
  () => [auth.role, route.path, auth.user?.address],
  async ([role, path]) => {
    if (!role || role === 'guest' || isPublicPage.value) return
    markRouteVisited(role, path)
    await loadWorkflowSnapshot()
  },
  { immediate: true }
)

onMounted(() => {
  document.documentElement.setAttribute('data-theme', themeMode.value)
})
</script>

<template>
  <div :class="appClasses">
    <div class="app-bg"></div>
    <div v-if="isPublicPage" class="public-shell">
      <RouterView v-slot="{ Component }">
        <Transition name="page-fade-slide" mode="out-in">
          <component :is="Component" />
        </Transition>
      </RouterView>
    </div>

    <div v-else class="layout">
      <AppSidebar
        :collapsed="sidebarCollapsed"
        :theme-mode="themeMode"
        @toggle-sidebar="toggleSidebar"
        @toggle-theme="toggleTheme"
      />
      <main class="content">
        <header class="topbar">
          <div>
            <strong>{{ auth.user?.username || '访客' }}</strong>
            <div class="sub">{{ auth.user?.address || '请先登录后访问控制台' }}</div>
          </div>
          <div class="row">
            <button class="btn mini ghost" @click="go('/dashboard')">总览</button>
            <button class="btn mini ghost" @click="go('/receivables')">应收池</button>
          </div>
        </header>

        <section class="workflow-banner card">
          <div class="workflow-copy">
            <div class="workflow-badge">流程导航</div>
            <strong>{{ workflowGuide?.title }}</strong>
            <div class="sub">{{ workflowGuide?.desc }}</div>
          </div>
          <div class="workflow-side">
            <span class="status-pill">已完成 {{ completedStepCount }} / {{ workflowSteps.length }}</span>
            <button class="btn mini secondary" @click="go(workflowGuide?.nextRoute)">{{ workflowGuide?.nextLabel }}</button>
          </div>
        </section>

        <section class="workflow-progress card">
          <div class="progress-row">
            <div v-for="(item, index) in workflowSteps" :key="item.title" class="progress-step" :class="`is-${item.state}`">
              <div class="progress-index">{{ index + 1 }}</div>
              <div class="progress-meta">
                <strong>{{ item.title }}</strong>
                <span>{{ item.state === 'done' ? '已完成' : item.state === 'current' ? '当前进行中' : item.state === 'next' ? '建议下一步' : '待处理' }}</span>
              </div>
            </div>
          </div>
          <div v-if="workflowLoading" class="sub">正在同步流程进度...</div>
        </section>

        <RouterView v-slot="{ Component }">
          <Transition name="page-fade-slide" mode="out-in">
            <component :is="Component" />
          </Transition>
        </RouterView>
      </main>
    </div>

    <el-dialog v-model="showOnboarding" width="720px" :show-close="false" class="onboarding-dialog" align-center>
      <template #header>
        <div class="onboarding-head">
          <div>
            <h3>{{ onboardingContent.title }}</h3>
            <div class="sub">{{ onboardingContent.intro }}</div>
          </div>
        </div>
      </template>
      <div class="onboarding-grid">
        <div v-for="item in onboardingContent.steps" :key="item.title" class="onboarding-step">
          <strong>{{ item.title }}</strong>
          <span>{{ item.desc }}</span>
        </div>
      </div>
      <template #footer>
        <div class="row between" style="width:100%;">
          <button class="btn mini ghost" @click="dismissOnboarding">稍后再看</button>
          <button class="btn mini secondary" @click="startOnboardingFlow">开始使用</button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>
