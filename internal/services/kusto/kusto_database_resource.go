// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/databases"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/migration"
	kustoValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKustoDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoDatabaseCreateUpdate,
		Read:   resourceKustoDatabaseRead,
		Update: resourceKustoDatabaseCreateUpdate,
		Delete: resourceKustoDatabaseDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KustoDatabaseV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseKustoDatabaseID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: kustoValidate.DatabaseName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"cluster_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: kustoValidate.ClusterName,
			},

			"soft_delete_period": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.ISO8601Duration,
			},

			"hot_cache_period": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.ISO8601Duration,
			},

			"size": {
				Type:     pluginsdk.TypeFloat,
				Computed: true,
			},
		},
	}
}

func resourceKustoDatabaseCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DatabasesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewKustoDatabaseID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cluster_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_kusto_database", id.ID())
		}
	}

	databaseProperties := expandKustoDatabaseProperties(d)

	readWriteDatabase := databases.ReadWriteDatabase{
		Location:   utils.String(location.Normalize(d.Get("location").(string))),
		Properties: databaseProperties,
	}

	err := client.CreateOrUpdateThenPoll(ctx, id, readWriteDatabase, databases.DefaultCreateOrUpdateOperationOptions())
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceKustoDatabaseRead(d, meta)
}

func resourceKustoDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseKustoDatabaseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: response was nil", *id)
	}

	database, ok := resp.Model.(databases.ReadWriteDatabase)
	if !ok {
		return fmt.Errorf("%s was not a Read/Write Database", *id)
	}

	d.Set("name", id.KustoDatabaseName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("cluster_name", id.KustoClusterName)
	d.Set("location", location.NormalizeNilable(database.Location))

	if props := database.Properties; props != nil {
		d.Set("hot_cache_period", props.HotCachePeriod)
		d.Set("soft_delete_period", props.SoftDeletePeriod)

		if statistics := props.Statistics; statistics != nil {
			d.Set("size", statistics.Size)
		}
	}

	return nil
}

func resourceKustoDatabaseDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DatabasesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseKustoDatabaseID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandKustoDatabaseProperties(d *pluginsdk.ResourceData) *databases.ReadWriteDatabaseProperties {
	databaseProperties := &databases.ReadWriteDatabaseProperties{}

	if softDeletePeriod, ok := d.GetOk("soft_delete_period"); ok {
		databaseProperties.SoftDeletePeriod = utils.String(softDeletePeriod.(string))
	}

	if hotCachePeriod, ok := d.GetOk("hot_cache_period"); ok {
		databaseProperties.HotCachePeriod = utils.String(hotCachePeriod.(string))
	}

	return databaseProperties
}
