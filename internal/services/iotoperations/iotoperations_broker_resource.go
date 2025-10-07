package iotoperations

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/broker"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type BrokerResource struct{}

var _ sdk.ResourceWithUpdate = BrokerResource{}

type BrokerModel struct {
	Name               string                    `tfschema:"name"`
	InstanceName       string                    `tfschema:"instance_name"`
	ResourceGroupName  string                    `tfschema:"resource_group_name"`
	Location           *string                   `tfschema:"location"`
	ExtendedLocation   *ExtendedLocationModel    `tfschema:"extended_location"`
	Properties         *BrokerPropertiesModel    `tfschema:"properties"`
	Tags               map[string]string         `tfschema:"tags"`
	ProvisioningState  *string                   `tfschema:"provisioning_state"`
}

type ExtendedLocationModel struct {
	Name *string `tfschema:"name"`
	Type *string `tfschema:"type"`
}

type BrokerPropertiesModel struct {
	Advanced                  *BrokerAdvancedModel               `tfschema:"advanced"`
	Cardinality               *BrokerCardinalityModel            `tfschema:"cardinality"`
	Diagnostics               *BrokerDiagnosticsModel            `tfschema:"diagnostics"`
	DiskBackedMessageBuffer   *BrokerDiskBackedModel             `tfschema:"disk_backed_message_buffer"`
	GenerateResourceLimits    *BrokerResourceLimitsModel         `tfschema:"generate_resource_limits"`
	MemoryProfile             *string                            `tfschema:"memory_profile"`
}

type BrokerAdvancedModel struct {
	Clients                 *BrokerClientsModel     `tfschema:"clients"`
	EncryptInternalTraffic  *string                 `tfschema:"encrypt_internal_traffic"`
	InternalCerts           *BrokerInternalCerts    `tfschema:"internal_certs"`
}

type BrokerClientsModel struct {
	MaxSessionExpirySeconds  *int                            `tfschema:"max_session_expiry_seconds"`
	MaxMessageExpirySeconds  *int                            `tfschema:"max_message_expiry_seconds"`
	MaxPacketSizeBytes       *int                            `tfschema:"max_packet_size_bytes"`
	SubscriberQueueLimit     *BrokerSubscriberQueueModel     `tfschema:"subscriber_queue_limit"`
	MaxReceiveMaximum        *int                            `tfschema:"max_receive_maximum"`
	MaxKeepAliveSeconds      *int                            `tfschema:"max_keep_alive_seconds"`
}

type BrokerSubscriberQueueModel struct {
	Length   *int    `tfschema:"length"`
	Strategy *string `tfschema:"strategy"`
}

type BrokerInternalCerts struct {
	Duration    *string                    `tfschema:"duration"`
	RenewBefore *string                    `tfschema:"renew_before"`
	PrivateKey  *BrokerPrivateKeyModel     `tfschema:"private_key"`
}

type BrokerPrivateKeyModel struct {
	Algorithm      *string `tfschema:"algorithm"`
	RotationPolicy *string `tfschema:"rotation_policy"`
}

type BrokerCardinalityModel struct {
	BackendChain *BrokerBackendChainModel `tfschema:"backend_chain"`
	Frontend     *BrokerFrontendModel     `tfschema:"frontend"`
}

//  backend_chain should be a single object, not a list
type BrokerBackendChainModel struct {
	Partitions       *int `tfschema:"partitions"`
	RedundancyFactor *int `tfschema:"redundancy_factor"`
	Workers          *int `tfschema:"workers"`
}

type BrokerFrontendModel struct {
	Replicas *int `tfschema:"replicas"`
	Workers  *int `tfschema:"workers"`
}

type BrokerDiagnosticsModel struct {
	Logs      *BrokerLogsModel      `tfschema:"logs"`
	Metrics   *BrokerMetricsModel   `tfschema:"metrics"`
	SelfCheck *BrokerSelfCheckModel `tfschema:"self_check"`
	Traces    *BrokerTracesModel    `tfschema:"traces"`
}

type BrokerLogsModel struct {
	Level *string `tfschema:"level"`
}

type BrokerMetricsModel struct {
	PrometheusPort *int `tfschema:"prometheus_port"`
}

type BrokerSelfCheckModel struct {
	Mode            *string `tfschema:"mode"`
	IntervalSeconds *int    `tfschema:"interval_seconds"`
	TimeoutSeconds  *int    `tfschema:"timeout_seconds"`
}

type BrokerTracesModel struct {
	Mode                *string                    `tfschema:"mode"`
	CacheSizeMegabytes  *int                       `tfschema:"cache_size_megabytes"`
	SelfTracing         *BrokerSelfTracingModel    `tfschema:"self_tracing"`
	SpanChannelCapacity *int                       `tfschema:"span_channel_capacity"`
}

