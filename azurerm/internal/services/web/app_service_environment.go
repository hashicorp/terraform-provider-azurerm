package web

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServiceEnvironmentResourceID struct {
	ResourceGroup string
	Name          string
}

func ParseAppServiceEnvironmentID(input string) (*AppServiceEnvironmentResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service Environment ID %q: %+v", input, err)
	}

	appServiceEnvironment := AppServiceEnvironmentResourceID{
		ResourceGroup: id.ResourceGroup,
	}

	if appServiceEnvironment.Name, err = id.PopSegment("hostingEnvironments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &appServiceEnvironment, nil
}

// ValidateAppServiceID validates that the specified ID is a valid App Service ID
func ValidateAppServiceEnvironmentID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := ParseAppServiceEnvironmentID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}

func validateAppServiceEnvironmentName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[0-9a-zA-Z][-0-9a-zA-Z]{0,61}[0-9a-zA-Z]$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes up to 60 characters in length, and must start and end in an alphanumeric", k))
	}

	return warnings, errors
}

func validateAppServiceEnvironmentPricingTier(v interface{}, k string) (warnings []string, errors []error) {
	tier := v.(string)

	valid := []string{"I1", "I2", "I3"}

	for _, val := range valid {
		if val == tier {
			return
		}
	}
	errors = append(errors, fmt.Errorf("pricing_tier must be one of %q", valid))
	return warnings, errors
}
