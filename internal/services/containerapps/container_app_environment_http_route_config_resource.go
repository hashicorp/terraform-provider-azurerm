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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ContainerAppEnvironmentHttpRouteConfigResource struct{}

type ContainerAppEnvironmentHttpRouteConfigModel struct {
	Name                      string                       `tfschema:"name"`
	ContainerAppEnvironmentId string                       `tfschema:"container_app_environment_id"`
	CustomDomains             []HttpRouteCustomDomainModel `tfschema:"custom_domains"`
	Rules                     []HttpRouteRuleModel         `tfschema:"rules"`
	Fqdn                      string                       `tfschema:"fqdn"`
}

type HttpRouteCustomDomainModel struct {
	BindingType   string `tfschema:"binding_type"`
	CertificateId string `tfschema:"certificate_id"`
	Name          string `tfschema:"name"`
}

type HttpRouteRuleModel struct {
	Description string                 `tfschema:"description"`
	Routes      []HttpRouteModel       `tfschema:"routes"`
	Targets     []HttpRouteTargetModel `tfschema:"targets"`
}

type HttpRouteModel struct {
	Action []HttpRouteActionModel `tfschema:"action"`
	Match  []HttpRouteMatchModel  `tfschema:"match"`
}

type HttpRouteActionModel struct {
	PrefixRewrite string `tfschema:"prefix_rewrite"`
}

type HttpRouteMatchModel struct {
	CaseSensitive       bool   `tfschema:"case_sensitive"`
	Path                string `tfschema:"path"`
	PathSeparatedPrefix string `tfschema:"path_separated_prefix"`
	Prefix              string `tfschema:"prefix"`
}

type HttpRouteTargetModel struct {
	ContainerApp string `tfschema:"container_app"`
	Label        string `tfschema:"label"`
	Revision     string `tfschema:"revision"`
}

var _ sdk.ResourceWithUpdate = ContainerAppEnvironmentHttpRouteConfigResource{}

func (r ContainerAppEnvironmentHttpRouteConfigResource) ModelObject() interface{} {
	return &ContainerAppEnvironmentHttpRouteConfigModel{}
}

func (r ContainerAppEnvironmentHttpRouteConfigResource) ResourceType() string {
	return "azurerm_container_app_environment_http_route_config"
}

func (r ContainerAppEnvironmentHttpRouteConfigResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return httprouteconfig.ValidateHTTPRouteConfigID
}

func (r ContainerAppEnvironmentHttpRouteConfigResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.HttpRouteConfigName,
			Description:  "The name for this HTTP Route Config.",
		},

		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: httprouteconfig.ValidateManagedEnvironmentID,
			Description:  "The ID of the Container App Environment.",
		},

		"custom_domains": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"binding_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(httprouteconfig.BindingTypeAuto),
							string(httprouteconfig.BindingTypeDisabled),
							string(httprouteconfig.BindingTypeSniEnabled),
						}, false),
						Description: "The Binding type. Possible values include `Auto`, `Disabled` and `SniEnabled`.",
					},

					"certificate_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						Description:  "The ID of the Certificate bound to this hostname.",
					},

					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						Description:  "The hostname.",
					},
				},
			},
		},

		"rules": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"description": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						Description:  "Description of the rule.",
					},

					"routes": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"action": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"prefix_rewrite": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
												Description:  "Rewrite prefix. Default is no rewrites.",
											},
										},
									},
								},

								"match": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"case_sensitive": {
												Type:        pluginsdk.TypeBool,
												Optional:    true,
												Default:     true,
												Description: "Path case sensitive. Defaults to `true`.",
											},

											"path": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
												Description:  "Match on exact path.",
											},

											"path_separated_prefix": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
												Description:  "Match on path separated prefix.",
											},

											"prefix": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
												Description:  "Match on prefix.",
											},
										},
									},
								},
							},
						},
					},

					"targets": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"container_app": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validate.ContainerAppName,
									Description:  "Container App Name to route requests to.",
								},

								"label": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
									Description:  "Label to route requests to.",
								},

								"revision": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
									Description:  "Revision to route requests to.",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r ContainerAppEnvironmentHttpRouteConfigResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"fqdn": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The FQDN of the HTTP Route Config.",
		},
	}
}

