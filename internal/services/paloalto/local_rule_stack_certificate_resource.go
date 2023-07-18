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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type LocalRulestackCertificate struct{}

var _ sdk.ResourceWithUpdate = LocalRulestackCertificate{}

type LocalRulestackCertificateModel struct {
	Name                string `tfschema:"name"`
	RuleStackID         string `tfschema:"rule_stack_id"`
	AuditComment        string `tfschema:"audit_comment"`
	CertificateSignerID string `tfschema:"certificate_signer_id"`
	Description         string `tfschema:"description"`
	SelfSigned          bool   `tfschema:"self_signed"`
}

func (r LocalRulestackCertificate) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return certificateobjectlocalrulestack.ValidateLocalRulestackCertificateID
}

func (r LocalRulestackCertificate) ResourceType() string {
	return "azurerm_palo_alto_local_rule_stack_certificate"
}

func (r LocalRulestackCertificate) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.LocalRulestackCertificateName,
		},

		"rule_stack_id": {
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

		"certificate_signer_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsURLWithHTTPS, // TODO - Need to investigate valid values for this - What resource Type is it actually?
			ExactlyOneOf: []string{"self_signed", "certificate_signer_id"},
		},

		"self_signed": {
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			ForceNew:     true,
			Default:      false,
			ExactlyOneOf: []string{"certificate_signer_id", "self_signed"},
		},
	}
}

func (r LocalRulestackCertificate) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LocalRulestackCertificate) ModelObject() interface{} {
	return &LocalRulestackCertificateModel{}
}

func (r LocalRulestackCertificate) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.CertificatesClient
			rulstackClient := metadata.Client.PaloAlto.LocalRulestacksClient

			model := LocalRulestackCertificateModel{}
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			ruleStackId, err := localrulestacks.ParseLocalRulestackID(model.RuleStackID)
			if err != nil {
				return err
			}

			id := certificateobjectlocalrulestack.NewLocalRulestackCertificateID(ruleStackId.SubscriptionId, ruleStackId.ResourceGroupName, ruleStackId.LocalRulestackName, model.Name)
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

			if _, err = rulstackClient.Commit(ctx, *ruleStackId); err != nil {
				return fmt.Errorf("committing Local Rulestack config for %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r LocalRulestackCertificate) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.CertificatesClient

			id, err := certificateobjectlocalrulestack.ParseLocalRulestackCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LocalRulestackCertificateModel

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state.Name = id.CertificateName
			state.RuleStackID = certificateobjectlocalrulestack.NewLocalRulestackID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName).ID()

			props := existing.Model.Properties

			state.AuditComment = pointer.From(props.AuditComment)
			state.CertificateSignerID = pointer.From(props.CertificateSignerResourceId)
			state.Description = pointer.From(props.Description)
			state.SelfSigned = boolEnumAsBoolCert(props.CertificateSelfSigned)

			return metadata.Encode(&state)
		},
	}
}

func (r LocalRulestackCertificate) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.CertificatesClient

			id, err := certificateobjectlocalrulestack.ParseLocalRulestackCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r LocalRulestackCertificate) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.CertificatesClient

			model := LocalRulestackCertificateModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id, err := certificateobjectlocalrulestack.ParseLocalRulestackCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

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

			if metadata.ResourceData.HasChanges("certificate_signer_id", "self_signed") {
				cert.Properties.CertificateSelfSigned = boolAsBooleanEnumCert(model.SelfSigned)
				cert.Properties.CertificateSignerResourceId = pointer.To(model.CertificateSignerID)
			}

			if _, err = client.CreateOrUpdate(ctx, *id, cert); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func boolAsBooleanEnumCert(input bool) certificateobjectlocalrulestack.BooleanEnum {
	var result certificateobjectlocalrulestack.BooleanEnum

	if input {
		result = certificateobjectlocalrulestack.BooleanEnumTRUE
	} else {
		result = certificateobjectlocalrulestack.BooleanEnumFALSE
	}

	return result
}

func boolEnumAsBoolCert(input certificateobjectlocalrulestack.BooleanEnum) bool {
	return input == certificateobjectlocalrulestack.BooleanEnumTRUE
}
