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
	Name              string                              `tfschema:"name"`
	ResourceGroupName string                              `tfschema:"resource_group_name"`
	InstanceName      string                              `tfschema:"instance_name"`
	EndpointType      string                              `tfschema:"endpoint_type"`
	DataExplorerSettings *DataflowEndpointDataExplorerModel `tfschema:"data_explorer_settings"`
	DataLakeStorageSettings *DataflowEndpointDataLakeStorageModel `tfschema:"data_lake_storage_settings"`
	FabricOneLakeSettings *DataflowEndpointFabricOneLakeModel `tfschema:"fabric_one_lake_settings"`
	KafkaSettings     *DataflowEndpointKafkaModel         `tfschema:"kafka_settings"`
	LocalStorageSettings *DataflowEndpointLocalStorageModel `tfschema:"local_storage_settings"`
	MqttSettings      *DataflowEndpointMqttModel          `tfschema:"mqtt_settings"`
	ExtendedLocation  *ExtendedLocationModel              `tfschema:"extended_location"`
	Tags              map[string]string                   `tfschema:"tags"`
	ProvisioningState *string                             `tfschema:"provisioning_state"`
}

type DataflowEndpointDataExplorerModel struct {
	Database string                                `tfschema:"database"`
	Host     string                                `tfschema:"host"`
	Batching *DataflowEndpointBatchingModel        `tfschema:"batching"`
	Authentication *DataflowEndpointAuthenticationModel `tfschema:"authentication"`
}

type DataflowEndpointDataLakeStorageModel struct {
	Host       string                                `tfschema:"host"`
	Batching   *DataflowEndpointBatchingModel        `tfschema:"batching"`
	Authentication *DataflowEndpointAuthenticationModel `tfschema:"authentication"`
}

type DataflowEndpointFabricOneLakeModel struct {
	Host       string                                `tfschema:"host"`
	Names      []string                              `tfschema:"names"`
	OneLakePathType string                           `tfschema:"one_lake_path_type"`
	Workspace  string                                `tfschema:"workspace"`
	Batching   *DataflowEndpointBatchingModel        `tfschema:"batching"`
	Authentication *DataflowEndpointAuthenticationModel `tfschema:"authentication"`
}

type DataflowEndpointKafkaModel struct {
	Host            string                                `tfschema:"host"`
	Batching        *DataflowEndpointBatchingModel        `tfschema:"batching"`
	Kafka           *DataflowEndpointKafkaSettingsModel   `tfschema:"kafka"`
	Authentication  *DataflowEndpointAuthenticationModel  `tfschema:"authentication"`
}

type DataflowEndpointKafkaSettingsModel struct {
	ConsumerGroupId *string `tfschema:"consumer_group_id"`
	Compression     *string `tfschema:"compression"`
	Batching        *DataflowEndpointKafkaBatchingModel `tfschema:"batching"`
}

type DataflowEndpointKafkaBatchingModel struct {
	Mode         *string `tfschema:"mode"`
	LatencyMs    *int    `tfschema:"latency_ms"`
	MaxBytes     *int    `tfschema:"max_bytes"`
	MaxMessages  *int    `tfschema:"max_messages"`
}

type DataflowEndpointLocalStorageModel struct {
	Path string `tfschema:"path"`
}

type DataflowEndpointMqttModel struct {
	Host                 string                                `tfschema:"host"`
	KeepAliveSeconds     *int                                  `tfschema:"keep_alive_seconds"`
	Retain               *string                               `tfschema:"retain"`
	SessionExpirySeconds *int                                  `tfschema:"session_expiry_seconds"`
	MaxInflightMessages  *int                                  `tfschema:"max_inflight_messages"`
	Qos                  *int                                  `tfschema:"qos"`
	Protocol             *string                               `tfschema:"protocol"`
	ClientIdPrefix       *string                               `tfschema:"client_id_prefix"`
	TlsSettings          *DataflowEndpointMqttTlsModel         `tfschema:"tls_settings"`
	Authentication       *DataflowEndpointAuthenticationModel  `tfschema:"authentication"`
}

