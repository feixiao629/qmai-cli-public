---
name: qmai-inventory
version: 1.0.0
description: "进销存：订货、调拨、退货、入出库与库存查询"
metadata:
  bins: [qmai]
  help: "qmai inventory --help"
---

> 前置条件: 请先阅读 ../qmai-shared/SKILL.md 了解认证和配置

## 核心概念

- **sales**: 报货单、订货单、调拨单、退货单
- **stock**: 入库单、出库单、门店库存台账、实时库存
- **master**: 盘点调整、物品、分类、单位、供应商、品项下发

## 命令概览

### 销售管理
```bash
qmai inventory sales declare-receive --from-json declare-receive.json --dry-run
qmai inventory sales declare-detail --declare-no BH20260407001
qmai inventory sales require-update --from-json require-update.json --dry-run
qmai inventory sales require-deliver --from-json require-deliver.json --dry-run
qmai inventory sales return-cancel --from-json return-cancel.json --dry-run
qmai inventory sales return-examine --from-json return-examine.json --dry-run
qmai inventory sales return-receipt --from-json return-receipt.json --dry-run
qmai inventory sales delivery-arrive --from-json delivery-arrive.json --dry-run
qmai inventory sales require-list --page 1 --page-size 10
qmai inventory sales require-create --from-json require-create.json --dry-run
qmai inventory sales require-detail-update --from-json require-detail-update.json --dry-run
qmai inventory sales transfer-list --page 1 --page-size 10
qmai inventory sales transfer-detail --transfer-no DB20260407001
qmai inventory sales return-list --page 1 --page-size 10
```

### 库存管理
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

### 基础资料管理
```bash
qmai inventory master inventory-adjust --from-json inventory-adjust.json --dry-run
qmai inventory master transfer-audit --from-json transfer-audit.json --dry-run
qmai inventory master product-create --from-json product-create.json --dry-run
qmai inventory master category-create --from-json category-create.json --dry-run
qmai inventory master unit-create --from-json unit-create.json --dry-run
qmai inventory master product-update --from-json product-update.json --dry-run
qmai inventory master supplier-create --from-json supplier-create.json --dry-run
qmai inventory master product-distribute --from-json product-distribute.json --dry-run
qmai inventory master supplier-update --from-json supplier-update.json --dry-run
qmai inventory master machining-card-batch-create --from-json machining-card-batch-create.json --dry-run
```

## 注意事项（Agent 必读）

- `declare-receive`、`require-update`、`require-deliver`、`return-cancel`、`return-examine`、`return-receipt`、`delivery-arrive`、`require-create` 以及库存/基础资料写接口都统一使用 `--from-json`
- 订货、调拨、退货、库存日志目前仍是显式分页，不会自动拉全量
- 半成品成本卡批量创建已经进入正式命令面，更新/套件/转化卡等待上线接口仍未暴露
- 文档中部分历史 ID 已失效或漂移，新增接口前应先重新核实文档地址和更新时间
