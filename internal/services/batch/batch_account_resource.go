package batch

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2022-01-01/batch"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/validate"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

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

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: storageValidate.StorageAccountID,
				RequiredWith: []string{"storage_account_authentication_mode"},
			},

			"storage_account_authentication_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(batch.AutoStorageAuthenticationModeStorageKeys),
					string(batch.AutoStorageAuthenticationModeBatchAccountManagedIdentity),
				}, false),
				RequiredWith: []string{"storage_account_id"},
			},

			"storage_account_node_identity": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: commonids.ValidateUserAssignedIdentityID,
				RequiredWith: []string{"storage_account_id"},
			},

			"allowed_authentication_modes": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(batch.AuthenticationModeSharedKey),
						string(batch.AuthenticationModeAAD),
						string(batch.AuthenticationModeTaskAuthenticationToken),
					}, false),
				},
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

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

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
			"encryption": {
				Type:       pluginsdk.TypeList,
				Optional:   true,
				MaxItems:   1,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_vault_key_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: keyVaultValidate.NestedItemId,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceBatchAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.AccountClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Batch account creation.")

	id := parse.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	location := azure.NormalizeLocation(d.Get("location").(string))
	storageAccountId := d.Get("storage_account_id").(string)
	poolAllocationMode := d.Get("pool_allocation_mode").(string)
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.BatchAccountName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_batch_account", id.ID())
		}
	}

	identity, err := expandBatchAccountIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf(`expanding "identity": %v`, err)
	}

	encryptionRaw := d.Get("encryption").([]interface{})
	encryption := expandEncryption(encryptionRaw)

	parameters := batch.AccountCreateParameters{
		Location: &location,
		AccountCreateProperties: &batch.AccountCreateProperties{
			PoolAllocationMode:         batch.PoolAllocationMode(poolAllocationMode),
			PublicNetworkAccess:        batch.PublicNetworkAccessTypeEnabled,
			Encryption:                 encryption,
			AllowedAuthenticationModes: expandAllowedAuthenticationModes(d.Get("allowed_authentication_modes").(*pluginsdk.Set).List()),
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
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		if keyVaultReference == nil {
			return fmt.Errorf("creating %s: When setting pool allocation mode to UserSubscription, a Key Vault reference needs to be set", id)
		}

		parameters.KeyVaultReference = keyVaultReference

		if v, ok := d.GetOk("allowed_authentication_modes"); ok {
			authModes := v.(*pluginsdk.Set).List()
			for _, mode := range authModes {
				if batch.AuthenticationMode(mode.(string)) == batch.AuthenticationModeSharedKey {
					return fmt.Errorf("creating %s: When setting pool allocation mode to UserSubscription, `allowed_authentication_modes=[StorageKeys]` is not allowed. ", id)
				}
			}
		}
	}

	authMode := d.Get("storage_account_authentication_mode").(string)
	if batch.AutoStorageAuthenticationMode(authMode) == batch.AutoStorageAuthenticationModeBatchAccountManagedIdentity &&
		identity.Type == batch.ResourceIdentityTypeNone {
		return fmt.Errorf(" storage_account_authentication_mode=`BatchAccountManagedIdentity` can only be set when identity.type is `SystemAssigned` or `UserAssigned`")
	}

	if storageAccountId != "" {
		if authMode == "" {
			return fmt.Errorf("`storage_account_authentication_mode` is required when `storage_account_id` ")
		}
		parameters.AccountCreateProperties.AutoStorage = &batch.AutoStorageBaseProperties{
			StorageAccountID:   &storageAccountId,
			AuthenticationMode: batch.AutoStorageAuthenticationMode(authMode),
		}
	}

	nodeIdentity := d.Get("storage_account_node_identity").(string)
	if nodeIdentity != "" {
		parameters.AccountCreateProperties.AutoStorage.NodeIdentityReference = &batch.ComputeNodeIdentityReference{
			ResourceID: utils.String(nodeIdentity),
		}
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.BatchAccountName, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

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

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	identity, err := flattenBatchAccountIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if props := resp.AccountProperties; props != nil {
		d.Set("account_endpoint", props.AccountEndpoint)
		if autoStorage := props.AutoStorage; autoStorage != nil {
			d.Set("storage_account_id", *autoStorage.StorageAccountID)
			d.Set("storage_account_authentication_mode", autoStorage.AuthenticationMode)

			if autoStorage.NodeIdentityReference != nil {
				d.Set("storage_account_node_identity", autoStorage.NodeIdentityReference.ResourceID)
			}
		} else {
			d.Set("storage_account_authentication_mode", "")
			d.Set("storage_account_id", "")
		}

		if props.PublicNetworkAccess != "" {
			d.Set("public_network_access_enabled", props.PublicNetworkAccess == batch.PublicNetworkAccessTypeEnabled)
		}

		d.Set("pool_allocation_mode", props.PoolAllocationMode)

		if err := d.Set("encryption", flattenEncryption(props.Encryption)); err != nil {
			return fmt.Errorf("setting `encryption`: %+v", err)
		}

		if err := d.Set("allowed_authentication_modes", flattenAllowedAuthenticationModes(props.AllowedAuthenticationModes)); err != nil {
			return fmt.Errorf("setting `allowed_authentication_modes`: %+v", err)
		}
	}

	if d.Get("pool_allocation_mode").(string) == string(batch.PoolAllocationModeBatchService) &&
		isShardKeyAllowed(d.Get("allowed_authentication_modes").(*pluginsdk.Set).List()) {
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

	t := d.Get("tags").(map[string]interface{})

	identity, err := expandBatchAccountIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf(`expanding "identity": %v`, err)
	}

	encryptionRaw := d.Get("encryption").([]interface{})
	encryption := expandEncryption(encryptionRaw)

	parameters := batch.AccountUpdateParameters{
		AccountUpdateProperties: &batch.AccountUpdateProperties{
			Encryption: encryption,
		},
		Identity: identity,
		Tags:     tags.Expand(t),
	}

	if d.HasChange("allowed_authentication_modes") {
		allowedAuthModes := d.Get("allowed_authentication_modes").(*pluginsdk.Set).List()
		if len(allowedAuthModes) == 0 {
			parameters.AllowedAuthenticationModes = &[]batch.AuthenticationMode{} // remove all modes need explicit set it to empty array not nil
		} else {
			parameters.AllowedAuthenticationModes = expandAllowedAuthenticationModes(d.Get("allowed_authentication_modes").(*pluginsdk.Set).List())
		}

	}

	if d.HasChange("storage_account_id") {
		if v, ok := d.GetOk("storage_account_id"); ok {
			parameters.AccountUpdateProperties.AutoStorage = &batch.AutoStorageBaseProperties{
				StorageAccountID: utils.String(v.(string)),
			}
		} else {
			// remove the storage account from the batch account
			parameters.AccountUpdateProperties.AutoStorage = &batch.AutoStorageBaseProperties{
				StorageAccountID: nil,
			}
		}
	}

	authMode := d.Get("storage_account_authentication_mode").(string)
	if batch.AutoStorageAuthenticationMode(authMode) == batch.AutoStorageAuthenticationModeBatchAccountManagedIdentity &&
		identity.Type == batch.ResourceIdentityTypeNone {
		return fmt.Errorf(" storage_account_authentication_mode=`BatchAccountManagedIdentity` can only be set when identity.type is `SystemAssigned` or `UserAssigned`")
	}

	storageAccountId := d.Get("storage_account_id").(string)
	if storageAccountId != "" {
		parameters.AutoStorage = &batch.AutoStorageBaseProperties{
			StorageAccountID:   &storageAccountId,
			AuthenticationMode: batch.AutoStorageAuthenticationMode(authMode),
		}
	}

	nodeIdentity := d.Get("storage_account_node_identity").(string)
	if nodeIdentity != "" {
		parameters.AutoStorage.NodeIdentityReference = &batch.ComputeNodeIdentityReference{
			ResourceID: utils.String(nodeIdentity),
		}
	}

	if _, err = client.Update(ctx, id.ResourceGroup, id.BatchAccountName, parameters); err != nil {
		return fmt.Errorf("updating Batch account %q (Resource Group %q): %+v", id.BatchAccountName, id.ResourceGroup, err)
	}

	d.SetId(id.ID())

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
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := batch.AccountIdentity{
		Type: batch.ResourceIdentityType(string(expanded.Type)),
	}
	if expanded.Type == identity.TypeUserAssigned || expanded.Type == identity.TypeSystemAssignedUserAssigned {
		out.UserAssignedIdentities = make(map[string]*batch.UserAssignedIdentities)
		for k := range expanded.IdentityIds {
			out.UserAssignedIdentities[k] = &batch.UserAssignedIdentities{
				// intentionally empty
			}
		}
	}
	return &out, nil
}

func flattenBatchAccountIdentity(input *batch.AccountIdentity) (*[]interface{}, error) {
	var transform *identity.SystemAndUserAssignedMap

	if input != nil {
		transform = &identity.SystemAndUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		if input.PrincipalID != nil {
			transform.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transform.TenantId = *input.TenantID
		}
		for k, v := range input.UserAssignedIdentities {
			transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				ClientId:    v.ClientID,
				PrincipalId: v.PrincipalID,
			}
		}
	}

	return identity.FlattenSystemAndUserAssignedMap(transform)
}

