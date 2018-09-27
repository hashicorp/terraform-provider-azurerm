package azurerm

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/apimanagement/mgmt/2018-06-01-preview/apimanagement"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApiManagementService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementServiceCreateUpdate,
		Read:   resourceArmApiManagementServiceRead,
		Update: resourceArmApiManagementServiceCreateUpdate,
		Delete: resourceArmApiManagementDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateApiManagementName,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

			"publisher_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateApiManagementPublisherName,
			},

			"publisher_email": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateApiManagementPublisherEmail,
			},

			"sku": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(apimanagement.SkuTypeDeveloper),
							ValidateFunc: validation.StringInSlice([]string{
								string(apimanagement.SkuTypeDeveloper),
								string(apimanagement.SkuTypeBasic),
								string(apimanagement.SkuTypeStandard),
								string(apimanagement.SkuTypePremium),
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},
						"capacity": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
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
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							ValidateFunc: validation.StringInSlice([]string{
								"SystemAssigned",
							}, true),
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
						"location": locationSchema(),

						"gateway_regional_url": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"static_ips": {
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
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},
					},
				},
			},

			"security": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disable_backend_ssl30": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"disable_backend_tls10": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"disable_backend_tls11": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"disable_triple_des_chipers": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"disable_frontend_ssl30": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"disable_frontend_tls10": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"disable_frontend_tls11": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"hostname_configurations": {
				Type:     schema.TypeList,
				Optional: true,
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
								Schema: apiManagementResourceHostnameProxySchema("proxy"),
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
	client := meta.(*ArmClient).apiManagementServiceClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM API Management creation.")

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})

	skuRaw := d.Get("sku").([]interface{})
	sku := expandAzureRmApiManagementSku(skuRaw)

	publisher_name := d.Get("publisher_name").(string)
	publisher_email := d.Get("publisher_email").(string)
	notification_sender_email := d.Get("notification_sender_email").(string)

	custom_properties := expandApiManagementCustomProperties(d)

	certificates := expandAzureRmApiManagementCertificates(d)
	hostname_configurations := expandAzureRmApiManagementHostnameConfigurations(d)

	apiManagement := apimanagement.ServiceResource{
		Location: &location,
		ServiceProperties: &apimanagement.ServiceProperties{
			PublisherName:          utils.String(publisher_name),
			PublisherEmail:         utils.String(publisher_email),
			CustomProperties:       custom_properties,
			Certificates:           certificates,
			HostnameConfigurations: hostname_configurations,
		},
		Tags: expandTags(tags),
		Sku:  sku,
	}

	if _, ok := d.GetOk("identity"); ok {
		apiManagement.Identity = expandAzureRmApiManagementIdentity(d)
	}

	if _, ok := d.GetOk("additional_location"); ok {
		apiManagement.ServiceProperties.AdditionalLocations = expandAzureRmApiManagementAdditionalLocations(d, sku)
	}

	if notification_sender_email != "" {
		apiManagement.ServiceProperties.NotificationSenderEmail = &notification_sender_email
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, apiManagement)
	if err != nil {
		return fmt.Errorf("Error creating API Management Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err := future.WaitForCompletion(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of API Management Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving API Management Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read AzureRM Api Management %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmApiManagementServiceRead(d, meta)
}

func resourceArmApiManagementServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	apiManagementClient := meta.(*ArmClient).apiManagementServiceClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["service"]

	ctx := client.StopContext
	resp, err := apiManagementClient.Get(ctx, resGroup, name)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on API Management Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	identity := flattenAzureRmApiManagementMachineIdentity(resp.Identity)
	if err := d.Set("identity", identity); err != nil {
		return err
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
		d.Set("static_ips", props.PublicIPAddresses)

		customProps, err := flattenApiManagementCustomProperties(props.CustomProperties)
		if err != nil {
			return err
		}

		if customProps != nil {
			if err := d.Set("security", customProps); err != nil {
				return fmt.Errorf("Error setting `security`: %+v", err)
			}
		}

		// Azure is not returning identifiable certs (no id), so there is no way to link
		// what is configured in terraform to Azure. Because of this, we just return the old state.
		if err := d.Set("certificate", d.Get("certificate").([]interface{})); err != nil {
			return fmt.Errorf("Error setting `certificate`: %+v", err)
		}

		if hostnameConfigs := flattenApiManagementHostnameConfigurations(d, props.HostnameConfigurations); hostnameConfigs != nil {
			if err := d.Set("hostname_configurations", hostnameConfigs); err != nil {
				return fmt.Errorf("Error setting `hostname_configurations`: %+v", err)
			}
		}

		if err := d.Set("additional_location", flattenApiManagementAdditionalLocations(props.AdditionalLocations)); err != nil {
			return fmt.Errorf("Error setting `additional_location`: %+v", err)
		}
	}

	if sku := resp.Sku; sku != nil {
		if err := d.Set("sku", flattenApiManagementServiceSku(sku)); err != nil {
			return fmt.Errorf("Error setting `sku`: %+v", err)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmApiManagementDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagementServiceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["service"]

	log.Printf("[DEBUG] Deleting api management %s: %s", resGroup, name)

	resp, err := client.Delete(ctx, resGroup, name)

	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return err
	}

	return nil
}

func expandAzureRmApiManagementIdentity(d *schema.ResourceData) *apimanagement.ServiceIdentity {
	identities := d.Get("identity").([]interface{})
	identity := identities[0].(map[string]interface{})
	identityType := identity["type"].(string)
	return &apimanagement.ServiceIdentity{
		Type: &identityType,
	}
}

func expandApiManagementCustomProperties(d *schema.ResourceData) map[string]*string {
	customProps := map[string]*string{
		"Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Ssl30": utils.String("false"),
		"Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls10": utils.String("false"),
		"Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls11": utils.String("false"),
		"Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TripleDes168":    utils.String("false"),
		"Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Ssl30":         utils.String("false"),
		"Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls10":         utils.String("false"),
		"Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls11":         utils.String("false"),
	}

	if v, ok := d.GetOk("security.0.disable_backend_ssl30"); ok {
		val := strings.Title(strconv.FormatBool(v.(bool)))
		customProps["Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Ssl30"] = utils.String(val)
	}

	if v, ok := d.GetOk("security.0.disable_backend_tls10"); ok {
		val := strings.Title(strconv.FormatBool(v.(bool)))
		customProps["Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls10"] = utils.String(val)
	}

	if v, ok := d.GetOk("security.0.disable_backend_tls11"); ok {
		val := strings.Title(strconv.FormatBool(v.(bool)))
		customProps["Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls11"] = utils.String(val)
	}

	if v, ok := d.GetOk("security.0.disable_triple_des_chipers"); ok {
		val := strings.Title(strconv.FormatBool(v.(bool)))
		customProps["Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TripleDes168"] = utils.String(val)
	}

	if v, ok := d.GetOk("security.0.disable_frontend_ssl30"); ok {
		val := strings.Title(strconv.FormatBool(v.(bool)))
		customProps["Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Ssl30"] = utils.String(val)
	}

	if v, ok := d.GetOk("security.0.disable_frontend_tls10"); ok {
		val := strings.Title(strconv.FormatBool(v.(bool)))
		customProps["Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls10"] = utils.String(val)
	}

	if v, ok := d.GetOk("security.0.disable_frontend_tls11"); ok {
		val := strings.Title(strconv.FormatBool(v.(bool)))
		customProps["Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls11"] = utils.String(val)
	}

	return customProps
}

func expandAzureRmApiManagementHostnameConfigurations(d *schema.ResourceData) *[]apimanagement.HostnameConfiguration {
	if _, ok := d.GetOk("hostname_configurations"); ok {
		log.Printf("Config for hostname_configuration found")

		hostnames := make([]apimanagement.HostnameConfiguration, 0)

		for _, hostNameType := range apimanagement.PossibleHostnameTypeValues() {
			sHostNameType := strings.ToLower(string(hostNameType))
			key := fmt.Sprintf("hostname_configurations.0.%s", sHostNameType)

			log.Printf("looking up key: %v", key)

			if v, ok := d.GetOk(key); ok {
				log.Printf("Config for %s found", key)
				configs := v.([]interface{})

				for _, v := range configs {
					config := v.(map[string]interface{})

					hostname := apimanagement.HostnameConfiguration{
						Type: hostNameType,
					}

					hostname.HostName = utils.String(config["host_name"].(string))
					hostname.EncodedCertificate = utils.String(config["certificate"].(string))
					hostname.CertificatePassword = utils.String(config["certificate_password"].(string))
					hostname.KeyVaultID = utils.String(config["key_vault_id"].(string))

					if v, ok := config["default_ssl_binding"]; ok {
						hostname.DefaultSslBinding = utils.Bool(v.(bool))
					}

					if v, ok := config["negotiate_client_certificate"]; ok {
						hostname.NegotiateClientCertificate = utils.Bool(v.(bool))
					}

					log.Printf("Here's the config: %v", &hostname)

					hostnames = append(hostnames, hostname)
				}
			}
		}
		return &hostnames
	}
	return nil
}

func expandAzureRmApiManagementCertificates(d *schema.ResourceData) *[]apimanagement.CertificateConfiguration {
	certConfigs := d.Get("certificate").([]interface{})

	certificates := make([]apimanagement.CertificateConfiguration, 0)

	for _, v := range certConfigs {
		config := v.(map[string]interface{})

		cert_base64 := config["encoded_certificate"].(string)
		certificate_password := config["certificate_password"].(string)
		store_name := apimanagement.StoreName(config["store_name"].(string))

		cert := apimanagement.CertificateConfiguration{
			EncodedCertificate:  &cert_base64,
			CertificatePassword: &certificate_password,
			StoreName:           store_name,
		}

		certificates = append(certificates, cert)
	}

	return &certificates
}

func expandAzureRmApiManagementAdditionalLocations(d *schema.ResourceData, sku *apimanagement.ServiceSkuProperties) *[]apimanagement.AdditionalLocation {
	inputLocations := d.Get("additional_location").([]interface{})

	additionalLocations := make([]apimanagement.AdditionalLocation, 0)

	for _, v := range inputLocations {
		config := v.(map[string]interface{})
		location := config["location"].(string)

		additionalLocation := apimanagement.AdditionalLocation{
			Location: &location,
			Sku:      sku,
		}

		additionalLocations = append(additionalLocations, additionalLocation)
	}

	return &additionalLocations
}

func expandAzureRmApiManagementSku(configs []interface{}) *apimanagement.ServiceSkuProperties {
	config := configs[0].(map[string]interface{})

	name := apimanagement.SkuType(config["name"].(string))
	capacity := int32(config["capacity"].(int))

	sku := &apimanagement.ServiceSkuProperties{
		Name:     name,
		Capacity: &capacity,
	}

	return sku
}

func flattenAzureRmApiManagementMachineIdentity(identity *apimanagement.ServiceIdentity) []interface{} {
	if identity == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})
	result["type"] = *identity.Type

	if identity.PrincipalID != nil {
		result["principal_id"] = identity.PrincipalID.String()
	}
	if identity.TenantID != nil {
		result["tenant_id"] = identity.TenantID.String()
	}

	return []interface{}{result}
}

