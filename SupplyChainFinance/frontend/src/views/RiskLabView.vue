<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import DOMPurify from 'dompurify'
import { marked } from 'marked'
import { request, requestPreferDB } from '../api/http'

const loading = ref(false)
const aiLoading = ref(false)
const receivables = ref([])
const selectedStatus = ref('')
const query = ref('')
const selectedId = ref('')
const aiMemo = ref('')
const aiMeta = ref(null)

const simulator = reactive({ principal: 100000, annualRateBp: 480, days: 90, reserveRatio: 12 })

const statusLabelMap = { PENDING: '待确认', CONFIRMED: '已确认', FINANCED: '已融资', SETTLED: '已结清' }
const statusToneMap = { PENDING: 'warning', CONFIRMED: 'success', FINANCED: 'primary', SETTLED: 'info' }

const filteredSamples = computed(() => {
  const keyword = query.value.trim().toLowerCase()
  return receivables.value
    .filter((item) => (selectedStatus.value ? item.statusText === selectedStatus.value : true))
    .filter((item) => {
      if (!keyword) return true
      return [item.receivableId, item.orderId, item.issuer, item.payer].filter(Boolean).some((field) => String(field).toLowerCase().includes(keyword))
    })
})

const selectedSample = computed(() => filteredSamples.value.find((item) => item.receivableId === selectedId.value) || filteredSamples.value[0] || null)
const sampleMetrics = computed(() => filteredSamples.value.map(buildRiskMetrics))
const allMetrics = computed(() => receivables.value.map(buildRiskMetrics))

const overview = computed(() => {
  const list = allMetrics.value
  return {
    count: list.length,
    avgRisk: list.length ? Math.round(list.reduce((sum, item) => sum + item.score, 0) / list.length) : 0,
    highRisk: list.filter((item) => item.score >= 70).length,
    totalExposure: list.reduce((sum, item) => sum + item.amount, 0)
  }
})

const buckets = computed(() => {
  const list = sampleMetrics.value
  const groups = [
    createBucket('低风险', '0-34', list.filter((item) => item.score < 35), '#2dd4a6'),
    createBucket('关注区', '35-59', list.filter((item) => item.score >= 35 && item.score < 60), '#37b9ff'),
    createBucket('审慎区', '60-79', list.filter((item) => item.score >= 60 && item.score < 80), '#ffb648'),
    createBucket('高风险', '80-100', list.filter((item) => item.score >= 80), '#ff5f87')
  ]
  const max = Math.max(...groups.map((item) => item.count), 1)
  return groups.map((item) => ({ ...item, width: `${Math.max(18, Math.round((item.count / max) * 100))}%` }))
})

const topRiskList = computed(() => sampleMetrics.value.slice().sort((a, b) => b.score - a.score).slice(0, 5))
const scatterPoints = computed(() => {
  const list = sampleMetrics.value
  if (!list.length) return []
  const maxAmount = Math.max(...list.map((item) => item.amount), 1)
  return list.map((item, index) => ({
    id: item.receivableId,
    x: 12 + Math.min(76, (item.amount / maxAmount) * 76),
    y: 88 - Math.min(76, (item.score / 100) * 76),
    size: 12 + Math.min(16, Math.round(item.integrity / 8)),
    label: item.receivableId,
    color: riskTone(item.score).accent,
    delay: `${index * 0.05}s`,
    selected: selectedSample.value?.receivableId === item.receivableId
  }))
})

const selectedProfile = computed(() => {
  const sample = selectedSample.value
  if (!sample) return null
  const metric = buildRiskMetrics(sample)
  return {
    ...metric,
    scoreTone: riskTone(metric.score),
    statusLabel: statusLabel(metric.statusText),
    overdueDays: Math.max(0, Math.round((Date.now() / 1000 - metric.dueTimestamp) / 86400)),
    radar: radarPolygon(metric),
    chips: [`订单 ${metric.orderId || '-'}`, `付款方 ${shortHash(metric.payer)}`, `发行方 ${shortHash(metric.issuer)}`]
  }
})