type BrokerSelfTracingModel struct {
	Mode            *string `tfschema:"mode"`
	IntervalSeconds *int    `tfschema:"interval_seconds"`
}

type BrokerDiskBackedModel struct {
	MaxSize                       *string                      `tfschema:"max_size"`
	EphemeralVolumeClaimSpec      *VolumeClaimSpecModel        `tfschema:"ephemeral_volume_claim_spec"`
	PersistentVolumeClaimSpec     *VolumeClaimSpecModel        `tfschema:"persistent_volume_claim_spec"`
}

type VolumeClaimSpecModel struct {
	VolumeName       *string                        `tfschema:"volume_name"`
	VolumeMode       *string                        `tfschema:"volume_mode"`
	StorageClassName *string                        `tfschema:"storage_class_name"`
	AccessModes      []string                       `tfschema:"access_modes"`
	DataSource       *DataSourceModel               `tfschema:"data_source"`
	DataSourceRef    *DataSourceRefModel            `tfschema:"data_source_ref"`
	Resources        *ResourceRequirementsModel     `tfschema:"resources"`
	Selector         *LabelSelectorModel            `tfschema:"selector"`
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

type BrokerResourceLimitsModel struct {
	Cpu *string `tfschema:"cpu"`
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
		"location": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
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
					"memory_profile": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},
		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			// NOTE: O+C Azure automatically assigns provisioning state during resource lifecycle
			Computed: true,
		},
		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
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

			// Build payload
			payload := broker.BrokerResource{
				Properties: expandBrokerProperties(model.Properties),
			}

			if model.Location != nil {
				payload.Location = model.Location
			}

			if model.ExtendedLocation != nil {
				payload.ExtendedLocation = expandExtendedLocation(model.ExtendedLocation)
			}

			if len(model.Tags) > 0 {
				payload.Tags = &model.Tags
			}

			if err := client.CreateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
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
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			model := BrokerModel{
				Name:              id.BrokerName,
				InstanceName:      id.InstanceName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if resp.Model != nil {
				if resp.Model.Location != nil {
					model.Location = resp.Model.Location
				}

				if resp.Model.ExtendedLocation != nil {
					model.ExtendedLocation = flattenExtendedLocation(resp.Model.ExtendedLocation)
				}

				if resp.Model.Properties != nil {
					model.Properties = flattenBrokerProperties(resp.Model.Properties)
					
					if resp.Model.Properties.ProvisioningState != nil {
						provState := string(*resp.Model.Properties.ProvisioningState)
						model.ProvisioningState = &provState
					}
				}

				if resp.Model.Tags != nil {
					model.Tags = *resp.Model.Tags
				}
			}

			return metadata.Encode(&model)
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

			// Check if anything actually changed before making API call
			if !metadata.ResourceData.HasChange("tags") && 
			   !metadata.ResourceData.HasChange("properties") {
				return nil
			}

			payload := broker.BrokerPatchModel{}
			hasChanges := false

			// Only include tags if they changed
			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = &model.Tags
				hasChanges = true
			}

			// Only include properties if they changed
			if metadata.ResourceData.HasChange("properties") {
				payload.Properties = expandBrokerPropertiesPatch(model.Properties)
				hasChanges = true
			}

			// Only make API call if something actually changed
			if !hasChanges {
				return nil
			}

			if _, err := client.Update(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
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
		props.MemoryProfile = input.MemoryProfile
	}

	if input.Advanced != nil {
		props.Advanced = expandBrokerAdvanced(input.Advanced)
	}

	if input.Cardinality != nil {
		props.Cardinality = expandBrokerCardinality(input.Cardinality)
	}

	if input.Diagnostics != nil {
		props.Diagnostics = expandBrokerDiagnostics(input.Diagnostics)
	}

	if input.DiskBackedMessageBuffer != nil {
		props.DiskBackedMessageBuffer = expandBrokerDiskBacked(input.DiskBackedMessageBuffer)
	}

	if input.GenerateResourceLimits != nil {
		props.GenerateResourceLimits = expandBrokerResourceLimits(input.GenerateResourceLimits)
	}

	return props
}

func expandBrokerPropertiesPatch(input *BrokerPropertiesModel) *broker.BrokerPropertiesPatch {
	if input == nil {
		return nil
	}

	props := &broker.BrokerPropertiesPatch{}

	if input.MemoryProfile != nil {
		props.MemoryProfile = input.MemoryProfile
	}

	return props
}

func expandExtendedLocation(input *ExtendedLocationModel) *broker.ExtendedLocation {
	if input == nil {
		return nil
	}

	result := &broker.ExtendedLocation{}

	if input.Name != nil {
		result.Name = input.Name
	}

	if input.Type != nil {
		extType := broker.ExtendedLocationType(*input.Type)
		result.Type = &extType
	}

	return result
}

