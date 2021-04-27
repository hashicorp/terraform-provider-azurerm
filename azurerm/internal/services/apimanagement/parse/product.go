package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ProductId struct {
	SubscriptionId string
	ResourceGroup  string
	ServiceName    string
	Name           string
}

func NewProductID(subscriptionId, resourceGroup, serviceName, name string) ProductId {
	return ProductId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServiceName:    serviceName,
		Name:           name,
	}
}

func (id ProductId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Service Name %q", id.ServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Product", segmentsStr)
}

func (id ProductId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/products/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.Name)
}

// ProductID parses a Product ID into an ProductId struct
func ProductID(input string) (*ProductId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ProductId{
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
	if resourceId.Name, err = id.PopSegment("products"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
