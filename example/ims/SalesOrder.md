# SalesOrder

## Description

SalesOrder model

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| id | uuid |  | false | [SalesOrderLineItem](SalesOrderLineItem.md) [Shipment](Shipment.md) |  |  |
| active | boolean |  | true |  |  | active |
| orderNumber | integer |  | true |  |  | Order number |
| customerID | uuid |  | true |  | [Contact](Contact.md) | Contact model. Contact and this model is n:1 |
| customer | Contact |  | true |  | [Contact](Contact.md) | Customer contact |
| shopifyID | string |  | true |  |  | Shopify order ID |
| cancelReason | string |  | true |  |  | Reason for cancellation |
| cancelledAt | datetime |  | true |  |  | Cancellation date |
| currency | string |  | true |  |  | Currency |
| currentSubtotalPrice | string |  | true |  |  | Current subtotal price |
| customerName | string |  | true |  |  | Customer name |
| customerEmail | string |  | true |  |  | Customer email |
| createdAt | datetime |  | true |  |  | createdAt |
| shipStationOrderStatus | enum |  | true |  |  | inventoryType |
| shippedAt | datetime |  | true |  |  | shipped at |
| updatedAt | datetime |  | true |  |  | updatedAt |

## Constraints

| Name | Type | Definition |
| ---- | ---- | ---------- |
| ForeignKey for customer to Contact | FOREIGN KEY | ForeignKeyType: Contact |

## Indexes

| Name | Definition |
| ---- | ---------- |
| Index for createdAt | Index: true |
| Index for updatedAt | Index: true |

## Relations

![er](SalesOrder.svg)

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
