package validate

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/digitaltwins/parse"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DigitaltwinsName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]{1,61}[A-Za-z0-9]$`),
		`Name contains invalid characters or exceeds allowed length.`,
	)
}

func DigitaltwinsID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := parse.DigitalTwinsID(v); err != nil {
		errors = append(errors, fmt.Errorf("can not parse %q as a Digital Twins resource id: %v", k, err))
	}

	return warnings, errors
}
