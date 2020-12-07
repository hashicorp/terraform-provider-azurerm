package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SqlVirtualMachineId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewSqlVirtualMachineID(subscriptionId, resourceGroup, name string) SqlVirtualMachineId {
	return SqlVirtualMachineId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id SqlVirtualMachineId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Name %q", id.Name),
	}
	return strings.Join(segments, " / ")
}

func (id SqlVirtualMachineId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachines/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// SqlVirtualMachineID parses a SqlVirtualMachine ID into an SqlVirtualMachineId struct
func SqlVirtualMachineID(input string) (*SqlVirtualMachineId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SqlVirtualMachineId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("sqlVirtualMachines"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
