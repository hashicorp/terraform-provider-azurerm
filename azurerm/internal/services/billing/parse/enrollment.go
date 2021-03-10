package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type EnrollmentId struct {
	EnrollmentAccountName string
}

func NewEnrollmentID(enrollmentAccountName string) EnrollmentId {
	return EnrollmentId{
		EnrollmentAccountName: enrollmentAccountName,
	}
}

func (id EnrollmentId) String() string {
	segments := []string{
		fmt.Sprintf("Enrollment Account Name %q", id.EnrollmentAccountName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Enrollment", segmentsStr)
}

func (id EnrollmentId) ID() string {
	fmtString := "/providers/Microsoft.Billing/enrollmentAccounts/%s"
	return fmt.Sprintf(fmtString, id.EnrollmentAccountName)
}

// EnrollmentID parses a Enrollment ID into an EnrollmentId struct
func EnrollmentID(input string) (*EnrollmentId, error) {
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

	resourceId := EnrollmentId{}

	if resourceId.EnrollmentAccountName, err = id.PopSegment("enrollmentAccounts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
