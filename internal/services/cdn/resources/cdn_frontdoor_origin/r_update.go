// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_origin

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigins"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorOriginUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginsClient
	profileClient := meta.(*clients.Client).Cdn.FrontDoorProfilesClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afdorigins.ParseOriginGroupOriginID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: properties was nil", id)
	}
	props := existing.Model.Properties

	if d.HasChange("certificate_name_check_enabled") {
		props.EnforceCertificateNameCheck = pointer.To(d.Get("certificate_name_check_enabled").(bool))
	}

	if d.HasChange("enabled") {
		props.EnabledState = pointer.To(afdorigins.EnabledStateDisabled)
		if d.Get("enabled").(bool) {
			props.EnabledState = pointer.To(afdorigins.EnabledStateEnabled)
		}
	}

	if d.HasChange("host_name") {
		props.HostName = pointer.To(d.Get("host_name").(string))
	}

	if d.HasChange("http_port") {
		props.HTTPPort = pointer.To(int64(d.Get("http_port").(int)))
	}

	if d.HasChange("https_port") {
		props.HTTPSPort = pointer.To(int64(d.Get("https_port").(int)))
	}

	if d.HasChange("origin_host_header") {
		props.OriginHostHeader = pointer.ToOrNil(d.Get("origin_host_header").(string))
	}

	if d.HasChange("private_link") {
		// I need to get the profile SKU so I know if it is valid or not to define a private link as
		// private links are only allowed in the premium sku...
		profileId := profiles.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName)

		profileResp, err := profileClient.Get(ctx, profileId)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", profileId, err)
		}

		profileModel := profileResp.Model
		if profileModel == nil {
			return fmt.Errorf("retrieving %s: model was nil", profileId)
		}

		if profileModel.Sku.Name == nil {
			return fmt.Errorf("retrieving %s: sku was nil", profileId)
		}

		privateLinkSettings, err := expandCdnFrontDoorOriginPrivateLinkSettings(d.Get("private_link").([]any), pointer.From(profileModel.Sku.Name), d.Get("certificate_name_check_enabled").(bool))
		if err != nil {
			return err
		}

		props.SharedPrivateLinkResource = privateLinkSettings
	}

	if d.HasChange("priority") {
		props.Priority = pointer.To(int64(d.Get("priority").(int)))
	}

	if d.HasChange("weight") {
		props.Weight = pointer.To(int64(d.Get("weight").(int)))
	}

	if err := client.CreateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceCdnFrontDoorOriginRead(d, meta)
}
