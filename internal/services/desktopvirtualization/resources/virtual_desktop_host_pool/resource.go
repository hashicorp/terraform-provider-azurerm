// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_host_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/resources/virtual_desktop_host_pool/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualDesktopHostPoolCreate,
		Read:   resourceVirtualDesktopHostPoolRead,
		Update: resourceVirtualDesktopHostPoolUpdate,
		Delete: resourceVirtualDesktopHostPoolDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := hostpool.ParseHostPoolID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.HostPoolV0ToV1{},
		}),
		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, meta any) error {
			if d.Get("type").(string) == string(hostpool.HostPoolTypePooled) {
				loadBalancerType := d.Get("load_balancer_type").(string)
				if loadBalancerType == string(hostpool.LoadBalancerTypePersistent) || loadBalancerType == string(hostpool.LoadBalancerTypeMultiplePersistent) {
					return fmt.Errorf("`type` must be \"Personal\" when `load_balancer_type` is \"Persistent\" or \"MultiplePersistent\"")
				}
			}
			return nil
		}),

		Schema: virtualDesktopHostPoolSchema(),
	}
}
