// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package availability_set

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/availabilitysets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAvailabilitySetCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.AvailabilitySetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewAvailabilitySetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_availability_set", id.ID())
			}
		}
	}

	updateDomainCount := d.Get("platform_update_domain_count").(int)
	faultDomainCount := d.Get("platform_fault_domain_count").(int)
	managed := d.Get("managed").(bool)
	t := d.Get("tags").(map[string]any)

	payload := availabilitysets.AvailabilitySet{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &availabilitysets.AvailabilitySetProperties{
			PlatformFaultDomainCount:  pointer.To(int64(faultDomainCount)),
			PlatformUpdateDomainCount: pointer.To(int64(updateDomainCount)),
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("proximity_placement_group_id"); ok {
		payload.Properties.ProximityPlacementGroup = &availabilitysets.SubResource{
			Id: pointer.To(v.(string)),
		}
	}

	if managed {
		payload.Sku = &availabilitysets.Sku{
			Name: pointer.To("Aligned"),
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceAvailabilitySetRead(d, meta)
}
