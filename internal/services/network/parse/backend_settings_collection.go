package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type BackendSettingsCollectionId struct {
	SubscriptionId                string
	ResourceGroup                 string
	ApplicationGatewayName        string
	BackendSettingsCollectionName string
}

func NewBackendSettingsCollectionID(subscriptionId, resourceGroup, applicationGatewayName, backendSettingsCollectionName string) BackendSettingsCollectionId {
	return BackendSettingsCollectionId{
		SubscriptionId:                subscriptionId,
		ResourceGroup:                 resourceGroup,
		ApplicationGatewayName:        applicationGatewayName,
		BackendSettingsCollectionName: backendSettingsCollectionName,
	}
}

func (id BackendSettingsCollectionId) String() string {
	segments := []string{
		fmt.Sprintf("Backend Settings Collection Name %q", id.BackendSettingsCollectionName),
		fmt.Sprintf("Application Gateway Name %q", id.ApplicationGatewayName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Backend Settings Collection", segmentsStr)
}

func (id BackendSettingsCollectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/applicationGateways/%s/backendSettingsCollection/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ApplicationGatewayName, id.BackendSettingsCollectionName)
}

// BackendSettingsCollectionID parses a BackendSettingsCollection ID into an BackendSettingsCollectionId struct
func BackendSettingsCollectionID(input string) (*BackendSettingsCollectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := BackendSettingsCollectionId{
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
	if resourceId.BackendSettingsCollectionName, err = id.PopSegment("backendSettingsCollection"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