func expandEncryption(e []interface{}) *batch.EncryptionProperties {
	defaultEnc := batch.EncryptionProperties{
		KeySource: batch.KeySourceMicrosoftBatch,
	}

	if len(e) == 0 || e[0] == nil {
		return &defaultEnc
	}

	v := e[0].(map[string]interface{})
	keyId := v["key_vault_key_id"].(string)
	encryptionProperty := batch.EncryptionProperties{
		KeySource: batch.KeySourceMicrosoftKeyVault,
		KeyVaultProperties: &batch.KeyVaultProperties{
			KeyIdentifier: &keyId,
		},
	}

	return &encryptionProperty
}

func expandAllowedAuthenticationModes(input []interface{}) *[]batch.AuthenticationMode {
	if len(input) == 0 {
		return nil
	}

	allowedAuthModes := make([]batch.AuthenticationMode, 0)
	for _, mode := range input {
		allowedAuthModes = append(allowedAuthModes, batch.AuthenticationMode(mode.(string)))
	}
	return &allowedAuthModes
}

func flattenAllowedAuthenticationModes(input *[]batch.AuthenticationMode) []string {
	if input == nil || len(*input) == 0 {
		return []string{}
	}

	allowedAuthModes := make([]string, 0)
	for _, mode := range *input {
		allowedAuthModes = append(allowedAuthModes, string(mode))
	}
	return allowedAuthModes
}

func flattenEncryption(encryptionProperties *batch.EncryptionProperties) []interface{} {
	if encryptionProperties == nil || encryptionProperties.KeySource == batch.KeySourceMicrosoftBatch {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"key_vault_key_id": *encryptionProperties.KeyVaultProperties.KeyIdentifier,
		},
	}
}

func isShardKeyAllowed(input []interface{}) bool {
	if len(input) == 0 {
		return false
	}
	for _, authMod := range input {
		if strings.EqualFold(authMod.(string), string(batch.AuthenticationModeSharedKey)) {
			return true
		}
	}
	return false
}
