package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

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

func (id AzureActiveDirectoryAdministratorId) String() string {
	segments := []string{
		fmt.Sprintf("Administrator Name %q", id.AdministratorName),
		fmt.Sprintf("Server Name %q", id.ServerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Azure Active Directory Administrator", segmentsStr)
}

func (id AzureActiveDirectoryAdministratorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforMySQL/servers/%s/administrators/%s"
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

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
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
