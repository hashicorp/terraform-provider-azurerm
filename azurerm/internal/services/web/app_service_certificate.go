package web

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServiceCertificateResourceID struct {
	ResourceGroup string
	Name          string
}

func ParseAppServiceCertificateID(input string) (*AppServiceCertificateResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service Certificate ID %q: %+v", input, err)
	}

	certificate := AppServiceCertificateResourceID{
		ResourceGroup: id.ResourceGroup,
	}

	if certificate.Name, err = id.PopSegment("certificates"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &certificate, nil
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
