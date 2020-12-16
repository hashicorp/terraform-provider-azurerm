package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SecurityCenterContactId struct {
	SubscriptionId      string
	SecurityContactName string
}

func NewSecurityCenterContactID(subscriptionId, securityContactName string) SecurityCenterContactId {
	return SecurityCenterContactId{
		SubscriptionId:      subscriptionId,
		SecurityContactName: securityContactName,
	}
}

func (id SecurityCenterContactId) String() string {
	segments := []string{
		fmt.Sprintf("Security Contact Name %q", id.SecurityContactName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Security Center Contact", segmentsStr)
}

func (id SecurityCenterContactId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Security/securityContacts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.SecurityContactName)
}

// SecurityCenterContactID parses a SecurityCenterContact ID into an SecurityCenterContactId struct
func SecurityCenterContactID(input string) (*SecurityCenterContactId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SecurityCenterContactId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.SecurityContactName, err = id.PopSegment("securityContacts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
