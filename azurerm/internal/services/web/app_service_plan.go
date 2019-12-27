package web

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServicePlanResourceID struct {
	ResourceGroup string
	Name          string
}

func ParseAppServicePlanID(input string) (*AppServicePlanResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service Plan ID %q: %+v", input, err)
	}

	appServicePlan := AppServicePlanResourceID{
		ResourceGroup: id.ResourceGroup,
	}

	if appServicePlan.Name, err = id.PopSegment("serverfarms"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &appServicePlan, nil
}

// ValidateAppServicePlanID validates that the specified ID is a valid App Service Plan ID
func ValidateAppServicePlanID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := ParseAppServicePlanID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}
