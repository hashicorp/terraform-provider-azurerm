package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LogAnalyticsStorageInsightConfigId struct {
	ResourceGroup string
	WorkspaceName string
	Name          string
}

func LogAnalyticsStorageInsightConfigID(input string) (*LogAnalyticsStorageInsightConfigId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing LogAnalyticsStorageInsightConfig ID %q: %+v", input, err)
	}

	logAnalyticsStorageInsightConfig := LogAnalyticsStorageInsightConfigId{
		ResourceGroup: id.ResourceGroup,
	}
	if logAnalyticsStorageInsightConfig.WorkspaceName, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if logAnalyticsStorageInsightConfig.Name, err = id.PopSegment("storageInsightConfigs"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &logAnalyticsStorageInsightConfig, nil
}
