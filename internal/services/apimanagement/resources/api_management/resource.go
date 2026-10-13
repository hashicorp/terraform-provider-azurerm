// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementServiceCreate,
		Read:   resourceApiManagementServiceRead,
		Update: resourceApiManagementServiceUpdate,
		Delete: resourceApiManagementServiceDelete,
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
}
