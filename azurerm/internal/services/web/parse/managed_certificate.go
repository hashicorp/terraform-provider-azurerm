package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ManagedCertificateId struct {
	SubscriptionId  string
	ResourceGroup   string
	CertificateName string
}

func NewManagedCertificateID(subscriptionId, resourceGroup, certificateName string) ManagedCertificateId {
	return ManagedCertificateId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		CertificateName: certificateName,
	}
}

func (id ManagedCertificateId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/certificates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.CertificateName)
}

func ManagedCertificateID(input string) (*ManagedCertificateId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ManagedCertificateId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.CertificateName, err = id.PopSegment("certificates"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
