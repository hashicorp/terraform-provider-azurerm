// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sre_agent

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-07-01/agents"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				// https://github.com/Azure/azure-rest-api-specs/blob/fdba4d6ef03d157da38bb8d7d6c3384e1fd9f25c/specification/app/resource-manager/Microsoft.App/SreAgent/agent.tsp#L21
				regexp.MustCompile(`^[A-Za-z]([-A-Za-z0-9]{0,30}[A-Za-z0-9])$`),
				"must be between 2 and 32 characters, start with a letter, end with a letter or digit, and contain only letters, digits, and hyphens",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityRequired(),

		"action": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"user_assigned_identity_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},

					"access_level": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(agents.PossibleValuesForAgentAccessLevel(), false),
					},

					"mode": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							// https://learn.microsoft.com/azure/sre-agent/run-modes
							string(agents.AgentModeAutonomous),
							string(agents.AgentModeReview),
						}, false),
					},
				},
			},
		},

		"application_insights": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"application_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},

					"connection_string": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"default_model": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			// NOTE: O+C Azure selects the model and provider when omitted.
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						// Model names vary by subscription, so Azure validates the selected model/provider pair.
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"provider": {
						Type:     pluginsdk.TypeString,
						Required: true,
						// Providers vary by subscription, so Azure validates the selected model/provider pair.
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"incident_management": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						// Valid incident management types vary by subscription, so Azure validates the selected value.
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"connection_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"connection_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"connection_url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},

					"on_behalf_of_user": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"knowledge_graph": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"user_assigned_identity_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},

					"managed_resource_ids": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},
		},

		"sandbox": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			// NOTE: O+C Azure supplies the sandbox configuration when omitted.
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"egress": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						AtLeastOneOf: []string{
							"sandbox.0.egress",
							"sandbox.0.package",
						},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"allowed_code_repositories": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									AtLeastOneOf: []string{
										"sandbox.0.egress.0.allowed_code_repositories",
										"sandbox.0.egress.0.allowed_hosts",
										"sandbox.0.egress.0.allowed_registries",
										"sandbox.0.egress.0.http_mcp_server_network_access_enabled",
										"sandbox.0.egress.0.http_proxy",
										"sandbox.0.egress.0.mode",
									},
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											// https://github.com/Azure/azure-rest-api-specs/blob/fdba4d6ef03d157da38bb8d7d6c3384e1fd9f25c/specification/app/resource-manager/Microsoft.App/SreAgent/agent.models.tsp#L197-L202
											"AzureDevOps",
											"GitLab",
											"Github",
											"GithubEnterprise",
										}, false),
									},
								},

								"allowed_hosts": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 100,
									AtLeastOneOf: []string{
										"sandbox.0.egress.0.allowed_code_repositories",
										"sandbox.0.egress.0.allowed_hosts",
										"sandbox.0.egress.0.allowed_registries",
										"sandbox.0.egress.0.http_mcp_server_network_access_enabled",
										"sandbox.0.egress.0.http_proxy",
										"sandbox.0.egress.0.mode",
									},
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.All(
											// https://github.com/Azure/azure-rest-api-specs/blob/fdba4d6ef03d157da38bb8d7d6c3384e1fd9f25c/specification/app/resource-manager/Microsoft.App/SreAgent/agent.models.tsp#L214-L219
											validation.StringLenBetween(1, 253),
											validation.StringMatch(
												regexp.MustCompile(`^(\*\.[A-Za-z0-9]([A-Za-z0-9-]{0,61}[A-Za-z0-9])?\.)?[A-Za-z0-9]([A-Za-z0-9-]{0,61}[A-Za-z0-9])?(\.[A-Za-z0-9]([A-Za-z0-9-]{0,61}[A-Za-z0-9])?)*$`),
												"must be a DNS hostname with labels of 1 to 63 letters, digits, or hyphens that start and end with a letter or digit, optionally prefixed with `*.` followed by at least two labels, for example `*.example.com`",
											),
										),
									},
								},

								"allowed_registries": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									AtLeastOneOf: []string{
										"sandbox.0.egress.0.allowed_code_repositories",
										"sandbox.0.egress.0.allowed_hosts",
										"sandbox.0.egress.0.allowed_registries",
										"sandbox.0.egress.0.http_mcp_server_network_access_enabled",
										"sandbox.0.egress.0.http_proxy",
										"sandbox.0.egress.0.mode",
									},
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											// https://github.com/Azure/azure-rest-api-specs/blob/fdba4d6ef03d157da38bb8d7d6c3384e1fd9f25c/specification/app/resource-manager/Microsoft.App/SreAgent/agent.models.tsp#L190-L195
											"npmjs",
											"nuget-org",
											"pypi",
											"ubuntu",
										}, false),
									},
								},

								"http_mcp_server_network_access_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
									AtLeastOneOf: []string{
										"sandbox.0.egress.0.allowed_code_repositories",
										"sandbox.0.egress.0.allowed_hosts",
										"sandbox.0.egress.0.allowed_registries",
										"sandbox.0.egress.0.http_mcp_server_network_access_enabled",
										"sandbox.0.egress.0.http_proxy",
										"sandbox.0.egress.0.mode",
									},
								},

								"http_proxy": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									AtLeastOneOf: []string{
										"sandbox.0.egress.0.allowed_code_repositories",
										"sandbox.0.egress.0.allowed_hosts",
										"sandbox.0.egress.0.allowed_registries",
										"sandbox.0.egress.0.http_mcp_server_network_access_enabled",
										"sandbox.0.egress.0.http_proxy",
										"sandbox.0.egress.0.mode",
									},
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											// https://github.com/Azure/azure-rest-api-specs/blob/fdba4d6ef03d157da38bb8d7d6c3384e1fd9f25c/specification/app/resource-manager/Microsoft.App/SreAgent/agent.models.tsp#L232-L250
											"https_url": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.IsURLWithScheme([]string{
													"http",
												}),
											},

											"http_url": {
												Type:     pluginsdk.TypeString,
												Optional: true,
												ValidateFunc: validation.IsURLWithScheme([]string{
													"http",
												}),
											},

											"no_proxy": {
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

								"mode": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringInSlice(agents.PossibleValuesForSandboxEgressMode(), false),
									AtLeastOneOf: []string{
										"sandbox.0.egress.0.allowed_code_repositories",
										"sandbox.0.egress.0.allowed_hosts",
										"sandbox.0.egress.0.allowed_registries",
										"sandbox.0.egress.0.http_mcp_server_network_access_enabled",
										"sandbox.0.egress.0.http_proxy",
										"sandbox.0.egress.0.mode",
									},
								},
							},
						},
					},

					"package": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 200,
						AtLeastOneOf: []string{
							"sandbox.0.egress",
							"sandbox.0.package",
						},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.All(
										validation.StringLenBetween(1, 214),
										validation.StringMatch(
											// https://github.com/Azure/azure-rest-api-specs/blob/fdba4d6ef03d157da38bb8d7d6c3384e1fd9f25c/specification/app/resource-manager/Microsoft.App/SreAgent/agent.models.tsp#L265-L267
											regexp.MustCompile(`^@?[A-Za-z0-9][A-Za-z0-9._-]*(/[A-Za-z0-9][A-Za-z0-9._-]*)?$`),
											"must match `name`, `@name`, `scope/name`, or `@scope/name`, where `scope` and `name` start with a letter or digit and contain only letters, digits, periods, underscores, or hyphens",
										),
									),
								},

								"package_manager": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(agents.PossibleValuesForSandboxPackageManager(), false),
								},

								"version": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.StringMatch(
										// https://github.com/Azure/azure-rest-api-specs/blob/fdba4d6ef03d157da38bb8d7d6c3384e1fd9f25c/specification/app/resource-manager/Microsoft.App/SreAgent/agent.models.tsp#L271-L273
										regexp.MustCompile(`^[A-Za-z0-9._+\-~^<>=,*\[\]()]{1,64}$`),
										"must be between 1 and 64 characters and contain only letters, digits, or the characters `.`, `_`, `+`, `-`, `~`, `^`, `<`, `>`, `=`, `,`, `*`, `[`, `]`, `(`, and `)`, without spaces",
									),
								},
							},
						},
					},
				},
			},
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"upgrade_channel": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(agents.UpgradeChannelStable),
			ValidateFunc: validation.StringInSlice(agents.PossibleValuesForUpgradeChannel(), false),
		},

		"tags": commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"outbound_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			config := metadata.ResourceDiff.GetRawConfig()
			if config.IsNull() || !config.IsKnown() {
				return nil
			}

			identityConfig := config.GetAttr("identity")
			if metadata.ResourceDiff.HasChanges("identity", "action", "knowledge_graph") &&
				identityConfig.IsKnown() && !identityConfig.IsNull() && identityConfig.LengthInt() > 0 {
				identityConfig = identityConfig.AsValueSlice()[0]

				if identityConfig.IsKnown() && !identityConfig.IsNull() && identityConfig.GetAttr("identity_ids").IsWhollyKnown() {
					ids := make([]string, 0)
					identityIds := identityConfig.GetAttr("identity_ids")
					if !identityIds.IsNull() {
						for _, id := range identityIds.AsValueSlice() {
							if !id.IsNull() {
								ids = append(ids, id.AsString())
							}
						}
					}

					action := config.GetAttr("action")
					if metadata.ResourceDiff.HasChanges("identity", "action") &&
						action.IsKnown() && !action.IsNull() && action.LengthInt() > 0 {
						action = action.AsValueSlice()[0]

						if action.IsKnown() && !action.IsNull() {
							value := action.GetAttr("user_assigned_identity_id")
							if value.IsKnown() && !value.IsNull() && value.AsString() != "" {
								found := false

								for _, id := range ids {
									if strings.EqualFold(value.AsString(), id) {
										found = true
										break
									}
								}

								if !found {
									return errors.New("`action.user_assigned_identity_id` must be included in `identity.identity_ids`")
								}
							}
						}
					}

					knowledgeGraph := config.GetAttr("knowledge_graph")
					if metadata.ResourceDiff.HasChanges("identity", "knowledge_graph") &&
						knowledgeGraph.IsKnown() && !knowledgeGraph.IsNull() && knowledgeGraph.LengthInt() > 0 {
						knowledgeGraph = knowledgeGraph.AsValueSlice()[0]

						if knowledgeGraph.IsKnown() && !knowledgeGraph.IsNull() {
							value := knowledgeGraph.GetAttr("user_assigned_identity_id")
							if value.IsKnown() && !value.IsNull() && value.AsString() != "" {
								found := false

								for _, id := range ids {
									if strings.EqualFold(value.AsString(), id) {
										found = true
										break
									}
								}

								if !found {
									return errors.New("`knowledge_graph.user_assigned_identity_id` must be included in `identity.identity_ids`")
								}
							}
						}
					}
				}
			}

			if metadata.ResourceDiff.HasChanges("subnet_id", "sandbox") {
				sandbox := config.GetAttr("sandbox")
				if !sandbox.IsKnown() {
					return nil
				}

				mode := ""
				if !sandbox.IsNull() && sandbox.LengthInt() > 0 {
					sandbox = sandbox.AsValueSlice()[0]
					if !sandbox.IsKnown() {
						return nil
					}

					if !sandbox.IsNull() {
						egress := sandbox.GetAttr("egress")
						if !egress.IsKnown() {
							return nil
						}

						if !egress.IsNull() && egress.LengthInt() > 0 {
							egress = egress.AsValueSlice()[0]
							if !egress.IsKnown() {
								return nil
							}

							if !egress.IsNull() {
								proxy := egress.GetAttr("http_proxy")
								if proxy.IsKnown() && !proxy.IsNull() && proxy.LengthInt() > 0 {
									proxy = proxy.AsValueSlice()[0]
									if proxy.IsKnown() && !proxy.IsNull() {
										httpURL := proxy.GetAttr("http_url")
										httpsURL := proxy.GetAttr("https_url")
										if httpURL.IsKnown() && !httpURL.IsNull() && httpsURL.IsKnown() && !httpsURL.IsNull() {
											httpProxyURL, err := url.Parse(httpURL.AsString())
											if err != nil {
												return fmt.Errorf("parsing `sandbox.egress.http_proxy.http_url`: %+v", err)
											}

											httpsProxyURL, err := url.Parse(httpsURL.AsString())
											if err != nil {
												return fmt.Errorf("parsing `sandbox.egress.http_proxy.https_url`: %+v", err)
											}

											if !strings.EqualFold(httpProxyURL.Scheme, httpsProxyURL.Scheme) ||
												!strings.EqualFold(httpProxyURL.Hostname(), httpsProxyURL.Hostname()) ||
												httpProxyURL.Port() != httpsProxyURL.Port() {
												return errors.New("`sandbox.egress.http_proxy.http_url` must use the same origin (scheme, host, and port) as `sandbox.egress.http_proxy.https_url`")
											}
										}
									}
								}

								value := egress.GetAttr("mode")
								if !value.IsKnown() {
									return nil
								}

								if !value.IsNull() {
									mode = value.AsString()
								}
							}
						}
					}
				} else {
					mode, _ = metadata.ResourceDiff.Get("sandbox.0.egress.0.mode").(string)
				}

				subnet := config.GetAttr("subnet_id")
				if !subnet.IsKnown() {
					return nil
				}

				hasSubnet := !subnet.IsNull() && subnet.AsString() != ""
				if hasSubnet != (mode == string(agents.SandboxEgressModeAzureVNet)) {
					return errors.New("`subnet_id` is required when `sandbox.egress.mode` is set to `AzureVNet` and must not be set otherwise")
				}
			}

			return nil
		},
	}
}
