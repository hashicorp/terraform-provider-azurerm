package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CustomerManagedKeyId struct {
	SubscriptionId          string
	ResourceGroup           string
	CustomerMangagedKeyName string
}

func NewCustomerManagedKeyID(subscriptionId, resourceGroup, customerMangagedKeyName string) CustomerManagedKeyId {
	return CustomerManagedKeyId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		CustomerMangagedKeyName: customerMangagedKeyName,
	}
}

func (id CustomerManagedKeyId) String() string {
	segments := []string{
		fmt.Sprintf("Customer Mangaged Key Name %q", id.CustomerMangagedKeyName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Customer Managed Key", segmentsStr)
}

func (id CustomerManagedKeyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Databricks/customerMangagedKey/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.CustomerMangagedKeyName)
}

// CustomerManagedKeyID parses a CustomerManagedKey ID into an CustomerManagedKeyId struct
func CustomerManagedKeyID(input string) (*CustomerManagedKeyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := CustomerManagedKeyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.CustomerMangagedKeyName, err = id.PopSegment("customerMangagedKey"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