type DataflowEndpointMqttTlsModel struct {
	Mode                             string  `tfschema:"mode"`
	TrustedCaCertificateConfigMapRef *string `tfschema:"trusted_ca_certificate_config_map_ref"`
}

type DataflowEndpointBatchingModel struct {
	LatencySeconds *int `tfschema:"latency_seconds"`
	MaxMessages    *int `tfschema:"max_messages"`
}

type DataflowEndpointAuthenticationModel struct {
	Method                                string                                                  `tfschema:"method"`
	SystemAssignedManagedIdentitySettings *DataflowEndpointSystemAssignedManagedIdentityModel     `tfschema:"system_assigned_managed_identity_settings"`
	UserAssignedManagedIdentitySettings   *DataflowEndpointUserAssignedManagedIdentityModel       `tfschema:"user_assigned_managed_identity_settings"`
	ServiceAccountTokenSettings           *DataflowEndpointServiceAccountTokenModel               `tfschema:"service_account_token_settings"`
	X509CertificateSettings               *DataflowEndpointX509CertificateModel                   `tfschema:"x509_certificate_settings"`
	AccessTokenSettings                   *DataflowEndpointAccessTokenModel                       `tfschema:"access_token_settings"`
	SaslSettings                          *DataflowEndpointSaslModel                              `tfschema:"sasl_settings"`
}

type DataflowEndpointSystemAssignedManagedIdentityModel struct {
	Audience string `tfschema:"audience"`
}

type DataflowEndpointUserAssignedManagedIdentityModel struct {
	ClientId string `tfschema:"client_id"`
	Audience string `tfschema:"audience"`
}

type DataflowEndpointServiceAccountTokenModel struct {
	Audience string `tfschema:"audience"`
}

type DataflowEndpointX509CertificateModel struct {
	SecretRef string `tfschema:"secret_ref"`
}

type DataflowEndpointAccessTokenModel struct {
	SecretRef string `tfschema:"secret_ref"`
}

type DataflowEndpointSaslModel struct {
	SaslType  string `tfschema:"sasl_type"`
	SecretRef string `tfschema:"secret_ref"`
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
		"extended_location": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
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
			Type:     pluginsdk.TypeString,
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

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := dataflowendpoint.NewDataflowEndpointID(subscriptionId, model.ResourceGroupName, model.InstanceName, model.Name)

			// Build payload
			payload := dataflowendpoint.DataflowEndpointResource{
				Properties: expandDataflowEndpointProperties(model),
			}

			if model.ExtendedLocation != nil {
				payload.ExtendedLocation = expandExtendedLocation(model.ExtendedLocation)
			}

			if len(model.Tags) > 0 {
				payload.Tags = &model.Tags
			}

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
				if respModel.ExtendedLocation != nil {
					model.ExtendedLocation = flattenExtendedLocation(respModel.ExtendedLocation)
				}

				if respModel.Tags != nil {
					model.Tags = *respModel.Tags
				}

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

			// Check if anything actually changed before making API call
			if !metadata.ResourceData.HasChange("tags") && 
			   !hasDataflowEndpointSettingsChanged(metadata) {
				return nil
			}

			payload := dataflowendpoint.DataflowEndpointPatchModel{}
			hasChanges := false

			// Only include tags if they changed
			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = &model.Tags
				hasChanges = true
			}

			// Only include properties if endpoint settings changed
			if hasDataflowEndpointSettingsChanged(metadata) {
				payload.Properties = &dataflowendpoint.DataflowEndpointPropertiesPatch{
					EndpointType:            dataflowendpoint.EndpointType(model.EndpointType),
					DataExplorerSettings:    expandDataflowEndpointDataExplorerPatch(model.DataExplorerSettings),
					DataLakeStorageSettings: expandDataflowEndpointDataLakeStoragePatch(model.DataLakeStorageSettings),
					FabricOneLakeSettings:   expandDataflowEndpointFabricOneLakePatch(model.FabricOneLakeSettings),
					KafkaSettings:           expandDataflowEndpointKafkaPatch(model.KafkaSettings),
					LocalStorageSettings:    expandDataflowEndpointLocalStoragePatch(model.LocalStorageSettings),
					MqttSettings:            expandDataflowEndpointMqttPatch(model.MqttSettings),
				}
				hasChanges = true
			}

			// Only make API call if something actually changed
			if !hasChanges {
				return nil
			}

			if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
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

