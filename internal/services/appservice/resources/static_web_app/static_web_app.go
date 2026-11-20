// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package static_web_app

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type BasicAuth struct {
	Password     string `tfschema:"password"`
	Environments string `tfschema:"environments"`
}

type BasicAuthComputed struct {
	Environments string `tfschema:"environments"`
}

const (
	EnvironmentsTypeAllEnvironments       string = "AllEnvironments"
	EnvironmentsTypeStagingEnvironments   string = "StagingEnvironments"
	EnvironmentsTypeSpecifiedEnvironments string = "SpecifiedEnvironments"
)

func BasicAuthSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:      pluginsdk.TypeList,
		MaxItems:  1,
		Optional:  true,
		Sensitive: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*schema.Schema{
				"password": {
					Description:  "The password for the basic authentication access.",
					Type:         pluginsdk.TypeString,
					Required:     true,
					Sensitive:    true,
					ValidateFunc: validate.StaticWebAppPassword,
				},

				"environments": {
					Description: "The Environment types to use the Basic Auth for access.",
					Type:        pluginsdk.TypeString,
					Required:    true,
					ValidateFunc: validation.StringInSlice([]string{
						EnvironmentsTypeAllEnvironments,
						EnvironmentsTypeStagingEnvironments,
					}, false),
				},
			},
		},
	}
}

func BasicAuthSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:      pluginsdk.TypeList,
		Computed:  true,
		Sensitive: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*schema.Schema{
				"environments": {
					Description: "The Environment types which are configured to use Basic Auth access.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},
			},
		},
	}
}
