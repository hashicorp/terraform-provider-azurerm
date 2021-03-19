package parse

import (
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ManagementGroupSubscriptionAssociationId struct {
	ManagementGroup string
	SubscriptionId  string
}

func NewManagementGroupSubscriptionAssociationID(managementGroupName string, subscriptionId string) ManagementGroupSubscriptionAssociationId {
	return ManagementGroupSubscriptionAssociationId{
		ManagementGroup: managementGroupName,
		SubscriptionId:  subscriptionId,
	}
}

func (r ManagementGroupSubscriptionAssociationId) ID() string {
	managementGroupSubscriptionAssociationFmt := "/managementGroup/%s/subscription/%s"
	return fmt.Sprintf(managementGroupSubscriptionAssociationFmt, r.ManagementGroup, r.SubscriptionId)
}

func ManagementGroupSubscriptionAssociationID(input string) (*ManagementGroupSubscriptionAssociationId, error) {
	id, err := azure.ParseAzureResourceIDWithoutSubscription(input)
	if err != nil {
		return nil, err
	}

	managementGroup, err := id.PopSegment("managementGroup")
	if err != nil {
		return nil, err
	}

	subscriptionId, err := id.PopSegment("subscription")
	if err != nil {
		return nil, err
	}
	if _, err := uuid.ParseUUID(subscriptionId); err != nil {
		return nil, fmt.Errorf("expected subscription ID to be UUID, got %q", subscriptionId)
	}

	return &ManagementGroupSubscriptionAssociationId{
		ManagementGroup: managementGroup,
		SubscriptionId:  subscriptionId,
	}, nil
}
