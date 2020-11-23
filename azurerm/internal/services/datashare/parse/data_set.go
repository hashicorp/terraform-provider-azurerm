package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DataShareDataSetId struct {
	SubscriptionId string
	ResourceGroup  string
	AccountName    string
	ShareName      string
	Name           string
}

func NewDataShareDataSetId(subscriptionId, resourceGroup, accountName, shareName, name string) DataShareDataSetId {
	return DataShareDataSetId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		AccountName:    accountName,
		ShareName:      shareName,
		Name:           name,
	}
}

func (id DataShareDataSetId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataShare/accounts/%s/shares/%s/dataSets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.AccountName, id.ShareName, id.Name)
}

func DataShareDataSetID(input string) (*DataShareDataSetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse DataShareDataSet ID %q: %+v", input, err)
	}

	dataShareDataSet := DataShareDataSetId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}
	if dataShareDataSet.AccountName, err = id.PopSegment("accounts"); err != nil {
		return nil, err
	}
	if dataShareDataSet.ShareName, err = id.PopSegment("shares"); err != nil {
		return nil, err
	}
	if dataShareDataSet.Name, err = id.PopSegment("dataSets"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &dataShareDataSet, nil
}
