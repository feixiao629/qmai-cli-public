# qmai-member asset 参考

## 建议使用顺序

1. 先用 `balance` / `points` 看会员当前资产
2. 再用 `balance-flow` / `points-flow` 查资产变化轨迹
3. 需要看券时再用 `coupons` 或 `coupon-list`
4. 资产写操作统一先走 `--dry-run`，确认 JSON 结构后再正式执行

## 示例

```bash
qmai member asset points --customer-id C001
qmai member asset balance --customer-id C001
qmai member asset coupons --customer-id C001 --use-status 1
qmai member asset coupon-list --customer-id C001 --coupon-type 2 --page 1 --page-size 20
qmai member asset coupon-detail-list --customer-id C001 --store-id 1001 --use-status 0 --page 1 --page-size 20
qmai member asset balance-flow --customer-id C001 --start-time 2026-04-01 00:00:00 --end-time 2026-04-07 23:59:59
qmai member asset points-flow --customer-id C001 --change-types REWARD,CONSUME --page 1 --page-size 20
qmai member asset deposit-rules --shop-code S001
qmai member asset inflate-status --customer-id C001
qmai member asset recharge --from-json recharge.json --dry-run
qmai member asset recharge-reverse --from-json recharge-reverse.json --dry-run
qmai member asset offline-balance-op --from-json offline-balance-op.json --dry-run
qmai member asset balance-debit --from-json balance-debit.json --dry-run
qmai member asset points-add --from-json points-add.json --dry-run
qmai member asset consume --from-json consume.json --dry-run
qmai member asset coupon-writeoff --from-json coupon-writeoff.json --dry-run
```

## 输出说明

- 表格模式默认显示核心字段，适合巡检
- `--format json` 输出完整响应结构，适合排查字段缺失和对照文档
- 写接口默认只回显提交结果，复杂入参请放在 JSON 文件中维护
