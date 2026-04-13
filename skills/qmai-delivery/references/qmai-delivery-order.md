# qmai-delivery 配送订单参考

## 适用场景

- 新建聚合配送单
- 查询配送详情和骑手位置
- 取消单笔或全部配送单

## 推荐命令

```bash
qmai delivery order create --from-json delivery-create.json --dry-run
qmai delivery order detail --order-source 1 --origin-order-no O20260407001
qmai delivery order rider-location --order-source 1 --origin-order-no O20260407001
qmai delivery order cancel --multi-mark S001 --order-source 1 --origin-order-no O20260407001 --dry-run
qmai delivery order cancel-all --order-source 1 --origin-order-no O20260407001 --dry-run
```

## 排查建议

- 创建配送单前先保留原始 JSON，请求失败时便于比对字段缺失
- `detail` 和 `rider-location` 都依赖 `order-source + origin-order-no`
- 如果同一笔单存在主单/补发单/拆单，补 `--order-flag`
