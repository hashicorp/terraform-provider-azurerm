// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack_rule

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/certificateobjectlocalrulestackresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulesresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulestackresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulesResources

			id, err := localrulesresources.ParseLocalRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LocalRuleModel

			existing, err := client.LocalRulesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state.RuleStackID = localrulestackresources.NewLocalRulestackID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName).ID()
			p, err := strconv.ParseInt(id.LocalRuleName, 10, 0)
			if err != nil {
				return fmt.Errorf("parsing Rule Priortiy for %s: %+v", *id, err)
			}
			state.Priority = p
			if model := existing.Model; model != nil {
				props := model.Properties
				state.Name = props.RuleName
				state.Action = pointer.FromEnum(props.ActionType)
				state.Applications = pointer.From(props.Applications)
				state.AuditComment = pointer.From(props.AuditComment)
				state.Category = schema.FlattenCategory(props.Category)
				state.DecryptionRuleType = pointer.FromEnum(props.DecryptionRuleType)
				state.Description = pointer.From(props.Description)
				state.Destination = schema.FlattenDestination(props.Destination, *id)
				state.LoggingEnabled = stateEnumAsBool(props.EnableLogging)
				if certName := pointer.From(props.InboundInspectionCertificate); certName != "" {
					state.InspectionCertificateID = certificateobjectlocalrulestackresources.NewLocalRulestackCertificateID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName, certName).ID()
				} else {
					state.InspectionCertificateID = certName
				}
				state.NegateDestination = boolEnumAsBoolRule(props.NegateDestination)
				state.NegateSource = boolEnumAsBoolRule(props.NegateSource)
				state.Protocol = pointer.From(props.Protocol)
				state.ProtocolPorts = pointer.From(props.ProtocolPortList)
				state.RuleEnabled = stateEnumAsBool(props.RuleState)
				state.Source = schema.FlattenSource(props.Source, *id)
				state.Tags = flattenTagsFromRule(props.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenTagsFromRule(input *[]localrulesresources.TagInfo) map[string]any {
	if input == nil {
		return map[string]any{}
	}

	result := make(map[string]any)
	for _, v := range *input {
		result[v.Key] = v.Value
	}

	return result
}

var protocolApplicationDefault = "application-default"

func boolAsStateEnum(input bool) *localrulesresources.StateEnum {
	var result localrulesresources.StateEnum

	if input {
		result = localrulesresources.StateEnumENABLED
	} else {
		result = localrulesresources.StateEnumDISABLED
	}

	return pointer.To(result)
}

func stateEnumAsBool(input *localrulesresources.StateEnum) bool {
	return pointer.From(input) == localrulesresources.StateEnumENABLED
}

func boolAsBooleanEnumRule(input bool) *localrulesresources.BooleanEnum {
	var result localrulesresources.BooleanEnum

	if input {
		result = localrulesresources.BooleanEnumTRUE
	} else {
		result = localrulesresources.BooleanEnumFALSE
	}

	return pointer.To(result)
}

func boolEnumAsBoolRule(input *localrulesresources.BooleanEnum) bool {
	return pointer.From(input) == localrulesresources.BooleanEnumTRUE
}
