package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2020-05-01/media"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	uuid "github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMediaContentKeyPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMediaContentKeyPolicyCreateUpdate,
		Read:   resourceMediaContentKeyPolicyRead,
		Update: resourceMediaContentKeyPolicyCreateUpdate,
		Delete: resourceMediaContentKeyPolicyDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ContentKeyPolicyID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9(_)]{1,128}$"),
					"Content Key Policy name must be 1 - 128 characters long, can contain letters, numbers, underscores, and hyphens (but the first and last character must be a letter or number).",
				),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"media_services_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateMediaServicesAccountName,
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"policy_option": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"clear_key_configuration_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"widevine_configuration_template": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"playready_configuration_license": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow_test_devices": {
										Type:     schema.TypeBool,
										Optional: true,
									},

									"begin_date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsRFC3339Time,
									},

									"content_key_location_from_header_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},

									"content_key_location_from_key_id": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsUUID,
									},

									"content_type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(media.ContentKeyPolicyPlayReadyContentTypeUltraVioletDownload),
											string(media.ContentKeyPolicyPlayReadyContentTypeUltraVioletStreaming),
											string(media.ContentKeyPolicyPlayReadyContentTypeUnspecified),
											string(media.ContentKeyPolicyPlayReadyContentTypeUnknown),
										}, false),
									},

									"expiration_date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsRFC3339Time,
									},

									"grace_period": {
										Type:         schema.TypeString,
										Optional:     true,
										Sensitive:    true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"license_type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(media.ContentKeyPolicyPlayReadyLicenseTypeNonPersistent),
											string(media.ContentKeyPolicyPlayReadyLicenseTypePersistent),
											string(media.ContentKeyPolicyPlayReadyLicenseTypeUnknown),
										}, false),
									},

									"play_right": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"agc_and_color_stripe_restriction": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntBetween(0, 3),
												},

												"allow_passing_video_content_to_unknown_output": {
													Type:     schema.TypeString,
													Optional: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(media.ContentKeyPolicyPlayReadyUnknownOutputPassingOptionAllowed),
														string(media.ContentKeyPolicyPlayReadyUnknownOutputPassingOptionAllowedWithVideoConstriction),
														string(media.ContentKeyPolicyPlayReadyUnknownOutputPassingOptionNotAllowed),
														string(media.ContentKeyPolicyPlayReadyUnknownOutputPassingOptionUnknown),
													}, false),
												},

												"analog_video_opl": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntInSlice([]int{100, 150, 200}),
												},

												"compressed_digital_audio_opl": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntInSlice([]int{100, 150, 200}),
												},

												"digital_video_only_content_restriction": {
													Type:     schema.TypeBool,
													Optional: true,
												},

												"first_play_expiration": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},

												"image_constraint_for_analog_component_video_restriction": {
													Type:     schema.TypeBool,
													Optional: true,
												},

												"image_constraint_for_analog_computer_monitor_restriction": {
													Type:     schema.TypeBool,
													Optional: true,
												},

												"scms_restriction": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntBetween(0, 3),
												},

												"uncompressed_digital_audio_opl": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntInSlice([]int{100, 150, 250, 300}),
												},

												"uncompressed_digital_video_opl": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntInSlice([]int{100, 250, 270, 300}),
												},
											},
										},
									},
									"relative_begin_date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsRFC3339Time,
									},

									"relative_expiration_date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsRFC3339Time,
									},
								},
							},
						},
						"fairplay_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ask": {
										Type:         schema.TypeString,
										Optional:     true,
										Sensitive:    true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"pfx": {
										Type:         schema.TypeString,
										Optional:     true,
										Sensitive:    true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"pfx_password": {
										Type:         schema.TypeString,
										Optional:     true,
										Sensitive:    true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"offline_rental_configuration": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"playback_duration_seconds": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntAtLeast(1),
												},
												"storage_duration_seconds": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntAtLeast(1),
												},
											},
										},
									},
									"rental_and_lease_key_type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(media.DualExpiry),
											string(media.PersistentLimited),
											string(media.PersistentUnlimited),
											string(media.Undefined),
											string(media.Unknown),
										}, false),
									},
									"rental_duration_seconds": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(1),
									},
								},
							},
						},
						"token_restriction": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"audience": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"issuer": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"token_type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(media.ContentKeyPolicyRestrictionTokenTypeJwt),
											string(media.ContentKeyPolicyRestrictionTokenTypeSwt),
											string(media.ContentKeyPolicyRestrictionTokenTypeUnknown),
										}, false),
									},
									"primary_symmetric_token_key": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsBase64,
										Sensitive:    true,
									},
									"primary_rsa_token_key_exponent": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
										Sensitive:    true,
									},
									"primary_rsa_token_key_modulus": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
										Sensitive:    true,
									},
									"primary_x509_token_key_raw": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
										Sensitive:    true,
									},
									"open_id_connect_discovery_document": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"required_claim": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"value": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
										},
									},
								},
							},
						},
						"open_restriction_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceMediaContentKeyPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.ContentKeyPoliciesClient
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceID := parse.NewContentKeyPolicyID(subscriptionID, d.Get("resource_group_name").(string), d.Get("media_services_account_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceID.ResourceGroup, resourceID.MediaserviceName, resourceID.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of %s: %+v", resourceID, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_media_content_key_policy", resourceID.ID())
		}
	}

	parameters := media.ContentKeyPolicy{
		ContentKeyPolicyProperties: &media.ContentKeyPolicyProperties{},
	}

	if description, ok := d.GetOk("description"); ok {
		parameters.ContentKeyPolicyProperties.Description = utils.String(description.(string))
	}

	if v, ok := d.GetOk("policy_option"); ok {
		options, err := expandPolicyOptions(v.(*schema.Set).List())
		if err != nil {
			return err
		}
		parameters.ContentKeyPolicyProperties.Options = options
	}

	_, err := client.CreateOrUpdate(ctx, resourceID.ResourceGroup, resourceID.MediaserviceName, resourceID.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", resourceID, err)
	}

	d.SetId(resourceID.ID())

	return resourceMediaContentKeyPolicyRead(d, meta)
}

func resourceMediaContentKeyPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.ContentKeyPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ContentKeyPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetPolicyPropertiesWithSecrets(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
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
	d.Set("description", resp.Description)

	options, err := flattenPolicyOptions(resp.Options)
	if err != nil {
		return err
	}

	d.Set("policy_option", options)

	return nil
}

func resourceMediaContentKeyPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.ContentKeyPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ContentKeyPolicyID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandPolicyOptions(input []interface{}) (*[]media.ContentKeyPolicyOption, error) {
	results := make([]media.ContentKeyPolicyOption, 0)

	for _, policyOptionRaw := range input {
		policyOption := policyOptionRaw.(map[string]interface{})

		restriction, err := expandRestriction(policyOption)
		if err != nil {
			return nil, err
		}

		configuration, err := expandConfiguration(policyOption)
		if err != nil {
			return nil, err
		}

		contentKeyPolicyOption := media.ContentKeyPolicyOption{
			Restriction:   restriction,
			Configuration: configuration,
		}

		if v := policyOption["name"]; v != nil {
			contentKeyPolicyOption.Name = utils.String(v.(string))
		}

		results = append(results, contentKeyPolicyOption)
	}

	return &results, nil
}

func flattenPolicyOptions(input *[]media.ContentKeyPolicyOption) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	results := make([]interface{}, 0)
	for _, option := range *input {
		policyOption := make(map[string]interface{})
		policyOption["name"] = option.Name
		policyOption = flattenConfiguration(option.Configuration, policyOption)

		restriction, err := flattenRestriction(option.Restriction, policyOption)
		if err != nil {
			return nil, err
		}
		policyOption = restriction

		results = append(results, policyOption)
	}

	return results, nil
}

