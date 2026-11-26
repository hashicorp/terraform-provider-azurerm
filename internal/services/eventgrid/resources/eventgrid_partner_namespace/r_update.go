// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_namespace

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/partnernamespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.PartnerNamespaces

			id, err := partnernamespaces.ParsePartnerNamespaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config EventGridPartnerNamespaceResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `model.Properties` was nil", *id)
			}

			model := existing.Model

			if metadata.ResourceData.HasChange("local_authentication_enabled") {
				model.Properties.DisableLocalAuth = pointer.To(!config.LocalAuthEnabled)
			}
			if metadata.ResourceData.HasChange("inbound_ip_rule") {
				model.Properties.InboundIPRules = expandPartnerInboundIPRules(config.InboundIPRules)
			}
			if metadata.ResourceData.HasChange("public_network_access") {
				model.Properties.PublicNetworkAccess = pointer.ToEnum[partnernamespaces.PublicNetworkAccess](config.PublicNetworkAccess)
			}
			if metadata.ResourceData.HasChange("tags") {
				model.Tags = pointer.To(config.Tags)
			}

			// endpoint is read-only and will throw an error if we keep it in the update
			model.Properties.Endpoint = nil

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
