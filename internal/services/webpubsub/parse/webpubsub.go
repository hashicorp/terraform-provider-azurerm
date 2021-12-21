package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type WebPubSubId struct {
	SubscriptionId  string
	ResourceGroupId string
	Name            string
}

func NewWebPubSubID(subscriptionId, resourceGroup, name string) WebPubSubId {
	return WebPubSubId{
		SubscriptionId:  subscriptionId,
		ResourceGroupId: resourceGroup,
		Name:            name,
	}
}

func (id WebPubSubId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroupId),
	}
	segmentsStr := strings.Join(segments, "/")
	return fmt.Sprintf("%s: (%s)", "Web PubSub", segmentsStr)
}

func (id WebPubSubId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.SignalRService/WebPubSub/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupId, id.Name)
}

func WebPubSubID(input string) (*WebPubSubId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := WebPubSubId{
		SubscriptionId:  id.SubscriptionID,
		ResourceGroupId: id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscription' element")
	}

	if resourceId.ResourceGroupId == "" {
		return nil, fmt.Errorf("ID was missiong the 'resourceGroup' element")
	}

	if resourceId.Name, err = id.PopSegment("WebPubSub"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
