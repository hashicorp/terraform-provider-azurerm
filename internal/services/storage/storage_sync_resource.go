// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/registeredserverresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/storagesyncservicesresource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

const storageSyncResourceName = "azurerm_storage_sync"

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name storage_sync -service-package-name storage -properties "name,resource_group_name" -known-values "subscription_id:data.Subscriptions.Primary"

func resourceStorageSync() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageSyncCreate,
		Read:   resourceStorageSyncRead,
		Update: resourceStorageSyncUpdate,
		Delete: resourceStorageSyncDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&storagesyncservicesresource.StorageSyncServiceId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&storagesyncservicesresource.StorageSyncServiceId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageSyncName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"incoming_traffic_policy": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(storagesyncservicesresource.IncomingTrafficPolicyAllowAllTraffic),
				ValidateFunc: validation.StringInSlice([]string{
					string(storagesyncservicesresource.IncomingTrafficPolicyAllowAllTraffic),
					string(storagesyncservicesresource.IncomingTrafficPolicyAllowVirtualNetworksOnly),
				}, false),
			},

			"registered_servers": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceStorageSyncCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.SyncServiceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := storagesyncservicesresource.NewStorageSyncServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.StorageSyncServicesGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError(storageSyncResourceName, id.ID())
	}

	parameters := storagesyncservicesresource.StorageSyncServiceCreateParameters{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &storagesyncservicesresource.StorageSyncServiceCreateParametersProperties{
			IncomingTrafficPolicy: pointer.To(storagesyncservicesresource.IncomingTrafficPolicy(d.Get("incoming_traffic_policy").(string))),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err = client.StorageSyncServicesCreateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err = pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceStorageSyncRead(d, meta)
}

func resourceStorageSyncRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.SyncServiceClient
	registeredServerClient := meta.(*clients.Client).Storage.SyncRegisteredServerClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := storagesyncservicesresource.ParseStorageSyncServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.StorageSyncServicesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return resourceStorageSyncFlatten(ctx, d, id, resp.Model, registeredServerClient, true)
}

func resourceStorageSyncFlatten(ctx context.Context, d *pluginsdk.ResourceData, id *storagesyncservicesresource.StorageSyncServiceId, model *storagesyncservicesresource.StorageSyncService, registeredServerClient *registeredserverresource.RegisteredServerResourceClient, includeRegisteredServers bool) error {
	d.Set("name", id.StorageSyncServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("incoming_traffic_policy", string(pointer.From(props.IncomingTrafficPolicy)))
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
			registeredServers := make([]interface{}, 0, len(*serverModel.Value))
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

func resourceStorageSyncUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.SyncServiceClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := storagesyncservicesresource.ParseStorageSyncServiceID(d.Id())
	if err != nil {
		return err
	}

	update := storagesyncservicesresource.StorageSyncServiceUpdateParameters{}

	if d.HasChange("tags") {
		update.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("incoming_traffic_policy") {
		update.Properties = &storagesyncservicesresource.StorageSyncServiceUpdateProperties{
			IncomingTrafficPolicy: pointer.To(storagesyncservicesresource.IncomingTrafficPolicy(d.Get("incoming_traffic_policy").(string))),
		}
	}

	if err = client.StorageSyncServicesUpdateThenPoll(ctx, *id, update); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceStorageSyncRead(d, meta)
}

func resourceStorageSyncDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.SyncServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := storagesyncservicesresource.ParseStorageSyncServiceID(d.Id())
	if err != nil {
		return err
	}

	if err = client.StorageSyncServicesDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
