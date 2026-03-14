package iotoperations

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/dataflowendpoint"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataflowEndpointResource struct{}

var _ sdk.ResourceWithUpdate = DataflowEndpointResource{}

type DataflowEndpointModel struct {
	Name                    string                                `tfschema:"name"`
	ResourceGroupName       string                                `tfschema:"resource_group_name"`
	InstanceName            string                                `tfschema:"instance_name"`
	EndpointType            string                                `tfschema:"endpoint_type"`
	ExtendedLocationName    *string                               `tfschema:"extended_location_name"`
	ExtendedLocationType    *string                               `tfschema:"extended_location_type"`
	DataExplorerSettings    *DataflowEndpointDataExplorerModel    `tfschema:"data_explorer_settings"`
	DataLakeStorageSettings *DataflowEndpointDataLakeStorageModel `tfschema:"data_lake_storage_settings"`
	FabricOneLakeSettings   *DataflowEndpointFabricOneLakeModel   `tfschema:"fabric_one_lake_settings"`
	KafkaSettings           *DataflowEndpointKafkaModel           `tfschema:"kafka_settings"`
	LocalStorageSettings    *DataflowEndpointLocalStorageModel    `tfschema:"local_storage_settings"`
	MqttSettings            *DataflowEndpointMqttModel            `tfschema:"mqtt_settings"`
	ProvisioningState       *string                               `tfschema:"provisioning_state"`
}

type DataflowEndpointDataExplorerModel struct {
	Authentication DataflowEndpointDataExplorerAuthenticationModel `tfschema:"authentication"`
	Batching       *BatchingConfigurationModel                     `tfschema:"batching"`
	Database       string                                          `tfschema:"database"`
	Host           string                                          `tfschema:"host"`
}

type DataflowEndpointDataExplorerAuthenticationModel struct {
	Method                                string                                                       `tfschema:"method"`
	SystemAssignedManagedIdentitySettings *DataflowEndpointAuthenticationSystemAssignedManagedIdentity `tfschema:"system_assigned_managed_identity_settings"`
	UserAssignedManagedIdentitySettings   *DataflowEndpointAuthenticationUserAssignedManagedIdentity   `tfschema:"user_assigned_managed_identity_settings"`
}

type DataflowEndpointDataLakeStorageModel struct {
	Authentication DataflowEndpointDataLakeStorageAuthenticationModel `tfschema:"authentication"`
	Batching       *BatchingConfigurationModel                        `tfschema:"batching"`
	Host           string                                             `tfschema:"host"`
}

type DataflowEndpointDataLakeStorageAuthenticationModel struct {
	Method                                string                                                       `tfschema:"method"`
	AccessTokenSettings                   *DataflowEndpointAuthenticationAccessTokenModel              `tfschema:"access_token_settings"`
	SystemAssignedManagedIdentitySettings *DataflowEndpointAuthenticationSystemAssignedManagedIdentity `tfschema:"system_assigned_managed_identity_settings"`
	UserAssignedManagedIdentitySettings   *DataflowEndpointAuthenticationUserAssignedManagedIdentity   `tfschema:"user_assigned_managed_identity_settings"`
}

type DataflowEndpointFabricOneLakeModel struct {
	Authentication  DataflowEndpointFabricOneLakeAuthenticationModel `tfschema:"authentication"`
	Batching        *BatchingConfigurationModel                      `tfschema:"batching"`
	Host            string                                           `tfschema:"host"`
	Names           DataflowEndpointFabricOneLakeNamesModel          `tfschema:"names"`
	OneLakePathType string                                           `tfschema:"one_lake_path_type"`
}

type DataflowEndpointFabricOneLakeAuthenticationModel struct {
	Method                                string                                                       `tfschema:"method"`
	SystemAssignedManagedIdentitySettings *DataflowEndpointAuthenticationSystemAssignedManagedIdentity `tfschema:"system_assigned_managed_identity_settings"`
	UserAssignedManagedIdentitySettings   *DataflowEndpointAuthenticationUserAssignedManagedIdentity   `tfschema:"user_assigned_managed_identity_settings"`
}

type DataflowEndpointFabricOneLakeNamesModel struct {
	LakehouseName string `tfschema:"lakehouse_name"`
	WorkspaceName string `tfschema:"workspace_name"`
}

type DataflowEndpointKafkaModel struct {
	Authentication       DataflowEndpointKafkaAuthenticationModel `tfschema:"authentication"`
	Batching             *DataflowEndpointKafkaBatchingModel      `tfschema:"batching"`
	CloudEventAttributes *string                                  `tfschema:"cloud_event_attributes"`
	Compression          *string                                  `tfschema:"compression"`
	ConsumerGroupId      *string                                  `tfschema:"consumer_group_id"`
	CopyMqttProperties   *string                                  `tfschema:"copy_mqtt_properties"`
	Host                 string                                   `tfschema:"host"`
	KafkaAcks            *string                                  `tfschema:"kafka_acks"`
	PartitionStrategy    *string                                  `tfschema:"partition_strategy"`
	Tls                  *TlsPropertiesModel                      `tfschema:"tls"`
}

type DataflowEndpointKafkaAuthenticationModel struct {
	Method                                string                                                       `tfschema:"method"`
	SaslSettings                          *DataflowEndpointAuthenticationSaslModel                     `tfschema:"sasl_settings"`
	SystemAssignedManagedIdentitySettings *DataflowEndpointAuthenticationSystemAssignedManagedIdentity `tfschema:"system_assigned_managed_identity_settings"`
	UserAssignedManagedIdentitySettings   *DataflowEndpointAuthenticationUserAssignedManagedIdentity   `tfschema:"user_assigned_managed_identity_settings"`
	X509CertificateSettings               *DataflowEndpointAuthenticationX509Model                     `tfschema:"x509_certificate_settings"`
}

type DataflowEndpointKafkaBatchingModel struct {
	LatencyMs   *int64  `tfschema:"latency_ms"`
	MaxBytes    *int64  `tfschema:"max_bytes"`
	MaxMessages *int64  `tfschema:"max_messages"`
	Mode        *string `tfschema:"mode"`
}