// Helper function to check if any endpoint settings changed
func hasDataflowEndpointSettingsChanged(metadata sdk.ResourceMetaData) bool {
	return metadata.ResourceData.HasChange("endpoint_type") ||
		   metadata.ResourceData.HasChange("data_explorer_settings") ||
		   metadata.ResourceData.HasChange("data_lake_storage_settings") ||
		   metadata.ResourceData.HasChange("fabric_one_lake_settings") ||
		   metadata.ResourceData.HasChange("kafka_settings") ||
		   metadata.ResourceData.HasChange("local_storage_settings") ||
		   metadata.ResourceData.HasChange("mqtt_settings")
}

// Helper functions for expand/flatten operations (simplified for brevity)
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

// Expand functions (simplified examples)
func expandDataflowEndpointDataExplorer(model DataflowEndpointDataExplorerModel) *dataflowendpoint.DataflowEndpointDataExplorer {
	result := &dataflowendpoint.DataflowEndpointDataExplorer{
		Database: model.Database,
		Host:     model.Host,
	}

	if model.Batching != nil {
		result.Batching = expandDataflowEndpointBatching(*model.Batching)
	}

	if model.Authentication != nil {
		result.Authentication = expandDataflowEndpointAuthentication(*model.Authentication)
	}

	return result
}

func expandDataflowEndpointDataLakeStorage(model DataflowEndpointDataLakeStorageModel) *dataflowendpoint.DataflowEndpointDataLakeStorage {
	result := &dataflowendpoint.DataflowEndpointDataLakeStorage{
		Host: model.Host,
	}

	if model.Batching != nil {
		result.Batching = expandDataflowEndpointBatching(*model.Batching)
	}

	if model.Authentication != nil {
		result.Authentication = expandDataflowEndpointAuthentication(*model.Authentication)
	}

	return result
}

func expandDataflowEndpointFabricOneLake(model DataflowEndpointFabricOneLakeModel) *dataflowendpoint.DataflowEndpointFabricOneLake {
	result := &dataflowendpoint.DataflowEndpointFabricOneLake{
		Host:            model.Host,
		Names:           model.Names,
		OneLakePathType: dataflowendpoint.FabricPathType(model.OneLakePathType),
		Workspace:       model.Workspace,
	}

	if model.Batching != nil {
		result.Batching = expandDataflowEndpointBatching(*model.Batching)
	}

	if model.Authentication != nil {
		result.Authentication = expandDataflowEndpointAuthentication(*model.Authentication)
	}

	return result
}

func expandDataflowEndpointKafka(model DataflowEndpointKafkaModel) *dataflowendpoint.DataflowEndpointKafka {
	result := &dataflowendpoint.DataflowEndpointKafka{
		Host: model.Host,
	}

	if model.Batching != nil {
		result.Batching = expandDataflowEndpointBatching(*model.Batching)
	}

	if model.Kafka != nil {
		result.Kafka = expandDataflowEndpointKafkaSettings(*model.Kafka)
	}

	if model.Authentication != nil {
		result.Authentication = expandDataflowEndpointAuthentication(*model.Authentication)
	}

	return result
}

func expandDataflowEndpointLocalStorage(model DataflowEndpointLocalStorageModel) *dataflowendpoint.DataflowEndpointLocalStorage {
	return &dataflowendpoint.DataflowEndpointLocalStorage{
		Path: model.Path,
	}
}

