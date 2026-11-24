// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_custom_domain

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/afddomains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/secrets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cdnFrontdoorCustomDomainSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Front Door Custom Domain.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FrontDoorCustomDomainName,
		},

		// NOTE: I need this fake field because I need the profile name during the create operation
		"cdn_frontdoor_profile_id": {
			Description:  "The ID of the Front Door Profile.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: profiles.ValidateProfileID,
		},

		"dns_zone_id": {
			Description:  "The ID of the Azure DNS Zone which should be used for this Front Door Custom Domain.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: zones.ValidateDnsZoneID,
		},

		"host_name": {
			Description:  "The host name of the domain.",
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validate.FrontDoorCustomDomainHostName,
		},

		"tls": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"certificate_type": {
						Description: "The source of the SSL certificate.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Default:     string(afddomains.AfdCertificateTypeManagedCertificate),
						ValidateFunc: validation.StringInSlice([]string{
							string(afddomains.AfdCertificateTypeCustomerCertificate),
							string(afddomains.AfdCertificateTypeManagedCertificate),
						}, false),
					},

					"minimum_version": {
						Description: "TLS protocol version that will be used for HTTPS.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Default:     string(afddomains.AfdMinimumTlsVersionTLSOneTwo),
						ValidateFunc: validation.StringInSlice([]string{
							string(afddomains.AfdMinimumTlsVersionTLSOneTwo),
						}, false),
					},

					"cdn_frontdoor_secret_id": {
						Description: "Resource ID of the Front Door Secret.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						// Note: O+C because if the secret is managed by FrontDoor this will cause a perpetual diff
						Computed:     true,
						ValidateFunc: secrets.ValidateSecretID,
					},

					"cipher_suite": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"type": {
									Description: "The cipher suite set type.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ValidateFunc: validation.StringInSlice([]string{
										string(afddomains.AfdCipherSuiteSetTypeCustomized),
										string(afddomains.AfdCipherSuiteSetTypeTLSOneTwoTwoZeroTwoThree),
										string(afddomains.AfdCipherSuiteSetTypeTLSOneTwoTwoZeroTwoTwo),
										// Explicitly exclude TLS10_2019 - TLS 1.0/1.1 support retired March 1, 2025
									}, false),
								},

								"custom_ciphers": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"tls12": {
												Description: "A set of TLS 1.2 cipher suites.",
												Type:        pluginsdk.TypeSet,
												Optional:    true,
												AtLeastOneOf: []string{
													"tls.0.cipher_suite.0.custom_ciphers.0.tls12",
													"tls.0.cipher_suite.0.custom_ciphers.0.tls13",
												},
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
													ValidateFunc: validation.StringInSlice(
														frontDoorCustomDomainCustomizedCipherSuitesForTls12,
														false,
													),
												},
											},

											"tls13": {
												Description: "A set of TLS 1.3 cipher suites.",
												Type:        pluginsdk.TypeSet,
												AtLeastOneOf: []string{
													"tls.0.cipher_suite.0.custom_ciphers.0.tls12",
													"tls.0.cipher_suite.0.custom_ciphers.0.tls13",
												},
												Optional: true,
												// NOTE: O+C Azure Front Door returns TLS 1.3 cipher suites even when `tls13` is not specified (see https://learn.microsoft.com/azure/frontdoor/standard-premium/tls-policy#custom-tls-policy)
												Computed: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
													ValidateFunc: validation.StringInSlice(
														afddomains.PossibleValuesForAfdCustomizedCipherSuiteForTls13(),
														false,
													),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},

		"expiration_date": {
			Description: "The date and time that the token expires.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"validation_token": {
			Description: "Challenge used for DNS TXT record or file based validation.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
