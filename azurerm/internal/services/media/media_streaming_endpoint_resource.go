package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2021-05-01/media"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMediaStreamingEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMediaStreamingEndpointCreate,
		Read:   resourceMediaStreamingEndpointRead,
		Update: resourceMediaStreamingEndpointUpdate,
		Delete: resourceMediaStreamingEndpointDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.StreamingEndpointID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StreamingEndpointName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"media_services_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName,
			},

			"auto_start_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},

			"location": azure.SchemaLocation(),

			"scale_units": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 10),
			},

			"access_control": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						//lintignore:XS003
						"akamai_signature_header_authentication_key": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"base64_key": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsBase64,
									},
									"expiration": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsRFC3339Time,
									},
									"identifier": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
							AtLeastOneOf: []string{"access_control.0.akamai_signature_header_authentication_key", "access_control.0.ip_allow"},
						},
						//lintignore:XS003
						"ip_allow": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"address": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"name": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"subnet_prefix_length": {
										Type:     pluginsdk.TypeInt,
										Optional: true,
									},
								},
							},
							AtLeastOneOf: []string{"access_control.0.akamai_signature_header_authentication_key", "access_control.0.ip_allow"},
						},
					},
				},
			},

			"cdn_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"cdn_profile": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-A-Za-z0-9]{1,120}$"),
					"CDN profile must be 1 - 120 characters long, can contain only letters, numbers, and hyphens.",
				),
			},

			"cdn_provider": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"StandardVerizon", "PremiumVerizon", "StandardAkamai",
				}, false),
			},

			"cross_site_access_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"client_access_policy": {
							Type:         pluginsdk.TypeString,
							Computed:     true,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{"cross_site_access_policy.0.client_access_policy", "cross_site_access_policy.0.cross_domain_policy"},
						},

						"cross_domain_policy": {
							Type:         pluginsdk.TypeString,
							Computed:     true,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{"cross_site_access_policy.0.client_access_policy", "cross_site_access_policy.0.cross_domain_policy"},
						},
					},
				},
			},

			"custom_host_names": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"max_cache_age_seconds": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 2147483647),
			},

			"host_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceMediaStreamingEndpointCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.StreamingEndpointsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	streamingEndpointName := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("media_services_account_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	scaleUnits := d.Get("scale_units").(int)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resourceId := parse.NewStreamingEndpointID(subscriptionId, d.Get("resource_group_name").(string), d.Get("media_services_account_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.MediaserviceName, resourceId.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of %s: %+v", resourceId, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_media_streaming_endpoint", resourceId.ID())
	}

	parameters := media.StreamingEndpoint{
		StreamingEndpointProperties: &media.StreamingEndpointProperties{
			ScaleUnits: utils.Int32(int32(scaleUnits)),
		},
		Location: utils.String(location),
	}

	autoStart := utils.Bool(false)
	if _, ok := d.GetOk("auto_start_enabled"); ok {
		autoStart = utils.Bool(d.Get("auto_start_enabled").(bool))
	}
	if _, ok := d.GetOk("access_control"); ok {
		accessControl, err := expandAccessControl(d)
		if err != nil {
			return err
		}
		parameters.StreamingEndpointProperties.AccessControl = accessControl
	}
	if cdnEnabled, ok := d.GetOk("cdn_enabled"); ok {
		parameters.StreamingEndpointProperties.CdnEnabled = utils.Bool(cdnEnabled.(bool))
	}

	if cdnProfile, ok := d.GetOk("cdn_profile"); ok {
		parameters.StreamingEndpointProperties.CdnProfile = utils.String(cdnProfile.(string))
	}

	if cdnProvider, ok := d.GetOk("cdn_provider"); ok {
		parameters.StreamingEndpointProperties.CdnProvider = utils.String(cdnProvider.(string))
	}

	if crossSite, ok := d.GetOk("cross_site_access_policy"); ok {
		parameters.StreamingEndpointProperties.CrossSiteAccessPolicies = expandCrossSiteAccessPolicies(crossSite.([]interface{}))
	}

	if _, ok := d.GetOk("custom_host_names"); ok {
		customHostNames := d.Get("custom_host_names").([]interface{})
		parameters.StreamingEndpointProperties.CustomHostNames = utils.ExpandStringSlice(customHostNames)
	}

	if description, ok := d.GetOk("description"); ok {
		parameters.StreamingEndpointProperties.Description = utils.String(description.(string))
	}

	if maxCacheAge, ok := d.GetOk("max_cache_age_seconds"); ok {
		parameters.StreamingEndpointProperties.MaxCacheAge = utils.Int64(int64(maxCacheAge.(int)))
	}

	future, err := client.Create(ctx, resourceGroup, accountName, streamingEndpointName, parameters, autoStart)
	if err != nil {
		return fmt.Errorf("Error creating Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", streamingEndpointName, accountName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", streamingEndpointName, accountName, resourceGroup, err)
	}

	d.SetId(resourceId.ID())

	return resourceMediaStreamingEndpointRead(d, meta)
}

func resourceMediaStreamingEndpointUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.StreamingEndpointsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StreamingEndpointID(d.Id())
	if err != nil {
		return err
	}
	location := azure.NormalizeLocation(d.Get("location").(string))
	scaleUnits := d.Get("scale_units").(int)

	parameters := media.StreamingEndpoint{
		StreamingEndpointProperties: &media.StreamingEndpointProperties{
			ScaleUnits: utils.Int32(int32(scaleUnits)),
		},
		Location: utils.String(location),
	}

	if d.HasChange("scale_units") {
		scaleParamaters := media.StreamingEntityScaleUnit{
			ScaleUnit: utils.Int32(int32(scaleUnits)),
		}

		future, err := client.Scale(ctx, id.ResourceGroup, id.MediaserviceName, id.Name, scaleParamaters)
		if err != nil {
			return fmt.Errorf("Error scaling units in Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for scaling of Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
		}
	}

	if _, ok := d.GetOk("access_control"); ok {
		accessControl, err := expandAccessControl(d)
		if err != nil {
			return err
		}
		parameters.StreamingEndpointProperties.AccessControl = accessControl
	}

	if cdnEnabled, ok := d.GetOk("cdn_enabled"); ok {
		parameters.StreamingEndpointProperties.CdnEnabled = utils.Bool(cdnEnabled.(bool))
	}

	if cdnProfile, ok := d.GetOk("cdn_profile"); ok {
		parameters.StreamingEndpointProperties.CdnProfile = utils.String(cdnProfile.(string))
	}

	if cdnProvider, ok := d.GetOk("cdn_provider"); ok {
		parameters.StreamingEndpointProperties.CdnProvider = utils.String(cdnProvider.(string))
	}

	if crossSitePolicies, ok := d.GetOk("cross_site_access_policy"); ok {
		parameters.StreamingEndpointProperties.CrossSiteAccessPolicies = expandCrossSiteAccessPolicies(crossSitePolicies.([]interface{}))
	}

	if _, ok := d.GetOk("custom_host_names"); ok {
		customHostNames := d.Get("custom_host_names").([]interface{})
		parameters.StreamingEndpointProperties.CustomHostNames = utils.ExpandStringSlice(customHostNames)
	}

	if description, ok := d.GetOk("description"); ok {
		parameters.StreamingEndpointProperties.Description = utils.String(description.(string))
	}

	if maxCacheAge, ok := d.GetOk("max_cache_age_seconds"); ok {
		parameters.StreamingEndpointProperties.MaxCacheAge = utils.Int64(int64(maxCacheAge.(int)))
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.MediaserviceName, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	return resourceMediaStreamingEndpointRead(d, meta)
}

func resourceMediaStreamingEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.StreamingEndpointsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StreamingEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Streaming Endpoint %q was not found in Media Services Account %q and Resource Group %q - removing from state", id.Name, id.MediaserviceName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("media_services_account_name", id.MediaserviceName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.StreamingEndpointProperties; props != nil {
		if scaleUnits := props.ScaleUnits; scaleUnits != nil {
			d.Set("scale_units", scaleUnits)
		}

		accessControl := flattenAccessControl(props.AccessControl)
		if err := d.Set("access_control", accessControl); err != nil {
			return fmt.Errorf("Error flattening `access_control`: %s", err)
		}

		d.Set("cdn_enabled", props.CdnEnabled)
		d.Set("cdn_profile", props.CdnProfile)
		d.Set("cdn_provider", props.CdnProvider)
		d.Set("host_name", props.HostName)

		crossSiteAccessPolicies := flattenCrossSiteAccessPolicies(resp.CrossSiteAccessPolicies)
		if err := d.Set("cross_site_access_policy", crossSiteAccessPolicies); err != nil {
			return fmt.Errorf("Error flattening `cross_site_access_policy`: %s", err)
		}

		d.Set("custom_host_names", props.CustomHostNames)
		d.Set("description", props.Description)

		maxCacheAge := 0
		if props.MaxCacheAge != nil {
			maxCacheAge = int(*props.MaxCacheAge)
		}
		d.Set("max_cache_age_seconds", maxCacheAge)
	}

	return nil
}

func resourceMediaStreamingEndpointDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.StreamingEndpointsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StreamingEndpointID(d.Id())
	if err != nil {
		return err
	}

	// Stop Streaming Endpoint before we attempt to delete it.
	resp, err := client.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		return fmt.Errorf("reading %s: %+v", id, err)
	}
	if props := resp.StreamingEndpointProperties; props != nil {
		if props.ResourceState == media.StreamingEndpointResourceStateRunning {
			stopFuture, err := client.Stop(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
			if err != nil {
				return fmt.Errorf("stopping %s: %+v", id, err)
			}

			if err = stopFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for %s to stop: %+v", id, err)
			}
		}
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	return nil
}

func expandAccessControl(d *pluginsdk.ResourceData) (*media.StreamingEndpointAccessControl, error) {
	accessControls := d.Get("access_control").([]interface{})
	if len(accessControls) == 0 {
		return nil, nil
	}
	accessControlResult := new(media.StreamingEndpointAccessControl)
	accessControl := accessControls[0].(map[string]interface{})
	// Get IP information
	if ipAllowsList := accessControl["ip_allow"].([]interface{}); len(ipAllowsList) > 0 {
		ipRanges := make([]media.IPRange, 0)
		for _, ipAllow := range ipAllowsList {
			if ipAllow == nil {
				continue
			}
			allow := ipAllow.(map[string]interface{})
			address := allow["address"].(string)
			name := allow["name"].(string)

			ipRange := media.IPRange{
				Name:    utils.String(name),
				Address: utils.String(address),
			}
			subnetPrefixLengthRaw := allow["subnet_prefix_length"]
			if subnetPrefixLengthRaw != "" {
				ipRange.SubnetPrefixLength = utils.Int32(int32(subnetPrefixLengthRaw.(int)))
			}
			ipRanges = append(ipRanges, ipRange)
		}
		accessControlResult.IP = &media.IPAccessControl{
			Allow: &ipRanges,
		}
	}
	// Get Akamai information
	if akamaiSignatureKeyList := accessControl["akamai_signature_header_authentication_key"].([]interface{}); len(akamaiSignatureKeyList) > 0 {
		akamaiSignatureHeaderAuthenticationKeyList := make([]media.AkamaiSignatureHeaderAuthenticationKey, 0)
		for _, akamaiSignatureKey := range akamaiSignatureKeyList {
			if akamaiSignatureKey == nil {
				continue
			}
			akamaiKey := akamaiSignatureKey.(map[string]interface{})
			base64Key := akamaiKey["base64_key"].(string)
			expirationRaw := akamaiKey["expiration"].(string)
			identifier := akamaiKey["identifier"].(string)

			akamaiSignatureHeaderAuthenticationKey := media.AkamaiSignatureHeaderAuthenticationKey{
				Base64Key:  utils.String(base64Key),
				Identifier: utils.String(identifier),
			}
			if expirationRaw != "" {
				expiration, err := date.ParseTime(time.RFC3339, expirationRaw)
				if err != nil {
					return nil, err
				}
				akamaiSignatureHeaderAuthenticationKey.Expiration = &date.Time{
					Time: expiration,
				}
			}
			akamaiSignatureHeaderAuthenticationKeyList = append(akamaiSignatureHeaderAuthenticationKeyList, akamaiSignatureHeaderAuthenticationKey)
		}
		accessControlResult.Akamai = &media.AkamaiAccessControl{
			AkamaiSignatureHeaderAuthenticationKeyList: &akamaiSignatureHeaderAuthenticationKeyList,
		}
	}

	return accessControlResult, nil
}

