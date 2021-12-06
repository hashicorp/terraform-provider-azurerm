package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type BackendAddressPoolId struct {
	SubscriptionId         string
	ResourceGroup          string
	ApplicationGatewayName string
	Name                   string
}

func NewBackendAddressPoolID(subscriptionId, resourceGroup, applicationGatewayName, name string) BackendAddressPoolId {
	return BackendAddressPoolId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		ApplicationGatewayName: applicationGatewayName,
		Name:                   name,
	}
}

func (id BackendAddressPoolId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Application Gateway Name %q", id.ApplicationGatewayName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Backend Address Pool", segmentsStr)
}

func (id BackendAddressPoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/applicationGateways/%s/backendAddressPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ApplicationGatewayName, id.Name)
}

// BackendAddressPoolID parses a BackendAddressPool ID into an BackendAddressPoolId struct
func BackendAddressPoolID(input string) (*BackendAddressPoolId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := BackendAddressPoolId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ApplicationGatewayName, err = id.PopSegment("applicationGateways"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("backendAddressPools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
