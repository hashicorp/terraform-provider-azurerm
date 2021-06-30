package iothub

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/iothub/mgmt/2020-03-01/devices"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iothub/parse"
	iothubValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iothub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// TODO: outside of this pr make this private

var IothubResourceName = "azurerm_iothub"

// nolint unparam
func suppressIfTypeIsNot(t string) pluginsdk.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *pluginsdk.ResourceData) bool {
		path := strings.Split(k, ".")
		path[len(path)-1] = "type"
		return d.Get(strings.Join(path, ".")).(string) != t
	}
}

// nolint unparam
func supressWhenAll(fs ...pluginsdk.SchemaDiffSuppressFunc) pluginsdk.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *pluginsdk.ResourceData) bool {
		for _, f := range fs {
			if !f(k, old, new, d) {
				return false
			}
		}
		return true
	}
}

func resourceIotHub() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubCreateUpdate,
		Read:   resourceIotHubRead,
		Update: resourceIotHubCreateUpdate,
		Delete: resourceIotHubDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.IotHubID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: iothubValidate.IoTHubName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:             pluginsdk.TypeString,
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
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 200),
						},
					},
				},
			},

			"shared_access_policy": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"primary_key": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"secondary_key": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"permissions": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"event_hub_partition_count": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(2, 128),
			},
			"event_hub_retention_in_days": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 7),
			},

			"file_upload": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"connection_string": {
							Type:     pluginsdk.TypeString,
							Required: true,
							DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
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
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"notifications": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
						"max_delivery_count": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      10,
							ValidateFunc: validation.IntBetween(1, 100),
						},
						"sas_ttl": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.ISO8601Duration,
						},
						"default_ttl": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.ISO8601Duration,
						},
						"lock_duration": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.ISO8601Duration,
						},
					},
				},
			},

			"endpoint": {
				Type:       pluginsdk.TypeList,
				Optional:   true,
				Computed:   true,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"AzureIotHub.StorageContainer",
								"AzureIotHub.ServiceBusQueue",
								"AzureIotHub.ServiceBusTopic",
								"AzureIotHub.EventHub",
							}, false),
						},

						"connection_string": {
							Type:     pluginsdk.TypeString,
							Required: true,
							DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
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
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: iothubValidate.IoTHubEndpointName,
						},

						"batch_frequency_in_seconds": {
							Type:             pluginsdk.TypeInt,
							Optional:         true,
							Default:          300,
							DiffSuppressFunc: suppressIfTypeIsNot("AzureIotHub.StorageContainer"),
							ValidateFunc:     validation.IntBetween(60, 720),
						},

						"max_chunk_size_in_bytes": {
							Type:             pluginsdk.TypeInt,
							Optional:         true,
							Default:          314572800,
							DiffSuppressFunc: suppressIfTypeIsNot("AzureIotHub.StorageContainer"),
							ValidateFunc:     validation.IntBetween(10485760, 524288000),
						},

						"container_name": {
							Type:             pluginsdk.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppressIfTypeIsNot("AzureIotHub.StorageContainer"),
						},

						"encoding": {
							Type:     pluginsdk.TypeString,
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
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: iothubValidate.FileNameFormat,
						},

						"resource_group_name": azure.SchemaResourceGroupNameOptional(),
					},
				},
			},

			"route": {
				Type:       pluginsdk.TypeList,
				Optional:   true,
				Computed:   true,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile("^[-_.a-zA-Z0-9]{1,64}$"),
								"Route Name name can only include alphanumeric characters, periods, underscores, hyphens, has a maximum length of 64 characters, and must be unique.",
							),
						},
						"source": {
							Type:     pluginsdk.TypeString,
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
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  "true",
						},
						"endpoint_names": {
							Type: pluginsdk.TypeList,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
							Required: true,
						},
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},
					},
				},
			},

			"enrichment": {
				Type: pluginsdk.TypeList,
				// Currently only 10 enrichments is allowed for standard or basic tier, 2 for Free tier.
				MaxItems:   10,
				Optional:   true,
				Computed:   true,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile("^[-_.a-zA-Z0-9]{1,64}$"),
								"Enrichment Key name can only include alphanumeric characters, periods, underscores, hyphens, has a maximum length of 64 characters, and must be unique.",
							),
						},
						"value": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"endpoint_names": {
							Type: pluginsdk.TypeList,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
							Required: true,
						},
					},
				},
			},

			"fallback_route": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"source": {
							Type:     pluginsdk.TypeString,
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
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  "true",
						},
						"endpoint_names": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringLenBetween(0, 64),
							},
						},
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"ip_filter_rule": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"ip_mask": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.CIDR,
						},
						"action": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(devices.Accept),
								string(devices.Reject),
							}, false),
						},
					},
				},
			},

			"min_tls_version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"1.2",
				}, false),
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"hostname": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"event_hub_events_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"event_hub_operations_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"event_hub_events_path": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"event_hub_operations_path": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceIotHubCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	locks.ByName(name, IothubResourceName)
	defer locks.UnlockByName(name, IothubResourceName)

	if d.IsNewResource() {
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
		if _, err = client.Get(ctx, resourceGroup, name); err != nil {
			return fmt.Errorf("An IoTHub already exists with the name %q - please choose an alternate name: %s", name, string(res.Reason))
		}
	}

	routingProperties := devices.RoutingProperties{}

	if _, ok := d.GetOk("route"); ok {
		routingProperties.Routes = expandIoTHubRoutes(d)
	}

	if _, ok := d.GetOk("enrichment"); ok {
		routingProperties.Enrichments = expandIoTHubEnrichments(d)
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

	// nolint staticcheck
	if v, ok := d.GetOkExists("public_network_access_enabled"); ok {
		enabled := devices.Disabled
		if v.(bool) {
			enabled = devices.Enabled
		}
		props.Properties.PublicNetworkAccess = enabled
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

	if v, ok := d.GetOk("min_tls_version"); ok {
		props.Properties.MinTLSVersion = utils.String(v.(string))
	}

	if _, err = client.CreateOrUpdate(ctx, resourceGroup, name, props, ""); err != nil {
		return fmt.Errorf("Error creating/updating IotHub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	timeout := pluginsdk.TimeoutUpdate
	if d.IsNewResource() {
		timeout = pluginsdk.TimeoutCreate
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"Activating", "Transitioning"},
		Target:  []string{"Succeeded"},
		Refresh: iothubStateRefreshFunc(ctx, client, resourceGroup, name),
		Timeout: d.Timeout(timeout),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("Error waiting for the completion of the creating/updating of IotHub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	d.SetId(*resp.ID)

	return resourceIotHubRead(d, meta)
}

func resourceIotHubRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IotHubID(d.Id())
	if err != nil {
		return err
	}

	hub, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(hub.Response) {
			log.Printf("[DEBUG] IoTHub %q (Resource Group %q) was not found!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving IotHub Client %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if keysResp, err := client.ListKeys(ctx, id.ResourceGroup, id.Name); err == nil {
		keyList := keysResp.Response()
		keys := flattenIoTHubSharedAccessPolicy(keyList.Value)

		if err := d.Set("shared_access_policy", keys); err != nil {
			return fmt.Errorf("setting `shared_access_policy` in IoTHub %q: %+v", id.Name, err)
		}
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
			return fmt.Errorf("setting `endpoint` in IoTHub %q: %+v", id.Name, err)
		}

		routes := flattenIoTHubRoute(properties.Routing)
		if err := d.Set("route", routes); err != nil {
			return fmt.Errorf("setting `route` in IoTHub %q: %+v", id.Name, err)
		}

		enrichments := flattenIoTHubEnrichment(properties.Routing)
		if err := d.Set("enrichment", enrichments); err != nil {
			return fmt.Errorf("setting `enrichment` in IoTHub %q: %+v", id.Name, err)
		}

		fallbackRoute := flattenIoTHubFallbackRoute(properties.Routing)
		if err := d.Set("fallback_route", fallbackRoute); err != nil {
			return fmt.Errorf("setting `fallbackRoute` in IoTHub %q: %+v", id.Name, err)
		}

		ipFilterRules := flattenIPFilterRules(properties.IPFilterRules)
		if err := d.Set("ip_filter_rule", ipFilterRules); err != nil {
			return fmt.Errorf("setting `ip_filter_rule` in IoTHub %q: %+v", id.Name, err)
		}

		fileUpload := flattenIoTHubFileUpload(properties.StorageEndpoints, properties.MessagingEndpoints, properties.EnableFileUploadNotifications)
		if err := d.Set("file_upload", fileUpload); err != nil {
			return fmt.Errorf("setting `file_upload` in IoTHub %q: %+v", id.Name, err)
		}

		if enabled := properties.PublicNetworkAccess; enabled != "" {
			d.Set("public_network_access_enabled", enabled == devices.Enabled)
		}

		d.Set("min_tls_version", properties.MinTLSVersion)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
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

func resourceIotHubDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.IotHubID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	locks.ByName(id.Name, IothubResourceName)
	defer locks.UnlockByName(id.Name, IothubResourceName)

	// when running acctest of `azurerm_iot_security_solution`, we found after delete the iot security solution, the iothub provisionState is `Transitioning`
	// if we delete directly, the func `client.Delete` will throw error
	// so first wait for the iotHub state become succeed
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"Activating", "Transitioning"},
		Target:  []string{"Succeeded"},
		Refresh: iothubStateRefreshFunc(ctx, client, id.ResourceGroup, id.Name),
		Timeout: d.Timeout(pluginsdk.TimeoutDelete),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for ProvisioningState of IotHub %q (Resource Group %q) to become `Succeeded`: %+v", id.Name, id.ResourceGroup, err)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return err
	}

	return waitForIotHubToBeDeleted(ctx, client, id.ResourceGroup, id.Name, d)
}

func waitForIotHubToBeDeleted(ctx context.Context, client *devices.IotHubResourceClient, resourceGroup, name string, d *pluginsdk.ResourceData) error {
	// we can't use the Waiter here since the API returns a 404 once it's deleted which is considered a polling status code..
	log.Printf("[DEBUG] Waiting for IotHub (%q in Resource Group %q) to be deleted", name, resourceGroup)
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"200"},
		Target:  []string{"404"},
		Refresh: iothubStateStatusCodeRefreshFunc(ctx, client, resourceGroup, name),
		Timeout: d.Timeout(pluginsdk.TimeoutDelete),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("Error waiting for IotHub (%q in Resource Group %q) to be deleted: %+v", name, resourceGroup, err)
	}

	return nil
}

