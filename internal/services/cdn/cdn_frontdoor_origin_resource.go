// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/privatelinkservices"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontDoorOrigin() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceCdnFrontDoorOriginCreate,
		Read:   resourceCdnFrontDoorOriginRead,
		Update: resourceCdnFrontDoorOriginUpdate,
		Delete: resourceCdnFrontDoorOriginDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FrontDoorOriginID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorOriginName,
			},

			"cdn_frontdoor_origin_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorOriginGroupID,
			},

			"host_name": {
				Type: pluginsdk.TypeString,
				// HostName cannot be null or empty.
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"certificate_name_check_enabled": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"http_port": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      80,
				ValidateFunc: validation.IntBetween(1, 65535),
			},

			"https_port": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      443,
				ValidateFunc: validation.IntBetween(1, 65535),
			},

			// Must be a valid domain name, IPv4 or IPv6 IP Address
			"origin_host_header": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.Any(validation.IsIPv6Address, validation.IsIPv4Address, validation.StringIsNotEmpty),
			},

			"priority": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 5),
			},

			"private_link": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"location": commonschema.Location(),

						"private_link_target_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"request_message": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "Access request for CDN FrontDoor Private Link Origin",
							ValidateFunc: validation.StringLenBetween(1, 140),
						},

						"target_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"blob",
								"blob_secondary",
								"Gateway",
								"managedEnvironments",
								"sites",
								"web",
								"web_secondary",
							}, false),
						},
					},
				},
			},

			"weight": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      500,
				ValidateFunc: validation.IntBetween(1, 1000),
			},
		},
	}

	return resource
}

func resourceCdnFrontDoorOriginCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginsClient
	profileClient := meta.(*clients.Client).Cdn.FrontDoorProfilesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	originGroupRaw := d.Get("cdn_frontdoor_origin_group_id").(string)
	originGroup, err := parse.FrontDoorOriginGroupID(originGroupRaw)
	if err != nil {
		return err
	}

	id := parse.NewFrontDoorOriginID(originGroup.SubscriptionId, originGroup.ResourceGroup, originGroup.ProfileName, originGroup.OriginGroupName, d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, id.OriginName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_cdn_frontdoor_origin", id.ID())
	}

	// I need to get the profile SKU so I know if it is valid or not to define a private link as
	// private links are only allowed in the premium sku...
	profileId := profiles.NewProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName)

	profileResp, err := profileClient.Get(ctx, profileId)
	if err != nil {
		if response.WasNotFound(profileResp.HttpResponse) {
			return fmt.Errorf("retrieving parent %s: not found", profileId)
		}

		return fmt.Errorf("retrieving parent %s: %+v", profileId, err)
	}

	profileModel := profileResp.Model

	if profileModel == nil {
		return fmt.Errorf("profileModel is 'nil'")
	}

	if profileModel.Properties == nil {
		return fmt.Errorf("profileModel.Properties is 'nil'")
	}

	if profileModel.Sku.Name == nil {
		return fmt.Errorf("profileModel.Sku.Name' is 'nil'")
	}
	skuName := string(pointer.From(profileModel.Sku.Name))

	var enabled bool
	if !pluginsdk.IsExplicitlyNullInConfig(d, "enabled") {
		enabled = d.Get("enabled").(bool)
	}

	enableCertNameCheck := d.Get("certificate_name_check_enabled").(bool)
	props := &cdn.AFDOriginProperties{
		EnabledState:                expandEnabledBool(enabled),
		EnforceCertificateNameCheck: utils.Bool(enableCertNameCheck),
		HostName:                    utils.String(d.Get("host_name").(string)),
		HTTPPort:                    utils.Int32(int32(d.Get("http_port").(int))),
		HTTPSPort:                   utils.Int32(int32(d.Get("https_port").(int))),
		Priority:                    utils.Int32(int32(d.Get("priority").(int))),
		Weight:                      utils.Int32(int32(d.Get("weight").(int))),
	}

	if originHostHeader := d.Get("origin_host_header").(string); originHostHeader != "" {
		props.OriginHostHeader = utils.String(originHostHeader)
	}

	expanded, err := expandPrivateLinkSettings(d.Get("private_link").([]interface{}), profiles.SkuName(skuName), enableCertNameCheck)
	if err != nil {
		return err
	}
	props.SharedPrivateLinkResource = expanded

	payload := cdn.AFDOrigin{
		AFDOriginProperties: props,
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, id.OriginName, payload)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnFrontDoorOriginRead(d, meta)
}

func resourceCdnFrontDoorOriginRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorOriginID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, id.OriginName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.OriginName)
	d.Set("cdn_frontdoor_origin_group_id", parse.NewFrontDoorOriginGroupID(id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.OriginGroupName).ID())

	if props := resp.AFDOriginProperties; props != nil {
		if err := d.Set("private_link", flattenPrivateLinkSettings(props.SharedPrivateLinkResource)); err != nil {
			return fmt.Errorf("setting 'private_link': %+v", err)
		}

		d.Set("certificate_name_check_enabled", props.EnforceCertificateNameCheck)
		d.Set("enabled", flattenEnabledBool(props.EnabledState))
		d.Set("host_name", props.HostName)
		d.Set("http_port", props.HTTPPort)
		d.Set("https_port", props.HTTPSPort)
		d.Set("origin_host_header", props.OriginHostHeader)
		d.Set("priority", props.Priority)
		d.Set("weight", props.Weight)
	}

	return nil
}

func resourceCdnFrontDoorOriginUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginsClient
	workaroundClient := azuresdkhacks.NewCdnFrontDoorOriginsWorkaroundClient(client)
	profileClient := meta.(*clients.Client).Cdn.FrontDoorProfilesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorOriginID(d.Id())
	if err != nil {
		return err
	}

	params := &azuresdkhacks.AFDOriginUpdatePropertiesParameters{}

	if d.HasChange("certificate_name_check_enabled") {
		params.EnforceCertificateNameCheck = utils.Bool(d.Get("certificate_name_check_enabled").(bool))
	}

	if d.HasChange("enabled") {
		params.EnabledState = expandEnabledBool(d.Get("enabled").(bool))
	}

	if d.HasChange("host_name") {
		params.HostName = utils.String(d.Get("host_name").(string))
	}

	if d.HasChange("http_port") {
		params.HTTPPort = utils.Int32(int32(d.Get("http_port").(int)))
	}

	if d.HasChange("https_port") {
		params.HTTPSPort = utils.Int32(int32(d.Get("https_port").(int)))
	}

	// The API requires that an explicit null be passed as the 'origin_host_header' value to remove the origin host header, see issue #20617
	// Since null is a valid value, we now have to always pass the value during update else we will inadvertently clear the value, see issue #20866
	params.OriginHostHeader = nil
	if d.Get("origin_host_header").(string) != "" {
		params.OriginHostHeader = utils.String(d.Get("origin_host_header").(string))
	}

	if d.HasChange("private_link") {
		// I need to get the profile SKU so I know if it is valid or not to define a private link as
		// private links are only allowed in the premium sku...
		profileId := profiles.NewProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName)

		profileResp, err := profileClient.Get(ctx, profileId)
		if err != nil {
			if response.WasNotFound(profileResp.HttpResponse) {
				return fmt.Errorf("retrieving parent %s: not found", profileId)
			}

			return fmt.Errorf("retrieving parent %s: %+v", profileId, err)
		}

		profileModel := profileResp.Model

		if profileModel == nil {
			return fmt.Errorf("profileModel is 'nil'")
		}

		if profileModel.Sku.Name == nil {
			return fmt.Errorf("retrieving parent %s: 'profileModel.Sku.Name' was 'nil'", profileId)
		}

		enableCertNameCheck := d.Get("certificate_name_check_enabled").(bool)
		privateLinkSettings, err := expandPrivateLinkSettings(d.Get("private_link").([]interface{}), pointer.From(profileModel.Sku.Name), enableCertNameCheck)
		if err != nil {
			return err
		}

		params.SharedPrivateLinkResource = privateLinkSettings
	}

	if d.HasChange("priority") {
		params.Priority = utils.Int32(int32(d.Get("priority").(int)))
	}

	if d.HasChange("weight") {
		params.Weight = utils.Int32(int32(d.Get("weight").(int)))
	}

	payload := &azuresdkhacks.AFDOriginUpdateParameters{
		AFDOriginUpdatePropertiesParameters: params,
	}

	future, err := workaroundClient.Update(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, id.OriginName, *payload)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	return resourceCdnFrontDoorOriginRead(d, meta)
}

func resourceCdnFrontDoorOriginDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorOriginID(d.Id())
	if err != nil {
		return err
	}

	// @tombuildsstuff: JC/WS to dig into if we need to conditionally remove the Private Link
	// via an Update before deletion - presumably we'd also need a Lock on the private link resource

	/*
		original:
			// TODO: Check to see if there is a Load Balancer Private Link connected,
			// if so disconnect the Private Link association with the Frontdoor Origin
			// else the destroy will fail because the Private Link Service has an active
			// Private Link Endpoint connection...

			// It looks like Frontdoor does remove the Private link, I just need to poll here until it is removed...
			// Investigate this further...
			// WS: There is a bug in the service code, for only the load balancer scenario, the private link connection is not removed until the
			// origin is totally destroyed. The workaround for this issue is to put a depends_on the private link service to the origin so the origin
			// will be deleted first before the private link service is destroyed.
	*/

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, id.OriginName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandPrivateLinkSettings(input []interface{}, skuName profiles.SkuName, enableCertNameCheck bool) (*cdn.SharedPrivateLinkResourceProperties, error) {
	if len(input) == 0 {
		// NOTE: This cannot return an empty object, the service team requires this to be set to nil else you will get the following error during creation:
		// Property 'AfdOrigin.SharedPrivateLinkResource.PrivateLink' is required but it was not set; Property 'AfdOrigin.SharedPrivateLinkResource.RequestMessage' is required but it was not set
		return nil, nil
	}

	if skuName != profiles.SkuNamePremiumAzureFrontDoor {
		return nil, fmt.Errorf("the 'private_link' field can only be configured when the Frontdoor Profile is using a 'Premium_AzureFrontDoor' SKU, got %q", skuName)
	}

	if !enableCertNameCheck {
		return nil, fmt.Errorf("the 'private_link' field can only be configured when 'certificate_name_check_enabled' is set to 'true'")
	}

	// Check if this a Load Balancer Private Link or not, the Load Balancer Private Link requires
	// that you stand up your own Private Link Service, which is why I am attempting to parse a
	// Private Link Service ID here...
	settings := input[0].(map[string]interface{})
	targetType := settings["target_type"].(string)
	_, err := privatelinkservices.ParsePrivateLinkServiceID(settings["private_link_target_id"].(string))
	if err != nil && targetType == "" {
		// It is not a Load Balancer and the Target Type is empty, which is invalid...
		return nil, fmt.Errorf("either 'private_link' or 'target_type' must be specified")
	}

	config := input[0].(map[string]interface{})

	resourceId := config["private_link_target_id"].(string)
	location := location.Normalize(config["location"].(string))
	groupId := config["target_type"].(string)
	requestMessage := config["request_message"].(string)

	return &cdn.SharedPrivateLinkResourceProperties{
		PrivateLink: &cdn.ResourceReference{
			ID: utils.String(resourceId),
		},
		GroupID:             utils.String(groupId),
		PrivateLinkLocation: utils.String(location),
		RequestMessage:      utils.String(requestMessage),
	}, nil
}

func flattenPrivateLinkSettings(input *cdn.SharedPrivateLinkResourceProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	privateLinkTargetId := ""
	if input.PrivateLink != nil && input.PrivateLink.ID != nil {
		privateLinkTargetId = *input.PrivateLink.ID
	}

	requestMessage := ""
	if input.RequestMessage != nil {
		requestMessage = *input.RequestMessage
	}

	targetType := ""
	if input.GroupID != nil {
		targetType = *input.GroupID
	}

	return []interface{}{
		map[string]interface{}{
			"location":               location.NormalizeNilable(input.PrivateLinkLocation),
			"private_link_target_id": privateLinkTargetId,
			"request_message":        requestMessage,
			"target_type":            targetType,
		},
	}
}
