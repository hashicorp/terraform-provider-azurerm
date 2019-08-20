package azurerm

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var apimBackendProtocolSsl3 = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Ssl30"
var apimBackendProtocolTls10 = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls10"
var apimBackendProtocolTls11 = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls11"
var apimFrontendProtocolSsl3 = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Ssl30"
var apimFrontendProtocolTls10 = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls10"
var apimFrontendProtocolTls11 = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls11"
var apimTripleDesCiphers = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TripleDes168"

func resourceArmApiManagementService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementServiceCreateUpdate,
		Read:   resourceArmApiManagementServiceRead,
		Update: resourceArmApiManagementServiceCreateUpdate,
		Delete: resourceArmApiManagementServiceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": azure.SchemaApiManagementName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"public_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"publisher_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ApiManagementServicePublisherName,
			},

			"publisher_email": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ApiManagementServicePublisherEmail,
			},

			"sku": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(apimanagement.SkuTypeDeveloper),
								string(apimanagement.SkuTypeBasic),
								string(apimanagement.SkuTypeStandard),
								string(apimanagement.SkuTypePremium),
							}, false),
						},
						"capacity": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
					},
				},
			},

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"SystemAssigned",
							}, false),
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"notification_sender_email": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"additional_location": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location": azure.SchemaLocation(),

						"gateway_regional_url": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"public_ip_addresses": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed: true,
						},
					},
				},
			},

			"certificate": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"encoded_certificate": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},

						"certificate_password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},

						"store_name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(apimanagement.CertificateAuthority),
								string(apimanagement.Root),
							}, false),
						},
					},
				},
			},

			"security": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true, // todo remove in 2.0 ?
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disable_backend_ssl30": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"disable_backend_tls10": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"disable_backend_tls11": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"disable_triple_des_chipers": {
							Type:          schema.TypeBool,
							Optional:      true,
							Computed:      true, // todo remove in 2.0
							Deprecated:    "This field has been deprecated in favour of the `disable_triple_des_ciphers` property to correct the spelling. it will be removed in version 2.0 of the provider",
							ConflictsWith: []string{"security.0.disable_triple_des_ciphers"},
						},
						"disable_triple_des_ciphers": {
							Type:     schema.TypeBool,
							Optional: true,
							// Default:       false, // todo remove in 2.0
							Computed:      true, // todo remove in 2.0
							ConflictsWith: []string{"security.0.disable_triple_des_chipers"},
						},
						"disable_frontend_ssl30": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"disable_frontend_tls10": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"disable_frontend_tls11": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"hostname_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"management": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: apiManagementResourceHostnameSchema("management"),
							},
						},
						"portal": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: apiManagementResourceHostnameSchema("portal"),
							},
						},
						"proxy": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: apiManagementResourceHostnameProxySchema(),
							},
						},
						"scm": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: apiManagementResourceHostnameSchema("scm"),
							},
						},
					},
				},
			},

			"policy": {
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				MaxItems:   1,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"xml_content": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ConflictsWith:    []string{"policy.0.xml_link"},
							DiffSuppressFunc: suppress.XmlDiff,
						},

						"xml_link": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"policy.0.xml_content"},
						},
					},
				},
			},

			"sign_in": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},

			"sign_up": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"terms_of_service": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"consent_required": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"text": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"tags": tagsSchema(),

			"gateway_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"gateway_regional_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"portal_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"management_api_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"scm_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmApiManagementServiceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.ServiceClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for API Management Service creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing API Management Service %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})

	sku := expandAzureRmApiManagementSku(d)

	publisherName := d.Get("publisher_name").(string)
	publisherEmail := d.Get("publisher_email").(string)
	notificationSenderEmail := d.Get("notification_sender_email").(string)

	customProperties := expandApiManagementCustomProperties(d)
	certificates := expandAzureRmApiManagementCertificates(d)
	hostnameConfigurations := expandAzureRmApiManagementHostnameConfigurations(d)

	properties := apimanagement.ServiceResource{
		Location: utils.String(location),
		ServiceProperties: &apimanagement.ServiceProperties{
			PublisherName:          utils.String(publisherName),
			PublisherEmail:         utils.String(publisherEmail),
			CustomProperties:       customProperties,
			Certificates:           certificates,
			HostnameConfigurations: hostnameConfigurations,
		},
		Tags: expandTags(tags),
		Sku:  sku,
	}

	if _, ok := d.GetOk("identity"); ok {
		properties.Identity = expandAzureRmApiManagementIdentity(d)
	}

	if _, ok := d.GetOk("additional_location"); ok {
		properties.ServiceProperties.AdditionalLocations = expandAzureRmApiManagementAdditionalLocations(d, sku)
	}

	if notificationSenderEmail != "" {
		properties.ServiceProperties.NotificationSenderEmail = &notificationSenderEmail
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, properties)
	if err != nil {
		return fmt.Errorf("Error creating/updating API Management Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation/update of API Management Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving API Management Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ID for API Management Service %q (Resource Group %q)", name, resourceGroup)
	}

	d.SetId(*read.ID)

	signInSettingsRaw := d.Get("sign_in").([]interface{})
	signInSettings := expandApiManagementSignInSettings(signInSettingsRaw)
	signInClient := meta.(*ArmClient).apiManagement.SignInClient
	if _, err := signInClient.CreateOrUpdate(ctx, resourceGroup, name, signInSettings); err != nil {
		return fmt.Errorf("Error setting Sign In settings for API Management Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	signUpSettingsRaw := d.Get("sign_up").([]interface{})
	signUpSettings := expandApiManagementSignUpSettings(signUpSettingsRaw)
	signUpClient := meta.(*ArmClient).apiManagement.SignUpClient
	if _, err := signUpClient.CreateOrUpdate(ctx, resourceGroup, name, signUpSettings); err != nil {
		return fmt.Errorf("Error setting Sign Up settings for API Management Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	policyClient := meta.(*ArmClient).apiManagement.PolicyClient
	policiesRaw := d.Get("policy").([]interface{})
	policy, err := expandApiManagementPolicies(policiesRaw)
	if err != nil {
		return err
	}

	if d.HasChange("policy") {
		// remove the existing policy
		if resp, err := policyClient.Delete(ctx, resourceGroup, name, ""); err != nil {
			if !utils.ResponseWasNotFound(resp) {
				return fmt.Errorf("Error removing Policies from API Management Service %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		// then add the new one, if it exists
		if policy != nil {
			if _, err := policyClient.CreateOrUpdate(ctx, resourceGroup, name, *policy); err != nil {
				return fmt.Errorf("Error setting Policies for API Management Service %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
	}

	return resourceArmApiManagementServiceRead(d, meta)
}

func resourceArmApiManagementServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.ServiceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["service"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("API Management Service %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on API Management Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	signInClient := meta.(*ArmClient).apiManagement.SignInClient
	signInSettings, err := signInClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Sign In Settings for API Management Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	signUpClient := meta.(*ArmClient).apiManagement.SignUpClient
	signUpSettings, err := signUpClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Sign Up Settings for API Management Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	policyClient := meta.(*ArmClient).apiManagement.PolicyClient
	policy, err := policyClient.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(policy.Response) {
			return fmt.Errorf("Error retrieving Policy for API Management Service %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	identity := flattenAzureRmApiManagementMachineIdentity(resp.Identity)
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("Error setting `identity`: %+v", err)
	}

	if props := resp.ServiceProperties; props != nil {
		d.Set("publisher_email", props.PublisherEmail)
		d.Set("publisher_name", props.PublisherName)
		d.Set("notification_sender_email", props.NotificationSenderEmail)
		d.Set("gateway_url", props.GatewayURL)
		d.Set("gateway_regional_url", props.GatewayRegionalURL)
		d.Set("portal_url", props.PortalURL)
		d.Set("management_api_url", props.ManagementAPIURL)
		d.Set("scm_url", props.ScmURL)
		d.Set("public_ip_addresses", props.PublicIPAddresses)

		if err := d.Set("security", flattenApiManagementCustomProperties(props.CustomProperties)); err != nil {
			return fmt.Errorf("Error setting `security`: %+v", err)
		}

		hostnameConfigs := flattenApiManagementHostnameConfigurations(props.HostnameConfigurations, d)
		if err := d.Set("hostname_configuration", hostnameConfigs); err != nil {
			return fmt.Errorf("Error setting `hostname_configuration`: %+v", err)
		}

		if err := d.Set("additional_location", flattenApiManagementAdditionalLocations(props.AdditionalLocations)); err != nil {
			return fmt.Errorf("Error setting `additional_location`: %+v", err)
		}
	}

	if err := d.Set("sku", flattenApiManagementServiceSku(resp.Sku)); err != nil {
		return fmt.Errorf("Error setting `sku`: %+v", err)
	}

	if err := d.Set("sign_in", flattenApiManagementSignInSettings(signInSettings)); err != nil {
		return fmt.Errorf("Error setting `sign_in`: %+v", err)
	}

	if err := d.Set("sign_up", flattenApiManagementSignUpSettings(signUpSettings)); err != nil {
		return fmt.Errorf("Error setting `sign_up`: %+v", err)
	}

	flattenAndSetTags(d, resp.Tags)

	if err := d.Set("policy", flattenApiManagementPolicies(d, policy)); err != nil {
		return fmt.Errorf("Error setting `policy`: %+v", err)
	}

	return nil
}

func resourceArmApiManagementServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.ServiceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["service"]

	log.Printf("[DEBUG] Deleting API Management Service %q (Resource Grouo %q)", name, resourceGroup)
	resp, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting API Management Service %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandAzureRmApiManagementHostnameConfigurations(d *schema.ResourceData) *[]apimanagement.HostnameConfiguration {
	results := make([]apimanagement.HostnameConfiguration, 0)
	hostnameVs := d.Get("hostname_configuration").([]interface{})

	for _, hostnameRawVal := range hostnameVs {
		hostnameV := hostnameRawVal.(map[string]interface{})

		managementVs := hostnameV["management"].([]interface{})
		for _, managementV := range managementVs {
			v := managementV.(map[string]interface{})
			output := expandApiManagementCommonHostnameConfiguration(v, apimanagement.Management)
			results = append(results, output)
		}

		portalVs := hostnameV["portal"].([]interface{})
		for _, portalV := range portalVs {
			v := portalV.(map[string]interface{})
			output := expandApiManagementCommonHostnameConfiguration(v, apimanagement.Portal)
			results = append(results, output)
		}

		proxyVs := hostnameV["proxy"].([]interface{})
		for _, proxyV := range proxyVs {
			v := proxyV.(map[string]interface{})
			output := expandApiManagementCommonHostnameConfiguration(v, apimanagement.Proxy)
			if value, ok := v["default_ssl_binding"]; ok {
				output.DefaultSslBinding = utils.Bool(value.(bool))
			}
			results = append(results, output)
		}

		scmVs := hostnameV["scm"].([]interface{})
		for _, scmV := range scmVs {
			v := scmV.(map[string]interface{})
			output := expandApiManagementCommonHostnameConfiguration(v, apimanagement.Scm)
			results = append(results, output)
		}
	}

	return &results
}

func expandApiManagementCommonHostnameConfiguration(input map[string]interface{}, hostnameType apimanagement.HostnameType) apimanagement.HostnameConfiguration {
	encodedCertificate := input["certificate"].(string)
	certificatePassword := input["certificate_password"].(string)
	hostName := input["host_name"].(string)
	keyVaultId := input["key_vault_id"].(string)

	output := apimanagement.HostnameConfiguration{
		EncodedCertificate:  utils.String(encodedCertificate),
		CertificatePassword: utils.String(certificatePassword),
		HostName:            utils.String(hostName),
		KeyVaultID:          utils.String(keyVaultId),
		Type:                hostnameType,
	}

	if v, ok := input["negotiate_client_certificate"]; ok {
		output.NegotiateClientCertificate = utils.Bool(v.(bool))
	}

	return output
}

func flattenApiManagementHostnameConfigurations(input *[]apimanagement.HostnameConfiguration, d *schema.ResourceData) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	managementResults := make([]interface{}, 0)
	portalResults := make([]interface{}, 0)
	proxyResults := make([]interface{}, 0)
	scmResults := make([]interface{}, 0)

	for _, config := range *input {
		output := make(map[string]interface{})

		if config.HostName != nil {
			output["host_name"] = *config.HostName
		}

		if config.NegotiateClientCertificate != nil {
			output["negotiate_client_certificate"] = *config.NegotiateClientCertificate
		}

		if config.KeyVaultID != nil {
			output["key_vault_id"] = *config.KeyVaultID
		}

		// Iterate through old state to find sensitive props not returned by API.
		// This must be done in order to avoid state diffs.
		// NOTE: this information won't be available during times like Import, so this is a best-effort.
		existingHostnames := d.Get("hostname_configuration").([]interface{})
		if len(existingHostnames) > 0 {
			v := existingHostnames[0].(map[string]interface{})

			if valsRaw, ok := v[strings.ToLower(string(config.Type))]; ok {
				vals := valsRaw.([]interface{})
				for _, val := range vals {
					oldConfig := val.(map[string]interface{})

					if oldConfig["host_name"] == *config.HostName {
						output["certificate_password"] = oldConfig["certificate_password"]
						output["certificate"] = oldConfig["certificate"]
					}
				}
			}
		}

		switch strings.ToLower(string(config.Type)) {
		case strings.ToLower(string(apimanagement.Proxy)):
			// only set SSL binding for proxy types
			if config.DefaultSslBinding != nil {
				output["default_ssl_binding"] = *config.DefaultSslBinding
			}
			proxyResults = append(proxyResults, output)

		case strings.ToLower(string(apimanagement.Management)):
			managementResults = append(managementResults, output)

		case strings.ToLower(string(apimanagement.Portal)):
			portalResults = append(portalResults, output)

		case strings.ToLower(string(apimanagement.Scm)):
			scmResults = append(scmResults, output)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"management": managementResults,
			"portal":     portalResults,
			"proxy":      proxyResults,
			"scm":        scmResults,
		},
	}
}

func expandAzureRmApiManagementCertificates(d *schema.ResourceData) *[]apimanagement.CertificateConfiguration {
	vs := d.Get("certificate").([]interface{})

	results := make([]apimanagement.CertificateConfiguration, 0)

	for _, v := range vs {
		config := v.(map[string]interface{})

		certBase64 := config["encoded_certificate"].(string)
		certificatePassword := config["certificate_password"].(string)
		storeName := apimanagement.StoreName(config["store_name"].(string))

		cert := apimanagement.CertificateConfiguration{
			EncodedCertificate:  utils.String(certBase64),
			CertificatePassword: utils.String(certificatePassword),
			StoreName:           storeName,
		}

		results = append(results, cert)
	}

	return &results
}

func expandAzureRmApiManagementAdditionalLocations(d *schema.ResourceData, sku *apimanagement.ServiceSkuProperties) *[]apimanagement.AdditionalLocation {
	inputLocations := d.Get("additional_location").([]interface{})

	additionalLocations := make([]apimanagement.AdditionalLocation, 0)

	for _, v := range inputLocations {
		config := v.(map[string]interface{})
		location := azure.NormalizeLocation(config["location"].(string))

		additionalLocation := apimanagement.AdditionalLocation{
			Location: utils.String(location),
			Sku:      sku,
		}

		additionalLocations = append(additionalLocations, additionalLocation)
	}

	return &additionalLocations
}

func flattenApiManagementAdditionalLocations(input *[]apimanagement.AdditionalLocation) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, prop := range *input {
		output := make(map[string]interface{})

		if prop.Location != nil {
			output["location"] = azure.NormalizeLocation(*prop.Location)
		}

		if prop.PublicIPAddresses != nil {
			output["public_ip_addresses"] = *prop.PublicIPAddresses
		}

		if prop.GatewayRegionalURL != nil {
			output["gateway_regional_url"] = *prop.GatewayRegionalURL
		}

		results = append(results, output)
	}

	return results
}

func expandAzureRmApiManagementIdentity(d *schema.ResourceData) *apimanagement.ServiceIdentity {
	vs := d.Get("identity").([]interface{})
	if len(vs) == 0 {
		return nil
	}

	v := vs[0].(map[string]interface{})
	identityType := v["type"].(string)
	return &apimanagement.ServiceIdentity{
		Type: utils.String(identityType),
	}
}

func flattenAzureRmApiManagementMachineIdentity(identity *apimanagement.ServiceIdentity) []interface{} {
	if identity == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if identity.Type != nil {
		result["type"] = *identity.Type
	}

	if identity.PrincipalID != nil {
		result["principal_id"] = identity.PrincipalID.String()
	}

	if identity.TenantID != nil {
		result["tenant_id"] = identity.TenantID.String()
	}

	return []interface{}{result}
}

func expandAzureRmApiManagementSku(d *schema.ResourceData) *apimanagement.ServiceSkuProperties {
	vs := d.Get("sku").([]interface{})
	// guaranteed by MinItems in the schema
	v := vs[0].(map[string]interface{})

	name := apimanagement.SkuType(v["name"].(string))
	capacity := int32(v["capacity"].(int))

	sku := &apimanagement.ServiceSkuProperties{
		Name:     name,
		Capacity: utils.Int32(capacity),
	}

	return sku
}

func flattenApiManagementServiceSku(input *apimanagement.ServiceSkuProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	sku := make(map[string]interface{})

	sku["name"] = string(input.Name)

	if input.Capacity != nil {
		sku["capacity"] = *input.Capacity
	}

	return []interface{}{sku}
}

func expandApiManagementCustomProperties(d *schema.ResourceData) map[string]*string {
	vs := d.Get("security").([]interface{})

	backendProtocolSsl3 := false
	backendProtocolTls10 := false
	backendProtocolTls11 := false
	frontendProtocolSsl3 := false
	frontendProtocolTls10 := false
	frontendProtocolTls11 := false
	tripleDesCiphers := false

	if len(vs) > 0 {
		v := vs[0].(map[string]interface{})
		backendProtocolSsl3 = v["disable_backend_ssl30"].(bool)
		backendProtocolTls10 = v["disable_backend_tls10"].(bool)
		backendProtocolTls11 = v["disable_backend_tls11"].(bool)
		frontendProtocolSsl3 = v["disable_frontend_ssl30"].(bool)
		frontendProtocolTls10 = v["disable_frontend_tls10"].(bool)
		frontendProtocolTls11 = v["disable_frontend_tls11"].(bool)
		//tripleDesCiphers = v["disable_triple_des_ciphers"].(bool) //restore in 2.0
	}

	if c, ok := d.GetOkExists("security.0.disable_triple_des_ciphers"); ok {
		tripleDesCiphers = c.(bool)
	} else if c, ok := d.GetOkExists("security.0.disable_triple_des_chipers"); ok {
		tripleDesCiphers = c.(bool)
	}

	return map[string]*string{
		apimBackendProtocolSsl3:   utils.String(strconv.FormatBool(backendProtocolSsl3)),
		apimBackendProtocolTls10:  utils.String(strconv.FormatBool(backendProtocolTls10)),
		apimBackendProtocolTls11:  utils.String(strconv.FormatBool(backendProtocolTls11)),
		apimFrontendProtocolSsl3:  utils.String(strconv.FormatBool(frontendProtocolSsl3)),
		apimFrontendProtocolTls10: utils.String(strconv.FormatBool(frontendProtocolTls10)),
		apimFrontendProtocolTls11: utils.String(strconv.FormatBool(frontendProtocolTls11)),
		apimTripleDesCiphers:      utils.String(strconv.FormatBool(tripleDesCiphers)),
	}
}

func flattenApiManagementCustomProperties(input map[string]*string) []interface{} {
	output := make(map[string]interface{})

	output["disable_backend_ssl30"] = parseApiManagementNilableDictionary(input, apimBackendProtocolSsl3)
	output["disable_backend_tls10"] = parseApiManagementNilableDictionary(input, apimBackendProtocolTls10)
	output["disable_backend_tls11"] = parseApiManagementNilableDictionary(input, apimBackendProtocolTls11)
	output["disable_frontend_ssl30"] = parseApiManagementNilableDictionary(input, apimFrontendProtocolSsl3)
	output["disable_frontend_tls10"] = parseApiManagementNilableDictionary(input, apimFrontendProtocolTls10)
	output["disable_frontend_tls11"] = parseApiManagementNilableDictionary(input, apimFrontendProtocolTls11)
	output["disable_triple_des_chipers"] = parseApiManagementNilableDictionary(input, apimTripleDesCiphers) // todo remove in 2.0
	output["disable_triple_des_ciphers"] = parseApiManagementNilableDictionary(input, apimTripleDesCiphers)

	return []interface{}{output}
}

func apiManagementResourceHostnameSchema(schemaName string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host_name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validate.NoEmptyStrings,
		},

		"key_vault_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: azure.ValidateKeyVaultChildId,
			ConflictsWith: []string{
				fmt.Sprintf("hostname_configuration.0.%s.0.certificate", schemaName),
				fmt.Sprintf("hostname_configuration.0.%s.0.certificate_password", schemaName),
			},
		},

		"certificate": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validate.NoEmptyStrings,
			ConflictsWith: []string{
				fmt.Sprintf("hostname_configuration.0.%s.0.key_vault_id", schemaName),
			},
		},

		"certificate_password": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validate.NoEmptyStrings,
			ConflictsWith: []string{
				fmt.Sprintf("hostname_configuration.0.%s.0.key_vault_id", schemaName),
			},
		},

		"negotiate_client_certificate": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func apiManagementResourceHostnameProxySchema() map[string]*schema.Schema {
	hostnameSchema := apiManagementResourceHostnameSchema("proxy")

	hostnameSchema["default_ssl_binding"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Computed: true, // Azure has certain logic to set this, which we cannot predict
	}

	return hostnameSchema
}

func parseApiManagementNilableDictionary(input map[string]*string, key string) bool {
	log.Printf("Parsing value for %q", key)

	v, ok := input[key]
	if !ok {
		log.Printf("%q was not found in the input - returning `false` as the default value", key)
		return false
	}

	val, err := strconv.ParseBool(*v)
	if err != nil {
		log.Printf("Error parsing %q (key %q) as bool: %+v - assuming false", key, *v, err)
		return false
	}

	return val
}

func expandApiManagementSignInSettings(input []interface{}) apimanagement.PortalSigninSettings {
	enabled := false

	if len(input) > 0 {
		vs := input[0].(map[string]interface{})
		enabled = vs["enabled"].(bool)
	}

	return apimanagement.PortalSigninSettings{
		PortalSigninSettingProperties: &apimanagement.PortalSigninSettingProperties{
			Enabled: utils.Bool(enabled),
		},
	}
}

func flattenApiManagementSignInSettings(input apimanagement.PortalSigninSettings) []interface{} {
	enabled := false

	if props := input.PortalSigninSettingProperties; props != nil {
		if props.Enabled != nil {
			enabled = *props.Enabled
		}
	}

	return []interface{}{
		map[string]interface{}{
			"enabled": enabled,
		},
	}
}

func expandApiManagementSignUpSettings(input []interface{}) apimanagement.PortalSignupSettings {
	if len(input) == 0 {
		return apimanagement.PortalSignupSettings{
			PortalSignupSettingsProperties: &apimanagement.PortalSignupSettingsProperties{
				Enabled: utils.Bool(false),
				TermsOfService: &apimanagement.TermsOfServiceProperties{
					ConsentRequired: utils.Bool(false),
					Enabled:         utils.Bool(false),
					Text:            utils.String(""),
				},
			},
		}
	}

	vs := input[0].(map[string]interface{})

	props := apimanagement.PortalSignupSettingsProperties{
		Enabled: utils.Bool(vs["enabled"].(bool)),
	}

	termsOfServiceRaw := vs["terms_of_service"].([]interface{})
	if len(termsOfServiceRaw) > 0 {
		termsOfServiceVs := termsOfServiceRaw[0].(map[string]interface{})
		props.TermsOfService = &apimanagement.TermsOfServiceProperties{
			Enabled:         utils.Bool(termsOfServiceVs["enabled"].(bool)),
			ConsentRequired: utils.Bool(termsOfServiceVs["consent_required"].(bool)),
			Text:            utils.String(termsOfServiceVs["text"].(string)),
		}
	}

	return apimanagement.PortalSignupSettings{
		PortalSignupSettingsProperties: &props,
	}
}

func flattenApiManagementSignUpSettings(input apimanagement.PortalSignupSettings) []interface{} {
	enabled := false
	termsOfService := make([]interface{}, 0)

	if props := input.PortalSignupSettingsProperties; props != nil {
		if props.Enabled != nil {
			enabled = *props.Enabled
		}

		if tos := props.TermsOfService; tos != nil {
			output := make(map[string]interface{})

			if tos.Enabled != nil {
				output["enabled"] = *tos.Enabled
			}

			if tos.ConsentRequired != nil {
				output["consent_required"] = *tos.ConsentRequired
			}

			if tos.Text != nil {
				output["text"] = *tos.Text
			}

			termsOfService = append(termsOfService, output)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"enabled":          enabled,
			"terms_of_service": termsOfService,
		},
	}
}

func expandApiManagementPolicies(input []interface{}) (*apimanagement.PolicyContract, error) {
	if len(input) == 0 {
		return nil, nil
	}

	vs := input[0].(map[string]interface{})
	xmlContent := vs["xml_content"].(string)
	xmlLink := vs["xml_link"].(string)

	if xmlContent != "" {
		return &apimanagement.PolicyContract{
			PolicyContractProperties: &apimanagement.PolicyContractProperties{
				ContentFormat: apimanagement.XML,
				PolicyContent: utils.String(xmlContent),
			},
		}, nil
	}

	if xmlLink != "" {
		return &apimanagement.PolicyContract{
			PolicyContractProperties: &apimanagement.PolicyContractProperties{
				ContentFormat: apimanagement.XMLLink,
				PolicyContent: utils.String(xmlLink),
			},
		}, nil
	}

	return nil, fmt.Errorf("Either `xml_content` or `xml_link` should be set if the `policy` block is defined.")
}

func flattenApiManagementPolicies(d *schema.ResourceData, input apimanagement.PolicyContract) []interface{} {
	xmlContent := ""
	if props := input.PolicyContractProperties; props != nil {
		if props.PolicyContent != nil {
			xmlContent = *props.PolicyContent
		}
	}

	// if there's no policy assigned, we set this to an empty list
	if xmlContent == "" {
		return []interface{}{}
	}

	output := map[string]interface{}{
		"xml_content": xmlContent,
		"xml_link":    "",
	}

	// when you submit an `xml_link` to the API, the API downloads this link and stores it as `xml_content`
	// as such we need to retrieve this value from the state if it's present
	if existing, ok := d.GetOk("policy"); ok {
		existingVs := existing.([]interface{})
		if len(existingVs) > 0 {
			existingV := existingVs[0].(map[string]interface{})
			output["xml_link"] = existingV["xml_link"].(string)
		}
	}

	return []interface{}{output}
}
