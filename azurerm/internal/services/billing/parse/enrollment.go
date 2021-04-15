package parse

import (
	"fmt"
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
	id, err := azure.ParseAzureResourceIDWithoutSubscription(input)
	if err != nil {
		return nil, err
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
