package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LinkedServiceId struct {
	SubscriptionId string
	ResourceGroup  string
	FactoryName    string
	Name           string
}

func NewLinkedServiceID(subscriptionId, resourceGroup, factoryName, name string) LinkedServiceId {
	return LinkedServiceId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		FactoryName:    factoryName,
		Name:           name,
	}
}

func (id LinkedServiceId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Factory Name %q", id.FactoryName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Linked Service", segmentsStr)
}

func (id LinkedServiceId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataFactory/factories/%s/linkedservices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FactoryName, id.Name)
}

// LinkedServiceID parses a LinkedService ID into an LinkedServiceId struct
func LinkedServiceID(input string) (*LinkedServiceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LinkedServiceId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.FactoryName, err = id.PopSegment("factories"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("linkedservices"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
