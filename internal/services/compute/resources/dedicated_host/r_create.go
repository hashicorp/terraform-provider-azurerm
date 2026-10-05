// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dedicated_host

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/dedicatedhosts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDedicatedHostCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.DedicatedHostsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	hostGroupId, err := commonids.ParseDedicatedHostGroupID(d.Get("dedicated_host_group_id").(string))
	if err != nil {
		return err
	}

	id := commonids.NewDedicatedHostID(hostGroupId.SubscriptionId, hostGroupId.ResourceGroupName, hostGroupId.HostGroupName, d.Get("name").(string))
	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, dedicatedhosts.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_dedicated_host", id.ID())
		}
	}

	payload := dedicatedhosts.DedicatedHost{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &dedicatedhosts.DedicatedHostProperties{
			AutoReplaceOnFailure: pointer.To(d.Get("auto_replace_on_failure").(bool)),
			LicenseType:          pointer.To(dedicatedhosts.DedicatedHostLicenseTypesNone),
			PlatformFaultDomain:  pointer.To(int64(d.Get("platform_fault_domain").(int))),
		},
		Sku: dedicatedhosts.Sku{
			Name: pointer.To(d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v := d.Get("license_type").(string); v != "" {
		payload.Properties.LicenseType = pointer.ToEnum[dedicatedhosts.DedicatedHostLicenseTypes](v)
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceDedicatedHostRead(d, meta)
}
