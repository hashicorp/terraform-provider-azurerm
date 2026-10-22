// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack_fqdn_list

type LocalRulestackFQDNListModel struct {
	Name         string   `tfschema:"name"`
	RuleStackID  string   `tfschema:"rulestack_id"`
	FQDNList     []string `tfschema:"fully_qualified_domain_names"`
	AuditComment string   `tfschema:"audit_comment"`
	Description  string   `tfschema:"description"`
}

func (r Resource) ModelObject() any {
	return &LocalRulestackFQDNListModel{}
}
