// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package batch

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2023-05-01/batchaccount"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/validate"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
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
			_, err := batchaccount.ParseBatchAccountID(id)
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
				ValidateFunc: commonids.ValidateStorageAccountID,
				RequiredWith: []string{"storage_account_authentication_mode"},
			},

			"storage_account_authentication_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(batchaccount.AutoStorageAuthenticationModeStorageKeys),
					string(batchaccount.AutoStorageAuthenticationModeBatchAccountManagedIdentity),
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
						string(batchaccount.AuthenticationModeSharedKey),
						string(batchaccount.AuthenticationModeAAD),
						string(batchaccount.AuthenticationModeTaskAuthenticationToken),
					}, false),
				},
			},

			"pool_allocation_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(batchaccount.PoolAllocationModeBatchService),
				ValidateFunc: validation.StringInSlice([]string{
					string(batchaccount.PoolAllocationModeBatchService),
					string(batchaccount.PoolAllocationModeUserSubscription),
				}, false),
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"network_profile": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"account_access": resourceBatchAccountEndpointAccessProfileSchema(),

						"node_management_access": resourceBatchAccountEndpointAccessProfileSchema(),
					},
				},
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
							ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceBatchAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.AccountClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Batch account creation.")

	id := batchaccount.NewBatchAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	location := azure.NormalizeLocation(d.Get("location").(string))
	storageAccountId := d.Get("storage_account_id").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_batch_account", id.ID())
		}
	}

	identity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf(`expanding "identity": %v`, err)
	}

	encryptionRaw := d.Get("encryption").([]interface{})
	encryption := expandEncryption(encryptionRaw)

	poolAllocationMode := batchaccount.PoolAllocationMode(d.Get("pool_allocation_mode").(string))
	parameters := batchaccount.BatchAccountCreateParameters{
		Location: location,
		Properties: &batchaccount.BatchAccountCreateProperties{
			PoolAllocationMode:         &poolAllocationMode,
			PublicNetworkAccess:        utils.ToPtr(batchaccount.PublicNetworkAccessTypeEnabled),
			Encryption:                 encryption,
			AllowedAuthenticationModes: expandAllowedAuthenticationModes(d.Get("allowed_authentication_modes").(*pluginsdk.Set).List()),
		},
		Identity: identity,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
		parameters.Properties.PublicNetworkAccess = utils.ToPtr(batchaccount.PublicNetworkAccessTypeDisabled)
	}

	if v, ok := d.GetOk("network_profile"); ok {
		parameters.Properties.NetworkProfile = expandBatchAccountNetworkProfile(v.([]interface{}))
	}

	// if pool allocation mode is UserSubscription, a key vault reference needs to be set
	if poolAllocationMode == batchaccount.PoolAllocationModeUserSubscription {
		keyVaultReferenceSet := d.Get("key_vault_reference").([]interface{})
		keyVaultReference, err := expandBatchAccountKeyVaultReference(keyVaultReferenceSet)
		if err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		if keyVaultReference == nil {
			return fmt.Errorf("creating %s: When setting pool allocation mode to UserSubscription, a Key Vault reference needs to be set", id)
		}

		parameters.Properties.KeyVaultReference = keyVaultReference

		if v, ok := d.GetOk("allowed_authentication_modes"); ok {
			authModes := v.(*pluginsdk.Set).List()
			for _, mode := range authModes {
				if batchaccount.AuthenticationMode(mode.(string)) == batchaccount.AuthenticationModeSharedKey {
					return fmt.Errorf("creating %s: When setting pool allocation mode to UserSubscription, `allowed_authentication_modes=[StorageKeys]` is not allowed. ", id)
				}
			}
		}
	}

	authMode := d.Get("storage_account_authentication_mode").(string)
	if batchaccount.AutoStorageAuthenticationMode(authMode) == batchaccount.AutoStorageAuthenticationModeBatchAccountManagedIdentity && identity.Type == "None" {
		return fmt.Errorf(" storage_account_authentication_mode=`BatchAccountManagedIdentity` can only be set when identity.type is `SystemAssigned` or `UserAssigned`")
	}

	if storageAccountId != "" {
		if authMode == "" {
			return fmt.Errorf("`storage_account_authentication_mode` is required when `storage_account_id` ")
		}
		parameters.Properties.AutoStorage = &batchaccount.AutoStorageBaseProperties{
			StorageAccountId:   &storageAccountId,
			AuthenticationMode: utils.ToPtr(batchaccount.AutoStorageAuthenticationMode(authMode)),
		}
	}

	nodeIdentity := d.Get("storage_account_node_identity").(string)
	if nodeIdentity != "" {
		parameters.Properties.AutoStorage.NodeIdentityReference = &batchaccount.ComputeNodeIdentityReference{
			ResourceId: utils.String(nodeIdentity),
		}
	}

	if err := client.CreateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceBatchAccountRead(d, meta)
}

func resourceBatchAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.AccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := batchaccount.ParseBatchAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			log.Printf("[DEBUG] Batch Account %s - removing from state!", *id)
			return nil
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.BatchAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if location := model.Location; location != nil {
			d.Set("location", azure.NormalizeLocation(*location))
		}

		identity, err := identity.FlattenSystemOrUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}

		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {

			d.Set("account_endpoint", props.AccountEndpoint)
			if autoStorage := props.AutoStorage; autoStorage != nil {
				d.Set("storage_account_id", autoStorage.StorageAccountId)
				d.Set("storage_account_authentication_mode", string(pointer.From(autoStorage.AuthenticationMode)))

				if autoStorage.NodeIdentityReference != nil {
					d.Set("storage_account_node_identity", autoStorage.NodeIdentityReference.ResourceId)
				}
			} else {
				d.Set("storage_account_authentication_mode", "")
				d.Set("storage_account_id", "")
			}

			if v := props.PublicNetworkAccess; v != nil {
				d.Set("public_network_access_enabled", *v == batchaccount.PublicNetworkAccessTypeEnabled)
			}

			if err := d.Set("network_profile", flattenBatchAccountNetworkProfile(props.NetworkProfile)); err != nil {
				return fmt.Errorf("setting `network_profile`: %+v", err)
			}

			d.Set("pool_allocation_mode", string(pointer.From(props.PoolAllocationMode)))

			if err := d.Set("encryption", flattenEncryption(props.Encryption)); err != nil {
				return fmt.Errorf("setting `encryption`: %+v", err)
			}

			if err := d.Set("allowed_authentication_modes", flattenAllowedAuthenticationModes(props.AllowedAuthenticationModes)); err != nil {
				return fmt.Errorf("setting `allowed_authentication_modes`: %+v", err)
			}

			if d.Get("pool_allocation_mode").(string) == string(batchaccount.PoolAllocationModeBatchService) &&
				isShardKeyAllowed(d.Get("allowed_authentication_modes").(*pluginsdk.Set).List()) {
				keys, err := client.GetKeys(ctx, *id)
				if err != nil {
					return fmt.Errorf("cannot read keys for Batch account %s: %v", *id, err)
				}

				if keysModel := keys.Model; keysModel != nil {
					d.Set("primary_access_key", keysModel.Primary)
					d.Set("secondary_access_key", keysModel.Secondary)
				}
			}
			return tags.FlattenAndSet(d, model.Tags)
		}
	}
	return nil
}

func resourceBatchAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.AccountClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Batch account update.")

	id, err := batchaccount.ParseBatchAccountID(d.Id())
	if err != nil {
		return err
	}

	t := d.Get("tags").(map[string]interface{})

	identity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf(`expanding "identity": %v`, err)
	}

	encryptionRaw := d.Get("encryption").([]interface{})
	encryption := expandEncryption(encryptionRaw)

	parameters := batchaccount.BatchAccountUpdateParameters{
		Properties: &batchaccount.BatchAccountUpdateProperties{
			Encryption: encryption,
		},
		Identity: identity,
		Tags:     tags.Expand(t),
	}

	if d.HasChange("allowed_authentication_modes") {
		allowedAuthModes := d.Get("allowed_authentication_modes").(*pluginsdk.Set).List()
		if len(allowedAuthModes) == 0 {
			parameters.Properties.AllowedAuthenticationModes = &[]batchaccount.AuthenticationMode{} // remove all modes need explicit set it to empty array not nil
		} else {
			parameters.Properties.AllowedAuthenticationModes = expandAllowedAuthenticationModes(d.Get("allowed_authentication_modes").(*pluginsdk.Set).List())
		}
	}

	if d.HasChange("public_network_access_enabled") {
		if d.Get("public_network_access_enabled").(bool) {
			parameters.Properties.PublicNetworkAccess = utils.ToPtr(batchaccount.PublicNetworkAccessTypeEnabled)
		} else {
			parameters.Properties.PublicNetworkAccess = utils.ToPtr(batchaccount.PublicNetworkAccessTypeDisabled)
		}
	}

	if d.HasChange("network_profile") {
		parameters.Properties.NetworkProfile = expandBatchAccountNetworkProfile(d.Get("network_profile").([]interface{}))
	}

	if d.HasChange("storage_account_id") {
		if v, ok := d.GetOk("storage_account_id"); ok {
			parameters.Properties.AutoStorage = &batchaccount.AutoStorageBaseProperties{
				StorageAccountId: utils.String(v.(string)),
			}
		} else {
			parameters.Properties.AutoStorage = &batchaccount.AutoStorageBaseProperties{
				StorageAccountId: nil,
			}
		}
	}

	authMode := d.Get("storage_account_authentication_mode").(string)
	if batchaccount.AutoStorageAuthenticationMode(authMode) == batchaccount.AutoStorageAuthenticationModeBatchAccountManagedIdentity && identity.Type == "None" {
		return fmt.Errorf(" storage_account_authentication_mode=`BatchAccountManagedIdentity` can only be set when identity.type is `SystemAssigned` or `UserAssigned`")
	}

	storageAccountId := d.Get("storage_account_id").(string)
	if storageAccountId != "" {
		parameters.Properties.AutoStorage = &batchaccount.AutoStorageBaseProperties{
			StorageAccountId:   &storageAccountId,
			AuthenticationMode: utils.ToPtr(batchaccount.AutoStorageAuthenticationMode(authMode)),
		}
	}

	nodeIdentity := d.Get("storage_account_node_identity").(string)
	if nodeIdentity != "" {
		parameters.Properties.AutoStorage.NodeIdentityReference = &batchaccount.ComputeNodeIdentityReference{
			ResourceId: utils.String(nodeIdentity),
		}
	}

	if _, err = client.Update(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	d.SetId(id.ID())

	return resourceBatchAccountRead(d, meta)
}

func resourceBatchAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.AccountClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := batchaccount.ParseBatchAccountID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandEncryption(e []interface{}) *batchaccount.EncryptionProperties {
	defaultEnc := batchaccount.EncryptionProperties{
		KeySource: utils.ToPtr(batchaccount.KeySourceMicrosoftPointBatch),
	}

	if len(e) == 0 || e[0] == nil {
		return &defaultEnc
	}

	v := e[0].(map[string]interface{})
	keyId := v["key_vault_key_id"].(string)
	encryptionProperty := batchaccount.EncryptionProperties{
		KeySource: utils.ToPtr(batchaccount.KeySourceMicrosoftPointKeyVault),
		KeyVaultProperties: &batchaccount.KeyVaultProperties{
			KeyIdentifier: &keyId,
		},
	}

	return &encryptionProperty
}

func expandAllowedAuthenticationModes(input []interface{}) *[]batchaccount.AuthenticationMode {
	if len(input) == 0 {
		return nil
	}

	allowedAuthModes := make([]batchaccount.AuthenticationMode, 0)
	for _, mode := range input {
		allowedAuthModes = append(allowedAuthModes, batchaccount.AuthenticationMode(mode.(string)))
	}
	return &allowedAuthModes
}

