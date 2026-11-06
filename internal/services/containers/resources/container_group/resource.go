// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerinstance/2025-09-01/containerinstance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceContainerGroupCreate,
		Read:   resourceContainerGroupRead,
		Delete: resourceContainerGroupDelete,
		Update: resourceContainerGroupUpdate,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := containerinstance.ParseContainerGroupID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: containerGroupSchema(),
		CustomizeDiff: func(ctx context.Context, d *pluginsdk.ResourceDiff, i any) error {
			if p := d.Get("priority").(string); p == string(containerinstance.ContainerGroupPrioritySpot) {
				if d.Get("ip_address_type").(string) != "None" {
					return fmt.Errorf("`ip_address_type` has to be `None` when `priority` is set to `Spot`")
				}
			}
			return nil
		},
	}
}
