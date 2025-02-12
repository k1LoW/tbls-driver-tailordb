package typecue

import (
	"github.com/tailor-platform/tailorctl/schema/v2/tailordb"
	"cue.example/types"
)

tailordb.#Spec & {
	Namespace: "my-tailordb",
	Types: [
		for model in types {
			model
		},
	]
}
