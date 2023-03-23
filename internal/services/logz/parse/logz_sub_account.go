package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type LogzSubAccountId struct {
	SubscriptionId string
	ResourceGroup  string
	MonitorName    string
	AccountName    string
}

func NewLogzSubAccountID(subscriptionId, resourceGroup, monitorName, accountName string) LogzSubAccountId {
	return LogzSubAccountId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		MonitorName:    monitorName,
		AccountName:    accountName,
	}
}

func (id LogzSubAccountId) String() string {
	segments := []string{
		fmt.Sprintf("Account Name %q", id.AccountName),
		fmt.Sprintf("Monitor Name %q", id.MonitorName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Logz Sub Account", segmentsStr)
}

func (id LogzSubAccountId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logz/monitors/%s/accounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.MonitorName, id.AccountName)
}

// LogzSubAccountID parses a LogzSubAccount ID into an LogzSubAccountId struct
func LogzSubAccountID(input string) (*LogzSubAccountId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LogzSubAccountId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.MonitorName, err = id.PopSegment("monitors"); err != nil {
		return nil, err
	}
	if resourceId.AccountName, err = id.PopSegment("accounts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
