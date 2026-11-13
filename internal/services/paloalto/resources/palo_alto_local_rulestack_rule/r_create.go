// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack_rule

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/certificateobjectlocalrulestackresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulesresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulestackresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulesResources
			rulestackClient := metadata.Client.PaloAlto.LocalRulestackResources

			model := LocalRuleModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			rulestackId, err := localrulestackresources.ParseLocalRulestackID(model.RuleStackID)
			if err != nil {
				return err
			}
			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			// API uses Priority not Name for ID, despite swagger defining `ruleName` as required, not Priority - https://github.com/Azure/azure-rest-api-specs/issues/24697
			id := localrulesresources.NewLocalRuleID(metadata.Client.Account.SubscriptionId, rulestackId.ResourceGroupName, rulestackId.LocalRulestackName, strconv.FormatInt(model.Priority, 10))

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.LocalRulesGet(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			destination, err := schema.ExpandDestination(model.Destination)
			if err != nil {
				return fmt.Errorf("expanding destination for %s, %+v", id, err)
			}

			source, err := schema.ExpandSource(model.Source)
			if err != nil {
				return fmt.Errorf("expanding source for %s: %+v", id, err)
			}

			props := localrulesresources.RuleEntry{
				Category:          schema.ExpandCategory(model.Category),
				Destination:       destination,
				EnableLogging:     boolAsStateEnum(model.LoggingEnabled),
				NegateDestination: boolAsBooleanEnumRule(model.NegateDestination),
				NegateSource:      boolAsBooleanEnumRule(model.NegateSource),
				RuleName:          model.Name,
				RuleState:         boolAsStateEnum(model.RuleEnabled),
				Source:            source,
				Tags:              expandTagsForRule(model.Tags),
			}

			if model.Action != "" {
				props.ActionType = pointer.ToEnum[localrulesresources.ActionEnum](model.Action)
			}

			if len(model.Applications) != 0 {
				props.Applications = pointer.To(model.Applications)
			}

			if model.AuditComment != "" {
				props.AuditComment = pointer.To(model.AuditComment)
			}

			if model.DecryptionRuleType != "" {
				props.DecryptionRuleType = pointer.ToEnum[localrulesresources.DecryptionRuleTypeEnum](model.DecryptionRuleType)
			}

			if model.Description != "" {
				props.Description = pointer.To(model.Description)
			}

			if model.InspectionCertificateID != "" {
				certID, err := certificateobjectlocalrulestackresources.ParseLocalRulestackCertificateID(model.InspectionCertificateID)
				if err != nil {
					return err
				}
				props.InboundInspectionCertificate = pointer.To(certID.CertificateName)
			}

			if model.Priority != 0 {
				props.Priority = pointer.To(model.Priority)
			}

			if len(model.ProtocolPorts) != 0 {
				props.ProtocolPortList = pointer.To(model.ProtocolPorts)
			}

			if model.Protocol != "" && !strings.EqualFold(model.Protocol, protocolApplicationDefault) && len(model.ProtocolPorts) == 0 {
				props.Protocol = pointer.To(model.Protocol)
			}

			if _, err := client.LocalRulesCreateOrUpdate(ctx, id, localrulesresources.LocalRulesResource{Properties: props}); err != nil {
				return err
			}

			metadata.SetID(id)

			if err := rulestackClient.LocalRulestackscommitThenPoll(ctx, *rulestackId); err != nil {
				return fmt.Errorf("committing Local Rulestack config for %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandTagsForRule(input map[string]any) *[]localrulesresources.TagInfo {
	result := make([]localrulesresources.TagInfo, 0)
	if len(input) == 0 {
		return pointer.To(result)
	}

	for k, v := range input {
		result = append(result, localrulesresources.TagInfo{
			Key:   k,
			Value: v.(string),
		})
	}

	return pointer.To(result)
}
