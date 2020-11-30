package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SpringCloudCertificateId struct {
	SubscriptionId  string
	ResourceGroup   string
	SpringName      string
	CertificateName string
}

func NewSpringCloudCertificateID(subscriptionId, resourceGroup, springName, certificateName string) SpringCloudCertificateId {
	return SpringCloudCertificateId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		SpringName:      springName,
		CertificateName: certificateName,
	}
}

func (id SpringCloudCertificateId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/Spring/%s/certificates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.CertificateName)
}

// SpringCloudCertificateID parses a SpringCloudCertificate ID into an SpringCloudCertificateId struct
func SpringCloudCertificateID(input string) (*SpringCloudCertificateId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudCertificateId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SpringName, err = id.PopSegment("Spring"); err != nil {
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
