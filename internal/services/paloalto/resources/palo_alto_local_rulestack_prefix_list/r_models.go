// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack_prefix_list

type LocalRuleStackPrefixListModel struct {
	Name         string   `tfschema:"name"`
	RuleStackID  string   `tfschema:"rulestack_id"`
	PrefixList   []string `tfschema:"prefix_list"`
	AuditComment string   `tfschema:"audit_comment"`
	Description  string   `tfschema:"description"`
}

func (r Resource) ModelObject() any {
	return &LocalRuleStackPrefixListModel{}
}
