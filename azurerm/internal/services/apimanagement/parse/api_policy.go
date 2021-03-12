package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ApiPolicyId struct {
	SubscriptionId string
	ResourceGroup  string
	ServiceName    string
	ApiName        string
	PolicyName     string
}

func NewApiPolicyID(subscriptionId, resourceGroup, serviceName, apiName, policyName string) ApiPolicyId {
	return ApiPolicyId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServiceName:    serviceName,
		ApiName:        apiName,
		PolicyName:     policyName,
	}
}

func (id ApiPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Policy Name %q", id.PolicyName),
		fmt.Sprintf("Api Name %q", id.ApiName),
		fmt.Sprintf("Service Name %q", id.ServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Api Policy", segmentsStr)
}

func (id ApiPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/apis/%s/policies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.ApiName, id.PolicyName)
}

// ApiPolicyID parses a ApiPolicy ID into an ApiPolicyId struct
func ApiPolicyID(input string) (*ApiPolicyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ApiPolicyId{
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
	if resourceId.ApiName, err = id.PopSegment("apis"); err != nil {
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
