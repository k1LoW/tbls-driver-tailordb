# ProductVariant

## Description

The variants of a product

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| id | uuid |  | false | [FinancialStockEvent](FinancialStockEvent.md) [FinancialStockSummary](FinancialStockSummary.md) [InventoryItem](InventoryItem.md) [InvoiceLineItem](InvoiceLineItem.md) [OperationalStockEvent](OperationalStockEvent.md) [PurchaseOrderLineItem](PurchaseOrderLineItem.md) [ReceiptLineItem](ReceiptLineItem.md) [SalesOrderLineItem](SalesOrderLineItem.md) [ShipmentLineItem](ShipmentLineItem.md) [StockSummary](StockSummary.md) |  |  |
| active | boolean |  | true |  |  | active |
| createdAt | datetime |  | true |  |  | createdAt |
| shopifyID | string |  | true |  |  | Shopify product ID |
| availableForSale | boolean |  | true |  |  | Is the product available for sale |
| barcode | string |  | true |  |  | The barcode of the product |
| sku | string |  | true |  |  | The sku(Stock Keeping Unit) of the product |
| displayName | string |  | true |  |  | The display name of the product |
| imageID | uuid |  | true |  | [ProductImage](ProductImage.md) | The image of the product variant |
| image | ProductImage |  | true |  | [ProductImage](ProductImage.md) | The image of the product variant |
| inventoryQuantity | integer |  | true |  |  | The inventory quantity of the product |
| price | float |  | true |  |  | The price of the product |
| inventoryItemID | uuid |  | true |  | [InventoryItem](InventoryItem.md) | The inventory item ID of the product deprecated |
| inventoryItem | InventoryItem |  | true |  | [InventoryItem](InventoryItem.md) | The inventory item of the product |
| productID | uuid |  | true |  | [Product](Product.md) | The product ID of the product variant |
| product | Product |  | true |  | [Product](Product.md) | The product of the product variant |
| taxable | boolean |  | true |  |  | Is the product taxable |
| inventoryType | enum |  | true |  |  | inventoryType |
| quickbookItemId | string |  | true |  |  | The quickbook item ID of the product |
| quickbookSyncToken | string |  | true |  |  | The quickbook sync token of the product |
| quickbookItemName | string |  | true |  |  | The quickbook item name of the product |
| updatedAt | datetime |  | true |  |  | updatedAt |

## Indexes

| Name | Definition |
| ---- | ---------- |
| Index for createdAt | Index: true |
| Index for updatedAt | Index: true |

## Relations

