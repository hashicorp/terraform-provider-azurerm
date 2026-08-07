// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package testdata

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/virtualwans"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name virtual_hub_ip -service-package-name network -properties "name"

func resourceVirtualHubIP() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualHubIPCreate,
		Read:   resourceVirtualHubIPRead,
		Update: resourceVirtualHubIPUpdate,
		Delete: resourceVirtualHubIPDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&commonids.VirtualHubIPConfigurationId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.VirtualHubIPConfigurationId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"private_ip_address": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceVirtualHubIPCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewVirtualHubIPConfigurationID("sub", "rg", "hub", d.Get("name").(string))

	if err := client.VirtualHubIPConfigurationCreateOrUpdateThenPoll(ctx, id, virtualwans.HubIPConfiguration{}); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubIPRead(d, meta)
}

func resourceVirtualHubIPUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceVirtualHubIPRead(d, meta)
}

func resourceVirtualHubIPRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVirtualHubIPConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.VirtualHubIPConfigurationGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.IpConfigName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("private_ip_address", props.PrivateIPAddress)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func resourceVirtualHubIPDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVirtualHubIPConfigurationID(d.Id())
	if err != nil {
		return err
	}

	if err := client.VirtualHubIPConfigurationDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
