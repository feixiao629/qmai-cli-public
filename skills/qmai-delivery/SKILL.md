---
name: qmai-delivery
version: 1.0.0
description: "聚合配送：创建、取消、查询配送订单与同步配送状态"
metadata:
  bins: [qmai]
  help: "qmai delivery --help"
---

> 前置条件: 请先阅读 ../qmai-shared/SKILL.md 了解认证和配置

## 核心概念

- **order**: 创建配送单、取消配送单、查详情、查骑手位置、取消全部配送单
- **status**: 同步商家自配送状态
- **create-and-update**: 聚合接口，一次提交配送单创建和状态更新

## 命令概览

### 配送订单
```bash
qmai delivery order create --from-json delivery-create.json --dry-run
qmai delivery order cancel --multi-mark S001 --order-source 1 --origin-order-no O20260407001 --dry-run
qmai delivery order detail --order-source 1 --origin-order-no O20260407001
qmai delivery order rider-location --order-source 1 --origin-order-no O20260407001
qmai delivery order cancel-all --order-source 1 --origin-order-no O20260407001 --dry-run
qmai delivery order create-and-update --from-json delivery-update.json --dry-run
```

### 状态同步
```bash
qmai delivery status update --origin-order-no O20260407001 --order-status 5 --dry-run
```

## 开放平台 API 端点

| 操作 | 端点 |
| --- | --- |
| 创建配送订单 | POST /v3/delivery/createDeliveryOrder |
| 取消配送订单 | POST /v3/delivery/cancelDeliveryOrder |
| 查询配送单详情 | POST /v3/delivery/getDeliveryOrderInfo |
| 查询骑手当前位置 | POST /v3/delivery/getRiderLocation |
| 更新订单配送状态 | POST /v3/delivery/updateSelfOrderStatus |
| 取消全部配送单 | POST /v3/delivery/cancelAllDeliveryOrder |
| 创建配送单并更新配送状态 | POST /v3/delivery/updateDeliveryStatus |

## 注意事项（Agent 必读）

- `create` 和 `create-and-update` 请求体很大，统一用 `--from-json`
- `cancel`、`cancel-all`、`status update` 都支持 `--dry-run`
- `create-and-update` 官方文档明确说明不能和 `create + updateSelfOrderStatus` 混用
- `order-source` 和 `order-status` 都是平台枚举值，传值前要先按企迈文档核实