```mermaid
erDiagram

"FinancialStockEvent" }o--o| "ProductVariant" : "Source: ProductVariant"
"FinancialStockSummary" }o--|| "ProductVariant" : "Source: ProductVariant"
"InventoryItem" }o--o| "ProductVariant" : "Source: ProductVariant"
"InvoiceLineItem" }o--o| "ProductVariant" : "Source: ProductVariant"
"OperationalStockEvent" }o--|| "ProductVariant" : "Source: ProductVariant"
"PurchaseOrderLineItem" }o--|| "ProductVariant" : "Source: ProductVariant"
"ReceiptLineItem" }o--|| "ProductVariant" : "Source: ProductVariant"
"SalesOrderLineItem" }o--o| "ProductVariant" : "Source: ProductVariant"
"ShipmentLineItem" }o--o| "ProductVariant" : "Source: ProductVariant"
"StockSummary" }o--|| "ProductVariant" : "Source: ProductVariant"
"ProductVariant" }o--o| "ProductImage" : "Source: ProductImage"
"ProductVariant" }o--o| "InventoryItem" : "Source: InventoryItem"
"ProductVariant" }o--o| "Product" : "Source: Product"

"ProductVariant" {
  uuid id
  boolean active
  datetime createdAt
  string shopifyID
  boolean availableForSale
  string barcode
  string sku
  string displayName
  uuid imageID FK
  ProductImage image FK
  integer inventoryQuantity
  float price
  uuid inventoryItemID FK
  InventoryItem inventoryItem FK
  uuid productID FK
  Product product FK
  boolean taxable
  enum inventoryType
  string quickbookItemId
  string quickbookSyncToken
  string quickbookItemName
  datetime updatedAt
}
"FinancialStockEvent" {
  uuid id
  boolean active
  datetime createdAt
  datetime updatedAt
  uuid variantID FK
  ProductVariant variant FK
  float incrementalQuantity
  float unitCost
  float transactionTotalCost
  boolean isOnHold
  float onHoldQuantity
  float availableQuantity
  float inStockQuantity
  float totalCost
  float averageCost
  uuid receiptLineItemID FK
  ReceiptLineItem receiptLineItem FK
  uuid shipmentLineItemID FK
  ShipmentLineItem shipmentLineItem FK
  uuid operationalStockEventID FK
  OperationalStockEvent operationalStockEvent FK
}
"FinancialStockSummary" {
  uuid id
  boolean active
  datetime createdAt
  datetime updatedAt
  uuid variantID FK
  ProductVariant variant FK
  float onHoldQuantity
  float availableQuantity
  float inStockQuantity
  float totalCost
  float averageCost
}
"InventoryItem" {
  uuid id
  boolean active
  datetime createdAt
  string shopifyID
  uuid productVariantID FK
  ProductVariant productVariant FK
  datetime updatedAt
}
"InvoiceLineItem" {
  uuid id
  boolean active
  datetime createdAt
  uuid invoiceID FK
  Invoice invoice FK
  float unitPrice
  datetime updatedAt
  uuid variantID FK
  ProductVariant variant FK
  float quantity
  boolean taxable
}
"OperationalStockEvent" {
  uuid id
  boolean active
  datetime createdAt
  datetime updatedAt
  uuid variantID FK
  ProductVariant variant FK
  float incrementalQuantity
  float unitCost
  float transactionTotalCost
  boolean isOnHold
  float onHoldQuantity
  float availableQuantity
  float inStockQuantity
  float totalCost
  float averageCost
  uuid receiptLineItemID FK
  ReceiptLineItem receiptLineItem FK
  uuid shipmentLineItemID FK
  ShipmentLineItem shipmentLineItem FK
  integer sequence
  boolean copiedToFinancialLedger
}
"PurchaseOrderLineItem" {
  uuid id
  boolean active
  datetime createdAt
  uuid purchaseOrderID FK
  PurchaseOrder purchaseOrder FK
  datetime updatedAt
  uuid variantID FK
  ProductVariant variant FK
  float quantity
  float unitCost
  float subtotalCost
}
"ReceiptLineItem" {
  uuid id
  boolean active
  datetime createdAt
  uuid receiptID FK
  Receipt receipt FK
  datetime updatedAt
  uuid variantID FK
  ProductVariant variant FK
  float quantity
  float subtotalUnitCost
  float subtotalCost
  float cubicMeters
  Array__nested__ costPools
  uuid costPools.costPoolID FK
  CostPool costPools.costPool FK
  float totalCostPoolAllocation
  float unitCostPoolAllocation
  float totalUnitCost
  enum receiptStatus
  datetime receivedAt
}
"SalesOrderLineItem" {
  uuid id
  boolean active
  datetime createdAt
  string shopifyID
  uuid salesOrderID FK
  SalesOrder salesOrder FK
  datetime updatedAt
  uuid variantID FK
  ProductVariant variant FK
  float quantity
  float unitPrice
  float subtotalPrice
  float unitCompareAtPrice
  float discount
  string name
  string sku
  boolean requiresShipping
  boolean taxable
  enum fulfillmentStatus
}
"ShipmentLineItem" {
  uuid id
  boolean active
  datetime createdAt
  uuid shipmentID FK
  Shipment shipment FK
  float unitCost
  float unitPrice
  datetime updatedAt
  uuid variantID FK
  ProductVariant variant FK
  float quantity
  boolean taxable
  enum shipmentStatus
  datetime shippedAt
}
"StockSummary" {
  uuid id
  boolean active
  datetime createdAt
  datetime updatedAt
  uuid variantID FK
  ProductVariant variant FK
  float onHoldQuantity
  float availableQuantity
  float inStockQuantity
  float totalCost
  float averageCost
}
"ProductImage" {
  uuid id
  boolean active
  datetime createdAt
  string shopifyID
  datetime updatedAt
  string url
  string altText
}
"Product" {
  uuid id
  boolean active
  datetime createdAt
  string shopifyID
  string title
  string handle
  string description
  uuid featuredImageID FK
  ProductImage featuredImage FK
  integer inStock
  datetime updatedAt
}
```

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
