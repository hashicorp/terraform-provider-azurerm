// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package media

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/contentkeypolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/media/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/media/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMediaContentKeyPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMediaContentKeyPolicyCreateUpdate,
		Read:   resourceMediaContentKeyPolicyRead,
		Update: resourceMediaContentKeyPolicyCreateUpdate,
		Delete: resourceMediaContentKeyPolicyDelete,

		DeprecationMessage: azureMediaRetirementMessage,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := contentkeypolicies.ParseContentKeyPolicyID(id)
			return err
		}),

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ContentKeyPolicyV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9(_)]{1,128}$"),
					"Content Key Policy name must be 1 - 128 characters long, can contain letters, numbers, underscores, and hyphens (but the first and last character must be a letter or number).",
				),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"media_services_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"policy_option": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"clear_key_configuration_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"widevine_configuration_template": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						// lintignore:XS003
						"playready_configuration_license": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"allow_test_devices": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
									},

									"begin_date": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsRFC3339Time,
									},

									"content_key_location_from_header_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
									},

									"content_key_location_from_key_id": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsUUID,
									},

									"content_type": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(contentkeypolicies.ContentKeyPolicyPlayReadyContentTypeUltraVioletDownload),
											string(contentkeypolicies.ContentKeyPolicyPlayReadyContentTypeUltraVioletStreaming),
											string(contentkeypolicies.ContentKeyPolicyPlayReadyContentTypeUnspecified),
										}, false),
									},

									"expiration_date": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsRFC3339Time,
									},

									"grace_period": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Sensitive:    true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"license_type": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(contentkeypolicies.ContentKeyPolicyPlayReadyLicenseTypeNonPersistent),
											string(contentkeypolicies.ContentKeyPolicyPlayReadyLicenseTypePersistent),
										}, false),
									},

									// lintignore:XS003
									"play_right": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"agc_and_color_stripe_restriction": {
													Type:         pluginsdk.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntBetween(0, 3),
												},

												"allow_passing_video_content_to_unknown_output": {
													Type:     pluginsdk.TypeString,
													Optional: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(contentkeypolicies.ContentKeyPolicyPlayReadyUnknownOutputPassingOptionAllowed),
														string(contentkeypolicies.ContentKeyPolicyPlayReadyUnknownOutputPassingOptionAllowedWithVideoConstriction),
														string(contentkeypolicies.ContentKeyPolicyPlayReadyUnknownOutputPassingOptionNotAllowed),
													}, false),
												},

												"analog_video_opl": {
													Type:         pluginsdk.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntInSlice([]int{100, 150, 200}),
												},

												"compressed_digital_audio_opl": {
													Type:         pluginsdk.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntInSlice([]int{100, 150, 200, 250, 300}),
												},

												"compressed_digital_video_opl": {
													Type:         pluginsdk.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntInSlice([]int{400, 500}),
												},

												"digital_video_only_content_restriction": {
													Type:     pluginsdk.TypeBool,
													Optional: true,
												},

												"explicit_analog_television_output_restriction": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*pluginsdk.Schema{
															"best_effort_enforced": {
																Type:     pluginsdk.TypeBool,
																Optional: true,
																Default:  false,
															},
															"control_bits": {
																Type:         pluginsdk.TypeInt,
																Required:     true,
																ValidateFunc: validation.IntBetween(0, 3),
															},
														},
													},
												},

												"first_play_expiration": {
													Type:         pluginsdk.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},

												"image_constraint_for_analog_component_video_restriction": {
													Type:     pluginsdk.TypeBool,
													Optional: true,
												},

												"image_constraint_for_analog_computer_monitor_restriction": {
													Type:     pluginsdk.TypeBool,
													Optional: true,
												},

												"scms_restriction": {
													Type:         pluginsdk.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntBetween(0, 3),
												},

												"uncompressed_digital_audio_opl": {
													Type:         pluginsdk.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntInSlice([]int{100, 150, 200, 250, 300}),
												},

												"uncompressed_digital_video_opl": {
													Type:         pluginsdk.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntInSlice([]int{100, 250, 270, 300}),
												},
											},
										},
									},

									"relative_begin_date": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsRFC3339Time,
									},

									"relative_expiration_date": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsRFC3339Time,
									},

									"security_level": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(contentkeypolicies.SecurityLevelSLOneFiveZero),
											string(contentkeypolicies.SecurityLevelSLTwoThousand),
											string(contentkeypolicies.SecurityLevelSLThreeThousand),
										}, false),
									},
								},
							},
						},

						"playready_response_custom_data": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						// lintignore:XS003
						"fairplay_configuration": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"ask": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Sensitive:    true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"pfx": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Sensitive:    true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"pfx_password": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Sensitive:    true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									// lintignore:XS003
									"offline_rental_configuration": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"playback_duration_seconds": {
													Type:         pluginsdk.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntAtLeast(1),
												},
												"storage_duration_seconds": {
													Type:         pluginsdk.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntAtLeast(1),
												},
											},
										},
									},
									"rental_and_lease_key_type": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(contentkeypolicies.ContentKeyPolicyFairPlayRentalAndLeaseKeyTypeDualExpiry),
											string(contentkeypolicies.ContentKeyPolicyFairPlayRentalAndLeaseKeyTypePersistentLimited),
											string(contentkeypolicies.ContentKeyPolicyFairPlayRentalAndLeaseKeyTypePersistentUnlimited),
											string(contentkeypolicies.ContentKeyPolicyFairPlayRentalAndLeaseKeyTypeUndefined),
										}, false),
									},
									"rental_duration_seconds": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(1),
									},
								},
							},
						},
						// lintignore:XS003
						"token_restriction": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									// lintignore:XS003
									"alternate_key": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"symmetric_token_key": {
													Type:         pluginsdk.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringIsBase64,
													Sensitive:    true,
												},
												"rsa_token_key_exponent": {
													Type:         pluginsdk.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
													Sensitive:    true,
												},
												"rsa_token_key_modulus": {
													Type:         pluginsdk.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
													Sensitive:    true,
												},
												"x509_token_key_raw": {
													Type:         pluginsdk.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
													Sensitive:    true,
												},
											},
										},
									},
									"audience": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"issuer": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"token_type": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(contentkeypolicies.ContentKeyPolicyRestrictionTokenTypeJwt),
											string(contentkeypolicies.ContentKeyPolicyRestrictionTokenTypeSwt),
										}, false),
									},
									"primary_symmetric_token_key": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsBase64,
										Sensitive:    true,
									},
									"primary_rsa_token_key_exponent": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
										Sensitive:    true,
									},
									"primary_rsa_token_key_modulus": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
										Sensitive:    true,
									},
									"primary_x509_token_key_raw": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
										Sensitive:    true,
									},
									"open_id_connect_discovery_document": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									// lintignore:XS003
									"required_claim": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"type": {
													Type:         pluginsdk.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"value": {
													Type:         pluginsdk.TypeString,
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
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceMediaContentKeyPolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220801Client.ContentKeyPolicies
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := contentkeypolicies.NewContentKeyPolicyID(subscriptionID, d.Get("resource_group_name").(string), d.Get("media_services_account_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_media_content_key_policy", id.ID())
		}
	}

	payload := contentkeypolicies.ContentKeyPolicy{
		Properties: &contentkeypolicies.ContentKeyPolicyProperties{},
	}

	if description, ok := d.GetOk("description"); ok {
		payload.Properties.Description = utils.String(description.(string))
	}

	if v, ok := d.GetOk("policy_option"); ok {
		options, err := expandPolicyOptions(v.(*pluginsdk.Set).List())
		if err != nil {
			return err
		}
		payload.Properties.Options = *options
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMediaContentKeyPolicyRead(d, meta)
}

func resourceMediaContentKeyPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220801Client.ContentKeyPolicies
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := contentkeypolicies.ParseContentKeyPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetPolicyPropertiesWithSecrets(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ContentKeyPolicyName)
	d.Set("media_services_account_name", id.MediaServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("description", model.Description)

		options, err := flattenPolicyOptions(model.Options)
		if err != nil {
			return fmt.Errorf("flattening `policy_option`: %+v", err)
		}
		if err := d.Set("policy_option", options); err != nil {
			return fmt.Errorf("setting `policy_option`: %+v", err)
		}
	}

	return nil
}

func resourceMediaContentKeyPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220801Client.ContentKeyPolicies
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := contentkeypolicies.ParseContentKeyPolicyID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandPolicyOptions(input []interface{}) (*[]contentkeypolicies.ContentKeyPolicyOption, error) {
	results := make([]contentkeypolicies.ContentKeyPolicyOption, 0)

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

		contentKeyPolicyOption := contentkeypolicies.ContentKeyPolicyOption{
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

func flattenPolicyOptions(input []contentkeypolicies.ContentKeyPolicyOption) ([]interface{}, error) {
	results := make([]interface{}, 0)
	for _, option := range input {
		name := ""
		if option.Name != nil {
			name = *option.Name
		}

		clearKeyConfigurationEnabled := false
		playReadyLicense := make([]interface{}, 0)
		playReadyResponseCustomData := ""
		widevineTemplate := ""
		fairplayConfiguration := make([]interface{}, 0)

		if _, ok := option.Configuration.(contentkeypolicies.ContentKeyPolicyClearKeyConfiguration); ok {
			clearKeyConfigurationEnabled = true
		}

		if v, ok := option.Configuration.(contentkeypolicies.ContentKeyPolicyFairPlayConfiguration); ok {
			flattened, err := flattenFairplayConfiguration(v)
			if err != nil {
				return nil, fmt.Errorf("flattening fairplay configuration: %+v", err)
			}
			fairplayConfiguration = *flattened
		}

		if v, ok := option.Configuration.(contentkeypolicies.ContentKeyPolicyPlayReadyConfiguration); ok {
			playReadyLicense = flattenPlayReadyLicenses(v.Licenses)
			if v.ResponseCustomData != nil {
				playReadyResponseCustomData = *v.ResponseCustomData
			}
		}

		if v, ok := option.Configuration.(contentkeypolicies.ContentKeyPolicyWidevineConfiguration); ok {
			widevineTemplate = v.WidevineTemplate
		}

		openRestrictionEnabled := false
		tokenRestriction := make([]interface{}, 0)
		if _, ok := option.Restriction.(contentkeypolicies.ContentKeyPolicyOpenRestriction); ok {
			openRestrictionEnabled = true
		}
		if v, ok := option.Restriction.(contentkeypolicies.ContentKeyPolicyTokenRestriction); ok {
			tokenRestriction = flattenTokenRestriction(v)
		}

		results = append(results, map[string]interface{}{
			"name":                            name,
			"clear_key_configuration_enabled": clearKeyConfigurationEnabled,
			"playready_configuration_license": playReadyLicense,
			"widevine_configuration_template": widevineTemplate,
			"fairplay_configuration":          fairplayConfiguration,
			"playready_response_custom_data":  playReadyResponseCustomData,
			"open_restriction_enabled":        openRestrictionEnabled,
			"token_restriction":               tokenRestriction,
		})
	}

	return results, nil
}

func expandRestriction(option map[string]interface{}) (contentkeypolicies.ContentKeyPolicyRestriction, error) {
	openRestrictionEnabled := option["open_restriction_enabled"].(bool)
	tokenRestrictions := option["token_restriction"].([]interface{})

	restrictionCount := 0
	if openRestrictionEnabled {
		restrictionCount++
	}
	if len(tokenRestrictions) > 0 {
		restrictionCount++
	}
	if restrictionCount == 0 {
		return nil, fmt.Errorf("policy_option must contain at least one type of restriction: open_restriction_enabled or token_restriction.")
	}
	if restrictionCount > 1 {
		return nil, fmt.Errorf("more than one type of restriction in the same policy_option is not allowed.")
	}

	if openRestrictionEnabled {
		openRestriction := &contentkeypolicies.ContentKeyPolicyOpenRestriction{}
		return openRestriction, nil
	}

	if len(tokenRestrictions) > 0 {
		tokenRestriction := tokenRestrictions[0].(map[string]interface{})
		contentKeyPolicyTokenRestriction := &contentkeypolicies.ContentKeyPolicyTokenRestriction{}
		if tokenRestriction["audience"] != nil && tokenRestriction["audience"].(string) != "" {
			contentKeyPolicyTokenRestriction.Audience = tokenRestriction["audience"].(string)
		}
		if tokenRestriction["issuer"] != nil && tokenRestriction["issuer"].(string) != "" {
			contentKeyPolicyTokenRestriction.Issuer = tokenRestriction["issuer"].(string)
		}
		if tokenRestriction["token_type"] != nil && tokenRestriction["token_type"].(string) != "" {
			contentKeyPolicyTokenRestriction.RestrictionTokenType = contentkeypolicies.ContentKeyPolicyRestrictionTokenType(tokenRestriction["token_type"].(string))
		}
		if tokenRestriction["open_id_connect_discovery_document"] != nil && tokenRestriction["open_id_connect_discovery_document"].(string) != "" {
			contentKeyPolicyTokenRestriction.OpenIdConnectDiscoveryDocument = utils.String(tokenRestriction["open_id_connect_discovery_document"].(string))
		}
		if v := tokenRestriction["required_claim"]; v != nil {
			contentKeyPolicyTokenRestriction.RequiredClaims = expandRequiredClaims(v.([]interface{}))
		}
		primaryVerificationKey, err := expandVerificationKey(tokenRestriction)
		if err != nil {
			return nil, err
		}
		contentKeyPolicyTokenRestriction.PrimaryVerificationKey = primaryVerificationKey

		alternateVerificationKeys, err := expandAlternateVerificationKeys(tokenRestriction["alternate_key"].([]interface{}))
		if err != nil {
			return nil, err
		}
		contentKeyPolicyTokenRestriction.AlternateVerificationKeys = alternateVerificationKeys
		return contentKeyPolicyTokenRestriction, nil
	}

	return nil, fmt.Errorf("policy_option must contain at least one type of restriction: open_restriction_enabled or token_restriction.")
}

func flattenTokenRestriction(input contentkeypolicies.ContentKeyPolicyTokenRestriction) []interface{} {
	openIDConnectDiscoveryDocument := ""
	if input.OpenIdConnectDiscoveryDocument != nil {
		openIDConnectDiscoveryDocument = *input.OpenIdConnectDiscoveryDocument
	}

	requiredClaims := make([]interface{}, 0)
	if input.RequiredClaims != nil {
		requiredClaims = flattenRequiredClaims(input.RequiredClaims)
	}

	symmetricToken := ""
	rsaTokenKeyExponent := ""
	rsaTokenKeyModulus := ""
	x509TokenBodyRaw := ""
	if v := input.PrimaryVerificationKey; v != nil {
		symmetricTokenKey, ok := v.(contentkeypolicies.ContentKeyPolicySymmetricTokenKey)
		if ok {
			symmetricToken = symmetricTokenKey.KeyValue
		}

		rsaTokenKey, ok := v.(contentkeypolicies.ContentKeyPolicyRsaTokenKey)
		if ok {
			rsaTokenKeyExponent = rsaTokenKey.Exponent
			rsaTokenKeyModulus = rsaTokenKey.Modulus
		}

		x509CertificateTokenKey, ok := v.(contentkeypolicies.ContentKeyPolicyX509CertificateTokenKey)
		if ok {
			x509TokenBodyRaw = x509CertificateTokenKey.RawBody
		}
	}

	return []interface{}{
		map[string]interface{}{
			"alternate_key":                      flattenAlternateVerificationKeys(input.AlternateVerificationKeys),
			"audience":                           input.Audience,
			"issuer":                             input.Issuer,
			"token_type":                         string(input.RestrictionTokenType),
			"open_id_connect_discovery_document": openIDConnectDiscoveryDocument,
			"required_claim":                     requiredClaims,
			"primary_symmetric_token_key":        symmetricToken,
			"primary_x509_token_key_raw":         x509TokenBodyRaw,
			"primary_rsa_token_key_exponent":     rsaTokenKeyExponent,
			"primary_rsa_token_key_modulus":      rsaTokenKeyModulus,
		},
	}
}

func expandConfiguration(input map[string]interface{}) (contentkeypolicies.ContentKeyPolicyConfiguration, error) {
	clearKeyConfigurationEnabled := input["clear_key_configuration_enabled"].(bool)
	fairPlayConfigurations := input["fairplay_configuration"].([]interface{})
	playReadyConfigurationLicences := input["playready_configuration_license"].([]interface{})
	widevineConfigurationTemplate := input["widevine_configuration_template"].(string)
	playReadyResponseCustomData := input["playready_response_custom_data"].(string)

	configurationCount := 0
	if clearKeyConfigurationEnabled {
		configurationCount++
	}
	if len(fairPlayConfigurations) > 0 {
		configurationCount++
	}
	if len(playReadyConfigurationLicences) > 0 {
		configurationCount++
	}
	if widevineConfigurationTemplate != "" {
		configurationCount++
	}
	if configurationCount == 0 {
		return nil, fmt.Errorf("policy_option must contain at least one type of configuration: clear_key_configuration_enabled , widevine_configuration_template, playready_configuration_license or fairplay_configuration.")
	}
	if configurationCount > 1 {
		return nil, fmt.Errorf("more than one type of configuration in the same policy_option is not allowed.")
	}

	if clearKeyConfigurationEnabled {
		clearKeyConfiguration := &contentkeypolicies.ContentKeyPolicyClearKeyConfiguration{}
		return clearKeyConfiguration, nil
	}
	if len(fairPlayConfigurations) > 0 {
		fairplayConfiguration, err := expandFairplayConfiguration(input["fairplay_configuration"].([]interface{}))
		if err != nil {
			return nil, err
		}
		return fairplayConfiguration, nil
	}
	if len(playReadyConfigurationLicences) > 0 {
		licenses, err := expandPlayReadyLicenses(input["playready_configuration_license"].([]interface{}))
		if err != nil {
			return nil, err
		}
		playReadyConfiguration := &contentkeypolicies.ContentKeyPolicyPlayReadyConfiguration{
			Licenses: *licenses,
		}
		if playReadyResponseCustomData != "" {
			playReadyConfiguration.ResponseCustomData = utils.String(playReadyResponseCustomData)
		}
		return playReadyConfiguration, nil
	}
	if widevineConfigurationTemplate != "" {
		wideVineConfiguration := &contentkeypolicies.ContentKeyPolicyWidevineConfiguration{
			WidevineTemplate: input["widevine_configuration_template"].(string),
		}
		return wideVineConfiguration, nil
	}

	return nil, fmt.Errorf("policy_option must contain at least one type of configuration: clear_key_configuration_enabled , widevine_configuration_template, playready_configuration_license or fairplay_configuration.")
}

func expandVerificationKey(input map[string]interface{}) (contentkeypolicies.ContentKeyPolicyRestrictionTokenKey, error) {
	primaryRsaTokenKeyExponent := input["primary_rsa_token_key_exponent"].(string)
	primaryRsaTokenKeyModulus := input["primary_rsa_token_key_modulus"].(string)
	primarySymmetricTokenKey := input["primary_symmetric_token_key"].(string)
	primaryX509TokenKeyRaw := input["primary_x509_token_key_raw"].(string)

	verificationKeyCount := 0
	if primaryRsaTokenKeyExponent != "" || primaryRsaTokenKeyModulus != "" {
		verificationKeyCount++
	}
	if primarySymmetricTokenKey != "" {
		verificationKeyCount++
	}
	if primaryX509TokenKeyRaw != "" {
		verificationKeyCount++
	}
	if verificationKeyCount > 1 {
		return nil, fmt.Errorf("more than one type of token key in the same token_restriction is not allowed.")
	}

	if primaryRsaTokenKeyExponent != "" || primaryRsaTokenKeyModulus != "" {
		rsaTokenKey := &contentkeypolicies.ContentKeyPolicyRsaTokenKey{
			Exponent: primaryRsaTokenKeyExponent,
			Modulus:  primaryRsaTokenKeyModulus,
		}
		return rsaTokenKey, nil
	}
	if primarySymmetricTokenKey != "" {
		symmetricTokenKey := &contentkeypolicies.ContentKeyPolicySymmetricTokenKey{
			KeyValue: primarySymmetricTokenKey,
		}
		return symmetricTokenKey, nil
	}
	if primaryX509TokenKeyRaw != "" {
		x509CertificateTokenKey := &contentkeypolicies.ContentKeyPolicyX509CertificateTokenKey{
			RawBody: primaryX509TokenKeyRaw,
		}
		return x509CertificateTokenKey, nil
	}

	return nil, nil
}

func expandAlternateVerificationKeys(input []interface{}) (*[]contentkeypolicies.ContentKeyPolicyRestrictionTokenKey, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	result := make([]contentkeypolicies.ContentKeyPolicyRestrictionTokenKey, 0)
	for _, v := range input {
		tokenKeyRaw := v.(map[string]interface{})
		symmetricTokenKey := tokenKeyRaw["symmetric_token_key"].(string)
		rsaTokenKeyExponent := tokenKeyRaw["rsa_token_key_exponent"].(string)
		rsaTokenKeyModulus := tokenKeyRaw["rsa_token_key_modulus"].(string)
		x509TokenKeyRaw := tokenKeyRaw["x509_token_key_raw"].(string)

		verificationKeyCount := 0
		if rsaTokenKeyExponent != "" || rsaTokenKeyModulus != "" {
			verificationKeyCount++
		}
		if symmetricTokenKey != "" {
			verificationKeyCount++
		}
		if x509TokenKeyRaw != "" {
			verificationKeyCount++
		}
		if verificationKeyCount != 1 {
			return nil, fmt.Errorf("exactlly one type of token key must be set in the alternate verificaton keys")
		}

		if rsaTokenKeyExponent != "" || rsaTokenKeyModulus != "" {
			result = append(result, &contentkeypolicies.ContentKeyPolicyRsaTokenKey{
				Exponent: rsaTokenKeyExponent,
				Modulus:  rsaTokenKeyModulus,
			})
		}
		if symmetricTokenKey != "" {
			result = append(result, &contentkeypolicies.ContentKeyPolicySymmetricTokenKey{
				KeyValue: symmetricTokenKey,
			})
		}
		if x509TokenKeyRaw != "" {
			result = append(result, &contentkeypolicies.ContentKeyPolicyX509CertificateTokenKey{
				RawBody: symmetricTokenKey,
			})
		}
	}

	return &result, nil
}

func flattenAlternateVerificationKeys(input *[]contentkeypolicies.ContentKeyPolicyRestrictionTokenKey) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make([]interface{}, 0)
	for _, v := range *input {
		symmetricToken := ""
		rsaTokenKeyExponent := ""
		rsaTokenKeyModulus := ""
		x509TokenBodyRaw := ""
		symmetricTokenKey, ok := v.(contentkeypolicies.ContentKeyPolicySymmetricTokenKey)
		if ok {
			symmetricToken = symmetricTokenKey.KeyValue
		}

		rsaTokenKey, ok := v.(contentkeypolicies.ContentKeyPolicyRsaTokenKey)
		if ok {
			rsaTokenKeyExponent = rsaTokenKey.Exponent
			rsaTokenKeyModulus = rsaTokenKey.Modulus
		}

		x509CertificateTokenKey, ok := v.(contentkeypolicies.ContentKeyPolicyX509CertificateTokenKey)
		if ok {
			x509TokenBodyRaw = x509CertificateTokenKey.RawBody
		}
		result = append(result, map[string]interface{}{
			"symmetric_token_key":    symmetricToken,
			"x509_token_key_raw":     x509TokenBodyRaw,
			"rsa_token_key_exponent": rsaTokenKeyExponent,
			"rsa_token_key_modulus":  rsaTokenKeyModulus,
		})
	}
	return result
}

func expandRequiredClaims(input []interface{}) *[]contentkeypolicies.ContentKeyPolicyTokenClaim {
	results := make([]contentkeypolicies.ContentKeyPolicyTokenClaim, 0)

	for _, tokenClaimRaw := range input {
		if tokenClaimRaw == nil {
			continue
		}
		tokenClaim := tokenClaimRaw.(map[string]interface{})

		claimType := ""
		if v := tokenClaim["type"]; v != nil {
			claimType = v.(string)
		}

		claimValue := ""
		if v := tokenClaim["value"]; v != nil {
			claimValue = v.(string)
		}

		contentPolicyTokenClaim := contentkeypolicies.ContentKeyPolicyTokenClaim{
			ClaimType:  &claimType,
			ClaimValue: &claimValue,
		}

		results = append(results, contentPolicyTokenClaim)
	}

	return &results
}

func flattenRequiredClaims(input *[]contentkeypolicies.ContentKeyPolicyTokenClaim) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, tokenClaim := range *input {
		claimValue := ""
		if tokenClaim.ClaimValue != nil {
			claimValue = *tokenClaim.ClaimValue
		}

		claimType := ""
		if tokenClaim.ClaimType != nil {
			claimType = *tokenClaim.ClaimType
		}

		results = append(results, map[string]interface{}{
			"value": claimValue,
			"type":  claimType,
		})
	}

	return results
}