func expandRestriction(option map[string]interface{}) (media.BasicContentKeyPolicyRestriction, error) {
	restrictionCount := 0
	restrictionType := ""
	if option["open_restriction_enabled"] != nil && option["open_restriction_enabled"].(bool) {
		restrictionCount++
		restrictionType = string(media.OdataTypeMicrosoftMediaContentKeyPolicyOpenRestriction)
	}
	if option["token_restriction"] != nil && len(option["token_restriction"].([]interface{})) > 0 {
		restrictionCount++
		restrictionType = string(media.OdataTypeMicrosoftMediaContentKeyPolicyTokenRestriction)
	}

	if restrictionCount == 0 {
		return nil, fmt.Errorf("policy_option must contain at least one type of restriction: open_restriction_enabled or token_restriction.")
	}

	if restrictionCount > 1 {
		return nil, fmt.Errorf("more than one type of restriction in the same policy_option is not allowed.")
	}

	switch restrictionType {
	case string(media.OdataTypeMicrosoftMediaContentKeyPolicyOpenRestriction):
		openRestriction := &media.ContentKeyPolicyOpenRestriction{
			OdataType: media.OdataTypeMicrosoftMediaContentKeyPolicyOpenRestriction,
		}
		return openRestriction, nil
	case string(media.OdataTypeMicrosoftMediaContentKeyPolicyTokenRestriction):
		tokenRestrictions := option["token_restriction"].([]interface{})
		tokenRestriction := tokenRestrictions[0].(map[string]interface{})
		contentKeyPolicyTokenRestriction := &media.ContentKeyPolicyTokenRestriction{
			OdataType: media.OdataTypeMicrosoftMediaContentKeyPolicyTokenRestriction,
		}
		if tokenRestriction["audience"] != nil && tokenRestriction["audience"].(string) != "" {
			contentKeyPolicyTokenRestriction.Audience = utils.String(tokenRestriction["audience"].(string))
		}
		if tokenRestriction["issuer"] != nil && tokenRestriction["issuer"].(string) != "" {
			contentKeyPolicyTokenRestriction.Issuer = utils.String(tokenRestriction["issuer"].(string))
		}
		if tokenRestriction["token_type"] != nil && tokenRestriction["token_type"].(string) != "" {
			contentKeyPolicyTokenRestriction.RestrictionTokenType = media.ContentKeyPolicyRestrictionTokenType(tokenRestriction["token_type"].(string))
		}
		if tokenRestriction["open_id_connect_discovery_document"] != nil && tokenRestriction["open_id_connect_discovery_document"].(string) != "" {
			contentKeyPolicyTokenRestriction.OpenIDConnectDiscoveryDocument = utils.String(tokenRestriction["open_id_connect_discovery_document"].(string))
		}
		if v := tokenRestriction["required_claim"]; v != nil {
			contentKeyPolicyTokenRestriction.RequiredClaims = expandRequiredClaims(v.([]interface{}))
		}
		primaryVerificationKey, err := expandVerificationKey(tokenRestriction)
		if err != nil {
			return nil, err
		}
		contentKeyPolicyTokenRestriction.PrimaryVerificationKey = primaryVerificationKey

		return contentKeyPolicyTokenRestriction, nil
	default:
		return nil, fmt.Errorf("policy_option must contain at least one type of restriction: open_restriction_enabled or token_restriction.")
	}
}

func flattenRestriction(input media.BasicContentKeyPolicyRestriction, option map[string]interface{}) (map[string]interface{}, error) {
	if input == nil {
		return option, nil
	}
	switch input.(type) {
	case media.ContentKeyPolicyOpenRestriction:
		option["open_restriction_enabled"] = true
	case media.ContentKeyPolicyTokenRestriction:
		token, _ := input.AsContentKeyPolicyTokenRestriction()

		audience := ""
		if token.Audience != nil {
			audience = *token.Audience
		}

		issuer := ""
		if token.Issuer != nil {
			issuer = *token.Issuer
		}

		openIDConnectDiscoveryDocument := ""
		if token.OpenIDConnectDiscoveryDocument != nil {
			openIDConnectDiscoveryDocument = *token.OpenIDConnectDiscoveryDocument
		}

		requiredClaims := flattenRequiredClaims(token.RequiredClaims)

		tokenRestriction := map[string]interface{}{
			"audience":                           audience,
			"issuer":                             issuer,
			"token_type":                         string(token.RestrictionTokenType),
			"open_id_connect_discovery_document": openIDConnectDiscoveryDocument,
			"required_claim":                     requiredClaims,
		}

		tokenRestriction = flattenVerificationKey(token.PrimaryVerificationKey, tokenRestriction)

		option["token_restriction"] = []interface{}{
			tokenRestriction,
		}
	}

	return option, nil
}

