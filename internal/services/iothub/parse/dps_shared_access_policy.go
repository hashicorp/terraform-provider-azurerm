package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type DpsSharedAccessPolicyId struct {
	SubscriptionId          string
	ResourceGroup           string
	ProvisioningServiceName string
	KeyName                 string
}

func NewDpsSharedAccessPolicyID(subscriptionId, resourceGroup, provisioningServiceName, keyName string) DpsSharedAccessPolicyId {
	return DpsSharedAccessPolicyId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		ProvisioningServiceName: provisioningServiceName,
		KeyName:                 keyName,
	}
}

func (id DpsSharedAccessPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Key Name %q", id.KeyName),
		fmt.Sprintf("Provisioning Service Name %q", id.ProvisioningServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Dps Shared Access Policy", segmentsStr)
}

func (id DpsSharedAccessPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Devices/provisioningServices/%s/keys/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProvisioningServiceName, id.KeyName)
}

// DpsSharedAccessPolicyID parses a DpsSharedAccessPolicy ID into an DpsSharedAccessPolicyId struct
func DpsSharedAccessPolicyID(input string) (*DpsSharedAccessPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DpsSharedAccessPolicyId{
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
	if resourceId.KeyName, err = id.PopSegment("keys"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
