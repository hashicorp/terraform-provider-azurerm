package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parsers"
)

func StorageContainerResourceManagerID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parsers.StorageContainerResourceManagerID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}

func StorageContainerLegalHoldTag() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-z0-9]{3,23}$"),
		`Each tag should be 3 to 23 alphanumeric characters and is normalized to lower case.`,
	)
}
