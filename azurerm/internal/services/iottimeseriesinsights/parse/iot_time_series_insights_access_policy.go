package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type TimeSeriesInsightsAccessPolicyId struct {
	SubscriptionId  string
	ResourceGroup   string
	EnvironmentName string
	Name            string
}

func NewTimeSeriesInsightsAccessPolicyID(subscriptionId, resourceGroup, environmentName, name string) TimeSeriesInsightsAccessPolicyId {
	return TimeSeriesInsightsAccessPolicyId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		EnvironmentName: environmentName,
		Name:            name,
	}
}

func (id TimeSeriesInsightsAccessPolicyId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.TimeSeriesInsights/environments/%s/accessPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.EnvironmentName, id.Name)
}

func TimeSeriesInsightsAccessPolicyID(input string) (*TimeSeriesInsightsAccessPolicyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Time Series Insights Access Policy ID %q: %+v", input, err)
	}

	service := TimeSeriesInsightsAccessPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if service.EnvironmentName, err = id.PopSegment("environments"); err != nil {
		return nil, err
	}

	if service.Name, err = id.PopSegment("accessPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
