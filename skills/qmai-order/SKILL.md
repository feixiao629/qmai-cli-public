---
name: qmai-order
version: 1.1.0
description: "订单服务：订单查询、评价回复、订单上报、批量订单同步、制作单查询"
metadata:
  bins: [qmai]
  help: "qmai order --help"
---

> 前置条件: 请先阅读 ../qmai-shared/SKILL.md 了解认证和配置

## 核心概念

- **query**: 订单列表、订单详情、储值单、订单状态、制作单和新老用户下单判断
- **review**: 门店回复用户评价
- **report**: 订单上报、退款上报、批量已完成订单同步、批量已退款订单同步

## 命令概览

### 订单查询
```bash
qmai order query user-orders --order-at-start "2026-04-01 00:00:00" --order-at-end "2026-04-07 23:59:59" --size 20 --user-id 1001
qmai order query detail --biz-type 5 --order-no O20260407001
qmai order query recharge-orders --shop-code S001 --start-time 20260401 --end-time 20260407 --page 1 --page-size 20
qmai order query recharge-refunds --shop-code S001 --start-time 20260401 --end-time 20260407 --page 1 --page-size 20
qmai order query status --order-no O20260407001
qmai order query member-ordered --biz-type 5 --user-id 1001
qmai order query production-records --user-id 1001 --order-no O20260407001
```

### 评价与上报
```bash
qmai order review reply --order-no O20260407001 --reply-at "2026-04-07 18:00:00" --seller-reply-info "感谢反馈" --dry-run
qmai order report upload --from-json order-upload.json
qmai order report refund-up --from-json refund-up.json
qmai order report completed-batch --from-json completed-batch.json --dry-run
qmai order report refunded-batch --from-json refunded-batch.json --dry-run
```

## 开放平台 API 端点

| 操作 | 端点 |
| --- | --- |
| 查询用户订单列表 | POST /v3/order/getUserOrderList |
| 查询订单详情 | POST /v3/order/getDetail |
| 查询门店储值消费订单 | POST /v3/order/standard/recharge/order |
| 查询门店储值消费退款单 | POST /v3/order/standard/recharge/refundOrder |
| 查询订单状态 | POST /v3/order/status |
| 查询用户是否下过订单 | POST /v3/order/checkMemberOrder |
| 回复订单评价 | POST /v3/order/comment/storeReplyUserComment |
| 订单上报 | POST /v3/bsns/order/orderUpload |
| 退款订单上报 | POST /v3/cy/order/refundOrderUp |
| 提交已完成订单 | POST /v3/bsns/order/orderBatchUpload |
| 提交已退款订单 | POST /v3/bsns/order/refundOrderBatchUpload |
| 查询商品制作单列表 | POST /v3/newPattern/orderCenter/post/order/item/production/record/list |

## 注意事项（Agent 必读）

- `report upload`、`report refund-up`、`report completed-batch`、`report refunded-batch` 的请求体较复杂，统一使用 `--from-json`
- `review reply` 为写操作，默认建议先用 `--dry-run`
- `query user-orders` 当前按接口原生条件查询，不会自动翻页补全历史区间内的全部订单
- `member-ordered` 返回值 `0` 表示新用户，`1` 表示非新用户
- 订单详情接口依赖 `bizType`，通常 `5=新饮食`、`4=新休闲`
- 售后、履约、配送骑手、退款审核等接口仍未进入正式命令面
