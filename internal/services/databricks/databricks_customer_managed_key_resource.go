package databricks

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/sdk/2021-04-01-preview/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/validate"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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
			workspace := workspaces.NewWorkspaceID(customManagedKey.SubscriptionId, customManagedKey.ResourceGroup, customManagedKey.CustomerMangagedKeyName)

			// validate that the workspace exists
			if _, err = client.Get(ctx, workspace); err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("retrieving the Databricks workspace customer managed key configuration(ID: %q) for workspace (ID: %q): %s", customManagedKey.ID(), workspace.ID(), err)
			}

			// set the new values for the CMK resource
			d.SetId(customManagedKey.ID())
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
	keyVaultsClient := meta.(*clients.Client).KeyVault
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
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
	locks.ByName(id.Name, "azurerm_databricks_workspace")
	defer locks.UnlockByName(id.Name, "azurerm_databricks_workspace")
	var encryptionEnabled, infrastructureEnabled bool

	workspace, err := workspaceClient.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if model := workspace.Model; model != nil {
		if parameters := model.Properties.Parameters; parameters != nil {
			if parameters.RequireInfrastructureEncryption != nil {
				infrastructureEnabled = model.Properties.Parameters.RequireInfrastructureEncryption.Value
			}
			if parameters.PrepareEncryption != nil {
				encryptionEnabled = model.Properties.Parameters.PrepareEncryption.Value
			}
		}
	}

	if infrastructureEnabled {
		return fmt.Errorf("%s: `infrastructure_encryption_enabled` must be set to `false`", *id)
	}
	if !encryptionEnabled {
		return fmt.Errorf("%s: `customer_managed_key_enabled` must be set to `true`", *id)
	}

	// make sure the key vault exists
	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, meta.(*clients.Client).Resource, key.KeyVaultBaseUrl)
	if err != nil || keyVaultIdRaw == nil {
		return fmt.Errorf("retrieving the Resource ID for the Key Vault at URL %q: %+v", key.KeyVaultBaseUrl, err)
	}

	resourceID := parse.NewCustomerManagedKeyID(subscriptionId, id.ResourceGroup, id.Name)

	if d.IsNewResource() {
		if workspace.Model != nil && workspace.Model.Properties.Parameters != nil && workspace.Model.Properties.Parameters.Encryption != nil {
			return tf.ImportAsExistsError("azurerm_databricks_workspace_customer_managed_key", resourceID.ID())
		}
	}

	// We need to pull all of the custom params from the parent
	// workspace resource and then add our new encryption values into the
	// structure, else the other values set in the parent workspace
	// resource will be lost and overwritten as nil. ¯\_(ツ)_/¯
	// NOTE: 'workspace.Parameters' will never be nil as 'customer_managed_key_enabled' and 'infrastructure_encryption_enabled'
	// fields have a default value in the parent workspace resource.
	keySource := workspaces.KeySourceMicrosoftPointKeyvault
	params := workspace.Model.Properties.Parameters
	params.Encryption = &workspaces.WorkspaceEncryptionParameter{
		Value: &workspaces.Encryption{
			KeySource:   &keySource,
			KeyName:     &key.Name,
			Keyversion:  &key.Version,
			Keyvaulturi: &key.KeyVaultBaseUrl,
		},
	}

	props := getProps(workspace.Model, params)

	if err = workspaceClient.CreateOrUpdateThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", resourceID, err)
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

	workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.CustomerMangagedKeyName)

	resp, err := client.Get(ctx, workspaceId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
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

	if model := resp.Model; model != nil {
		if model.Properties.Parameters != nil {
			if props := model.Properties.Parameters.Encryption; props != nil {
				if props.Value.KeySource != nil {
					keySource = string(*props.Value.KeySource)
				}
				if props.Value.KeyName != nil {
					keyName = *props.Value.KeyName
				}
				if props.Value.Keyversion != nil {
					keyVersion = *props.Value.Keyversion
				}
				if props.Value.Keyvaulturi != nil {
					keyVaultURI = *props.Value.Keyvaulturi
				}
			}
		}
	}

	if strings.EqualFold(keySource, string(workspaces.KeySourceMicrosoftPointKeyvault)) && (keyName == "" || keyVersion == "" || keyVaultURI == "") {
		return fmt.Errorf("%s: `Workspace.WorkspaceProperties.Parameters.Encryption.Value(s)` were nil", *id)
	}

	d.SetId(id.ID())
	d.Set("workspace_id", workspaceId.ID())

	if keyVaultURI != "" {
		key, err := keyVaultParse.NewNestedItemID(keyVaultURI, "keys", keyName, keyVersion)
		if err == nil {
			d.Set("key_vault_key_id", key.ID())
		}
	}

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

	workspaceID := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.CustomerMangagedKeyName)

	// Not sure if I should also lock the key vault here too
	locks.ByName(workspaceID.Name, "azurerm_databricks_workspace")
	defer locks.UnlockByName(workspaceID.Name, "azurerm_databricks_workspace")

	workspace, err := client.Get(ctx, workspaceID)
	if err != nil {
		if response.WasNotFound(workspace.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	// Since this isn't real and you cannot turn off CMK without destroying the
	// workspace and recreating it the best I can do is to set the workspace
	// back to using Microsoft managed keys and removing the CMK fields
	// also need to pull all of the custom params from the parent
	// workspace resource and then add our new encryption values into the
	// structure, else the other values set in the parent workspace
	// resource will be lost and overwritten as nil. ¯\_(ツ)_/¯
	keySource := workspaces.KeySourceDefault
	params := workspace.Model.Properties.Parameters
	params.Encryption = &workspaces.WorkspaceEncryptionParameter{
		Value: &workspaces.Encryption{
			KeySource: &keySource,
		},
	}

	props := getProps(workspace.Model, params)

	if err = client.CreateOrUpdateThenPoll(ctx, workspaceID, props); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", workspaceID, err)
	}

	return nil
}

func getProps(workspace *workspaces.Workspace, params *workspaces.WorkspaceCustomParameters) workspaces.Workspace {
	props := workspaces.Workspace{
		Location: workspace.Location,
		Sku:      workspace.Sku,
		Properties: workspaces.WorkspaceProperties{
			PublicNetworkAccess:    workspace.Properties.PublicNetworkAccess,
			ManagedResourceGroupId: workspace.Properties.ManagedResourceGroupId,
			Parameters:             params,
		},
		Tags: workspace.Tags,
	}

	// If notebook encryption exists add it to the properties
	if workspace.Properties.Encryption != nil {
		props.Properties.Encryption = workspace.Properties.Encryption
	}

	// This is only valid if Private Link only is set
	if workspace.Properties.PublicNetworkAccess != nil && *workspace.Properties.PublicNetworkAccess == workspaces.PublicNetworkAccessDisabled {
		props.Properties.RequiredNsgRules = workspace.Properties.RequiredNsgRules
	}

	return props
}
