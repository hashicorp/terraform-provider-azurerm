// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceKeyVaultCertificateRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			// TODO: Change this back to 5min, once https://github.com/hashicorp/terraform-provider-azurerm/issues/11059 is addressed.
			Read: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: keyvault.ValidateNestedItemName,
			},

			"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),

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
				Optional: true,
				// Note: O+C because Azure returns a version in the ID, which is set into state
				Computed: true,
			},

			// Computed
			"certificate_policy": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"issuer_parameters": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
						"key_properties": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"curve": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"exportable": {
										Type:     pluginsdk.TypeBool,
										Computed: true,
									},
									"key_size": {
										Type:     pluginsdk.TypeInt,
										Computed: true,
									},
									"key_type": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"reuse_key": {
										Type:     pluginsdk.TypeBool,
										Computed: true,
									},
								},
							},
						},

						"lifetime_action": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"action": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"action_type": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
											},
										},
									},
									"trigger": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"days_before_expiry": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"lifetime_percentage": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},

						"secret_properties": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"content_type": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},

						"x509_certificate_properties": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"extended_key_usage": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
									"key_usage": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
									"subject": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"subject_alternative_names": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"emails": {
													Type:     pluginsdk.TypeList,
													Computed: true,
													Elem: &pluginsdk.Schema{
														Type: pluginsdk.TypeString,
													},
												},
												"dns_names": {
													Type:     pluginsdk.TypeList,
													Computed: true,
													Elem: &pluginsdk.Schema{
														Type: pluginsdk.TypeString,
													},
												},
												"upns": {
													Type:     pluginsdk.TypeList,
													Computed: true,
													Elem: &pluginsdk.Schema{
														Type: pluginsdk.TypeString,
													},
												},
											},
										},
									},
									"validity_in_months": {
										Type:     pluginsdk.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
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

			"expires": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"not_before": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceKeyVaultCertificateRead(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}
	version := d.Get("version").(string)

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("looking up base uri for Key %q in %s: %+v", name, *keyVaultId, err)
	}

	client := keyVaultsClient.DataPlaneKeyVaultClient.Certificates.Clone(*keyVaultBaseUri)
	resp, err := client.GetCertificate(ctx, certificates.NewCertificateversionID(*keyVaultBaseUri, name, version))
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("a Certificate named %q was not found in Key Vault at URI %q", name, *keyVaultBaseUri)
		}

		return fmt.Errorf("reading Key Vault Certificate: %+v", err)
	}

	cert := resp.Model
	if cert == nil || cert.Id == nil || *cert.Id == "" {
		return fmt.Errorf("failure reading Key Vault Certificate ID for %q", name)
	}

	id, err := keyvault.ParseNestedItemID(*cert.Id, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeCertificate)
	if err != nil {
		return err
	}
	d.SetId(id.ID())

	d.Set("name", id.Name)

	if err := d.Set("certificate_policy", flattenKeyVaultCertificatePolicyForDataSource(cert.Policy)); err != nil {
		return fmt.Errorf("setting Key Vault Certificate Policy: %+v", err)
	}

	d.Set("version", id.Version)
	d.Set("secret_id", cert.Sid)
	d.Set("versionless_id", id.VersionlessID())

	d.Set("resource_manager_id", parse.NewCertificateID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, id.Name, id.Version).ID())
	d.Set("resource_manager_versionless_id", parse.NewCertificateVersionlessID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, id.Name).ID())

	if cert.Sid != nil {
		secretId, err := keyvault.ParseNestedItemID(*cert.Sid, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeSecret)
		if err != nil {
			return err
		}
		d.Set("versionless_secret_id", secretId.VersionlessID())
	}

	var certBytes []byte
	if cert.Cer != nil {
		certBytes, err = base64.StdEncoding.DecodeString(*cert.Cer)
		if err != nil {
			return fmt.Errorf("decoding certificate data for %q: %+v", name, err)
		}
	}

	certificateData := ""
	if cert.Cer != nil {
		certificateData = strings.ToUpper(hex.EncodeToString(certBytes))
	}
	d.Set("certificate_data", certificateData)

	certificateDataBase64 := ""
	if cert.Cer != nil {
		certificateDataBase64 = base64.StdEncoding.EncodeToString(certBytes)
	}
	d.Set("certificate_data_base64", certificateDataBase64)

	thumbprint := ""
	if v := cert.X5t; v != nil {
		x509Thumbprint, err := base64.RawURLEncoding.DecodeString(*v)
		if err != nil {
			return err
		}

		thumbprint = strings.ToUpper(hex.EncodeToString(x509Thumbprint))
	}
	d.Set("thumbprint", thumbprint)

	expires := ""
	notBefore := ""
	if attributes := cert.Attributes; attributes != nil {
		if attributes.Exp != nil {
			expires = time.Unix(*attributes.Exp, 0).UTC().Format(time.RFC3339)
		}
		if attributes.Nbf != nil {
			notBefore = time.Unix(*attributes.Nbf, 0).UTC().Format(time.RFC3339)
		}
	}
	d.Set("expires", expires)
	d.Set("not_before", notBefore)

	return tags.FlattenAndSet(d, cert.Tags)
}

