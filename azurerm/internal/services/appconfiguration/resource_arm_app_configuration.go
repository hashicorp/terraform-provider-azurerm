package appconfiguration

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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appconfiguration/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppConfigurationCreate,
		Read:   resourceArmAppConfigurationRead,
		Update: resourceArmAppConfigurationUpdate,
		Delete: resourceArmAppConfigurationDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.AppConfigurationID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateAppConfigurationName,
			},

			"location": azure.SchemaLocation(),

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

func resourceArmAppConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppConfiguration.AppConfigurationsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
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

	parameters := appconf.ConfigurationStore{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Sku: &appconf.Sku{
			Name: utils.String(d.Get("sku").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
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
	client := meta.(*clients.Client).AppConfiguration.AppConfigurationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM App Configuration update.")
	id, err := parse.AppConfigurationID(d.Id())
	if err != nil {
		return err
	}

	parameters := appconf.ConfigurationStoreUpdateParameters{
		Sku: &appconf.Sku{
			Name: utils.String(d.Get("sku").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
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

	return resourceArmAppConfigurationRead(d, meta)
}

func resourceArmAppConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppConfiguration.AppConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppConfigurationID(d.Id())
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

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", sku.Name)
	}

	if endpoint := resp.Endpoint; endpoint != nil {
		d.Set("endpoint", resp.Endpoint)
	}

	resultPage, err := client.ListKeys(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		return fmt.Errorf("Failed to receive access keys for App Configuration %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
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
		if n := value.Name; n != nil {
			if strings.HasPrefix(*n, "Primary") {
				if readOnly := value.ReadOnly; readOnly != nil {
					if *readOnly {
						d.Set("primary_read_key", accessKey)
					} else {
						d.Set("primary_write_key", accessKey)
					}
				}
			} else if strings.HasPrefix(*n, "Secondary") {
				if readOnly := value.ReadOnly; readOnly != nil {
					if *readOnly {
						d.Set("secondary_read_key", accessKey)
					} else {
						d.Set("secondary_write_key", accessKey)
					}
				}
			} else {
				log.Printf("[WARN] Received unknown App Configuration access key '%s', ignoring...", *n)
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmAppConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppConfiguration.AppConfigurationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppConfigurationID(d.Id())
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

func ValidateAppConfigurationName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if matched := regexp.MustCompile(`^[a-zA-Z0-9-]{5,50}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes and must be between 5-50 chars", k))
	}

	return warnings, errors
}
