// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package private_link_service

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privatelinkservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePrivateLinkServiceUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.PrivateLinkServices
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privatelinkservices.ParsePrivateLinkServiceID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, privatelinkservices.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("auto_approval_subscription_ids") {
		payload.Properties.AutoApproval = &privatelinkservices.ResourceSet{
			Subscriptions: pluginsdk.ExpandStringSlice(d.Get("auto_approval_subscription_ids").(*pluginsdk.Set).List()),
		}
	}

	if d.HasChange("proxy_protocol_enabled") {
		payload.Properties.EnableProxyProtocol = pointer.To(d.Get("proxy_protocol_enabled").(bool))
	}

	if d.HasChange("visibility_subscription_ids") {
		payload.Properties.Visibility = &privatelinkservices.ResourceSet{
			Subscriptions: pluginsdk.ExpandStringSlice(d.Get("visibility_subscription_ids").(*pluginsdk.Set).List()),
		}
	}

	if d.HasChange("fqdns") {
		payload.Properties.Fqdns = pluginsdk.ExpandStringSlice(d.Get("fqdns").([]any))
	}

	if d.HasChange("nat_ip_configuration") {
		payload.Properties.IPConfigurations = expandPrivateLinkServiceIPConfiguration(d.Get("nat_ip_configuration").([]any))
	}

	if d.HasChange("destination_ip_address") {
		payload.Properties.DestinationIPAddress = pointer.To(d.Get("destination_ip_address").(string))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	// we can't rely on the use of the Future here due to the resource being successfully completed but now the service is applying those values.
	// currently being tracked with issue #6466: https://github.com/Azure/azure-sdk-for-go/issues/6466
	log.Printf("[DEBUG] Waiting for %s to finish applying", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Pending", "Updating", "Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    privateLinkServiceWaitForReadyRefreshFunc(ctx, client, *id),
		MinTimeout: 15 * time.Second,
	}

	stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become available: %s", id, err)
	}

	d.SetId(id.ID())

	return resourcePrivateLinkServiceRead(d, meta)
}
