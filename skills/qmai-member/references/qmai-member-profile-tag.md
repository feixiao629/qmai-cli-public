# qmai-member profile/tag 参考

## 常见查询

```bash
qmai member profile info --identifier-number 13800138000 --type 2
qmai member profile id-by-code --code 13800138000 --type 2
qmai member profile ids-by-phone --phones 13800138000,13900139000
qmai member profile open-id --customer-id C001 --type 1
qmai member profile search --phone 13800138000 --page 1 --page-size 20
qmai member profile dynamic-code --customer-id 1001
qmai member profile risk-level --code 13800138000 --type 2
qmai member profile register-phone --mobile-phone 13800138000 --reg-app-type 1 --username 张三 --dry-run
qmai member profile register-third-id --from-json register-third-id.json --dry-run
qmai member profile send-captcha --from-json send-captcha.json --dry-run
qmai member profile check-captcha --from-json check-captcha.json --dry-run
qmai member profile blacklist-add --from-json blacklist-add.json --dry-run
qmai member profile blacklist-remove --from-json blacklist-remove.json --dry-run
qmai member profile blacklist-status --customer-id C001
qmai member profile freeze --customer-id C001 --reason "风险处理" --dry-run
qmai member profile unfreeze --customer-id C001 --dry-run
qmai member profile logoff --customer-id C001 --reason "用户申请注销" --dry-run
qmai member profile update-phone --customer-id 1001 --phone 13800138001 --reason "换绑手机号" --dry-run
qmai member profile freeze-record --customer-id C001
qmai member profile has-open-order --customer-id C001 --biz-type 5
qmai member profile wecom-info --customer-id 1001
qmai member profile sign-status --user-id 1001 --activity-id 2001
qmai member profile account-level --customer-id 1001
qmai member profile condition-query --customer-ids 1001,1002 --conditions 10,20,40,60
qmai member profile update --id C001 --nickname 新昵称 --qm-from wechat --dry-run
qmai member tag list --customer-id C001 --page 1 --page-size 20
qmai member tag detail --id 1001
qmai member tag brand-list
qmai member tag groups --label-attributed 2
qmai member tag group-labels --label-group-id 1001 --page 1 --page-size 20
qmai member tag create --label-code VIP_TAG --label-name "会员标签" --dry-run
qmai member tag delete --label-code VIP_TAG --dry-run
qmai member tag clear-members --label-code VIP_TAG --customer-ids C001,C002 --dry-run
qmai member tag delete-customer-label --customer-id C001 --panorama-label-id 1001 --dry-run
qmai member tag settings
qmai member tag mark --customer-ids C001,C002 --label-code VIP_TAG --mark-date "2026-04-07 12:00:00" --dry-run
qmai member tag ids --customer-id C001
```

## 约束

- `profile info` 至少传一个查询条件
- `id-by-code` 和 `risk-level` 都必须带 `--type`
- `open-id` 的 `--type` 不传时按企迈默认值处理
- `tag list` 是分页查询，不会自动取全量
- `register-phone`、`register-third-id`、短信验证码、标签写接口和 `tag mark` 都建议先用 `--dry-run`
