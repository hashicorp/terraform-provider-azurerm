package media

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2020-05-01/media"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMediaLiveEvent() *schema.Resource {
	return &schema.Resource{
		Create: resourceMediaLiveEventCreateUpdate,
		Read:   resourceMediaLiveEventRead,
		Update: resourceMediaLiveEventCreateUpdate,
		Delete: resourceMediaLiveEventDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.LiveEventID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LiveEventName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"media_services_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName,
			},

			"auto_start_enabled": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
			},

			"location": azure.SchemaLocation(),

			"input": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						//lintignore:XS003
						"ip_access_control_allow": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"name": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"subnet_prefix_length": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
									},
								},
							},
							AtLeastOneOf: []string{"input.0.ip_access_control_allow", "input.0.access_token",
								"input.0.key_frame_interval_duration", "input.0.streaming_protocol",
							},
						},

						"access_token": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{"input.0.ip_access_control_allow", "input.0.access_token",
								"input.0.key_frame_interval_duration", "input.0.streaming_protocol",
							},
						},

						"endpoint": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"url": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},

						"key_frame_interval_duration": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{"input.0.ip_access_control_allow", "input.0.access_token",
								"input.0.key_frame_interval_duration", "input.0.streaming_protocol",
							},
						},

						"streaming_protocol": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(media.RTMP),
								string(media.FragmentedMP4),
							}, false),
							AtLeastOneOf: []string{"input.0.ip_access_control_allow", "input.0.access_token",
								"input.0.key_frame_interval_duration", "input.0.streaming_protocol",
							},
						},
					},
				},
			},

			"cross_site_access_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_access_policy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{"cross_site_access_policy.0.client_access_policy", "cross_site_access_policy.0.cross_domain_policy"},
						},

						"cross_domain_policy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{"cross_site_access_policy.0.client_access_policy", "cross_site_access_policy.0.cross_domain_policy"},
						},
					},
				},
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"encoding": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(media.LiveEventEncodingTypeNone),
								string(media.LiveEventEncodingTypePremium1080p),
								string(media.LiveEventEncodingTypeStandard),
							}, false),
							Default: string(media.LiveEventEncodingTypeNone),
						},

						"key_frame_interval": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "PT2S",
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"preset_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"stretch_mode": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(media.StretchModeAutoFit),
								string(media.StretchModeAutoSize),
								string(media.StretchModeNone),
							}, false),
							Default: string(media.StretchModeNone),
						},
					},
				},
			},

			"hostname_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"preview": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						//lintignore:XS003
						"ip_access_control_allow": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"name": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"subnet_prefix_length": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
									},
								},
							},
							AtLeastOneOf: []string{"preview.0.ip_access_control_allow", "preview.0.alternative_media_id",
								"preview.0.preview_locator", "preview.0.streaming_policy_name",
							},
						},

						"alternative_media_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
							AtLeastOneOf: []string{"preview.0.ip_access_control_allow", "preview.0.alternative_media_id",
								"preview.0.preview_locator", "preview.0.streaming_policy_name",
							},
						},

						"endpoint": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"url": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},

						"preview_locator": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Computed:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{"preview.0.ip_access_control_allow", "preview.0.alternative_media_id",
								"preview.0.preview_locator", "preview.0.streaming_policy_name",
							},
						},

						"streaming_policy_name": {
							Type:         schema.TypeString,
							Computed:     true,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{"preview.0.ip_access_control_allow", "preview.0.alternative_media_id",
								"preview.0.preview_locator", "preview.0.streaming_policy_name",
							},
						},
					},
				},
			},

			"transcription_languages": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"use_static_hostname": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceMediaLiveEventCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.LiveEventsClient
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceID := parse.NewLiveEventID(subscriptionID, d.Get("resource_group_name").(string), d.Get("media_services_account_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceID.ResourceGroup, resourceID.MediaserviceName, resourceID.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", resourceID, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_media_live_event", resourceID.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	parameters := media.LiveEvent{
		LiveEventProperties: &media.LiveEventProperties{},
		Location:            utils.String(location),
		Tags:                tags.Expand(t),
	}

	autoStart := utils.Bool(false)
	if _, ok := d.GetOk("auto_start_enabled"); ok {
		autoStart = utils.Bool(d.Get("auto_start_enabled").(bool))
	}

	if input, ok := d.GetOk("input"); ok {
		parameters.LiveEventProperties.Input = expandLiveEventInput(input.([]interface{}))
	}

	if crossSitePolicies, ok := d.GetOk("cross_site_access_policy"); ok {
		parameters.LiveEventProperties.CrossSiteAccessPolicies = expandCrossSiteAccessPolicies(crossSitePolicies.([]interface{}))
	}

	if description, ok := d.GetOk("description"); ok {
		parameters.LiveEventProperties.Description = utils.String(description.(string))
	}

	if encoding, ok := d.GetOk("encoding"); ok {
		parameters.LiveEventProperties.Encoding = expandEncoding(encoding.([]interface{}))
	}

	if hostNamePrefix, ok := d.GetOk("hostname_prefix"); ok {
		parameters.LiveEventProperties.HostnamePrefix = utils.String(hostNamePrefix.(string))
	}

	if preview, ok := d.GetOk("preview"); ok {
		parameters.LiveEventProperties.Preview = expandPreview(preview.([]interface{}))
	}

	if transcriptionLanguages, ok := d.GetOk("transcription_languages"); ok {
		parameters.LiveEventProperties.Transcriptions = expandTranscriptions(transcriptionLanguages.([]interface{}))
	}

	if useStaticHostName, ok := d.GetOk("use_static_hostname"); ok {
		parameters.LiveEventProperties.UseStaticHostname = utils.Bool(useStaticHostName.(bool))
	}

	if d.IsNewResource() {
		future, err := client.Create(ctx, resourceID.ResourceGroup, resourceID.MediaserviceName, resourceID.Name, parameters, autoStart)
		if err != nil {
			return fmt.Errorf("creating %s: %+v", resourceID, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for creation %s: %+v", resourceID, err)
		}
	} else {
		future, err := client.Update(ctx, resourceID.ResourceGroup, resourceID.MediaserviceName, resourceID.Name, parameters)
		if err != nil {
			return fmt.Errorf("updating %s: %+v", resourceID, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for %s to update: %+v", resourceID, err)
		}
	}

	d.SetId(resourceID.ID())

	return resourceMediaLiveEventRead(d, meta)
}

func resourceMediaLiveEventRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.LiveEventsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LiveEventID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("media_services_account_name", id.MediaserviceName)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.LiveEventProperties; props != nil {
		input := flattenLiveEventInput(props.Input)
		if err := d.Set("input", input); err != nil {
			return fmt.Errorf("Error flattening `input`: %s", err)
		}

		crossSiteAccessPolicies := flattenLiveEventCrossSiteAccessPolicies(resp.CrossSiteAccessPolicies)
		if err := d.Set("cross_site_access_policy", crossSiteAccessPolicies); err != nil {
			return fmt.Errorf("Error flattening `cross_site_access_policy`: %s", err)
		}

		encoding := flattenEncoding(resp.Encoding)
		if err := d.Set("encoding", encoding); err != nil {
			return fmt.Errorf("Error flattening `encoding`: %s", err)
		}

		d.Set("description", props.Description)
		d.Set("hostname_prefix", props.HostnamePrefix)

		preview := flattenPreview(resp.Preview)
		if err := d.Set("preview", preview); err != nil {
			return fmt.Errorf("Error flattening `preview`: %s", err)
		}

		transcriptions := flattenTranscriptions(resp.Transcriptions)
		if err := d.Set("transcription_languages", transcriptions); err != nil {
			return fmt.Errorf("Error flattening `transcription_languages`: %s", err)
		}

		useStaticHostName := false
		if props.UseStaticHostname != nil {
			useStaticHostName = *props.UseStaticHostname
		}
		d.Set("use_static_hostname", useStaticHostName)
	}

	return nil
}

func resourceMediaLiveEventDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.LiveEventsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LiveEventID(d.Id())
	if err != nil {
		return err
	}

	// Stop Live Event before we attempt to delete it.
	resp, err := client.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		return fmt.Errorf("reading %s: %+v", id, err)
	}
	if props := resp.LiveEventProperties; props != nil {
		if props.ResourceState == media.Running {
			stopFuture, err := client.Stop(ctx, id.ResourceGroup, id.MediaserviceName, id.Name, media.LiveEventActionInput{RemoveOutputsOnStop: utils.Bool(false)})
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
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to delete: %+v", id, err)
	}

	return nil
}

