// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sreagents

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-07-01/agents"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity

type SREAgentResource struct{}

var (
	_ sdk.ResourceWithIdentity      = SREAgentResource{}
	_ sdk.ResourceWithUpdate        = SREAgentResource{}
	_ sdk.ResourceWithCustomizeDiff = SREAgentResource{}
)

type SREAgentModel struct {
	Name                  string                                     `tfschema:"name"`
	ResourceGroupName     string                                     `tfschema:"resource_group_name"`
	Location              string                                     `tfschema:"location"`
	Identity              []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Action                []ActionModel                              `tfschema:"action"`
	ApplicationInsights   []ApplicationInsightsModel                 `tfschema:"application_insights"`
	DefaultModel          []DefaultModel                             `tfschema:"default_model"`
	IncidentManagement    []IncidentManagementModel                  `tfschema:"incident_management"`
	InitialSponsorGroupId string                                     `tfschema:"initial_sponsor_group_id"`
	KnowledgeGraph        []KnowledgeGraphModel                      `tfschema:"knowledge_graph"`
	Sandbox               []SandboxModel                             `tfschema:"sandbox"`
	SubnetId              string                                     `tfschema:"subnet_id"`
	UpgradeChannel        string                                     `tfschema:"upgrade_channel"`
	Tags                  map[string]string                          `tfschema:"tags"`
	Endpoint              string                                     `tfschema:"endpoint"`
	OutboundIPAddresses   []string                                   `tfschema:"outbound_ip_addresses"`
}

type ActionModel struct {
	UserAssignedIdentityId string `tfschema:"user_assigned_identity_id"`
	AccessLevel            string `tfschema:"access_level"`
	Mode                   string `tfschema:"mode"`
}

type ApplicationInsightsModel struct {
	AppId            string `tfschema:"application_id"`
	ConnectionString string `tfschema:"connection_string"`
}

type DefaultModel struct {
	Name     string `tfschema:"name"`
	Provider string `tfschema:"provider"`
}

type IncidentManagementModel struct {
	Type           string `tfschema:"type"`
	ConnectionKey  string `tfschema:"connection_key"`
	ConnectionName string `tfschema:"connection_name"`
	ConnectionURL  string `tfschema:"connection_url"`
	OboUser        string `tfschema:"obo_user"`
}

type KnowledgeGraphModel struct {
	UserAssignedIdentityId string   `tfschema:"user_assigned_identity_id"`
	ManagedResourceIds     []string `tfschema:"managed_resource_ids"`
}

type SandboxModel struct {
	Egress  []SandboxEgressModel  `tfschema:"egress"`
	Package []SandboxPackageModel `tfschema:"package"`
}

type SandboxEgressModel struct {
	AllowedCodeRepositories           []string                `tfschema:"allowed_code_repositories"`
	AllowedHosts                      []string                `tfschema:"allowed_hosts"`
	AllowedRegistries                 []string                `tfschema:"allowed_registries"`
	HTTPMCPServerNetworkAccessEnabled bool                    `tfschema:"http_mcp_server_network_access_enabled"`
	HTTPProxy                         []SandboxHTTPProxyModel `tfschema:"http_proxy"`
	Mode                              string                  `tfschema:"mode"`
}

type SandboxHTTPProxyModel struct {
	HTTPSURL string   `tfschema:"https_url"`
	HTTPURL  string   `tfschema:"http_url"`
	NoProxy  []string `tfschema:"no_proxy"`
}

type SandboxPackageModel struct {
	Name           string `tfschema:"name"`
	PackageManager string `tfschema:"package_manager"`
	Version        string `tfschema:"version"`
}

func (SREAgentResource) ResourceType() string {
	return "azurerm_sre_agent"
}

func (SREAgentResource) ModelObject() any {
	return &SREAgentModel{}
}

func (SREAgentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return agents.ValidateAgentID
}

func (SREAgentResource) Identity() resourceids.ResourceId {
	return &agents.AgentId{}
}

