package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type MaintenanceConfigurationId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewMaintenanceConfigurationID(subscriptionId, resourceGroup, name string) MaintenanceConfigurationId {
	return MaintenanceConfigurationId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id MaintenanceConfigurationId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Maintenance Configuration", segmentsStr)
}

func (id MaintenanceConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Maintenance/maintenanceConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// MaintenanceConfigurationID parses a MaintenanceConfiguration ID into an MaintenanceConfigurationId struct
func MaintenanceConfigurationID(input string) (*MaintenanceConfigurationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := MaintenanceConfigurationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("maintenanceConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// MaintenanceConfigurationIDInsensitively parses an MaintenanceConfiguration ID into an MaintenanceConfigurationId struct, insensitively
// This should only be used to parse an ID for rewriting, the MaintenanceConfigurationID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func MaintenanceConfigurationIDInsensitively(input string) (*MaintenanceConfigurationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := MaintenanceConfigurationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'maintenanceConfigurations' segment
	maintenanceConfigurationsKey := "maintenanceConfigurations"
	for key := range id.Path {
		if strings.EqualFold(key, maintenanceConfigurationsKey) {
			maintenanceConfigurationsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(maintenanceConfigurationsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
