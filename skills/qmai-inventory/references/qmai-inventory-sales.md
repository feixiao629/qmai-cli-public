# qmai-inventory 销售管理参考

## 适用场景

- 报货单接单或拒单
- 查询订货单、调拨单、退货单
- 创建、完成、发货订货单

## 推荐命令

```bash
qmai inventory sales declare-detail --declare-no BH20260407001
qmai inventory sales require-list --page 1 --page-size 10
qmai inventory sales transfer-list --page 1 --page-size 10
qmai inventory sales return-list --page 1 --page-size 10
qmai inventory sales return-examine --from-json return-examine.json --dry-run
qmai inventory sales return-receipt --from-json return-receipt.json --dry-run
qmai inventory sales require-create --from-json require-create.json --dry-run
```

## 排查建议

- 高副作用动作优先保留 JSON 请求体，便于审计和重试
- `require-list`、`return-list`、`transfer-list` 默认只查单页
- 调拨详情需使用 `transfer-no`，不是内部 id
