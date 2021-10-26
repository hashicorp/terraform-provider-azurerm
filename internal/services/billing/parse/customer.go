package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type CustomerId struct {
	BillingAccountName string
	Name               string
}

func NewCustomerID(billingAccountName, name string) CustomerId {
	return CustomerId{
		BillingAccountName: billingAccountName,
		Name:               name,
	}
}

func (id CustomerId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Billing Account Name %q", id.BillingAccountName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Customer", segmentsStr)
}

func (id CustomerId) ID() string {
	fmtString := "/providers/Microsoft.Billing/billingAccounts/%s/customers/%s"
	return fmt.Sprintf(fmtString, id.BillingAccountName, id.Name)
}

// CustomerID parses a Customer ID into an CustomerId struct
func CustomerID(input string) (*CustomerId, error) {
	id, err := azure.ParseAzureResourceIDWithoutSubscription(input)
	if err != nil {
		return nil, err
	}

	resourceId := CustomerId{}

	if resourceId.BillingAccountName, err = id.PopSegment("billingAccounts"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("customers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
