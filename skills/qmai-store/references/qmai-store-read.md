# 门店与组织读取指南

## 门店详情

```bash
# 按门店编码
qmai store get S001

# 按门店 ID
qmai store get-by-id 10001

# 门店编码转门店 ID
qmai store id S001
```

## 门店列表

```bash
qmai store list --keyword 咖啡 --page 1 --page-size 10
qmai store list --contain-close 1
```

## 外卖映射

```bash
# 4=美团 5=饿了么
qmai store takeout-map-list --platform-type 4
```

## 配置、扩展字段、标签

```bash
qmai store config --store-id 10001 --field-codes OrderSettingConfig:autoCancelMinutes
qmai store brand-config --field-codes OrderSettingConfig:autoCancelMinutes
qmai store ext-data --shop-id 10001
qmai store labels --shop-id 10001
```

## 组织结构

```bash
qmai store org-tree
qmai store team-list
qmai store dept-tree
```
