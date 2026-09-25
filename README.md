# 供应链金融系统（SupplyChainFinance）运行使用说明

## 一、项目简介

本项目是一个**基于区块链的供应链金融服务系统**，面向供应链上下游企业（供应商、核心企业）与金融机构，把传统供应链金融中"应收账款确权难、流转难、融资难"的业务流程搬到链上，借助区块链的不可篡改特性实现多方之间的可信协作。

系统围绕**应收账款（Receivable）**这一核心资产展开业务：供应商向核心企业交付货物后形成应收账款，供应商在平台上签发这笔应收账款，核心企业在链上完成确权（即承认这笔债务），供应商凭确权后的应收账款向金融机构申请融资放款，最终在到期日由核心企业完成结算清偿。整个过程中每一次状态变更都会作为一笔真实交易写入区块链，交易哈希回写到本地数据库，从而形成"链上留痕 + 链下可查"的双重保障。

围绕这一主线，系统提供以下能力：

- **应收账款全生命周期管理**：支持应收账款的签发、确权、融资、结算四个关键动作，每一步都调用智能合约完成链上状态变更，并记录交易哈希与执行结果
- **多角色协同**：区分 `business`（企业用户：供应商/核心企业）、`finance`（金融机构用户：银行/保理公司）、`admin`（管理员）三类角色，不同角色拥有不同的操作权限与工作台页面，业务流程中各方各司其职
- **企业注册与身份管理**：企业信息（名称、钱包地址、是否核心企业、是否金融机构）注册上链，形成链上的企业信用主体
- **文件存证**：合同、发票、票据等业务单据上传至 IPFS 分布式存储，返回的内容标识 CID 与业务单据类型、关联订单等元数据一并落库，实现文件的防篡改存证与快速检索
- **AI 智能助手**：接入豆包大模型（火山方舟 API），既支持业务知识问答，也支持对指定链上交易哈希进行风险解读与异常分析，降低审计与风控门槛
- **数据看板与审计**：提供链上数据统计、应收账款状态分布、金额汇总、最近交易日志等可视化视图，管理员可查看完整交易轨迹

## 二、技术架构

系统采用前后端分离架构，前端负责交互与展示，后端作为统一服务入口，向下对接三类外部基础设施：关系型数据库（业务数据的链下镜像与查询加速）、区块链网络（可信存证与状态变更）、IPFS 网络（文件分布式存储），同时向上通过 HTTP 接口向 AI 大模型服务发起调用。

```
┌─────────────────┐     HTTP (代理/直连)      ┌──────────────────────┐
│  前端 frontend   │ ──────────────────────►  │  后端 backend2 (Go)   │
│  Vue3 + Vite     │   /api /workflow /health │  Gin + MySQL 驱动     │
│  Element Plus    │                          └──────────┬───────────┘
│  Pinia + Router  │                                     │
└─────────────────┘                        ┌─────────────┼─────────────┐
                                           ▼             ▼             ▼
                                     MySQL 数据库   FISCO-BCOS 链   IPFS 节点
                                     (finance库)   (Channel连接)  (localhost:5001)
                                                       │
                                                       ▼
                                              豆包 AI (火山方舟 API)
```

### 模块构成

| 模块 | 目录 | 技术栈与职责 |
|---|---|---|
| 后端 | `backend2/` | Go 1.25 + Gin 框架；通过 `go-sql-driver/mysql` 访问数据库，通过 `FISCO-BCOS/go-sdk` 调用智能合约，通过 `go-ipfs-api` 上传文件，使用 `bcrypt` 对用户密码做哈希 |
| 前端 | `frontend/` | Vue 3 + Vite 5 构建，Element Plus 提供 UI 组件，Pinia 做状态管理，vue-router 做路由与权限守卫，marked + dompurify 渲染 AI 返回的 Markdown |
| 智能合约 | `backend2/SupplyChainFinance/` | Solidity 源码，已编译为 ABI 与 Go 绑定代码，后端通过生成的 `SupplyChainFinance.go` 直接调用合约方法 |

### 数据流转说明

一次典型的"签发应收账款"操作，数据是这样流转的：前端表单提交 → 后端校验角色与参数 → 调用智能合约 `issue` 方法 → 区块链接收并执行交易，返回交易哈希与回执 → 后端把该笔应收账款的核心字段写入 MySQL `receivables` 表，同时把交易哈希、操作人、执行状态写入 `tx_logs` 表 → 前端刷新列表展示最新状态。

