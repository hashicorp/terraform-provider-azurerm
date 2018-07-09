package azurerm

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2017-03-01/apimanagement"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceApiManagementService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApiManagementRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"location": locationForDataSourceSchema(),

			"publisher_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"publisher_email": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sku": apiManagementDataSourceSkuSchema(),

			"notification_sender_email": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

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

			"additional_location": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location": locationForDataSourceSchema(),

						"sku": apiManagementDataSourceSkuSchema(),

						"gateway_regional_url": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"static_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"certificate": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"store_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"certificate_info": apiManagementDataSourceCertificateInfoSchema(),
					},
				},
			},

			"custom_properties": {
				Type:     schema.TypeMap,
				Computed: true,
			},

			"hostname_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"certificate_info": apiManagementDataSourceCertificateInfoSchema(),

						"default_ssl_binding": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"negotiate_client_certificate": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceApiManagementRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagementServiceClient

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	ctx := meta.(*ArmClient).StopContext
	resp, err := client.Get(ctx, resourceGroup, name)

	if err != nil {
		return fmt.Errorf("Error making Read request on API Management Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if utils.ResponseWasNotFound(resp.Response) {
		return fmt.Errorf("Error: API Management Service %q (Resource Group %q) was not found", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.ServiceProperties; props != nil {
		d.Set("publisher_email", props.PublisherEmail)
		d.Set("publisher_name", props.PublisherName)

		d.Set("notification_sender_email", props.NotificationSenderEmail)
		d.Set("created", props.CreatedAtUtc.Format(time.RFC3339))
		d.Set("gateway_url", props.GatewayURL)
		d.Set("gateway_regional_url", props.GatewayRegionalURL)
		d.Set("portal_url", props.PortalURL)
		d.Set("management_api_url", props.ManagementAPIURL)
		d.Set("scm_url", props.ScmURL)
		d.Set("static_ips", props.StaticIps)
		d.Set("custom_properties", &props.CustomProperties)

		hostnameConfigurations, err := flattenDataSourceApiManagementHostnameConfigurations(d, props.HostnameConfigurations)
		if err != nil {
			return err
		}

		if err := d.Set("hostname_configuration", hostnameConfigurations); err != nil {
			return fmt.Errorf("Error setting `hostname_configuration`: %+v", err)
		}

		additionalLocations := flattenDataSourceApiManagementAdditionalLocations(props.AdditionalLocations)
		if err := d.Set("additional_location", additionalLocations); err != nil {
			return fmt.Errorf("Error setting `additional_location`: %+v", err)
		}

		certificates, err := flattenDataSourceApiManagementCertificates(d, props.Certificates)

		if err != nil {
			return err
		}

		if err := d.Set("certificate", certificates); err != nil {
			return fmt.Errorf("Error setting `certificate`: %+v", err)
		}
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", flattenDataSourceApiManagementServiceSku(sku))
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func apiManagementDataSourceSkuSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"capacity": {
					Type:     schema.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func apiManagementDataSourceCertificateInfoSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"expiry": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"thumbprint": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"subject": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func flattenDataSourceApiManagementHostnameConfigurations(d *schema.ResourceData, configs *[]apimanagement.HostnameConfiguration) ([]interface{}, error) {
	host_configs := make([]interface{}, 0, 1)

	if configs != nil {
		for i, config := range *configs {
			host_config := make(map[string]interface{}, 2)

			if config.Type != "" {
				host_config["type"] = string(config.Type)
			}

			if config.HostName != nil {
				host_config["host_name"] = *config.HostName
			}

			if config.DefaultSslBinding != nil {
				host_config["default_ssl_binding"] = *config.DefaultSslBinding
			}

			if config.NegotiateClientCertificate != nil {
				host_config["negotiate_client_certificate"] = bool(*config.NegotiateClientCertificate)
			}

			if config.Certificate != nil {
				host_config["certificate_info"] = flattenDataSourceApiManagementCertificate(config.Certificate)
			}

			// certificate password isn't returned, so let's look it up
			passKey := fmt.Sprintf("hostname_configuration.%d.certificate_password", i)
			if v, ok := d.GetOk(passKey); ok {
				password := v.(string)
				host_config["certificate_password"] = password
			} else {
				return nil, fmt.Errorf("Error getting `certificate_password` from key %s", passKey)
			}

			// encoded certificate isn't returned, so let's look it up
			certKey := fmt.Sprintf("hostname_configuration.%d.certificate", i)
			if v, ok := d.GetOk(certKey); ok {
				cert := v.(string)
				host_config["certificate"] = cert
			} else {
				return nil, fmt.Errorf("Error getting `certificate` from key %s", certKey)
			}

			host_configs = append(host_configs, host_config)
		}
	}

	return host_configs, nil
}

func flattenDataSourceApiManagementCertificate(cert *apimanagement.CertificateInformation) []interface{} {
	certificate := make(map[string]interface{}, 2)
	certInfos := make([]interface{}, 0, 1)

	if cert != nil {
		if cert.Expiry != nil {
			certificate["expiry"] = cert.Expiry.Format(time.RFC3339)
		}

		certificate["thumbprint"] = *cert.Thumbprint
		certificate["subject"] = *cert.Subject

		certInfos = append(certInfos, certificate)
	}

	return certInfos
}

func flattenDataSourceApiManagementAdditionalLocations(props *[]apimanagement.AdditionalLocation) []interface{} {
	additional_locations := make([]interface{}, 0, 1)

	if props != nil {
		for _, prop := range *props {
			additional_location := make(map[string]interface{}, 2)

			if prop.Location != nil {
				additional_location["location"] = *prop.Location
			}

			if prop.StaticIps != nil {
				additional_location["static_ips"] = *prop.StaticIps
			}

			if prop.GatewayRegionalURL != nil {
				additional_location["gateway_regional_url"] = *prop.GatewayRegionalURL
			}

			if prop.Sku != nil {
				if sku := flattenDataSourceApiManagementServiceSku(prop.Sku); sku != nil {
					additional_location["sku"] = sku
				}
			}

			additional_locations = append(additional_locations, additional_location)
		}
	}

	return additional_locations
}

func flattenDataSourceApiManagementCertificates(d *schema.ResourceData, props *[]apimanagement.CertificateConfiguration) ([]interface{}, error) {
	certificates := make([]interface{}, 0, 1)

	if props != nil {
		for i, prop := range *props {
			certificate := make(map[string]interface{}, 2)

			if prop.StoreName != "" {
				certificate["store_name"] = string(prop.StoreName)
			}

			if cert := flattenDataSourceApiManagementCertificate(prop.Certificate); cert != nil {
				certificate["certificate_info"] = cert
			}

			// certificate password isn't returned, so let's look it up
			passwKey := fmt.Sprintf("certificate.%d.certificate_password", i)
			if v, ok := d.GetOk(passwKey); ok {
				password := v.(string)
				certificate["certificate_password"] = password
			} else {
				return nil, fmt.Errorf("Error getting `certificate_password` from key %s", passwKey)
			}
			// encoded certificate isn't returned, so let's look it up
			certKey := fmt.Sprintf("certificate.%d.encoded_certificate", i)
			if v, ok := d.GetOk(certKey); ok {
				cert := v.(string)
				certificate["encoded_certificate"] = cert
			} else {
				return nil, fmt.Errorf("Error getting `encoded_certificate` from key %s", certKey)
			}

			certificates = append(certificates, certificate)
		}
	}

	return certificates, nil
}

func flattenDataSourceApiManagementServiceSku(profile *apimanagement.ServiceSkuProperties) []interface{} {
	skus := make([]interface{}, 0, 1)
	sku := make(map[string]interface{}, 2)

	if profile != nil {
		sku["name"] = string(profile.Name)
		sku["capacity"] = *profile.Capacity
	}

	skus = append(skus, sku)

	return skus
}
