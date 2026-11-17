// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_port

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressrouteports"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmExpressRoutePortRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ExpressRoute.ExpressRoutePorts
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressrouteports.ParseExpressRoutePortID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ExpressRoutePortName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("peering_location", props.PeeringLocation)
			d.Set("bandwidth_in_gbps", props.BandwidthInGbps)
			d.Set("encapsulation", pointer.FromEnum(props.Encapsulation))
			d.Set("billing_type", pointer.FromEnum(props.BillingType))
			link1, link2, err := flattenExpressRoutePortLinks(props.Links)
			if err != nil {
				return fmt.Errorf("flattening links: %v", err)
			}
			if err := d.Set("link1", link1); err != nil {
				return fmt.Errorf("setting `link1`: %v", err)
			}
			if err := d.Set("link2", link2); err != nil {
				return fmt.Errorf("setting `link2`: %v", err)
			}
			d.Set("ethertype", props.EtherType)
			d.Set("guid", props.ResourceGuid)
			d.Set("mtu", props.Mtu)
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return nil
}

func flattenExpressRoutePortLinks(links *[]expressrouteports.ExpressRouteLink) ([]any, []any, error) {
	if links == nil {
		return []any{}, []any{}, nil
	}
	length := len(*links)
	if length != 2 {
		return nil, nil, fmt.Errorf("expected two links, but got %d", length)
	}

	return flattenExpressRoutePortLink((*links)[0]), flattenExpressRoutePortLink((*links)[1]), nil
}

func flattenExpressRoutePortLink(link expressrouteports.ExpressRouteLink) []any {
	var (
		routerName    string
		interfaceName string
		patchPanelId  string
		rackId        string
		connectorType string
		adminState    bool
		cknSecretId   string
		cakSecretId   string
		cipher        string
		sciState      bool
	)

	if props := link.Properties; props != nil {
		if props.RouterName != nil {
			routerName = *props.RouterName
		}
		if props.InterfaceName != nil {
			interfaceName = *props.InterfaceName
		}
		if props.PatchPanelId != nil {
			patchPanelId = *props.PatchPanelId
		}
		if props.RackId != nil {
			rackId = *props.RackId
		}
		connectorType = pointer.FromEnum(props.ConnectorType)
		adminState = pointer.From(props.AdminState) == expressrouteports.ExpressRouteLinkAdminStateEnabled
		sciState = pointer.From(props.MacSecConfig.SciState) == expressrouteports.ExpressRouteLinkMacSecSciStateEnabled
		if cfg := props.MacSecConfig; cfg != nil {
			if cfg.CknSecretIdentifier != nil {
				cknSecretId = *cfg.CknSecretIdentifier
			}
			if cfg.CakSecretIdentifier != nil {
				cakSecretId = *cfg.CakSecretIdentifier
			}
			cipher = pointer.FromEnum(cfg.Cipher)
		}
	}

	return []any{
		map[string]any{
			"id":                            pointer.From(link.Id),
			"router_name":                   routerName,
			"interface_name":                interfaceName,
			"patch_panel_id":                patchPanelId,
			"rack_id":                       rackId,
			"connector_type":                connectorType,
			"admin_enabled":                 adminState,
			"macsec_ckn_keyvault_secret_id": cknSecretId,
			"macsec_cak_keyvault_secret_id": cakSecretId,
			"macsec_cipher":                 cipher,
			"macsec_sci_enabled":            sciState,
		},
	}
}
