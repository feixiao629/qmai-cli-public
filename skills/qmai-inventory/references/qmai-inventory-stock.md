# qmai-inventory 库存管理参考

## 适用场景

- 批量查询入库单和出库单
- 查询门店库存台账
- 查询仓库实时库存

## 推荐命令

```bash
qmai inventory stock inbound-create --from-json inbound-create.json --dry-run
qmai inventory stock inbound-list --created-start-at "2026-04-01 00:00:00" --created-end-at "2026-04-07 23:59:59"
qmai inventory stock inbound-update --from-json inbound-update.json --dry-run
qmai inventory stock inbound-finish --from-json inbound-finish.json --dry-run
qmai inventory stock outbound-create --from-json outbound-create.json --dry-run
qmai inventory stock outbound-list --created-start-at "2026-04-01 00:00:00" --created-end-at "2026-04-07 23:59:59"
qmai inventory stock outbound-update --from-json outbound-update.json --dry-run
qmai inventory stock occupy --from-json stock-occupy.json --dry-run
qmai inventory stock release --from-json stock-release.json --dry-run
qmai inventory stock store-ledger --start-date 2026-04-01 --end-date 2026-04-07 --page 1 --page-size 10
qmai inventory stock realtime-list --warehouse-nos CK001 --page 1 --page-size 20
```

## 排查建议

- 入库单和出库单查询都要求创建时间范围
- `realtime-list` 的 `page-size` 受官方接口限制，最大 100
- 如果库存类数据需要完整字段，优先加 `--format json`
