// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package static_web_app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.StaticWebAppName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"configuration_file_changes_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"preview_environments_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"sku_tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(resourceproviders.SkuNameFree),
			ValidateFunc: validation.StringInSlice([]string{
				string(resourceproviders.SkuNameStandard),
				string(resourceproviders.SkuNameFree),
			}, false),
		},

		"sku_size": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(resourceproviders.SkuNameFree),
			ValidateFunc: validation.StringInSlice([]string{
				string(resourceproviders.SkuNameStandard),
				string(resourceproviders.SkuNameFree),
			}, false),
		},

		"app_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"basic_auth": BasicAuthSchema(),

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"repository_url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			RequiredWith: []string{"repository_token", "repository_branch"},
		},

		"repository_token": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"repository_url", "repository_branch"},
		},

		"repository_branch": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"repository_url", "repository_token"},
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"default_host_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			skuTier := rd.Get("sku_tier").(string)
			skuSize := rd.Get("sku_size").(string)

			if strings.EqualFold(skuTier, string(resourceproviders.SkuNameFree)) && strings.EqualFold(skuSize, string(resourceproviders.SkuNameFree)) {
				basicAuth, authOk := rd.GetOk("basic_auth")
				if authOk && len(basicAuth.([]any)) > 0 {
					return fmt.Errorf("basic_auth cannot be used with the Free tier of Static Web Apps")
				}
				ident, identOk := rd.GetOk("identity")
				if identOk && len(ident.([]any)) > 0 {
					return fmt.Errorf("identities cannot be used with the Free tier of Static Web Apps")
				}
			}

			return nil
		},
	}
}
