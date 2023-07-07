// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/streamingendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/media/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/media/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMediaStreamingEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMediaStreamingEndpointCreate,
		Read:   resourceMediaStreamingEndpointRead,
		Update: resourceMediaStreamingEndpointUpdate,
		Delete: resourceMediaStreamingEndpointDelete,

		DeprecationMessage: azureMediaRetirementMessage,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := streamingendpoints.ParseStreamingEndpointID(id)
			return err
		}),

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.StreamingEndpointV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StreamingEndpointName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

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

			"location": commonschema.Location(),

			"scale_units": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 50),
			},

			"access_control": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						// lintignore:XS003
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
						// lintignore:XS003
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

			"sku": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"capacity": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
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

			"tags": commonschema.Tags(),
		},
	}
}

func resourceMediaStreamingEndpointCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Media.V20220801Client.StreamingEndpoints
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := streamingendpoints.NewStreamingEndpointID(subscriptionId, d.Get("resource_group_name").(string), d.Get("media_services_account_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_media_streaming_endpoint", id.ID())
	}

	payload := streamingendpoints.StreamingEndpoint{
		Properties: &streamingendpoints.StreamingEndpointProperties{
			ScaleUnits: int64(d.Get("scale_units").(int)),
		},
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
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
		payload.Properties.AccessControl = accessControl
	}
	if cdnEnabled, ok := d.GetOk("cdn_enabled"); ok {
		payload.Properties.CdnEnabled = utils.Bool(cdnEnabled.(bool))
	}

	if cdnProfile, ok := d.GetOk("cdn_profile"); ok {
		payload.Properties.CdnProfile = utils.String(cdnProfile.(string))
	}

	if cdnProvider, ok := d.GetOk("cdn_provider"); ok {
		payload.Properties.CdnProvider = utils.String(cdnProvider.(string))
	}

	if crossSite, ok := d.GetOk("cross_site_access_policy"); ok {
		payload.Properties.CrossSiteAccessPolicies = expandStreamingEndpointCrossSiteAccessPolicies(crossSite.([]interface{}))
	}

	if _, ok := d.GetOk("custom_host_names"); ok {
		customHostNames := d.Get("custom_host_names").([]interface{})
		payload.Properties.CustomHostNames = utils.ExpandStringSlice(customHostNames)
	}

	if description, ok := d.GetOk("description"); ok {
		payload.Properties.Description = utils.String(description.(string))
	}

	if maxCacheAge, ok := d.GetOk("max_cache_age_seconds"); ok {
		payload.Properties.MaxCacheAge = utils.Int64(int64(maxCacheAge.(int)))
	}

	options := streamingendpoints.CreateOperationOptions{
		AutoStart: autoStart,
	}
	if err := client.CreateThenPoll(ctx, id, payload, options); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMediaStreamingEndpointRead(d, meta)
}

func resourceMediaStreamingEndpointUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220801Client.StreamingEndpoints
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := streamingendpoints.ParseStreamingEndpointID(d.Id())
	if err != nil {
		return err
	}
	scaleUnits := d.Get("scale_units").(int)

	if d.HasChange("scale_units") {
		scaleParameters := streamingendpoints.StreamingEntityScaleUnit{
			ScaleUnit: pointer.To(int64(scaleUnits)),
		}
		if err := client.ScaleThenPoll(ctx, *id, scaleParameters); err != nil {
			return fmt.Errorf("scaling units for %s: %+v", *id, err)
		}
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if resp.Model == nil || resp.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: model was nil", *id)
	}

	existing := resp.Model

	if d.HasChange("access_control") {
		accessControl, err := expandAccessControl(d)
		if err != nil {
			return err
		}
		existing.Properties.AccessControl = accessControl
	}

	if d.HasChange("cdn_enabled") {
		existing.Properties.CdnEnabled = utils.Bool(d.Get("cdn_enabled").(bool))
	}

	if d.HasChange("cdn_profile") {
		existing.Properties.CdnProfile = utils.String(d.Get("cdn_profile").(string))
	}

	if d.HasChange("cdn_provider") {
		existing.Properties.CdnProvider = utils.String(d.Get("cdn_provider").(string))
	}

	if d.HasChange("cross_site_access_policy") {
		existing.Properties.CrossSiteAccessPolicies = expandStreamingEndpointCrossSiteAccessPolicies(d.Get("cross_site_access_policy").([]interface{}))
	}

	if d.HasChange("custom_host_names") {
		existing.Properties.CustomHostNames = utils.ExpandStringSlice(d.Get("custom_host_names").([]interface{}))
	}

	if d.HasChange("description") {
		existing.Properties.Description = utils.String(d.Get("description").(string))
	}

	if d.HasChange("max_cache_age_seconds") {
		existing.Properties.MaxCacheAge = utils.Int64(int64(d.Get("max_cache_age_seconds").(int)))
	}

	if d.HasChange("tags") {
		existing.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.UpdateThenPoll(ctx, *id, *existing); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceMediaStreamingEndpointRead(d, meta)
}

func resourceMediaStreamingEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220801Client.StreamingEndpoints
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := streamingendpoints.ParseStreamingEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.StreamingEndpointName)
	d.Set("media_services_account_name", id.MediaServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		if props := model.Properties; props != nil {
			d.Set("scale_units", props.ScaleUnits)

			accessControl, err := flattenAccessControl(props.AccessControl)
			if err != nil {
				return fmt.Errorf("flattening `access_control`: %+v", err)
			}
			if err := d.Set("access_control", accessControl); err != nil {
				return fmt.Errorf("setting `access_control`: %+v", err)
			}

			d.Set("cdn_enabled", props.CdnEnabled)
			d.Set("cdn_profile", props.CdnProfile)
			d.Set("cdn_provider", props.CdnProvider)
			d.Set("host_name", props.HostName)

			crossSiteAccessPolicies := flattenStreamingEndpointCrossSiteAccessPolicies(props.CrossSiteAccessPolicies)
			if err := d.Set("cross_site_access_policy", crossSiteAccessPolicies); err != nil {
				return fmt.Errorf("flattening `cross_site_access_policy`: %s", err)
			}

			d.Set("custom_host_names", props.CustomHostNames)
			d.Set("description", props.Description)

			maxCacheAge := 0
			if props.MaxCacheAge != nil {
				maxCacheAge = int(*props.MaxCacheAge)
			}
			d.Set("max_cache_age_seconds", maxCacheAge)
		}

		d.Set("sku", flattenEndpointSku(model.Sku))

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceMediaStreamingEndpointDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220801Client.StreamingEndpoints
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := streamingendpoints.ParseStreamingEndpointID(d.Id())
	if err != nil {
		return err
	}

	// Stop Streaming Endpoint before we attempt to delete it.
	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("reading %s: %+v", id, err)
	}
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.ResourceState != nil && *props.ResourceState == streamingendpoints.StreamingEndpointResourceStateRunning {
				if err := client.StopThenPoll(ctx, *id); err != nil {
					return fmt.Errorf("stopping %s: %+v", id, err)
				}
			}
		}
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandAccessControl(d *pluginsdk.ResourceData) (*streamingendpoints.StreamingEndpointAccessControl, error) {
	accessControls := d.Get("access_control").([]interface{})
	if len(accessControls) == 0 {
		return nil, nil
	}
	accessControlResult := new(streamingendpoints.StreamingEndpointAccessControl)
	accessControl := accessControls[0].(map[string]interface{})
	// Get IP information
	if ipAllowsList := accessControl["ip_allow"].([]interface{}); len(ipAllowsList) > 0 {
		ipRanges := make([]streamingendpoints.IPRange, 0)
		for _, ipAllow := range ipAllowsList {
			if ipAllow == nil {
				continue
			}
			allow := ipAllow.(map[string]interface{})
			address := allow["address"].(string)
			name := allow["name"].(string)

			ipRange := streamingendpoints.IPRange{
				Name:    utils.String(name),
				Address: utils.String(address),
			}
			subnetPrefixLengthRaw := allow["subnet_prefix_length"]
			if subnetPrefixLengthRaw != "" {
				ipRange.SubnetPrefixLength = utils.Int64(int64(subnetPrefixLengthRaw.(int)))
			}
			ipRanges = append(ipRanges, ipRange)
		}
		accessControlResult.IP = &streamingendpoints.IPAccessControl{
			Allow: &ipRanges,
		}
	}
	// Get Akamai information
	if akamaiSignatureKeyList := accessControl["akamai_signature_header_authentication_key"].([]interface{}); len(akamaiSignatureKeyList) > 0 {
		akamaiSignatureHeaderAuthenticationKeyList := make([]streamingendpoints.AkamaiSignatureHeaderAuthenticationKey, 0)
		for _, akamaiSignatureKey := range akamaiSignatureKeyList {
			if akamaiSignatureKey == nil {
				continue
			}
			akamaiKey := akamaiSignatureKey.(map[string]interface{})
			base64Key := akamaiKey["base64_key"].(string)
			expirationRaw := akamaiKey["expiration"].(string)
			identifier := akamaiKey["identifier"].(string)

			akamaiSignatureHeaderAuthenticationKey := streamingendpoints.AkamaiSignatureHeaderAuthenticationKey{
				Base64Key:  utils.String(base64Key),
				Identifier: utils.String(identifier),
			}
			if expirationRaw != "" {
				expiration, err := time.Parse(time.RFC3339, expirationRaw)
				if err != nil {
					return nil, err
				}
				akamaiSignatureHeaderAuthenticationKey.SetExpirationAsTime(expiration)
			}
			akamaiSignatureHeaderAuthenticationKeyList = append(akamaiSignatureHeaderAuthenticationKeyList, akamaiSignatureHeaderAuthenticationKey)
		}
		accessControlResult.Akamai = &streamingendpoints.AkamaiAccessControl{
			AkamaiSignatureHeaderAuthenticationKeyList: &akamaiSignatureHeaderAuthenticationKeyList,
		}
	}

	return accessControlResult, nil
}

