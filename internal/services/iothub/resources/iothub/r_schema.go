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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInEnumSlice(devices.PossibleIotHubSkuValues(), false),
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
			Default:      4,
			ValidateFunc: validation.IntBetween(2, 128),
		},
		"event_hub_retention_in_days": {
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
						Type:             pluginsdk.TypeString,
						Required:         true,
						DiffSuppressFunc: helpers.FileUploadConnectionStringDiffSuppress,
						Sensitive:        true,
					},
					"container_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"authentication_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(devices.AuthenticationTypeKeyBased),
						ValidateFunc: validation.StringInEnumSlice(devices.PossibleAuthenticationTypeValues(), false),
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
						ValidateFunc: validation.ISO8601Duration,
						Default:      "PT1H",
					},
					"default_ttl": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.ISO8601Duration,
						Default:      "PT1H",
					},
					"lock_duration": {
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
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(devices.AuthenticationTypeKeyBased),
						ValidateFunc: validation.StringInEnumSlice(devices.PossibleAuthenticationTypeValues(), false),
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
						Type:             pluginsdk.TypeString,
						Optional:         true,
						DiffSuppressFunc: IothubConnectionStringSuppress,
						Sensitive:        true,
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
						ValidateFunc:     validation.StringInEnumSlice(devices.PossibleEncodingValues(), false),
					},

					"file_name_format": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						Default:          "{iothub}/{partition}/{YYYY}/{MM}/{DD}/{HH}/{mm}",
						DiffSuppressFunc: suppressIfTypeIsNot("AzureIotHub.StorageContainer"),
						ValidateFunc:     iothubValidate.FileNameFormat,
					},

					"resource_group_name": commonschema.ResourceGroupNameOptional(),

					"subscription_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
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
						// https://docs.microsoft.com/azure/iot-hub/iot-hub-devguide-query-language#device-to-cloud-message-routes-query-expressions
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
			Computed:   true, // azignore:AZS007 - pre-existing violation
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
			Computed: true, // azignore:AZS007 - pre-existing violation
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
						// https://docs.microsoft.com/azure/iot-hub/iot-hub-devguide-query-language#device-to-cloud-message-routes-query-expressions
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "true",
					},
					"endpoint_names": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringLenBetween(0, 64),
						},
					},
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
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
									ValidateFunc: validation.IsCIDRIPv4,
								},
								"action": {
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
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      10,
						ValidateFunc: validation.IntBetween(1, 100),
					},
					"default_ttl": {
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
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      "PT1H",
									ValidateFunc: validation.ISO8601DurationBetween("PT15M", "P2D"),
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
									ValidateFunc: validation.ISO8601DurationBetween("PT5S", "PT300S"),
								},
							},
						},
					},
				},
			},
		},

		"local_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"min_tls_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "1.2",
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

		"tags": commonschema.Tags(),
	}
}
