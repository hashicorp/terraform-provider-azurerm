// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package media

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/accounts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/media/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMediaServicesAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMediaServicesAccountCreateUpdate,
		Read:   resourceMediaServicesAccountRead,
		Update: resourceMediaServicesAccountCreateUpdate,
		Delete: resourceMediaServicesAccountDelete,

		DeprecationMessage: azureMediaRetirementMessage,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := accounts.ParseMediaServiceID(id)
			return err
		}),

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ServiceV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-z0-9]{3,24}$"),
					"Media Services Account name must be 3 - 24 characters long, contain only lowercase letters and numbers.",
				),
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"storage_account": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: commonids.ValidateStorageAccountID,
						},

						"is_primary": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"managed_identity": mediaServicesAccountUseManagedIdentity(),
					},
				},
			},

			"encryption": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      string(accounts.AccountEncryptionKeyTypeSystemKey),
							ValidateFunc: validation.StringInSlice(accounts.PossibleValuesForAccountEncryptionKeyType(), false),
						},

						"key_vault_key_identifier": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsURLWithHTTPorHTTPS,
						},

						"current_key_identifier": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"managed_identity": mediaServicesAccountUseManagedIdentity(),
					},
				},
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"storage_authentication_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(accounts.StorageAuthenticationSystem),
					string(accounts.StorageAuthenticationManagedIdentity),
				}, false),
			},

			"key_delivery_access_control": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"default_action": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(accounts.DefaultActionDeny),
								string(accounts.DefaultActionAllow),
							}, false),
						},

						"ip_allow_list": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceMediaServicesAccountCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20211101Client.Accounts
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := accounts.NewMediaServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.MediaservicesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_media_services_account", id.ID())
		}
	}

	t := d.Get("tags").(map[string]interface{})

	storageAccountsRaw := d.Get("storage_account").(*pluginsdk.Set).List()
	storageAccounts, err := expandMediaServicesAccountStorageAccounts(storageAccountsRaw)
	if err != nil {
		return err
	}

	identity, err := expandMediaServicesAccountIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	publicNetworkAccess := accounts.PublicNetworkAccessDisabled
	if d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = accounts.PublicNetworkAccessEnabled
	}

	payload := accounts.MediaService{
		Location: location.Normalize(d.Get("location").(string)),
		Identity: identity,
		Properties: &accounts.MediaServiceProperties{
			StorageAccounts:     storageAccounts,
			PublicNetworkAccess: &publicNetworkAccess,
		},
		Tags: tags.Expand(t),
	}

	if encryptionRaw, ok := d.GetOk("encryption"); ok {
		encryption, err := expandMediaServicesAccountEncryption(encryptionRaw.([]interface{}))
		if err != nil {
			return err
		}
		payload.Properties.Encryption = encryption
	}

	if keyDelivery, ok := d.GetOk("key_delivery_access_control"); ok {
		payload.Properties.KeyDelivery = expandMediaServicesAccountKeyDelivery(keyDelivery.([]interface{}))
	}

	if v, ok := d.GetOk("storage_authentication_type"); ok {
		payload.Properties.StorageAuthentication = pointer.To(accounts.StorageAuthentication(v.(string)))
	}

	if err := client.MediaservicesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMediaServicesAccountRead(d, meta)
}

func resourceMediaServicesAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20211101Client.Accounts
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := accounts.ParseMediaServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.MediaservicesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %q was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.MediaServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		accountIdentity, err := flattenMediaServicesAccountIdentity(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %s", err)
		}
		if err := d.Set("identity", accountIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %s", err)
		}

		if props := model.Properties; props != nil {
			storageAccounts, err := flattenMediaServicesAccountStorageAccounts(props.StorageAccounts)
			if err != nil {
				return fmt.Errorf("flattening `storage_account`: %s", err)
			}
			if err := d.Set("storage_account", storageAccounts); err != nil {
				return fmt.Errorf("setting `storage_account`: %s", err)
			}

			if err := d.Set("encryption", flattenMediaServicesAccountEncryption(props.Encryption)); err != nil {
				return fmt.Errorf("setting `encryption`: %s", err)
			}

			publicNetworkAccess := false
			if props.PublicNetworkAccess != nil && *props.PublicNetworkAccess == accounts.PublicNetworkAccessEnabled {
				publicNetworkAccess = true
			}
			d.Set("public_network_access_enabled", publicNetworkAccess)

			storageAuthenticationType := ""
			if props.StorageAuthentication != nil {
				storageAuthenticationType = string(*props.StorageAuthentication)
			}
			d.Set("storage_authentication_type", storageAuthenticationType)

			if err := d.Set("key_delivery_access_control", flattenMediaServicesAccountKeyDelivery(props.KeyDelivery)); err != nil {
				return fmt.Errorf("flattening `key_delivery_access_control`: %s", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceMediaServicesAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20211101Client.Accounts
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := accounts.ParseMediaServiceID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.MediaservicesDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func mediaServicesAccountUseManagedIdentity() *schema.Schema {
	return &schema.Schema{
		Type:     pluginsdk.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"user_assigned_identity_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: commonids.ValidateUserAssignedIdentityID,
				},

				"use_system_assigned_identity": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
			},
		},
	}
}

