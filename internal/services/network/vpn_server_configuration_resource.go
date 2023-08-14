// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVPNServerConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVPNServerConfigurationCreateUpdate,
		Read:   resourceVPNServerConfigurationRead,
		Update: resourceVPNServerConfigurationCreateUpdate,
		Delete: resourceVPNServerConfigurationDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualwans.ParseVpnServerConfigurationID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"vpn_authentication_types": {
				Type:     pluginsdk.TypeList,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(virtualwans.VpnAuthenticationTypeAAD),
						string(virtualwans.VpnAuthenticationTypeCertificate),
						string(virtualwans.VpnAuthenticationTypeRadius),
					}, false),
				},
			},

			// Optional
			"azure_active_directory_authentication": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"audience": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"issuer": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"tenant": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					},
				},
			},

			"client_revoked_certificate": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"thumbprint": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					},
				},
			},

			"client_root_certificate": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"public_cert_data": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					},
				},
			},

			"ipsec_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"dh_group": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForDhGroup(), false),
						},

						"ike_encryption": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForIkeEncryption(), false),
						},

						"ike_integrity": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForIkeIntegrity(), false),
						},

						"ipsec_encryption": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForIPsecEncryption(), false),
						},

						"ipsec_integrity": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForIPsecIntegrity(), false),
						},

						"pfs_group": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForPfsGroup(), false),
						},

						"sa_lifetime_seconds": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},

						"sa_data_size_kilobytes": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},
					},
				},
			},

			"radius": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"server": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"address": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"secret": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
										Sensitive:    true,
									},

									"score": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 30),
									},
								},
							},
						},

						"client_root_certificate": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:     pluginsdk.TypeString,
										Required: true,
									},

									"thumbprint": {
										Type:     pluginsdk.TypeString,
										Required: true,
									},
								},
							},
						},

						"server_root_certificate": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:     pluginsdk.TypeString,
										Required: true,
									},

									"public_cert_data": {
										Type:     pluginsdk.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"vpn_protocols": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForVpnGatewayTunnelingProtocol(), false),
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceVPNServerConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualwans.NewVpnServerConfigurationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.VpnServerConfigurationsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_vpn_server_configuration", id.ID())
		}
	}

	aadAuthenticationRaw := d.Get("azure_active_directory_authentication").([]interface{})
	aadAuthentication := expandVpnServerConfigurationAADAuthentication(aadAuthenticationRaw)

	clientRevokedCertsRaw := d.Get("client_revoked_certificate").(*pluginsdk.Set).List()
	clientRevokedCerts := expandVpnServerConfigurationClientRevokedCertificates(clientRevokedCertsRaw)

	clientRootCertsRaw := d.Get("client_root_certificate").(*pluginsdk.Set).List()
	clientRootCerts := expandVpnServerConfigurationClientRootCertificates(clientRootCertsRaw)

	ipSecPoliciesRaw := d.Get("ipsec_policy").([]interface{})
	ipSecPolicies := expandVpnServerConfigurationIPSecPolicies(ipSecPoliciesRaw)

	radius := expandVpnServerConfigurationRadius(d.Get("radius").([]interface{}))

	vpnProtocolsRaw := d.Get("vpn_protocols").(*pluginsdk.Set).List()
	vpnProtocols := expandVpnServerConfigurationVPNProtocols(vpnProtocolsRaw)

	supportsAAD := false
	supportsCertificates := false
	supportsRadius := false

	vpnAuthenticationTypesRaw := d.Get("vpn_authentication_types").([]interface{})
	vpnAuthenticationTypes := make([]virtualwans.VpnAuthenticationType, 0)
	for _, v := range vpnAuthenticationTypesRaw {
		authType := virtualwans.VpnAuthenticationType(v.(string))

		switch authType {
		case virtualwans.VpnAuthenticationTypeAAD:
			supportsAAD = true

		case virtualwans.VpnAuthenticationTypeCertificate:
			supportsCertificates = true

		case virtualwans.VpnAuthenticationTypeRadius:
			supportsRadius = true

		default:
			return fmt.Errorf("Unsupported `vpn_authentication_type`: %q", authType)
		}

		vpnAuthenticationTypes = append(vpnAuthenticationTypes, authType)
	}

	props := virtualwans.VpnServerConfigurationProperties{
		AadAuthenticationParameters:  aadAuthentication,
		VpnAuthenticationTypes:       &vpnAuthenticationTypes,
		VpnClientRootCertificates:    clientRootCerts,
		VpnClientRevokedCertificates: clientRevokedCerts,
		VpnClientIPsecPolicies:       ipSecPolicies,
		VpnProtocols:                 vpnProtocols,
	}

	if supportsAAD && aadAuthentication == nil {
		return fmt.Errorf("`azure_active_directory_authentication` must be specified when `vpn_authentication_type` is set to `AAD`")
	}

	// parameter:VpnServerConfigVpnClientRootCertificates is not specified when VpnAuthenticationType as Certificate is selected.
	if supportsCertificates && len(clientRootCertsRaw) == 0 {
		return fmt.Errorf("`client_root_certificate` must be specified when `vpn_authentication_type` is set to `Certificate`")
	}

	if supportsRadius {
		if radius == nil {
			return fmt.Errorf("`radius` must be specified when `vpn_authentication_type` is set to `Radius`")
		}

		if radius.servers != nil && len(*radius.servers) != 0 {
			props.RadiusServers = radius.servers
		}

		props.RadiusServerAddress = utils.String(radius.address)
		props.RadiusServerSecret = utils.String(radius.secret)

		props.RadiusClientRootCertificates = radius.clientRootCertificates
		props.RadiusServerRootCertificates = radius.serverRootCertificates
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	parameters := virtualwans.VpnServerConfiguration{
		Location:   utils.String(location),
		Properties: &props,
		Tags:       tags.Expand(t),
	}

	if err := client.VpnServerConfigurationsCreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVPNServerConfigurationRead(d, meta)
}

func resourceVPNServerConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVpnServerConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.VpnServerConfigurationsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.VpnServerConfigurationName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			flattenedAADAuthentication := flattenVpnServerConfigurationAADAuthentication(props.AadAuthenticationParameters)
			if err := d.Set("azure_active_directory_authentication", flattenedAADAuthentication); err != nil {
				return fmt.Errorf("setting `azure_active_directory_authentication`: %+v", err)
			}

			flattenedClientRootCerts := flattenVpnServerConfigurationClientRootCertificates(props.VpnClientRootCertificates)
			if err := d.Set("client_root_certificate", flattenedClientRootCerts); err != nil {
				return fmt.Errorf("setting `client_root_certificate`: %+v", err)
			}

			flattenedClientRevokedCerts := flattenVpnServerConfigurationClientRevokedCertificates(props.VpnClientRevokedCertificates)
			if err := d.Set("client_revoked_certificate", flattenedClientRevokedCerts); err != nil {
				return fmt.Errorf("setting `client_revoked_certificate`: %+v", err)
			}

			flattenedIPSecPolicies := flattenVpnServerConfigurationIPSecPolicies(props.VpnClientIPsecPolicies)
			if err := d.Set("ipsec_policy", flattenedIPSecPolicies); err != nil {
				return fmt.Errorf("setting `ipsec_policy`: %+v", err)
			}

			flattenedRadius := flattenVpnServerConfigurationRadius(props)
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

			flattenedVpnProtocols := flattenVpnServerConfigurationVPNProtocols(props.VpnProtocols)
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

func resourceVPNServerConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVpnServerConfigurationID(d.Id())
	if err != nil {
		return err
	}

	if err := client.VpnServerConfigurationsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandVpnServerConfigurationAADAuthentication(input []interface{}) *virtualwans.AadAuthenticationParameters {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &virtualwans.AadAuthenticationParameters{
		AadAudience: utils.String(v["audience"].(string)),
		AadIssuer:   utils.String(v["issuer"].(string)),
		AadTenant:   utils.String(v["tenant"].(string)),
	}
}

func flattenVpnServerConfigurationAADAuthentication(input *virtualwans.AadAuthenticationParameters) []interface{} {
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

func expandVpnServerConfigurationClientRootCertificates(input []interface{}) *[]virtualwans.VpnServerConfigVpnClientRootCertificate {
	clientRootCertificates := make([]virtualwans.VpnServerConfigVpnClientRootCertificate, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})
		clientRootCertificates = append(clientRootCertificates, virtualwans.VpnServerConfigVpnClientRootCertificate{
			Name:           utils.String(raw["name"].(string)),
			PublicCertData: utils.String(raw["public_cert_data"].(string)),
		})
	}

	return &clientRootCertificates
}

func flattenVpnServerConfigurationClientRootCertificates(input *[]virtualwans.VpnServerConfigVpnClientRootCertificate) []interface{} {
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

func expandVpnServerConfigurationClientRevokedCertificates(input []interface{}) *[]virtualwans.VpnServerConfigVpnClientRevokedCertificate {
	clientRevokedCertificates := make([]virtualwans.VpnServerConfigVpnClientRevokedCertificate, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})
		clientRevokedCertificates = append(clientRevokedCertificates, virtualwans.VpnServerConfigVpnClientRevokedCertificate{
			Name:       utils.String(raw["name"].(string)),
			Thumbprint: utils.String(raw["thumbprint"].(string)),
		})
	}

	return &clientRevokedCertificates
}

