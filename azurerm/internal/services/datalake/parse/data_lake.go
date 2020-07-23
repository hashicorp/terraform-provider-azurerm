package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DataLakeStoreId struct {
	Subscription  string
	ResourceGroup string
	Name          string
}

func DataLakeStoreID(input string) (*DataLakeStoreId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Data Lake Store ID %q: %+v", input, err)
	}

	dataLakeStore := DataLakeStoreId{
		ResourceGroup: id.ResourceGroup,
		Subscription:  id.SubscriptionID,
	}
	if dataLakeStore.Name, err = id.PopSegment("accounts"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &dataLakeStore, nil
}
