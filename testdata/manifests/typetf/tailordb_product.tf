resource "tailor_tailordb_type" "product" {
  workspace_id = tailor_workspace.my.id
  namespace    = tailor_tailordb.my.namespace
  name         = "Product"
  description  = "Product model"

  indexes = {
    shopifyID = {
      field_names = ["shopifyID", "title"]
      unique = true
    }
  }
  
  fields = {
    shopifyID = {
      type        = "string"
      description = "Shopify product ID"
      required    = true
      index       = true
      unique      = true
    }
    title = {
      type        = "string"
      description = "Title of the product"
    }
    handle = {
      type        = "string"
      description = "Product handle"
    }
    description = {
      type        = "string"
      description = "Product description"
    }
    featuredImageID = {
      type        = "uuid"
      description = "Featured image ID"
    }
    featuredImage = {
      type        = tailor_tailordb_type.product_image.name
      description = "Featured image of the product"
      source = "featuredImageID"
      foreign_key = {
        type = tailor_tailordb_type.product_image.name
      }
    }
  }

  type_permission = local.permission_everyone
}

resource "tailor_tailordb_type" "product_image" {
  workspace_id = tailor_workspace.my.id
  namespace    = tailor_tailordb.my.namespace
  name         = "ProductImage"
  description  = "ProductImage model"

  fields = {
    url = {
      type        = "string"
      description = "Image URL"
      required    = true
      index       = true
    }
    description = {
      type        = "string"
      description = "Product description"
    }
  }

  type_permission = local.permission_everyone
}

