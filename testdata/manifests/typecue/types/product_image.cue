package types

import(
	"github.com/tailor-platform/tailorctl/schema/v2/tailordb"
)

ProductImage: tailordb.#Type & {
	Name:        "ProductImage"
	Description: "ProductImage model"
	Fields: {
		url: {
			Type:        tailordb.#TypeString
			Description: "Image URL"
		}
		description: {
			Type:        tailordb.#TypeString
			Description: "Product description"
		}
	}
}