const factorRows = computed(() => {
  const sample = selectedProfile.value
  if (!sample) return []
  return [
    { label: '状态压力', value: sample.breakdown.statusPressure, color: '#4f8cff' },
    { label: '金额暴露', value: sample.breakdown.exposurePressure, color: '#32c5ff' },
    { label: '期限压力', value: sample.breakdown.tenorPressure, color: '#7f8dff' },
    { label: '逾期压力', value: sample.breakdown.overduePressure, color: '#ff8f6b' },
    { label: '集中度', value: sample.breakdown.concentrationPressure, color: '#ff5f87' },
    { label: '凭证完整度', value: sample.breakdown.integrityPressure, color: '#56d39b' }
  ]
})

const simulatorSummary = computed(() => {
  const principal = Number(simulator.principal || 0)
  const rate = Number(simulator.annualRateBp || 0) / 10000
  const days = Number(simulator.days || 0)
  const reserve = Number(simulator.reserveRatio || 0) / 100
  const interest = principal * rate * (days / 360)
  const reserveAmount = principal * reserve
  const netAmount = Math.max(0, principal - reserveAmount)
  const leverageIndex = Math.min(100, Math.round((principal / 200000) * 45 + reserve * 100 + rate * 800))
  const selectedRisk = selectedProfile.value?.score || 0
  const combinedRisk = Math.min(100, Math.round(selectedRisk * 0.65 + leverageIndex * 0.35))
  return {
    principal,
    interest,
    reserveAmount,
    netAmount,
    annualizedCost: principal ? ((interest + reserveAmount * 0.02) / principal) * (360 / Math.max(days, 1)) * 100 : 0,
    leverageIndex,
    combinedRisk,
    decision: combinedRisk >= 75 ? '建议缩量或补充担保' : combinedRisk >= 55 ? '附条件准入' : '可进入授信讨论'
  }
})

const reviewMemo = computed(() => {
  const profile = selectedProfile.value
  const scheme = simulatorSummary.value
  if (!profile) return []
  return [
    { label: '审查立场', value: scheme.decision },
    { label: '当前样本', value: profile.receivableId },
    { label: '综合风险', value: `${profile.score} / 100` },
    { label: '净放款', value: formatMoney(scheme.netAmount) },
    { label: '保证金', value: formatMoney(scheme.reserveAmount) },
    { label: '期限', value: `${simulator.days} 天` }
  ]
})

const aiRendered = computed(() => DOMPurify.sanitize(marked.parse(aiMemo.value || '')))
const currentModel = computed(() => aiMeta.value?.model || 'doubao')
const currentUsage = computed(() => {
  const usage = aiMeta.value?.usage
  if (!usage) return ''
  return `Tokens ${usage.total_tokens || 0} · Prompt ${usage.prompt_tokens || 0} · Completion ${usage.completion_tokens || 0}`
})

watch(filteredSamples, (list) => {
  if (!list.length) return (selectedId.value = '')
  if (!list.find((item) => item.receivableId === selectedId.value)) selectedId.value = list[0].receivableId
})

watch(() => selectedSample.value?.amount, (amount) => {
  if (amount) simulator.principal = Math.max(20000, Math.round(Number(amount) * 0.82))
}, { immediate: true })

function createBucket(label, range, items, color) {
  return { label, range, items, color, count: items.length, exposure: items.reduce((sum, item) => sum + item.amount, 0) }
}

function buildRiskMetrics(item) {
  const amount = Number(item.amount || 0)
  const dueTimestamp = Number(item.dueTimestamp || 0)
  const now = Date.now() / 1000
  const overdueDays = Math.max(0, Math.round((now - dueTimestamp) / 86400))
  const daysLeft = dueTimestamp ? Math.round((dueTimestamp - now) / 86400) : 0
  const statusPressure = item.statusText === 'PENDING' ? 26 : item.statusText === 'CONFIRMED' ? 8 : item.statusText === 'FINANCED' ? 18 : 10
  const exposurePressure = Math.min(24, Math.round(amount / 12000))
  const tenorPressure = daysLeft < 0 ? 26 : daysLeft <= 30 ? 18 : daysLeft <= 90 ? 11 : 6
  const overduePressure = Math.min(28, overdueDays * 4)
  const concentrationPressure = item.payer ? Math.min(16, (String(item.payer).length % 8) * 2 + 4) : 8
  const integrity = computeIntegrity(item)
  const integrityPressure = Math.max(4, 18 - Math.round(integrity / 8))
  return {
    ...item,
    amount,
    dueTimestamp,
    integrity,
    score: Math.min(100, statusPressure + exposurePressure + tenorPressure + overduePressure + concentrationPressure + integrityPressure),
    breakdown: { statusPressure, exposurePressure, tenorPressure, overduePressure, concentrationPressure, integrityPressure }
  }
}

