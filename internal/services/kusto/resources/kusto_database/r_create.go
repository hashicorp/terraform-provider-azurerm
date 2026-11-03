// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_database

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/databases"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoDatabaseCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.DatabasesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewKustoDatabaseID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cluster_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
	}

	databaseProperties := expandKustoDatabaseProperties(d)

	readWriteDatabase := databases.ReadWriteDatabase{
		Location:   pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: databaseProperties,
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, readWriteDatabase, databases.DefaultCreateOrUpdateOperationOptions(), sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, readWriteDatabase, databases.DefaultCreateOrUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceKustoDatabaseRead(d, meta)
}

func expandKustoDatabaseProperties(d *pluginsdk.ResourceData) *databases.ReadWriteDatabaseProperties {
	databaseProperties := &databases.ReadWriteDatabaseProperties{}

	if softDeletePeriod, ok := d.GetOk("soft_delete_period"); ok {
		databaseProperties.SoftDeletePeriod = pointer.To(softDeletePeriod.(string))
	}

	if hotCachePeriod, ok := d.GetOk("hot_cache_period"); ok {
		databaseProperties.HotCachePeriod = pointer.To(hotCachePeriod.(string))
	}

	return databaseProperties
}
