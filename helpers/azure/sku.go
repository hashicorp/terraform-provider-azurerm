package azure

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-07-01/compute"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

func ExpandOrchestratedVirtualMachineScaleSetSku(input string) (*compute.Sku, error) {
	skuParts := strings.Split(input, "_")

	if len(skuParts) < 3 || strings.Contains(input, "__") || strings.Contains(input, " ") {
		return nil, fmt.Errorf("'sku_name'(%q) is not formatted properly.", input)
	}

	// capacity is always the last argument in the string
	index := (len(skuParts) - 1)

	capacity, err := strconv.Atoi(skuParts[index])
	if err != nil || capacity < 0 || capacity > 1000 {
		return nil, fmt.Errorf("%q in 'sku_name' is not a valid value for capacity.", skuParts[index])
	}

	skuName := input[:len(input)-(len(skuParts[index])+1)]

	sku := &compute.Sku{
		Name:     utils.String(skuName),
		Capacity: utils.Int64(int64(capacity)),
		Tier:     utils.String("Standard"),
	}

	return sku, nil
}

func FlattenOrchestratedVirtualMachineScaleSetSku(input *compute.Sku) (*string, error) {
	if input != nil {
		if input.Name != nil && input.Capacity != nil {
			skuName := fmt.Sprintf("Standard_%s_%d", *input.Name, *input.Capacity)
			if strings.HasPrefix(strings.ToLower(*input.Name), "standard") {
				skuName = fmt.Sprintf("%s_%d", *input.Name, *input.Capacity)
			}

			return &skuName, nil
		}
	}

	return nil, fmt.Errorf("Sku struct 'name' or 'capacity' are nil")
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
