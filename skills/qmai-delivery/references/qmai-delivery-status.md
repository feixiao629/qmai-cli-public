# qmai-delivery 状态同步参考

## 适用场景

- 商家自配送状态回传
- 三方配送聚合接口简化对接

## 推荐命令

```bash
qmai delivery status update --origin-order-no O20260407001 --order-status 5 --dry-run
qmai delivery order create-and-update --from-json delivery-update.json --dry-run
```

## 排查建议

- `status update` 适用于商家自配送状态同步
- `create-and-update` 是独立聚合接口，不应与普通创建接口混合使用
- 当状态为取消或异常重发时，先确认企迈要求的状态推进顺序
