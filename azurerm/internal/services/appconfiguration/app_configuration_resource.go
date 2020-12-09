package appconfiguration

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/appconfiguration/mgmt/2019-10-01/appconfiguration"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appconfiguration/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appconfiguration/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceAppConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppConfigurationCreate,
		Read:   resourceAppConfigurationRead,
		Update: resourceAppConfigurationUpdate,
		Delete: resourceAppConfigurationDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ConfigurationStoreID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ConfigurationStoreName,
			},

			"location": azure.SchemaLocation(),

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(appconfiguration.SystemAssigned),
							}, false),
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			// the API changed and now returns the rg in lowercase
			// revert when https://github.com/Azure/azure-sdk-for-go/issues/6606 is fixed
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"sku": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "free",
				ValidateFunc: validation.StringInSlice([]string{
					"free",
					"standard",
				}, false),
			},

			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_read_key": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"secret": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"connection_string": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"secondary_read_key": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"secret": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"connection_string": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"primary_write_key": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"secret": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"connection_string": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"secondary_write_key": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"secret": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"connection_string": {
							Type:      schema.TypeString,
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

func resourceAppConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppConfiguration.AppConfigurationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM App Configuration creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	resourceId := parse.NewConfigurationStoreID(subscriptionId, resourceGroup, name).ID("")
	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of existing App Configuration %q (Resource Group %q): %s", name, resourceGroup, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_app_configuration", resourceId)
	}

	parameters := appconfiguration.ConfigurationStore{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Sku: &appconfiguration.Sku{
			Name: utils.String(d.Get("sku").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	parameters.Identity = expandAppConfigurationIdentity(d.Get("identity").([]interface{}))

	future, err := client.Create(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating App Configuration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of App Configuration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(resourceId)
	return resourceAppConfigurationRead(d, meta)
}

func resourceAppConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppConfiguration.AppConfigurationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM App Configuration update.")
	id, err := parse.ConfigurationStoreID(d.Id())
	if err != nil {
		return err
	}

	parameters := appconfiguration.ConfigurationStoreUpdateParameters{
		Sku: &appconfiguration.Sku{
			Name: utils.String(d.Get("sku").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if d.HasChange("identity") {
		parameters.Identity = expandAppConfigurationIdentity(d.Get("identity").([]interface{}))
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("Error updating App Configuration %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of App Configuration %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error retrieving App Configuration %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read App Configuration %s (resource Group %q) ID", id.Name, id.ResourceGroup)
	}

	d.SetId(*read.ID)

	return resourceAppConfigurationRead(d, meta)
}

func resourceAppConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppConfiguration.AppConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ConfigurationStoreID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Configuration %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on App Configuration %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	resultPage, err := client.ListKeys(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		return fmt.Errorf("Failed to receive access keys for App Configuration %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	skuName := ""
	if resp.Sku != nil && resp.Sku.Name != nil {
		skuName = *resp.Sku.Name
	}
	d.Set("sku", skuName)

	if props := resp.ConfigurationStoreProperties; props != nil {
		d.Set("endpoint", props.Endpoint)
	}

	accessKeys := flattenAppConfigurationAccessKeys(resultPage.Values())
	d.Set("primary_read_key", accessKeys.primaryReadKey)
	d.Set("primary_write_key", accessKeys.primaryWriteKey)
	d.Set("secondary_read_key", accessKeys.secondaryReadKey)
	d.Set("secondary_write_key", accessKeys.secondaryWriteKey)

	if err := d.Set("identity", flattenAppConfigurationIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("Error setting `identity`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceAppConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppConfiguration.AppConfigurationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ConfigurationStoreID(d.Id())
	if err != nil {
		return err
	}

	fut, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(fut.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting App Configuration %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = fut.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(fut.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting App Configuration %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

type flattenedAccessKeys struct {
	primaryReadKey    []interface{}
	primaryWriteKey   []interface{}
	secondaryReadKey  []interface{}
	secondaryWriteKey []interface{}
}

func flattenAppConfigurationAccessKeys(values []appconfiguration.APIKey) flattenedAccessKeys {
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

func flattenAppConfigurationAccessKey(input appconfiguration.APIKey) []interface{} {
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

func expandAppConfigurationIdentity(identities []interface{}) *appconfiguration.ResourceIdentity {
	if len(identities) == 0 {
		return &appconfiguration.ResourceIdentity{
			Type: appconfiguration.None,
		}
	}
	identity := identities[0].(map[string]interface{})
	identityType := appconfiguration.IdentityType(identity["type"].(string))
	return &appconfiguration.ResourceIdentity{
		Type: identityType,
	}
}

func flattenAppConfigurationIdentity(identity *appconfiguration.ResourceIdentity) []interface{} {
	if identity == nil || identity.Type == appconfiguration.None {
		return []interface{}{}
	}

	principalId := ""
	if identity.PrincipalID != nil {
		principalId = *identity.PrincipalID
	}

	tenantId := ""
	if identity.TenantID != nil {
		tenantId = *identity.TenantID
	}

	return []interface{}{
		map[string]interface{}{
			"type":         string(identity.Type),
			"principal_id": principalId,
			"tenant_id":    tenantId,
		},
	}
}
