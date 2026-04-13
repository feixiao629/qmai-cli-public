# qmai-inventory 基础资料参考

## 适用场景

- 创建盘点调整单
- 审核调拨申请单
- 新增或修改物品、分类、单位
- 新增或修改供应商
- 批量下发品项

## 推荐命令

```bash
qmai inventory master inventory-adjust --from-json inventory-adjust.json --dry-run
qmai inventory master transfer-audit --from-json transfer-audit.json --dry-run
qmai inventory master product-create --from-json product-create.json --dry-run
qmai inventory master category-create --from-json category-create.json --dry-run
qmai inventory master unit-create --from-json unit-create.json --dry-run
qmai inventory master supplier-create --from-json supplier-create.json --dry-run
qmai inventory master product-distribute --from-json product-distribute.json --dry-run
qmai inventory master machining-card-batch-create --from-json machining-card-batch-create.json --dry-run
```

## 排查建议

- 物品和供应商接口字段很多，优先维护可复用 JSON 模板
- `unit-create` 返回的是 `unitCode`，不是业务 id
- `changeProduct` 与 `update` 字段结构相近，但新增与修改的必填项不同
- 半成品成本卡的 `list` 为批量结构，单次最多内容较多时优先先跑 `--dry-run`
