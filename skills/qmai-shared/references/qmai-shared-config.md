# 配置管理详情

## 配置文件结构

```yaml
# ~/.config/qmai/config.yaml
active_profile: default
default_format: table     # json | table | csv
debug: false

profiles:
  default:
    name: default
    shop_code: "S001"
    open_id: "d14c1559e87b747d577c834b275a4310"
    grant_code: "ba67d4fa46"
    base_url: ""          # 空则使用默认值

  test-store:
    name: test-store
    shop_code: "S002"
    open_id: "..."
    grant_code: "..."
    base_url: "https://openapi.qmai.co/"  # 沙箱环境
```

## 初始化流程

```bash
$ qmai config init
qmai CLI 配置初始化
------------------------------
Profile 名称 (default): my-store
门店编码 shopCode (可选): S001
API Base URL (https://openapi.qmai.cn/):
默认输出格式 [table/json/csv] (table): json

✓ 配置已保存到 ~/.config/qmai/config.yaml
  Active profile: my-store

下一步: 运行 'qmai auth login' 配置开放平台凭证
```

## 配置键列表

| 键 | 说明 | 默认值 |
|---|------|--------|
| `active_profile` | 当前活动 profile | `default` |
| `default_format` | 默认输出格式 | `table` |
| `debug` | 调试模式 | `false` |
| `base_url` | API 基础 URL（per profile） | `https://openapi.qmai.cn/` |
| `shop_code` | 门店编码（per profile） | - |
| `open_id` | 开放平台 openId（per profile） | - |
| `grant_code` | 开放平台 grantCode（per profile） | - |

## Profile 管理

```bash
# 添加 profile
qmai config profile add new-store --shop-code S003 --base-url https://...

# 列出所有 profiles
qmai config profile list
# * default (shop: S001)
#   test-store (shop: S002) [https://openapi.qmai.co/]

# 删除 profile
qmai config profile remove test-store
```
