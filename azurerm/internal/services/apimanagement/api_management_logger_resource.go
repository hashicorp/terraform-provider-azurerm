package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementLogger() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiManagementLoggerCreate,
		Read:   resourceApiManagementLoggerRead,
		Update: resourceApiManagementLoggerUpdate,
		Delete: resourceApiManagementLoggerDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": schemaz.SchemaApiManagementChildName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"resource_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"eventhub": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"application_insights"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.ValidateEventHubName(),
						},

						"connection_string": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"application_insights": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"eventhub"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instrumentation_key": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"buffered": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceApiManagementLoggerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.LoggerClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)

	eventHubRaw := d.Get("eventhub").([]interface{})
	appInsightsRaw := d.Get("application_insights").([]interface{})

	if len(eventHubRaw) == 0 && len(appInsightsRaw) == 0 {
		return fmt.Errorf("Either `eventhub` or `application_insights` is required")
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Logger %q (API Management Service %q / Resource Group %q): %s", name, serviceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_logger", *existing.ID)
		}
	}

	parameters := apimanagement.LoggerContract{
		LoggerContractProperties: &apimanagement.LoggerContractProperties{
			IsBuffered:  utils.Bool(d.Get("buffered").(bool)),
			Description: utils.String(d.Get("description").(string)),
		},
	}

	if len(eventHubRaw) > 0 {
		parameters.LoggerType = apimanagement.AzureEventHub
		parameters.Credentials = expandApiManagementLoggerEventHub(eventHubRaw)
	} else if len(appInsightsRaw) > 0 {
		parameters.LoggerType = apimanagement.ApplicationInsights
		parameters.Credentials = expandApiManagementLoggerApplicationInsights(appInsightsRaw)
	}

	if resourceId := d.Get("resource_id").(string); resourceId != "" {
		parameters.ResourceID = utils.String(resourceId)
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, name, parameters, ""); err != nil {
		return fmt.Errorf("creating Logger %q (Resource Group %q / API Management Service %q): %+v", name, resourceGroup, serviceName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		return fmt.Errorf("retrieving Logger %q (Resource Group %q / API Management Service %q): %+v", name, resourceGroup, serviceName, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Logger %q (Resource Group %q / API Management Service %q) ID", name, resourceGroup, serviceName)
	}
	d.SetId(*resp.ID)

	return resourceApiManagementLoggerRead(d, meta)
}

func resourceApiManagementLoggerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.LoggerClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	name := id.Path["loggers"]

	resp, err := client.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Logger %q (API Management Service %q / Resource Group %q) was not found - removing from state", name, serviceName, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Logger %q (API Management Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)
	d.Set("resource_id", resp.ResourceID)

	if properties := resp.LoggerContractProperties; properties != nil {
		d.Set("buffered", properties.IsBuffered)
		d.Set("description", properties.Description)
		if err := d.Set("eventhub", flattenApiManagementLoggerEventHub(d, properties)); err != nil {
			return fmt.Errorf("setting `eventhub`: %s", err)
		}
	}

	return nil
}

func resourceApiManagementLoggerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.LoggerClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	name := d.Get("name").(string)

	eventHubRaw, hasEventHub := d.GetOk("eventhub")
	appInsightsRaw, hasAppInsights := d.GetOk("application_insights")

	parameters := apimanagement.LoggerUpdateContract{
		LoggerUpdateParameters: &apimanagement.LoggerUpdateParameters{
			IsBuffered:  utils.Bool(d.Get("buffered").(bool)),
			Description: utils.String(d.Get("description").(string)),
		},
	}

	if hasEventHub {
		parameters.LoggerType = apimanagement.AzureEventHub
		parameters.Credentials = expandApiManagementLoggerEventHub(eventHubRaw.([]interface{}))
	} else if hasAppInsights {
		parameters.LoggerType = apimanagement.ApplicationInsights
		parameters.Credentials = expandApiManagementLoggerApplicationInsights(appInsightsRaw.([]interface{}))
	}

	if _, err := client.Update(ctx, resourceGroup, serviceName, name, parameters, ""); err != nil {
		return fmt.Errorf("updating Logger %q (Resource Group %q / API Management Service %q): %+v", name, resourceGroup, serviceName, err)
	}

	return resourceApiManagementLoggerRead(d, meta)
}

func resourceApiManagementLoggerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.LoggerClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	name := id.Path["loggers"]

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, name, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting Logger %q (Resource Group %q / API Management Service %q): %+v", name, resourceGroup, serviceName, err)
		}
	}

	return nil
}

func expandApiManagementLoggerEventHub(input []interface{}) map[string]*string {
	credentials := make(map[string]*string)
	eventHub := input[0].(map[string]interface{})
	credentials["name"] = utils.String(eventHub["name"].(string))
	credentials["connectionString"] = utils.String(eventHub["connection_string"].(string))
	return credentials
}

func expandApiManagementLoggerApplicationInsights(input []interface{}) map[string]*string {
	credentials := make(map[string]*string)
	ai := input[0].(map[string]interface{})
	credentials["instrumentationKey"] = utils.String(ai["instrumentation_key"].(string))
	return credentials
}

func flattenApiManagementLoggerEventHub(d *schema.ResourceData, properties *apimanagement.LoggerContractProperties) []interface{} {
	result := make([]interface{}, 0)
	if name := properties.Credentials["name"]; name != nil {
		eventHub := make(map[string]interface{})
		eventHub["name"] = *name
		if existing := d.Get("eventhub").([]interface{}); len(existing) > 0 {
			existingEventHub := existing[0].(map[string]interface{})
			if conn, ok := existingEventHub["connection_string"]; ok {
				eventHub["connection_string"] = conn.(string)
			}
		}
		result = append(result, eventHub)
	}
	return result
}
