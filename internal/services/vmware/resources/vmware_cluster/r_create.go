// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vmware_cluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/privateclouds"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVmwareClusterCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Vmware.ClusterClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	privateCloudId, err := privateclouds.ParsePrivateCloudID(d.Get("vmware_cloud_id").(string))
	if err != nil {
		return err
	}

	id := clusters.NewClusterID(subscriptionId, privateCloudId.ResourceGroupName, privateCloudId.PrivateCloudName, name)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_vmware_cluster", id.ID())
		}
	}

	cluster := clusters.Cluster{
		Sku: clusters.Sku{
			Name: d.Get("sku_name").(string),
		},
		Properties: &clusters.CommonClusterProperties{
			ClusterSize: pointer.To(int64(d.Get("cluster_node_count").(int))),
		},
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, cluster, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVmwareClusterRead(d, meta)
}