可以看到，链上保存的是**权威的业务状态**，数据库保存的是**便于查询与统计的结构化副本**，两者通过 `tx_hash` 关联，任何一笔业务都可以从数据库记录追溯到链上的原始交易。这种设计既保证了数据可信，又避免了每次查询都要遍历链上事件的性能问题。

## 三、运行环境要求

| 依赖 | 版本要求 | 用途与说明 |
|---|---|---|
| Go | ≥ 1.25 | 编译与运行后端服务 |
| Node.js | ≥ 18（建议 20+） | 前端依赖安装、开发服务器、生产构建 |
| npm | 随 Node 一起安装 | 前端包管理 |
| MySQL | 8.x | 持久化业务数据，库名固定为 `finance` |
| FISCO-BCOS 节点 | 版本与链配置匹配，Channel 端口可访问 | 区块链网络，承载智能合约与交易 |
| IPFS 节点 | 本机 `localhost:5001` | 文件分布式存储（仅上传文件功能需要） |
| 网络 | 可访问 `ark.cn-beijing.volces.com` | AI 助手功能（仅 AI 相关页面需要） |

> **关于启动顺序与外部依赖的说明**
>
> - MySQL 的数据表会由后端在启动时自动创建（见 `dbs/schema.go` 中的 `EnsureSchema` 方法，包含建表与字段迁移逻辑），因此**空数据库也能直接启动**，不必手动逐条执行建表语句。但数据库本身（`finance` 库）需要先手工创建。
> - 区块链节点的连接是在 `router/bcos.go` 的 `init()` 函数里完成的，一旦连接失败会直接 `log.Fatal` 终止进程。也就是说**区块链节点不可达时后端根本起不来**，这是本系统最关键的外部依赖。
> - IPFS 与 AI 服务属于**可降级依赖**：不启动 IPFS 节点时，除文件上传与文档存证外的其他功能均正常；AI 服务不可用时，仅 AI 助手和风险分析两个页面受影响。

## 四、配置说明（backend2/config.toml）

后端在启动时通过 `appconfig.Load("config.toml")` 读取配置文件，路径是**相对于当前工作目录**的相对路径，因此**必须在 `backend2/` 目录下执行启动命令**，换目录运行会直接报 `Load app config failed` 并退出。

```toml
[App]
listen = ":8888"                          # 后端 HTTP 监听端口
log_file = "gin-runner.log"               # Gin 访问日志与运行日志输出文件
frontend_dist = "../frontend/dist"        # 前端构建产物目录，存在则由后端直接托管
legacy_index = "./web/index.html"         # 兜底静态页：前端未构建时使用的旧版页面

[Database]                                # MySQL 连接信息，按实际环境修改
user = "root"
password = "123456"
host = "localhost"
port = 3306
name = "finance"                          # 库名，需先手动执行 CREATE DATABASE finance;
charset = "utf8mb4"
parse_time = true
loc = "Local"

[Contract]
address = "0x0a1f..."                     # 已部署的 SupplyChainFinance 智能合约地址

[AI]                                      # 豆包大模型（火山方舟）配置
base_url = "https://ark.cn-beijing.volces.com/api/v3"
api_key = "<你的API Key>"                 # ⚠️ 请替换为自己的 Key
model = "doubao-1-5-pro-32k-250115"
timeout_sec = 30                          # 单次请求超时时间（秒）
temperature = 0.2                         # 生成温度，值越小输出越稳定

[Network]                                 # FISCO-BCOS Channel 连接配置，证书位于 sdk/ 目录
Type="channel"
CAFile="sdk/ca.crt"
Cert="sdk/sdk.crt"
Key="sdk/sdk.key"
[[Network.Connection]]
NodeURL="111.228.57.198:20200"            # 节点 Channel 端口，按实际环境修改
GroupID=1                                 # 群组 ID，需与合约部署的群组一致
# [[Network.Connection]]                  # 如需连接多群组/多节点，可取消注释继续追加
# NodeURL="127.0.0.1:20200"
# GroupID=2

[Account]
KeyFile="ci/0x....pem"                    # 交易签名用的链上账户私钥文件

[Chain]
ChainID=1
SMCrypto=false                            # 非国密链填 false，国密链填 true
```

### 需要重点核对的配置项

首次部署时，以下几项最容易因为环境差异而需要调整：

