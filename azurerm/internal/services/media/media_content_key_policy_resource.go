package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2020-05-01/media"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				Type:     schema.TypeList,
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
										ValidateFunc: validation.ValidateRFC3339TimeString,
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
										ValidateFunc: validation.ValidateRFC3339TimeString,
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
									"play_right": { //TODO:Complete definition
										Type:     schema.TypeList,
										Optional: true,
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
									"relative_begin_date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.ValidateRFC3339TimeString,
									},
									"relative_expiration_date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.ValidateRFC3339TimeString,
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
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"pfx": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"pfx_password": {
										Type:         schema.TypeString,
										Optional:     true,
										Sensitive:    true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"offline_rental_configuration": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(media.ContentKeyPolicyRestrictionTokenTypeJwt),
											string(media.ContentKeyPolicyRestrictionTokenTypeSwt),
											string(media.ContentKeyPolicyRestrictionTokenTypeUnknown),
										}, false),
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
									"rental_duration": {
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
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"primary_rsa_token_key_exponent": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"primary_rsa_token_key_modulus": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"primary_x509_token_key_raw": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
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
		options, err := expandPolicyOptions(v.([]interface{}))
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

	if props := resp.ContentKeyPolicyProperties; props != nil {
		d.Set("description", props.Description)

		options := flattenPolicyOptions(resp.Options)
		if err := d.Set("policy_option", options); err != nil {
			return fmt.Errorf("Error flattening `policy_option`: %s", err)
		}
	}

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

func flattenPolicyOptions(input *[]media.ContentKeyPolicyOption) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, option := range *input {
		policyOption := make(map[string]interface{})
		policyOption["name"] = option.Name

		configuration, err := flattenConfiguration(option.Configuration, policyOption)
		if err != nil {

		}
		policyOption = configuration

		restriction, err := flattenRestriction(option.Restriction, policyOption)
		if err != nil {

		}
		policyOption = restriction

		results = append(results, policyOption)
	}

	return results
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

		option["token_restriction"] = []interface{}{
			map[string]interface{}{
				"audience":                           audience,
				"issuer":                             issuer,
				"token_type":                         string(token.RestrictionTokenType),
				"open_id_connect_discovery_document": openIDConnectDiscoveryDocument,
				"required_claim":                     requiredClaims,
			},
		}

		option = flattenVerificationKey(token.PrimaryVerificationKey, option)
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
	default:
		return nil, fmt.Errorf("policy_option must contain at least one type of configuration: clear_key_configuration_enabled , widevine_configuration_template, playready_configuration_license or fairplay_configuration.")
	}
}

func flattenConfiguration(input media.BasicContentKeyPolicyConfiguration, option map[string]interface{}) (map[string]interface{}, error) {
	if input == nil {
		return option, nil
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
	}

	return option, nil
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
		return nil, fmt.Errorf("more than one type of token key in the same token_restriction is not allowed.", verificationKeyType)
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
