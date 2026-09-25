<script setup>
import { computed, onMounted, ref } from 'vue'
import DOMPurify from 'dompurify'
import { marked } from 'marked'
import { request, requestPreferDB } from '../api/http'

const overview = ref(null)
const chainStats = ref(null)
const txPage = ref(1)
const txPageSize = 30
const txTotal = ref(0)
const txHistory = ref([])
const txAnalysis = ref(null)
const txAnalyzing = ref(false)
const txDialogVisible = ref(false)
const txInsight = ref('')
const txInsightLoading = ref(false)
const err = ref('')
const loading = ref(false)

marked.setOptions({
  gfm: true,
  breaks: true
})

const statusNameMap = {
  PENDING: '待确认',
  CONFIRMED: '已确认',
  FINANCED: '已融资',
  SETTLED: '已结清'
}

const actionNameMap = {
  register_enterprise: '企业登记上链',
  issue_receivable: '创建应收凭证',
  confirm_receivable: '应收确认',
  record_financing: '融资登记',
  record_settlement: '应收结算'
}

function normalizeChainNumber(value) {
  if (value === null || value === undefined || value === '') return '0'
  const raw = String(value).trim()
  if (!raw) return '0'
  if (raw.startsWith('0x') || raw.startsWith('0X')) {
    try {
      return BigInt(raw).toString(10)
    } catch (_) {
      return raw
    }
  }
  return raw
}

const total = computed(() => Number(overview.value?.receivableSize || 0))

const statusList = computed(() => {
  const counters = overview.value?.statusCounters || {}
  return ['PENDING', 'CONFIRMED', 'FINANCED', 'SETTLED'].map((key) => {
    const count = Number(counters[key] || 0)
    const ratio = total.value ? (count / total.value) * 100 : 0
    return {
      key,
      label: statusNameMap[key],
      count,
      ratio: ratio.toFixed(1)
    }
  })
})

