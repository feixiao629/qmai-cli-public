---
name: qmai-queue
version: 1.0.0
description: "排队服务：查询门店排队进度、订单排队进度和叫号列表"
metadata:
  bins: [qmai]
  help: "qmai queue --help"
---

> 前置条件: 请先阅读 ../qmai-shared/SKILL.md 了解认证和配置

## 核心概念

- **shop-progress**: 批量查询门店当前排队进度
- **order-progress**: 查询单笔订单的排队进度
- **shop-queue-nos**: 查询门店叫号列表

## 命令概览

```bash
qmai queue shop-progress --shop-type 1 --shop-ids 1001,1002
qmai queue order-progress --order-no O20260407001
qmai queue shop-queue-nos --shop-code S001 --page 1 --size 10
```

## 开放平台 API 端点

| 操作 | 端点 |
| --- | --- |
| 查询门店排队进度 | POST /v3/queuing/queryShopQueueCup |
| 查询订单排队进度 | POST /v3/queuing/queryOrderQueueCup |
| 查询门店排队叫号列表 | POST /v3/queuing/queueNo/queryShopQueueNoList |

## 注意事项（Agent 必读）

- `shop-progress` 的 `shopIdList` 官方最多支持 20 个
- `shop-queue-nos` 的 `size` 官方限制在 1 到 20 之间
- 这个模块目前只有读接口，没有写操作
