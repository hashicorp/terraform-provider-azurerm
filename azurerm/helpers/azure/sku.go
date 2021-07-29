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

func SplitOrchestratedVirtualMachineScaleSetSku(sku string) (string, int32, error) {
	skuParts := strings.Split(sku, "_")

	if len(skuParts) < 3 || strings.Contains(sku, "__") || strings.Contains(sku, " ") {
		return "", -1, fmt.Errorf("'sku_name'(%q) is not formatted properly.", sku)
	}

	// capacity is always the last argument in the string
	index := (len(skuParts) - 1)

	capacity, err := strconv.Atoi(skuParts[index])
	if err != nil || capacity < 0 || capacity > 1000 {
		return "", -1, fmt.Errorf("%q in 'sku_name' is not a valid value for capacity.", skuParts[index])
	}

	skuName := sku[:len(sku)-(len(skuParts[index])+1)]

	return skuName, int32(capacity), nil
}

func ValidateOrchestratedVirtualMachineScaleSetSku(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	skuParts := strings.Split(v, "_")

	if len(skuParts) < 3 || strings.Contains(v, "__") || strings.Contains(v, " ") {
		errors = append(errors, fmt.Errorf("%q(%q) is not formatted properly.", key, v))
	}

	index := (len(skuParts) - 1)

	capacity, err := strconv.Atoi(skuParts[index])
	if err != nil || capacity < 0 || capacity > 1000 {
		errors = append(errors, fmt.Errorf("%q in %q is not a valid value for capacity.", skuParts[index], key))
	}

	return
}
