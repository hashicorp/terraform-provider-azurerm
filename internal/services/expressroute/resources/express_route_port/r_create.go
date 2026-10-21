// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package express_route_port

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressrouteports"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmExpressRoutePortCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ExpressRoute.ExpressRoutePorts
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := expressrouteports.NewExpressRoutePortID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_express_route_port", id.ID())
		}
	}

	expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	param := expressrouteports.ExpressRoutePort{
		Name:     pointer.To(id.ExpressRoutePortName),
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &expressrouteports.ExpressRoutePortPropertiesFormat{
			PeeringLocation: pointer.To(d.Get("peering_location").(string)),
			BandwidthInGbps: pointer.To(int64(d.Get("bandwidth_in_gbps").(int))),
			Encapsulation:   pointer.ToEnum[expressrouteports.ExpressRoutePortsEncapsulation](d.Get("encapsulation").(string)),
		},
		Identity: expandedIdentity,
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("billing_type"); ok {
		param.Properties.BillingType = pointer.ToEnum[expressrouteports.ExpressRoutePortsBillingType](v.(string))
	}

	// a lock is needed here for subresource express_route_port_authorization needs a lock.
	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, param, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	// The link properties can't be specified in first creation. It will result into either error (e.g. setting `adminState`) or being ignored (e.g. setting MACSec)
	param.Properties.Links = expandExpressRoutePortLinks(d.Get("link1").([]any), d.Get("link2").([]any))

	if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	return resourceArmExpressRoutePortRead(d, meta)
}

func expandExpressRoutePortLinks(link1, link2 []any) *[]expressrouteports.ExpressRouteLink {
	var out []expressrouteports.ExpressRouteLink
	if link := expandExpressRoutePortLink(1, link1); link != nil {
		out = append(out, *link)
	}
	if link := expandExpressRoutePortLink(2, link2); link != nil {
		out = append(out, *link)
	}
	if len(out) == 0 {
		return nil
	}
	return &out
}

func expandExpressRoutePortLink(idx int, input []any) *expressrouteports.ExpressRouteLink {
	if len(input) == 0 {
		return nil
	}

	b := input[0].(map[string]any)
	adminState := expressrouteports.ExpressRouteLinkAdminStateDisabled
	if b["admin_enabled"].(bool) {
		adminState = expressrouteports.ExpressRouteLinkAdminStateEnabled
	}

	sciState := expressrouteports.ExpressRouteLinkMacSecSciStateDisabled
	if b["macsec_sci_enabled"].(bool) {
		sciState = expressrouteports.ExpressRouteLinkMacSecSciStateEnabled
	}

	link := expressrouteports.ExpressRouteLink{
		// The link name is fixed
		Name: pointer.To(fmt.Sprintf("link%d", idx)),
		Properties: &expressrouteports.ExpressRouteLinkPropertiesFormat{
			AdminState: pointer.To(adminState),
			MacSecConfig: &expressrouteports.ExpressRouteLinkMacSecConfig{
				Cipher:   pointer.ToEnum[expressrouteports.ExpressRouteLinkMacSecCipher](b["macsec_cipher"].(string)),
				SciState: pointer.To(sciState),
			},
		},
	}

	if cknSecretId := b["macsec_ckn_keyvault_secret_id"].(string); cknSecretId != "" {
		link.Properties.MacSecConfig.CknSecretIdentifier = &cknSecretId
	}
	if cakSecretId := b["macsec_cak_keyvault_secret_id"].(string); cakSecretId != "" {
		link.Properties.MacSecConfig.CakSecretIdentifier = &cakSecretId
	}
	return &link
}
