// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/iotdpsresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceIotHubDPSSharedAccessPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceIotHubDPSSharedAccessPolicyRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.IotHubSharedAccessPolicyName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"iothub_dps_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"primary_key": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"primary_connection_string": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"secondary_key": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"secondary_connection_string": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

func dataSourceIotHubDPSSharedAccessPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewProvisioningServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iothub_dps_name").(string))

	iothubDps, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(iothubDps.HttpResponse) {
			return fmt.Errorf("Error: IotHub DPS %q was not found", id)
		}

		return fmt.Errorf("retrieving IotHub DPS %q: %+v", id, err)
	}

	keyId := iotdpsresource.NewKeyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iothub_dps_name").(string), d.Get("name").(string))
	accessPolicy, err := client.ListKeysForKeyName(ctx, keyId)
	if err != nil {
		if response.WasNotFound(accessPolicy.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("loading %s: %+v", id, err)
	}

	d.Set("name", keyId.KeyName)
	d.Set("resource_group_name", keyId.ResourceGroupName)

	d.SetId(id.ID())
	if model := accessPolicy.Model; model != nil {
		d.Set("primary_key", model.PrimaryKey)
		d.Set("secondary_key", model.SecondaryKey)

		if dpsModel := iothubDps.Model; dpsModel != nil {
			properties := dpsModel.Properties
			primaryConnectionString := ""
			secondaryConnectionString := ""
			if properties.ServiceOperationsHostName != nil {
				hostname := properties.ServiceOperationsHostName
				if primary := model.PrimaryKey; primary != nil {
					primaryConnectionString = getSAPConnectionString(*hostname, keyId.KeyName, *primary)
				}
				if secondary := model.SecondaryKey; secondary != nil {
					secondaryConnectionString = getSAPConnectionString(*hostname, keyId.KeyName, *secondary)
				}
			}
			d.Set("primary_connection_string", primaryConnectionString)
			d.Set("secondary_connection_string", secondaryConnectionString)
		}
	}

	return nil
}