type DataflowEndpointLocalStorageModel struct {
	PersistentVolumeClaimRef string `tfschema:"persistent_volume_claim_ref"`
}

type DataflowEndpointMqttModel struct {
	Authentication       DataflowEndpointMqttAuthenticationModel `tfschema:"authentication"`
	ClientIdPrefix       *string                                 `tfschema:"client_id_prefix"`
	CloudEventAttributes *string                                 `tfschema:"cloud_event_attributes"`
	Host                 *string                                 `tfschema:"host"`
	KeepAliveSeconds     *int64                                  `tfschema:"keep_alive_seconds"`
	MaxInflightMessages  *int64                                  `tfschema:"max_inflight_messages"`
	Protocol             *string                                 `tfschema:"protocol"`
	Qos                  *int64                                  `tfschema:"qos"`
	Retain               *string                                 `tfschema:"retain"`
	SessionExpirySeconds *int64                                  `tfschema:"session_expiry_seconds"`
	Tls                  *TlsPropertiesModel                     `tfschema:"tls"`
}

type DataflowEndpointMqttAuthenticationModel struct {
	Method                                string                                                       `tfschema:"method"`
	ServiceAccountTokenSettings           *DataflowEndpointAuthenticationServiceAccountTokenModel      `tfschema:"service_account_token_settings"`
	SystemAssignedManagedIdentitySettings *DataflowEndpointAuthenticationSystemAssignedManagedIdentity `tfschema:"system_assigned_managed_identity_settings"`
	UserAssignedManagedIdentitySettings   *DataflowEndpointAuthenticationUserAssignedManagedIdentity   `tfschema:"user_assigned_managed_identity_settings"`
	X509CertificateSettings               *DataflowEndpointAuthenticationX509Model                     `tfschema:"x509_certificate_settings"`
}

// Common authentication models
type DataflowEndpointAuthenticationAccessTokenModel struct {
	SecretRef string `tfschema:"secret_ref"`
}

type DataflowEndpointAuthenticationSaslModel struct {
	SaslType  string `tfschema:"sasl_type"`
	SecretRef string `tfschema:"secret_ref"`
}

type DataflowEndpointAuthenticationServiceAccountTokenModel struct {
	Audience string `tfschema:"audience"`
}

type DataflowEndpointAuthenticationSystemAssignedManagedIdentity struct {
	Audience *string `tfschema:"audience"`
}

type DataflowEndpointAuthenticationUserAssignedManagedIdentity struct {
	ClientId string  `tfschema:"client_id"`
	Scope    *string `tfschema:"scope"`
	TenantId string  `tfschema:"tenant_id"`
}

type DataflowEndpointAuthenticationX509Model struct {
	SecretRef string `tfschema:"secret_ref"`
}

type BatchingConfigurationModel struct {
	LatencySeconds *int64 `tfschema:"latency_seconds"`
	MaxMessages    *int64 `tfschema:"max_messages"`
}

type TlsPropertiesModel struct {
	Mode                             *string `tfschema:"mode"`
	TrustedCaCertificateConfigMapRef *string `tfschema:"trusted_ca_certificate_config_map_ref"`
}

func (r DataflowEndpointResource) ModelObject() interface{} {
	return &DataflowEndpointModel{}
}

func (r DataflowEndpointResource) ResourceType() string {
	return "azurerm_iotoperations_dataflow_endpoint"
}

func (r DataflowEndpointResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dataflowendpoint.ValidateDataflowEndpointID
}

func (r DataflowEndpointResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 63),
				validation.StringMatch(regexp.MustCompile("^[a-z0-9][a-z0-9-]*[a-z0-9]$"), "must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$"),
			),
		},
		"resource_group_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 90),
		},
		"instance_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 63),
				validation.StringMatch(regexp.MustCompile("^[a-z0-9][a-z0-9-]*[a-z0-9]$"), "must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$"),
			),
		},
		"endpoint_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"DataExplorer",
				"DataLakeStorage",
				"FabricOneLake",
				"Kafka",
				"LocalStorage",
				"Mqtt",
			}, false),
		},
		"extended_location": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"CustomLocation",
						}, false),
					},
				},
			},
		},
		"data_explorer_settings": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"data_lake_storage_settings", "fabric_one_lake_settings", "kafka_settings", "local_storage_settings", "mqtt_settings"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"database": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 253),
					},
					"host": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 253),
					},
					"batching": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem:     dataflowEndpointBatchingSchema(),
					},
					"authentication": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem:     dataflowEndpointAuthenticationSchema(),
					},
				},
			},
		},
		"data_lake_storage_settings": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"data_explorer_settings", "fabric_one_lake_settings", "kafka_settings", "local_storage_settings", "mqtt_settings"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"host": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 253),
					},
					"batching": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem:     dataflowEndpointBatchingSchema(),
					},
					"authentication": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem:     dataflowEndpointAuthenticationSchema(),
					},
				},
			},
		},
		"fabric_one_lake_settings": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"data_explorer_settings", "data_lake_storage_settings", "kafka_settings", "local_storage_settings", "mqtt_settings"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"host": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 253),
					},
					"names": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringLenBetween(1, 253),
						},
					},
					"one_lake_path_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Files",
							"Tables",
						}, false),
					},
					"workspace": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 253),
					},
					"batching": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem:     dataflowEndpointBatchingSchema(),
					},
					"authentication": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem:     dataflowEndpointAuthenticationSchema(),
					},
				},
			},
		},
		"kafka_settings": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"data_explorer_settings", "data_lake_storage_settings", "fabric_one_lake_settings", "local_storage_settings", "mqtt_settings"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"host": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 253),
					},
					"batching": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem:     dataflowEndpointBatchingSchema(),
					},
					"kafka": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"consumer_group_id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringLenBetween(1, 253),
								},
								"compression": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.StringInSlice([]string{
										"None",
										"Gzip",
										"Snappy",
										"Lz4",
									}, false),
								},
								"batching": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"mode": {
												Type:     pluginsdk.TypeString,
												Optional: true,
												ValidateFunc: validation.StringInSlice([]string{
													"Enabled",
													"Disabled",
												}, false),
											},
											"latency_ms": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(0, 3600000),
											},
											"max_bytes": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(1, 1073741824),
											},
											"max_messages": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(1, 1000000),
											},
										},
									},
								},
							},
						},
					},
					"authentication": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem:     dataflowEndpointAuthenticationSchema(),
					},
				},
			},
		},
		"local_storage_settings": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"data_explorer_settings", "data_lake_storage_settings", "fabric_one_lake_settings", "kafka_settings", "mqtt_settings"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"path": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 1000),
					},
				},
			},
		},
		"mqtt_settings": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"data_explorer_settings", "data_lake_storage_settings", "fabric_one_lake_settings", "kafka_settings", "local_storage_settings"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"host": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 253),
					},
					"keep_alive_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 65535),
					},
					"retain": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Keep",
							"Never",
						}, false),
					},
					"session_expiry_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 4294967295),
					},
					"max_inflight_messages": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 65535),
					},
					"qos": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 2),
					},
					"protocol": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Mqtt",
							"WebSockets",
						}, false),
					},
					"client_id_prefix": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringLenBetween(1, 253),
					},
					"tls_settings": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"mode": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										"Enabled",
										"Disabled",
									}, false),
								},
								"trusted_ca_certificate_config_map_ref": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringLenBetween(1, 253),
								},
							},
						},
					},
					"authentication": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem:     dataflowEndpointAuthenticationSchema(),
					},
				},
			},
		},
		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r DataflowEndpointResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"provisioning_state": {
			Type: pluginsdk.TypeString,
			// NOTE: O+C Azure automatically assigns provisioning state during resource lifecycle
			Computed: true,
		},
	}
}

