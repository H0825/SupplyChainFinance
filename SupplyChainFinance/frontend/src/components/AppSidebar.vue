<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  Aim,
  ArrowLeftBold,
  ArrowRightBold,
  ChatDotRound,
  Collection,
  DataAnalysis,
  Document,
  Files,
  House,
  Moon,
  OfficeBuilding,
  Opportunity,
  Reading,
  Sunny
} from '@element-plus/icons-vue'
import { useAuthStore } from '../stores/auth'

const props = defineProps({
  collapsed: {
    type: Boolean,
    default: false
  },
  themeMode: {
    type: String,
    default: 'light'
  }
})

const emit = defineEmits(['toggle-sidebar', 'toggle-theme'])

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()

const roleGuide = {
  business: {
    title: '企业端引导',
    tips: ['先完成企业登记，再创建应收与上传凭证', '应收确认后再进入融资协同', 'AI 助手适合整理融资说明与凭证检查清单']
  },
  finance: {
    title: '金融端引导',
    tips: ['先看金融总览，再进入融资处理', '风控模型页适合做贷前筛查与授信判断', 'AI 助手适合生成审查清单与隐私合规建议']
  },
  admin: {
    title: '管理端引导',
    tips: ['先看平台总览，再查看链上运行状态', '管理控制台用于用户与企业映射核查', '审计中心适合追踪异常交易与权限使用']
  }
}

const menus = computed(() => {
  const role = auth.role

  if (role === 'business') {
    return [
      { to: '/dashboard', label: '运营总览', icon: DataAnalysis },
      { to: '/business', label: '企业业务', icon: OfficeBuilding },
      { to: '/receivables', label: '应收资产', icon: Collection },
      { to: '/documents', label: '凭证文档', icon: Document },
      { to: '/ai-assistant', label: 'AI助手', icon: ChatDotRound }
    ]
  }

  if (role === 'finance') {
    return [
      { to: '/dashboard', label: '金融总览', icon: DataAnalysis },
      { to: '/finance', label: '融资处理', icon: Aim },
      { to: '/risk-lab', label: '风控模型', icon: Opportunity },
      { to: '/receivables', label: '应收池', icon: Collection },
      { to: '/ai-assistant', label: 'AI助手', icon: ChatDotRound }
    ]
  }

  return [
    { to: '/dashboard', label: '平台总览', icon: House },
    { to: '/audit', label: '审计中心', icon: Reading },
    { to: '/control-center', label: '管理控制台', icon: DataAnalysis },
    { to: '/documents', label: '文档中心', icon: Files },
    { to: '/ai-assistant', label: 'AI助手', icon: ChatDotRound }
  ]
})

const guide = computed(() => roleGuide[auth.role] || roleGuide.admin)

function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<template>
  <aside class="sidebar" :class="{ collapsed }">
    <div class="sidebar-top row between">
      <div v-if="!collapsed">
        <div class="brand">SCF CONSOLE</div>
        <div class="role-chip">ROLE · {{ auth.role.toUpperCase() }}</div>
      </div>
      <el-tooltip :content="collapsed ? '展开侧边栏' : '折叠侧边栏'" placement="right">
        <button class="sidebar-toggle" @click="emit('toggle-sidebar')">
          <el-icon><ArrowRightBold v-if="collapsed" /><ArrowLeftBold v-else /></el-icon>
        </button>
      </el-tooltip>
    </div>

    <div class="sidebar-actions">
      <el-tooltip :content="themeMode === 'dark' ? '切换为浅色模式' : '切换为暗色模式'" placement="right">
        <button class="side-tool icon-only" @click="emit('toggle-theme')">
          <el-icon><Sunny v-if="themeMode === 'dark'" /><Moon v-else /></el-icon>
        </button>
      </el-tooltip>
    </div>

    <nav class="sidebar-nav">
      <el-tooltip v-for="menu in menus" :key="menu.to" :content="menu.label" placement="right" :disabled="!collapsed">
        <RouterLink :to="menu.to" class="nav-btn" :class="{ active: route.path === menu.to }">
          <span class="nav-icon"><el-icon><component :is="menu.icon" /></el-icon></span>
          <span v-if="!collapsed">{{ menu.label }}</span>
        </RouterLink>
      </el-tooltip>
    </nav>

    <div v-if="!collapsed" class="guide-card">
      <div class="guide-title">{{ guide.title }}</div>
      <div v-for="tip in guide.tips" :key="tip" class="guide-tip">{{ tip }}</div>
    </div>

    <div class="sidebar-foot">
      <div v-if="!collapsed" class="sub sidebar-user">{{ auth.user?.username || '未登录' }}</div>
      <el-tooltip content="退出当前账号" placement="right" :disabled="!collapsed">
        <button class="btn mini danger" @click="logout">退出登录</button>
      </el-tooltip>
    </div>
  </aside>
</template>