const hotActions = computed(() => {
  const logs = overview.value?.recentTxLogs || []
  const actionMap = {}
  logs.forEach((item) => {
    const key = item.action || 'unknown'
    actionMap[key] = (actionMap[key] || 0) + 1
  })
  return Object.entries(actionMap)
    .map(([key, count]) => ({
      key,
      label: actionNameMap[key] || key,
      count
    }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 5)
})

const topMetrics = computed(() => {
  const data = overview.value || {}
  const chain = chainStats.value || {}
  return [
    { label: '平台用户', value: data.users || 0 },
    { label: '应收总量', value: data.receivableSize || 0 },
    { label: '融资余额', value: data.totalAmount || 0 },
    { label: '链上高度', value: chain.blockHeight || 0 }
  ]
})

const chainCards = computed(() => {
  const chain = chainStats.value || {}
  return [
    { label: '累计交易数', value: normalizeChainNumber(chain.totalTxCount) },
    { label: '失败交易数', value: normalizeChainNumber(chain.failedTxCount) },
    { label: '节点数量', value: chain.nodeCount || 0 },
    { label: '群组编号', value: chain.groupId || '-' }
  ]
})

const nodeIds = computed(() => chainStats.value?.nodeIds || [])
const rollingNodeIds = computed(() => {
  const list = nodeIds.value || []
  if (!list.length) return []
  if (list.length === 1) return [list[0], list[0]]
  return [list[0], list[1]]
})

function actionLabel(key) {
  return actionNameMap[key] || key
}

function statusLabel(key) {
  return statusNameMap[key] || key
}

function txStatusLabel(status) {
  const num = Number(status)
  if (Number.isNaN(num)) return String(status || '-')
  return num === 0 ? '成功' : `异常(${num})`
}

async function load() {
  loading.value = true
  try {
    const [overviewResp, chainResp, txResp] = await Promise.all([
      requestPreferDB('/api/db/dashboard/overview', '/api/dashboard/overview'),
      request('/api/chain/stats', { withRole: true }).catch(() => null),
      request(`/api/tx-logs?page=${txPage.value}&pageSize=${txPageSize}`, { withRole: true }).catch(() => ({ data: [], total: 0 }))
    ])
    overview.value = overviewResp
    chainStats.value = chainResp
    txHistory.value = txResp.data || []
    txTotal.value = Number(txResp.total || 0)
    err.value = ''
  } catch (error) {
    err.value = JSON.stringify(error, null, 2)
  } finally {
    loading.value = false
  }
}

async function changeTxPage(next) {
  if (next < 1) return
  txPage.value = next
  await load()
}

async function analyzeTx(hash) {
  if (!hash) return
  txDialogVisible.value = true
  txAnalyzing.value = true
  txInsight.value = ''
  try {
    txAnalysis.value = await request(`/api/tx-logs/${hash}/analyze`, { withRole: true })
  } catch (error) {
    txAnalysis.value = { error: error?.error || '交易解析失败' }
  } finally {
    txAnalyzing.value = false
  }
}

async function generateTxInsight() {
  if (!txAnalysis.value || txAnalysis.value.error) return
  txInsightLoading.value = true
  try {
    const prompt = [
      '你是供应链金融链上交易审查助手。',
      `交易哈希: ${txAnalysis.value.hash}`,
      `业务编号: ${txAnalysis.value.summary?.businessId || '未知'}`,
      `执行状态: ${txAnalysis.value.summary?.statusText || '未知'}`,
      `区块高度: ${txAnalysis.value.summary?.blockNumber || '-'}`,
      `发起地址: ${txAnalysis.value.summary?.from || '-'}`,
      `目标地址: ${txAnalysis.value.summary?.to || '-'}`,
      `Gas Used: ${txAnalysis.value.summary?.gasUsed || '-'}`,
      `Receipt Logs 数量: ${txAnalysis.value.summary?.logCount || 0}`,
      `审计动作: ${txAnalysis.value.audit?.action || '-'}`,
      `审计消息: ${txAnalysis.value.audit?.message || '-'}`,
      '请输出三部分：1. 交易摘要 2. 风险或异常判断 3. 建议人工复核点。使用 markdown。'
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
    txInsight.value = DOMPurify.sanitize(marked.parse(resp?.reply || ''))
  } catch (error) {
    txInsight.value = DOMPurify.sanitize(marked.parse(error?.error || 'AI 分析失败'))
  } finally {
    txInsightLoading.value = false
  }
}

onMounted(load)
</script>

<template>
  <section class="dashboard-shell">
    <article class="card dashboard-hero">
      <div class="hero-copy">
        <span class="hero-badge">运营总览</span>
        <h2>供应链金融链上运营总览</h2>
      </div>
      <div class="hero-actions">
        <button class="btn neon" @click="load">{{ loading ? '刷新中...' : '刷新总览' }}</button>
      </div>
    </article>

    <section class="dashboard-metrics">
      <article v-for="item in topMetrics" :key="item.label" class="metric-card">
        <span>{{ item.label }}</span>
        <strong>{{ item.value }}</strong>
      </article>
    </section>

    <section v-if="overview" class="dashboard-grid">
      <article class="card overview-panel">
        <div class="panel-head">
          <div>
            <h3>业务状态流转</h3>
          </div>
          <span class="panel-tag">总计 {{ total }}</span>
        </div>

        <div class="status-stack">
          <div v-for="item in statusList" :key="item.key" class="status-line">
            <div class="row between">
              <div class="status-meta">
                <strong>{{ item.label }}</strong>
                <span>{{ item.count }} 笔</span>
              </div>
              <span class="status-ratio">{{ item.ratio }}%</span>
            </div>
            <div class="status-track">
              <div class="status-fill" :style="{ width: item.ratio + '%' }"></div>
            </div>
          </div>
        </div>
      </article>

      <article class="card overview-panel">
        <div class="panel-head">
          <div>
            <h3>区块链运行面板</h3>
          </div>
          <span class="panel-tag">FISCO BCOS</span>
        </div>

        <div class="chain-metrics">
          <div v-for="item in chainCards" :key="item.label" class="chain-card-mini">
            <span>{{ item.label }}</span>
            <strong>{{ item.value }}</strong>
          </div>
        </div>

        <div class="chain-meta row wrap">
          <span class="status-pill">链 ID · {{ chainStats?.chainId || '-' }}</span>
          <span class="status-pill">版本 · {{ chainStats?.fiscoVersion || '-' }}</span>
        </div>

        <div class="node-box">
          <div class="node-head">
            <strong>节点 ID</strong>
            <span class="sub">{{ nodeIds.length }} 个节点</span>
          </div>
          <div v-if="rollingNodeIds.length" class="node-ticker">
            <div class="node-track">
              <code v-for="(id, index) in [...rollingNodeIds, ...rollingNodeIds]" :key="`${id}-${index}`">{{ id }}</code>
            </div>
          </div>
          <div v-else class="sub">暂无节点数据</div>
        </div>
      </article>

      <article class="card overview-panel">
        <div class="panel-head">
          <div>
            <h3>近期高频动作</h3>
          </div>
        </div>
        <div class="action-list" v-if="hotActions.length">
          <div v-for="item in hotActions" :key="item.key" class="action-row">
            <div>
              <strong>{{ item.label }}</strong>
              <div class="sub">动作标识：{{ item.key }}</div>
            </div>
            <span class="action-count">{{ item.count }}</span>
          </div>
        </div>
        <div v-else class="sub">暂无近期动作数据</div>
      </article>
    </section>

    <article class="card full trade-panel">
      <div class="panel-head">
        <div>
          <h3>最近交易记录</h3>
        </div>
        <span class="panel-tag">最新 30 条</span>
      </div>
      <div class="trade-table-wrap">
        <table class="table trade-table">
          <thead>
            <tr>
              <th>时间</th>
              <th>动作</th>
              <th>业务编号</th>
              <th>执行结果</th>
              <th>交易哈希</th>
              <th>分析</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in txHistory" :key="item.id">
              <td>{{ item.createdAt }}</td>
              <td>{{ actionLabel(item.action) }}</td>
              <td>{{ item.businessId }}</td>
              <td>{{ txStatusLabel(item.status) }}</td>
              <td class="hash-cell">{{ item.txHash }}</td>
              <td><button class="btn mini secondary" @click="analyzeTx(item.txHash)">解析</button></td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="row between" style="margin-top: 12px;">
        <span class="sub">历史总数 {{ txTotal }} · 当前第 {{ txPage }} 页</span>
        <div class="row wrap">
          <button class="btn mini ghost" :disabled="txPage <= 1 || loading" @click="changeTxPage(txPage - 1)">上一页</button>
          <button class="btn mini ghost" :disabled="txPage * txPageSize >= txTotal || loading" @click="changeTxPage(txPage + 1)">下一页</button>
        </div>
      </div>
      <pre v-if="err" class="log">{{ err }}</pre>
    </article>

    <el-dialog v-model="txDialogVisible" width="860px" title="交易哈希分析" class="tx-dialog">
      <div v-if="txAnalyzing" class="sub">解析中...</div>
      <div v-else-if="txAnalysis?.error" class="sub">{{ txAnalysis.error }}</div>
      <div v-else-if="txAnalysis" class="tx-analysis-card">
        <div class="analysis-grid">
          <div class="analysis-item">
            <span>交易哈希</span>
            <strong>{{ txAnalysis.hash }}</strong>
          </div>
          <div class="analysis-item">
            <span>业务编号</span>
            <strong>{{ txAnalysis.summary?.businessId || txAnalysis.audit?.businessId || '-' }}</strong>
          </div>
          <div class="analysis-item">
            <span>执行状态</span>
            <strong>{{ txAnalysis.summary?.statusText || '-' }}</strong>
          </div>
          <div class="analysis-item">
            <span>审计动作</span>
            <strong>{{ actionLabel(txAnalysis.audit?.action) || '-' }}</strong>
          </div>
          <div class="analysis-item">
            <span>区块高度</span>
            <strong>{{ txAnalysis.summary?.blockNumber || '-' }}</strong>
          </div>
          <div class="analysis-item">
            <span>Gas Used</span>
            <strong>{{ txAnalysis.summary?.gasUsed || '-' }}</strong>
          </div>
        </div>

        <div v-if="txAnalysis.audit" class="analysis-block">
          <div class="panel-head">
            <div><h3>关联审计记录</h3></div>
          </div>
          <div class="analysis-grid">
            <div class="analysis-item">
              <span>操作账户</span>
              <strong>{{ txAnalysis.audit.operatorAddr || '-' }}</strong>
            </div>
            <div class="analysis-item">
              <span>审计消息</span>
              <strong>{{ txAnalysis.audit.message || '-' }}</strong>
            </div>
          </div>
        </div>

        <div class="analysis-block">
          <div class="panel-head">
            <div><h3>Receipt Logs 摘要</h3></div>
            <span class="panel-tag">{{ txAnalysis.receipt?.logCount || 0 }} 条</span>
          </div>
          <div v-if="txAnalysis.receipt?.logsSummary?.length" class="result-grid">
            <div v-for="(item, index) in txAnalysis.receipt.logsSummary" :key="index" class="analysis-item">
              <span>{{ item.address }}</span>
              <strong>Topics {{ item.topicCount }}</strong>
              <div class="sub">{{ item.dataPreview }}</div>
            </div>
          </div>
          <div v-else class="sub">当前交易没有 receipt logs</div>
        </div>

        <el-collapse class="analysis-block">
          <el-collapse-item title="交易 Input / Output 折叠查看" name="io">
            <div class="analysis-grid">
              <div class="analysis-item">
                <span>Input</span>
                <strong>{{ txAnalysis.transaction?.inputPreview || txAnalysis.receipt?.inputPreview || '-' }}</strong>
              </div>
              <div class="analysis-item">
                <span>Output</span>
                <strong>{{ txAnalysis.receipt?.outputPreview || '-' }}</strong>
              </div>
            </div>
          </el-collapse-item>
        </el-collapse>

        <div class="analysis-block">
          <div class="panel-head">
            <div><h3>AI 交易分析</h3></div>
            <button class="btn mini secondary" :disabled="txInsightLoading" @click="generateTxInsight">
              {{ txInsightLoading ? '分析中...' : '生成 AI 分析' }}
            </button>
          </div>
          <div v-if="txInsight" class="ai-markdown" v-html="txInsight"></div>
          <div v-else class="sub">点击按钮生成 AI 审查结论</div>
        </div>
      </div>
    </el-dialog>
  </section>
</template>

<style scoped>
.dashboard-shell {
  display: grid;
  gap: 14px;
}

.dashboard-hero {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: center;
  padding: 20px;
  background:
    radial-gradient(circle at right top, rgba(39, 193, 255, 0.18), transparent 30%),
    linear-gradient(135deg, rgba(20, 67, 166, 0.96), rgba(46, 117, 255, 0.88));
  color: #eef5ff;
  overflow: hidden;
}

.hero-badge {
  display: inline-flex;
  padding: 6px 10px;
  border-radius: 999px;
  border: 1px solid rgba(255, 255, 255, 0.22);
  background: rgba(255, 255, 255, 0.12);
  font-size: 12px;
  margin-bottom: 10px;
}

.dashboard-hero h2 {
  margin: 0 0 8px;
  font-size: 34px;
  color: #fff;
}

.dashboard-hero p {
  margin: 0;
  max-width: 720px;
  color: rgba(238, 245, 255, 0.86);
}

.dashboard-metrics {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(210px, 1fr));
  gap: 14px;
}

.metric-card {
  position: relative;
  display: grid;
  gap: 8px;
  padding: 18px;
  border: 1px solid #d9e8ff;
  border-radius: 18px;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.92), rgba(241, 248, 255, 0.88));
  box-shadow: 0 14px 30px rgba(40, 101, 203, 0.12);
}

