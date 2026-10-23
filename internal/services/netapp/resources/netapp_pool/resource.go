// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/capacitypools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetAppPoolCreate,
		Read:   resourceNetAppPoolRead,
		Update: resourceNetAppPoolUpdate,
		Delete: resourceNetAppPoolDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := capacitypools.ParseCapacityPoolID(id)
			return err
		}),

		Schema: netappPoolSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			// `cool_access_enabled` cannot be disabled
			pluginsdk.ForceNewIfChange("cool_access_enabled", func(ctx context.Context, old, new, meta any) bool {
				return old.(bool) && !new.(bool)
			}),
			// custom_throughput_mibps validation
			func(ctx context.Context, d *pluginsdk.ResourceDiff, i any) error {
				customThroughput := d.Get("custom_throughput_mibps").(int)
				qosType := d.Get("qos_type").(string)
				serviceLevel := d.Get("service_level").(string)

				if customThroughput > 0 {
					// customThroughputMibps should only be accepted if manual qosType is set
					if qosType != string(capacitypools.QosTypeManual) {
						return fmt.Errorf("`custom_throughput_mibps` can only be set when `qos_type` is set to `Manual`")
					}

					// customThroughputMibps should only be accepted for Flexible service level
					if serviceLevel != string(capacitypools.ServiceLevelFlexible) {
						return fmt.Errorf("`custom_throughput_mibps` can only be set when `service_level` is set to `Flexible`")
					}
				}

				return nil
			},
		),
	}
}
