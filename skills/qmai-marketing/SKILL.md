---
name: qmai-marketing
version: 1.1.0
description: "营销服务：券查询与发券、活动查询与撤销、礼品卡查询与写入、结算页算价"
metadata:
  bins: [qmai]
  help: "qmai marketing --help"
---

> 前置条件: 请先阅读 ../qmai-shared/SKILL.md 了解认证和配置

## 核心概念

- **coupon**: 券状态、券详情、券模板、不记名券、筛券、发券和模板启用
- **campaign**: 发券活动、撤销发券、兑换码下发/作废、兑换码状态、任务活动参与记录和任务领取
- **gift-card**: 卡信息、消费明细、卡模板、会员卡列表，以及单卡/批量余额扣减、冲正和发放
- **pricing**: 门店促销活动查询与结算页算价

## 命令概览

### 券管理
```bash
qmai marketing coupon status --id CPN001
qmai marketing coupon detail --user-coupon-code CPN001
qmai marketing coupon template --id 1001
qmai marketing coupon template-enable --id 1001 --dry-run
qmai marketing coupon template-batch --ids 1001,1002
qmai marketing coupon anonymous --code ANON001
qmai marketing coupon template-by-third-code --seller-type 2 --third-biz-code EXT001
qmai marketing coupon choose --from-json choose-coupon.json
qmai marketing coupon grant-activity --from-json grant-activity.json --dry-run
qmai marketing coupon grant-activity-async --from-json grant-activity-async.json --dry-run
qmai marketing coupon grant-template --from-json grant-template.json --dry-run
```

### 活动管理
```bash
qmai marketing campaign list --channel-id 1 --page 1 --page-size 20
qmai marketing campaign revoke-grant --from-json revoke-grant.json --dry-run
qmai marketing campaign exchange-dispatch --from-json exchange-dispatch.json --dry-run
qmai marketing campaign exchange-disable --from-json exchange-disable.json --dry-run
qmai marketing campaign exchange-status --code EXCHANGE001
qmai marketing campaign task-records --activity-ids 1001,1002 --customer-id 2001
qmai marketing campaign task-claim --from-json task-claim.json --dry-run
qmai marketing campaign recycle-coupons --from-json recycle-coupons.json --dry-run
```

### 卡管理
```bash
qmai marketing gift-card info --card-no CARD001 --take-card-template 1
qmai marketing gift-card flow --card-no CARD001 --page 1 --page-size 20
qmai marketing gift-card template --id 3001
qmai marketing gift-card template-batch --ids 3001,3002
qmai marketing gift-card list --customer-id 2001 --shop-code S001
qmai marketing gift-card consume --from-json gift-card-consume.json --dry-run
qmai marketing gift-card consume-batch --from-json gift-card-consume-batch.json --dry-run
qmai marketing gift-card reverse --from-json gift-card-reverse.json --dry-run
qmai marketing gift-card reverse-batch --from-json gift-card-reverse-batch.json --dry-run
qmai marketing gift-card part-reverse --from-json gift-card-part-reverse.json --dry-run
qmai marketing gift-card issue --from-json gift-card-issue.json --dry-run
qmai marketing gift-card report-loss --from-json gift-card-report-loss.json --dry-run
qmai marketing gift-card relieve-loss --from-json gift-card-relieve-loss.json --dry-run
qmai marketing gift-card recycle --from-json gift-card-recycle.json --dry-run
qmai marketing gift-card exchange --from-json gift-card-exchange.json --dry-run
qmai marketing gift-card grant-template --from-json gift-card-grant-template.json --dry-run
```

### 营销计算
```bash
qmai marketing pricing activities --from-json promotion-activities.json
qmai marketing pricing confirm --from-json confirm.json
```

## 注意事项（Agent 必读）

- `coupon choose`、发券命令和 `pricing activities/confirm` 入参复杂，统一用 `--from-json`
- `pricing confirm` 目前默认表格只展示价格摘要，若要看完整算价结果请加 `--format json`
- 活动实时发券、异步发券、模板发券、模板启用、撤销发券、兑换码下发/作废、任务领取、回收券，以及礼品卡扣减/冲正/挂失/回收/兑换/发放已经进入正式命令面
- 礼品卡写接口同样统一用 `--from-json`，建议先 `--dry-run`；其中 `part-reverse` 只能冲正统一消费接口，不能和全额冲正混用
- `campaign list`、`gift-card flow` 是显式分页，不会自动拉全量
