package iotoperations

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/broker"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type BrokerResource struct{}

var _ sdk.ResourceWithUpdate = BrokerResource{}

type BrokerModel struct {
	Name              string                 `tfschema:"name"`
	InstanceName      string                 `tfschema:"instance_name"`
	ResourceGroupName string                 `tfschema:"resource_group_name"`
	ExtendedLocation  *ExtendedLocationModel `tfschema:"extended_location"`
	Properties        *BrokerPropertiesModel `tfschema:"properties"`
	ProvisioningState *string                `tfschema:"provisioning_state"`
}

type ExtendedLocationModel struct {
	Name *string `tfschema:"name"`
	Type *string `tfschema:"type"`
}

type BrokerPropertiesModel struct {
	Advanced                *AdvancedSettingsModel        `tfschema:"advanced"`
	Cardinality             *CardinalityModel             `tfschema:"cardinality"`
	Diagnostics             *BrokerDiagnosticsModel       `tfschema:"diagnostics"`
	DiskBackedMessageBuffer *DiskBackedMessageBufferModel `tfschema:"disk_backed_message_buffer"`
	GenerateResourceLimits  *GenerateResourceLimitsModel  `tfschema:"generate_resource_limits"`
	MemoryProfile           *string                       `tfschema:"memory_profile"`
}

type AdvancedSettingsModel struct {
	Clients                *ClientConfigModel           `tfschema:"clients"`
	EncryptInternalTraffic *string                      `tfschema:"encrypt_internal_traffic"`
	InternalCerts          *CertManagerCertOptionsModel `tfschema:"internal_certs"`
}

type ClientConfigModel struct {
	MaxSessionExpirySeconds *int                       `tfschema:"max_session_expiry_seconds"`
	MaxMessageExpirySeconds *int                       `tfschema:"max_message_expiry_seconds"`
	MaxPacketSizeBytes      *int                       `tfschema:"max_packet_size_bytes"`
	SubscriberQueueLimit    *SubscriberQueueLimitModel `tfschema:"subscriber_queue_limit"`
	MaxReceiveMaximum       *int                       `tfschema:"max_receive_maximum"`
	MaxKeepAliveSeconds     *int                       `tfschema:"max_keep_alive_seconds"`
}

type SubscriberQueueLimitModel struct {
	Length   *int    `tfschema:"length"`
	Strategy *string `tfschema:"strategy"`
}

type CertManagerCertOptionsModel struct {
	Duration    *string                     `tfschema:"duration"`
	RenewBefore *string                     `tfschema:"renew_before"`
	PrivateKey  *CertManagerPrivateKeyModel `tfschema:"private_key"`
}

type CertManagerPrivateKeyModel struct {
	Algorithm      *string `tfschema:"algorithm"`
	RotationPolicy *string `tfschema:"rotation_policy"`
}

type CardinalityModel struct {
	BackendChain BackendChainModel `tfschema:"backend_chain"`
	Frontend     FrontendModel     `tfschema:"frontend"`
}

// Note: BackendChain fields are required in SDK, not optional
type BackendChainModel struct {
	Partitions       int  `tfschema:"partitions"`        // Required in SDK
	RedundancyFactor int  `tfschema:"redundancy_factor"` // Required in SDK
	Workers          *int `tfschema:"workers"`           // Optional in SDK
}

type FrontendModel struct {
	Replicas int  `tfschema:"replicas"` // Required in SDK
	Workers  *int `tfschema:"workers"`  // Optional in SDK
}

type BrokerDiagnosticsModel struct {
	Logs      *DiagnosticsLogsModel `tfschema:"logs"`
	Metrics   *MetricsModel         `tfschema:"metrics"`
	SelfCheck *SelfCheckModel       `tfschema:"self_check"`
	Traces    *TracesModel          `tfschema:"traces"`
}

type DiagnosticsLogsModel struct {
	Level *string `tfschema:"level"`
}

type MetricsModel struct {
	PrometheusPort *int `tfschema:"prometheus_port"`
}

type SelfCheckModel struct {
	Mode            *string `tfschema:"mode"`
	IntervalSeconds *int    `tfschema:"interval_seconds"`
	TimeoutSeconds  *int    `tfschema:"timeout_seconds"`
}

type TracesModel struct {
	Mode                *string           `tfschema:"mode"`
	CacheSizeMegabytes  *int              `tfschema:"cache_size_megabytes"`
	SelfTracing         *SelfTracingModel `tfschema:"self_tracing"`
	SpanChannelCapacity *int              `tfschema:"span_channel_capacity"`
}

type SelfTracingModel struct {
	Mode            *string `tfschema:"mode"`
	IntervalSeconds *int    `tfschema:"interval_seconds"`
}

