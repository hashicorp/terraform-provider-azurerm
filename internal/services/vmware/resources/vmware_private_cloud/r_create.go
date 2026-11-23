// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vmware_private_cloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/privateclouds"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVmwarePrivateCloudCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Vmware.PrivateCloudClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := privateclouds.NewPrivateCloudID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_vmware_private_cloud", id.ID())
		}
	}

	internet := privateclouds.InternetEnumDisabled
	if d.Get("internet_connection_enabled").(bool) {
		internet = privateclouds.InternetEnumEnabled
	}

	privateCloud := privateclouds.PrivateCloud{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Sku: privateclouds.Sku{
			Name: d.Get("sku_name").(string),
		},
		Properties: &privateclouds.PrivateCloudProperties{
			ManagementCluster: &privateclouds.CommonClusterProperties{
				ClusterSize: pointer.To(int64(d.Get("management_cluster.0.size").(int))),
			},
			NetworkBlock:    d.Get("network_subnet_cidr").(string),
			Internet:        &internet,
			NsxtPassword:    pointer.To(d.Get("nsxt_password").(string)),
			VcenterPassword: pointer.To(d.Get("vcenter_password").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if _, err := client.CreateOrUpdate(ctx, id, privateCloud); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// TODO:
	// @sreallymatt: this can likely be removed in favour of using polling methods from the SDK, however,
	// we have no capacity to test this resource at the moment so leaving it be for now.
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(privateclouds.PrivateCloudProvisioningStateBuilding)},
		Target:     []string{string(privateclouds.PrivateCloudProvisioningStateSucceeded)},
		Refresh:    privateCloudStateRefreshFunc(ctx, client, id),
		MinTimeout: 15 * time.Second,
		Timeout:    time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for vmware private cloud %s provisioning state to become available error: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVmwarePrivateCloudRead(d, meta)
}
