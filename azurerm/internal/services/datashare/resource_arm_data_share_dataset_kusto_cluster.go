package datashare

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datashare/mgmt/2019-11-01/datashare"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/helper"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataShareDataSetKustoCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDataShareDataSetKustoClusterCreate,
		Read:   resourceArmDataShareDataSetKustoClusterRead,
		Delete: resourceArmDataShareDataSetKustoClusterDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DataShareDataSetID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DatashareDataSetName(),
			},

			"share_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataShareID,
			},

			"kusto_cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"kusto_cluster_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceArmDataShareDataSetKustoClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	shareId, err := parse.DataShareID(d.Get("share_id").(string))
	if err != nil {
		return err
	}

	existingModel, err := client.Get(ctx, shareId.ResourceGroup, shareId.AccountName, shareId.Name, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existingModel.Response) {
			return fmt.Errorf("checking for presence of existing  DataShare Kusto Cluster DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name, err)
		}
	}
	existingId := helper.GetAzurermDataShareDataSetId(existingModel.Value)
	if existingId != nil && *existingId != "" {
		return tf.ImportAsExistsError("azurerm_data_share_dataset_kusto_cluster", *existingId)
	}

	dataSet := datashare.KustoClusterDataSet{
		Kind: datashare.KindKustoCluster,
		KustoClusterDataSetProperties: &datashare.KustoClusterDataSetProperties{
			KustoClusterResourceID: utils.String(d.Get("kusto_cluster_id").(string)),
		},
	}

	if _, err := client.Create(ctx, shareId.ResourceGroup, shareId.AccountName, shareId.Name, name, dataSet); err != nil {
		return fmt.Errorf("creating DataShare Kusto Cluster DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name, err)
	}

	respModel, err := client.Get(ctx, shareId.ResourceGroup, shareId.AccountName, shareId.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving DataShare Kusto Cluster DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name, err)
	}

	respId := helper.GetAzurermDataShareDataSetId(respModel.Value)
	if respId == nil || *respId == "" {
		return fmt.Errorf("empty or nil ID returned for DataShare Kusto Cluster DataSet %q (Resource Group %q / accountName %q / shareName %q)", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name)
	}

	d.SetId(*respId)
	return resourceArmDataShareDataSetKustoClusterRead(d, meta)
}

func resourceArmDataShareDataSetKustoClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	shareClient := meta.(*clients.Client).DataShare.SharesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataShareDataSetID(d.Id())
	if err != nil {
		return err
	}

	respModel, err := client.Get(ctx, id.ResourceGroup, id.AccountName, id.ShareName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(respModel.Response) {
			log.Printf("[INFO] DataShare %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving DataShare Kusto Cluster DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", id.Name, id.ResourceGroup, id.AccountName, id.ShareName, err)
	}

	d.Set("name", id.Name)
	shareResp, err := shareClient.Get(ctx, id.ResourceGroup, id.AccountName, id.ShareName)
	if err != nil {
		return fmt.Errorf("retrieving DataShare %q (Resource Group %q / accountName %q): %+v", id.ShareName, id.ResourceGroup, id.AccountName, err)
	}
	if shareResp.ID == nil || *shareResp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for DataShare %q (Resource Group %q / accountName %q)", id.ShareName, id.ResourceGroup, id.AccountName)
	}

	d.Set("share_id", shareResp.ID)

	resp, ok := respModel.Value.AsKustoClusterDataSet()
	if !ok {
		return fmt.Errorf("dataShare dataset %q (Resource Group %q / accountName %q / shareName %q) is not kusto cluster dataset", id.Name, id.ResourceGroup, id.AccountName, id.ShareName)
	}
	if props := resp.KustoClusterDataSetProperties; props != nil {
		d.Set("kusto_cluster_id", props.KustoClusterResourceID)
		d.Set("display_name", props.DataSetID)
		d.Set("kusto_cluster_location", props.Location)
	}

	return nil
}

func resourceArmDataShareDataSetKustoClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataShareDataSetID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.AccountName, id.ShareName, id.Name); err != nil {
		return fmt.Errorf("deleting DataShare Kusto Cluster DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", id.Name, id.ResourceGroup, id.AccountName, id.ShareName, err)
	}
	return nil
}
