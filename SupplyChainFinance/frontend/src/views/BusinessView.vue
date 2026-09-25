<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'
import { request, requestPreferDB, uploadFiles } from '../api/http'

const users = ref([])
const enterprises = ref([])
const receivables = ref([])
const logs = ref([])
const currentStep = ref(1)
const trackingId = ref('')
const trackingStatus = ref('')

const enterprise = reactive({ enterpriseId: '', name: '', isCore: false, isFinance: false })
const issue = reactive({ receivableId: '', orderId: '', payer: '', amount: '', dueTimestamp: '' })
const confirmForm = reactive({ receivableId: '' })
const settleForm = reactive({ receivableId: '' })

const issueFiles = ref([])
const issueCids = ref([])
const docsUploading = ref(false)
const loading = ref(false)
const currentUser = JSON.parse(localStorage.getItem('scf_user') || 'null')

const pendingReceivables = computed(() => receivables.value.filter((item) => item.statusText === 'PENDING'))
const settleCandidates = computed(() => receivables.value.filter((item) => ['CONFIRMED', 'FINANCED'].includes(item.statusText)))
const latestReceivables = computed(() => receivables.value.slice(0, 8))
const lastResult = computed(() => logs.value[0]?.data || null)
const lastLog = computed(() => logs.value[0] || null)
const hasBusinessData = computed(() => users.value.length > 0 || receivables.value.length > 0)
const registeredEnterprise = computed(() => {
  const addr = String(currentUser?.address || '').toLowerCase()
  if (!addr) return null
  return enterprises.value.find((item) => String(item.wallet || '').toLowerCase() === addr) || null
})
const selectedEnterpriseRole = computed(() => (enterprise.isCore ? 'core' : 'supplier'))
const registeredEnterpriseRole = computed(() => {
  if (!registeredEnterprise.value) return ''
  return registeredEnterprise.value.isCore ? 'core' : 'supplier'
})
const effectiveEnterpriseRole = computed(() => registeredEnterpriseRole.value || selectedEnterpriseRole.value)
const canCreateReceivable = computed(() => effectiveEnterpriseRole.value === 'supplier')
const canConfirmReceivable = computed(() => effectiveEnterpriseRole.value === 'core')
const canSettleReceivable = computed(() => effectiveEnterpriseRole.value === 'core')
const resultCards = computed(() => {
  if (!lastLog.value) return []
  const data = lastLog.value.data || {}
  return [
    { label: '最近操作', value: lastLog.value.title },
    { label: '处理状态', value: data.error ? '失败' : data.alreadyRegistered ? '已登记' : '成功' },
    { label: '业务编号', value: lastLog.value.businessId || data?.data?.enterpriseId || trackingId.value || '-' },
    { label: '链上哈希', value: data.hash || '-' }
  ]
})
const resultMessage = computed(() => {
  const data = lastResult.value || {}
  if (!lastLog.value) return '执行操作后展示最近处理摘要'
  if (data.message) return data.message
  if (data.error?.error) return data.error.error
  if (data.error) return typeof data.error === 'string' ? data.error : '最近一次操作失败，请查看时间线'
  return '最近一次业务处理已完成'
})

const steps = [
  { no: 1, title: '企业注册', desc: '登记企业身份和属性' },
  { no: 2, title: '应收创建', desc: '供应商填写应收并上传凭证文档' },
  { no: 3, title: '应收确权', desc: '核心企业确认应收有效性' },
  { no: 4, title: '应收结算', desc: '核心企业完成链上结算归档' }
]

function getRoleLabel(role) {
  return role === 'core' ? '核心企业' : role === 'supplier' ? '供应商' : '未选择'
}

function syncStepByRole(role = effectiveEnterpriseRole.value) {
  if (role === 'core') currentStep.value = Math.max(currentStep.value, 3)
  else if (role === 'supplier') currentStep.value = Math.max(currentStep.value, 2)
}

function appendLog(title, data) {
  const businessId = data?.data?.enterpriseId || data?.receivableId || trackingId.value || ''
  logs.value.unshift({ title, at: new Date().toLocaleString(), data, businessId })
  if (logs.value.length > 12) logs.value = logs.value.slice(0, 12)
}

