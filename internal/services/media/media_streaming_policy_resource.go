// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/streamingpoliciesandstreaminglocators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/media/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/media/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMediaStreamingPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMediaStreamingPolicyCreate,
		Read:   resourceMediaStreamingPolicyRead,
		Delete: resourceMediaStreamingPolicyDelete,

		DeprecationMessage: azureMediaRetirementMessage,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := streamingpoliciesandstreaminglocators.ParseStreamingPolicyID(id)
			return err
		}),

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.StreamingPolicyV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9(_)]{1,128}$"),
					"Streaming Policy name must be 1 - 128 characters long, can contain letters, numbers, underscores, and hyphens (but the first and last character must be a letter or number).",
				),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"media_services_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName,
			},

			"no_encryption_enabled_protocols": enabledProtocolsSchema(),

			"common_encryption_cenc": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"clear_key_encryption": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							ConflictsWith: []string{
								"common_encryption_cenc.0.drm_widevine_custom_license_acquisition_url_template",
								"common_encryption_cenc.0.drm_playready",
							},
							AtLeastOneOf: []string{
								"common_encryption_cenc.0.drm_widevine_custom_license_acquisition_url_template",
								"common_encryption_cenc.0.drm_playready",
								"common_encryption_cenc.0.clear_key_encryption",
							},
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"custom_keys_acquisition_url_template": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.IsURLWithHTTPS,
									},
								},
							},
						},

						"clear_track": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*schema.Schema{
									"condition": {
										Type:     pluginsdk.TypeSet,
										Required: true,
										ForceNew: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"operation": {
													Type:     pluginsdk.TypeString,
													Required: true,
													ForceNew: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(streamingpoliciesandstreaminglocators.TrackPropertyCompareOperationEqual),
													}, false),
												},

												"property": {
													Type:     pluginsdk.TypeString,
													Required: true,
													ForceNew: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(streamingpoliciesandstreaminglocators.TrackPropertyTypeFourCC),
													}, false),
												},

												"value": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ForceNew:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
										},
									},
								},
							},
						},

						"enabled_protocols": enabledProtocolsSchema(),

						"drm_widevine_custom_license_acquisition_url_template": {
							Type:          pluginsdk.TypeString,
							Optional:      true,
							ForceNew:      true,
							ValidateFunc:  validation.IsURLWithHTTPS,
							ConflictsWith: []string{"common_encryption_cenc.0.clear_key_encryption"},
						},

						"drm_playready": {
							Type:          pluginsdk.TypeList,
							Optional:      true,
							ForceNew:      true,
							MaxItems:      1,
							ConflictsWith: []string{"common_encryption_cenc.0.clear_key_encryption"},
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"custom_license_acquisition_url_template": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ForceNew:     true,
										ValidateFunc: validation.IsURLWithHTTPS,
										AtLeastOneOf: []string{"common_encryption_cenc.0.drm_playready.0.custom_license_acquisition_url_template", "common_encryption_cenc.0.drm_playready.0.custom_attributes"},
									},

									"custom_attributes": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ForceNew:     true,
										ValidateFunc: validation.StringIsNotEmpty,
										AtLeastOneOf: []string{"common_encryption_cenc.0.drm_playready.0.custom_license_acquisition_url_template", "common_encryption_cenc.0.drm_playready.0.custom_attributes"},
									},
								},
							},
						},

						"default_content_key": defaultContentKeySchema(),

						"content_key_to_track_mapping": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"label": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ForceNew:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"policy_name": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ForceNew:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"track": {
										Type:     pluginsdk.TypeSet,
										Required: true,
										ForceNew: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*schema.Schema{
												"condition": {
													Type:     pluginsdk.TypeSet,
													Required: true,
													ForceNew: true,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*pluginsdk.Schema{
															"operation": {
																Type:     pluginsdk.TypeString,
																Required: true,
																ForceNew: true,
																ValidateFunc: validation.StringInSlice([]string{
																	string(streamingpoliciesandstreaminglocators.TrackPropertyCompareOperationEqual),
																}, false),
															},

															"property": {
																Type:     pluginsdk.TypeString,
																Required: true,
																ForceNew: true,
																ValidateFunc: validation.StringInSlice([]string{
																	string(streamingpoliciesandstreaminglocators.TrackPropertyTypeFourCC),
																}, false),
															},

															"value": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ForceNew:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"common_encryption_cbcs": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"clear_key_encryption": {
							Type:         pluginsdk.TypeList,
							Optional:     true,
							ForceNew:     true,
							MaxItems:     1,
							ExactlyOneOf: []string{"common_encryption_cbcs.0.drm_fairplay", "common_encryption_cbcs.0.clear_key_encryption"},
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"custom_keys_acquisition_url_template": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.IsURLWithHTTPS,
									},
								},
							},
						},

						"enabled_protocols": enabledProtocolsSchema(),

						"drm_fairplay": {
							Type:         pluginsdk.TypeList,
							Optional:     true,
							ForceNew:     true,
							MaxItems:     1,
							ExactlyOneOf: []string{"common_encryption_cbcs.0.drm_fairplay", "common_encryption_cbcs.0.clear_key_encryption"},
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"custom_license_acquisition_url_template": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ForceNew:     true,
										ValidateFunc: validation.IsURLWithHTTPS,
										AtLeastOneOf: []string{"common_encryption_cbcs.0.drm_fairplay.0.custom_license_acquisition_url_template", "common_encryption_cbcs.0.drm_fairplay.0.allow_persistent_license"},
									},

									"allow_persistent_license": {
										Type:         pluginsdk.TypeBool,
										Optional:     true,
										ForceNew:     true,
										AtLeastOneOf: []string{"common_encryption_cbcs.0.drm_fairplay.0.custom_license_acquisition_url_template", "common_encryption_cbcs.0.drm_fairplay.0.allow_persistent_license"},
									},
								},
							},
						},

						"default_content_key": defaultContentKeySchema(),
					},
				},
			},

			"envelope_encryption": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"custom_keys_acquisition_url_template": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IsURLWithHTTPS,
						},

						"default_content_key": defaultContentKeySchema(),

						"enabled_protocols": enabledProtocolsSchema(),
					},
				},
			},

			"default_content_key_policy_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceMediaStreamingPolicyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220801Client.StreamingPoliciesAndStreamingLocators
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := streamingpoliciesandstreaminglocators.NewStreamingPolicyID(subscriptionID, d.Get("resource_group_name").(string), d.Get("media_services_account_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.StreamingPoliciesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_media_streaming_policy", id.ID())
		}
	}

	payload := streamingpoliciesandstreaminglocators.StreamingPolicy{
		Properties: &streamingpoliciesandstreaminglocators.StreamingPolicyProperties{},
	}

	if noEncryption, ok := d.GetOk("no_encryption_enabled_protocols"); ok {
		payload.Properties.NoEncryption = &streamingpoliciesandstreaminglocators.NoEncryption{
			EnabledProtocols: expandEnabledProtocols(noEncryption.([]interface{})),
		}
	}

	if commonEncryptionCENC, ok := d.GetOk("common_encryption_cenc"); ok {
		payload.Properties.CommonEncryptionCenc = expandCommonEncryptionCenc(commonEncryptionCENC.([]interface{}))
	}

	if commonEncryptionCBCS, ok := d.GetOk("common_encryption_cbcs"); ok {
		payload.Properties.CommonEncryptionCbcs = expandCommonEncryptionCbcs(commonEncryptionCBCS.([]interface{}))
	}

	if contentKeyPolicyName, ok := d.GetOk("default_content_key_policy_name"); ok {
		payload.Properties.DefaultContentKeyPolicyName = utils.String(contentKeyPolicyName.(string))
	}

	if envelopeEncryption, ok := d.GetOk("envelope_encryption"); ok {
		payload.Properties.EnvelopeEncryption = expandEnvelopeEncryption(envelopeEncryption.([]interface{}))
	}

	if _, err := client.StreamingPoliciesCreate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMediaStreamingPolicyRead(d, meta)
}

func resourceMediaStreamingPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220801Client.StreamingPoliciesAndStreamingLocators
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := streamingpoliciesandstreaminglocators.ParseStreamingPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.StreamingPoliciesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.StreamingPolicyName)
	d.Set("media_services_account_name", id.MediaServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			noEncryption := flattenNoEncryption(props.NoEncryption)
			if err := d.Set("no_encryption_enabled_protocols", noEncryption); err != nil {
				return fmt.Errorf("flattening `no_encryption_enabled_protocols`: %s", err)
			}

			commonEncryptionCENC := flattenCommonEncryptionCenc(props.CommonEncryptionCenc)
			if err := d.Set("common_encryption_cenc", commonEncryptionCENC); err != nil {
				return fmt.Errorf("flattening `common_encryption_cenc`: %s", err)
			}

			commonEncryptionCBCS := flattenCommonEncryptionCbcs(props.CommonEncryptionCbcs)
			if err := d.Set("common_encryption_cbcs", commonEncryptionCBCS); err != nil {
				return fmt.Errorf("flattening `common_encryption_cbcs`: %s", err)
			}

			envelopeEncryption := flattenEnvelopeEncryption(props.EnvelopeEncryption)
			if err := d.Set("envelope_encryption", envelopeEncryption); err != nil {
				return fmt.Errorf("flattening `envelope_encryption`: %s", err)
			}

			d.Set("default_content_key_policy_name", props.DefaultContentKeyPolicyName)
		}
	}

	return nil
}

func resourceMediaStreamingPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220801Client.StreamingPoliciesAndStreamingLocators
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := streamingpoliciesandstreaminglocators.ParseStreamingPolicyID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.StreamingPoliciesDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func enabledProtocolsSchema() *pluginsdk.Schema {
	// lintignore:XS003
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"dash": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
				},

				"download": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
				},

				"hls": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
				},

				"smooth_streaming": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
				},
			},
		},
	}
}

func defaultContentKeySchema() *pluginsdk.Schema {
	// lintignore:XS003
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"label": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"policy_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func expandEnabledProtocols(input []interface{}) *streamingpoliciesandstreaminglocators.EnabledProtocols {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	protocols := input[0].(map[string]interface{})

	dash := false
	if v := protocols["dash"]; v != nil {
		dash = v.(bool)
	}

	download := false
	if v := protocols["download"]; v != nil {
		download = v.(bool)
	}

	hls := false
	if v := protocols["hls"]; v != nil {
		hls = v.(bool)
	}

	smoothStreaming := false
	if v := protocols["smooth_streaming"]; v != nil {
		smoothStreaming = v.(bool)
	}

	return &streamingpoliciesandstreaminglocators.EnabledProtocols{
		Dash:            dash,
		Download:        download,
		Hls:             hls,
		SmoothStreaming: smoothStreaming,
	}
}

