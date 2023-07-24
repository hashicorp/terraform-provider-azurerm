package paloalto

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	certificates "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/certificateobjectlocalrulestack"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrulestacks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LocalRulestackOutboundTrustCertificateResource struct{}

type LocalRulestackOutboundTrustCertificateResourceModel struct {
	RulestackID   string `tfschema:"rulestack_id"`
	CertificateID string `tfschema:"certificate_id"`
}

var _ sdk.Resource = LocalRulestackOutboundTrustCertificateResource{}

func (l LocalRulestackOutboundTrustCertificateResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return certificates.ValidateLocalRulestackCertificateID
}

func (l LocalRulestackOutboundTrustCertificateResource) ModelObject() interface{} {
	return &LocalRulestackOutboundTrustCertificateResourceModel{}
}

func (l LocalRulestackOutboundTrustCertificateResource) ResourceType() string {
	return "azurerm_local_rulestack_outbound_trust_certificate_association"
}

func (l LocalRulestackOutboundTrustCertificateResource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"rulestack_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: localrulestacks.ValidateLocalRulestackID,
		},

		"certificate_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: certificates.ValidateLocalRulestackCertificateID,
		},
	}
}

func (l LocalRulestackOutboundTrustCertificateResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (l LocalRulestackOutboundTrustCertificateResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestacksClient

			model := LocalRulestackOutboundTrustCertificateResourceModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			ruleStackId, err := localrulestacks.ParseLocalRulestackID(model.RulestackID)
			if err != nil {
				return err
			}

			certificateId, err := certificates.ParseLocalRulestackCertificateID(model.CertificateID)
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *ruleStackId)
			if err != nil {
				return fmt.Errorf("retrieving the local Rulestack to associate the Outbound Trust Certificate on %s: %+v", *ruleStackId, err)
			}

			rulestack := *existing.Model
			props := rulestack.Properties

			secUpdate := pointer.From(props.SecurityServices)
			secUpdate.OutboundTrustCertificate = pointer.To(certificateId.CertificateName)

			props.SecurityServices = pointer.To(secUpdate)

			rulestack.Properties = props

			if err = client.CreateOrUpdateThenPoll(ctx, *ruleStackId, rulestack); err != nil {
				return fmt.Errorf("creating Outbound Trust Certificate Association for %s: %+v", ruleStackId, err)
			}

			if err = client.CommitThenPoll(ctx, *ruleStackId); err != nil {
				return fmt.Errorf("committing Local Rulestack configurtion for UnTrust Certificate for %s: %+v", ruleStackId, err)
			}

			metadata.SetID(certificateId)

			return nil
		},
	}
}

func (l LocalRulestackOutboundTrustCertificateResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestacksClient

			certificateId, err := certificates.ParseLocalRulestackCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			rulestackId := localrulestacks.NewLocalRulestackID(certificateId.SubscriptionId, certificateId.ResourceGroupName, certificateId.LocalRulestackName)

			var state LocalRulestackOutboundTrustCertificateResourceModel

			existing, err := client.Get(ctx, rulestackId)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(rulestackId)
				}
				return fmt.Errorf("reading %s for Outbound Trust Association: %+v", rulestackId, err)
			}

			props := existing.Model.Properties
			secServices := pointer.From(props.SecurityServices)

			state.RulestackID = rulestackId.ID()
			state.CertificateID = certificates.NewLocalRulestackCertificateID(certificateId.SubscriptionId, certificateId.ResourceGroupName, certificateId.LocalRulestackName, pointer.From(secServices.OutboundTrustCertificate)).ID()

			return metadata.Encode(&state)
		},
	}
}

func (l LocalRulestackOutboundTrustCertificateResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestacksClient

			certId, err := certificates.ParseLocalRulestackCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			rulestackId := localrulestacks.NewLocalRulestackID(certId.SubscriptionId, certId.ResourceGroupName, certId.LocalRulestackName)

			existing, err := client.Get(ctx, rulestackId)
			if err != nil {
				return fmt.Errorf("retrieving the local Rulestack to disassociate the Outbound Trust Certificate on %s: %+v", rulestackId, err)
			}

			rulestack := *existing.Model
			props := rulestack.Properties
			secServices := pointer.From(props.SecurityServices)

			secServices.OutboundTrustCertificate = nil
			props.SecurityServices = pointer.To(secServices)

			rulestack.Properties = props

			if err = client.CreateOrUpdateThenPoll(ctx, rulestackId, rulestack); err != nil {
				return fmt.Errorf("deleting Local Rulestack Outbound Trust Certificate Association for %s: %+v", rulestackId, err)
			}

			if err = client.CommitThenPoll(ctx, rulestackId); err != nil {
				return fmt.Errorf("committing rulestack config for removal of Trust Certificate from %s: %+v", rulestackId, err)
			}

			return nil
		},
	}
}
