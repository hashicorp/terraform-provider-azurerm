// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package privatedns

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/virtualnetworklinks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name private_dns_zone_virtual_network_link -properties "name" -compare-values "subscription_id:private_dns_zone_id,resource_group_name:private_dns_zone_id,private_dns_zone_name:private_dns_zone_id"

func resourcePrivateDnsZoneVirtualNetworkLink() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create:   resourcePrivateDnsZoneVirtualNetworkLinkCreateUpdate,
		Read:     resourcePrivateDnsZoneVirtualNetworkLinkRead,
		Update:   resourcePrivateDnsZoneVirtualNetworkLinkCreateUpdate,
		Delete:   resourcePrivateDnsZoneVirtualNetworkLinkDelete,
		Importer: pluginsdk.ImporterValidatingIdentity(&virtualnetworklinks.VirtualNetworkLinkId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&virtualnetworklinks.VirtualNetworkLinkId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"private_dns_zone_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: privatezones.ValidatePrivateDnsZoneID,
			},

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

			"resolution_policy": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				// Note: O+C because when the `name` of `azurerm_private_dns_zone` is a private link endpoint, the service will set default value for this.
				Computed:     true,
				ValidateFunc: validation.StringInSlice(virtualnetworklinks.PossibleValuesForResolutionPolicy(), false),
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourcePrivateDnsZoneVirtualNetworkLinkCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.VirtualNetworkLinksClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	privateDNSZoneID, err := privatezones.ParsePrivateDnsZoneID(d.Get("private_dns_zone_id").(string))
	if err != nil {
		return err
	}

	id := virtualnetworklinks.NewVirtualNetworkLinkID(privateDNSZoneID.SubscriptionId, privateDNSZoneID.ResourceGroupName, privateDNSZoneID.PrivateDnsZoneName, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
	}

	parameters := virtualnetworklinks.VirtualNetworkLink{
		Location: pointer.To("global"),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: &virtualnetworklinks.VirtualNetworkLinkProperties{
			VirtualNetwork: &virtualnetworklinks.SubResource{
				Id: pointer.To(d.Get("virtual_network_id").(string)),
			},
			RegistrationEnabled: pointer.To(d.Get("registration_enabled").(bool)),
		},
	}

	if v, ok := d.GetOk("resolution_policy"); ok {
		parameters.Properties.ResolutionPolicy = pointer.ToEnum[virtualnetworklinks.ResolutionPolicy](v.(string))
	}

	options := virtualnetworklinks.CreateOrUpdateOperationOptions{
		IfMatch:     pointer.To(""),
		IfNoneMatch: pointer.To(""),
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, options, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
		if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
			return err
		}
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, parameters, options); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

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
	return resourcePrivateDnsZoneVirtualNetworkLinkFlatten(d, id, resp.Model)
}

func resourcePrivateDnsZoneVirtualNetworkLinkFlatten(d *pluginsdk.ResourceData, id *virtualnetworklinks.VirtualNetworkLinkId, model *virtualnetworklinks.VirtualNetworkLink) error {
	d.Set("name", id.VirtualNetworkLinkName)
	d.Set("private_dns_zone_id", privatezones.NewPrivateDnsZoneID(id.SubscriptionId, id.ResourceGroupName, id.PrivateDnsZoneName).ID())

	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("registration_enabled", props.RegistrationEnabled)
			d.Set("resolution_policy", pointer.From(props.ResolutionPolicy))

			if network := props.VirtualNetwork; network != nil {
				d.Set("virtual_network_id", network.Id)
			}
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func resourcePrivateDnsZoneVirtualNetworkLinkDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.VirtualNetworkLinksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworklinks.ParseVirtualNetworkLinkID(d.Id())
	if err != nil {
		return err
	}

	options := virtualnetworklinks.DeleteOperationOptions{IfMatch: pointer.To("")}

	if err = client.DeleteThenPoll(ctx, *id, options); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	// whilst the Delete above returns a Future, the Azure API's broken such that even though it's marked as "gone"
	// it's still kicking around - so we have to poll until this is actually gone
	poller := custompollers.NewEventualConsistencyPoller(10, func(pollerCtx context.Context) (*http.Response, error) {
		resp, err := client.Get(pollerCtx, *id)
		return resp.HttpResponse, err
	}, custompollers.DefaultDeletionEventualConsistencyPollerOptions())
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}
