// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_share_dataset_kusto_database

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/dataset"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/share"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	datashareClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Data Share Kusto Database Dataset.",
		Read:        dataSourceDataShareDatasetKustoDatabaseRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of this Data Share Kusto Database Dataset.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.DataSetName(),
			},

			"share_id": {
				Description:  "The resource ID of the Data Share where this Data Share Kusto Database Dataset should be created.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: share.ValidateShareID,
			},

			"kusto_database_id": {
				Description: "The resource ID of the Kusto Cluster Database to be shared with the receiver.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"display_name": {
				Description: "The name of the Data Share Dataset.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"kusto_cluster_location": {
				Description: "The location of the Kusto Cluster.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceDataShareDatasetKustoDatabaseRead(d *pluginsdk.ResourceData, meta any) error {
	client := datashareClient.From(meta).DataSetClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	shareId, err := share.ParseShareID(d.Get("share_id").(string))
	if err != nil {
		return err
	}
	id := dataset.NewDataSetID(subscriptionId, shareId.ResourceGroupName, shareId.AccountName, shareId.ShareName, d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.DataSetName)
	d.Set("share_id", shareId.ID())

	if model := resp.Model; model != nil {
		if ds, ok := model.(dataset.KustoDatabaseDataSet); ok {
			props := ds.Properties
			d.Set("kusto_database_id", props.KustoDatabaseResourceId)
			d.Set("display_name", props.DataSetId)
			d.Set("kusto_cluster_location", props.Location)
		} else {
			return fmt.Errorf("%s is not kusto database dataset", id)
		}
	}

	return nil
}
