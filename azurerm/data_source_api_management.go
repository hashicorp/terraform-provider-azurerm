package azurerm

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/apimanagement/mgmt/2018-06-01-preview/apimanagement"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceApiManagementService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApiManagementRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateApiManagementName,
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

			"sku": {
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
			},

			"notification_sender_email": {
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

			"security": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disable_backend_ssl30": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"disable_backend_tls10": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"disable_backend_tls11": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"disable_triple_des_chipers": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"disable_frontend_ssl30": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"disable_frontend_tls10": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"disable_frontend_tls11": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
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

		if err := d.Set("hostname_configuration", flattenApiManagementHostnameConfigurations(d, props.HostnameConfigurations)); err != nil {
			return fmt.Errorf("Error setting `hostname_configuration`: %+v", err)
		}

		if err := d.Set("additional_location", flattenApiManagementAdditionalLocations(props.AdditionalLocations)); err != nil {
			return fmt.Errorf("Error setting `additional_location`: %+v", err)
		}

		if err := d.Set("certificate", flattenApiManagementCertificates(d, props.Certificates)); err != nil {
			return fmt.Errorf("Error setting `certificate`: %+v", err)
		}
	}

	if sku := resp.Sku; sku != nil {
		if err := d.Set("sku", flattenApiManagementServiceSku(sku)); err != nil {
			return fmt.Errorf("Error flattening `sku`: %+v", err)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
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
	host_configs := make([]interface{}, 0)

	if configs != nil {
		for _, config := range *configs {
			host_config := make(map[string]interface{}, 2)

			host_config["type"] = string(config.Type)

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

			host_configs = append(host_configs, host_config)
		}
	}

	return host_configs, nil
}

func flattenDataSourceApiManagementCertificate(cert *apimanagement.CertificateInformation) []interface{} {
	certificate := make(map[string]interface{}, 2)

	if cert != nil {
		if cert.Expiry != nil {
			certificate["expiry"] = cert.Expiry.Format(time.RFC3339)
		}

		if cert.Thumbprint != nil {
			certificate["thumbprint"] = *cert.Thumbprint
		}

		if cert.Subject != nil {
			certificate["subject"] = *cert.Subject
		}
	}

	return []interface{}{certificate}
}

func flattenDataSourceApiManagementAdditionalLocations(props *[]apimanagement.AdditionalLocation) []interface{} {
	additional_locations := make([]interface{}, 0)

	if props != nil {
		for _, prop := range *props {
			additional_location := make(map[string]interface{}, 2)

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

func flattenDataSourceApiManagementCertificates(d *schema.ResourceData, props *[]apimanagement.CertificateConfiguration) ([]interface{}, error) {
	certificates := make([]interface{}, 0)

	if props != nil {
		for _, prop := range *props {
			certificate := make(map[string]interface{}, 2)

			certificate["store_name"] = string(prop.StoreName)

			if cert := flattenDataSourceApiManagementCertificate(prop.Certificate); cert != nil {
				certificate["certificate_info"] = cert
			}

			certificates = append(certificates, certificate)
		}
	}

	return certificates, nil
}

func flattenDataSourceApiManagementServiceSku(profile *apimanagement.ServiceSkuProperties) []interface{} {
	sku := make(map[string]interface{}, 2)

	if profile != nil {
		sku["name"] = string(profile.Name)

		if profile.Capacity != nil {
			sku["capacity"] = *profile.Capacity
		}
	}

	return []interface{}{sku}
}
