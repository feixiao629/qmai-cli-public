# qmai-marketing card/pricing 参考

```bash
qmai marketing campaign list --channel-id 1
qmai marketing campaign task-records --activity-ids 1001,1002 --customer-id 2001
qmai marketing gift-card info --card-no CARD001
qmai marketing gift-card flow --card-no CARD001 --page 1 --page-size 20
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
qmai marketing pricing activities --from-json promotion-activities.json
qmai marketing pricing confirm --from-json confirm.json --format json
```

## 说明

- `pricing` 相关接口字段很多，建议先保留一份请求 JSON 模板
- `pricing confirm` 用表格模式看摘要，用 `json` 看完整算价细节
- 礼品卡冲正分为全额冲正和部分冲正两种，部分冲正要额外提供 `subBizId`
- 礼品卡挂失/取消挂失一次最多处理数量以开放平台文档为准，CLI 不放宽上限
