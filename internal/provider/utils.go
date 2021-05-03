package improvmx

import (
	"hash/fnv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func hashSetValue(key string) schema.SchemaSetFunc {
	return func(v interface{}) int {
		alias := v.(map[string]interface{})[key].(string)
		return hash(alias)
	}
}

func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
}