func expandDataflowEndpointMqtt(model DataflowEndpointMqttModel) *dataflowendpoint.DataflowEndpointMqtt {
	result := &dataflowendpoint.DataflowEndpointMqtt{
		Host: model.Host,
	}

	if model.KeepAliveSeconds != nil {
		result.KeepAliveSeconds = func(i int) *int64 { v := int64(i); return &v }(*model.KeepAliveSeconds)
	}

	if model.Retain != nil {
		retain := dataflowendpoint.MqttRetainType(*model.Retain)
		result.Retain = &retain
	}

	if model.SessionExpirySeconds != nil {
		result.SessionExpirySeconds = func(i int) *int64 { v := int64(i); return &v }(*model.SessionExpirySeconds)
	}

	if model.MaxInflightMessages != nil {
		result.MaxInflightMessages = func(i int) *int64 { v := int64(i); return &v }(*model.MaxInflightMessages)
	}

	if model.Qos != nil {
		result.Qos = func(i int) *int64 { v := int64(i); return &v }(*model.Qos)
	}

	if model.Protocol != nil {
		protocol := dataflowendpoint.BrokerProtocolType(*model.Protocol)
		result.Protocol = &protocol
	}

	if model.ClientIdPrefix != nil {
		result.ClientIdPrefix = model.ClientIdPrefix
	}

	if model.TlsSettings != nil {
		result.TlsSettings = expandDataflowEndpointMqttTls(*model.TlsSettings)
	}

	if model.Authentication != nil {
		result.Authentication = expandDataflowEndpointAuthentication(*model.Authentication)
	}

	return result
}

func expandDataflowEndpointBatching(model DataflowEndpointBatchingModel) *dataflowendpoint.BatchingConfiguration {
	result := &dataflowendpoint.BatchingConfiguration{}

	if model.LatencySeconds != nil {
		result.LatencySeconds = func(i int) *int64 { v := int64(i); return &v }(*model.LatencySeconds)
	}

	if model.MaxMessages != nil {
		result.MaxMessages = func(i int) *int64 { v := int64(i); return &v }(*model.MaxMessages)
	}

	return result
}

func expandDataflowEndpointAuthentication(model DataflowEndpointAuthenticationModel) *dataflowendpoint.DataflowEndpointAuthentication {
	result := &dataflowendpoint.DataflowEndpointAuthentication{
		Method: dataflowendpoint.DataflowEndpointAuthenticationMethod(model.Method),
	}

	if model.SystemAssignedManagedIdentitySettings != nil {
		result.SystemAssignedManagedIdentitySettings = &dataflowendpoint.DataflowEndpointAuthenticationSystemAssignedManagedIdentity{
			Audience: model.SystemAssignedManagedIdentitySettings.Audience,
		}
	}

	if model.UserAssignedManagedIdentitySettings != nil {
		result.UserAssignedManagedIdentitySettings = &dataflowendpoint.DataflowEndpointAuthenticationUserAssignedManagedIdentity{
			ClientId: model.UserAssignedManagedIdentitySettings.ClientId,
			Audience: model.UserAssignedManagedIdentitySettings.Audience,
		}
	}

	if model.ServiceAccountTokenSettings != nil {
		result.ServiceAccountTokenSettings = &dataflowendpoint.DataflowEndpointAuthenticationServiceAccountToken{
			Audience: model.ServiceAccountTokenSettings.Audience,
		}
	}

	if model.X509CertificateSettings != nil {
		result.X509CertificateSettings = &dataflowendpoint.DataflowEndpointAuthenticationX509{
			SecretRef: model.X509CertificateSettings.SecretRef,
		}
	}

	if model.AccessTokenSettings != nil {
		result.AccessTokenSettings = &dataflowendpoint.DataflowEndpointAuthenticationAccessToken{
			SecretRef: model.AccessTokenSettings.SecretRef,
		}
	}

	if model.SaslSettings != nil {
		result.SaslSettings = &dataflowendpoint.DataflowEndpointAuthenticationSasl{
			SaslType:  dataflowendpoint.DataflowEndpointAuthenticationSaslType(model.SaslSettings.SaslType),
			SecretRef: model.SaslSettings.SecretRef,
		}
	}

	return result
}