func expandConfiguration(input map[string]interface{}) (media.BasicContentKeyPolicyConfiguration, error) {
	configurationCount := 0
	configurationType := ""
	if input["clear_key_configuration_enabled"] != nil && input["clear_key_configuration_enabled"].(bool) {
		configurationCount++
		configurationType = string(media.OdataTypeMicrosoftMediaContentKeyPolicyClearKeyConfiguration)
	}
	if input["widevine_configuration_template"] != nil && input["widevine_configuration_template"].(string) != "" {
		configurationCount++
		configurationType = string(media.OdataTypeMicrosoftMediaContentKeyPolicyWidevineConfiguration)
	}
	if input["fairplay_configuration"] != nil && len(input["fairplay_configuration"].([]interface{})) > 0 {
		configurationCount++
		configurationType = string(media.OdataTypeMicrosoftMediaContentKeyPolicyFairPlayConfiguration)
	}

	if input["playready_configuration_license"] != nil && len(input["playready_configuration_license"].([]interface{})) > 0 {
		configurationCount++
		configurationType = string(media.OdataTypeMicrosoftMediaContentKeyPolicyPlayReadyConfiguration)
	}

	if configurationCount == 0 {
		return nil, fmt.Errorf("policy_option must contain at least one type of configuration: clear_key_configuration_enabled , widevine_configuration_template, playready_configuration_license or fairplay_configuration.")
	}

	if configurationCount > 1 {
		return nil, fmt.Errorf("more than one type of configuration in the same policy_option is not allowed.")
	}

	switch configurationType {
	case string(media.OdataTypeMicrosoftMediaContentKeyPolicyClearKeyConfiguration):
		clearKeyConfiguration := &media.ContentKeyPolicyClearKeyConfiguration{
			OdataType: media.OdataTypeMicrosoftMediaContentKeyPolicyClearKeyConfiguration,
		}
		return clearKeyConfiguration, nil
	case string(media.OdataTypeMicrosoftMediaContentKeyPolicyWidevineConfiguration):
		wideVineConfiguration := &media.ContentKeyPolicyWidevineConfiguration{
			OdataType:        media.OdataTypeMicrosoftMediaContentKeyPolicyWidevineConfiguration,
			WidevineTemplate: utils.String(input["widevine_configuration_template"].(string)),
		}
		return wideVineConfiguration, nil
	case string(media.OdataTypeMicrosoftMediaContentKeyPolicyFairPlayConfiguration):
		fairplayConfiguration := expandFairplayConfiguration(input["fairplay_configuration"].([]interface{}))
		return fairplayConfiguration, nil
	case string(media.OdataTypeMicrosoftMediaContentKeyPolicyPlayReadyConfiguration):
		playReadyConfiguration := &media.ContentKeyPolicyPlayReadyConfiguration{
			OdataType: media.OdataTypeMicrosoftMediaContentKeyPolicyPlayReadyConfiguration,
		}

		if input["playready_configuration_license"] != nil {
			licenses, err := expandPlayReadyLicenses(input["playready_configuration_license"].([]interface{}))
			if err != nil {
				return nil, err
			}
			playReadyConfiguration.Licenses = licenses
		}
		return playReadyConfiguration, nil

	default:
		return nil, fmt.Errorf("policy_option must contain at least one type of configuration: clear_key_configuration_enabled , widevine_configuration_template, playready_configuration_license or fairplay_configuration.")
	}
}

func flattenConfiguration(input media.BasicContentKeyPolicyConfiguration, option map[string]interface{}) map[string]interface{} {
	if input == nil {
		return option
	}

	switch input.(type) {
	case media.ContentKeyPolicyClearKeyConfiguration:
		option["clear_key_configuration_enabled"] = true
	case media.ContentKeyPolicyWidevineConfiguration:
		wideVineConfiguration, _ := input.AsContentKeyPolicyWidevineConfiguration()

		template := ""
		if wideVineConfiguration.WidevineTemplate != nil {
			template = *wideVineConfiguration.WidevineTemplate
		}

		option["widevine_configuration_template"] = template

	case media.ContentKeyPolicyFairPlayConfiguration:
		fairPlayConfiguration, _ := input.AsContentKeyPolicyFairPlayConfiguration()
		option["fairplay_configuration"] = flattenFairplayConfiguration(fairPlayConfiguration)

	case media.ContentKeyPolicyPlayReadyConfiguration:
		playReadyConfiguration, _ := input.AsContentKeyPolicyPlayReadyConfiguration()
		if playReadyConfiguration.Licenses != nil {
			option["playready_configuration_license"] = flattenPlayReadyLicenses(playReadyConfiguration.Licenses)
		}
	}

	return option
}

