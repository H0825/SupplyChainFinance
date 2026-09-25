<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { request, requestPreferDB } from '../api/http'

const mode = ref('single')
const currentStep = ref(1)
const trackingStatus = ref('')
const loading = ref(false)

const form = reactive({ receivableId: '', amount: 0, interestRate: 420 })
const batchText = ref('')
const detailId = ref('')
const detail = ref(null)
const receivables = ref([])
const logs = ref([])
const screenQuery = reactive({ keyword: '', sort: 'amount_desc' })

const statusTextMap = {
  PENDING: '待确认',
  CONFIRMED: '已确认',
  FINANCED: '已融资',
  SETTLED: '已结清'
}

const availableReceivables = computed(() => receivables.value.filter((item) => item.statusText === 'CONFIRMED'))
const screenedReceivables = computed(() => {
  const keyword = screenQuery.keyword.trim().toLowerCase()
  const rows = availableReceivables.value.filter((item) => {
    if (!keyword) return true
    return [item.receivableId, item.orderId, item.payer]
      .some((value) => String(value || '').toLowerCase().includes(keyword))
  })
  const sorted = [...rows]
  sorted.sort((a, b) => {
    if (screenQuery.sort === 'amount_asc') return Number(a.amount || 0) - Number(b.amount || 0)
    if (screenQuery.sort === 'due_asc') return Number(a.dueTime || 0) - Number(b.dueTime || 0)
    return Number(b.amount || 0) - Number(a.amount || 0)
  })
  return sorted
})
const detailOptions = computed(() => receivables.value.slice(0, 12))
const selectedReceivable = computed(() => receivables.value.find((item) => item.receivableId === form.receivableId) || null)
const lastLog = computed(() => logs.value[0] || null)
const lastResult = computed(() => lastLog.value?.data || null)
const isEmptyWorkspace = computed(() => !loading.value && receivables.value.length === 0)

const financingStats = computed(() => ({
  available: availableReceivables.value.length,
  financed: receivables.value.filter((item) => item.statusText === 'FINANCED').length,
  settled: receivables.value.filter((item) => item.statusText === 'SETTLED').length,
  total: receivables.value.length
}))

const recommendedBatch = computed(() =>
  availableReceivables.value
    .slice(0, 3)
    .map((item) => `${item.receivableId},${Number(item.amount || 0)},420`)
    .join('\n')
)

const detailSummary = computed(() => {
  const payload = detail.value?.data
  if (!payload) return []
  return [
    { label: '应收编号', value: payload.receivableId || '-' },
    { label: '订单编号', value: payload.orderId || '-' },
    { label: '当前状态', value: statusLabel(payload.statusText) },
    { label: '应收金额', value: payload.amount ?? '-' }
  ]
})

const resultCards = computed(() => {
  if (!lastLog.value) return []
  const data = lastLog.value.data || {}
  return [
    { label: '最近操作', value: lastLog.value.title },
    { label: '处理状态', value: data.error ? '失败' : '成功' },
    { label: '应收编号', value: lastLog.value.businessId || detailId.value || form.receivableId || '-' },
    { label: '链上哈希', value: data.hash || '-' }
  ]
})

const resultMessage = computed(() => {
  const data = lastResult.value || {}
  if (!lastLog.value) return '执行融资操作后展示最近处理摘要'
  if (data.message) return data.message
  if (data.error?.error) return data.error.error
  if (data.error) return typeof data.error === 'string' ? data.error : '最近一次融资处理失败，请查看时间线'
  return '最近一次融资处理已完成'
})

const steps = [
  { no: 1, title: '授信筛选', desc: '优先从已确权应收中选择可融资资产' },
  { no: 2, title: '融资登记', desc: '登记融资金额与利率，并同步业务状态' },
  { no: 3, title: '详情复核', desc: '查看应收详情与融资记录，完成复核' }
]

function pushLog(title, data, businessId = '') {
  logs.value.unshift({ title, at: new Date().toLocaleString(), data, businessId })
  if (logs.value.length > 12) logs.value = logs.value.slice(0, 12)
}

function statusLabel(statusText) {
  return statusTextMap[statusText] || statusText || '未识别'
}

function applyStatusToStep(statusText) {
  trackingStatus.value = statusLabel(statusText)
  if (statusText === 'CONFIRMED') currentStep.value = 1
  else if (statusText === 'FINANCED') currentStep.value = 2
  else if (statusText === 'SETTLED') currentStep.value = 3
}

