// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management

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
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/tenantaccess"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing API Management Service.",
		Read:        dataSourceApiManagementRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": schemaz.SchemaApiManagementDataSourceName(),

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"public_ip_addresses": {
				Description: "The Public IP addresses of the API Management Service.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"public_ip_address_id": {
				Description: "ID of the standard SKU IPv4 Public IP. Available only for Premium and Developer SKU deployed in a Virtual Network.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"private_ip_addresses": {
				Description: "The Private IP addresses of the API Management Service.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"publisher_name": {
				Description: "The name of the Publisher/Company of the API Management Service.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"publisher_email": {
				Description: "The email of Publisher/Company of the API Management Service.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"sku_name": {
				Description: "The SKU of the API Management Service consisting of the name and capacity, separated by an underscore.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

			"notification_sender_email": {
				Description: "The email address from which the notification will be sent.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"gateway_url": {
				Description: "The URL for the API Management Service's Gateway.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"gateway_regional_url": {
				Description: "The URL for the Gateway in the Default Region.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"portal_url": {
				Description: "The URL of the Publisher Portal.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"developer_portal_url": {
				Description: "The URL for the Developer Portal associated with this API Management service.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"management_api_url": {
				Description: "The URL for the Management API.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"scm_url": {
				Description: "The SCM (Source Code Management) endpoint.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"additional_location": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"location": commonschema.LocationComputed(),

						"capacity": {
							Description: "The number of compute units in this region.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
						},

						"zones": commonschema.ZonesMultipleComputed(),

						"gateway_regional_url": {
							Description: "Gateway URL of the API Management service in the Region.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"public_ip_addresses": {
							Description: "Public Static Load Balanced IP addresses of the API Management service in the additional location. Available only for Basic, Standard and Premium SKU.",
							Type:        pluginsdk.TypeList,
							Computed:    true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"public_ip_address_id": {
							Description: "ID of the standard SKU IPv4 Public IP. Available only for Premium SKU deployed in a Virtual Network.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"private_ip_addresses": {
							Description: "Private IP addresses of the API Management service in the additional location, for instances using Virtual Network mode.",
							Type:        pluginsdk.TypeList,
							Computed:    true,
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
							Description: "Is access to the Management API enabled (presumably \"for this Tenant\")?",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},
						"tenant_id": {
							Description: "The ID of the Tenant which has access to this API Management instance.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"primary_key": {
							Description: "Primary access key for the tenant access information contract.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
							Sensitive:   true,
						},
						"secondary_key": {
							Description: "Secondary access key for the tenant access information contract.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
							Sensitive:   true,
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func dataSourceApiManagementRead(d *pluginsdk.ResourceData, meta any) error {
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
		d.Set("location", location.Normalize(model.Location))

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

		tenantAccess := make([]any, 0)
		if model.Sku.Name != apimanagementservice.SkuTypeConsumption && !strings.Contains(string(model.Sku.Name), "V2") {
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

func flattenDataSourceApiManagementHostnameConfigurations(input *[]apimanagementservice.HostnameConfiguration) []any {
	if input == nil {
		return []any{}
	}

	// management, portal, proxy, scm
	managementResults := make([]any, 0)
	proxyResults := make([]any, 0)
	portalResults := make([]any, 0)
	developerPortalResults := make([]any, 0)
	scmResults := make([]any, 0)

	for _, config := range *input {
		output := make(map[string]any)

		output["host_name"] = config.HostName

		output["negotiate_client_certificate"] = pointer.From(config.NegotiateClientCertificate)

		output["key_vault_certificate_id"] = pointer.From(config.KeyVaultId)

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

	return []any{
		map[string]any{
			"management":       managementResults,
			"portal":           portalResults,
			"developer_portal": developerPortalResults,
			"proxy":            proxyResults,
			"scm":              scmResults,
		},
	}
}

func flattenDataSourceApiManagementAdditionalLocations(input *[]apimanagementservice.AdditionalLocation) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, prop := range *input {
		results = append(results, map[string]any{
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
	return map[string]*pluginsdk.Schema{
		"host_name": {
			Description: "The hostname.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"key_vault_certificate_id": {
			Description: "The ID of the Key Vault Secret which contains the SSL Certificate.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"negotiate_client_certificate": {
			Description: "Is Client Certificate Negotiation enabled?",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
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
