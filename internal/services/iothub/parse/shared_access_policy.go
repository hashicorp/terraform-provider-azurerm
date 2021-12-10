package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SharedAccessPolicyId struct {
	SubscriptionId string
	ResourceGroup  string
	IotHubName     string
	IotHubKeyName  string
}

func NewSharedAccessPolicyID(subscriptionId, resourceGroup, iotHubName, iotHubKeyName string) SharedAccessPolicyId {
	return SharedAccessPolicyId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		IotHubName:     iotHubName,
		IotHubKeyName:  iotHubKeyName,
	}
}

func (id SharedAccessPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Iot Hub Key Name %q", id.IotHubKeyName),
		fmt.Sprintf("Iot Hub Name %q", id.IotHubName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Shared Access Policy", segmentsStr)
}

func (id SharedAccessPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Devices/IotHubs/%s/IotHubKeys/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IotHubName, id.IotHubKeyName)
}

// SharedAccessPolicyID parses a SharedAccessPolicy ID into an SharedAccessPolicyId struct
func SharedAccessPolicyID(input string) (*SharedAccessPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SharedAccessPolicyId{
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
	if resourceId.IotHubKeyName, err = id.PopSegment("IotHubKeys"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