func flattenVpnServerConfigurationClientRevokedCertificates(input *[]virtualwans.VpnServerConfigVpnClientRevokedCertificate) []interface{} {
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

func expandVpnServerConfigurationIPSecPolicies(input []interface{}) *[]virtualwans.IPsecPolicy {
	ipSecPolicies := make([]virtualwans.IPsecPolicy, 0)

	for _, raw := range input {
		v := raw.(map[string]interface{})
		ipSecPolicies = append(ipSecPolicies, virtualwans.IPsecPolicy{
			DhGroup:             virtualwans.DhGroup(v["dh_group"].(string)),
			IkeEncryption:       virtualwans.IkeEncryption(v["ike_encryption"].(string)),
			IkeIntegrity:        virtualwans.IkeIntegrity(v["ike_integrity"].(string)),
			IPsecEncryption:     virtualwans.IPsecEncryption(v["ipsec_encryption"].(string)),
			IPsecIntegrity:      virtualwans.IPsecIntegrity(v["ipsec_integrity"].(string)),
			PfsGroup:            virtualwans.PfsGroup(v["pfs_group"].(string)),
			SaLifeTimeSeconds:   int64(v["sa_lifetime_seconds"].(int)),
			SaDataSizeKilobytes: int64(v["sa_data_size_kilobytes"].(int)),
		})
	}

	return &ipSecPolicies
}

func flattenVpnServerConfigurationIPSecPolicies(input *[]virtualwans.IPsecPolicy) []interface{} {
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

type vpnServerConfigurationRadius struct {
	address                string
	secret                 string
	servers                *[]virtualwans.RadiusServer
	clientRootCertificates *[]virtualwans.VpnServerConfigRadiusClientRootCertificate
	serverRootCertificates *[]virtualwans.VpnServerConfigRadiusServerRootCertificate
}

func expandVpnServerConfigurationRadius(input []interface{}) *vpnServerConfigurationRadius {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	val := input[0].(map[string]interface{})

	clientRootCertificates := make([]virtualwans.VpnServerConfigRadiusClientRootCertificate, 0)
	clientRootCertsRaw := val["client_root_certificate"].(*pluginsdk.Set).List()
	for _, raw := range clientRootCertsRaw {
		v := raw.(map[string]interface{})
		clientRootCertificates = append(clientRootCertificates, virtualwans.VpnServerConfigRadiusClientRootCertificate{
			Name:       utils.String(v["name"].(string)),
			Thumbprint: utils.String(v["thumbprint"].(string)),
		})
	}

	serverRootCertificates := make([]virtualwans.VpnServerConfigRadiusServerRootCertificate, 0)
	serverRootCertsRaw := val["server_root_certificate"].(*pluginsdk.Set).List()
	for _, raw := range serverRootCertsRaw {
		v := raw.(map[string]interface{})
		serverRootCertificates = append(serverRootCertificates, virtualwans.VpnServerConfigRadiusServerRootCertificate{
			Name:           utils.String(v["name"].(string)),
			PublicCertData: utils.String(v["public_cert_data"].(string)),
		})
	}

	radiusServers := make([]virtualwans.RadiusServer, 0)
	address := ""
	secret := ""

	if val["server"] != nil {
		radiusServersRaw := val["server"].([]interface{})
		for _, raw := range radiusServersRaw {
			v := raw.(map[string]interface{})
			radiusServers = append(radiusServers, virtualwans.RadiusServer{
				RadiusServerAddress: v["address"].(string),
				RadiusServerSecret:  utils.String(v["secret"].(string)),
				RadiusServerScore:   utils.Int64(int64(v["score"].(int))),
			})
		}
	}

	return &vpnServerConfigurationRadius{
		address:                address,
		secret:                 secret,
		servers:                &radiusServers,
		clientRootCertificates: &clientRootCertificates,
		serverRootCertificates: &serverRootCertificates,
	}
}

func flattenVpnServerConfigurationRadius(input *virtualwans.VpnServerConfigurationProperties) []interface{} {
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

func expandVpnServerConfigurationVPNProtocols(input []interface{}) *[]virtualwans.VpnGatewayTunnelingProtocol {
	vpnProtocols := make([]virtualwans.VpnGatewayTunnelingProtocol, 0)

	for _, v := range input {
		vpnProtocols = append(vpnProtocols, virtualwans.VpnGatewayTunnelingProtocol(v.(string)))
	}

	return &vpnProtocols
}

func flattenVpnServerConfigurationVPNProtocols(input *[]virtualwans.VpnGatewayTunnelingProtocol) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, v := range *input {
		output = append(output, string(v))
	}

	return output
}