func expandLiveEventInput(input []interface{}) *media.LiveEventInput {
	if len(input) == 0 {
		return nil
	}

	liveInput := input[0].(map[string]interface{})

	var inputAccessControl *media.LiveEventInputAccessControl
	if v := liveInput["ip_access_control_allow"]; v != nil {
		ipRanges := expandIPRanges(v.([]interface{}))
		inputAccessControl = &media.LiveEventInputAccessControl{
			IP: &media.IPAccessControl{
				Allow: &ipRanges,
			},
		}
	}

	accessToken := ""
	if v := liveInput["access_token"]; v != nil {
		accessToken = v.(string)
	}

	keyFrameInterval := ""
	if v := liveInput["key_frame_interval_duration"]; v != nil {
		keyFrameInterval = v.(string)
	}

	streamingProtocol := ""
	if v := liveInput["streaming_protocol"]; v != nil {
		streamingProtocol = v.(string)
	}

	return &media.LiveEventInput{
		AccessControl:            inputAccessControl,
		AccessToken:              utils.String(accessToken),
		KeyFrameIntervalDuration: utils.String(keyFrameInterval),
		StreamingProtocol:        media.LiveEventInputProtocol(streamingProtocol),
	}
}

func expandIPRanges(input []interface{}) []media.IPRange {
	if len(input) == 0 {
		return nil
	}

	ipRanges := make([]media.IPRange, 0)
	for _, ipAllow := range input {
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

	return ipRanges
}

func expandEncoding(input []interface{}) *media.LiveEventEncoding {
	if len(input) == 0 {
		return nil
	}

	liveEncoding := input[0].(map[string]interface{})

	encodingType := ""
	if v := liveEncoding["type"]; v != nil {
		encodingType = v.(string)
	}

	stretchMode := ""
	if v := liveEncoding["stretch_mode"]; v != nil {
		stretchMode = v.(string)
	}

	liveEventEncoding := &media.LiveEventEncoding{
		EncodingType: media.LiveEventEncodingType(encodingType),
		StretchMode:  media.StretchMode(stretchMode),
	}

	if v := liveEncoding["key_frame_interval"]; v != nil && v.(string) != "" {
		liveEventEncoding.KeyFrameInterval = utils.String(v.(string))
	}

	if v := liveEncoding["preset_name"]; v != nil {
		liveEventEncoding.PresetName = utils.String(v.(string))
	}

	return liveEventEncoding
}

func expandPreview(input []interface{}) *media.LiveEventPreview {
	if len(input) == 0 {
		return nil
	}

	livePreview := input[0].(map[string]interface{})
	var inputAccessControl *media.LiveEventPreviewAccessControl
	if v := livePreview["ip_access_control_allow"]; v != nil {
		ipRanges := expandIPRanges(v.([]interface{}))
		inputAccessControl = &media.LiveEventPreviewAccessControl{
			IP: &media.IPAccessControl{
				Allow: &ipRanges,
			},
		}
	}

	alternativeMediaID := ""
	if v := livePreview["alternative_media_id"]; v != nil {
		alternativeMediaID = v.(string)
	}

	previewLocator := ""
	if v := livePreview["preview_locator"]; v != nil {
		previewLocator = v.(string)
	}

	streamingPolicyName := ""
	if v := livePreview["streaming_policy_name"]; v != nil {
		streamingPolicyName = v.(string)
	}

	return &media.LiveEventPreview{
		AccessControl:       inputAccessControl,
		AlternativeMediaID:  utils.String(alternativeMediaID),
		PreviewLocator:      utils.String(previewLocator),
		StreamingPolicyName: utils.String(streamingPolicyName),
	}
}

func expandTranscriptions(input []interface{}) *[]media.LiveEventTranscription {
	transcriptions := make([]media.LiveEventTranscription, 0)
	for _, v := range input {
		transcriptions = append(transcriptions, media.LiveEventTranscription{
			Language: utils.String(v.(string)),
		})
	}
	return &transcriptions
}

func flattenLiveEventInput(input *media.LiveEventInput) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	ipAccessControlAllow := flattenEventAccessControl(input.AccessControl)

	accessToken := ""
	if input.AccessToken != nil {
		accessToken = *input.AccessToken
	}

	endpoints := flattenEndpoints(input.Endpoints)

	keyFrameInterval := ""
	if input.KeyFrameIntervalDuration != nil {
		keyFrameInterval = *input.KeyFrameIntervalDuration
	}

	return []interface{}{
		map[string]interface{}{
			"ip_access_control_allow":     ipAccessControlAllow,
			"access_token":                accessToken,
			"endpoint":                    endpoints,
			"key_frame_interval_duration": keyFrameInterval,
			"streaming_protocol":          string(input.StreamingProtocol),
		},
	}
}

