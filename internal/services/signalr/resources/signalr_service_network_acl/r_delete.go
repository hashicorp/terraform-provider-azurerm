// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_service_network_acl

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSignalRServiceNetworkACLDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SignalR.SignalRClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := signalr.ParseSignalRID(d.Id())
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

	defaultRequestTypes := []signalr.SignalRRequestType{
		signalr.SignalRRequestTypeClientConnection,
		signalr.SignalRRequestTypeRESTAPI,
		signalr.SignalRRequestTypeServerConnection,
		signalr.SignalRRequestTypeTrace,
	}
	networkACL := &signalr.SignalRNetworkACLs{
		DefaultAction: pointer.To(signalr.ACLActionDeny),
		PublicNetwork: &signalr.NetworkACL{
			Allow: &defaultRequestTypes,
		},
	}

	if model.Properties != nil && model.Properties.NetworkACLs != nil && model.Properties.NetworkACLs.PrivateEndpoints != nil {
		privateEndpoints := make([]signalr.PrivateEndpointACL, 0)
		for _, item := range *model.Properties.NetworkACLs.PrivateEndpoints {
			privateEndpoints = append(privateEndpoints, signalr.PrivateEndpointACL{
				Allow: &defaultRequestTypes,
				Name:  item.Name,
			})
		}
		networkACL.PrivateEndpoints = &privateEndpoints
	}

	if model.Properties != nil {
		model.Properties.NetworkACLs = networkACL
	}

	if err := client.UpdateThenPoll(ctx, *id, model); err != nil {
		return fmt.Errorf("resetting the default Network ACL configuration for %s: %+v", *id, err)
	}

	return nil
}
