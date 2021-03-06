package parse

import (
	"fmt"
	"net/url"
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
