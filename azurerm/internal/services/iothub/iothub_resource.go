package iothub

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/iothub/mgmt/2019-03-22-preview/devices"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// TODO: outside of this pr make this private

var IothubResourceName = "azurerm_iothub"

// nolint unparam
func suppressIfTypeIsNot(t string) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		path := strings.Split(k, ".")
		path[len(path)-1] = "type"
		return d.Get(strings.Join(path, ".")).(string) != t
	}
}

// nolint unparam
func supressWhenAll(fs ...schema.SchemaDiffSuppressFunc) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		for _, f := range fs {
			if !f(k, old, new, d) {
				return false
			}
		}
		return true
	}
}

func resourceArmIotHub() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmIotHubCreateUpdate,
		Read:   resourceArmIotHubRead,
		Update: resourceArmIotHubCreateUpdate,
		Delete: resourceArmIotHubDelete,
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
				ValidateFunc: validate.IoTHubName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

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
							}, false),
						},

						"capacity": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 200),
						},
					},
				},
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

			"event_hub_partition_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(2, 128),
			},
			"event_hub_retention_in_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 7),
			},

			"file_upload": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection_string": {
							Type:     schema.TypeString,
							Required: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								secretKeyRegex := regexp.MustCompile("(SharedAccessKey|AccountKey)=[^;]+")
								sbProtocolRegex := regexp.MustCompile("sb://([^:]+)(:5671)?/;")

								// Azure will always mask the Access Keys and will include the port number in the GET response
								// 5671 is the default port for Azure Service Bus connections
								maskedNew := sbProtocolRegex.ReplaceAllString(new, "sb://$1:5671/;")
								maskedNew = secretKeyRegex.ReplaceAllString(maskedNew, "$1=****")
								return (new == d.Get(k).(string)) && (maskedNew == old)
							},
							Sensitive: true,
						},
						"container_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"notifications": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"max_delivery_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      10,
							ValidateFunc: validation.IntBetween(1, 100),
						},
						"sas_ttl": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.ISO8601Duration,
						},
						"default_ttl": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.ISO8601Duration,
						},
						"lock_duration": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.ISO8601Duration,
						},
					},
				},
			},

			"endpoint": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
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
								secretKeyRegex := regexp.MustCompile("(SharedAccessKey|AccountKey)=[^;]+")
								sbProtocolRegex := regexp.MustCompile("sb://([^:]+)(:5671)?/;")

								// Azure will always mask the Access Keys and will include the port number in the GET response
								// 5671 is the default port for Azure Service Bus connections
								maskedNew := sbProtocolRegex.ReplaceAllString(new, "sb://$1:5671/;")
								maskedNew = secretKeyRegex.ReplaceAllString(maskedNew, "$1=****")
								return (new == d.Get(k).(string)) && (maskedNew == old)
							},
							Sensitive: true,
						},
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.IoTHubEndpointName,
						},
						"batch_frequency_in_seconds": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          300,
							DiffSuppressFunc: suppressIfTypeIsNot("AzureIotHub.StorageContainer"),
							ValidateFunc:     validation.IntBetween(60, 720),
						},
						"max_chunk_size_in_bytes": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          314572800,
							DiffSuppressFunc: suppressIfTypeIsNot("AzureIotHub.StorageContainer"),
							ValidateFunc:     validation.IntBetween(10485760, 524288000),
						},
						"container_name": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppressIfTypeIsNot("AzureIotHub.StorageContainer"),
						},
						"encoding": {
							Type:     schema.TypeString,
							Optional: true,
							DiffSuppressFunc: supressWhenAll(
								suppressIfTypeIsNot("AzureIotHub.StorageContainer"),
								suppress.CaseDifference),
							ValidateFunc: validation.StringInSlice([]string{
								string(devices.Avro),
								string(devices.AvroDeflate),
								string(devices.JSON),
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
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile("^[-_.a-zA-Z0-9]{1,64}$"),
								"Route Name name can only include alphanumeric characters, periods, underscores, hyphens, has a maximum length of 64 characters, and must be unique.",
							),
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
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:     schema.TypeString,
							Optional: true,
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
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringLenBetween(0, 64),
							},
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"ip_filter_rule": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"ip_mask": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.CIDR,
						},
						"action": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(devices.Accept),
								string(devices.Reject),
							}, false),
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

			"tags": tags.Schema(),
		},
	}
}

func resourceArmIotHubCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	locks.ByName(name, IothubResourceName)
	defer locks.UnlockByName(name, IothubResourceName)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
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

	routingProperties := devices.RoutingProperties{}

	if _, ok := d.GetOk("route"); ok {
		routingProperties.Routes = expandIoTHubRoutes(d)
	}

	if _, ok := d.GetOk("fallback_route"); ok {
		routingProperties.FallbackRoute = expandIoTHubFallbackRoute(d)
	}

	if _, ok := d.GetOk("endpoint"); ok {
		routingProperties.Endpoints = expandIoTHubEndpoints(d, subscriptionID)
	}

	storageEndpoints, messagingEndpoints, enableFileUploadNotifications := expandIoTHubFileUpload(d)
	if err != nil {
		return fmt.Errorf("Error expanding `file_upload`: %+v", err)
	}

	props := devices.IotHubDescription{
		Name:     utils.String(name),
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Sku:      expandIoTHubSku(d),
		Properties: &devices.IotHubProperties{
			IPFilterRules:                 expandIPFilterRules(d),
			Routing:                       &routingProperties,
			StorageEndpoints:              storageEndpoints,
			MessagingEndpoints:            messagingEndpoints,
			EnableFileUploadNotifications: &enableFileUploadNotifications,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	retention, retentionOk := d.GetOk("event_hub_retention_in_days")
	partition, partitionOk := d.GetOk("event_hub_partition_count")
	if partitionOk || retentionOk {
		eh := devices.EventHubProperties{}
		if retentionOk {
			eh.RetentionTimeInDays = utils.Int64(int64(retention.(int)))
		}
		if partitionOk {
			eh.PartitionCount = utils.Int32(int32(partition.(int)))
		}

		props.Properties.EventHubEndpoints = map[string]*devices.EventHubProperties{
			"events": &eh,
		}
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, props, "")
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
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
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
				d.Set("event_hub_partition_count", v.PartitionCount)
				d.Set("event_hub_retention_in_days", v.RetentionTimeInDays)
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

		ipFilterRules := flattenIPFilterRules(properties.IPFilterRules)
		if err := d.Set("ip_filter_rule", ipFilterRules); err != nil {
			return fmt.Errorf("Error setting `ip_filter_rule` in IoTHub %q: %+v", name, err)
		}

		fileUpload := flattenIoTHubFileUpload(properties.StorageEndpoints, properties.MessagingEndpoints, properties.EnableFileUploadNotifications)
		if err := d.Set("file_upload", fileUpload); err != nil {
			return fmt.Errorf("Error setting `file_upload` in IoTHub %q: %+v", name, err)
		}
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := hub.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	sku := flattenIoTHubSku(hub.Sku)
	if err := d.Set("sku", sku); err != nil {
		return fmt.Errorf("Error setting `sku`: %+v", err)
	}
	d.Set("type", hub.Type)
	return tags.FlattenAndSet(d, hub.Tags)
}

func resourceArmIotHubDelete(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := id.Path["IotHubs"]
	resourceGroup := id.ResourceGroup

	locks.ByName(name, IothubResourceName)
	defer locks.UnlockByName(name, IothubResourceName)

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return err
	}

	return waitForIotHubToBeDeleted(ctx, client, resourceGroup, name, d)
}