func expandDataflowEndpointKafkaSettings(model DataflowEndpointKafkaSettingsModel) *dataflowendpoint.DataflowEndpointKafkaSettings {
	result := &dataflowendpoint.DataflowEndpointKafkaSettings{}

	if model.ConsumerGroupId != nil {
		result.ConsumerGroupId = model.ConsumerGroupId
	}

	if model.Compression != nil {
		compression := dataflowendpoint.DataflowEndpointKafkaCompression(*model.Compression)
		result.Compression = &compression
	}

	if model.Batching != nil {
		result.Batching = expandDataflowEndpointKafkaBatching(*model.Batching)
	}

	return result
}

func expandDataflowEndpointKafkaBatching(model DataflowEndpointKafkaBatchingModel) *dataflowendpoint.DataflowEndpointKafkaBatching {
	result := &dataflowendpoint.DataflowEndpointKafkaBatching{}

	if model.Mode != nil {
		mode := dataflowendpoint.OperationalMode(*model.Mode)
		result.Mode = &mode
	}

	if model.LatencyMs != nil {
		result.LatencyMs = func(i int) *int64 { v := int64(i); return &v }(*model.LatencyMs)
	}

	if model.MaxBytes != nil {
		result.MaxBytes = func(i int) *int64 { v := int64(i); return &v }(*model.MaxBytes)
	}

	if model.MaxMessages != nil {
		result.MaxMessages = func(i int) *int64 { v := int64(i); return &v }(*model.MaxMessages)
	}

	return result
}

func expandDataflowEndpointMqttTls(model DataflowEndpointMqttTlsModel) *dataflowendpoint.TlsProperties {
	result := &dataflowendpoint.TlsProperties{
		Mode: dataflowendpoint.OperationalMode(model.Mode),
	}

	if model.TrustedCaCertificateConfigMapRef != nil {
		result.TrustedCaCertificateConfigMapRef = model.TrustedCaCertificateConfigMapRef
	}

	return result
}

// Flatten functions (simplified examples)
func flattenDataflowEndpointDataExplorer(settings dataflowendpoint.DataflowEndpointDataExplorer) *DataflowEndpointDataExplorerModel {
	result := &DataflowEndpointDataExplorerModel{
		Database: settings.Database,
		Host:     settings.Host,
	}

	if settings.Batching != nil {
		result.Batching = flattenDataflowEndpointBatching(*settings.Batching)
	}

	if settings.Authentication != nil {
		result.Authentication = flattenDataflowEndpointAuthentication(*settings.Authentication)
	}

	return result
}

func flattenDataflowEndpointDataLakeStorage(settings dataflowendpoint.DataflowEndpointDataLakeStorage) *DataflowEndpointDataLakeStorageModel {
	result := &DataflowEndpointDataLakeStorageModel{
		Host: settings.Host,
	}

	if settings.Batching != nil {
		result.Batching = flattenDataflowEndpointBatching(*settings.Batching)
	}

	if settings.Authentication != nil {
		result.Authentication = flattenDataflowEndpointAuthentication(*settings.Authentication)
	}

	return result
}

func flattenDataflowEndpointFabricOneLake(settings dataflowendpoint.DataflowEndpointFabricOneLake) *DataflowEndpointFabricOneLakeModel {
	result := &DataflowEndpointFabricOneLakeModel{
		Host:            settings.Host,
		Names:           settings.Names,
		OneLakePathType: string(settings.OneLakePathType),
		Workspace:       settings.Workspace,
	}

	if settings.Batching != nil {
		result.Batching = flattenDataflowEndpointBatching(*settings.Batching)
	}

	if settings.Authentication != nil {
		result.Authentication = flattenDataflowEndpointAuthentication(*settings.Authentication)
	}

	return result
}

