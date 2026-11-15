// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb

import (
	"context"
	"reflect"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmLoadBalancerCreate,
		Read:   resourceArmLoadBalancerRead,
		Update: resourceArmLoadBalancerUpdate,
		Delete: resourceArmLoadBalancerDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := loadbalancers.ParseLoadBalancerID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: lbSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIf("frontend_ip_configuration", func(ctx context.Context, d *schema.ResourceDiff, meta any) bool {
				old, new := d.GetChange("frontend_ip_configuration")
				switch {
				case len(old.([]any)) == 0 && len(new.([]any)) > 0:
					return false
				case len(old.([]any)) > 0 && len(new.([]any)) == 0:
					// Azure does not allow an LB to have all frontend removed, it results in the following error:
					// Error: "Deleting all frontendIPConfigs from load balancer is not supported."
					// If the old config had frontends, and new config has none, need to force new LB
					return true
				default:
					for i, nc := range new.([]any) {
						dataNew := nc.(map[string]any)
						for _, oc := range old.([]any) {
							dataOld := oc.(map[string]any)
							if dataOld["name"].(string) == dataNew["name"].(string) {
								if !reflect.DeepEqual(dataOld["zones"].(*pluginsdk.Set).List(), dataNew["zones"].(*pluginsdk.Set).List()) {
									// set ForceNew to true when the `frontend_ip_configuration.#.zones` is changed.
									d.ForceNew("frontend_ip_configuration." + strconv.Itoa(i) + ".zones")
									break
								}
							}
						}
					}
				}
				return false
			}),
		),
	}
}
