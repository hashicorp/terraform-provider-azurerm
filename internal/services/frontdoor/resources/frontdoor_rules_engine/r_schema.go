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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"frontdoor_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FrontDoorName,
		},
		"location": commonschema.LocationComputed(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"rule": {
			Type:     pluginsdk.TypeList,
			MaxItems: 100,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"priority": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"match_condition": {
						Type:     pluginsdk.TypeList,
						MaxItems: 100,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"variable": {
									Type:     pluginsdk.TypeString,
									Optional: true,
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
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"operator": {
									Type:     pluginsdk.TypeString,
									Required: true,
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
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 6,
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
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"value": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 25,
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
												Type:         pluginsdk.TypeString,
												ValidateFunc: validation.StringInSlice(frontdoors.PossibleValuesForHeaderActionType(), false),
												Optional:     true,
											},

											"header_name": {
												Type:         pluginsdk.TypeString,
												ValidateFunc: validation.StringIsNotEmpty,
												Optional:     true,
											},

											"value": {
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
												Type:         pluginsdk.TypeString,
												ValidateFunc: validation.StringInSlice(frontdoors.PossibleValuesForHeaderActionType(), false),
												Optional:     true,
											},

											"header_name": {
												Type:         pluginsdk.TypeString,
												ValidateFunc: validation.StringIsNotEmpty,
												Optional:     true,
											},

											"value": {
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