func expandCommonEncryptionCenc(input []interface{}) *streamingpoliciesandstreaminglocators.CommonEncryptionCenc {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	item := input[0].(map[string]interface{})

	result := &streamingpoliciesandstreaminglocators.CommonEncryptionCenc{
		ClearKeyEncryptionConfiguration: expandClearKeyEncryptionConfiguration(item["clear_key_encryption"].([]interface{})),
		ClearTracks:                     expandTrackSelections(item["clear_track"].(*pluginsdk.Set).List()),
		EnabledProtocols:                expandEnabledProtocols(item["enabled_protocols"].([]interface{})),
	}

	if v := item["default_content_key"].([]interface{}); len(v) != 0 && v[0] != nil {
		if result.ContentKeys == nil {
			result.ContentKeys = &streamingpoliciesandstreaminglocators.StreamingPolicyContentKeys{}
		}
		result.ContentKeys.DefaultKey = expandDefaultKey(v)
	}

	if v := item["content_key_to_track_mapping"].(*pluginsdk.Set).List(); len(v) != 0 && v[0] != nil {
		if result.ContentKeys == nil {
			result.ContentKeys = &streamingpoliciesandstreaminglocators.StreamingPolicyContentKeys{}
		}
		result.ContentKeys.KeyToTrackMappings = expandKeyToTrackMappings(v)
	}

	if v := item["drm_widevine_custom_license_acquisition_url_template"].(string); v != "" {
		if result.Drm == nil {
			result.Drm = &streamingpoliciesandstreaminglocators.CencDrmConfiguration{}
		}
		result.Drm.Widevine = &streamingpoliciesandstreaminglocators.StreamingPolicyWidevineConfiguration{
			CustomLicenseAcquisitionUrlTemplate: utils.String(v),
		}
	}

	if v := item["drm_playready"].([]interface{}); len(v) != 0 && v[0] != nil {
		if result.Drm == nil {
			result.Drm = &streamingpoliciesandstreaminglocators.CencDrmConfiguration{}
		}
		result.Drm.PlayReady = expandPlayReady(v)
	}

	return result
}

func expandCommonEncryptionCbcs(input []interface{}) *streamingpoliciesandstreaminglocators.CommonEncryptionCbcs {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	item := input[0].(map[string]interface{})

	var enabledProtocols *streamingpoliciesandstreaminglocators.EnabledProtocols
	if v := item["enabled_protocols"]; v != nil {
		enabledProtocols = expandEnabledProtocols(v.([]interface{}))
	}

	result := &streamingpoliciesandstreaminglocators.CommonEncryptionCbcs{
		ClearKeyEncryptionConfiguration: expandClearKeyEncryptionConfiguration(item["clear_key_encryption"].([]interface{})),
		EnabledProtocols:                enabledProtocols,
	}

	if v := item["default_content_key"].([]interface{}); len(v) != 0 && v[0] != nil {
		if result.ContentKeys == nil {
			result.ContentKeys = &streamingpoliciesandstreaminglocators.StreamingPolicyContentKeys{}
		}
		result.ContentKeys.DefaultKey = expandDefaultKey(v)
	}

	if v := item["drm_fairplay"].([]interface{}); len(v) != 0 && v[0] != nil {
		if result.Drm == nil {
			result.Drm = &streamingpoliciesandstreaminglocators.CbcsDrmConfiguration{}
		}
		result.Drm.FairPlay = expandFairPlay(v)
	}

	return result
}

func expandEnvelopeEncryption(input []interface{}) *streamingpoliciesandstreaminglocators.EnvelopeEncryption {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	item := input[0].(map[string]interface{})
	result := &streamingpoliciesandstreaminglocators.EnvelopeEncryption{
		EnabledProtocols: expandEnabledProtocols(item["enabled_protocols"].([]interface{})),
	}

	if v := item["default_content_key"].([]interface{}); len(v) != 0 && v[0] != nil {
		if result.ContentKeys == nil {
			result.ContentKeys = &streamingpoliciesandstreaminglocators.StreamingPolicyContentKeys{}
		}
		result.ContentKeys.DefaultKey = expandDefaultKey(v)
	}

	if v := item["custom_keys_acquisition_url_template"].(string); v != "" {
		result.CustomKeyAcquisitionUrlTemplate = utils.String(v)
	}

	return result
}

func expandPlayReady(input []interface{}) *streamingpoliciesandstreaminglocators.StreamingPolicyPlayReadyConfiguration {
	if len(input) == 0 {
		return nil
	}

	playReady := input[0].(map[string]interface{})

	result := &streamingpoliciesandstreaminglocators.StreamingPolicyPlayReadyConfiguration{}

	if v := playReady["custom_license_acquisition_url_template"].(string); v != "" {
		result.CustomLicenseAcquisitionUrlTemplate = utils.String(v)
	}

	if v := playReady["custom_attributes"].(string); v != "" {
		result.PlayReadyCustomAttributes = utils.String(v)
	}

	return result
}