func expandVerificationKey(input map[string]interface{}) (media.BasicContentKeyPolicyRestrictionTokenKey, error) {
	verificationKeyCount := 0
	verificationKeyType := ""
	if input["primary_symmetric_token_key"] != nil && input["primary_symmetric_token_key"].(string) != "" {
		verificationKeyCount++
		verificationKeyType = string(media.OdataTypeMicrosoftMediaContentKeyPolicySymmetricTokenKey)
	}
	if (input["primary_rsa_token_key_exponent"] != nil && input["primary_rsa_token_key_exponent"].(string) != "") || (input["primary_rsa_token_key_modulus"] != nil && input["primary_rsa_token_key_modulus"].(string) != "") {
		verificationKeyCount++
		verificationKeyType = string(media.OdataTypeMicrosoftMediaContentKeyPolicyRsaTokenKey)
	}

	if input["primary_x509_token_key_raw"] != nil && input["primary_x509_token_key_raw"].(string) != "" {
		verificationKeyCount++
		verificationKeyType = string(media.OdataTypeMicrosoftMediaContentKeyPolicyX509CertificateTokenKey)
	}

	if verificationKeyCount > 1 {
		return nil, fmt.Errorf("more than one type of token key in the same token_restriction is not allowed.")
	}

	switch verificationKeyType {
	case string(media.OdataTypeMicrosoftMediaContentKeyPolicySymmetricTokenKey):
		symmetricTokenKey := &media.ContentKeyPolicySymmetricTokenKey{
			OdataType: media.OdataTypeMicrosoftMediaContentKeyPolicySymmetricTokenKey,
		}

		if input["primary_symmetric_token_key"] != nil && input["primary_symmetric_token_key"].(string) != "" {
			keyValue := []byte(input["primary_symmetric_token_key"].(string))
			symmetricTokenKey.KeyValue = &keyValue
		}
		return symmetricTokenKey, nil
	case string(media.OdataTypeMicrosoftMediaContentKeyPolicyRsaTokenKey):
		rsaTokenKey := &media.ContentKeyPolicyRsaTokenKey{
			OdataType: media.OdataTypeMicrosoftMediaContentKeyPolicyRsaTokenKey,
		}
		if input["primary_rsa_token_key_exponent"] != nil && input["primary_rsa_token_key_exponent"].(string) != "" {
			exponent := []byte(input["primary_rsa_token_key_exponent"].(string))
			rsaTokenKey.Exponent = &exponent
		}
		if input["primary_rsa_token_key_modulus"] != nil && input["primary_rsa_token_key_modulus"].(string) != "" {
			modulus := []byte(input["primary_rsa_token_key_modulus"].(string))
			rsaTokenKey.Modulus = &modulus
		}
		return rsaTokenKey, nil
	case string(media.OdataTypeMicrosoftMediaContentKeyPolicyX509CertificateTokenKey):
		x509CertificateTokenKey := &media.ContentKeyPolicyX509CertificateTokenKey{
			OdataType: media.OdataTypeMicrosoftMediaContentKeyPolicyX509CertificateTokenKey,
		}

		if input["primary_x509_token_key_raw"] != nil && input["primary_x509_token_key_raw"].(string) != "" {
			rawBody := []byte(input["primary_x509_token_key_raw"].(string))
			x509CertificateTokenKey.RawBody = &rawBody
		}
		return x509CertificateTokenKey, nil
	default:
		return nil, nil
	}
}

func flattenVerificationKey(input media.BasicContentKeyPolicyRestrictionTokenKey, key map[string]interface{}) map[string]interface{} {
	if input == nil {
		return key
	}

	switch input.(type) {
	case media.ContentKeyPolicySymmetricTokenKey:
		symmetricTokenKey, _ := input.AsContentKeyPolicySymmetricTokenKey()

		keyValue := ""
		if symmetricTokenKey.KeyValue != nil {
			keyValue = string(*symmetricTokenKey.KeyValue)
		}
		key["primary_symmetric_token_key"] = keyValue
	case media.ContentKeyPolicyRsaTokenKey:
		rsaTokenKey, _ := input.AsContentKeyPolicyRsaTokenKey()
		exponent := ""
		if rsaTokenKey.Exponent != nil {
			exponent = string(*rsaTokenKey.Exponent)
		}
		modulus := ""
		if rsaTokenKey.Modulus != nil {
			modulus = string(*rsaTokenKey.Modulus)
		}
		key["primary_rsa_token_key_exponent"] = exponent
		key["primary_rsa_token_key_modulus"] = modulus
	case media.ContentKeyPolicyX509CertificateTokenKey:
		x509CertificateTokenKey, _ := input.AsContentKeyPolicyX509CertificateTokenKey()
		rawBody := ""
		if x509CertificateTokenKey.RawBody != nil {
			rawBody = string(*x509CertificateTokenKey.RawBody)
		}
		key["primary_x509_token_key_raw"] = rawBody
	}
	return key
}

func expandRequiredClaims(input []interface{}) *[]media.ContentKeyPolicyTokenClaim {
	results := make([]media.ContentKeyPolicyTokenClaim, 0)

	for _, tokenClaimRaw := range input {
		tokenClaim := tokenClaimRaw.(map[string]interface{})

		claimType := ""
		if v := tokenClaim["type"]; v != nil {
			claimType = v.(string)
		}

		claimValue := ""
		if v := tokenClaim["value"]; v != nil {
			claimValue = v.(string)
		}

		contentPolicyTokenClaim := media.ContentKeyPolicyTokenClaim{
			ClaimType:  &claimType,
			ClaimValue: &claimValue,
		}

		results = append(results, contentPolicyTokenClaim)
	}

	return &results
}

func flattenRequiredClaims(input *[]media.ContentKeyPolicyTokenClaim) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, tokenClaim := range *input {
		claim := make(map[string]interface{})
		claim["value"] = tokenClaim.ClaimValue
		claim["type"] = tokenClaim.ClaimType
		results = append(results, claim)
	}

	return results
}

