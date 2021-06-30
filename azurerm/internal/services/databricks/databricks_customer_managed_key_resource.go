package databricks

import (
	"fmt"
	"log"
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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
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

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.WorkspaceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			// Make this key vault key id and abstract everything from the string...
			"key_vault_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: keyVaultValidate.VaultID,
			},

			"key_source": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(databricks.Default),
					string(databricks.MicrosoftKeyvault),
				}, true),
			},

			"key_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"key_version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
		},
	}
}

func DatabricksWorkspaceCustomerManagedKeyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	workspaceClient := meta.(*clients.Client).DataBricks.WorkspacesClient
	keyVaultClient := meta.(*clients.Client).KeyVault
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceIDRaw := d.Get("workspace_id").(string)
	workspaceID, err := parse.WorkspaceID(workspaceIDRaw)
	if err != nil {
		return err
	}

	// Not sure if I should also lock the key vault here too
	locks.ByName(workspaceID.Name, "azurerm_databricks_workspace")
	defer locks.UnlockByName(workspaceID.Name, "azurerm_databricks_workspace")

	workspace, err := workspaceClient.Get(ctx, workspaceID.ResourceGroup, workspaceID.Name)
	if err != nil {
		return fmt.Errorf("retrieving Databricks Workspace %q (Resource Group %q): %+v", workspaceID.Name, workspaceID.ResourceGroup, err)
	}
	if workspace.Parameters == nil {
		return fmt.Errorf("retrieving Databricks Workspace %q (Resource Group %q): `WorkspaceCustomParameters` was nil", workspaceID.Name, workspaceID.ResourceGroup)
	}

	resourceID := parse.NewCustomerManagedKeyID(subscriptionId, workspaceID.ResourceGroup, workspaceID.Name)

	if d.IsNewResource() {
		if workspace.Parameters.Encryption != nil {
			return tf.ImportAsExistsError("azurerm_databricks_workspace_customer_managed_key", resourceID.ID())
		}
	}

	keyVaultIDRaw := d.Get("key_vault_id").(string)
	keyVaultID, err := keyVaultParse.VaultID(keyVaultIDRaw)
	if err != nil {
		return err
	}

	keyVaultBaseURL, err := keyVaultClient.BaseUriForKeyVault(ctx, *keyVaultID)
	if err != nil {
		return fmt.Errorf("looking up Key Vault URI from Key Vault %q (Resource Group %q): %+v", keyVaultID.Name, keyVaultID.ResourceGroup, err)
	}

	keySource := d.Get("key_source").(string)
	keyName := d.Get("key_name").(string)
	keyVersion := d.Get("key_version").(string)

	props := databricks.Workspace{
		WorkspaceProperties: &databricks.WorkspaceProperties{
			Parameters: &databricks.WorkspaceCustomParameters{
				Encryption: &databricks.WorkspaceEncryptionParameter{
					Value: &databricks.Encryption{
						KeySource:   databricks.KeySource(keySource),
						KeyName:     &keyName,
						KeyVersion:  &keyVersion,
						KeyVaultURI: keyVaultBaseURL,
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
	keyVaultsClient := meta.(*clients.Client).KeyVault
	resourcesClient := meta.(*clients.Client).Resource
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

	if keyVaultURI == "" {
		return fmt.Errorf("Error retrieving Databricks Workspace %q (Resource Group %q): `Workspace.WorkspaceProperties.Encryption.Value.KeyVaultURI` was nil", id.CustomerMangagedKeyName, id.ResourceGroup)
	}

	// now we have the key vault uri we can look up the ID
	keyVaultID, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, keyVaultURI)
	if err != nil {
		return fmt.Errorf("Error retrieving Key Vault ID from the Base URI %q: %+v", keyVaultURI, err)
	}

	d.Set("workspace_id", workspaceId.ID())
	d.Set("key_vault_id", keyVaultID)
	d.Set("key_source", keySource)
	d.Set("key_name", keyName)
	d.Set("key_version", keyVersion)
	d.Set("key_vault_uri", keyVaultURI)

	return nil
}

func DatabricksWorkspaceCustomerManagedKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	resourcesClient := meta.(*clients.Client).Resource
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

	if keyVaultURI == "" {
		return fmt.Errorf("Error retrieving Databricks Workspace %q (Resource Group %q): `Workspace.WorkspaceProperties.Encryption.Value.KeyVaultURI` was nil", id.CustomerMangagedKeyName, id.ResourceGroup)
	}

	// now we have the key vault uri we can look up the ID
	keyVaultID, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, keyVaultURI)
	if err != nil {
		return fmt.Errorf("Error retrieving Key Vault ID from the Base URI %q: %+v", keyVaultURI, err)
	}

	d.Set("workspace_id", workspaceId.ID())
	d.Set("key_vault_id", keyVaultID)
	d.Set("key_source", keySource)
	d.Set("key_name", keyName)
	d.Set("key_version", keyVersion)
	d.Set("key_vault_uri", keyVaultURI)

	return nil
}

// func DatabricksWorkspaceCustomerManagedKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
// 	client := meta.(*clients.Client).DataBricks.WorkspacesClient
// 	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
// 	defer cancel()

// 	id, err := parse.WorkspaceID(d.Id())
// 	if err != nil {
// 		return err
// 	}

// 	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
// 	if err != nil {
// 		return fmt.Errorf("deleting %s: %+v", *id, err)
// 	}

// 	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
// 		if !response.WasNotFound(future.Response()) {
// 			return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
// 		}
// 	}

// 	return nil

// client := meta.(*clients.Client).Kusto.ClustersClient
// ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
// defer cancel()

// clusterID, err := parse.ClusterID(d.Id())
// if err != nil {
// 	return err
// }

// locks.ByName(clusterID.Name, "azurerm_kusto_cluster")
// defer locks.UnlockByName(clusterID.Name, "azurerm_kusto_cluster")

// // confirm it still exists prior to trying to update it, else we'll get an error
// cluster, err := client.Get(ctx, clusterID.ResourceGroup, clusterID.Name)
// if err != nil {
// 	if utils.ResponseWasNotFound(cluster.Response) {
// 		return nil
// 	}

// 	return fmt.Errorf("Error retrieving Kusto Cluster %q (Resource Group %q): %+v", clusterID.Name, clusterID.ResourceGroup, err)
// }

// // Since this isn't a real object, just modifying an existing object
// // "Delete" doesn't really make sense it should really be a "Revert to Default"
// // So instead of the Delete func actually deleting the Kusto Cluster I am
// // making it reset the Kusto cluster to its default state
// props := kusto.ClusterUpdate{
// 	ClusterProperties: &kusto.ClusterProperties{
// 		KeyVaultProperties: &kusto.KeyVaultProperties{},
// 	},
// }

// future, err := client.Update(ctx, clusterID.ResourceGroup, clusterID.Name, props)
// if err != nil {
// 	return fmt.Errorf("Error removing Customer Managed Key for Kusto Cluster %q (Resource Group %q): %+v", clusterID.Name, clusterID.ResourceGroup, err)
// }
// if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
// 	return fmt.Errorf("Error waiting for completion of Kusto Cluster Update %q (Resource Group %q): %+v", clusterID.Name, clusterID.ResourceGroup, err)
// }

// return nil
// }
