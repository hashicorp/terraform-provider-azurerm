// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack_rule

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
)

type LocalRuleModel struct {
	Name        string `tfschema:"name"`
	RuleStackID string `tfschema:"rulestack_id"`
	Priority    int64  `tfschema:"priority"`

	Action                  string               `tfschema:"action"`
	Applications            []string             `tfschema:"applications"`
	AuditComment            string               `tfschema:"audit_comment"`
	Category                []schema.Category    `tfschema:"category"`
	DecryptionRuleType      string               `tfschema:"decryption_rule_type"`
	Description             string               `tfschema:"description"`
	Destination             []schema.Destination `tfschema:"destination"`
	LoggingEnabled          bool                 `tfschema:"logging_enabled"`
	InspectionCertificateID string               `tfschema:"inspection_certificate_id"` // This is the name of a Certificate resource belonging to the SAME LocalRuleStack as this rule
	NegateDestination       bool                 `tfschema:"negate_destination"`
	NegateSource            bool                 `tfschema:"negate_source"`
	Protocol                string               `tfschema:"protocol"`
	ProtocolPorts           []string             `tfschema:"protocol_ports"`
	RuleEnabled             bool                 `tfschema:"enabled"`
	Source                  []schema.Source      `tfschema:"source"`
	Tags                    map[string]any       `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &LocalRuleModel{}
}
