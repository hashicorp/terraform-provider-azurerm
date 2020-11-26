package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ApiManagementApiDiagnosticId struct {
	ResourceGroup string
	ServiceName   string
	ApiName       string
	Name          string
}

func ApiManagementApiDiagnosticID(input string) (*ApiManagementApiDiagnosticId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Api Management Diagnostic ID %q: %+v", input, err)
	}

	diagnostic := ApiManagementApiDiagnosticId{
		ResourceGroup: id.ResourceGroup,
	}

	if diagnostic.ServiceName, err = id.PopSegment("service"); err != nil {
		return nil, err
	}

	if diagnostic.ApiName, err = id.PopSegment("apis"); err != nil {
		return nil, err
	}

	if diagnostic.Name, err = id.PopSegment("diagnostics"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &diagnostic, nil
}
