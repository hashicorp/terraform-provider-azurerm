// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceStorageAccountQueueProperties() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Read: dataSourceStorageAccountQueuePropertiesRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"primary_queue_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_queue_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_queue_microsoft_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_queue_microsoft_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_queue_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_queue_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_queue_microsoft_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_queue_microsoft_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}

	return resource
}

func dataSourceStorageAccountQueuePropertiesRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ResourceManager.StorageAccounts
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	resp, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	var primaryEndpoints *storageaccounts.Endpoints
	var secondaryEndpoints *storageaccounts.Endpoints
	var routingPreference *storageaccounts.RoutingPreference

	if model := resp.Model; model != nil && model.Properties != nil {
		primaryEndpoints = model.Properties.PrimaryEndpoints
		routingPreference = model.Properties.RoutingPreference
		secondaryEndpoints = model.Properties.SecondaryEndpoints
	}

	queueEndpoints := flattenAccountQueueEndpoints(primaryEndpoints, secondaryEndpoints, routingPreference)
	if err := queueEndpoints.set(d); err != nil {
		return err
	}

	return nil
}
