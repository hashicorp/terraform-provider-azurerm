// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Description: "Manages an API Management Service.",
		Create:      resourceApiManagementServiceCreate,
		Read:        resourceApiManagementServiceRead,
		Update:      resourceApiManagementServiceUpdate,
		Delete:      resourceApiManagementServiceDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := apimanagementservice.ParseServiceID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(3 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(3 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(3 * time.Hour),
		},

		Schema: resourceApiManagementSchema(),

		// we can only change `virtual_network_type` from None to Internal Or External, Else the subnet can not be destroyed cause “InUseSubnetCannotBeDeleted” for 3 hours
		// we can not change the subnet from subnet1 to subnet2 either, Else the subnet1 can not be destroyed cause “InUseSubnetCannotBeDeleted” for 3 hours
		// Issue: https://github.com/Azure/azure-rest-api-specs/issues/10395
		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("virtual_network_type", func(ctx context.Context, old, new, meta any) bool {
				return old.(string) != string(apimanagementservice.VirtualNetworkTypeNone) ||
					(new.(string) != string(apimanagementservice.VirtualNetworkTypeInternal) &&
						new.(string) != string(apimanagementservice.VirtualNetworkTypeExternal))
			}),

			pluginsdk.ForceNewIfChange("virtual_network_configuration", func(ctx context.Context, old, new, meta any) bool {
				return len(old.([]any)) != 0 || len(new.([]any)) == 0
			}),

			pluginsdk.ForceNewIfChange("sku_name", func(ctx context.Context, old, new, meta any) bool {
				return (strings.Contains(old.(string), "V2") && !strings.Contains(new.(string), "V2")) || (strings.Contains(new.(string), "V2") && !strings.Contains(old.(string), "V2"))
			}),
		),
	}

	if !features.SixPointOh() {
		resource.Schema["tenant_access"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeList,
			Optional:   true,
			Computed:   true, // azignore:AZS007 - pre-existing violation
			MaxItems:   1,
			Deprecated: "`tenant_access` has been deprecated and will be removed in v6.0 of the AzureRM provider as the direct management API of API Management has been retired.",
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Description: "Whether the access to the management API should be enabled.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},
					"tenant_id": {
						Description: "The identifier for the tenant access information contract.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"primary_key": {
						Description: "Primary access key for the tenant access information contract.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"secondary_key": {
						Description: "Secondary access key for the tenant access information contract.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
		}
	}

	return resource
}
