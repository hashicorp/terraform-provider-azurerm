// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package communication_service

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/communicationservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.ServiceClient

			state := CommunicationServiceResourceModel{}

			id, err := communicationservices.ParseCommunicationServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state.Name = id.CommunicationServiceName
			state.ResourceGroupName = id.ResourceGroupName

			keysResp, err := client.ListKeys(ctx, *id)
			if err != nil {
				return fmt.Errorf("listing keys for %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.DataLocation = props.DataLocation
					state.HostName = pointer.From(props.HostName)
				}

				state.Tags = pointer.From(model.Tags)
			}

			if model := keysResp.Model; model != nil {
				state.PrimaryConnectionString = pointer.From(model.PrimaryConnectionString)
				state.SecondaryConnectionString = pointer.From(model.SecondaryConnectionString)
				state.PrimaryKey = pointer.From(model.PrimaryKey)
				state.SecondaryKey = pointer.From(model.SecondaryKey)
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}