func expandDefaultKey(input []interface{}) *streamingpoliciesandstreaminglocators.DefaultKey {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	defaultKey := input[0].(map[string]interface{})
	defaultKeyResult := &streamingpoliciesandstreaminglocators.DefaultKey{}

	if v := defaultKey["policy_name"].(string); v != "" {
		defaultKeyResult.PolicyName = utils.String(v)
	}

	if v := defaultKey["label"].(string); v != "" {
		defaultKeyResult.Label = utils.String(v)
	}

	return defaultKeyResult
}

func expandClearKeyEncryptionConfiguration(input []interface{}) *streamingpoliciesandstreaminglocators.ClearKeyEncryptionConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	clearKeyEncryptionConfiguration := input[0].(map[string]interface{})
	result := streamingpoliciesandstreaminglocators.ClearKeyEncryptionConfiguration{}

	if v := clearKeyEncryptionConfiguration["custom_keys_acquisition_url_template"].(string); v != "" {
		result.CustomKeysAcquisitionUrlTemplate = utils.String(v)
	}

	return &result
}

func expandTrackSelections(input []interface{}) *[]streamingpoliciesandstreaminglocators.TrackSelection {
	if len(input) == 0 {
		return nil
	}

	result := make([]streamingpoliciesandstreaminglocators.TrackSelection, 0)
	for _, v := range input {
		selections := v.(map[string]interface{})
		conditions := expandTrackPropertyConditions(selections["condition"].(*pluginsdk.Set).List())

		result = append(result, streamingpoliciesandstreaminglocators.TrackSelection{
			TrackSelections: conditions,
		})
	}

	return &result
}

func expandTrackPropertyConditions(input []interface{}) *[]streamingpoliciesandstreaminglocators.TrackPropertyCondition {
	if len(input) == 0 {
		return nil
	}

	result := make([]streamingpoliciesandstreaminglocators.TrackPropertyCondition, 0)
	for _, c := range input {
		conditionRaw := c.(map[string]interface{})
		condition := streamingpoliciesandstreaminglocators.TrackPropertyCondition{
			Operation: streamingpoliciesandstreaminglocators.TrackPropertyCompareOperation(conditionRaw["operation"].(string)),
			Property:  streamingpoliciesandstreaminglocators.TrackPropertyType(conditionRaw["property"].(string)),
			Value:     utils.String(conditionRaw["value"].(string)),
		}

		result = append(result, condition)
	}
	return &result
}

func expandKeyToTrackMappings(input []interface{}) *[]streamingpoliciesandstreaminglocators.StreamingPolicyContentKey {
	if len(input) == 0 {
		return nil
	}

	result := make([]streamingpoliciesandstreaminglocators.StreamingPolicyContentKey, 0)
	for _, v := range input {
		mappings := v.(map[string]interface{})
		contentKey := streamingpoliciesandstreaminglocators.StreamingPolicyContentKey{
			Tracks: expandTrackSelections(mappings["track"].(*pluginsdk.Set).List()),
		}

		if label := mappings["label"].(string); label != "" {
			contentKey.Label = utils.String(label)
		}

		if policyName := mappings["policy_name"].(string); policyName != "" {
			contentKey.PolicyName = utils.String(policyName)
		}

		result = append(result, contentKey)
	}

	return &result
}

func expandFairPlay(input []interface{}) *streamingpoliciesandstreaminglocators.StreamingPolicyFairPlayConfiguration {
	if len(input) == 0 {
		return nil
	}

	fairPlay := input[0].(map[string]interface{})

	allowPersistentLicense := false
	if v := fairPlay["allow_persistent_license"]; v != nil {
		allowPersistentLicense = v.(bool)
	}

	result := &streamingpoliciesandstreaminglocators.StreamingPolicyFairPlayConfiguration{
		AllowPersistentLicense: allowPersistentLicense,
	}

	if v := fairPlay["custom_license_acquisition_url_template"].(string); v != "" {
		result.CustomLicenseAcquisitionUrlTemplate = utils.String(v)
	}

	return result
}

