package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type IntegrationAccountPartnerId struct {
	SubscriptionId         string
	ResourceGroup          string
	IntegrationAccountName string
	PartnerName            string
}

func NewIntegrationAccountPartnerID(subscriptionId, resourceGroup, integrationAccountName, partnerName string) IntegrationAccountPartnerId {
	return IntegrationAccountPartnerId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		IntegrationAccountName: integrationAccountName,
		PartnerName:            partnerName,
	}
}

func (id IntegrationAccountPartnerId) String() string {
	segments := []string{
		fmt.Sprintf("Partner Name %q", id.PartnerName),
		fmt.Sprintf("Integration Account Name %q", id.IntegrationAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Integration Account Partner", segmentsStr)
}

func (id IntegrationAccountPartnerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/integrationAccounts/%s/partners/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IntegrationAccountName, id.PartnerName)
}

// IntegrationAccountPartnerID parses a IntegrationAccountPartner ID into an IntegrationAccountPartnerId struct
func IntegrationAccountPartnerID(input string) (*IntegrationAccountPartnerId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := IntegrationAccountPartnerId{
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
	if resourceId.PartnerName, err = id.PopSegment("partners"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
