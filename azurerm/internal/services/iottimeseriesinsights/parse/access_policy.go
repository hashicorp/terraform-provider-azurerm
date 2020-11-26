package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AccessPolicyId struct {
	ResourceGroup   string
	EnvironmentName string
	Name            string
}

func NewAccessPolicyID(resourceGroup, environmentName, name string) AccessPolicyId {
	return AccessPolicyId{
		ResourceGroup:   resourceGroup,
		EnvironmentName: environmentName,
		Name:            name,
	}
}

func (id AccessPolicyId) ID(subscriptionId string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.TimeSeriesInsights/environments/%s/accessPolicies/%s"
	return fmt.Sprintf(fmtString, subscriptionId, id.ResourceGroup, id.EnvironmentName, id.Name)
}

func AccessPolicyID(input string) (*AccessPolicyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Time Series Insights Access Policy ID %q: %+v", input, err)
	}

	service := AccessPolicyId{
		ResourceGroup: id.ResourceGroup,
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
