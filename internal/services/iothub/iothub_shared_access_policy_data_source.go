// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceIotHubSharedAccessPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceIotHubSharedAccessPolicyRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`[a-zA-Z0-9!._-]{1,64}`), ""+
					"The shared access policy key name must not be empty, and must not exceed 64 characters in length.  The shared access policy key name can only contain alphanumeric characters, exclamation marks, periods, underscores and hyphens."),
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"iothub_name": {
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

func dataSourceIotHubSharedAccessPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	subscriptionId := meta.(*clients.Client).IoTHub.ResourceClient.SubscriptionID
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewSharedAccessPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iothub_name").(string), d.Get("name").(string))

	accessPolicy, err := client.GetKeysForKeyName(ctx, id.ResourceGroup, id.IotHubName, id.IotHubKeyName)
	if err != nil {
		if utils.ResponseWasNotFound(accessPolicy.Response) {
			return fmt.Errorf("Error: %s was not found", id)
		}

		return fmt.Errorf("loading %s: %+v", id, err)
	}

	iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil {
		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", id.IotHubName, id.ResourceGroup, err)
	}

	d.Set("name", id.IotHubKeyName)
	d.Set("iothub_name", id.IotHubName)
	d.Set("resource_group_name", id.ResourceGroup)

	d.SetId(id.ID())

	d.Set("primary_key", accessPolicy.PrimaryKey)
	if err := d.Set("primary_connection_string", getSharedAccessPolicyConnectionString(*iothub.Properties.HostName, id.IotHubKeyName, *accessPolicy.PrimaryKey)); err != nil {
		return fmt.Errorf("setting `primary_connection_string`: %v", err)
	}
	d.Set("secondary_key", accessPolicy.SecondaryKey)
	if err := d.Set("secondary_connection_string", getSharedAccessPolicyConnectionString(*iothub.Properties.HostName, id.IotHubKeyName, *accessPolicy.SecondaryKey)); err != nil {
		return fmt.Errorf("setting `secondary_connection_string`: %v", err)
	}

	return nil
}