func flattenEventAccessControl(input *media.LiveEventInputAccessControl) []interface{} {
	if input == nil || input.IP == nil || input.IP.Allow == nil {
		return make([]interface{}, 0)
	}

	return flattenIPAllow(input.IP.Allow)
}

func flattenEndpoints(input *[]media.LiveEventEndpoint) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	endpoints := make([]interface{}, 0)
	for _, v := range *input {
		protocol := ""
		if v.Protocol != nil {
			protocol = *v.Protocol
		}

		url := ""
		if v.URL != nil {
			url = *v.URL
		}

		endpoints = append(endpoints, map[string]interface{}{
			"protocol": protocol,
			"url":      url,
		})
	}

	return endpoints
}

func flattenEncoding(input *media.LiveEventEncoding) []interface{} {
	if input == nil || (input.KeyFrameInterval == nil && input.PresetName == nil && input.EncodingType == media.LiveEventEncodingTypeNone) {
		return make([]interface{}, 0)
	}

	keyFrameInterval := ""
	if input.KeyFrameInterval != nil {
		keyFrameInterval = *input.KeyFrameInterval
	}

	presetName := ""
	if input.PresetName != nil {
		presetName = *input.PresetName
	}

	return []interface{}{
		map[string]interface{}{
			"type":               string(input.EncodingType),
			"key_frame_interval": keyFrameInterval,
			"preset_name":        presetName,
			"stretch_mode":       string(input.StretchMode),
		},
	}
}

