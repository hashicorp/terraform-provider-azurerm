// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack_outbound_trust_certificate_association

type LocalRulestackOutboundTrustCertificateResourceModel struct {
	CertificateID string `tfschema:"certificate_id"`
}

func (l Resource) ModelObject() any {
	return &LocalRulestackOutboundTrustCertificateResourceModel{}
}
