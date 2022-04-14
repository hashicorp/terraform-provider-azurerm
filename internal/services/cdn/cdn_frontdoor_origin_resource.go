package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
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
				ValidateFunc: IsValidDomain,
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

						"approval_message": { // A custom message to be included in the approval request to connect to the Private Link.
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"location": { // The location of the Private Link resource.
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"private_link_service_id": { // The Resource Id of the Private Link resource.
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"cdn_frontdoor_origin_group_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCdnFrontdoorOriginCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginsClient
	profileClient := meta.(*clients.Client).Cdn.FrontdoorProfileClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	originGroupId, err := parse.FrontdoorOriginGroupID(d.Get("cdn_frontdoor_origin_group_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewFrontdoorOriginID(originGroupId.SubscriptionId, originGroupId.ResourceGroup, originGroupId.ProfileName, originGroupId.OriginGroupName, d.Get("name").(string))

	// I need to get the profile SKU so I know if it is valid or not to define a private link as
	// private links are only allowed in the premium sku...
	profileId := parse.NewFrontdoorProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName)

	profile, err := profileClient.Get(ctx, profileId.ResourceGroup, profileId.ProfileName)
	if err != nil {
		if utils.ResponseWasNotFound(profile.Response) {
			return fmt.Errorf("%s does not exist: %+v", profileId, err)
		}

		return fmt.Errorf("retrieving SKU information from %s: %+v", profileId, err)
	}

	var sku string

	if profileSku := profile.Sku; profileSku != nil {
		sku = string(profileSku.Name)
	} else {
		return fmt.Errorf("retrieving SKU information from %s: %+v", profileId, err)
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, id.OriginName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_cdn_frontdoor_origin", id.ID())
		}
	}

	originHostHeader := d.Get("origin_host_header").(string)
	enableCertNameCheck := d.Get("certificate_name_check_enabled").(bool)

	props := track1.AFDOrigin{
		AFDOriginProperties: &track1.AFDOriginProperties{
			// AzureOrigin:                 expandResourceReference(d.Get("cdn_frontdoor_origin_id").(string)),
			EnabledState:                ConvertBoolToEnabledState(d.Get("health_probes_enabled").(bool)),
			EnforceCertificateNameCheck: utils.Bool(enableCertNameCheck),
			HostName:                    utils.String(d.Get("host_name").(string)),
			HTTPPort:                    utils.Int32(int32(d.Get("http_port").(int))),
			HTTPSPort:                   utils.Int32(int32(d.Get("https_port").(int))),
			Priority:                    utils.Int32(int32(d.Get("priority").(int))),
			Weight:                      utils.Int32(int32(d.Get("weight").(int))),
		},
	}

	if originHostHeader != "" {
		props.OriginHostHeader = utils.String(originHostHeader)
	}

	privateLinkSettings := d.Get("private_link").([]interface{})
	if len(privateLinkSettings) > 0 {
		if sku == string(track1.SkuNamePremiumAzureFrontDoor) {
			if !enableCertNameCheck {
				return fmt.Errorf("%q requires that the %q field be set to %q, got %q", "private_link", "certificate_name_check_enabled", "true", "false")
			} else {
				props.SharedPrivateLinkResource = expandPrivateLinkSettings(privateLinkSettings)
			}
		} else {
			return fmt.Errorf("the %q field is only valid if the %q SKU is set to %q, got %q", "private_link", "Frontdoor Profile", track1.SkuNamePremiumAzureFrontDoor, sku)
		}
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, id.OriginName, props)
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
		return fmt.Errorf("retrieving %s: %+v", id, err)
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

		d.Set("health_probes_enabled", ConvertEnabledStateToBool(&props.EnabledState))
		d.Set("certificate_name_check_enabled", props.EnforceCertificateNameCheck)
		d.Set("host_name", props.HostName)
		d.Set("http_port", props.HTTPPort)
		d.Set("https_port", props.HTTPSPort)
		d.Set("cdn_frontdoor_origin_group_name", props.OriginGroupName)
		d.Set("origin_host_header", props.OriginHostHeader)
		d.Set("priority", props.Priority)
		d.Set("weight", props.Weight)
	}

	return nil
}

func resourceCdnFrontdoorOriginUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginsClient
	profileClient := meta.(*clients.Client).Cdn.FrontdoorProfileClient
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
			return fmt.Errorf("%s does not exist: %+v", profileId, err)
		}

		return fmt.Errorf("retrieving SKU information from %s: %+v", profileId, err)
	}

	var sku string

	if profileSku := profile.Sku; profileSku != nil {
		sku = string(profileSku.Name)
	} else {
		return fmt.Errorf("retrieving SKU information from %s: %+v", profileId, err)
	}

	originHostHeader := d.Get("origin_host_header").(string)
	enableCertNameCheck := d.Get("certificate_name_check_enabled").(bool)

	props := track1.AFDOriginUpdateParameters{
		AFDOriginUpdatePropertiesParameters: &track1.AFDOriginUpdatePropertiesParameters{
			// AzureOrigin:                 expandResourceReference(d.Get("cdn_frontdoor_origin_id").(string)),
			EnabledState:                ConvertBoolToEnabledState(d.Get("health_probes_enabled").(bool)),
			EnforceCertificateNameCheck: utils.Bool(enableCertNameCheck),
			HostName:                    utils.String(d.Get("host_name").(string)),
			HTTPPort:                    utils.Int32(int32(d.Get("http_port").(int))),
			HTTPSPort:                   utils.Int32(int32(d.Get("https_port").(int))),
			Priority:                    utils.Int32(int32(d.Get("priority").(int))),
			Weight:                      utils.Int32(int32(d.Get("weight").(int))),
		},
	}

	if d.HasChange("private_link") {
		privateLinkSettings := d.Get("private_link").([]interface{})
		if len(privateLinkSettings) > 0 {
			if sku == string(track1.SkuNamePremiumAzureFrontDoor) {
				if !enableCertNameCheck {
					return fmt.Errorf("%q requires that the %q field be set to %q, got %q", "private_link", "certificate_name_check_enabled", "true", "false")
				} else {
					props.SharedPrivateLinkResource = expandPrivateLinkSettings(privateLinkSettings)
				}
			} else {
				return fmt.Errorf("the %q field is only valid if the %q SKU is set to %q, got %q", "private_link", "Frontdoor Profile", track1.SkuNamePremiumAzureFrontDoor, sku)
			}
		}
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

func expandPrivateLinkSettings(input []interface{}) *track1.SharedPrivateLinkResourceProperties {
	if len(input) == 0 {
		return &track1.SharedPrivateLinkResourceProperties{}
	}

	config := input[0].(map[string]interface{})

	resourceId := config["private_link_service_id"].(string)
	location := config["location"].(string)
	approvalMessage := config["approval_message"].(string)

	privateLinkResource := track1.SharedPrivateLinkResourceProperties{
		PrivateLink: &track1.ResourceReference{
			ID: &resourceId,
		},
		PrivateLinkLocation: &location,
		RequestMessage:      &approvalMessage,
	}

	return &privateLinkResource
}

func flattenPrivateLinkSettings(input *track1.SharedPrivateLinkResourceProperties) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	// resourceId := config["resource_id"].(string)
	// location := config["location"].(string)
	// approvalMessage := config["approval_message"].(string)

	// privateLinkResource := track1.SharedPrivateLinkResourceProperties{
	// 	PrivateLink: &track1.ResourceReference{
	// 		ID: &resourceId,
	// 	},
	// 	PrivateLinkLocation: &location,
	// 	RequestMessage:      &approvalMessage,
	// }

	return results
}
