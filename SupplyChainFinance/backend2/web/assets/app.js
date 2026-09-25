const { createApp, reactive, computed, onMounted } = Vue;

createApp({
  setup() {
    const state = reactive({
      current: "auth",
      user: null,
      users: [],
      receivables: [],
      detail: null,
      overview: null,
      log: "",
      forms: {
        register: { username: "", password: "", role: "business" },
        login: { username: "", password: "" },
        enterprise: { enterpriseId: "", name: "", isCore: false, isFinance: false },
        issue: { receivableId: "", orderId: "", payer: "", amount: "", dueTimestamp: "" },
        confirm: { receivableId: "" },
        finance: { receivableId: "", amount: "", interestRate: "" },
        settle: { receivableId: "" },
        queryId: ""
      }
    });

    const role = computed(() => (state.user?.role || "guest").toLowerCase());

    const menus = computed(() => {
      const list = [{ key: "auth", label: "登录/注册" }, { key: "dashboard", label: "总览仪表盘" }, { key: "receivables", label: "应收账款池" }];
      if (["business", "admin"].includes(role.value)) list.push({ key: "business", label: "企业端流程" });
      if (["finance", "admin"].includes(role.value)) list.push({ key: "finance", label: "金融端流程" });
      return list;
    });

    function roleHeaders() {
      return state.user ? { "X-User-Role": state.user.role } : {};
    }

    async function api(path, method = "GET", body, withRole = false) {
      const headers = { "Content-Type": "application/json" };
      if (withRole) Object.assign(headers, roleHeaders());
      const res = await fetch(path, { method, headers, body: body ? JSON.stringify(body) : undefined });
      const data = await res.json().catch(() => ({}));
      if (!res.ok) throw data;
      return data;
    }

    function setLog(title, payload) {
      state.log = `${title}\n${JSON.stringify(payload, null, 2)}`;
    }

    async function loadUsers() {
      const r = await api("/api/auth/users");
      state.users = r.data || [];
      if (!state.forms.issue.payer && state.users.length) state.forms.issue.payer = state.users[0].address;
    }

    async function loadReceivables() {
      const r = await api("/api/receivable");
      state.receivables = r.data || [];
    }

    async function loadOverview() {
      const r = await api("/api/dashboard/overview");
      state.overview = r;
    }

    async function registerUser() {
      try {
        const r = await api("/api/auth/register", "POST", state.forms.register);
        setLog("注册成功", r);
        await loadUsers();
      } catch (e) { setLog("注册失败", e); }
    }

    async function login() {
      try {
        const r = await api("/api/auth/login", "POST", state.forms.login);
        state.user = r;
        setLog("登录成功", r);
        if (!menus.value.find((x) => x.key === state.current)) state.current = "dashboard";
      } catch (e) { setLog("登录失败", e); }
    }

    async function registerEnterprise() {
      try {
        const payload = { ...state.forms.enterprise, isFinance: state.forms.enterprise.isFinance };
        const r = await api("/api/business/enterprise/register", "POST", payload, true);
        setLog("企业注册上链成功", r);
      } catch (e) { setLog("企业注册失败", e); }
    }

    async function issueReceivable() {
      try {
        const f = state.forms.issue;
        const r = await api("/api/business/receivable/issue", "POST", {
          receivableId: f.receivableId,
          orderId: f.orderId,
          payer: f.payer,
          amount: Number(f.amount),
          dueTimestamp: Number(f.dueTimestamp)
        }, true);
        setLog("应收创建成功", r);
        await loadReceivables();
        await loadOverview();
      } catch (e) { setLog("应收创建失败", e); }
    }

    async function confirmReceivable() {
      try {
        const r = await api("/api/business/receivable/confirm", "POST", { receivableId: state.forms.confirm.receivableId }, true);
        setLog("确权成功", r);
        await loadReceivables();
        await loadOverview();
      } catch (e) { setLog("确权失败", e); }
    }

    async function settleReceivable() {
      try {
        const r = await api("/api/business/receivable/settle", "POST", { receivableId: state.forms.settle.receivableId }, true);
        setLog("结算成功", r);
        await loadReceivables();
        await loadOverview();
      } catch (e) { setLog("结算失败", e); }
    }

    async function financeReceivable() {
      try {
        const f = state.forms.finance;
        const r = await api("/api/finance/receivable/finance", "POST", {
          receivableId: f.receivableId,
          amount: Number(f.amount),
          interestRate: Number(f.interestRate)
        }, true);
        setLog("融资登记成功", r);
        await loadReceivables();
        await loadOverview();
      } catch (e) { setLog("融资登记失败", e); }
    }

    async function queryDetail() {
      const id = state.forms.queryId.trim();
      if (!id) return;
      try {
        const r = await api(`/api/receivable/${id}`);
        state.detail = r;
      } catch (e) {
        state.detail = e;
      }
    }

    onMounted(async () => {
      try {
        await Promise.all([loadUsers(), loadReceivables(), loadOverview()]);
      } catch (_) {}
    });

    return {
      state, role, menus,
      registerUser, login,
      registerEnterprise, issueReceivable, confirmReceivable, settleReceivable, financeReceivable,
      loadReceivables, loadOverview, queryDetail
    };
  },
  template: `
  <div class="app">
    <aside class="sidebar">
      <div class="brand">SCF Vue3 Console</div>
      <div class="role-chip">角色: {{ role }}</div>
      <button v-for="m in menus" :key="m.key" class="nav-btn" :class="{active: state.current===m.key}" @click="state.current=m.key">{{m.label}}</button>
    </aside>

    <section class="main">
      <div class="topbar">
        <div>
          <strong>{{ state.user ? state.user.username : '未登录' }}</strong>
          <div class="sub">{{ state.user ? state.user.address : '请先登录后再执行链上写操作' }}</div>
        </div>
        <div class="row">
          <button class="btn secondary" @click="loadReceivables">刷新应收池</button>
          <button class="btn secondary" @click="loadOverview">刷新仪表盘</button>
        </div>
      </div>

      <div class="cards" v-if="state.current==='auth'">
        <div class="card">
          <h2>用户注册</h2>
          <div class="form">
            <input class="input" v-model="state.forms.register.username" placeholder="用户名" />
            <input class="input" type="password" v-model="state.forms.register.password" placeholder="密码" />
            <select v-model="state.forms.register.role">
              <option value="business">business</option>
              <option value="finance">finance</option>
              <option value="admin">admin</option>
            </select>
            <button class="btn" @click="registerUser">注册</button>
          </div>
        </div>
        <div class="card">
          <h2>用户登录</h2>
          <div class="form">
            <input class="input" v-model="state.forms.login.username" placeholder="用户名" />
            <input class="input" type="password" v-model="state.forms.login.password" placeholder="密码" />
            <button class="btn" @click="login">登录并切换链上账户</button>
          </div>
        </div>
      </div>

      <div class="cards" v-if="state.current==='dashboard'">
        <div class="card wide">
          <h2>总览指标</h2>
          <div class="kpi" v-if="state.overview">
            <div><div class="sub">用户数</div><strong>{{state.overview.users}}</strong></div>
            <div><div class="sub">应收总量</div><strong>{{state.overview.receivableSize}}</strong></div>
            <div><div class="sub">总金额</div><strong>{{state.overview.totalAmount}}</strong></div>
            <div><div class="sub">待确权</div><strong>{{state.overview.statusCounters?.PENDING}}</strong></div>
            <div><div class="sub">已确权</div><strong>{{state.overview.statusCounters?.CONFIRMED}}</strong></div>
            <div><div class="sub">已融资</div><strong>{{state.overview.statusCounters?.FINANCED}}</strong></div>
            <div><div class="sub">已结算</div><strong>{{state.overview.statusCounters?.SETTLED}}</strong></div>
          </div>
        </div>
        <div class="card wide">
          <h3>最近交易</h3>
          <table class="table">
            <thead><tr><th>时间</th><th>动作</th><th>业务ID</th><th>状态</th><th>哈希</th></tr></thead>
            <tbody>
              <tr v-for="x in (state.overview?.recentTxLogs || [])" :key="x.id">
                <td>{{x.createdAt}}</td><td>{{x.action}}</td><td>{{x.businessId}}</td><td>{{x.status}}</td><td>{{x.txHash}}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="cards" v-if="state.current==='business' && ['business','admin'].includes(role)">
        <div class="card">
          <h2>企业注册上链</h2>
          <div class="form">
            <input class="input" v-model="state.forms.enterprise.enterpriseId" placeholder="enterpriseId" />
            <input class="input" v-model="state.forms.enterprise.name" placeholder="企业名称" />
            <label><input type="checkbox" v-model="state.forms.enterprise.isCore" /> 核心企业</label>
            <label><input type="checkbox" v-model="state.forms.enterprise.isFinance" /> 金融机构</label>
            <button class="btn" @click="registerEnterprise">注册</button>
          </div>
        </div>

        <div class="card">
          <h2>应收创建</h2>
          <div class="form">
            <input class="input" v-model="state.forms.issue.receivableId" placeholder="receivableId" />
            <input class="input" v-model="state.forms.issue.orderId" placeholder="orderId" />
            <select v-model="state.forms.issue.payer">
              <option v-for="u in state.users" :value="u.address">{{u.username}} | {{u.address}}</option>
            </select>
            <input class="input" type="number" v-model="state.forms.issue.amount" placeholder="amount" />
            <input class="input" type="number" v-model="state.forms.issue.dueTimestamp" placeholder="due timestamp" />
            <button class="btn" @click="issueReceivable">创建应收</button>
          </div>
        </div>

        <div class="card">
          <h2>确权与结算</h2>
          <div class="form">
            <input class="input" v-model="state.forms.confirm.receivableId" placeholder="确权 receivableId" />
            <button class="btn warn" @click="confirmReceivable">确权</button>
            <input class="input" v-model="state.forms.settle.receivableId" placeholder="结算 receivableId" />
            <button class="btn danger" @click="settleReceivable">结算</button>
          </div>
        </div>
      </div>

      <div class="cards" v-if="state.current==='finance' && ['finance','admin'].includes(role)">
        <div class="card">
          <h2>融资登记</h2>
          <div class="form">
            <input class="input" v-model="state.forms.finance.receivableId" placeholder="receivableId" />
            <input class="input" type="number" v-model="state.forms.finance.amount" placeholder="融资金额" />
            <input class="input" type="number" v-model="state.forms.finance.interestRate" placeholder="利率(万分比)" />
            <button class="btn" @click="financeReceivable">提交融资</button>
          </div>
        </div>
      </div>

      <div class="cards" v-if="state.current==='receivables'">
        <div class="card wide">
          <h2>应收账款池</h2>
          <div class="row">
            <input class="input" v-model="state.forms.queryId" placeholder="输入应收ID查询详情" />
            <button class="btn secondary" @click="queryDetail">查询详情</button>
          </div>
          <table class="table">
            <thead><tr><th>ID</th><th>订单</th><th>发行方</th><th>付款方</th><th>金额</th><th>状态</th></tr></thead>
            <tbody>
              <tr v-for="x in state.receivables" :key="x.receivableId">
                <td>{{x.receivableId}}</td><td>{{x.orderId}}</td><td>{{x.issuer}}</td><td>{{x.payer}}</td><td>{{x.amount}}</td><td>{{x.statusText}}</td>
              </tr>
            </tbody>
          </table>
          <div class="log" v-if="state.detail">{{ JSON.stringify(state.detail, null, 2) }}</div>
        </div>
      </div>

      <div class="card wide">
        <h3>操作日志</h3>
        <div class="log">{{ state.log || '暂无日志' }}</div>
      </div>
    </section>
  </div>
  `
}).mount("#app");