function computeIntegrity(item) {
  let score = 50
  if (item.receivableId) score += 10
  if (item.orderId) score += 10
  if (item.issuer) score += 10
  if (item.payer) score += 10
  if (Number(item.amount || 0) > 0) score += 10
  return Math.min(100, score)
}

function radarPolygon(metric) {
  const values = [100 - metric.breakdown.statusPressure * 3, 100 - metric.breakdown.exposurePressure * 4, 100 - metric.breakdown.tenorPressure * 4, 100 - metric.breakdown.overduePressure * 3, 100 - metric.breakdown.concentrationPressure * 5, metric.integrity].map((value) => Math.max(12, Math.min(100, value)))
  const center = 90
  const radius = 68
  return values.map((value, index) => {
    const angle = (-90 + index * 60) * (Math.PI / 180)
    const currentRadius = (value / 100) * radius
    return `${center + Math.cos(angle) * currentRadius},${center + Math.sin(angle) * currentRadius}`
  }).join(' ')
}

function riskTone(score) {
  if (score >= 80) return { label: '高风险', accent: '#ff5f87', glow: 'rgba(255,95,135,0.28)' }
  if (score >= 60) return { label: '审慎区', accent: '#ff9f43', glow: 'rgba(255,159,67,0.24)' }
  if (score >= 35) return { label: '关注区', accent: '#37b9ff', glow: 'rgba(55,185,255,0.22)' }
  return { label: '低风险', accent: '#2dd4a6', glow: 'rgba(45,212,166,0.2)' }
}

function riskText(score) {
  if (score >= 80) return '当前样本风险暴露较高，建议暂停准入并补充担保或凭证。'
  if (score >= 60) return '建议附条件准入，重点复核付款能力、期限和凭证完整度。'
  if (score >= 35) return '可继续观察，优先补足凭证并控制放款规模。'
  return '样本表现稳定，可进入授信讨论与额度测算。'
}

function statusLabel(status) { return statusLabelMap[status] || status || '未识别' }
function shortHash(value) { if (!value) return '-'; const text = String(value); return text.length <= 16 ? text : `${text.slice(0, 8)}...${text.slice(-6)}` }
function formatMoney(value) { return `¥${Number(value || 0).toLocaleString('zh-CN', { maximumFractionDigits: 2 })}` }
function percent(value) { return `${Math.round(Number(value || 0) * 100)}%` }
function focusSample(id) { selectedId.value = id }

async function loadSamples() {
  loading.value = true
  try {
    const resp = await requestPreferDB('/api/db/receivables', '/api/receivable', { withRole: true })
    receivables.value = Array.isArray(resp?.data) ? resp.data : []
  } catch (error) {
    receivables.value = []
    ElMessage.error(error?.message || '风控样本加载失败')
  } finally {
    loading.value = false
  }
}

