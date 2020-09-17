package avs

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/avs/mgmt/2020-03-20/avs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/avs/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"log"
	"time"
)

func resourceArmAvsCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAvsClusterCreate,
		Read:   resourceArmAvsClusterRead,
		Update: resourceArmAvsClusterUpdate,
		Delete: resourceArmAvsClusterDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(240 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(240 * time.Minute),
			Delete: schema.DefaultTimeout(240 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.AvsClusterID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"private_cloud_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cluster_size": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"sku_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"av20",
					"av36",
					"av36t",
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"cluster_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"hosts": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}
func resourceArmAvsClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Avs.ClusterClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	privateCloudId, _ := parse.AvsPrivateCloudID(d.Get("private_cloud_id").(string))

	existing, err := client.Get(ctx, privateCloudId.ResourceGroup, privateCloudId.Name, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing Avs Cluster %q (Resource Group %q / privateCloudName %q): %+v", name, privateCloudId.ResourceGroup, privateCloudId.Name, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_avs_cluster", *existing.ID)
	}

	cluster := avs.Cluster{
		Sku: &avs.Sku{
			Name: utils.String(d.Get("sku_name").(string)),
		},
		ClusterProperties: &avs.ClusterProperties{
			ClusterSize: utils.Int32(int32(d.Get("cluster_size").(int))),
		},
	}
	future, err := client.CreateOrUpdate(ctx, privateCloudId.ResourceGroup, privateCloudId.Name, name, cluster)
	if err != nil {
		return fmt.Errorf("creating Avs Cluster %q (Resource Group %q / privateCloudName %q): %+v", name, privateCloudId.ResourceGroup, privateCloudId.Name, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating future for Avs Cluster %q (Resource Group %q / privateCloudName %q): %+v", name, privateCloudId.ResourceGroup, privateCloudId.Name, err)
	}

	resp, err := client.Get(ctx, privateCloudId.ResourceGroup, privateCloudId.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Avs Cluster %q (Resource Group %q / privateCloudName %q): %+v", name, privateCloudId.ResourceGroup, privateCloudId.Name, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Avs Cluster %q (Resource Group %q / privateCloudName %q) ID", name, privateCloudId.ResourceGroup, privateCloudId.Name)
	}

	d.SetId(*resp.ID)
	return resourceArmAvsClusterRead(d, meta)
}

func resourceArmAvsClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Avs.ClusterClient
	pcClient := meta.(*clients.Client).Avs.PrivateCloudClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AvsClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.PrivateCloudName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] avs %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Avs Cluster %q (Resource Group %q / privateCloudName %q): %+v", id.Name, id.ResourceGroup, id.PrivateCloudName, err)
	}

	pcResp, err := pcClient.Get(ctx, id.ResourceGroup, id.PrivateCloudName)
	if err != nil {
		return fmt.Errorf("retrieving Avs PrivateCloud %q (Resource Group %q): %+v", id.PrivateCloudName, id.ResourceGroup, err)
	}

	if pcResp.ID == nil || *pcResp.ID == "" {
		return fmt.Errorf("avs PrivateCloud %q (Resource Group %q) ID is empty or nil", id.PrivateCloudName, id.ResourceGroup)
	}
	d.Set("name", id.Name)
	d.Set("private_cloud_id", pcResp.ID)
	d.Set("sku_name", resp.Sku.Name)
	if props := resp.ClusterProperties; props != nil {
		d.Set("cluster_size", props.ClusterSize)
		d.Set("cluster_id", props.ClusterID)
		d.Set("hosts", utils.FlattenStringSlice(props.Hosts))
	}
	return nil
}

func resourceArmAvsClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Avs.ClusterClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AvsClusterID(d.Id())
	if err != nil {
		return err
	}

	clusterUpdate := avs.ClusterUpdate{
		ClusterUpdateProperties: &avs.ClusterUpdateProperties{},
	}
	if d.HasChange("cluster_size") {
		clusterUpdate.ClusterUpdateProperties.ClusterSize = utils.Int32(int32(d.Get("cluster_size").(int)))
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.PrivateCloudName, id.Name, clusterUpdate)
	if err != nil {
		return fmt.Errorf("updating Avs Cluster %q (Resource Group %q / privateCloudName %q): %+v", id.Name, id.ResourceGroup, id.PrivateCloudName, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on updating future for Avs Cluster %q (Resource Group %q / privateCloudName %q): %+v", id.Name, id.ResourceGroup, id.PrivateCloudName, err)
	}
	return resourceArmAvsClusterRead(d, meta)
}

func resourceArmAvsClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Avs.ClusterClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AvsClusterID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.PrivateCloudName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Avs Cluster %q (Resource Group %q / privateCloudName %q): %+v", id.Name, id.ResourceGroup, id.PrivateCloudName, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on deleting future for Avs Cluster %q (Resource Group %q / privateCloudName %q): %+v", id.Name, id.ResourceGroup, id.PrivateCloudName, err)
	}
	return nil
}
