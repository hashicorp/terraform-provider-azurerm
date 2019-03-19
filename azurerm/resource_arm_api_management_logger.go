package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
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

			"resource_group_name": resourceGroupNameSchema(),

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

		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			_, hasEventHub := d.GetOk("eventhub")
			_, hasAppInsights := d.GetOk("application_insights")
			if !hasEventHub && !hasAppInsights {
				return fmt.Errorf("Either `eventhub` or `application_insights` is required")
			}
			return nil
		},
	}
}

func resourceArmApiManagementLoggerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagementLoggerClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)

	eventHubRaw, hasEventHub := d.GetOk("eventhub")
	appInsightsRaw, hasAppInsights := d.GetOk("application_insights")

	parameters := apimanagement.LoggerContract{
		LoggerContractProperties: &apimanagement.LoggerContractProperties{
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
	client := meta.(*ArmClient).apiManagementLoggerClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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
		return fmt.Errorf("Error reading Logger %q (Resource Group %q / API Management Name %q): %+v", name, resourceGroup, serviceName, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)

	if properties := resp.LoggerContractProperties; properties != nil {
		d.Set("buffered", properties.IsBuffered)
		d.Set("description", properties.Description)

		if properties.LoggerType == apimanagement.AzureEventHub {
			if err := d.Set("eventhub", flattenArmApiManagementLoggerEventHub(d, properties.Credentials)); err != nil {
				return fmt.Errorf("Error setting `eventhub` for Logger %q (Resource Group %q / API Management Name %q): %s", name, resourceGroup, serviceName, err)
			}
		} else if properties.LoggerType == apimanagement.ApplicationInsights {
			if err := d.Set("application_insights", flattenArmApiManagementLoggerApplicationInsights(d, properties.Credentials)); err != nil {
				return fmt.Errorf("Error setting `application_insights` for Logger %q (Resource Group %q / API Management Name %q): %s", name, resourceGroup, serviceName, err)
			}
		}
	}

	return nil
}

func resourceArmApiManagementLoggerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagementLoggerClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing API Management Logger ID %q: %+v", d.Id(), err)
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	name := id.Path["loggers"]

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
	client := meta.(*ArmClient).apiManagementLoggerClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing API Management Logger ID %q: %+v", d.Id(), err)
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

func flattenArmApiManagementLoggerEventHub(d *schema.ResourceData, credentials map[string]*string) []interface{} {
	eventHub := make(map[string]interface{})
	if name := credentials["name"]; name != nil {
		eventHub["name"] = *name
	}
	if conn, ok := d.GetOk("eventhub.0.connection_string"); ok {
		eventHub["connection_string"] = conn.(string)
	}
	return []interface{}{eventHub}
}

func flattenArmApiManagementLoggerApplicationInsights(d *schema.ResourceData, credentials map[string]*string) []interface{} {
	appInsights := make(map[string]interface{})
	if conn, ok := d.GetOk("application_insights.0.instrumentation_key"); ok {
		appInsights["instrumentation_key"] = conn.(string)
	}
	return []interface{}{appInsights}
}
