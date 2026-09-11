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
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity

type MonitorPipelineResource struct{}

const (
	defaultPipelineBatchSize             = 8192
	defaultPipelineTimeoutInMilliseconds = 300000
)

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
	Name                                string                                             `tfschema:"name"`
	ResourceGroupName                   string                                             `tfschema:"resource_group_name"`
	Location                            string                                             `tfschema:"location"`
	CustomLocationId                    string                                             `tfschema:"custom_location_id"`
	LogPipeline                         []PipelineLogPipelineModel                         `tfschema:"log_pipeline"`
	AzureMonitorWorkspaceLogExporter    []PipelineAzureMonitorWorkspaceLogExporterModel    `tfschema:"azure_monitor_workspace_log_exporter"`
	BatchProcessor                      []PipelineBatchProcessorModel                      `tfschema:"batch_processor"`
	ExecutionPlacementConstraint        []PipelinePlacementConstraintModel                 `tfschema:"execution_placement_constraint"`
	MicrosoftCommonSecurityLogProcessor []PipelineMicrosoftCommonSecurityLogProcessorModel `tfschema:"microsoft_common_security_log_processor"`
	MicrosoftSyslogProcessor            []PipelineMicrosoftSyslogProcessorModel            `tfschema:"microsoft_syslog_processor"`
	OtlpReceiver                        []PipelineOtlpReceiverModel                        `tfschema:"otlp_receiver"`
	PersistentVolumeName                string                                             `tfschema:"persistent_volume_name"`
	Replicas                            int64                                              `tfschema:"replicas"`
	SyslogReceiver                      []PipelineSyslogReceiverModel                      `tfschema:"syslog_receiver"`
	TlsConfiguration                    []PipelineTlsConfigurationModel                    `tfschema:"tls_configuration"`
	TransformLanguageProcessor          []PipelineTransformLanguageProcessorModel          `tfschema:"transform_language_processor"`
	Tags                                map[string]string                                  `tfschema:"tags"`
}

type PipelineLogPipelineModel struct {
	Name       string   `tfschema:"name"`
	Exporters  []string `tfschema:"exporters"`
	Receivers  []string `tfschema:"receivers"`
	Processors []string `tfschema:"processors"`
}

type PipelineAzureMonitorWorkspaceLogExporterModel struct {
	Name                                string                                      `tfschema:"name"`
	Api                                 []PipelineAzureMonitorWorkspaceLogsApiModel `tfschema:"api"`
	PersistenceMaximumStorageUsageInGB  int64                                       `tfschema:"persistence_maximum_storage_usage_in_gb"`
	PersistenceRetentionPeriodInMinutes int64                                       `tfschema:"persistence_retention_period_in_minutes"`
}

type PipelineBatchProcessorModel struct {
	Name                  string `tfschema:"name"`
	BatchSize             int64  `tfschema:"batch_size"`
	TimeoutInMilliseconds int64  `tfschema:"timeout_in_milliseconds"`
}

type PipelinePlacementConstraintModel struct {
	Capability string   `tfschema:"capability"`
	Operator   string   `tfschema:"operator"`
	Values     []string `tfschema:"values"`
}

type PipelineMicrosoftCommonSecurityLogProcessorModel struct {
	Name string `tfschema:"name"`
}

type PipelineMicrosoftSyslogProcessorModel struct {
	Name string `tfschema:"name"`
}

type PipelineOtlpReceiverModel struct {
	Name                 string `tfschema:"name"`
	Endpoint             string `tfschema:"endpoint"`
	TlsConfigurationName string `tfschema:"tls_configuration_name"`
}

type PipelineSyslogReceiverModel struct {
	Name                    string   `tfschema:"name"`
	Endpoint                string   `tfschema:"endpoint"`
	AllowSkipPriorityHeader bool     `tfschema:"allow_skip_priority_header"`
	AllowedFormats          []string `tfschema:"allowed_formats"`
	TlsConfigurationName    string   `tfschema:"tls_configuration_name"`
	TransportProtocol       string   `tfschema:"transport_protocol"`
}

type PipelineTlsConfigurationModel struct {
	Name                       string                           `tfschema:"name"`
	ClientCertificateAuthority []PipelineCertificateSourceModel `tfschema:"client_certificate_authority"`
	Mode                       string                           `tfschema:"mode"`
	TlsCertificate             []PipelineTlsCertificateModel    `tfschema:"tls_certificate"`
}

type PipelineTransformLanguageProcessorModel struct {
	Name               string `tfschema:"name"`
	TransformStatement string `tfschema:"transform_statement"`
}

