# 开放平台认证详情

## 凭证配置流程

```
$ qmai auth login
请输入 openId: d14c1559e87b747d577c834b275a4310
请输入 grantCode: ba67d4fa46
请输入 openKey: LyvrkvkxRkG2R6aM55bXpPwjYAbkEXTbVnKwfDYvVHjNwNFAmx
请输入门店编码 (shopCode): S001
✓ 凭证已保存 (profile: default)
```

## 凭证说明

| 凭证 | 说明 | 存储位置 |
|------|------|----------|
| openId | 开放平台应用 ID | config yaml |
| grantCode | 授权码 | config yaml |
| openKey | 签名密钥（敏感） | OS keychain |
| shopCode | 门店编码 | config yaml |

## 请求签名算法 (HmacSHA1)

1. 取 openId, grantCode, nonce, timestamp 四个字段
2. 按 key 字典序升序排列
3. 拼接为 `grantCode=xx&nonce=xx&openId=xx&timestamp=xx`
4. HmacSHA1(拼接串, openKey) → Base64 → URL Encode
5. 得到 token，一次性使用

示例签名验证：
- openKey: `LyvrkvkxRkG2R6aM55bXpPwjYAbkEXTbVnKwfDYvVHjNwNFAmx`
- openId: `d14c1559e87b747d577c834b275a4310`
- grantCode: `ba67d4fa46`
- timestamp: `1465185768`, nonce: `11886`
- 期望 token: `cFw0t9IuvL9jVo9qAzk0qMcw5BM%3D`

## AES 加解密（敏感字段）

部分接口的敏感字段需要 AES 加密传输：

- 算法: AES-256/CBC/PKCS5Padding
- Key: MD5(openKey) → 32 hex chars 作为字节
- IV: MD5(openId)[8:24] → 16 hex chars 作为字节

## 多门店支持

每个 profile 独立存储凭证：

```bash
qmai auth login --profile store1
qmai auth login --profile store2
qmai config set active_profile store1   # 切换到 store1
```
