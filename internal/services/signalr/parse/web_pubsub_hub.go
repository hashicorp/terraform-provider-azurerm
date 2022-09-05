package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type WebPubsubHubId struct {
	SubscriptionId string
	ResourceGroup  string
	WebPubSubName  string
	HubName        string
}

func NewWebPubsubHubID(subscriptionId, resourceGroup, webPubSubName, hubName string) WebPubsubHubId {
	return WebPubsubHubId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		WebPubSubName:  webPubSubName,
		HubName:        hubName,
	}
}

func (id WebPubsubHubId) String() string {
	segments := []string{
		fmt.Sprintf("Hub Name %q", id.HubName),
		fmt.Sprintf("Web Pub Sub Name %q", id.WebPubSubName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Web Pubsub Hub", segmentsStr)
}

func (id WebPubsubHubId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.SignalRService/webPubSub/%s/hubs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WebPubSubName, id.HubName)
}

// WebPubsubHubID parses a WebPubsubHub ID into an WebPubsubHubId struct
func WebPubsubHubID(input string) (*WebPubsubHubId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
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
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.WebPubSubName, err = id.PopSegment("webPubSub"); err != nil {
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

// WebPubsubHubIDInsensitively parses an WebPubsubHub ID into an WebPubsubHubId struct, insensitively
// This should only be used to parse an ID for rewriting, the WebPubsubHubID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func WebPubsubHubIDInsensitively(input string) (*WebPubsubHubId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
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
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'webPubSub' segment
	webPubSubKey := "webPubSub"
	for key := range id.Path {
		if strings.EqualFold(key, webPubSubKey) {
			webPubSubKey = key
			break
		}
	}
	if resourceId.WebPubSubName, err = id.PopSegment(webPubSubKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'hubs' segment
	hubsKey := "hubs"
	for key := range id.Path {
		if strings.EqualFold(key, hubsKey) {
			hubsKey = key
			break
		}
	}
	if resourceId.HubName, err = id.PopSegment(hubsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