1. `[Database]` 段的连接账号密码、主机地址与端口；
2. `[Network]` 段中 `[[Network.Connection]]` 的 `NodeURL` 与 `GroupID`，必须与实际的 FISCO-BCOS 节点地址和群组一致；
3. `sdk/` 目录下的三份证书（`ca.crt`、`sdk.crt`、`sdk.key`）必须是与上述节点配套签发的，否则 Channel 握手会失败；
4. `[Contract]` 段的合约地址必须是当前群组里已部署的合约，换链或重新部署合约后都需要同步更新；
5. `[AI]` 段的 `api_key` 需要替换为你自己的火山方舟 API Key，否则 AI 相关功能会调用失败。

> **⚠️ 安全提示**
>
> 1. `config.toml` 中同时包含数据库密码、AI 服务 API Key 与链上账户私钥，属于**高度敏感文件**，切勿提交到公开代码仓库或在公开场合传播；
> 2. 当前仓库中已写入了一个明文的 API Key，建议尽快到火山方舟控制台将其作废并更换为自己的 Key，避免额度被盗用；
> 3. 正式部署时，建议把敏感项改为通过环境变量注入，而不是硬编码在配置文件里。

## 五、快速启动步骤

### 第 1 步：准备数据库

先确保本机 MySQL 服务已经启动，然后使用有建库权限的账号登录，创建名为 `finance` 的数据库：

```sql
-- 用 root 登录 MySQL 后执行
CREATE DATABASE finance DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
```

创建完成后，后端启动时会自动在库中建立 `users`、`tx_logs`、`enterprises`、`receivables`、`financing_records`、`documents` 六张表，并对 `documents` 表做必要的字段迁移，无需手工干预。

如果希望系统一上来就有可查看的业务数据（推荐用于演示与功能验证），可以导入项目自带的演示数据：

```bash
mysql -uroot -p finance < backend2/dbs/finance-data.sql
```

该脚本包含企业信息、应收账款、融资记录、文档记录与 27 条链上交易日志，并预置了 3 个演示账号（密码以 bcrypt 密文形式存储）：

| 用户名 | 角色 | 可进入的页面 |
|---|---|---|
| user001 | business（企业） | 企业工作台、应收账款、文档存证、AI 助手 |
| jr001 | finance（金融机构） | 金融工作台、应收账款、风控实验室、文档存证、AI 助手 |
| admin001 | admin（管理员） | 全部页面，含审计视图与控制中心 |

> 如果演示账号的密码未知，完全不必纠结——直接打开登录页注册一个新用户即可。注册时系统会自动在 `keystore/` 目录下为该用户生成一份链上账户密钥文件，并分配钱包地址，注册成功后即可正常登录使用。

当然，不导入演示数据也完全可以：表结构会自动创建，注册账号后从零开始跑一遍完整业务流程，反而更有助于理解系统的运作方式。

### 第 2 步：启动后端

打开终端，进入后端目录并运行：

```bash
cd backend2

# 首次运行需要下载依赖，请确保可访问 Go 模块代理（国内建议设置 GOPROXY）
go mod download

# 启动服务（务必在 backend2 目录内执行，否则读不到 config.toml）
go run .
```

首次启动时，后端会依次完成这些动作：加载配置 → 连接 MySQL 并自动建表 → 打开日志文件 → 注册所有路由与中间件 → 初始化区块链客户端（连节点、绑定合约）→ 检测前端构建产物是否存在并决定静态资源目录 → 在 `8888` 端口监听请求。

当终端出现类似下面的输出时，说明启动成功：

```
[GIN-debug] Listening and serving HTTP on :8888
```

此时可以做一个最简单的连通性验证：

```bash
curl http://localhost:8888/health
# 返回 {"message":"ok"} 表示服务正常
```

> Windows 环境下，也可以直接运行目录内已编译好的 `backend.exe`，效果与 `go run .` 完全一致，同样要求工作目录为 `backend2/`。

### 第 3 步：启动前端

前端提供两种使用方式，按需选择即可。

**方式 A：开发模式（带热更新，适合调试与二次开发）**

```bash
cd frontend
npm install        # 仅首次需要，安装项目依赖
npm run dev        # 启动开发服务器，默认 http://localhost:5173
```

开发模式下，Vite 已在 `vite.config.js` 中配置好代理，会把 `/api`、`/workflow`、`/health` 三类请求自动转发到 `http://localhost:8888`，因此前端代码里可以直接写相对路径，无需关心跨域问题。修改源码后浏览器会自动热更新，便于边改边看。

**方式 B：生产模式（构建产物由后端统一托管，适合演示与部署）**

