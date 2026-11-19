// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_service_network_acl

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSignalRServiceNetworkACLCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SignalR.SignalRClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := signalr.ParseSignalRID(d.Get("signalr_service_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(id.SignalRName, "azurerm_signalr_service")
	defer locks.UnlockByName(id.SignalRName, "azurerm_signalr_service")

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", *id)
	}

	model := *resp.Model
	if props := model.Properties; props != nil {
		defaultAction := signalr.ACLAction(d.Get("default_action").(string))
		networkACL := signalr.SignalRNetworkACLs{
			DefaultAction: &defaultAction,
			PublicNetwork: expandSignalRServicePublicNetwork(d.Get("public_network").([]any)),
		}

		if v, ok := d.GetOk("private_endpoint"); ok {
			networkACL.PrivateEndpoints = expandSignalRServicePrivateEndpoint(v.(*pluginsdk.Set).List(), props.PrivateEndpointConnections)
		}

		if defaultAction == signalr.ACLActionAllow && len(*networkACL.PublicNetwork.Allow) != 0 {
			return fmt.Errorf("when `default_action` is `Allow` for `public_network`, `allowed_request_types` cannot be specified")
		} else if defaultAction == signalr.ACLActionDeny && len(*networkACL.PublicNetwork.Deny) != 0 {
			return fmt.Errorf("when `default_action` is `Deny` for `public_network`, `denied_request_types` cannot be specified")
		}

		if networkACL.PrivateEndpoints != nil {
			for _, privateEndpoint := range *networkACL.PrivateEndpoints {
				if len(*privateEndpoint.Allow) != 0 && len(*privateEndpoint.Deny) != 0 {
					return fmt.Errorf("`allowed_request_types` and `denied_request_types` cannot be set together for `private_endpoint`")
				}

				if defaultAction == signalr.ACLActionAllow && len(*privateEndpoint.Allow) != 0 {
					return fmt.Errorf("when `default_action` is `Allow` for `private_endpoint`, `allowed_request_types` cannot be specified")
				} else if defaultAction == signalr.ACLActionDeny && len(*privateEndpoint.Deny) != 0 {
					return fmt.Errorf("when `default_action` is `Deny` for `private_endpoint`, `denied_request_types` cannot be specified")
				}
			}
		}

		model.Properties.NetworkACLs = &networkACL
	}

	if d.IsNewResource() {
		if err := client.UpdateCallbackThenPoll(ctx, *id, model, sdk.SetIDCallback(meta, id, d)); err != nil {
			return fmt.Errorf("creating Network ACL for %s: %v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.UpdateThenPoll(ctx, *id, model); err != nil {
			return fmt.Errorf("updating Network ACL for %s: %v", id, err)
		}
	}

	return resourceSignalRServiceNetworkACLRead(d, meta)
}

func expandSignalRServicePublicNetwork(input []any) *signalr.NetworkACL {
	allowedRTs := make([]signalr.SignalRRequestType, 0)
	deniedRTs := make([]signalr.SignalRRequestType, 0)

	if len(input) != 0 && input[0] != nil {
		v := input[0].(map[string]any)

		for _, item := range *pluginsdk.ExpandStringSlice(v["allowed_request_types"].(*pluginsdk.Set).List()) {
			allowedRTs = append(allowedRTs, signalr.SignalRRequestType(item))
		}

		for _, item := range *pluginsdk.ExpandStringSlice(v["denied_request_types"].(*pluginsdk.Set).List()) {
			deniedRTs = append(deniedRTs, signalr.SignalRRequestType(item))
		}
	}

	return &signalr.NetworkACL{
		Allow: &allowedRTs,
		Deny:  &deniedRTs,
	}
}

func expandSignalRServicePrivateEndpoint(input []any, privateEndpointConnections *[]signalr.PrivateEndpointConnection) *[]signalr.PrivateEndpointACL {
	results := make([]signalr.PrivateEndpointACL, 0)
	if privateEndpointConnections == nil {
		return &results
	}

	for _, privateEndpointConnection := range *privateEndpointConnections {
		result := signalr.PrivateEndpointACL{
			Allow: &[]signalr.SignalRRequestType{},
			Deny:  &[]signalr.SignalRRequestType{},
		}

		if privateEndpointConnection.Name != nil {
			result.Name = *privateEndpointConnection.Name
		}

		for _, item := range input {
			v := item.(map[string]any)
			privateEndpointId := v["id"].(string)

			if props := privateEndpointConnection.Properties; props != nil {
				if props.PrivateEndpoint == nil || props.PrivateEndpoint.Id == nil || privateEndpointId != *props.PrivateEndpoint.Id {
					continue
				}

				allowedRTs := make([]signalr.SignalRRequestType, 0)
				for _, item := range *pluginsdk.ExpandStringSlice(v["allowed_request_types"].(*pluginsdk.Set).List()) {
					allowedRTs = append(allowedRTs, signalr.SignalRRequestType(item))
				}
				result.Allow = &allowedRTs

				deniedRTs := make([]signalr.SignalRRequestType, 0)
				for _, item := range *pluginsdk.ExpandStringSlice(v["denied_request_types"].(*pluginsdk.Set).List()) {
					deniedRTs = append(deniedRTs, signalr.SignalRRequestType(item))
				}
				result.Deny = &deniedRTs

				break
			}
		}

		results = append(results, result)
	}

	return &results
}
