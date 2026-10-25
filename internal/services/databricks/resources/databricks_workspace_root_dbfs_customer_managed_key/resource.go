// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package databricks_workspace_root_dbfs_customer_managed_key

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -parent-id "workspace_id"

func RegisterResource() *pluginsdk.Resource {
	r := &pluginsdk.Resource{
		Create: databricksWorkspaceRootDbfsCustomerManagedKeyCreate,
		Read:   databricksWorkspaceRootDbfsCustomerManagedKeyRead,
		Update: databricksWorkspaceRootDbfsCustomerManagedKeyUpdate,
		Delete: databricksWorkspaceRootDbfsCustomerManagedKeyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingIdentityThen(&workspaces.WorkspaceId{}, func(ctx context.Context, d *pluginsdk.ResourceData, meta any) ([]*pluginsdk.ResourceData, error) {
			// validate that the passed ID is a valid CMK configuration ID
			id, err := workspaces.ParseWorkspaceID(d.Id())
			if err != nil {
				return []*pluginsdk.ResourceData{d}, err
			}

			// set the new values for the CMK resource
			d.SetId(id.ID())
			d.Set("workspace_id", id.ID())

			return []*pluginsdk.ResourceData{d}, nil
		}),
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&workspaces.WorkspaceId{}),
		},

		Schema: databricksWorkspaceRootDbfsCustomerManagedKeySchema(),
	}

	if !features.SixPointOh() {
		r.Schema["key_vault_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateKeyVaultID,
			DiffSuppressFunc: func(_, o, n string, _ *pluginsdk.ResourceData) bool {
				// Suppress removal diff for 5.x since that does not require an update
				return o != "" && n == ""
			},
			Deprecated: "`key_vault_id` has been deprecated and will be removed in v6.0 of the AzureRM provider. This property is no longer required for cross-subscription scenarios.",
		}
	}

	return r
}
