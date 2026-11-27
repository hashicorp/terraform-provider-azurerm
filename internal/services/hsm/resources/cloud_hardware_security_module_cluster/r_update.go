// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cloud_hardware_security_module_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hardwaresecuritymodules/2025-03-31/cloudhsmclusters"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hsm/custompollers"
)

func (Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HSM.CloudHsmClustersClient

			id, err := cloudhsmclusters.ParseCloudHsmClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model CloudHardwareSecurityModuleClusterModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			exists, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if exists.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			cluster := *exists.Model
			// systemData has to be nil for update operations or the API will reject the request
			cluster.SystemData = nil

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}

				cluster.Identity = &identity.LegacySystemAndUserAssignedMap{
					Type:        expandedIdentity.Type,
					IdentityIds: expandedIdentity.IdentityIds,
				}
			}

			// the `tags` field should always be set or it will be removed by the API
			if metadata.ResourceData.HasChange("tags") {
				if model.Tags == nil {
					cluster.Tags = pointer.To(map[string]string{})
				} else {
					cluster.Tags = pointer.To(model.Tags)
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, cluster); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			poller := custompollers.NewCloudHsmClusterStatePoller(client, *id)
			pollerType := pollers.NewPoller(poller, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := pollerType.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for update of %s to complete: %+v", *id, err)
			}

			return nil
		},
	}
}
