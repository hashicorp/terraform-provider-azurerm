package media

import (
	"fmt"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/storageaccounts"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2021-05-01/media"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-05-01/accounts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/media/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
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
							ValidateFunc: storageaccounts.ValidateStorageAccountID,
						},

						"is_primary": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"identity": commonschema.SystemAssignedIdentityOptional(),

			"storage_authentication_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(media.StorageAuthenticationSystem),
					string(media.StorageAuthenticationManagedIdentity),
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
								string(media.DefaultActionDeny),
								string(media.DefaultActionAllow),
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

			"tags": tags.Schema(),
		},
	}
}

func resourceMediaServicesAccountCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.ServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := accounts.NewMediaServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroupName, id.AccountName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_media_services_account", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	storageAccountsRaw := d.Get("storage_account").(*pluginsdk.Set).List()
	storageAccounts, err := expandMediaServicesAccountStorageAccounts(storageAccountsRaw)
	if err != nil {
		return err
	}

	expandedIdentity, err := expandAccountIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	parameters := media.Service{
		Location: utils.String(location),
		Identity: expandedIdentity,
		ServiceProperties: &media.ServiceProperties{
			StorageAccounts: storageAccounts,
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("storage_authentication_type"); ok {
		parameters.StorageAuthentication = media.StorageAuthentication(v.(string))
	}

	if keyDelivery, ok := d.GetOk("key_delivery_access_control"); ok {
		parameters.KeyDelivery = expandKeyDelivery(keyDelivery.([]interface{}))
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.AccountName, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMediaServicesAccountRead(d, meta)
}

func resourceMediaServicesAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.ServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := accounts.ParseMediaServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.AccountName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %q was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AccountName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.ServiceProperties; props != nil {
		accounts, err := flattenMediaServicesAccountStorageAccounts(props.StorageAccounts)
		if err != nil {
			return fmt.Errorf("flattening `storage_account`: %s", err)
		}
		if err := d.Set("storage_account", accounts); err != nil {
			return fmt.Errorf("setting `storage_account`: %s", err)
		}
		d.Set("storage_authentication_type", string(props.StorageAuthentication))
	}

	if err := d.Set("identity", flattenAccountIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("flattening `identity`: %s", err)
	}

	if err := d.Set("key_delivery_access_control", flattenKeyDelivery(resp.KeyDelivery)); err != nil {
		return fmt.Errorf("flattening `key_delivery_access_control`: %s", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceMediaServicesAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.ServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := accounts.ParseMediaServiceID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroupName, id.AccountName); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandMediaServicesAccountStorageAccounts(input []interface{}) (*[]media.StorageAccount, error) {
	results := make([]media.StorageAccount, 0)

	foundPrimary := false
	for _, accountMapRaw := range input {
		accountMap := accountMapRaw.(map[string]interface{})

		id := accountMap["id"].(string)

		storageType := media.StorageAccountTypeSecondary
		if accountMap["is_primary"].(bool) {
			if foundPrimary {
				return nil, fmt.Errorf("Only one Storage Account can be set as Primary")
			}

			storageType = media.StorageAccountTypePrimary
			foundPrimary = true
		}

		storageAccount := media.StorageAccount{
			ID:   utils.String(id),
			Type: storageType,
		}

		results = append(results, storageAccount)
	}

	return &results, nil
}

func flattenMediaServicesAccountStorageAccounts(input *[]media.StorageAccount) (*[]interface{}, error) {
	if input == nil {
		return &[]interface{}{}, nil
	}

	results := make([]interface{}, 0)
	for _, storageAccount := range *input {
		storageAccountId := ""
		if storageAccount.ID != nil {
			id, err := storageaccounts.ParseStorageAccountIDInsensitively(*storageAccount.ID)
			if err != nil {
				return nil, fmt.Errorf("parsing %q as a Storage Account ID: %+v", *storageAccount.ID, err)
			}
			storageAccountId = id.ID()
		}

		results = append(results, map[string]interface{}{
			"id":         storageAccountId,
			"is_primary": storageAccount.Type == media.StorageAccountTypePrimary,
		})
	}

	return &results, nil
}

func expandAccountIdentity(input []interface{}) (*media.ServiceIdentity, error) {
	expanded, err := identity.ExpandSystemAssigned(input)
	if err != nil {
		return nil, err
	}

	return &media.ServiceIdentity{
		Type: media.ManagedIdentityType(string(expanded.Type)),
	}, nil
}

func flattenAccountIdentity(input *media.ServiceIdentity) []interface{} {
	var transform *identity.SystemAssigned

	if input != nil {
		transform = &identity.SystemAssigned{
			Type: identity.Type(string(input.Type)),
		}
		if input.PrincipalID != nil {
			transform.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transform.TenantId = *input.TenantID
		}
	}

	return identity.FlattenSystemAssigned(transform)
}

func expandKeyDelivery(input []interface{}) *media.KeyDelivery {
	if len(input) == 0 {
		return nil
	}

	keyDelivery := input[0].(map[string]interface{})
	defaultAction := keyDelivery["default_action"].(string)

	var ipAllowList *[]string
	if v := keyDelivery["ip_allow_list"]; v != nil {
		ips := keyDelivery["ip_allow_list"].(*pluginsdk.Set).List()
		ipAllowList = utils.ExpandStringSlice(ips)
	}

	return &media.KeyDelivery{
		AccessControl: &media.AccessControl{
			DefaultAction: media.DefaultAction(defaultAction),
			IPAllowList:   ipAllowList,
		},
	}
}

func flattenKeyDelivery(input *media.KeyDelivery) []interface{} {
	if input == nil && input.AccessControl != nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"default_action": string(input.AccessControl.DefaultAction),
			"ip_allow_list":  utils.FlattenStringSlice(input.AccessControl.IPAllowList),
		},
	}
}