function applyReceivableToForm(receivableId) {
  form.receivableId = receivableId
  const selected = receivables.value.find((item) => item.receivableId === receivableId)
  if (!selected) return
  form.amount = Number(selected.amount || 0)
  detailId.value = receivableId
  applyStatusToStep(selected.statusText)
}

function moveToScreeningStep() {
  currentStep.value = 1
}

function moveToRegistrationStep() {
  if (!selectedReceivable.value) {
    ElMessage.info('请先在授信筛选中选择一条可融资应收')
    return
  }
  currentStep.value = 2
}

function fillBatchTemplate() {
  if (!recommendedBatch.value) {
    ElMessage.info('当前没有可直接批量融资的已确权应收')
    return
  }
  batchText.value = recommendedBatch.value
  ElMessage.success('已填入推荐批量模板')
}

async function loadReceivables() {
  loading.value = true
  try {
    const res = await requestPreferDB('/api/db/receivables', '/api/finance/receivables', { withRole: true })
    receivables.value = res.data || []
    if (!form.receivableId && availableReceivables.value.length) {
      applyReceivableToForm(availableReceivables.value[0].receivableId)
    }
  } finally {
    loading.value = false
  }
}

watch(
  () => form.receivableId,
  (id) => {
    if (!id) return
    applyReceivableToForm(id)
  }
)

async function financeOne(payload) {
  return request('/api/finance/receivable/finance', { method: 'POST', withRole: true, body: payload })
}

async function submitOne() {
  try {
    if (!form.receivableId) {
      ElMessage.info('请先从已确权应收中选择一条记录')
      return
    }
    const res = await financeOne({ receivableId: form.receivableId, amount: Number(form.amount), interestRate: Number(form.interestRate) })
    pushLog('单笔融资登记成功', res, form.receivableId)
    currentStep.value = 3
    mode.value = 'query'
    detailId.value = form.receivableId
    await loadReceivables()
    await queryDetail(false)
  } catch (e) {
    pushLog('单笔融资登记失败', e, form.receivableId)
  }
}

async function submitBatch() {
  const lines = batchText.value.split('\n').map((item) => item.trim()).filter(Boolean)
  if (!lines.length) {
    ElMessage.info('请先填入批量融资模板')
    return
  }
  for (const line of lines) {
    const [receivableId, amount, interestRate] = line.split(',').map((item) => item.trim())
    try {
      const res = await financeOne({ receivableId, amount: Number(amount), interestRate: Number(interestRate) })
      pushLog(`批量融资成功: ${receivableId}`, res, receivableId)
    } catch (e) {
      pushLog(`批量融资失败: ${receivableId}`, e, receivableId)
    }
  }
  currentStep.value = 3
  mode.value = 'query'
  await loadReceivables()
}

async function queryDetail(logOnSuccess = true) {
  if (!detailId.value.trim()) {
    ElMessage.info('请先从已有应收中选择一条记录')
    return
  }
  try {
    const id = detailId.value.trim()
    const resp = await requestPreferDB(`/api/db/receivables/${id}`, `/api/finance/receivable/${id}`, { withRole: true })
    if (!resp?.data) {
      detail.value = { message: '数据库暂时没有该应收记录，可能尚未同步。', financing: [] }
      if (logOnSuccess) pushLog(`详情查询: ${id}`, detail.value, id)
      return
    }
    detail.value = resp
    applyStatusToStep(detail.value?.data?.statusText || '')
    if (logOnSuccess) pushLog(`详情查询: ${id}`, detail.value, id)
  } catch (e) {
    detail.value = { message: '当前没有查询到对应数据' }
    pushLog('详情查询失败', e, detailId.value.trim())
  }
}

onMounted(loadReceivables)
</script>

