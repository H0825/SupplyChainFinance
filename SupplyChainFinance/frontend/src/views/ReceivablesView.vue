<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { request, requestPreferDB } from '../api/http'

const list = ref([])
const detail = ref(null)
const err = ref('')
const loading = ref(false)
const enterprises = ref([])
const query = reactive({ id: '', status: 'ALL', keyword: '', sort: 'amount_desc' })

const currentUser = JSON.parse(localStorage.getItem('scf_user') || 'null')
const role = (currentUser?.role || 'guest').toLowerCase()

const statusLabelMap = {
  PENDING: '待确认',
  CONFIRMED: '已确认',
  FINANCED: '已融资',
  SETTLED: '已结清'
}

const filtered = computed(() => {
  const rows = list.value.filter((item) => {
    const statusOk = query.status === 'ALL' ? true : item.statusText === query.status
    const key = query.keyword.trim().toLowerCase()
    const keyOk = !key
      ? true
      : [item.receivableId, item.orderId, item.issuer, item.payer].some((value) => String(value || '').toLowerCase().includes(key))
    return statusOk && keyOk
  })

  const sorted = [...rows]
  sorted.sort((a, b) => {
    if (query.sort === 'amount_desc') return Number(b.amount) - Number(a.amount)
    if (query.sort === 'amount_asc') return Number(a.amount) - Number(b.amount)
    return String(a.receivableId).localeCompare(String(b.receivableId))
  })
  return sorted
})

const queryOptions = computed(() => list.value.slice(0, 20))
const registeredEnterprise = computed(() => {
  const addr = String(currentUser?.address || '').toLowerCase()
  if (!addr) return null
  return enterprises.value.find((item) => String(item.wallet || '').toLowerCase() === addr) || null
})
const enterpriseRole = computed(() => {
  if (!registeredEnterprise.value) return ''
  return registeredEnterprise.value.isCore ? 'core' : 'supplier'
})
const canQuickConfirm = computed(() => ['admin'].includes(role) || (role === 'business' && enterpriseRole.value === 'core'))
const canQuickSettle = computed(() => ['admin'].includes(role) || (role === 'business' && enterpriseRole.value === 'core'))
const roleHint = computed(() => {
  if (role === 'admin') return '管理员可查看并处理全部应收。'
  if (role !== 'business') return '当前账号仅可查看应收详情。'
  if (!registeredEnterprise.value) return '当前账号尚未完成企业登记，暂不能执行确权或结算。'
  return enterpriseRole.value === 'core'
    ? '当前账号角色为核心企业，可在这里快捷进行确权与结算。'
    : '当前账号角色为供应商，可查看应收详情；确权与结算需由核心企业处理。'
})
const detailCards = computed(() => {
  const payload = detail.value?.data
  if (!payload) return []
  return [
    { label: '应收编号', value: payload.receivableId || '-' },
    { label: '订单编号', value: payload.orderId || '-' },
    { label: '当前状态', value: statusLabel(payload.statusText) },
    { label: '应收金额', value: payload.amount || '-' },
    { label: '发行方', value: payload.issuer || '-' },
    { label: '付款方', value: payload.payer || '-' }
  ]
})

function statusLabel(status) {
  return statusLabelMap[status] || status
}

async function load() {
  loading.value = true
  try {
    const [receivableRes, enterpriseRes] = await Promise.all([
      requestPreferDB('/api/db/receivables', '/api/receivable'),
      requestPreferDB('/api/db/enterprises', '/api/enterprise', { withRole: true }).catch(() => ({ data: [] }))
    ])
    list.value = receivableRes.data || []
    enterprises.value = enterpriseRes.data || []
    err.value = ''
  } catch (e) {
    err.value = JSON.stringify(e, null, 2)
  } finally {
    loading.value = false
  }
}

async function queryDetail(id) {
  const rid = (id || query.id).trim()
  if (!rid) {
    ElMessage.info('请先从列表中选择一条应收记录')
    return
  }
  try {
    const resp = await requestPreferDB(`/api/db/receivables/${rid}`, `/api/receivable/${rid}`)
    if (!resp?.data) {
      detail.value = { message: '当前数据库没有该应收记录，可能尚未同步。', financing: [] }
      return
    }
    detail.value = resp
  } catch (e) {
    detail.value = { message: '当前没有查到这条应收', financing: [] }
  }
}

async function quickConfirm(id) {
  try {
    if (!canQuickConfirm.value) {
      ElMessage.warning(role === 'business' ? '仅核心企业可以进行应收确权' : '当前角色无权执行应收确权')
      return
    }
    await request('/api/receivable/confirm', { method: 'POST', withRole: true, body: { receivableId: id } })
    await load()
    await queryDetail(id)
  } catch (e) {
    err.value = JSON.stringify(e, null, 2)
  }
}

