// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceVPNServerConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceVPNServerConfigurationRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"vpn_authentication_types": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"azure_active_directory_authentication": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"audience": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"issuer": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"tenant": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"client_revoked_certificate": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"thumbprint": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"client_root_certificate": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"public_cert_data": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"ipsec_policy": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"dh_group": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"ike_encryption": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"ike_integrity": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"ipsec_encryption": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"ipsec_integrity": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"pfs_group": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"sa_lifetime_seconds": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"sa_data_size_kilobytes": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"radius": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"server": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"address": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"secret": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"score": {
										Type:     pluginsdk.TypeInt,
										Computed: true,
									},
								},
							},
						},

						"client_root_certificate": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"thumbprint": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},

						"server_root_certificate": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"public_cert_data": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"vpn_protocols": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceVPNServerConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	defer cancel()

	id := virtualwans.NewVpnServerConfigurationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.VpnServerConfigurationsGet(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.VpnServerConfigurationName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			flattenedAADAuthentication := dataSourceFlattenVpnServerConfigurationAADAuthentication(props.AadAuthenticationParameters)
			if err := d.Set("azure_active_directory_authentication", flattenedAADAuthentication); err != nil {
				return fmt.Errorf("setting `azure_active_directory_authentication`: %+v", err)
			}

			flattenedClientRootCerts := dataSourceFlattenVpnServerConfigurationClientRootCertificates(props.VpnClientRootCertificates)
			if err := d.Set("client_root_certificate", flattenedClientRootCerts); err != nil {
				return fmt.Errorf("setting `client_root_certificate`: %+v", err)
			}

			flattenedClientRevokedCerts := dataSourceFlattenVpnServerConfigurationClientRevokedCertificates(props.VpnClientRevokedCertificates)
			if err := d.Set("client_revoked_certificate", flattenedClientRevokedCerts); err != nil {
				return fmt.Errorf("setting `client_revoked_certificate`: %+v", err)
			}

			flattenedIPSecPolicies := dataSourceFlattenVpnServerConfigurationIPSecPolicies(props.VpnClientIPsecPolicies)
			if err := d.Set("ipsec_policy", flattenedIPSecPolicies); err != nil {
				return fmt.Errorf("setting `ipsec_policy`: %+v", err)
			}

			flattenedRadius := dataSourceFlattenVpnServerConfigurationRadius(props)
			if err := d.Set("radius", flattenedRadius); err != nil {
				return fmt.Errorf("setting `radius`: %+v", err)
			}

			vpnAuthenticationTypes := make([]interface{}, 0)
			if props.VpnAuthenticationTypes != nil {
				for _, v := range *props.VpnAuthenticationTypes {
					vpnAuthenticationTypes = append(vpnAuthenticationTypes, string(v))
				}
			}
			if err := d.Set("vpn_authentication_types", vpnAuthenticationTypes); err != nil {
				return fmt.Errorf("setting `vpn_authentication_types`: %+v", err)
			}

			flattenedVpnProtocols := dataSourceFlattenVpnServerConfigurationVPNProtocols(props.VpnProtocols)
			if err := d.Set("vpn_protocols", pluginsdk.NewSet(pluginsdk.HashString, flattenedVpnProtocols)); err != nil {
				return fmt.Errorf("setting `vpn_protocols`: %+v", err)
			}

			if err := tags.FlattenAndSet(d, model.Tags); err != nil {
				return err
			}
		}
	}

	return nil
}

func dataSourceFlattenVpnServerConfigurationAADAuthentication(input *virtualwans.AadAuthenticationParameters) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	audience := ""
	if input.AadAudience != nil {
		audience = *input.AadAudience
	}

	issuer := ""
	if input.AadIssuer != nil {
		issuer = *input.AadIssuer
	}

	tenant := ""
	if input.AadTenant != nil {
		tenant = *input.AadTenant
	}

	return []interface{}{
		map[string]interface{}{
			"audience": audience,
			"issuer":   issuer,
			"tenant":   tenant,
		},
	}
}

