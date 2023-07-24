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

type LocalRulestackOutboundUnTrustCertificateResource struct{}

type LocalRulestackOutboundUnTrustCertificateResourceModel struct {
	RulestackID   string `tfschema:"rulestack_id"`
	CertificateID string `tfschema:"certificate_id"`
}

var _ sdk.Resource = LocalRulestackOutboundUnTrustCertificateResource{}

func (l LocalRulestackOutboundUnTrustCertificateResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return certificates.ValidateLocalRulestackCertificateID
}

func (l LocalRulestackOutboundUnTrustCertificateResource) ModelObject() interface{} {
	return &LocalRulestackOutboundUnTrustCertificateResourceModel{}
}

func (l LocalRulestackOutboundUnTrustCertificateResource) ResourceType() string {
	return "azurerm_local_rulestack_outbound_untrust_certificate_association"
}

func (l LocalRulestackOutboundUnTrustCertificateResource) Arguments() map[string]*schema.Schema {
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

func (l LocalRulestackOutboundUnTrustCertificateResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (l LocalRulestackOutboundUnTrustCertificateResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestacksClient

			model := LocalRulestackOutboundUnTrustCertificateResourceModel{}

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
				return fmt.Errorf("retrieving the local Rulestack to associate the Outbound UnTrust Certificate on %s: %+v", *ruleStackId, err)
			}

			rulestack := *existing.Model
			props := rulestack.Properties

			secUpdate := pointer.From(props.SecurityServices)
			secUpdate.OutboundUnTrustCertificate = pointer.To(certificateId.CertificateName)

			props.SecurityServices = pointer.To(secUpdate)

			rulestack.Properties = props

			if err = client.CreateOrUpdateThenPoll(ctx, *ruleStackId, rulestack); err != nil {
				return fmt.Errorf("creating Outbound UnTrust association for %s: %+v", ruleStackId, err)
			}

			if err = client.CommitThenPoll(ctx, *ruleStackId); err != nil {
				return fmt.Errorf("committing rulestack config for UnTrust Certificate for %s: %+v", ruleStackId, err)
			}

			metadata.SetID(certificateId)

			return nil
		},
	}
}

func (l LocalRulestackOutboundUnTrustCertificateResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestacksClient

			certificateId, err := certificates.ParseLocalRulestackCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			rulestackId := localrulestacks.NewLocalRulestackID(certificateId.SubscriptionId, certificateId.ResourceGroupName, certificateId.LocalRulestackName)

			var state LocalRulestackOutboundUnTrustCertificateResourceModel

			existing, err := client.Get(ctx, rulestackId)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(rulestackId)
				}
				return fmt.Errorf("reading %s for Outbound UnTrust Association: %+v", rulestackId, err)
			}

			props := existing.Model.Properties
			secServices := pointer.From(props.SecurityServices)

			state.RulestackID = rulestackId.ID()
			state.CertificateID = certificates.NewLocalRulestackCertificateID(certificateId.SubscriptionId, certificateId.ResourceGroupName, certificateId.LocalRulestackName, pointer.From(secServices.OutboundUnTrustCertificate)).ID()

			return metadata.Encode(&state)
		},
	}
}

func (l LocalRulestackOutboundUnTrustCertificateResource) Delete() sdk.ResourceFunc {
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
				return fmt.Errorf("retrieving the local Rulestack to disassociate the Outbound UnTrust Certificate on %s: %+v", rulestackId, err)
			}

			rulestack := *existing.Model
			props := rulestack.Properties
			secServices := pointer.From(props.SecurityServices)

			secServices.OutboundUnTrustCertificate = nil
			props.SecurityServices = pointer.To(secServices)
			rulestack.Properties = props

			if err = client.CreateOrUpdateThenPoll(ctx, rulestackId, rulestack); err != nil {
				return fmt.Errorf("deleting Local Rulestack Outbound UnTrust Certificate Association for %s: %+v", rulestackId, err)
			}

			if err = client.CommitThenPoll(ctx, rulestackId); err != nil {
				return fmt.Errorf("committing rulestack config for removing UnTrust Certificate for %s: %+v", rulestackId, err)
			}

			return nil
		},
	}
}
