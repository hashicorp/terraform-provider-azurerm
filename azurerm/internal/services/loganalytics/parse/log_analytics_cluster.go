package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LogAnalyticsClusterId struct {
	ResourceGroup string
	Name          string
}

func NewLogAnalyticsClusterId(name, resourceGroup string) LogAnalyticsClusterId {
	return LogAnalyticsClusterId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func (id LogAnalyticsClusterId) ID(subscriptionId string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/clusters/%s"
	return fmt.Sprintf(fmtString, subscriptionId, id.ResourceGroup, id.Name)
}

func LogAnalyticsClusterID(input string) (*LogAnalyticsClusterId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing operationalinsightsCluster ID %q: %+v", input, err)
	}

	logAnalyticsCluster := LogAnalyticsClusterId{
		ResourceGroup: id.ResourceGroup,
	}

	if logAnalyticsCluster.Name, err = id.PopSegment("clusters"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &logAnalyticsCluster, nil
}
