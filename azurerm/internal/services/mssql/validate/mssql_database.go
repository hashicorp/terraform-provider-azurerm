package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
)

func ValidateMsSqlDatabaseID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := parse.MsSqlDatabaseID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a MsSql Database resource id: %v", k, err))
	}

	return warnings, errors
}

func ValidateMsSqlDatabaseAutoPauseDelay(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be integer", k))
		return warnings, errors
	}
	min := 60
	max := 10080
	if (v < min || v > max) && v%10 != 0 && v != -1 {
		errors = append(errors, fmt.Errorf("expected %s to be in the range (%d - %d) and divisible by 10 or -1, got %d", k, min, max, v))
		return warnings, errors
	}

	return warnings, errors
}

func ValidateMsSqlDBMinCapacity(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(float64)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be float", k))
		return warnings, errors
	}

	valid := []float64{0.5, 0.75, 1, 1.25, 1.5, 1.75, 2}

	for _, validValue := range valid {
		if v == validValue {
			return warnings, errors
		}
	}

	errors = append(errors, fmt.Errorf("expected %s to be one of %v, got %f", k, valid, v))
	return warnings, errors
}

func ValidateMsSqlDBSkuName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`(?i)(^((GP(_S)?|BC)_(Gen4|Gen5)_(2|4|6|8|10|12|14|16|18|20|24|32|40|80))|(HS_(Gen4|Gen5)_(1|2|3|4|5|6|7|8|9|10|16|24))|Basic|Standard|Premium|ElasticPool|S(0|1|2|3|4|6|7|9|12)|P(1|2|4|6|11|15)$)`),
		`This is not a valid sku name. For example, a valid sku name is 'GP_Gen5_2','HS_Gen4_1','BC_Gen5_2', 'ElasticPool', 'Basic', 'Standard', 'Premium' and etc.`,
	)
}