type PipelineAzureMonitorWorkspaceLogsApiModel struct {
	DataCollectionEndpointUrl     string                   `tfschema:"data_collection_endpoint_url"`
	DataCollectionRuleImmutableId string                   `tfschema:"data_collection_rule_immutable_id"`
	Schema                        []PipelineSchemaMapModel `tfschema:"schema"`
	Stream                        string                   `tfschema:"stream"`
}

type PipelineCertificateSourceModel struct {
	Location    string `tfschema:"location"`
	SubLocation string `tfschema:"sub_location"`
	Type        string `tfschema:"type"`
}

type PipelineTlsCertificateModel struct {
	Certificate []PipelineCertificateSourceModel `tfschema:"certificate"`
	PrivateKey  []PipelinePrivateKeySourceModel  `tfschema:"private_key"`
}

type PipelineSchemaMapModel struct {
	RecordMap   []PipelineFieldMapModel `tfschema:"record_map"`
	ResourceMap []PipelineFieldMapModel `tfschema:"resource_map"`
	ScopeMap    []PipelineFieldMapModel `tfschema:"scope_map"`
}

type PipelinePrivateKeySourceModel struct {
	Location    string `tfschema:"location"`
	SubLocation string `tfschema:"sub_location"`
}

type PipelineFieldMapModel struct {
	From string `tfschema:"from"`
	To   string `tfschema:"to"`
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

		"log_pipeline": {
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

		"azure_monitor_workspace_log_exporter": {
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
												Elem:     pipelineFieldMapSchema(),
											},

											"resource_map": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem:     pipelineFieldMapSchema(),
											},

											"scope_map": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem:     pipelineFieldMapSchema(),
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

					"persistence_maximum_storage_usage_in_gb": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},

					// The service caps retention at 2 days (2880 minutes).
					"persistence_retention_period_in_minutes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 2880),
					},
				},
			},
		},

		"batch_processor": {
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

					"batch_size": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      defaultPipelineBatchSize,
						ValidateFunc: validation.IntBetween(10, 100000),
					},

					"timeout_in_milliseconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      defaultPipelineTimeoutInMilliseconds,
						ValidateFunc: validation.IntBetween(10, 300000),
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

		"microsoft_common_security_log_processor": {
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
				},
			},
		},

		"microsoft_syslog_processor": {
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
				},
			},
		},

		"otlp_receiver": {
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

					"endpoint": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^(?:[a-zA-Z][a-zA-Z0-9+.-]*://)?(?:\[[0-9a-fA-F:.]+\]|[^:/?#[:space:]]*):(?:[1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`),
							"must be in the form `<host>:<port>` with a numeric port between `1` and `65535`",
						),
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

		"persistent_volume_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"replicas": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"syslog_receiver": {
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
						// MaxItems matches the number of AllowedFormats enum values; no combination can exceed this.
						MaxItems: len(pipelinegroups.PossibleValuesForAllowedFormats()),
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(pipelinegroups.PossibleValuesForAllowedFormats(), false),
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

					"transport_protocol": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(pipelinegroups.TransportProtocolTcp),
						ValidateFunc: validation.StringInSlice(pipelinegroups.PossibleValuesForTransportProtocol(), false),
					},
				},
			},
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
						Elem:     pipelineCertificateSourceSchema(),
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
									Elem:     pipelineCertificateSourceSchema(),
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

		"transform_language_processor": {
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

					"transform_statement": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 10000),
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

