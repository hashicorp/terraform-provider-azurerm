// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package paloalto

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/certificateobjectlocalrulestack"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrulestacks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	keyvaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LocalRuleStackCertificate struct{}

var _ sdk.ResourceWithUpdate = LocalRuleStackCertificate{}

type LocalRuleStackCertificateModel struct {
	Name                string `tfschema:"name"`
	RuleStackID         string `tfschema:"rulestack_id"`
	AuditComment        string `tfschema:"audit_comment"`
	CertificateSignerID string `tfschema:"key_vault_certificate_id"`
	Description         string `tfschema:"description"`
	SelfSigned          bool   `tfschema:"self_signed"`
}

func (r LocalRuleStackCertificate) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return certificateobjectlocalrulestack.ValidateLocalRulestackCertificateID
}

func (r LocalRuleStackCertificate) ResourceType() string {
	return "azurerm_palo_alto_local_rulestack_certificate"
}

func (r LocalRuleStackCertificate) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.LocalRuleStackCertificateName,
		},

		"rulestack_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: certificateobjectlocalrulestack.ValidateLocalRulestackID,
		},

		"audit_comment": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"key_vault_certificate_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: keyvaultValidate.VersionlessNestedItemId,
			ExactlyOneOf: []string{"self_signed", "key_vault_certificate_id"},
		},

		"self_signed": {
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			ForceNew:     true,
			Default:      false,
			ExactlyOneOf: []string{"key_vault_certificate_id", "self_signed"},
		},
	}
}

func (r LocalRuleStackCertificate) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LocalRuleStackCertificate) ModelObject() interface{} {
	return &LocalRuleStackCertificateModel{}
}

func (r LocalRuleStackCertificate) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.Client.CertificateObjectLocalRulestack
			rulestackClient := metadata.Client.PaloAlto.Client.LocalRulestacks

			model := LocalRuleStackCertificateModel{}
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			rulestackId, err := localrulestacks.ParseLocalRulestackID(model.RuleStackID)
			if err != nil {
				return err
			}

			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			id := certificateobjectlocalrulestack.NewLocalRulestackCertificateID(rulestackId.SubscriptionId, rulestackId.ResourceGroupName, rulestackId.LocalRulestackName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			props := certificateobjectlocalrulestack.CertificateObject{
				CertificateSelfSigned: boolAsBooleanEnumCert(model.SelfSigned),
			}

			if model.AuditComment != "" {
				props.AuditComment = pointer.To(model.AuditComment)
			}

			if model.CertificateSignerID != "" {
				props.CertificateSignerResourceId = pointer.To(model.CertificateSignerID)
			}

			if model.Description != "" {
				props.Description = pointer.To(model.Description)
			}

			cert := certificateobjectlocalrulestack.CertificateObjectLocalRulestackResource{
				Properties: props,
			}

			if _, err = client.CreateOrUpdate(ctx, id, cert); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			if err = rulestackClient.CommitThenPoll(ctx, *rulestackId); err != nil {
				return fmt.Errorf("committing Local RuleStack config for %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r LocalRuleStackCertificate) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.Client.CertificateObjectLocalRulestack

			id, err := certificateobjectlocalrulestack.ParseLocalRulestackCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LocalRuleStackCertificateModel

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state.Name = id.CertificateName
			state.RuleStackID = certificateobjectlocalrulestack.NewLocalRulestackID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName).ID()

			if model := existing.Model; model != nil {
				props := model.Properties

				state.AuditComment = pointer.From(props.AuditComment)
				state.CertificateSignerID = pointer.From(props.CertificateSignerResourceId)
				state.Description = pointer.From(props.Description)
				state.SelfSigned = boolEnumAsBoolCert(props.CertificateSelfSigned)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LocalRuleStackCertificate) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.Client.CertificateObjectLocalRulestack

			id, err := certificateobjectlocalrulestack.ParseLocalRulestackCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			rulestackId := localrulestacks.NewLocalRulestackID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName)
			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r LocalRuleStackCertificate) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.Client.CertificateObjectLocalRulestack
			rulestackClient := metadata.Client.PaloAlto.Client.LocalRulestacks
			model := LocalRuleStackCertificateModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id, err := certificateobjectlocalrulestack.ParseLocalRulestackCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			rulestackId := localrulestacks.NewLocalRulestackID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName)
			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retreiving %s: %+v", *id, err)
			}

			cert := *existing.Model

			if metadata.ResourceData.HasChange("description") {
				cert.Properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("audit_comment") {
				cert.Properties.AuditComment = pointer.To(model.AuditComment)
			}

			if metadata.ResourceData.HasChanges("key_vault_certificate_id", "self_signed") {
				cert.Properties.CertificateSelfSigned = boolAsBooleanEnumCert(model.SelfSigned)
				cert.Properties.CertificateSignerResourceId = pointer.To(model.CertificateSignerID)
			}

			if err = client.CreateOrUpdateThenPoll(ctx, *id, cert); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			if err = rulestackClient.CommitThenPoll(ctx, rulestackId); err != nil {
				return fmt.Errorf("committing Local RuleStack config for %s: %+v", id, err)
			}

			return nil
		},
	}
}

func boolAsBooleanEnumCert(input bool) certificateobjectlocalrulestack.BooleanEnum {
	if input {
		return certificateobjectlocalrulestack.BooleanEnumTRUE
	}

	return certificateobjectlocalrulestack.BooleanEnumFALSE
}

func boolEnumAsBoolCert(input certificateobjectlocalrulestack.BooleanEnum) bool {
	return input == certificateobjectlocalrulestack.BooleanEnumTRUE
}
