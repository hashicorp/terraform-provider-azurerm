package signalr

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SignalRId struct {
	SubscriptionId string
	ResourceGroup  string
	SignalRName    string
}

func NewSignalRID(subscriptionId, resourceGroup, signalRName string) SignalRId {
	return SignalRId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SignalRName:    signalRName,
	}
}

func (id SignalRId) String() string {
	segments := []string{
		fmt.Sprintf("Signal R Name %q", id.SignalRName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Signal R", segmentsStr)
}

func (id SignalRId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.SignalRService/SignalR/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SignalRName)
}

// ParseSignalRID parses a SignalR ID into an SignalRId struct
func ParseSignalRID(input string) (*SignalRId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SignalRId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SignalRName, err = id.PopSegment("SignalR"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseSignalRIDInsensitively parses an SignalR ID into an SignalRId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseSignalRID method should be used instead for validation etc.
func ParseSignalRIDInsensitively(input string) (*SignalRId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SignalRId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'SignalR' segment
	SignalRKey := "SignalR"
	for key := range id.Path {
		if strings.EqualFold(key, SignalRKey) {
			SignalRKey = key
			break
		}
	}
	if resourceId.SignalRName, err = id.PopSegment(SignalRKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