type DiskBackedMessageBufferModel struct {
	MaxSize                   *string               `tfschema:"max_size"`
	EphemeralVolumeClaimSpec  *VolumeClaimSpecModel `tfschema:"ephemeral_volume_claim_spec"`
	PersistentVolumeClaimSpec *VolumeClaimSpecModel `tfschema:"persistent_volume_claim_spec"`
}

type GenerateResourceLimitsModel struct {
	Cpu *string `tfschema:"cpu"`
}

type VolumeClaimSpecModel struct {
	VolumeName       *string                    `tfschema:"volume_name"`
	VolumeMode       *string                    `tfschema:"volume_mode"`
	StorageClassName *string                    `tfschema:"storage_class_name"`
	AccessModes      []string                   `tfschema:"access_modes"`
	DataSource       *DataSourceModel           `tfschema:"data_source"`
	DataSourceRef    *DataSourceRefModel        `tfschema:"data_source_ref"`
	Resources        *ResourceRequirementsModel `tfschema:"resources"`
	Selector         *LabelSelectorModel        `tfschema:"selector"`
}

type DataSourceModel struct {
	ApiGroup *string `tfschema:"api_group"`
	Kind     *string `tfschema:"kind"`
	Name     *string `tfschema:"name"`
}

type DataSourceRefModel struct {
	ApiGroup  *string `tfschema:"api_group"`
	Kind      *string `tfschema:"kind"`
	Name      *string `tfschema:"name"`
	Namespace *string `tfschema:"namespace"`
}

type ResourceRequirementsModel struct {
	Limits   map[string]string `tfschema:"limits"`
	Requests map[string]string `tfschema:"requests"`
}

type LabelSelectorModel struct {
	MatchExpressions []MatchExpressionModel `tfschema:"match_expressions"`
	MatchLabels      map[string]string      `tfschema:"match_labels"`
}

type MatchExpressionModel struct {
	Key      *string  `tfschema:"key"`
	Operator *string  `tfschema:"operator"`
	Values   []string `tfschema:"values"`
}

func (r BrokerResource) ModelObject() interface{} {
	return &BrokerModel{}
}

func (r BrokerResource) ResourceType() string {
	return "azurerm_iotoperations_broker"
}

func (r BrokerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return broker.ValidateBrokerID
}

func (r BrokerResource) Arguments() map[string]*pluginsdk.Schema {
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
		"instance_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"resource_group_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 90),
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
					},
				},
			},
		},
		"properties": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"memory_profile": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Tiny",
							"Low",
							"Medium",
							"High",
						}, false),
					},
					"advanced": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"clients": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"max_session_expiry_seconds": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
											},
											"max_message_expiry_seconds": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
											},
											"max_packet_size_bytes": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
											},
											"subscriber_queue_limit": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"length": {
															Type:     pluginsdk.TypeInt,
															Optional: true,
														},
														"strategy": {
															Type:     pluginsdk.TypeString,
															Optional: true,
														},
													},
												},
											},
											"max_receive_maximum": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
											},
											"max_keep_alive_seconds": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
											},
										},
									},
								},
								"encrypt_internal_traffic": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"internal_certs": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"duration": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},
											"renew_before": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},
											"private_key": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"algorithm": {
															Type:     pluginsdk.TypeString,
															Optional: true,
														},
														"rotation_policy": {
															Type:     pluginsdk.TypeString,
															Optional: true,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					"cardinality": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"backend_chain": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"partitions": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(1, 16),
											},
											"redundancy_factor": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(1, 5),
											},
											"workers": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(1, 16),
											},
										},
									},
								},
								"frontend": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"replicas": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(1, 16),
											},
											"workers": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(1, 16),
											},
										},
									},
								},
							},
						},
					},
					"diagnostics": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"logs": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"level": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},
										},
									},
								},
								"metrics": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"prometheus_port": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(0, 65535),
											},
										},
									},
								},
								"self_check": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"mode": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},
											"interval_seconds": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
											},
											"timeout_seconds": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
											},
										},
									},
								},
								"traces": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"mode": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},
											"cache_size_megabytes": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
											},
											"self_tracing": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"mode": {
															Type:     pluginsdk.TypeString,
															Optional: true,
														},
														"interval_seconds": {
															Type:     pluginsdk.TypeInt,
															Optional: true,
														},
													},
												},
											},
											"span_channel_capacity": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"disk_backed_message_buffer": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"max_size": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"ephemeral_volume_claim_spec": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem:     volumeClaimSpecSchema(),
								},
								"persistent_volume_claim_spec": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem:     volumeClaimSpecSchema(),
								},
							},
						},
					},
					"generate_resource_limits": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"cpu": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r BrokerResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func volumeClaimSpecSchema() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"volume_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
			"volume_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
			"storage_class_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
			"access_modes": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},
			"data_source": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"api_group": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"kind": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"name": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
					},
				},
			},
			"data_source_ref": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"api_group": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"kind": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"name": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"namespace": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
					},
				},
			},
			"resources": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"limits": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"requests": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
					},
				},
			},
			"selector": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"match_expressions": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"key": {
										Type:     pluginsdk.TypeString,
										Optional: true,
									},
									"operator": {
										Type:     pluginsdk.TypeString,
										Optional: true,
									},
									"values": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
									},
								},
							},
						},
						"match_labels": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
					},
				},
			},
		},
	}
}

