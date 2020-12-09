package parse

import (
	"fmt"
	"strings"
)

// This is a special case ID for a meta resource that links a web certificate to a Web App or Function App

type CertificateBindingId struct {
	HostnameBindingId
	CertificateId
}

func NewCertificateBindingId(hostnameBindingId HostnameBindingId, certificateId CertificateId) *CertificateBindingId {
	return &CertificateBindingId{
		HostnameBindingId: hostnameBindingId,
		CertificateId:     certificateId,
	}
}

func (id CertificateBindingId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/hostNameBindings/%s|/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/certificates/%s"
	return fmt.Sprintf(fmtString, id.HostnameBindingId.SubscriptionId, id.HostnameBindingId.ResourceGroup, id.HostnameBindingId.SiteName, id.HostnameBindingId.Name, id.CertificateId.SubscriptionId, id.CertificateId.ResourceGroup, id.CertificateId.Name)
}

func CertificateBindingID(input string) (*CertificateBindingId, error) {
	certificateBindingId := CertificateBindingId{}
	idParts := strings.Split(input, "|")
	if len(idParts) != 2 {
		return nil, fmt.Errorf("could not parse Certificate Binding ID, expected two resource IDs joined by `|`")
	}

	hostnameBindingId, err := HostnameBindingID(idParts[0])
	if err != nil {
		return nil, fmt.Errorf("could not parse Hostname Binding portion of Certificate Binding ID: %+v", err)
	}
	certificateId, err := CertificateID(idParts[1])
	if err != nil {
		return nil, fmt.Errorf("could not parse Certificate ID portion of Certificate Binding ID: %+v", err)
	}

	certificateBindingId.HostnameBindingId = *hostnameBindingId
	certificateBindingId.CertificateId = *certificateId

	return &certificateBindingId, nil
}