func expandBrokerAdvanced(input *BrokerAdvancedModel) *broker.BrokerAdvanced {
	if input == nil {
		return nil
	}

	result := &broker.BrokerAdvanced{}

	if input.EncryptInternalTraffic != nil {
		result.EncryptInternalTraffic = input.EncryptInternalTraffic
	}


	return result
}

func expandBrokerCardinality(input *BrokerCardinalityModel) *broker.BrokerCardinality {
	if input == nil {
		return nil
	}

	result := &broker.BrokerCardinality{}

	if input.BackendChain != nil {
		result.BackendChain = expandBrokerBackendChain(input.BackendChain)
	}

	if input.Frontend != nil {
		result.Frontend = expandBrokerFrontend(input.Frontend)
	}

	return result
}

func expandBrokerBackendChain(input *BrokerBackendChainModel) *broker.BrokerBackendChain {
	if input == nil {
		return nil
	}

	result := &broker.BrokerBackendChain{}

	if input.Partitions != nil {
		partitions := int64(*input.Partitions)
		result.Partitions = &partitions
	}

	if input.RedundancyFactor != nil {
		redundancy := int64(*input.RedundancyFactor)
		result.RedundancyFactor = &redundancy
	}

	if input.Workers != nil {
		workers := int64(*input.Workers)
		result.Workers = &workers
	}

	return result
}

func expandBrokerFrontend(input *BrokerFrontendModel) *broker.BrokerFrontend {
	if input == nil {
		return nil
	}

	result := &broker.BrokerFrontend{}

	if input.Replicas != nil {
		replicas := int64(*input.Replicas)
		result.Replicas = &replicas
	}

	if input.Workers != nil {
		workers := int64(*input.Workers)
		result.Workers = &workers
	}

	return result
}

func expandBrokerDiagnostics(input *BrokerDiagnosticsModel) *broker.BrokerDiagnostics {
	if input == nil {
		return nil
	}

	result := &broker.BrokerDiagnostics{}

	return result
}

func expandBrokerDiskBacked(input *BrokerDiskBackedModel) *broker.BrokerDiskBacked {
	if input == nil {
		return nil
	}

	result := &broker.BrokerDiskBacked{}

	if input.MaxSize != nil {
		result.MaxSize = input.MaxSize
	}

	return result
}

func expandBrokerResourceLimits(input *BrokerResourceLimitsModel) *broker.BrokerResourceLimits {
	if input == nil {
		return nil
	}

	result := &broker.BrokerResourceLimits{}

	if input.Cpu != nil {
		result.Cpu = input.Cpu
	}

	return result
}

// Flatten functions for Read operations
func flattenExtendedLocation(input *broker.ExtendedLocation) *ExtendedLocationModel {
	if input == nil {
		return nil
	}

	result := &ExtendedLocationModel{}

	if input.Name != nil {
		result.Name = input.Name
	}

	if input.Type != nil {
		extType := string(*input.Type)
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
		result.MemoryProfile = input.MemoryProfile
	}

	if input.Advanced != nil {
		result.Advanced = flattenBrokerAdvanced(input.Advanced)
	}

	if input.Cardinality != nil {
		result.Cardinality = flattenBrokerCardinality(input.Cardinality)
	}

	return result
}

func flattenBrokerAdvanced(input *broker.BrokerAdvanced) *BrokerAdvancedModel {
	if input == nil {
		return nil
	}

	result := &BrokerAdvancedModel{}

	if input.EncryptInternalTraffic != nil {
		result.EncryptInternalTraffic = input.EncryptInternalTraffic
	}



	return result
}

func flattenBrokerCardinality(input *broker.BrokerCardinality) *BrokerCardinalityModel {
	if input == nil {
		return nil
	}

	result := &BrokerCardinalityModel{}

	if input.BackendChain != nil {
		result.BackendChain = flattenBrokerBackendChain(input.BackendChain)
	}

	if input.Frontend != nil {
		result.Frontend = flattenBrokerFrontend(input.Frontend)
	}

	return result
}

func flattenBrokerBackendChain(input *broker.BrokerBackendChain) *BrokerBackendChainModel {
	if input == nil {
		return nil
	}

	result := &BrokerBackendChainModel{}

	if input.Partitions != nil {
		partitions := int(*input.Partitions)
		result.Partitions = &partitions
	}

	if input.RedundancyFactor != nil {
		redundancy := int(*input.RedundancyFactor)
		result.RedundancyFactor = &redundancy
	}

	if input.Workers != nil {
		workers := int(*input.Workers)
		result.Workers = &workers
	}

	return result
}

func flattenBrokerFrontend(input *broker.BrokerFrontend) *BrokerFrontendModel {
	if input == nil {
		return nil
	}

	result := &BrokerFrontendModel{}

	if input.Replicas != nil {
		replicas := int(*input.Replicas)
		result.Replicas = &replicas
	}

	if input.Workers != nil {
		workers := int(*input.Workers)
		result.Workers = &workers
	}

	return result
}
