package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type IntegrationAccountMapId struct {
	SubscriptionId         string
	ResourceGroup          string
	IntegrationAccountName string
	MapName                string
}

func NewIntegrationAccountMapID(subscriptionId, resourceGroup, integrationAccountName, mapName string) IntegrationAccountMapId {
	return IntegrationAccountMapId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		IntegrationAccountName: integrationAccountName,
		MapName:                mapName,
	}
}

func (id IntegrationAccountMapId) String() string {
	segments := []string{
		fmt.Sprintf("Map Name %q", id.MapName),
		fmt.Sprintf("Integration Account Name %q", id.IntegrationAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Integration Account Map", segmentsStr)
}

func (id IntegrationAccountMapId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/integrationAccounts/%s/maps/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IntegrationAccountName, id.MapName)
}

// IntegrationAccountMapID parses a IntegrationAccountMap ID into an IntegrationAccountMapId struct
func IntegrationAccountMapID(input string) (*IntegrationAccountMapId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := IntegrationAccountMapId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.IntegrationAccountName, err = id.PopSegment("integrationAccounts"); err != nil {
		return nil, err
	}
	if resourceId.MapName, err = id.PopSegment("maps"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