func flattenEndpointSku(input *streamingendpoints.ArmStreamingEndpointCurrentSku) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	skuName := ""
	if input.Name != nil {
		skuName = *input.Name
	}

	skuCapacity := 0
	if input.Capacity != nil {
		skuCapacity = int(*input.Capacity)
	}

	return []interface{}{
		map[string]interface{}{
			"name":     skuName,
			"capacity": skuCapacity,
		},
	}

}

func flattenAccessControl(input *streamingendpoints.StreamingEndpointAccessControl) (*[]interface{}, error) {
	if input == nil {
		return &[]interface{}{}, nil
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

			var subnetPrefixLength int64
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
				t, err := v.GetExpirationAsTime()
				if err != nil {
					return nil, fmt.Errorf("parsing expiration: %+v", err)
				}
				expiration = t.Format(time.RFC3339)
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

	return &[]interface{}{
		map[string]interface{}{
			"akamai_signature_header_authentication_key": akamaiRules,
			"ip_allow": ipAllowRules,
		},
	}, nil
}

func expandStreamingEndpointCrossSiteAccessPolicies(input []interface{}) *streamingendpoints.CrossSiteAccessPolicies {
	if len(input) == 0 {
		return nil
	}

	crossSiteAccessPolicy := input[0].(map[string]interface{})
	clientAccessPolicy := crossSiteAccessPolicy["client_access_policy"].(string)
	crossDomainPolicy := crossSiteAccessPolicy["cross_domain_policy"].(string)
	return &streamingendpoints.CrossSiteAccessPolicies{
		ClientAccessPolicy: &clientAccessPolicy,
		CrossDomainPolicy:  &crossDomainPolicy,
	}
}

func flattenStreamingEndpointCrossSiteAccessPolicies(input *streamingendpoints.CrossSiteAccessPolicies) []interface{} {
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