async function generateAiMemo() {
  if (!selectedProfile.value) return ElMessage.info('先选择一个应收样本再生成 AI 审查意见')
  aiLoading.value = true
  try {
    const profile = selectedProfile.value
    const scheme = simulatorSummary.value
    const prompt = [
      '请以供应链金融风控审批官身份输出简明中文审查意见。',
      '输出结构：一、结论；二、核心风险；三、建议措施；四、是否建议人工复核。',
      `样本数：${overview.value.count}`,
      `平均风险分：${overview.value.avgRisk}`,
      `高风险样本：${overview.value.highRisk}`,
      `当前样本：${profile.receivableId}`,
      `当前状态：${profile.statusLabel}`,
      `应收金额：${profile.amount}`,
      `风险分：${profile.score}`,
      `凭证完整度：${profile.integrity}`,
      `逾期天数：${profile.overdueDays}`,
      `融资本金：${scheme.principal}`,
      `预计利息：${scheme.interest.toFixed(2)}`,
      `保证金：${scheme.reserveAmount.toFixed(2)}`,
      `放款净额：${scheme.netAmount.toFixed(2)}`,
      `综合风险：${scheme.combinedRisk}`,
      `系统建议：${scheme.decision}`
    ].join('\n')
    const resp = await request('/api/ai/chat', { method: 'POST', withRole: true, body: { messages: [{ role: 'user', content: prompt }], scene: 'finance' } })
    aiMemo.value = resp?.reply || '未获取到审查意见。'
    aiMeta.value = resp?.meta || null
  } catch {
    aiMemo.value = ['### 审查意见', `- 当前样本：${selectedProfile.value.receivableId}`, `- 综合风险：${selectedProfile.value.score}/100`, `- 结论：${simulatorSummary.value.decision}`, '', '### 建议措施', '- 优先核验付款方履约能力与凭证链完整度。', '- 根据风险分层控制放款净额和保证金比例。', '- 对逾期或高暴露样本建议人工复核。'].join('\n')
    aiMeta.value = { model: 'local-fallback', usage: null }
  } finally {
    aiLoading.value = false
  }
}

onMounted(loadSamples)
</script>

