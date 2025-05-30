// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package paloalto

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	certificates "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/certificateobjectlocalrulestack"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrulestacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type LocalRuleStackRule struct{}

var _ sdk.ResourceWithUpdate = LocalRuleStackRule{}

var protocolApplicationDefault = "application-default"

type LocalRuleModel struct {
	Name        string `tfschema:"name"`
	RuleStackID string `tfschema:"rulestack_id"`
	Priority    int64  `tfschema:"priority"`

	Action                  string                 `tfschema:"action"`
	Applications            []string               `tfschema:"applications"`
	AuditComment            string                 `tfschema:"audit_comment"`
	Category                []schema.Category      `tfschema:"category"`
	DecryptionRuleType      string                 `tfschema:"decryption_rule_type"`
	Description             string                 `tfschema:"description"`
	Destination             []schema.Destination   `tfschema:"destination"`
	LoggingEnabled          bool                   `tfschema:"logging_enabled"`
	InspectionCertificateID string                 `tfschema:"inspection_certificate_id"` // This is the name of a Certificate resource belonging to the SAME LocalRuleStack as this rule
	NegateDestination       bool                   `tfschema:"negate_destination"`
	NegateSource            bool                   `tfschema:"negate_source"`
	Protocol                string                 `tfschema:"protocol"`
	ProtocolPorts           []string               `tfschema:"protocol_ports"`
	RuleEnabled             bool                   `tfschema:"enabled"`
	Source                  []schema.Source        `tfschema:"source"`
	Tags                    map[string]interface{} `tfschema:"tags"`
}

func (r LocalRuleStackRule) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return localrules.ValidateLocalRuleID
}

func (r LocalRuleStackRule) ResourceType() string {
	return "azurerm_palo_alto_local_rulestack_rule"
}

func (r LocalRuleStackRule) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.LocalRuleStackRuleName, // TODO - Check this
		},

		"rulestack_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: localrules.ValidateLocalRulestackID,
		},

		"priority": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 10000),
		},

		"action": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(localrules.PossibleValuesForActionEnum(), false),
		},

		// Optional

		"applications": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"audit_comment": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"category": schema.CategorySchema(),

		"decryption_rule_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      localrules.DecryptionRuleTypeEnumNone,
			ValidateFunc: validation.StringInSlice(localrules.PossibleValuesForDecryptionRuleTypeEnum(), false),
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"destination": schema.DestinationSchema(),

		"logging_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"inspection_certificate_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: certificates.ValidateLocalRulestackCertificateID,
		},

		"negate_destination": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"negate_source": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"protocol": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.Any(
				validate.ProtocolWithPort,
				validation.StringInSlice([]string{protocolApplicationDefault}, false),
			),
			ExactlyOneOf: []string{"protocol", "protocol_ports"},
		},

		"protocol_ports": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.ProtocolWithPort,
			},
			ExactlyOneOf: []string{"protocol", "protocol_ports"},
		},

		"source": schema.SourceSchema(),

		"tags": commonschema.Tags(),
	}

	return schema
}

func (r LocalRuleStackRule) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LocalRuleStackRule) ModelObject() interface{} {
	return &LocalRuleModel{}
}

