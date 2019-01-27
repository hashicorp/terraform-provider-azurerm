package azurerm

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"strings"

	"github.com/Azure/azure-sdk-for-go/services/eventhub/mgmt/2017-04-01/eventhub"
	"github.com/Azure/azure-sdk-for-go/services/iothub/mgmt/2018-04-01/devices"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmIotHub() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmIotHubCreateUpdate,
		Read:   resourceArmIotHubRead,
		Update: resourceArmIotHubCreateUpdate,
		Delete: resourceArmIotHubDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"sku": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(devices.B1),
								string(devices.B2),
								string(devices.B3),
								string(devices.F1),
								string(devices.S1),
								string(devices.S2),
								string(devices.S3),
							}, true),
						},

						"tier": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(devices.Basic),
								string(devices.Free),
								string(devices.Standard),
							}, true),
						},

						"capacity": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
					},
				},
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"event_hub_events_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"event_hub_operations_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"event_hub_events_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"event_hub_operations_path": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"shared_access_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"primary_key": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"secondary_key": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"permissions": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"endpoint": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"AzureIotHub.StorageContainer",
								"AzureIotHub.ServiceBusQueue",
								"AzureIotHub.ServiceBusTopic",
								"AzureIotHub.EventHub",
							}, false),
						},
						"connection_string": {
							Type:     schema.TypeString,
							Required: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								// As Azure API masks the connection string key suppress diff for this property
								if old != "" && strings.HasSuffix(old, "****") {
									return true
								}

								return false
							},
						},
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIoTHubEndpointName,
						},
						"batch_frequency_in_seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      300,
							ValidateFunc: validation.IntBetween(60, 720),
						},
						"max_chunk_size_in_bytes": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      314572800,
							ValidateFunc: validation.IntBetween(10485760, 524288000),
						},
						"container_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"encoding": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(eventhub.Avro),
								string(eventhub.AvroDeflate),
							}, true),
						},
						"file_name_format": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateIoTHubFileNameFormat,
						},
					},
				},
			},

			"route": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(0, 64),
						},
						"source": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"DeviceJobLifecycleEvents",
								"DeviceLifecycleEvents",
								"DeviceMessages",
								"Invalid",
								"TwinChangeEvents",
							}, false),
						},
						"condition": {
							// The condition is a string value representing device-to-cloud message routes query expression
							// https://docs.microsoft.com/en-us/azure/iot-hub/iot-hub-devguide-query-language#device-to-cloud-message-routes-query-expressions
							Type:     schema.TypeString,
							Optional: true,
							Default:  "true",
						},
						"endpoint_names": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},

			"fallback_route": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							Default:      "$fallback",
							ValidateFunc: validation.StringLenBetween(0, 64),
						},
						"source": {
							Type:     schema.TypeString,
							Required: true,
							Default:  "DeviceMessages",
							ValidateFunc: validation.StringInSlice([]string{
								"DeviceJobLifecycleEvents",
								"DeviceLifecycleEvents",
								"DeviceMessages",
								"Invalid",
								"TwinChangeEvents",
							}, false),
						},
						"condition": {
							// The condition is a string value representing device-to-cloud message routes query expression
							// https://docs.microsoft.com/en-us/azure/iot-hub/iot-hub-devguide-query-language#device-to-cloud-message-routes-query-expressions
							Type:     schema.TypeString,
							Optional: true,
							Default:  "true",
						},
						"endpoint_names": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							MaxItems: 1,
							Required: true,
							Default:  &[]string{"events"},
						},
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}

}

func resourceArmIotHubCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).iothubResourceClient
	ctx := meta.(*ArmClient).StopContext
	subscriptionID := meta.(*ArmClient).subscriptionId

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing IoTHub %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_iothub", *existing.ID)
		}
	}

	res, err := client.CheckNameAvailability(ctx, devices.OperationInputs{
		Name: &name,
	})
	if err != nil {
		return fmt.Errorf("An error occurred checking if the IoTHub name was unique: %+v", err)
	}

	if !*res.NameAvailable {
		_, err = client.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("An IoTHub already exists with the name %q - please choose an alternate name: %s", name, string(res.Reason))
		}
	}

	location := azureRMNormalizeLocation(d.Get("location").(string))
	skuInfo := expandIoTHubSku(d)
	tags := d.Get("tags").(map[string]interface{})
	fallbackRoute := expandIoTHubFallbackRoute(d)

	endpoints, err := expandIoTHubEndpoints(d, subscriptionID)
	if err != nil {
		return fmt.Errorf("Error expanding `endpoint`: %+v", err)
	}

	routes := expandIoTHubRoutes(d)

	properties := devices.IotHubDescription{
		Name:     utils.String(name),
		Location: utils.String(location),
		Sku:      skuInfo,
		Properties: &devices.IotHubProperties{
			Routing: &devices.RoutingProperties{
				Endpoints:     endpoints,
				Routes:        routes,
				FallbackRoute: fallbackRoute,
			},
		},
		Tags: expandTags(tags),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, properties, "")
	if err != nil {
		return fmt.Errorf("Error creating/updating IotHub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the completion of the creating/updating of IotHub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	d.SetId(*resp.ID)
	return resourceArmIotHubRead(d, meta)
}

func resourceArmIotHubRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).iothubResourceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["IotHubs"]
	resourceGroup := id.ResourceGroup
	hub, err := client.Get(ctx, id.ResourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(hub.Response) {
			log.Printf("[DEBUG] IoTHub %q (Resource Group %q) was not found!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving IotHub Client %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	keysResp, err := client.ListKeys(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error listing keys for IoTHub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	keyList := keysResp.Response()
	keys := flattenIoTHubSharedAccessPolicy(keyList.Value)

	if err := d.Set("shared_access_policy", keys); err != nil {
		return fmt.Errorf("Error setting `shared_access_policy` in IoTHub %q: %+v", name, err)
	}

	if properties := hub.Properties; properties != nil {

		for k, v := range properties.EventHubEndpoints {
			if v == nil {
				continue
			}

			if k == "events" {
				d.Set("event_hub_events_endpoint", v.Endpoint)
				d.Set("event_hub_events_path", v.Path)
			} else if k == "operationsMonitoringEvents" {
				d.Set("event_hub_operations_endpoint", v.Endpoint)
				d.Set("event_hub_operations_path", v.Path)
			}

		}

		d.Set("hostname", properties.HostName)

		endpoints := flattenIoTHubEndpoint(properties.Routing)
		if err := d.Set("endpoint", endpoints); err != nil {
			return fmt.Errorf("Error setting `endpoint` in IoTHub %q: %+v", name, err)
		}

		routes := flattenIoTHubRoute(properties.Routing)
		if err := d.Set("route", routes); err != nil {
			return fmt.Errorf("Error setting `route` in IoTHub %q: %+v", name, err)
		}

		fallbackRoute := flattenIoTHubFallbackRoute(properties.Routing)
		if err := d.Set("fallback_route", fallbackRoute); err != nil {
			return fmt.Errorf("Error setting `fallbackRoute` in IoTHub %q: %+v", name, err)
		}
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := hub.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}
	sku := flattenIoTHubSku(hub.Sku)
	if err := d.Set("sku", sku); err != nil {
		return fmt.Errorf("Error setting `sku`: %+v", err)
	}
	d.Set("type", hub.Type)
	flattenAndSetTags(d, hub.Tags)

	return nil
}

func resourceArmIotHubDelete(d *schema.ResourceData, meta interface{}) error {
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*ArmClient).iothubResourceClient
	ctx := meta.(*ArmClient).StopContext

	name := id.Path["IotHubs"]
	resourceGroup := id.ResourceGroup

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return err
	}

	return waitForIotHubToBeDeleted(ctx, client, resourceGroup, name)
}

