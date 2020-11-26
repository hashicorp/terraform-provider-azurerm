package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ApplicationGroupId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewApplicationGroupID(subscriptionId, resourceGroup, name string) ApplicationGroupId {
	return ApplicationGroupId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id ApplicationGroupId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DesktopVirtualization/applicationGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

func ApplicationGroupID(input string) (*ApplicationGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ApplicationGroupId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	applicationGroupsKey := "applicationGroups"
	// find the correct casing for the `applicationGroups` segment
	for key := range id.Path {
		if strings.EqualFold(key, applicationGroupsKey) {
			applicationGroupsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(applicationGroupsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
