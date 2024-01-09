// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LocalRulestackOutboundTrustCertificateAssociationResource struct{}

type LocalRulestackOutboundTrustCertificateResourceModel struct {
	CertificateID string `tfschema:"certificate_id"`
}

var _ sdk.Resource = LocalRulestackOutboundTrustCertificateAssociationResource{}

func (l LocalRulestackOutboundTrustCertificateAssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return certificates.ValidateLocalRulestackCertificateID
}

func (l LocalRulestackOutboundTrustCertificateAssociationResource) ModelObject() interface{} {
	return &LocalRulestackOutboundTrustCertificateResourceModel{}
}

func (l LocalRulestackOutboundTrustCertificateAssociationResource) ResourceType() string {
	return "azurerm_palo_alto_local_rulestack_outbound_trust_certificate_association"
}

func (l LocalRulestackOutboundTrustCertificateAssociationResource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"certificate_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: certificates.ValidateLocalRulestackCertificateID,
		},
	}
}

func (l LocalRulestackOutboundTrustCertificateAssociationResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (l LocalRulestackOutboundTrustCertificateAssociationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.Client.LocalRulestacks

			model := LocalRulestackOutboundTrustCertificateResourceModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			certificateId, err := certificates.ParseLocalRulestackCertificateID(model.CertificateID)
			if err != nil {
				return err
			}

			locks.ByID(certificateId.ID())
			defer locks.UnlockByID(certificateId.ID())
			rulestackId := localrulestacks.NewLocalRulestackID(certificateId.SubscriptionId, certificateId.ResourceGroupName, certificateId.LocalRulestackName)

			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			existing, err := client.Get(ctx, rulestackId)
			if err != nil {
				return fmt.Errorf("retrieving the local Rulestack to associate the Outbound Trust Certificate on %s: %+v", rulestackId, err)
			}

			rulestack := *existing.Model
			props := rulestack.Properties

			secUpdate := pointer.From(props.SecurityServices)
			secUpdate.OutboundTrustCertificate = pointer.To(certificateId.CertificateName)

			props.SecurityServices = pointer.To(secUpdate)

			rulestack.Properties = props

			if err = client.CreateOrUpdateThenPoll(ctx, rulestackId, rulestack); err != nil {
				return fmt.Errorf("creating Outbound Trust Certificate Association for %s: %+v", rulestackId, err)
			}

			if err = client.CommitThenPoll(ctx, rulestackId); err != nil {
				return fmt.Errorf("committing Local Rulestack configuration for Outbound Trust Certificate for %s: %+v", rulestackId, err)
			}

			metadata.SetID(certificateId)

			return nil
		},
	}
}

func (l LocalRulestackOutboundTrustCertificateAssociationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.Client.LocalRulestacks

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
			if model := existing.Model; model != nil {
				props := model.Properties
				secServices := pointer.From(props.SecurityServices)

				state.CertificateID = certificates.NewLocalRulestackCertificateID(certificateId.SubscriptionId, certificateId.ResourceGroupName, certificateId.LocalRulestackName, pointer.From(secServices.OutboundTrustCertificate)).ID()
			}

			return metadata.Encode(&state)
		},
	}
}

func (l LocalRulestackOutboundTrustCertificateAssociationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.Client.LocalRulestacks

			certId, err := certificates.ParseLocalRulestackCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			locks.ByID(certId.ID())
			defer locks.UnlockByID(certId.ID())

			rulestackId := localrulestacks.NewLocalRulestackID(certId.SubscriptionId, certId.ResourceGroupName, certId.LocalRulestackName)
			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

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
