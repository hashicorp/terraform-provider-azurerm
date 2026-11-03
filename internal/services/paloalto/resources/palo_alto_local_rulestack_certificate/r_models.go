// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack_certificate

type LocalRuleStackCertificateModel struct {
	Name                string `tfschema:"name"`
	RuleStackID         string `tfschema:"rulestack_id"`
	AuditComment        string `tfschema:"audit_comment"`
	CertificateSignerID string `tfschema:"key_vault_certificate_id"`
	Description         string `tfschema:"description"`
	SelfSigned          bool   `tfschema:"self_signed"`
}

func (r Resource) ModelObject() any {
	return &LocalRuleStackCertificateModel{}
}
