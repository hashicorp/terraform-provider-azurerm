// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// This is a special case ID for a meta resource that links a web certificate to a Web App or Function App

var _ resourceids.Id = CertificateBindingId{}

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

func (id CertificateBindingId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/hostNameBindings/%s|/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/certificates/%s"
	return fmt.Sprintf(fmtString, id.HostnameBindingId.SubscriptionId, id.HostnameBindingId.ResourceGroup, id.HostnameBindingId.SiteName, id.HostnameBindingId.Name, id.CertificateId.SubscriptionId, id.CertificateId.ResourceGroup, id.CertificateId.Name)
}

func (id CertificateBindingId) String() string {
	components := []string{
		fmt.Sprintf("Hostname Binding ID %q", id.HostnameBindingId.String()),
		fmt.Sprintf("Certificate ID %q", id.CertificateId.String()),
	}
	return fmt.Sprintf("Certificate Binding %s", strings.Join(components, " / "))
}

func CertificateBindingID(input string) (*CertificateBindingId, error) {
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

	return &CertificateBindingId{
		*hostnameBindingId,
		*certificateId,
	}, nil
}
