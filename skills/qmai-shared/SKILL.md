---
name: qmai-shared
version: 3.0.0
description: "基础能力：开放平台认证、配置管理、安全规则、通用约定"
metadata:
  bins: [qmai]
  help: "qmai --help"
---

## 概述

qmai CLI 是门店经营命令行工具。通过企迈开放平台 (`openapi.qmai.cn`) 管理门店商品、门店与组织、会员、营销、订单、财务、聚合配送、进销存和排队。本 Skill 描述认证、配置和安全相关的基础能力。

## 认证流程

### 凭证配置

```bash
qmai auth login [--profile <name>]
```

流程：
1. 交互式输入 openId、grantCode、openKey、shopCode
2. openKey 存入 OS keychain（敏感凭证）
3. openId/grantCode/shopCode 写入 config profile
4. 设为 active_profile

### 认证管理命令

| 命令 | 说明 |
|------|------|
| `qmai auth login` | 配置开放平台凭证 |
| `qmai auth logout` | 清除凭证（keychain + config） |

### 凭证存储

- **openKey**: macOS Keychain（service: `qmai-cli`），Linux: libsecret
- **openId/grantCode/shopCode**: `~/.config/qmai/config.yaml` profile 中

### 请求签名

每个 API 请求自动签名，无需手动操作：
1. 生成随机 nonce + 当前时间戳
2. 按字典序拼接 `grantCode=xx&nonce=xx&openId=xx&timestamp=xx`
3. HmacSHA1(拼接串, openKey) → Base64 → URL Encode → token
4. 请求体: `{openId, grantCode, nonce, timestamp, token, params}`

## 配置管理

### 配置文件

位置：`~/.config/qmai/config.yaml`

```yaml
active_profile: default
default_format: table
debug: false
profiles:
  default:
    name: default
    shop_code: "S001"
    open_id: "d14c1559..."
    grant_code: "ba67d4fa46"
  store2:
    name: store2
    shop_code: "S002"
    open_id: "..."
    grant_code: "..."
```

### 配置优先级

flags → profile 配置 → 默认值

### 配置命令

```bash
qmai config init                                    # 交互式初始化
qmai config set <key> <val>                         # 设置配置项
qmai config get <key>                               # 读取配置项
qmai config list                                    # 查看所有配置
qmai config profile add <name> --shop-code xxx      # 添加 profile
qmai config profile remove <name>                   # 删除 profile
qmai config profile list                            # 列出 profiles
```

## 输出格式

所有命令支持 `--format` flag：

| 格式 | 说明 |
|------|------|
| `table` | 对齐的文本表格（默认，终端友好） |
| `json` | JSON 输出（脚本友好） |
| `csv` | CSV 输出（Excel/导入友好） |

## Debug 模式

```bash
qmai product list --debug
```

启用后在 `Client.Call` 层输出完整日志（一行包含全部信息）：
```
[DEBUG] POST <url> | req=<完整请求JSON> | resp=<完整响应JSON> | <耗时>ms
```

用于排查 API 对接问题，查看实际请求参数和响应结构。

## 安全规则（Agent 必读）

1. **openKey 安全**: openKey 仅存储在 OS keychain，不写入配置文件或日志
2. **不存储敏感信息**: config 中不保存 openKey
3. **Dry-run 优先**: 所有变更操作建议先使用 `--dry-run` 预览
4. **确认机制**: 删除/下架操作需要 `--force` 确认
5. **批量操作安全**: 批量操作前务必备份（`export`），建议先 `--dry-run`

## API 基础

- Base URL: `https://openapi.qmai.cn/`
- 协议: 全部 HTTP POST，Content-Type: application/json
- QPS 限制: 单接口最高 10
- 分页限制: pageSize 最大 50
- 认证: 请求体签名（自动注入）
- 响应格式:
  ```json
  {
    "status": true,
    "code": 0,
    "message": "请求成功",
    "data": { ... },
    "traceId": "..."
  }
  ```

## 可用命令

```
qmai
├── auth        认证管理
├── config      配置管理
├── product     商品管理
├── store       门店与组织管理
├── member      会员服务
├── marketing   营销服务
├── order       订单服务
├── finance     财务服务
├── delivery    聚合配送
├── inventory   进销存
├── queue       排队服务
├── api         Raw API 请求透传
├── doctor      诊断检查
├── completion  Shell 自动补全
└── version     版本信息
```
