package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ShareId struct {
	SubscriptionId string
	ResourceGroup  string
	AccountName    string
	Name           string
}

func NewShareID(subscriptionId, resourceGroup, accountName, name string) ShareId {
	return ShareId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		AccountName:    accountName,
		Name:           name,
	}
}

func (id ShareId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Account Name %q", id.AccountName),
		fmt.Sprintf("Name %q", id.Name),
	}
	return strings.Join(segments, " / ")
}

func (id ShareId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataShare/accounts/%s/shares/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.AccountName, id.Name)
}

// ShareID parses a Share ID into an ShareId struct
func ShareID(input string) (*ShareId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ShareId{
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
	if resourceId.Name, err = id.PopSegment("shares"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
