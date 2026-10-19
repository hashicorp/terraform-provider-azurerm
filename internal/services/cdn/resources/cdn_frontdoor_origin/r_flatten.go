// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_origin

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigins"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceCdnFrontDoorOriginFlatten(d *pluginsdk.ResourceData, id *afdorigins.OriginGroupOriginId, model *afdorigins.AFDOrigin) error {
	d.Set("name", id.OriginName)
	d.Set("cdn_frontdoor_origin_group_id", afdorigins.NewOriginGroupID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.OriginGroupName).ID())

	if model != nil {
		if props := model.Properties; props != nil {
			if err := d.Set("private_link", flattenCdnFrontDoorOriginPrivateLinkSettings(props.SharedPrivateLinkResource)); err != nil {
				return fmt.Errorf("setting `private_link`: %+v", err)
			}

			d.Set("certificate_name_check_enabled", props.EnforceCertificateNameCheck)
			d.Set("enabled", pointer.From(props.EnabledState) == afdorigins.EnabledStateEnabled)
			d.Set("host_name", props.HostName)
			d.Set("http_port", props.HTTPPort)
			d.Set("https_port", props.HTTPSPort)
			d.Set("origin_host_header", props.OriginHostHeader)
			d.Set("priority", props.Priority)
			d.Set("weight", props.Weight)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenCdnFrontDoorOriginPrivateLinkSettings(input *afdorigins.SharedPrivateLinkResourceProperties) []any {
	if input == nil {
		return []any{}
	}

	privateLinkTargetId := ""
	if input.PrivateLink != nil && input.PrivateLink.Id != nil {
		privateLinkTargetId = *input.PrivateLink.Id
	}

	return []any{
		map[string]any{
			"location":               location.NormalizeNilable(input.PrivateLinkLocation),
			"private_link_target_id": privateLinkTargetId,
			"request_message":        pointer.From(input.RequestMessage),
			"target_type":            pointer.From(input.GroupId),
		},
	}
}
