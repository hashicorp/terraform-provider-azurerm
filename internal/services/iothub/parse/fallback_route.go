package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FallbackRouteId struct {
	SubscriptionId    string
	ResourceGroup     string
	IotHubName        string
	FallbackRouteName string
}

func NewFallbackRouteID(subscriptionId, resourceGroup, iotHubName, fallbackRouteName string) FallbackRouteId {
	return FallbackRouteId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		IotHubName:        iotHubName,
		FallbackRouteName: fallbackRouteName,
	}
}

func (id FallbackRouteId) String() string {
	segments := []string{
		fmt.Sprintf("Fallback Route Name %q", id.FallbackRouteName),
		fmt.Sprintf("Iot Hub Name %q", id.IotHubName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Fallback Route", segmentsStr)
}

func (id FallbackRouteId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Devices/IotHubs/%s/FallbackRoute/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IotHubName, id.FallbackRouteName)
}

// FallbackRouteID parses a FallbackRoute ID into an FallbackRouteId struct
func FallbackRouteID(input string) (*FallbackRouteId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FallbackRouteId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.IotHubName, err = id.PopSegment("IotHubs"); err != nil {
		return nil, err
	}
	if resourceId.FallbackRouteName, err = id.PopSegment("FallbackRoute"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
