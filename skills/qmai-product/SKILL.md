---
name: qmai-product
version: 3.1.0
description: "商品管理：列表、同步、上下架、估清、售罄、置满、做法启停、实时查询、批量导入导出、批量调价"
metadata:
  bins: [qmai]
  help: "qmai product --help"
---

> 前置条件: 请先阅读 ../qmai-shared/SKILL.md 了解认证和配置

## 核心概念

- **列表商品 (OpenProduct)**: API 返回的商品结构，字段包括 id（数字ID）、name（名称）、status（10=上架/20=下架）、showPriceLow（价格，单位：分）、categoryNameList（分类名数组）、goodsSkuList（SKU列表）
- **同步商品 (ShopGoods)**: 创建/更新用的输入，用户只需填写 tradeName/tradePrice/tradeNo/className/stock 五个字段，Sync() 方法会自动补全 20+ 个 API 必填字段
- **SaleChannel**: 1（堂食）| 2（外卖）| 11（堂食+外卖）
- **SaleType**: 1（普通商品）| 2（套餐）
- **shopCode**: 门店编码，从 profile config 自动读取
- **tradeMark**: SKU 级别的商品标识，用于上下架操作（BatchUp/BatchDown）

## Shortcuts（推荐优先使用）

| Shortcut | 说明 | 示例 |
|----------|------|------|
| +quick-add | 快速添加商品 | `qmai product +quick-add "拿铁" 32.00` |
| +on-sale | 批量上架 | `qmai product +on-sale --trade-marks TM001,TM002` |
| +off-sale | 批量下架 | `qmai product +off-sale --trade-marks TM001` |
| +price-adjust | 快速调价 | `qmai product +price-adjust +10%` |

## API 命令

### 商品列表
```bash
qmai product list [--name 咖啡] [--sale-channel 1] [--sale-type 1] [--page 1] [--page-size 20]
```
> 注意: page-size 最大 50

### 商品详情
```bash
qmai product get <id|name>
```

### 创建/同步商品（通过 Sync API）
```bash
qmai product create --name "美式咖啡" --price 28.00 --trade-no P001 [--class 饮品]
qmai product create --from-json goods.json
qmai product create --name "拿铁" --price 32.00 --dry-run
```

### 更新商品（通过 Sync API）
```bash
qmai product update <tradeNo> --name "新名称" --price 30.00 [--class 饮品]
```

### 下架商品
```bash
qmai product delete <tradeMark> --force [--sale-channel 1]
```

### 批量上架/下架
```bash
qmai product batch-status --action up --trade-marks TM001,TM002 [--sale-channel 1]
qmai product batch-status --action down --trade-marks TM001 --dry-run
```

### 估清 / 取消估清
```bash
qmai product estimate-clear --trade-mark TM001 --dry-run
qmai product cancel-estimate-clear --trade-mark TM001
qmai product cancel-estimate-clear --trade-marks TM001,TM002 --sale-channel 1 --sale-type 2
```

### 售罄 / 置满
```bash
qmai product sold-out --trade-marks TM001,TM002 --sale-channels 3 --sale-types 1
qmai product fill-full --trade-marks TM001 --sale-channels 3 --sale-types 1 --dry-run
```

### 做法启停
```bash
qmai product practice-status --practice-values 少糖,去冰 --status 1
qmai product practice-status --practice-values 热饮 --status 0 --dry-run
```

### 实时数据 / 加料 / 支持做法的商品列表
```bash
qmai product attach-list --sale-channel 3 --sale-type 1
qmai product list-with-practice --sale-channel 3 --sale-type 1 --status 10 --include-properties PRACTICE,SKU
qmai product energy --goods-id 20001 --store-id 10001
qmai product realtime --store-id 10001 --goods-ids 20001,20002 --sale-channel 3 --sale-type 1
```

### 删除任务
```bash
qmai product delete-task --sale-channel 3 --sale-type 1 --trade-marks TM001,TM002 --dry-run
qmai product delete-task --sale-channel 4 --sale-type 0 --spec-codes SKU001,SKU002
```