func (SREAgentResource) Arguments() map[string]*pluginsdk.Schema {
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

					"obo_user": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"initial_sponsor_group_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
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
											// https://github.com/Azure/azure-rest-api-specs/blob/fdba4d6ef03d157da38bb8d7d6c3384e1fd9f25c/specification/app/resource-manager/Microsoft.App/SreAgent/agent.models.tsp#L180-L188
											validation.StringLenBetween(1, 253),
											validation.StringMatch(
												regexp.MustCompile(`^(\*\.[A-Za-z0-9][A-Za-z0-9.-]*\.)?[A-Za-z0-9]([A-Za-z0-9.-]*[A-Za-z0-9])?$`),
												"must be a hostname or a hostname with a leading wildcard `*.` followed by at least two hostname labels, for example `*.example.com`",
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
											"https_url": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.IsURLWithHTTPorHTTPS,
											},

											"http_url": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.IsURLWithHTTPorHTTPS,
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

func (SREAgentResource) Attributes() map[string]*pluginsdk.Schema {
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

func (SREAgentResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			config := metadata.ResourceDiff.GetRawConfig()
			if config.IsNull() || !config.IsKnown() {
				return nil
			}

			identity := config.GetAttr("identity")
			if metadata.ResourceDiff.HasChanges("identity", "action", "knowledge_graph") &&
				identity.IsKnown() && !identity.IsNull() && identity.LengthInt() > 0 {
				identity = identity.AsValueSlice()[0]

				if identity.IsKnown() && !identity.IsNull() && identity.GetAttr("identity_ids").IsWhollyKnown() {
					ids := make([]string, 0)
					identityIds := identity.GetAttr("identity_ids")
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
				}

				subnet := config.GetAttr("subnet_id")
				if !subnet.IsKnown() {
					return nil
				}

				hasSubnet := !subnet.IsNull() && subnet.AsString() != ""
				if hasSubnet != (mode == string(agents.SandboxEgressModeAzureVNet)) {
					return errors.New("`subnet_id` and `sandbox.egress.mode` set to `AzureVNet` must be configured together")
				}
			}

			return nil
		},
	}
}

