# qmai-finance 统计参考

## 适用场景

- 查询单店营业统计
- 查询订单渠道和结账方式字典
- 查询门店商品销售汇总

## 推荐命令

```bash
qmai finance stats business-summary --shop-code S001 --start-date 2026-04-01 --end-date 2026-04-07
qmai finance stats order-types
qmai finance stats settle-scenes
qmai finance stats item-turnover --shop-code S001 --start-date 2026-04-01 --page 1 --page-size 20
```

## 排查建议

- `business-summary` 必须带 `shop-code` 和日期范围
- `item-turnover` 需要 `shop-code` 或 `shop-id` 二选一
- 若要保留完整财务字段，优先用 `--format json`
