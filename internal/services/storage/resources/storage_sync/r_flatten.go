// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_sync

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/registeredserverresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/storagesyncservicesresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceStorageSyncFlatten(ctx context.Context, d *pluginsdk.ResourceData, id *storagesyncservicesresource.StorageSyncServiceId, model *storagesyncservicesresource.StorageSyncService, registeredServerClient *registeredserverresource.RegisteredServerResourceClient, includeRegisteredServers bool) error {
	d.Set("name", id.StorageSyncServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("incoming_traffic_policy", pointer.FromEnum(props.IncomingTrafficPolicy))
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	if includeRegisteredServers {
		storageSyncId := registeredserverresource.NewStorageSyncServiceID(id.SubscriptionId, id.ResourceGroupName, id.StorageSyncServiceName)
		registeredServersResp, err := registeredServerClient.RegisteredServersListByStorageSyncService(ctx, storageSyncId)
		if err != nil {
			if !response.WasNotFound(registeredServersResp.HttpResponse) {
				return fmt.Errorf("retrieving registered servers for %s: %+v", id, err)
			}
		}

		if serverModel := registeredServersResp.Model; serverModel != nil && serverModel.Value != nil {
			registeredServers := make([]any, 0, len(*serverModel.Value))
			for _, registeredServer := range *serverModel.Value {
				if registeredServer.Id != nil {
					registeredServers = append(registeredServers, *registeredServer.Id)
				}
			}
			if err := d.Set("registered_servers", registeredServers); err != nil {
				return fmt.Errorf("setting `registered_servers`: %+v", err)
			}
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
