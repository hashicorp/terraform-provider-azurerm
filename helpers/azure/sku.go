package azure

import (
	"fmt"
	"strconv"
	"strings"
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
