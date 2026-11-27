// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_network_acl

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	signalrClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceWebPubsubNetworkACLCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := signalrClient.From(meta).WebPubSubClient.WebPubSub
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := webpubsub.ParseWebPubSubID(d.Get("web_pubsub_id").(string))
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("checking for present of existing %q: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving existing %s: `model` was nil", *id)
	}

	locks.ByName(id.WebPubSubName, "azurerm_web_pubsub")
	defer locks.UnlockByName(id.WebPubSubName, "azurerm_web_pubsub")

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			if !isNewNetworkACL(*existing.Model) {
				return tf.ImportAsExistsError("azurerm_web_pubsub_network_acl", id.ID())
			}
		}
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving existing %s: `model.Properties` was nil", *id)
	}

	payload := *existing.Model
	defaultAction := webpubsub.ACLAction(d.Get("default_action").(string))
	networkACL := webpubsub.WebPubSubNetworkACLs{
		DefaultAction:    &defaultAction,
		PublicNetwork:    expandWebpubsubPublicNetwork(d.Get("public_network").([]any)),
		PrivateEndpoints: expandWebpubsubPrivateEndpoint(d.Get("private_endpoint").(*pluginsdk.Set).List(), payload.Properties.PrivateEndpointConnections),
	}

	if defaultAction == webpubsub.ACLActionAllow && networkACL.PublicNetwork.Allow != nil && len(*networkACL.PublicNetwork.Allow) != 0 {
		return errors.New("when `default_action` is `Allow` for `public_network`, `allowed_request_types` cannot be specified")
	} else if defaultAction == webpubsub.ACLActionDeny && networkACL.PublicNetwork.Deny != nil && len(*networkACL.PublicNetwork.Deny) != 0 {
		return errors.New("when `default_action` is `Deny` for `public_network`, `denied_request_types` cannot be specified")
	}

	if networkACL.PrivateEndpoints != nil {
		for _, privateEndpoint := range *networkACL.PrivateEndpoints {
			if privateEndpoint.Allow != nil && privateEndpoint.Deny != nil && len(*privateEndpoint.Allow) != 0 && len(*privateEndpoint.Deny) != 0 {
				return errors.New("`allowed_request_types` and `denied_request_types` cannot be set together for `private_endpoint`")
			}

			if defaultAction == webpubsub.ACLActionAllow && privateEndpoint.Allow != nil && len(*privateEndpoint.Allow) != 0 {
				return errors.New("when `default_action` is `Allow` for `private_endpoint`, `allowed_request_types` cannot be specified")
			} else if defaultAction == webpubsub.ACLActionDeny && privateEndpoint.Deny != nil && len(*privateEndpoint.Deny) != 0 {
				return errors.New("when `default_action` is `Deny` for `private_endpoint`, `denied_request_types` cannot be specified")
			}
		}
	}
	payload.Properties.NetworkACLs = &networkACL

	if err := client.UpdateCallbackThenPoll(ctx, *id, payload, sdk.SetIDCallback(meta, id, d)); err != nil {
		return fmt.Errorf("updating Network ACL configuration for %q: %+v", id, err)
	}
	d.SetId(id.ID())

	return resourceWebPubsubNetworkACLRead(d, meta)
}

func expandWebpubsubPublicNetwork(input []any) *webpubsub.NetworkACL {
	allowRTs := make([]webpubsub.WebPubSubRequestType, 0)
	deniedRTs := make([]webpubsub.WebPubSubRequestType, 0)

	if len(input) == 0 || input[0] == nil {
		return &webpubsub.NetworkACL{}
	}

	v := input[0].(map[string]any)

	for _, item := range *pluginsdk.ExpandStringSlice(v["allowed_request_types"].(*pluginsdk.Set).List()) {
		allowRTs = append(allowRTs, webpubsub.WebPubSubRequestType(item))
	}

	for _, item := range *pluginsdk.ExpandStringSlice(v["denied_request_types"].(*pluginsdk.Set).List()) {
		deniedRTs = append(deniedRTs, webpubsub.WebPubSubRequestType(item))
	}

	return &webpubsub.NetworkACL{
		Allow: &allowRTs,
		Deny:  &deniedRTs,
	}
}

func expandWebpubsubPrivateEndpoint(input []any, privateEndpointConnections *[]webpubsub.PrivateEndpointConnection) *[]webpubsub.PrivateEndpointACL {
	results := make([]webpubsub.PrivateEndpointACL, 0)
	if privateEndpointConnections == nil {
		return &results
	}

	for _, privateEndpointConnection := range *privateEndpointConnections {
		result := webpubsub.PrivateEndpointACL{
			Allow: &[]webpubsub.WebPubSubRequestType{},
			Deny:  &[]webpubsub.WebPubSubRequestType{},
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

				allowedRTs := make([]webpubsub.WebPubSubRequestType, 0)
				for _, item := range *pluginsdk.ExpandStringSlice(v["allowed_request_types"].(*pluginsdk.Set).List()) {
					allowedRTs = append(allowedRTs, webpubsub.WebPubSubRequestType(item))
				}
				result.Allow = &allowedRTs

				deniedRTs := make([]webpubsub.WebPubSubRequestType, 0)
				for _, item := range *pluginsdk.ExpandStringSlice(v["denied_request_types"].(*pluginsdk.Set).List()) {
					deniedRTs = append(deniedRTs, webpubsub.WebPubSubRequestType(item))
				}
				result.Deny = &deniedRTs

				break
			}
		}
		results = append(results, result)
	}
	return &results
}
