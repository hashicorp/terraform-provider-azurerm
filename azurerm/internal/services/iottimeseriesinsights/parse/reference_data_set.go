package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ReferenceDataSetId struct {
	SubscriptionId  string
	ResourceGroup   string
	EnvironmentName string
	Name            string
}

func NewReferenceDataSetID(subscriptionId, resourceGroup, environmentName, name string) ReferenceDataSetId {
	return ReferenceDataSetId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		EnvironmentName: environmentName,
		Name:            name,
	}
}

func (id ReferenceDataSetId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Environment Name %q", id.EnvironmentName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Reference Data Set", segmentsStr)
}

func (id ReferenceDataSetId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.TimeSeriesInsights/environments/%s/referenceDataSets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.EnvironmentName, id.Name)
}

// ReferenceDataSetID parses a ReferenceDataSet ID into an ReferenceDataSetId struct
func ReferenceDataSetID(input string) (*ReferenceDataSetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ReferenceDataSetId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.EnvironmentName, err = id.PopSegment("environments"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("referenceDataSets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
