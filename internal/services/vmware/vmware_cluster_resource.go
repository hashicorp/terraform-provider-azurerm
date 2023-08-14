// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vmware

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/privateclouds"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/vmware/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVmwareCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVmwareClusterCreate,
		Read:   resourceVmwareClusterRead,
		Update: resourceVmwareClusterUpdate,
		Delete: resourceVmwareClusterDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(5 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(5 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(5 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := clusters.ParseClusterID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"vmware_cloud_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.PrivateCloudID,
			},

			"cluster_node_count": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(3, 16),
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"av20",
					"av36",
					"av36t",
					"av36p",
					"av52",
				}, false),
			},

			"cluster_number": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"hosts": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func resourceVmwareClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
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
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_vmware_cluster", id.ID())
	}

	cluster := clusters.Cluster{
		Sku: clusters.Sku{
			Name: d.Get("sku_name").(string),
		},
		Properties: &clusters.CommonClusterProperties{
			ClusterSize: pointer.To(int64(d.Get("cluster_node_count").(int))),
		},
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, cluster); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVmwareClusterRead(d, meta)
}

func resourceVmwareClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Vmware.ClusterClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ClusterName)
	d.Set("vmware_cloud_id", privateclouds.NewPrivateCloudID(id.SubscriptionId, id.ResourceGroupName, id.PrivateCloudName).ID())

	if model := resp.Model; model != nil {
		d.Set("cluster_node_count", model.Properties.ClusterSize)
		d.Set("cluster_number", model.Properties.ClusterId)
		d.Set("hosts", utils.FlattenStringSlice(model.Properties.Hosts))
		d.Set("sku_name", model.Sku.Name)
	}

	return nil
}

func resourceVmwareClusterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Vmware.ClusterClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	clusterUpdate := clusters.ClusterUpdate{
		Properties: &clusters.ClusterUpdateProperties{},
	}
	if d.HasChange("cluster_node_count") {
		clusterUpdate.Properties.ClusterSize = utils.Int64(int64(d.Get("cluster_node_count").(int)))
	}

	if err := client.UpdateThenPoll(ctx, *id, clusterUpdate); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}
	return resourceVmwareClusterRead(d, meta)
}

func resourceVmwareClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Vmware.ClusterClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
