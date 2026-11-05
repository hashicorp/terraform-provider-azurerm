// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_standard_web_test

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			// SSLCheck conditions
			url, ok := rd.GetOk("request.0.url")
			if ok {
				if !strings.HasPrefix(strings.ToLower(url.(string)), "https://") {
					if v, ok := rd.GetOkExists("validation_rules.0.ssl_check_enabled"); ok && v.(bool) {
						return fmt.Errorf("cannot set ssl_check_enabled to true if request.0.url is not https")
					}
					if v, ok := rd.GetOkExists("validation_rules.0.ssl_cert_remaining_lifetime"); ok && v.(int) != 0 {
						return fmt.Errorf("cannot set ssl_cert_remaining_lifetime if request.0.url is not https")
					}
				}
			}

			return nil
		},
	}
}

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"application_insights_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: componentsapis.ValidateComponentID,
		},

		"location": commonschema.Location(),

		"frequency": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  300,
			ValidateFunc: validation.IntInSlice([]int{
				300,
				600,
				900,
			}),
		},

		"timeout": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  30,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"retry_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"request": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"follow_redirects_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"http_verb": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "GET",
						ValidateFunc: validation.StringInSlice([]string{
							"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
						}, false),
					},

					"parse_dependent_requests_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"header": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"value": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},

					"body": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"url": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"validation_rules": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"expected_status_code": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  200,
					},

					// Typo in API spec, issue: https://github.com/Azure/azure-rest-api-specs/issues/22136
					// "ignore_status_code": {
					// 	Type:   pluginsdk.TypeBool,
					// 	Optional: true,
					// 	Default: false,
					// },

					"ssl_cert_remaining_lifetime": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 365),
					},

					"ssl_check_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"content": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"content_match": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"ignore_case": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
								"pass_if_text_found": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
							},
						},
					},
				},
			},
		},

		"geo_locations": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:             pluginsdk.TypeString,
				ValidateFunc:     validation.StringIsNotEmpty,
				StateFunc:        location.StateFunc,
				DiffSuppressFunc: location.DiffSuppressFunc,
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"tags": commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"synthetic_monitor_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