.metric-card span,
.chain-card-mini span {
  color: #6780a3;
  font-size: 13px;
}

.metric-card strong,
.chain-card-mini strong {
  font-size: 30px;
  color: #12356f;
}

.metric-card small {
  color: #7f95b3;
  font-size: 12px;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: 1.2fr 1.15fr 0.9fr;
  gap: 14px;
}

.overview-panel {
  min-height: 360px;
}

.panel-head {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: flex-start;
  margin-bottom: 14px;
}

.panel-head p {
  margin: 0;
  color: #6f86a8;
  font-size: 13px;
}

.panel-tag {
  display: inline-flex;
  align-items: center;
  padding: 6px 10px;
  border-radius: 999px;
  background: #edf5ff;
  border: 1px solid #d6e6ff;
  color: #2f62b8;
  font-size: 12px;
}

.status-stack {
  display: grid;
  gap: 14px;
}

.status-line {
  display: grid;
  gap: 8px;
}

.status-meta {
  display: grid;
  gap: 2px;
}

.status-meta strong {
  font-size: 15px;
}

.status-meta span,
.status-ratio {
  color: #6d82a0;
  font-size: 12px;
}

.status-track {
  height: 10px;
  background: linear-gradient(90deg, #eef4ff, #e5f1ff);
  border-radius: 999px;
  overflow: hidden;
}

.status-fill {
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, #2563eb, #41c0ff);
  box-shadow: 0 0 16px rgba(65, 192, 255, 0.28);
}

.chain-metrics {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
  margin-bottom: 12px;
}

.chain-card-mini {
  display: grid;
  gap: 6px;
  padding: 14px;
  border-radius: 14px;
  border: 1px solid #dbe9ff;
  background: linear-gradient(180deg, rgba(242, 248, 255, 0.96), rgba(255, 255, 255, 0.94));
}

.chain-meta {
  margin-bottom: 12px;
}

.node-box {
  padding: 12px;
  border-radius: 14px;
  background: rgba(243, 249, 255, 0.92);
  border: 1px solid #dbe9ff;
}

.node-head {
  display: flex;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 8px;
}

.node-ticker {
  position: relative;
  overflow: hidden;
  height: 96px;
}

.node-track {
  display: grid;
  gap: 8px;
  animation: node-scroll 12s linear infinite;
}

.node-track code,
.hash-cell {
  word-break: break-all;
}

.node-track code {
  display: block;
  padding: 8px 10px;
  border-radius: 10px;
  background: #edf5ff;
  border: 1px solid #d6e6ff;
  color: #34557b;
  font-size: 12px;
}

.action-list {
  display: grid;
  gap: 10px;
}

.action-row {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  align-items: center;
  padding: 12px;
  border-radius: 14px;
  background: rgba(245, 250, 255, 0.95);
  border: 1px solid #dce9ff;
}

.action-count {
  min-width: 40px;
  height: 40px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 999px;
  background: linear-gradient(135deg, #2f72ef, #43c3ff);
  color: white;
  font-weight: 700;
}

.trade-panel {
  padding-top: 16px;
}

.trade-table-wrap {
  overflow: auto;
}

.trade-table th {
  white-space: nowrap;
}

.tx-analysis-card {
  display: grid;
  gap: 14px;
}

.analysis-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.result-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 10px;
}

.analysis-item {
  display: grid;
  gap: 4px;
  padding: 12px;
  border-radius: 14px;
  border: 1px solid #d8e8ff;
  background: rgba(255, 255, 255, 0.9);
  min-width: 0;
}

.analysis-item span {
  color: #6a80a0;
  font-size: 12px;
}

.analysis-item strong {
  color: #16386f;
  font-size: 15px;
  word-break: break-all;
  overflow-wrap: anywhere;
}

.analysis-item .sub,
.analysis-item span {
  overflow-wrap: anywhere;
  word-break: break-all;
}

.analysis-block {
  padding: 14px;
  border-radius: 16px;
  border: 1px solid #d9e8ff;
  background: linear-gradient(180deg, rgba(248, 252, 255, 0.97), rgba(240, 247, 255, 0.94));
}

.ai-markdown :deep(h2),
.ai-markdown :deep(h3) {
  margin: 0 0 10px;
  color: #153a77;
}

.ai-markdown :deep(p) { margin: 0 0 10px; line-height: 1.65; }
.ai-markdown :deep(ul) { margin: 0 0 10px; padding-left: 18px; }
.ai-markdown :deep(li) { margin: 5px 0; }

@keyframes node-scroll {
  0% {
    transform: translateY(0);
  }
  100% {
    transform: translateY(calc(-50% - 4px));
  }
}

@media (max-width: 1360px) {
  .dashboard-grid {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 980px) {
  .dashboard-hero {
    flex-direction: column;
    align-items: flex-start;
  }

  .dashboard-grid,
  .dashboard-metrics,
  .analysis-grid,
  .result-grid {
    grid-template-columns: 1fr;
  }
}
</style>