func waitForIotHubToBeDeleted(ctx context.Context, client *devices.IotHubResourceClient, resourceGroup, name string, d *schema.ResourceData) error {
	// we can't use the Waiter here since the API returns a 404 once it's deleted which is considered a polling status code..
	log.Printf("[DEBUG] Waiting for IotHub (%q in Resource Group %q) to be deleted", name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"200"},
		Target:  []string{"404"},
		Refresh: iothubStateStatusCodeRefreshFunc(ctx, client, resourceGroup, name),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for IotHub (%q in Resource Group %q) to be deleted: %+v", name, resourceGroup, err)
	}

	return nil
}

func iothubStateStatusCodeRefreshFunc(ctx context.Context, client *devices.IotHubResourceClient, resourceGroup, name string) resource.StateRefreshFunc {
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

		isEnabled := route["enabled"].(bool)

		routeProperties = append(routeProperties, devices.RouteProperties{
			Name:          &name,
			Source:        source,
			Condition:     &condition,
			EndpointNames: utils.ExpandStringSlice(endpointNamesRaw),
			IsEnabled:     &isEnabled,
		})
	}

	return &routeProperties
}

func expandIoTHubFileUpload(d *schema.ResourceData) (map[string]*devices.StorageEndpointProperties, map[string]*devices.MessagingEndpointProperties, bool) {
	fileUploadList := d.Get("file_upload").([]interface{})

	storageEndpointProperties := make(map[string]*devices.StorageEndpointProperties)
	messagingEndpointProperties := make(map[string]*devices.MessagingEndpointProperties)
	notifications := false

	if len(fileUploadList) > 0 {
		fileUploadMap := fileUploadList[0].(map[string]interface{})

		connectionStr := fileUploadMap["connection_string"].(string)
		containerName := fileUploadMap["container_name"].(string)
		notifications = fileUploadMap["notifications"].(bool)
		maxDeliveryCount := int32(fileUploadMap["max_delivery_count"].(int))
		sasTTL := fileUploadMap["sas_ttl"].(string)
		defaultTTL := fileUploadMap["default_ttl"].(string)
		lockDuration := fileUploadMap["lock_duration"].(string)

		storageEndpointProperties["$default"] = &devices.StorageEndpointProperties{
			SasTTLAsIso8601:  &sasTTL,
			ConnectionString: &connectionStr,
			ContainerName:    &containerName,
		}

		messagingEndpointProperties["fileNotifications"] = &devices.MessagingEndpointProperties{
			LockDurationAsIso8601: &lockDuration,
			TTLAsIso8601:          &defaultTTL,
			MaxDeliveryCount:      &maxDeliveryCount,
		}
	}

	return storageEndpointProperties, messagingEndpointProperties, notifications
}

func expandIoTHubEndpoints(d *schema.ResourceData, subscriptionId string) *devices.RoutingEndpoints {
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
				Encoding:                devices.Encoding(encoding),
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
	}
}

func expandIoTHubFallbackRoute(d *schema.ResourceData) *devices.FallbackRouteProperties {
	fallbackRouteList := d.Get("fallback_route").([]interface{})
	if len(fallbackRouteList) == 0 {
		return nil
	}

	fallbackRouteMap := fallbackRouteList[0].(map[string]interface{})

	source := fallbackRouteMap["source"].(string)
	condition := fallbackRouteMap["condition"].(string)
	isEnabled := fallbackRouteMap["enabled"].(bool)

	return &devices.FallbackRouteProperties{
		Source:        &source,
		Condition:     &condition,
		EndpointNames: utils.ExpandStringSlice(fallbackRouteMap["endpoint_names"].([]interface{})),
		IsEnabled:     &isEnabled,
	}
}

func expandIoTHubSku(d *schema.ResourceData) *devices.IotHubSkuInfo {
	skuList := d.Get("sku").([]interface{})
	skuMap := skuList[0].(map[string]interface{})

	return &devices.IotHubSkuInfo{
		Name:     devices.IotHubSku(skuMap["name"].(string)),
		Capacity: utils.Int64(int64(skuMap["capacity"].(int))),
	}
}

