// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package media

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/liveevents"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/media/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/media/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMediaLiveEvent() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMediaLiveEventCreate,
		Read:   resourceMediaLiveEventRead,
		Update: resourceMediaLiveEventUpdate,
		Delete: resourceMediaLiveEventDelete,

		DeprecationMessage: azureMediaRetirementMessage,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := liveevents.ParseLiveEventID(id)
			return err
		}),

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.LiveEventV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LiveEventName,
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
				ForceNew: true,
				Optional: true,
			},

			"location": commonschema.Location(),

			"input": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						// lintignore:XS003
						"ip_access_control_allow": {
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
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
									},
								},
							},
							AtLeastOneOf: []string{
								"input.0.ip_access_control_allow", "input.0.access_token",
								"input.0.key_frame_interval_duration", "input.0.streaming_protocol",
							},
						},

						"access_token": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{
								"input.0.ip_access_control_allow", "input.0.access_token",
								"input.0.key_frame_interval_duration", "input.0.streaming_protocol",
							},
						},

						"endpoint": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"protocol": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"url": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},

						"key_frame_interval_duration": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{
								"input.0.ip_access_control_allow", "input.0.access_token",
								"input.0.key_frame_interval_duration", "input.0.streaming_protocol",
							},
						},

						"streaming_protocol": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(liveevents.LiveEventInputProtocolRTMP),
								string(liveevents.LiveEventInputProtocolFragmentedMPFour),
							}, false),
							AtLeastOneOf: []string{
								"input.0.ip_access_control_allow", "input.0.access_token",
								"input.0.key_frame_interval_duration", "input.0.streaming_protocol",
							},
						},
					},
				},
			},

			"cross_site_access_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"client_access_policy": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{"cross_site_access_policy.0.client_access_policy", "cross_site_access_policy.0.cross_domain_policy"},
						},

						"cross_domain_policy": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{"cross_site_access_policy.0.client_access_policy", "cross_site_access_policy.0.cross_domain_policy"},
						},
					},
				},
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"encoding": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(liveevents.LiveEventEncodingTypeNone),
								string(liveevents.LiveEventEncodingTypePremiumOneZeroEightZerop),
								string(liveevents.LiveEventEncodingTypePassthroughBasic),
								string(liveevents.LiveEventEncodingTypePassthroughStandard),
								string(liveevents.LiveEventEncodingTypeStandard),
							}, false),
							Default: string(liveevents.LiveEventEncodingTypeNone),
						},

						"key_frame_interval": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "PT2S",
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"preset_name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"stretch_mode": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(liveevents.StretchModeAutoFit),
								string(liveevents.StretchModeAutoSize),
								string(liveevents.StretchModeNone),
							}, false),
							Default: string(liveevents.StretchModeNone),
						},
					},
				},
			},

			"hostname_prefix": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"preview": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						// lintignore:XS003
						"ip_access_control_allow": {
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
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
									},
								},
							},
							AtLeastOneOf: []string{
								"preview.0.ip_access_control_allow", "preview.0.alternative_media_id",
								"preview.0.preview_locator", "preview.0.streaming_policy_name",
							},
						},

						"alternative_media_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
							AtLeastOneOf: []string{
								"preview.0.ip_access_control_allow", "preview.0.alternative_media_id",
								"preview.0.preview_locator", "preview.0.streaming_policy_name",
							},
						},

						"endpoint": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"protocol": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"url": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},

						"preview_locator": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							Computed:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{
								"preview.0.ip_access_control_allow", "preview.0.alternative_media_id",
								"preview.0.preview_locator", "preview.0.streaming_policy_name",
							},
						},

						"streaming_policy_name": {
							Type:         pluginsdk.TypeString,
							Computed:     true,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{
								"preview.0.ip_access_control_allow", "preview.0.alternative_media_id",
								"preview.0.preview_locator", "preview.0.streaming_policy_name",
							},
						},
					},
				},
			},

			"stream_options": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice(liveevents.PossibleValuesForStreamOptionsFlag(), false),
				},
			},

			"transcription_languages": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"use_static_hostname": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceMediaLiveEventCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220801Client.LiveEvents
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := liveevents.NewLiveEventID(subscriptionID, d.Get("resource_group_name").(string), d.Get("media_services_account_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_media_live_event", id.ID())
	}

	t := d.Get("tags").(map[string]interface{})

	payload := liveevents.LiveEvent{
		Properties: &liveevents.LiveEventProperties{},
		Location:   location.Normalize(d.Get("location").(string)),
		Tags:       tags.Expand(t),
	}

	autoStart := utils.Bool(false)
	if _, ok := d.GetOk("auto_start_enabled"); ok {
		autoStart = utils.Bool(d.Get("auto_start_enabled").(bool))
	}

	if input, ok := d.GetOk("input"); ok {
		payload.Properties.Input = expandLiveEventInput(input.([]interface{}))
	}

	if crossSitePolicies, ok := d.GetOk("cross_site_access_policy"); ok {
		payload.Properties.CrossSiteAccessPolicies = expandLiveEventCrossSiteAccessPolicies(crossSitePolicies.([]interface{}))
	}

	if description, ok := d.GetOk("description"); ok {
		payload.Properties.Description = utils.String(description.(string))
	}

	if encoding, ok := d.GetOk("encoding"); ok {
		payload.Properties.Encoding = expandEncoding(encoding.([]interface{}))
	}

	if hostNamePrefix, ok := d.GetOk("hostname_prefix"); ok {
		payload.Properties.HostnamePrefix = utils.String(hostNamePrefix.(string))
	}

	if preview, ok := d.GetOk("preview"); ok {
		payload.Properties.Preview = expandPreview(preview.([]interface{}))
	}

	if streamOptions, ok := d.GetOk("stream_options"); ok {
		payload.Properties.StreamOptions = expandStreamOptions(streamOptions.([]interface{}))
	}

	if transcriptionLanguages, ok := d.GetOk("transcription_languages"); ok {
		payload.Properties.Transcriptions = expandTranscriptions(transcriptionLanguages.([]interface{}))
	}

	if useStaticHostName, ok := d.GetOk("use_static_hostname"); ok {
		payload.Properties.UseStaticHostname = utils.Bool(useStaticHostName.(bool))
	}

	options := liveevents.CreateOperationOptions{
		AutoStart: autoStart,
	}
	if err := client.CreateThenPoll(ctx, id, payload, options); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMediaLiveEventRead(d, meta)
}

func resourceMediaLiveEventUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220801Client.LiveEvents
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := liveevents.ParseLiveEventID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil {
		return fmt.Errorf("unexpected null model of %s", id)
	}
	existing := resp.Model

	if d.HasChange("input") {
		existing.Properties.Input = expandLiveEventInput(d.Get("input").([]interface{}))
	}

	if d.HasChange("cross_site_access_policy") {
		existing.Properties.CrossSiteAccessPolicies = expandLiveEventCrossSiteAccessPolicies(d.Get("cross_site_access_policy").([]interface{}))
	}

	if d.HasChange("description") {
		existing.Properties.Description = utils.String(d.Get("description").(string))
	}

	if d.HasChange("encoding") {
		existing.Properties.Encoding = expandEncoding(d.Get("encoding").([]interface{}))
	}

	if d.HasChange("hostname_prefix") {
		existing.Properties.HostnamePrefix = utils.String(d.Get("hostname_prefix").(string))
	}

	if d.HasChange("preview") {
		existing.Properties.Preview = expandPreview(d.Get("preview").([]interface{}))
	}

	if d.HasChange("transcription_languages") {
		existing.Properties.Transcriptions = expandTranscriptions(d.Get("transcription_languages").([]interface{}))
	}

	if d.HasChange("tags") {
		existing.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.UpdateThenPoll(ctx, *id, *existing); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceMediaLiveEventRead(d, meta)
}

func resourceMediaLiveEventRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220801Client.LiveEvents
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := liveevents.ParseLiveEventID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.LiveEventName)
	d.Set("media_services_account_name", id.MediaServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			input := flattenLiveEventInput(props.Input)
			if err := d.Set("input", input); err != nil {
				return fmt.Errorf("flattening `input`: %s", err)
			}

			crossSiteAccessPolicies := flattenLiveEventCrossSiteAccessPolicies(props.CrossSiteAccessPolicies)
			if err := d.Set("cross_site_access_policy", crossSiteAccessPolicies); err != nil {
				return fmt.Errorf("flattening `cross_site_access_policy`: %s", err)
			}

			encoding := flattenEncoding(props.Encoding)
			if err := d.Set("encoding", encoding); err != nil {
				return fmt.Errorf("flattening `encoding`: %s", err)
			}

			d.Set("description", props.Description)
			d.Set("hostname_prefix", props.HostnamePrefix)

			preview := flattenPreview(props.Preview)
			if err := d.Set("preview", preview); err != nil {
				return fmt.Errorf("flattening `preview`: %s", err)
			}

			streamOptions := flattenStreamOptions(props.StreamOptions)
			if err := d.Set("stream_options", streamOptions); err != nil {
				return fmt.Errorf("flattening `stream_options`: %s", err)
			}

			transcriptions := flattenTranscriptions(props.Transcriptions)
			if err := d.Set("transcription_languages", transcriptions); err != nil {
				return fmt.Errorf("flattening `transcription_languages`: %s", err)
			}

			useStaticHostName := false
			if props.UseStaticHostname != nil {
				useStaticHostName = *props.UseStaticHostname
			}
			d.Set("use_static_hostname", useStaticHostName)
		}
	}

	return nil
}

func resourceMediaLiveEventDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220801Client.LiveEvents
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := liveevents.ParseLiveEventID(d.Id())
	if err != nil {
		return err
	}

	// Stop Live Event before we attempt to delete it.
	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.ResourceState != nil && *props.ResourceState == liveevents.LiveEventResourceStateRunning {
				payload := liveevents.LiveEventActionInput{
					RemoveOutputsOnStop: utils.Bool(false),
				}
				if err := client.StopThenPoll(ctx, *id, payload); err != nil {
					return fmt.Errorf("stopping %s: %+v", *id, err)
				}
			}
		}
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandLiveEventInput(input []interface{}) liveevents.LiveEventInput {
	liveInput := input[0].(map[string]interface{})

	var inputAccessControl *liveevents.LiveEventInputAccessControl
	if v := liveInput["ip_access_control_allow"]; v != nil {
		ipRanges := expandIPRanges(v.([]interface{}))
		inputAccessControl = &liveevents.LiveEventInputAccessControl{
			IP: &liveevents.IPAccessControl{
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

	return liveevents.LiveEventInput{
		AccessControl:            inputAccessControl,
		AccessToken:              utils.String(accessToken),
		KeyFrameIntervalDuration: utils.String(keyFrameInterval),
		StreamingProtocol:        liveevents.LiveEventInputProtocol(streamingProtocol),
	}
}

func expandIPRanges(input []interface{}) []liveevents.IPRange {
	if len(input) == 0 {
		return nil
	}

	ipRanges := make([]liveevents.IPRange, 0)
	for _, ipAllow := range input {
		if ipAllow == nil {
			continue
		}
		allow := ipAllow.(map[string]interface{})
		address := allow["address"].(string)
		name := allow["name"].(string)

		ipRange := liveevents.IPRange{
			Name:    utils.String(name),
			Address: utils.String(address),
		}
		subnetPrefixLengthRaw := allow["subnet_prefix_length"]
		if subnetPrefixLengthRaw != "" {
			ipRange.SubnetPrefixLength = pointer.To(int64(subnetPrefixLengthRaw.(int)))
		}
		ipRanges = append(ipRanges, ipRange)
	}

	return ipRanges
}

func expandEncoding(input []interface{}) *liveevents.LiveEventEncoding {
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

	liveEventEncoding := &liveevents.LiveEventEncoding{
		EncodingType: pointer.To(liveevents.LiveEventEncodingType(encodingType)),
		StretchMode:  pointer.To(liveevents.StretchMode(stretchMode)),
	}

	if v := liveEncoding["key_frame_interval"]; v != nil && v.(string) != "" {
		liveEventEncoding.KeyFrameInterval = utils.String(v.(string))
	}

	if v := liveEncoding["preset_name"]; v != nil {
		liveEventEncoding.PresetName = utils.String(v.(string))
	}

	return liveEventEncoding
}

func expandPreview(input []interface{}) *liveevents.LiveEventPreview {
	if len(input) == 0 {
		return nil
	}

	livePreview := input[0].(map[string]interface{})
	var inputAccessControl *liveevents.LiveEventPreviewAccessControl
	if v := livePreview["ip_access_control_allow"]; v != nil {
		ipRanges := expandIPRanges(v.([]interface{}))
		inputAccessControl = &liveevents.LiveEventPreviewAccessControl{
			IP: &liveevents.IPAccessControl{
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

	return &liveevents.LiveEventPreview{
		AccessControl:       inputAccessControl,
		AlternativeMediaId:  utils.String(alternativeMediaID),
		PreviewLocator:      utils.String(previewLocator),
		StreamingPolicyName: utils.String(streamingPolicyName),
	}
}

func expandLiveEventCrossSiteAccessPolicies(input []interface{}) *liveevents.CrossSiteAccessPolicies {
	if len(input) == 0 {
		return nil
	}

	crossSiteAccessPolicy := input[0].(map[string]interface{})
	clientAccessPolicy := crossSiteAccessPolicy["client_access_policy"].(string)
	crossDomainPolicy := crossSiteAccessPolicy["cross_domain_policy"].(string)
	return &liveevents.CrossSiteAccessPolicies{
		ClientAccessPolicy: &clientAccessPolicy,
		CrossDomainPolicy:  &crossDomainPolicy,
	}
}

func expandStreamOptions(input []interface{}) *[]liveevents.StreamOptionsFlag {
	streamOptions := make([]liveevents.StreamOptionsFlag, 0)
	for _, v := range input {
		streamOptions = append(streamOptions, liveevents.StreamOptionsFlag(v.(string)))
	}
	return &streamOptions
}

func expandTranscriptions(input []interface{}) *[]liveevents.LiveEventTranscription {
	transcriptions := make([]liveevents.LiveEventTranscription, 0)
	for _, v := range input {
		transcriptions = append(transcriptions, liveevents.LiveEventTranscription{
			Language: utils.String(v.(string)),
		})
	}
	return &transcriptions
}

func flattenLiveEventInput(input liveevents.LiveEventInput) []interface{} {
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

func flattenEventAccessControl(input *liveevents.LiveEventInputAccessControl) []interface{} {
	if input == nil || input.IP == nil || input.IP.Allow == nil {
		return make([]interface{}, 0)
	}

	return flattenIPAllow(input.IP.Allow)
}

func flattenEndpoints(input *[]liveevents.LiveEventEndpoint) []interface{} {
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
		if v.Url != nil {
			url = *v.Url
		}

		endpoints = append(endpoints, map[string]interface{}{
			"protocol": protocol,
			"url":      url,
		})
	}

	return endpoints
}

func flattenEncoding(input *liveevents.LiveEventEncoding) []interface{} {
	if input == nil || (input.KeyFrameInterval == nil && input.PresetName == nil && (input.EncodingType == nil || *input.EncodingType == liveevents.LiveEventEncodingTypeNone)) {
		return make([]interface{}, 0)
	}

	encodingType := ""
	if input.EncodingType != nil {
		encodingType = string(*input.EncodingType)
	}

	keyFrameInterval := ""
	if input.KeyFrameInterval != nil {
		keyFrameInterval = *input.KeyFrameInterval
	}

	presetName := ""
	if input.PresetName != nil {
		presetName = *input.PresetName
	}

	stretchMode := ""
	if input.StretchMode != nil {
		stretchMode = string(*input.StretchMode)
	}

	return []interface{}{
		map[string]interface{}{
			"type":               encodingType,
			"key_frame_interval": keyFrameInterval,
			"preset_name":        presetName,
			"stretch_mode":       stretchMode,
		},
	}
}

func flattenPreview(input *liveevents.LiveEventPreview) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	iPAccessControlAllow := flattenPreviewAccessControl(input.AccessControl)

	alternativeMediaID := ""
	if input.AlternativeMediaId != nil {
		alternativeMediaID = *input.AlternativeMediaId
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

func flattenPreviewAccessControl(input *liveevents.LiveEventPreviewAccessControl) []interface{} {
	if input == nil || input.IP == nil || input.IP.Allow == nil {
		return make([]interface{}, 0)
	}

	return flattenIPAllow(input.IP.Allow)
}

func flattenIPAllow(input *[]liveevents.IPRange) []interface{} {
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

		var subnetPrefixLength int64
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

func flattenLiveEventCrossSiteAccessPolicies(input *liveevents.CrossSiteAccessPolicies) []interface{} {
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

func flattenStreamOptions(input *[]liveevents.StreamOptionsFlag) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	streamOptions := make([]interface{}, 0)
	for _, v := range *input {
		streamOptions = append(streamOptions, string(v))
	}

	return streamOptions
}

func flattenTranscriptions(input *[]liveevents.LiveEventTranscription) []string {
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