func flattenNoEncryption(input *streamingpoliciesandstreaminglocators.NoEncryption) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return flattenEnabledProtocols(input.EnabledProtocols)
}

func flattenEnabledProtocols(input *streamingpoliciesandstreaminglocators.EnabledProtocols) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"dash":             input.Dash,
			"download":         input.Download,
			"hls":              input.Hls,
			"smooth_streaming": input.SmoothStreaming,
		},
	}
}

func flattenCommonEncryptionCenc(input *streamingpoliciesandstreaminglocators.CommonEncryptionCenc) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	enabledProtocols := make([]interface{}, 0)
	if input.EnabledProtocols != nil {
		enabledProtocols = flattenEnabledProtocols(input.EnabledProtocols)
	}

	widevineTemplate := ""
	if input.Drm != nil && input.Drm.Widevine != nil && input.Drm.Widevine.CustomLicenseAcquisitionUrlTemplate != nil {
		widevineTemplate = *input.Drm.Widevine.CustomLicenseAcquisitionUrlTemplate
	}

	keyToTrackMappings := make([]interface{}, 0)
	if input.ContentKeys != nil && input.ContentKeys.KeyToTrackMappings != nil {
		keyToTrackMappings = flattenKeyToTrackMappings(input.ContentKeys.KeyToTrackMappings)
	}

	drmPlayReady := make([]interface{}, 0)
	if input.Drm != nil && input.Drm.PlayReady != nil {
		drmPlayReady = flattenPlayReady(input.Drm.PlayReady)
	}

	defaultContentKey := make([]interface{}, 0)
	if input.ContentKeys != nil && input.ContentKeys.DefaultKey != nil {
		defaultContentKey = flattenContentKey(input.ContentKeys.DefaultKey)
	}

	return []interface{}{
		map[string]interface{}{
			"clear_key_encryption":         flattenClearKeyEncryptionConfiguration(input.ClearKeyEncryptionConfiguration),
			"clear_track":                  flattenTrackSelections(input.ClearTracks),
			"content_key_to_track_mapping": keyToTrackMappings,
			"enabled_protocols":            enabledProtocols,
			"drm_widevine_custom_license_acquisition_url_template": widevineTemplate,
			"drm_playready":       drmPlayReady,
			"default_content_key": defaultContentKey,
		},
	}
}

func flattenCommonEncryptionCbcs(input *streamingpoliciesandstreaminglocators.CommonEncryptionCbcs) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	enabledProtocols := make([]interface{}, 0)
	if input.EnabledProtocols != nil {
		enabledProtocols = flattenEnabledProtocols(input.EnabledProtocols)
	}

	defaultContentKey := make([]interface{}, 0)
	if input.ContentKeys != nil && input.ContentKeys.DefaultKey != nil {
		defaultContentKey = flattenContentKey(input.ContentKeys.DefaultKey)
	}

	drmFairPlay := make([]interface{}, 0)
	if input.Drm != nil && input.Drm.FairPlay != nil {
		drmFairPlay = flattenFairPlay(input.Drm.FairPlay)
	}

	return []interface{}{
		map[string]interface{}{
			"clear_key_encryption": flattenClearKeyEncryptionConfiguration(input.ClearKeyEncryptionConfiguration),
			"enabled_protocols":    enabledProtocols,
			"default_content_key":  defaultContentKey,
			"drm_fairplay":         drmFairPlay,
		},
	}
}