func expandRentalConfiguration(input []interface{}) *contentkeypolicies.ContentKeyPolicyFairPlayOfflineRentalConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	rentalConfiguration := input[0].(map[string]interface{})
	playbackDuration := int64(rentalConfiguration["playback_duration_seconds"].(int))
	storageDuration := int64(rentalConfiguration["storage_duration_seconds"].(int))
	return &contentkeypolicies.ContentKeyPolicyFairPlayOfflineRentalConfiguration{
		PlaybackDurationSeconds: playbackDuration,
		StorageDurationSeconds:  storageDuration,
	}
}

func flattenRentalConfiguration(input *contentkeypolicies.ContentKeyPolicyFairPlayOfflineRentalConfiguration) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{map[string]interface{}{
		"playback_duration_seconds": input.PlaybackDurationSeconds,
		"storage_duration_seconds":  input.StorageDurationSeconds,
	}}
}

func expandFairplayConfiguration(input []interface{}) (*contentkeypolicies.ContentKeyPolicyFairPlayConfiguration, error) {
	fairplayConfiguration := &contentkeypolicies.ContentKeyPolicyFairPlayConfiguration{}

	fairplay := input[0].(map[string]interface{})
	if fairplay["rental_duration_seconds"] != nil {
		fairplayConfiguration.RentalDuration = int64(fairplay["rental_duration_seconds"].(int))
	}

	if fairplay["offline_rental_configuration"] != nil {
		fairplayConfiguration.OfflineRentalConfiguration = expandRentalConfiguration(fairplay["offline_rental_configuration"].([]interface{}))
	}

	if fairplay["rental_and_lease_key_type"] != nil {
		fairplayConfiguration.RentalAndLeaseKeyType = contentkeypolicies.ContentKeyPolicyFairPlayRentalAndLeaseKeyType(fairplay["rental_and_lease_key_type"].(string))
	}

	if fairplay["ask"] != nil && fairplay["ask"].(string) != "" {
		askBytes, err := hex.DecodeString(fairplay["ask"].(string))
		if err != nil {
			return nil, err
		}
		fairplayConfiguration.Ask = base64.StdEncoding.EncodeToString(askBytes)
	}

	if fairplay["pfx"] != nil && fairplay["pfx"].(string) != "" {
		fairplayConfiguration.FairPlayPfx = fairplay["pfx"].(string)
	}

	if fairplay["pfx_password"] != nil && fairplay["pfx_password"].(string) != "" {
		fairplayConfiguration.FairPlayPfxPassword = fairplay["pfx_password"].(string)
	}

	return fairplayConfiguration, nil
}

