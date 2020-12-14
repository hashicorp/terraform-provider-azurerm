package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DataSetId struct {
	SubscriptionId string
	ResourceGroup  string
	AccountName    string
	ShareName      string
	Name           string
}

func NewDataSetID(subscriptionId, resourceGroup, accountName, shareName, name string) DataSetId {
	return DataSetId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		AccountName:    accountName,
		ShareName:      shareName,
		Name:           name,
	}
}

func (id DataSetId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Share Name %q", id.ShareName),
		fmt.Sprintf("Account Name %q", id.AccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Data Set", segmentsStr)
}

func (id DataSetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataShare/accounts/%s/shares/%s/dataSets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.AccountName, id.ShareName, id.Name)
}

// DataSetID parses a DataSet ID into an DataSetId struct
func DataSetID(input string) (*DataSetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DataSetId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.AccountName, err = id.PopSegment("accounts"); err != nil {
		return nil, err
	}
	if resourceId.ShareName, err = id.PopSegment("shares"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("dataSets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
