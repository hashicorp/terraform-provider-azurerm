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

func (l Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestackResources

			certId, err := certificateobjectlocalrulestackresources.ParseLocalRulestackCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			locks.ByID(certId.ID())
			defer locks.UnlockByID(certId.ID())

			rulestackId := localrulestackresources.NewLocalRulestackID(certId.SubscriptionId, certId.ResourceGroupName, certId.LocalRulestackName)
			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			existing, err := client.LocalRulestacksGet(ctx, rulestackId)
			if err != nil {
				return fmt.Errorf("retrieving the local Rulestack to disassociate the Outbound Trust Certificate on %s: %+v", rulestackId, err)
			}

			rulestack := *existing.Model
			props := rulestack.Properties
			secServices := pointer.From(props.SecurityServices)

			secServices.OutboundTrustCertificate = nil
			props.SecurityServices = pointer.To(secServices)

			rulestack.Properties = props

			if err = client.LocalRulestacksCreateOrUpdateThenPoll(ctx, rulestackId, rulestack); err != nil {
				return fmt.Errorf("deleting Local Rulestack Outbound Trust Certificate Association for %s: %+v", rulestackId, err)
			}

			if err = client.LocalRulestackscommitThenPoll(ctx, rulestackId); err != nil {
				return fmt.Errorf("committing rulestack config for removal of Trust Certificate from %s: %+v", rulestackId, err)
			}

			return nil
		},
	}
}
