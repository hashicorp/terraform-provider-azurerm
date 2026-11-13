// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack_rule

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/certificateobjectlocalrulestackresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulesresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulestackresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulesResources
			rulestackClient := metadata.Client.PaloAlto.LocalRulestackResources

			model := LocalRuleModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id, err := localrulesresources.ParseLocalRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			rulestackId := localrulestackresources.NewLocalRulestackID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName)
			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			existing, err := client.LocalRulesGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			ruleEntry := *existing.Model

			if metadata.ResourceData.HasChange("action") {
				ruleEntry.Properties.ActionType = pointer.ToEnum[localrulesresources.ActionEnum](model.Action)
			}

			if metadata.ResourceData.HasChange("applications") {
				ruleEntry.Properties.Applications = pointer.To(model.Applications)
			}

			if metadata.ResourceData.HasChange("audit_comment") {
				ruleEntry.Properties.AuditComment = pointer.To(model.AuditComment)
			}

			if metadata.ResourceData.HasChange("category") {
				ruleEntry.Properties.Category = schema.ExpandCategory(model.Category)
			}

			if metadata.ResourceData.HasChange("decryption_rule_type") {
				ruleEntry.Properties.DecryptionRuleType = pointer.ToEnum[localrulesresources.DecryptionRuleTypeEnum](model.DecryptionRuleType)
			}

			if metadata.ResourceData.HasChange("description") {
				ruleEntry.Properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("destination") {
				destination, err := schema.ExpandDestination(model.Destination)
				if err != nil {
					return fmt.Errorf("expanding destination for %s, %+v", id, err)
				}
				ruleEntry.Properties.Destination = destination
			}

			if metadata.ResourceData.HasChange("logging_enabled") {
				ruleEntry.Properties.EnableLogging = boolAsStateEnum(model.LoggingEnabled)
			}

			if metadata.ResourceData.HasChange("inspection_certificate_id") {
				if model.InspectionCertificateID != "" {
					certID, err := certificateobjectlocalrulestackresources.ParseLocalRulestackCertificateID(model.InspectionCertificateID)
					if err != nil {
						return err
					}
					ruleEntry.Properties.InboundInspectionCertificate = pointer.To(certID.CertificateName)
				} else {
					ruleEntry.Properties.InboundInspectionCertificate = pointer.To("")
				}
			}

			if metadata.ResourceData.HasChange("negate_destination") {
				ruleEntry.Properties.NegateDestination = boolAsBooleanEnumRule(model.NegateDestination)
			}

			if metadata.ResourceData.HasChange("negate_source") {
				ruleEntry.Properties.NegateSource = boolAsBooleanEnumRule(model.NegateSource)
			}

			if metadata.ResourceData.HasChange("protocol") {
				if model.Protocol != "" && !strings.EqualFold(model.Protocol, protocolApplicationDefault) && len(model.ProtocolPorts) == 0 {
					ruleEntry.Properties.Protocol = pointer.To(model.Protocol)
				} else {
					ruleEntry.Properties.Protocol = nil
				}
			}

			if metadata.ResourceData.HasChange("protocol_ports") {
				if len(model.ProtocolPorts) != 0 {
					ruleEntry.Properties.ProtocolPortList = pointer.To(model.ProtocolPorts)
				} else {
					ruleEntry.Properties.ProtocolPortList = nil
				}
			}

			if metadata.ResourceData.HasChange("enabled") {
				ruleEntry.Properties.RuleState = boolAsStateEnum(model.RuleEnabled)
			}

			if metadata.ResourceData.HasChange("source") {
				source, err := schema.ExpandSource(model.Source)
				if err != nil {
					return fmt.Errorf("expanding source for %s: %+v", id, err)
				}
				ruleEntry.Properties.Source = source
			}

			if metadata.ResourceData.HasChange("tags") {
				ruleEntry.Properties.Tags = expandTagsForRule(model.Tags)
			}

			if _, err = client.LocalRulesCreateOrUpdate(ctx, *id, ruleEntry); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			if err = rulestackClient.LocalRulestackscommitThenPoll(ctx, rulestackId); err != nil {
				return fmt.Errorf("committing Local Rulestack config for %s: %+v", id, err)
			}

			return nil
		},
	}
}
