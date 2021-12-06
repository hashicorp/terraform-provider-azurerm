package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AnalyticsAccountId struct {
	SubscriptionId string
	ResourceGroup  string
	AccountName    string
}

func NewAnalyticsAccountID(subscriptionId, resourceGroup, accountName string) AnalyticsAccountId {
	return AnalyticsAccountId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		AccountName:    accountName,
	}
}

func (id AnalyticsAccountId) String() string {
	segments := []string{
		fmt.Sprintf("Account Name %q", id.AccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Analytics Account", segmentsStr)
}

func (id AnalyticsAccountId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataLakeAnalytics/accounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.AccountName)
}

// AnalyticsAccountID parses a AnalyticsAccount ID into an AnalyticsAccountId struct
func AnalyticsAccountID(input string) (*AnalyticsAccountId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AnalyticsAccountId{
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

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
