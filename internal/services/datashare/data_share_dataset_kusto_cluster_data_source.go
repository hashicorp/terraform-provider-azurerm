// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datashare

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/dataset"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/share"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceDataShareDatasetKustoCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDataShareDatasetKustoClusterRead,

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
				ValidateFunc: share.ValidateShareID,
			},

			"kusto_cluster_id": {
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

func dataSourceDataShareDatasetKustoClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	shareId, err := share.ParseShareID(d.Get("share_id").(string))
	if err != nil {
		return err
	}
	id := dataset.NewDataSetID(shareId.SubscriptionId, shareId.ResourceGroupName, shareId.AccountName, shareId.ShareName, d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.DataSetName)
	d.Set("share_id", shareId.ID())

	if model := resp.Model; model != nil {
		if ds, ok := model.(dataset.KustoClusterDataSet); ok {
			props := ds.Properties
			// TODO parse kusto id
			d.Set("kusto_cluster_id", props.KustoClusterResourceId)
			d.Set("display_name", props.DataSetId)
			d.Set("kusto_cluster_location", props.Location)
		} else {
			return fmt.Errorf("%s is not kusto cluster dataset", id)
		}
	}

	return nil
}
