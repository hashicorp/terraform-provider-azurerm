package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type TimeSeriesInsightsReferenceDataSetId struct {
	SubscriptionId  string
	ResourceGroup   string
	EnvironmentName string
	Name            string
}

func NewTimeSeriesInsightsReferenceDataSetID(subscriptionId, resourceGroup, environmentName, name string) TimeSeriesInsightsReferenceDataSetId {
	return TimeSeriesInsightsReferenceDataSetId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		EnvironmentName: environmentName,
		Name:            name,
	}
}

func (id TimeSeriesInsightsReferenceDataSetId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.TimeSeriesInsights/environments/%s/referenceDataSets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.EnvironmentName, id.Name)
}

func TimeSeriesInsightsReferenceDataSetID(input string) (*TimeSeriesInsightsReferenceDataSetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Time Series Insights Reference Dataset ID %q: %+v", input, err)
	}

	service := TimeSeriesInsightsReferenceDataSetId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if service.EnvironmentName, err = id.PopSegment("environments"); err != nil {
		return nil, err
	}

	if service.Name, err = id.PopSegment("referenceDataSets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
