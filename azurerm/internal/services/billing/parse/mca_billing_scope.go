package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type MicrosoftCustomerAccountBillingScopeId struct {
	BillingAccountName string
	BillingProfileName string
	InvoiceSectionName string
}

func NewMCABillingScopeID(billingAccountName, billingProfileName, invoiceSectionName string) MicrosoftCustomerAccountBillingScopeId {
	return MicrosoftCustomerAccountBillingScopeId{
		BillingAccountName: billingAccountName,
		BillingProfileName: billingProfileName,
		InvoiceSectionName: invoiceSectionName,
	}
}

func (id MicrosoftCustomerAccountBillingScopeId) String() string {
	segments := []string{
		fmt.Sprintf("Invoice Section Name %q", id.InvoiceSectionName),
		fmt.Sprintf("Billing Profile Name %q", id.BillingProfileName),
		fmt.Sprintf("Billing Account Name %q", id.BillingAccountName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "M C A Billing Scope", segmentsStr)
}

func (id MicrosoftCustomerAccountBillingScopeId) ID() string {
	fmtString := "/providers/Microsoft.Billing/billingAccounts/%s/billingProfiles/%s/invoiceSections/%s"
	return fmt.Sprintf(fmtString, id.BillingAccountName, id.BillingProfileName, id.InvoiceSectionName)
}

// MicrosoftCustomerAccountBillingScopeID parses a MCABillingScope ID into an MicrosoftCustomerAccountBillingScopeId struct
func MicrosoftCustomerAccountBillingScopeID(input string) (*MicrosoftCustomerAccountBillingScopeId, error) {
	id, err := azure.ParseAzureResourceIDWithoutSubscription(input)
	if err != nil {
		return nil, err
	}

	resourceId := MicrosoftCustomerAccountBillingScopeId{}

	if resourceId.BillingAccountName, err = id.PopSegment("billingAccounts"); err != nil {
		return nil, err
	}
	if resourceId.BillingProfileName, err = id.PopSegment("billingProfiles"); err != nil {
		return nil, err
	}
	if resourceId.InvoiceSectionName, err = id.PopSegment("invoiceSections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