func dataflowEndpointBatchingSchema() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"latency_seconds": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 3600),
			},
			"max_messages": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 1000000),
			},
		},
	}
}

func dataflowEndpointAuthenticationSchema() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"method": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SystemAssignedManagedIdentity",
					"UserAssignedManagedIdentity",
					"ServiceAccountToken",
					"X509Certificate",
					"AccessToken",
					"Sasl",
				}, false),
			},
			"system_assigned_managed_identity_settings": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"audience": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 253),
						},
					},
				},
			},
			"user_assigned_managed_identity_settings": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"client_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"audience": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 253),
						},
					},
				},
			},
			"service_account_token_settings": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"audience": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 253),
						},
					},
				},
			},
			"x509_certificate_settings": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"secret_ref": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 253),
						},
					},
				},
			},
			"access_token_settings": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"secret_ref": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 253),
						},
					},
				},
			},
			"sasl_settings": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"sasl_type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Plain",
								"ScramSha256",
								"ScramSha512",
							}, false),
						},
						"secret_ref": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 253),
						},
					},
				},
			},
		},
	}
}

func (r DataflowEndpointResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.DataflowEndpointClient

			var model DataflowEndpointModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Extract extended_location block from ResourceData and set model fields (value types)
			if v, ok := metadata.ResourceData.GetOk("extended_location"); ok {
				if list, ok := v.([]interface{}); ok && len(list) > 0 && list[0] != nil {
					m := list[0].(map[string]interface{})
					if name, ok := m["name"].(string); ok && name != "" {
						model.ExtendedLocationName = &name
					}
					if typ, ok := m["type"].(string); ok && typ != "" {
						model.ExtendedLocationType = &typ
					}
				}
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := dataflowendpoint.NewDataflowEndpointID(subscriptionId, model.ResourceGroupName, model.InstanceName, model.Name)

			// Build payload
			var extendedLocation *dataflowendpoint.ExtendedLocation
			if model.ExtendedLocationName != nil && model.ExtendedLocationType != nil {
				extendedLocation = &dataflowendpoint.ExtendedLocation{
					Name: *model.ExtendedLocationName,
					Type: dataflowendpoint.ExtendedLocationType(*model.ExtendedLocationType),
				}
			}
			var payload dataflowendpoint.DataflowEndpointResource
			if extendedLocation != nil {
				payload.ExtendedLocation = *extendedLocation
			}
			payload.Properties = expandDataflowEndpointProperties(model)

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DataflowEndpointResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.DataflowEndpointClient

			id, err := dataflowendpoint.ParseDataflowEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			model := DataflowEndpointModel{
				Name:              id.DataflowEndpointName,
				ResourceGroupName: id.ResourceGroupName,
				InstanceName:      id.InstanceName,
			}

			if respModel := resp.Model; respModel != nil {
				model.ExtendedLocationName = &respModel.ExtendedLocation.Name
				extendedLocationType := string(respModel.ExtendedLocation.Type)
				model.ExtendedLocationType = &extendedLocationType

				if respModel.Properties != nil {
					flattenDataflowEndpointProperties(respModel.Properties, &model)

					if respModel.Properties.ProvisioningState != nil {
						provisioningState := string(*respModel.Properties.ProvisioningState)
						model.ProvisioningState = &provisioningState
					}
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r DataflowEndpointResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.DataflowEndpointClient

			id, err := dataflowendpoint.ParseDataflowEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model DataflowEndpointModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Extract extended_location block from ResourceData and set model fields (value types)
			if v, ok := metadata.ResourceData.GetOk("extended_location"); ok {
				if list, ok := v.([]interface{}); ok && len(list) > 0 && list[0] != nil {
					m := list[0].(map[string]interface{})
					if name, ok := m["name"].(string); ok && name != "" {
						model.ExtendedLocationName = &name
					}
					if typ, ok := m["type"].(string); ok && typ != "" {
						model.ExtendedLocationType = &typ
					}
				}
			}

			// For dataflow endpoint, we use CreateOrUpdate for updates since there's no dedicated Update method
			var extendedLocation *dataflowendpoint.ExtendedLocation
			if model.ExtendedLocationName != nil && model.ExtendedLocationType != nil {
				extendedLocation = &dataflowendpoint.ExtendedLocation{
					Name: *model.ExtendedLocationName,
					Type: dataflowendpoint.ExtendedLocationType(*model.ExtendedLocationType),
				}
			}
			var payload dataflowendpoint.DataflowEndpointResource
			if extendedLocation != nil {
				payload.ExtendedLocation = *extendedLocation
			}
			payload.Properties = expandDataflowEndpointProperties(model)

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r DataflowEndpointResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.DataflowEndpointClient

			id, err := dataflowendpoint.ParseDataflowEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandDataflowEndpointProperties(model DataflowEndpointModel) *dataflowendpoint.DataflowEndpointProperties {
	props := &dataflowendpoint.DataflowEndpointProperties{
		EndpointType: dataflowendpoint.EndpointType(model.EndpointType),
	}

	if model.DataExplorerSettings != nil {
		props.DataExplorerSettings = expandDataflowEndpointDataExplorer(*model.DataExplorerSettings)
	}

	if model.DataLakeStorageSettings != nil {
		props.DataLakeStorageSettings = expandDataflowEndpointDataLakeStorage(*model.DataLakeStorageSettings)
	}

	if model.FabricOneLakeSettings != nil {
		props.FabricOneLakeSettings = expandDataflowEndpointFabricOneLake(*model.FabricOneLakeSettings)
	}

	if model.KafkaSettings != nil {
		props.KafkaSettings = expandDataflowEndpointKafka(*model.KafkaSettings)
	}

	if model.LocalStorageSettings != nil {
		props.LocalStorageSettings = expandDataflowEndpointLocalStorage(*model.LocalStorageSettings)
	}

	if model.MqttSettings != nil {
		props.MqttSettings = expandDataflowEndpointMqtt(*model.MqttSettings)
	}

	return props
}

func expandDataflowEndpointKafka(kafka DataflowEndpointKafkaModel) *dataflowendpoint.DataflowEndpointKafka {
	result := &dataflowendpoint.DataflowEndpointKafka{
		Host:           kafka.Host,
		Authentication: expandDataflowEndpointKafkaAuthentication(kafka.Authentication),
	}

	if kafka.Batching != nil {
		result.Batching = expandDataflowEndpointKafkaBatching(*kafka.Batching)
	}

	if kafka.CloudEventAttributes != nil {
		cloudEventType := dataflowendpoint.CloudEventAttributeType(*kafka.CloudEventAttributes)
		result.CloudEventAttributes = &cloudEventType
	}

	if kafka.Compression != nil {
		compression := dataflowendpoint.DataflowEndpointKafkaCompression(*kafka.Compression)
		result.Compression = &compression
	}

	if kafka.ConsumerGroupId != nil {
		result.ConsumerGroupId = kafka.ConsumerGroupId
	}

	if kafka.CopyMqttProperties != nil {
		copyMqtt := dataflowendpoint.OperationalMode(*kafka.CopyMqttProperties)
		result.CopyMqttProperties = &copyMqtt
	}

	if kafka.KafkaAcks != nil {
		acks := dataflowendpoint.DataflowEndpointKafkaAcks(*kafka.KafkaAcks)
		result.KafkaAcks = &acks
	}

	if kafka.PartitionStrategy != nil {
		strategy := dataflowendpoint.DataflowEndpointKafkaPartitionStrategy(*kafka.PartitionStrategy)
		result.PartitionStrategy = &strategy
	}

	if kafka.Tls != nil {
		result.Tls = expandTlsProperties(*kafka.Tls)
	}

	return result
}

func expandDataflowEndpointKafkaAuthentication(auth DataflowEndpointKafkaAuthenticationModel) dataflowendpoint.DataflowEndpointKafkaAuthentication {
	result := dataflowendpoint.DataflowEndpointKafkaAuthentication{
		Method: dataflowendpoint.KafkaAuthMethod(auth.Method),
	}

	if auth.SaslSettings != nil {
		result.SaslSettings = &dataflowendpoint.DataflowEndpointAuthenticationSasl{
			SaslType:  dataflowendpoint.DataflowEndpointAuthenticationSaslType(auth.SaslSettings.SaslType),
			SecretRef: auth.SaslSettings.SecretRef,
		}
	}

	if auth.SystemAssignedManagedIdentitySettings != nil {
		result.SystemAssignedManagedIdentitySettings = &dataflowendpoint.DataflowEndpointAuthenticationSystemAssignedManagedIdentity{
			Audience: auth.SystemAssignedManagedIdentitySettings.Audience,
		}
	}

	if auth.UserAssignedManagedIdentitySettings != nil {
		result.UserAssignedManagedIdentitySettings = &dataflowendpoint.DataflowEndpointAuthenticationUserAssignedManagedIdentity{
			ClientId: auth.UserAssignedManagedIdentitySettings.ClientId,
			TenantId: auth.UserAssignedManagedIdentitySettings.TenantId,
			Scope:    auth.UserAssignedManagedIdentitySettings.Scope,
		}
	}

	if auth.X509CertificateSettings != nil {
		result.X509CertificateSettings = &dataflowendpoint.DataflowEndpointAuthenticationX509{
			SecretRef: auth.X509CertificateSettings.SecretRef,
		}
	}

	return result
}

func expandDataflowEndpointKafkaBatching(batching DataflowEndpointKafkaBatchingModel) *dataflowendpoint.DataflowEndpointKafkaBatching {
	result := &dataflowendpoint.DataflowEndpointKafkaBatching{}

	if batching.LatencyMs != nil {
		result.LatencyMs = batching.LatencyMs
	}

	if batching.MaxBytes != nil {
		result.MaxBytes = batching.MaxBytes
	}

	if batching.MaxMessages != nil {
		result.MaxMessages = batching.MaxMessages
	}

	if batching.Mode != nil {
		mode := dataflowendpoint.OperationalMode(*batching.Mode)
		result.Mode = &mode
	}

	return result
}

func expandTlsProperties(tls TlsPropertiesModel) *dataflowendpoint.TlsProperties {
	result := &dataflowendpoint.TlsProperties{}

	if tls.Mode != nil {
		mode := dataflowendpoint.OperationalMode(*tls.Mode)
		result.Mode = &mode
	}

	if tls.TrustedCaCertificateConfigMapRef != nil {
		result.TrustedCaCertificateConfigMapRef = tls.TrustedCaCertificateConfigMapRef
	}

	return result
}

// Additional expand functions for other endpoint types would go here...
func expandDataflowEndpointDataExplorer(dataExplorer DataflowEndpointDataExplorerModel) *dataflowendpoint.DataflowEndpointDataExplorer {
	result := &dataflowendpoint.DataflowEndpointDataExplorer{
		Database:       dataExplorer.Database,
		Host:           dataExplorer.Host,
		Authentication: expandDataflowEndpointDataExplorerAuthentication(dataExplorer.Authentication),
	}

	if dataExplorer.Batching != nil {
		result.Batching = &dataflowendpoint.BatchingConfiguration{
			LatencySeconds: dataExplorer.Batching.LatencySeconds,
			MaxMessages:    dataExplorer.Batching.MaxMessages,
		}
	}

	return result
}

func expandDataflowEndpointDataExplorerAuthentication(auth DataflowEndpointDataExplorerAuthenticationModel) dataflowendpoint.DataflowEndpointDataExplorerAuthentication {
	result := dataflowendpoint.DataflowEndpointDataExplorerAuthentication{
		Method: dataflowendpoint.ManagedIdentityMethod(auth.Method),
	}

	if auth.SystemAssignedManagedIdentitySettings != nil {
		result.SystemAssignedManagedIdentitySettings = &dataflowendpoint.DataflowEndpointAuthenticationSystemAssignedManagedIdentity{
			Audience: auth.SystemAssignedManagedIdentitySettings.Audience,
		}
	}

	if auth.UserAssignedManagedIdentitySettings != nil {
		result.UserAssignedManagedIdentitySettings = &dataflowendpoint.DataflowEndpointAuthenticationUserAssignedManagedIdentity{
			ClientId: auth.UserAssignedManagedIdentitySettings.ClientId,
			TenantId: auth.UserAssignedManagedIdentitySettings.TenantId,
			Scope:    auth.UserAssignedManagedIdentitySettings.Scope,
		}
	}

	return result
}

func expandDataflowEndpointDataLakeStorage(dataLakeStorage DataflowEndpointDataLakeStorageModel) *dataflowendpoint.DataflowEndpointDataLakeStorage {
	result := &dataflowendpoint.DataflowEndpointDataLakeStorage{
		Host:           dataLakeStorage.Host,
		Authentication: expandDataflowEndpointDataLakeStorageAuthentication(dataLakeStorage.Authentication),
	}

	if dataLakeStorage.Batching != nil {
		result.Batching = &dataflowendpoint.BatchingConfiguration{
			LatencySeconds: dataLakeStorage.Batching.LatencySeconds,
			MaxMessages:    dataLakeStorage.Batching.MaxMessages,
		}
	}

	return result
}

func expandDataflowEndpointDataLakeStorageAuthentication(auth DataflowEndpointDataLakeStorageAuthenticationModel) dataflowendpoint.DataflowEndpointDataLakeStorageAuthentication {
	result := dataflowendpoint.DataflowEndpointDataLakeStorageAuthentication{
		Method: dataflowendpoint.DataLakeStorageAuthMethod(auth.Method),
	}

	if auth.AccessTokenSettings != nil {
		result.AccessTokenSettings = &dataflowendpoint.DataflowEndpointAuthenticationAccessToken{
			SecretRef: auth.AccessTokenSettings.SecretRef,
		}
	}

	if auth.SystemAssignedManagedIdentitySettings != nil {
		result.SystemAssignedManagedIdentitySettings = &dataflowendpoint.DataflowEndpointAuthenticationSystemAssignedManagedIdentity{
			Audience: auth.SystemAssignedManagedIdentitySettings.Audience,
		}
	}

	if auth.UserAssignedManagedIdentitySettings != nil {
		result.UserAssignedManagedIdentitySettings = &dataflowendpoint.DataflowEndpointAuthenticationUserAssignedManagedIdentity{
			ClientId: auth.UserAssignedManagedIdentitySettings.ClientId,
			TenantId: auth.UserAssignedManagedIdentitySettings.TenantId,
			Scope:    auth.UserAssignedManagedIdentitySettings.Scope,
		}
	}

	return result
}

func expandDataflowEndpointFabricOneLake(fabricOneLake DataflowEndpointFabricOneLakeModel) *dataflowendpoint.DataflowEndpointFabricOneLake {
	result := &dataflowendpoint.DataflowEndpointFabricOneLake{
		Host:            fabricOneLake.Host,
		OneLakePathType: dataflowendpoint.DataflowEndpointFabricPathType(fabricOneLake.OneLakePathType),
		Authentication:  expandDataflowEndpointFabricOneLakeAuthentication(fabricOneLake.Authentication),
		Names: dataflowendpoint.DataflowEndpointFabricOneLakeNames{
			LakehouseName: fabricOneLake.Names.LakehouseName,
			WorkspaceName: fabricOneLake.Names.WorkspaceName,
		},
	}

	if fabricOneLake.Batching != nil {
		result.Batching = &dataflowendpoint.BatchingConfiguration{
			LatencySeconds: fabricOneLake.Batching.LatencySeconds,
			MaxMessages:    fabricOneLake.Batching.MaxMessages,
		}
	}

	return result
}

func expandDataflowEndpointFabricOneLakeAuthentication(auth DataflowEndpointFabricOneLakeAuthenticationModel) dataflowendpoint.DataflowEndpointFabricOneLakeAuthentication {
	result := dataflowendpoint.DataflowEndpointFabricOneLakeAuthentication{
		Method: dataflowendpoint.ManagedIdentityMethod(auth.Method),
	}

	if auth.SystemAssignedManagedIdentitySettings != nil {
		result.SystemAssignedManagedIdentitySettings = &dataflowendpoint.DataflowEndpointAuthenticationSystemAssignedManagedIdentity{
			Audience: auth.SystemAssignedManagedIdentitySettings.Audience,
		}
	}

	if auth.UserAssignedManagedIdentitySettings != nil {
		result.UserAssignedManagedIdentitySettings = &dataflowendpoint.DataflowEndpointAuthenticationUserAssignedManagedIdentity{
			ClientId: auth.UserAssignedManagedIdentitySettings.ClientId,
			TenantId: auth.UserAssignedManagedIdentitySettings.TenantId,
			Scope:    auth.UserAssignedManagedIdentitySettings.Scope,
		}
	}

	return result
}

func expandDataflowEndpointLocalStorage(localStorage DataflowEndpointLocalStorageModel) *dataflowendpoint.DataflowEndpointLocalStorage {
	return &dataflowendpoint.DataflowEndpointLocalStorage{
		PersistentVolumeClaimRef: localStorage.PersistentVolumeClaimRef,
	}
}

func expandDataflowEndpointMqtt(mqtt DataflowEndpointMqttModel) *dataflowendpoint.DataflowEndpointMqtt {
	result := &dataflowendpoint.DataflowEndpointMqtt{
		Authentication: expandDataflowEndpointMqttAuthentication(mqtt.Authentication),
	}

	if mqtt.ClientIdPrefix != nil {
		result.ClientIdPrefix = mqtt.ClientIdPrefix
	}

	if mqtt.CloudEventAttributes != nil {
		cloudEventType := dataflowendpoint.CloudEventAttributeType(*mqtt.CloudEventAttributes)
		result.CloudEventAttributes = &cloudEventType
	}

	if mqtt.Host != nil {
		result.Host = mqtt.Host
	}

	if mqtt.KeepAliveSeconds != nil {
		result.KeepAliveSeconds = mqtt.KeepAliveSeconds
	}

	if mqtt.MaxInflightMessages != nil {
		result.MaxInflightMessages = mqtt.MaxInflightMessages
	}

	if mqtt.Protocol != nil {
		protocol := dataflowendpoint.BrokerProtocolType(*mqtt.Protocol)
		result.Protocol = &protocol
	}

	if mqtt.Qos != nil {
		result.Qos = mqtt.Qos
	}

	if mqtt.Retain != nil {
		retain := dataflowendpoint.MqttRetainType(*mqtt.Retain)
		result.Retain = &retain
	}

	if mqtt.SessionExpirySeconds != nil {
		result.SessionExpirySeconds = mqtt.SessionExpirySeconds
	}

	if mqtt.Tls != nil {
		result.Tls = expandTlsProperties(*mqtt.Tls)
	}

	return result
}

func expandDataflowEndpointMqttAuthentication(auth DataflowEndpointMqttAuthenticationModel) dataflowendpoint.DataflowEndpointMqttAuthentication {
	result := dataflowendpoint.DataflowEndpointMqttAuthentication{
		Method: dataflowendpoint.MqttAuthMethod(auth.Method),
	}

	if auth.ServiceAccountTokenSettings != nil {
		result.ServiceAccountTokenSettings = &dataflowendpoint.DataflowEndpointAuthenticationServiceAccountToken{
			Audience: auth.ServiceAccountTokenSettings.Audience,
		}
	}

	if auth.SystemAssignedManagedIdentitySettings != nil {
		result.SystemAssignedManagedIdentitySettings = &dataflowendpoint.DataflowEndpointAuthenticationSystemAssignedManagedIdentity{
			Audience: auth.SystemAssignedManagedIdentitySettings.Audience,
		}
	}

	if auth.UserAssignedManagedIdentitySettings != nil {
		result.UserAssignedManagedIdentitySettings = &dataflowendpoint.DataflowEndpointAuthenticationUserAssignedManagedIdentity{
			ClientId: auth.UserAssignedManagedIdentitySettings.ClientId,
			TenantId: auth.UserAssignedManagedIdentitySettings.TenantId,
			Scope:    auth.UserAssignedManagedIdentitySettings.Scope,
		}
	}

	if auth.X509CertificateSettings != nil {
		result.X509CertificateSettings = &dataflowendpoint.DataflowEndpointAuthenticationX509{
			SecretRef: auth.X509CertificateSettings.SecretRef,
		}
	}

	return result
}

func flattenDataflowEndpointProperties(props *dataflowendpoint.DataflowEndpointProperties, model *DataflowEndpointModel) {
	if props == nil {
		return
	}

	model.EndpointType = string(props.EndpointType)

	if props.DataExplorerSettings != nil {
		model.DataExplorerSettings = flattenDataflowEndpointDataExplorer(*props.DataExplorerSettings)
	}

	if props.DataLakeStorageSettings != nil {
		model.DataLakeStorageSettings = flattenDataflowEndpointDataLakeStorage(*props.DataLakeStorageSettings)
	}

	if props.FabricOneLakeSettings != nil {
		model.FabricOneLakeSettings = flattenDataflowEndpointFabricOneLake(*props.FabricOneLakeSettings)
	}

	if props.KafkaSettings != nil {
		model.KafkaSettings = flattenDataflowEndpointKafka(*props.KafkaSettings)
	}

	if props.LocalStorageSettings != nil {
		model.LocalStorageSettings = flattenDataflowEndpointLocalStorage(*props.LocalStorageSettings)
	}

	if props.MqttSettings != nil {
		model.MqttSettings = flattenDataflowEndpointMqtt(*props.MqttSettings)
	}
}

// Flatten functions would follow similar patterns...
func flattenDataflowEndpointKafka(kafka dataflowendpoint.DataflowEndpointKafka) *DataflowEndpointKafkaModel {
	result := &DataflowEndpointKafkaModel{
		Host:           kafka.Host,
		Authentication: flattenDataflowEndpointKafkaAuthentication(kafka.Authentication),
	}

	if kafka.Batching != nil {
		result.Batching = &DataflowEndpointKafkaBatchingModel{
			LatencyMs:   kafka.Batching.LatencyMs,
			MaxBytes:    kafka.Batching.MaxBytes,
			MaxMessages: kafka.Batching.MaxMessages,
		}
		if kafka.Batching.Mode != nil {
			mode := string(*kafka.Batching.Mode)
			result.Batching.Mode = &mode
		}
	}

	if kafka.CloudEventAttributes != nil {
		cloudEvent := string(*kafka.CloudEventAttributes)
		result.CloudEventAttributes = &cloudEvent
	}

	if kafka.Compression != nil {
		compression := string(*kafka.Compression)
		result.Compression = &compression
	}

	if kafka.ConsumerGroupId != nil {
		result.ConsumerGroupId = kafka.ConsumerGroupId
	}

	if kafka.CopyMqttProperties != nil {
		copyMqtt := string(*kafka.CopyMqttProperties)
		result.CopyMqttProperties = &copyMqtt
	}

	if kafka.KafkaAcks != nil {
		acks := string(*kafka.KafkaAcks)
		result.KafkaAcks = &acks
	}

	if kafka.PartitionStrategy != nil {
		strategy := string(*kafka.PartitionStrategy)
		result.PartitionStrategy = &strategy
	}

	if kafka.Tls != nil {
		result.Tls = flattenTlsProperties(*kafka.Tls)
	}

	return result
}

func flattenDataflowEndpointKafkaAuthentication(auth dataflowendpoint.DataflowEndpointKafkaAuthentication) DataflowEndpointKafkaAuthenticationModel {
	result := DataflowEndpointKafkaAuthenticationModel{
		Method: string(auth.Method),
	}

	if auth.SaslSettings != nil {
		result.SaslSettings = &DataflowEndpointAuthenticationSaslModel{
			SaslType:  string(auth.SaslSettings.SaslType),
			SecretRef: auth.SaslSettings.SecretRef,
		}
	}

	if auth.SystemAssignedManagedIdentitySettings != nil {
		result.SystemAssignedManagedIdentitySettings = &DataflowEndpointAuthenticationSystemAssignedManagedIdentity{
			Audience: auth.SystemAssignedManagedIdentitySettings.Audience,
		}
	}

	if auth.UserAssignedManagedIdentitySettings != nil {
		result.UserAssignedManagedIdentitySettings = &DataflowEndpointAuthenticationUserAssignedManagedIdentity{
			ClientId: auth.UserAssignedManagedIdentitySettings.ClientId,
			TenantId: auth.UserAssignedManagedIdentitySettings.TenantId,
			Scope:    auth.UserAssignedManagedIdentitySettings.Scope,
		}
	}

	if auth.X509CertificateSettings != nil {
		result.X509CertificateSettings = &DataflowEndpointAuthenticationX509Model{
			SecretRef: auth.X509CertificateSettings.SecretRef,
		}
	}

	return result
}

func flattenTlsProperties(tls dataflowendpoint.TlsProperties) *TlsPropertiesModel {
	result := &TlsPropertiesModel{}

	if tls.Mode != nil {
		mode := string(*tls.Mode)
		result.Mode = &mode
	}

	if tls.TrustedCaCertificateConfigMapRef != nil {
		result.TrustedCaCertificateConfigMapRef = tls.TrustedCaCertificateConfigMapRef
	}

	return result
}

// Additional flatten functions for other endpoint types would follow...
func flattenDataflowEndpointDataExplorer(dataExplorer dataflowendpoint.DataflowEndpointDataExplorer) *DataflowEndpointDataExplorerModel {
	result := &DataflowEndpointDataExplorerModel{
		Database:       dataExplorer.Database,
		Host:           dataExplorer.Host,
		Authentication: flattenDataflowEndpointDataExplorerAuthentication(dataExplorer.Authentication),
	}

	if dataExplorer.Batching != nil {
		result.Batching = &BatchingConfigurationModel{
			LatencySeconds: dataExplorer.Batching.LatencySeconds,
			MaxMessages:    dataExplorer.Batching.MaxMessages,
		}
	}

	return result
}

func flattenDataflowEndpointDataExplorerAuthentication(auth dataflowendpoint.DataflowEndpointDataExplorerAuthentication) DataflowEndpointDataExplorerAuthenticationModel {
	result := DataflowEndpointDataExplorerAuthenticationModel{
		Method: string(auth.Method),
	}

	if auth.SystemAssignedManagedIdentitySettings != nil {
		result.SystemAssignedManagedIdentitySettings = &DataflowEndpointAuthenticationSystemAssignedManagedIdentity{
			Audience: auth.SystemAssignedManagedIdentitySettings.Audience,
		}
	}

	if auth.UserAssignedManagedIdentitySettings != nil {
		result.UserAssignedManagedIdentitySettings = &DataflowEndpointAuthenticationUserAssignedManagedIdentity{
			ClientId: auth.UserAssignedManagedIdentitySettings.ClientId,
			TenantId: auth.UserAssignedManagedIdentitySettings.TenantId,
			Scope:    auth.UserAssignedManagedIdentitySettings.Scope,
		}
	}

	return result
}

func flattenDataflowEndpointDataLakeStorage(dataLakeStorage dataflowendpoint.DataflowEndpointDataLakeStorage) *DataflowEndpointDataLakeStorageModel {
	result := &DataflowEndpointDataLakeStorageModel{
		Host:           dataLakeStorage.Host,
		Authentication: flattenDataflowEndpointDataLakeStorageAuthentication(dataLakeStorage.Authentication),
	}

	if dataLakeStorage.Batching != nil {
		result.Batching = &BatchingConfigurationModel{
			LatencySeconds: dataLakeStorage.Batching.LatencySeconds,
			MaxMessages:    dataLakeStorage.Batching.MaxMessages,
		}
	}

	return result
}

func flattenDataflowEndpointDataLakeStorageAuthentication(auth dataflowendpoint.DataflowEndpointDataLakeStorageAuthentication) DataflowEndpointDataLakeStorageAuthenticationModel {
	result := DataflowEndpointDataLakeStorageAuthenticationModel{ // Add "Model" suffix
		Method: string(auth.Method),
	}

	if auth.AccessTokenSettings != nil {
		result.AccessTokenSettings = &DataflowEndpointAuthenticationAccessTokenModel{
			SecretRef: auth.AccessTokenSettings.SecretRef,
		}
	}

	if auth.SystemAssignedManagedIdentitySettings != nil {
		result.SystemAssignedManagedIdentitySettings = &DataflowEndpointAuthenticationSystemAssignedManagedIdentity{
			Audience: auth.SystemAssignedManagedIdentitySettings.Audience,
		}
	}

	if auth.UserAssignedManagedIdentitySettings != nil {
		result.UserAssignedManagedIdentitySettings = &DataflowEndpointAuthenticationUserAssignedManagedIdentity{
			ClientId: auth.UserAssignedManagedIdentitySettings.ClientId,
			TenantId: auth.UserAssignedManagedIdentitySettings.TenantId,
			Scope:    auth.UserAssignedManagedIdentitySettings.Scope,
		}
	}

	return result
}

func flattenDataflowEndpointFabricOneLake(fabricOneLake dataflowendpoint.DataflowEndpointFabricOneLake) *DataflowEndpointFabricOneLakeModel {
	result := &DataflowEndpointFabricOneLakeModel{
		Host:            fabricOneLake.Host,
		OneLakePathType: string(fabricOneLake.OneLakePathType),
		Authentication:  flattenDataflowEndpointFabricOneLakeAuthentication(fabricOneLake.Authentication),
		Names: DataflowEndpointFabricOneLakeNamesModel{
			LakehouseName: fabricOneLake.Names.LakehouseName,
			WorkspaceName: fabricOneLake.Names.WorkspaceName,
		},
	}

	if fabricOneLake.Batching != nil {
		result.Batching = &BatchingConfigurationModel{
			LatencySeconds: fabricOneLake.Batching.LatencySeconds,
			MaxMessages:    fabricOneLake.Batching.MaxMessages,
		}
	}

	return result
}

func flattenDataflowEndpointFabricOneLakeAuthentication(auth dataflowendpoint.DataflowEndpointFabricOneLakeAuthentication) DataflowEndpointFabricOneLakeAuthenticationModel {
	result := DataflowEndpointFabricOneLakeAuthenticationModel{
		Method: string(auth.Method),
	}

	if auth.SystemAssignedManagedIdentitySettings != nil {
		result.SystemAssignedManagedIdentitySettings = &DataflowEndpointAuthenticationSystemAssignedManagedIdentity{
			Audience: auth.SystemAssignedManagedIdentitySettings.Audience,
		}
	}

	if auth.UserAssignedManagedIdentitySettings != nil {
		result.UserAssignedManagedIdentitySettings = &DataflowEndpointAuthenticationUserAssignedManagedIdentity{
			ClientId: auth.UserAssignedManagedIdentitySettings.ClientId,
			TenantId: auth.UserAssignedManagedIdentitySettings.TenantId,
			Scope:    auth.UserAssignedManagedIdentitySettings.Scope,
		}
	}

	return result
}

func flattenDataflowEndpointLocalStorage(localStorage dataflowendpoint.DataflowEndpointLocalStorage) *DataflowEndpointLocalStorageModel {
	return &DataflowEndpointLocalStorageModel{
		PersistentVolumeClaimRef: localStorage.PersistentVolumeClaimRef,
	}
}

func flattenDataflowEndpointMqtt(mqtt dataflowendpoint.DataflowEndpointMqtt) *DataflowEndpointMqttModel {
	result := &DataflowEndpointMqttModel{
		Authentication: flattenDataflowEndpointMqttAuthentication(mqtt.Authentication),
	}

	if mqtt.ClientIdPrefix != nil {
		result.ClientIdPrefix = mqtt.ClientIdPrefix
	}

	if mqtt.CloudEventAttributes != nil {
		cloudEvent := string(*mqtt.CloudEventAttributes)
		result.CloudEventAttributes = &cloudEvent
	}

	if mqtt.Host != nil {
		result.Host = mqtt.Host
	}

	if mqtt.KeepAliveSeconds != nil {
		result.KeepAliveSeconds = mqtt.KeepAliveSeconds
	}

	if mqtt.MaxInflightMessages != nil {
		result.MaxInflightMessages = mqtt.MaxInflightMessages
	}

	if mqtt.Protocol != nil {
		protocol := string(*mqtt.Protocol)
		result.Protocol = &protocol
	}

	if mqtt.Qos != nil {
		result.Qos = mqtt.Qos
	}

	if mqtt.Retain != nil {
		retain := string(*mqtt.Retain)
		result.Retain = &retain
	}

	if mqtt.SessionExpirySeconds != nil {
		result.SessionExpirySeconds = mqtt.SessionExpirySeconds
	}

	if mqtt.Tls != nil {
		result.Tls = flattenTlsProperties(*mqtt.Tls)
	}

	return result
}

func flattenDataflowEndpointMqttAuthentication(auth dataflowendpoint.DataflowEndpointMqttAuthentication) DataflowEndpointMqttAuthenticationModel {
	result := DataflowEndpointMqttAuthenticationModel{
		Method: string(auth.Method),
	}

	if auth.ServiceAccountTokenSettings != nil {
		result.ServiceAccountTokenSettings = &DataflowEndpointAuthenticationServiceAccountTokenModel{
			Audience: auth.ServiceAccountTokenSettings.Audience,
		}
	}

	if auth.SystemAssignedManagedIdentitySettings != nil {
		result.SystemAssignedManagedIdentitySettings = &DataflowEndpointAuthenticationSystemAssignedManagedIdentity{
			Audience: auth.SystemAssignedManagedIdentitySettings.Audience,
		}
	}

	if auth.UserAssignedManagedIdentitySettings != nil {
		result.UserAssignedManagedIdentitySettings = &DataflowEndpointAuthenticationUserAssignedManagedIdentity{
			ClientId: auth.UserAssignedManagedIdentitySettings.ClientId,
			TenantId: auth.UserAssignedManagedIdentitySettings.TenantId,
			Scope:    auth.UserAssignedManagedIdentitySettings.Scope,
		}
	}

	if auth.X509CertificateSettings != nil {
		result.X509CertificateSettings = &DataflowEndpointAuthenticationX509Model{
			SecretRef: auth.X509CertificateSettings.SecretRef,
		}
	}

	return result
}