func flattenApiManagementCustomProperties(input map[string]*string) ([]interface{}, error) {
	output := make(map[string]interface{}, 0)

	if err := azure.SetCustomPropertyFrom(input, "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Ssl30", output, "disable_backend_ssl30"); err != nil {
		return nil, err
	}

	if err := azure.SetCustomPropertyFrom(input, "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls10", output, "disable_backend_tls10"); err != nil {
		return nil, err
	}

	if err := azure.SetCustomPropertyFrom(input, "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls11", output, "disable_backend_tls11"); err != nil {
		return nil, err
	}

	if err := azure.SetCustomPropertyFrom(input, "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Ssl30", output, "disable_frontend_ssl30"); err != nil {
		return nil, err
	}

	if err := azure.SetCustomPropertyFrom(input, "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls10", output, "disable_frontend_tls10"); err != nil {
		return nil, err
	}

	if err := azure.SetCustomPropertyFrom(input, "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls11", output, "disable_frontend_tls11"); err != nil {
		return nil, err
	}

	if err := azure.SetCustomPropertyFrom(input, "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TripleDes168", output, "disable_triple_des_chipers"); err != nil {
		return nil, err
	}

	if len(output) > 0 {
		return []interface{}{output}, nil
	}

	return nil, nil
}

