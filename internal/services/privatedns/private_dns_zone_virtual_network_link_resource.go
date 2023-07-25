// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatedns

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2020-06-01/virtualnetworklinks"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourcePrivateDnsZoneVirtualNetworkLink() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePrivateDnsZoneVirtualNetworkLinkCreateUpdate,
		Read:   resourcePrivateDnsZoneVirtualNetworkLinkRead,
		Update: resourcePrivateDnsZoneVirtualNetworkLinkCreateUpdate,
		Delete: resourcePrivateDnsZoneVirtualNetworkLinkDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualnetworklinks.ParseVirtualNetworkLinkID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			// TODO: these can become case-sensitive with a state migration
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				// TODO: make this case sensitive once the API's fixed https://github.com/Azure/azure-rest-api-specs/issues/10933
				DiffSuppressFunc: suppress.CaseDifference,
			},

			// TODO: in 4.0 switch this to `private_dns_zone_id`
			"private_dns_zone_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// TODO: make this case sensitive once the API's fixed https://github.com/Azure/azure-rest-api-specs/issues/10933
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"virtual_network_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateVirtualNetworkID,
			},

			"registration_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourcePrivateDnsZoneVirtualNetworkLinkCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.VirtualNetworkLinksClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualnetworklinks.NewVirtualNetworkLinkID(subscriptionId, d.Get("resource_group_name").(string), d.Get("private_dns_zone_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_private_dns_zone_virtual_network_link", id.ID())
		}
	}

	parameters := virtualnetworklinks.VirtualNetworkLink{
		Location: utils.String("global"),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: &virtualnetworklinks.VirtualNetworkLinkProperties{
			VirtualNetwork: &virtualnetworklinks.SubResource{
				Id: utils.String(d.Get("virtual_network_id").(string)),
			},
			RegistrationEnabled: utils.Bool(d.Get("registration_enabled").(bool)),
		},
	}

	options := virtualnetworklinks.CreateOrUpdateOperationOptions{
		IfMatch:     utils.String(""),
		IfNoneMatch: utils.String(""),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters, options); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourcePrivateDnsZoneVirtualNetworkLinkRead(d, meta)
}

func resourcePrivateDnsZoneVirtualNetworkLinkRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.VirtualNetworkLinksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworklinks.ParseVirtualNetworkLinkID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.VirtualNetworkLinkName)
	d.Set("private_dns_zone_name", id.PrivateDnsZoneName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("registration_enabled", props.RegistrationEnabled)

			if network := props.VirtualNetwork; network != nil {
				d.Set("virtual_network_id", network.Id)
			}
		}
		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourcePrivateDnsZoneVirtualNetworkLinkDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.VirtualNetworkLinksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworklinks.ParseVirtualNetworkLinkID(d.Id())
	if err != nil {
		return err
	}

	options := virtualnetworklinks.DeleteOperationOptions{IfMatch: utils.String("")}

	if err = client.DeleteThenPoll(ctx, *id, options); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	// whilst the Delete above returns a Future, the Azure API's broken such that even though it's marked as "gone"
	// it's still kicking around - so we have to poll until this is actually gone
	log.Printf("[DEBUG] Waiting for %s to be deleted", *id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"Available"},
		Target:  []string{"NotFound"},
		Refresh: func() (interface{}, string, error) {
			log.Printf("[DEBUG] Checking to see if %s is still available", *id)
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					log.Printf("[DEBUG] %s was not found", *id)
					return "NotFound", "NotFound", nil
				}

				return "", "error", err
			}

			log.Printf("[DEBUG] %s still exists", *id)
			return "Available", "Available", nil
		},
		Delay:                     30 * time.Second,
		PollInterval:              10 * time.Second,
		ContinuousTargetOccurence: 10,
		Timeout:                   d.Timeout(pluginsdk.TimeoutDelete),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}