func flattenFairplayConfiguration(input contentkeypolicies.ContentKeyPolicyFairPlayConfiguration) (*[]interface{}, error) {
	offlineRentalConfiguration := make([]interface{}, 0)
	if input.OfflineRentalConfiguration != nil {
		offlineRentalConfiguration = flattenRentalConfiguration(input.OfflineRentalConfiguration)
	}

	ask := ""
	if input.Ask != "" {
		decodedAsk, err := base64.StdEncoding.DecodeString(input.Ask)
		if err != nil {
			return nil, fmt.Errorf("base64-decoding %q: %+v", input.Ask, err)
		}
		ask = hex.EncodeToString(decodedAsk)
	}

	return &[]interface{}{
		map[string]interface{}{
			"rental_duration_seconds":      input.RentalDuration,
			"offline_rental_configuration": offlineRentalConfiguration,
			"rental_and_lease_key_type":    string(input.RentalAndLeaseKeyType),
			"pfx":                          input.FairPlayPfx,
			"pfx_password":                 input.FairPlayPfxPassword,
			"ask":                          ask,
		},
	}, nil
}

func expandPlayReadyLicenses(input []interface{}) (*[]contentkeypolicies.ContentKeyPolicyPlayReadyLicense, error) {
	results := make([]contentkeypolicies.ContentKeyPolicyPlayReadyLicense, 0)

	for _, licenseRaw := range input {
		if licenseRaw == nil {
			continue
		}
		license := licenseRaw.(map[string]interface{})
		playReadyLicense := contentkeypolicies.ContentKeyPolicyPlayReadyLicense{}

		if v := license["allow_test_devices"]; v != nil {
			playReadyLicense.AllowTestDevices = v.(bool)
		}

		if v := license["begin_date"]; v != nil && v != "" {
			beginDate, err := time.Parse(time.RFC3339, v.(string))
			if err != nil {
				return nil, err
			}
			playReadyLicense.SetBeginDateAsTime(beginDate)
		}

		locationFromHeader := false
		if v := license["content_key_location_from_header_enabled"]; v != nil && v != "" {
			playReadyLicense.ContentKeyLocation = contentkeypolicies.ContentKeyPolicyPlayReadyContentEncryptionKeyFromHeader{}
			locationFromHeader = true
		}

		if v := license["content_key_location_from_key_id"]; v != nil && v != "" {
			if locationFromHeader {
				return nil, fmt.Errorf("playready_configuration_license only support one key location at time, you must to specify content_key_location_from_header_enabled or content_key_location_from_key_id but not both at the same time")
			}

			playReadyLicense.ContentKeyLocation = contentkeypolicies.ContentKeyPolicyPlayReadyContentEncryptionKeyFromKeyIdentifier{
				KeyId: v.(string),
			}
		}

		if v := license["content_type"]; v != nil && v != "" {
			playReadyLicense.ContentType = contentkeypolicies.ContentKeyPolicyPlayReadyContentType(v.(string))
		}

		if v := license["expiration_date"]; v != nil && v != "" {
			expirationDate, err := time.Parse(time.RFC3339, v.(string))
			if err != nil {
				return nil, err
			}
			playReadyLicense.SetExpirationDateAsTime(expirationDate)
		}

		if v := license["grace_period"]; v != nil && v != "" {
			playReadyLicense.GracePeriod = utils.String(v.(string))
		}

		if v := license["license_type"]; v != nil && v != "" {
			playReadyLicense.LicenseType = contentkeypolicies.ContentKeyPolicyPlayReadyLicenseType(v.(string))
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

		if v := license["security_level"]; v != nil && v != "" {
			securityLevel := contentkeypolicies.SecurityLevel(v.(string))
			playReadyLicense.SecurityLevel = &securityLevel
		}

		results = append(results, playReadyLicense)
	}

	return &results, nil
}

func flattenPlayReadyLicenses(input []contentkeypolicies.ContentKeyPolicyPlayReadyLicense) []interface{} {
	results := make([]interface{}, 0)
	for _, v := range input {
		beginDate := ""
		if t, err := v.GetBeginDateAsTime(); t != nil && err == nil {
			beginDate = t.Format(time.RFC3339)
		}

		locationFromHeaderEnabled := false
		locationFromKeyID := ""
		if v.ContentKeyLocation != nil {
			if _, ok := v.ContentKeyLocation.(contentkeypolicies.ContentKeyPolicyPlayReadyContentEncryptionKeyFromHeader); ok {
				locationFromHeaderEnabled = true
			}
			if val, ok := v.ContentKeyLocation.(*contentkeypolicies.ContentKeyPolicyPlayReadyContentEncryptionKeyFromKeyIdentifier); ok {
				locationFromKeyID = val.KeyId
			}
		}

		expirationDate := ""
		if t, err := v.GetExpirationDateAsTime(); t != nil && err == nil {
			expirationDate = t.Format(time.RFC3339)
		}

		gracePeriod := ""
		if v.GracePeriod != nil {
			gracePeriod = *v.GracePeriod
		}

		playRight := make([]interface{}, 0)
		if v.PlayRight != nil {
			playRight = flattenPlayRight(v.PlayRight)
		}

		relativeBeginDate := ""
		if v.RelativeBeginDate != nil {
			relativeBeginDate = *v.RelativeBeginDate
		}

		relativeExpirationDate := ""
		if v.RelativeExpirationDate != nil {
			relativeExpirationDate = *v.RelativeExpirationDate
		}

		securityLevel := ""
		if v.SecurityLevel != nil {
			securityLevel = string(*v.SecurityLevel)
		}

		results = append(results, map[string]interface{}{
			"allow_test_devices": v.AllowTestDevices,
			"begin_date":         beginDate,
			"content_key_location_from_header_enabled": locationFromHeaderEnabled,
			"content_key_location_from_key_id":         locationFromKeyID,
			"content_type":                             string(v.ContentType),
			"expiration_date":                          expirationDate,
			"grace_period":                             gracePeriod,
			"license_type":                             string(v.LicenseType),
			"relative_begin_date":                      relativeBeginDate,
			"relative_expiration_date":                 relativeExpirationDate,
			"play_right":                               playRight,
			"security_level":                           securityLevel,
		})
	}

	return results
}

func expandPlayRight(input []interface{}) *contentkeypolicies.ContentKeyPolicyPlayReadyPlayRight {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	playRightConfiguration := input[0].(map[string]interface{})
	playRight := &contentkeypolicies.ContentKeyPolicyPlayReadyPlayRight{
		ExplicitAnalogTelevisionOutputRestriction: expandExplicitAnalogTelevisionOutputRestriction(playRightConfiguration["explicit_analog_television_output_restriction"].([]interface{})),
	}

	if v := playRightConfiguration["agc_and_color_stripe_restriction"]; v != nil {
		playRight.AgcAndColorStripeRestriction = utils.Int64(int64(v.(int)))
	}

	if v := playRightConfiguration["allow_passing_video_content_to_unknown_output"]; v != nil {
		playRight.AllowPassingVideoContentToUnknownOutput = contentkeypolicies.ContentKeyPolicyPlayReadyUnknownOutputPassingOption(v.(string))
	}

	if v := playRightConfiguration["analog_video_opl"]; v != nil && v != 0 {
		playRight.AnalogVideoOpl = utils.Int64(int64(v.(int)))
	}

	if v := playRightConfiguration["compressed_digital_audio_opl"]; v != nil && v != 0 {
		playRight.CompressedDigitalAudioOpl = utils.Int64(int64(v.(int)))
	}

	if v := playRightConfiguration["compressed_digital_video_opl"]; v != nil && v != 0 {
		playRight.CompressedDigitalVideoOpl = utils.Int64(int64(v.(int)))
	}

	if v := playRightConfiguration["digital_video_only_content_restriction"]; v != nil {
		playRight.DigitalVideoOnlyContentRestriction = v.(bool)
	}

	if v := playRightConfiguration["first_play_expiration"]; v != nil && v != "" {
		playRight.FirstPlayExpiration = utils.String(v.(string))
	}

	if v := playRightConfiguration["image_constraint_for_analog_component_video_restriction"]; v != nil {
		playRight.ImageConstraintForAnalogComponentVideoRestriction = v.(bool)
	}

	if v := playRightConfiguration["image_constraint_for_analog_computer_monitor_restriction"]; v != nil {
		playRight.ImageConstraintForAnalogComputerMonitorRestriction = v.(bool)
	}

	if v := playRightConfiguration["scms_restriction"]; v != nil {
		playRight.ScmsRestriction = utils.Int64(int64(v.(int)))
	}
	if v := playRightConfiguration["uncompressed_digital_audio_opl"]; v != nil && v != 0 {
		playRight.UncompressedDigitalAudioOpl = utils.Int64(int64(v.(int)))
	}

	if v := playRightConfiguration["uncompressed_digital_video_opl"]; v != nil && v != 0 {
		playRight.UncompressedDigitalVideoOpl = utils.Int64(int64(v.(int)))
	}

	return playRight
}

func flattenPlayRight(input *contentkeypolicies.ContentKeyPolicyPlayReadyPlayRight) []interface{} {
	agcStripeRestriction := 0
	if input.AgcAndColorStripeRestriction != nil {
		agcStripeRestriction = int(*input.AgcAndColorStripeRestriction)
	}

	analogVideoOpl := 0
	if input.AnalogVideoOpl != nil {
		analogVideoOpl = int(*input.AnalogVideoOpl)
	}

	compressedDigitalAudioOpl := 0
	if input.CompressedDigitalAudioOpl != nil {
		compressedDigitalAudioOpl = int(*input.CompressedDigitalAudioOpl)
	}

	compressedDigitalVideoOpl := 0
	if input.CompressedDigitalVideoOpl != nil {
		compressedDigitalVideoOpl = int(*input.CompressedDigitalVideoOpl)
	}

	firstPlayExpiration := ""
	if input.FirstPlayExpiration != nil {
		firstPlayExpiration = *input.FirstPlayExpiration
	}

	scmsRestriction := 0
	if input.ScmsRestriction != nil {
		scmsRestriction = int(*input.ScmsRestriction)
	}

	uncompressedDigitalAudioOpl := 0
	if input.UncompressedDigitalAudioOpl != nil {
		uncompressedDigitalAudioOpl = int(*input.UncompressedDigitalAudioOpl)
	}

	uncompressedDigitalVideoOpl := 0
	if input.UncompressedDigitalVideoOpl != nil {
		uncompressedDigitalVideoOpl = int(*input.UncompressedDigitalVideoOpl)
	}

	return []interface{}{
		map[string]interface{}{
			"agc_and_color_stripe_restriction":                         agcStripeRestriction,
			"allow_passing_video_content_to_unknown_output":            string(input.AllowPassingVideoContentToUnknownOutput),
			"analog_video_opl":                                         analogVideoOpl,
			"compressed_digital_audio_opl":                             compressedDigitalAudioOpl,
			"compressed_digital_video_opl":                             compressedDigitalVideoOpl,
			"digital_video_only_content_restriction":                   input.DigitalVideoOnlyContentRestriction,
			"explicit_analog_television_output_restriction":            flattenExplicitAnalogTelevisionOutputRestriction(input.ExplicitAnalogTelevisionOutputRestriction),
			"first_play_expiration":                                    firstPlayExpiration,
			"image_constraint_for_analog_component_video_restriction":  input.ImageConstraintForAnalogComponentVideoRestriction,
			"image_constraint_for_analog_computer_monitor_restriction": input.ImageConstraintForAnalogComputerMonitorRestriction,
			"scms_restriction":                                         scmsRestriction,
			"uncompressed_digital_audio_opl":                           uncompressedDigitalAudioOpl,
			"uncompressed_digital_video_opl":                           uncompressedDigitalVideoOpl,
		},
	}
}

func expandExplicitAnalogTelevisionOutputRestriction(input []interface{}) *contentkeypolicies.ContentKeyPolicyPlayReadyExplicitAnalogTelevisionRestriction {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	restriction := input[0].(map[string]interface{})
	result := &contentkeypolicies.ContentKeyPolicyPlayReadyExplicitAnalogTelevisionRestriction{
		BestEffort:        restriction["best_effort_enforced"].(bool),
		ConfigurationData: int64(restriction["control_bits"].(int)),
	}

	return result
}

func flattenExplicitAnalogTelevisionOutputRestriction(input *contentkeypolicies.ContentKeyPolicyPlayReadyExplicitAnalogTelevisionRestriction) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"best_effort_enforced": input.BestEffort,
			"control_bits":         input.ConfigurationData,
		},
	}
}