<template>
  <section class="workspace-3col">
    <aside class="card panel">
      <h3>融资步骤</h3>
      <el-skeleton v-if="loading" :rows="6" animated />
      <div v-else class="steps">
        <div v-for="step in steps" :key="step.no" class="step-item" :class="{ active: currentStep === step.no, done: currentStep > step.no }">
          <div class="step-no">{{ step.no }}</div>
          <div>
            <div><strong>{{ step.title }}</strong></div>
            <div class="sub">{{ step.desc }}</div>
          </div>
        </div>
      </div>

      <div v-if="!loading" class="form-block finance-side-card">
        <h3>业务概览</h3>
        <div class="finance-mini-grid">
          <div><span>可融资应收</span><strong>{{ financingStats.available }}</strong></div>
          <div><span>已融资</span><strong>{{ financingStats.financed }}</strong></div>
          <div><span>已结清</span><strong>{{ financingStats.settled }}</strong></div>
          <div><span>总样本</span><strong>{{ financingStats.total }}</strong></div>
        </div>
        <el-alert type="info" :closable="false" show-icon title="金融端仅对已确权应收发起融资，避免把待确认资产直接送入授信流程。" />
        <div class="sub">当前识别状态: <span class="status-pill">{{ trackingStatus || '未识别' }}</span></div>
      </div>
    </aside>

    <section class="card panel">
      <div class="row between">
        <h2>融资处理工作台</h2>
        <div class="row wrap">
          <el-button size="small" plain @click="mode='single'; currentStep=1">单笔融资</el-button>
          <el-button size="small" plain @click="mode='batch'; currentStep=1">批量融资</el-button>
          <el-button size="small" plain @click="mode='query'; currentStep=3">详情复核</el-button>
        </div>
      </div>

      <div v-if="loading" class="finance-loading-stack">
        <el-skeleton animated :rows="4" />
        <el-skeleton animated :rows="8" />
      </div>

      <div v-else-if="isEmptyWorkspace" class="finance-empty-wrap">
        <el-empty description="当前没有可处理的融资数据">
          <template #image>
            <div class="empty-orbit empty-finance">
              <div class="empty-orbit-ring"></div>
              <div class="empty-orbit-core">¥</div>
            </div>
          </template>
          <p class="sub">先由企业端完成企业登记、应收创建和确权，金融端才会出现可融资资产。</p>
          <el-button type="primary" plain @click="mode='query'">先看详情复核</el-button>
        </el-empty>
      </div>

      <div v-else-if="mode==='single'" class="form-block">
        <div class="row between">
          <h3>{{ currentStep === 1 ? 'STEP 1 · 授信筛选' : 'STEP 2 · 融资登记' }}</h3>
          <div class="row wrap">
            <el-button type="primary" plain @click="moveToScreeningStep">授信筛选</el-button>
            <el-button type="primary" plain :disabled="!selectedReceivable" @click="moveToRegistrationStep">进入登记</el-button>
            <el-button v-if="currentStep === 2" type="primary" plain @click="selectedReceivable && applyReceivableToForm(selectedReceivable.receivableId)">刷新建议金额</el-button>
          </div>
        </div>
        <el-alert
          :type="currentStep === 1 ? 'info' : 'success'"
          :closable="false"
          show-icon
          :title="currentStep === 1 ? '先从已确权应收中筛选并锁定融资对象，再进入融资登记。' : '已锁定融资对象，系统会自动带入建议融资金额。'"
        />

        <div v-if="currentStep === 1" class="finance-screening-stack">
          <div class="row wrap">
            <input v-model="screenQuery.keyword" class="input" placeholder="搜索应收ID / 订单号 / 付款方" />
            <select v-model="screenQuery.sort" class="input mini">
              <option value="amount_desc">金额降序</option>
              <option value="amount_asc">金额升序</option>
              <option value="due_asc">到期时间升序</option>
            </select>
          </div>
          <div class="sub">先筛出候选应收，再点击“进入登记”继续下一步。</div>
          <div v-if="screenedReceivables.length" class="finance-screening-stack">
            <el-select
              v-model="form.receivableId"
              filterable
              clearable
              placeholder="请选择筛选后的可融资应收"
              style="width:100%;"
            >
              <el-option
                v-for="item in screenedReceivables"
                :key="item.receivableId"
                :label="`${item.receivableId} · ${item.orderId || '-'} · ${item.amount || 0}`"
                :value="item.receivableId"
              />
            </el-select>
            <div v-if="selectedReceivable" class="finance-focus-card">
              <div class="focus-head">
                <strong>{{ selectedReceivable.receivableId }}</strong>
                <span class="status-pill">{{ statusLabel(selectedReceivable.statusText) }}</span>
              </div>
              <div class="focus-grid">
                <div><span>订单编号</span><strong>{{ selectedReceivable.orderId || '-' }}</strong></div>
                <div><span>应收金额</span><strong>{{ selectedReceivable.amount || 0 }}</strong></div>
                <div><span>付款方</span><strong>{{ selectedReceivable.payer || '-' }}</strong></div>
                <div><span>到期时间</span><strong>{{ selectedReceivable.dueTime || '-' }}</strong></div>
              </div>
            </div>
          </div>
          <div v-else class="sub">当前筛选条件下没有可融资的已确权应收。</div>
        </div>

        <div v-if="currentStep === 2 && selectedReceivable" class="finance-focus-card">
          <div class="focus-head">
            <strong>{{ selectedReceivable.receivableId }}</strong>
            <span class="status-pill">{{ statusLabel(selectedReceivable.statusText) }}</span>
          </div>
          <div class="focus-grid">
            <div><span>订单编号</span><strong>{{ selectedReceivable.orderId || '-' }}</strong></div>
            <div><span>建议金额</span><strong>{{ selectedReceivable.amount || 0 }}</strong></div>
            <div><span>付款方</span><strong>{{ selectedReceivable.payer || '-' }}</strong></div>
            <div><span>到期时间</span><strong>{{ selectedReceivable.dueTime || '-' }}</strong></div>
          </div>
        </div>

        <el-form v-if="currentStep === 2" label-width="120px" class="form-el">
          <el-form-item label="可融资应收">
            <el-select v-model="form.receivableId" filterable clearable placeholder="从已筛选应收中确认" style="width:100%;">
              <el-option v-for="item in availableReceivables" :key="item.receivableId" :label="`${item.receivableId} · ${item.amount} · ${statusLabel(item.statusText)}`" :value="item.receivableId" />
            </el-select>
          </el-form-item>
          <el-form-item label="融资金额">
            <el-input-number v-model="form.amount" :min="1" :precision="0" style="width:100%;" />
          </el-form-item>
          <el-form-item label="利率(万分比)">
            <el-input-number v-model="form.interestRate" :min="0" :precision="0" style="width:100%;" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="submitOne">提交融资</el-button>
          </el-form-item>
        </el-form>

        <div v-if="currentStep === 2" class="row wrap">
          <el-tag v-for="item in availableReceivables.slice(0, 4)" :key="item.receivableId" effect="light" class="click-tag" @click="applyReceivableToForm(item.receivableId)">
            {{ item.receivableId }}
          </el-tag>
        </div>
        <el-empty v-else-if="!availableReceivables.length" description="当前没有可融资的已确权应收" />
      </div>

      <div v-if="mode==='batch'" class="form-block">
        <div class="row between">
          <h3>批量融资登记</h3>
          <el-button type="primary" plain @click="fillBatchTemplate">填入推荐模板</el-button>
        </div>
        <el-alert type="success" :closable="false" show-icon title="系统会优先给出前几条已确权应收的批量模板，减少手工拼接格式。" />
        <el-input v-model="batchText" type="textarea" :rows="8" placeholder="R-1001,50000,420&#10;R-1002,30000,390" />
        <div style="margin-top:8px;">
          <el-button type="primary" plain @click="submitBatch">批量执行</el-button>
        </div>
        <div class="row wrap" style="margin-top:10px;" v-if="availableReceivables.length">
          <el-tag v-for="item in availableReceivables.slice(0, 6)" :key="item.receivableId" effect="plain" class="click-tag" @click="batchText += `${batchText ? '\n' : ''}${item.receivableId},${Number(item.amount || 0)},420`">
            {{ item.receivableId }}
          </el-tag>
        </div>
        <el-empty v-else description="当前没有可用于批量融资的样本" />
      </div>

      <div v-if="mode==='query'" class="form-block">
        <h3>融资详情查询</h3>
        <el-alert type="info" :closable="false" show-icon title="详情查询优先走数据库，无数据时只展示友好提示，不会中断页面。" />
        <div class="row">
          <el-select v-model="detailId" filterable clearable placeholder="选择应收ID" style="width:100%;">
            <el-option v-for="item in detailOptions" :key="item.receivableId" :label="`${item.receivableId} · ${statusLabel(item.statusText)}`" :value="item.receivableId" />
          </el-select>
          <el-button type="primary" @click="queryDetail(true)">查询</el-button>
        </div>

        <div v-if="detailSummary.length" class="finance-summary-grid">
          <div v-for="item in detailSummary" :key="item.label" class="finance-summary-card">
            <span>{{ item.label }}</span>
            <strong>{{ item.value }}</strong>
          </div>
        </div>

        <div v-if="detail?.message" class="result-card" style="margin-top:10px;">
          <span>查询提示</span>
          <strong>{{ detail.message }}</strong>
        </div>
        <div v-else-if="detail?.financing?.length" class="result-grid" style="margin-top:10px;">
          <div v-for="(item, index) in detail.financing" :key="index" class="result-card">
            <span>融资记录 {{ index + 1 }}</span>
            <strong>{{ item.amount }} / 利率 {{ item.interestRate }}</strong>
          </div>
        </div>
        <el-empty v-else description="请选择一条应收查看详情" />
      </div>
    </section>

    <aside class="card panel">
      <h3>结果面板</h3>
      <div v-if="lastLog" class="result-grid">
        <div v-for="item in resultCards" :key="item.label" class="result-card">
          <span>{{ item.label }}</span>
          <strong>{{ item.value }}</strong>
        </div>
        <el-alert :type="lastResult?.error ? 'error' : 'success'" :closable="false" show-icon :title="resultMessage" />
      </div>
      <div v-else class="sub">执行融资操作后展示最近处理摘要</div>

      <h3 style="margin-top:10px;">处理时间线</h3>
      <div class="timeline" v-if="logs.length">
        <div class="tl-item" v-for="item in logs" :key="item.at + item.title">
          <div class="tl-dot"></div>
          <div>
            <div><strong>{{ item.title }}</strong></div>
            <div class="sub">{{ item.at }}</div>
          </div>
        </div>
      </div>
      <div v-else class="sub">暂无记录</div>
    </aside>
  </section>
