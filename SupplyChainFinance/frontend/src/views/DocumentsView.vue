<script setup>
import { computed, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import DOMPurify from 'dompurify'
import { marked } from 'marked'
import { request, requestPreferDB, uploadFiles } from '../api/http'

const files = ref([])
const uploading = ref(false)
const loading = ref(false)
const aiLoading = ref(false)
const log = ref('')
const query = ref('')
const typeFilter = ref('ALL')
const sourceFilter = ref('ALL')
const dbDocs = ref([])
const receivables = ref([])
const aiAdvice = ref('')

marked.setOptions({
  gfm: true,
  breaks: true
})

const sourceOptions = {
  ALL: '全部来源',
  db: '数据库',
  local: '本地索引',
  ipfs: 'IPFS'
}

const typeNameMap = {
  contract: '合同',
  invoice: '发票',
  logistics: '签收/运输',
  receipt: '签收/运输',
  voucher: '业务凭证',
  other: '其他'
}

const mergedDocs = computed(() => {
  return dbDocs.value.map((item) => {
    const relation = item.receivableId || item.orderId
      ? receivables.value.find((row) => row.receivableId === item.receivableId || row.orderId === item.orderId) || {
          receivableId: item.receivableId,
          orderId: item.orderId
        }
      : null
    return {
      ...item,
      typeLabel: typeNameMap[item.docType] || '其他',
      sourceLabel: sourceOptions[item.source] || item.source || '数据库',
      relation,
      createdLabel: formatDate(item.createdAt)
    }
  })
})

const filteredDocs = computed(() => {
  const key = query.value.trim().toLowerCase()
  return mergedDocs.value.filter((item) => {
    const typeOK = typeFilter.value === 'ALL' ? true : item.docType === typeFilter.value
    const sourceOK = sourceFilter.value === 'ALL' ? true : item.source === sourceFilter.value
    const keyOK = !key
      ? true
      : [item.name, item.cid, item.tag, item.operator, item.receivableId, item.orderId].some((value) =>
          String(value || '').toLowerCase().includes(key)
        )
    return typeOK && sourceOK && keyOK
  })
})

const docStats = computed(() => {
  const docs = mergedDocs.value
  return {
    total: docs.length,
    linked: docs.filter((item) => item.receivableId || item.orderId).length,
    contract: docs.filter((item) => item.docType === 'contract').length,
    invoice: docs.filter((item) => item.docType === 'invoice').length
  }
})

const completenessList = computed(() => {
  return receivables.value.slice(0, 12).map((item) => {
    const linked = mergedDocs.value.filter((doc) => doc.receivableId === item.receivableId || doc.orderId === item.orderId)
    const present = new Set(linked.map((doc) => doc.docType))
    const required = ['contract', 'invoice', 'logistics']
    const missing = required.filter((type) => !present.has(type)).map((type) => typeNameMap[type])
    const ratio = Math.round(((required.length - missing.length) / required.length) * 100)
    return {
      receivableId: item.receivableId,
      orderId: item.orderId,
      statusText: item.statusText,
      ratio,
      missing,
      linkedCount: linked.length
    }
  })
})

function formatDate(value) {
  if (!value) return '-'
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? String(value) : date.toLocaleString()
}

function riskToneByRatio(ratio) {
  if (ratio >= 100) return 'safe'
  if (ratio >= 67) return 'watch'
  return 'danger'
}

function renderMarkdown(content) {
  return DOMPurify.sanitize(marked.parse(String(content || '')))
}

async function load() {
  loading.value = true
  try {
    const [docsResp, recResp] = await Promise.all([
      requestPreferDB('/api/db/documents', '/api/db/documents').catch(() => ({ data: [] })),
      requestPreferDB('/api/db/receivables', '/api/receivable').catch(() => ({ data: [] }))
    ])
    dbDocs.value = docsResp.data || []
    receivables.value = recResp.data || []
  } finally {
    loading.value = false
  }
}

async function submit() {
  if (!files.value.length) return
  uploading.value = true
  try {
    const user = JSON.parse(localStorage.getItem('scf_user') || 'null')
    const resp = await uploadFiles('/api/ipfs/upload', files.value, true, {
      tag: '合同/凭证',
      operator: user?.username || user?.address || '',
      docType: 'contract',
      status: 'uploaded',
      source: 'ipfs'
    })
    log.value = `上传成功\nCID 数量: ${(resp.cids || []).length}`
    files.value = []
    ElMessage.success('文档已上传到 IPFS，并同步写入数据库')
    await load()
  } catch (e) {
    log.value = `上传失败\n${JSON.stringify(e, null, 2)}`
    ElMessage.error('文档上传失败')
  } finally {
    uploading.value = false
  }
}

async function generateAIAdvice() {
  if (!mergedDocs.value.length) {
    ElMessage.info('当前没有可分析的文档')
    return
  }
  aiLoading.value = true
  try {
    const summary = completenessList.value
      .slice(0, 6)
      .map((item) => `${item.receivableId} 完整度 ${item.ratio}% 缺失 ${item.missing.join('/') || '无'}`)
      .join('\n')
    const prompt = [
      '你是供应链金融凭证审查助手。',
      `文档总数: ${docStats.value.total}`,
      `已关联应收: ${docStats.value.linked}`,
      `合同数: ${docStats.value.contract}`,
      `发票数: ${docStats.value.invoice}`,
      '以下是文档完整性摘要：',
      summary,
      '请输出三部分：1. 文档治理判断 2. 缺失项风险 3. 补件建议。使用 markdown。'
    ].join('\n')
    const resp = await request('/api/ai/chat', {
      method: 'POST',
      withRole: true,
      body: {
        message: prompt,
        scene: 'finance',
        guard: {
          privacyMode: true,
          localMasking: true,
          noRetention: true,
          minimalDisclosure: true
        }
      }
    })
    aiAdvice.value = resp?.reply || ''
  } catch (_) {
    aiAdvice.value = [
      '## 文档治理判断',
      `当前共有 **${docStats.value.total}** 份文档，已关联应收 **${docStats.value.linked}** 份。`,
      '## 缺失项风险',
      '- 缺合同会影响债权真实性判断。',
      '- 缺发票会影响金额与税务一致性审核。',
      '- 缺签收/运输信息会削弱履约真实性证明。',
      '## 补件建议',
      '- 优先补齐高金额应收的合同与发票。',
      '- 文档命名中加入应收ID或订单ID，便于系统自动关联。'
    ].join('\n')
  } finally {
    aiLoading.value = false
  }
}

function openCID(cid) {
  window.open(`http://localhost:8081/ipfs/${cid}`, '_blank')
}

function onFileChange(e) {
  files.value = Array.from(e.target.files || [])
}

onMounted(load)
</script>

<template>
  <section class="docs-shell">
    <article class="card docs-hero">
      <div>
        <span class="docs-badge">文档中心</span>
        <h2>凭证治理台</h2>
      </div>
      <div class="row wrap">
        <button class="btn secondary" @click="load">{{ loading ? '刷新中...' : '刷新文档' }}</button>
        <button class="btn" :disabled="uploading" @click="submit">{{ uploading ? '上传中...' : '上传到 IPFS' }}</button>
        <input type="file" multiple @change="onFileChange" />
      </div>
    </article>

    <section v-if="loading" class="docs-stats">
      <article v-for="i in 4" :key="i" class="doc-stat-card"><el-skeleton animated :rows="3" /></article>
    </section>
    <section v-else class="docs-stats">
      <article class="doc-stat-card"><span>文档总量</span><strong>{{ docStats.total }}</strong></article>
      <article class="doc-stat-card"><span>已关联应收</span><strong>{{ docStats.linked }}</strong></article>
      <article class="doc-stat-card"><span>合同数量</span><strong>{{ docStats.contract }}</strong></article>
      <article class="doc-stat-card"><span>发票数量</span><strong>{{ docStats.invoice }}</strong></article>
    </section>

    <section class="docs-grid">
      <article class="card">
        <div class="panel-head compact">
          <h3>文档检索</h3>
        </div>
        <div class="doc-filters">
          <el-input v-model="query" clearable placeholder="按文件名 / CID / 标签 / 应收ID 搜索" />
          <el-select v-model="typeFilter">
            <el-option label="全部类型" value="ALL" />
            <el-option label="合同" value="contract" />
            <el-option label="发票" value="invoice" />
            <el-option label="签收/运输" value="logistics" />
            <el-option label="业务凭证" value="voucher" />
            <el-option label="其他" value="other" />
          </el-select>
          <el-select v-model="sourceFilter">
            <el-option v-for="(label, key) in sourceOptions" :key="key" :label="label" :value="key" />
          </el-select>
        </div>

        <div v-if="!filteredDocs.length" class="docs-empty-wrap">
          <el-empty description="当前筛选条件下没有文档" />
        </div>
        <div v-else class="table-wrap">
          <table class="table">
            <thead>
              <tr>
                <th>文件</th>
                <th>类型</th>
                <th>关联应收</th>
                <th>CID</th>
                <th>来源</th>
                <th>时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="doc in filteredDocs" :key="doc.id">
                <td>{{ doc.name }}</td>
                <td><span class="status-pill">{{ doc.typeLabel }}</span></td>
                <td>{{ doc.relation?.receivableId || '-' }}</td>
                <td class="mono-cell">{{ doc.cid }}</td>
                <td>{{ doc.sourceLabel }}</td>
                <td>{{ doc.createdLabel }}</td>
                <td><button class="btn mini" @click="openCID(doc.cid)">查看</button></td>
              </tr>
            </tbody>
          </table>
        </div>
        <pre class="log">{{ log || '上传后的文档会直接写入数据库，并按应收字段做严格关联。' }}</pre>
      </article>

      <article class="card">
        <div class="panel-head compact">
          <h3>完整性校验</h3>
        </div>
        <div v-if="!completenessList.length" class="docs-empty-wrap">
          <el-empty description="当前没有可校验的应收样本" />
        </div>
        <div v-else class="completeness-list">
          <div v-for="item in completenessList" :key="item.receivableId" class="completeness-card">
            <div class="row between">
              <strong>{{ item.receivableId }}</strong>
              <span class="status-pill" :class="`pill-${riskToneByRatio(item.ratio)}`">{{ item.ratio }}%</span>
            </div>
            <div class="sub">{{ item.orderId }} · {{ item.statusText }}</div>
            <div class="bucket-track">
              <div class="bucket-fill" :style="{ width: `${item.ratio}%` }"></div>
            </div>
            <div class="sub">已关联 {{ item.linkedCount }} 份文档</div>
            <div class="missing-list">
              <span v-for="miss in (item.missing.length ? item.missing : ['已齐备'])" :key="miss" class="status-pill">{{ miss }}</span>
            </div>
          </div>
        </div>
      </article>
    </section>

    <article class="card">
      <div class="panel-head compact">
        <h3>AI 文档建议</h3>
        <button class="btn secondary" :disabled="aiLoading" @click="generateAIAdvice">{{ aiLoading ? '生成中...' : '生成建议' }}</button>
      </div>
      <div v-if="aiLoading"><el-skeleton animated :rows="8" /></div>
      <div v-else-if="aiAdvice" class="ai-markdown" v-html="renderMarkdown(aiAdvice)"></div>
      <el-empty v-else description="点击生成建议，获取文档治理和补件建议" />
    </article>
  </section>
</template>

<style scoped>
.docs-shell { display: grid; gap: 14px; }
.docs-hero {
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
.docs-badge {
  display: inline-flex;
  padding: 6px 10px;
  border-radius: 999px;
  border: 1px solid rgba(255,255,255,.22);
  background: rgba(255,255,255,.12);
  font-size: 12px;
  margin-bottom: 10px;
}
.docs-hero h2 { margin: 0; color: #fff; }
.docs-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 14px;
}
.doc-stat-card {
  display: grid;
  gap: 6px;
  padding: 16px;
  border-radius: 16px;
  border: 1px solid #d9e8ff;
  background: linear-gradient(180deg, rgba(250, 252, 255, 0.96), rgba(239, 247, 255, 0.9));
}
.doc-stat-card span { color: #6e84a4; font-size: 12px; }
.doc-stat-card strong { color: #14366e; font-size: 28px; }
.docs-grid {
  display: grid;
  grid-template-columns: 1.15fr 0.85fr;
  gap: 14px;
}
.panel-head.compact {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 14px;
}
.doc-filters {
  display: grid;
  grid-template-columns: 1.2fr 180px 160px;
  gap: 10px;
  margin-bottom: 12px;
}
.table-wrap { overflow: auto; }
.mono-cell { word-break: break-all; }
.docs-empty-wrap {
  min-height: 220px;
  display: grid;
  place-items: center;
}
.completeness-list {
  display: grid;
  gap: 10px;
}
.completeness-card {
  display: grid;
  gap: 8px;
  padding: 12px;
  border-radius: 14px;
  border: 1px solid #d9e8ff;
  background: linear-gradient(180deg, rgba(250, 252, 255, 0.96), rgba(242, 248, 255, 0.92));
}
.missing-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.bucket-track {
  height: 10px;
  border-radius: 999px;
  background: #e8f0ff;
  overflow: hidden;
}
.bucket-fill {
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, #2f72ef, #49c7ff);
}
.pill-safe {
  background: rgba(46, 188, 139, 0.1);
  border-color: rgba(46, 188, 139, 0.3);
  color: #159468;
}
.pill-watch {
  background: rgba(69, 176, 255, 0.1);
  border-color: rgba(69, 176, 255, 0.28);
  color: #217fb4;
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
  .docs-grid,
  .doc-filters,
  .docs-hero { grid-template-columns: 1fr; flex-direction: column; align-items: flex-start; }
}
</style>
