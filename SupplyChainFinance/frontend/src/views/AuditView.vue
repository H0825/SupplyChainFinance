<script setup>
import { computed, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import DOMPurify from 'dompurify'
import { marked } from 'marked'
import { request } from '../api/http'

const logs = ref([])
const stats = ref(null)
const actionFilter = ref('')
const keyword = ref('')
const statusFilter = ref('all')
const traceId = ref('')
const loading = ref(false)
const aiLoading = ref(false)
const err = ref('')
const aiSummary = ref('')

marked.setOptions({
  gfm: true,
  breaks: true
})

const statusMap = {
  '0': '执行成功',
  '1': '业务异常',
  '22': '交易回滚',
  '-1': '系统异常'
}

const filtered = computed(() => {
  return logs.value.filter((item) => {
    const matchAction = actionFilter.value ? item.action === actionFilter.value : true
    const matchStatus = statusFilter.value === 'all' ? true : String(item.status) === statusFilter.value
    const key = keyword.value.trim().toLowerCase()
    const matchKeyword = !key
      ? true
      : [item.businessId, item.txHash, item.operatorAddr, item.message].some((value) =>
          String(value || '').toLowerCase().includes(key)
        )
    return matchAction && matchStatus && matchKeyword
  })
})

const actions = computed(() => Array.from(new Set(logs.value.map((item) => item.action))).filter(Boolean).sort())

const traceLogs = computed(() => {
  const id = traceId.value.trim().toLowerCase()
  if (!id) return []
  return logs.value.filter((item) => String(item.businessId || '').toLowerCase().includes(id))
})

const summaryCards = computed(() => {
  const total = logs.value.length
  const abnormal = logs.value.filter((item) => Number(item.status) !== 0).length
  const rollback = logs.value.filter((item) => Number(item.status) === 22).length
  const uniqueOps = new Set(logs.value.map((item) => item.operatorAddr).filter(Boolean)).size
  return [
    { label: '审计日志', value: total },
    { label: '异常记录', value: abnormal },
    { label: '回滚交易', value: rollback },
    { label: '操作账户', value: uniqueOps }
  ]
})

const anomalyCards = computed(() => {
  const byAction = {}
  logs.value.forEach((item) => {
    const key = item.action || '未命名动作'
    if (!byAction[key]) byAction[key] = { total: 0, abnormal: 0 }
    byAction[key].total += 1
    if (Number(item.status) !== 0) byAction[key].abnormal += 1
  })
  return Object.entries(byAction)
    .map(([action, row]) => ({ action, ...row }))
    .sort((a, b) => b.abnormal - a.abnormal || b.total - a.total)
    .slice(0, 6)
})

const recentWarnings = computed(() => filtered.value.filter((item) => Number(item.status) !== 0).slice(0, 8))
const isEmptyAudit = computed(() => !loading.value && logs.value.length === 0)

function statusLabel(status) {
  return statusMap[String(status)] || `状态 ${status}`
}

function renderMarkdown(content) {
  return DOMPurify.sanitize(marked.parse(String(content || '')))
}

function triggerDownload(content, fileName, type) {
  const blob = new Blob([content], { type })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = fileName
  a.click()
  URL.revokeObjectURL(url)
}

function escapeCsv(value) {
  const text = String(value ?? '')
  return `"${text.replace(/"/g, '""')}"`
}

async function load() {
  loading.value = true
  try {
    const [logResp, statResp] = await Promise.all([request('/api/tx-logs?limit=200'), request('/workflow/public/stats')])
    logs.value = logResp.data || []
    stats.value = statResp
    err.value = ''
  } catch (e) {
    err.value = JSON.stringify(e, null, 2)
  } finally {
    loading.value = false
  }
}

async function generateAISummary() {
  if (!filtered.value.length) return
  aiLoading.value = true
  try {
    const excerpt = filtered.value
      .slice(0, 12)
      .map((item) => `${item.createdAt} | ${item.action} | ${item.businessId} | ${statusLabel(item.status)} | ${item.operatorAddr}`)
      .join('\n')
    const prompt = [
      '你是供应链金融审计分析助手。',
      `总日志数: ${logs.value.length}`,
      `当前筛选结果: ${filtered.value.length}`,
      `异常记录: ${summaryCards.value[1].value}`,
      '以下是审计摘录：',
      excerpt,
      '请输出三部分：1. 审计判断 2. 异常模式 3. 建议人工复核项。使用 markdown。'
    ].join('\n')
    const resp = await request('/api/ai/chat', {
      method: 'POST',
      withRole: true,
      body: {
        message: prompt,
        scene: 'admin',
        guard: {
          privacyMode: true,
          localMasking: true,
          noRetention: true,
          minimalDisclosure: true
        }
      }
    })
    aiSummary.value = resp?.reply || ''
  } catch (_) {
    aiSummary.value = [
      '## 审计判断',
      `当前筛选结果共 **${filtered.value.length}** 条，异常记录 **${summaryCards.value[1].value}** 条。`,
      '## 异常模式',
      '- 重复失败动作需要复核操作账户与参数来源。',
      '- 回滚交易应核查上游业务状态是否与链上状态一致。',
      '## 建议人工复核项',
      '- 优先检查高频异常动作和同一业务编号的连续失败记录。'
    ].join('\n')
  } finally {
    aiLoading.value = false
  }
}

function exportJson() {
  triggerDownload(JSON.stringify(filtered.value, null, 2), `audit-${Date.now()}.json`, 'application/json;charset=utf-8')
}

function exportCsv() {
  const headers = ['时间', '动作', '业务ID', '状态', '操作账户', '交易哈希', '返回信息']
  const rows = filtered.value.map((item) =>
    [
      item.createdAt,
      item.action,
      item.businessId,
      statusLabel(item.status),
      item.operatorAddr,
      item.txHash,
      item.message
    ].map(escapeCsv).join(',')
  )
  const csv = [headers.map(escapeCsv).join(','), ...rows].join('\n')
  triggerDownload(`\uFEFF${csv}`, `audit-${Date.now()}.csv`, 'text/csv;charset=utf-8')
}

function exportSummary() {
  const filters = [
    `动作筛选：${actionFilter.value || '全部'}`,
    `状态筛选：${statusFilter.value === 'all' ? '全部' : statusLabel(statusFilter.value)}`,
    `关键词：${keyword.value || '无'}`,
    `业务追溯：${traceId.value || '无'}`
  ]
  const anomalies = anomalyCards.value
    .slice(0, 5)
    .map((item, index) => `${index + 1}. ${item.action}：异常 ${item.abnormal} 次 / 总计 ${item.total} 次`)
    .join('\n')
  const warnings = recentWarnings.value
    .slice(0, 5)
    .map((item, index) => `${index + 1}. ${item.businessId || '未标识业务'} | ${item.action} | ${statusLabel(item.status)} | ${item.createdAt}`)
    .join('\n')
  const content = [
    '# 审计摘要',
    '',
    `导出时间：${new Date().toLocaleString()}`,
    '',
    '## 当前筛选条件',
    ...filters.map((line) => `- ${line}`),
    '',
    '## 总览',
    ...summaryCards.value.map((item) => `- ${item.label}：${item.value}`),
    '',
    '## 高频异常动作',
    anomalies || '- 当前没有异常动作聚类',
    '',
    '## 建议复核记录',
    warnings || '- 当前筛选条件下没有异常记录'
  ].join('\n')
  triggerDownload(content, `audit-summary-${Date.now()}.md`, 'text/markdown;charset=utf-8')
}

function exportTrace() {
  if (!traceId.value.trim() || !traceLogs.value.length) {
    ElMessage.info('请先输入业务ID并命中链路数据')
    return
  }
  const content = [
    '# 业务链路审计导出',
    '',
    `业务ID关键字：${traceId.value.trim()}`,
    `导出时间：${new Date().toLocaleString()}`,
    '',
    '## 业务链路',
    ...traceLogs.value.map(
      (item, index) =>
        `${index + 1}. ${item.createdAt} | ${item.action} | ${statusLabel(item.status)} | ${item.operatorAddr || '未知账户'} | ${item.txHash || '无交易哈希'}`
    )
  ].join('\n')
  triggerDownload(content, `audit-trace-${traceId.value.trim()}-${Date.now()}.md`, 'text/markdown;charset=utf-8')
}

function handleExport(command) {
  if (command === 'json') exportJson()
  if (command === 'csv') exportCsv()
  if (command === 'summary') exportSummary()
  if (command === 'trace') exportTrace()
}

onMounted(load)
</script>

<template>
  <section class="audit-shell">
    <article class="card audit-hero">
      <div>
        <span class="audit-badge">审计中心</span>
        <h2>审计分析台</h2>
      </div>
      <div class="row wrap">
        <button class="btn secondary" @click="load">{{ loading ? '刷新中...' : '刷新审计' }}</button>
        <el-dropdown trigger="click" @command="handleExport">
          <button class="btn">导出审计</button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="json">导出 JSON</el-dropdown-item>
              <el-dropdown-item command="csv">导出 CSV</el-dropdown-item>
              <el-dropdown-item command="summary">导出审计摘要</el-dropdown-item>
              <el-dropdown-item command="trace">导出当前业务链路</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </article>

    <section v-if="loading" class="audit-summary">
      <article v-for="i in 4" :key="i" class="audit-stat-card audit-stat-skeleton"><el-skeleton animated :rows="3" /></article>
    </section>
    <section v-else class="audit-summary">
      <article v-for="item in summaryCards" :key="item.label" class="audit-stat-card">
        <span>{{ item.label }}</span>
        <strong>{{ item.value }}</strong>
      </article>
    </section>

    <section class="audit-grid">
      <article class="card">
        <div class="panel-head compact">
          <h3>异常监测</h3>
        </div>
        <div v-if="loading"><el-skeleton animated :rows="8" /></div>
        <div v-else-if="anomalyCards.length" class="anomaly-list">
          <div v-for="item in anomalyCards" :key="item.action" class="anomaly-card">
            <div>
              <strong>{{ item.action }}</strong>
              <div class="sub">总次数 {{ item.total }}</div>
            </div>
            <span class="status-pill" :class="item.abnormal ? 'pill-danger' : 'pill-safe'">异常 {{ item.abnormal }}</span>
          </div>
        </div>
        <el-empty v-else description="暂无异常动作聚类" />
      </article>

      <article class="card">
        <div class="panel-head compact">
          <h3>业务追溯</h3>
        </div>
        <div class="trace-filter">
          <el-input v-model="traceId" clearable placeholder="输入业务ID追踪完整链路" />
        </div>
        <div v-if="traceLogs.length" class="trace-list">
          <div v-for="item in traceLogs" :key="item.id" class="trace-item">
            <div class="trace-dot"></div>
            <div>
              <strong>{{ item.action }}</strong>
              <div class="sub">{{ item.createdAt }} · {{ statusLabel(item.status) }}</div>
              <div class="sub mono-cell">{{ item.txHash }}</div>
            </div>
          </div>
        </div>
        <el-empty v-else description="输入业务ID后显示链路追溯结果" />
      </article>
    </section>

    <section class="audit-grid">
      <article class="card">
        <div class="panel-head compact">
          <h3>筛选条件</h3>
        </div>
        <div class="audit-filters">
          <el-select v-model="actionFilter" clearable placeholder="全部动作">
            <el-option label="全部动作" value="" />
            <el-option v-for="a in actions" :key="a" :label="a" :value="a" />
          </el-select>
          <el-select v-model="statusFilter" placeholder="全部状态">
            <el-option label="全部状态" value="all" />
            <el-option label="执行成功" value="0" />
            <el-option label="业务异常" value="1" />
            <el-option label="交易回滚" value="22" />
            <el-option label="系统异常" value="-1" />
          </el-select>
          <el-input v-model="keyword" clearable placeholder="搜索业务ID / 哈希 / 操作账户 / 消息" />
        </div>

        <div class="warning-block">
          <strong>近期异常</strong>
          <div v-if="recentWarnings.length" class="warning-list">
            <div v-for="item in recentWarnings" :key="item.id" class="warning-item">
              <span>{{ item.businessId || '未标识业务' }}</span>
              <span class="status-pill pill-danger">{{ statusLabel(item.status) }}</span>
            </div>
          </div>
          <el-empty v-else description="当前筛选条件下没有异常记录" />
        </div>
      </article>

      <article class="card">
        <div class="panel-head compact">
          <h3>AI 审计摘要</h3>
          <button class="btn secondary" :disabled="aiLoading" @click="generateAISummary">{{ aiLoading ? '生成中...' : '生成摘要' }}</button>
        </div>
        <div v-if="aiLoading"><el-skeleton animated :rows="8" /></div>
        <div v-else-if="aiSummary" class="ai-markdown" v-html="renderMarkdown(aiSummary)"></div>
        <el-empty v-else description="点击生成摘要，输出异常模式和复核建议" />
      </article>
    </section>

    <article class="card">
      <div class="panel-head compact">
        <h3>审计明细</h3>
      </div>
      <div v-if="isEmptyAudit" class="audit-empty-wrap">
        <el-empty description="当前还没有审计日志数据" />
      </div>
      <div v-else class="table-wrap">
        <table class="table">
          <thead>
            <tr>
              <th>时间</th>
              <th>动作</th>
              <th>业务ID</th>
              <th>状态</th>
              <th>操作账户</th>
              <th>交易哈希</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in filtered" :key="item.id">
              <td>{{ item.createdAt }}</td>
              <td>{{ item.action }}</td>
              <td>{{ item.businessId }}</td>
              <td><span class="status-pill" :class="Number(item.status) === 0 ? 'pill-safe' : 'pill-danger'">{{ statusLabel(item.status) }}</span></td>
              <td class="mono-cell">{{ item.operatorAddr }}</td>
              <td class="mono-cell">{{ item.txHash }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <pre v-if="err" class="log">{{ err }}</pre>
      <div v-if="stats" class="sub" style="margin-top:10px;">工作流累计请求 {{ stats.totalRequests || 0 }}</div>
    </article>
  </section>
</template>

<style scoped>
.audit-shell { display: grid; gap: 14px; }
.audit-hero {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background:
    radial-gradient(circle at right top, rgba(39, 193, 255, 0.15), transparent 32%),
    linear-gradient(135deg, rgba(20, 67, 166, 0.96), rgba(46, 117, 255, 0.88));
  color: #edf5ff;
}
.audit-badge {
  display: inline-flex;
  padding: 6px 10px;
  border-radius: 999px;
  border: 1px solid rgba(255,255,255,.22);
  background: rgba(255,255,255,.12);
  font-size: 12px;
  margin-bottom: 10px;
}
.audit-hero h2 { margin: 0; color: #fff; }
.audit-summary {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 14px;
}
.audit-stat-card {
  display: grid;
  gap: 6px;
  padding: 16px;
  border-radius: 16px;
  border: 1px solid #d9e8ff;
  background: linear-gradient(180deg, rgba(250, 252, 255, 0.96), rgba(239, 247, 255, 0.9));
}
.audit-stat-card span { color: #6e84a4; font-size: 12px; }
.audit-stat-card strong { color: #14366e; font-size: 28px; }
.audit-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}
.panel-head.compact {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 14px;
}
.anomaly-list,
.trace-list,
.warning-list {
  display: grid;
  gap: 10px;
}
.anomaly-card,
.warning-item {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  align-items: center;
  padding: 12px;
  border-radius: 14px;
  border: 1px solid #d9e8ff;
  background: linear-gradient(180deg, rgba(250, 252, 255, 0.96), rgba(242, 248, 255, 0.92));
}
.trace-item {
  display: grid;
  grid-template-columns: 16px 1fr;
  gap: 10px;
  padding: 10px 0;
}
.trace-dot {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  background: linear-gradient(135deg, #2f72ef, #49c7ff);
  margin-top: 6px;
}
.audit-filters {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
  margin-bottom: 14px;
}
.trace-filter {
  margin-bottom: 12px;
}
.warning-block {
  display: grid;
  gap: 10px;
}
.table-wrap { overflow: auto; }
.mono-cell { word-break: break-all; }
.audit-stat-skeleton { min-height: 128px; }
.audit-empty-wrap {
  min-height: 240px;
  display: grid;
  place-items: center;
}
.pill-safe {
  background: rgba(46, 188, 139, 0.1);
  border-color: rgba(46, 188, 139, 0.3);
  color: #159468;
}
.pill-danger {
  background: rgba(245, 90, 122, 0.12);
  border-color: rgba(245, 90, 122, 0.3);
  color: #c63f60;
}
.ai-markdown :deep(h2),
.ai-markdown :deep(h3) {
  margin: 0 0 10px;
  color: #153a77;
}
.ai-markdown :deep(p) { margin: 0 0 10px; line-height: 1.65; }
.ai-markdown :deep(ul) { margin: 0 0 10px; padding-left: 18px; }
.ai-markdown :deep(li) { margin: 5px 0; }
@media (max-width: 980px) {
  .audit-hero,
  .audit-grid,
  .audit-filters {
    grid-template-columns: 1fr;
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