func flattenApiManagementAdditionalLocations(props *[]apimanagement.AdditionalLocation) []interface{} {
	additional_locations := make([]interface{}, 0)

	if props != nil {
		for _, prop := range *props {
			additional_location := make(map[string]interface{}, 0)

			if prop.Location != nil {
				additional_location["location"] = *prop.Location
			}

			if prop.PublicIPAddresses != nil {
				additional_location["static_ips"] = *prop.PublicIPAddresses
			}

			if prop.GatewayRegionalURL != nil {
				additional_location["gateway_regional_url"] = *prop.GatewayRegionalURL
			}

			additional_locations = append(additional_locations, additional_location)
		}
	}

	return additional_locations
}

func flattenApiManagementHostnameConfigurations(d *schema.ResourceData, configs *[]apimanagement.HostnameConfiguration) []interface{} {
	if configs != nil && len(*configs) > 0 {
		hostTypes := make(map[string]interface{}) // protal, proxy etc.

		for _, hostNameType := range apimanagement.PossibleHostnameTypeValues() {
			v := strings.ToLower(string(hostNameType))
			hostTypes[v] = make([]interface{}, 0)
		}

		for _, config := range *configs {
			host_config := make(map[string]interface{}, 0)

			configType := strings.ToLower(string(config.Type))

			if config.HostName != nil {
				host_config["host_name"] = *config.HostName
			}

			// only set SSL binding for proxy types
			hostnameTypeProxy := strings.ToLower(string(apimanagement.Proxy))
			if configType == hostnameTypeProxy && config.DefaultSslBinding != nil {
				host_config["default_ssl_binding"] = *config.DefaultSslBinding
			}

			if config.NegotiateClientCertificate != nil {
				host_config["negotiate_client_certificate"] = *config.NegotiateClientCertificate
			}

			if config.KeyVaultID != nil {
				host_config["key_vault_id"] = *config.KeyVaultID
			}

			// Iterate through old state to find sensitive props not returned by API.
			// This must be done in order to avoid state diffs.
			key := fmt.Sprintf("hostname_configurations.0.%s", configType)
			if oldState, ok := d.GetOk(key); ok {
				for _, v := range oldState.([]interface{}) {
					oldConfig := v.(map[string]interface{})

					if oldConfig["host_name"] == *config.HostName {
						host_config["certificate_password"] = oldConfig["certificate_password"]
						host_config["certificate"] = oldConfig["certificate"]
					}
				}
			}

			hostTypes[configType] = append(hostTypes[configType].([]interface{}), host_config)
		}

		return []interface{}{hostTypes}
	}

	return nil
}

