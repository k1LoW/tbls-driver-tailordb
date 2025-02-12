# ims

## Tables

| Name | Columns | Comment | Type |
| ---- | ------- | ------- | ---- |
| [Category](Category.md) | 5 | Contains information about categories of products. | TailorDB.Type |
| [Contact](Contact.md) | 20 | Contact model | TailorDB.Type |
| [CostPool](CostPool.md) | 7 | CostPool model | TailorDB.Type |
| [CostPoolLineItem](CostPoolLineItem.md) | 9 | CostPoolLineItem model | TailorDB.Type |
| [FinancialStockEvent](FinancialStockEvent.md) | 21 | DO NOT UPDATE FROM THE FRONT END. FinancialStockEvent model. Holds StockEvents that will not change anymore. | TailorDB.Type |
| [FinancialStockSummary](FinancialStockSummary.md) | 11 | StockSummary model for the financial ledger | TailorDB.Type |
| [InventoryItem](InventoryItem.md) | 7 | Product variant's inventory item model | TailorDB.Type |
| [InventoryLevel](InventoryLevel.md) | 16 | Inventory Level of an inventory item | TailorDB.Type |
| [Invoice](Invoice.md) | 13 | Invoice model | TailorDB.Type |
| [InvoiceLineItem](InvoiceLineItem.md) | 11 | InvoiceLineItem model | TailorDB.Type |
| [Location](Location.md) | 6 | Inventory location on Shopify | TailorDB.Type |
| [OperationalStockEvent](OperationalStockEvent.md) | 21 | OperationalStockEvent model. Holds the stock event data that can change. | TailorDB.Type |
| [Product](Product.md) | 11 | Product model | TailorDB.Type |
| [ProductImage](ProductImage.md) | 7 | Product Image | TailorDB.Type |
| [ProductVariant](ProductVariant.md) | 22 | The variants of a product | TailorDB.Type |
| [PurchaseOrder](PurchaseOrder.md) | 24 | PurchaseOrder model | TailorDB.Type |
| [PurchaseOrderLineItem](PurchaseOrderLineItem.md) | 11 | PurchaseOrderLineItem model | TailorDB.Type |
| [Receipt](Receipt.md) | 11 | Receipt model | TailorDB.Type |
| [ReceiptLineItem](ReceiptLineItem.md) | 20 | ReceiptLineItem model | TailorDB.Type |
| [Role](Role.md) | 4 | User Roles | TailorDB.Type |
| [SalesOrder](SalesOrder.md) | 16 | SalesOrder model | TailorDB.Type |
| [SalesOrderLineItem](SalesOrderLineItem.md) | 19 | SalesOrderLineItem model | TailorDB.Type |
| [Shipment](Shipment.md) | 10 | Shipment model | TailorDB.Type |
| [ShipmentLineItem](ShipmentLineItem.md) | 14 | ShipmentLineItem model | TailorDB.Type |
| [StockSummary](StockSummary.md) | 11 | StockSummary model | TailorDB.Type |
| [User](User.md) | 6 | User of the system | TailorDB.Type |

## Relations

