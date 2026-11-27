// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_api

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_api/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "resource_group_name,service_name:api_management_name" -compare-values "api_id:id"

const azureApiManagementApiResourceName = "azurerm_api_management_api"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages an API within an API Management Service.",
		Create:      resourceApiManagementApiCreate,
		Read:        resourceApiManagementApiRead,
		Update:      resourceApiManagementApiUpdate,
		Delete:      resourceApiManagementApiDelete,
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&api.ApiId{}),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&api.ApiId{}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ApiV0ToV1{},
		}),

		Schema: apiManagementApiSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
				values := d.GetRawConfig().AsValueMap()
				if d.Get("version").(string) != "" && values["version_set_id"].IsNull() {
					return errors.New("setting `version` without the required `version_set_id`")
				}

				protocols := expandApiManagementApiProtocols(d.Get("protocols").(*pluginsdk.Set).List())
				if values["source_api_id"].IsNull() && (values["display_name"].IsNull() || protocols == nil || len(*protocols) == 0) {
					return errors.New("`display_name`, `protocols` are required when `source_api_id` is not set")
				}

				if d.Get("api_type").(string) == string(api.ApiTypeWebsocket) && d.Get("service_url").(string) == "" {
					return errors.New("`service_url` is required when `api_type` is `websocket`")
				}
				return nil
			}),
		),
	}
}
