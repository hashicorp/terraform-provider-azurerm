package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ProductPolicyId struct {
	SubscriptionId string
	ResourceGroup  string
	ServiceName    string
	ProductName    string
	PolicyName     string
}

func NewProductPolicyID(subscriptionId, resourceGroup, serviceName, productName, policyName string) ProductPolicyId {
	return ProductPolicyId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServiceName:    serviceName,
		ProductName:    productName,
		PolicyName:     policyName,
	}
}

func (id ProductPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Policy Name %q", id.PolicyName),
		fmt.Sprintf("Product Name %q", id.ProductName),
		fmt.Sprintf("Service Name %q", id.ServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Product Policy", segmentsStr)
}

func (id ProductPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/products/%s/policies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.ProductName, id.PolicyName)
}

// ProductPolicyID parses a ProductPolicy ID into an ProductPolicyId struct
func ProductPolicyID(input string) (*ProductPolicyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ProductPolicyId{
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
	if resourceId.PolicyName, err = id.PopSegment("policies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
