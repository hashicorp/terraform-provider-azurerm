// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_network_acl

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	signalrClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceWebpubsubNetworkACLDelete(d *pluginsdk.ResourceData, meta any) error {
	client := signalrClient.From(meta).WebPubSubClient.WebPubSub
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := webpubsub.ParseWebPubSubID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	var denyRequestTypes []webpubsub.WebPubSubRequestType
	networkACL := &webpubsub.WebPubSubNetworkACLs{
		DefaultAction: pointer.To(webpubsub.ACLActionDeny),
		PublicNetwork: &webpubsub.NetworkACL{
			Allow: &defaultRequestTypes,
			Deny:  &denyRequestTypes,
		},
	}

	if resp.Model == nil {
		return fmt.Errorf("retrieving existing %s: `model` was nil", *id)
	}
	if resp.Model.Properties == nil {
		return fmt.Errorf("retrieving existing %s: `model.Properties` was nil", *id)
	}

	payload := *resp.Model
	if payload.Properties.NetworkACLs != nil && payload.Properties.NetworkACLs.PrivateEndpoints != nil {
		privateEndpoints := make([]webpubsub.PrivateEndpointACL, 0)
		for _, item := range *payload.Properties.NetworkACLs.PrivateEndpoints {
			privateEndpoints = append(privateEndpoints, webpubsub.PrivateEndpointACL{
				Allow: &defaultRequestTypes,
				Name:  item.Name,
			})
		}
		networkACL.PrivateEndpoints = &privateEndpoints
	}

	payload.Properties.NetworkACLs = networkACL

	if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("resetting the default Network ACL configuration for %q: %+v", id, err)
	}

	return nil
}
