package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ProductApiId struct {
	SubscriptionId string
	ResourceGroup  string
	ServiceName    string
	ProductName    string
	ApiName        string
}

func NewProductApiID(subscriptionId, resourceGroup, serviceName, productName, apiName string) ProductApiId {
	return ProductApiId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServiceName:    serviceName,
		ProductName:    productName,
		ApiName:        apiName,
	}
}

func (id ProductApiId) String() string {
	segments := []string{
		fmt.Sprintf("Api Name %q", id.ApiName),
		fmt.Sprintf("Product Name %q", id.ProductName),
		fmt.Sprintf("Service Name %q", id.ServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Product Api", segmentsStr)
}

func (id ProductApiId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/products/%s/apis/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.ProductName, id.ApiName)
}

// ProductApiID parses a ProductApi ID into an ProductApiId struct
func ProductApiID(input string) (*ProductApiId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ProductApiId{
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
	if resourceId.ApiName, err = id.PopSegment("apis"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