func (r SREAgentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SREAgents.AgentsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config SREAgentModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := agents.NewAgentID(subscriptionId, config.ResourceGroupName, config.Name)
			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			identity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(config.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			param := agents.Agent{
				Location: location.Normalize(config.Location),
				Identity: identity,
				Properties: &agents.AgentProperties{
					ActionConfiguration:             expandSREAgentAction(config.Action),
					LogConfiguration:                expandSREAgentApplicationInsights(config.ApplicationInsights),
					DefaultModel:                    expandSREAgentDefaultModel(config.DefaultModel),
					IncidentManagementConfiguration: expandSREAgentIncidentManagement(config.IncidentManagement),
					KnowledgeGraphConfiguration:     expandSREAgentKnowledgeGraph(config.KnowledgeGraph),
					SandboxConfiguration:            expandSREAgentSandbox(config.Sandbox),
				},
				Tags: pointer.To(config.Tags),
			}

			if config.InitialSponsorGroupId != "" {
				param.Properties.AgentIdentity = &agents.AgentIdentity{
					InitialSponsorGroupId: config.InitialSponsorGroupId,
				}
			}

			if config.SubnetId != "" {
				param.Properties.VnetConfiguration = &agents.VnetConfiguration{
					SubnetResourceId: pointer.To(config.SubnetId),
				}
			}

			if config.UpgradeChannel != "" {
				param.Properties.UpgradeChannel = pointer.ToEnum[agents.UpgradeChannel](config.UpgradeChannel)
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, param, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func (r SREAgentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SREAgents.AgentsClient

			id, err := agents.ParseAgentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (SREAgentResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SREAgents.AgentsClient

			id, err := agents.ParseAgentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config SREAgentModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			// Azure omits the connection string from GET responses, but PUT requires it alongside AppId.
			existing.Model.Properties.LogConfiguration = expandSREAgentApplicationInsights(config.ApplicationInsights)

			// Azure omits the connection key from GET responses, so resend the configured secret.
			existing.Model.Properties.IncidentManagementConfiguration = expandSREAgentIncidentManagement(config.IncidentManagement)

			if metadata.ResourceData.HasChange("identity") {
				identity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(config.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				existing.Model.Identity = identity
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = pointer.To(config.Tags)
			}

			if metadata.ResourceData.HasChange("action") {
				existing.Model.Properties.ActionConfiguration = expandSREAgentAction(config.Action)
			}

			if metadata.ResourceData.HasChange("default_model") {
				existing.Model.Properties.DefaultModel = expandSREAgentDefaultModel(config.DefaultModel)
			}

			if metadata.ResourceData.HasChange("initial_sponsor_group_id") {
				existing.Model.Properties.AgentIdentity = nil
				if config.InitialSponsorGroupId != "" {
					existing.Model.Properties.AgentIdentity = &agents.AgentIdentity{
						InitialSponsorGroupId: config.InitialSponsorGroupId,
					}
				}
			}

			if metadata.ResourceData.HasChange("knowledge_graph") {
				existing.Model.Properties.KnowledgeGraphConfiguration = expandSREAgentKnowledgeGraph(config.KnowledgeGraph)
			}

			if metadata.ResourceData.HasChange("sandbox") {
				existing.Model.Properties.SandboxConfiguration = expandSREAgentSandbox(config.Sandbox)
			}

			if metadata.ResourceData.HasChange("subnet_id") {
				existing.Model.Properties.VnetConfiguration = nil
				if config.SubnetId != "" {
					existing.Model.Properties.VnetConfiguration = &agents.VnetConfiguration{
						SubnetResourceId: pointer.To(config.SubnetId),
					}
				}
			}

			if metadata.ResourceData.HasChange("upgrade_channel") {
				existing.Model.Properties.UpgradeChannel = nil
				if config.UpgradeChannel != "" {
					existing.Model.Properties.UpgradeChannel = pointer.ToEnum[agents.UpgradeChannel](config.UpgradeChannel)
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (SREAgentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SREAgents.AgentsClient

			id, err := agents.ParseAgentID(metadata.ResourceData.Id())
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

func (SREAgentResource) flatten(metadata sdk.ResourceMetaData, id *agents.AgentId, input *agents.Agent) error {
	var config SREAgentModel
	if err := metadata.Decode(&config); err != nil {
		return fmt.Errorf("decoding: %+v", err)
	}

	state := SREAgentModel{
		Name:              id.AgentName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model := input; model != nil {
		state.Location = location.Normalize(model.Location)
		state.Tags = pointer.From(model.Tags)

		identity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		state.Identity = identity

		if props := model.Properties; props != nil {
			state.Endpoint = pointer.From(props.AgentEndpoint)
			state.OutboundIPAddresses = pointer.From(props.OutboundIPAddresses)
			state.UpgradeChannel = pointer.FromEnum(props.UpgradeChannel)
			state.ApplicationInsights = flattenSREAgentApplicationInsights(props.LogConfiguration, config.ApplicationInsights)
			state.DefaultModel = flattenSREAgentDefaultModel(props.DefaultModel)
			state.IncidentManagement = flattenSREAgentIncidentManagement(props.IncidentManagementConfiguration, config.IncidentManagement)
			state.Sandbox = flattenSREAgentSandbox(props.SandboxConfiguration)

			if props.AgentIdentity != nil {
				state.InitialSponsorGroupId = props.AgentIdentity.InitialSponsorGroupId
			}

			state.Action, err = flattenSREAgentAction(props.ActionConfiguration)
			if err != nil {
				return err
			}

			state.KnowledgeGraph, err = flattenSREAgentKnowledgeGraph(props.KnowledgeGraphConfiguration)
			if err != nil {
				return err
			}

			if props.VnetConfiguration != nil && pointer.From(props.VnetConfiguration.SubnetResourceId) != "" {
				parsed, err := commonids.ParseSubnetIDInsensitively(*props.VnetConfiguration.SubnetResourceId)
				if err != nil {
					return err
				}

				state.SubnetId = parsed.ID()
			}
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func expandSREAgentAction(input []ActionModel) *agents.ActionConfiguration {
	if len(input) == 0 {
		return nil
	}

	action := input[0]
	result := &agents.ActionConfiguration{}
	if action.AccessLevel != "" {
		result.AccessLevel = pointer.ToEnum[agents.AgentAccessLevel](action.AccessLevel)
	}

	if action.Mode != "" {
		result.Mode = pointer.ToEnum[agents.AgentMode](action.Mode)
	}

	if action.UserAssignedIdentityId != "" {
		result.Identity = pointer.To(action.UserAssignedIdentityId)
	}

	return result
}

func flattenSREAgentAction(input *agents.ActionConfiguration) ([]ActionModel, error) {
	if input == nil {
		return make([]ActionModel, 0), nil
	}

	action := ActionModel{
		AccessLevel: pointer.FromEnum(input.AccessLevel),
		Mode:        pointer.FromEnum(input.Mode),
	}
	if pointer.From(input.Identity) != "" {
		parsed, err := commonids.ParseUserAssignedIdentityIDInsensitively(*input.Identity)
		if err != nil {
			return nil, err
		}

		action.UserAssignedIdentityId = parsed.ID()
	}

	return []ActionModel{
		action,
	}, nil
}

func expandSREAgentApplicationInsights(input []ApplicationInsightsModel) *agents.LogConfiguration {
	if len(input) == 0 {
		return nil
	}

	applicationInsights := input[0]
	return &agents.LogConfiguration{
		ApplicationInsightsConfiguration: &agents.ApplicationInsightsConfiguration{
			AppId:            pointer.To(applicationInsights.AppId),
			ConnectionString: pointer.To(applicationInsights.ConnectionString),
		},
	}
}

func flattenSREAgentApplicationInsights(input *agents.LogConfiguration, config []ApplicationInsightsModel) []ApplicationInsightsModel {
	if input == nil || input.ApplicationInsightsConfiguration == nil {
		return make([]ApplicationInsightsModel, 0)
	}

	applicationInsights := ApplicationInsightsModel{
		AppId: pointer.From(input.ApplicationInsightsConfiguration.AppId),
	}

	// Azure does not return the create/update-only connection string.
	if len(config) > 0 {
		applicationInsights.ConnectionString = config[0].ConnectionString
	}

	return []ApplicationInsightsModel{
		applicationInsights,
	}
}

func expandSREAgentDefaultModel(input []DefaultModel) *agents.DefaultModel {
	if len(input) == 0 {
		return nil
	}

	defaultModel := input[0]
	return &agents.DefaultModel{
		Name:     pointer.To(defaultModel.Name),
		Provider: pointer.To(defaultModel.Provider),
	}
}

func flattenSREAgentDefaultModel(input *agents.DefaultModel) []DefaultModel {
	if input == nil {
		return make([]DefaultModel, 0)
	}

	return []DefaultModel{
		{
			Name:     pointer.From(input.Name),
			Provider: pointer.From(input.Provider),
		},
	}
}

func expandSREAgentIncidentManagement(input []IncidentManagementModel) *agents.IncidentManagementConfiguration {
	if len(input) == 0 {
		return nil
	}

	incidentManagement := input[0]
	result := &agents.IncidentManagementConfiguration{
		Type:          pointer.To(incidentManagement.Type),
		ConnectionKey: pointer.To(incidentManagement.ConnectionKey),
		ConnectionURL: pointer.To(incidentManagement.ConnectionURL),
		OboUser:       pointer.To(incidentManagement.OboUser),
	}
	if incidentManagement.ConnectionName != "" {
		result.ConnectionName = pointer.To(incidentManagement.ConnectionName)
	}

	return result
}

func flattenSREAgentIncidentManagement(input *agents.IncidentManagementConfiguration, config []IncidentManagementModel) []IncidentManagementModel {
	if input == nil {
		return make([]IncidentManagementModel, 0)
	}

	incidentManagement := IncidentManagementModel{
		ConnectionName: pointer.From(input.ConnectionName),
		ConnectionURL:  pointer.From(input.ConnectionURL),
		OboUser:        pointer.From(input.OboUser),
		Type:           pointer.From(input.Type),
	}

	// Azure masks the connection key on GET, so retain the configured secret.
	if len(config) > 0 {
		incidentManagement.ConnectionKey = config[0].ConnectionKey
	}

	return []IncidentManagementModel{
		incidentManagement,
	}
}

func expandSREAgentKnowledgeGraph(input []KnowledgeGraphModel) *agents.KnowledgeGraphConfiguration {
	if len(input) == 0 {
		return nil
	}

	knowledgeGraph := input[0]
	result := &agents.KnowledgeGraphConfiguration{
		ManagedResources: pointer.To(knowledgeGraph.ManagedResourceIds),
	}
	if knowledgeGraph.UserAssignedIdentityId != "" {
		result.Identity = pointer.To(knowledgeGraph.UserAssignedIdentityId)
	}

	return result
}

func flattenSREAgentKnowledgeGraph(input *agents.KnowledgeGraphConfiguration) ([]KnowledgeGraphModel, error) {
	if input == nil {
		return make([]KnowledgeGraphModel, 0), nil
	}

	knowledgeGraph := KnowledgeGraphModel{
		ManagedResourceIds: pointer.From(input.ManagedResources),
	}
	if pointer.From(input.Identity) != "" {
		parsed, err := commonids.ParseUserAssignedIdentityIDInsensitively(*input.Identity)
		if err != nil {
			return nil, err
		}

		knowledgeGraph.UserAssignedIdentityId = parsed.ID()
	}

	return []KnowledgeGraphModel{
		knowledgeGraph,
	}, nil
}

func expandSREAgentSandbox(input []SandboxModel) *agents.SandboxConfiguration {
	if len(input) == 0 {
		return nil
	}

	sandbox := input[0]
	packages := make([]agents.SandboxPackageSpec, 0, len(sandbox.Package))
	for _, item := range sandbox.Package {
		spec := agents.SandboxPackageSpec{
			Name:           item.Name,
			PackageManager: agents.SandboxPackageManager(item.PackageManager),
		}
		if item.Version != "" {
			spec.Version = pointer.To(item.Version)
		}

		packages = append(packages, spec)
	}

	result := &agents.SandboxConfiguration{
		Packages: pointer.To(packages),
	}
	if len(sandbox.Egress) > 0 {
		egress := sandbox.Egress[0]
		result.Egress = &agents.SandboxEgressConfiguration{
			AllowHTTPMcpServerNetworkAccess: pointer.To(egress.HTTPMCPServerNetworkAccessEnabled),
			AllowedCodeRepositories:         pointer.To(egress.AllowedCodeRepositories),
			AllowedHosts:                    pointer.To(egress.AllowedHosts),
			AllowedRegistries:               pointer.To(egress.AllowedRegistries),
		}
		if egress.Mode != "" {
			result.Egress.Mode = pointer.ToEnum[agents.SandboxEgressMode](egress.Mode)
		}

		if len(egress.HTTPProxy) > 0 {
			proxy := egress.HTTPProxy[0]
			result.Egress.HTTPProxy = &agents.SandboxHTTPProxyConfiguration{
				HTTPSURL: pointer.To(proxy.HTTPSURL),
				HTTPURL:  pointer.To(proxy.HTTPURL),
				NoProxy:  pointer.To(proxy.NoProxy),
			}
		}
	}

	return result
}

func flattenSREAgentSandbox(input *agents.SandboxConfiguration) []SandboxModel {
	if input == nil {
		return make([]SandboxModel, 0)
	}

	sandbox := SandboxModel{
		Package: make([]SandboxPackageModel, 0),
	}
	for _, item := range pointer.From(input.Packages) {
		sandbox.Package = append(sandbox.Package, SandboxPackageModel{
			Name:           item.Name,
			PackageManager: string(item.PackageManager),
			Version:        pointer.From(item.Version),
		})
	}

	if egress := input.Egress; egress != nil {
		sandbox.Egress = []SandboxEgressModel{
			{
				Mode:                              pointer.FromEnum(egress.Mode),
				HTTPMCPServerNetworkAccessEnabled: pointer.From(egress.AllowHTTPMcpServerNetworkAccess),
				AllowedCodeRepositories:           pointer.From(egress.AllowedCodeRepositories),
				AllowedHosts:                      pointer.From(egress.AllowedHosts),
				AllowedRegistries:                 pointer.From(egress.AllowedRegistries),
			},
		}
		if proxy := egress.HTTPProxy; proxy != nil {
			sandbox.Egress[0].HTTPProxy = []SandboxHTTPProxyModel{
				{
					HTTPSURL: pointer.From(proxy.HTTPSURL),
					HTTPURL:  pointer.From(proxy.HTTPURL),
					NoProxy:  pointer.From(proxy.NoProxy),
				},
			}
		}
	}

	return []SandboxModel{
		sandbox,
	}
}