```bash
cd frontend
npm run build      # 构建产物输出到 frontend/dist 目录
```

构建完成后不再需要单独启动前端服务，**直接访问 http://localhost:8888 即可**。因为后端启动时会检测 `../frontend/dist/index.html` 是否存在，若存在则自动托管前端页面，同时挂载 `/assets` 静态资源目录，并对所有未匹配的 GET 请求回落到 `index.html`，从而支持 Vue 的 history 路由（刷新子页面不会 404）。仓库中已经附带了一份构建好的 `dist/`，可以直接使用。

> 两种方式的取舍很简单：改代码用方式 A，跑演示或部署用方式 B。若同时启动，请注意 5173 与 8888 是两个独立入口，页面功能一致。

### 第 4 步：登录使用

1. 浏览器打开 `http://localhost:5173`（开发模式）或 `http://localhost:8888`（生产模式）；
2. 系统默认跳转到 `/welcome` 欢迎页，点击登录进入 `/login` 页面；
3. 使用已有账号登录，或切换到注册页签注册新用户（注册时可选择 business / finance / admin 三种角色）；
4. 登录成功后，用户信息（用户名、钱包地址、角色）会保存在浏览器 `localStorage` 的 `scf_user` 键中，后续所有接口请求都会自动携带对应的角色标识；
5. 退出登录时清除该键即可，刷新页面后未登录状态会被路由守卫拦截并跳回登录页。

## 六、角色与页面功能

不同角色登录后看到的侧边栏菜单是不同的，这是由前端路由的 `meta.roles` 与后端的接口级校验共同控制的——前端负责隐藏无权限的入口，后端负责拦截越权的请求，两层防护缺一不可。

| 页面 | 路由 | 允许角色 | 功能说明 |
|---|---|---|---|
| 欢迎页 | `/welcome` | 公开访问 | 项目整体介绍与功能概览 |
| 登录/注册 | `/login` | 公开访问 | 账号登录与新用户注册（注册时分配链上账户） |
| 数据看板 | `/dashboard` | 所有登录用户 | 应收账款总量、状态分布、最近链上交易日志等统计视图 |
| 应收账款 | `/receivables` | 所有登录用户 | 应收账款列表查询、详情查看、关联融资记录查看 |
| 企业工作台 | `/business` | business、admin | 企业注册上链、签发应收账款、确权、到期结算 |
| 金融机构工作台 | `/finance` | finance、admin | 查看待融资应收账款、执行融资放款操作 |
| 审计视图 | `/audit` | admin | 链上交易日志全量审计与追溯 |
| 文档存证 | `/documents` | business、finance、admin | 业务单据上传至 IPFS，查看历史存证记录 |
| 风控实验室 | `/risk-lab` | finance、admin | 输入链上交易哈希，调用 AI 进行风险解读 |
| AI 助手 | `/ai-assistant` | business、finance、admin | 基于豆包大模型的业务问答对话 |
| 控制中心 | `/control-center` | admin | 系统层面的管理与配置入口 |

### 业务流程（应收账款生命周期）

一笔应收账款从产生到消亡，会经历四个阶段，每个阶段对应一次链上交易与一个状态码：

```
企业签发(issue) ──► 核心企业确权(confirm) ──► 金融机构融资(finance) ──► 到期结算(settle)
   status=0 PENDING      status=1 CONFIRMED       status=2 FINANCED       status=3 SETTLED
```

- **签发（PENDING）**：供应商基于真实贸易背景，在平台上登记一笔应收账款，录入订单号、付款方、金额、到期日等信息，上链后状态为待确权；
- **确权（CONFIRMED）**：核心企业（付款方）核对贸易背景无误后，在链上确认这笔债务。确权是关键一步，它让应收账款从"供应商单方面声明"变成"核心企业承认的债务"，也是后续融资的信用基础；
- **融资（FINANCED）**：供应商急需资金时，可将已确权的应收账款转让给金融机构，金融机构按一定利率放款，系统记录融资金额、利率、放款时间等信息；
- **结算（SETTLED）**：到期日由核心企业向金融机构或供应商清偿款项，应收账款状态置为已结清，整个生命周期结束。

上述每一步操作都会在链上产生一笔交易，后端会把交易哈希、执行状态、操作人地址写入 MySQL 的 `tx_logs` 表，供看板展示与审计追溯使用。若某一步调用失败（例如角色不对、状态不合法），链上不会改变状态，后端会在交易日志中记录失败原因，前端也会给出相应提示。

