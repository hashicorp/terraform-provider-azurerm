// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datashare

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/dataset"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/share"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name data_share_dataset_kusto_database -service-package-name datashare -properties "name" -compare-values "resource_group_name:share_id,account_name:share_id,share_name:share_id" -known-values "subscription_id:data.Subscriptions.Primary"

func resourceDataShareDataSetKustoDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataShareDataSetKustoDatabaseCreate,
		Read:   resourceDataShareDataSetKustoDatabaseRead,
		Delete: resourceDataShareDataSetKustoDatabaseDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&dataset.DataSetId{}),
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&dataset.DataSetId{}),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataSetName(),
			},

			"share_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: share.ValidateShareID,
			},

			"kusto_database_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateKustoDatabaseID,
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

func resourceDataShareDataSetKustoDatabaseCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	shareId, err := share.ParseShareID(d.Get("share_id").(string))
	if err != nil {
		return err
	}
	id := dataset.NewDataSetID(shareId.SubscriptionId, shareId.ResourceGroupName, shareId.AccountName, shareId.ShareName, d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for present of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_data_share_dataset_kusto_database", id.ID())
	}

	dataSet := dataset.KustoDatabaseDataSet{
		Properties: dataset.KustoDatabaseDataSetProperties{
			KustoDatabaseResourceId: d.Get("kusto_database_id").(string),
		},
	}

	if _, err := client.Create(ctx, id, dataSet); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceDataShareDataSetKustoDatabaseRead(d, meta)
}

func resourceDataShareDataSetKustoDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataset.ParseDataSetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] DataShare %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.DataSetName)

	shareId := share.NewShareID(id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ShareName)
	d.Set("share_id", shareId.ID())

	if model := resp.Model; model != nil {
		if ds, ok := model.(dataset.KustoDatabaseDataSet); ok {
			props := ds.Properties
			d.Set("kusto_database_id", props.KustoDatabaseResourceId)
			d.Set("display_name", props.DataSetId)
			d.Set("kusto_cluster_location", props.Location)
		} else {
			return fmt.Errorf("%s is not kusto database dataset", *id)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func resourceDataShareDataSetKustoDatabaseDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataset.ParseDataSetID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
