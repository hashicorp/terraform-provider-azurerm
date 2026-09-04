// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/monitor/2026-04-01/pipelinegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity

type MonitorPipelineResource struct{}

var (
	_ sdk.ResourceWithUpdate        = MonitorPipelineResource{}
	_ sdk.ResourceWithIdentity      = MonitorPipelineResource{}
	_ sdk.ResourceWithCustomizeDiff = MonitorPipelineResource{}
)

func (r MonitorPipelineResource) ResourceType() string {
	return "azurerm_monitor_pipeline"
}

func (r MonitorPipelineResource) Identity() resourceids.ResourceId {
	return &pipelinegroups.PipelineGroupId{}
}

func (r MonitorPipelineResource) ModelObject() interface{} {
	return &MonitorPipelineResourceModel{}
}

func (r MonitorPipelineResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return pipelinegroups.ValidatePipelineGroupID
}

type MonitorPipelineResourceModel struct {
	Name                         string                                  `tfschema:"name"`
	ResourceGroupName            string                                  `tfschema:"resource_group_name"`
	Location                     string                                  `tfschema:"location"`
	CustomLocationId             string                                  `tfschema:"custom_location_id"`
	Service                      []PipelineGroupServiceModel             `tfschema:"service"`
	ExecutionPlacementConstraint []PipelineGroupPlacementConstraintModel `tfschema:"execution_placement_constraint"`
	Exporter                     []PipelineGroupExporterModel            `tfschema:"exporter"`
	Processor                    []PipelineGroupProcessorModel           `tfschema:"processor"`
	Receiver                     []PipelineGroupReceiverModel            `tfschema:"receiver"`
	Replicas                     int64                                   `tfschema:"replicas"`
	TlsConfiguration             []PipelineGroupTlsConfigurationModel    `tfschema:"tls_configuration"`
	Tags                         map[string]string                       `tfschema:"tags"`
}

type PipelineGroupPlacementConstraintModel struct {
	Capability string   `tfschema:"capability"`
	Operator   string   `tfschema:"operator"`
	Values     []string `tfschema:"values"`
}

type PipelineGroupExporterModel struct {
	Name                      string                                                `tfschema:"name"`
	AzureMonitorWorkspaceLogs []PipelineGroupAzureMonitorWorkspaceLogsExporterModel `tfschema:"azure_monitor_workspace_logs"`
}

type PipelineGroupAzureMonitorWorkspaceLogsExporterModel struct {
	Api         []PipelineGroupAzureMonitorWorkspaceLogsApiModel `tfschema:"api"`
	Persistence []PipelineGroupExporterPersistenceModel          `tfschema:"persistence"`
}

type PipelineGroupAzureMonitorWorkspaceLogsApiModel struct {
	DataCollectionEndpointUrl     string                        `tfschema:"data_collection_endpoint_url"`
	DataCollectionRuleImmutableId string                        `tfschema:"data_collection_rule_immutable_id"`
	Schema                        []PipelineGroupSchemaMapModel `tfschema:"schema"`
	Stream                        string                        `tfschema:"stream"`
}

type PipelineGroupSchemaMapModel struct {
	RecordMap   []PipelineGroupFieldMapModel `tfschema:"record_map"`
	ResourceMap []PipelineGroupFieldMapModel `tfschema:"resource_map"`
	ScopeMap    []PipelineGroupFieldMapModel `tfschema:"scope_map"`
}

type PipelineGroupFieldMapModel struct {
	From string `tfschema:"from"`
	To   string `tfschema:"to"`
}

type PipelineGroupExporterPersistenceModel struct {
	MaximumStorageUsageInGB  int64 `tfschema:"maximum_storage_usage_in_gb"`
	RetentionPeriodInMinutes int64 `tfschema:"retention_period_in_minutes"`
}

type PipelineGroupProcessorModel struct {
	Name               string                             `tfschema:"name"`
	Type               string                             `tfschema:"type"`
	Batch              []PipelineGroupBatchProcessorModel `tfschema:"batch"`
	TransformStatement string                             `tfschema:"transform_statement"`
}

type PipelineGroupBatchProcessorModel struct {
	BatchSize             int64 `tfschema:"batch_size"`
	TimeoutInMilliseconds int64 `tfschema:"timeout_in_milliseconds"`
}

type PipelineGroupReceiverModel struct {
	Name                 string                             `tfschema:"name"`
	Type                 string                             `tfschema:"type"`
	Otlp                 []PipelineGroupOtlpReceiverModel   `tfschema:"otlp"`
	Syslog               []PipelineGroupSyslogReceiverModel `tfschema:"syslog"`
	TlsConfigurationName string                             `tfschema:"tls_configuration_name"`
}

type PipelineGroupOtlpReceiverModel struct {
	Endpoint string `tfschema:"endpoint"`
}

type PipelineGroupSyslogReceiverModel struct {
	AllowSkipPriorityHeader bool     `tfschema:"allow_skip_priority_header"`
	AllowedFormats          []string `tfschema:"allowed_formats"`
	Endpoint                string   `tfschema:"endpoint"`
	TransportProtocol       string   `tfschema:"transport_protocol"`
}

type PipelineGroupServiceModel struct {
	PersistentVolumeName string                       `tfschema:"persistent_volume_name"`
	Pipeline             []PipelineGroupPipelineModel `tfschema:"pipeline"`
}