## 七、主要 API 一览

系统接口的权限校验通过请求头 `X-User-Role` 携带角色标识（取值为 `business`、`finance`、`admin`）实现，前端已在 `src/api/http.js` 中统一封装，业务代码调用时传入 `withRole: true` 即可自动附加，无需手工设置。

| 方法 | 路径 | 说明 | 角色要求 |
|---|---|---|---|
| GET | `/health` | 健康检查，用于验证服务是否可用 | 无 |
| GET | `/workflow/public/stats` | 公开统计数据 | 无 |
| POST | `/api/auth/register` | 注册用户（自动生成链上账户） | 无 |
| POST | `/api/auth/login` | 用户登录 | 无 |
| GET | `/api/auth/users` | 用户列表查询 | 无 |
| POST | `/api/enterprise/register` | 企业注册上链 | business/admin |
| GET | `/api/enterprise` | 企业列表查询 | 无 |
| GET | `/api/enterprise/:address` | 按钱包地址查询企业 | 无 |
| POST | `/api/receivable/issue` | 签发应收账款 | business/admin |
| POST | `/api/receivable/confirm` | 应收账款确权 | business/admin |
| POST | `/api/receivable/finance` | 应收账款融资放款 | finance/admin |
| POST | `/api/receivable/settle` | 应收账款结算 | business/admin |
| GET | `/api/receivable` | 应收账款列表 | 无 |
| GET | `/api/receivable/:id` | 应收账款详情 | 无 |
| GET | `/api/receivable/:id/financing` | 某笔应收账款的融资记录 | 无 |
| GET | `/api/dashboard/overview` | 看板总览数据 | 无 |
| GET | `/api/tx-logs` | 链上交易日志 | 无 |
| GET | `/api/tx-logs/:hash/analyze` | AI 分析指定交易哈希 | business/finance/admin |
| GET | `/api/chain/stats` | 链上统计信息 | business/finance/admin |
| POST | `/api/ai/chat` | AI 对话问答 | business/finance/admin |
| POST | `/api/ipfs/upload` | 上传文件到 IPFS 并存证 | business/finance/admin |
| GET | `/api/db/dashboard/overview` | 数据库维度的看板统计 | 无 |
| GET | `/api/db/enterprises` | 数据库企业列表 | 无 |
| GET | `/api/db/receivables` | 数据库应收账款列表 | 无 |
| GET | `/api/db/receivables/:id` | 数据库应收账款详情 | 无 |
| GET | `/api/db/receivables/:id/financing` | 数据库融资记录 | 无 |
| GET | `/api/db/documents` | 文档存证列表 | 无 |

> 此外还有 `/api/business/*` 与 `/api/finance/*` 两组按角色聚合的路由，功能与上面的 `/api/receivable/*`、`/api/enterprise/*` 等价，只是路径前缀不同。

## 八、常见问题排查

| 现象 | 原因与解决办法 |
|---|---|
| 后端启动即退出，提示 `client.Dial failed` | FISCO-BCOS 节点 Channel 端口不可达。检查 `config.toml` 的 `NodeURL` 与 `GroupID` 是否正确、`sdk/` 下三份证书是否与节点配套、节点所在服务器防火墙是否放行该端口 |
| 后端启动即退出，提示 `Load app config failed` | 启动目录不对。**必须先 `cd backend2` 再运行**，因为配置文件按相对路径读取 |
| 启动报 `dial tcp ...:3306` 或数据库相关错误 | MySQL 服务未启动，或 `[Database]` 段账号密码/主机端口配置有误；同时确认已执行 `CREATE DATABASE finance` |
| 后端启动报 `invalid contract address` | `config.toml` 中 `[Contract] address` 不是合法的十六进制地址，或该合约并未部署在当前群组 |
| 前端登录或请求提示网络错误 | 后端未启动或不在 8888 端口。开发模式下还需确认 `vite.config.js` 的代理配置未被改动 |
| 文件上传报错 | 本机未运行 IPFS 节点。需先启动 `ipfs daemon` 并确保 API 监听在 `localhost:5001` |
| AI 助手无响应或返回错误 | `[AI]` 段的 `api_key` 无效或已过期；也可能是服务器无法访问火山方舟 API，检查网络与超时配置 |
| 生产模式下刷新子页面出现 404 | 后端需能找到前端构建产物。确认 `frontend/dist/index.html` 存在，且 `config.toml` 中 `frontend_dist` 相对路径正确；必要时重新执行 `npm run build` |
| 提示端口已被占用 | 修改 `config.toml` 的 `listen`（后端）或 `vite.config.js` 的 `server.port`（前端），注意同步更新代理目标地址 |
| 操作提示 `role not allowed` | 当前登录用户的角色不具备该操作的权限。例如融资放款只有 finance/admin 能做，签发与结算只有 business/admin 能做 |
| 操作提示 `missing role header: X-User-Role` | 请求未携带角色请求头。前端登录后会自动附加，若直接调用接口需手工添加该请求头 |

