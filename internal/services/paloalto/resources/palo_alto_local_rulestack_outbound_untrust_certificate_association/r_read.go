// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack_outbound_untrust_certificate_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/certificateobjectlocalrulestackresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulestackresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (l Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestackResources

			certificateId, err := certificateobjectlocalrulestackresources.ParseLocalRulestackCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			rulestackId := localrulestackresources.NewLocalRulestackID(certificateId.SubscriptionId, certificateId.ResourceGroupName, certificateId.LocalRulestackName)

			var state LocalRulestackOutboundUnTrustCertificateResourceModel

			existing, err := client.LocalRulestacksGet(ctx, rulestackId)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(rulestackId)
				}
				return fmt.Errorf("reading %s for Outbound UnTrust Association: %+v", rulestackId, err)
			}
			if model := existing.Model; model != nil {
				props := model.Properties
				secServices := pointer.From(props.SecurityServices)

				state.CertificateID = certificateobjectlocalrulestackresources.NewLocalRulestackCertificateID(certificateId.SubscriptionId, certificateId.ResourceGroupName, certificateId.LocalRulestackName, pointer.From(secServices.OutboundUnTrustCertificate)).ID()
			}

			return metadata.Encode(&state)
		},
	}
}
