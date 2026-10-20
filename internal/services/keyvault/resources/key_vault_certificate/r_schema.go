// Copyright IBM Corp. 2014, 2025
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

func keyVaultCertificateSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: keyvault.ValidateNestedItemName,
		},

		"key_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KeyVaultId{}),

		"certificate": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			AtLeastOneOf: []string{
				"certificate_policy",
				"certificate",
			},
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"contents": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"password": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						Sensitive: true,
					},
				},
			},
		},

		"certificate_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			AtLeastOneOf: []string{
				"certificate_policy",
				"certificate",
			},
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"issuer_parameters": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
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
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Computed:     true, // azignore:AZS007 - pre-existing violation
									ValidateFunc: validation.StringInEnumSlice(kv.PossibleJSONWebKeyCurveNameValues(), false),
								},
								"exportable": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
								"key_size": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									Computed: true, // azignore:AZS007 - pre-existing violation
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
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(kv.JSONWebKeyTypeEC),
										string(kv.JSONWebKeyTypeECHSM),
										string(kv.JSONWebKeyTypeRSA),
										string(kv.JSONWebKeyTypeRSAHSM),
										string(kv.JSONWebKeyTypeOct),
									}, false),
								},
								"reuse_key": {
									Type:     pluginsdk.TypeBool,
									Required: true,
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
												Type:     pluginsdk.TypeInt,
												Optional: true,
											},
											"lifetime_percentage": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
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
									Type:     pluginsdk.TypeString,
									Required: true,
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
									Type:     pluginsdk.TypeList,
									Optional: true,
									Computed: true, // azignore:AZS007 - pre-existing violation
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"key_usage": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInEnumSlice(kv.PossibleKeyUsageTypeValues(), false),
									},
								},
								"subject": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"subject_alternative_names": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Computed: true, // azignore:AZS007 - pre-existing violation
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"emails": {
												Type:     pluginsdk.TypeSet,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
												Set: pluginsdk.HashString,
												AtLeastOneOf: []string{
													"certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.emails",
													"certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.dns_names",
													"certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.upns",
												},
											},
											"dns_names": {
												Type:     pluginsdk.TypeSet,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
												Set: pluginsdk.HashString,
												AtLeastOneOf: []string{
													"certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.emails",
													"certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.dns_names",
													"certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.upns",
												},
											},
											"upns": {
												Type:     pluginsdk.TypeSet,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
												Set: pluginsdk.HashString,
												AtLeastOneOf: []string{
													"certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.emails",
													"certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.dns_names",
													"certificate_policy.0.x509_certificate_properties.0.subject_alternative_names.0.upns",
												},
											},
										},
									},
								},
								"validity_in_months": {
									Type:     pluginsdk.TypeInt,
									Required: true,
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
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"expires": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"not_before": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"recovery_level": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"updated": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"resource_manager_id": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},

		"resource_manager_versionless_id": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},

		"version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secret_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"versionless_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"versionless_secret_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"certificate_data": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"certificate_data_base64": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"thumbprint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.Tags(),
	}
}
