---
name: qmai-member
version: 1.2.0
description: "会员服务：资产查询与写入、手机号注册、会员更新、会员标签查询与打标"
metadata:
  bins: [qmai]
  help: "qmai member --help"
---

> 前置条件: 请先阅读 ../qmai-shared/SKILL.md 了解认证和配置

## 核心概念

- **customerId**: 会员主键，大多数会员接口都以它作为入参
- **会员资产**: 当前命令覆盖积分、储值余额、优惠券、优惠券明细、积分明细、余额明细、储值规则、资产膨胀状态，以及充值、扣减、冲正、消费、核销等写接口
- **会员信息**: 当前首批命令覆盖会员查询、会员基础信息、会员 ID 映射、等级与风险等级，以及手机号注册、基本信息更新
- **会员标签**: 当前首批命令覆盖标签列表、标签详情、品牌标签、标签组、标签 ID 列表和打标

## 命令概览

### 会员资产
```bash
qmai member asset points --customer-id C001
qmai member asset balance --customer-id C001
qmai member asset coupons --customer-id C001
qmai member asset coupon-list --customer-id C001 --page 1 --page-size 20
qmai member asset balance-flow --customer-id C001 --page 1 --page-size 20
qmai member asset points-flow --customer-id C001 --page 1 --page-size 20
qmai member asset level-experience --customer-id C001
qmai member asset discount --customer-id C001
qmai member asset personal-asset --customer-id C001
qmai member asset deposit-rules --shop-code S001
qmai member asset inflate-status --customer-id C001
qmai member asset coupon-detail-list --customer-id C001 --store-id 1001
qmai member asset recharge --from-json recharge.json --dry-run
qmai member asset recharge-reverse --from-json recharge-reverse.json --dry-run
qmai member asset balance-debit --from-json balance-debit.json --dry-run
qmai member asset balance-reverse --from-json balance-reverse.json --dry-run
qmai member asset offline-balance-op --from-json offline-balance-op.json --dry-run
qmai member asset points-debit --from-json points-debit.json --dry-run
qmai member asset points-reverse --from-json points-reverse.json --dry-run
qmai member asset points-add --from-json points-add.json --dry-run
qmai member asset consume --from-json consume.json --dry-run
qmai member asset consume-reverse --from-json consume-reverse.json --dry-run
qmai member asset coupon-writeoff --from-json coupon-writeoff.json --dry-run
qmai member asset coupon-reverse --from-json coupon-reverse.json --dry-run
```

### 会员信息
```bash
qmai member profile info --identifier-number 13800138000 --type 2
qmai member profile id-by-code --code 13800138000 --type 2
qmai member profile ids-by-phone --phones 13800138000,13900139000
qmai member profile base-info --customer-id C001
qmai member profile open-id --customer-id C001 --type 1
qmai member profile search --phone 13800138000
qmai member profile dynamic-code --customer-id 1001
qmai member profile risk-level --code 13800138000 --type 2
qmai member profile level --mobile-phone 13800138000
qmai member profile register-phone --mobile-phone 13800138000 --reg-app-type 1 --dry-run
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
qmai member profile update --id C001 --nickname "新昵称" --dry-run
```

### 会员标签
```bash
qmai member tag list --customer-id C001 --page 1 --page-size 20
qmai member tag detail --id 1001
qmai member tag brand-list
qmai member tag ids --customer-id C001
qmai member tag groups --label-attributed 2
qmai member tag group-labels --label-group-id 1001
qmai member tag create --label-code VIP_TAG --label-name "会员标签" --dry-run
qmai member tag delete --label-code VIP_TAG --dry-run
qmai member tag clear-members --label-code VIP_TAG --customer-ids C001,C002 --dry-run
qmai member tag delete-customer-label --customer-id C001 --panorama-label-id 1001 --dry-run
qmai member tag settings
qmai member tag mark --customer-ids C001,C002 --label-code VIP_TAG --mark-date "2026-04-07 12:00:00" --dry-run
```

## 开放平台 API 端点

