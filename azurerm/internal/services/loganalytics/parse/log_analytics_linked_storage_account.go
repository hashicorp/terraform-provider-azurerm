package parse

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2020-03-01-preview/operationalinsights"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

// TODO: @tombuildsstuff pass back through and generate these

type LogAnalyticsLinkedStorageAccountId struct {
	ResourceGroup string
	WorkspaceName string
	WorkspaceID   string
	Name          string
}

func LogAnalyticsLinkedStorageAccountID(input string) (*LogAnalyticsLinkedStorageAccountId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Log Analytics Linked Storage Account ID %q: %+v", input, err)
	}

	logAnalyticsLinkedStorageAccount := LogAnalyticsLinkedStorageAccountId{
		ResourceGroup: id.ResourceGroup,
	}
	if logAnalyticsLinkedStorageAccount.WorkspaceName, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if logAnalyticsLinkedStorageAccount.WorkspaceID = fmt.Sprintf("/subscriptions/%s/resourcegroups/%s/providers/%s/workspaces/%s", id.SubscriptionID, id.ResourceGroup, id.Provider, logAnalyticsLinkedStorageAccount.WorkspaceName); err != nil {
		return nil, fmt.Errorf("formatting Log Analytics Data Export Rule workspace ID %q", input)
	}
	var name string
	if name, err = id.PopSegment("linkedStorageAccounts"); err == nil {
		logAnalyticsLinkedStorageAccount.Name = string(operationalinsights.DataSourceType(name))
	} else {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &logAnalyticsLinkedStorageAccount, nil
}
