<script setup>
import { computed, nextTick, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import DOMPurify from 'dompurify'
import { marked } from 'marked'
import { requestPreferDB, request } from '../api/http'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const input = ref('')
const sending = ref(false)
const chatBox = ref(null)

const guard = ref({
  privacyMode: true,
  localMasking: true,
  noRetention: true,
  minimalDisclosure: true
})

const messages = ref([
  {
    role: 'assistant',
    content: '你好，我是供应链金融对话助手。可协助你进行应收融资分析、流程诊断、风控提示与合规建议。',
    time: nowTime(),
    meta: {
      model: 'system'
    }
  }
])

marked.setOptions({
  gfm: true,
  breaks: true
})

function nowTime() {
  return new Date().toLocaleTimeString()
}

const roleLabel = computed(() => {
  if (auth.role === 'business') return '企业端模式'
  if (auth.role === 'finance') return '金融机构模式'
  return '管理模式'
})

const currentModel = computed(() => {
  const latestAssistant = [...messages.value].reverse().find((message) => message.role === 'assistant' && message.meta?.model)
  return latestAssistant?.meta?.model || 'doubao'
})

const quickPrompts = computed(() => {
  if (auth.role === 'business') {
    return [
      '请给我一个应收创建到结算的标准流程检查清单',
      '如何提高应收确权成功率？',
      '凭证文档应包含哪些字段才便于融资审核？',
      '帮我写一段给金融机构的融资申请说明'
    ]
  }
  return [
    '基于当前应收池，给出融资审查优先级建议',
    '给我一份贷前数据核验清单（链上+数据库）',
    '对可能的逾期风险给出早期预警信号',
    '如何在不暴露隐私数据前提下完成风控建模？'
  ]
})

function pushMessage(role, content, meta = null) {
  messages.value.push({ role, content, time: nowTime(), meta })
}

function renderMarkdown(content) {
  return DOMPurify.sanitize(marked.parse(String(content || '')))
}

function applyMasking(text) {
  let result = text
  result = result.replace(/0x[a-fA-F0-9]{40}/g, (match) => `${match.slice(0, 6)}****${match.slice(-4)}`)
  result = result.replace(/[A-Z]-\d{8,}/g, (match) => `${match.slice(0, 3)}****`)
  result = result.replace(/\b\d{11,18}\b/g, (match) => `${match.slice(0, 3)}****${match.slice(-2)}`)
  return result
}

function usageText(usage) {
  if (!usage) return ''
  return `Tokens ${usage.total_tokens || 0} · Prompt ${usage.prompt_tokens || 0} · Completion ${usage.completion_tokens || 0}`
}

async function copyMessage(content) {
  try {
    await navigator.clipboard.writeText(String(content || ''))
    ElMessage.success('已复制对话内容')
  } catch (_) {
    ElMessage.error('复制失败')
  }
}

async function buildDataSnapshot() {
  const [overview, chainStats] = await Promise.all([
    requestPreferDB('/api/db/dashboard/overview', '/api/dashboard/overview').catch(() => null),
    request('/api/chain/stats', { withRole: true }).catch(() => null)
  ])
  return {
    users: overview?.users ?? '-',
    receivableSize: overview?.receivableSize ?? '-',
    totalAmount: overview?.totalAmount ?? '-',
    statusCounters: overview?.statusCounters ?? null,
    blockHeight: chainStats?.blockHeight ?? '-',
    totalTxCount: chainStats?.totalTxCount ?? '-',
    nodeCount: chainStats?.nodeCount ?? '-'
  }
}

async function genAssistantReply(rawText) {
  const text = rawText.toLowerCase()
  const snapshot = await buildDataSnapshot()

  if (text.includes('流程') || text.includes('清单')) {
    return [
      '建议流程：企业注册 -> 应收创建(含凭证CID) -> 应收确权 -> 融资登记 -> 到期结算。',
      '关键检查：ID唯一性、金额与到期日合法性、付款方地址有效、凭证与应收ID关联一致。',
      `当前系统参考：应收总量 ${snapshot.receivableSize}，总金额 ${snapshot.totalAmount}。`
    ].join('\n')
  }

  if (text.includes('风险') || text.includes('预警') || text.includes('逾期')) {
    return [
      '风险预警建议：',
      '1. 关注临近到期且未确权/未融资的应收。',
      '2. 对同付款方短期内集中新增应收设置阈值告警。',
      '3. 将链上状态变化频率与数据库修改频率进行比对，识别异常操作。',
      `链上参考：区块高度 ${snapshot.blockHeight}，累计交易 ${snapshot.totalTxCount}。`
    ].join('\n')
  }

  if (text.includes('隐私') || text.includes('安全') || text.includes('合规')) {
    return [
      '金融数据安全设计建议：',
      '1. 默认最小化披露：仅传输风控所需字段，不上传原始证件与明文账户号。',
      '2. 本地脱敏优先：地址、订单号、手机号等先在前端掩码后再进入对话。',
      '3. 结果留痕分层：对话日志仅保留摘要，原始问答不落库或设置短周期过期。',
      '4. 链上存证只放哈希/CID，不直接上链敏感原文。',
      `当前节点规模 ${snapshot.nodeCount}，建议按节点与角色做分级访问控制。`
    ].join('\n')
  }

  return [
    `已收到你的问题（${roleLabel.value}）。`,
    '你可以继续告诉我具体目标，例如：融资审查、应收异常排查、凭证合规检查、链上数据核验。',
    `系统快照：用户 ${snapshot.users}，应收 ${snapshot.receivableSize}，交易 ${snapshot.totalTxCount}。`
  ].join('\n')
}

async function callRemoteAssistant(message) {
  const resp = await request('/api/ai/chat', {
    method: 'POST',
    withRole: true,
    body: {
      message,
      scene: auth.role,
      guard: guard.value
    }
  })
  return {
    reply: String(resp?.reply || '').trim(),
    meta: {
      model: resp?.meta?.model || 'doubao',
      usage: resp?.meta?.usage || null
    }
  }
}

async function sendMessage(prefill) {
  const raw = (prefill || input.value || '').trim()
  if (!raw || sending.value) return

  const outbound = guard.value.localMasking ? applyMasking(raw) : raw
  pushMessage('user', outbound)
  input.value = ''
  sending.value = true

  try {
    let reply = ''
    let meta = null
    try {
      const remote = await callRemoteAssistant(outbound)
      reply = remote.reply
      meta = remote.meta
    } catch (_) {
      reply = await genAssistantReply(outbound)
      meta = { model: 'local-fallback', usage: null }
    }

    let finalReply = reply
    if (guard.value.privacyMode) {
      finalReply += '\n\n[Privacy Guard] 已启用最小化披露与本地脱敏策略。'
    }
    if (guard.value.noRetention) {
      finalReply += '\n[No-Retention] 当前会话建议仅在本地短期保存。'
    }

    pushMessage('assistant', finalReply, meta)
  } catch (_) {
    pushMessage('assistant', '当前助手服务繁忙，请稍后重试。', { model: 'error' })
  } finally {
    sending.value = false
    nextTick(() => {
      if (chatBox.value) chatBox.value.scrollTop = chatBox.value.scrollHeight
    })
  }
}

function clearChat() {
  messages.value = [
    {
      role: 'assistant',
      content: '会话已清空。你可以继续咨询供应链金融业务问题。',
      time: nowTime(),
      meta: {
        model: currentModel.value
      }
    }
  ]
  ElMessage.success('会话已清空')
}

onMounted(() => {
  nextTick(() => {
    if (chatBox.value) chatBox.value.scrollTop = chatBox.value.scrollHeight
  })
})
</script>

<template>
  <section class="card ai-assistant-page">
    <div class="row between">
      <div>
        <h2>供应链金融对话助手</h2>
        <div class="sub">{{ roleLabel }} · Built on FISCO BCOS</div>
      </div>
      <div class="row wrap">
        <span class="model-chip">Model · {{ currentModel }}</span>
        <el-button size="small" plain @click="clearChat">清空会话</el-button>
      </div>
    </div>

    <div class="ai-grid">
      <aside class="form-block">
        <h3>安全与隐私策略</h3>
        <div class="sub">金融场景默认开启，优先保护敏感数据。</div>
        <div class="guard-list">
          <el-switch v-model="guard.privacyMode" active-text="隐私模式" />
          <el-switch v-model="guard.localMasking" active-text="本地脱敏" />
          <el-switch v-model="guard.noRetention" active-text="不留存原文" />
          <el-switch v-model="guard.minimalDisclosure" active-text="最小化披露" />
        </div>
        <el-alert
          type="info"
          :closable="false"
          show-icon
          title="建议：链上仅存证哈希/CID，明文凭证放受控存储并设置权限与到期策略。"
        />

        <h3 style="margin-top: 12px;">快捷问题</h3>
        <div class="quick-list">
          <el-button
            v-for="q in quickPrompts"
            :key="q"
            plain
            class="quick-btn"
            @click="sendMessage(q)"
          >
            {{ q }}
          </el-button>
        </div>
      </aside>

      <section class="form-block">
        <h3>对话区</h3>
        <div ref="chatBox" class="chat-box">
          <div
            v-for="(message, index) in messages"
            :key="index"
            class="msg"
            :class="message.role === 'user' ? 'msg-user' : 'msg-ai'"
          >
            <div class="msg-head">
              <div class="msg-role">{{ message.role === 'user' ? '你' : '助手' }}</div>
              <div class="msg-tools">
                <span v-if="message.meta?.model" class="msg-chip">{{ message.meta.model }}</span>
                <span v-if="message.meta?.usage" class="msg-chip usage-chip">{{ usageText(message.meta.usage) }}</span>
                <el-button size="small" plain @click="copyMessage(message.content)">复制</el-button>
              </div>
            </div>
            <div
              v-if="message.role === 'assistant'"
              class="msg-content md-content"
              v-html="renderMarkdown(message.content)"
            ></div>
            <div v-else class="msg-content user-content">{{ message.content }}</div>
            <div class="msg-time">{{ message.time }}</div>
          </div>
          <div v-if="sending" class="sub">助手正在生成回复...</div>
        </div>

        <div class="input-wrap">
          <el-input
            v-model="input"
            type="textarea"
            :rows="4"
            placeholder="输入你的业务问题，例如：如何降低融资违约风险？"
            @keyup.ctrl.enter="sendMessage()"
          />
          <div class="row between" style="margin-top:8px;">
            <div class="sub">Ctrl + Enter 发送</div>
            <el-button type="primary" :loading="sending" @click="sendMessage()">发送</el-button>
          </div>
        </div>
      </section>
    </div>
  </section>
</template>

<style scoped>
.ai-assistant-page {
  display: grid;
  gap: 12px;
}

.ai-grid {
  display: grid;
  grid-template-columns: 320px 1fr;
  gap: 12px;
}

.guard-list {
  display: grid;
  gap: 10px;
  margin: 10px 0 12px;
}

.quick-list {
  display: grid;
  gap: 8px;
  margin-top: 8px;
}

.quick-btn {
  justify-content: flex-start;
  margin: 0;
  white-space: normal;
  height: auto;
  padding: 8px 10px;
}

.model-chip,
.msg-chip {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  border-radius: 999px;
  border: 1px solid #d2e3ff;
  background: #eef5ff;
  color: #2b5eb6;
  font-size: 12px;
}

.usage-chip {
  background: #f5f9ff;
  color: #55739c;
}

.chat-box {
  height: 460px;
  overflow: auto;
  border: 1px solid #d8e8ff;
  border-radius: 12px;
  background: rgba(248, 252, 255, 0.9);
  padding: 10px;
  display: grid;
  gap: 8px;
}

.msg {
  border-radius: 10px;
  padding: 8px 10px;
  border: 1px solid #d8e6ff;
}

.msg-user {
  background: #eef5ff;
}

.msg-ai {
  background: #f7fbff;
}

.msg-head {
  display: flex;
  justify-content: space-between;
  gap: 8px;
  align-items: flex-start;
  margin-bottom: 4px;
}

.msg-tools {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  justify-content: flex-end;
}

.msg-role {
  font-size: 12px;
  color: #5a6d86;
  margin-bottom: 4px;
}

.msg-content {
  line-height: 1.6;
}

.user-content {
  white-space: pre-wrap;
}

.md-content :deep(p) {
  margin: 0 0 8px;
}

.md-content :deep(p:last-child) {
  margin-bottom: 0;
}

.md-content :deep(strong) {
  color: #163d84;
  font-weight: 700;
}

.md-content :deep(ul),
.md-content :deep(ol) {
  margin: 8px 0;
  padding-left: 20px;
}

.md-content :deep(li) {
  margin: 4px 0;
}

.md-content :deep(code) {
  padding: 2px 6px;
  border-radius: 6px;
  background: #eaf3ff;
  border: 1px solid #d4e6ff;
  font-size: 12px;
}

.md-content :deep(blockquote) {
  margin: 8px 0;
  padding: 8px 12px;
  border-left: 3px solid #5f9dff;
  background: rgba(233, 243, 255, 0.85);
  color: #36557f;
  border-radius: 0 10px 10px 0;
}

.msg-time {
  margin-top: 6px;
  font-size: 11px;
  color: #7387a5;
}

.input-wrap {
  margin-top: 10px;
}

@media (max-width: 1100px) {
  .ai-grid {
    grid-template-columns: 1fr;
  }

  .chat-box {
    height: 360px;
  }

  .msg-head {
    flex-direction: column;
  }

  .msg-tools {
    justify-content: flex-start;
  }
}
</style>