func flattenPreview(input *media.LiveEventPreview) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	iPAccessControlAllow := flattenPreviewAccessControl(input.AccessControl)

	alternativeMediaID := ""
	if input.AlternativeMediaID != nil {
		alternativeMediaID = *input.AlternativeMediaID
	}

	endpoints := flattenEndpoints(input.Endpoints)

	previewLocator := ""
	if input.PreviewLocator != nil {
		previewLocator = *input.PreviewLocator
	}

	streamingPolicyName := ""
	if input.StreamingPolicyName != nil {
		streamingPolicyName = *input.StreamingPolicyName
	}

	return []interface{}{
		map[string]interface{}{
			"ip_access_control_allow": iPAccessControlAllow,
			"alternative_media_id":    alternativeMediaID,
			"endpoint":                endpoints,
			"preview_locator":         previewLocator,
			"streaming_policy_name":   streamingPolicyName,
		},
	}
}

func flattenPreviewAccessControl(input *media.LiveEventPreviewAccessControl) []interface{} {
	if input == nil || input.IP == nil || input.IP.Allow == nil {
		return make([]interface{}, 0)
	}

	return flattenIPAllow(input.IP.Allow)
}

func flattenIPAllow(input *[]media.IPRange) []interface{} {
	ipAllow := make([]interface{}, 0)

	for _, v := range *input {
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

		ipAllow = append(ipAllow, map[string]interface{}{
			"name":                 name,
			"address":              address,
			"subnet_prefix_length": subnetPrefixLength,
		})
	}

	return ipAllow
}

func flattenLiveEventCrossSiteAccessPolicies(input *media.CrossSiteAccessPolicies) []interface{} {
	if input == nil || (input.ClientAccessPolicy == nil && input.CrossDomainPolicy == nil) {
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

func flattenTranscriptions(input *[]media.LiveEventTranscription) []string {
	if input == nil {
		return make([]string, 0)
	}

	transcriptionLanguages := make([]string, 0)
	for _, v := range *input {
		if v.Language != nil && *v.Language != "" {
			transcriptionLanguages = append(transcriptionLanguages, *v.Language)
		}
	}

	return transcriptionLanguages
}
