package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApiManagementLogger() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementLoggerCreate,
		Read:   resourceArmApiManagementLoggerRead,
		Update: resourceArmApiManagementLoggerUpdate,
		Delete: resourceArmApiManagementLoggerDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": azure.SchemaApiManagementChildName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": azure.SchemaApiManagementName(),

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
							ValidateFunc: azure.ValidateEventHubName(),
						},

						"connection_string": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validate.NoEmptyStrings,
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
							ValidateFunc: validate.NoEmptyStrings,
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

func resourceArmApiManagementLoggerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.LoggerClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)

	eventHubRaw := d.Get("eventhub").([]interface{})
	appInsightsRaw := d.Get("application_insights").([]interface{})

	if len(eventHubRaw) == 0 && len(appInsightsRaw) == 0 {
		return fmt.Errorf("Either `eventhub` or `application_insights` is required")
	}

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Logger %q (API Management Service %q / Resource Group %q): %s", name, serviceName, resourceGroup, err)
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
		parameters.Credentials = expandArmApiManagementLoggerEventHub(eventHubRaw)
	} else if len(appInsightsRaw) > 0 {
		parameters.LoggerType = apimanagement.ApplicationInsights
		parameters.Credentials = expandArmApiManagementLoggerApplicationInsights(appInsightsRaw)
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, name, parameters, ""); err != nil {
		return fmt.Errorf("Error creating Logger %q (Resource Group %q / API Management Service %q): %+v", name, resourceGroup, serviceName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Logger %q (Resource Group %q / API Management Service %q): %+v", name, resourceGroup, serviceName, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Logger %q (Resource Group %q / API Management Service %q) ID", name, resourceGroup, serviceName)
	}
	d.SetId(*resp.ID)

	return resourceArmApiManagementLoggerRead(d, meta)
}

func resourceArmApiManagementLoggerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.LoggerClient
	ctx := meta.(*ArmClient).StopContext

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
		return fmt.Errorf("Error reading Logger %q (API Management Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)

	if properties := resp.LoggerContractProperties; properties != nil {
		d.Set("buffered", properties.IsBuffered)
		d.Set("description", properties.Description)
		if err := d.Set("eventhub", flattenArmApiManagementLoggerEventHub(d, properties)); err != nil {
			return fmt.Errorf("Error setting `eventhub`: %s", err)
		}
	}

	return nil
}

func resourceArmApiManagementLoggerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.LoggerClient
	ctx := meta.(*ArmClient).StopContext

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
		parameters.Credentials = expandArmApiManagementLoggerEventHub(eventHubRaw.([]interface{}))
	} else if hasAppInsights {
		parameters.LoggerType = apimanagement.ApplicationInsights
		parameters.Credentials = expandArmApiManagementLoggerApplicationInsights(appInsightsRaw.([]interface{}))
	}

	if _, err := client.Update(ctx, resourceGroup, serviceName, name, parameters, ""); err != nil {
		return fmt.Errorf("Error updating Logger %q (Resource Group %q / API Management Service %q): %+v", name, resourceGroup, serviceName, err)
	}

	return resourceArmApiManagementLoggerRead(d, meta)
}

func resourceArmApiManagementLoggerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.LoggerClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	name := id.Path["loggers"]

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, name, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting Logger %q (Resource Group %q / API Management Service %q): %+v", name, resourceGroup, serviceName, err)
		}
	}

	return nil
}

func expandArmApiManagementLoggerEventHub(input []interface{}) map[string]*string {
	credentials := make(map[string]*string)
	eventHub := input[0].(map[string]interface{})
	credentials["name"] = utils.String(eventHub["name"].(string))
	credentials["connectionString"] = utils.String(eventHub["connection_string"].(string))
	return credentials
}

func expandArmApiManagementLoggerApplicationInsights(input []interface{}) map[string]*string {
	credentials := make(map[string]*string)
	ai := input[0].(map[string]interface{})
	credentials["instrumentationKey"] = utils.String(ai["instrumentation_key"].(string))
	return credentials
}

func flattenArmApiManagementLoggerEventHub(d *schema.ResourceData, properties *apimanagement.LoggerContractProperties) []interface{} {
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