func flattenDataflowEndpointKafka(settings dataflowendpoint.DataflowEndpointKafka) *DataflowEndpointKafkaModel {
	result := &DataflowEndpointKafkaModel{
		Host: settings.Host,
	}

	if settings.Batching != nil {
		result.Batching = flattenDataflowEndpointBatching(*settings.Batching)
	}

	if settings.Kafka != nil {
		result.Kafka = flattenDataflowEndpointKafkaSettings(*settings.Kafka)
	}

	if settings.Authentication != nil {
		result.Authentication = flattenDataflowEndpointAuthentication(*settings.Authentication)
	}

	return result
}

func flattenDataflowEndpointLocalStorage(settings dataflowendpoint.DataflowEndpointLocalStorage) *DataflowEndpointLocalStorageModel {
	return &DataflowEndpointLocalStorageModel{
		Path: settings.Path,
	}
}

func flattenDataflowEndpointMqtt(settings dataflowendpoint.DataflowEndpointMqtt) *DataflowEndpointMqttModel {
	result := &DataflowEndpointMqttModel{
		Host: settings.Host,
	}

	if settings.KeepAliveSeconds != nil {
		keepAlive := int(*settings.KeepAliveSeconds)
		result.KeepAliveSeconds = &keepAlive
	}

	if settings.Retain != nil {
		retain := string(*settings.Retain)
		result.Retain = &retain
	}

	if settings.SessionExpirySeconds != nil {
		sessionExpiry := int(*settings.SessionExpirySeconds)
		result.SessionExpirySeconds = &sessionExpiry
	}

	if settings.MaxInflightMessages != nil {
		maxInflight := int(*settings.MaxInflightMessages)
		result.MaxInflightMessages = &maxInflight
	}

	if settings.Qos != nil {
		qos := int(*settings.Qos)
		result.Qos = &qos
	}

	if settings.Protocol != nil {
		protocol := string(*settings.Protocol)
		result.Protocol = &protocol
	}

	if settings.ClientIdPrefix != nil {
		result.ClientIdPrefix = settings.ClientIdPrefix
	}

	if settings.TlsSettings != nil {
		result.TlsSettings = flattenDataflowEndpointMqttTls(*settings.TlsSettings)
	}

	if settings.Authentication != nil {
		result.Authentication = flattenDataflowEndpointAuthentication(*settings.Authentication)
	}

	return result
}

func flattenDataflowEndpointBatching(batching dataflowendpoint.BatchingConfiguration) *DataflowEndpointBatchingModel {
	result := &DataflowEndpointBatchingModel{}

	if batching.LatencySeconds != nil {
		latency := int(*batching.LatencySeconds)
		result.LatencySeconds = &latency
	}

	if batching.MaxMessages != nil {
		maxMessages := int(*batching.MaxMessages)
		result.MaxMessages = &maxMessages
	}

	return result
}

func flattenDataflowEndpointAuthentication(auth dataflowendpoint.DataflowEndpointAuthentication) *DataflowEndpointAuthenticationModel {
	result := &DataflowEndpointAuthenticationModel{
		Method: string(auth.Method),
	}

	if auth.SystemAssignedManagedIdentitySettings != nil {
		result.SystemAssignedManagedIdentitySettings = &DataflowEndpointSystemAssignedManagedIdentityModel{
			Audience: auth.SystemAssignedManagedIdentitySettings.Audience,
		}
	}

	if auth.UserAssignedManagedIdentitySettings != nil {
		result.UserAssignedManagedIdentitySettings = &DataflowEndpointUserAssignedManagedIdentityModel{
			ClientId: auth.UserAssignedManagedIdentitySettings.ClientId,
			Audience: auth.UserAssignedManagedIdentitySettings.Audience,
		}
	}

	if auth.ServiceAccountTokenSettings != nil {
		result.ServiceAccountTokenSettings = &DataflowEndpointServiceAccountTokenModel{
			Audience: auth.ServiceAccountTokenSettings.Audience,
		}
	}

	if auth.X509CertificateSettings != nil {
		result.X509CertificateSettings = &DataflowEndpointX509CertificateModel{
			SecretRef: auth.X509CertificateSettings.SecretRef,
		}
	}

	if auth.AccessTokenSettings != nil {
		result.AccessTokenSettings = &DataflowEndpointAccessTokenModel{
			SecretRef: auth.AccessTokenSettings.SecretRef,
		}
	}

	if auth.SaslSettings != nil {
		result.SaslSettings = &DataflowEndpointSaslModel{
			SaslType:  string(auth.SaslSettings.SaslType),
			SecretRef: auth.SaslSettings.SecretRef,
		}
	}

	return result
}

