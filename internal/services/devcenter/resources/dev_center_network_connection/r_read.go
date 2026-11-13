// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_network_connection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/networkconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.NetworkConnections

			id, err := networkconnections.ParseNetworkConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := DevCenterNetworkConnectionResourceModel{
				Name:              id.NetworkConnectionName,
				ResourceGroupName: id.ResourceGroupName,
				DomainPassword:    metadata.ResourceData.Get("domain_password").(string),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.SubnetId = pointer.From(props.SubnetId)
					state.DomainName = pointer.From(props.DomainName)
					state.DomainUsername = pointer.From(props.DomainUsername)
					state.OrganizationUnit = pointer.From(props.OrganizationUnit)

					if v := props.DomainJoinType; v != "" {
						state.DomainJoinType = string(v)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}
