package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type MariaDBConfigurationId struct {
	SubscriptionId    string
	ResourceGroup     string
	ServerName        string
	ConfigurationName string
}

func NewMariaDBConfigurationID(subscriptionId, resourceGroup, serverName, configurationName string) MariaDBConfigurationId {
	return MariaDBConfigurationId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		ServerName:        serverName,
		ConfigurationName: configurationName,
	}
}

func (id MariaDBConfigurationId) String() string {
	segments := []string{
		fmt.Sprintf("Configuration Name %q", id.ConfigurationName),
		fmt.Sprintf("Server Name %q", id.ServerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Maria D B Configuration", segmentsStr)
}

func (id MariaDBConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforMariaDB/servers/%s/configurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServerName, id.ConfigurationName)
}

// MariaDBConfigurationID parses a MariaDBConfiguration ID into an MariaDBConfigurationId struct
func MariaDBConfigurationID(input string) (*MariaDBConfigurationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := MariaDBConfigurationId{
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
	if resourceId.ConfigurationName, err = id.PopSegment("configurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