function generateEnterpriseId(showToast = false) {
  const now = new Date()
  const stamp = `${now.getFullYear()}${String(now.getMonth() + 1).padStart(2, '0')}${String(now.getDate()).padStart(2, '0')}${String(now.getHours()).padStart(2, '0')}${String(now.getMinutes()).padStart(2, '0')}${String(now.getSeconds()).padStart(2, '0')}`
  const suffix = String(Math.floor(Math.random() * 900) + 100)
  enterprise.enterpriseId = `ENT-${stamp}-${suffix}`
  if (showToast) ElMessage.success('企业ID已重新生成')
}

function autoFillIssue() {
  const now = Date.now()
  issue.receivableId = `R-${now}`
  issue.orderId = `O-${now}`
  issue.amount = 120000
  issue.dueTimestamp = Math.floor((now + 1000 * 3600 * 24 * 60) / 1000)
}

async function loadUsers() {
  const res = await request('/api/auth/users')
  users.value = res.data || []
  if (!issue.payer && users.value.length) issue.payer = users.value[0].address
}

async function loadReceivables() {
  const res = await requestPreferDB('/api/db/receivables', '/api/business/receivables', { withRole: true })
  receivables.value = res.data || []
}

async function loadEnterprises() {
  const res = await requestPreferDB('/api/db/enterprises', '/api/enterprise', { withRole: true })
  enterprises.value = res.data || []
}

function useReceivableForConfirm(id) {
  confirmForm.receivableId = id
  trackingId.value = id
  currentStep.value = 3
}

function useReceivableForSettle(id) {
  settleForm.receivableId = id
  trackingId.value = id
  currentStep.value = 4
}

function applyStatusToStep(statusText) {
  trackingStatus.value = statusText || ''
  if (statusText === 'PENDING') currentStep.value = 3
  else if (statusText === 'CONFIRMED') currentStep.value = 4
  else if (statusText === 'FINANCED' || statusText === 'SETTLED') currentStep.value = 4
}

async function syncStepByReceivableID(id) {
  const rid = (id || trackingId.value).trim()
  if (!rid) {
    ElMessage.info('请选择已有应收，或输入应收ID后再同步')
    return
  }
  try {
    const detail = await requestPreferDB(`/api/db/receivables/${rid}`, `/api/business/receivable/${rid}`, { withRole: true })
    const statusText = detail?.data?.statusText || ''
    applyStatusToStep(statusText)
    confirmForm.receivableId = rid
    settleForm.receivableId = rid
    trackingId.value = rid
    appendLog(`状态同步: ${rid}`, { statusText })
  } catch (e) {
    appendLog('状态同步失败', e)
  }
}

function onUploadChange(file, fileList) {
  issueFiles.value = fileList
}

async function uploadIssueDocs() {
  const raws = issueFiles.value.map((item) => item.raw).filter(Boolean)
  if (!raws.length) {
    ElMessage.warning('请先选择凭证文件')
    return
  }
  docsUploading.value = true
  try {
    const resp = await uploadFiles('/api/ipfs/upload', raws, true, {
      tag: issue.receivableId || issue.orderId || 'receivable-doc',
      operator: currentUser?.username || currentUser?.address || '',
      receivableId: issue.receivableId || '',
      orderId: issue.orderId || '',
      docType: 'contract',
      status: 'uploaded',
      source: 'ipfs'
    })
    issueCids.value = resp.cids || []
    ElMessage.success(`凭证上传成功，共 ${issueCids.value.length} 个CID`)
  } catch (e) {
    appendLog('凭证上传失败', e)
    ElMessage.error('凭证上传失败，请检查 IPFS 服务')
  } finally {
    docsUploading.value = false
  }
}

async function registerEnterprise() {
  try {
    if (registeredEnterprise.value) {
      appendLog('企业已登记', { alreadyRegistered: true, data: registeredEnterprise.value, message: '当前账户已完成企业登记，无需重复注册' })
      syncStepByRole(registeredEnterpriseRole.value)
      return
    }
    if (!enterprise.enterpriseId) generateEnterpriseId()
    const res = await request('/api/business/enterprise/register', { method: 'POST', body: enterprise, withRole: true })
    appendLog('企业注册成功', res)
    await loadEnterprises()
    syncStepByRole(selectedEnterpriseRole.value)
  } catch (e) {
    appendLog('企业注册失败', e)
  }
}