func flattenAccessControl(input *media.StreamingEndpointAccessControl) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	ipAllowRules := make([]interface{}, 0)
	if input.IP != nil && input.IP.Allow != nil {
		for _, v := range *input.IP.Allow {
			name := ""
			if v.Name != nil {
				name = *v.Name
			}

			address := ""
			if v.Address != nil {
				address = *v.Address
			}

			var subnetPrefixLength int32
			if v.SubnetPrefixLength != nil {
				subnetPrefixLength = *v.SubnetPrefixLength
			}

			ipAllowRules = append(ipAllowRules, map[string]interface{}{
				"name":                 name,
				"address":              address,
				"subnet_prefix_length": subnetPrefixLength,
			})
		}
	}

	akamaiRules := make([]interface{}, 0)
	if input.Akamai != nil && input.Akamai.AkamaiSignatureHeaderAuthenticationKeyList != nil {
		for _, v := range *input.Akamai.AkamaiSignatureHeaderAuthenticationKeyList {
			base64Key := ""
			if v.Base64Key != nil {
				base64Key = *v.Base64Key
			}

			expiration := ""
			if v.Expiration != nil {
				expiration = v.Expiration.Format(time.RFC3339)
			}

			identifier := ""
			if v.Identifier != nil {
				identifier = *v.Identifier
			}

			akamaiRules = append(akamaiRules, map[string]interface{}{
				"base64_key": base64Key,
				"expiration": expiration,
				"identifier": identifier,
			})
		}
	}

	return []interface{}{
		map[string]interface{}{
			"akamai_signature_header_authentication_key": akamaiRules,
			"ip_allow": ipAllowRules,
		},
	}
}

func expandCrossSiteAccessPolicies(input []interface{}) *media.CrossSiteAccessPolicies {
	if len(input) == 0 {
		return nil
	}

	crossSiteAccessPolicy := input[0].(map[string]interface{})
	clientAccessPolicy := crossSiteAccessPolicy["client_access_policy"].(string)
	crossDomainPolicy := crossSiteAccessPolicy["cross_domain_policy"].(string)
	return &media.CrossSiteAccessPolicies{
		ClientAccessPolicy: &clientAccessPolicy,
		CrossDomainPolicy:  &crossDomainPolicy,
	}
}

func flattenCrossSiteAccessPolicies(input *media.CrossSiteAccessPolicies) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	clientAccessPolicy := ""
	if input.ClientAccessPolicy != nil {
		clientAccessPolicy = *input.ClientAccessPolicy
	}

	crossDomainPolicy := ""
	if input.CrossDomainPolicy != nil {
		crossDomainPolicy = *input.CrossDomainPolicy
	}

	return []interface{}{
		map[string]interface{}{
			"client_access_policy": clientAccessPolicy,
			"cross_domain_policy":  crossDomainPolicy,
		},
	}
}
