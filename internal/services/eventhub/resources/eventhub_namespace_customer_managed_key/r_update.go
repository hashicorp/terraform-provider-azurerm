// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace_customer_managed_key

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2026-01-01/ehnamespaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventHubNamespaceCustomerManagedKeyCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ehnamespaces.ParseNamespaceID(d.Get("eventhub_namespace_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(id.NamespaceName, "azurerm_eventhub_namespace")
	defer locks.UnlockByName(id.NamespaceName, "azurerm_eventhub_namespace")

	resp, err := client.NamespacesGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			if resp.Model.Properties != nil && resp.Model.Properties.Encryption != nil {
				return tf.ImportAsExistsError("azurerm_eventhub_namespace_customer_managed_key", id.ID())
			}
		}
	}

	namespace := resp.Model
	namespace.Properties.Encryption = &ehnamespaces.Encryption{
		KeySource: pointer.To(ehnamespaces.KeySourceMicrosoftPointKeyVault),
	}

	keyVaultProps, err := expandEventHubNamespaceKeyVaultKeyIds(d.Get("key_vault_key_ids").(*pluginsdk.Set).List())
	if err != nil {
		return err
	}

	userAssignedIdentity := d.Get("user_assigned_identity_id").(string)
	if userAssignedIdentity != "" && keyVaultProps != nil {
		userAssignedIdentityId, err := commonids.ParseUserAssignedIdentityID(userAssignedIdentity)
		if err != nil {
			return err
		}

		// this provides a more helpful error message than the API response
		if namespace.Identity == nil {
			return fmt.Errorf("user assigned identity '%s' must also be assigned to the parent event hub - currently no user assigned identities are assigned to the parent event hub", userAssignedIdentity)
		}

		isIdentityAssignedToParent := false
		for item := range namespace.Identity.IdentityIds {
			parentEhnUaiId, err := commonids.ParseUserAssignedIdentityIDInsensitively(item)
			if err != nil {
				return fmt.Errorf("parsing %q as a User Assigned Identity ID: %+v", item, err)
			}
			if resourceids.Match(parentEhnUaiId, userAssignedIdentityId) {
				isIdentityAssignedToParent = true
			}
		}

		// this provides a more helpful error message than the API response
		if !isIdentityAssignedToParent {
			return fmt.Errorf("user assigned identity '%s' must also be assigned to the parent event hub", userAssignedIdentity)
		}

		for i := 0; i < len(*keyVaultProps); i++ {
			(*keyVaultProps)[i].Identity = &ehnamespaces.UserAssignedIdentityProperties{
				UserAssignedIdentity: &userAssignedIdentity,
			}
		}
	}

	namespace.Properties.Encryption.KeyVaultProperties = keyVaultProps
	namespace.Properties.Encryption.RequireInfrastructureEncryption = pointer.To(d.Get("infrastructure_encryption_enabled").(bool))

	if d.IsNewResource() {
		if err := client.NamespacesCreateOrUpdateCallbackThenPoll(ctx, *id, *namespace, sdk.SetIDCallback(meta, id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", *id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.NamespacesCreateOrUpdateThenPoll(ctx, *id, *namespace); err != nil {
			return fmt.Errorf("updating %s: %+v", *id, err)
		}
	}

	return resourceEventHubNamespaceCustomerManagedKeyRead(d, meta)
}