async function issueReceivable() {
  try {
    if (!canCreateReceivable.value) {
      ElMessage.warning('当前角色为核心企业，仅供应商可以创建应收')
      return
    }
    if (!issue.receivableId || !issue.orderId || !issue.payer) {
      ElMessage.warning('请完整填写应收字段')
      return
    }
    const res = await request('/api/business/receivable/issue', {
      method: 'POST',
      withRole: true,
      body: {
        receivableId: issue.receivableId,
        orderId: issue.orderId,
        payer: issue.payer,
        amount: Number(issue.amount),
        dueTimestamp: Number(issue.dueTimestamp)
      }
    })
    appendLog('应收创建成功', { ...res, docCids: issueCids.value })
    trackingId.value = issue.receivableId
    currentStep.value = Math.max(currentStep.value, 3)
    await loadReceivables()
    await syncStepByReceivableID(issue.receivableId)
  } catch (e) {
    appendLog('应收创建失败', e)
  }
}

async function confirmReceivable() {
  try {
    if (!canConfirmReceivable.value) {
      ElMessage.warning('当前角色为供应商，仅核心企业可以进行应收确权')
      return
    }
    if (!confirmForm.receivableId) {
      ElMessage.info('请从下方待确认列表中选择一条应收')
      return
    }
    const res = await request('/api/business/receivable/confirm', { method: 'POST', withRole: true, body: confirmForm })
    appendLog('确权成功', res)
    await loadReceivables()
    await syncStepByReceivableID(confirmForm.receivableId)
  } catch (e) {
    appendLog('确权失败', e)
  }
}

async function settleReceivable() {
  try {
    if (!canSettleReceivable.value) {
      ElMessage.warning('当前角色为供应商，仅核心企业可以进行应收结算')
      return
    }
    if (!settleForm.receivableId) {
      ElMessage.info('请从下方可结算列表中选择一条应收')
      return
    }
    const res = await request('/api/business/receivable/settle', { method: 'POST', withRole: true, body: settleForm })
    appendLog('结算成功', res)
    await loadReceivables()
    await syncStepByReceivableID(settleForm.receivableId)
  } catch (e) {
    appendLog('结算失败', e)
  }
}