<template>
  <section class="risk-lab-page">
    <div class="risk-hero card panel">
      <div>
        <div class="hero-badge">风控模型</div>
        <h1>金融风控画像实验室</h1>
        <div class="hero-metrics">
          <div class="hero-metric"><span>样本总量</span><strong>{{ overview.count }}</strong></div>
          <div class="hero-metric"><span>平均风险分</span><strong>{{ overview.avgRisk }}</strong></div>
          <div class="hero-metric"><span>高风险样本</span><strong>{{ overview.highRisk }}</strong></div>
          <div class="hero-metric"><span>风险暴露</span><strong>{{ formatMoney(overview.totalExposure) }}</strong></div>
        </div>
      </div>
      <div class="hero-side">
        <div class="hero-orb" :style="{ '--orb': selectedProfile?.scoreTone?.accent || '#4f8cff' }"><span>{{ selectedProfile?.score || 0 }}</span><small>{{ selectedProfile?.scoreTone?.label || '待分析' }}</small></div>
        <el-button type="primary" :loading="aiLoading" @click="generateAiMemo">生成 AI 审查意见</el-button>
      </div>
    </div>

    <div v-if="loading" class="risk-grid"><div class="card panel"><el-skeleton animated :rows="10" /></div><div class="card panel"><el-skeleton animated :rows="10" /></div><div class="card panel"><el-skeleton animated :rows="8" /></div><div class="card panel"><el-skeleton animated :rows="8" /></div></div>

    <div v-else-if="!receivables.length" class="card panel risk-empty-panel">
      <el-empty description="当前没有可分析的风控样本">
        <template #image><div class="empty-lab"><div class="empty-grid"></div><div class="empty-pulse"></div><div class="empty-core">Risk</div></div></template>
        <p class="sub">先让企业端生成并同步应收数据，风控实验室才有样本可做画像、分层和 AI 解读。</p>
      </el-empty>
    </div>

    <template v-else>
      <div class="risk-grid">
        <section class="card panel risk-section">
          <div class="section-head">
            <h2>风险分层视图</h2>
            <div class="filters-row">
              <el-select v-model="selectedStatus" clearable placeholder="全部状态" style="width:140px"><el-option label="待确认" value="PENDING" /><el-option label="已确认" value="CONFIRMED" /><el-option label="已融资" value="FINANCED" /><el-option label="已结清" value="SETTLED" /></el-select>
              <el-input v-model="query" clearable placeholder="搜索应收ID / 订单 / 地址" style="width:240px" />
              <el-button plain @click="loadSamples">刷新</el-button>
            </div>
          </div>
          <div class="bucket-grid">
            <button v-for="bucket in buckets" :key="bucket.label" class="bucket-card" :style="{ '--bucket': bucket.color }" type="button">
              <div class="bucket-top"><strong>{{ bucket.label }}</strong><span>{{ bucket.range }}</span></div>
              <div class="bucket-value">{{ bucket.count }}</div>
              <div class="bucket-sub">{{ formatMoney(bucket.exposure) }}</div>
              <div class="bucket-bar"><span :style="{ width: bucket.width }"></span></div>
            </button>
          </div>
          <div class="matrix-panel">
            <div class="bucket-top"><strong>风险热区矩阵</strong><span>金额暴露 / 风险分</span></div>
            <div class="matrix-board">
              <div class="matrix-grid"></div>
              <button v-for="point in scatterPoints" :key="point.id" class="matrix-point" :class="{ selected: point.selected }" :style="{ left: `${point.x}%`, top: `${point.y}%`, width: `${point.size}px`, height: `${point.size}px`, background: point.color, animationDelay: point.delay }" type="button" @click="focusSample(point.id)"><span>{{ point.label }}</span></button>
              <div class="axis axis-y">风险分</div><div class="axis axis-x">金额暴露</div>
            </div>
          </div>
          <div class="top-risk-list">
            <div class="bucket-top"><strong>高风险样本雷达</strong><span>Top 5</span></div>
            <button v-for="item in topRiskList" :key="item.receivableId" class="risk-rank-item" type="button" @click="focusSample(item.receivableId)"><div><strong>{{ item.receivableId }}</strong><div class="sub">{{ statusLabel(item.statusText) }} · {{ formatMoney(item.amount) }}</div></div><div class="rank-score" :style="{ color: riskTone(item.score).accent }">{{ item.score }}</div></button>
          </div>
        </section>

        <section class="card panel risk-section">
          <div class="section-head"><h2>风险画像剖面</h2><el-select v-model="selectedId" filterable placeholder="选择样本" style="width:220px"><el-option v-for="item in filteredSamples" :key="item.receivableId" :label="`${item.receivableId} · ${statusLabel(item.statusText)}`" :value="item.receivableId" /></el-select></div>
          <div v-if="selectedProfile" class="profile-shell">
            <div class="profile-score" :style="{ '--tone': selectedProfile.scoreTone.accent, '--glow': selectedProfile.scoreTone.glow }">
              <div class="score-ring"><svg viewBox="0 0 180 180"><circle cx="90" cy="90" r="68" class="ring-bg" /><circle cx="90" cy="90" r="68" class="ring-value" :stroke-dasharray="`${selectedProfile.score * 4.27} 427`" /></svg><div class="score-center"><strong>{{ selectedProfile.score }}</strong><span>{{ selectedProfile.scoreTone.label }}</span></div></div>
              <p>{{ riskText(selectedProfile.score) }}</p>
              <div class="chip-row"><span v-for="chip in selectedProfile.chips" :key="chip" class="chip">{{ chip }}</span></div>
            </div>
            <div class="profile-radar">
              <svg viewBox="0 0 180 180" class="radar-svg"><polygon points="90,18 150,54 150,126 90,162 30,126 30,54" class="radar-grid-shape" /><polygon points="90,36 135,63 135,117 90,144 45,117 45,63" class="radar-grid-shape inner" /><polygon :points="selectedProfile.radar" class="radar-value-shape" /><text x="90" y="12">状态</text><text x="152" y="58">暴露</text><text x="150" y="132">期限</text><text x="90" y="176">逾期</text><text x="18" y="132">集中</text><text x="18" y="58">完整</text></svg>
              <div class="stats-grid"><div><span>当前状态</span><strong>{{ selectedProfile.statusLabel }}</strong></div><div><span>凭证完整度</span><strong>{{ selectedProfile.integrity }}/100</strong></div><div><span>应收金额</span><strong>{{ formatMoney(selectedProfile.amount) }}</strong></div><div><span>逾期天数</span><strong>{{ selectedProfile.overdueDays }}</strong></div></div>
            </div>
            <div class="factor-stack"><div v-for="factor in factorRows" :key="factor.label" class="factor-row"><div class="factor-labels"><span>{{ factor.label }}</span><strong>{{ factor.value }}</strong></div><div class="factor-track"><span :style="{ width: `${Math.min(100, factor.value)}%`, background: factor.color }"></span></div></div></div>
          </div>
        </section>
        <section class="card panel risk-section">
          <div class="section-head"><h2>融资方案模拟器</h2><div class="scheme-chip">{{ simulatorSummary.decision }}</div></div>
          <div class="simulator-grid">
            <el-form label-position="top" class="sim-form">
              <el-form-item label="融资本金"><el-input-number v-model="simulator.principal" :min="10000" :step="10000" style="width:100%" /></el-form-item>
              <el-form-item label="年化利率(BP)"><el-input-number v-model="simulator.annualRateBp" :min="0" :step="10" style="width:100%" /></el-form-item>
              <el-form-item label="融资期限(天)"><el-input-number v-model="simulator.days" :min="1" :step="5" style="width:100%" /></el-form-item>
              <el-form-item label="保证金比例(%)"><el-input-number v-model="simulator.reserveRatio" :min="0" :max="50" :step="1" style="width:100%" /></el-form-item>
            </el-form>
            <div class="scheme-board">
              <div class="kpi-grid"><div class="info-card"><span>预计利息</span><strong>{{ formatMoney(simulatorSummary.interest) }}</strong></div><div class="info-card"><span>保证金</span><strong>{{ formatMoney(simulatorSummary.reserveAmount) }}</strong></div><div class="info-card"><span>净放款</span><strong>{{ formatMoney(simulatorSummary.netAmount) }}</strong></div><div class="info-card"><span>综合风险</span><strong>{{ simulatorSummary.combinedRisk }}</strong></div></div>
              <div class="capital-stack"><div class="capital-bar"><span class="principal-bar" :style="{ width: percent(simulatorSummary.netAmount / Math.max(simulatorSummary.principal, 1)) }">净放款</span><span class="reserve-bar" :style="{ width: percent(simulatorSummary.reserveAmount / Math.max(simulatorSummary.principal, 1)) }">保证金</span></div><div class="capital-meta"><span>综合资金成本 {{ simulatorSummary.annualizedCost.toFixed(2) }}%</span><span>杠杆指数 {{ simulatorSummary.leverageIndex }}</span></div></div>
              <div class="kpi-grid"><div v-for="item in reviewMemo" :key="item.label" class="info-card"><span>{{ item.label }}</span><strong>{{ item.value }}</strong></div></div>
            </div>
          </div>
        </section>

        <section class="card panel risk-section">
          <div class="section-head"><h2>AI 审查纪要</h2><div class="filters-row"><span class="scheme-chip">{{ currentModel }}</span><span v-if="currentUsage" class="scheme-chip">{{ currentUsage }}</span></div></div>
          <div class="kpi-grid"><div class="info-card"><span>重点样本</span><strong>{{ selectedProfile?.receivableId || '-' }}</strong></div><div class="info-card"><span>样本状态</span><strong>{{ selectedProfile?.statusLabel || '-' }}</strong></div><div class="info-card"><span>AI 立场</span><strong>{{ simulatorSummary.decision }}</strong></div><div class="info-card"><span>人工复核</span><strong>{{ (selectedProfile?.score || 0) >= 70 ? '建议' : '按需' }}</strong></div></div>
          <div class="ai-panel" v-loading="aiLoading">
            <div v-if="aiMemo" class="markdown-body" v-html="aiRendered"></div>
            <div v-else class="ai-placeholder"><div class="placeholder-orb"></div><strong>等待生成审查意见</strong><span>系统会结合样本风险、期限压力、暴露规模和融资方案给出更像审批纪要的结果。</span></div>
          </div>
        </section>
      </div>
    </template>
  </section>
