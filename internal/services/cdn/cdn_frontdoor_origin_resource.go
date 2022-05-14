package cdn

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	privateLinkServiceParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontdoorOrigin() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontdoorOriginCreate,
		Read:   resourceCdnFrontdoorOriginRead,
		Update: resourceCdnFrontdoorOriginUpdate,
		Delete: resourceCdnFrontdoorOriginDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FrontdoorOriginID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CdnFrontdoorOriginName,
			},

			"cdn_frontdoor_origin_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontdoorOriginGroupID,
			},

			// HostName cannot be null or empty.;
			"host_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// NOTE: In swagger this is the Azure Origin... this is not currently used and is not what I thought is was
			// This value can be: storage (Azure Blobs), Storage (Classic), Storage (Static Website), Cloud Service,
			// App Service, Static Web App, API Management, Application Gateway, Public IP Address or a Traffic Manager.
			// Currently, this functionality is being exposed via the origin_host_header field.

			// "cdn_frontdoor_origin_id": {
			// 	Type:     pluginsdk.TypeString,
			// 	Optional: true,
			// },

			"health_probes_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"certificate_name_check_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
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

			"weight": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      500,
				ValidateFunc: validation.IntBetween(1, 1000),
			},

			"private_link": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						// TODO: does this want to be normalized?
						// WS: Fixed
						"location": azure.SchemaLocation(),

						"private_link_target_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"request_message": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "Access request for CDN Frontdoor Private Link Origin",
							ValidateFunc: validation.StringLenBetween(1, 140),
						},

						"target_type": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(cdnFrontdoorPrivateLinkTargetTypes(), false),
						},
					},
				},
			},
		},
	}
}

func resourceCdnFrontdoorOriginCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginsClient
	profileClient := meta.(*clients.Client).Cdn.FrontDoorProfileClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	originGroupId, err := parse.FrontdoorOriginGroupID(d.Get("cdn_frontdoor_origin_group_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewFrontdoorOriginID(originGroupId.SubscriptionId, originGroupId.ResourceGroup, originGroupId.ProfileName, originGroupId.OriginGroupName, d.Get("name").(string))
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
	profileId := parse.NewFrontdoorProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName)

	profile, err := profileClient.Get(ctx, profileId.ResourceGroup, profileId.ProfileName)
	if err != nil {
		if utils.ResponseWasNotFound(profile.Response) {
			return fmt.Errorf("retrieving parent %s: not found", profileId)
		}

		return fmt.Errorf("retrieving parent %s: %+v", profileId, err)
	}

	if profile.Sku == nil {
		return fmt.Errorf("retrieving parent %s: `sku` was nil", profileId)
	}
	skuName := profile.Sku.Name

	enableCertNameCheck := d.Get("certificate_name_check_enabled").(bool)
	props := &cdn.AFDOriginProperties{
		// AzureOrigin is currently not used, service team asked me to temporarily remove it from the resource
		// AzureOrigin:                 expandResourceReference(d.Get("cdn_frontdoor_origin_id").(string)),
		EnabledState:                convertCdnFrontdoorBoolToEnabledState(d.Get("health_probes_enabled").(bool)),
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

	expanded, err := expandPrivateLinkSettings(d.Get("private_link").([]interface{}), skuName, enableCertNameCheck)
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
	return resourceCdnFrontdoorOriginRead(d, meta)
}

func resourceCdnFrontdoorOriginRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorOriginID(d.Id())
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
	d.Set("cdn_frontdoor_origin_group_id", parse.NewFrontdoorOriginGroupID(id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.OriginGroupName).ID())

	if props := resp.AFDOriginProperties; props != nil {

		// if err := d.Set("cdn_frontdoor_origin_id", flattenResourceReference(props.AzureOrigin)); err != nil {
		// 	return fmt.Errorf("setting `cdn_frontdoor_origin_id`: %+v", err)
		// }

		if props.SharedPrivateLinkResource != nil {
			d.Set("private_link", flattenPrivateLinkSettings(props.SharedPrivateLinkResource))
		}

		d.Set("health_probes_enabled", convertCdnFrontdoorEnabledStateToBool(&props.EnabledState))
		d.Set("certificate_name_check_enabled", props.EnforceCertificateNameCheck)
		d.Set("host_name", props.HostName)
		d.Set("http_port", props.HTTPPort)
		d.Set("https_port", props.HTTPSPort)
		d.Set("origin_host_header", props.OriginHostHeader)
		d.Set("priority", props.Priority)
		d.Set("weight", props.Weight)
	}

	return nil
}

func resourceCdnFrontdoorOriginUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginsClient
	profileClient := meta.(*clients.Client).Cdn.FrontDoorProfileClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorOriginID(d.Id())
	if err != nil {
		return err
	}

	// I need to get the profile SKU so I know if it is valid or not to define a private link as
	// private links are only allowed in the premium sku...
	profileId := parse.NewFrontdoorProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName)
	profile, err := profileClient.Get(ctx, profileId.ResourceGroup, profileId.ProfileName)
	if err != nil {
		if utils.ResponseWasNotFound(profile.Response) {
			return fmt.Errorf("retrieving parent %s: not found", profileId)
		}

		return fmt.Errorf("retrieving parent %s: %+v", profileId, err)
	}
	if profile.Sku == nil {
		return fmt.Errorf("retrieving parent %s: `sku` was nil", profileId)
	}
	skuName := profile.Sku.Name

	originHostHeader := d.Get("origin_host_header").(string)
	enableCertNameCheck := d.Get("certificate_name_check_enabled").(bool)

	props := cdn.AFDOriginUpdateParameters{
		AFDOriginUpdatePropertiesParameters: &cdn.AFDOriginUpdatePropertiesParameters{
			// AzureOrigin:                 expandResourceReference(d.Get("cdn_frontdoor_origin_id").(string)),
			EnabledState:                convertCdnFrontdoorBoolToEnabledState(d.Get("health_probes_enabled").(bool)),
			EnforceCertificateNameCheck: utils.Bool(enableCertNameCheck),
			HostName:                    utils.String(d.Get("host_name").(string)),
			HTTPPort:                    utils.Int32(int32(d.Get("http_port").(int))),
			HTTPSPort:                   utils.Int32(int32(d.Get("https_port").(int))),
			Priority:                    utils.Int32(int32(d.Get("priority").(int))),
			Weight:                      utils.Int32(int32(d.Get("weight").(int))),
		},
	}

	if d.HasChange("private_link") {
		privateLinkSettings, err := expandPrivateLinkSettings(d.Get("private_link").([]interface{}), skuName, enableCertNameCheck)
		if err != nil {
			return err
		}
		props.SharedPrivateLinkResource = privateLinkSettings
	}

	if originHostHeader != "" {
		props.OriginHostHeader = utils.String(originHostHeader)
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, id.OriginName, props)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	return resourceCdnFrontdoorOriginRead(d, meta)
}

func resourceCdnFrontdoorOriginDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// TODO: Check to see if there is a Load Balancer Private Link connected,
	// if so disconnect the Private Link association with the Frontdoor Origin
	// else the destroy will fail because the Private Link Service has an active
	// Private Link Endpoint connection...

	// It looks like Frontdoor does remove the Private link, I just need to poll here until it is removed...
	// Investigate this further...
	// TODO: what's the deal with this? ^
	// WS: There is a bug in the service code, for only the load balancer scenario, the private link connection is not removed until the
	// origin is totally destroyed. The workaround for this issue is to put a depends_on the private link service to the origin so the origin
	// will be deleted first before the private link service is destroyed.
	id, err := parse.FrontdoorOriginID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, id.OriginName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandPrivateLinkSettings(input []interface{}, skuName cdn.SkuName, enableCertNameCheck bool) (*cdn.SharedPrivateLinkResourceProperties, error) {
	if len(input) == 0 {
		// TODO: Should this return an empty object?
		// WS: This cannot return an empty object, the service team requires this to be set to nil else you will get the following error during creation:
		// Property 'AfdOrigin.SharedPrivateLinkResource.PrivateLink' is required but it was not set; Property 'AfdOrigin.SharedPrivateLinkResource.RequestMessage' is required but it was not set
		return nil, nil
	}

	if skuName != cdn.SkuNamePremiumAzureFrontDoor {
		return nil, fmt.Errorf("the `private_link` field can only be configured when the Frontdoor Profile is using a %q SKU, got %q", cdn.SkuNamePremiumAzureFrontDoor, skuName)
	}

	if !enableCertNameCheck {
		return nil, fmt.Errorf("the `private_link` field can only be configured when `certificate_name_check_enabled` is set to `true`")
	}

	// Check if this a Load Balancer Private Link or not, the Load Balancer Private Link requires
	// that you stand up your own Private Link Service, which is why I am attempting to parse a
	// Private Link Service ID here...
	settings := input[0].(map[string]interface{})
	targetType := settings["target_type"].(string)
	_, err := privateLinkServiceParse.PrivateLinkServiceID(settings["private_link_target_id"].(string))
	if err != nil && targetType == "" {
		// It is not a Load Balancer and the Target Type is empty, which is invalid...
		return nil, fmt.Errorf("the %[1]q block requires that you define the %[2]q field if the %[1]q is not a Load Balancer, expected %[3]s but got %[4]q", "private_link", "target_type", strings.Join(cdnFrontdoorPrivateLinkTargetTypes(), ", "), targetType)
	}

	config := input[0].(map[string]interface{})

	resourceId := config["private_link_target_id"].(string)
	location := azure.NormalizeLocation(config["location"].(string))
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

	location := ""
	if input.PrivateLinkLocation != nil {
		location = *input.PrivateLinkLocation
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
			"location":               location,
			"private_link_target_id": privateLinkTargetId,
			"request_message":        requestMessage,
			"target_type":            targetType,
		},
	}
}
