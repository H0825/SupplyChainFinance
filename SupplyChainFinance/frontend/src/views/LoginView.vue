<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { CopyDocument } from '@element-plus/icons-vue'
import { useAuthStore } from '../stores/auth'
import { request } from '../api/http'

const auth = useAuthStore()
const router = useRouter()
const mode = ref('login')
const loading = ref(false)
const booting = ref(true)
const log = ref('')

const reg = reactive({ username: '', password: '', role: 'business' })
const loginForm = reactive({ username: '', password: '' })
const registerResult = ref(null)

const pwdStrength = computed(() => {
  const pwd = reg.password || ''
  let score = 0
  if (pwd.length >= 8) score++
  if (/[A-Z]/.test(pwd) && /[a-z]/.test(pwd)) score++
  if (/\d/.test(pwd)) score++
  if (/[^\w]/.test(pwd)) score++
  return score
})

const pwdLabel = computed(() => ['弱', '一般', '中', '强', '很强'][pwdStrength.value])
const pwdType = computed(() => ['danger', 'warning', '', 'success', 'success'][pwdStrength.value] || '')

const logAlertType = computed(() => {
  if (!log.value) return 'info'
  if (log.value.includes('成功') || log.value.toLowerCase().includes('success')) return 'success'
  if (log.value.includes('失败') || log.value.toLowerCase().includes('fail')) return 'error'
  return 'info'
})

function jumpByRole(user) {
  const role = (user.role || '').toLowerCase()
  if (role === 'business') router.push('/business')
  else if (role === 'finance') router.push('/finance')
  else router.push('/dashboard')
}

async function copyText(text, label) {
  await navigator.clipboard.writeText(String(text || ''))
  ElMessage.success(`${label} 已复制`)
}

async function onRegister() {
  loading.value = true
  try {
    const res = await request('/api/auth/register', { method: 'POST', body: reg })
    registerResult.value = res
    log.value = `注册成功：链账户 ${res.address}`
    mode.value = 'login'
    loginForm.username = reg.username
    loginForm.password = reg.password
  } catch (e) {
    registerResult.value = null
    log.value = `注册失败
${JSON.stringify(e, null, 2)}`
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  window.setTimeout(() => {
    booting.value = false
  }, 260)
})

async function onLogin() {
  loading.value = true
  try {
    const res = await request('/api/auth/login', { method: 'POST', body: loginForm })
    auth.setUser(res)
    log.value = `登录成功，当前角色：${res.role}`
    jumpByRole(res)
  } catch (e) {
    log.value = `登录失败
${JSON.stringify(e, null, 2)}`
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <section class="auth-page">
    <div class="auth-frame">
      <div class="auth-shell">
        <aside class="auth-side">
          <div class="auth-side-tag">IDENTITY GATEWAY</div>
          <h2>区块链账户准入中心</h2>
          <p>完成注册后自动创建链账户，登录后按角色进入对应控制台。</p>

          <div class="auth-side-grid">
            <div class="auth-metric">
              <strong>3</strong>
              <span>角色体系</span>
            </div>
            <div class="auth-metric">
              <strong>4+</strong>
              <span>业务工作台</span>
            </div>
            <div class="auth-metric">
              <strong>100%</strong>
              <span>链路可追溯</span>
            </div>
            <div class="auth-metric">
              <strong>24/7</strong>
              <span>实时查询</span>
            </div>
          </div>

          <el-alert
            title="注册成功后请妥善保存 PrivateKey，避免账户风险。"
            type="warning"
            :closable="false"
            show-icon
          />

          <button class="btn mini ghost" @click="router.push('/welcome')">返回欢迎页</button>
        </aside>

        <section class="auth-form-pane auth-tech">
          <div v-if="booting" class="auth-boot-skeleton">
            <el-skeleton animated :rows="4" />
            <el-skeleton animated :rows="6" />
          </div>

          <template v-else>
            <div class="auth-head row between">
              <h3>{{ mode === 'login' ? '登录系统' : '注册账户' }}</h3>
              <div class="mode-switch">
                <button class="btn mini" :class="{ ghost: mode !== 'login' }" @click="mode='login'">登录</button>
                <button class="btn mini" :class="{ ghost: mode !== 'register' }" @click="mode='register'">注册</button>
              </div>
            </div>

            <div v-if="mode === 'login'" class="form">
              <input v-model="loginForm.username" class="input" placeholder="用户名" />
              <input v-model="loginForm.password" class="input" type="password" placeholder="密码" />
              <button class="btn neon" :disabled="loading" @click="onLogin">{{ loading ? '提交中...' : '登录并进入控制台' }}</button>
            </div>

            <div v-else class="form">
              <input v-model="reg.username" class="input" placeholder="用户名" />
              <input v-model="reg.password" class="input" type="password" placeholder="密码" />
              <div class="row">
                <span class="sub">密码强度：</span>
                <el-tag size="small" :type="pwdType">{{ pwdLabel }}</el-tag>
              </div>
              <select v-model="reg.role" class="input">
                <option value="business">business（企业）</option>
                <option value="finance">finance（金融机构）</option>
                <option value="admin">admin（管理）</option>
              </select>
              <button class="btn neon" :disabled="loading" @click="onRegister">{{ loading ? '提交中...' : '注册用户并创建链账户' }}</button>
            </div>

            <div class="chain-card" v-if="registerResult">
              <h3>注册成功 · 链上账户信息</h3>
              <div class="chain-row">
                <span class="sub">Address</span>
                <code>{{ registerResult.address }}</code>
                <el-tooltip content="复制地址" placement="top">
                  <button class="btn mini ghost" @click="copyText(registerResult.address, '地址')">
                    <el-icon><CopyDocument /></el-icon>
                  </button>
                </el-tooltip>
              </div>
              <div class="chain-row">
                <span class="sub">PrivateKey</span>
                <code>{{ registerResult.privateKey }}</code>
                <el-tooltip content="复制私钥" placement="top">
                  <button class="btn mini ghost" @click="copyText(registerResult.privateKey, '私钥')">
                    <el-icon><CopyDocument /></el-icon>
                  </button>
                </el-tooltip>
              </div>
            </div>

            <el-alert v-if="log" :title="log" :type="logAlertType" :closable="false" show-icon />
            <div v-else class="auth-empty-slot">
              <el-empty description="注册成功后这里会展示区块链账户信息">
                <template #image>
                  <div class="auth-empty-orbit">
                    <div class="auth-empty-ring"></div>
                    <div class="auth-empty-core">Key</div>
                  </div>
                </template>
                <div class="sub">基于 FISCO BCOS 区块链开发</div>
              </el-empty>
            </div>
          </template>
        </section>
      </div>
    </div>
  </section>
</template>

<style scoped>
.auth-boot-skeleton {
  display: grid;
  gap: 14px;
}

.auth-empty-slot {
  padding-top: 8px;
}

.auth-empty-orbit {
  position: relative;
  width: 120px;
  height: 120px;
  display: grid;
  place-items: center;
}

.auth-empty-ring {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  border: 1px solid rgba(76, 151, 255, 0.28);
  box-shadow: 0 0 0 12px rgba(103, 181, 255, 0.08), inset 0 0 24px rgba(57, 128, 239, 0.08);
}

.auth-empty-core {
  width: 56px;
  height: 56px;
  border-radius: 18px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #2f72ef, #49c7ff);
  color: white;
  font-weight: 700;
  box-shadow: 0 16px 30px rgba(58, 122, 233, 0.22);
}
</style>
