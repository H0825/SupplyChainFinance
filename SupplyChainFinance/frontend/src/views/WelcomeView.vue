<script setup>
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const welcomeRef = ref(null)
const pageReady = ref(false)
let readyTimer = null

function enter() {
  router.push('/login')
}

function onMouseMove(e) {
  const el = welcomeRef.value
  if (!el) return
  const rect = el.getBoundingClientRect()
  const x = e.clientX - rect.left
  const y = e.clientY - rect.top
  el.style.setProperty('--mx', `${x}px`)
  el.style.setProperty('--my', `${y}px`)
}

onMounted(() => {
  const el = welcomeRef.value
  if (el) el.addEventListener('mousemove', onMouseMove)
  readyTimer = window.setTimeout(() => {
    pageReady.value = true
  }, 260)
})

onBeforeUnmount(() => {
  const el = welcomeRef.value
  if (el) el.removeEventListener('mousemove', onMouseMove)
  if (readyTimer) window.clearTimeout(readyTimer)
})
</script>

<template>
  <section ref="welcomeRef" class="welcome-page">
    <div class="glow orb-a"></div>
    <div class="glow orb-b"></div>
    <div class="mouse-light"></div>

    <div v-if="!pageReady" class="hero-card hero-tech welcome-skeleton">
      <el-skeleton animated :rows="6" />
      <div class="welcome-skeleton-grid">
        <el-skeleton v-for="i in 4" :key="i" animated :rows="3" />
      </div>
    </div>

    <div v-else class="hero-card hero-tech">
      <div class="tag">NEXT-GEN SCF PLATFORM</div>
      <h1>供应链金融智能协同平台</h1>
      <p>
        区块链可信流转 + 多角色控制台，支撑企业、金融机构与管理方的全流程协同。
      </p>

      <div class="hero-grid">
        <div class="hero-block">
          <h3>可信链路</h3>
          <span>应收创建、确权、融资、结算全流程留痕</span>
        </div>
        <div class="hero-block">
          <h3>协同控制台</h3>
          <span>按角色进入专属工作台，流程清晰可控</span>
        </div>
        <div class="hero-block">
          <h3>角色控制台</h3>
          <span>business / finance / admin 分角色导航</span>
        </div>
        <div class="hero-block">
          <h3>审计可追溯</h3>
          <span>交易日志、文档索引、风控实验室</span>
        </div>
      </div>

      <div class="row wrap" style="margin-top: 14px;">
        <button class="btn neon" @click="enter">进入系统</button>
        <button class="btn ghost" @click="router.push('/login')">登录 / 注册</button>
      </div>
      <div class="sub" style="margin-top: 8px;">基于FISCO BCOS区块链开发</div>
    </div>

    <div v-if="pageReady" class="feature-strip">
      <div class="f-item">
        <strong>SCF INSIGHT</strong>
        <span>供应链金融链上流转与企业信用协同平台</span>
      </div>
      <div class="f-item">
        <strong>企业端</strong>
        <span>企业注册、应收登记、确权、结算</span>
      </div>
      <div class="f-item">
        <strong>金融端</strong>
        <span>融资登记、批量处理、风控建模</span>
      </div>
      <div class="f-item">
        <strong>管理端</strong>
        <span>审计中心、控制台、文档治理</span>
      </div>
    </div>

    <footer v-if="pageReady" class="welcome-footer">
      <div>SCF Insight Platform · Version 1.0.0</div>
      <div>Copyright © 2026 SCF Insight Team. All rights reserved.</div>
      <div>Powered by FISCO BCOS</div>
    </footer>
  </section>
</template>

<style scoped>
.welcome-skeleton {
  display: grid;
  gap: 18px;
}

.welcome-skeleton-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

@media (max-width: 760px) {
  .welcome-skeleton-grid {
    grid-template-columns: 1fr;
  }
}
</style>