func (r BrokerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.BrokerClient

			var model BrokerModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := broker.NewBrokerID(subscriptionId, model.ResourceGroupName, model.InstanceName, model.Name)

			// Check if resource already exists
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("A resource with the ID %q already exists - to be managed via Terraform this resource needs to be imported into the State", id.ID())
			}

			// Build FULL payload for Create
			payload := broker.BrokerResource{
				Properties: expandBrokerProperties(model.Properties),
			}

			if model.ExtendedLocation != nil {
				payload.ExtendedLocation = *expandExtendedLocation(model.ExtendedLocation)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r BrokerResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.BrokerClient

			id, err := broker.ParseBrokerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model BrokerModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			//Check what actually changed using d.HasChange()
			hasChanges := false
			payload := broker.BrokerResource{}

			// Only include properties if they changed
			if metadata.ResourceData.HasChange("properties") {
				payload.Properties = expandBrokerProperties(model.Properties)
				hasChanges = true
			}

			// Only make API call if something actually changed
			if !hasChanges {
				return nil
			}

			// Get existing resource to preserve unchanged fields
			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", *id, err)
			}

			if existing.Model != nil {
				if payload.ExtendedLocation.Name == "" && existing.Model.ExtendedLocation.Name != "" {
					payload.ExtendedLocation = existing.Model.ExtendedLocation
				}
				// Preserve unchanged properties
				if payload.Properties == nil && existing.Model.Properties != nil {
					payload.Properties = existing.Model.Properties
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r BrokerResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.BrokerClient

			id, err := broker.ParseBrokerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			model := BrokerModel{
				Name:              id.BrokerName,
				InstanceName:      id.InstanceName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if resp.Model != nil {

				if resp.Model.ExtendedLocation.Name != "" {
					model.ExtendedLocation = flattenExtendedLocation(&resp.Model.ExtendedLocation)
				}

				if resp.Model.Properties != nil {
					model.Properties = flattenBrokerProperties(resp.Model.Properties)

					if resp.Model.Properties.ProvisioningState != nil {
						provState := string(*resp.Model.Properties.ProvisioningState)
						model.ProvisioningState = &provState
					}
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r BrokerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.BrokerClient

			id, err := broker.ParseBrokerID(metadata.ResourceData.Id())
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

// Helper functions for expanding and flattening data structures
func expandBrokerProperties(input *BrokerPropertiesModel) *broker.BrokerProperties {
	if input == nil {
		return nil
	}

	props := &broker.BrokerProperties{}

	if input.MemoryProfile != nil {
		memProfile := broker.BrokerMemoryProfile(*input.MemoryProfile)
		props.MemoryProfile = &memProfile
	}

	if input.Advanced != nil {
		props.Advanced = expandAdvancedSettings(input.Advanced)
	}

	if input.Cardinality != nil {
		props.Cardinality = expandCardinality(input.Cardinality)
	}

	if input.Diagnostics != nil {
		props.Diagnostics = expandBrokerDiagnostics(input.Diagnostics)
	}

	if input.DiskBackedMessageBuffer != nil {
		props.DiskBackedMessageBuffer = expandDiskBackedMessageBuffer(input.DiskBackedMessageBuffer)
	}

	if input.GenerateResourceLimits != nil {
		props.GenerateResourceLimits = expandGenerateResourceLimits(input.GenerateResourceLimits)
	}

	return props
}

func expandExtendedLocation(input *ExtendedLocationModel) *broker.ExtendedLocation {
	if input == nil {
		return nil
	}

	result := &broker.ExtendedLocation{}

	if input.Name != nil {
		result.Name = *input.Name
	}

	if input.Type != nil {
		extType := broker.ExtendedLocationType(*input.Type)
		result.Type = extType
	}

	return result
}

func expandAdvancedSettings(input *AdvancedSettingsModel) *broker.AdvancedSettings {
	if input == nil {
		return nil
	}

	result := &broker.AdvancedSettings{}

	if input.EncryptInternalTraffic != nil {
		opMode := broker.OperationalMode(*input.EncryptInternalTraffic)
		result.EncryptInternalTraffic = &opMode
	}

	if input.Clients != nil {
		result.Clients = expandClientConfig(input.Clients)
	}

	if input.InternalCerts != nil {
		result.InternalCerts = expandCertManagerCertOptions(input.InternalCerts)
	}

	return result
}

func expandCardinality(input *CardinalityModel) *broker.Cardinality {
	if input == nil {
		return nil
	}

	result := &broker.Cardinality{
		BackendChain: expandBackendChain(&input.BackendChain),
		Frontend:     expandFrontend(&input.Frontend),
	}

	return result
}

func expandBackendChain(input *BackendChainModel) broker.BackendChain {
	result := broker.BackendChain{
		Partitions:       int64(input.Partitions),
		RedundancyFactor: int64(input.RedundancyFactor),
	}

	if input.Workers != nil {
		workers := int64(*input.Workers)
		result.Workers = &workers
	}

	return result
}

func expandFrontend(input *FrontendModel) broker.Frontend {
	result := broker.Frontend{
		Replicas: int64(input.Replicas),
	}

	if input.Workers != nil {
		workers := int64(*input.Workers)
		result.Workers = &workers
	}

	return result
}

func expandClientConfig(input *ClientConfigModel) *broker.ClientConfig {
	if input == nil {
		return nil
	}
	// Implement based on the SDK's ClientConfig model
	return &broker.ClientConfig{}
}

func expandCertManagerCertOptions(input *CertManagerCertOptionsModel) *broker.CertManagerCertOptions {
	if input == nil {
		return nil
	}
	// Implement based on the SDK's CertManagerCertOptions model
	return &broker.CertManagerCertOptions{}
}

func expandBrokerDiagnostics(input *BrokerDiagnosticsModel) *broker.BrokerDiagnostics {
	if input == nil {
		return nil
	}
	// Implement based on the SDK's BrokerDiagnostics model
	return &broker.BrokerDiagnostics{}
}

func expandDiskBackedMessageBuffer(input *DiskBackedMessageBufferModel) *broker.DiskBackedMessageBuffer {
	if input == nil {
		return nil
	}
	// Implement based on the SDK's DiskBackedMessageBuffer model
	return &broker.DiskBackedMessageBuffer{}
}

func expandGenerateResourceLimits(input *GenerateResourceLimitsModel) *broker.GenerateResourceLimits {
	if input == nil {
		return nil
	}
	// Implement based on the SDK's GenerateResourceLimits model
	return &broker.GenerateResourceLimits{}
}

// Flatten functions for Read operations
func flattenExtendedLocation(input *broker.ExtendedLocation) *ExtendedLocationModel {
	if input == nil {
		return nil
	}

	result := &ExtendedLocationModel{}

	if input.Name != "" {
		result.Name = &input.Name
	}

	if input.Type != "" {
		extType := string(input.Type)
		result.Type = &extType
	}

	return result
}

func flattenBrokerProperties(input *broker.BrokerProperties) *BrokerPropertiesModel {
	if input == nil {
		return nil
	}

	result := &BrokerPropertiesModel{}

	if input.MemoryProfile != nil {
		memProfile := string(*input.MemoryProfile)
		result.MemoryProfile = &memProfile
	}

	if input.Advanced != nil {
		result.Advanced = flattenAdvancedSettings(input.Advanced)
	}

	if input.Cardinality != nil {
		result.Cardinality = flattenCardinality(input.Cardinality)
	}

	return result
}

func flattenAdvancedSettings(input *broker.AdvancedSettings) *AdvancedSettingsModel {
	if input == nil {
		return nil
	}

	result := &AdvancedSettingsModel{}

	if input.EncryptInternalTraffic != nil {
		opMode := string(*input.EncryptInternalTraffic)
		result.EncryptInternalTraffic = &opMode
	}

	return result
}

func flattenCardinality(input *broker.Cardinality) *CardinalityModel {
	if input == nil {
		return nil
	}

	result := &CardinalityModel{
		BackendChain: flattenBackendChain(&input.BackendChain),
		Frontend:     flattenFrontend(&input.Frontend),
	}

	return result
}

func flattenBackendChain(input *broker.BackendChain) BackendChainModel {
	result := BackendChainModel{
		Partitions:       int(input.Partitions),
		RedundancyFactor: int(input.RedundancyFactor),
	}

	if input.Workers != nil {
		workers := int(*input.Workers)
		result.Workers = &workers
	}

	return result
}

func flattenFrontend(input *broker.Frontend) FrontendModel {
	result := FrontendModel{
		Replicas: int(input.Replicas),
	}

	if input.Workers != nil {
		workers := int(*input.Workers)
		result.Workers = &workers
	}

	return result
}
