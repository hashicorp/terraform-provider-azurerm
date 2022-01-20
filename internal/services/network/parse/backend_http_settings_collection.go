package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type BackendHttpSettingsCollectionId struct {
	SubscriptionId                    string
	ResourceGroup                     string
	ApplicationGatewayName            string
	BackendHttpSettingsCollectionName string
}

func NewBackendHttpSettingsCollectionID(subscriptionId, resourceGroup, applicationGatewayName, backendHttpSettingsCollectionName string) BackendHttpSettingsCollectionId {
	return BackendHttpSettingsCollectionId{
		SubscriptionId:                    subscriptionId,
		ResourceGroup:                     resourceGroup,
		ApplicationGatewayName:            applicationGatewayName,
		BackendHttpSettingsCollectionName: backendHttpSettingsCollectionName,
	}
}

func (id BackendHttpSettingsCollectionId) String() string {
	segments := []string{
		fmt.Sprintf("Backend Http Settings Collection Name %q", id.BackendHttpSettingsCollectionName),
		fmt.Sprintf("Application Gateway Name %q", id.ApplicationGatewayName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Backend Http Settings Collection", segmentsStr)
}

func (id BackendHttpSettingsCollectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/applicationGateways/%s/backendHttpSettingsCollection/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ApplicationGatewayName, id.BackendHttpSettingsCollectionName)
}

// BackendHttpSettingsCollectionID parses a BackendHttpSettingsCollection ID into an BackendHttpSettingsCollectionId struct
func BackendHttpSettingsCollectionID(input string) (*BackendHttpSettingsCollectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := BackendHttpSettingsCollectionId{
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
	if resourceId.BackendHttpSettingsCollectionName, err = id.PopSegment("backendHttpSettingsCollection"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
