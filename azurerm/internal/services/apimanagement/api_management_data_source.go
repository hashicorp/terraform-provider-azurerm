package apimanagement

import (
	"fmt"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"
	msiparse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/parse"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceApiManagementService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApiManagementRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": schemaz.SchemaApiManagementDataSourceName(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"public_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"private_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"publisher_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"publisher_email": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sku_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"identity": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"identity_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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

			"developer_portal_url": {
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
						"location": azure.SchemaLocationForDataSource(),

						"gateway_regional_url": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"public_ip_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"private_ip_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"hostname_configuration": {
				Type:     schema.TypeList,
				Computed: true,
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
						"developer_portal": {
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

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceApiManagementRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("API Management Service %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("retrieving API Management Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	id, err := parse.ApiManagementID(*resp.ID)
	if err != nil {
		return fmt.Errorf("Error parsing API Management ID: %q", *resp.ID)
	}

	d.SetId(id.ID())

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	identity, err := flattenApiManagementDataSourceIdentity(resp.Identity)
	if err != nil {
		return err
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if props := resp.ServiceProperties; props != nil {
		d.Set("publisher_email", props.PublisherEmail)
		d.Set("publisher_name", props.PublisherName)
		d.Set("notification_sender_email", props.NotificationSenderEmail)
		d.Set("gateway_url", props.GatewayURL)
		d.Set("gateway_regional_url", props.GatewayRegionalURL)
		d.Set("portal_url", props.PortalURL)
		d.Set("developer_portal_url", props.DeveloperPortalURL)
		d.Set("management_api_url", props.ManagementAPIURL)
		d.Set("scm_url", props.ScmURL)
		d.Set("public_ip_addresses", props.PublicIPAddresses)
		d.Set("private_ip_addresses", props.PrivateIPAddresses)

		if err := d.Set("hostname_configuration", flattenDataSourceApiManagementHostnameConfigurations(props.HostnameConfigurations)); err != nil {
			return fmt.Errorf("setting `hostname_configuration`: %+v", err)
		}

		if err := d.Set("additional_location", flattenDataSourceApiManagementAdditionalLocations(props.AdditionalLocations)); err != nil {
			return fmt.Errorf("setting `additional_location`: %+v", err)
		}
	}

	d.Set("sku_name", flattenApiManagementServiceSkuName(resp.Sku))

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenDataSourceApiManagementHostnameConfigurations(input *[]apimanagement.HostnameConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	// management, portal, proxy, scm
	managementResults := make([]interface{}, 0)
	proxyResults := make([]interface{}, 0)
	portalResults := make([]interface{}, 0)
	developerPortalResults := make([]interface{}, 0)
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

		switch strings.ToLower(string(config.Type)) {
		case strings.ToLower(string(apimanagement.HostnameTypeProxy)):
			// only set SSL binding for proxy types
			if config.DefaultSslBinding != nil {
				output["default_ssl_binding"] = *config.DefaultSslBinding
			}
			proxyResults = append(proxyResults, output)

		case strings.ToLower(string(apimanagement.HostnameTypeManagement)):
			managementResults = append(managementResults, output)

		case strings.ToLower(string(apimanagement.HostnameTypePortal)):
			portalResults = append(portalResults, output)

		case strings.ToLower(string(apimanagement.HostnameTypeDeveloperPortal)):
			developerPortalResults = append(developerPortalResults, output)

		case strings.ToLower(string(apimanagement.HostnameTypeScm)):
			scmResults = append(scmResults, output)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"management":       managementResults,
			"portal":           portalResults,
			"developer_portal": developerPortalResults,
			"proxy":            proxyResults,
			"scm":              scmResults,
		},
	}
}

func flattenDataSourceApiManagementAdditionalLocations(input *[]apimanagement.AdditionalLocation) []interface{} {
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

		if prop.PrivateIPAddresses != nil {
			output["private_ip_addresses"] = *prop.PrivateIPAddresses
		}

		if prop.GatewayRegionalURL != nil {
			output["gateway_regional_url"] = *prop.GatewayRegionalURL
		}

		results = append(results, output)
	}

	return results
}

func flattenApiManagementDataSourceIdentity(identity *apimanagement.ServiceIdentity) ([]interface{}, error) {
	if identity == nil || identity.Type == apimanagement.None {
		return make([]interface{}, 0), nil
	}

	result := make(map[string]interface{})
	result["type"] = string(identity.Type)

	if identity.PrincipalID != nil {
		result["principal_id"] = identity.PrincipalID.String()
	}

	if identity.TenantID != nil {
		result["tenant_id"] = identity.TenantID.String()
	}

	identityIds := make([]interface{}, 0)
	if identity.UserAssignedIdentities != nil {
		for key := range identity.UserAssignedIdentities {
			parsedId, err := msiparse.UserAssignedIdentityID(key)
			if err != nil {
				return nil, err
			}
			identityIds = append(identityIds, parsedId.ID())
		}
		result["identity_ids"] = identityIds
	}

	return []interface{}{result}, nil
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

		"negotiate_client_certificate": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}

func apiManagementDataSourceHostnameProxySchema() map[string]*schema.Schema {
	hostnameSchema := apiManagementDataSourceHostnameSchema()

	hostnameSchema["default_ssl_binding"] = &schema.Schema{
		Type:     schema.TypeBool,
		Computed: true,
	}

	return hostnameSchema
}
