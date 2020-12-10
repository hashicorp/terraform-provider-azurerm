package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SpringCloudServiceId struct {
	SubscriptionId string
	ResourceGroup  string
	SpringName     string
}

func NewSpringCloudServiceID(subscriptionId, resourceGroup, springName string) SpringCloudServiceId {
	return SpringCloudServiceId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SpringName:     springName,
	}
}

func (id SpringCloudServiceId) String() string {
	segments := []string{
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud Service", segmentsStr)
}

func (id SpringCloudServiceId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/Spring/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName)
}

// SpringCloudServiceID parses a SpringCloudService ID into an SpringCloudServiceId struct
func SpringCloudServiceID(input string) (*SpringCloudServiceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudServiceId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SpringName, err = id.PopSegment("Spring"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