func (r ContainerAppEnvironmentHttpRouteConfigResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.HttpRouteConfigClient

			var model ContainerAppEnvironmentHttpRouteConfigModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			envId, err := httprouteconfig.ParseManagedEnvironmentID(model.ContainerAppEnvironmentId)
			if err != nil {
				return err
			}

			id := httprouteconfig.NewHTTPRouteConfigID(metadata.Client.Account.SubscriptionId, envId.ResourceGroupName, envId.ManagedEnvironmentName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			httpRouteConfig := httprouteconfig.HTTPRouteConfig{
				Properties: &httprouteconfig.HTTPRouteConfigProperties{
					CustomDomains: expandHttpRouteCustomDomains(model.CustomDomains),
					Rules:         expandHttpRouteRules(model.Rules),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, httpRouteConfig); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ContainerAppEnvironmentHttpRouteConfigResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.HttpRouteConfigClient

			id, err := httprouteconfig.ParseHTTPRouteConfigID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			var state ContainerAppEnvironmentHttpRouteConfigModel

			state.Name = id.HttpRouteConfigName
			state.ContainerAppEnvironmentId = httprouteconfig.NewManagedEnvironmentID(id.SubscriptionId, id.ResourceGroupName, id.ManagedEnvironmentName).ID()

			if model := existing.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Fqdn = pointer.From(props.Fqdn)
					state.CustomDomains = flattenHttpRouteCustomDomains(props.CustomDomains)
					state.Rules = flattenHttpRouteRules(props.Rules)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ContainerAppEnvironmentHttpRouteConfigResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.HttpRouteConfigClient

			id, err := httprouteconfig.ParseHTTPRouteConfigID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ContainerAppEnvironmentHttpRouteConfigResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.HttpRouteConfigClient

			id, err := httprouteconfig.ParseHTTPRouteConfigID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state ContainerAppEnvironmentHttpRouteConfigModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil || existing.Model == nil || existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s for update: %+v", *id, err)
			}

			if metadata.ResourceData.HasChange("custom_domains") {
				existing.Model.Properties.CustomDomains = expandHttpRouteCustomDomains(state.CustomDomains)
			}

			if metadata.ResourceData.HasChange("rules") {
				existing.Model.Properties.Rules = expandHttpRouteRules(state.Rules)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandHttpRouteCustomDomains(input []HttpRouteCustomDomainModel) *[]httprouteconfig.CustomDomain {
	if len(input) == 0 {
		return nil
	}

	result := make([]httprouteconfig.CustomDomain, 0)
	for _, v := range input {
		cd := httprouteconfig.CustomDomain{
			Name: v.Name,
		}
		if v.BindingType != "" {
			bindingType := httprouteconfig.BindingType(v.BindingType)
			cd.BindingType = &bindingType
		}
		if v.CertificateId != "" {
			cd.CertificateId = pointer.To(v.CertificateId)
		}
		result = append(result, cd)
	}

	return &result
}

func flattenHttpRouteCustomDomains(input *[]httprouteconfig.CustomDomain) []HttpRouteCustomDomainModel {
	if input == nil {
		return []HttpRouteCustomDomainModel{}
	}

	result := make([]HttpRouteCustomDomainModel, 0)
	for _, v := range *input {
		cd := HttpRouteCustomDomainModel{
			Name: v.Name,
		}
		if v.BindingType != nil {
			cd.BindingType = string(*v.BindingType)
		}
		if v.CertificateId != nil {
			cd.CertificateId = *v.CertificateId
		}
		result = append(result, cd)
	}

	return result
}

func expandHttpRouteRules(input []HttpRouteRuleModel) *[]httprouteconfig.HTTPRouteRule {
	if len(input) == 0 {
		return nil
	}

	result := make([]httprouteconfig.HTTPRouteRule, 0)
	for _, v := range input {
		rule := httprouteconfig.HTTPRouteRule{
			Routes:  expandHttpRoutes(v.Routes),
			Targets: expandHttpRouteTargets(v.Targets),
		}
		if v.Description != "" {
			rule.Description = pointer.To(v.Description)
		}
		result = append(result, rule)
	}

	return &result
}

func flattenHttpRouteRules(input *[]httprouteconfig.HTTPRouteRule) []HttpRouteRuleModel {
	if input == nil {
		return []HttpRouteRuleModel{}
	}

	result := make([]HttpRouteRuleModel, 0)
	for _, v := range *input {
		rule := HttpRouteRuleModel{
			Description: pointer.From(v.Description),
			Routes:      flattenHttpRoutes(v.Routes),
			Targets:     flattenHttpRouteTargets(v.Targets),
		}
		result = append(result, rule)
	}

	return result
}

func expandHttpRoutes(input []HttpRouteModel) *[]httprouteconfig.HTTPRoute {
	if len(input) == 0 {
		return nil
	}

	result := make([]httprouteconfig.HTTPRoute, 0)
	for _, v := range input {
		route := httprouteconfig.HTTPRoute{}
		if len(v.Action) > 0 {
			route.Action = &httprouteconfig.HTTPRouteAction{}
			if v.Action[0].PrefixRewrite != "" {
				route.Action.PrefixRewrite = pointer.To(v.Action[0].PrefixRewrite)
			}
		}
		if len(v.Match) > 0 {
			m := v.Match[0]
			route.Match = &httprouteconfig.HTTPRouteMatch{
				CaseSensitive: pointer.To(m.CaseSensitive),
			}
			if m.Path != "" {
				route.Match.Path = pointer.To(m.Path)
			}
			if m.PathSeparatedPrefix != "" {
				route.Match.PathSeparatedPrefix = pointer.To(m.PathSeparatedPrefix)
			}
			if m.Prefix != "" {
				route.Match.Prefix = pointer.To(m.Prefix)
			}
		}
		result = append(result, route)
	}

	return &result
}

func flattenHttpRoutes(input *[]httprouteconfig.HTTPRoute) []HttpRouteModel {
	if input == nil {
		return []HttpRouteModel{}
	}

	result := make([]HttpRouteModel, 0)
	for _, v := range *input {
		route := HttpRouteModel{}
		if v.Action != nil {
			route.Action = []HttpRouteActionModel{
				{
					PrefixRewrite: pointer.From(v.Action.PrefixRewrite),
				},
			}
		}
		if v.Match != nil {
			route.Match = []HttpRouteMatchModel{
				{
					CaseSensitive:       pointer.From(v.Match.CaseSensitive),
					Path:                pointer.From(v.Match.Path),
					PathSeparatedPrefix: pointer.From(v.Match.PathSeparatedPrefix),
					Prefix:              pointer.From(v.Match.Prefix),
				},
			}
		}
		result = append(result, route)
	}

	return result
}

func expandHttpRouteTargets(input []HttpRouteTargetModel) *[]httprouteconfig.HTTPRouteTarget {
	if len(input) == 0 {
		return nil
	}

	result := make([]httprouteconfig.HTTPRouteTarget, 0)
	for _, v := range input {
		target := httprouteconfig.HTTPRouteTarget{
			ContainerApp: v.ContainerApp,
		}
		if v.Label != "" {
			target.Label = pointer.To(v.Label)
		}
		if v.Revision != "" {
			target.Revision = pointer.To(v.Revision)
		}
		result = append(result, target)
	}

	return &result
}

func flattenHttpRouteTargets(input *[]httprouteconfig.HTTPRouteTarget) []HttpRouteTargetModel {
	if input == nil {
		return []HttpRouteTargetModel{}
	}

	result := make([]HttpRouteTargetModel, 0)
	for _, v := range *input {
		target := HttpRouteTargetModel{
			ContainerApp: v.ContainerApp,
			Label:        pointer.From(v.Label),
			Revision:     pointer.From(v.Revision),
		}
		result = append(result, target)
	}

	return result
}
