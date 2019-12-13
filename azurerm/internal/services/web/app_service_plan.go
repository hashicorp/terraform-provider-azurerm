package web

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServicePlanResourceID struct {
	Base azure.ResourceID

	Name string
}

func ParseAppServicePlanID(input string) (*AppServicePlanResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service Plan ID %q: %+v", input, err)
	}

	group := AppServicePlanResourceID{
		Base: *id,
		Name: id.Path["serverfarms"],
	}

	if group.Name == "" {
		return nil, fmt.Errorf("ID was missing the `serverfarms` element")
	}

	pathWithoutElements := group.Base.Path
	delete(pathWithoutElements, "serverfarms")
	if len(pathWithoutElements) != 0 {
		return nil, fmt.Errorf("ID contained more segments than a Resource ID requires: %q", input)
	}

	return &group, nil
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