type PipelineGroupPipelineModel struct {
	Name       string   `tfschema:"name"`
	Exporters  []string `tfschema:"exporters"`
	Receivers  []string `tfschema:"receivers"`
	Processors []string `tfschema:"processors"`
}

type PipelineGroupTlsConfigurationModel struct {
	Name                       string                                `tfschema:"name"`
	ClientCertificateAuthority []PipelineGroupCertificateSourceModel `tfschema:"client_certificate_authority"`
	Mode                       string                                `tfschema:"mode"`
	TlsCertificate             []PipelineGroupTlsCertificateModel    `tfschema:"tls_certificate"`
}

type PipelineGroupCertificateSourceModel struct {
	Location    string `tfschema:"location"`
	SubLocation string `tfschema:"sub_location"`
	Type        string `tfschema:"type"`
}

type PipelineGroupTlsCertificateModel struct {
	Certificate []PipelineGroupCertificateSourceModel `tfschema:"certificate"`
	PrivateKey  []PipelineGroupPrivateKeySourceModel  `tfschema:"private_key"`
}

type PipelineGroupPrivateKeySourceModel struct {
	Location    string `tfschema:"location"`
	SubLocation string `tfschema:"sub_location"`
}

func (r MonitorPipelineResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{2,31}[a-zA-Z0-9]$`),
				"must be between 4 and 33 characters, contain only letters, numbers and hyphens, and must not start or end with a hyphen",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"custom_location_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: customlocations.ValidateCustomLocationID,
		},

		"service": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"pipeline": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{2,31}[a-zA-Z0-9]$`),
										"must be between 4 and 33 characters, contain only letters, numbers and hyphens, and must not start or end with a hyphen",
									),
								},

								"exporters": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringMatch(
											regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{2,31}[a-zA-Z0-9]$`),
											"must be between 4 and 33 characters, contain only letters, numbers and hyphens, and must not start or end with a hyphen",
										),
									},
								},

								"receivers": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringMatch(
											regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{2,31}[a-zA-Z0-9]$`),
											"must be between 4 and 33 characters, contain only letters, numbers and hyphens, and must not start or end with a hyphen",
										),
									},
								},

								"processors": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringMatch(
											regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{2,31}[a-zA-Z0-9]$`),
											"must be between 4 and 33 characters, contain only letters, numbers and hyphens, and must not start or end with a hyphen",
										),
									},
								},
							},
						},
					},

					"persistent_volume_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"execution_placement_constraint": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"capability": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"operator": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(pipelinegroups.PossibleValuesForCapabilityOperator(), false),
					},

					"values": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},

		"exporter": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{2,31}[a-zA-Z0-9]$`),
							"must be between 4 and 33 characters, contain only letters, numbers and hyphens, and must not start or end with a hyphen",
						),
					},

					"azure_monitor_workspace_logs": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"api": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"data_collection_endpoint_url": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.IsURLWithHTTPorHTTPS,
											},

											"data_collection_rule_immutable_id": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringMatch(
													regexp.MustCompile(`^dcr-[0-9a-fA-F]{32}$`),
													"must be a Data Collection Rule immutable ID in the form `dcr-<guid>`",
												),
											},

											"schema": {
												Type:     pluginsdk.TypeList,
												Required: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"record_map": {
															Type:     pluginsdk.TypeList,
															Required: true,
															MinItems: 1,
															Elem:     pipelineGroupFieldMapSchema(),
														},

														"resource_map": {
															Type:     pluginsdk.TypeList,
															Optional: true,
															Elem:     pipelineGroupFieldMapSchema(),
														},

														"scope_map": {
															Type:     pluginsdk.TypeList,
															Optional: true,
															Elem:     pipelineGroupFieldMapSchema(),
														},
													},
												},
											},

											"stream": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"persistence": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"maximum_storage_usage_in_gb": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntAtLeast(1),
											},

											// The service caps retention at 2 days (2880 minutes).
											"retention_period_in_minutes": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(1, 2880),
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

		"processor": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{2,31}[a-zA-Z0-9]$`),
							"must be between 4 and 33 characters, contain only letters, numbers and hyphens, and must not start or end with a hyphen",
						),
					},

					"type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(pipelinegroups.PossibleValuesForProcessorType(), false),
					},

					"batch": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"batch_size": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntBetween(10, 100000),
								},

								"timeout_in_milliseconds": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntBetween(10, 300000),
								},
							},
						},
					},

					"transform_statement": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringLenBetween(1, 10000),
					},
				},
			},
		},

		"receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{2,31}[a-zA-Z0-9]$`),
							"must be between 4 and 33 characters, contain only letters, numbers and hyphens, and must not start or end with a hyphen",
						),
					},

					"type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(pipelinegroups.PossibleValuesForReceiverType(), false),
					},

					"otlp": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"endpoint": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile(`^(?:[a-zA-Z][a-zA-Z0-9+.-]*://)?(?:\[[0-9a-fA-F:.]+\]|[^:/?#[:space:]]*):(?:[1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`),
										"must be in the form `<host>:<port>` with a numeric port between `1` and `65535`",
									),
								},
							},
						},
					},

					"syslog": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"endpoint": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile(`^(?:[a-zA-Z][a-zA-Z0-9+.-]*://)?(?:\[[0-9a-fA-F:.]+\]|[^:/?#[:space:]]*):(?:[1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`),
										"must be in the form `<host>:<port>` with a numeric port between `1` and `65535`",
									),
								},

								"allow_skip_priority_header": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"allowed_formats": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									// NOTE: O+C - Azure defaults this to ["all"] when omitted.
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(pipelinegroups.PossibleValuesForAllowedFormats(), false),
									},
								},

								"transport_protocol": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      string(pipelinegroups.TransportProtocolTcp),
									ValidateFunc: validation.StringInSlice(pipelinegroups.PossibleValuesForTransportProtocol(), false),
								},
							},
						},
					},

					"tls_configuration_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{2,31}[a-zA-Z0-9]$`),
							"must be between 4 and 33 characters, contain only letters, numbers and hyphens, and must not start or end with a hyphen",
						),
					},
				},
			},
		},

		"replicas": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"tls_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{2,31}[a-zA-Z0-9]$`),
							"must be between 4 and 33 characters, contain only letters, numbers and hyphens, and must not start or end with a hyphen",
						),
					},

					"client_certificate_authority": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem:     pipelineGroupCertificateSourceSchema(),
					},

					"mode": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(pipelinegroups.TlsModeMutualTls),
						ValidateFunc: validation.StringInSlice(pipelinegroups.PossibleValuesForTlsMode(), false),
					},

					"tls_certificate": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"certificate": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
									Elem:     pipelineGroupCertificateSourceSchema(),
								},

								"private_key": {
									Type:      pluginsdk.TypeList,
									Required:  true,
									MaxItems:  1,
									Sensitive: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"location": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"sub_location": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
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

		"tags": commonschema.Tags(),
	}
}

