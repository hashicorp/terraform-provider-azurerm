package vmware

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/avs/mgmt/2020-03-20/avs"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/vmware/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/vmware/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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
			_, err := parse.ClusterID(id)
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
	privateCloudId, err := parse.PrivateCloudID(d.Get("vmware_cloud_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewClusterID(subscriptionId, privateCloudId.ResourceGroup, privateCloudId.Name, name)

	existing, err := client.Get(ctx, id.ResourceGroup, id.PrivateCloudName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing %q : %+v", id.ID(), err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_vmware_cluster", id.ID())
	}

	cluster := avs.Cluster{
		Sku: &avs.Sku{
			Name: utils.String(d.Get("sku_name").(string)),
		},
		ClusterProperties: &avs.ClusterProperties{
			ClusterSize: utils.Int32(int32(d.Get("cluster_node_count").(int))),
		},
	}
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.PrivateCloudName, id.Name, cluster)
	if err != nil {
		return fmt.Errorf("creating Vmware  Cluster %q (Resource Group %q / privateCloudName %q): %+v", id.Name, id.ResourceGroup, id.PrivateCloudName, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of the Vmware  Cluster %q (Resource Group %q / privateCloudName %q): %+v", id.Name, id.ResourceGroup, id.PrivateCloudName, err)
	}

	d.SetId(id.ID())
	return resourceVmwareClusterRead(d, meta)
}

func resourceVmwareClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Vmware.ClusterClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.PrivateCloudName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Vmware Cluster %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Vmware  Cluster %q (Resource Group %q / privateCloudName %q): %+v", id.Name, id.ResourceGroup, id.PrivateCloudName, err)
	}

	d.Set("name", id.Name)
	d.Set("vmware_cloud_id", parse.NewPrivateCloudID(subscriptionId, id.ResourceGroup, id.PrivateCloudName).ID())
	d.Set("sku_name", resp.Sku.Name)
	if props := resp.ClusterProperties; props != nil {
		d.Set("cluster_node_count", props.ClusterSize)
		d.Set("cluster_number", props.ClusterID)
		d.Set("hosts", utils.FlattenStringSlice(props.Hosts))
	}
	return nil
}

func resourceVmwareClusterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Vmware.ClusterClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	clusterUpdate := avs.ClusterUpdate{
		ClusterUpdateProperties: &avs.ClusterUpdateProperties{},
	}
	if d.HasChange("cluster_node_count") {
		clusterUpdate.ClusterUpdateProperties.ClusterSize = utils.Int32(int32(d.Get("cluster_node_count").(int)))
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.PrivateCloudName, id.Name, clusterUpdate)
	if err != nil {
		return fmt.Errorf("updating Vmware  Cluster %q (Resource Group %q / privateCloudName %q): %+v", id.Name, id.ResourceGroup, id.PrivateCloudName, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of the Vmware  Cluster %q (Resource Group %q / privateCloudName %q): %+v", id.Name, id.ResourceGroup, id.PrivateCloudName, err)
	}
	return resourceVmwareClusterRead(d, meta)
}

func resourceVmwareClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Vmware.ClusterClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.PrivateCloudName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Vmware  Cluster %q (Resource Group %q / privateCloudName %q): %+v", id.Name, id.ResourceGroup, id.PrivateCloudName, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of the Vmware  Cluster %q (Resource Group %q / privateCloudName %q): %+v", id.Name, id.ResourceGroup, id.PrivateCloudName, err)
	}
	return nil
}
