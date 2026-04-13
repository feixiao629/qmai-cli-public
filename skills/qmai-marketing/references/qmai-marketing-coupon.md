# qmai-marketing coupon 参考

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
qmai marketing campaign revoke-grant --from-json revoke-grant.json --dry-run
qmai marketing campaign exchange-dispatch --from-json exchange-dispatch.json --dry-run
qmai marketing campaign exchange-disable --from-json exchange-disable.json --dry-run
qmai marketing campaign task-claim --from-json task-claim.json --dry-run
qmai marketing campaign recycle-coupons --from-json recycle-coupons.json --dry-run
```

## 说明

- `choose-coupon.json` 直接按企迈文档请求结构组织
- 若只是查券模板，优先用 `template` / `template-batch`
- 发券与撤销发券都建议先 `--dry-run`，确认 `activityId`、`channelId`、`orderNo` 和会员标识结构
- 兑换码下发与作废一次最多处理的数量要按企迈文档约束，CLI 不额外放宽
