package types

import(
	"github.com/tailor-platform/tailorctl/schema/v2/tailordb"
)

Product: tailordb.#Type & {
	Name:        "Product"
	Description: "Product model"
	Fields: {
		shopifyID: {
			Type:        tailordb.#TypeString
			Description: "Shopify product ID"
		}
		title: {
			Type:        tailordb.#TypeString
			Description: "Product title"
		}
		handle: {
			Type:        tailordb.#TypeString
			Description: "Product handle"
		}
		description: {
			Type:        tailordb.#TypeString
			Description: "Product description"
		}
		featuredImageID: {
			Type:        tailordb.#TypeUUID
			Description: "Featured image ID"
		}
		featuredImage: {
			Type:        "ProductImage"
			Description: "Featured image of the product"
			SourceId:    "featuredImageID"
		}
	}
}
