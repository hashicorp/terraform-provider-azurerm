// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type Category struct {
	Feeds      []string `tfschema:"feeds"`
	CustomUrls []string `tfschema:"custom_urls"`
}

func CategorySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"feeds": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MinItems: 1,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"custom_urls": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MinItems: 1,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validate.CategoryNames,
					},
				},
			},
		},
	}
}

func ExpandCategory(input []Category) *localrules.Category {
	if len(input) == 0 {
		return nil
	}

	c := input[0]

	return &localrules.Category{
		Feeds:     c.Feeds,
		UrlCustom: c.CustomUrls,
	}
}

func FlattenCategory(input *localrules.Category) []Category {
	if input == nil {
		return []Category{}
	}

	return []Category{{
		Feeds:      input.Feeds,
		CustomUrls: input.UrlCustom,
	}}
}
