// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containerapps

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/httprouteconfig"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ContainerAppEnvironmentHttpRouteConfigDataSource struct{}

type ContainerAppEnvironmentHttpRouteConfigDataSourceModel struct {
	Name                      string                       `tfschema:"name"`
	ContainerAppEnvironmentId string                       `tfschema:"container_app_environment_id"`
	CustomDomains             []HttpRouteCustomDomainModel `tfschema:"custom_domains"`
	Rules                     []HttpRouteRuleModel         `tfschema:"rules"`
	Fqdn                      string                       `tfschema:"fqdn"`
}

var _ sdk.DataSource = ContainerAppEnvironmentHttpRouteConfigDataSource{}

func (r ContainerAppEnvironmentHttpRouteConfigDataSource) ModelObject() interface{} {
	return &ContainerAppEnvironmentHttpRouteConfigDataSourceModel{}
}

func (r ContainerAppEnvironmentHttpRouteConfigDataSource) ResourceType() string {
	return "azurerm_container_app_environment_http_route_config"
}

func (r ContainerAppEnvironmentHttpRouteConfigDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.HttpRouteConfigName,
			Description:  "The name of the HTTP Route Config.",
		},

		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: httprouteconfig.ValidateManagedEnvironmentID,
			Description:  "The ID of the Container App Environment.",
		},
	}
}

func (r ContainerAppEnvironmentHttpRouteConfigDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"custom_domains": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"binding_type": {
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Description: "The Binding type.",
					},

					"certificate_id": {
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Description: "The ID of the Certificate bound to this hostname.",
					},

					"name": {
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Description: "The hostname.",
					},
				},
			},
		},

		"rules": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"description": {
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Description: "Description of the rule.",
					},

					"routes": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"action": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"prefix_rewrite": {
												Type:        pluginsdk.TypeString,
												Computed:    true,
												Description: "Rewrite prefix.",
											},
										},
									},
								},

								"match": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"case_sensitive": {
												Type:        pluginsdk.TypeBool,
												Computed:    true,
												Description: "Path case sensitive.",
											},

											"path": {
												Type:        pluginsdk.TypeString,
												Computed:    true,
												Description: "Match on exact path.",
											},

											"path_separated_prefix": {
												Type:        pluginsdk.TypeString,
												Computed:    true,
												Description: "Match on path separated prefix.",
											},

											"prefix": {
												Type:        pluginsdk.TypeString,
												Computed:    true,
												Description: "Match on prefix.",
											},
										},
									},
								},
							},
						},
					},

					"targets": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"container_app": {
									Type:        pluginsdk.TypeString,
									Computed:    true,
									Description: "Container App Name to route requests to.",
								},

								"label": {
									Type:        pluginsdk.TypeString,
									Computed:    true,
									Description: "Label to route requests to.",
								},

								"revision": {
									Type:        pluginsdk.TypeString,
									Computed:    true,
									Description: "Revision to route requests to.",
								},
							},
						},
					},
				},
			},
		},

		"fqdn": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The FQDN of the HTTP Route Config.",
		},
	}
}

func (r ContainerAppEnvironmentHttpRouteConfigDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.HttpRouteConfigClient

			var state ContainerAppEnvironmentHttpRouteConfigDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			envId, err := httprouteconfig.ParseManagedEnvironmentID(state.ContainerAppEnvironmentId)
			if err != nil {
				return err
			}

			id := httprouteconfig.NewHTTPRouteConfigID(envId.SubscriptionId, envId.ResourceGroupName, envId.ManagedEnvironmentName, state.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			state.Name = id.HttpRouteConfigName
			state.ContainerAppEnvironmentId = envId.ID()

			if model := existing.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Fqdn = pointer.From(props.Fqdn)
					state.CustomDomains = flattenHttpRouteCustomDomains(props.CustomDomains)
					state.Rules = flattenHttpRouteRules(props.Rules)
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
