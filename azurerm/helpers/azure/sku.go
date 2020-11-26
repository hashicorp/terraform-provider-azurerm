package azure

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func SplitSku(sku string) (string, int32, error) {
	skuParts := strings.Split(sku, "_")

	if len(skuParts) != 2 {
		return "", -1, fmt.Errorf("sku_name(%s) is not formatted properly.", sku)
	}

	capacity, err := strconv.Atoi(skuParts[1])
	if err != nil {
		return "", -1, fmt.Errorf("%s in sku_name is not a valid value.", skuParts[1])
	}

	return skuParts[0], int32(capacity), nil
}

// MinCapacitySkuNameInSlice returns a SchemaValidateFunc which tests if the provided value
// is of type string and matches the value of an element in the valid slice
// will test with in lower case if ignoreCase is true will also validate if the
// capacity if above passed minCapacity value
func MinCapacitySkuNameInSlice(valid []string, minCapacity int32, ignoreCase bool) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		name, capacity, err := SplitSku(v)
		if err != nil {
			es = append(es, err)
			return
		}

		for _, str := range valid {
			if name == str || (ignoreCase && strings.EqualFold(name, str)) {
				if capacity < minCapacity {
					es = append(es, fmt.Errorf("expected %s capacity value to be greater that %d, got %d", k, minCapacity, capacity))
				}
				return
			}
		}

		es = append(es, fmt.Errorf("expected %s to be one of %v, got %s", k, valid, name))
		return
	}
}
