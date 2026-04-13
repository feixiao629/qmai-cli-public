---
name: qmai-finance
version: 1.0.0
description: "财务服务：分账明细、支付账单、营业统计、商品销售汇总"
metadata:
  bins: [qmai]
  help: "qmai finance --help"
---

> 前置条件: 请先阅读 ../qmai-shared/SKILL.md 了解认证和配置

## 核心概念

- **statement**: 分账明细、微信账单、支付宝账单、微信账单下载地址
- **stats**: 营业统计、订单渠道、结账方式、商品销售汇总

## 命令概览

### 对账与账单
```bash
qmai finance statement split-flows --created-at-start "2026-04-01 00:00:00" --created-at-end "2026-04-07 23:59:59" --page 1 --page-size 20
qmai finance statement wechat-bills --bill-date 2026-04-07 --page 1 --page-size 20
qmai finance statement alipay-bills --bill-date 2026-04-07 --page 1 --page-size 20
qmai finance statement wechat-bill-url --bill-date 2026-04-07 --shop-code S001
```

### 统计与字典
```bash
qmai finance stats business-summary --shop-code S001 --start-date 2026-04-01 --end-date 2026-04-07
qmai finance stats order-types
qmai finance stats settle-scenes
qmai finance stats item-turnover --shop-code S001 --start-date 2026-04-01 --page 1 --page-size 20
```

## 开放平台 API 端点

| 操作 | 端点 |
| --- | --- |
| 查询分账明细 | POST /v3/payCenter/split/orderFlow |
| 查询微信支付账单 | POST /v3/payCenter/custom/wechatBill |
| 查询支付宝支付账单 | POST /v3/payCenter/offlinePayments/pullAlipayBillData |
| 查询营业统计数据 | POST /v3/dataone/finance/summary/businessRecord |
| 查询订单渠道 | POST /v3/dataone/finance/detail/orderType |
| 查询订单结账方式 | POST /v3/dataone/finance/detail/settleScene |
| 查询微信支付账单URL | POST /v3/pay/getWechatMerchantBillUrl |
| 查询门店商品销售汇总 | POST /v3/dataone/item/store/turnover |

## 注意事项（Agent 必读）

- `alipay-bills` 的官方文档没有给出明确业务字段明细，CLI 先保留原始记录列；需要完整结构时优先用 `--format json`
- `split-flows`、`wechat-bills`、`business-summary`、`item-turnover` 当前都是显式分页，不会自动拉全量
- `item-turnover` 的 `page-size` 官方最大支持 200
- `7.2 交易支付` 和 `7.3 门店账户` 在官方文档仍标记待上线，当前只保留在覆盖矩阵，不进入正式命令面