func expandRentalConfiguration(input []interface{}) *media.ContentKeyPolicyFairPlayOfflineRentalConfiguration {
	if len(input) == 0 {
		return nil
	}

	rentalConfiguration := input[0].(map[string]interface{})
	playbackDuration := utils.Int64(int64(rentalConfiguration["playback_duration_seconds"].(int)))
	storageDuration := utils.Int64(int64(rentalConfiguration["storage_duration_seconds"].(int)))
	return &media.ContentKeyPolicyFairPlayOfflineRentalConfiguration{
		PlaybackDurationSeconds: playbackDuration,
		StorageDurationSeconds:  storageDuration,
	}
}

func flattenRentalConfiguration(input *media.ContentKeyPolicyFairPlayOfflineRentalConfiguration) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})
	if input.PlaybackDurationSeconds != nil {
		result["playback_duration_seconds"] = int(*input.PlaybackDurationSeconds)
	}
	if input.StorageDurationSeconds != nil {
		result["storage_duration_seconds"] = int(*input.StorageDurationSeconds)
	}

	return []interface{}{result}
}

func expandFairplayConfiguration(input []interface{}) *media.ContentKeyPolicyFairPlayConfiguration {
	fairplayConfiguration := &media.ContentKeyPolicyFairPlayConfiguration{
		OdataType: media.OdataTypeMicrosoftMediaContentKeyPolicyWidevineConfiguration,
	}

	fairplay := input[0].(map[string]interface{})
	if fairplay["rental_duration_seconds"] != nil {
		fairplayConfiguration.RentalDuration = utils.Int64(int64(fairplay["rental_duration_seconds"].(int)))
	}

	if fairplay["offline_rental_configuration"] != nil {
		fairplayConfiguration.OfflineRentalConfiguration = expandRentalConfiguration(fairplay["offline_rental_configuration"].([]interface{}))
	}

	if fairplay["rental_and_lease_key_type"] != nil {
		fairplayConfiguration.RentalAndLeaseKeyType = media.ContentKeyPolicyFairPlayRentalAndLeaseKeyType(fairplay["rental_and_lease_key_type"].(string))
	}

	if fairplay["ask"] != nil && fairplay["ask"].(string) != "" {
		ask := []byte(fairplay["ask"].(string))
		fairplayConfiguration.Ask = &ask
	}

	if fairplay["pfx"] != nil && fairplay["pfx"].(string) != "" {
		fairplayConfiguration.FairPlayPfx = utils.String(fairplay["pfx"].(string))
	}

	if fairplay["pfx_password"] != nil && fairplay["pfx_password"].(string) != "" {
		fairplayConfiguration.FairPlayPfxPassword = utils.String(fairplay["pfx_password"].(string))
	}

	return fairplayConfiguration
}

func flattenFairplayConfiguration(input *media.ContentKeyPolicyFairPlayConfiguration) []interface{} {
	fairPlay := make(map[string]interface{})
	rentalDuration := 0
	if input.RentalDuration != nil {
		rentalDuration = int(*input.RentalDuration)
	}
	fairPlay["rental_duration_seconds"] = rentalDuration

	offlineRentalConfiguration := make([]interface{}, 0)
	if input.OfflineRentalConfiguration != nil {
		offlineRentalConfiguration = flattenRentalConfiguration(input.OfflineRentalConfiguration)
	}
	fairPlay["offline_rental_configuration"] = offlineRentalConfiguration
	fairPlay["rental_and_lease_key_type"] = string(input.RentalAndLeaseKeyType)

	pfx := ""
	if input.FairPlayPfx != nil {
		pfx = *input.FairPlayPfx
	}
	fairPlay["pfx"] = pfx

	pfxPassword := ""
	if input.FairPlayPfxPassword != nil {
		pfxPassword = *input.FairPlayPfxPassword
	}
	fairPlay["pfx_password"] = pfxPassword

	ask := ""
	if input.Ask != nil {
		ask = string(*input.Ask)
	}
	fairPlay["ask"] = ask

	return []interface{}{
		fairPlay,
	}
}

