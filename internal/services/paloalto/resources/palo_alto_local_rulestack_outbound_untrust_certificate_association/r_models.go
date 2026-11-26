// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack_outbound_untrust_certificate_association

type LocalRulestackOutboundUnTrustCertificateResourceModel struct {
	CertificateID string `tfschema:"certificate_id"`
}

func (l Resource) ModelObject() any {
	return &LocalRulestackOutboundUnTrustCertificateResourceModel{}
}
