// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/deletedvaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	managementClient := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseKeyVaultID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	vaultUri := ""
	if model := resp.Model; model != nil {
		if model.Properties.VaultUri != nil {
			vaultUri = *model.Properties.VaultUri
		}
	}

	if vaultUri != "" {
		meta.(*clients.Client).KeyVault.AddToCache(*id, vaultUri)
	}

	return resourceKeyVaultFlatten(ctx, managementClient, d, id, resp.Model, true)
}

func keyVaultRefreshFunc(ctx context.Context, vaultUri string) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		log.Printf("[DEBUG] Checking to see if KeyVault %q is available..", vaultUri)

		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, vaultUri, nil)
		if err != nil {
			return nil, "pending", fmt.Errorf("building request to check %q: %s", vaultUri, err)
		}

		conn, err := client.Do(req)
		if err != nil {
			log.Printf("[DEBUG] Didn't find KeyVault at %q", vaultUri)
			return nil, "pending", fmt.Errorf("connecting to %q: %s", vaultUri, err)
		}

		defer conn.Body.Close()

		log.Printf("[DEBUG] Found KeyVault at %q", vaultUri)
		return "available", "available", nil
	}
}

func optedOutOfRecoveringSoftDeletedKeyVaultErrorFmt(name, location string) string {
	return fmt.Sprintf(`
An existing soft-deleted Key Vault exists with the Name %q in the location %q, however
automatically recovering this KeyVault has been disabled via the "features" block.

Terraform can automatically recover the soft-deleted Key Vault when this behaviour is
enabled within the "features" block (located within the "provider" block) - more
information can be found here:

https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block

Alternatively you can manually recover this (e.g. using the Azure CLI) and then import
this into Terraform via "terraform import", or pick a different name/location.
`, name, location)
}

type keyVaultDeletionStatus struct {
	deleteDate string
	purgeDate  string
}

func getSoftDeletedStateForKeyVault(ctx context.Context, deletedVaultsClient *deletedvaults.DeletedVaultsClient, deletedVaultId deletedvaults.DeletedVaultId) (*keyVaultDeletionStatus, error) {
	resp, err := deletedVaultsClient.VaultsGetDeleted(ctx, deletedVaultId)
	if err != nil {
		return nil, err
	}

	if model := resp.Model; model != nil {
		// the logic is this way because the GetDeleted call will return an existing key vault
		// that is not soft deleted, but the Deleted Vault properties will be nil
		result := keyVaultDeletionStatus{}
		if props := model.Properties; props != nil {
			if props.DeletionDate != nil {
				t, err := props.GetDeletionDateAsTime()
				if err != nil {
					return nil, fmt.Errorf("parsing `DeletionDate`: %+v", err)
				}
				result.deleteDate = t.Format(time.RFC3339)
			}
			if props.ScheduledPurgeDate != nil {
				t, err := props.GetScheduledPurgeDateAsTime()
				if err != nil {
					return nil, fmt.Errorf("parsing `ScheduledPurgeDate`: %+v", err)
				}
				result.purgeDate = t.Format(time.RFC3339)
			}
		}

		if result.deleteDate != "" && result.purgeDate != "" {
			return &result, nil
		}
	}

	// otherwise we've found an existing key vault that is not soft deleted
	return nil, nil
}
