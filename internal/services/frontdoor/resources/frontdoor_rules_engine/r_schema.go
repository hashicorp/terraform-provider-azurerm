// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package frontdoor_rules_engine

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func frontdoorRulesEngineSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the Rules engine configuration.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"frontdoor_name": {
			Description:  "The name of the Front Door instance.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FrontDoorName,
		},
		"location": commonschema.LocationComputed(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"enabled": {
			Description: "Whether this Rules engine configuration is enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"rule": {
			Type:     pluginsdk.TypeList,
			MaxItems: 100,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the rule.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"priority": {
						Description: "Priority of the rule, must be unique per rules engine definition.",
						Type:        pluginsdk.TypeInt,
						Required:    true,
					},

					"match_condition": {
						Type:     pluginsdk.TypeList,
						MaxItems: 100,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"variable": {
									Description: "The condition variable.",
									Type:        pluginsdk.TypeString,
									Optional:    true,
									ValidateFunc: validation.StringInSlice([]string{
										"IsMobile",
										"RemoteAddr",
										"RequestMethod",
										"QueryString",
										"PostArgs",
										"RequestUri",
										"RequestPath",
										"RequestFilename",
										"RequestFilenameExtension",
										"RequestHeader",
										"RequestBody",
										"RequestScheme",
									}, false),
								},

								"selector": {
									Description:  "Match against a specific key when `variable` is set to `PostArgs` or `RequestHeader`. It cannot be used with `QueryString` and `RequestMethod`.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"operator": {
									Description: "The condition operator.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ValidateFunc: validation.StringInSlice([]string{
										"Any",
										"IPMatch",
										"GeoMatch",
										"Equal",
										"Contains",
										"LessThan",
										"GreaterThan",
										"LessThanOrEqual",
										"GreaterThanOrEqual",
										"BeginsWith",
										"EndsWith",
									}, false),
								},

								"transform": {
									Description: "The condition transform.",
									Type:        pluginsdk.TypeList,
									Optional:    true,
									MaxItems:    6,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											"Lowercase",
											"RemoveNulls",
											"Trim",
											"Uppercase",
											"UrlDecode",
											"UrlEncode",
										}, false),
									},
								},

								"negate_condition": {
									Description: "Whether to negate the condition.",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Default:     false,
								},

								"value": {
									Description: "A list of values.",
									Type:        pluginsdk.TypeList,
									Optional:    true,
									MaxItems:    25,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},

					"action": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"request_header": {
									Type:     pluginsdk.TypeList,
									MaxItems: 100,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"header_action_type": {
												Description:  "The action to take on the header.",
												Type:         pluginsdk.TypeString,
												ValidateFunc: validation.StringInSlice(frontdoors.PossibleValuesForHeaderActionType(), false),
												Optional:     true,
											},

											"header_name": {
												Description:  "The header name.",
												Type:         pluginsdk.TypeString,
												ValidateFunc: validation.StringIsNotEmpty,
												Optional:     true,
											},

											"value": {
												Description:  "The header value.",
												Type:         pluginsdk.TypeString,
												ValidateFunc: validation.StringIsNotEmpty,
												Optional:     true,
											},
										},
									},
								},

								"response_header": {
									Type:     pluginsdk.TypeList,
									MaxItems: 100,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"header_action_type": {
												Description:  "The action to take on the header.",
												Type:         pluginsdk.TypeString,
												ValidateFunc: validation.StringInSlice(frontdoors.PossibleValuesForHeaderActionType(), false),
												Optional:     true,
											},

											"header_name": {
												Description:  "The header name.",
												Type:         pluginsdk.TypeString,
												ValidateFunc: validation.StringIsNotEmpty,
												Optional:     true,
											},

											"value": {
												Description:  "The header value.",
												Type:         pluginsdk.TypeString,
												ValidateFunc: validation.StringIsNotEmpty,
												Optional:     true,
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
	}
}
