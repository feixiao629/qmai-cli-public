# qmai-cli

企迈开放平台命令行工具，面向商品、门店组织、会员、营销、订单、财务、聚合配送、进销存和排队等业务场景，提供统一的命令行调用入口，并支持通过 npm 全局安装。

为什么选 qmai-cli？ · 工具能力项 · 安装与快速开始 · Agent Skills · 常用命令 · 认证与配置 · 安全与声明

## 为什么选 qmai-cli？

- 统一命令入口：围绕企迈开放平台常见业务域提供一致的 CLI 结构
- 适合自动化：支持脚本调用、批量处理和结构化输出
- 凭证分层存储：敏感密钥走系统 keychain，非敏感配置走本地 profile
- 面向实际运维和联调：支持查询、变更、批量导入导出和原始 API 调试

## 工具能力项

### 商品

- 商品列表、详情查询
- 商品创建、更新、删除
- 商品上下架、估清、售罄、库存相关操作
- 批量调价、批量状态变更
- 商品导入、导出

### 门店与组织

- 门店详情、门店列表查询
- 门店编码与门店 ID 转换
- 组织树、团队列表、部门树查询
- 门店状态修改、门店归属团队调整
- 门店资料同步

### 会员

- 积分、余额、优惠券等资产查询
- 会员检索、基础资料查询
- openId、动态码、条件查询等辅助能力
- 标签查询、标签创建、会员打标
- 冻结、注销、换绑手机号等会员操作入口

### 营销

- 优惠券状态、营销活动查询
- 礼品卡信息、发卡、核销、兑换等操作
- 优惠券模板启用、发券活动相关能力
- 营销定价确认

### 订单

- 用户订单、订单详情、订单状态查询
- 评价回复
- 订单上报、已完成订单上报、已退款订单上报

### 财务

- 分账明细查询
- 微信账单与账单下载地址查询
- 营业汇总、商品销售汇总等统计能力

### 聚合配送

- 配送单创建、取消、详情查询
- 骑手位置查询
- 配送状态同步

### 进销存

- 订货、调拨、退货单查询与处理
- 入库、出库、库存占用等库存操作
- 商品、供应商、半成品成本卡等基础资料维护

### 排队

- 门店排队进度查询
- 订单排队进度查询
- 门店叫号列表查询

### 平台调试与辅助能力

- 原始 API 透传
- 环境诊断
- shell 自动补全

### Agent Skills（AI 助手）

本仓库在 `skills/` 下提供与业务域对应的 Agent Skills（如 `qmai-product`、`qmai-order` 等），便于在 Cursor、Codex 等支持 Agent Skills 的环境中复用 `qmai` 子命令与参数说明。安装 CLI 后，可用 [Skills CLI](https://github.com/vercel-labs/skills) 从本仓库拉取技能包：

```bash
# 安装本仓库中的全部技能（需 Node.js 18+）
npx skills add feixiao629/qmai-cli-public

# 仅安装某一技能
npx skills add feixiao629/qmai-cli-public --skill qmai-product

# 使用完整 Git URL 亦可
npx skills add https://github.com/feixiao629/qmai-cli-public
```

全局安装到当前用户目录时可按需加 CLI 支持的参数（例如部分环境为 `-g`）。技能安装位置取决于所用工具（常见为项目或用户下的 `.cursor/skills/`、`.agents/skills/` 等），请以 Skills CLI 与编辑器文档为准。

## 安装与快速开始

### 环境要求

- Node.js `18+`
- 可访问企迈开放平台的合法账号与授权信息
- 本机支持系统 keychain 或等价安全存储

### 通过 npm 安装

```bash
npm install -g qmai-cli-public
qmai version
```

说明：

- npm 安装会自动下载当前平台对应的预编译二进制
- 普通使用者不需要额外安装 Go

### 从源码构建

源码安装需要：

- Go `1.23.0+`

```bash
make build
./bin/qmai version
```

### 从源码安装

```bash
make install
qmai version
```

### 初始化与认证

```bash
qmai config init
qmai auth login
qmai doctor
```

## 常用命令

### 商品

```bash
qmai product list
qmai product get <id|name>
qmai product create --name "美式咖啡" --price 28.00 --trade-no P001
qmai product batch-status --action up --trade-marks TM001,TM002
```

### 门店与组织

```bash
qmai store get <shop-code>
qmai store list --keyword 门店名
qmai store org-tree
qmai store team-list
```

### 会员

```bash
qmai member asset points --customer-id C001
qmai member profile search --phone 13800138000
qmai member tag list --customer-id C001
qmai member profile freeze --customer-id C001 --reason "风险处理" --dry-run
```

### 订单与财务

```bash
qmai order query status --order-no O20260407001
qmai order report upload --from-json order-upload.json
qmai finance statement split-flows --created-at-start "2026-04-01 00:00:00" --created-at-end "2026-04-07 23:59:59"
qmai finance stats business-summary --shop-code S001 --start-date 2026-04-01 --end-date 2026-04-07
```

### 配送、进销存与排队

```bash
qmai delivery order detail --order-source 1 --origin-order-no O20260407001
qmai inventory stock realtime-list --warehouse-nos CK001 --page 1 --page-size 20
qmai queue order-progress --order-no O20260407001
```

### 原始 API 与补全

```bash
qmai api v3/goods/item/getItemList --body '{"shopCode":"S001"}'
qmai completion zsh
```

## 认证与配置

- 配置文件：`~/.config/qmai/config.yaml`
- 敏感信息：`openKey` 存储在操作系统 keychain 中
- 非敏感信息：`openId`、`grantCode`、`shopCode` 存在 profile 配置中
- 全局参数：`--profile`、`--format`、`--debug`

常用配置命令：

```bash
qmai config list
qmai config profile add store2 --shop-code S002
qmai config profile list
qmai config set active_profile store2
```

## 安全与声明

- 请使用你自己的企迈开放平台凭证，并确保拥有目标资源的合法访问权限
- 本仓库不提供真实凭证、真实门店信息或内部文档摘录
- 示例参数仅用于演示，不代表真实生产数据
- 这不是企迈官方项目
- `Qmai` 及相关产品名称归其各自权利人所有
- 使用本软件时，你需要自行遵守对应平台的服务条款、接口文档和授权要求