### 批量导入（通过 Sync API）
```bash
qmai product import --file products.csv --dry-run    # 预览
qmai product import --file products.csv              # 执行
qmai product import --file products.json
```

### 批量导出
```bash
qmai product export --file products.csv
qmai product export --file products.json --name 咖啡
```

### 批量调价
```bash
qmai product batch-price --adjust +10% --dry-run
qmai product batch-price --adjust -5%
qmai product batch-price --adjust +2.00
```

### 分类
```bash
qmai product category   # 提示通过 --class 参数在创建/同步时指定分类
```

## 开放平台 API 端点

| 操作 | 端点 |
|------|------|
| 商品列表 | POST /v3/goods/item/getItemList |
| 商品同步 | POST /v3/goods/sync/shopGoodsSync |
| 批量上架 | POST /v3/goods/item/externalUp |
| 批量下架 | POST /v3/goods/item/externalDown |
| 估清 | POST /v3/goods/item/sellOut |
| 取消估清 | POST /v3/goods/item/fillUp |
| 批量售罄 | POST /v3/goods/item/externalEmpty |
| 批量置满 | POST /v3/goods/item/externalFull |
| 查询门店加料列表 | POST /v3/goods/item/shopGoodsAttachList |
| 查询门店商品列表(支持做法) | POST /v3/goods/item/getShopGoodsList |
| 查询商品能量值 | POST /v3/goods/nutritional/energy |
| 门店商品实时数据查询 | POST /v3/newPattern/goodsCenter/post/v2/item/real-time/list |
| 门店商品做法启用停用 | POST /v3/goods/item/practiceOnOff |
| 提交商品删除任务 | POST /v3/goods/sync/tripartiteShopGoodsDel |

## API 响应结构

### 列表响应 (getItemList)
```json
{
  "data": {
    "data": [ ... ],   // 商品数组（注意：嵌套在 data.data 中）
    "total": 380
  }
}
```

### 商品对象字段
| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 商品列表 ID |
| goodsId | int64 | 底层商品 ID |
| name | string | 商品名称 |
| status | int | 10=上架, 20=下架 |
| showPriceLow | int | 价格（**单位：分**，显示时 /100） |
| categoryNameList | []string | 分类名称数组 |
| goodsSkuList | []object | SKU 列表，含 skuId/tradeMark/salePrice/inventory |
| saleChannel | int | 销售渠道 |

## 注意事项（Agent 必读）

- **价格单位**: 列表 API 返回的 showPriceLow/showPriceHigh 单位是**分**，显示时需除以 100
- **状态码**: 10=上架, 20=下架（不是 1/0）
- **库存类型**: API 返回的 inventory 可能是 int 或 float，代码中用 float64 接收
- **分页限制**: page-size 最大 50
- **Sync API 必填字段多**: 用户只需填 5 个字段，Sync() 自动补全 pictureUrlList、categoryList（含 name/mark/isRequired/isBackend/isFront）、skuList（含 tradeMark/marketPrice/stock/clearStatus）、saleTime 等
- **同步是异步的**: Sync API 返回成功后，商品不会立即出现在列表中
- **删除任务不可撤销**: `delete-task` 只是提交删除任务，但业务效果不可撤销
- **实时数据查询**: `realtime` 需要 `storeId`，不是 `shopCode`
- **支持做法列表**: `list-with-practice` 需要显式传 `saleChannel`、`saleType` 和 `status`
- 批量操作前务必使用 `--dry-run` 预览
- 价格调整不可撤回，建议先 `export` 备份
- 下架操作需要 `--force` 确认
- 导入 CSV 首行必须为表头: trade_name,trade_price,trade_no,class_name,stock
- shopCode 从 profile config 自动读取，无需手动指定
- Debug 模式 (`--debug`) 会在一行日志中打印完整请求和响应 JSON
