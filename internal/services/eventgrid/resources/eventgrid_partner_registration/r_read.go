// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_registration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/partnerregistrations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.PartnerRegistrations

			id, err := partnerregistrations.ParsePartnerRegistrationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := EventGridPartnerRegistrationResourceModel{
				Name:          id.PartnerRegistrationName,
				ResourceGroup: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Tags = pointer.From(model.Tags)
				if props := model.Properties; props != nil {
					state.PartnerRegistrationID = pointer.From(props.PartnerRegistrationImmutableId)
				}
			}
			return metadata.Encode(&state)
		},
	}
}