func (r LocalRuleStackRule) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.Client.LocalRules
			rulestackClient := metadata.Client.PaloAlto.Client.LocalRulestacks

			model := LocalRuleModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			rulestackId, err := localrulestacks.ParseLocalRulestackID(model.RuleStackID)
			if err != nil {
				return err
			}
			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			// API uses Priority not Name for ID, despite swagger defining `ruleName` as required, not Priority - https://github.com/Azure/azure-rest-api-specs/issues/24697
			id := localrules.NewLocalRuleID(metadata.Client.Account.SubscriptionId, rulestackId.ResourceGroupName, rulestackId.LocalRulestackName, strconv.FormatInt(model.Priority, 10))

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			destination, err := schema.ExpandDestination(model.Destination)
			if err != nil {
				return fmt.Errorf("expanding destination for %s, %+v", id, err)
			}

			source, err := schema.ExpandSource(model.Source)
			if err != nil {
				return fmt.Errorf("expanding source for %s: %+v", id, err)
			}

			props := localrules.RuleEntry{
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
				props.ActionType = pointer.To(localrules.ActionEnum(model.Action))
			}

			if len(model.Applications) != 0 {
				props.Applications = pointer.To(model.Applications)
			}

			if model.AuditComment != "" {
				props.AuditComment = pointer.To(model.AuditComment)
			}

			if model.DecryptionRuleType != "" {
				props.DecryptionRuleType = pointer.To(localrules.DecryptionRuleTypeEnum(model.DecryptionRuleType))
			}

			if model.Description != "" {
				props.Description = pointer.To(model.Description)
			}

			if model.InspectionCertificateID != "" {
				certID, err := certificates.ParseLocalRulestackCertificateID(model.InspectionCertificateID)
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

			if _, err = client.CreateOrUpdate(ctx, id, localrules.LocalRulesResource{Properties: props}); err != nil {
				return err
			}

			metadata.SetID(id)

			if err = rulestackClient.CommitThenPoll(ctx, *rulestackId); err != nil {
				return fmt.Errorf("committing Local Rulestack config for %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r LocalRuleStackRule) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.Client.LocalRules

			id, err := localrules.ParseLocalRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LocalRuleModel

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state.RuleStackID = localrulestacks.NewLocalRulestackID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName).ID()
			p, err := strconv.ParseInt(id.LocalRuleName, 10, 0)
			if err != nil {
				return fmt.Errorf("parsing Rule Priortiy for %s: %+v", *id, err)
			}
			state.Priority = p
			if model := existing.Model; model != nil {
				props := model.Properties
				state.Name = props.RuleName
				state.Action = string(pointer.From(props.ActionType))
				state.Applications = pointer.From(props.Applications)
				state.AuditComment = pointer.From(props.AuditComment)
				state.Category = schema.FlattenCategory(props.Category)
				state.DecryptionRuleType = string(pointer.From(props.DecryptionRuleType))
				state.Description = pointer.From(props.Description)
				state.Destination = schema.FlattenDestination(props.Destination, *id)
				state.LoggingEnabled = stateEnumAsBool(props.EnableLogging)
				if certName := pointer.From(props.InboundInspectionCertificate); certName != "" {
					state.InspectionCertificateID = certificates.NewLocalRulestackCertificateID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName, certName).ID()
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

func (r LocalRuleStackRule) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.Client.LocalRules

			id, err := localrules.ParseLocalRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			rulestackId := localrulestacks.NewLocalRulestackID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName)
			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r LocalRuleStackRule) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.Client.LocalRules
			rulestackClient := metadata.Client.PaloAlto.Client.LocalRulestacks

			model := LocalRuleModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id, err := localrules.ParseLocalRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			rulestackId := localrulestacks.NewLocalRulestackID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName)
			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retreiving %s: %+v", *id, err)
			}

			ruleEntry := *existing.Model

			if metadata.ResourceData.HasChange("action") {
				ruleEntry.Properties.ActionType = pointer.To(localrules.ActionEnum(model.Action))
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
				ruleEntry.Properties.DecryptionRuleType = pointer.To(localrules.DecryptionRuleTypeEnum(model.DecryptionRuleType))
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
					certID, err := certificates.ParseLocalRulestackCertificateID(model.InspectionCertificateID)
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

			if _, err = client.CreateOrUpdate(ctx, *id, ruleEntry); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			if err = rulestackClient.CommitThenPoll(ctx, rulestackId); err != nil {
				return fmt.Errorf("committing Local Rulestack config for %s: %+v", id, err)
			}

			return nil
		},
	}
}

func boolAsStateEnum(input bool) *localrules.StateEnum {
	var result localrules.StateEnum

	if input {
		result = localrules.StateEnumENABLED
	} else {
		result = localrules.StateEnumDISABLED
	}

	return pointer.To(result)
}

func stateEnumAsBool(input *localrules.StateEnum) bool {
	return pointer.From(input) == localrules.StateEnumENABLED
}

func boolAsBooleanEnumRule(input bool) *localrules.BooleanEnum {
	var result localrules.BooleanEnum

	if input {
		result = localrules.BooleanEnumTRUE
	} else {
		result = localrules.BooleanEnumFALSE
	}

	return pointer.To(result)
}

func boolEnumAsBoolRule(input *localrules.BooleanEnum) bool {
	return pointer.From(input) == localrules.BooleanEnumTRUE
}

func expandTagsForRule(input map[string]interface{}) *[]localrules.TagInfo {
	result := make([]localrules.TagInfo, 0)
	if len(input) == 0 {
		return pointer.To(result)
	}

	for k, v := range input {
		result = append(result, localrules.TagInfo{
			Key:   k,
			Value: v.(string),
		})
	}

	return pointer.To(result)
}

func flattenTagsFromRule(input *[]localrules.TagInfo) map[string]interface{} {
	if input == nil {
		return map[string]interface{}{}
	}

	result := make(map[string]interface{})
	for _, v := range *input {
		result[v.Key] = v.Value
	}

	return result
}
