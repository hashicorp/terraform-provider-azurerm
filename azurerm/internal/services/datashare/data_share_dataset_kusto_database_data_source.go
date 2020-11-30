package datashare

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/helper"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceDataShareDatasetKustoDatabase() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDataShareDatasetKustoDatabaseRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.DataSetName(),
			},

			"share_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ShareID,
			},

			"kusto_database_id": {
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceDataShareDatasetKustoDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	shareID := d.Get("share_id").(string)
	shareId, err := parse.ShareID(shareID)
	if err != nil {
		return err
	}

	respModel, err := client.Get(ctx, shareId.ResourceGroup, shareId.AccountName, shareId.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving DataShare Kusto Database DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name, err)
	}

	respId := helper.GetAzurermDataShareDataSetId(respModel.Value)
	if respId == nil || *respId == "" {
		return fmt.Errorf("empty or nil ID returned for DataShare Kusto Database DataSet %q (Resource Group %q / accountName %q / shareName %q)", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name)
	}

	d.SetId(*respId)
	d.Set("name", name)
	d.Set("share_id", shareID)

	resp, ok := respModel.Value.AsKustoDatabaseDataSet()
	if !ok {
		return fmt.Errorf("dataShare %q (Resource Group %q / accountName %q / shareName %q) is not kusto database dataset", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name)
	}
	if props := resp.KustoDatabaseDataSetProperties; props != nil {
		d.Set("kusto_database_id", props.KustoDatabaseResourceID)
		d.Set("display_name", props.DataSetID)
		d.Set("kusto_cluster_location", props.Location)
	}

	return nil
}
