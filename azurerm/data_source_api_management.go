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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"location": locationForDataSourceSchema(),

			"kind": {
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

			"properties": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"publisher_email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"publisher_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notification_sender_email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provisioning_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_provisioning_state": {
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
						"hostname_configurations": {
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
									"key_vault_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"encoded_certificate": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"certificate_password": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"default_ssl_binding": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"negotiate_client_certificate": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"certificate": {
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
									},
								},
							},
						},
						"static_ips": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed: true,
						},
						"virtual_network_configuration": {
							Type:     schema.TypeString,
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
	// d.Set("kind", resp.Kind)

	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.ServiceProperties; props != nil {
		d.Set("properties", flattenApiManagementServiceProperties(props))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", flattenApiManagementServiceSku(sku))
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func flattenApiManagementServiceProperties(props *apimanagement.ServiceProperties) []interface{} {
	result := make([]interface{}, 0, 1)
	properties := make(map[string]interface{}, 0)

	if props.PublisherEmail != nil {
		properties["publisher_email"] = *props.PublisherEmail
	}

	if props.PublisherName != nil {
		properties["publisher_name"] = *props.PublisherName
	}

	if props.NotificationSenderEmail != nil {
		properties["notification_sender_email"] = *props.NotificationSenderEmail
	}

	if props.ProvisioningState != nil {
		properties["provisioning_state"] = *props.ProvisioningState
	}

	if props.TargetProvisioningState != nil {
		properties["target_provisioning_state"] = *props.TargetProvisioningState
	}

	if props.CreatedAtUtc != nil {
		properties["created"] = props.CreatedAtUtc.Format(time.RFC3339)
	}

	if props.GatewayURL != nil {
		properties["gateway_url"] = *props.GatewayURL
	}

	if props.GatewayRegionalURL != nil {
		properties["gateway_regional_url"] = *props.GatewayRegionalURL
	}

	if props.PortalURL != nil {
		properties["portal_url"] = *props.PortalURL
	}

	if props.ManagementAPIURL != nil {
		properties["management_api_url"] = *props.ManagementAPIURL
	}

	if props.ScmURL != nil {
		properties["scm_url"] = *props.ScmURL
	}

	if props.HostnameConfigurations != nil {
		properties["hostname_configurations"] = flattenApiHostnameConfigurations(props.HostnameConfigurations)
	}

	result = append(result, properties)
	return result
}

func flattenApiHostnameConfigurations(configs *[]apimanagement.HostnameConfiguration) []interface{} {
	host_configs := make([]interface{}, 0, 1)
	host_config := make(map[string]interface{}, 2)

	for _, config := range *configs {
		host_config["type"] = string(config.Type)

		if config.HostName != nil {
			host_config["host_name"] = *config.HostName
		}

		if config.KeyVaultID != nil {
			host_config["key_vault_id"] = *config.KeyVaultID
		}

		if config.EncodedCertificate != nil {
			host_config["encoded_certificate"] = *config.EncodedCertificate
		}

		if config.CertificatePassword != nil {
			host_config["certificate_password"] = *config.CertificatePassword
		}

		if config.DefaultSslBinding != nil {
			host_config["default_ssl_binding"] = *config.DefaultSslBinding
		}

		if config.NegotiateClientCertificate != nil {
			host_config["negotiate_client_certificate"] = *config.NegotiateClientCertificate
		}

		if config.Certificate != nil {
			host_config["certificate"] = flattenApiManagementCertificate(config.Certificate)
		}

		host_configs = append(host_configs, host_config)
	}

	return host_configs
}

func flattenApiManagementCertificate(cert *apimanagement.CertificateInformation) interface{} {
	certificate := make(map[string]interface{}, 2)

	if cert.Expiry != nil {
		certificate["expiry"] = cert.Expiry.Format(time.RFC3339)
	}

	if cert.Thumbprint != nil {
		certificate["thumbprint"] = cert.Thumbprint
	}

	if cert.Subject != nil {
		certificate["subject"] = cert.Subject
	}

	return certificate
}

func flattenApiManagementServiceSku(profile *apimanagement.ServiceSkuProperties) []interface{} {
	skus := make([]interface{}, 0, 1)
	sku := make(map[string]interface{}, 2)

	sku["name"] = string(profile.Name)

	if profile.Capacity != nil {
		sku["capacity"] = *profile.Capacity
	}

	skus = append(skus, sku)

	return skus
}
