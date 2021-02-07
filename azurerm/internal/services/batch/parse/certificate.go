package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CertificateId struct {
	SubscriptionId   string
	ResourceGroup    string
	BatchAccountName string
	Name             string
}

func NewCertificateID(subscriptionId, resourceGroup, batchAccountName, name string) CertificateId {
	return CertificateId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		BatchAccountName: batchAccountName,
		Name:             name,
	}
}

func (id CertificateId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Batch Account Name %q", id.BatchAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Certificate", segmentsStr)
}

func (id CertificateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Batch/batchAccounts/%s/certificates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.BatchAccountName, id.Name)
}

// CertificateID parses a Certificate ID into an CertificateId struct
func CertificateID(input string) (*CertificateId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := CertificateId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.BatchAccountName, err = id.PopSegment("batchAccounts"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("certificates"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
