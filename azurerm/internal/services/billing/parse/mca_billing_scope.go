package parse

import (
	"fmt"
	"net/url"
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
	idURL, err := url.ParseRequestURI(input)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse Azure ID: %s", err)
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")

	if len(components)%2 != 0 {
		return nil, fmt.Errorf("The number of path segments is not divisible by 2 in %q", path)
	}

	componentMap := make(map[string]string, len(components)/2)
	for current := 0; current < len(components); current += 2 {
		key := components[current]
		value := components[current+1]

		// Check key/value for empty strings.
		if key == "" || value == "" {
			return nil, fmt.Errorf("Key/Value cannot be empty strings. Key: '%s', Value: '%s'", key, value)
		}
		componentMap[key] = value
	}

	// Build up a TargetResourceID from the map
	id := &azure.ResourceID{}
	id.Path = componentMap

	if provider, ok := componentMap["providers"]; ok {
		id.Provider = provider
		delete(componentMap, "providers")
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
