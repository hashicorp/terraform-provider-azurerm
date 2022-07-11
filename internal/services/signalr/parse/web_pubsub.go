package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type WebPubsubId struct {
	SubscriptionId string
	ResourceGroup  string
	WebPubSubName  string
}

func NewWebPubsubID(subscriptionId, resourceGroup, webPubSubName string) WebPubsubId {
	return WebPubsubId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		WebPubSubName:  webPubSubName,
	}
}

func (id WebPubsubId) String() string {
	segments := []string{
		fmt.Sprintf("Web Pub Sub Name %q", id.WebPubSubName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Web Pubsub", segmentsStr)
}

func (id WebPubsubId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.SignalRService/WebPubSub/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WebPubSubName)
}

// WebPubsubID parses a WebPubsub ID into an WebPubsubId struct
func WebPubsubID(input string) (*WebPubsubId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := WebPubsubId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.WebPubSubName, err = id.PopSegment("WebPubSub"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
