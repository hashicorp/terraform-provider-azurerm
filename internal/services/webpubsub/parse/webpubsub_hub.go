package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type WebPubsubHubId struct {
	SubscriptionId string
	ResourceGroup  string
	WebPubsubName  string
	HubName        string
}

func NewWebPubsubHubID(subscriptionId, resourceGroup, webPubsubName, hubname string) WebPubsubHubId {
	return WebPubsubHubId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		WebPubsubName:  webPubsubName,
		HubName:        hubname,
	}
}

func (id WebPubsubHubId) String() string {
	segments := []string{
		fmt.Sprintf("HubName %q", id.HubName),
		fmt.Sprintf("WebPubsub Name %q", id.WebPubsubName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Web Pubsub Hub", segmentsStr)
}

func (id WebPubsubHubId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.SignalRService/webPubSub/%s/hubs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WebPubsubName, id.HubName)
}

func WebPubsubHubID(input string) (*WebPubsubHubId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := WebPubsubHubId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missiong the 'resourceGroups' element")
	}

	if resourceId.WebPubsubName, err = id.PopSegment("webPubSub"); err != nil {
		return nil, err
	}
	if resourceId.HubName, err = id.PopSegment("hubs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
