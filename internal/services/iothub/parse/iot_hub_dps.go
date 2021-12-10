package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type IotHubDpsId struct {
	SubscriptionId          string
	ResourceGroup           string
	ProvisioningServiceName string
}

func NewIotHubDpsID(subscriptionId, resourceGroup, provisioningServiceName string) IotHubDpsId {
	return IotHubDpsId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		ProvisioningServiceName: provisioningServiceName,
	}
}

func (id IotHubDpsId) String() string {
	segments := []string{
		fmt.Sprintf("Provisioning Service Name %q", id.ProvisioningServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Iot Hub Dps", segmentsStr)
}

func (id IotHubDpsId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Devices/provisioningServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProvisioningServiceName)
}

// IotHubDpsID parses a IotHubDps ID into an IotHubDpsId struct
func IotHubDpsID(input string) (*IotHubDpsId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := IotHubDpsId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ProvisioningServiceName, err = id.PopSegment("provisioningServices"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
