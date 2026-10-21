// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package fluid_relay_server

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/fluidrelay/2022-05-26/fluidrelayservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (s Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.FluidRelay.FluidRelayServers
			id, err := fluidrelayservers.ParseFluidRelayServerID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ServerModel
			if err = meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			payload := fluidrelayservers.FluidRelayServerUpdate{}
			if meta.ResourceData.HasChange("tags") {
				payload.Tags = &model.Tags
			}
			if meta.ResourceData.HasChange("identity") {
				payload.Identity, err = identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding user identities: %+v", err)
				}
			}
			if meta.ResourceData.HasChange("customer_managed_key") {
				payload.Properties = &fluidrelayservers.FluidRelayServerUpdateProperties{
					Encryption: expandFluidRelayServerCustomerManagedKey(model.CustomerManagedKey),
				}
			}

			if _, err = client.Update(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}
