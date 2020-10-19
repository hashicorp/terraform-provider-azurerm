package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LogAnalyticsClusterId struct {
	ResourceGroup string
	Name          string
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
