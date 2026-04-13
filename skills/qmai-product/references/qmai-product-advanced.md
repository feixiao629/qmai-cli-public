# 商品高级命令指南

## 估清与取消估清

```bash
# 商品估清
qmai product estimate-clear --trade-mark TM001

# 取消估清
qmai product cancel-estimate-clear --trade-mark TM001

# 指定渠道和售卖类型取消估清
qmai product cancel-estimate-clear \
  --trade-marks TM001,TM002 \
  --sale-channel 1 \
  --sale-type 2 \
  --is-all-empty 1
```

> `estimate-clear` 走 `sellOut`，`cancel-estimate-clear` 走 `fillUp`。

## 售罄与置满

```bash
# 批量售罄
qmai product sold-out --trade-marks TM001,TM002 --sale-channels 3 --sale-types 1

# 批量置满
qmai product fill-full --trade-marks TM001 --sale-channels 3 --sale-types 1 --dry-run
```

> `sold-out` 走 `externalEmpty`，`fill-full` 走 `externalFull`。

## 做法启停

```bash
# 启用做法
qmai product practice-status --practice-values 少糖,去冰 --status 1

# 停用做法
qmai product practice-status --practice-values 热饮 --status 0
```

## 支持做法商品列表

```bash
qmai product list-with-practice \
  --sale-channel 3 \
  --sale-type 1 \
  --status 10 \
  --include-properties PRACTICE,SKU,ATTACH
```

必填项：
- `--sale-channel`
- `--sale-type`
- `--status`

## 门店加料列表

```bash
qmai product attach-list --sale-channel 3 --sale-type 1
```

## 门店商品实时数据

```bash
qmai product energy --goods-id 20001 --store-id 10001

qmai product realtime \
  --store-id 10001 \
  --goods-ids 20001,20002 \
  --sale-channel 3 \
  --sale-type 1
```

> `realtime` 用的是 `storeId`，不是 `shopCode`。

## 删除任务

```bash
# 按商品标识删除
qmai product delete-task --sale-channel 3 --sale-type 1 --trade-marks TM001,TM002 --dry-run

# 按规格码删除
qmai product delete-task --sale-channel 4 --sale-type 0 --spec-codes SKU001,SKU002
```

注意：
- 删除任务只会返回 `taskId`
- 该操作不可撤销
- 必须先确认渠道和售卖类型是否正确