func dataSourceFlattenVpnServerConfigurationClientRootCertificates(input *[]virtualwans.VpnServerConfigVpnClientRootCertificate) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, v := range *input {
		name := ""
		if v.Name != nil {
			name = *v.Name
		}

		publicCertData := ""
		if v.PublicCertData != nil {
			publicCertData = *v.PublicCertData
		}

		output = append(output, map[string]interface{}{
			"name":             name,
			"public_cert_data": publicCertData,
		})
	}

	return output
}

func dataSourceFlattenVpnServerConfigurationClientRevokedCertificates(input *[]virtualwans.VpnServerConfigVpnClientRevokedCertificate) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)
	for _, v := range *input {
		name := ""
		if v.Name != nil {
			name = *v.Name
		}

		thumbprint := ""
		if v.Thumbprint != nil {
			thumbprint = *v.Thumbprint
		}

		output = append(output, map[string]interface{}{
			"name":       name,
			"thumbprint": thumbprint,
		})
	}
	return output
}

func dataSourceFlattenVpnServerConfigurationIPSecPolicies(input *[]virtualwans.IPsecPolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)
	for _, v := range *input {
		output = append(output, map[string]interface{}{
			"dh_group":               string(v.DhGroup),
			"ipsec_encryption":       string(v.IPsecEncryption),
			"ipsec_integrity":        string(v.IPsecIntegrity),
			"ike_encryption":         string(v.IkeEncryption),
			"ike_integrity":          string(v.IkeIntegrity),
			"pfs_group":              string(v.PfsGroup),
			"sa_data_size_kilobytes": int(v.SaDataSizeKilobytes),
			"sa_lifetime_seconds":    int(v.SaLifeTimeSeconds),
		})
	}
	return output
}

func dataSourceFlattenVpnServerConfigurationRadius(input *virtualwans.VpnServerConfigurationProperties) []interface{} {
	if input == nil || (input.RadiusServerAddress == nil && (input.RadiusServers == nil || len(*input.RadiusServers) == 0)) {
		return []interface{}{}
	}

	clientRootCertificates := make([]interface{}, 0)
	if input.RadiusClientRootCertificates != nil {
		for _, v := range *input.RadiusClientRootCertificates {
			name := ""
			if v.Name != nil {
				name = *v.Name
			}

			thumbprint := ""
			if v.Thumbprint != nil {
				thumbprint = *v.Thumbprint
			}

			clientRootCertificates = append(clientRootCertificates, map[string]interface{}{
				"name":       name,
				"thumbprint": thumbprint,
			})
		}
	}

	serverRootCertificates := make([]interface{}, 0)
	if input.RadiusServerRootCertificates != nil {
		for _, v := range *input.RadiusServerRootCertificates {
			name := ""
			if v.Name != nil {
				name = *v.Name
			}

			publicCertData := ""
			if v.PublicCertData != nil {
				publicCertData = *v.PublicCertData
			}

			serverRootCertificates = append(serverRootCertificates, map[string]interface{}{
				"name":             name,
				"public_cert_data": publicCertData,
			})
		}
	}

	servers := make([]interface{}, 0)
	if input.RadiusServers != nil && len(*input.RadiusServers) > 0 {
		for _, v := range *input.RadiusServers {
			servers = append(servers, map[string]interface{}{
				"address": v.RadiusServerAddress,
				"secret":  pointer.From(v.RadiusServerSecret),
				"score":   pointer.From(v.RadiusServerScore),
			})
		}
	}

	return []interface{}{
		map[string]interface{}{
			"client_root_certificate": clientRootCertificates,
			"server_root_certificate": serverRootCertificates,
			"server":                  servers,
		},
	}
}

func dataSourceFlattenVpnServerConfigurationVPNProtocols(input *[]virtualwans.VpnGatewayTunnelingProtocol) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, v := range *input {
		output = append(output, string(v))
	}

	return output
}