func flattenApiManagementServiceSku(profile *apimanagement.ServiceSkuProperties) []interface{} {
	sku := make(map[string]interface{}, 0)

	if profile != nil {
		if profile.Name != "" {
			sku["name"] = string(profile.Name)
		}

		if profile.Capacity != nil {
			sku["capacity"] = *profile.Capacity
		}
	}

	return []interface{}{sku}
}

func apiManagementResourceHostnameSchema(schemaName string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host_name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
		},

		"key_vault_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: azure.ValidateResourceID,
			ConflictsWith: []string{
				fmt.Sprintf("%s.0.certificate", schemaName),
				fmt.Sprintf("%s.0.certificate_password", schemaName),
			},
		},

		"certificate": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.NoZeroValues,
			ConflictsWith: []string{
				fmt.Sprintf("%s.0.key_vault_id", schemaName),
			},
		},

		"certificate_password": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.NoZeroValues,
			ConflictsWith: []string{
				fmt.Sprintf("%s.0.key_vault_id", schemaName),
			},
		},

		"negotiate_client_certificate": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func apiManagementResourceHostnameProxySchema(schemaName string) map[string]*schema.Schema {
	hostnameSchema := apiManagementResourceHostnameSchema(schemaName)

	hostnameSchema["default_ssl_binding"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Computed: true, // Azure has certain logic to set this, which we cannot predict
	}

	return hostnameSchema
}
