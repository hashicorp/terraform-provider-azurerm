// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databricks_workspace_root_dbfs_customer_managed_key

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func databricksWorkspaceRootDbfsCustomerManagedKeyRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if model := resp.Model; model != nil {
		if params := model.Properties.Parameters; params != nil {
			if encryption := params.Encryption; encryption != nil {
				if value := encryption.Value; value != nil {
					if strings.EqualFold(pointer.FromEnum(value.KeySource), string(workspaces.KeySourceDefault)) && value.Keyvaulturi == nil && value.KeyName == nil {
						d.SetId("")
						return nil
					}

					key, err := keyvault.NewNestedItemID(pointer.From(value.Keyvaulturi), keyvault.NestedItemTypeKey, pointer.From(value.KeyName), pointer.From(value.Keyversion))
					if err != nil {
						return err
					}
					d.Set("key_vault_key_id", key.ID())
				}
			}
		}
	}

	d.Set("workspace_id", id.ID())
	if !features.SixPointOh() {
		d.Set("key_vault_id", d.Get("key_vault_id").(string))
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