func (r MonitorPipelineResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r MonitorPipelineResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model MonitorPipelineResourceModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if err := validatePipelineGroupPlacementConstraints(model.ExecutionPlacementConstraint); err != nil {
				return err
			}

			if err := validatePipelineGroupTlsConfigurations(model.TlsConfiguration); err != nil {
				return err
			}

			if err := validatePipelineGroupReceivers(model.Receiver, model.TlsConfiguration); err != nil {
				return err
			}

			if err := validatePipelineGroupExporters(model.Exporter, model.Service); err != nil {
				return err
			}

			if err := validatePipelineGroupProcessors(model.Processor); err != nil {
				return err
			}

			return nil
		},
	}
}

func validatePipelineGroupPlacementConstraints(input []PipelineGroupPlacementConstraintModel) error {
	for i, constraint := range input {
		switch constraint.Operator {
		case string(pipelinegroups.CapabilityOperatorIn), string(pipelinegroups.CapabilityOperatorNotIn):
			if len(constraint.Values) == 0 {
				return fmt.Errorf("`values` must be set and non-empty for `execution_placement_constraint.%d` when `operator` is `%s`", i, constraint.Operator)
			}
		case string(pipelinegroups.CapabilityOperatorExists), string(pipelinegroups.CapabilityOperatorDoesNotExist):
			if len(constraint.Values) > 0 {
				return fmt.Errorf("`values` must not be set for `execution_placement_constraint.%d` when `operator` is `%s`", i, constraint.Operator)
			}
		}
	}

	return nil
}

func validatePipelineGroupTlsConfigurations(input []PipelineGroupTlsConfigurationModel) error {
	for i, tlsConfiguration := range input {
		switch tlsConfiguration.Mode {
		case string(pipelinegroups.TlsModeDisabled):
			if len(tlsConfiguration.TlsCertificate) > 0 || len(tlsConfiguration.ClientCertificateAuthority) > 0 {
				return fmt.Errorf("`tls_certificate` and `client_certificate_authority` must not be set for `tls_configuration.%d` when `mode` is `%s`", i, tlsConfiguration.Mode)
			}
		case string(pipelinegroups.TlsModeServerOnly):
			if len(tlsConfiguration.ClientCertificateAuthority) > 0 {
				return fmt.Errorf("`client_certificate_authority` must not be set for `tls_configuration.%d` when `mode` is `%s`", i, tlsConfiguration.Mode)
			}
		}
	}

	return nil
}