func flattenIoTHubSku(input *devices.IotHubSkuInfo) []interface{} {
	output := make(map[string]interface{})

	output["name"] = string(input.Name)
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

func flattenIoTHubFileUpload(storageEndpoints map[string]*devices.StorageEndpointProperties, messagingEndpoints map[string]*devices.MessagingEndpointProperties, enableFileUploadNotifications *bool) []interface{} {
	results := make([]interface{}, 0)
	output := make(map[string]interface{})

	if storageEndpointProperties, ok := storageEndpoints["$default"]; ok {
		if connString := storageEndpointProperties.ConnectionString; connString != nil {
			output["connection_string"] = *connString
		}
		if containerName := storageEndpointProperties.ContainerName; containerName != nil {
			output["container_name"] = *containerName
		}
		if sasTTLAsIso8601 := storageEndpointProperties.SasTTLAsIso8601; sasTTLAsIso8601 != nil {
			output["sas_ttl"] = *sasTTLAsIso8601
		}

		if messagingEndpointProperties, ok := messagingEndpoints["fileNotifications"]; ok {
			if lockDurationAsIso8601 := messagingEndpointProperties.LockDurationAsIso8601; lockDurationAsIso8601 != nil {
				output["lock_duration"] = *lockDurationAsIso8601
			}
			if ttlAsIso8601 := messagingEndpointProperties.TTLAsIso8601; ttlAsIso8601 != nil {
				output["default_ttl"] = *ttlAsIso8601
			}
			if maxDeliveryCount := messagingEndpointProperties.MaxDeliveryCount; maxDeliveryCount != nil {
				output["max_delivery_count"] = *maxDeliveryCount
			}
		}

		if enableFileUploadNotifications != nil {
			output["notifications"] = *enableFileUploadNotifications
		}

		results = append(results, output)
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

				output["encoding"] = string(container.Encoding)
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

				output["type"] = "AzureIotHub.ServiceBusQueue"

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

				output["type"] = "AzureIotHub.ServiceBusTopic"

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

				output["type"] = "AzureIotHub.EventHub"

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

func flattenIoTHubFallbackRoute(input *devices.RoutingProperties) []interface{} {
	if input.FallbackRoute == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})
	route := input.FallbackRoute

	if condition := route.Condition; condition != nil {
		output["condition"] = *condition
	}
	if isEnabled := route.IsEnabled; isEnabled != nil {
		output["enabled"] = *isEnabled
	}
	if source := route.Source; source != nil {
		output["source"] = *source
	}

	output["endpoint_names"] = utils.FlattenStringSlice(route.EndpointNames)

	return []interface{}{output}
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
func expandIPFilterRules(d *schema.ResourceData) *[]devices.IPFilterRule {
	ipFilterRuleList := d.Get("ip_filter_rule").(*schema.Set).List()
	if len(ipFilterRuleList) == 0 {
		return nil
	}

	rules := make([]devices.IPFilterRule, 0)

	for _, r := range ipFilterRuleList {
		rawRule := r.(map[string]interface{})
		rule := &devices.IPFilterRule{
			FilterName: utils.String(rawRule["name"].(string)),
			Action:     devices.IPFilterActionType(rawRule["action"].(string)),
			IPMask:     utils.String(rawRule["ip_mask"].(string)),
		}

		rules = append(rules, *rule)
	}
	return &rules
}

func flattenIPFilterRules(in *[]devices.IPFilterRule) []interface{} {
	rules := make([]interface{}, 0)
	if in == nil {
		return rules
	}

	for _, r := range *in {
		rawRule := make(map[string]interface{})

		if r.FilterName != nil {
			rawRule["name"] = *r.FilterName
		}

		rawRule["action"] = string(r.Action)

		if r.IPMask != nil {
			rawRule["ip_mask"] = *r.IPMask
		}
		rules = append(rules, rawRule)
	}
	return rules
}
