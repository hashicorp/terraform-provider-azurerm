// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	eventhubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/helpers"
	iothubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	servicebusValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func iothubSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the IotHub resource.",
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
						Description:  "The name of the SKU.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInEnumSlice(devices.PossibleIotHubSkuValues(), false),
					},

					"capacity": {
						Description:  "The number of provisioned IoT Hub units.",
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
						Description: "The name of the shared access policy.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"primary_key": {
						Description: "The primary key.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"secondary_key": {
						Description: "The secondary key.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"permissions": {
						Description: "The permissions assigned to the shared access policy.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"event_hub_partition_count": {
			Description:  "The number of device-to-cloud partitions used by backing Event Hubs.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      4,
			ValidateFunc: validation.IntBetween(2, 128),
		},
		"event_hub_retention_in_days": {
			Description:  "The Event Hub retention to use in days.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntBetween(1, 7),
		},

		"file_upload": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"connection_string": {
						Description:      "The connection string for the Azure Storage Account to which files are uploaded.",
						Type:             pluginsdk.TypeString,
						Required:         true,
						DiffSuppressFunc: helpers.FileUploadConnectionStringDiffSuppress,
						Sensitive:        true,
					},
					"container_name": {
						Description: "The name of the root container where the files should be uploaded to. The container need not exist but should be creatable using the connection_string specified.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"authentication_type": {
						Description:  "The type used to authenticate against the Storage Account.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(devices.AuthenticationTypeKeyBased),
						ValidateFunc: validation.StringInEnumSlice(devices.PossibleAuthenticationTypeValues(), false),
					},
					"identity_id": {
						Description:  "The ID of the User Managed Identity used to authenticate against the Storage Account.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
					"notifications": {
						Description: "Used to specify whether file notifications are sent to IoT Hub on upload.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"max_delivery_count": {
						Description:  "The number of times the IoT Hub attempts to deliver a file upload notification message.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      10,
						ValidateFunc: validation.IntBetween(1, 100),
					},
					"sas_ttl": {
						Description:  "The period of time for which the SAS URI generated by IoT Hub for file upload is valid, specified as an [ISO 8601 timespan duration](https://en.wikipedia.org/wiki/ISO_8601#Durations). This value must be between 1 minute and 24 hours.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.ISO8601Duration,
						Default:      "PT1H",
					},
					"default_ttl": {
						Description:  "The period of time for which a file upload notification message is available to consume before it expires, specified as an [ISO 8601 timespan duration](https://en.wikipedia.org/wiki/ISO_8601#Durations). This value must be between 1 minute and 48 hours.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.ISO8601Duration,
						Default:      "PT1H",
					},
					"lock_duration": {
						Description:  "The lock duration for the file upload notifications queue, specified as an [ISO 8601 timespan duration](https://en.wikipedia.org/wiki/ISO_8601#Durations). This value must be between 5 and 300 seconds.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.ISO8601Duration,
						Default:      "PT1M",
					},
				},
			},
		},

		"endpoint": {
			Type:       pluginsdk.TypeList,
			Optional:   true,
			Computed:   true, // azignore:AZS007 - pre-existing violation
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Description: "The type of the endpoint.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"AzureIotHub.StorageContainer",
							"AzureIotHub.ServiceBusQueue",
							"AzureIotHub.ServiceBusTopic",
							"AzureIotHub.EventHub",
						}, false),
					},

					"authentication_type": {
						Description:  "The type used to authenticate against the endpoint.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(devices.AuthenticationTypeKeyBased),
						ValidateFunc: validation.StringInEnumSlice(devices.PossibleAuthenticationTypeValues(), false),
					},

					"identity_id": {
						Description:  "The ID of the User Managed Identity used to authenticate against the endpoint.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},

					"endpoint_uri": {
						Description:  "URI of the Service Bus or Event Hubs Namespace endpoint. This attribute can only be specified and is mandatory when `authentication_type` is `identityBased` for endpoint type `AzureIotHub.ServiceBusQueue`, `AzureIotHub.ServiceBusTopic` or `AzureIotHub.EventHub`.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"entity_path": {
						Description:      "Name of the Service Bus Queue/Topic or Event Hub. This attribute can only be specified and is mandatory when `authentication_type` is `identityBased` for endpoint type `AzureIotHub.ServiceBusQueue`, `AzureIotHub.ServiceBusTopic` or `AzureIotHub.EventHub`.",
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
						Description:      "The connection string for the endpoint. This attribute is mandatory and can only be specified when `authentication_type` is `keyBased`.",
						Type:             pluginsdk.TypeString,
						Optional:         true,
						DiffSuppressFunc: IothubConnectionStringSuppress,
						Sensitive:        true,
					},

					"name": {
						Description:  "The name of the endpoint. The name must be unique across endpoint types. The following names are reserved: `events`, `operationsMonitoringEvents`, `fileNotifications` and `$default`.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: iothubValidate.IoTHubEndpointName,
					},

					"batch_frequency_in_seconds": {
						Description:      "Time interval at which blobs are written to storage. This attribute is applicable for endpoint type `AzureIotHub.StorageContainer`.",
						Type:             pluginsdk.TypeInt,
						Optional:         true,
						Default:          300,
						DiffSuppressFunc: suppressIfTypeIsNot("AzureIotHub.StorageContainer"),
						ValidateFunc:     validation.IntBetween(60, 720),
					},

					"max_chunk_size_in_bytes": {
						Description:      "Maximum number of bytes for each blob written to storage. This attribute is applicable for endpoint type `AzureIotHub.StorageContainer`.",
						Type:             pluginsdk.TypeInt,
						Optional:         true,
						Default:          314572800,
						DiffSuppressFunc: suppressIfTypeIsNot("AzureIotHub.StorageContainer"),
						ValidateFunc:     validation.IntBetween(10485760, 524288000),
					},

					"container_name": {
						Description:      "The name of Storage Container in the Storage Account. This attribute is mandatory for endpoint type `AzureIotHub.StorageContainer`.",
						Type:             pluginsdk.TypeString,
						Optional:         true,
						DiffSuppressFunc: suppressIfTypeIsNot("AzureIotHub.StorageContainer"),
					},

					"encoding": {
						Description:      "Encoding that is used to serialize messages to blobs. This attribute is applicable for endpoint type `AzureIotHub.StorageContainer`.",
						Type:             pluginsdk.TypeString,
						Optional:         true,
						ForceNew:         true,
						Default:          string(devices.EncodingAvro),
						DiffSuppressFunc: suppressIfTypeIsNot("AzureIotHub.StorageContainer"),
						ValidateFunc:     validation.StringInEnumSlice(devices.PossibleEncodingValues(), false),
					},

					"file_name_format": {
						Description:      "File name format for the blob. All parameters are mandatory but can be reordered. This attribute is applicable for endpoint type `AzureIotHub.StorageContainer`.",
						Type:             pluginsdk.TypeString,
						Optional:         true,
						Default:          "{iothub}/{partition}/{YYYY}/{MM}/{DD}/{HH}/{mm}",
						DiffSuppressFunc: suppressIfTypeIsNot("AzureIotHub.StorageContainer"),
						ValidateFunc:     iothubValidate.FileNameFormat,
					},

					"resource_group_name": commonschema.ResourceGroupNameOptional(),

					"subscription_id": {
						Description: "The subscription ID for the endpoint.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						// NOTE: O+C : required since this property would always be set even if it isn't specified in the tf config, otherwise it would cause a diff and break existing users
						Computed:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},

		"route": {
			Type:       pluginsdk.TypeList,
			Optional:   true,
			Computed:   true, // azignore:AZS007 - pre-existing violation
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description: "The name of the route.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^[-_.a-zA-Z0-9]{1,64}$"),
							"Route Name name can only include alphanumeric characters, periods, underscores, hyphens, has a maximum length of 64 characters, and must be unique.",
						),
					},
					"source": {
						Description: "The source that the routing rule is to be applied to, such as `DeviceMessages`.",
						Type:        pluginsdk.TypeString,
						Required:    true,
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
						Description: "The condition that is evaluated to apply the routing rule. For grammar, see: <https://learn.microsoft.com/azure/iot-hub/iot-hub-devguide-query-language>.",
						// The condition is a string value representing device-to-cloud message routes query expression
						// https://learn.microsoft.com/azure/iot-hub/iot-hub-devguide-query-language#device-to-cloud-message-routes-query-expressions
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "true",
					},
					"endpoint_names": {
						Description: "The list of endpoints to which messages that satisfy the condition are routed.",
						Type:        pluginsdk.TypeList,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Required: true,
					},
					"enabled": {
						Description: "Used to specify whether a route is enabled.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},
				},
			},
		},

		"enrichment": {
			Type: pluginsdk.TypeList,
			// Currently only 10 enrichments is allowed for standard or basic tier, 2 for Free tier.
			MaxItems:   10,
			Optional:   true,
			Computed:   true, // azignore:AZS007 - pre-existing violation
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key": {
						Description: "The key of the enrichment.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^[-_.a-zA-Z0-9]{1,64}$"),
							"Enrichment Key name can only include alphanumeric characters, periods, underscores, hyphens, has a maximum length of 64 characters, and must be unique.",
						),
					},
					"value": {
						Description:  "The value of the enrichment.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"endpoint_names": {
						Description: "The list of endpoints which will be enriched.",
						Type:        pluginsdk.TypeList,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Required: true,
					},
				},
			},
		},

		"fallback_route": {
			Description: "If the fallback route is enabled, messages that don't match any of the supplied routes are automatically sent to this route.",
			Type:        pluginsdk.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"source": {
						Description: "The source that the routing rule is to be applied to, such as `DeviceMessages`.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Default:     string(devices.RoutingSourceDeviceMessages),
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
						Description: "The condition that is evaluated to apply the routing rule. For grammar, see: <https://learn.microsoft.com/azure/iot-hub/iot-hub-devguide-query-language>.",
						// The condition is a string value representing device-to-cloud message routes query expression
						// https://learn.microsoft.com/azure/iot-hub/iot-hub-devguide-query-language#device-to-cloud-message-routes-query-expressions
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "true",
					},
					"endpoint_names": {
						Description: "The endpoints to which messages that satisfy the condition are routed. Currently only 1 endpoint is allowed.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringLenBetween(0, 64),
						},
					},
					"enabled": {
						Description: "Used to specify whether the fallback route is enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
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
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(devices.DefaultActionDeny),
						ValidateFunc: validation.StringInEnumSlice(devices.PossibleDefaultActionValues(), false),
					},
					"apply_to_builtin_eventhub_endpoint": {
						Description: "Whether Network Rule Set is also applied to the BuiltIn EventHub EndPoint of the IotHub.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"ip_rule": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description:  "The name of the IP rule.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: iothubValidate.IoTHubIpRuleName,
								},
								"ip_mask": {
									Description:  "The IP address range in CIDR notation for the IP rule.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.IsCIDRIPv4,
								},
								"action": {
									Description:  "The desired action for requests captured by this rule.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      string(devices.NetworkRuleIPActionAllow),
									ValidateFunc: validation.StringInEnumSlice(devices.PossibleNetworkRuleIPActionValues(), false),
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
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"max_delivery_count": {
						Description:  "The maximum delivery count for cloud-to-device per-device queues.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      10,
						ValidateFunc: validation.IntBetween(1, 100),
					},
					"default_ttl": {
						Description:  "This value must be between 1 minute and 48 hours.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "PT1H",
						ValidateFunc: validation.ISO8601DurationBetween("PT15M", "P2D"),
					},
					"feedback": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"time_to_live": {
									Description:  "The retention time for service-bound feedback messages, specified as an [ISO 8601 timespan duration](https://en.wikipedia.org/wiki/ISO_8601#Durations). This value must be between 1 minute and 48 hours.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      "PT1H",
									ValidateFunc: validation.ISO8601DurationBetween("PT15M", "P2D"),
								},
								"max_delivery_count": {
									Description:  "The maximum delivery count for the feedback queue.",
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      10,
									ValidateFunc: validation.IntBetween(1, 100),
								},
								"lock_duration": {
									Description:  "The lock duration for the feedback queue, specified as an [ISO 8601 timespan duration](https://en.wikipedia.org/wiki/ISO_8601#Durations). This value must be between 5 and 300 seconds.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      "PT60S",
									ValidateFunc: validation.ISO8601DurationBetween("PT5S", "PT300S"),
								},
							},
						},
					},
				},
			},
		},

		"local_authentication_enabled": {
			Description: "If false, SAS tokens with Iot hub scoped SAS keys cannot be used for authentication.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"min_tls_version": {
			Description: "The minimum TLS version to support for this hub.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     "1.2",
			ValidateFunc: validation.StringInSlice([]string{
				"1.2",
			}, false),
		},

		"public_network_access_enabled": {
			Description: "Is the IotHub resource accessible from a public network?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"hostname": {
			Description: "The hostname of the IotHub Resource.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"event_hub_events_endpoint": {
			Description: "The EventHub compatible endpoint for events data.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"event_hub_events_namespace": {
			Description: "The EventHub namespace for events data.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"event_hub_operations_endpoint": {
			Description: "The EventHub compatible endpoint for operational data.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"event_hub_events_path": {
			Description: "The EventHub compatible path for events data.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"event_hub_operations_path": {
			Description: "The EventHub compatible path for operational data.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"tags": commonschema.Tags(),
	}
}
