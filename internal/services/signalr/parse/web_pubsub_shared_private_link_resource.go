package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type WebPubsubSharedPrivateLinkResourceId struct {
	SubscriptionId                string
	ResourceGroup                 string
	WebPubSubName                 string
	SharedPrivateLinkResourceName string
}

func NewWebPubsubSharedPrivateLinkResourceID(subscriptionId, resourceGroup, webPubSubName, sharedPrivateLinkResourceName string) WebPubsubSharedPrivateLinkResourceId {
	return WebPubsubSharedPrivateLinkResourceId{
		SubscriptionId:                subscriptionId,
		ResourceGroup:                 resourceGroup,
		WebPubSubName:                 webPubSubName,
		SharedPrivateLinkResourceName: sharedPrivateLinkResourceName,
	}
}

func (id WebPubsubSharedPrivateLinkResourceId) String() string {
	segments := []string{
		fmt.Sprintf("Shared Private Link Resource Name %q", id.SharedPrivateLinkResourceName),
		fmt.Sprintf("Web Pub Sub Name %q", id.WebPubSubName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Web Pubsub Shared Private Link Resource", segmentsStr)
}

func (id WebPubsubSharedPrivateLinkResourceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.SignalRService/webPubSub/%s/sharedPrivateLinkResources/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WebPubSubName, id.SharedPrivateLinkResourceName)
}

// WebPubsubSharedPrivateLinkResourceID parses a WebPubsubSharedPrivateLinkResource ID into an WebPubsubSharedPrivateLinkResourceId struct
func WebPubsubSharedPrivateLinkResourceID(input string) (*WebPubsubSharedPrivateLinkResourceId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := WebPubsubSharedPrivateLinkResourceId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.WebPubSubName, err = id.PopSegment("webPubSub"); err != nil {
		return nil, err
	}
	if resourceId.SharedPrivateLinkResourceName, err = id.PopSegment("sharedPrivateLinkResources"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