// shared by the recordMap, resourceMap and scopeMap blocks.
func pipelineFieldMapSchema() *pluginsdk.Resource {
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
func pipelineCertificateSourceSchema() *pluginsdk.Resource {
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

func (r MonitorPipelineResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r MonitorPipelineResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			rawConfig := metadata.ResourceDiff.GetRawConfig()
			if !rawConfig.IsKnown() || rawConfig.IsNull() {
				return nil
			}

			constraints := rawConfig.GetAttr("execution_placement_constraint")
			if constraints.IsKnown() && !constraints.IsNull() {
				for i, constraint := range constraints.AsValueSlice() {
					if !constraint.IsKnown() || constraint.IsNull() {
						continue
					}

					operator := constraint.GetAttr("operator")
					values := constraint.GetAttr("values")
					if !operator.IsKnown() || operator.IsNull() || !values.IsKnown() {
						continue
					}

					switch operator.AsString() {
					case string(pipelinegroups.CapabilityOperatorIn), string(pipelinegroups.CapabilityOperatorNotIn):
						if values.IsNull() || values.LengthInt() == 0 {
							return fmt.Errorf("`values` must be set and non-empty for `execution_placement_constraint.%d` when `operator` is `%s`", i, operator.AsString())
						}
					case string(pipelinegroups.CapabilityOperatorExists), string(pipelinegroups.CapabilityOperatorDoesNotExist):
						if !values.IsNull() && values.LengthInt() > 0 {
							return fmt.Errorf("`values` must not be set for `execution_placement_constraint.%d` when `operator` is `%s`", i, operator.AsString())
						}
					}
				}
			}

			exporters := rawConfig.GetAttr("azure_monitor_workspace_log_exporter")
			requiresPersistentVolume := false
			if exporters.IsKnown() && !exporters.IsNull() {
				for i, exporter := range exporters.AsValueSlice() {
					if !exporter.IsKnown() || exporter.IsNull() {
						continue
					}

					maximumStorageUsage := exporter.GetAttr("persistence_maximum_storage_usage_in_gb")
					retentionPeriod := exporter.GetAttr("persistence_retention_period_in_minutes")
					if (maximumStorageUsage.IsKnown() && !maximumStorageUsage.IsNull()) || (retentionPeriod.IsKnown() && !retentionPeriod.IsNull()) {
						requiresPersistentVolume = true
					}

					api := exporter.GetAttr("api")
					if !api.IsKnown() || api.IsNull() || api.LengthInt() == 0 {
						continue
					}

					apiConfig := api.Index(cty.NumberIntVal(0))
					if !apiConfig.IsKnown() || apiConfig.IsNull() {
						continue
					}

					schemaConfig := apiConfig.GetAttr("schema")
					if !schemaConfig.IsKnown() || schemaConfig.IsNull() || schemaConfig.LengthInt() == 0 {
						continue
					}

					schema := schemaConfig.Index(cty.NumberIntVal(0))
					if !schema.IsKnown() || schema.IsNull() {
						continue
					}

					recordMaps := schema.GetAttr("record_map")
					if !recordMaps.IsKnown() || recordMaps.IsNull() {
						continue
					}

					hasTimeGenerated := false
					recordMapsKnown := true
					for _, recordMap := range recordMaps.AsValueSlice() {
						if !recordMap.IsKnown() || recordMap.IsNull() {
							recordMapsKnown = false
							break
						}

						to := recordMap.GetAttr("to")
						if !to.IsKnown() || to.IsNull() {
							recordMapsKnown = false
							break
						}
						if to.AsString() == "TimeGenerated" {
							hasTimeGenerated = true
							break
						}
					}

					if recordMapsKnown && !hasTimeGenerated {
						return fmt.Errorf("`azure_monitor_workspace_log_exporter.%d.api.0.schema.0.record_map` must define an entry with `to` set to `TimeGenerated`", i)
					}
				}
			}

			persistentVolumeName := rawConfig.GetAttr("persistent_volume_name")
			if requiresPersistentVolume && persistentVolumeName.IsKnown() && (persistentVolumeName.IsNull() || persistentVolumeName.AsString() == "") {
				return errors.New("`persistent_volume_name` must be set when `persistence_maximum_storage_usage_in_gb` or `persistence_retention_period_in_minutes` is configured")
			}

			tlsConfigurations := rawConfig.GetAttr("tls_configuration")
			tlsConfigurationNames := make(map[string]struct{})
			tlsConfigurationNamesKnown := tlsConfigurations.IsKnown()
			if tlsConfigurationNamesKnown && !tlsConfigurations.IsNull() {
				for _, tlsConfiguration := range tlsConfigurations.AsValueSlice() {
					if !tlsConfiguration.IsKnown() || tlsConfiguration.IsNull() {
						tlsConfigurationNamesKnown = false
						break
					}

					name := tlsConfiguration.GetAttr("name")
					if !name.IsKnown() || name.IsNull() {
						tlsConfigurationNamesKnown = false
						break
					}
					tlsConfigurationNames[name.AsString()] = struct{}{}
				}
			}

			otlpReceivers := rawConfig.GetAttr("otlp_receiver")
			if otlpReceivers.IsKnown() && !otlpReceivers.IsNull() {
				for i, receiver := range otlpReceivers.AsValueSlice() {
					if !receiver.IsKnown() || receiver.IsNull() {
						continue
					}

					tlsConfigurationName := receiver.GetAttr("tls_configuration_name")
					if tlsConfigurationName.IsKnown() && !tlsConfigurationName.IsNull() && tlsConfigurationNamesKnown {
						if _, ok := tlsConfigurationNames[tlsConfigurationName.AsString()]; !ok {
							return fmt.Errorf("`otlp_receiver.%d.tls_configuration_name` references unknown `tls_configuration` `%s`", i, tlsConfigurationName.AsString())
						}
					}
				}
			}

			syslogReceivers := rawConfig.GetAttr("syslog_receiver")
			if syslogReceivers.IsKnown() && !syslogReceivers.IsNull() {
				for i, receiver := range syslogReceivers.AsValueSlice() {
					if !receiver.IsKnown() || receiver.IsNull() {
						continue
					}

					tlsConfigurationName := receiver.GetAttr("tls_configuration_name")
					if tlsConfigurationName.IsKnown() && !tlsConfigurationName.IsNull() && tlsConfigurationNamesKnown {
						if _, ok := tlsConfigurationNames[tlsConfigurationName.AsString()]; !ok {
							return fmt.Errorf("`syslog_receiver.%d.tls_configuration_name` references unknown `tls_configuration` `%s`", i, tlsConfigurationName.AsString())
						}
					}

					transportProtocol := receiver.GetAttr("transport_protocol")
					if tlsConfigurationName.IsKnown() && !tlsConfigurationName.IsNull() && transportProtocol.IsKnown() && !transportProtocol.IsNull() && transportProtocol.AsString() == string(pipelinegroups.TransportProtocolUdp) {
						return fmt.Errorf("`tls_configuration_name` is not supported for `syslog_receiver.%d` when `transport_protocol` is `%s`", i, pipelinegroups.TransportProtocolUdp)
					}

					allowedFormats := receiver.GetAttr("allowed_formats")
					if !allowedFormats.IsKnown() || allowedFormats.IsNull() || allowedFormats.LengthInt() == 0 {
						continue
					}

					hasAll := false
					hasSkipPriorityHeaderCompatibleFormat := false
					allowedFormatsKnown := true
					for _, format := range allowedFormats.AsValueSlice() {
						if !format.IsKnown() || format.IsNull() {
							allowedFormatsKnown = false
							break
						}

						switch format.AsString() {
						case string(pipelinegroups.AllowedFormatsAll):
							hasAll = true
							hasSkipPriorityHeaderCompatibleFormat = true
						case string(pipelinegroups.AllowedFormatsSyslogRfcThreeOneSixFour), string(pipelinegroups.AllowedFormatsCefRfcThreeOneSixFour):
							hasSkipPriorityHeaderCompatibleFormat = true
						}
					}

					if !allowedFormatsKnown {
						continue
					}
					if hasAll && allowedFormats.LengthInt() > 1 {
						return fmt.Errorf("`allowed_formats` for `syslog_receiver.%d` must not combine `%s` with other formats", i, pipelinegroups.AllowedFormatsAll)
					}

					allowSkipPriorityHeader := receiver.GetAttr("allow_skip_priority_header")
					if allowSkipPriorityHeader.IsKnown() && !allowSkipPriorityHeader.IsNull() && allowSkipPriorityHeader.True() && !hasSkipPriorityHeaderCompatibleFormat {
						return fmt.Errorf("`allow_skip_priority_header` for `syslog_receiver.%d` requires `allowed_formats` to include `%s`, `%s`, or `%s`", i, pipelinegroups.AllowedFormatsAll, pipelinegroups.AllowedFormatsSyslogRfcThreeOneSixFour, pipelinegroups.AllowedFormatsCefRfcThreeOneSixFour)
					}
				}
			}

			if tlsConfigurations.IsKnown() && !tlsConfigurations.IsNull() {
				for i, tlsConfiguration := range tlsConfigurations.AsValueSlice() {
					if !tlsConfiguration.IsKnown() || tlsConfiguration.IsNull() {
						continue
					}

					mode := tlsConfiguration.GetAttr("mode")
					if !mode.IsKnown() || mode.IsNull() {
						continue
					}

					clientCertificateAuthority := tlsConfiguration.GetAttr("client_certificate_authority")
					tlsCertificate := tlsConfiguration.GetAttr("tls_certificate")
					switch mode.AsString() {
					case string(pipelinegroups.TlsModeDisabled):
						clientCertificateAuthorityConfigured := clientCertificateAuthority.IsKnown() && !clientCertificateAuthority.IsNull() && clientCertificateAuthority.LengthInt() > 0
						tlsCertificateConfigured := tlsCertificate.IsKnown() && !tlsCertificate.IsNull() && tlsCertificate.LengthInt() > 0
						if clientCertificateAuthorityConfigured || tlsCertificateConfigured {
							return fmt.Errorf("`tls_certificate` and `client_certificate_authority` must not be set for `tls_configuration.%d` when `mode` is `%s`", i, mode.AsString())
						}
					case string(pipelinegroups.TlsModeServerOnly):
						if clientCertificateAuthority.IsKnown() && !clientCertificateAuthority.IsNull() && clientCertificateAuthority.LengthInt() > 0 {
							return fmt.Errorf("`client_certificate_authority` must not be set for `tls_configuration.%d` when `mode` is `%s`", i, mode.AsString())
						}
					}
				}
			}

			return nil
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
					ExecutionPlacement: expandPipelineExecutionPlacement(model.ExecutionPlacementConstraint),
					Exporters:          expandPipelineAzureMonitorWorkspaceLogExporters(model.AzureMonitorWorkspaceLogExporter),
					Processors:         expandPipelineProcessors(model),
					Receivers:          expandPipelineReceivers(model),
					Replicas:           pointer.ToOrNil(model.Replicas),
					Service:            expandPipelineService(model.LogPipeline, model.PersistentVolumeName),
				},
				Tags: pointer.To(model.Tags),
			}

			if len(model.TlsConfiguration) > 0 {
				payload.Properties.TlsConfigurations = expandPipelineTlsConfigurations(model.TlsConfiguration)
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
				payload.Properties.ExecutionPlacement = expandPipelineExecutionPlacement(model.ExecutionPlacementConstraint)
			}

			if metadata.ResourceData.HasChange("azure_monitor_workspace_log_exporter") {
				payload.Properties.Exporters = expandPipelineAzureMonitorWorkspaceLogExporters(model.AzureMonitorWorkspaceLogExporter)
			}

			if metadata.ResourceData.HasChanges("batch_processor", "microsoft_common_security_log_processor", "microsoft_syslog_processor", "transform_language_processor") {
				payload.Properties.Processors = expandPipelineProcessors(model)
			}

			if metadata.ResourceData.HasChanges("otlp_receiver", "syslog_receiver") {
				payload.Properties.Receivers = expandPipelineReceivers(model)
			}

			if metadata.ResourceData.HasChange("replicas") {
				payload.Properties.Replicas = pointer.ToOrNil(model.Replicas)
			}

			if metadata.ResourceData.HasChanges("log_pipeline", "persistent_volume_name") {
				payload.Properties.Service = expandPipelineService(model.LogPipeline, model.PersistentVolumeName)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if metadata.ResourceData.HasChange("tls_configuration") {
				payload.Properties.TlsConfigurations = expandPipelineTlsConfigurations(model.TlsConfiguration)
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
			state.ExecutionPlacementConstraint = flattenPipelineExecutionPlacement(props.ExecutionPlacement)
			state.AzureMonitorWorkspaceLogExporter = flattenPipelineAzureMonitorWorkspaceLogExporters(props.Exporters)
			flattenPipelineProcessors(props.Processors, &state)
			flattenPipelineReceivers(props.Receivers, &state)
			state.Replicas = pointer.From(props.Replicas)
			flattenPipelineService(props.Service, &state)
			state.TlsConfiguration = flattenPipelineTlsConfigurations(props.TlsConfigurations)
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (r MonitorPipelineResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: time.Hour,
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

func expandPipelineExecutionPlacement(input []PipelinePlacementConstraintModel) *pipelinegroups.ExecutionPlacement {
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

func flattenPipelineExecutionPlacement(input *pipelinegroups.ExecutionPlacement) []PipelinePlacementConstraintModel {
	if input == nil || input.Constraints == nil {
		return []PipelinePlacementConstraintModel{}
	}

	output := make([]PipelinePlacementConstraintModel, 0, len(*input.Constraints))
	for _, constraint := range *input.Constraints {
		output = append(output, PipelinePlacementConstraintModel{
			Capability: constraint.Capability,
			Operator:   string(constraint.Operator),
			Values:     pointer.From(constraint.Values),
		})
	}

	return output
}

func expandPipelineAzureMonitorWorkspaceLogExporters(input []PipelineAzureMonitorWorkspaceLogExporterModel) []pipelinegroups.Exporter {
	output := make([]pipelinegroups.Exporter, 0, len(input))

	for _, v := range input {
		output = append(output, pipelinegroups.Exporter{
			Name:                      v.Name,
			Type:                      pipelinegroups.ExporterTypeAzureMonitorWorkspaceLogs,
			AzureMonitorWorkspaceLogs: expandPipelineAzureMonitorWorkspaceLogExporter(v),
		})
	}

	return output
}

func expandPipelineAzureMonitorWorkspaceLogExporter(input PipelineAzureMonitorWorkspaceLogExporterModel) *pipelinegroups.AzureMonitorWorkspaceLogsExporter {
	output := pipelinegroups.AzureMonitorWorkspaceLogsExporter{}

	if len(input.Api) > 0 {
		api := input.Api[0]
		output.Api = pipelinegroups.AzureMonitorWorkspaceLogsApiConfig{
			DataCollectionEndpointURL: api.DataCollectionEndpointUrl,
			DataCollectionRule:        api.DataCollectionRuleImmutableId,
			Stream:                    api.Stream,
		}

		if len(api.Schema) > 0 {
			output.Api.Schema = expandPipelineSchemaMap(api.Schema[0])
		}
	}

	if input.PersistenceMaximumStorageUsageInGB != 0 || input.PersistenceRetentionPeriodInMinutes != 0 {
		output.Persistence = &pipelinegroups.ExporterPersistenceConfiguration{
			MaxStorageUsage: pointer.ToOrNil(input.PersistenceMaximumStorageUsageInGB),
			RetentionPeriod: pointer.ToOrNil(input.PersistenceRetentionPeriodInMinutes),
		}
	}

	return &output
}

func expandPipelineSchemaMap(input PipelineSchemaMapModel) pipelinegroups.SchemaMap {
	output := pipelinegroups.SchemaMap{
		RecordMap: expandPipelineFieldMaps(input.RecordMap),
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

func expandPipelineFieldMaps(input []PipelineFieldMapModel) []pipelinegroups.RecordMap {
	output := make([]pipelinegroups.RecordMap, 0, len(input))
	for _, v := range input {
		output = append(output, pipelinegroups.RecordMap{From: v.From, To: v.To})
	}
	return output
}

func flattenPipelineAzureMonitorWorkspaceLogExporters(input []pipelinegroups.Exporter) []PipelineAzureMonitorWorkspaceLogExporterModel {
	output := make([]PipelineAzureMonitorWorkspaceLogExporterModel, 0, len(input))

	for _, v := range input {
		if v.Type != pipelinegroups.ExporterTypeAzureMonitorWorkspaceLogs || v.AzureMonitorWorkspaceLogs == nil {
			continue
		}

		exporter := flattenPipelineAzureMonitorWorkspaceLogExporter(*v.AzureMonitorWorkspaceLogs)
		exporter.Name = v.Name
		output = append(output, exporter)
	}

	return output
}

func flattenPipelineAzureMonitorWorkspaceLogExporter(input pipelinegroups.AzureMonitorWorkspaceLogsExporter) PipelineAzureMonitorWorkspaceLogExporterModel {
	output := PipelineAzureMonitorWorkspaceLogExporterModel{
		Api: []PipelineAzureMonitorWorkspaceLogsApiModel{
			{
				DataCollectionEndpointUrl:     input.Api.DataCollectionEndpointURL,
				DataCollectionRuleImmutableId: input.Api.DataCollectionRule,
				Stream:                        input.Api.Stream,
				Schema:                        flattenPipelineSchemaMap(input.Api.Schema),
			},
		},
	}

	if input.Persistence != nil {
		output.PersistenceMaximumStorageUsageInGB = pointer.From(input.Persistence.MaxStorageUsage)
		output.PersistenceRetentionPeriodInMinutes = pointer.From(input.Persistence.RetentionPeriod)
	}

	return output
}

func flattenPipelineSchemaMap(input pipelinegroups.SchemaMap) []PipelineSchemaMapModel {
	output := PipelineSchemaMapModel{
		RecordMap: flattenPipelineFieldMaps(input.RecordMap),
	}

	if input.ResourceMap != nil {
		resourceMaps := make([]PipelineFieldMapModel, 0, len(*input.ResourceMap))
		for _, v := range *input.ResourceMap {
			resourceMaps = append(resourceMaps, PipelineFieldMapModel{From: v.From, To: v.To})
		}
		output.ResourceMap = resourceMaps
	}

	if input.ScopeMap != nil {
		scopeMaps := make([]PipelineFieldMapModel, 0, len(*input.ScopeMap))
		for _, v := range *input.ScopeMap {
			scopeMaps = append(scopeMaps, PipelineFieldMapModel{From: v.From, To: v.To})
		}
		output.ScopeMap = scopeMaps
	}

	return []PipelineSchemaMapModel{output}
}

func flattenPipelineFieldMaps(input []pipelinegroups.RecordMap) []PipelineFieldMapModel {
	output := make([]PipelineFieldMapModel, 0, len(input))
	for _, v := range input {
		output = append(output, PipelineFieldMapModel{From: v.From, To: v.To})
	}
	return output
}

func expandPipelineProcessors(input MonitorPipelineResourceModel) []pipelinegroups.Processor {
	output := make([]pipelinegroups.Processor, 0, len(input.BatchProcessor)+len(input.MicrosoftCommonSecurityLogProcessor)+len(input.MicrosoftSyslogProcessor)+len(input.TransformLanguageProcessor))

	for _, v := range input.BatchProcessor {
		output = append(output, pipelinegroups.Processor{
			Name: v.Name,
			Type: pipelinegroups.ProcessorTypeBatch,
			Batch: &pipelinegroups.BatchProcessor{
				BatchSize: pointer.To(v.BatchSize),
				Timeout:   pointer.To(v.TimeoutInMilliseconds),
			},
		})
	}

	for _, v := range input.MicrosoftCommonSecurityLogProcessor {
		output = append(output, pipelinegroups.Processor{Name: v.Name, Type: pipelinegroups.ProcessorTypeMicrosoftCommonSecurityLog})
	}

	for _, v := range input.MicrosoftSyslogProcessor {
		output = append(output, pipelinegroups.Processor{Name: v.Name, Type: pipelinegroups.ProcessorTypeMicrosoftSyslog})
	}

	for _, v := range input.TransformLanguageProcessor {
		output = append(output, pipelinegroups.Processor{
			Name: v.Name,
			Type: pipelinegroups.ProcessorTypeTransformLanguage,
			TransformLanguage: &pipelinegroups.TransformLanguageProcessor{
				TransformStatement: v.TransformStatement,
			},
		})
	}

	return output
}

func flattenPipelineProcessors(input []pipelinegroups.Processor, output *MonitorPipelineResourceModel) {
	for _, v := range input {
		switch v.Type {
		case pipelinegroups.ProcessorTypeBatch:
			processor := PipelineBatchProcessorModel{Name: v.Name, BatchSize: defaultPipelineBatchSize, TimeoutInMilliseconds: defaultPipelineTimeoutInMilliseconds}
			if v.Batch != nil {
				processor.BatchSize = pointer.From(v.Batch.BatchSize)
				processor.TimeoutInMilliseconds = pointer.From(v.Batch.Timeout)
			}
			output.BatchProcessor = append(output.BatchProcessor, processor)
		case pipelinegroups.ProcessorTypeMicrosoftCommonSecurityLog:
			output.MicrosoftCommonSecurityLogProcessor = append(output.MicrosoftCommonSecurityLogProcessor, PipelineMicrosoftCommonSecurityLogProcessorModel{Name: v.Name})
		case pipelinegroups.ProcessorTypeMicrosoftSyslog:
			output.MicrosoftSyslogProcessor = append(output.MicrosoftSyslogProcessor, PipelineMicrosoftSyslogProcessorModel{Name: v.Name})
		case pipelinegroups.ProcessorTypeTransformLanguage:
			processor := PipelineTransformLanguageProcessorModel{Name: v.Name}
			if v.TransformLanguage != nil {
				processor.TransformStatement = v.TransformLanguage.TransformStatement
			}
			output.TransformLanguageProcessor = append(output.TransformLanguageProcessor, processor)
		}
	}
}

func expandPipelineReceivers(input MonitorPipelineResourceModel) []pipelinegroups.Receiver {
	output := make([]pipelinegroups.Receiver, 0, len(input.OtlpReceiver)+len(input.SyslogReceiver))

	for _, v := range input.OtlpReceiver {
		output = append(output, pipelinegroups.Receiver{
			Name:             v.Name,
			Type:             pipelinegroups.ReceiverTypeOTLP,
			Otlp:             &pipelinegroups.OtlpReceiver{Endpoint: v.Endpoint},
			TlsConfiguration: pointer.ToOrNil(v.TlsConfigurationName),
		})
	}

	for _, v := range input.SyslogReceiver {
		syslog := &pipelinegroups.SyslogReceiver{
			Endpoint:           v.Endpoint,
			AllowSkipPriHeader: &v.AllowSkipPriorityHeader,
			TransportProtocol:  pointer.ToEnum[pipelinegroups.TransportProtocol](v.TransportProtocol),
		}
		if len(v.AllowedFormats) > 0 {
			syslog.AllowedFormats = pointer.ToEnumSlice[pipelinegroups.AllowedFormats](v.AllowedFormats)
		}

		output = append(output, pipelinegroups.Receiver{
			Name:             v.Name,
			Type:             pipelinegroups.ReceiverTypeSyslog,
			Syslog:           syslog,
			TlsConfiguration: pointer.ToOrNil(v.TlsConfigurationName),
		})
	}

	return output
}

func flattenPipelineReceivers(input []pipelinegroups.Receiver, output *MonitorPipelineResourceModel) {
	for _, v := range input {
		switch v.Type {
		case pipelinegroups.ReceiverTypeOTLP:
			receiver := PipelineOtlpReceiverModel{Name: v.Name, TlsConfigurationName: pointer.From(v.TlsConfiguration)}
			if v.Otlp != nil {
				receiver.Endpoint = v.Otlp.Endpoint
			}
			output.OtlpReceiver = append(output.OtlpReceiver, receiver)
		case pipelinegroups.ReceiverTypeSyslog:
			receiver := PipelineSyslogReceiverModel{Name: v.Name, TlsConfigurationName: pointer.From(v.TlsConfiguration)}
			if v.Syslog != nil {
				receiver.Endpoint = v.Syslog.Endpoint
				receiver.AllowSkipPriorityHeader = pointer.From(v.Syslog.AllowSkipPriHeader)
				receiver.AllowedFormats = pointer.FromEnumSlice(v.Syslog.AllowedFormats)
				receiver.TransportProtocol = pointer.FromEnum(v.Syslog.TransportProtocol)
			}
			output.SyslogReceiver = append(output.SyslogReceiver, receiver)
		}
	}
}

func expandPipelineService(input []PipelineLogPipelineModel, persistentVolumeName string) pipelinegroups.Service {
	output := pipelinegroups.Service{
		Pipelines: []pipelinegroups.Pipeline{},
	}

	if persistentVolumeName != "" {
		output.Persistence = &pipelinegroups.PersistenceConfigurations{
			PersistentVolumeName: persistentVolumeName,
		}
	}

	pipelines := make([]pipelinegroups.Pipeline, 0, len(input))
	for _, p := range input {
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

func flattenPipelineService(input pipelinegroups.Service, output *MonitorPipelineResourceModel) {
	output.PersistentVolumeName = pointer.From(input.Persistence).PersistentVolumeName
	output.LogPipeline = make([]PipelineLogPipelineModel, 0, len(input.Pipelines))
	for _, p := range input.Pipelines {
		if p.Type != pipelinegroups.PipelineTypeLogs {
			continue
		}
		output.LogPipeline = append(output.LogPipeline, PipelineLogPipelineModel{
			Name:       p.Name,
			Exporters:  p.Exporters,
			Receivers:  p.Receivers,
			Processors: pointer.From(p.Processors),
		})
	}
}

func expandPipelineTlsConfigurations(input []PipelineTlsConfigurationModel) *[]pipelinegroups.TlsConfiguration {
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
			tlsConfiguration.ClientCa = expandPipelineCertificateSource(v.ClientCertificateAuthority[0])
		}

		if len(v.TlsCertificate) > 0 {
			certificate := v.TlsCertificate[0]
			certificateWithKey := pipelinegroups.CertificateWithKey{}

			if len(certificate.Certificate) > 0 {
				certificateWithKey.Certificate = *expandPipelineCertificateSource(certificate.Certificate[0])
			}
			if len(certificate.PrivateKey) > 0 {
				certificateWithKey.PrivateKey = pipelinegroups.PrivateKeySource{
					Location:    certificate.PrivateKey[0].Location,
					SubLocation: certificate.PrivateKey[0].SubLocation,
					Type:        pipelinegroups.PrivateKeySourceTypeKubernetesSecret,
				}
			}

			tlsConfiguration.TlsCertificate = &certificateWithKey
		}

		output = append(output, tlsConfiguration)
	}

	return pointer.To(output)
}

func expandPipelineCertificateSource(input PipelineCertificateSourceModel) *pipelinegroups.CertificateSource {
	return &pipelinegroups.CertificateSource{
		Location:    input.Location,
		SubLocation: input.SubLocation,
		Type:        pipelinegroups.CertificateSourceType(input.Type),
	}
}

func flattenPipelineTlsConfigurations(input *[]pipelinegroups.TlsConfiguration) []PipelineTlsConfigurationModel {
	if input == nil {
		return []PipelineTlsConfigurationModel{}
	}

	output := make([]PipelineTlsConfigurationModel, 0, len(*input))

	for _, v := range *input {
		tlsConfiguration := PipelineTlsConfigurationModel{
			Name: v.Name,
			Mode: pointer.FromEnum(v.Mode),
		}

		tlsConfiguration.ClientCertificateAuthority = flattenPipelineCertificateSource(v.ClientCa)

		if v.TlsCertificate != nil {
			tlsConfiguration.TlsCertificate = []PipelineTlsCertificateModel{{
				Certificate: flattenPipelineCertificateSource(&v.TlsCertificate.Certificate),
				PrivateKey: []PipelinePrivateKeySourceModel{{
					Location:    v.TlsCertificate.PrivateKey.Location,
					SubLocation: v.TlsCertificate.PrivateKey.SubLocation,
				}},
			}}
		}

		output = append(output, tlsConfiguration)
	}

	return output
}

func flattenPipelineCertificateSource(input *pipelinegroups.CertificateSource) []PipelineCertificateSourceModel {
	if input == nil {
		return []PipelineCertificateSourceModel{}
	}

	return []PipelineCertificateSourceModel{{
		Location:    input.Location,
		SubLocation: input.SubLocation,
		Type:        string(input.Type),
	}}
}
