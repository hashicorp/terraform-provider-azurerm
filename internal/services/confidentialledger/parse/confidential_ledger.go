package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ConfidentialLedgerId struct {
	SubscriptionId string
	ResourceGroup  string
	LedgerName     string
}

func NewConfidentialLedgerID(subscriptionId, resourceGroup, ledgerName string) ConfidentialLedgerId {
	return ConfidentialLedgerId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		LedgerName:     ledgerName,
	}
}

func (id ConfidentialLedgerId) String() string {
	segments := []string{
		fmt.Sprintf("Ledger Name %q", id.LedgerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Confidential Ledger", segmentsStr)
}

func (id ConfidentialLedgerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ConfidentialLedger/Ledgers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LedgerName)
}

// ConfidentialLedgerID parses a ConfidentialLedger ID into an ConfidentialLedgerId struct
func ConfidentialLedgerID(input string) (*ConfidentialLedgerId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ConfidentialLedgerId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.LedgerName, err = id.PopSegment("Ledgers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