## 九、目录结构速览

```
SupplyChainFinance/
├── backend2/                     # Go 后端服务
│   ├── main.go                   # 入口：加载配置 → 初始化DB → 注册路由 → 托管前端
│   ├── config.toml               # 核心配置文件（★ 启动前务必逐项核对）
│   ├── appconfig/                # 配置加载与 DSN 拼接
│   ├── dbs/                      # 数据库初始化、表结构定义与数据访问
│   │   ├── schema.go             # 自动建表与字段迁移逻辑
│   │   ├── finance.sql           # 完整表结构 SQL（可手工执行）
│   │   └── finance-data.sql      # 演示数据 SQL（含 3 个演示账号）
│   ├── router/                   # 业务路由
│   │   ├── sign.go               # 用户注册/登录、企业注册、应收账务操作
│   │   ├── bcos.go               # 区块链客户端初始化与链上交互
│   │   ├── ipfs.go               # 文件上传至 IPFS
│   │   ├── ai.go                 # 豆包大模型调用
│   │   ├── db_query.go           # 数据库查询类接口
│   │   └── routes.go             # 路由总注册
│   ├── middleware/               # CORS、访问统计、角色校验中间件
│   ├── sdk/                      # FISCO-BCOS Channel 通信证书（★ 敏感）
│   ├── keystore/                 # 注册用户时自动生成的链上账户（★ 需备份）
│   ├── ci/                       # 交易签名账户私钥 .pem（★ 需妥善保管）
│   ├── SupplyChainFinance/       # 智能合约源码 + ABI + Go 绑定代码
│   ├── uploadsIPFS/              # 上传文件的本地缓存副本
│   ├── web/                      # 兜底静态页面（前端未构建时使用）
│   ├── gin-runner.log            # 运行日志文件
│   └── backend.exe               # Windows 平台已编译的可执行文件
└── frontend/                     # Vue3 前端
    ├── src/api/http.js           # fetch 封装，自动附加 X-User-Role 请求头
    ├── src/router/index.js       # 路由表与登录/角色守卫
    ├── src/views/                # 各功能页面组件
    ├── src/components/           # 侧边栏、AI 助手等公共组件
    ├── src/stores/               # Pinia 状态管理
    ├── vite.config.js            # Vite 配置（开发端口 5173 与后端代理）
    ├── index.html                # 前端入口 HTML
    └── dist/                     # 构建产物，可由后端直接托管
```

## 十、注意事项

1. **密钥与证书资产需妥善保管**：`keystore/` 目录中保存着注册用户时生成的链上账户密钥文件，`ci/*.pem` 是后端发起交易时的签名私钥，`sdk/` 目录是与区块链节点通信的证书。这些文件一旦丢失，对应的链上资产与操作能力将无法恢复；一旦泄露，他人可冒用身份发起交易。迁移或备份部署环境时，请务必一并处理且不要外传。
2. **重新部署合约后需同步更新配置**：如果更换了区块链网络或重新部署了 `SupplyChainFinance` 合约，必须更新 `config.toml` 中 `[Contract]` 的地址，并确认 `backend2/SupplyChainFinance/` 中的 ABI 与 Go 绑定代码与新合约接口保持一致，否则调用会因方法不存在而失败。
3. **合约地址与群组强相关**：`GroupID` 决定后端连接哪个群组，合约只存在于它被部署的那个群组中。切换群组时，除了修改 `GroupID`，也要确认合约地址在该群组中确实存在。
4. **Windows 下运行的注意点**：`go run .` 与直接运行 `backend.exe` 效果相同，但无论哪种方式，工作目录都必须是 `backend2/`，否则相对路径形式的证书、私钥、配置文件都会找不到。
5. **日志文件会持续增长**：`gin-runner.log` 以追加方式记录所有请求，长时间运行后体积会变大，建议定期清理或在非调试场景下降低日志级别。
