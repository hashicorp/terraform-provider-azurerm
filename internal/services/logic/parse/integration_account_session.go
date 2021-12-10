package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type IntegrationAccountSessionId struct {
	SubscriptionId         string
	ResourceGroup          string
	IntegrationAccountName string
	SessionName            string
}

func NewIntegrationAccountSessionID(subscriptionId, resourceGroup, integrationAccountName, sessionName string) IntegrationAccountSessionId {
	return IntegrationAccountSessionId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		IntegrationAccountName: integrationAccountName,
		SessionName:            sessionName,
	}
}

func (id IntegrationAccountSessionId) String() string {
	segments := []string{
		fmt.Sprintf("Session Name %q", id.SessionName),
		fmt.Sprintf("Integration Account Name %q", id.IntegrationAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Integration Account Session", segmentsStr)
}

func (id IntegrationAccountSessionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/integrationAccounts/%s/sessions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IntegrationAccountName, id.SessionName)
}

// IntegrationAccountSessionID parses a IntegrationAccountSession ID into an IntegrationAccountSessionId struct
func IntegrationAccountSessionID(input string) (*IntegrationAccountSessionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := IntegrationAccountSessionId{
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
	if resourceId.SessionName, err = id.PopSegment("sessions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
