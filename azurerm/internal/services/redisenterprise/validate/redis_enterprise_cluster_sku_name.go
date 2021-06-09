package validate

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/redisenterprise/mgmt/2021-03-01/redisenterprise"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

// RedisEnterpriseClusterSkuName - validates if passed input string contains a valid Redis Enterprise Cluster Sku
func RedisEnterpriseClusterSkuName(v interface{}, k string) (warnings []string, errors []error) {
	validSku := false
	validCapacity := false

	value, ok := v.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	skuParts := strings.Split(value, "-")
	validSkus := getValidRedisEnterpriseClusterSkus()
	validValues := "2, 4, 6, ..."
	// Validate the SKU Name section
	for _, str := range validSkus {
		if skuParts[0] == str {
			validSku = true
			break
		}
	}

	isFlash := strings.Contains(skuParts[0], "Flash")

	if isFlash {
		validValues = "3, 9, 15, ..."
	}

	if len(skuParts) > 1 {
		i, err := strconv.ParseInt(skuParts[1], 10, 32)
		if err != nil {
			errors = append(errors, fmt.Errorf("expected %q %q segment to fit into an int32 value, got %q", k, "capacity", skuParts[1]))
			return warnings, errors
		}

		skuCapacity := int32(i)

		if isFlash {
			if skuCapacity%6 == 3 {
				validCapacity = true
			}
		} else if skuCapacity%2 == 0 {
			validCapacity = true
		}
	}

	if !validSku {
		errors = append(errors, fmt.Errorf("expected %q %q segment to be one of [%s], got %q", k, "name", azure.QuotedStringSlice(validSkus), value))
	}

	if !validCapacity {
		errors = append(errors, fmt.Errorf("expected %q %q segment to be one of [%s], got %q", k, "capacity", validValues, value))
	}

	return warnings, errors
}

func getValidRedisEnterpriseClusterSkus() []string {
	return []string{
		string(redisenterprise.EnterpriseE10),
		string(redisenterprise.EnterpriseE20),
		string(redisenterprise.EnterpriseE50),
		string(redisenterprise.EnterpriseE100),
		string(redisenterprise.EnterpriseFlashF300),
		string(redisenterprise.EnterpriseFlashF700),
		string(redisenterprise.EnterpriseFlashF1500),
	}
}
