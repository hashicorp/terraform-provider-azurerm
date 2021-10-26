package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type BillingProfileId struct {
	BillingAccountName string
	Name               string
}

func NewBillingProfileID(billingAccountName, name string) BillingProfileId {
	return BillingProfileId{
		BillingAccountName: billingAccountName,
		Name:               name,
	}
}

func (id BillingProfileId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Billing Account Name %q", id.BillingAccountName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Billing Profile", segmentsStr)
}

func (id BillingProfileId) ID() string {
	fmtString := "/providers/Microsoft.Billing/billingAccounts/%s/billingProfiles/%s"
	return fmt.Sprintf(fmtString, id.BillingAccountName, id.Name)
}

// BillingProfileID parses a BillingProfile ID into an BillingProfileId struct
func BillingProfileID(input string) (*BillingProfileId, error) {
	id, err := azure.ParseAzureResourceIDWithoutSubscription(input)
	if err != nil {
		return nil, err
	}

	resourceId := BillingProfileId{}

	if resourceId.BillingAccountName, err = id.PopSegment("billingAccounts"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("billingProfiles"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