</template>

<style scoped>
.finance-side-card {
  margin-top: 10px;
  display: grid;
  gap: 12px;
}

.finance-screening-stack {
  display: grid;
  gap: 12px;
  margin-top: 12px;
}

.finance-mini-grid,
.finance-summary-grid,
.focus-grid,
.result-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.finance-candidate-grid {
  display: grid;
  gap: 12px;
}

.finance-candidate-card {
  padding: 14px;
  border-radius: 16px;
  border: 1px solid #d9e8ff;
  background:
    radial-gradient(circle at right top, rgba(79, 197, 255, 0.14), transparent 28%),
    linear-gradient(145deg, rgba(244, 250, 255, 0.98), rgba(235, 244, 255, 0.92));
  text-align: left;
  cursor: pointer;
  transition: transform .18s ease, box-shadow .18s ease, border-color .18s ease;
}

.finance-candidate-card:hover {
  transform: translateY(-2px);
  border-color: #8ec0ff;
  box-shadow: 0 12px 24px rgba(58, 122, 233, 0.12);
}

.finance-candidate-card.active {
  border-color: #57a4ff;
  box-shadow: 0 0 0 1px rgba(73, 146, 255, 0.16), 0 12px 28px rgba(60, 127, 236, 0.14);
}

.finance-mini-grid > div,
.finance-summary-card,
.focus-grid > div,
.result-card {
  display: grid;
  gap: 4px;
  padding: 12px;
  border-radius: 14px;
  border: 1px solid #d8e8ff;
  background: linear-gradient(180deg, rgba(247, 251, 255, 0.96), rgba(238, 246, 255, 0.92));
}

