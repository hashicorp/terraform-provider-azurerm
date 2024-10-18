// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/logger"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementLogger() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementLoggerCreate,
		Read:   resourceApiManagementLoggerRead,
		Update: resourceApiManagementLoggerUpdate,
		Delete: resourceApiManagementLoggerDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := logger.ParseLoggerID(id)
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

			"resource_group_name": commonschema.ResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"resource_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
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
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{
								"eventhub.0.connection_string",
								"eventhub.0.endpoint_uri",
							},
							ConflictsWith: []string{
								"eventhub.0.endpoint_uri",
								"eventhub.0.user_assigned_identity_client_id",
							},
						},
						"endpoint_uri": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{
								"eventhub.0.connection_string",
								"eventhub.0.endpoint_uri",
							},
							ConflictsWith: []string{
								"eventhub.0.connection_string",
							},
						},
						"user_assigned_identity_client_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
							ConflictsWith: []string{
								"eventhub.0.connection_string",
							},
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
						"connection_string": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{
								"application_insights.0.connection_string",
								"application_insights.0.instrumentation_key",
							},
							ConflictsWith: []string{
								"application_insights.0.instrumentation_key",
							},
						},
						"instrumentation_key": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{
								"application_insights.0.connection_string",
								"application_insights.0.instrumentation_key",
							},
							ConflictsWith: []string{
								"application_insights.0.connection_string",
							},
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

	id := logger.NewLoggerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	eventHubRaw := d.Get("eventhub").([]interface{})
	appInsightsRaw := d.Get("application_insights").([]interface{})

	if len(eventHubRaw) == 0 && len(appInsightsRaw) == 0 {
		return fmt.Errorf("Either `eventhub` or `application_insights` is required")
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_logger", id.ID())
		}
	}

	parameters := logger.LoggerContract{
		Properties: &logger.LoggerContractProperties{
			IsBuffered:  pointer.To(d.Get("buffered").(bool)),
			Description: pointer.To(d.Get("description").(string)),
		},
	}

	if len(eventHubRaw) > 0 {
		parameters.Properties.LoggerType = logger.LoggerTypeAzureEventHub
		credentials := expandApiManagementLoggerEventHub(eventHubRaw)
		parameters.Properties.Credentials = credentials
	} else if len(appInsightsRaw) > 0 {
		parameters.Properties.LoggerType = logger.LoggerTypeApplicationInsights
		parameters.Properties.Credentials = expandApiManagementLoggerApplicationInsights(appInsightsRaw)
	}

	if resourceId := d.Get("resource_id").(string); resourceId != "" {
		parameters.Properties.ResourceId = pointer.To(resourceId)
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, logger.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementLoggerRead(d, meta)
}

func resourceApiManagementLoggerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.LoggerClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := logger.ParseLoggerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("api_management_name", id.ServiceName)

	if model := resp.Model; model != nil {
		d.Set("name", pointer.From(model.Name))
		if props := model.Properties; props != nil {
			d.Set("resource_id", pointer.From(props.ResourceId))
			d.Set("buffered", pointer.From(props.IsBuffered))
			d.Set("description", pointer.From(props.Description))
			if err := d.Set("eventhub", flattenApiManagementLoggerEventHub(d, props)); err != nil {
				return fmt.Errorf("setting `eventhub`: %s", err)
			}
		}
	}

	return nil
}

func resourceApiManagementLoggerUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.LoggerClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := logger.NewLoggerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	eventHubRaw, hasEventHub := d.GetOk("eventhub")
	appInsightsRaw, hasAppInsights := d.GetOk("application_insights")

	parameters := logger.LoggerUpdateContract{
		Properties: &logger.LoggerUpdateParameters{
			IsBuffered:  pointer.To(d.Get("buffered").(bool)),
			Description: pointer.To(d.Get("description").(string)),
		},
	}

	if hasEventHub {
		parameters.Properties.LoggerType = pointer.To(logger.LoggerTypeAzureEventHub)
		credentials := expandApiManagementLoggerEventHub(eventHubRaw.([]interface{}))
		parameters.Properties.Credentials = credentials
	} else if hasAppInsights {
		parameters.Properties.LoggerType = pointer.To(logger.LoggerTypeApplicationInsights)
		parameters.Properties.Credentials = expandApiManagementLoggerApplicationInsights(appInsightsRaw.([]interface{}))
	}

	if _, err := client.Update(ctx, id, parameters, logger.UpdateOperationOptions{}); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceApiManagementLoggerRead(d, meta)
}

func resourceApiManagementLoggerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.LoggerClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := logger.ParseLoggerID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id, logger.DeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}

func expandApiManagementLoggerEventHub(input []interface{}) *map[string]string {
	credentials := make(map[string]string)
	eventHub := input[0].(map[string]interface{})

	connectionString := eventHub["connection_string"].(string)
	endpointAddress := eventHub["endpoint_uri"].(string)
	clientId := eventHub["user_assigned_identity_client_id"].(string)

	credentials["name"] = eventHub["name"].(string)
	if len(connectionString) > 0 {
		credentials["connectionString"] = connectionString
	} else if len(endpointAddress) > 0 {
		credentials["endpointAddress"] = endpointAddress
		// This field is required by the API and only accepts either a valid UUID or `SystemAssigned` as a value, so we default this to `SystemAssigned` in the create if the field is omitted
		credentials["identityClientId"] = "SystemAssigned"
		if clientId != "" {
			credentials["identityClientId"] = clientId
		}

	}

	return &credentials
}

func expandApiManagementLoggerApplicationInsights(input []interface{}) *map[string]string {
	credentials := make(map[string]string)
	ai := input[0].(map[string]interface{})
	if ai["instrumentation_key"].(string) != "" {
		credentials["instrumentationKey"] = ai["instrumentation_key"].(string)
	}
	if ai["connection_string"].(string) != "" {
		credentials["connectionString"] = ai["connection_string"].(string)
	}
	return &credentials
}

func flattenApiManagementLoggerEventHub(d *pluginsdk.ResourceData, properties *logger.LoggerContractProperties) []interface{} {
	result := make([]interface{}, 0)
	if c := properties.Credentials; c != nil && (*c)["name"] != "" {
		eventHub := make(map[string]interface{})
		eventHub["name"] = (*c)["name"]
		if existing := d.Get("eventhub").([]interface{}); len(existing) > 0 {
			existingEventHub := existing[0].(map[string]interface{})
			if conn, ok := existingEventHub["connection_string"]; ok {
				eventHub["connection_string"] = conn.(string)
			}
			if endpoint, ok := existingEventHub["endpoint_uri"]; ok {
				eventHub["endpoint_uri"] = endpoint
			}
			if clientId, ok := existingEventHub["user_assigned_identity_client_id"]; ok {
				if clientId != "SystemAssigned" {
					eventHub["user_assigned_identity_client_id"] = clientId
				}
			}
		}
		result = append(result, eventHub)
	}
	return result
}