func flattenKeyVaultCertificatePolicyForDataSource(input *certificates.CertificatePolicy) []any {
	if input == nil {
		return []any{}
	}

	policy := make(map[string]any)

	if params := input.Issuer; params != nil {
		policy["issuer_parameters"] = []any{
			map[string]any{
				"name": pointer.From(params.Name),
			},
		}
	}

	// key properties
	if props := input.KeyProps; props != nil {
		var exportable, reuseKey bool
		var keySize int
		curve := pointer.FromEnum(props.Crv)
		keyType := pointer.FromEnum(props.Kty)
		if props.Exportable != nil {
			exportable = *props.Exportable
		}
		if props.ReuseKey != nil {
			reuseKey = *props.ReuseKey
		}
		if props.KeySize != nil {
			keySize = int(*props.KeySize)
		}

		policy["key_properties"] = []any{
			map[string]any{
				"curve":      curve,
				"exportable": exportable,
				"key_size":   keySize,
				"key_type":   keyType,
				"reuse_key":  reuseKey,
			},
		}
	}

	// lifetime actions
	lifetimeActions := make([]any, 0)
	if actions := input.LifetimeActions; actions != nil {
		for _, action := range *actions {
			lifetimeAction := make(map[string]any)

			actionOutput := make(map[string]any)
			if act := action.Action; act != nil {
				actionOutput["action_type"] = pointer.FromEnum(act.ActionType)
			}
			lifetimeAction["action"] = []any{actionOutput}

			triggerOutput := make(map[string]any)
			if trigger := action.Trigger; trigger != nil {
				if days := trigger.DaysBeforeExpiry; days != nil {
					triggerOutput["days_before_expiry"] = int(*trigger.DaysBeforeExpiry)
				}

				if days := trigger.LifetimePercentage; days != nil {
					triggerOutput["lifetime_percentage"] = int(*trigger.LifetimePercentage)
				}
			}
			lifetimeAction["trigger"] = []any{triggerOutput}
			lifetimeActions = append(lifetimeActions, lifetimeAction)
		}
	}
	policy["lifetime_action"] = lifetimeActions

	// secret properties
	if props := input.SecretProps; props != nil {
		policy["secret_properties"] = []any{
			map[string]any{
				"content_type": pointer.From(props.ContentType),
			},
		}
	}

	// x509 Certificate Properties
	if props := input.X509Props; props != nil {
		var subject string
		var validityInMonths int
		if props.Subject != nil {
			subject = *props.Subject
		}
		if props.ValidityMonths != nil {
			validityInMonths = int(*props.ValidityMonths)
		}

		usages := make([]string, 0)
		if props.KeyUsage != nil {
			for _, usage := range *props.KeyUsage {
				usages = append(usages, string(usage))
			}
		}

		sanOutputs := make([]any, 0)
		if san := props.Sans; san != nil {
			sanOutputs = append(sanOutputs, map[string]any{
				"emails":    pluginsdk.FlattenSlice(san.Emails),
				"dns_names": pluginsdk.FlattenSlice(san.DnsNames),
				"upns":      pluginsdk.FlattenSlice(san.Upns),
			})
		}

		policy["x509_certificate_properties"] = []any{
			map[string]any{
				"key_usage":                 usages,
				"subject":                   subject,
				"validity_in_months":        validityInMonths,
				"extended_key_usage":        pluginsdk.FlattenSlice(props.Ekus),
				"subject_alternative_names": sanOutputs,
			},
		}
	}

	return []any{policy}
}
