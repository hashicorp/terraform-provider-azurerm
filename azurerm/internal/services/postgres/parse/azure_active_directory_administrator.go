package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AzureActiveDirectoryAdministratorId struct {
	SubscriptionId    string
	ResourceGroup     string
	ServerName        string
	AdministratorName string
}

func NewAzureActiveDirectoryAdministratorID(subscriptionId, resourceGroup, serverName, administratorName string) AzureActiveDirectoryAdministratorId {
	return AzureActiveDirectoryAdministratorId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		ServerName:        serverName,
		AdministratorName: administratorName,
	}
}

func (id AzureActiveDirectoryAdministratorId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforPostgreSQL/servers/%s/administrators/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServerName, id.AdministratorName)
}

// AzureActiveDirectoryAdministratorID parses a AzureActiveDirectoryAdministrator ID into an AzureActiveDirectoryAdministratorId struct
func AzureActiveDirectoryAdministratorID(input string) (*AzureActiveDirectoryAdministratorId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AzureActiveDirectoryAdministratorId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.ServerName, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}
	if resourceId.AdministratorName, err = id.PopSegment("administrators"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
