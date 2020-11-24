package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CertificateId struct {
	SubscriptionId   string
	ResourceGroup    string
	BatchAccountName string
	CertificateName  string
}

func NewCertificateID(subscriptionId, resourceGroup, batchAccountName, certificateName string) CertificateId {
	return CertificateId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		BatchAccountName: batchAccountName,
		CertificateName:  certificateName,
	}
}

func (id CertificateId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Batch/batchAccounts/%s/certificates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.BatchAccountName, id.CertificateName)
}

func CertificateID(input string) (*CertificateId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := CertificateId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.BatchAccountName, err = id.PopSegment("batchAccounts"); err != nil {
		return nil, err
	}
	if resourceId.CertificateName, err = id.PopSegment("certificates"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
