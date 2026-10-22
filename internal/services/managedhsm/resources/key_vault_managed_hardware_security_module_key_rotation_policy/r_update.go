// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module_key_rotation_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedHSMs.DataPlaneRoleAssignmentsClient
			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			var config MHSMKeyRotationPolicyResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := parse.ManagedHSMDataPlaneVersionlessKeyID(metadata.ResourceData.Id(), domainSuffix)
			if err != nil {
				return err
			}

			resp, err := client.UpdateKeyRotationPolicy(ctx, id.BaseUri(), id.KeyName, expandKeyRotationPolicy(config))
			if err != nil {
				return fmt.Errorf("updating HSM Key Rotation Policy for Key %q: %v", id, err)
			}

			// Managed HSM serves the data plane from multiple partitions; a read immediately after an update
			// may be routed to a stale replica, so we poll until the read-back `updated` timestamp matches
			// the write response to ensure consistency before the framework performs its read
			if resp.Attributes != nil && resp.Attributes.Updated != nil {
				pollerType := custompollers.NewKeyRotationPolicyUpdatePoller(metadata.Client.ManagedHSMs.DataPlaneKeysClient, id.BaseUri(), id.KeyName, *resp.Attributes.Updated)
				poller := pollers.NewPoller(pollerType, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
				if err := poller.PollUntilDone(ctx); err != nil {
					return fmt.Errorf("waiting for key rotation policy for %q to propagate: %+v", id, err)
				}
			}

			return nil
		},
	}
}
