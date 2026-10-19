// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_origin

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigingroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigins"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privatelinkservices"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorOriginCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginsClient
	profileClient := meta.(*clients.Client).Cdn.FrontDoorProfilesClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	originGroup, err := afdorigingroups.ParseOriginGroupID(d.Get("cdn_frontdoor_origin_group_id").(string))
	if err != nil {
		return err
	}

	id := afdorigins.NewOriginGroupOriginID(originGroup.SubscriptionId, originGroup.ResourceGroupName, originGroup.ProfileName, originGroup.OriginGroupName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError(azureCdnFrontDoorOriginResourceName, id.ID())
		}
	}

	// I need to get the profile SKU so I know if it is valid or not to define a private link as
	// private links are only allowed in the premium sku...
	profileId := profiles.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName)
	profileResp, err := profileClient.Get(ctx, profileId)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", profileId, err)
	}

	if profileResp.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", profileId)
	}

	if profileResp.Model.Sku.Name == nil {
		return fmt.Errorf("retrieving %s: sku was nil", profileId)
	}

	enabledState := afdorigins.EnabledStateDisabled
	if d.Get("enabled").(bool) {
		enabledState = afdorigins.EnabledStateEnabled
	}

	enforceCertificateNameCheck := d.Get("certificate_name_check_enabled").(bool)
	expandedPrivateLink, err := expandCdnFrontDoorOriginPrivateLinkSettings(d.Get("private_link").([]any), pointer.From(profileResp.Model.Sku.Name), enforceCertificateNameCheck)
	if err != nil {
		return err
	}

	payload := afdorigins.AFDOrigin{
		Properties: &afdorigins.AFDOriginProperties{
			EnabledState:                &enabledState,
			EnforceCertificateNameCheck: pointer.To(enforceCertificateNameCheck),
			HostName:                    pointer.To(d.Get("host_name").(string)),
			HTTPPort:                    pointer.To(int64(d.Get("http_port").(int))),
			HTTPSPort:                   pointer.To(int64(d.Get("https_port").(int))),
			OriginHostHeader:            pointer.ToOrNil(d.Get("origin_host_header").(string)),
			Priority:                    pointer.To(int64(d.Get("priority").(int))),
			SharedPrivateLinkResource:   expandedPrivateLink,
			Weight:                      pointer.To(int64(d.Get("weight").(int))),
		},
	}

	if err := client.CreateCallbackThenPoll(ctx, id, payload, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceCdnFrontDoorOriginRead(d, meta)
}

func expandCdnFrontDoorOriginPrivateLinkSettings(input []any, skuName profiles.SkuName, enableCertNameCheck bool) (*afdorigins.SharedPrivateLinkResourceProperties, error) {
	if len(input) == 0 {
		// NOTE: This cannot return an empty object, the service team requires this to be set to nil else you will get the following error during creation:
		// Property 'AfdOrigin.SharedPrivateLinkResource.PrivateLink' is required but it was not set; Property 'AfdOrigin.SharedPrivateLinkResource.RequestMessage' is required but it was not set
		return nil, nil
	}

	if skuName != profiles.SkuNamePremiumAzureFrontDoor {
		return nil, fmt.Errorf("the `private_link` field can only be configured when the Frontdoor Profile is using a `Premium_AzureFrontDoor` SKU, got %q", skuName)
	}

	if !enableCertNameCheck {
		return nil, fmt.Errorf("the `private_link` field can only be configured when `certificate_name_check_enabled` is set to `true`")
	}

	// Check if this a Load Balancer Private Link or not, the Load Balancer Private Link requires
	// that you stand up your own Private Link Service, which is why I am attempting to parse a
	// Private Link Service ID here...
	config := input[0].(map[string]any)
	targetType := config["target_type"].(string)
	if _, err := privatelinkservices.ParsePrivateLinkServiceID(config["private_link_target_id"].(string)); err != nil && targetType == "" {
		// It is not a Load Balancer and the Target Type is empty, which is invalid...
		return nil, fmt.Errorf("either `private_link_target_id` or `target_type` must be specified")
	}

	return &afdorigins.SharedPrivateLinkResourceProperties{
		PrivateLink: &afdorigins.ResourceReference{
			Id: pointer.To(config["private_link_target_id"].(string)),
		},
		GroupId:             pointer.To(config["target_type"].(string)),
		PrivateLinkLocation: pointer.To(location.Normalize(config["location"].(string))),
		RequestMessage:      pointer.To(config["request_message"].(string)),
	}, nil
}