func expandPlayReadyLicenses(input []interface{}) (*[]media.ContentKeyPolicyPlayReadyLicense, error) {
	results := make([]media.ContentKeyPolicyPlayReadyLicense, 0)

	for _, licenseRaw := range input {
		license := licenseRaw.(map[string]interface{})
		playReadyLicense := media.ContentKeyPolicyPlayReadyLicense{}

		if v := license["allow_test_devices"]; v != nil {
			playReadyLicense.AllowTestDevices = utils.Bool(v.(bool))
		}

		if v := license["begin_date"]; v != nil && v != "" {
			beginDate, err := date.ParseTime(time.RFC3339, v.(string))
			if err != nil {
				return nil, err
			}
			playReadyLicense.BeginDate = &date.Time{
				Time: beginDate,
			}
		}

		locationFromHeader := false
		if v := license["content_key_location_from_header_enabled"]; v != nil && v != "" {
			playReadyLicense.ContentKeyLocation = media.ContentKeyPolicyPlayReadyContentEncryptionKeyFromHeader{
				OdataType: media.OdataTypeMicrosoftMediaContentKeyPolicyPlayReadyContentEncryptionKeyFromHeader,
			}
			locationFromHeader = true
		}

		if v := license["content_key_location_from_key_id"]; v != nil && v != "" {
			if locationFromHeader {
				return nil, fmt.Errorf("playready_configuration_license only support one key location at time, you must to specify content_key_location_from_header_enabled or content_key_location_from_key_id but not both at the same time")
			}

			keyID := uuid.FromStringOrNil(v.(string))
			playReadyLicense.ContentKeyLocation = media.ContentKeyPolicyPlayReadyContentEncryptionKeyFromKeyIdentifier{
				OdataType: media.OdataTypeMicrosoftMediaContentKeyPolicyPlayReadyContentEncryptionKeyFromHeader,
				KeyID:     &keyID,
			}
		}

		if v := license["content_type"]; v != nil && v != "" {
			playReadyLicense.ContentType = media.ContentKeyPolicyPlayReadyContentType(v.(string))
		}

		if v := license["expiration_date"]; v != nil && v != "" {
			expirationDate, err := date.ParseTime(time.RFC3339, v.(string))
			if err != nil {
				return nil, err
			}
			playReadyLicense.ExpirationDate = &date.Time{
				Time: expirationDate,
			}
		}

		if v := license["grace_period"]; v != nil && v != "" {
			playReadyLicense.GracePeriod = utils.String(v.(string))
		}

		if v := license["license_type"]; v != nil && v != "" {
			playReadyLicense.LicenseType = media.ContentKeyPolicyPlayReadyLicenseType(v.(string))
		}

		if v := license["play_right"]; v != nil {
			playReadyLicense.PlayRight = expandPlayRight(v.([]interface{}))
		}

		if v := license["relative_begin_date"]; v != nil && v != "" {
			playReadyLicense.RelativeBeginDate = utils.String(v.(string))
		}

		if v := license["relative_expiration_date"]; v != nil && v != "" {
			playReadyLicense.RelativeExpirationDate = utils.String(v.(string))
		}

		results = append(results, playReadyLicense)
	}

	return &results, nil
}

func flattenPlayReadyLicenses(input *[]media.ContentKeyPolicyPlayReadyLicense) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, v := range *input {
		license := make(map[string]interface{})

		if v.AllowTestDevices != nil {
			license["allow_test_devices"] = *v.AllowTestDevices
		}

		if v.BeginDate != nil {
			license["begin_date"] = v.BeginDate.Format(time.RFC3339)
		}

		if v.ContentKeyLocation != nil {
			switch v.ContentKeyLocation.(type) {
			case media.ContentKeyPolicyPlayReadyContentEncryptionKeyFromHeader:
				license["content_key_location_from_header_enabled"] = true
			case media.ContentKeyPolicyPlayReadyContentEncryptionKeyFromKeyIdentifier:
				keyLocation, _ := v.ContentKeyLocation.AsContentKeyPolicyPlayReadyContentEncryptionKeyFromKeyIdentifier()
				license["content_key_location_from_key_id"] = keyLocation.KeyID.String()
			}
		}

		license["content_type"] = string(v.ContentType)

		if v.ExpirationDate != nil {
			license["expiration_date"] = v.ExpirationDate.Format(time.RFC3339)
		}

		if v.GracePeriod != nil {
			license["grace_period"] = *v.GracePeriod
		}

		license["license_type"] = string(v.LicenseType)

		if v.PlayRight != nil {
			license["play_right"] = flattenPlayRight(v.PlayRight)
		}

		if v.RelativeBeginDate != nil {
			license["relative_begin_date"] = *v.RelativeBeginDate
		}

		if v.RelativeExpirationDate != nil {
			license["relative_expiration_date"] = *v.RelativeExpirationDate
		}

		results = append(results, license)
	}

	return results
}

