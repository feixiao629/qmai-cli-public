# qmai-order 写操作参考

## 适用场景

- 订单评价回复
- 第三方订单上报
- 第三方退款单上报

## 推荐命令

```bash
qmai order review reply --order-no O20260407001 --reply-at "2026-04-07 18:00:00" --seller-reply-info "感谢反馈" --dry-run
qmai order report upload --from-json order-upload.json
qmai order report refund-up --from-json refund-up.json
```

## 风险控制

- `review reply` 先用 `--dry-run` 校验参数，再执行真实回复
- `order upload` 和 `refund-up` 先保留原始 JSON 文件，便于重试和审计
- 上报类接口成功返回的是受理结果，不代表下游业务已经全部完成