.finance-mini-grid span,
.finance-summary-card span,
.focus-grid span,
.result-card span {
  color: #6a80a0;
  font-size: 12px;
}

.finance-mini-grid strong,
.finance-summary-card strong,
.focus-grid strong,
.result-card strong {
  color: #16386f;
  font-size: 18px;
  word-break: break-all;
}

.finance-focus-card {
  margin: 12px 0;
  padding: 14px;
  border-radius: 16px;
  border: 1px solid #d9e8ff;
  background:
    radial-gradient(circle at right top, rgba(79, 197, 255, 0.14), transparent 28%),
    linear-gradient(145deg, rgba(244, 250, 255, 0.98), rgba(235, 244, 255, 0.92));
}

.focus-head {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  align-items: center;
  margin-bottom: 10px;
}

.focus-head strong {
  font-size: 16px;
  color: #14366e;
}

.click-tag {
  cursor: pointer;
}

.finance-loading-stack {
  display: grid;
  gap: 12px;
}

.finance-empty-wrap {
  min-height: 420px;
  display: grid;
  place-items: center;
}

.empty-orbit {
  position: relative;
  width: 132px;
  height: 132px;
  display: grid;
  place-items: center;
}

.empty-orbit-ring {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  border: 1px solid rgba(76, 151, 255, 0.28);
  box-shadow: 0 0 0 12px rgba(103, 181, 255, 0.08), inset 0 0 24px rgba(57, 128, 239, 0.08);
}

.empty-orbit-core {
  width: 58px;
  height: 58px;
  border-radius: 18px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #2f72ef, #49c7ff);
  color: white;
  font-size: 28px;
  font-weight: 700;
  box-shadow: 0 16px 30px rgba(58, 122, 233, 0.22);
}

@media (max-width: 960px) {
  .finance-mini-grid,
  .finance-summary-grid,
  .focus-grid,
  .result-grid {
    grid-template-columns: 1fr;
  }
}
</style>
