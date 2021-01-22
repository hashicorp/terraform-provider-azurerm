package keyvault

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceKeyVaultCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeyVaultCertificateRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateKeyVaultChildName,
			},

			"key_vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			// Computed
			"certificate_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"issuer_parameters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"key_properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"exportable": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"key_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"key_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"reuse_key": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},

						"lifetime_action": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"trigger": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"days_before_expiry": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"lifetime_percentage": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},

						"secret_properties": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},

						"x509_certificate_properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"extended_key_usage": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"key_usage": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"subject": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"subject_alternative_names": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"emails": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"dns_names": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"upns": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"validity_in_months": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"secret_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"certificate_data": {
				Type:     schema.TypeString,
				Computed: true,
				Deprecated:	"Deprecated in favour of `certificate_data_base64`",
			},

			"certificate_data_base64": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"thumbprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceKeyVaultCertificateRead(d *schema.ResourceData, meta interface{}) error {
	vaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId := d.Get("key_vault_id").(string)
	version := d.Get("version").(string)

	keyVaultBaseUri, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
	if err != nil {
		return fmt.Errorf("Error looking up Key %q vault url from id %q: %+v", name, keyVaultId, err)
	}

	cert, err := client.GetCertificate(ctx, keyVaultBaseUri, name, version)
	if err != nil {
		if utils.ResponseWasNotFound(cert.Response) {
			log.Printf("[DEBUG] Certificate %q was not found in Key Vault at URI %q - removing from state", name, keyVaultBaseUri)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Key Vault Certificate: %+v", err)
	}

	if cert.ID == nil || *cert.ID == "" {
		return fmt.Errorf("failure reading Key Vault Certificate ID for %q", name)
	}

	d.SetId(*cert.ID)

	id, err := azure.ParseKeyVaultChildID(*cert.ID)
	if err != nil {
		return err
	}

	d.Set("name", id.Name)

	certificatePolicy := flattenKeyVaultCertificatePolicyForDataSource(cert.Policy)
	if err := d.Set("certificate_policy", certificatePolicy); err != nil {
		return fmt.Errorf("Error setting Key Vault Certificate Policy: %+v", err)
	}

	d.Set("version", id.Version)
	d.Set("secret_id", cert.Sid)

	certificateData := ""
	if contents := cert.Cer; contents != nil {
		certificateData = strings.ToUpper(hex.EncodeToString(*contents))
	}
	d.Set("certificate_data", certificateData)

	certificateDataBase64 := ""
	if contents := cert.Cer; contents != nil {
		certificateDataBase64 = base64.StdEncoding.EncodeToString(*contents)
	}
	d.Set("certificate_data_base64", certificateDataBase64)

	thumbprint := ""
	if v := cert.X509Thumbprint; v != nil {
		x509Thumbprint, err := base64.RawURLEncoding.DecodeString(*v)
		if err != nil {
			return err
		}

		thumbprint = strings.ToUpper(hex.EncodeToString(x509Thumbprint))
	}
	d.Set("thumbprint", thumbprint)

	return tags.FlattenAndSet(d, cert.Tags)
}

func flattenKeyVaultCertificatePolicyForDataSource(input *keyvault.CertificatePolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	policy := make(map[string]interface{})

	if params := input.IssuerParameters; params != nil {
		var name string
		if params.Name != nil {
			name = *params.Name
		}
		policy["issuer_parameters"] = []interface{}{
			map[string]interface{}{
				"name": name,
			},
		}
	}

	// key properties
	if props := input.KeyProperties; props != nil {
		var exportable, reuseKey bool
		var keySize int
		var keyType string
		if props.Exportable != nil {
			exportable = *props.Exportable
		}
		if props.ReuseKey != nil {
			reuseKey = *props.ReuseKey
		}
		if props.KeySize != nil {
			keySize = int(*props.KeySize)
		}
		if props.KeyType != nil {
			keyType = *props.KeyType
		}

		policy["key_properties"] = []interface{}{
			map[string]interface{}{
				"exportable": exportable,
				"key_size":   keySize,
				"key_type":   keyType,
				"reuse_key":  reuseKey,
			},
		}
	}

	// lifetime actions
	lifetimeActions := make([]interface{}, 0)
	if actions := input.LifetimeActions; actions != nil {
		for _, action := range *actions {
			lifetimeAction := make(map[string]interface{})

			actionOutput := make(map[string]interface{})
			if act := action.Action; act != nil {
				actionOutput["action_type"] = string(act.ActionType)
			}
			lifetimeAction["action"] = []interface{}{actionOutput}

			triggerOutput := make(map[string]interface{})
			if trigger := action.Trigger; trigger != nil {
				if days := trigger.DaysBeforeExpiry; days != nil {
					triggerOutput["days_before_expiry"] = int(*trigger.DaysBeforeExpiry)
				}

				if days := trigger.LifetimePercentage; days != nil {
					triggerOutput["lifetime_percentage"] = int(*trigger.LifetimePercentage)
				}
			}
			lifetimeAction["trigger"] = []interface{}{triggerOutput}
			lifetimeActions = append(lifetimeActions, lifetimeAction)
		}
	}
	policy["lifetime_action"] = lifetimeActions

	// secret properties
	if props := input.SecretProperties; props != nil {
		var contentType string
		if props.ContentType != nil {
			contentType = *props.ContentType
		}
		policy["secret_properties"] = []interface{}{
			map[string]interface{}{
				"content_type": contentType,
			},
		}
	}

	// x509 Certificate Properties
	if props := input.X509CertificateProperties; props != nil {
		var subject string
		var validityInMonths int
		if props.Subject != nil {
			subject = *props.Subject
		}
		if props.ValidityInMonths != nil {
			validityInMonths = int(*props.ValidityInMonths)
		}

		usages := make([]string, 0)
		if props.KeyUsage != nil {
			for _, usage := range *props.KeyUsage {
				usages = append(usages, string(usage))
			}
		}

		sanOutputs := make([]interface{}, 0)
		if san := props.SubjectAlternativeNames; san != nil {
			sanOutputs = append(sanOutputs, map[string]interface{}{
				"emails":    utils.FlattenStringSlice(san.Emails),
				"dns_names": utils.FlattenStringSlice(san.DNSNames),
				"upns":      utils.FlattenStringSlice(san.Upns),
			})
		}

		policy["x509_certificate_properties"] = []interface{}{
			map[string]interface{}{
				"key_usage":                 usages,
				"subject":                   subject,
				"validity_in_months":        validityInMonths,
				"extended_key_usage":        utils.FlattenStringSlice(props.Ekus),
				"subject_alternative_names": sanOutputs,
			},
		}
	}

	return []interface{}{policy}
}