func waitForIotHubToBeDeleted(ctx context.Context, client devices.IotHubResourceClient, resourceGroup, name string) error {
	// we can't use the Waiter here since the API returns a 404 once it's deleted which is considered a polling status code..
	log.Printf("[DEBUG] Waiting for IotHub (%q in Resource Group %q) to be deleted", name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"200"},
		Target:  []string{"404"},
		Refresh: iothubStateStatusCodeRefreshFunc(ctx, client, resourceGroup, name),
		Timeout: 40 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for IotHub (%q in Resource Group %q) to be deleted: %+v", name, resourceGroup, err)
	}

	return nil
}

func iothubStateStatusCodeRefreshFunc(ctx context.Context, client devices.IotHubResourceClient, resourceGroup, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroup, name)

		log.Printf("Retrieving IoTHub %q (Resource Group %q) returned Status %d", resourceGroup, name, res.StatusCode)

		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return res, strconv.Itoa(res.StatusCode), nil
			}
			return nil, "", fmt.Errorf("Error polling for the status of the IotHub %q (RG: %q): %+v", name, resourceGroup, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}

func expandIoTHubRoutes(d *schema.ResourceData) *[]devices.RouteProperties {
	routeList := d.Get("route").([]interface{})

	routeProperties := make([]devices.RouteProperties, 0)

	for _, routeRaw := range routeList {
		route := routeRaw.(map[string]interface{})

		name := route["name"].(string)
		source := devices.RoutingSource(route["source"].(string))
		condition := route["condition"].(string)

		endpointNamesRaw := route["endpoint_names"].([]interface{})
		endpointsNames := make([]string, 0)
		for _, n := range endpointNamesRaw {
			endpointsNames = append(endpointsNames, n.(string))
		}

		isEnabled := route["enabled"].(bool)

		routeProperties = append(routeProperties, devices.RouteProperties{
			Name:          &name,
			Source:        source,
			Condition:     &condition,
			EndpointNames: &endpointsNames,
			IsEnabled:     &isEnabled,
		})
	}

	return &routeProperties
}

func expandIoTHubEndpoints(d *schema.ResourceData, subscriptionId string) (*devices.RoutingEndpoints, error) {
	routeEndpointList := d.Get("endpoint").([]interface{})

	serviceBusQueueEndpointProperties := make([]devices.RoutingServiceBusQueueEndpointProperties, 0)
	serviceBusTopicEndpointProperties := make([]devices.RoutingServiceBusTopicEndpointProperties, 0)
	eventHubProperties := make([]devices.RoutingEventHubProperties, 0)
	storageContainerProperties := make([]devices.RoutingStorageContainerProperties, 0)

	for _, endpointRaw := range routeEndpointList {
		endpoint := endpointRaw.(map[string]interface{})

		t := endpoint["type"]
		connectionStr := endpoint["connection_string"].(string)
		name := endpoint["name"].(string)
		subscriptionID := subscriptionId
		resourceGroup := d.Get("resource_group_name").(string)

		switch t {
		case "AzureIotHub.StorageContainer":
			containerName := endpoint["container_name"].(string)
			fileNameFormat := endpoint["file_name_format"].(string)
			batchFrequencyInSeconds := int32(endpoint["batch_frequency_in_seconds"].(int))
			maxChunkSizeInBytes := int32(endpoint["max_chunk_size_in_bytes"].(int))
			encoding := endpoint["encoding"].(string)

			storageContainer := devices.RoutingStorageContainerProperties{
				ConnectionString:        &connectionStr,
				Name:                    &name,
				SubscriptionID:          &subscriptionID,
				ResourceGroup:           &resourceGroup,
				ContainerName:           &containerName,
				FileNameFormat:          &fileNameFormat,
				BatchFrequencyInSeconds: &batchFrequencyInSeconds,
				MaxChunkSizeInBytes:     &maxChunkSizeInBytes,
				Encoding:                &encoding,
			}
			storageContainerProperties = append(storageContainerProperties, storageContainer)

		case "AzureIotHub.ServiceBusQueue":
			sbQueue := devices.RoutingServiceBusQueueEndpointProperties{
				ConnectionString: &connectionStr,
				Name:             &name,
				SubscriptionID:   &subscriptionID,
				ResourceGroup:    &resourceGroup,
			}
			serviceBusQueueEndpointProperties = append(serviceBusQueueEndpointProperties, sbQueue)

		case "AzureIotHub.ServiceBusTopic":
			sbTopic := devices.RoutingServiceBusTopicEndpointProperties{
				ConnectionString: &connectionStr,
				Name:             &name,
				SubscriptionID:   &subscriptionID,
				ResourceGroup:    &resourceGroup,
			}
			serviceBusTopicEndpointProperties = append(serviceBusTopicEndpointProperties, sbTopic)

		case "AzureIotHub.EventHub":
			eventHub := devices.RoutingEventHubProperties{
				ConnectionString: &connectionStr,
				Name:             &name,
				SubscriptionID:   &subscriptionID,
				ResourceGroup:    &resourceGroup,
			}
			eventHubProperties = append(eventHubProperties, eventHub)
		}
	}

	return &devices.RoutingEndpoints{
		ServiceBusQueues:  &serviceBusQueueEndpointProperties,
		ServiceBusTopics:  &serviceBusTopicEndpointProperties,
		EventHubs:         &eventHubProperties,
		StorageContainers: &storageContainerProperties,
	}, nil
}