```mermaid
erDiagram

"Invoice" }o--o| "Contact" : "ForeignKeyType: Contact"
"Receipt" }o--o| "Contact" : "ForeignKeyType: Contact"
"SalesOrder" }o--o| "Contact" : "ForeignKeyType: Contact"
"Shipment" }o--o| "Contact" : "ForeignKeyType: Contact"
"FinancialStockEvent" }o--|| "OperationalStockEvent" : "ForeignKeyType: OperationalStockEvent"
"FinancialStockEvent" }o--o| "ReceiptLineItem" : "ForeignKeyType: ReceiptLineItem"
"OperationalStockEvent" }o--o| "ReceiptLineItem" : "ForeignKeyType: ReceiptLineItem"
"FinancialStockEvent" }o--o| "ShipmentLineItem" : "ForeignKeyType: ShipmentLineItem"
"OperationalStockEvent" }o--o| "ShipmentLineItem" : "ForeignKeyType: ShipmentLineItem"
"PurchaseOrder" }o--o| "Contact" : "Source: Contact"
"PurchaseOrder" }o--o| "Contact" : "Source: Contact"
"PurchaseOrder" }o--|| "Contact" : "Source: Contact"
"CostPoolLineItem" }o--|| "CostPool" : "Source: CostPool"
"ReceiptLineItem" }o--o| "CostPool" : "Source: CostPool"
"InventoryLevel" }o--o| "InventoryItem" : "Source: InventoryItem"
"ProductVariant" }o--o| "InventoryItem" : "Source: InventoryItem"
"InvoiceLineItem" }o--|| "Invoice" : "Source: Invoice"
"InventoryLevel" }o--o| "Location" : "Source: Location"
"ProductVariant" }o--o| "Product" : "Source: Product"
"Product" }o--o| "ProductImage" : "Source: ProductImage"
"ProductVariant" }o--o| "ProductImage" : "Source: ProductImage"
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
"PurchaseOrderLineItem" }o--|| "PurchaseOrder" : "Source: PurchaseOrder"
"Receipt" }o--o| "PurchaseOrder" : "Source: PurchaseOrder"
"ReceiptLineItem" }o--|| "Receipt" : "Source: Receipt"
"SalesOrderLineItem" }o--o| "SalesOrder" : "Source: SalesOrder"
"Shipment" }o--o| "SalesOrder" : "Source: SalesOrder"
"Invoice" }o--o| "Shipment" : "Source: Shipment"
"ShipmentLineItem" }o--|| "Shipment" : "Source: Shipment"

"Category" {
  uuid id
  string name
  string description
  datetime createdAt
  datetime updatedAt
}
"Contact" {
  uuid id
  boolean active
  datetime createdAt
  string name
  string email
  string phone
  string address1
  string address2
  string city
  string province
  string country
  datetime updatedAt
  string zip
  string countryCode
  string company
  string provinceCode
  string quickBookCustomerId
  float openBalance
  float creditLimit
  float availableCredit
}
"CostPool" {
  uuid id
  boolean active
  datetime createdAt
  string name
  boolean isClosed
  datetime closedAt
  datetime updatedAt
}
"CostPoolLineItem" {
  uuid id
  boolean active
  uuid costPoolID FK
  CostPool costPool FK
  datetime createdAt
  string name
  float amount
  enum allocationBase
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
"InventoryLevel" {
  uuid id
  boolean active
  datetime createdAt
  uuid locationID FK
  Location location FK
  uuid inventoryItemID FK
  InventoryItem inventoryItem FK
  integer available
  integer committed
  integer reserved
  integer damaged
  integer safetyStock
  integer qualityControl
  integer onHand
  integer incoming
  datetime updatedAt
}
"Invoice" {
  uuid id
  boolean active
  datetime createdAt
  string invoiceNumber
  uuid customerID FK
  Contact customer FK
  uuid shipmentID FK
  Shipment shipment FK
  datetime date
  string quickbookInvoiceId
  enum invoiceStatus
  datetime pushedToQBAt
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
"Location" {
  uuid id
  boolean active
  datetime createdAt
  string name
  string shopifyID
  datetime updatedAt
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
"ProductImage" {
  uuid id
  boolean active
  datetime createdAt
  string shopifyID
  datetime updatedAt
  string url
  string altText
}
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
"PurchaseOrder" {
  uuid id
  boolean active
  datetime createdAt
  string documentNumber
  uuid supplierID FK
  Contact supplier FK
  datetime date
  uuid billToID FK
  Contact billTo FK
  uuid shipToID FK
  Contact shipTo FK
  string shipVia
  string trackingNumber
  string shippingContactPhone
  string shippingContactName
  datetime exFactoryDate
  datetime dueDate
  string terms
  string shippingInstructions
  string notes
  string approvedBy
  string pulledBy
  string receivedBy
  datetime updatedAt
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
"Receipt" {
  uuid id
  boolean active
  datetime createdAt
  string receiptNumber
  uuid supplierID FK
  Contact supplier FK
  uuid purchaseOrderID FK
  PurchaseOrder purchaseOrder FK
  datetime date
  enum receiptStatus
  datetime updatedAt
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
"Role" {
  uuid id
  string name
  datetime createdAt
  datetime updatedAt
}
"SalesOrder" {
  uuid id
  boolean active
  integer orderNumber
  uuid customerID FK
  Contact customer FK
  string shopifyID
  string cancelReason
  datetime cancelledAt
  string currency
  string currentSubtotalPrice
  string customerName
  string customerEmail
  datetime createdAt
  enum shipStationOrderStatus
  datetime shippedAt
  datetime updatedAt
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
"Shipment" {
  uuid id
  boolean active
  datetime createdAt
  string shipmentNumber
  uuid customerID FK
  Contact customer FK
  uuid salesOrderID FK
  SalesOrder salesOrder FK
  datetime date
  datetime updatedAt
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
"User" {
  uuid id
  string name
  string email
  Array__uuid__ roles
  datetime createdAt
  datetime updatedAt
}
```

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