func iothubStateRefreshFunc(ctx context.Context, client *devices.IotHubResourceClient, resourceGroup, name string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroup, name)

		log.Printf("Retrieving IoTHub %q (Resource Group %q) returned Status %d", resourceGroup, name, res.StatusCode)

		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return res, "NotFound", nil
			}
			return nil, "", fmt.Errorf("polling for the Provisioning State of the IotHub %q (RG: %q): %+v", name, resourceGroup, err)
		}

		if res.Properties == nil || res.Properties.ProvisioningState == nil {
			return res, "", fmt.Errorf("polling for the Provisioning State of the IotHub %q (RG: %q): %+v", name, resourceGroup, err)
		}

		return res, *res.Properties.ProvisioningState, nil
	}
}

func iothubStateStatusCodeRefreshFunc(ctx context.Context, client *devices.IotHubResourceClient, resourceGroup, name string) pluginsdk.StateRefreshFunc {
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

func expandIoTHubRoutes(d *pluginsdk.ResourceData) *[]devices.RouteProperties {
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

func expandIoTHubEnrichments(d *pluginsdk.ResourceData) *[]devices.EnrichmentProperties {
	enrichmentList := d.Get("enrichment").([]interface{})

	enrichmentProperties := make([]devices.EnrichmentProperties, 0)

	for _, enrichmentRaw := range enrichmentList {
		enrichment := enrichmentRaw.(map[string]interface{})

		key := enrichment["key"].(string)
		value := enrichment["value"].(string)

		endpointNamesRaw := enrichment["endpoint_names"].([]interface{})

		enrichmentProperties = append(enrichmentProperties, devices.EnrichmentProperties{
			Key:           &key,
			Value:         &value,
			EndpointNames: utils.ExpandStringSlice(endpointNamesRaw),
		})
	}

	return &enrichmentProperties
}

func expandIoTHubFileUpload(d *pluginsdk.ResourceData) (map[string]*devices.StorageEndpointProperties, map[string]*devices.MessagingEndpointProperties, bool) {
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

func expandIoTHubEndpoints(d *pluginsdk.ResourceData, subscriptionId string) *devices.RoutingEndpoints {
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
		resourceGroup := endpoint["resource_group_name"].(string)
		subscriptionID := subscriptionId

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

func expandIoTHubFallbackRoute(d *pluginsdk.ResourceData) *devices.FallbackRouteProperties {
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

func expandIoTHubSku(d *pluginsdk.ResourceData) *devices.IotHubSkuInfo {
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
				if resourceGroup := container.ResourceGroup; resourceGroup != nil {
					output["resource_group_name"] = *resourceGroup
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
				if resourceGroup := queue.ResourceGroup; resourceGroup != nil {
					output["resource_group_name"] = *resourceGroup
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
				if resourceGroup := topic.ResourceGroup; resourceGroup != nil {
					output["resource_group_name"] = *resourceGroup
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
				if resourceGroup := eventHub.ResourceGroup; resourceGroup != nil {
					output["resource_group_name"] = *resourceGroup
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

func flattenIoTHubEnrichment(input *devices.RoutingProperties) []interface{} {
	results := make([]interface{}, 0)

	if input != nil && input.Enrichments != nil {
		for _, enrichment := range *input.Enrichments {
			output := make(map[string]interface{})

			if key := enrichment.Key; key != nil {
				output["key"] = *key
			}
			if value := enrichment.Value; value != nil {
				output["value"] = *value
			}
			if endpointNames := enrichment.EndpointNames; endpointNames != nil {
				output["endpoint_names"] = *endpointNames
			}

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

func expandIPFilterRules(d *pluginsdk.ResourceData) *[]devices.IPFilterRule {
	ipFilterRuleList := d.Get("ip_filter_rule").([]interface{})
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
