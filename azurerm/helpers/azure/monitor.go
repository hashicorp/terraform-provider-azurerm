package azure

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// ValidateThreshold checks that a threshold value is between 0 and 10000
// and is a whole number. The azure-sdk-for-go expects this value to be a float64
// but the user validation rules want an integer.
func ValidateThreshold(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(float64)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be float64", k))
	}

	if v != float64(int64(v)) {
		errors = append(errors, fmt.Errorf("%q must be a whole number", k))
	}

	if v < 0 || v > 10000 {
		errors = append(errors, fmt.Errorf("%q must be between 0 and 10000 inclusive", k))
	}

	return warnings, errors
}
