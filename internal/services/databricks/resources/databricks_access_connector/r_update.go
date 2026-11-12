// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package databricks_access_connector

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/accessconnector"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataBricks.AccessConnectorClient
			id, err := accessconnector.ParseAccessConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state AccessConnectorResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Databricks %s: %v", id, err)
			}

			if metadata.ResourceData.HasChange("identity") {
				// TODO: Switch this to 'identity.ExpandSystemOrSingleUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))'
				// once SDK Helpers PR #164 has been merged and integrated into the provider...
				identityValue, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]any))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}

				existing.Model.Identity = identityValue
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = &state.Tags
			}

			if err = client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating Databricks %s: %+v", id, err)
			}

			return nil
		},
	}
}
