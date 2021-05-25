package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2021-05-01/media"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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
			_, err := parse.StreamingPolicyID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9(_)]{1,128}$"),
					"Steraming Policy name must be 1 - 128 characters long, can contain letters, numbers, underscores, and hyphens (but the first and last character must be a letter or number).",
				),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

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
	client := meta.(*clients.Client).Media.StreamingPoliciesClient
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceID := parse.NewStreamingPolicyID(subscriptionID, d.Get("resource_group_name").(string), d.Get("media_services_account_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceID.ResourceGroup, resourceID.MediaserviceName, resourceID.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", resourceID, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_media_streaming_policy", resourceID.ID())
		}
	}

	parameters := media.StreamingPolicy{
		StreamingPolicyProperties: &media.StreamingPolicyProperties{},
	}

	if noEncryption, ok := d.GetOk("no_encryption_enabled_protocols"); ok {
		parameters.NoEncryption = expandNoEncryption(noEncryption.([]interface{}))
	}

	if commonEncryptionCENC, ok := d.GetOk("common_encryption_cenc"); ok {
		parameters.CommonEncryptionCenc = expandCommonEncryptionCenc(commonEncryptionCENC.([]interface{}))
	}

	if commonEncryptionCBCS, ok := d.GetOk("common_encryption_cbcs"); ok {
		parameters.CommonEncryptionCbcs = expandCommonEncryptionCbcs(commonEncryptionCBCS.([]interface{}))
	}

	if contentKeyPolicyName, ok := d.GetOk("default_content_key_policy_name"); ok {
		parameters.DefaultContentKeyPolicyName = utils.String(contentKeyPolicyName.(string))
	}

	if _, err := client.Create(ctx, resourceID.ResourceGroup, resourceID.MediaserviceName, resourceID.Name, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", resourceID, err)
	}

	d.SetId(resourceID.ID())

	return resourceMediaStreamingPolicyRead(d, meta)
}

func resourceMediaStreamingPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.StreamingPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StreamingPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("media_services_account_name", id.MediaserviceName)

	if props := resp.StreamingPolicyProperties; props != nil {
		noEncryption := flattenNoEncryption(resp.NoEncryption)
		if err := d.Set("no_encryption_enabled_protocols", noEncryption); err != nil {
			return fmt.Errorf("Error flattening `no_encryption_enabled_protocols`: %s", err)
		}

		commonEncryptionCENC := flattenCommonEncryptionCenc(resp.CommonEncryptionCenc)
		if err := d.Set("common_encryption_cenc", commonEncryptionCENC); err != nil {
			return fmt.Errorf("Error flattening `common_encryption_cenc`: %s", err)
		}

		commonEncryptionCBCS := flattenCommonEncryptionCbcs(resp.CommonEncryptionCbcs)
		if err := d.Set("common_encryption_cbcs", commonEncryptionCBCS); err != nil {
			return fmt.Errorf("Error flattening `common_encryption_cbcs`: %s", err)
		}

		d.Set("default_content_key_policy_name", props.DefaultContentKeyPolicyName)
	}

	return nil
}

func resourceMediaStreamingPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.StreamingPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StreamingPolicyID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, id.ResourceGroup, id.MediaserviceName, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func enabledProtocolsSchema() *pluginsdk.Schema {
	//lintignore:XS003
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
	//lintignore:XS003
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

func expandNoEncryption(input []interface{}) *media.NoEncryption {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	noEncryption := input[0].(map[string]interface{})

	return &media.NoEncryption{
		EnabledProtocols: expandEnabledProtocols(noEncryption),
	}
}

func expandEnabledProtocols(input map[string]interface{}) *media.EnabledProtocols {
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

	return &media.EnabledProtocols{
		Dash:            utils.Bool(dash),
		Download:        utils.Bool(download),
		Hls:             utils.Bool(hls),
		SmoothStreaming: utils.Bool(smoothStreaming),
	}
}

func expandCommonEncryptionCenc(input []interface{}) *media.CommonEncryptionCenc {
	if len(input) == 0 {
		return nil
	}

	CommonEncryptionCenc := input[0].(map[string]interface{})

	var enabledProtocols *media.EnabledProtocols
	if v := CommonEncryptionCenc["enabled_protocols"]; v != nil {
		protocols := v.([]interface{})
		if len(protocols) != 0 && protocols[0] != nil {
			enabledProtocols = expandEnabledProtocols(protocols[0].(map[string]interface{}))
		}
	}

	drmWidevineTemplate := ""
	if v := CommonEncryptionCenc["drm_widevine_custom_license_acquisition_url_template"]; v != nil {
		drmWidevineTemplate = v.(string)
	}

	var drmPlayReady *media.StreamingPolicyPlayReadyConfiguration
	if v := CommonEncryptionCenc["drm_playready"]; v != nil {
		drmPlayReady = expandPlayReady(v.([]interface{}))
	}

	var defaultKey *media.DefaultKey
	if v := CommonEncryptionCenc["default_content_key"]; v != nil {
		defaultKey = expandDefaultKey(v.([]interface{}))
	}

	return &media.CommonEncryptionCenc{
		EnabledProtocols: enabledProtocols,
		Drm: &media.CencDrmConfiguration{
			Widevine: &media.StreamingPolicyWidevineConfiguration{
				CustomLicenseAcquisitionURLTemplate: utils.String(drmWidevineTemplate),
			},
			PlayReady: drmPlayReady,
		},
		ContentKeys: &media.StreamingPolicyContentKeys{
			DefaultKey: defaultKey,
		},
	}
}

func expandCommonEncryptionCbcs(input []interface{}) *media.CommonEncryptionCbcs {
	if len(input) == 0 {
		return nil
	}

	CommonEncryptionCenc := input[0].(map[string]interface{})

	var enabledProtocols *media.EnabledProtocols
	if v := CommonEncryptionCenc["enabled_protocols"]; v != nil {
		protocols := v.([]interface{})
		if len(protocols) != 0 && protocols[0] != nil {
			enabledProtocols = expandEnabledProtocols(protocols[0].(map[string]interface{}))
		}
	}

	var defaultKey *media.DefaultKey
	if v := CommonEncryptionCenc["default_content_key"]; v != nil {
		defaultKey = expandDefaultKey(v.([]interface{}))
	}

	var drmFairPlay *media.StreamingPolicyFairPlayConfiguration
	if v := CommonEncryptionCenc["drm_fairplay"]; v != nil {
		drmFairPlay = expandFairPlay(v.([]interface{}))
	}

	return &media.CommonEncryptionCbcs{
		EnabledProtocols: enabledProtocols,
		Drm: &media.CbcsDrmConfiguration{
			FairPlay: drmFairPlay,
		},
		ContentKeys: &media.StreamingPolicyContentKeys{
			DefaultKey: defaultKey,
		},
	}
}

func expandPlayReady(input []interface{}) *media.StreamingPolicyPlayReadyConfiguration {
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

	return &media.StreamingPolicyPlayReadyConfiguration{
		CustomLicenseAcquisitionURLTemplate: utils.String(customLicenseURLTemplate),
		PlayReadyCustomAttributes:           utils.String(customAttributes),
	}
}

func expandDefaultKey(input []interface{}) *media.DefaultKey {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	defaultKey := input[0].(map[string]interface{})
	defaultKeyResult := &media.DefaultKey{}

	if v := defaultKey["policy_name"]; v != nil {
		defaultKeyResult.PolicyName = utils.String(v.(string))
	}

	if v := defaultKey["label"]; v != nil {
		defaultKeyResult.Label = utils.String(v.(string))
	}

	return defaultKeyResult
}

func expandFairPlay(input []interface{}) *media.StreamingPolicyFairPlayConfiguration {
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

	return &media.StreamingPolicyFairPlayConfiguration{
		CustomLicenseAcquisitionURLTemplate: utils.String(customLicenseURLTemplate),
		AllowPersistentLicense:              utils.Bool(allowPersistentLicense),
	}
}

func flattenNoEncryption(input *media.NoEncryption) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return flattenEnabledProtocols(input.EnabledProtocols)
}

func flattenEnabledProtocols(input *media.EnabledProtocols) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	dash := false
	if input.Dash != nil {
		dash = *input.Dash
	}

	download := false
	if input.Download != nil {
		download = *input.Download
	}

	hls := false
	if input.Hls != nil {
		hls = *input.Hls
	}

	smoothStreaming := false
	if input.SmoothStreaming != nil {
		smoothStreaming = *input.SmoothStreaming
	}

	return []interface{}{
		map[string]interface{}{
			"dash":             dash,
			"download":         download,
			"hls":              hls,
			"smooth_streaming": smoothStreaming,
		},
	}
}

func flattenCommonEncryptionCenc(input *media.CommonEncryptionCenc) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	enabledProtocols := make([]interface{}, 0)
	if input.EnabledProtocols != nil {
		enabledProtocols = flattenEnabledProtocols(input.EnabledProtocols)
	}

	widevineTemplate := ""
	if input.Drm != nil && input.Drm.Widevine != nil && input.Drm.Widevine.CustomLicenseAcquisitionURLTemplate != nil {
		widevineTemplate = *input.Drm.Widevine.CustomLicenseAcquisitionURLTemplate
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

func flattenCommonEncryptionCbcs(input *media.CommonEncryptionCbcs) []interface{} {
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

func flattenPlayReady(input *media.StreamingPolicyPlayReadyConfiguration) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	customAttributes := ""
	if input.PlayReadyCustomAttributes != nil {
		customAttributes = *input.PlayReadyCustomAttributes
	}

	customLicenseURLTemplate := ""
	if input.CustomLicenseAcquisitionURLTemplate != nil {
		customLicenseURLTemplate = *input.CustomLicenseAcquisitionURLTemplate
	}

	return []interface{}{
		map[string]interface{}{
			"custom_attributes":                       customAttributes,
			"custom_license_acquisition_url_template": customLicenseURLTemplate,
		},
	}
}

func flattenFairPlay(input *media.StreamingPolicyFairPlayConfiguration) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	allowPersistentLicense := false
	if input.AllowPersistentLicense != nil {
		allowPersistentLicense = *input.AllowPersistentLicense
	}

	customLicenseURLTemplate := ""
	if input.CustomLicenseAcquisitionURLTemplate != nil {
		customLicenseURLTemplate = *input.CustomLicenseAcquisitionURLTemplate
	}

	return []interface{}{
		map[string]interface{}{
			"allow_persistent_license":                allowPersistentLicense,
			"custom_license_acquisition_url_template": customLicenseURLTemplate,
		},
	}
}

func flattenContentKey(input *media.DefaultKey) []interface{} {
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
