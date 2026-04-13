# 门店与组织写操作指南

## 设置门店营业状态

```bash
# 1=开启 2=关闭
qmai store set-status S001 --status 1 --dry-run
qmai store set-status S001 --status 2
```

## 修改门店组别

```bash
qmai store move-team --shop-id 10001 --team-id 20001 --dry-run
qmai store move-team --shop-id 10001 --team-id 20001
```

## 门店信息同步

推荐使用 JSON 文件：

```bash
qmai store sync --from-json shop.json --dry-run
qmai store sync --from-json shop.json
```

`shop.json` 最少应包含：

```json
{
  "name": "测试门店",
  "provinceId": 310000,
  "cityId": 310100,
  "districtId": 310101,
  "lat": "31.2304",
  "lng": "121.4737",
  "managerStatus": 1
}
```

注意：
- `sync` 副作用大，先 `--dry-run`
- `shopInfoSync` 字段很多，避免直接在命令行拼长 JSON
- 提交成功不等于立即可查到变更结果
