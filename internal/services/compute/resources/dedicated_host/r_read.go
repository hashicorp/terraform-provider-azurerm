// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dedicated_host

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/dedicatedhosts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDedicatedHostRead(d *pluginsdk.ResourceData, meta any) error {
	hostsClient := meta.(*clients.Client).Compute.DedicatedHostsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseDedicatedHostID(d.Id())
	if err != nil {
		return err
	}

	resp, err := hostsClient.Get(ctx, *id, dedicatedhosts.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.HostName)
	d.Set("dedicated_host_group_id", commonids.NewDedicatedHostGroupID(id.SubscriptionId, id.ResourceGroupName, id.HostGroupName).ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("sku_name", model.Sku.Name)
		if props := model.Properties; props != nil {
			d.Set("auto_replace_on_failure", props.AutoReplaceOnFailure)

			licenseType := pointer.FromEnum(props.LicenseType)
			if licenseType == string(dedicatedhosts.DedicatedHostLicenseTypesNone) {
				licenseType = ""
			}
			d.Set("license_type", licenseType)

			platformFaultDomain := 0
			if props.PlatformFaultDomain != nil {
				platformFaultDomain = int(*props.PlatformFaultDomain)
			}
			d.Set("platform_fault_domain", platformFaultDomain)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