func validatePipelineGroupReceivers(input []PipelineGroupReceiverModel, tlsConfigurations []PipelineGroupTlsConfigurationModel) error {
	tlsConfigurationNames := make(map[string]struct{}, len(tlsConfigurations))
	for _, tlsConfiguration := range tlsConfigurations {
		tlsConfigurationNames[tlsConfiguration.Name] = struct{}{}
	}

	for i, receiver := range input {
		if receiver.TlsConfigurationName != "" {
			if _, ok := tlsConfigurationNames[receiver.TlsConfigurationName]; !ok {
				return fmt.Errorf("`receiver.%d.tls_configuration_name` references unknown `tls_configuration` `%s`", i, receiver.TlsConfigurationName)
			}
		}

		switch receiver.Type {
		case string(pipelinegroups.ReceiverTypeOTLP):
			if len(receiver.Otlp) == 0 {
				return fmt.Errorf("`receiver.%d.otlp` must be set when `type` is `%s`", i, pipelinegroups.ReceiverTypeOTLP)
			}
			if len(receiver.Syslog) > 0 {
				return fmt.Errorf("`receiver.%d.syslog` must not be set when `type` is `%s`", i, pipelinegroups.ReceiverTypeOTLP)
			}
		case string(pipelinegroups.ReceiverTypeSyslog):
			if len(receiver.Syslog) == 0 {
				return fmt.Errorf("`receiver.%d.syslog` must be set when `type` is `%s`", i, pipelinegroups.ReceiverTypeSyslog)
			}
			if len(receiver.Otlp) > 0 {
				return fmt.Errorf("`receiver.%d.otlp` must not be set when `type` is `%s`", i, pipelinegroups.ReceiverTypeSyslog)
			}
		}

		if len(receiver.Syslog) == 0 {
			continue
		}

		syslog := receiver.Syslog[0]

		if receiver.TlsConfigurationName != "" && syslog.TransportProtocol == string(pipelinegroups.TransportProtocolUdp) {
			return fmt.Errorf("`tls_configuration_name` is not supported for `receiver.%d` when `syslog.0.transport_protocol` is `%s`", i, pipelinegroups.TransportProtocolUdp)
		}

		if err := validatePipelineGroupSyslogAllowedFormats(syslog, i); err != nil {
			return err
		}
	}

	return nil
}

func validatePipelineGroupSyslogAllowedFormats(syslog PipelineGroupSyslogReceiverModel, index int) error {
	if len(syslog.AllowedFormats) == 0 {
		return nil
	}

	hasAll := false
	hasSkipPriorityHeaderCompatibleFormat := false
	for _, format := range syslog.AllowedFormats {
		if format == string(pipelinegroups.AllowedFormatsAll) {
			hasAll = true
			hasSkipPriorityHeaderCompatibleFormat = true
		}
		if format == string(pipelinegroups.AllowedFormatsSyslogRfcThreeOneSixFour) || format == string(pipelinegroups.AllowedFormatsCefRfcThreeOneSixFour) {
			hasSkipPriorityHeaderCompatibleFormat = true
		}
	}

	if hasAll && len(syslog.AllowedFormats) > 1 {
		return fmt.Errorf("`allowed_formats` for `receiver.%d.syslog` must not combine `%s` with other formats", index, pipelinegroups.AllowedFormatsAll)
	}

	if syslog.AllowSkipPriorityHeader && !hasSkipPriorityHeaderCompatibleFormat {
		return fmt.Errorf("`allow_skip_priority_header` for `receiver.%d.syslog` requires `allowed_formats` to include `%s`, `%s`, or `%s`", index, pipelinegroups.AllowedFormatsAll, pipelinegroups.AllowedFormatsSyslogRfcThreeOneSixFour, pipelinegroups.AllowedFormatsCefRfcThreeOneSixFour)
	}

	return nil
}

func validatePipelineGroupProcessors(input []PipelineGroupProcessorModel) error {
	for i, processor := range input {
		if len(processor.Batch) == 0 {
			continue
		}

		batch := processor.Batch[0]
		if batch.BatchSize == 0 && batch.TimeoutInMilliseconds == 0 {
			return fmt.Errorf("`processor.%d.batch` must have at least one of `batch_size` or `timeout_in_milliseconds` set", i)
		}
	}

	return nil
}

func validatePipelineGroupExporters(input []PipelineGroupExporterModel, service []PipelineGroupServiceModel) error {
	requiresPersistentVolume := false

	for i, exporter := range input {
		if len(exporter.AzureMonitorWorkspaceLogs) == 0 {
			continue
		}

		azureMonitorWorkspaceLogs := exporter.AzureMonitorWorkspaceLogs[0]

		if len(azureMonitorWorkspaceLogs.Persistence) > 0 {
			requiresPersistentVolume = true

			persistence := azureMonitorWorkspaceLogs.Persistence[0]
			if persistence.MaximumStorageUsageInGB == 0 && persistence.RetentionPeriodInMinutes == 0 {
				return fmt.Errorf("`exporter.%d.azure_monitor_workspace_logs.0.persistence` must have at least one of `maximum_storage_usage_in_gb` or `retention_period_in_minutes` set", i)
			}
		}

		if len(azureMonitorWorkspaceLogs.Api) == 0 || len(azureMonitorWorkspaceLogs.Api[0].Schema) == 0 {
			continue
		}

		hasTimeGenerated := false
		for _, entry := range azureMonitorWorkspaceLogs.Api[0].Schema[0].RecordMap {
			if entry.To == "TimeGenerated" {
				hasTimeGenerated = true
				break
			}
		}

		if !hasTimeGenerated {
			return fmt.Errorf("`exporter.%d` must define a `record_map` entry with `to` set to `TimeGenerated`", i)
		}
	}

	if requiresPersistentVolume && (len(service) == 0 || service[0].PersistentVolumeName == "") {
		return errors.New("`service.0.persistent_volume_name` must be set when an exporter has `persistence` configured")
	}

	return nil
}

// shared by the recordMap, resourceMap and scopeMap blocks.
func pipelineGroupFieldMapSchema() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"from": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"to": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

// shared by the clientCertificateAuthority and tlsCertificate.certificate blocks.
func pipelineGroupCertificateSourceSchema() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"location": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"sub_location": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"type": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(pipelinegroups.PossibleValuesForCertificateSourceType(), false),
			},
		},
	}
}

