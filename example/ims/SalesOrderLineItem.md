# SalesOrderLineItem

## Description

SalesOrderLineItem model

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| id | uuid |  | false |  |  |  |
| active | boolean |  | true |  |  | active |
| createdAt | datetime |  | true |  |  | createdAt |
| shopifyID | string |  | true |  |  | Shopify product ID |
| salesOrderID | uuid |  | true |  | [SalesOrder](SalesOrder.md) | salesOrder ID |
| salesOrder | SalesOrder |  | true |  | [SalesOrder](SalesOrder.md) | SalesOrder model. SalesOrder and this model is n:1. |
| updatedAt | datetime |  | true |  |  | updatedAt |
| variantID | uuid |  | true |  | [ProductVariant](ProductVariant.md) | Variant ID |
| variant | ProductVariant |  | true |  | [ProductVariant](ProductVariant.md) | Variant |
| quantity | float |  | true |  |  | quantity |
| unitPrice | float |  | true |  |  | unitPrice |
| subtotalPrice | float |  | true |  |  | subtotalCost |
| unitCompareAtPrice | float |  | true |  |  | unitCompareAtPrice |
| discount | float |  | true |  |  | discount |
| name | string |  | true |  |  | name |
| sku | string |  | true |  |  | sku |
| requiresShipping | boolean |  | true |  |  | requiresShipping |
| taxable | boolean |  | true |  |  | taxable |
| fulfillmentStatus | enum |  | true |  |  | fulfillmentStatus |

## Indexes

| Name | Definition |
| ---- | ---------- |
| Index for createdAt | Index: true |
| Index for updatedAt | Index: true |

## Relations

![er](SalesOrderLineItem.svg)

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