</template>

<style scoped>
.risk-lab-page,.risk-grid,.risk-section,.profile-shell{display:grid;gap:18px}.risk-grid{grid-template-columns:repeat(2,minmax(0,1fr))}.risk-hero{position:relative;overflow:hidden;display:grid;grid-template-columns:minmax(0,1.6fr) 260px;gap:20px;padding:26px 28px;background:radial-gradient(circle at top left,rgba(105,180,255,.24),transparent 42%),linear-gradient(135deg,rgba(255,255,255,.96),rgba(240,247,255,.94))}.risk-hero::after{content:'';position:absolute;inset:auto -80px -90px auto;width:260px;height:260px;border-radius:50%;background:radial-gradient(circle,rgba(83,150,255,.22),transparent 68%)}.hero-badge,.scheme-chip{display:inline-flex;align-items:center;padding:8px 12px;border-radius:999px;font-size:12px;color:#3565d9;background:rgba(79,140,255,.1);border:1px solid rgba(79,140,255,.18)}.hero-metrics,.bucket-grid,.kpi-grid,.chip-row,.stats-grid{display:grid;gap:12px}.hero-metrics{grid-template-columns:repeat(4,minmax(0,1fr));margin-top:18px}.hero-metric,.bucket-card,.info-card,.sim-form,.scheme-board,.profile-score,.profile-radar,.factor-stack,.matrix-panel,.ai-panel{border-radius:22px;background:linear-gradient(180deg,rgba(255,255,255,.96),rgba(245,249,255,.94));border:1px solid rgba(98,147,255,.14)}.hero-metric,.bucket-card,.info-card,.matrix-panel,.scheme-board,.profile-score,.profile-radar,.factor-stack,.ai-panel{padding:18px}.hero-metric span,.bucket-top span,.info-card span,.factor-labels span,.capital-meta span,.sub{font-size:13px;color:#61769b}.hero-metric strong,.bucket-value,.rank-score{display:block;margin-top:8px;font-size:28px;color:#0f1f38}.hero-side{position:relative;z-index:1;display:grid;align-content:center;justify-items:end;gap:18px}.hero-orb{--orb:#4f8cff;width:168px;height:168px;border-radius:50%;display:grid;place-items:center;text-align:center;color:#fff;background:radial-gradient(circle at 30% 30%,rgba(255,255,255,.36),transparent 36%),radial-gradient(circle at 50% 50%,color-mix(in srgb,var(--orb) 72%,#fff 28%),var(--orb));box-shadow:0 24px 54px color-mix(in srgb,var(--orb) 32%,transparent);animation:orbFloat 5.6s ease-in-out infinite}.hero-orb span{display:block;font-size:42px;font-weight:800}.bucket-grid,.kpi-grid{grid-template-columns:repeat(2,minmax(0,1fr))}.bucket-grid{grid-template-columns:repeat(4,minmax(0,1fr))}.section-head,.filters-row,.bucket-top,.capital-meta,.factor-labels{display:flex;align-items:center;justify-content:space-between;gap:12px;flex-wrap:wrap}.bucket-card{position:relative;text-align:left;border:0;background:linear-gradient(160deg,rgba(255,255,255,.92),rgba(246,250,255,.98));box-shadow:inset 0 0 0 1px rgba(99,151,255,.14)}.bucket-value{font-size:34px}.bucket-bar,.factor-track{margin-top:12px;height:8px;border-radius:999px;background:rgba(82,124,196,.08);overflow:hidden}.bucket-bar span,.factor-track span{display:block;height:100%;border-radius:inherit}.matrix-board{position:relative;margin-top:14px;height:280px;border-radius:22px;background:linear-gradient(180deg,rgba(255,255,255,.84),rgba(244,248,255,.98));overflow:hidden}.matrix-grid{position:absolute;inset:0;background-image:linear-gradient(rgba(100,146,255,.08) 1px,transparent 1px),linear-gradient(90deg,rgba(100,146,255,.08) 1px,transparent 1px);background-size:25% 25%}.matrix-point{position:absolute;border:0;border-radius:50%;cursor:pointer;transform:translate(-50%,-50%);box-shadow:0 0 0 8px rgba(255,255,255,.35);animation:pulseIn 1.8s ease-out infinite}.matrix-point.selected{box-shadow:0 0 0 10px rgba(255,255,255,.48),0 0 22px rgba(79,140,255,.28)}.matrix-point span{position:absolute;top:calc(100% + 8px);left:50%;transform:translateX(-50%);padding:4px 8px;font-size:12px;white-space:nowrap;color:#2f4f85;background:rgba(255,255,255,.9);border-radius:999px}.axis{position:absolute;font-size:12px;color:#6f84ab}.axis-y{left:12px;top:12px}.axis-x{right:16px;bottom:12px}.risk-rank-item{display:flex;align-items:center;justify-content:space-between;gap:12px;width:100%;border:0;cursor:pointer;text-align:left;padding:14px 16px;border-radius:18px;background:rgba(248,251,255,.9);box-shadow:inset 0 0 0 1px rgba(89,139,245,.1);transition:transform .24s ease,box-shadow .24s ease}.risk-rank-item:hover{transform:translateY(-2px);box-shadow:inset 0 0 0 1px rgba(89,139,245,.2),0 14px 28px rgba(75,118,188,.12)}.profile-score{box-shadow:0 20px 36px var(--glow,rgba(79,140,255,.12))}.score-ring{position:relative;width:180px;height:180px;margin:0 auto 12px}.score-ring svg{width:100%;height:100%;transform:rotate(-90deg)}.ring-bg{fill:none;stroke:rgba(89,127,194,.12);stroke-width:14}.ring-value{fill:none;stroke:var(--tone,#4f8cff);stroke-width:14;stroke-linecap:round}.score-center{position:absolute;inset:0;display:grid;place-items:center;text-align:center}.score-center strong{font-size:38px}.chip-row{grid-template-columns:repeat(3,minmax(0,1fr));margin-top:16px}.chip{padding:10px 12px;text-align:center;border-radius:14px;background:rgba(241,247,255,.9);color:#345282}.profile-radar{display:grid;grid-template-columns:220px minmax(0,1fr);align-items:center;gap:18px}.radar-svg{width:100%;height:auto;overflow:visible}.radar-svg text{font-size:11px;fill:#5f7396}.radar-grid-shape{fill:rgba(88,148,255,.06);stroke:rgba(88,148,255,.18)}.radar-grid-shape.inner{fill:rgba(88,148,255,.03)}.radar-value-shape{fill:rgba(79,140,255,.24);stroke:#4f8cff;stroke-width:2}.stats-grid{grid-template-columns:repeat(2,minmax(0,1fr))}.simulator-grid{display:grid;grid-template-columns:280px minmax(0,1fr);gap:16px}.capital-stack{padding:14px;border-radius:18px;background:rgba(244,248,255,.94)}.capital-bar{display:flex;overflow:hidden;border-radius:999px;height:42px;background:rgba(90,132,194,.08)}.capital-bar span{display:flex;align-items:center;justify-content:center;color:#fff;font-size:13px;white-space:nowrap}.principal-bar{background:linear-gradient(90deg,#3975ff,#33b8ff)}.reserve-bar{background:linear-gradient(90deg,#ff7f7f,#ffb648)}.ai-panel{min-height:360px}.ai-placeholder{min-height:320px;display:grid;place-items:center;align-content:center;gap:12px;text-align:center;color:#607699}.placeholder-orb,.empty-pulse{width:88px;height:88px;border-radius:50%;background:radial-gradient(circle,rgba(80,143,255,.36),rgba(80,143,255,.1));box-shadow:0 0 0 18px rgba(80,143,255,.08)}.markdown-body :deep(h1),.markdown-body :deep(h2),.markdown-body :deep(h3){margin:0 0 12px;color:#10203a}.markdown-body :deep(p),.markdown-body :deep(li){color:#4c6387;line-height:1.8}.markdown-body :deep(ul),.markdown-body :deep(ol){margin:0;padding-left:18px}.risk-empty-panel{min-height:560px;display:grid;place-items:center}.empty-lab{position:relative;width:180px;height:180px;display:grid;place-items:center}.empty-grid{position:absolute;inset:18px;border-radius:30px;background-image:linear-gradient(rgba(90,142,255,.12) 1px,transparent 1px),linear-gradient(90deg,rgba(90,142,255,.12) 1px,transparent 1px);background-size:25% 25%}.empty-pulse{animation:pulseIn 2s ease-in-out infinite}.empty-core{position:absolute;font-size:28px;font-weight:800;color:#3565d9}@keyframes pulseIn{0%,100%{transform:translate(-50%,-50%) scale(1)}50%{transform:translate(-50%,-50%) scale(1.14)}}@keyframes orbFloat{0%,100%{transform:translateY(0)}50%{transform:translateY(-10px)}}@media (max-width:1280px){.risk-grid,.risk-hero,.profile-radar,.simulator-grid{grid-template-columns:1fr}.hero-side{justify-items:start}}@media (max-width:900px){.hero-metrics,.bucket-grid,.chip-row,.stats-grid,.kpi-grid{grid-template-columns:1fr 1fr}}@media (max-width:720px){.hero-metrics,.bucket-grid,.chip-row,.stats-grid,.kpi-grid{grid-template-columns:1fr}.section-head,.filters-row,.capital-meta{flex-direction:column;align-items:stretch}}
</style>
