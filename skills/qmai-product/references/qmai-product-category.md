# 分类管理

## 概述

开放平台的分类通过 `className` 字段在商品同步（Sync API）时自动关联，**无独立的分类 CRUD API**。

## 工作方式

1. 在创建/同步商品时通过 `--class` 参数或 CSV 的 `class_name` 列指定分类名称
2. `Sync()` 方法自动构建 `categoryList`，包含 `name`/`mark`/`categoryName`/`isRequired`/`isBackend`/`isFront`/`sort`/`type` 等必填字段
3. 平台会自动匹配已有分类或创建新分类
4. 分类信息可在商品列表返回结果的 `categoryNameList` 字段中查看

## Sync API 的分类结构

用户只需提供 `className`（字符串），代码自动生成完整的分类对象：

```json
{
  "categoryList": [{
    "name": "饮品",
    "categoryName": "饮品",
    "mark": "饮品",
    "isRequired": 0,
    "isBackend": 0,
    "isFront": 0,
    "sort": 0,
    "type": 0
  }]
}
```

## 使用示例

### 同步时指定分类

```bash
# 创建商品时指定分类
qmai product create --name "美式咖啡" --price 28.00 --class 饮品

# 批量导入时在 CSV 中指定
trade_name,trade_price,trade_no,class_name,stock
美式咖啡,28.00,P001,饮品,100
拿铁,32.00,P002,饮品,100
三明治,15.00,P003,小食,50
```

### 查看分类

```bash
# 列表中查看分类列
qmai product list

# JSON 输出查看分类
qmai product list --format json

# 分类提示命令
qmai product category
# → 提示: 通过 --class 参数在创建/同步时指定分类
```

## 注意事项

- `className` 为字符串，平台按名称匹配
- 不支持分类的独立创建、删除、排序操作
- 修改商品分类 = 同步时传入新的 `className`
- 不指定分类时，`categoryList` 传空数组 `[]`
