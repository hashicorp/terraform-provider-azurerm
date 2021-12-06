package batch

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2021-06-01/batch"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/validate"
	msiparse "github.com/hashicorp/terraform-provider-azurerm/internal/services/msi/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type batchAccountIdentity = identity.SystemAssignedUserAssigned

func resourceBatchAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBatchAccountCreate,
		Read:   resourceBatchAccountRead,
		Update: resourceBatchAccountUpdate,
		Delete: resourceBatchAccountDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AccountID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName,
			},

			// TODO: make this case sensitive once this API bug has been fixed:
			// https://github.com/Azure/azure-rest-api-specs/issues/5574
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"location": azure.SchemaLocation(),

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},

			"pool_allocation_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(batch.PoolAllocationModeBatchService),
				ValidateFunc: validation.StringInSlice([]string{
					string(batch.PoolAllocationModeBatchService),
					string(batch.PoolAllocationModeUserSubscription),
				}, false),
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"key_vault_reference": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"url": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsURLWithHTTPS,
						},
					},
				},
			},

			"identity": batchAccountIdentity{}.Schema(),

			"primary_access_key": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"secondary_access_key": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"account_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"tags": tags.Schema(),
		},
	}
}

func resourceBatchAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.AccountClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Batch account creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	storageAccountId := d.Get("storage_account_id").(string)
	poolAllocationMode := d.Get("pool_allocation_mode").(string)
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Batch Account %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_batch_account", *existing.ID)
		}
	}

	identity, err := expandBatchAccountIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf(`expanding "identity": %v`, err)
	}

	parameters := batch.AccountCreateParameters{
		Location: &location,
		AccountCreateProperties: &batch.AccountCreateProperties{
			PoolAllocationMode:  batch.PoolAllocationMode(poolAllocationMode),
			PublicNetworkAccess: batch.PublicNetworkAccessTypeEnabled,
		},
		Identity: identity,
		Tags:     tags.Expand(t),
	}

	if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
		parameters.AccountCreateProperties.PublicNetworkAccess = batch.PublicNetworkAccessTypeDisabled
	}

	// if pool allocation mode is UserSubscription, a key vault reference needs to be set
	if poolAllocationMode == string(batch.PoolAllocationModeUserSubscription) {
		keyVaultReferenceSet := d.Get("key_vault_reference").([]interface{})
		keyVaultReference, err := expandBatchAccountKeyVaultReference(keyVaultReferenceSet)
		if err != nil {
			return fmt.Errorf("creating Batch account %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		if keyVaultReference == nil {
			return fmt.Errorf("creating Batch account %q (Resource Group %q): When setting pool allocation mode to UserSubscription, a Key Vault reference needs to be set", name, resourceGroup)
		}

		parameters.KeyVaultReference = keyVaultReference
	}

	if storageAccountId != "" {
		parameters.AccountCreateProperties.AutoStorage = &batch.AutoStorageBaseProperties{
			StorageAccountID: &storageAccountId,
		}
	}

	future, err := client.Create(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("creating Batch account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Batch account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Batch account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Batch account %q (resource group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceBatchAccountRead(d, meta)
}

func resourceBatchAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.AccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.BatchAccountName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			log.Printf("[DEBUG] Batch Account %q was not found in Resource Group %q - removing from state!", id.BatchAccountName, id.ResourceGroup)
			return nil
		}
		return fmt.Errorf("reading the state of Batch account %q: %+v", id.BatchAccountName, err)
	}

	d.Set("name", id.BatchAccountName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_endpoint", resp.AccountEndpoint)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.AccountProperties; props != nil {
		if autoStorage := props.AutoStorage; autoStorage != nil {
			d.Set("storage_account_id", autoStorage.StorageAccountID)
		}

		if props.PublicNetworkAccess != "" {
			d.Set("public_network_access_enabled", props.PublicNetworkAccess == batch.PublicNetworkAccessTypeEnabled)
		}

		d.Set("pool_allocation_mode", props.PoolAllocationMode)

		identity, err := flattenBatchAccountIdentity(resp.Identity)
		if err != nil {
			return err
		}
		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}
	}

	if d.Get("pool_allocation_mode").(string) == string(batch.PoolAllocationModeBatchService) {
		keys, err := client.GetKeys(ctx, id.ResourceGroup, id.BatchAccountName)
		if err != nil {
			return fmt.Errorf("Cannot read keys for Batch account %q (resource group %q): %v", id.BatchAccountName, id.ResourceGroup, err)
		}

		d.Set("primary_access_key", keys.Primary)
		d.Set("secondary_access_key", keys.Secondary)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceBatchAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.AccountClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Batch account update.")

	id, err := parse.AccountID(d.Id())
	if err != nil {
		return err
	}

	storageAccountId := d.Get("storage_account_id").(string)
	t := d.Get("tags").(map[string]interface{})

	identity, err := expandBatchAccountIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf(`expanding "identity": %v`, err)
	}

	parameters := batch.AccountUpdateParameters{
		AccountUpdateProperties: &batch.AccountUpdateProperties{
			AutoStorage: &batch.AutoStorageBaseProperties{
				StorageAccountID: &storageAccountId,
			},
		},
		Identity: identity,
		Tags:     tags.Expand(t),
	}

	if _, err = client.Update(ctx, id.ResourceGroup, id.BatchAccountName, parameters); err != nil {
		return fmt.Errorf("updating Batch account %q (Resource Group %q): %+v", id.BatchAccountName, id.ResourceGroup, err)
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.BatchAccountName)
	if err != nil {
		return fmt.Errorf("retrieving Batch account %q (Resource Group %q): %+v", id.BatchAccountName, id.ResourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Batch account %q (resource group %q) ID", id.BatchAccountName, id.ResourceGroup)
	}

	d.SetId(*read.ID)

	return resourceBatchAccountRead(d, meta)
}

func resourceBatchAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.AccountClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccountID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.BatchAccountName)
	if err != nil {
		return fmt.Errorf("deleting Batch account %q (Resource Group %q): %+v", id.BatchAccountName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of Batch account %q (Resource Group %q): %+v", id.BatchAccountName, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandBatchAccountIdentity(input []interface{}) (*batch.AccountIdentity, error) {
	config, err := batchAccountIdentity{}.Expand(input)
	if err != nil {
		return nil, err
	}

	var identityMaps map[string]*batch.UserAssignedIdentities
	if len(config.UserAssignedIdentityIds) != 0 {
		identityMaps = make(map[string]*batch.UserAssignedIdentities, len(config.UserAssignedIdentityIds))
		for _, id := range config.UserAssignedIdentityIds {
			identityMaps[id] = &batch.UserAssignedIdentities{}
		}
	}

	return &batch.AccountIdentity{
		Type:                   batch.ResourceIdentityType(config.Type),
		TenantID:               &config.TenantId,
		PrincipalID:            &config.PrincipalId,
		UserAssignedIdentities: identityMaps,
	}, nil
}

func flattenBatchAccountIdentity(input *batch.AccountIdentity) ([]interface{}, error) {
	var config *identity.ExpandedConfig

	if input == nil {
		return []interface{}{}, nil
	}

	var identityIds []string
	for id := range input.UserAssignedIdentities {
		parsedId, err := msiparse.UserAssignedIdentityIDInsensitively(id)
		if err != nil {
			return nil, err
		}
		identityIds = append(identityIds, parsedId.ID())
	}

	principalId := ""
	if input.PrincipalID != nil {
		principalId = *input.PrincipalID
	}

	tenantId := ""
	if input.TenantID != nil {
		tenantId = *input.TenantID
	}

	config = &identity.ExpandedConfig{
		Type:                    identity.Type(string(input.Type)),
		PrincipalId:             principalId,
		TenantId:                tenantId,
		UserAssignedIdentityIds: identityIds,
	}
	return batchAccountIdentity{}.Flatten(config), nil
}
