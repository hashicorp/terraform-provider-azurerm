package databricks

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/databricks/mgmt/2018-04-01/databricks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databricks/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databricks/validate"
	keyVaultParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDatabricksWorkspaceCustomerManagedKey() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: DatabricksWorkspaceCustomerManagedKeyCreateUpdate,
		Read:   DatabricksWorkspaceCustomerManagedKeyRead,
		Update: DatabricksWorkspaceCustomerManagedKeyCreateUpdate,
		Delete: DatabricksWorkspaceCustomerManagedKeyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.CustomerManagedKeyID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
			client := meta.(*clients.Client).DataBricks.WorkspacesClient

			// validate that the passed ID is a valid CMK configuration ID
			customManagedKey, err := parse.CustomerManagedKeyID(d.Id())
			if err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("parsing Databricks workspace customer managed key ID %q for import: %v", d.Id(), err)
			}

			// convert the passed custom Managed Key ID to a valid workspace ID
			workspace := parse.NewWorkspaceID(customManagedKey.SubscriptionId, customManagedKey.ResourceGroup, customManagedKey.CustomerMangagedKeyName)

			// validate that the workspace exists
			if _, err = client.Get(ctx, workspace.ResourceGroup, workspace.Name); err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("retrieving the Databricks workspace customer managed key configuration(ID: %q) for workspace (ID: %q): %s", customManagedKey.ID(), workspace.ID(), err)
			}

			// set the new values for the CMK resource
			d.Set("id", customManagedKey.ID())
			d.Set("workspace_id", workspace.ID())

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Schema: map[string]*pluginsdk.Schema{
			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			// Make this key vault key id and abstract everything from the string...
			"key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: keyVaultValidate.KeyVaultChildID,
			},
		},
	}
}

func DatabricksWorkspaceCustomerManagedKeyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	workspaceClient := meta.(*clients.Client).DataBricks.WorkspacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceIDRaw := d.Get("workspace_id").(string)
	workspaceID, err := parse.WorkspaceID(workspaceIDRaw)
	if err != nil {
		return err
	}

	keyIdRaw := d.Get("key_vault_key_id").(string)
	key, err := keyVaultParse.ParseNestedItemID(keyIdRaw)
	if err != nil {
		return err
	}

	// Not sure if I should also lock the key vault here too
	// or at the very least the key?
	locks.ByName(workspaceID.Name, "azurerm_databricks_workspace")
	defer locks.UnlockByName(workspaceID.Name, "azurerm_databricks_workspace")

	workspace, err := workspaceClient.Get(ctx, workspaceID.ResourceGroup, workspaceID.Name)
	if err != nil {
		return fmt.Errorf("retrieving Databricks Workspace %q (Resource Group %q): %+v", workspaceID.Name, workspaceID.ResourceGroup, err)
	}
	if workspace.Parameters == nil {
		return fmt.Errorf("retrieving Databricks Workspace %q (Resource Group %q): `WorkspaceCustomParameters` was nil", workspaceID.Name, workspaceID.ResourceGroup)
	}

	infrastructureEnabled := *workspace.Parameters.RequireInfrastructureEncryption.Value
	encryptionEnabled := *workspace.Parameters.PrepareEncryption.Value
	if infrastructureEnabled {
		return fmt.Errorf("Databricks Workspace %q (Resource Group %q): `infrastructure_encryption_enabled` must be set to `false`", workspaceID.Name, workspaceID.ResourceGroup)
	}
	if !encryptionEnabled {
		return fmt.Errorf("Databricks Workspace %q (Resource Group %q): `customer_managed_key_enabled` must be set to `true`", workspaceID.Name, workspaceID.ResourceGroup)
	}

	resourceID := parse.NewCustomerManagedKeyID(subscriptionId, workspaceID.ResourceGroup, workspaceID.Name)

	if d.IsNewResource() {
		if workspace.Parameters.Encryption != nil {
			return tf.ImportAsExistsError("azurerm_databricks_workspace_customer_managed_key", resourceID.ID())
		}
	}

	keySource := databricks.MicrosoftKeyvault
	keyName := key.Name
	keyVersion := key.Version
	keyVaultBaseURI := key.KeyVaultBaseUrl

	props := databricks.Workspace{
		Location: workspace.Location,
		Sku:      workspace.Sku,
		WorkspaceProperties: &databricks.WorkspaceProperties{
			ManagedResourceGroupID: workspace.WorkspaceProperties.ManagedResourceGroupID,
			Parameters: &databricks.WorkspaceCustomParameters{
				Encryption: &databricks.WorkspaceEncryptionParameter{
					Value: &databricks.Encryption{
						KeySource:   databricks.KeySource(keySource),
						KeyName:     &keyName,
						KeyVersion:  &keyVersion,
						KeyVaultURI: &keyVaultBaseURI,
					},
				},
			},
		},
	}

	future, err := workspaceClient.CreateOrUpdate(ctx, props, resourceID.ResourceGroup, resourceID.CustomerMangagedKeyName)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", resourceID, err)
	}

	if err = future.WaitForCompletionRef(ctx, workspaceClient.Client); err != nil {
		return fmt.Errorf("waiting for create/update of %s: %+v", resourceID, err)
	}

	d.SetId(resourceID.ID())
	return DatabricksWorkspaceCustomerManagedKeyRead(d, meta)
}

func DatabricksWorkspaceCustomerManagedKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CustomerManagedKeyID(d.Id())
	if err != nil {
		return err
	}

	workspaceId := parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.CustomerMangagedKeyName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.CustomerMangagedKeyName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	keySource := ""
	keyName := ""
	keyVersion := ""
	keyVaultURI := ""

	if props := resp.WorkspaceProperties; props != nil {
		if props.Parameters.Encryption.Value.KeySource != "" {
			keySource = string(props.Parameters.Encryption.Value.KeySource)
		}
		if props.Parameters.Encryption.Value.KeyName != nil {
			keyName = *props.Parameters.Encryption.Value.KeyName
		}
		if props.Parameters.Encryption.Value.KeyVersion != nil {
			keyVersion = *props.Parameters.Encryption.Value.KeyVersion
		}
		if props.Parameters.Encryption.Value.KeyVaultURI != nil {
			keyVaultURI = *props.Parameters.Encryption.Value.KeyVaultURI
		}
	}

	if !strings.EqualFold(keySource, string(databricks.MicrosoftKeyvault)) {
		return fmt.Errorf("retrieving Databricks Workspace %q (Resource Group %q): `Workspace.WorkspaceProperties.Encryption.Value.KeySource` was expected to be %q, got %q", id.CustomerMangagedKeyName, id.ResourceGroup, string(databricks.MicrosoftKeyvault), keySource)
	}

	if keyName == "" || keyVersion == "" || keyVaultURI == "" {
		return fmt.Errorf("Databricks Workspace %q (Resource Group %q): `Workspace.WorkspaceProperties.Encryption.Value(s)` was nil", id.CustomerMangagedKeyName, id.ResourceGroup)
	}

	key, err := keyVaultParse.NewNestedItemID(keyVaultURI, "keys", keyName, keyVersion)

	d.Set("workspace_id", workspaceId.ID())
	d.Set("key_vault_key_id", key.ID())

	return nil
}

func DatabricksWorkspaceCustomerManagedKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CustomerManagedKeyID(d.Id())
	if err != nil {
		return err
	}

	workspaceID := parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.CustomerMangagedKeyName)

	// Not sure if I should also lock the key vault here too
	locks.ByName(workspaceID.Name, "azurerm_databricks_workspace")
	defer locks.UnlockByName(workspaceID.Name, "azurerm_databricks_workspace")

	resp, err := client.Get(ctx, id.ResourceGroup, id.CustomerMangagedKeyName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	// Since this isn't real and you cannot turn off CMK without destroying the
	// workspace and recreating it the best I can do is to set the workspace
	// back to using Microsoft managed keys and removing the CMK fields
	props := databricks.Workspace{
		Location: resp.Location,
		Sku:      resp.Sku,
		WorkspaceProperties: &databricks.WorkspaceProperties{
			ManagedResourceGroupID: resp.WorkspaceProperties.ManagedResourceGroupID,
			Parameters: &databricks.WorkspaceCustomParameters{
				Encryption: &databricks.WorkspaceEncryptionParameter{
					Value: &databricks.Encryption{
						KeySource: databricks.Default,
					},
				},
			},
		},
	}

	future, err := client.CreateOrUpdate(ctx, props, workspaceID.ResourceGroup, workspaceID.Name)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", workspaceID, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for create/update of %s: %+v", workspaceID, err)
	}

	return nil
}
