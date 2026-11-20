// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_pipeline

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/monitor/2026-04-01/pipelinegroups"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
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
										"must start with `dcr-` followed by 32 hexadecimal characters",
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

		"replica_count": {
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

								"private_key_source": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
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

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
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
