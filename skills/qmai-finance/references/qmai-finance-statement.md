# qmai-finance 对账参考

## 适用场景

- 查询分账明细
- 查询微信或支付宝支付账单
- 获取微信账单下载地址

## 推荐命令

```bash
qmai finance statement split-flows --created-at-start "2026-04-01 00:00:00" --created-at-end "2026-04-07 23:59:59" --page 1 --page-size 20
qmai finance statement wechat-bills --bill-date 2026-04-07 --page 1 --page-size 20
qmai finance statement alipay-bills --bill-date 2026-04-07 --page 1 --page-size 20
qmai finance statement wechat-bill-url --bill-date 2026-04-07 --shop-code S001
```

## 排查建议

- 对账时间格式要严格按接口要求填写，`split-flows` 用完整时间，账单接口用日期
- 微信账单查询支持多门店编码列表；需要跨店汇总时优先用 `--shop-codes`
- 支付宝账单字段文档不完整，默认表格只回显原始记录
