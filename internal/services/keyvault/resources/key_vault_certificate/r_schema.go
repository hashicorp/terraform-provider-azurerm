// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	kv "github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

var validDigiCertCertificateTypes = []string{
	"Basic-SSL",
	"BasicEV-SSL",
	"OV-SSL",
	"EV-SSL",
	"ProEV-SSL",
	"ProOV-SSL",
}

func keyVaultCertificateSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Key Vault Certificate.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: keyvault.ValidateNestedItemName,
		},

		"key_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KeyVaultId{}, "Key Vault"),

		"certificate": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			AtLeastOneOf: []string{
				"certificate_policy",
				"certificate",
			},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"contents": {
						Description:  "The base64-encoded certificate contents.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"password": {
						Description: "The password associated with the certificate.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Sensitive:   true,
					},
				},
			},
		},

		"certificate_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			AtLeastOneOf: []string{
				"certificate_policy",
				"certificate",
			},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"issuer_parameters": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description: "The name of the Certificate Issuer.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
								"certificate_type": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringInSlice(validDigiCertCertificateTypes, false),
								},
							},
						},
					},
					"key_properties": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"curve": {
									Description:  "The curve to use when creating an `EC` key. This field will be required in a future release if `key_type` is `EC` or `EC-HSM`.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Computed:     true, // azignore:AZS007 - pre-existing violation
									ValidateFunc: validation.StringInEnumSlice(kv.PossibleJSONWebKeyCurveNameValues(), false),
								},
								"exportable": {
									Description: "Whether this certificate is exportable.",
									Type:        pluginsdk.TypeBool,
									Required:    true,
								},
								"key_size": {
									Description: "The size of the key used in the certificate. This property is required when using RSA keys.",
									Type:        pluginsdk.TypeInt,
									Optional:    true,
									Computed:    true, // azignore:AZS007 - pre-existing violation
									ValidateFunc: validation.IntInSlice([]int{
										256,
										384,
										521,
										2048,
										3072,
										4096,
									}),
								},
								"key_type": {
									Description: "The type of key.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ValidateFunc: validation.StringInSlice([]string{
										string(kv.JSONWebKeyTypeEC),
										string(kv.JSONWebKeyTypeECHSM),
										string(kv.JSONWebKeyTypeRSA),
										string(kv.JSONWebKeyTypeRSAHSM),
										string(kv.JSONWebKeyTypeOct),
									}, false),
								},
								"reuse_key": {
									Description: "Whether the key is reusable.",
									Type:        pluginsdk.TypeBool,
									Required:    true,
								},
							},
						},
					},
					"lifetime_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"action": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"action_type": {
												Description:  "The Type of action to be performed when the lifetime trigger is triggered.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInEnumSlice(kv.PossibleCertificatePolicyActionValues(), false),
											},
										},
									},
								},
								// lintignore:XS003
								"trigger": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"days_before_expiry": {
												Description: "The number of days before the Certificate expires that the action associated with this Trigger should run.",
												Type:        pluginsdk.TypeInt,
												Optional:    true,
											},
											"lifetime_percentage": {
												Description: "The percentage at which during the Certificates Lifetime the action associated with this Trigger should run.",
												Type:        pluginsdk.TypeInt,
												Optional:    true,
											},
										},
									},
								},
							},
						},
					},
					"secret_properties": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"content_type": {
									Description: "The Content-Type of the Certificate, such as `application/x-pkcs12` for a PFX or `application/x-pem-file` for a PEM.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
							},
						},
					},

					"x509_certificate_properties": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"extended_key_usage": {
									Description: "A list of Extended/Enhanced Key Usages.",
									Type:        pluginsdk.TypeList,
									Optional:    true,
									Computed:    true, // azignore:AZS007 - pre-existing violation
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"key_usage": {
									Description: "A list of uses associated with this Key.",
									Type:        pluginsdk.TypeSet,
									Required:    true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInEnumSlice(kv.PossibleKeyUsageTypeValues(), false),
									},
								},
								"subject": {
									Description: "The Certificate's Subject.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
								"subject_alternative_names": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Computed: true, // azignore:AZS007 - pre-existing violation
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"emails": {
												Description: "A list of email addresses identified by this Certificate.",
												Type:        pluginsdk.TypeSet,
												Optional:    true,
												Set:         pluginsdk.HashString,
												AtLeastOneOf: []string{
													"certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.emails",
													"certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.dns_names",
													"certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.upns",
												},
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},
											"dns_names": {
												Description: "A list of alternative DNS names (FQDNs) identified by the Certificate.",
												Type:        pluginsdk.TypeSet,
												Optional:    true,
												Set:         pluginsdk.HashString,
												AtLeastOneOf: []string{
													"certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.emails",
													"certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.dns_names",
													"certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.upns",
												},
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},
											"upns": {
												Description: "A list of User Principal Names identified by the Certificate.",
												Type:        pluginsdk.TypeSet,
												Optional:    true,
												Set:         pluginsdk.HashString,
												AtLeastOneOf: []string{
													"certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.emails",
													"certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.dns_names",
													"certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.upns",
												},
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},
										},
									},
								},
								"validity_in_months": {
									Description: "The Certificates Validity Period in Months.",
									Type:        pluginsdk.TypeInt,
									Required:    true,
								},
							},
						},
					},
				},
			},
		},

		// Computed
		"certificate_attribute": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"created": {
						Description: "The create time of the Key Vault Certificate.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"enabled": {
						Description: "Whether the Key Vault Certificate is enabled.",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},

					"expires": {
						Description: "The expires time of the Key Vault Certificate.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"not_before": {
						Description: "The not before valid time of the Key Vault Certificate.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"recovery_level": {
						Description: "The deletion recovery level of the Key Vault Certificate.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"updated": {
						Description: "The recent update time of the Key Vault Certificate.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"resource_manager_id": {
			Description: "The (Versioned) ID for this Key Vault Certificate. This property points to a specific version of a Key Vault Certificate, as such using this won't auto-rotate values if used in other Azure Services.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"resource_manager_versionless_id": {
			Description: "The Versionless ID of the Key Vault Certificate. This property allows other Azure Services (that support it) to auto-rotate their value when the Key Vault Certificate is updated.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"version": {
			Description: "The current version of the Key Vault Certificate.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"secret_id": {
			Description: "The ID of the associated Key Vault Secret.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"versionless_id": {
			Description: "The Base ID of the Key Vault Certificate.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"versionless_secret_id": {
			Description: "The Base ID of the Key Vault Secret.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"certificate_data": {
			Description: "The raw Key Vault Certificate data represented as a hexadecimal string.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"certificate_data_base64": {
			Description: "The Base64 encoded Key Vault Certificate data.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"thumbprint": {
			Description: "The X509 Thumbprint of the Key Vault Certificate represented as a hexadecimal string.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.Tags(),
	}
}