func expandPlayRight(input []interface{}) *media.ContentKeyPolicyPlayReadyPlayRight {
	if len(input) == 0 {
		return nil
	}

	playRight := &media.ContentKeyPolicyPlayReadyPlayRight{}
	playRightConfiguration := input[0].(map[string]interface{})

	if v := playRightConfiguration["agc_and_color_stripe_restriction"]; v != nil {
		playRight.AgcAndColorStripeRestriction = utils.Int32(int32(v.(int)))
	}

	if v := playRightConfiguration["allow_passing_video_content_to_unknown_output"]; v != nil {
		playRight.AllowPassingVideoContentToUnknownOutput = media.ContentKeyPolicyPlayReadyUnknownOutputPassingOption(v.(string))
	}

	if v := playRightConfiguration["analog_video_opl"]; v != nil && v != 0 {
		playRight.AnalogVideoOpl = utils.Int32(int32(v.(int)))
	}

	if v := playRightConfiguration["compressed_digital_audio_opl"]; v != nil && v != 0 {
		playRight.CompressedDigitalAudioOpl = utils.Int32(int32(v.(int)))
	}

	if v := playRightConfiguration["digital_video_only_content_restriction"]; v != nil {
		playRight.DigitalVideoOnlyContentRestriction = utils.Bool(v.(bool))
	}

	if v := playRightConfiguration["first_play_expiration"]; v != nil && v != "" {
		playRight.FirstPlayExpiration = utils.String(v.(string))
	}

	if v := playRightConfiguration["image_constraint_for_analog_component_video_restriction"]; v != nil {
		playRight.ImageConstraintForAnalogComponentVideoRestriction = utils.Bool(v.(bool))
	}

	if v := playRightConfiguration["image_constraint_for_analog_computer_monitor_restriction"]; v != nil {
		playRight.ImageConstraintForAnalogComputerMonitorRestriction = utils.Bool(v.(bool))
	}

	if v := playRightConfiguration["scms_restriction"]; v != nil {
		playRight.ScmsRestriction = utils.Int32(int32(v.(int)))
	}
	if v := playRightConfiguration["uncompressed_digital_audio_opl"]; v != nil && v != 0 {
		playRight.UncompressedDigitalAudioOpl = utils.Int32(int32(v.(int)))
	}

	if v := playRightConfiguration["uncompressed_digital_video_opl"]; v != nil && v != 0 {
		playRight.UncompressedDigitalVideoOpl = utils.Int32(int32(v.(int)))
	}

	return playRight
}

func flattenPlayRight(input *media.ContentKeyPolicyPlayReadyPlayRight) []interface{} {
	playRight := make(map[string]interface{})

	agcStripeRestriction := 0
	if input.AgcAndColorStripeRestriction != nil {
		agcStripeRestriction = int(*input.AgcAndColorStripeRestriction)
	}
	playRight["agc_and_color_stripe_restriction"] = agcStripeRestriction

	playRight["allow_passing_video_content_to_unknown_output"] = string(input.AllowPassingVideoContentToUnknownOutput)

	analogVideoOpl := 0
	if input.AnalogVideoOpl != nil {
		analogVideoOpl = int(*input.AnalogVideoOpl)
	}
	playRight["analog_video_opl"] = analogVideoOpl

	compressedDigitalAudioOpl := 0
	if input.AnalogVideoOpl != nil {
		compressedDigitalAudioOpl = int(*input.CompressedDigitalAudioOpl)
	}
	playRight["compressed_digital_audio_opl"] = compressedDigitalAudioOpl

	digitalVideoOnlyContentRestriction := false
	if input.DigitalVideoOnlyContentRestriction != nil {
		digitalVideoOnlyContentRestriction = *input.DigitalVideoOnlyContentRestriction
	}
	playRight["digital_video_only_content_restriction"] = digitalVideoOnlyContentRestriction

	firstPlayExpiration := ""
	if input.FirstPlayExpiration != nil {
		firstPlayExpiration = *input.FirstPlayExpiration
	}
	playRight["first_play_expiration"] = firstPlayExpiration

	imageConstraintForAnalogComponentVideoRestriction := false
	if input.ImageConstraintForAnalogComponentVideoRestriction != nil {
		imageConstraintForAnalogComponentVideoRestriction = *input.ImageConstraintForAnalogComponentVideoRestriction
	}
	playRight["image_constraint_for_analog_component_video_restriction"] = imageConstraintForAnalogComponentVideoRestriction

	imageConstraintForAnalogComputerMonitorRestriction := false
	if input.ImageConstraintForAnalogComputerMonitorRestriction != nil {
		imageConstraintForAnalogComputerMonitorRestriction = *input.ImageConstraintForAnalogComputerMonitorRestriction
	}
	playRight["image_constraint_for_analog_computer_monitor_restriction"] = imageConstraintForAnalogComputerMonitorRestriction

	scmsRestriction := 0
	if input.ScmsRestriction != nil {
		scmsRestriction = int(*input.ScmsRestriction)
	}
	playRight["scms_restriction"] = scmsRestriction

	if input.UncompressedDigitalAudioOpl != nil {
		playRight["uncompressed_digital_audio_opl"] = int(*input.UncompressedDigitalAudioOpl)
	}

	if input.UncompressedDigitalVideoOpl != nil {
		playRight["uncompressed_digital_video_opl"] = int(*input.UncompressedDigitalVideoOpl)
	}

	return []interface{}{
		playRight,
	}
}
