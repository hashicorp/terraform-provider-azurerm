package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type IntegrationAccountAgreementId struct {
	SubscriptionId         string
	ResourceGroup          string
	IntegrationAccountName string
	AgreementName          string
}

func NewIntegrationAccountAgreementID(subscriptionId, resourceGroup, integrationAccountName, agreementName string) IntegrationAccountAgreementId {
	return IntegrationAccountAgreementId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		IntegrationAccountName: integrationAccountName,
		AgreementName:          agreementName,
	}
}

func (id IntegrationAccountAgreementId) String() string {
	segments := []string{
		fmt.Sprintf("Agreement Name %q", id.AgreementName),
		fmt.Sprintf("Integration Account Name %q", id.IntegrationAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Integration Account Agreement", segmentsStr)
}

func (id IntegrationAccountAgreementId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/integrationAccounts/%s/agreements/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IntegrationAccountName, id.AgreementName)
}

// IntegrationAccountAgreementID parses a IntegrationAccountAgreement ID into an IntegrationAccountAgreementId struct
func IntegrationAccountAgreementID(input string) (*IntegrationAccountAgreementId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := IntegrationAccountAgreementId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.IntegrationAccountName, err = id.PopSegment("integrationAccounts"); err != nil {
		return nil, err
	}
	if resourceId.AgreementName, err = id.PopSegment("agreements"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
