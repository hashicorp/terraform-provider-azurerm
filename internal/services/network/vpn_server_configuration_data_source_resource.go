// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type VPNServerConfigurationDataSource struct{}

var _ sdk.DataSource = VPNServerConfigurationDataSource{}

type VPNServerConfigurationDataSourceModel struct {
	Name                               string                                    `tfschema:"name"`
	ResourceGroup                      string                                    `tfschema:"resource_group_name"`
	Location                           string                                    `tfschema:"location"`
	VpnAuthenticationTypes             []string                                  `tfschema:"vpn_authentication_types"`
	AzureActiveDirectoryAuthentication []AzureActiveDirectoryAuthenticationModel `tfschema:"azure_active_directory_authentication"`
	ClientRevokedCertificate           []ClientRevokedCertificateModel           `tfschema:"client_revoked_certificate"`
	ClientRootCertificate              []ClientRootCertificateModel              `tfschema:"client_root_certificate"`
	IpsecPolicy                        []IpsecPolicyModel                        `tfschema:"ipsec_policy"`
	Radius                             []RadiusModel                             `tfschema:"radius"`
	VpnProtocols                       []string                                  `tfschema:"vpn_protocols"`
	Tags                               map[string]string                         `tfschema:"tags"`
}

type AzureActiveDirectoryAuthenticationModel struct {
	Audience string `tfschema:"audience"`
	Issuer   string `tfschema:"issuer"`
	Tenant   string `tfschema:"tenant"`
}

type ClientRevokedCertificateModel struct {
	Name       string `tfschema:"name"`
	Thumbprint string `tfschema:"thumbprint"`
}

type ClientRootCertificateModel struct {
	Name           string `tfschema:"name"`
	PublicCertData string `tfschema:"public_cert_data"`
}

type IpsecPolicyModel struct {
	DhGroup             string `tfschema:"dh_group"`
	IkeEncryption       string `tfschema:"ike_encryption"`
	IkeIntegrity        string `tfschema:"ike_integrity"`
	IpsecEncryption     string `tfschema:"ipsec_encryption"`
	IpsecIntegrity      string `tfschema:"ipsec_integrity"`
	PfsGroup            string `tfschema:"pfs_group"`
	SaLifetimeSeconds   int64  `tfschema:"sa_lifetime_seconds"`
	SaDataSizeKilobytes int64  `tfschema:"sa_data_size_kilobytes"`
}

type RadiusModel struct {
	Server                []ServerModel                      `tfschema:"server"`
	ClientRootCertificate []RadiusClientRootCertificateModel `tfschema:"client_root_certificate"`
	ServerRootCertificate []ClientRootCertificateModel       `tfschema:"server_root_certificate"`
}

type ServerModel struct {
	Address string `tfschema:"address"`
	Secret  string `tfschema:"secret"`
	Score   int64  `tfschema:"score"`
}

type RadiusClientRootCertificateModel struct {
	Name       string `tfschema:"name"`
	Thumbprint string `tfschema:"thumbprint"`
}

func (d VPNServerConfigurationDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d VPNServerConfigurationDataSource) ModelObject() interface{} {
	return &VPNServerConfigurationDataSource{}
}

func (d VPNServerConfigurationDataSource) ResourceType() string {
	return "azurerm_vpn_server_configuration"
}

func (d VPNServerConfigurationDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
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
	}
}

func (d VPNServerConfigurationDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.VirtualWANs
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model VPNServerConfigurationDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := virtualwans.NewVpnServerConfigurationID(subscriptionId, model.ResourceGroup, model.Name)

			resp, err := client.VpnServerConfigurationsGet(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			m := VPNServerConfigurationDataSourceModel{
				Name:          id.VpnServerConfigurationName,
				ResourceGroup: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				m.Location = pointer.ToString(model.Location)
				if tags := model.Tags; tags != nil {
					m.Tags = pointer.ToMapOfStringStrings(tags)
				}

				if props := resp.Model.Properties; props != nil {
					m.AzureActiveDirectoryAuthentication = dataSourceFlattenVpnServerConfigurationAADAuthentication(props.AadAuthenticationParameters)
					m.ClientRootCertificate = dataSourceFlattenVpnServerConfigurationClientRootCertificates(props.VpnClientRootCertificates)
					m.ClientRevokedCertificate = dataSourceFlattenVpnServerConfigurationClientRevokedCertificates(props.VpnClientRevokedCertificates)
					m.IpsecPolicy = dataSourceFlattenVpnServerConfigurationIPSecPolicies(props.VpnClientIPsecPolicies)
					m.Radius = dataSourceFlattenVpnServerConfigurationRadius(props)
					m.VpnAuthenticationTypes = dataSourceFlattenVpnServerConfigurationVpnAuthenticationTypes(props.VpnAuthenticationTypes)
					m.VpnProtocols = dataSourceFlattenVpnServerConfigurationVPNProtocols(props.VpnProtocols)
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&m)
		},
	}
}

