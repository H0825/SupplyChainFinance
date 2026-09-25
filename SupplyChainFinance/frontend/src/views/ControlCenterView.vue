<script setup>
import { computed, onMounted, ref } from 'vue'
import { request, requestPreferDB } from '../api/http'

const users = ref([])
const enterprises = ref([])
const err = ref('')
const keyword = ref('')
const loading = ref(false)
const chainStats = ref({
  groupId: '-',
  chainId: '-',
  fiscoVersion: '-',
  blockHeight: 0,
  totalTxCount: '0',
  failedTxCount: '0',
  nodeCount: 0,
  nodeIds: [],
  warnings: []
})

const merged = computed(() => {
  const eMap = new Map(enterprises.value.map((e) => [String(e.wallet).toLowerCase(), e]))
  return users.value
    .map((u) => {
      const e = eMap.get(String(u.address).toLowerCase())
      return {
        ...u,
        enterpriseId: e?.enterpriseId || '-',
        enterpriseName: e?.name || '-'
      }
    })
    .filter((x) => {
      const key = keyword.value.trim().toLowerCase()
      if (!key) return true
      return [x.username, x.address, x.role, x.enterpriseId, x.enterpriseName].some((v) => String(v).toLowerCase().includes(key))
    })
})

const roleStats = computed(() => {
  const m = { business: 0, finance: 0, admin: 0 }
  users.value.forEach((u) => {
    const r = String(u.role || '').toLowerCase()
    if (m[r] !== undefined) m[r]++
  })
  return m
})

const summaryCards = computed(() => [
  { label: '企业用户', value: roleStats.value.business, hint: '企业与核心方账户' },
  { label: '金融用户', value: roleStats.value.finance, hint: '授信与资金方账户' },
  { label: '管理员', value: roleStats.value.admin, hint: '控制台管理账户' },
  { label: '已登记企业', value: enterprises.value.length, hint: '数据库企业主数据' }
])

const chainCards = computed(() => [
  { label: '区块高度', value: chainStats.value.blockHeight || 0 },
  { label: '累计交易', value: chainStats.value.totalTxCount || '0' },
  { label: '失败交易', value: chainStats.value.failedTxCount || '0' },
  { label: '节点数量', value: chainStats.value.nodeCount || 0 }
])

const nodeTicker = computed(() => {
  const list = chainStats.value.nodeIds || []
  if (!list.length) return []
  if (list.length === 1) return [list[0], list[0]]
  return list.slice(0, 2)
})

const isEmptyMappings = computed(() => !loading.value && merged.value.length === 0)