func expandMediaServicesAccountStorageAccounts(input []interface{}) (*[]accounts.StorageAccount, error) {
	results := make([]accounts.StorageAccount, 0)

	foundPrimary := false
	for _, accountMapRaw := range input {
		accountMap := accountMapRaw.(map[string]interface{})

		id := accountMap["id"].(string)

		storageType := accounts.StorageAccountTypeSecondary
		if accountMap["is_primary"].(bool) {
			if foundPrimary {
				return nil, fmt.Errorf("Only one Storage Account can be set as Primary")
			}

			storageType = accounts.StorageAccountTypePrimary
			foundPrimary = true
		}

		resourceIdentity, err := expandMediaServicesAccountManagedIdentity(accountMap["managed_identity"].([]interface{}))
		if err != nil {
			return nil, err
		}
		results = append(results, accounts.StorageAccount{
			Id:       utils.String(id),
			Type:     storageType,
			Identity: resourceIdentity,
		})
	}

	return &results, nil
}

func flattenMediaServicesAccountStorageAccounts(input *[]accounts.StorageAccount) (*[]interface{}, error) {
	if input == nil {
		return &[]interface{}{}, nil
	}

	results := make([]interface{}, 0)
	for _, storageAccount := range *input {
		storageAccountId := ""
		if storageAccount.Id != nil {
			id, err := commonids.ParseStorageAccountIDInsensitively(*storageAccount.Id)
			if err != nil {
				return nil, fmt.Errorf("parsing %q as a Storage Account ID: %+v", *storageAccount.Id, err)
			}
			storageAccountId = id.ID()
		}

		results = append(results, map[string]interface{}{
			"id":               storageAccountId,
			"is_primary":       storageAccount.Type == accounts.StorageAccountTypePrimary,
			"managed_identity": flattenMediaServicesAccountManagedIdentity(storageAccount.Identity),
		})
	}

	return &results, nil
}

func expandMediaServicesAccountEncryption(input []interface{}) (*accounts.AccountEncryption, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}
	val := input[0].(map[string]interface{})

	resourceIdentity, err := expandMediaServicesAccountManagedIdentity(val["managed_identity"].([]interface{}))
	if err != nil {
		return nil, err
	}

	accountsEncryption := accounts.AccountEncryption{
		Type:     accounts.AccountEncryptionKeyType(val["type"].(string)),
		Identity: resourceIdentity,
	}

	if keyIdentifier, ok := val["key_vault_key_identifier"].(string); ok && keyIdentifier != "" {
		if accountsEncryption.Type != accounts.AccountEncryptionKeyTypeCustomerKey {
			return nil, fmt.Errorf("key_vault_key_identifier can only be set when encryption type is 'CustomerKey'")
		}
		accountsEncryption.KeyVaultProperties = &accounts.KeyVaultProperties{
			KeyIdentifier: &keyIdentifier,
		}
	}

	return &accountsEncryption, nil
}

func flattenMediaServicesAccountEncryption(input *accounts.AccountEncryption) *[]interface{} {
	if input == nil {
		return &[]interface{}{}
	}

	var keyIdentifier, currentKeyIdentifier string

	if input.KeyVaultProperties != nil {
		if input.KeyVaultProperties.KeyIdentifier != nil {
			keyIdentifier = *input.KeyVaultProperties.KeyIdentifier
		}
		if input.KeyVaultProperties.CurrentKeyIdentifier != nil {
			currentKeyIdentifier = *input.KeyVaultProperties.CurrentKeyIdentifier
		}
	}

	return &[]interface{}{
		map[string]interface{}{
			"type":                     string(input.Type),
			"key_vault_key_identifier": keyIdentifier,
			"current_key_identifier":   currentKeyIdentifier,
			"managed_identity":         flattenMediaServicesAccountManagedIdentity(input.Identity),
		},
	}
}

