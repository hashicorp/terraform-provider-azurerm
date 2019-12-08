package web

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServiceCertificateResourceID struct {
	Base azure.ResourceID

	Name string
}

func ParseAppServiceCertificateID(input string) (*AppServiceCertificateResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service Certificate ID %q: %+v", input, err)
	}

	group := AppServiceCertificateResourceID{
		Base: *id,
		Name: id.Path["certificates"],
	}

	if group.Name == "" {
		return nil, fmt.Errorf("ID was missing the `certificates` element")
	}

	pathWithoutElements := group.Base.Path
	delete(pathWithoutElements, "certificates")
	if len(pathWithoutElements) != 0 {
		return nil, fmt.Errorf("ID contained more segments than a Resource ID requires: %q", input)
	}

	return &group, nil
}

// ValidateAppServiceCertificateID validates that the specified ID is a valid App Service Certificate ID
func ValidateAppServiceCertificateID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := ParseAppServiceCertificateID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}