func (r MonitorPipelineResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.PipelineGroupsClient

			var model MonitorPipelineResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := pipelinegroups.NewPipelineGroupID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := pipelinegroups.PipelineGroup{
				ExtendedLocation: &pipelinegroups.AzureResourceManagerCommonTypesExtendedLocation{
					Name: model.CustomLocationId,
					Type: pipelinegroups.ExtendedLocationTypeCustomLocation,
				},
				Location: location.Normalize(model.Location),
				Properties: &pipelinegroups.PipelineGroupProperties{
					ExecutionPlacement: expandPipelineGroupExecutionPlacement(model.ExecutionPlacementConstraint),
					Exporters:          expandPipelineGroupExporters(model.Exporter),
					Processors:         expandPipelineGroupProcessors(model.Processor),
					Receivers:          expandPipelineGroupReceivers(model.Receiver),
					Replicas:           pointer.ToOrNil(model.Replicas),
					Service:            expandPipelineGroupService(model.Service),
				},
				Tags: pointer.To(model.Tags),
			}

			if len(model.TlsConfiguration) > 0 {
				payload.Properties.TlsConfigurations = expandPipelineGroupTlsConfigurations(model.TlsConfiguration)
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.ResourceData.SetId(id.ID())
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r MonitorPipelineResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.PipelineGroupsClient

			id, err := pipelinegroups.ParsePipelineGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model MonitorPipelineResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp.Model == nil || resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			payload := *resp.Model

			if metadata.ResourceData.HasChange("execution_placement_constraint") {
				payload.Properties.ExecutionPlacement = expandPipelineGroupExecutionPlacement(model.ExecutionPlacementConstraint)
			}

			if metadata.ResourceData.HasChange("exporter") {
				payload.Properties.Exporters = expandPipelineGroupExporters(model.Exporter)
			}

			if metadata.ResourceData.HasChange("processor") {
				payload.Properties.Processors = expandPipelineGroupProcessors(model.Processor)
			}

			if metadata.ResourceData.HasChange("receiver") {
				payload.Properties.Receivers = expandPipelineGroupReceivers(model.Receiver)
			}

			if metadata.ResourceData.HasChange("replicas") {
				payload.Properties.Replicas = pointer.ToOrNil(model.Replicas)
			}

			if metadata.ResourceData.HasChange("service") {
				payload.Properties.Service = expandPipelineGroupService(model.Service)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if metadata.ResourceData.HasChange("tls_configuration") {
				payload.Properties.TlsConfigurations = expandPipelineGroupTlsConfigurations(model.TlsConfiguration)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r MonitorPipelineResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.PipelineGroupsClient

			id, err := pipelinegroups.ParsePipelineGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r MonitorPipelineResource) flatten(metadata sdk.ResourceMetaData, id *pipelinegroups.PipelineGroupId, model *pipelinegroups.PipelineGroup) error {
	state := MonitorPipelineResourceModel{
		Name:              id.PipelineGroupName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		state.Tags = pointer.From(model.Tags)

		if model.ExtendedLocation != nil && model.ExtendedLocation.Name != "" {
			customLocationId, err := customlocations.ParseCustomLocationIDInsensitively(model.ExtendedLocation.Name)
			if err != nil {
				return err
			}
			state.CustomLocationId = customLocationId.ID()
		}

		if props := model.Properties; props != nil {
			state.ExecutionPlacementConstraint = flattenPipelineGroupExecutionPlacement(props.ExecutionPlacement)
			state.Exporter = flattenPipelineGroupExporters(props.Exporters)
			state.Processor = flattenPipelineGroupProcessors(props.Processors)
			state.Receiver = flattenPipelineGroupReceivers(props.Receivers)
			state.Replicas = pointer.From(props.Replicas)
			state.Service = flattenPipelineGroupService(props.Service)
			state.TlsConfiguration = flattenPipelineGroupTlsConfigurations(props.TlsConfigurations)
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (r MonitorPipelineResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.PipelineGroupsClient

			id, err := pipelinegroups.ParsePipelineGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandPipelineGroupExecutionPlacement(input []PipelineGroupPlacementConstraintModel) *pipelinegroups.ExecutionPlacement {
	if len(input) == 0 {
		return nil
	}

	constraints := make([]pipelinegroups.PlacementConstraint, 0, len(input))
	for _, constraint := range input {
		constraints = append(constraints, pipelinegroups.PlacementConstraint{
			Capability: constraint.Capability,
			Operator:   pipelinegroups.CapabilityOperator(constraint.Operator),
			Values:     pointer.To(constraint.Values),
		})
	}

	return &pipelinegroups.ExecutionPlacement{
		Constraints: &constraints,
	}
}

func flattenPipelineGroupExecutionPlacement(input *pipelinegroups.ExecutionPlacement) []PipelineGroupPlacementConstraintModel {
	if input == nil || input.Constraints == nil {
		return []PipelineGroupPlacementConstraintModel{}
	}

	output := make([]PipelineGroupPlacementConstraintModel, 0, len(*input.Constraints))
	for _, constraint := range *input.Constraints {
		output = append(output, PipelineGroupPlacementConstraintModel{
			Capability: constraint.Capability,
			Operator:   string(constraint.Operator),
			Values:     pointer.From(constraint.Values),
		})
	}

	return output
}

func expandPipelineGroupExporters(input []PipelineGroupExporterModel) []pipelinegroups.Exporter {
	output := make([]pipelinegroups.Exporter, 0, len(input))

	for _, v := range input {
		exporter := pipelinegroups.Exporter{
			Name: v.Name,
			Type: pipelinegroups.ExporterTypeAzureMonitorWorkspaceLogs,
		}

		if len(v.AzureMonitorWorkspaceLogs) > 0 {
			exporter.AzureMonitorWorkspaceLogs = expandPipelineGroupAzureMonitorWorkspaceLogsExporter(v.AzureMonitorWorkspaceLogs[0])
		}

		output = append(output, exporter)
	}

	return output
}

func expandPipelineGroupAzureMonitorWorkspaceLogsExporter(input PipelineGroupAzureMonitorWorkspaceLogsExporterModel) *pipelinegroups.AzureMonitorWorkspaceLogsExporter {
	output := pipelinegroups.AzureMonitorWorkspaceLogsExporter{}

	if len(input.Api) > 0 {
		api := input.Api[0]
		output.Api = pipelinegroups.AzureMonitorWorkspaceLogsApiConfig{
			DataCollectionEndpointURL: api.DataCollectionEndpointUrl,
			DataCollectionRule:        api.DataCollectionRuleImmutableId,
			Stream:                    api.Stream,
		}

		if len(api.Schema) > 0 {
			output.Api.Schema = expandPipelineGroupSchemaMap(api.Schema[0])
		}
	}

	if len(input.Persistence) > 0 {
		persistence := input.Persistence[0]
		output.Persistence = &pipelinegroups.ExporterPersistenceConfiguration{
			MaxStorageUsage: pointer.ToOrNil(persistence.MaximumStorageUsageInGB),
			RetentionPeriod: pointer.ToOrNil(persistence.RetentionPeriodInMinutes),
		}
	}

	return &output
}

func expandPipelineGroupSchemaMap(input PipelineGroupSchemaMapModel) pipelinegroups.SchemaMap {
	output := pipelinegroups.SchemaMap{
		RecordMap: expandPipelineGroupFieldMaps(input.RecordMap),
	}

	if len(input.ResourceMap) > 0 {
		resourceMaps := make([]pipelinegroups.ResourceMap, 0, len(input.ResourceMap))
		for _, v := range input.ResourceMap {
			resourceMaps = append(resourceMaps, pipelinegroups.ResourceMap{From: v.From, To: v.To})
		}
		output.ResourceMap = &resourceMaps
	}

	if len(input.ScopeMap) > 0 {
		scopeMaps := make([]pipelinegroups.ScopeMap, 0, len(input.ScopeMap))
		for _, v := range input.ScopeMap {
			scopeMaps = append(scopeMaps, pipelinegroups.ScopeMap{From: v.From, To: v.To})
		}
		output.ScopeMap = &scopeMaps
	}

	return output
}

func expandPipelineGroupFieldMaps(input []PipelineGroupFieldMapModel) []pipelinegroups.RecordMap {
	output := make([]pipelinegroups.RecordMap, 0, len(input))
	for _, v := range input {
		output = append(output, pipelinegroups.RecordMap{From: v.From, To: v.To})
	}
	return output
}

func flattenPipelineGroupExporters(input []pipelinegroups.Exporter) []PipelineGroupExporterModel {
	output := make([]PipelineGroupExporterModel, 0, len(input))

	for _, v := range input {
		exporter := PipelineGroupExporterModel{
			Name: v.Name,
		}

		if v.AzureMonitorWorkspaceLogs != nil {
			exporter.AzureMonitorWorkspaceLogs = []PipelineGroupAzureMonitorWorkspaceLogsExporterModel{
				flattenPipelineGroupAzureMonitorWorkspaceLogsExporter(*v.AzureMonitorWorkspaceLogs),
			}
		}

		output = append(output, exporter)
	}

	return output
}

func flattenPipelineGroupAzureMonitorWorkspaceLogsExporter(input pipelinegroups.AzureMonitorWorkspaceLogsExporter) PipelineGroupAzureMonitorWorkspaceLogsExporterModel {
	output := PipelineGroupAzureMonitorWorkspaceLogsExporterModel{
		Api: []PipelineGroupAzureMonitorWorkspaceLogsApiModel{
			{
				DataCollectionEndpointUrl:     input.Api.DataCollectionEndpointURL,
				DataCollectionRuleImmutableId: input.Api.DataCollectionRule,
				Stream:                        input.Api.Stream,
				Schema:                        flattenPipelineGroupSchemaMap(input.Api.Schema),
			},
		},
	}

	if input.Persistence != nil {
		output.Persistence = []PipelineGroupExporterPersistenceModel{
			{
				MaximumStorageUsageInGB:  pointer.From(input.Persistence.MaxStorageUsage),
				RetentionPeriodInMinutes: pointer.From(input.Persistence.RetentionPeriod),
			},
		}
	}

	return output
}

func flattenPipelineGroupSchemaMap(input pipelinegroups.SchemaMap) []PipelineGroupSchemaMapModel {
	output := PipelineGroupSchemaMapModel{
		RecordMap: flattenPipelineGroupFieldMaps(input.RecordMap),
	}

	if input.ResourceMap != nil {
		resourceMaps := make([]PipelineGroupFieldMapModel, 0, len(*input.ResourceMap))
		for _, v := range *input.ResourceMap {
			resourceMaps = append(resourceMaps, PipelineGroupFieldMapModel{From: v.From, To: v.To})
		}
		output.ResourceMap = resourceMaps
	}

	if input.ScopeMap != nil {
		scopeMaps := make([]PipelineGroupFieldMapModel, 0, len(*input.ScopeMap))
		for _, v := range *input.ScopeMap {
			scopeMaps = append(scopeMaps, PipelineGroupFieldMapModel{From: v.From, To: v.To})
		}
		output.ScopeMap = scopeMaps
	}

	return []PipelineGroupSchemaMapModel{output}
}

func flattenPipelineGroupFieldMaps(input []pipelinegroups.RecordMap) []PipelineGroupFieldMapModel {
	output := make([]PipelineGroupFieldMapModel, 0, len(input))
	for _, v := range input {
		output = append(output, PipelineGroupFieldMapModel{From: v.From, To: v.To})
	}
	return output
}

func expandPipelineGroupProcessors(input []PipelineGroupProcessorModel) []pipelinegroups.Processor {
	output := make([]pipelinegroups.Processor, 0, len(input))

	for _, v := range input {
		processor := pipelinegroups.Processor{
			Name: v.Name,
			Type: pipelinegroups.ProcessorType(v.Type),
		}

		if len(v.Batch) > 0 {
			batch := v.Batch[0]
			processor.Batch = &pipelinegroups.BatchProcessor{
				BatchSize: pointer.ToOrNil(batch.BatchSize),
				Timeout:   pointer.ToOrNil(batch.TimeoutInMilliseconds),
			}
		}

		if v.TransformStatement != "" {
			processor.TransformLanguage = &pipelinegroups.TransformLanguageProcessor{
				TransformStatement: v.TransformStatement,
			}
		}

		output = append(output, processor)
	}

	return output
}

func flattenPipelineGroupProcessors(input []pipelinegroups.Processor) []PipelineGroupProcessorModel {
	output := make([]PipelineGroupProcessorModel, 0, len(input))

	for _, v := range input {
		processor := PipelineGroupProcessorModel{
			Name: v.Name,
			Type: string(v.Type),
		}

		if v.Batch != nil {
			processor.Batch = []PipelineGroupBatchProcessorModel{
				{
					BatchSize:             pointer.From(v.Batch.BatchSize),
					TimeoutInMilliseconds: pointer.From(v.Batch.Timeout),
				},
			}
		}

		if v.TransformLanguage != nil {
			processor.TransformStatement = v.TransformLanguage.TransformStatement
		}

		output = append(output, processor)
	}

	return output
}

func expandPipelineGroupReceivers(input []PipelineGroupReceiverModel) []pipelinegroups.Receiver {
	output := make([]pipelinegroups.Receiver, 0, len(input))

	for _, v := range input {
		receiver := pipelinegroups.Receiver{
			Name: v.Name,
			Type: pipelinegroups.ReceiverType(v.Type),
		}

		receiver.TlsConfiguration = pointer.ToOrNil(v.TlsConfigurationName)

		switch {
		case len(v.Otlp) > 0:
			receiver.Otlp = &pipelinegroups.OtlpReceiver{
				Endpoint: v.Otlp[0].Endpoint,
			}
		case len(v.Syslog) > 0:
			syslog := v.Syslog[0]
			receiver.Syslog = &pipelinegroups.SyslogReceiver{
				Endpoint:           syslog.Endpoint,
				AllowSkipPriHeader: &syslog.AllowSkipPriorityHeader,
				TransportProtocol:  pointer.ToEnum[pipelinegroups.TransportProtocol](syslog.TransportProtocol),
			}

			if len(syslog.AllowedFormats) > 0 {
				receiver.Syslog.AllowedFormats = pointer.ToEnumSlice[pipelinegroups.AllowedFormats](syslog.AllowedFormats)
			}
		}

		output = append(output, receiver)
	}

	return output
}

func flattenPipelineGroupReceivers(input []pipelinegroups.Receiver) []PipelineGroupReceiverModel {
	output := make([]PipelineGroupReceiverModel, 0, len(input))

	for _, v := range input {
		receiver := PipelineGroupReceiverModel{
			Name:                 v.Name,
			Type:                 string(v.Type),
			TlsConfigurationName: pointer.From(v.TlsConfiguration),
		}

		if v.Otlp != nil {
			receiver.Otlp = []PipelineGroupOtlpReceiverModel{
				{Endpoint: v.Otlp.Endpoint},
			}
		}

		if v.Syslog != nil {
			receiver.Syslog = []PipelineGroupSyslogReceiverModel{
				{
					Endpoint:                v.Syslog.Endpoint,
					AllowSkipPriorityHeader: pointer.From(v.Syslog.AllowSkipPriHeader),
					AllowedFormats:          pointer.FromEnumSlice(v.Syslog.AllowedFormats),
					TransportProtocol:       pointer.FromEnum(v.Syslog.TransportProtocol),
				},
			}
		}

		output = append(output, receiver)
	}

	return output
}

func expandPipelineGroupService(input []PipelineGroupServiceModel) pipelinegroups.Service {
	output := pipelinegroups.Service{
		Pipelines: []pipelinegroups.Pipeline{},
	}

	if len(input) == 0 {
		return output
	}

	v := input[0]

	if v.PersistentVolumeName != "" {
		output.Persistence = &pipelinegroups.PersistenceConfigurations{
			PersistentVolumeName: v.PersistentVolumeName,
		}
	}

	pipelines := make([]pipelinegroups.Pipeline, 0, len(v.Pipeline))
	for _, p := range v.Pipeline {
		pipeline := pipelinegroups.Pipeline{
			Name:      p.Name,
			Type:      pipelinegroups.PipelineTypeLogs,
			Exporters: p.Exporters,
			Receivers: p.Receivers,
		}

		if len(p.Processors) > 0 {
			pipeline.Processors = &p.Processors
		}

		pipelines = append(pipelines, pipeline)
	}
	output.Pipelines = pipelines

	return output
}

func flattenPipelineGroupService(input pipelinegroups.Service) []PipelineGroupServiceModel {
	output := PipelineGroupServiceModel{
		PersistentVolumeName: pointer.From(input.Persistence).PersistentVolumeName,
	}

	pipelines := make([]PipelineGroupPipelineModel, 0, len(input.Pipelines))
	for _, p := range input.Pipelines {
		pipelines = append(pipelines, PipelineGroupPipelineModel{
			Name:       p.Name,
			Exporters:  p.Exporters,
			Receivers:  p.Receivers,
			Processors: pointer.From(p.Processors),
		})
	}
	output.Pipeline = pipelines

	return []PipelineGroupServiceModel{output}
}

func expandPipelineGroupTlsConfigurations(input []PipelineGroupTlsConfigurationModel) *[]pipelinegroups.TlsConfiguration {
	if len(input) == 0 {
		return pointer.To([]pipelinegroups.TlsConfiguration{})
	}

	output := make([]pipelinegroups.TlsConfiguration, 0, len(input))

	for _, v := range input {
		tlsConfiguration := pipelinegroups.TlsConfiguration{
			Name: v.Name,
			Mode: pointer.ToEnum[pipelinegroups.TlsMode](v.Mode),
		}

		if len(v.ClientCertificateAuthority) > 0 {
			tlsConfiguration.ClientCa = expandPipelineGroupCertificateSource(v.ClientCertificateAuthority[0])
		}

		if len(v.TlsCertificate) > 0 {
			tlsCertificate := v.TlsCertificate[0]
			certificateWithKey := pipelinegroups.CertificateWithKey{}

			if len(tlsCertificate.Certificate) > 0 {
				certificateWithKey.Certificate = *expandPipelineGroupCertificateSource(tlsCertificate.Certificate[0])
			}

			if len(tlsCertificate.PrivateKey) > 0 {
				privateKey := tlsCertificate.PrivateKey[0]
				certificateWithKey.PrivateKey = pipelinegroups.PrivateKeySource{
					Location:    privateKey.Location,
					SubLocation: privateKey.SubLocation,
					Type:        pipelinegroups.PrivateKeySourceTypeKubernetesSecret,
				}
			}

			tlsConfiguration.TlsCertificate = &certificateWithKey
		}

		output = append(output, tlsConfiguration)
	}

	return pointer.To(output)
}

func expandPipelineGroupCertificateSource(input PipelineGroupCertificateSourceModel) *pipelinegroups.CertificateSource {
	return &pipelinegroups.CertificateSource{
		Location:    input.Location,
		SubLocation: input.SubLocation,
		Type:        pipelinegroups.CertificateSourceType(input.Type),
	}
}

func flattenPipelineGroupTlsConfigurations(input *[]pipelinegroups.TlsConfiguration) []PipelineGroupTlsConfigurationModel {
	if input == nil {
		return []PipelineGroupTlsConfigurationModel{}
	}

	output := make([]PipelineGroupTlsConfigurationModel, 0, len(*input))

	for _, v := range *input {
		tlsConfiguration := PipelineGroupTlsConfigurationModel{
			Name: v.Name,
			Mode: pointer.FromEnum(v.Mode),
		}

		tlsConfiguration.ClientCertificateAuthority = flattenPipelineGroupCertificateSourceList(v.ClientCa)

		if v.TlsCertificate != nil {
			tlsConfiguration.TlsCertificate = []PipelineGroupTlsCertificateModel{
				{
					Certificate: flattenPipelineGroupCertificateSourceList(&v.TlsCertificate.Certificate),
					PrivateKey: []PipelineGroupPrivateKeySourceModel{
						{
							Location:    v.TlsCertificate.PrivateKey.Location,
							SubLocation: v.TlsCertificate.PrivateKey.SubLocation,
						},
					},
				},
			}
		}

		output = append(output, tlsConfiguration)
	}

	return output
}

func flattenPipelineGroupCertificateSourceList(input *pipelinegroups.CertificateSource) []PipelineGroupCertificateSourceModel {
	if input == nil {
		return []PipelineGroupCertificateSourceModel{}
	}

	return []PipelineGroupCertificateSourceModel{
		{
			Location:    input.Location,
			SubLocation: input.SubLocation,
			Type:        string(input.Type),
		},
	}
}
