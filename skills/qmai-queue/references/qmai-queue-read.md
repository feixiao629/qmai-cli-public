# qmai-queue 查询参考

## 适用场景

- 查询门店当前排队压力
- 查询单笔订单还要等多久
- 查询门店当前叫号列表

## 推荐命令

```bash
qmai queue shop-progress --shop-type 1 --shop-ids 1001,1002
qmai queue order-progress --order-no O20260407001
qmai queue shop-queue-nos --shop-code S001 --page 1 --size 10
```

## 排查建议

- `order-progress` 中 `order-no` 和 `source-no` 至少传一个，同时传时以 `order-no` 为准
- `shop-queue-nos` 默认可不传 `status-list`，官方会按待制作、制作中、待取餐查询
- 如果要保留完整原始字段，优先使用 `--format json`
