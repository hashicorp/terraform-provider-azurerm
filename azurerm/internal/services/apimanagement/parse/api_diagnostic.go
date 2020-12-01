package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ApiDiagnosticId struct {
	ResourceGroup  string
	ServiceName    string
	ApiName        string
	DiagnosticName string
}

func ApiDiagnosticID(input string) (*ApiDiagnosticId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Api Management Diagnostic ID %q: %+v", input, err)
	}

	diagnostic := ApiDiagnosticId{
		ResourceGroup: id.ResourceGroup,
	}

	if diagnostic.ServiceName, err = id.PopSegment("service"); err != nil {
		return nil, err
	}

	if diagnostic.ApiName, err = id.PopSegment("apis"); err != nil {
		return nil, err
	}

	if diagnostic.DiagnosticName, err = id.PopSegment("diagnostics"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &diagnostic, nil
}
