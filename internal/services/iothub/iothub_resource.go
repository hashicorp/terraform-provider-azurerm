// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	eventhubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	iothubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	servicebusValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	devices "github.com/tombuildsstuff/kermit/sdk/iothub/2022-04-30-preview/iothub"
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
func suppressIfTypeIs(t string) pluginsdk.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *pluginsdk.ResourceData) bool {
		path := strings.Split(k, ".")
		path[len(path)-1] = "type"
		return d.Get(strings.Join(path, ".")).(string) == t
	}
}

// nolint unparam
func suppressWhenAny(fs ...pluginsdk.SchemaDiffSuppressFunc) pluginsdk.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *pluginsdk.ResourceData) bool {
		for _, f := range fs {
			if f(k, old, new, d) {
				return true
			}
		}
		return false
	}
}

func resourceIotHub() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubCreateUpdate,
		Read:   resourceIotHubRead,
		Update: resourceIotHubCreateUpdate,
		Delete: resourceIotHubDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.IoTHubV0ToV1{},
		}),

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

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"sku": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(devices.IotHubSkuB1),
								string(devices.IotHubSkuB2),
								string(devices.IotHubSkuB3),
								string(devices.IotHubSkuF1),
								string(devices.IotHubSkuS1),
								string(devices.IotHubSkuS2),
								string(devices.IotHubSkuS3),
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
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: fileUploadConnectionStringDiffSuppress,
							Sensitive:        true,
						},
						"container_name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"authentication_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(devices.AuthenticationTypeKeyBased),
							ValidateFunc: validation.StringInSlice([]string{
								string(devices.AuthenticationTypeKeyBased),
								string(devices.AuthenticationTypeIdentityBased),
							}, false),
						},
						"identity_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: commonids.ValidateUserAssignedIdentityID,
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
							Computed:     !features.FourPointOhBeta(),
							ValidateFunc: validate.ISO8601Duration,
							Default: func() interface{} {
								if !features.FourPointOhBeta() {
									return nil
								}
								return "PT1H"
							}(),
						},
						"default_ttl": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     !features.FourPointOhBeta(),
							ValidateFunc: validate.ISO8601Duration,
							Default: func() interface{} {
								if !features.FourPointOhBeta() {
									return nil
								}
								return "PT1H"
							}(),
						},
						"lock_duration": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     !features.FourPointOhBeta(),
							ValidateFunc: validate.ISO8601Duration,
							Default: func() interface{} {
								if !features.FourPointOhBeta() {
									return nil
								}
								return "PT1M"
							}(),
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

						"authentication_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(devices.AuthenticationTypeKeyBased),
							ValidateFunc: validation.StringInSlice([]string{
								string(devices.AuthenticationTypeKeyBased),
								string(devices.AuthenticationTypeIdentityBased),
							}, false),
						},

						"identity_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: commonids.ValidateUserAssignedIdentityID,
						},

						"endpoint_uri": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"entity_path": {
							Type:             pluginsdk.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppressIfTypeIs("AzureIotHub.StorageContainer"),
							ValidateFunc: validation.Any(
								servicebusValidate.QueueName(),
								servicebusValidate.TopicName(),
								eventhubValidate.ValidateEventHubName(),
							),
						},

						"connection_string": {
							Type:     pluginsdk.TypeString,
							Optional: true,
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
							Type:             pluginsdk.TypeString,
							Optional:         true,
							ForceNew:         true,
							Default:          string(devices.EncodingAvro),
							DiffSuppressFunc: suppressIfTypeIsNot("AzureIotHub.StorageContainer"),
							ValidateFunc: validation.StringInSlice([]string{
								string(devices.EncodingAvro),
								string(devices.EncodingAvroDeflate),
								string(devices.EncodingJSON),
							}, false),
						},

						"file_name_format": {
							Type:             pluginsdk.TypeString,
							Optional:         true,
							Default:          "{iothub}/{partition}/{YYYY}/{MM}/{DD}/{HH}/{mm}",
							DiffSuppressFunc: suppressIfTypeIsNot("AzureIotHub.StorageContainer"),
							ValidateFunc:     iothubValidate.FileNameFormat,
						},

						"resource_group_name": commonschema.ResourceGroupNameOptional(),
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
								string(devices.RoutingSourceDeviceConnectionStateEvents),
								string(devices.RoutingSourceDeviceJobLifecycleEvents),
								string(devices.RoutingSourceDeviceLifecycleEvents),
								string(devices.RoutingSourceDeviceMessages),
								string(devices.RoutingSourceDigitalTwinChangeEvents),
								string(devices.RoutingSourceInvalid),
								string(devices.RoutingSourceTwinChangeEvents),
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
							Default:  string(devices.RoutingSourceDeviceMessages),
							ValidateFunc: validation.StringInSlice([]string{
								string(devices.RoutingSourceDeviceConnectionStateEvents),
								string(devices.RoutingSourceDeviceJobLifecycleEvents),
								string(devices.RoutingSourceDeviceLifecycleEvents),
								string(devices.RoutingSourceDeviceMessages),
								string(devices.RoutingSourceDigitalTwinChangeEvents),
								string(devices.RoutingSourceInvalid),
								string(devices.RoutingSourceTwinChangeEvents),
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

			"network_rule_set": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"default_action": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(devices.DefaultActionDeny),
							ValidateFunc: validation.StringInSlice([]string{
								string(devices.DefaultActionAllow),
								string(devices.DefaultActionDeny),
							}, false),
						},
						"apply_to_builtin_eventhub_endpoint": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
						"ip_rule": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: iothubValidate.IoTHubIpRuleName,
									},
									"ip_mask": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validate.CIDR,
									},
									"action": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  string(devices.NetworkRuleIPActionAllow),
										ValidateFunc: validation.StringInSlice([]string{
											string(devices.NetworkRuleIPActionAllow),
										}, false),
									},
								},
							},
						},
					},
				},
			},

			"cloud_to_device": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"max_delivery_count": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      10,
							ValidateFunc: validation.IntBetween(1, 100),
						},
						"default_ttl": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "PT1H",
							ValidateFunc: validate.ISO8601DurationBetween("PT15M", "P2D"),
						},
						"feedback": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"time_to_live": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Default:      "PT1H",
										ValidateFunc: validate.ISO8601DurationBetween("PT15M", "P2D"),
									},
									"max_delivery_count": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										Default:      10,
										ValidateFunc: validation.IntBetween(1, 100),
									},
									"lock_duration": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Default:      "PT60S",
										ValidateFunc: validate.ISO8601DurationBetween("PT5S", "PT300S"),
									},
								},
							},
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
			"event_hub_events_namespace": {
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

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"tags": tags.Schema(),
		},
	}
}

func resourceIotHubCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	id := parse.NewIotHubID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByName(id.Name, IothubResourceName)
	defer locks.UnlockByName(id.Name, IothubResourceName)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_iothub", id.ID())
		}
	}

	res, err := client.CheckNameAvailability(ctx, devices.OperationInputs{
		Name: &id.Name,
	})
	if err != nil {
		return fmt.Errorf("An error occurred checking if the IoTHub name was unique: %+v", err)
	}

	if !*res.NameAvailable {
		if _, err = client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			return fmt.Errorf("An IoTHub already exists with the name %q - please choose an alternate name: %s", id.Name, string(res.Reason))
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
	} else {
		routingProperties.FallbackRoute = &devices.FallbackRouteProperties{
			Source:        utils.String(string(devices.RoutingSourceDeviceMessages)),
			Condition:     utils.String("true"),
			EndpointNames: &[]string{"events"},
			IsEnabled:     utils.Bool(true),
		}
	}

	if _, ok := d.GetOk("endpoint"); ok {
		routingProperties.Endpoints, err = expandIoTHubEndpoints(d, subscriptionId)
		if err != nil {
			return fmt.Errorf("expanding `endpoint`: %+v", err)
		}
	}

	storageEndpoints, messagingEndpoints, enableFileUploadNotifications, err := expandIoTHubFileUpload(d)
	if err != nil {
		return fmt.Errorf("expanding `file_upload`: %+v", err)
	}

	cloudToDeviceProperties := &devices.CloudToDeviceProperties{}
	if _, ok := d.GetOk("cloud_to_device"); ok {
		cloudToDeviceProperties = expandIoTHubCloudToDevice(d)
	}

	identity, err := expandIotHubIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	props := devices.IotHubDescription{
		Name:     utils.String(id.Name),
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Sku:      expandIoTHubSku(d),
		Properties: &devices.IotHubProperties{
			Routing:                       &routingProperties,
			StorageEndpoints:              storageEndpoints,
			MessagingEndpoints:            messagingEndpoints,
			EnableFileUploadNotifications: &enableFileUploadNotifications,
			CloudToDevice:                 cloudToDeviceProperties,
		},
		Identity: identity,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, ok := d.GetOk("network_rule_set"); ok {
		props.Properties.NetworkRuleSets = expandNetworkRuleSetProperties(d)
	}

	// nolint staticcheck
	if v, ok := d.GetOkExists("public_network_access_enabled"); ok {
		enabled := devices.PublicNetworkAccessDisabled
		if v.(bool) {
			enabled = devices.PublicNetworkAccessEnabled
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

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, props, "")
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %q: %+v", id, err)
	}

	d.SetId(id.ID())

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
			log.Printf("[DEBUG] %s was not found!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
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

				if *v.Endpoint != "" {
					uri, err := url.Parse(*v.Endpoint)
					if err == nil {
						d.Set("event_hub_events_namespace", strings.Split(uri.Hostname(), ".")[0])
					}
				}

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

		networkRuleSet := flattenNetworkRuleSetProperties(properties.NetworkRuleSets)
		if err := d.Set("network_rule_set", networkRuleSet); err != nil {
			return fmt.Errorf("setting `network_rule_set` in IoTHub %q: %+v", id.Name, err)
		}

		fileUpload := flattenIoTHubFileUpload(properties.StorageEndpoints, properties.MessagingEndpoints, properties.EnableFileUploadNotifications)
		if err := d.Set("file_upload", fileUpload); err != nil {
			return fmt.Errorf("setting `file_upload` in IoTHub %q: %+v", id.Name, err)
		}

		if enabled := properties.PublicNetworkAccess; enabled != "" {
			d.Set("public_network_access_enabled", enabled == devices.PublicNetworkAccessEnabled)
		}

		cloudToDevice := flattenIoTHubCloudToDevice(properties.CloudToDevice)
		if err := d.Set("cloud_to_device", cloudToDevice); err != nil {
			return fmt.Errorf("setting `cloudToDevice` in IoTHub %q: %+v", id.Name, err)
		}

		d.Set("min_tls_version", properties.MinTLSVersion)
	}

	identity, err := flattenIotHubIdentity(hub.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := hub.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	sku := flattenIoTHubSku(hub.Sku)
	if err := d.Set("sku", sku); err != nil {
		return fmt.Errorf("setting `sku`: %+v", err)
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

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return err
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %q: %+v", id, err)
	}

	return nil
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

func expandIoTHubFileUpload(d *pluginsdk.ResourceData) (map[string]*devices.StorageEndpointProperties, map[string]*devices.MessagingEndpointProperties, bool, error) {
	fileUploadList := d.Get("file_upload").([]interface{})

	storageEndpointProperties := make(map[string]*devices.StorageEndpointProperties)
	messagingEndpointProperties := make(map[string]*devices.MessagingEndpointProperties)
	notifications := false

	if len(fileUploadList) > 0 {
		fileUploadMap := fileUploadList[0].(map[string]interface{})

		authenticationType := devices.AuthenticationType(fileUploadMap["authentication_type"].(string))
		identityId := fileUploadMap["identity_id"].(string)
		connectionStr := fileUploadMap["connection_string"].(string)
		containerName := fileUploadMap["container_name"].(string)
		notifications = fileUploadMap["notifications"].(bool)
		maxDeliveryCount := int32(fileUploadMap["max_delivery_count"].(int))
		sasTTL := fileUploadMap["sas_ttl"].(string)
		defaultTTL := fileUploadMap["default_ttl"].(string)
		lockDuration := fileUploadMap["lock_duration"].(string)

		storageEndpointProperties["$default"] = &devices.StorageEndpointProperties{
			AuthenticationType: authenticationType,
			ConnectionString:   &connectionStr,
			ContainerName:      &containerName,
		}

		if sasTTL != "" {
			storageEndpointProperties["$default"].SasTTLAsIso8601 = &sasTTL
		}

		if identityId != "" {
			if authenticationType != devices.AuthenticationTypeIdentityBased {
				return nil, nil, false, fmt.Errorf("`identity_id` can only be specified when `authentication_type` is `identityBased`")
			}
			storageEndpointProperties["$default"].Identity = &devices.ManagedIdentity{
				UserAssignedIdentity: &identityId,
			}
		}

		messagingEndpointProperties["fileNotifications"] = &devices.MessagingEndpointProperties{
			LockDurationAsIso8601: &lockDuration,
			TTLAsIso8601:          &defaultTTL,
			MaxDeliveryCount:      &maxDeliveryCount,
		}
	}

	return storageEndpointProperties, messagingEndpointProperties, notifications, nil
}

func expandIoTHubEndpoints(d *pluginsdk.ResourceData, subscriptionId string) (*devices.RoutingEndpoints, error) {
	routeEndpointList := d.Get("endpoint").([]interface{})

	serviceBusQueueEndpointProperties := make([]devices.RoutingServiceBusQueueEndpointProperties, 0)
	serviceBusTopicEndpointProperties := make([]devices.RoutingServiceBusTopicEndpointProperties, 0)
	eventHubProperties := make([]devices.RoutingEventHubProperties, 0)
	storageContainerProperties := make([]devices.RoutingStorageContainerProperties, 0)

	for _, endpointRaw := range routeEndpointList {
		endpoint := endpointRaw.(map[string]interface{})

		t := endpoint["type"]
		name := endpoint["name"].(string)
		resourceGroup := endpoint["resource_group_name"].(string)
		authenticationType := devices.AuthenticationType(endpoint["authentication_type"].(string))
		subscriptionID := subscriptionId

		var identity *devices.ManagedIdentity
		var endpointUri *string
		var entityPath *string
		var connectionStr *string
		if v := endpoint["identity_id"].(string); v != "" {
			identity = &devices.ManagedIdentity{
				UserAssignedIdentity: utils.String(v),
			}
		}
		if v := endpoint["endpoint_uri"].(string); v != "" {
			endpointUri = utils.String(v)
		}
		if v := endpoint["entity_path"].(string); v != "" {
			entityPath = utils.String(v)
		}
		if v := endpoint["connection_string"].(string); v != "" {
			connectionStr = utils.String(v)
		}

		if authenticationType == devices.AuthenticationTypeKeyBased {
			if connectionStr == nil {
				return nil, fmt.Errorf("`connection_string` must be specified when `authentication_type` is `keyBased`")
			}
			if identity != nil || endpointUri != nil || entityPath != nil {
				return nil, fmt.Errorf("`identity_id`, `endpoint_uri` or `entity_path` cannot be specified when `authentication_type` is `keyBased`")
			}
		} else {
			if endpointUri == nil {
				return nil, fmt.Errorf("`endpoint_uri` must be specified when `authentication_type` is `identityBased`")
			}

			if entityPath == nil && t != "AzureIotHub.StorageContainer" {
				return nil, fmt.Errorf("`entity_path` must be specified when `authentication_type` is `identityBased` and `type` is `%s`", t)
			}

			if connectionStr != nil {
				return nil, fmt.Errorf("`connection_string` cannot be specified when `authentication_type` is `identityBased`")
			}
		}

		switch t {
		case "AzureIotHub.StorageContainer":
			containerName := endpoint["container_name"].(string)
			if containerName == "" {
				return nil, fmt.Errorf("`container_name` must be specified when `type` is `AzureIotHub.StorageContainer`")
			}

			fileNameFormat := endpoint["file_name_format"].(string)
			batchFrequencyInSeconds := int32(endpoint["batch_frequency_in_seconds"].(int))
			maxChunkSizeInBytes := int32(endpoint["max_chunk_size_in_bytes"].(int))
			encoding := endpoint["encoding"].(string)

			storageContainer := devices.RoutingStorageContainerProperties{
				AuthenticationType:      authenticationType,
				Identity:                identity,
				EndpointURI:             endpointUri,
				ConnectionString:        connectionStr,
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
				AuthenticationType: authenticationType,
				Identity:           identity,
				EndpointURI:        endpointUri,
				EntityPath:         entityPath,
				ConnectionString:   connectionStr,
				Name:               &name,
				SubscriptionID:     &subscriptionID,
				ResourceGroup:      &resourceGroup,
			}
			serviceBusQueueEndpointProperties = append(serviceBusQueueEndpointProperties, sbQueue)

		case "AzureIotHub.ServiceBusTopic":
			sbTopic := devices.RoutingServiceBusTopicEndpointProperties{
				AuthenticationType: authenticationType,
				Identity:           identity,
				EndpointURI:        endpointUri,
				EntityPath:         entityPath,
				ConnectionString:   connectionStr,
				Name:               &name,
				SubscriptionID:     &subscriptionID,
				ResourceGroup:      &resourceGroup,
			}
			serviceBusTopicEndpointProperties = append(serviceBusTopicEndpointProperties, sbTopic)

		case "AzureIotHub.EventHub":
			eventHub := devices.RoutingEventHubProperties{
				AuthenticationType: authenticationType,
				Identity:           identity,
				EndpointURI:        endpointUri,
				EntityPath:         entityPath,
				ConnectionString:   connectionStr,
				Name:               &name,
				SubscriptionID:     &subscriptionID,
				ResourceGroup:      &resourceGroup,
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

func expandIoTHubCloudToDevice(d *pluginsdk.ResourceData) *devices.CloudToDeviceProperties {
	ctdList := d.Get("cloud_to_device").([]interface{})
	if len(ctdList) == 0 {
		return nil
	}
	cloudToDevice := devices.CloudToDeviceProperties{}
	ctdMap := ctdList[0].(map[string]interface{})
	defaultTimeToLive := ctdMap["default_ttl"].(string)

	cloudToDevice.DefaultTTLAsIso8601 = &defaultTimeToLive
	cloudToDevice.MaxDeliveryCount = utils.Int32(int32(ctdMap["max_delivery_count"].(int)))
	feedback := ctdMap["feedback"].([]interface{})

	cloudToDeviceFeedback := devices.FeedbackProperties{}
	if len(feedback) > 0 {
		feedbackMap := feedback[0].(map[string]interface{})

		lockDuration := feedbackMap["lock_duration"].(string)
		timeToLive := feedbackMap["time_to_live"].(string)

		cloudToDeviceFeedback.TTLAsIso8601 = &timeToLive
		cloudToDeviceFeedback.LockDurationAsIso8601 = &lockDuration
		cloudToDeviceFeedback.MaxDeliveryCount = utils.Int32(int32(feedbackMap["max_delivery_count"].(int)))
	}

	cloudToDevice.Feedback = &cloudToDeviceFeedback

	return &cloudToDevice
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

		authenticationType := string(devices.AuthenticationTypeKeyBased)
		if v := string(storageEndpointProperties.AuthenticationType); v != "" {
			authenticationType = v
		}
		output["authentication_type"] = authenticationType

		identityId := ""
		if storageEndpointProperties.Identity != nil && storageEndpointProperties.Identity.UserAssignedIdentity != nil {
			identityId = *storageEndpointProperties.Identity.UserAssignedIdentity
		}
		output["identity_id"] = identityId

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

				authenticationType := string(devices.AuthenticationTypeKeyBased)
				if string(container.AuthenticationType) != "" {
					authenticationType = string(container.AuthenticationType)
				}
				output["authentication_type"] = authenticationType

				connectionStr := ""
				if container.ConnectionString != nil {
					connectionStr = *container.ConnectionString
				}
				output["connection_string"] = connectionStr

				endpointUri := ""
				if container.EndpointURI != nil {
					endpointUri = *container.EndpointURI
				}
				output["endpoint_uri"] = endpointUri

				identityId := ""
				if container.Identity != nil && container.Identity.UserAssignedIdentity != nil {
					identityId = *container.Identity.UserAssignedIdentity
				}
				output["identity_id"] = identityId

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

				authenticationType := string(devices.AuthenticationTypeKeyBased)
				if string(queue.AuthenticationType) != "" {
					authenticationType = string(queue.AuthenticationType)
				}
				output["authentication_type"] = authenticationType

				connectionStr := ""
				if queue.ConnectionString != nil {
					connectionStr = *queue.ConnectionString
				}
				output["connection_string"] = connectionStr

				endpointUri := ""
				if queue.EndpointURI != nil {
					endpointUri = *queue.EndpointURI
				}
				output["endpoint_uri"] = endpointUri

				entityPath := ""
				if queue.EntityPath != nil {
					entityPath = *queue.EntityPath
				}
				output["entity_path"] = entityPath

				identityId := ""
				if queue.Identity != nil && queue.Identity.UserAssignedIdentity != nil {
					identityId = *queue.Identity.UserAssignedIdentity
				}
				output["identity_id"] = identityId

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

				authenticationType := string(devices.AuthenticationTypeKeyBased)
				if string(topic.AuthenticationType) != "" {
					authenticationType = string(topic.AuthenticationType)
				}
				output["authentication_type"] = authenticationType

				connectionStr := ""
				if topic.ConnectionString != nil {
					connectionStr = *topic.ConnectionString
				}
				output["connection_string"] = connectionStr

				endpointUri := ""
				if topic.EndpointURI != nil {
					endpointUri = *topic.EndpointURI
				}
				output["endpoint_uri"] = endpointUri

				entityPath := ""
				if topic.EntityPath != nil {
					entityPath = *topic.EntityPath
				}
				output["entity_path"] = entityPath

				identityId := ""
				if topic.Identity != nil && topic.Identity.UserAssignedIdentity != nil {
					identityId = *topic.Identity.UserAssignedIdentity
				}
				output["identity_id"] = identityId

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

				authenticationType := string(devices.AuthenticationTypeKeyBased)
				if string(eventHub.AuthenticationType) != "" {
					authenticationType = string(eventHub.AuthenticationType)
				}
				output["authentication_type"] = authenticationType

				connectionStr := ""
				if eventHub.ConnectionString != nil {
					connectionStr = *eventHub.ConnectionString
				}
				output["connection_string"] = connectionStr

				endpointUri := ""
				if eventHub.EndpointURI != nil {
					endpointUri = *eventHub.EndpointURI
				}
				output["endpoint_uri"] = endpointUri

				entityPath := ""
				if eventHub.EntityPath != nil {
					entityPath = *eventHub.EntityPath
				}
				output["entity_path"] = entityPath

				identityId := ""
				if eventHub.Identity != nil && eventHub.Identity.UserAssignedIdentity != nil {
					identityId = *eventHub.Identity.UserAssignedIdentity
				}
				output["identity_id"] = identityId

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

func flattenIoTHubCloudToDevice(input *devices.CloudToDeviceProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})

	if maxDeliveryCount := input.MaxDeliveryCount; maxDeliveryCount != nil {
		output["max_delivery_count"] = *maxDeliveryCount
	}
	if defaultTimeToLive := input.DefaultTTLAsIso8601; defaultTimeToLive != nil {
		output["default_ttl"] = *defaultTimeToLive
	}

	output["feedback"] = flattenIoTHubCloudToDeviceFeedback(input.Feedback)

	return []interface{}{output}
}

func flattenIoTHubCloudToDeviceFeedback(input *devices.FeedbackProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	feedback := make(map[string]interface{})
	if feedbackMaxDeliveryCount := input.MaxDeliveryCount; feedbackMaxDeliveryCount != nil {
		feedback["max_delivery_count"] = *feedbackMaxDeliveryCount
	}
	if feedbackTimeToLive := input.TTLAsIso8601; feedbackTimeToLive != nil {
		feedback["time_to_live"] = *feedbackTimeToLive
	}
	if feedbackLockDuration := input.LockDurationAsIso8601; feedbackLockDuration != nil {
		feedback["lock_duration"] = *feedbackLockDuration
	}

	return []interface{}{feedback}
}

func expandNetworkRuleSetProperties(d *pluginsdk.ResourceData) *devices.NetworkRuleSetProperties {
	networkRuleSet := d.Get("network_rule_set").([]interface{})
	networkRuleSetProps := devices.NetworkRuleSetProperties{}
	nrsMap := networkRuleSet[0].(map[string]interface{})

	networkRuleSetProps.DefaultAction = devices.DefaultAction(nrsMap["default_action"].(string))
	networkRuleSetProps.ApplyToBuiltInEventHubEndpoint = utils.Bool(nrsMap["apply_to_builtin_eventhub_endpoint"].(bool))
	ipRules := nrsMap["ip_rule"].([]interface{})

	if len(ipRules) != 0 {
		rules := make([]devices.NetworkRuleSetIPRule, 0)

		for _, r := range ipRules {
			rawRule := r.(map[string]interface{})
			rule := &devices.NetworkRuleSetIPRule{
				FilterName: utils.String(rawRule["name"].(string)),
				Action:     devices.NetworkRuleIPAction(rawRule["action"].(string)),
				IPMask:     utils.String(rawRule["ip_mask"].(string)),
			}
			rules = append(rules, *rule)
		}
		networkRuleSetProps.IPRules = &rules
	}
	return &networkRuleSetProps
}

func flattenNetworkRuleSetProperties(input *devices.NetworkRuleSetProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})
	output["default_action"] = input.DefaultAction
	output["apply_to_builtin_eventhub_endpoint"] = input.ApplyToBuiltInEventHubEndpoint
	rules := make([]interface{}, 0)

	for _, r := range *input.IPRules {
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

	output["ip_rule"] = rules
	return []interface{}{output}
}

func expandIotHubIdentity(input []interface{}) (*devices.ArmIdentity, error) {
	config, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	identity := devices.ArmIdentity{
		Type: devices.ResourceIdentityType(config.Type),
	}

	if len(config.IdentityIds) != 0 {
		identityIds := make(map[string]*devices.ArmUserIdentity, len(config.IdentityIds))
		for id := range config.IdentityIds {
			identityIds[id] = &devices.ArmUserIdentity{}
		}
		identity.UserAssignedIdentities = identityIds
	}

	return &identity, nil
}

func flattenIotHubIdentity(input *devices.ArmIdentity) (*[]interface{}, error) {
	var transform *identity.SystemAndUserAssignedMap

	if input != nil {
		transform = &identity.SystemAndUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		for k, v := range input.UserAssignedIdentities {
			transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				ClientId:    v.ClientID,
				PrincipalId: v.PrincipalID,
			}
		}
		if input.PrincipalID != nil {
			transform.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transform.TenantId = *input.TenantID
		}
	}

	return identity.FlattenSystemAndUserAssignedMap(transform)
}

func fileUploadConnectionStringDiffSuppress(k, old, new string, d *pluginsdk.ResourceData) bool {
	// The access keys are always masked by Azure and the ordering of the parameters in the connection string
	// differs across services, so we will compare the fields individually instead.
	secretKeyRegex := regexp.MustCompile("(SharedAccessKey|AccountKey)=[^;]+")

	if secretKeyRegex.MatchString(new) {
		maskedNew := secretKeyRegex.ReplaceAllString(new, "$1=****")

		oldSplit := strings.Split(old, ";")
		newSplit := strings.Split(maskedNew, ";")

		sort.Strings(oldSplit)
		sort.Strings(newSplit)

		if len(oldSplit) != len(newSplit) {
			return false
		}

		for i := range oldSplit {
			if !strings.EqualFold(oldSplit[i], newSplit[i]) {
				return false
			}
		}
		return true
	}
	return false
}
