---
name: qmai-store
version: 1.0.0
description: "门店与组织管理：门店详情、门店列表、组织树、门店组、配置查询、门店状态与门店同步"
metadata:
  bins: [qmai]
  help: "qmai store --help"
---

> 前置条件: 请先阅读 ../qmai-shared/SKILL.md 了解认证和配置

## 核心概念

- **shopCode**: 门店编码，很多查询和状态修改命令以此为主键
- **shopId / storeId**: 门店 ID，配置、扩展字段、标签和部分批量查询接口会使用
- **teamId**: 门店组 ID，用于门店组别调整
- **fieldCodes**: 门店配置或品牌配置的编码列表，一次最多建议传 10 个
- **shopInfoSync**: 门店主数据同步接口，字段多、副作用大，建议只通过 JSON 文件驱动

## 命令概览

### 门店详情与列表
```bash
qmai store get <shop-code>
qmai store get-by-id <shop-id>
qmai store list [--keyword 门店名] [--page 1] [--page-size 10]
qmai store id <shop-code>
```

### 外卖映射
```bash
qmai store takeout-map-list --platform-type 4 [--page 1] [--page-size 10]
```

### 门店状态与组别
```bash
qmai store set-status <shop-code> --status 1 --dry-run
qmai store move-team --shop-id 1001 --team-id 2001
```

### 配置、扩展字段与标签
```bash
qmai store config --store-id 1001 --field-codes OrderSettingConfig:autoCancelMinutes
qmai store brand-config --field-codes OrderSettingConfig:autoCancelMinutes
qmai store ext-data --shop-id 1001
qmai store labels --shop-id 1001
```

### 组织结构
```bash
qmai store org-tree [--contain-close 1]
qmai store team-list [--name 华东区]
qmai store dept-tree
```

### 门店同步
```bash
qmai store sync --from-json shop.json --dry-run
qmai store sync --from-json shop.json
```

## 开放平台 API 端点

| 操作 | 端点 |
| --- | --- |
| 根据门店编码查询门店详情 | POST /v3/org/shop/getShopDetail |
| 批量查询门店 | POST /v3/org/shop/getShopList |
| 查询平台外卖门店映射明细 | POST /v3/dist/meTakeoutShopMapPage |
| 设置门店营业状态 | POST /v3/org/shop/changeStatus |
| 门店信息同步 | POST /v3/org/shop/shopInfoSync |
| 根据门店编码查询门店id | POST /v3/org/shop/shopCode2Id |
| 根据ID查询门店详情 | POST /v3/org/shop/getShopDetailById |
| 修改门店的组别 | POST /v3/org/shop/modifyShopTeam |
| 查询门店配置(支持批量) | POST /v3/storeConfig/queryStoreConfigBatch |
| 查询品牌配置(支持批量) | POST /v3/sellerConfig/querySellerConfigBatch |
| 查询门店扩展字段信息 | POST /v3/newPattern/orgCenter/post/shop/get-shop-ext-data |
| 查询门店标签列表 | POST /v3/newPattern/orgCenter/get/shop/label-by-id |
| 查询组织机构列表 | POST /v3/org/shop/getOrgTree |
| 查询门店组列表 | POST /v3/org/shop/shopTeamList |
| 查询门店树结构数据 | POST /v3/org/shop/shopDeptTree |

## 注意事项（Agent 必读）

- `store get` 和 `store id` 用 `shopCode`
- `store config`、`store ext-data`、`store labels` 常用 `storeId` 或 `shopId`
- `set-status` 是变更命令，建议优先 `--dry-run`
- `sync` 字段复杂，优先用 `--from-json`
- `shopInfoSync` 可能不是立即一致，不要假设提交成功后列表立刻更新
- `getShopList` 和 `takeout-map-list` 当前是显式分页，不会自动拉全量
- `2.3 人员管理` 尚未纳入当前 `qmai store` 命令面
