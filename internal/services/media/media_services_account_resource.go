package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-05-01/accounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/storageaccounts"
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
	client := meta.(*clients.Client).Media.V20210501Client.Accounts
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

	identity, err := identity.ExpandSystemAssigned(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	payload := accounts.MediaService{
		Location: location.Normalize(d.Get("location").(string)),
		Identity: identity,
		Properties: &accounts.MediaServiceProperties{
			StorageAccounts: storageAccounts,
		},
		Tags: tags.Expand(t),
	}

	if keyDelivery, ok := d.GetOk("key_delivery_access_control"); ok {
		payload.Properties.KeyDelivery = expandKeyDelivery(keyDelivery.([]interface{}))
	}

	if v, ok := d.GetOk("storage_authentication_type"); ok {
		payload.Properties.StorageAuthentication = pointer.To(accounts.StorageAuthentication(v.(string)))
	}

	if _, err := client.MediaservicesCreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMediaServicesAccountRead(d, meta)
}

func resourceMediaServicesAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20210501Client.Accounts
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

	d.Set("name", id.AccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		if err := d.Set("identity", identity.FlattenSystemAssigned(model.Identity)); err != nil {
			return fmt.Errorf("flattening `identity`: %s", err)
		}

		if props := model.Properties; props != nil {
			accounts, err := flattenMediaServicesAccountStorageAccounts(props.StorageAccounts)
			if err != nil {
				return fmt.Errorf("flattening `storage_account`: %s", err)
			}
			if err := d.Set("storage_account", accounts); err != nil {
				return fmt.Errorf("setting `storage_account`: %s", err)
			}

			storageAuthenticationType := ""
			if props.StorageAuthentication != nil {
				storageAuthenticationType = string(*props.StorageAuthentication)
			}
			d.Set("storage_authentication_type", storageAuthenticationType)

			if err := d.Set("key_delivery_access_control", flattenKeyDelivery(props.KeyDelivery)); err != nil {
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
	client := meta.(*clients.Client).Media.V20210501Client.Accounts
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

		results = append(results, accounts.StorageAccount{
			Id:   utils.String(id),
			Type: storageType,
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
			id, err := storageaccounts.ParseStorageAccountIDInsensitively(*storageAccount.Id)
			if err != nil {
				return nil, fmt.Errorf("parsing %q as a Storage Account ID: %+v", *storageAccount.Id, err)
			}
			storageAccountId = id.ID()
		}

		results = append(results, map[string]interface{}{
			"id":         storageAccountId,
			"is_primary": storageAccount.Type == accounts.StorageAccountTypePrimary,
		})
	}

	return &results, nil
}

func expandKeyDelivery(input []interface{}) *accounts.KeyDelivery {
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

	return &accounts.KeyDelivery{
		AccessControl: &accounts.AccessControl{
			DefaultAction: pointer.To(accounts.DefaultAction(defaultAction)),
			IPAllowList:   ipAllowList,
		},
	}
}

func flattenKeyDelivery(input *accounts.KeyDelivery) []interface{} {
	if input == nil && input.AccessControl != nil {
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
