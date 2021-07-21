package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type FlexibleServerConfigurationId struct {
	SubscriptionId     string
	ResourceGroup      string
	FlexibleServerName string
	ConfigurationName  string
}

func NewFlexibleServerConfigurationID(subscriptionId, resourceGroup, flexibleServerName, configurationName string) FlexibleServerConfigurationId {
	return FlexibleServerConfigurationId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		FlexibleServerName: flexibleServerName,
		ConfigurationName:  configurationName,
	}
}

func (id FlexibleServerConfigurationId) String() string {
	segments := []string{
		fmt.Sprintf("Configuration Name %q", id.ConfigurationName),
		fmt.Sprintf("Flexible Server Name %q", id.FlexibleServerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Flexible Server Configuration", segmentsStr)
}

func (id FlexibleServerConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforPostgreSQL/flexibleServers/%s/configurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FlexibleServerName, id.ConfigurationName)
}

// FlexibleServerConfigurationID parses a FlexibleServerConfiguration ID into an FlexibleServerConfigurationId struct
func FlexibleServerConfigurationID(input string) (*FlexibleServerConfigurationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FlexibleServerConfigurationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.FlexibleServerName, err = id.PopSegment("flexibleServers"); err != nil {
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
