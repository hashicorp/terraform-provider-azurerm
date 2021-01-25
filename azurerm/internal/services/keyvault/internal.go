package keyvault

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type deleteAndPurgeNestedItem interface {
	DeleteNestedItem(ctx context.Context) (autorest.Response, error)
	NestedItemHasBeenDeleted(ctx context.Context) (autorest.Response, error)

	PurgeNestedItem(ctx context.Context) (autorest.Response, error)
	NestedItemHasBeenPurged(ctx context.Context) (autorest.Response, error)
}

func deleteAndOptionallyPurge(ctx context.Context, description string, shouldPurge bool, helper deleteAndPurgeNestedItem) error {
	timeout, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context is missing a timeout")
	}

	log.Printf("[DEBUG] Deleting %s..", description)
	if resp, err := helper.DeleteNestedItem(ctx); err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", description, err)
	}
	log.Printf("[DEBUG] Waiting for %s to finish deleting..", description)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"InProgress"},
		Target:  []string{"NotFound"},
		Refresh: func() (interface{}, string, error) {
			item, err := helper.NestedItemHasBeenDeleted(ctx)
			if err != nil {
				if utils.ResponseWasNotFound(item) {
					return item, "NotFound", nil
				}

				return nil, "Error", err
			}

			return item, "InProgress", nil
		},
		ContinuousTargetOccurence: 3,
		PollInterval:              5 * time.Second,
		Timeout:                   time.Until(timeout),
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", description, err)
	}
	log.Printf("[DEBUG] Deleted %s.", description)

	if !shouldPurge {
		log.Printf("[DEBUG] Skipping purging of %s as opted-out..", description)
		return nil
	}

	log.Printf("[DEBUG] Purging %s..", description)
	if _, err := helper.PurgeNestedItem(ctx); err != nil {
		return fmt.Errorf("purging %s: %+v", description, err)
	}

	log.Printf("[DEBUG] Waiting for %s to finish purging..", description)
	stateConf = &resource.StateChangeConf{
		Pending: []string{"InProgress"},
		Target:  []string{"NotFound"},
		Refresh: func() (interface{}, string, error) {
			item, err := helper.NestedItemHasBeenPurged(ctx)
			if err != nil {
				if utils.ResponseWasNotFound(item) {
					return item, "NotFound", nil
				}

				return nil, "Error", err
			}

			return item, "InProgress", nil
		},
		ContinuousTargetOccurence: 3,
		PollInterval:              5 * time.Second,
		Timeout:                   time.Until(timeout),
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for %s to finish purging: %+v", description, err)
	}
	log.Printf("[DEBUG] Purged %s.", description)

	return nil
}

func keyVaultChildItemRefreshFunc(secretUri string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking to see if KeyVault Secret %q is available..", secretUri)

		PTransport := &http.Transport{Proxy: http.ProxyFromEnvironment}

		client := &http.Client{
			Transport: PTransport,
		}

		conn, err := client.Get(secretUri)
		if err != nil {
			log.Printf("[DEBUG] Didn't find KeyVault secret at %q", secretUri)
			return nil, "pending", fmt.Errorf("Error checking secret at %q: %s", secretUri, err)
		}

		defer conn.Body.Close()

		log.Printf("[DEBUG] Found KeyVault Secret %q", secretUri)
		return "available", "available", nil
	}
}

func nestedItemResourceImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseNestedItemID(d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("parsing ID %q for Key Vault Child import: %v", d.Id(), err)
	}

	keyVaultId, err := azure.GetKeyVaultIDFromBaseUrl(ctx, client, id.KeyVaultBaseUrl)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	d.Set("key_vault_id", keyVaultId)

	return []*schema.ResourceData{d}, nil
}
