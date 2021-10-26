package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type DepartmentId struct {
	BillingAccountName string
	Name               string
}

func NewDepartmentID(billingAccountName, name string) DepartmentId {
	return DepartmentId{
		BillingAccountName: billingAccountName,
		Name:               name,
	}
}

func (id DepartmentId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Billing Account Name %q", id.BillingAccountName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Department", segmentsStr)
}

func (id DepartmentId) ID() string {
	fmtString := "/providers/Microsoft.Billing/billingAccounts/%s/departments/%s"
	return fmt.Sprintf(fmtString, id.BillingAccountName, id.Name)
}

// DepartmentID parses a Department ID into an DepartmentId struct
func DepartmentID(input string) (*DepartmentId, error) {
	id, err := azure.ParseAzureResourceIDWithoutSubscription(input)
	if err != nil {
		return nil, err
	}

	resourceId := DepartmentId{}

	if resourceId.BillingAccountName, err = id.PopSegment("billingAccounts"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("departments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