func expandIoTHubFallbackRoute(d *schema.ResourceData) *devices.FallbackRouteProperties {
	fallbackRouteList := d.Get("fallback_route").([]interface{})
	fallbackRouteMap := fallbackRouteList[0].(map[string]interface{})

	name := fallbackRouteMap["name"].(string)
	source := fallbackRouteMap["source"].(string)
	endpointNamesRaw := route["endpoint_names"].([]interface{})
	endpointsNames := make([]string, 0)
	for _, n := range endpointNamesRaw {
		endpointsNames = append(endpointsNames, n.(string))
	}
	isEnabled := fallbackRouteMap["enabled"].(bool)

	return &devices.FallbackRouteProperties{
		Name:          &name,
		Source:        &source,
		EndpointNames: &endpointNames,
		IsEnabled:     &isEnabled,
	}
}

func expandIoTHubSku(d *schema.ResourceData) *devices.IotHubSkuInfo {
	skuList := d.Get("sku").([]interface{})
	skuMap := skuList[0].(map[string]interface{})
	capacity := int64(skuMap["capacity"].(int))

	name := skuMap["name"].(string)
	tier := skuMap["tier"].(string)

	return &devices.IotHubSkuInfo{
		Name:     devices.IotHubSku(name),
		Tier:     devices.IotHubSkuTier(tier),
		Capacity: utils.Int64(capacity),
	}
}

func flattenIoTHubSku(input *devices.IotHubSkuInfo) []interface{} {
	output := make(map[string]interface{})

	output["name"] = string(input.Name)
	output["tier"] = string(input.Tier)
	if capacity := input.Capacity; capacity != nil {
		output["capacity"] = int(*capacity)
	}

	return []interface{}{output}
}

func flattenIoTHubSharedAccessPolicy(input *[]devices.SharedAccessSignatureAuthorizationRule) []interface{} {
	results := make([]interface{}, 0)

	if keys := input; keys != nil {
		for _, key := range *keys {
			keyMap := make(map[string]interface{})

			if keyName := key.KeyName; keyName != nil {
				keyMap["key_name"] = *keyName
			}

			if primaryKey := key.PrimaryKey; primaryKey != nil {
				keyMap["primary_key"] = *primaryKey
			}

			if secondaryKey := key.SecondaryKey; secondaryKey != nil {
				keyMap["secondary_key"] = *secondaryKey
			}

			keyMap["permissions"] = string(key.Rights)
			results = append(results, keyMap)
		}
	}

	return results
}