func flattenDataflowEndpointKafkaSettings(settings dataflowendpoint.DataflowEndpointKafkaSettings) *DataflowEndpointKafkaSettingsModel {
	result := &DataflowEndpointKafkaSettingsModel{}

	if settings.ConsumerGroupId != nil {
		result.ConsumerGroupId = settings.ConsumerGroupId
	}

	if settings.Compression != nil {
		compression := string(*settings.Compression)
		result.Compression = &compression
	}

	if settings.Batching != nil {
		result.Batching = flattenDataflowEndpointKafkaBatching(*settings.Batching)
	}

	return result
}

func flattenDataflowEndpointKafkaBatching(batching dataflowendpoint.DataflowEndpointKafkaBatching) *DataflowEndpointKafkaBatchingModel {
	result := &DataflowEndpointKafkaBatchingModel{}

	if batching.Mode != nil {
		mode := string(*batching.Mode)
		result.Mode = &mode
	}

	if batching.LatencyMs != nil {
		latency := int(*batching.LatencyMs)
		result.LatencyMs = &latency
	}

	if batching.MaxBytes != nil {
		maxBytes := int(*batching.MaxBytes)
		result.MaxBytes = &maxBytes
	}

	if batching.MaxMessages != nil {
		maxMessages := int(*batching.MaxMessages)
		result.MaxMessages = &maxMessages
	}

	return result
}

func flattenDataflowEndpointMqttTls(tls dataflowendpoint.TlsProperties) *DataflowEndpointMqttTlsModel {
	result := &DataflowEndpointMqttTlsModel{
		Mode: string(tls.Mode),
	}

	if tls.TrustedCaCertificateConfigMapRef != nil {
		result.TrustedCaCertificateConfigMapRef = tls.TrustedCaCertificateConfigMapRef
	}

	return result
}

// Simplified patch functions for update operations
func expandDataflowEndpointDataExplorerPatch(model *DataflowEndpointDataExplorerModel) *dataflowendpoint.DataflowEndpointDataExplorer {
	if model == nil {
		return nil
	}
	return expandDataflowEndpointDataExplorer(*model)
}

func expandDataflowEndpointDataLakeStoragePatch(model *DataflowEndpointDataLakeStorageModel) *dataflowendpoint.DataflowEndpointDataLakeStorage {
	if model == nil {
		return nil
	}
	return expandDataflowEndpointDataLakeStorage(*model)
}

func expandDataflowEndpointFabricOneLakePatch(model *DataflowEndpointFabricOneLakeModel) *dataflowendpoint.DataflowEndpointFabricOneLake {
	if model == nil {
		return nil
	}
	return expandDataflowEndpointFabricOneLake(*model)
}

func expandDataflowEndpointKafkaPatch(model *DataflowEndpointKafkaModel) *dataflowendpoint.DataflowEndpointKafka {
	if model == nil {
		return nil
	}
	return expandDataflowEndpointKafka(*model)
}

func expandDataflowEndpointLocalStoragePatch(model *DataflowEndpointLocalStorageModel) *dataflowendpoint.DataflowEndpointLocalStorage {
	if model == nil {
		return nil
	}
	return expandDataflowEndpointLocalStorage(*model)
}

func expandDataflowEndpointMqttPatch(model *DataflowEndpointMqttModel) *dataflowendpoint.DataflowEndpointMqtt {
	if model == nil {
		return nil
	}
	return expandDataflowEndpointMqtt(*model)
}