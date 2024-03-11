// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

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
					Type:         pluginsdk.TypeString,
					Required:     true,
					Sensitive:    true,
					ValidateFunc: validate.StaticWebAppPassword,
				},

				"environments": {
					Type:     pluginsdk.TypeString,
					Required: true,
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
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}
