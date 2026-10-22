// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module_key

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedHSMs.DataPlaneKeysClient
			hsmClient := metadata.Client.ManagedHSMs.ManagedHsmClient

			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
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

			managedHSM, err := hsmClient.Get(ctx, *resourceManagerId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", resourceManagerId, err)
			}

			shouldPurge := metadata.Client.Features.KeyVault.PurgeSoftDeletedHSMKeysOnDestroy
			if shouldPurge && managedHSM.Model != nil && managedHSM.Model.Properties != nil && pointer.From(managedHSM.Model.Properties.EnablePurgeProtection) {
				log.Printf("[DEBUG] cannot purge key %q because Managed HSM %q has purge protection enabled", id.KeyName, id.ManagedHSMName)
				shouldPurge = false
			}

			description := fmt.Sprintf("Key %q (Managed HSM %q)", id.KeyName, id.ManagedHSMName)
			deleter := deleteAndPurgeKey{
				client:      client,
				keyVaultUri: id.BaseUri(),
				name:        id.KeyName,
			}

			return deleteAndOptionallyPurge(ctx, description, shouldPurge, deleter)
		},
	}
}
