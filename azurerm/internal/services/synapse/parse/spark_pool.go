package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SparkPoolId struct {
	SubscriptionId  string
	ResourceGroup   string
	WorkspaceName   string
	BigDataPoolName string
}

func NewSparkPoolID(subscriptionId, resourceGroup, workspaceName, bigDataPoolName string) SparkPoolId {
	return SparkPoolId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		WorkspaceName:   workspaceName,
		BigDataPoolName: bigDataPoolName,
	}
}

func (id SparkPoolId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Synapse/workspaces/%s/bigDataPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.BigDataPoolName)
}

// SparkPoolID parses a SparkPool ID into an SparkPoolId struct
func SparkPoolID(input string) (*SparkPoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SparkPoolId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.WorkspaceName, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if resourceId.BigDataPoolName, err = id.PopSegment("bigDataPools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