func dataSourceFlattenVpnServerConfigurationAADAuthentication(input *virtualwans.AadAuthenticationParameters) []AzureActiveDirectoryAuthenticationModel {
	if input == nil {
		return []AzureActiveDirectoryAuthenticationModel{}
	}

	return []AzureActiveDirectoryAuthenticationModel{
		{
			Audience: pointer.ToString(input.AadAudience),
			Issuer:   pointer.ToString(input.AadIssuer),
			Tenant:   pointer.ToString(input.AadTenant),
		},
	}
}

func dataSourceFlattenVpnServerConfigurationClientRootCertificates(input *[]virtualwans.VpnServerConfigVpnClientRootCertificate) []ClientRootCertificateModel {
	if input == nil {
		return []ClientRootCertificateModel{}
	}

	output := make([]ClientRootCertificateModel, 0)

	for _, v := range *input {
		if v.Name == nil {
			continue
		}
		output = append(output, ClientRootCertificateModel{
			Name:           pointer.ToString(v.Name),
			PublicCertData: pointer.ToString(v.PublicCertData),
		})
	}

	return output
}

func dataSourceFlattenVpnServerConfigurationClientRevokedCertificates(input *[]virtualwans.VpnServerConfigVpnClientRevokedCertificate) []ClientRevokedCertificateModel {
	if input == nil {
		return []ClientRevokedCertificateModel{}
	}

	output := make([]ClientRevokedCertificateModel, 0)
	for _, v := range *input {
		if v.Name == nil {
			continue
		}

		output = append(output, ClientRevokedCertificateModel{
			Name:       pointer.ToString(v.Name),
			Thumbprint: pointer.ToString(v.Thumbprint),
		})
	}
	return output
}

func dataSourceFlattenVpnServerConfigurationIPSecPolicies(input *[]virtualwans.IPsecPolicy) []IpsecPolicyModel {
	if input == nil {
		return []IpsecPolicyModel{}
	}

	output := make([]IpsecPolicyModel, 0)
	for _, v := range *input {
		output = append(output, IpsecPolicyModel{
			DhGroup:             string(v.DhGroup),
			IkeEncryption:       string(v.IPsecEncryption),
			IkeIntegrity:        string(v.IPsecIntegrity),
			IpsecEncryption:     string(v.IkeEncryption),
			IpsecIntegrity:      string(v.IkeIntegrity),
			PfsGroup:            string(v.PfsGroup),
			SaLifetimeSeconds:   v.SaDataSizeKilobytes,
			SaDataSizeKilobytes: v.SaLifeTimeSeconds,
		})
	}
	return output
}

func dataSourceFlattenVpnServerConfigurationRadius(input *virtualwans.VpnServerConfigurationProperties) []RadiusModel {
	if input == nil || (input.RadiusServerAddress == nil && (input.RadiusServers == nil || len(*input.RadiusServers) == 0)) {
		return []RadiusModel{}
	}

	clientRootCertificates := make([]RadiusClientRootCertificateModel, 0)
	if input.RadiusClientRootCertificates != nil {
		for _, v := range *input.RadiusClientRootCertificates {
			if v.Name == nil {
				continue
			}
			clientRootCertificates = append(clientRootCertificates, RadiusClientRootCertificateModel{
				Name:       pointer.ToString(v.Name),
				Thumbprint: pointer.ToString(v.Thumbprint),
			})
		}
	}

	serverRootCertificates := make([]ClientRootCertificateModel, 0)
	if input.RadiusServerRootCertificates != nil {
		for _, v := range *input.RadiusServerRootCertificates {
			if v.Name == nil {
				continue
			}

			serverRootCertificates = append(serverRootCertificates, ClientRootCertificateModel{
				Name:           pointer.ToString(v.Name),
				PublicCertData: pointer.ToString(v.PublicCertData),
			})
		}
	}

	servers := make([]ServerModel, 0)
	if input.RadiusServers != nil && len(*input.RadiusServers) > 0 {
		for _, v := range *input.RadiusServers {
			servers = append(servers, ServerModel{
				Address: v.RadiusServerAddress,
				Secret:  pointer.ToString(v.RadiusServerSecret),
				Score:   pointer.ToInt64(v.RadiusServerScore),
			})
		}
	}

	return []RadiusModel{
		{
			Server:                servers,
			ClientRootCertificate: clientRootCertificates,
			ServerRootCertificate: serverRootCertificates,
		},
	}
}

func dataSourceFlattenVpnServerConfigurationVpnAuthenticationTypes(input *[]virtualwans.VpnAuthenticationType) []string {
	if input == nil {
		return []string{}
	}

	output := make([]string, 0)

	for _, v := range *input {
		output = append(output, string(v))
	}

	return output
}

func dataSourceFlattenVpnServerConfigurationVPNProtocols(input *[]virtualwans.VpnGatewayTunnelingProtocol) []string {
	if input == nil {
		return []string{}
	}

	output := make([]string, 0)

	for _, v := range *input {
		output = append(output, string(v))
	}

	return output
}
