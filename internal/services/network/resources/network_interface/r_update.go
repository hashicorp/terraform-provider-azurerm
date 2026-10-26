// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_interface

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkInterfaceUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseNetworkInterfaceID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(strings.ToLower(id.ID()))
	defer locks.UnlockByID(strings.ToLower(id.ID()))

	// first get the existing one so that we can pull things as needed
	existing, err := client.Get(ctx, *id, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	// then pull out things we need to lock on
	info := helpers.ParseFieldsFromNetworkInterface(*existing.Model.Properties)

	payload := existing.Model

	// For NIC attached to private endpoint, tags cannot be updated using PUT.
	// It has to use the specific update tag PATCH API.
	// An issue has been raised to improve the API design: https://github.com/Azure/azure-rest-api-specs/issues/34437
	propsOtherThanTagsUpdated := false
	attachedToPrivateEndpoint := payload.Properties.PrivateEndpoint != nil && pointer.From(payload.Properties.PrivateEndpoint.Id) != ""

	if d.HasChange("auxiliary_mode") {
		propsOtherThanTagsUpdated = true
		if auxiliaryMode := d.Get("auxiliary_mode").(string); auxiliaryMode != "" {
			payload.Properties.AuxiliaryMode = pointer.ToEnum[networkinterfaces.NetworkInterfaceAuxiliaryMode](auxiliaryMode)
		} else {
			payload.Properties.AuxiliaryMode = nil
		}
	}

	if d.HasChange("auxiliary_sku") {
		propsOtherThanTagsUpdated = true
		if auxiliarySku := d.Get("auxiliary_sku").(string); auxiliarySku != "" {
			payload.Properties.AuxiliarySku = pointer.ToEnum[networkinterfaces.NetworkInterfaceAuxiliarySku](auxiliarySku)
		} else {
			payload.Properties.AuxiliarySku = nil
		}
	}

	if d.HasChange("dns_servers") {
		propsOtherThanTagsUpdated = true
		dnsServersRaw := d.Get("dns_servers").([]any)
		dnsServers := expandNetworkInterfaceDnsServers(dnsServersRaw)

		payload.Properties.DnsSettings.DnsServers = &dnsServers
	}

	if d.HasChange("accelerated_networking_enabled") {
		propsOtherThanTagsUpdated = true
		payload.Properties.EnableAcceleratedNetworking = pointer.To(d.Get("accelerated_networking_enabled").(bool))
	}

	if d.HasChange("ip_forwarding_enabled") {
		propsOtherThanTagsUpdated = true
		payload.Properties.EnableIPForwarding = pointer.To(d.Get("ip_forwarding_enabled").(bool))
	}

	if d.HasChange("internal_dns_name_label") {
		propsOtherThanTagsUpdated = true
		payload.Properties.DnsSettings.InternalDnsNameLabel = pointer.To(d.Get("internal_dns_name_label").(string))
	}

	if d.HasChange("ip_configuration") {
		propsOtherThanTagsUpdated = true
		ipConfigsRaw := d.Get("ip_configuration").([]any)
		ipConfigs, err := expandNetworkInterfaceIPConfigurations(ipConfigsRaw)
		if err != nil {
			return fmt.Errorf("expanding `ip_configuration`: %+v", err)
		}
		lockingDetails, err := determineResourcesToLockFromIPConfiguration(ipConfigs)
		if err != nil {
			return fmt.Errorf("determining locking details: %+v", err)
		}

		lockingDetails.lock()
		defer lockingDetails.unlock()

		// then map the fields managed in other resources back
		ipConfigs = helpers.MapFieldsToNetworkInterface(ipConfigs, info)

		payload.Properties.IPConfigurations = ipConfigs
	}

	if d.HasChange("tags") && !attachedToPrivateEndpoint {
		tagsRaw := d.Get("tags").(map[string]any)
		payload.Tags = tags.Expand(tagsRaw)
	}

	if propsOtherThanTagsUpdated || !attachedToPrivateEndpoint {
		if err = client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
			return fmt.Errorf("updating %s: %+v", *id, err)
		}
	}

	if d.HasChange("tags") && attachedToPrivateEndpoint {
		tagsRaw := d.Get("tags").(map[string]any)
		tags := networkinterfaces.TagsObject{
			Tags: tags.Expand(tagsRaw),
		}
		if _, err = client.UpdateTags(ctx, *id, tags); err != nil {
			return fmt.Errorf("updating tags for %s: %+v", *id, err)
		}
	}

	return resourceNetworkInterfaceRead(d, meta)
}
