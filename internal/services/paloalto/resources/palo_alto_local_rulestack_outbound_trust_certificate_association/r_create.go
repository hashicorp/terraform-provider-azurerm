// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack_outbound_trust_certificate_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/certificateobjectlocalrulestackresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulestackresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (l Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestackResources

			model := LocalRulestackOutboundTrustCertificateResourceModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			certificateId, err := certificateobjectlocalrulestackresources.ParseLocalRulestackCertificateID(model.CertificateID)
			if err != nil {
				return err
			}

			locks.ByID(certificateId.ID())
			defer locks.UnlockByID(certificateId.ID())
			rulestackId := localrulestackresources.NewLocalRulestackID(certificateId.SubscriptionId, certificateId.ResourceGroupName, certificateId.LocalRulestackName)

			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			existing, err := client.LocalRulestacksGet(ctx, rulestackId)
			if err != nil {
				return fmt.Errorf("retrieving the local Rulestack to associate the Outbound Trust Certificate on %s: %+v", rulestackId, err)
			}

			rulestack := *existing.Model
			props := rulestack.Properties

			secUpdate := pointer.From(props.SecurityServices)
			secUpdate.OutboundTrustCertificate = pointer.To(certificateId.CertificateName)

			props.SecurityServices = pointer.To(secUpdate)

			rulestack.Properties = props

			if err := client.LocalRulestacksCreateOrUpdateCallbackThenPoll(ctx, rulestackId, rulestack, metadata.SetIDCallback(certificateId)); err != nil {
				return fmt.Errorf("creating Outbound Trust Certificate Association for %s: %+v", rulestackId, err)
			}

			metadata.SetID(certificateId)

			if err := client.LocalRulestackscommitThenPoll(ctx, rulestackId); err != nil {
				return fmt.Errorf("committing Local Rulestack configuration for Outbound Trust Certificate for %s: %+v", rulestackId, err)
			}

			return nil
		},
	}
}
