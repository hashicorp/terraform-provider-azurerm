// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_network_connection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/networkconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.NetworkConnections

			id, err := networkconnections.ParseNetworkConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model DevCenterNetworkConnectionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := networkconnections.NetworkConnectionUpdate{
				Properties: &networkconnections.NetworkConnectionUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("subnet_id") {
				parameters.Properties.SubnetId = pointer.To(model.SubnetId)
			}

			if metadata.ResourceData.HasChange("domain_name") {
				parameters.Properties.DomainName = pointer.To(model.DomainName)
			}

			if metadata.ResourceData.HasChange("domain_password") {
				parameters.Properties.DomainPassword = pointer.To(model.DomainPassword)
			}

			if metadata.ResourceData.HasChange("domain_username") {
				parameters.Properties.DomainUsername = pointer.To(model.DomainUsername)
			}

			if metadata.ResourceData.HasChange("organization_unit") {
				parameters.Properties.OrganizationUnit = pointer.To(model.OrganizationUnit)
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = pointer.To(model.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
