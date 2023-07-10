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

type LocalRuleStackCertificate struct{}

var _ sdk.ResourceWithUpdate = LocalRuleStackCertificate{}

type LocalRuleStackCertificateModel struct {
	Name                string `tfschema:"name"`
	RuleStackID         string `tfschema:"rule_stack_id"`
	AuditComment        string `tfschema:"audit_comment"`
	CertificateSignerID string `tfschema:"certificate_signer_id"`
	Description         string `tfschema:"description"`
	SelfSigned          bool   `tfschema:"self_signed"`
}

func (r LocalRuleStackCertificate) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return certificateobjectlocalrulestack.ValidateLocalRuleStackCertificateID
}

func (r LocalRuleStackCertificate) ResourceType() string {
	return "azurerm_palo_alto_local_rule_stack_certificate"
}

func (r LocalRuleStackCertificate) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.LocalRuleStackCertificateName,
		},

		"rule_stack_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: certificateobjectlocalrulestack.ValidateLocalRuleStackID,
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
			ValidateFunc: validation.StringIsNotEmpty, // TODO - Need to investigate valid values for this - What resource Type is it actually?
			ExactlyOneOf: []string{"self_signed", "certificate_signer_id"},
		},

		"self_signed": {
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			Default:      false,
			ExactlyOneOf: []string{"certificate_signer_id", "self_signed"},
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
			client := metadata.Client.PaloAlto.CertificatesClient

			model := LocalRuleStackCertificateModel{}
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			ruleStackId, err := localrulestacks.ParseLocalRuleStackID(model.RuleStackID)
			if err != nil {
				return err
			}

			id := certificateobjectlocalrulestack.NewLocalRuleStackCertificateID(ruleStackId.SubscriptionId, ruleStackId.ResourceGroupName, ruleStackId.LocalRuleStackName, model.Name)
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

			return nil
		},
	}
}

func (r LocalRuleStackCertificate) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.CertificatesClient

			id, err := certificateobjectlocalrulestack.ParseLocalRuleStackCertificateID(metadata.ResourceData.Id())
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
			state.RuleStackID = certificateobjectlocalrulestack.NewLocalRuleStackID(id.SubscriptionId, id.ResourceGroupName, id.LocalRuleStackName).ID()

			props := existing.Model.Properties

			state.AuditComment = pointer.From(props.AuditComment)
			state.CertificateSignerID = pointer.From(props.CertificateSignerResourceId)
			state.Description = pointer.From(props.Description)
			state.SelfSigned = boolEnumAsBoolCert(props.CertificateSelfSigned)

			return metadata.Encode(&state)
		},
	}
}

func (r LocalRuleStackCertificate) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.CertificatesClient

			id, err := certificateobjectlocalrulestack.ParseLocalRuleStackCertificateID(metadata.ResourceData.Id())
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

func (r LocalRuleStackCertificate) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO -
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
