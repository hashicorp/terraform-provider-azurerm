package azurerm

import (
	"fmt"
	"strings"

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

			"hostname_configurations": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"management": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: apiManagementDataSourceHostnameSchema(),
							},
						},
						"portal": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: apiManagementDataSourceHostnameSchema(),
							},
						},
						"proxy": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: apiManagementDataSourceHostnameProxySchema(),
							},
						},
						"scm": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: apiManagementDataSourceHostnameSchema(),
							},
						},
					},
				},
			},

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func apiManagementDataSourceHostnameSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host_name": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"key_vault_id": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"certificate": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"certificate_password": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"negotiate_client_certificate": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}

func apiManagementDataSourceHostnameProxySchema() map[string]*schema.Schema {
	hostnameSchema := apiManagementResourceHostnameSchema()

	hostnameSchema["default_ssl_binding"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Computed: true, //Azure has certain logic to set this, which we cannot predict
	}

	return hostnameSchema
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

		if err := d.Set("hostname_configurations", flattenDataSourceApiManagementHostnameConfigurations(d, props.HostnameConfigurations)); err != nil {
			return fmt.Errorf("Error setting `hostname_configurations`: %+v", err)
		}

		if err := d.Set("additional_location", flattenDataSourceApiManagementAdditionalLocations(props.AdditionalLocations)); err != nil {
			return fmt.Errorf("Error setting `additional_location`: %+v", err)
		}
	}

	if sku := resp.Sku; sku != nil {
		if err := d.Set("sku", flattenDataSourceApiManagementServiceSku(sku)); err != nil {
			return fmt.Errorf("Error flattening `sku`: %+v", err)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func flattenDataSourceApiManagementHostnameConfigurations(d *schema.ResourceData, configs *[]apimanagement.HostnameConfiguration) []interface{} {
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

			hostTypes[configType] = append(hostTypes[configType].([]interface{}), host_config)
		}

		return []interface{}{hostTypes}
	}

	return nil
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
