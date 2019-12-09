package azurerm

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	appconf "github.com/Azure/azure-sdk-for-go/services/appconfiguration/mgmt/2019-10-01/appconfiguration"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppConfigurationCreate,
		Read:   resourceArmAppConfigurationRead,
		Update: resourceArmAppConfigurationUpdate,
		Delete: resourceArmAppConfigurationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAppConfigurationName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

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

func resourceArmAppConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppConfiguration.AppConfigurationsClient
	ctx, cancel := timeouts.ForCreate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM App Configuration creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing App Configuration %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_app_configuration", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	skuName := d.Get("sku").(string)
	sku := appconf.Sku{
		Name: &skuName,
	}

	parameters := appconf.ConfigurationStore{
		Location: &location,
		Sku:      &sku,
		Tags:     tags.Expand(t),
	}

	future, err := client.Create(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating App Configuration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of App Configuration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving App Configuration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read App Configuration %s (resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAppConfigurationRead(d, meta)
}

func resourceArmAppConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppConfiguration.AppConfigurationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM App Configuration update.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	t := d.Get("tags").(map[string]interface{})
	skuName := d.Get("sku").(string)
	sku := appconf.Sku{
		Name: &skuName,
	}

	parameters := appconf.ConfigurationStoreUpdateParameters{
		Sku:  &sku,
		Tags: tags.Expand(t),
	}

	future, err := client.Update(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error updating App Configuration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of App Configuration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving App Configuration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read App Configuration %s (resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAppConfigurationRead(d, meta)
}

func resourceArmAppConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppConfiguration.AppConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["configurationStores"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Configuration %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on App Configuration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", sku.Name)
	}

	if endpoint := resp.Endpoint; endpoint != nil {
		d.Set("endpoint", resp.Endpoint)
	}

	resultPage, err := client.ListKeys(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("Failed to receive access keys for App Configuration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	values := resultPage.Values()
	for _, value := range values {
		accessKeyParams := map[string]string{}
		if id := value.ID; id != nil {
			accessKeyParams["id"] = *id
		}
		if value := value.Value; value != nil {
			accessKeyParams["secret"] = *value
		}
		if connectionString := value.ConnectionString; connectionString != nil {
			accessKeyParams["connection_string"] = *value.ConnectionString
		}

		accessKey := []interface{}{accessKeyParams}
		if name := value.Name; name != nil {
			if strings.HasPrefix(*name, "Primary") {
				if readOnly := value.ReadOnly; readOnly != nil {
					if *readOnly {
						d.Set("primary_read_key", accessKey)
					} else {
						d.Set("primary_write_key", accessKey)
					}
				}
			} else if strings.HasPrefix(*name, "Secondary") {
				if readOnly := value.ReadOnly; readOnly != nil {
					if *readOnly {
						d.Set("secondary_read_key", accessKey)
					} else {
						d.Set("secondary_write_key", accessKey)
					}
				}
			} else {
				log.Printf("[WARN] Received unknown App Configuration access key '%s', ignoring...", *name)
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmAppConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppConfiguration.AppConfigurationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["configurationStores"]

	fut, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(fut.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting App Configuration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = fut.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(fut.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting App Configuration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func validateAppConfigurationName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if matched := regexp.MustCompile(`^[a-zA-Z0-9-]{5,50}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes and must be between 5-50 chars", k))
	}

	return warnings, errors
}
