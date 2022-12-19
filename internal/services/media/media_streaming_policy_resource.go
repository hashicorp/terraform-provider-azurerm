package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/streamingpoliciesandstreaminglocators"
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
						"enabled_protocols": enabledProtocolsSchema(),

						"drm_widevine_custom_license_acquisition_url_template": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IsURLWithHTTPS,
						},

						"drm_playready": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
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
						"enabled_protocols": enabledProtocolsSchema(),

						"drm_fairplay": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
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
		payload.Properties.NoEncryption = expandNoEncryption(noEncryption.([]interface{}))
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
	d.Set("media_services_account_name", id.AccountName)
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

func expandNoEncryption(input []interface{}) *streamingpoliciesandstreaminglocators.NoEncryption {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	noEncryption := input[0].(map[string]interface{})

	return &streamingpoliciesandstreaminglocators.NoEncryption{
		EnabledProtocols: expandEnabledProtocols(noEncryption),
	}
}

func expandEnabledProtocols(input map[string]interface{}) *streamingpoliciesandstreaminglocators.EnabledProtocols {
	if len(input) == 0 {
		return nil
	}

	dash := false
	if v := input["dash"]; v != nil {
		dash = v.(bool)
	}

	download := false
	if v := input["download"]; v != nil {
		download = v.(bool)
	}

	hls := false
	if v := input["hls"]; v != nil {
		hls = v.(bool)
	}

	smoothStreaming := false
	if v := input["smooth_streaming"]; v != nil {
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
	if len(input) == 0 {
		return nil
	}

	item := input[0].(map[string]interface{})

	var enabledProtocols *streamingpoliciesandstreaminglocators.EnabledProtocols
	if v := item["enabled_protocols"]; v != nil {
		protocols := v.([]interface{})
		if len(protocols) != 0 && protocols[0] != nil {
			enabledProtocols = expandEnabledProtocols(protocols[0].(map[string]interface{}))
		}
	}

	drmWidevineTemplate := ""
	if v := item["drm_widevine_custom_license_acquisition_url_template"]; v != nil {
		drmWidevineTemplate = v.(string)
	}

	var drmPlayReady *streamingpoliciesandstreaminglocators.StreamingPolicyPlayReadyConfiguration
	if v := item["drm_playready"]; v != nil {
		drmPlayReady = expandPlayReady(v.([]interface{}))
	}

	var defaultKey *streamingpoliciesandstreaminglocators.DefaultKey
	if v := item["default_content_key"]; v != nil {
		defaultKey = expandDefaultKey(v.([]interface{}))
	}

	return &streamingpoliciesandstreaminglocators.CommonEncryptionCenc{
		EnabledProtocols: enabledProtocols,
		Drm: &streamingpoliciesandstreaminglocators.CencDrmConfiguration{
			Widevine: &streamingpoliciesandstreaminglocators.StreamingPolicyWidevineConfiguration{
				CustomLicenseAcquisitionUrlTemplate: utils.String(drmWidevineTemplate),
			},
			PlayReady: drmPlayReady,
		},
		ContentKeys: &streamingpoliciesandstreaminglocators.StreamingPolicyContentKeys{
			DefaultKey: defaultKey,
		},
	}
}

func expandCommonEncryptionCbcs(input []interface{}) *streamingpoliciesandstreaminglocators.CommonEncryptionCbcs {
	if len(input) == 0 {
		return nil
	}

	item := input[0].(map[string]interface{})

	var enabledProtocols *streamingpoliciesandstreaminglocators.EnabledProtocols
	if v := item["enabled_protocols"]; v != nil {
		protocols := v.([]interface{})
		if len(protocols) != 0 && protocols[0] != nil {
			enabledProtocols = expandEnabledProtocols(protocols[0].(map[string]interface{}))
		}
	}

	var defaultKey *streamingpoliciesandstreaminglocators.DefaultKey
	if v := item["default_content_key"]; v != nil {
		defaultKey = expandDefaultKey(v.([]interface{}))
	}

	var drmFairPlay *streamingpoliciesandstreaminglocators.StreamingPolicyFairPlayConfiguration
	if v := item["drm_fairplay"]; v != nil {
		drmFairPlay = expandFairPlay(v.([]interface{}))
	}

	return &streamingpoliciesandstreaminglocators.CommonEncryptionCbcs{
		EnabledProtocols: enabledProtocols,
		Drm: &streamingpoliciesandstreaminglocators.CbcsDrmConfiguration{
			FairPlay: drmFairPlay,
		},
		ContentKeys: &streamingpoliciesandstreaminglocators.StreamingPolicyContentKeys{
			DefaultKey: defaultKey,
		},
	}
}

func expandPlayReady(input []interface{}) *streamingpoliciesandstreaminglocators.StreamingPolicyPlayReadyConfiguration {
	if len(input) == 0 {
		return nil
	}

	playReady := input[0].(map[string]interface{})

	customLicenseURLTemplate := ""
	if v := playReady["custom_license_acquisition_url_template"]; v != nil {
		customLicenseURLTemplate = v.(string)
	}

	customAttributes := ""
	if v := playReady["custom_attributes"]; v != nil {
		customAttributes = v.(string)
	}

	return &streamingpoliciesandstreaminglocators.StreamingPolicyPlayReadyConfiguration{
		CustomLicenseAcquisitionUrlTemplate: utils.String(customLicenseURLTemplate),
		PlayReadyCustomAttributes:           utils.String(customAttributes),
	}
}

func expandDefaultKey(input []interface{}) *streamingpoliciesandstreaminglocators.DefaultKey {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	defaultKey := input[0].(map[string]interface{})
	defaultKeyResult := &streamingpoliciesandstreaminglocators.DefaultKey{}

	if v := defaultKey["policy_name"]; v != nil {
		defaultKeyResult.PolicyName = utils.String(v.(string))
	}

	if v := defaultKey["label"]; v != nil {
		defaultKeyResult.Label = utils.String(v.(string))
	}

	return defaultKeyResult
}

func expandFairPlay(input []interface{}) *streamingpoliciesandstreaminglocators.StreamingPolicyFairPlayConfiguration {
	if len(input) == 0 {
		return nil
	}

	fairPlay := input[0].(map[string]interface{})

	customLicenseURLTemplate := ""
	if v := fairPlay["custom_license_acquisition_url_template"]; v != nil {
		customLicenseURLTemplate = v.(string)
	}

	allowPersistentLicense := false
	if v := fairPlay["allow_persistent_license"]; v != nil {
		allowPersistentLicense = v.(bool)
	}

	return &streamingpoliciesandstreaminglocators.StreamingPolicyFairPlayConfiguration{
		CustomLicenseAcquisitionUrlTemplate: utils.String(customLicenseURLTemplate),
		AllowPersistentLicense:              allowPersistentLicense,
	}
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
			"enabled_protocols": enabledProtocols,
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
			"enabled_protocols":   enabledProtocols,
			"default_content_key": defaultContentKey,
			"drm_fairplay":        drmFairPlay,
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