func expandMediaServicesAccountIdentity(input []interface{}) (*accounts.MediaServiceIdentity, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	identityType := string(expanded.Type)
	// https://github.com/Azure/azure-rest-api-specs/issues/21905
	if identityType == string(identity.TypeSystemAssignedUserAssigned) {
		identityType = "SystemAssigned,UserAssigned"
	}
	out := accounts.MediaServiceIdentity{
		Type: identityType,
	}
	if expanded.Type == identity.TypeUserAssigned || expanded.Type == identity.TypeSystemAssignedUserAssigned {
		userAssignedIdentities := make(map[string]accounts.UserAssignedManagedIdentity)
		for k := range expanded.IdentityIds {
			userAssignedIdentities[k] = accounts.UserAssignedManagedIdentity{
				// intentionally empty
			}
		}
		out.UserAssignedIdentities = &userAssignedIdentities
	}
	return &out, nil
}

func flattenMediaServicesAccountIdentity(input *accounts.MediaServiceIdentity) (*[]interface{}, error) {
	var transform *identity.SystemAndUserAssignedMap

	if input != nil {
		identityType := identity.Type(input.Type)
		if strings.EqualFold(input.Type, "SystemAssigned,UserAssigned") {
			identityType = identity.TypeSystemAssignedUserAssigned
		}
		transform = &identity.SystemAndUserAssignedMap{
			Type: identityType,
		}
		if input.PrincipalId != nil {
			transform.PrincipalId = *input.PrincipalId
		}
		if input.TenantId != nil {
			transform.TenantId = *input.TenantId
		}
		if input.UserAssignedIdentities != nil {
			transform.IdentityIds = make(map[string]identity.UserAssignedIdentityDetails)
			for k, v := range *input.UserAssignedIdentities {
				transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
					ClientId:    v.ClientId,
					PrincipalId: v.PrincipalId,
				}
			}
		}
	}

	return identity.FlattenSystemAndUserAssignedMap(transform)
}

func expandMediaServicesAccountKeyDelivery(input []interface{}) *accounts.KeyDelivery {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	keyDelivery := input[0].(map[string]interface{})
	defaultAction := keyDelivery["default_action"].(string)

	var ipAllowList *[]string
	if v := keyDelivery["ip_allow_list"]; v != nil {
		ips := keyDelivery["ip_allow_list"].(*pluginsdk.Set).List()
		ipAllowList = utils.ExpandStringSlice(ips)
	}

	return &accounts.KeyDelivery{
		AccessControl: &accounts.AccessControl{
			DefaultAction: pointer.To(accounts.DefaultAction(defaultAction)),
			IPAllowList:   ipAllowList,
		},
	}
}

func flattenMediaServicesAccountKeyDelivery(input *accounts.KeyDelivery) []interface{} {
	if input == nil || input.AccessControl == nil {
		return make([]interface{}, 0)
	}

	defaultAction := ""
	if input.AccessControl.DefaultAction != nil {
		defaultAction = string(*input.AccessControl.DefaultAction)
	}

	return []interface{}{
		map[string]interface{}{
			"default_action": defaultAction,
			"ip_allow_list":  utils.FlattenStringSlice(input.AccessControl.IPAllowList),
		},
	}
}

func expandMediaServicesAccountManagedIdentity(input []interface{}) (*accounts.ResourceIdentity, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	ManagedIdentity := input[0].(map[string]interface{})

	result := &accounts.ResourceIdentity{
		UseSystemAssignedIdentity: ManagedIdentity["use_system_assigned_identity"].(bool),
	}
	if userAssignedIdentityId := ManagedIdentity["user_assigned_identity_id"].(string); userAssignedIdentityId != "" {
		if result.UseSystemAssignedIdentity {
			return nil, fmt.Errorf("use either of user assigned identity or system assigned identity for ecryption")
		}
		result.UserAssignedIdentity = &userAssignedIdentityId
	}
	return result, nil
}

func flattenMediaServicesAccountManagedIdentity(input *accounts.ResourceIdentity) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var userAssignedIdentity string
	if input.UserAssignedIdentity != nil {
		userAssignedIdentity = *input.UserAssignedIdentity
	}

	return []interface{}{
		map[string]interface{}{
			"use_system_assigned_identity": input.UseSystemAssignedIdentity,
			"user_assigned_identity_id":    userAssignedIdentity,
		},
	}
}
