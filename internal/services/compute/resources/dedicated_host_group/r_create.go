// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dedicated_host_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/dedicatedhostgroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDedicatedHostGroupCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.DedicatedHostGroupsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewDedicatedHostGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, dedicatedhostgroups.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_dedicated_host_group", id.ID())
		}
	}

	platformFaultDomainCount := d.Get("platform_fault_domain_count").(int)
	t := d.Get("tags").(map[string]any)

	payload := dedicatedhostgroups.DedicatedHostGroup{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &dedicatedhostgroups.DedicatedHostGroupProperties{
			PlatformFaultDomainCount: int64(platformFaultDomainCount),
		},
		Tags: tags.Expand(t),
	}

	if zone, ok := d.GetOk("zone"); ok {
		payload.Zones = &[]string{
			zone.(string),
		}
	}

	if v, ok := d.GetOk("automatic_placement_enabled"); ok {
		payload.Properties.SupportAutomaticPlacement = pointer.To(v.(bool))
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceDedicatedHostGroupRead(d, meta)
}
