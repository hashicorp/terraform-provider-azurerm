package apimanagement

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceApiManagementService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceApiManagementRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": schemaz.SchemaApiManagementDataSourceName(),

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"public_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"public_ip_address_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"private_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"publisher_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"publisher_email": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

			"notification_sender_email": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"gateway_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"gateway_regional_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"portal_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"developer_portal_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"management_api_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"scm_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"additional_location": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"location": azure.SchemaLocationForDataSource(),

						"capacity": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"zones": commonschema.ZonesMultipleComputed(),

						"gateway_regional_url": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"public_ip_addresses": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"public_ip_address_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"private_ip_addresses": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"hostname_configuration": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"management": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: apiManagementDataSourceHostnameSchema(),
							},
						},
						"portal": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: apiManagementDataSourceHostnameSchema(),
							},
						},
						"developer_portal": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: apiManagementDataSourceHostnameSchema(),
							},
						},
						"proxy": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: apiManagementDataSourceHostnameProxySchema(),
							},
						},
						"scm": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
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

func dataSourceApiManagementRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("parsing API Management ID: %q", *resp.ID)
	}

	d.SetId(id.ID())

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	identity, err := flattenIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
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
		d.Set("public_ip_address_id", props.PublicIPAddressID)
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
		var capacity *int32
		if prop.Sku.Capacity != nil {
			capacity = prop.Sku.Capacity
		}

		var publicIpAddresses []string
		if prop.PublicIPAddresses != nil {
			publicIpAddresses = *prop.PublicIPAddresses
		}

		publicIpAddressId := ""
		if prop.PublicIPAddressID != nil {
			publicIpAddressId = *prop.PublicIPAddressID
		}

		var privateIpAddresses []string
		if prop.PrivateIPAddresses != nil {
			privateIpAddresses = *prop.PrivateIPAddresses
		}

		gatewayRegionalUrl := ""
		if prop.GatewayRegionalURL != nil {
			gatewayRegionalUrl = *prop.GatewayRegionalURL
		}

		results = append(results, map[string]interface{}{
			"capacity":             capacity,
			"gateway_regional_url": gatewayRegionalUrl,
			"location":             location.NormalizeNilable(prop.Location),
			"private_ip_addresses": privateIpAddresses,
			"public_ip_address_id": publicIpAddressId,
			"public_ip_addresses":  publicIpAddresses,
			"zones":                zones.Flatten(prop.Zones),
		})
	}

	return results
}

func apiManagementDataSourceHostnameSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"host_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"key_vault_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"negotiate_client_certificate": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
	}
}

func apiManagementDataSourceHostnameProxySchema() map[string]*pluginsdk.Schema {
	hostnameSchema := apiManagementDataSourceHostnameSchema()

	hostnameSchema["default_ssl_binding"] = &pluginsdk.Schema{
		Type:     pluginsdk.TypeBool,
		Computed: true,
	}

	return hostnameSchema
}
