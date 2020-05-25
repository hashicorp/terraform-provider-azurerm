package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ApiManagementDiagnosticId struct {
	ResourceGroup string
	ServiceName   string
	Name          string
}

func ApiManagementDiagnosticID(input string) (*ApiManagementDiagnosticId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Api Management Diagnostic ID %q: %+v", input, err)
	}

	diagnostic := ApiManagementDiagnosticId{
		ResourceGroup: id.ResourceGroup,
	}

	if diagnostic.ServiceName, err = id.PopSegment("service"); err != nil {
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