async function load() {
  loading.value = true
  try {
    const [u, e, s] = await Promise.all([
      request('/api/auth/users'),
      requestPreferDB('/api/db/enterprises', '/api/enterprise'),
      request('/api/chain/stats', { withRole: true })
    ])
    users.value = u.data || []
    enterprises.value = e.data || []
    chainStats.value = {
      ...chainStats.value,
      ...(s || {}),
      nodeIds: Array.isArray(s?.nodeIds) ? s.nodeIds : [],
      warnings: Array.isArray(s?.warnings) ? s.warnings : []
    }
    err.value = ''
  } catch (x) {
    err.value = JSON.stringify(x, null, 2)
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <section class="control-shell">
    <article class="card control-hero">
      <div>
        <span class="hero-badge">管理控制台</span>
        <h2>用户、企业映射与链上运行一体化视图</h2>
        <p>把角色分布、企业映射关系和区块链节点运行数据统一放到一个页面里，便于你快速定位系统状态。</p>
      </div>
      <button class="btn secondary" @click="load">刷新控制台</button>
    </article>

    <section v-if="loading" class="control-summary">
      <article v-for="i in 4" :key="i" class="summary-card summary-skeleton">
        <el-skeleton animated :rows="3" />
      </article>
    </section>

    <section v-else class="control-summary">
      <article v-for="item in summaryCards" :key="item.label" class="summary-card">
        <span>{{ item.label }}</span>
        <strong>{{ item.value }}</strong>
        <small>{{ item.hint }}</small>
      </article>
    </section>

    <section v-if="loading" class="control-grid">
      <article class="card"><el-skeleton animated :rows="10" /></article>
      <article class="card"><el-skeleton animated :rows="6" /></article>
    </section>

    <section v-else class="control-grid">
      <article class="card">
        <div class="panel-head">
          <div>
            <h3>区块链运行状态</h3>
            <p>直接查看链高度、交易量、节点规模和版本信息。</p>
          </div>
          <span class="panel-tag">FISCO BCOS</span>
        </div>

        <div class="chain-grid">
          <div v-for="item in chainCards" :key="item.label" class="chain-stat-card">
            <span>{{ item.label }}</span>
            <strong>{{ item.value }}</strong>
          </div>
        </div>

        <div class="row wrap" style="margin-bottom: 12px;">
          <span class="status-pill">群组 {{ chainStats.groupId || '-' }}</span>
          <span class="status-pill">链 ID {{ chainStats.chainId || '-' }}</span>
          <span class="status-pill">版本 {{ chainStats.fiscoVersion || '-' }}</span>
        </div>

        <div class="node-box-admin">
          <div class="row between">
            <strong>节点 ID</strong>
            <span class="sub">{{ chainStats.nodeIds.length }} 个节点</span>
          </div>
          <div v-if="nodeTicker.length" class="node-ticker-admin">
            <div class="node-track-admin">
              <code v-for="(id, index) in [...nodeTicker, ...nodeTicker]" :key="`${id}-${index}`">{{ id }}</code>
            </div>
          </div>
          <div v-else class="sub">暂无节点ID数据</div>
        </div>

        <el-alert
          v-for="warning in chainStats.warnings"
          :key="warning"
          type="warning"
          :title="warning"
          :closable="false"
          show-icon
          style="margin-top: 10px;"
        />
      </article>

      <article class="card">
        <div class="panel-head">
          <div>
            <h3>账户映射检索</h3>
            <p>按用户、地址、企业编号或企业名称直接筛选。</p>
          </div>
        </div>
        <el-input v-model="keyword" placeholder="搜索用户/地址/企业信息" clearable />
        <div class="mapping-tags row wrap" style="margin-top: 12px;">
          <span class="status-pill">企业端 {{ roleStats.business }}</span>
          <span class="status-pill">金融端 {{ roleStats.finance }}</span>
          <span class="status-pill">管理端 {{ roleStats.admin }}</span>
        </div>
        <div class="sub" style="margin-top: 12px;">当前共匹配 {{ merged.length }} 条映射记录</div>
      </article>
    </section>

    <article class="card">
      <div class="panel-head">
        <div>
          <h3>用户与企业映射表</h3>
          <p>核对链账户、企业身份与角色归属关系。</p>
        </div>
      </div>

      <div v-if="isEmptyMappings" class="empty-panel">
        <el-empty description="当前没有可展示的用户或企业映射数据">
          <template #image>
            <div class="empty-orbit empty-admin">
              <div class="empty-orbit-ring"></div>
              <div class="empty-orbit-core">ID</div>
            </div>
          </template>
          <p class="sub">先插入用户和企业基础数据，控制台会自动把映射关系组织出来。</p>
        </el-empty>
      </div>
      <div v-else class="table-wrap">
        <table class="table">
          <thead>
            <tr>
              <th>用户</th>
              <th>角色</th>
              <th>地址</th>
              <th>企业ID</th>
              <th>企业名</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="x in merged" :key="x.id">
              <td>{{ x.username }}</td>
              <td><span class="status-pill">{{ x.role }}</span></td>
              <td class="mono-cell">{{ x.address }}</td>
              <td>{{ x.enterpriseId }}</td>
              <td>{{ x.enterpriseName }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <pre v-if="err" class="log">{{ err }}</pre>
    </article>
  </section>
</template>

<style scoped>
.control-shell {
  display: grid;
  gap: 14px;
}

.control-hero {
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

.hero-badge {
  display: inline-flex;
  padding: 6px 10px;
  border-radius: 999px;
  border: 1px solid rgba(255,255,255,.22);
  background: rgba(255,255,255,.12);
  font-size: 12px;
  margin-bottom: 10px;
}

.control-hero h2 {
  margin: 0 0 8px;
  color: #fff;
}

.control-hero p {
  margin: 0;
  color: rgba(237, 245, 255, 0.84);
}

.control-summary {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 14px;
}

.summary-card,
.chain-stat-card {
  display: grid;
  gap: 6px;
  padding: 16px;
  border-radius: 16px;
  border: 1px solid #d9e8ff;
  background: linear-gradient(180deg, rgba(250, 252, 255, 0.96), rgba(239, 247, 255, 0.9));
}

.summary-card span,
.chain-stat-card span {
  color: #6e84a4;
  font-size: 12px;
}

.summary-card strong,
.chain-stat-card strong {
  color: #14366e;
  font-size: 28px;
}

.summary-card small {
  color: #8095b1;
}

.control-grid {
  display: grid;
  grid-template-columns: 1.2fr 0.8fr;
  gap: 14px;
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

.chain-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
  margin-bottom: 12px;
}

.node-box-admin {
  padding: 12px;
  border-radius: 14px;
  background: rgba(243, 249, 255, 0.92);
  border: 1px solid #dbe9ff;
}

.node-ticker-admin {
  position: relative;
  overflow: hidden;
  height: 96px;
  margin-top: 10px;
}

.node-track-admin {
  display: grid;
  gap: 8px;
  animation: node-scroll-admin 12s linear infinite;
}

.node-track-admin code {
  display: block;
  padding: 8px 10px;
  border-radius: 10px;
  background: #edf5ff;
  border: 1px solid #d6e6ff;
  color: #34557b;
  font-size: 12px;
  word-break: break-all;
}

.table-wrap {
  overflow: auto;
}

.mono-cell {
  word-break: break-all;
}

.summary-skeleton {
  min-height: 128px;
}

.empty-panel {
  min-height: 320px;
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
  font-size: 18px;
  font-weight: 700;
  box-shadow: 0 16px 30px rgba(58, 122, 233, 0.22);
}

@keyframes node-scroll-admin {
  0% { transform: translateY(0); }
  100% { transform: translateY(calc(-50% - 4px)); }
}

@media (max-width: 980px) {
  .control-hero,
  .control-grid {
    grid-template-columns: 1fr;
    flex-direction: column;
    align-items: flex-start;
  }

  .chain-grid {
    grid-template-columns: 1fr;
  }
}
</style>
