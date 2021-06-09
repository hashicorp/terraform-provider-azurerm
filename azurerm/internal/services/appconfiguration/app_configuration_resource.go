package appconfiguration

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/appconfiguration/mgmt/2020-06-01/appconfiguration"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appconfiguration/sdk/configurationstores"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appconfiguration/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func resourceAppConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppConfigurationCreate,
		Read:   resourceAppConfigurationRead,
		Update: resourceAppConfigurationUpdate,
		Delete: resourceAppConfigurationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := configurationstores.ConfigurationStoreID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ConfigurationStoreName,
			},

			"location": azure.SchemaLocation(),

			"identity": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(appconfiguration.IdentityTypeSystemAssigned),
							}, false),
						},
						"principal_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			// the API changed and now returns the rg in lowercase
			// revert when https://github.com/Azure/azure-sdk-for-go/issues/6606 is fixed
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "free",
				ValidateFunc: validation.StringInSlice([]string{
					"free",
					"standard",
				}, false),
			},

			"endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_read_key": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"secret": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"connection_string": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"secondary_read_key": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"secret": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"connection_string": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"primary_write_key": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"secret": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"connection_string": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"secondary_write_key": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"secret": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"connection_string": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceAppConfigurationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppConfiguration.ConfigurationStoresClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM App Configuration creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	resourceId := configurationstores.NewConfigurationStoreID(subscriptionId, resourceGroup, name)
	existing, err := client.Get(ctx, resourceId)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", resourceId, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_app_configuration", resourceId.ID())
	}

	parameters := configurationstores.ConfigurationStore{
		Location: azure.NormalizeLocation(d.Get("location").(string)),
		Sku: configurationstores.Sku{
			Name: d.Get("sku").(string),
		},
		Tags: expandTags(d.Get("tags").(map[string]interface{})),
	}

	parameters.Identity = expandAppConfigurationIdentity(d.Get("identity").([]interface{}))

	if err := client.CreateThenPoll(ctx, resourceId, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", resourceId, err)
	}

	d.SetId(resourceId.ID())
	return resourceAppConfigurationRead(d, meta)
}

func resourceAppConfigurationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppConfiguration.ConfigurationStoresClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM App Configuration update.")
	id, err := configurationstores.ConfigurationStoreID(d.Id())
	if err != nil {
		return err
	}

	parameters := configurationstores.ConfigurationStoreUpdateParameters{
		Sku: &configurationstores.Sku{
			Name: d.Get("sku").(string),
		},
		Tags: expandTags(d.Get("tags").(map[string]interface{})),
	}

	if d.HasChange("identity") {
		parameters.Identity = expandAppConfigurationIdentity(d.Get("identity").([]interface{}))
	}

	if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceAppConfigurationRead(d, meta)
}

func resourceAppConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppConfiguration.ConfigurationStoresClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := configurationstores.ConfigurationStoreID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	resultPage, err := client.ListKeysComplete(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving access keys for %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("sku", model.Sku.Name)

		if props := model.Properties; props != nil {
			d.Set("endpoint", props.Endpoint)
		}

		accessKeys := flattenAppConfigurationAccessKeys(resultPage.Items)
		d.Set("primary_read_key", accessKeys.primaryReadKey)
		d.Set("primary_write_key", accessKeys.primaryWriteKey)
		d.Set("secondary_read_key", accessKeys.secondaryReadKey)
		d.Set("secondary_write_key", accessKeys.secondaryWriteKey)

		if err := d.Set("identity", flattenAppConfigurationIdentity(model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		return tags.FlattenAndSet(d, flattenTags(model.Tags))
	}

	return nil
}

func resourceAppConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppConfiguration.ConfigurationStoresClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := configurationstores.ConfigurationStoreID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

type flattenedAccessKeys struct {
	primaryReadKey    []interface{}
	primaryWriteKey   []interface{}
	secondaryReadKey  []interface{}
	secondaryWriteKey []interface{}
}

func flattenAppConfigurationAccessKeys(values []configurationstores.AccessKey) flattenedAccessKeys {
	result := flattenedAccessKeys{
		primaryReadKey:    make([]interface{}, 0),
		primaryWriteKey:   make([]interface{}, 0),
		secondaryReadKey:  make([]interface{}, 0),
		secondaryWriteKey: make([]interface{}, 0),
	}

	for _, value := range values {
		if value.Name == nil || value.ReadOnly == nil {
			continue
		}

		accessKey := flattenAppConfigurationAccessKey(value)
		name := *value.Name
		readOnly := *value.ReadOnly

		if strings.HasPrefix(strings.ToLower(name), "primary") {
			if readOnly {
				result.primaryReadKey = accessKey
			} else {
				result.primaryWriteKey = accessKey
			}
		}

		if strings.HasPrefix(strings.ToLower(name), "secondary") {
			if readOnly {
				result.secondaryReadKey = accessKey
			} else {
				result.secondaryWriteKey = accessKey
			}
		}
	}

	return result
}

func flattenAppConfigurationAccessKey(input configurationstores.AccessKey) []interface{} {
	connectionString := ""

	if input.ConnectionString != nil {
		connectionString = *input.ConnectionString
	}

	id := ""
	if input.ID != nil {
		id = *input.ID
	}

	secret := ""
	if input.Value != nil {
		secret = *input.Value
	}

	return []interface{}{
		map[string]interface{}{
			"connection_string": connectionString,
			"id":                id,
			"secret":            secret,
		},
	}
}

func expandAppConfigurationIdentity(identities []interface{}) *configurationstores.ResourceIdentity {
	var out = func(in configurationstores.IdentityType) *configurationstores.ResourceIdentity {
		return &configurationstores.ResourceIdentity{
			Type: &in,
		}
	}

	if len(identities) == 0 {
		return out(configurationstores.IdentityTypeNone)
	}
	identity := identities[0].(map[string]interface{})
	identityType := configurationstores.IdentityType(identity["type"].(string))
	return out(identityType)
}

func flattenAppConfigurationIdentity(identity *configurationstores.ResourceIdentity) []interface{} {
	if identity == nil || identity.Type == nil || *identity.Type == configurationstores.IdentityTypeNone {
		return []interface{}{}
	}

	identityType := ""
	if identity.Type != nil {
		identityType = string(*identity.Type)
	}

	principalId := ""
	if identity.PrincipalId != nil {
		principalId = *identity.PrincipalId
	}

	tenantId := ""
	if identity.TenantId != nil {
		tenantId = *identity.TenantId
	}

	return []interface{}{
		map[string]interface{}{
			"type":         identityType,
			"principal_id": principalId,
			"tenant_id":    tenantId,
		},
	}
}
