// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vmware_private_cloud

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/privateclouds"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVmwarePrivateCloudUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Vmware.PrivateCloudClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privateclouds.ParsePrivateCloudID(d.Id())
	if err != nil {
		return err
	}

	privateCloudUpdate := privateclouds.PrivateCloudUpdate{
		Properties: &privateclouds.PrivateCloudUpdateProperties{},
	}

	if d.HasChange("management_cluster") && d.HasChange("internet_connection_enabled") {
		return fmt.Errorf("`management_cluster.0.size` and `internet_connection_enabled` could not be changed together")
	}

	if d.HasChange("management_cluster") {
		privateCloudUpdate.Properties.ManagementCluster = &privateclouds.CommonClusterProperties{
			ClusterSize: pointer.To(int64(d.Get("management_cluster.0.size").(int))),
		}
	}

	if d.HasChange("internet_connection_enabled") {
		internet := privateclouds.InternetEnumDisabled
		if d.Get("internet_connection_enabled").(bool) {
			internet = privateclouds.InternetEnumEnabled
		}
		privateCloudUpdate.Properties.Internet = &internet
	}

	if d.HasChange("tags") {
		privateCloudUpdate.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.UpdateThenPoll(ctx, *id, privateCloudUpdate); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceVmwarePrivateCloudRead(d, meta)
}
