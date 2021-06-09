package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2021-05-01/media"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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
			_, err := parse.MediaServiceID(id)
			return err
		}),

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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"storage_account": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"is_primary": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			//lintignore:XS003
			"identity": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"principal_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"type": {
							Type:             pluginsdk.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(media.ManagedIdentityTypeSystemAssigned),
							}, true),
						},
					},
				},
			},

			"storage_authentication_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(media.StorageAuthenticationSystem),
					string(media.StorageAuthenticationManagedIdentity),
				}, true),
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
							}, true),
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

	resourceId := parse.NewMediaServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", resourceId, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_media_services_account", resourceId.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	storageAccountsRaw := d.Get("storage_account").(*pluginsdk.Set).List()
	storageAccounts, err := expandMediaServicesAccountStorageAccounts(storageAccountsRaw)
	if err != nil {
		return err
	}
	parameters := media.Service{
		ServiceProperties: &media.ServiceProperties{
			StorageAccounts: storageAccounts,
		},
		Location: utils.String(location),
		Tags:     tags.Expand(t),
	}

	if _, ok := d.GetOk("identity"); ok {
		parameters.Identity = expandAzureRmMediaServiceIdentity(d)
	}

	if v, ok := d.GetOk("storage_authentication_type"); ok {
		parameters.StorageAuthentication = media.StorageAuthentication(v.(string))
	}

	if keyDelivery, ok := d.GetOk("key_delivery_access_control"); ok {
		parameters.KeyDelivery = expandKeyDelivery(keyDelivery.([]interface{}))
	}

	if _, err := client.CreateOrUpdate(ctx, resourceId.ResourceGroup, resourceId.Name, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", resourceId, err)
	}

	d.SetId(resourceId.ID())
	return resourceMediaServicesAccountRead(d, meta)
}

func resourceMediaServicesAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.ServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MediaServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Media Services Account %q was not found in Resource Group %q - removing from state", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Media Services Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	props := resp.ServiceProperties
	if props != nil {
		accounts := flattenMediaServicesAccountStorageAccounts(props.StorageAccounts)
		if e := d.Set("storage_account", accounts); e != nil {
			return fmt.Errorf("flattening `storage_account`: %s", e)
		}
		d.Set("storage_authentication_type", string(props.StorageAuthentication))
	}

	if err := d.Set("identity", flattenAzureRmMediaServicedentity(resp.Identity)); err != nil {
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

	id, err := parse.MediaServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("issuing AzureRM delete request for Media Services Account '%s': %+v", id.Name, err)
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

func flattenMediaServicesAccountStorageAccounts(input *[]media.StorageAccount) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, storageAccount := range *input {
		output := make(map[string]interface{})

		if storageAccount.ID != nil {
			output["id"] = *storageAccount.ID
		}

		output["is_primary"] = storageAccount.Type == media.StorageAccountTypePrimary

		results = append(results, output)
	}

	return results
}

func expandAzureRmMediaServiceIdentity(d *pluginsdk.ResourceData) *media.ServiceIdentity {
	identities := d.Get("identity").([]interface{})
	if identities[0] == nil {
		return nil
	}
	identity := identities[0].(map[string]interface{})
	identityType := identity["type"].(string)
	return &media.ServiceIdentity{
		Type: media.ManagedIdentityType(identityType),
	}
}

func flattenAzureRmMediaServicedentity(identity *media.ServiceIdentity) []interface{} {
	if identity == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})
	result["type"] = string(identity.Type)
	if identity.PrincipalID != nil {
		result["principal_id"] = *identity.PrincipalID
	}
	if identity.TenantID != nil {
		result["tenant_id"] = *identity.TenantID
	}

	return []interface{}{result}
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