func expandBatchAccountNetworkProfile(input []interface{}) *batchaccount.NetworkProfile {
	if len(input) == 0 || input[0] == nil {
		return &batchaccount.NetworkProfile{}
	}

	networkProfile := input[0].(map[string]interface{})
	return &batchaccount.NetworkProfile{
		AccountAccess:        expandBatchAccountEndpointAccessProfile(networkProfile["account_access"].([]interface{})),
		NodeManagementAccess: expandBatchAccountEndpointAccessProfile(networkProfile["node_management_access"].([]interface{})),
	}
}

func expandBatchAccountEndpointAccessProfile(input []interface{}) *batchaccount.EndpointAccessProfile {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	accessProfile := input[0].(map[string]interface{})

	ipRulesRaw := accessProfile["ip_rule"].([]interface{})
	ipRules := make([]batchaccount.IPRule, 0)
	for _, ipRule := range ipRulesRaw {
		ipRuleRaw := ipRule.(map[string]interface{})
		ipRules = append(ipRules, batchaccount.IPRule{
			Action: batchaccount.IPRuleAction(ipRuleRaw["action"].(string)),
			Value:  ipRuleRaw["ip_range"].(string),
		})
	}

	return &batchaccount.EndpointAccessProfile{
		DefaultAction: batchaccount.EndpointAccessDefaultAction(accessProfile["default_action"].(string)),
		IPRules:       pointer.To(ipRules),
	}
}

func flattenAllowedAuthenticationModes(input *[]batchaccount.AuthenticationMode) []string {
	if input == nil || len(*input) == 0 {
		return []string{}
	}

	allowedAuthModes := make([]string, 0)
	for _, mode := range *input {
		allowedAuthModes = append(allowedAuthModes, string(mode))
	}
	return allowedAuthModes
}

func flattenEncryption(encryptionProperties *batchaccount.EncryptionProperties) []interface{} {
	if encryptionProperties == nil || *encryptionProperties.KeySource == batchaccount.KeySourceMicrosoftPointBatch {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"key_vault_key_id": *encryptionProperties.KeyVaultProperties.KeyIdentifier,
		},
	}
}

func flattenBatchAccountNetworkProfile(input *batchaccount.NetworkProfile) []interface{} {
	if input == nil || input.AccountAccess == nil && input.NodeManagementAccess == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"account_access":         flattenBatchAccountEndpointAccessProfile(input.AccountAccess),
			"node_management_access": flattenBatchAccountEndpointAccessProfile(input.NodeManagementAccess),
		},
	}
}

func flattenBatchAccountEndpointAccessProfile(input *batchaccount.EndpointAccessProfile) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	ipRules := make([]interface{}, 0)
	if input.IPRules != nil {
		for _, ipRule := range *input.IPRules {
			flattenedIpRule := map[string]interface{}{
				"action":   string(ipRule.Action),
				"ip_range": ipRule.Value,
			}
			ipRules = append(ipRules, flattenedIpRule)
		}

	}

	return []interface{}{
		map[string]interface{}{
			"default_action": string(input.DefaultAction),
			"ip_rule":        ipRules,
		},
	}
}

func isShardKeyAllowed(input []interface{}) bool {
	if len(input) == 0 {
		return false
	}
	for _, authMod := range input {
		if strings.EqualFold(authMod.(string), string(batchaccount.AuthenticationModeSharedKey)) {
			return true
		}
	}
	return false
}

func resourceBatchAccountEndpointAccessProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeList,
		Optional:     true,
		MaxItems:     1,
		AtLeastOneOf: []string{"network_profile.0.account_access", "network_profile.0.node_management_access"},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"default_action": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  batchaccount.EndpointAccessDefaultActionDeny,
					ValidateFunc: validation.StringInSlice([]string{
						string(batchaccount.EndpointAccessDefaultActionAllow),
						string(batchaccount.EndpointAccessDefaultActionDeny),
					}, false),
				},

				"ip_rule": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"ip_range": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validate.BatchAccountIpRange,
							},

							"action": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								Default:  string(batchaccount.IPRuleActionAllow),
								ValidateFunc: validation.StringInSlice([]string{
									string(batchaccount.IPRuleActionAllow),
								}, false),
							},
						},
					},
				},
			},
		},
	}
}
