package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type DpsCertificateId struct {
	SubscriptionId          string
	ResourceGroup           string
	ProvisioningServiceName string
	CertificateName         string
}

func NewDpsCertificateID(subscriptionId, resourceGroup, provisioningServiceName, certificateName string) DpsCertificateId {
	return DpsCertificateId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		ProvisioningServiceName: provisioningServiceName,
		CertificateName:         certificateName,
	}
}

func (id DpsCertificateId) String() string {
	segments := []string{
		fmt.Sprintf("Certificate Name %q", id.CertificateName),
		fmt.Sprintf("Provisioning Service Name %q", id.ProvisioningServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Dps Certificate", segmentsStr)
}

func (id DpsCertificateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Devices/provisioningServices/%s/certificates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProvisioningServiceName, id.CertificateName)
}

// DpsCertificateID parses a DpsCertificate ID into an DpsCertificateId struct
func DpsCertificateID(input string) (*DpsCertificateId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DpsCertificateId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ProvisioningServiceName, err = id.PopSegment("provisioningServices"); err != nil {
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