async function quickSettle(id) {
  try {
    if (!canQuickSettle.value) {
      ElMessage.warning(role === 'business' ? '仅核心企业可以进行应收结算' : '当前角色无权执行应收结算')
      return
    }
    await request('/api/receivable/settle', { method: 'POST', withRole: true, body: { receivableId: id } })
    await load()
    await queryDetail(id)
  } catch (e) {
    err.value = JSON.stringify(e, null, 2)
  }
}

onMounted(load)
</script>

<template>
  <section class="receivable-shell">
    <article class="card">
      <div class="row between">
        <h2>应收账款池</h2>
        <button class="btn secondary" @click="load">{{ loading ? '加载中...' : '刷新' }}</button>
      </div>

      <el-alert type="info" :closable="false" show-icon title="先筛选状态，再从候选列表中选择应收查看详情，不需要手动记ID。" style="margin-bottom:10px;" />
      <el-alert :title="roleHint" :type="canQuickConfirm || canQuickSettle ? 'success' : 'warning'" :closable="false" show-icon style="margin-bottom:10px;" />

      <div class="row wrap">
        <el-select v-model="query.id" filterable clearable placeholder="选择应收ID查询详情" style="width:260px;">
          <el-option v-for="item in queryOptions" :key="item.receivableId" :label="`${item.receivableId} · ${statusLabel(item.statusText)}`" :value="item.receivableId" />
        </el-select>
        <button class="btn mini" @click="queryDetail()">查询详情</button>
        <select v-model="query.status" class="input mini">
          <option value="ALL">全部状态</option>
          <option value="PENDING">待确认</option>
          <option value="CONFIRMED">已确认</option>
          <option value="FINANCED">已融资</option>
          <option value="SETTLED">已结清</option>
        </select>
        <select v-model="query.sort" class="input mini">
          <option value="amount_desc">金额降序</option>
          <option value="amount_asc">金额升序</option>
          <option value="id_asc">ID升序</option>
        </select>
        <input v-model="query.keyword" class="input" placeholder="搜索订单号 / 发行方 / 付款方" />
      </div>

      <table class="table">
        <thead>
          <tr>
            <th>ID</th><th>订单</th><th>发行方</th><th>付款方</th><th>金额</th><th>状态</th><th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in filtered" :key="item.receivableId">
            <td>{{ item.receivableId }}</td>
            <td>{{ item.orderId }}</td>
            <td>{{ item.issuer }}</td>
            <td>{{ item.payer }}</td>
            <td>{{ item.amount }}</td>
            <td><span class="status-pill">{{ statusLabel(item.statusText) }}</span></td>
            <td>
              <div class="row wrap">
                <button class="btn mini secondary" @click="queryDetail(item.receivableId)">详情</button>
                <button v-if="canQuickConfirm && item.statusText === 'PENDING'" class="btn mini warn" @click="quickConfirm(item.receivableId)">确权</button>
                <button v-if="canQuickSettle && ['CONFIRMED','FINANCED'].includes(item.statusText)" class="btn mini danger" @click="quickSettle(item.receivableId)">结算</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </article>

    <article class="card">
      <div class="panel-head compact">
        <h3>详情面板</h3>
      </div>
      <div v-if="detailCards.length" class="detail-grid">
        <div v-for="item in detailCards" :key="item.label" class="detail-card">
          <span>{{ item.label }}</span>
          <strong>{{ item.value }}</strong>
        </div>
      </div>
      <el-alert v-if="detail?.message" :title="detail.message" type="info" :closable="false" show-icon style="margin-top:10px;" />
      <div v-if="detail?.financing?.length" class="finance-records">
        <h4>融资记录</h4>
        <div v-for="(item, index) in detail.financing" :key="index" class="detail-card">
          <span>记录 {{ index + 1 }}</span>
          <strong>{{ item.amount }} / 利率 {{ item.interestRate }} / {{ item.financingTime }}</strong>
        </div>
      </div>
      <el-empty v-else-if="!detail" description="请选择一条应收查看详情" />
      <pre v-if="err" class="log">{{ err }}</pre>
    </article>
  </section>
</template>

<style scoped>
.receivable-shell {
  display: grid;
  gap: 14px;
}

.detail-grid,
.finance-records {
  display: grid;
  gap: 10px;
}

.detail-card {
  display: grid;
  gap: 4px;
  padding: 12px;
  border-radius: 14px;
  border: 1px solid #d8e8ff;
  background: linear-gradient(180deg, rgba(247, 251, 255, 0.96), rgba(238, 246, 255, 0.92));
}

.detail-card span {
  color: #6a80a0;
  font-size: 12px;
}

.detail-card strong {
  color: #16386f;
  font-size: 16px;
  word-break: break-all;
}
</style>
