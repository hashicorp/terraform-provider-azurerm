package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type BillingAccountId struct {
	Name string
}

func NewBillingAccountID(name string) BillingAccountId {
	return BillingAccountId{
		Name: name,
	}
}

func (id BillingAccountId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Billing Account", segmentsStr)
}

func (id BillingAccountId) ID() string {
	fmtString := "/providers/Microsoft.Billing/billingAccounts/%s"
	return fmt.Sprintf(fmtString, id.Name)
}

// BillingAccountID parses a BillingAccount ID into an BillingAccountId struct
func BillingAccountID(input string) (*BillingAccountId, error) {
	id, err := azure.ParseAzureResourceIDWithoutSubscription(input)
	if err != nil {
		return nil, err
	}

	resourceId := BillingAccountId{}

	if resourceId.Name, err = id.PopSegment("billingAccounts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
