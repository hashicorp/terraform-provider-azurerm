// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apimanagementservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/tenantaccess"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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
						"location": commonschema.LocationComputed(),

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

			"tenant_access": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"primary_key": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"secondary_key": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func dataSourceApiManagementRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ServiceClient
	tenantAccessClient := meta.(*clients.Client).ApiManagement.TenantAccessClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := apimanagementservice.NewServiceID(meta.(*clients.Client).Account.SubscriptionId, resourceGroup, name)

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("checking for presence of an existing %s: %+v", id, err)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	if model := resp.Model; model != nil {
		d.Set("location", azure.NormalizeLocation(model.Location))

		identity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		d.Set("publisher_email", model.Properties.PublisherEmail)
		d.Set("publisher_name", model.Properties.PublisherName)
		d.Set("notification_sender_email", pointer.From(model.Properties.NotificationSenderEmail))
		d.Set("gateway_url", pointer.From(model.Properties.GatewayURL))
		d.Set("gateway_regional_url", pointer.From(model.Properties.GatewayRegionalURL))
		d.Set("portal_url", pointer.From(model.Properties.PortalURL))
		d.Set("developer_portal_url", pointer.From(model.Properties.DeveloperPortalURL))
		d.Set("management_api_url", pointer.From(model.Properties.ManagementApiURL))
		d.Set("scm_url", pointer.From(model.Properties.ScmURL))
		d.Set("public_ip_addresses", pointer.From(model.Properties.PublicIPAddresses))
		d.Set("public_ip_address_id", pointer.From(model.Properties.PublicIPAddressId))
		d.Set("private_ip_addresses", pointer.From(model.Properties.PrivateIPAddresses))

		if err := d.Set("hostname_configuration", flattenDataSourceApiManagementHostnameConfigurations(model.Properties.HostnameConfigurations)); err != nil {
			return fmt.Errorf("setting `hostname_configuration`: %+v", err)
		}

		if err := d.Set("additional_location", flattenDataSourceApiManagementAdditionalLocations(model.Properties.AdditionalLocations)); err != nil {
			return fmt.Errorf("setting `additional_location`: %+v", err)
		}

		d.Set("sku_name", flattenApiManagementServiceSkuName(&model.Sku))

		tenantAccess := make([]interface{}, 0)
		if model.Sku.Name != apimanagementservice.SkuTypeConsumption {
			tenantAccessServiceId := tenantaccess.NewAccessID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, "access")
			tenantAccessInformationContract, err := tenantAccessClient.ListSecrets(ctx, tenantAccessServiceId)
			if err != nil {
				if !response.WasForbidden(tenantAccessInformationContract.HttpResponse) {
					return fmt.Errorf("retrieving tenant access properties for %s: %+v", id, err)
				}
			} else {
				tenantAccess = flattenApiManagementTenantAccessSettings(*tenantAccessInformationContract.Model)
			}
		}
		if err := d.Set("tenant_access", tenantAccess); err != nil {
			return fmt.Errorf("setting `tenant_access`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func flattenDataSourceApiManagementHostnameConfigurations(input *[]apimanagementservice.HostnameConfiguration) []interface{} {
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

		output["host_name"] = config.HostName

		output["negotiate_client_certificate"] = pointer.From(config.NegotiateClientCertificate)

		output["key_vault_certificate_id"] = pointer.From(config.KeyVaultId)
		if !features.FivePointOh() {
			output["key_vault_id"] = pointer.From(config.KeyVaultId)
		}

		switch strings.ToLower(string(config.Type)) {
		case strings.ToLower(string(apimanagementservice.HostnameTypeProxy)):
			// only set SSL binding for proxy types
			if config.DefaultSslBinding != nil {
				output["default_ssl_binding"] = *config.DefaultSslBinding
			}
			proxyResults = append(proxyResults, output)

		case strings.ToLower(string(apimanagementservice.HostnameTypeManagement)):
			managementResults = append(managementResults, output)

		case strings.ToLower(string(apimanagementservice.HostnameTypePortal)):
			portalResults = append(portalResults, output)

		case strings.ToLower(string(apimanagementservice.HostnameTypeDeveloperPortal)):
			developerPortalResults = append(developerPortalResults, output)

		case strings.ToLower(string(apimanagementservice.HostnameTypeScm)):
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

func flattenDataSourceApiManagementAdditionalLocations(input *[]apimanagementservice.AdditionalLocation) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, prop := range *input {
		results = append(results, map[string]interface{}{
			"capacity":             int32(prop.Sku.Capacity),
			"gateway_regional_url": pointer.From(prop.GatewayRegionalURL),
			"location":             location.NormalizeNilable(pointer.To(prop.Location)),
			"private_ip_addresses": pointer.From(prop.PrivateIPAddresses),
			"public_ip_address_id": pointer.From(prop.PublicIPAddressId),
			"public_ip_addresses":  pointer.From(prop.PublicIPAddresses),
			"zones":                zones.FlattenUntyped(prop.Zones),
		})
	}

	return results
}

func apiManagementDataSourceHostnameSchema() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"host_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"key_vault_certificate_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"negotiate_client_certificate": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
	}

	if !features.FivePointOh() {
		s["key_vault_id"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Computed: true,
		}
	}

	return s
}

func apiManagementDataSourceHostnameProxySchema() map[string]*pluginsdk.Schema {
	hostnameSchema := apiManagementDataSourceHostnameSchema()

	hostnameSchema["default_ssl_binding"] = &pluginsdk.Schema{
		Type:     pluginsdk.TypeBool,
		Computed: true,
	}

	return hostnameSchema
}
