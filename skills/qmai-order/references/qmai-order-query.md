# qmai-order 查询参考

## 适用场景

- 按用户查询订单列表
- 查看单笔订单详情和状态
- 拉取储值消费订单和退款单
- 判断会员是否属于“下过单用户”
- 查询制作单记录

## 推荐命令

```bash
qmai order query user-orders --order-at-start "2026-04-01 00:00:00" --order-at-end "2026-04-07 23:59:59" --size 20 --user-id 1001
qmai order query detail --biz-type 5 --order-no O20260407001
qmai order query status --order-no O20260407001
qmai order query member-ordered --biz-type 5 --user-id 1001
qmai order query production-records --user-id 1001 --order-no O20260407001
```

## 排查建议

- 查询不到订单时，先核对时间范围格式是否符合接口要求
- `detail` 查不到时，先确认 `biz-type` 是否与该订单所属业务一致
- 储值订单接口时间参数使用 `yyyyMMdd`，不要误传完整时间戳
- 若需要完整原始字段，优先加 `--format json`
