# 商品增删改查操作指南

## 数据结构

### 列表返回 (OpenProduct)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 商品列表 ID |
| goodsId | int64 | 底层商品 ID |
| name | string | 商品名称 |
| status | int | 10=上架, 20=下架 |
| showPriceLow | int | 最低价（**单位：分**） |
| showPriceHigh | int | 最高价（**单位：分**） |
| saleChannel | int | 销售渠道: 1=堂食, 2=外卖, 11=堂食+外卖 |
| categoryNameList | []string | 分类名称数组 |
| goodsSkuList | []GoodsSku | SKU 列表 |

### SKU (GoodsSku)

| 字段 | 类型 | 说明 |
|------|------|------|
| skuId | string | SKU ID |
| tradeMark | string | 商品编码（用于上下架） |
| salePrice | float64 | 销售价（元） |
| inventory | float64 | 库存（API 可能返回 int 或 float） |

### 同步输入 (ShopGoods)

用户只需填写以下 5 个字段，`Sync()` 自动补全 20+ 个 API 必填字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| tradeName | string | 商品名称（必填） |
| tradePrice | float64 | 价格（必填，> 0） |
| tradeNo | string | 商户自定义编号 |
| className | string | 分类名称 |
| stock | int | 库存（默认 9999） |

## 创建商品（通过 Sync API）

### 命令行方式
```bash
qmai product create --name "美式咖啡" --price 28.00 --trade-no P001 --class 饮品
```

### JSON 文件方式
```bash
qmai product create --from-json goods.json
```

`goods.json` 示例（只需包含用户关心的字段）:
```json
[
  {"tradeName": "美式咖啡", "tradePrice": 28.00, "tradeNo": "P001", "className": "饮品"},
  {"tradeName": "拿铁", "tradePrice": 32.00, "tradeNo": "P002", "className": "饮品"}
]
```

### Dry-run 预览
```bash
$ qmai product create --name "拿铁" --price 32.00 --dry-run
[dry-run] 将要同步 1 个商品:
  拿铁 (32.00)
```

> **注意**: Sync API 是异步的，创建成功后商品不会立即出现在列表中

## 查询商品

```bash
# 全部商品（默认每页 20 条，最大 50）
qmai product list

# 按名称搜索
qmai product list --name 咖啡

# 按销售渠道
qmai product list --sale-channel 1

# 按类型
qmai product list --sale-type 1

# 分页
qmai product list --page 1 --page-size 50

# JSON 输出
qmai product list --format json

# 商品详情（按 ID 或名称）
qmai product get <id|name>
```

### 列表输出示例
```
ID                   名称       价格    状态    库存    分类
--                   ------   ------  ------  ------  ------
1185299121892380796  桃花落拿铁  21.00   上架    9999    咖啡系列
1185299121892380733  春樱拿铁   21.00   上架    9999    咖啡系列
```

## 更新商品（通过 Sync API）

```bash
# 改名（--name 必填）
qmai product update <tradeNo> --name "新品名" --price 35.00

# 改分类
qmai product update <tradeNo> --name "原名称" --class 饮品

# 多字段更新
qmai product update <tradeNo> --name "新品" --price 30.00 --class 小食
```

## 下架商品

```bash
# 需要 --force 确认（通过 BatchDown API）
qmai product delete <tradeMark> --force --sale-channel 1

# 不带 --force 只提示
$ qmai product delete <tradeMark>
确定要下架商品 <tradeMark> 吗？使用 --force 确认
```

> 下架通过 BatchDown API 实现，将商品 status 从 10 改为 20
