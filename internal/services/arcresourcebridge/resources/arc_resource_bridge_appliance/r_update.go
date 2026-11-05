// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_resource_bridge_appliance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resourceconnector/2022-10-27/appliances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcResourceBridge.AppliancesClient

			id, err := appliances.ParseApplianceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ApplianceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			parameters := resp.Model
			if parameters == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if metadata.ResourceData.HasChange("identity") {
				identity, err := identity.ExpandSystemAssignedFromModel(model.Identity)
				if err != nil {
					return err
				}

				parameters.Identity = identity
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = tags.Expand(model.Tags)
			}

			if metadata.ResourceData.HasChanges("public_key_base64") {
				if parameters.Properties == nil {
					parameters.Properties = &appliances.ApplianceProperties{}
				}
				parameters.Properties.PublicKey = pointer.To(model.PublicKeyBase64)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