onMounted(async () => {
  generateEnterpriseId()
  autoFillIssue()
  loading.value = true
  try {
    await Promise.all([loadUsers(), loadReceivables(), loadEnterprises()])
    if (registeredEnterprise.value) syncStepByRole(registeredEnterpriseRole.value)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <section class="workspace-3col">
    <aside class="card panel">
      <h3>流程步骤</h3>
      <div class="steps">
        <div v-for="step in steps" :key="step.no" class="step-item" :class="{ active: currentStep === step.no, done: currentStep > step.no }">
          <div class="step-no">{{ step.no }}</div>
          <div>
            <div><strong>{{ step.title }}</strong></div>
            <div class="sub">{{ step.desc }}</div>
          </div>
        </div>
      </div>

      <div v-if="!loading" class="form-block" style="margin-top:10px; display:grid; gap:8px;">
        <h3>状态联动</h3>
        <el-alert type="info" :closable="false" show-icon title="可以直接从下方最近应收中点选，不必手动记忆ID。" />
        <el-select v-model="trackingId" filterable clearable placeholder="选择已有应收同步步骤" style="width:100%;">
          <el-option v-for="item in latestReceivables" :key="item.receivableId" :label="`${item.receivableId} · ${item.statusText}`" :value="item.receivableId" />
        </el-select>
        <el-button type="primary" @click="syncStepByReceivableID('')">同步状态</el-button>
        <div class="sub">当前识别状态: <span class="status-pill">{{ trackingStatus || '未识别' }}</span></div>
      </div>
    </aside>

    <section class="card panel">
      <div class="row between">
        <h2>企业业务工作台</h2>
        <div class="row wrap">
          <el-button size="small" plain @click="currentStep=1">企业注册</el-button>
          <el-button size="small" plain :disabled="!canCreateReceivable" @click="currentStep=2">应收创建</el-button>
          <el-button size="small" plain :disabled="!canConfirmReceivable" @click="currentStep=3">应收确权</el-button>
          <el-button size="small" plain :disabled="!canSettleReceivable" @click="currentStep=4">应收结算</el-button>
        </div>
      </div>

      <div v-if="loading" class="business-loading-stack">
        <el-skeleton animated :rows="4" />
        <el-skeleton animated :rows="9" />
      </div>

      <div v-else-if="!hasBusinessData && currentStep!==1" class="business-empty-wrap">
        <el-empty description="No enterprise workflow data yet">
          <template #image>
            <div class="business-empty-orbit">
              <div class="business-empty-ring"></div>
              <div class="business-empty-core">SCF</div>
            </div>
          </template>
          <p class="sub">Complete enterprise registration and receivable creation first. Confirmation, settlement, and status sync candidates will appear automatically.</p>
          <el-button type="primary" plain @click="currentStep=1">Go to Enterprise Registration</el-button>
        </el-empty>
      </div>

      <div v-else-if="currentStep===1" class="form-block">
        <h3>STEP 1 · 企业注册</h3>
        <el-alert v-if="!registeredEnterprise" type="info" :closable="false" show-icon title="系统已自动生成企业ID，你只需补充企业名称并选择企业属性。" />
        <div v-if="registeredEnterprise" class="result-grid">
          <div class="result-card">
            <span>企业编号</span>
            <strong>{{ registeredEnterprise.enterpriseId }}</strong>
          </div>
          <div class="result-card">
            <span>企业名称</span>
            <strong>{{ registeredEnterprise.name }}</strong>
          </div>
          <div class="result-card">
            <span>企业角色</span>
            <strong>{{ getRoleLabel(registeredEnterpriseRole) }}</strong>
          </div>
          <div class="result-card">
            <span>可执行操作</span>
            <strong>{{ registeredEnterpriseRole === 'supplier' ? '应收创建' : '应收确权、应收结算' }}</strong>
          </div>
          <el-alert :type="registeredEnterpriseRole === 'supplier' ? 'success' : 'info'" :closable="false" show-icon :title="registeredEnterpriseRole === 'supplier' ? '当前账户已完成企业登记，可直接进入应收创建。' : '当前账户已完成企业登记，可直接进入应收确权与结算。'" />
        </div>
        <el-form v-else label-width="110px" class="form-el">
          <el-form-item label="企业ID">
            <div class="auto-id-row">
              <el-input v-model="enterprise.enterpriseId" readonly placeholder="系统自动生成企业ID" />
              <el-button type="primary" plain @click="generateEnterpriseId(true)">重新生成</el-button>
            </div>
          </el-form-item>
          <el-form-item label="企业名称">
            <el-input v-model="enterprise.name" placeholder="输入企业全称" />
          </el-form-item>
          <el-form-item label="企业属性">
            <div class="property-grid">
              <button type="button" class="property-card" :class="{ active: selectedEnterpriseRole === 'core' }" @click="enterprise.isCore = true; enterprise.isFinance = false">
                <span class="property-badge">A</span>
                <div>
                  <strong>核心企业</strong>
                  <span>可进行应收确权与应收结算</span>
                </div>
              </button>
              <button type="button" class="property-card" :class="{ active: selectedEnterpriseRole === 'supplier' }" @click="enterprise.isCore = false; enterprise.isFinance = false">
                <span class="property-badge">S</span>
                <div>
                  <strong>供应商</strong>
                  <span>可进行应收创建与凭证上传</span>
                </div>
              </button>
            </div>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="registerEnterprise">提交注册</el-button>
          </el-form-item>
        </el-form>
      </div>

      <div v-if="currentStep===2" class="form-block">
        <div class="row between">
          <h3>STEP 2 · 应收创建（含凭证）</h3>
          <el-button type="primary" plain @click="autoFillIssue">自动填充</el-button>
        </div>
        <el-alert v-if="!canCreateReceivable" type="warning" :closable="false" show-icon title="当前角色无权创建应收。请选择“供应商”角色后再操作。" />
        <el-alert type="success" :closable="false" show-icon title="系统会自动生成应收ID和订单ID。只需选付款方、填金额和到期时间。" />

        <el-form label-width="110px" class="form-el">
          <el-form-item label="应收ID">
            <el-input v-model="issue.receivableId" placeholder="例如 R-20260309001" />
          </el-form-item>
          <el-form-item label="订单ID">
            <el-input v-model="issue.orderId" placeholder="例如 O-20260309001" />
          </el-form-item>
          <el-form-item label="付款方">
            <el-select v-model="issue.payer" placeholder="选择付款方地址" style="width:100%;">
              <el-option v-for="user in users" :key="user.id" :label="`${user.username} | ${user.address}`" :value="user.address" />
            </el-select>
          </el-form-item>
          <el-form-item label="应收金额">
            <el-input-number v-model="issue.amount" :min="1" :precision="0" style="width:100%;" />
          </el-form-item>
          <el-form-item label="到期时间戳">
            <el-input-number v-model="issue.dueTimestamp" :min="1" :precision="0" style="width:100%;" />
          </el-form-item>

          <el-form-item label="凭证文档">
            <div style="width:100%; display:grid; gap:8px;">
              <el-upload drag multiple :auto-upload="false" :on-change="onUploadChange" :show-file-list="true" :limit="10">
                <el-icon><upload-filled /></el-icon>
                <div>拖拽文件到此或点击选择（合同、发票、签收单等）</div>
              </el-upload>
              <div class="row wrap">
                <el-button :loading="docsUploading" type="primary" plain @click="uploadIssueDocs">上传凭证到IPFS</el-button>
                <el-tag v-for="cid in issueCids" :key="cid" size="small" effect="plain">{{ cid }}</el-tag>
              </div>
              <el-alert type="warning" :closable="false" show-icon title="建议将 CID 与 receivable_id 一并插入 documents 表，后续查询会更顺畅。" />
            </div>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" :disabled="!canCreateReceivable" @click="issueReceivable">创建应收</el-button>
          </el-form-item>
        </el-form>
      </div>

      <div v-if="currentStep===3" class="form-block" style="display:grid; gap:10px;">
        <h3>STEP 3 · 应收确权</h3>
        <el-alert v-if="!canConfirmReceivable" type="warning" :closable="false" show-icon title="当前角色无权进行应收确权。请选择“核心企业”角色后再操作。" />
        <el-alert type="info" :closable="false" show-icon title="不用记ID，直接从待确认应收里选择即可。" />
        <el-form label-width="110px" class="form-el">
          <el-form-item label="待确认应收">
            <el-select v-model="confirmForm.receivableId" filterable clearable placeholder="从待确认列表中选择" style="width:100%;">
              <el-option v-for="item in pendingReceivables" :key="item.receivableId" :label="`${item.receivableId} · ${item.orderId} · ${item.amount}`" :value="item.receivableId" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="warning" :disabled="!canConfirmReceivable" @click="confirmReceivable">确认应收</el-button>
          </el-form-item>
        </el-form>
        <div class="row wrap" v-if="pendingReceivables.length">
          <el-tag v-for="item in pendingReceivables.slice(0, 4)" :key="item.receivableId" effect="light" class="click-tag" @click="useReceivableForConfirm(item.receivableId)">
            {{ item.receivableId }}
          </el-tag>
        </div>
        <div v-else class="sub">当前没有待确认的应收，先去“应收创建”新增一笔。</div>
      </div>

      <div v-if="currentStep===4" class="form-block" style="display:grid; gap:10px;">
        <h3>STEP 4 · 应收结算</h3>
        <el-alert v-if="!canSettleReceivable" type="warning" :closable="false" show-icon title="当前角色无权进行应收结算。请选择“核心企业”角色后再操作。" />
        <el-alert type="info" :closable="false" show-icon title="可结算列表会自动筛出已确认或已融资的应收。" />
        <el-form label-width="110px" class="form-el">
          <el-form-item label="可结算应收">
            <el-select v-model="settleForm.receivableId" filterable clearable placeholder="从可结算列表中选择" style="width:100%;">
              <el-option v-for="item in settleCandidates" :key="item.receivableId" :label="`${item.receivableId} · ${item.statusText} · ${item.amount}`" :value="item.receivableId" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="danger" :disabled="!canSettleReceivable" @click="settleReceivable">登记结算</el-button>
          </el-form-item>
        </el-form>
        <div class="row wrap" v-if="settleCandidates.length">
          <el-tag v-for="item in settleCandidates.slice(0, 4)" :key="item.receivableId" effect="light" class="click-tag" @click="useReceivableForSettle(item.receivableId)">
            {{ item.receivableId }}
          </el-tag>
        </div>
        <div v-else class="sub">当前没有可结算的应收，先完成确权或融资流程。</div>
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
        <div v-if="lastResult?.docCids?.length" class="result-list">
          <span class="sub">凭证 CID</span>
          <el-tag v-for="cid in lastResult.docCids" :key="cid" effect="plain">{{ cid }}</el-tag>
        </div>
      </div>
      <div v-else class="sub">执行操作后展示最近处理摘要</div>

      <h3 style="margin-top: 10px;">时间线</h3>
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
.click-tag {
  cursor: pointer;
  transition: transform .18s ease, box-shadow .18s ease;
}

.click-tag:hover {
  transform: translateY(-1px);
  box-shadow: 0 8px 16px rgba(63, 126, 234, 0.12);
}

.auto-id-row {
  width: 100%;
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 10px;
}

.property-grid {
  width: 100%;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.property-card {
  display: grid;
  grid-template-columns: 42px 1fr;
  gap: 12px;
  align-items: center;
  padding: 14px;
  border-radius: 14px;
  border: 1px solid #d6e6ff;
  background: linear-gradient(180deg, rgba(245, 250, 255, 0.96), rgba(236, 245, 255, 0.9));
  cursor: pointer;
  transition: transform .2s ease, border-color .2s ease, box-shadow .2s ease;
  text-align: left;
}

.property-card:hover {
  transform: translateY(-2px);
  border-color: #8ec0ff;
  box-shadow: 0 12px 24px rgba(58, 122, 233, 0.12);
}

.property-card.active {
  border-color: #57a4ff;
  background: linear-gradient(135deg, rgba(227, 241, 255, 0.98), rgba(213, 235, 255, 0.92));
  box-shadow: 0 0 0 1px rgba(73, 146, 255, 0.16), 0 12px 28px rgba(60, 127, 236, 0.14);
}

.property-card.active .business-loading-stack {
  display: grid;
  gap: 12px;
}

.business-empty-wrap {
  min-height: 420px;
  display: grid;
  place-items: center;
}

.business-empty-orbit {
  position: relative;
  width: 132px;
  height: 132px;
  display: grid;
  place-items: center;
}

.business-empty-ring {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  border: 1px solid rgba(76, 151, 255, 0.28);
  box-shadow: 0 0 0 12px rgba(103, 181, 255, 0.08), inset 0 0 24px rgba(57, 128, 239, 0.08);
}

.business-empty-core {
  width: 62px;
  height: 62px;
  border-radius: 18px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #2f72ef, #49c7ff);
  color: white;
  font-size: 16px;
  font-weight: 700;
  box-shadow: 0 16px 30px rgba(58, 122, 233, 0.22);
}

.property-badge {
  transform: scale(1.05);
  box-shadow: 0 10px 18px rgba(52, 129, 245, 0.22);
}

.property-card strong {
  display: block;
  margin-bottom: 4px;
  color: #173c7b;
}

.property-card span {
  color: #6a7f9a;
  font-size: 12px;
  line-height: 1.5;
}

.property-badge {
  width: 42px;
  height: 42px;
  border-radius: 12px;
  background: linear-gradient(135deg, #2f72ef, #49c7ff);
  color: white;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 16px;
  transition: transform .2s ease, box-shadow .2s ease;
}

.result-grid {
  display: grid;
  gap: 10px;
}

.result-card {
  display: grid;
  gap: 4px;
  padding: 12px;
  border-radius: 14px;
  border: 1px solid #d8e8ff;
  background: linear-gradient(180deg, rgba(247, 251, 255, 0.96), rgba(238, 246, 255, 0.92));
}

.result-card span {
  color: #6a80a0;
  font-size: 12px;
}

.result-card strong {
  color: #16386f;
  font-size: 16px;
  word-break: break-all;
}

.result-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

@media (max-width: 960px) {
  .auto-id-row,
  .property-grid {
    grid-template-columns: 1fr;
  }
}
</style>
