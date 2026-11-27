// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package proximity_placement_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	computeClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceProximityPlacementGroupCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := computeClient.From(meta).ProximityPlacementGroupsClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	id := proximityplacementgroups.NewProximityPlacementGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id, proximityplacementgroups.DefaultGetOperationOptions())
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_proximity_placement_group", id.ID())
			}
		}
	}

	payload := proximityplacementgroups.ProximityPlacementGroup{
		Location:   location.Normalize(d.Get("location").(string)),
		Properties: &proximityplacementgroups.ProximityPlacementGroupProperties{},
		Tags:       tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("allowed_vm_sizes"); ok {
		if payload.Properties.Intent == nil {
			payload.Properties.Intent = &proximityplacementgroups.ProximityPlacementGroupPropertiesIntent{}
		}
		payload.Properties.Intent.VMSizes = pluginsdk.ExpandStringSlice(v.(*pluginsdk.Set).List())
	}

	if v, ok := d.GetOk("zone"); ok {
		zones := zones.Expand([]string{v.(string)})
		payload.Zones = &zones
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
		if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
			return err
		}
	}

	return resourceProximityPlacementGroupRead(d, meta)
}
