// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceIotHubDPS() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceIotHubDPSRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"allocation_policy": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"device_provisioning_host_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"id_scope": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"service_operations_host_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func dataSourceIotHubDPSRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewProvisioningServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("Error: IoT Device Provisioning Service %s was not found", id)
		}

		return fmt.Errorf("retrieving IoT Device Provisioning Service %s: %+v", id, err)
	}

	d.Set("name", id.ProvisioningServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		props := model.Properties
		d.Set("service_operations_host_name", props.ServiceOperationsHostName)
		d.Set("device_provisioning_host_name", props.DeviceProvisioningHostName)
		d.Set("id_scope", props.IdScope)
		d.Set("allocation_policy", string(pointer.From(props.AllocationPolicy)))
		d.Set("tags", flattenTags(model.Tags))
	}

	return nil
}
