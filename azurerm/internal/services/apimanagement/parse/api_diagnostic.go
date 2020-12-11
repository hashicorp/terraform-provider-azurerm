package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ApiDiagnosticId struct {
	SubscriptionId string
	ResourceGroup  string
	ServiceName    string
	ApiName        string
	DiagnosticName string
}

func NewApiDiagnosticID(subscriptionId, resourceGroup, serviceName, apiName, diagnosticName string) ApiDiagnosticId {
	return ApiDiagnosticId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServiceName:    serviceName,
		ApiName:        apiName,
		DiagnosticName: diagnosticName,
	}
}

func (id ApiDiagnosticId) String() string {
	segments := []string{
		fmt.Sprintf("Diagnostic Name %q", id.DiagnosticName),
		fmt.Sprintf("Api Name %q", id.ApiName),
		fmt.Sprintf("Service Name %q", id.ServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Api Diagnostic", segmentsStr)
}

func (id ApiDiagnosticId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/apis/%s/diagnostics/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.ApiName, id.DiagnosticName)
}

// ApiDiagnosticID parses a ApiDiagnostic ID into an ApiDiagnosticId struct
func ApiDiagnosticID(input string) (*ApiDiagnosticId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ApiDiagnosticId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServiceName, err = id.PopSegment("service"); err != nil {
		return nil, err
	}
	if resourceId.ApiName, err = id.PopSegment("apis"); err != nil {
		return nil, err
	}
	if resourceId.DiagnosticName, err = id.PopSegment("diagnostics"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