| 操作 | 端点 |
| --- | --- |
| 会员积分查询 | POST /v3/crm/points/getCustomerPoints |
| 查询会员优惠券 | POST /v3/crm/coupon/getCustomerCoupons |
| 查询会员储值余额扣减明细 | POST /v3/crm/account/getDecreaseBalanceBiz |
| 查询会员储值余额明细 | POST /v3/crm/account/queryAccountFlow |
| 查询会员实时折扣 | POST /v3/crm/customer/getLevelConsumeDiscount |
| 查询会员等级和经验值 | POST /v3/crm/customer/getLevelExperience |
| 查询会员储值余额 | POST /v3/crm/customer/getCustomerBalance |
| 查询会员优惠券列表 | POST /v3/crm/coupon/getCustomerCouponList |
| 查询会员积分明细 | POST /v3/crm/customer/getCrmPointsFlow |
| 查询用户资产 | POST /v3/crm/customerInfo/personalAsset |
| 查询门店储值规则 | POST /v3/marketing/deposit/getStoreOrShopDepositList |
| 查询资产膨胀状态 | POST /v3/crm/customer/assetInflateDetail |
| 优惠券列表明细 | POST /v3/newPattern/crmCenter/get/customer-coupons/get_my_coupons_v2 |
| 会员储值充值 | POST /v3/bsns/customerBalance/userRecharge |
| 会员充值冲正 | POST /v3/bsns/customerBalance/userRechargeReverse |
| 会员储值余额扣减 | POST /v3/crm/account/decreaseBalance |
| 会员储值余额冲正 | POST /v3/crm/account/decreaseBalanceReverse |
| 储值账户充值扣减 | POST /v3/newPattern/crmCenter/post/customer-account/v2/offline-operation-balance |
| 会员积分扣减 | POST /v3/crm/points/reducePoints |
| 会员积分冲正 | POST /v3/crm/points/reversePoints |
| 会员积分发放 | POST /v3/crm/points/addCustomerPoints |
| 会员资产消费 | POST /v3/crm/customer/customerConsume |
| 会员资产冲正 | POST /v3/crm/customer/consumeReverse |
| 优惠券核销 | POST /v3/crm/coupon/writeOffCoupon |
| 优惠券冲正 | POST /v3/crm/coupon/couponReverse |
| 会员信息查询 | POST /v3/crm/customer/getCustomerInfo |
| 通过会员标识查询会员ID | POST /v3/crm/customer/getCustomerIdByCode |
| 通过手机号批量查询会员ID | POST /v3/crm/customer/getCustomerIdByPhone |
| 通过会员ID获取微信支付宝用户唯一标识 | POST /v3/crm/customer/getCustomerOpenId |
| 通过手机号注册会员 | POST /v3/crm/customer/phoneRegister |
| 查询会员的动态码 | POST /v3/crm/customer/getCustomerCode |
| 通过三方ID注册会员 | POST /v3/crm/customer/customerRegister |
| 发送短信验证码 | POST /v3/message/commonSendSmsCaptcha |
| 校验短信验证码 | POST /v3/message/commonCheckSmsCaptcha |
| 设置黑名单 | POST /v3/crm/customer/batchAddBlackList |
| 取消黑名单 | POST /v3/crm/customer/batchCancelBlackList |
| 黑名单校验 | POST /v3/crm/customer/queryMemberBlack |
| 冻结用户 | POST /v3/crm/customer/logoffFreeze |
| 冻结用户-撤销 | POST /v3/crm/customer/logoffUnfreeze |
| 会员注销 | POST /v3/crm/customer/logOff |
| 会员换绑手机号 | POST /v3/crm/customer/updatePhone |
| 查询用户冻结记录 | POST /v3/crm/customer/queryFreezeRecord |
| 查询用户是否存在未完成订单 | POST /v3/newPattern/accountCenter/post/account/can-logoff |
| 查询会员企微信息 | POST /v3/crm/customer/weComCustomerInfo |
| 查询会员签到状态 | POST /v3/crm/customer/queryActivitySign |
| 查询会员储值等级 | POST /v3/newPattern/crmCenter/get/customer-account-card/query-account-level |
| 按需查询会员信息 | POST /v3/newPattern/crmCenter/post/customer-condition/query |
| 更新会员基本信息 | POST /v3/crm/customer/updateMember |
| 会员基础信息查询 | POST /v3/crm/customer/getBaseInfo |
| 会员查询 | POST /v3/crm/customer/allCustomerInfo |
| 查询会员风险等级 | POST /v3/crm/customer/riskLevel |
| 查询会员等级 | POST /v3/crm/customer/member/detail |
| 通过会员ID查询用户标签 | POST /v3/cdp/panoramaLabel/pageByCustomerId |
| 查询标签详情 | POST /v3/cdp/panoramaLabel/detail |
| 查询静态用户标签列表（根据标签组ID） | POST /v3/cdp/panoramaLabel/list |
| 查询静态用户标签列表（根据标签类别） | POST /v3/cdp/panoramaLabel/listGroup |
| 创建会员标签 | POST /v3/cdp/panoramaLabel/panoramaLabelCreate |
| 删除会员标签 | POST /v3/cdp/panoramaLabel/panoramaLabelDelete |
| 清除会员标签下的会员 | POST /v3/cdp/panoramaLabel/panoramaLabelDelete |
| 清除会员的指定标签 | POST /v3/cdp/panoramaLabel/deleteCustomerLabel |
| 给指定会员打标签 | POST /v3/cdp/panoramaLabel/panoramaLabelMark |
| 查询标签设置 | POST /v3/crm/customer/brandInfo |
| 查询品牌下标签列表 | POST /v3/cdp/panoramaLabel/panoramaLabelList |
| 查询会员的标签id | POST /v3/newPattern/cdpCenter/get/panorama-label/list-by-customerId |

## 注意事项（Agent 必读）

- 当前 `qmai member` 已补资产写接口、手机号注册、会员更新和会员打标，但冻结、注销、短信验证码等接口仍未进入正式命令面
- 大多数会员接口不依赖 `shopCode`，但仍要求 profile 中已完成开放平台认证
- 分页命令当前是显式分页，`coupon-list`、`coupon-detail-list`、`balance-flow`、`points-flow`、`tag list`、`tag group-labels` 不会自动拉全量
- `personal-asset` 对应文档带有 C 端属性，CLI 目前只做查询封装，不默认承诺和小程序端完全一致
- 风险等级和 ID 映射类接口都依赖 `code + type` 组合，类型值要按企迈文档填写
- 资产写接口统一使用 `--from-json`，先 `--dry-run` 预览再执行，避免在 flag 层散落大量不稳定字段
- 标签写接口优先使用明确 flag，不走 `--from-json`；单次清除/打标人数上限仍以企迈文档为准
- 三方注册和短信验证码相关接口统一使用 `--from-json`，避免手机号、渠道、验证码事件配置散落在 flag 中
- `4.2.7 查询会员是否开启密码校验` 的文档未给出明确业务响应字段，当前仍保留在矩阵中但未暴露命令
