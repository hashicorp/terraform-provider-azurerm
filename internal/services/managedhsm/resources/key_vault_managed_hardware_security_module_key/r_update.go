// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module_key

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
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

			var config KeyVaultMHSMKeyResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := parse.ManagedHSMDataPlaneVersionlessKeyID(metadata.ResourceData.Id(), domainSuffix)
			if err != nil {
				return err
			}

			subscriptionId := commonids.NewSubscriptionID(metadata.Client.Account.SubscriptionId)
			resourceManagerId, err := metadata.Client.ManagedHSMs.ManagedHSMIDFromBaseUrl(ctx, subscriptionId, id.BaseUri(), domainSuffix)
			if err != nil {
				return fmt.Errorf("determining Resource Manager ID for %q: %+v", id, err)
			}
			if resourceManagerId == nil {
				return fmt.Errorf("unable to determine the Resource Manager ID for %s", id)
			}

			parameters := keyvault.KeyUpdateParameters{
				KeyOps: expandKeyVaultKeyOptions(config.KeyOpts),
				KeyAttributes: &keyvault.KeyAttributes{
					Enabled: pointer.To(true),
				},

				Tags: tags.Expand(config.Tags),
			}

			if config.NotBeforeDate != "" {
				notBeforeDate, _ := time.Parse(time.RFC3339, config.NotBeforeDate) // validated by schema
				parameters.KeyAttributes.NotBefore = pointer.To(date.UnixTime(notBeforeDate))
			}

			if config.ExpirationDate != "" {
				expirationDate, _ := time.Parse(time.RFC3339, config.ExpirationDate) // validated by schema
				parameters.KeyAttributes.Expires = pointer.To(date.UnixTime(expirationDate))
			}

			resp, err := client.UpdateKey(ctx, id.BaseUri(), config.Name, "", parameters)
			if err != nil {
				return err
			}

			// Managed HSM serves the data plane from multiple partitions; a read immediately after an update
			// may be routed to a stale replica, so we poll until the read-back `updated` timestamp matches
			// the write response to ensure consistency before the framework performs its read
			if resp.Attributes != nil && resp.Attributes.Updated != nil {
				pollerType := custompollers.NewKeyUpdatePoller(metadata.Client.ManagedHSMs.DataPlaneKeysClient, id.BaseUri(), config.Name, *resp.Attributes.Updated)
				poller := pollers.NewPoller(pollerType, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
				if err := poller.PollUntilDone(ctx); err != nil {
					return fmt.Errorf("waiting for key %q update to propagate: %+v", config.Name, err)
				}
			}

			return nil
		},
	}
}
