package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ProductGroupId struct {
	SubscriptionId string
	ResourceGroup  string
	ServiceName    string
	ProductName    string
	GroupName      string
}

func NewProductGroupID(subscriptionId, resourceGroup, serviceName, productName, groupName string) ProductGroupId {
	return ProductGroupId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServiceName:    serviceName,
		ProductName:    productName,
		GroupName:      groupName,
	}
}

func (id ProductGroupId) String() string {
	segments := []string{
		fmt.Sprintf("Group Name %q", id.GroupName),
		fmt.Sprintf("Product Name %q", id.ProductName),
		fmt.Sprintf("Service Name %q", id.ServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Product Group", segmentsStr)
}

func (id ProductGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/products/%s/groups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.ProductName, id.GroupName)
}

// ProductGroupID parses a ProductGroup ID into an ProductGroupId struct
func ProductGroupID(input string) (*ProductGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ProductGroupId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServiceName, err = id.PopSegment("service"); err != nil {
		return nil, err
	}
	if resourceId.ProductName, err = id.PopSegment("products"); err != nil {
		return nil, err
	}
	if resourceId.GroupName, err = id.PopSegment("groups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
