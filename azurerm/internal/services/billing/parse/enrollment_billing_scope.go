package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type EnrollmentBillingScopeId struct {
	BillingAccountName    string
	EnrollmentAccountName string
}

func NewEnrollmentBillingScopeID(billingAccountName, enrollmentAccountName string) EnrollmentBillingScopeId {
	return EnrollmentBillingScopeId{
		BillingAccountName:    billingAccountName,
		EnrollmentAccountName: enrollmentAccountName,
	}
}

func (id EnrollmentBillingScopeId) String() string {
	segments := []string{
		fmt.Sprintf("Enrollment Account Name %q", id.EnrollmentAccountName),
		fmt.Sprintf("Billing Account Name %q", id.BillingAccountName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Enrollment Billing Scope", segmentsStr)
}

func (id EnrollmentBillingScopeId) ID() string {
	fmtString := "/providers/Microsoft.Billing/billingAccounts/%s/enrollmentAccounts/%s"
	return fmt.Sprintf(fmtString, id.BillingAccountName, id.EnrollmentAccountName)
}

// EnrollmentBillingScopeID parses a EnrollmentBillingScope ID into an EnrollmentBillingScopeId struct
func EnrollmentBillingScopeID(input string) (*EnrollmentBillingScopeId, error) {
	id, err := azure.ParseAzureResourceIDWithoutSubscription(input)
	if err != nil {
		return nil, err
	}

	resourceId := EnrollmentBillingScopeId{}

	if resourceId.BillingAccountName, err = id.PopSegment("billingAccounts"); err != nil {
		return nil, err
	}
	if resourceId.EnrollmentAccountName, err = id.PopSegment("enrollmentAccounts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
