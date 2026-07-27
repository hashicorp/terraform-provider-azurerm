// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigingroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigins"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/privatelinkservices"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorOrigin() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceCdnFrontDoorOriginCreate,
		Read:   resourceCdnFrontDoorOriginRead,
		Update: resourceCdnFrontDoorOriginUpdate,
		Delete: resourceCdnFrontDoorOriginDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(4 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(4 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(6 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := afdorigins.ParseOriginGroupOriginID(id)
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
				ValidateFunc: afdorigins.ValidateOriginGroupID,
			},

			"host_name": {
				Type:         pluginsdk.TypeString,
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
							ValidateFunc: privatelinkservices.ValidatePrivateLinkServiceID,
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
			return tf.ImportAsExistsError("azurerm_cdn_frontdoor_origin", id.ID())
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
	expandedPrivateLink, err := expandCdnFrontDoorOriginPrivateLinkSettings(d.Get("private_link").([]interface{}), pointer.From(profileResp.Model.Sku.Name), enforceCertificateNameCheck)
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

	if err := client.CreateCallbackThenPoll(ctx, id, payload, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnFrontDoorOriginRead(d, meta)
}

func resourceCdnFrontDoorOriginRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginsClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afdorigins.ParseOriginGroupOriginID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.OriginName)
	d.Set("cdn_frontdoor_origin_group_id", afdorigins.NewOriginGroupID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.OriginGroupName).ID())

	if model := resp.Model; model != nil {
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

	return nil
}

func resourceCdnFrontDoorOriginUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("retreiving %s: model was nil", profileId)
		}

		if profileModel.Sku.Name == nil {
			return fmt.Errorf("retrieving %s: sku was nil", profileId)
		}

		privateLinkSettings, err := expandCdnFrontDoorOriginPrivateLinkSettings(d.Get("private_link").([]interface{}), pointer.From(profileModel.Sku.Name), d.Get("certificate_name_check_enabled").(bool))
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

func resourceCdnFrontDoorOriginDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginsClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afdorigins.ParseOriginGroupOriginID(d.Id())
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

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandCdnFrontDoorOriginPrivateLinkSettings(input []interface{}, skuName profiles.SkuName, enableCertNameCheck bool) (*afdorigins.SharedPrivateLinkResourceProperties, error) {
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
	config := input[0].(map[string]interface{})
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

func flattenCdnFrontDoorOriginPrivateLinkSettings(input *afdorigins.SharedPrivateLinkResourceProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	privateLinkTargetId := ""
	if input.PrivateLink != nil && input.PrivateLink.Id != nil {
		privateLinkTargetId = *input.PrivateLink.Id
	}

	return []interface{}{
		map[string]interface{}{
			"location":               location.NormalizeNilable(input.PrivateLinkLocation),
			"private_link_target_id": privateLinkTargetId,
			"request_message":        pointer.From(input.RequestMessage),
			"target_type":            pointer.From(input.GroupId),
		},
	}
}
