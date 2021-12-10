package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApiManagementLogger() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementLoggerCreate,
		Read:   resourceApiManagementLoggerRead,
		Update: resourceApiManagementLoggerUpdate,
		Delete: resourceApiManagementLoggerDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.LoggerID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": schemaz.SchemaApiManagementChildName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"resource_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"eventhub": {
				Type:          pluginsdk.TypeList,
				MaxItems:      1,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"application_insights"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.ValidateEventHubName(),
						},

						"connection_string": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"application_insights": {
				Type:          pluginsdk.TypeList,
				MaxItems:      1,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"eventhub"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"instrumentation_key": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"buffered": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceApiManagementLoggerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.LoggerClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewLoggerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	eventHubRaw := d.Get("eventhub").([]interface{})
	appInsightsRaw := d.Get("application_insights").([]interface{})

	if len(eventHubRaw) == 0 && len(appInsightsRaw) == 0 {
		return fmt.Errorf("Either `eventhub` or `application_insights` is required")
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
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

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, id.Name, parameters, ""); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementLoggerRead(d, meta)
}

func resourceApiManagementLoggerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.LoggerClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoggerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("api_management_name", id.ServiceName)
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

func resourceApiManagementLoggerUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.LoggerClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewLoggerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

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

	if _, err := client.Update(ctx, id.ResourceGroup, id.ServiceName, id.Name, parameters, ""); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceApiManagementLoggerRead(d, meta)
}

func resourceApiManagementLoggerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.LoggerClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoggerID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, id.ResourceGroup, id.ServiceName, id.Name, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
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

func flattenApiManagementLoggerEventHub(d *pluginsdk.ResourceData, properties *apimanagement.LoggerContractProperties) []interface{} {
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
