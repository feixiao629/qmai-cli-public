# 批量操作指南

## 批量导入（通过 Sync API）

### CSV 格式

首行表头，后续为数据行：
```csv
trade_name,trade_price,trade_no,class_name,stock
美式咖啡,28.00,P001,饮品,100
拿铁,32.00,P002,饮品,100
抹茶拿铁,35.00,P003,饮品,50
```

### JSON 格式

```json
[
  {"tradeName": "美式咖啡", "tradePrice": 28.00, "tradeNo": "P001", "className": "饮品"},
  {"tradeName": "拿铁", "tradePrice": 32.00, "tradeNo": "P002", "className": "饮品"}
]
```

> 用户只需提供 5 个字段，`Sync()` 方法会自动补全所有 API 必填字段（pictureUrlList、categoryList、skuList、saleTime 等）

### 导入命令

```bash
# 先预览
qmai product import --file products.csv --dry-run

# 确认后执行（调用 shopGoodsSync）
qmai product import --file products.csv

# 跳过错误行继续
qmai product import --file products.csv --skip-errors
```

> **注意**: Sync API 是异步的，导入成功后商品不会立即出现在列表中

## 批量导出

```bash
# 导出为 CSV
qmai product export --file backup.csv

# 导出为 JSON
qmai product export --file backup.json

# 按名称筛选
qmai product export --file drinks.csv --name 咖啡
```

导出 CSV 表头为: `id,name,price,status,inventory,category`

> 注意: 导出最多一次获取 50 条（API 分页限制）

## 批量调价

### 百分比调价

```bash
# 全品涨价 10%（先预览）
qmai product batch-price --adjust +10% --dry-run

# 降价 5%
qmai product batch-price --adjust -5%
```

### 固定金额调价

```bash
# 所有商品涨 2 元
qmai product batch-price --adjust +2.00

# 降 1.5 元
qmai product batch-price --adjust -1.50
```

### Dry-run 输出示例

```
[dry-run] 将要调整 5 个商品价格:
  美式咖啡: 28.00 → 30.80
  拿铁: 32.00 → 35.20
  卡布奇诺: 30.00 → 33.00
  摩卡: 35.00 → 38.50
  冰美式: 25.00 → 27.50
```

> 批量调价内部流程：先通过 List API 查询当前价格（单位：分，/100 转为元），计算新价格后通过 Sync API 更新

## 批量上下架

```bash
# 批量上架（调用 externalUp API）
qmai product batch-status --action up --trade-marks TM001,TM002 [--sale-channel 1]

# 批量下架（调用 externalDown API）
qmai product batch-status --action down --trade-marks TM001 --dry-run

# 预览
qmai product batch-status --action up --trade-marks TM001,TM002 --dry-run
```

> `--trade-marks` 的值来自商品 SKU 的 tradeMark 字段（可通过 `qmai product list --format json` 查看 goodsSkuList[].tradeMark）

## 安全建议

1. **操作前备份**: `qmai product export --file backup.json`
2. **先 Dry-run**: 所有批量操作先 `--dry-run` 预览
3. **调价不可撤回**: 价格修改后无法自动回滚，保留备份
4. **分步执行**: 大批量操作建议按分类分步执行