func flattenEnvelopeEncryption(input *streamingpoliciesandstreaminglocators.EnvelopeEncryption) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	customKeyAcquisitionUrlTemplate := ""
	if input.CustomKeyAcquisitionUrlTemplate != nil {
		customKeyAcquisitionUrlTemplate = *input.CustomKeyAcquisitionUrlTemplate
	}

	defaultContentKey := make([]interface{}, 0)
	if input.ContentKeys != nil && input.ContentKeys.DefaultKey != nil {
		defaultContentKey = flattenContentKey(input.ContentKeys.DefaultKey)
	}

	return []interface{}{
		map[string]interface{}{
			"custom_keys_acquisition_url_template": customKeyAcquisitionUrlTemplate,
			"default_content_key":                  defaultContentKey,
			"enabled_protocols":                    flattenEnabledProtocols(input.EnabledProtocols),
		},
	}
}

func flattenPlayReady(input *streamingpoliciesandstreaminglocators.StreamingPolicyPlayReadyConfiguration) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	customAttributes := ""
	if input.PlayReadyCustomAttributes != nil {
		customAttributes = *input.PlayReadyCustomAttributes
	}

	customLicenseURLTemplate := ""
	if input.CustomLicenseAcquisitionUrlTemplate != nil {
		customLicenseURLTemplate = *input.CustomLicenseAcquisitionUrlTemplate
	}

	return []interface{}{
		map[string]interface{}{
			"custom_attributes":                       customAttributes,
			"custom_license_acquisition_url_template": customLicenseURLTemplate,
		},
	}
}

func flattenClearKeyEncryptionConfiguration(input *streamingpoliciesandstreaminglocators.ClearKeyEncryptionConfiguration) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	customKeysAcquisitionUrlTemplate := ""
	if input.CustomKeysAcquisitionUrlTemplate != nil {
		customKeysAcquisitionUrlTemplate = *input.CustomKeysAcquisitionUrlTemplate
	}

	return []interface{}{
		map[string]interface{}{
			"custom_keys_acquisition_url_template": customKeysAcquisitionUrlTemplate,
		},
	}
}

func flattenTrackSelections(input *[]streamingpoliciesandstreaminglocators.TrackSelection) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make([]interface{}, 0)
	for _, v := range *input {
		result = append(result, map[string]interface{}{
			"condition": flattenTrackPropertyConditions(v.TrackSelections),
		})
	}

	return result
}

func flattenTrackPropertyConditions(input *[]streamingpoliciesandstreaminglocators.TrackPropertyCondition) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make([]interface{}, 0)
	for _, v := range *input {
		value := ""
		if v.Value != nil {
			value = *v.Value
		}
		result = append(result, map[string]interface{}{
			"operation": string(v.Operation),
			"property":  string(v.Property),
			"value":     value,
		})
	}

	return result
}

func flattenKeyToTrackMappings(input *[]streamingpoliciesandstreaminglocators.StreamingPolicyContentKey) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}
	result := make([]interface{}, 0)
	for _, v := range *input {
		label := ""
		if v.Label != nil {
			label = *v.Label
		}

		policyName := ""
		if v.PolicyName != nil {
			policyName = *v.PolicyName
		}

		result = append(result, map[string]interface{}{
			"label":       label,
			"policy_name": policyName,
			"track":       flattenTrackSelections(v.Tracks),
		})
	}

	return result
}

func flattenFairPlay(input *streamingpoliciesandstreaminglocators.StreamingPolicyFairPlayConfiguration) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	customLicenseURLTemplate := ""
	if input.CustomLicenseAcquisitionUrlTemplate != nil {
		customLicenseURLTemplate = *input.CustomLicenseAcquisitionUrlTemplate
	}

	return []interface{}{
		map[string]interface{}{
			"allow_persistent_license":                input.AllowPersistentLicense,
			"custom_license_acquisition_url_template": customLicenseURLTemplate,
		},
	}
}

func flattenContentKey(input *streamingpoliciesandstreaminglocators.DefaultKey) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	policyName := ""
	if input.PolicyName != nil {
		policyName = *input.PolicyName
	}

	label := ""
	if input.Label != nil {
		label = *input.Label
	}

	if label == "" && policyName == "" {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"policy_name": policyName,
			"label":       label,
		},
	}
}
