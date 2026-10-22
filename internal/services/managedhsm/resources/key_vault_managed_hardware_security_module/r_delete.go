// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/deletedmanagedhsms"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/managedhsms"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmKeyVaultManagedHardwareSecurityModuleDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ManagedHSMs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managedhsms.ParseManagedHSMID(d.Id())
	if err != nil {
		return err
	}

	// We need to grab the managed hsm to see if purge protection is enabled prior to deletion
	resp, err := client.ManagedHsmClient.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	loc := ""
	purgeProtectionEnabled := false
	if model := resp.Model; model != nil {
		loc = location.NormalizeNilable(model.Location)
		if props := model.Properties; props != nil && props.EnablePurgeProtection != nil {
			purgeProtectionEnabled = *props.EnablePurgeProtection
		}
	}

	if err := client.ManagedHsmClient.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if meta.(*clients.Client).Features.KeyVault.PurgeSoftDeletedHSMsOnDestroy {
		if purgeProtectionEnabled {
			log.Printf("[DEBUG] cannot purge %s because purge protection is enabled", id)
			return nil
		}
	}

	// the polling operation of purge can not terminate correctly, so we use the custom polling operation of polling delete
	// try to purge again if managed HSM still exists after 1 minute
	// for API issue: https://github.com/Azure/azure-rest-api-specs/issues/27138
	purgeId := deletedmanagedhsms.NewDeletedManagedHSMID(id.SubscriptionId, loc, id.ManagedHSMName)
	if _, err := client.DeletedManagedHsmClient.ManagedHsmsPurgeDeleted(ctx, purgeId); err != nil {
		return fmt.Errorf("purging %s: %+v", id, err)
	}

	purgePoller := custompollers.NewHSMPurgePoller(client.DeletedManagedHsmClient, purgeId)
	poller := pollers.NewPoller(purgePoller, time.Second*30, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be purged: %+v", id, err)
	}

	client.Purge(*id)

	return nil
}