func flattenIoTHubEndpoint(input *devices.RoutingProperties) []interface{} {
	results := make([]interface{}, 0)

	if input != nil && input.Endpoints != nil {

		if containers := input.Endpoints.StorageContainers; containers != nil {
			for _, container := range *containers {
				output := make(map[string]interface{})

				if connString := container.ConnectionString; connString != nil {
					output["connection_string"] = *connString
				}
				if name := container.Name; name != nil {
					output["name"] = *name
				}
				if containerName := container.ContainerName; containerName != nil {
					output["container_name"] = *containerName
				}
				if fileNameFmt := container.FileNameFormat; fileNameFmt != nil {
					output["file_name_format"] = *fileNameFmt
				}
				if batchFreq := container.BatchFrequencyInSeconds; batchFreq != nil {
					output["batch_frequency_in_seconds"] = *batchFreq
				}
				if chunkSize := container.MaxChunkSizeInBytes; chunkSize != nil {
					output["max_chunk_size_in_bytes"] = *chunkSize
				}
				if encoding := container.Encoding; encoding != nil {
					output["encoding"] = *encoding
				}
				output["type"] = "AzureIotHub.StorageContainer"

				results = append(results, output)
			}
		}

		if queues := input.Endpoints.ServiceBusQueues; queues != nil {
			for _, queue := range *queues {
				output := make(map[string]interface{})

				if connString := queue.ConnectionString; connString != nil {
					output["connection_string"] = *connString
				}
				if name := queue.Name; name != nil {
					output["name"] = *name
				}

				results = append(results, output)
			}
		}

		if topics := input.Endpoints.ServiceBusTopics; topics != nil {
			for _, topic := range *topics {
				output := make(map[string]interface{})

				if connString := topic.ConnectionString; connString != nil {
					output["connection_string"] = *connString
				}
				if name := topic.Name; name != nil {
					output["name"] = *name
				}

				results = append(results, output)
			}
		}

		if eventHubs := input.Endpoints.EventHubs; eventHubs != nil {
			for _, eventHub := range *eventHubs {
				output := make(map[string]interface{})

				if connString := eventHub.ConnectionString; connString != nil {
					output["connection_string"] = *connString
				}
				if name := eventHub.Name; name != nil {
					output["name"] = *name
				}

				results = append(results, output)
			}
		}
	}

	return results
}

func flattenIoTHubRoute(input *devices.RoutingProperties) []interface{} {
	results := make([]interface{}, 0)

	if input != nil && input.Routes != nil {
		for _, route := range *input.Routes {
			output := make(map[string]interface{})

			if name := route.Name; name != nil {
				output["name"] = *name
			}
			if condition := route.Condition; condition != nil {
				output["condition"] = *condition
			}
			if endpointNames := route.EndpointNames; endpointNames != nil {
				output["endpoint_names"] = *endpointNames
			}
			if isEnabled := route.IsEnabled; isEnabled != nil {
				output["enabled"] = *isEnabled
			}
			output["source"] = route.Source

			results = append(results, output)
		}
	}

	return results
}

func flattenIoTHubFallbackRoute(input *devices.FallbackRouteProperties) []interface{} {
	output := make(map[string]interface{})

	if name := input.Name; name != nil {
		output["name"] = *name
	}
	if condition := input.Condition; condition != nil {
		output["condition"] = *condition
	}
	if endpointNames := input.EndpointNames; endpointNames != nil {
		output["endpoint_names"] = *endpointNames
	}
	if isEnabled := input.IsEnabled; isEnabled != nil {
		output["enabled"] = *isEnabled
	}
	output["source"] = input.Source

	return []interface{}{output}
}

func validateIoTHubEndpointName(v interface{}, _ string) (warnings []string, errors []error) {
	value := v.(string)

	reservedNames := []string{
		"events",
		"operationsMonitoringEvents",
		"fileNotifications",
		"$default",
	}

	for _, name := range reservedNames {
		if name == value {
			errors = append(errors, fmt.Errorf("The reserved endpoint name %s could not be used as a name for a custom endpoint", name))
		}
	}

	return warnings, errors
}

func validateIoTHubFileNameFormat(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	requiredComponents := []string{
		"{iothub}",
		"{partition}",
		"{YYYY}",
		"{MM}",
		"{DD}",
		"{HH}",
		"{mm}",
	}

	for _, component := range requiredComponents {
		if !strings.Contains(value, component) {
			errors = append(errors, fmt.Errorf("%s needs to contain %q", k, component))
		}
	}

	return warnings, errors
}
