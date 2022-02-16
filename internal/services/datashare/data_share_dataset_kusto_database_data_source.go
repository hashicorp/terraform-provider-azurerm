package datashare

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceDataShareDatasetKustoDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDataShareDatasetKustoDatabaseRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.DataSetName(),
			},

			"share_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ShareID,
			},

			"kusto_database_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"display_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"kusto_cluster_location": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDataShareDatasetKustoDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	shareId, err := parse.ShareID(d.Get("share_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewDataSetID(subscriptionId, shareId.ResourceGroup, shareId.AccountName, shareId.Name, d.Get("name").(string))

	respModel, err := client.Get(ctx, id.ResourceGroup, id.AccountName, id.ShareName, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	respId := helper.GetAzurermDataShareDataSetId(respModel.Value)
	if respId == nil || *respId == "" {
		return fmt.Errorf("empty or nil ID returned for %s", id)
	}

	d.SetId(id.ID())
	d.Set("name", id.Name)
	d.Set("share_id", shareId.ID())

	resp, ok := respModel.Value.AsKustoDatabaseDataSet()
	if !ok {
		return fmt.Errorf("%s is not kusto database dataset", id)
	}
	if props := resp.KustoDatabaseDataSetProperties; props != nil {
		d.Set("kusto_database_id", props.KustoDatabaseResourceID)
		d.Set("display_name", props.DataSetID)
		d.Set("kusto_cluster_location", props.Location)
	}

	return nil
}
