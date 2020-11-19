package keyvault

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// todo refactor and find a home for this wayward func
func resourceArmKeyVaultChildResourceImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseKeyVaultChildID(d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("Error Unable to parse ID (%s) for Key Vault Child import: %v", d.Id(), err)
	}

	kvid, err := azure.GetKeyVaultIDFromBaseUrl(ctx, client, id.KeyVaultBaseUrl)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("Error retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}

	d.Set("key_vault_id", kvid)

	return []*schema.ResourceData{d}, nil
}

func resourceArmKeyVaultCertificate() *schema.Resource {
	return &schema.Resource{
		// TODO: support Updating once we have more information about what can be updated
		Create: resourceArmKeyVaultCertificateCreate,
		Read:   resourceArmKeyVaultCertificateRead,
		Delete: resourceArmKeyVaultCertificateDelete,

		Importer: &schema.ResourceImporter{
			State: resourceArmKeyVaultChildResourceImporter,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateKeyVaultChildName,
			},

			"key_vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"certificate": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"contents": {
							Type:      schema.TypeString,
							Required:  true,
							ForceNew:  true,
							Sensitive: true,
						},
						"password": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},
					},
				},
			},

			"certificate_policy": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"issuer_parameters": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
						"key_properties": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"exportable": {
										Type:     schema.TypeBool,
										Required: true,
										ForceNew: true,
									},
									"key_size": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: true,
										ValidateFunc: validation.IntInSlice([]int{
											2048,
											4096,
										}),
									},
									"key_type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"reuse_key": {
										Type:     schema.TypeBool,
										Required: true,
										ForceNew: true,
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
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action_type": {
													Type:     schema.TypeString,
													Required: true,
													ForceNew: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(keyvault.AutoRenew),
														string(keyvault.EmailContacts),
													}, false),
												},
											},
										},
									},
									"trigger": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"days_before_expiry": {
													Type:     schema.TypeInt,
													Optional: true,
													ForceNew: true,
												},
												"lifetime_percentage": {
													Type:     schema.TypeInt,
													Optional: true,
													ForceNew: true,
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
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content_type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},

						"x509_certificate_properties": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"extended_key_usage": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										ForceNew: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
									"key_usage": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												string(keyvault.CRLSign),
												string(keyvault.DataEncipherment),
												string(keyvault.DecipherOnly),
												string(keyvault.DigitalSignature),
												string(keyvault.EncipherOnly),
												string(keyvault.KeyAgreement),
												string(keyvault.KeyCertSign),
												string(keyvault.KeyEncipherment),
												string(keyvault.NonRepudiation),
											}, false),
										},
									},
									"subject": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"subject_alternative_names": {
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: true,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"emails": {
													Type:     schema.TypeSet,
													Optional: true,
													ForceNew: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Set: schema.HashString,
												},
												"dns_names": {
													Type:     schema.TypeSet,
													Optional: true,
													ForceNew: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Set: schema.HashString,
												},
												"upns": {
													Type:     schema.TypeSet,
													Optional: true,
													ForceNew: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Set: schema.HashString,
												},
											},
										},
									},
									"validity_in_months": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},

			// Computed
			"certificate_attribute": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"expires": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"not_before": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"recovery_level": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"updated": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secret_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"certificate_data": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"thumbprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.ForceNewSchema(),
		},
	}
}

func resourceArmKeyVaultCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	vaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId := d.Get("key_vault_id").(string)

	keyVaultBaseUrl, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
	if err != nil {
		return fmt.Errorf("Error looking up Certificate %q vault url from id %q: %+v", name, keyVaultId, err)
	}

	existing, err := client.GetCertificate(ctx, keyVaultBaseUrl, name, "")
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of existing Certificate %q (Key Vault %q): %s", name, keyVaultBaseUrl, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_key_vault_certificate", *existing.ID)
	}

	t := d.Get("tags").(map[string]interface{})
	policy := expandKeyVaultCertificatePolicy(d)

	if v, ok := d.GetOk("certificate"); ok {
		// Import
		certificate := expandKeyVaultCertificate(v)
		importParameters := keyvault.CertificateImportParameters{
			Base64EncodedCertificate: utils.String(certificate.CertificateData),
			Password:                 utils.String(certificate.CertificatePassword),
			CertificatePolicy:        &policy,
			Tags:                     tags.Expand(t),
		}
		if _, err := client.ImportCertificate(ctx, keyVaultBaseUrl, name, importParameters); err != nil {
			return err
		}
	} else {
		// Generate new
		parameters := keyvault.CertificateCreateParameters{
			CertificatePolicy: &policy,
			Tags:              tags.Expand(t),
		}
		if resp, err := client.CreateCertificate(ctx, keyVaultBaseUrl, name, parameters); err != nil {
			if meta.(*clients.Client).Features.KeyVault.RecoverSoftDeletedKeyVaults && utils.ResponseWasConflict(resp.Response) {
				recoveredCertificate, err := client.RecoverDeletedCertificate(ctx, keyVaultBaseUrl, name)
				if err != nil {
					return err
				}
				log.Printf("[DEBUG] Recovering Secret %q with ID: %q", name, *recoveredCertificate.ID)
				if certificate := recoveredCertificate.ID; certificate != nil {
					stateConf := &resource.StateChangeConf{
						Pending:                   []string{"pending"},
						Target:                    []string{"available"},
						Refresh:                   keyVaultChildItemRefreshFunc(*certificate),
						Delay:                     30 * time.Second,
						PollInterval:              10 * time.Second,
						ContinuousTargetOccurence: 10,
						Timeout:                   d.Timeout(schema.TimeoutCreate),
					}

					if _, err := stateConf.WaitForState(); err != nil {
						return fmt.Errorf("Error waiting for Key Vault Secret %q to become available: %s", name, err)
					}
					log.Printf("[DEBUG] Secret %q recovered with ID: %q", name, *recoveredCertificate.ID)
				}
			} else {
				return err
			}
		}

		log.Printf("[DEBUG] Waiting for Key Vault Certificate %q in Vault %q to be provisioned", name, keyVaultBaseUrl)
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"Provisioning"},
			Target:     []string{"Ready"},
			Refresh:    keyVaultCertificateCreationRefreshFunc(ctx, client, keyVaultBaseUrl, name),
			MinTimeout: 15 * time.Second,
			Timeout:    d.Timeout(schema.TimeoutCreate),
		}
		// It has been observed that at least one certificate issuer responds to a request with manual processing by issuer staff. SLA's may differ among issuers.
		// The total create timeout duration is divided by a modified poll interval of 30s to calculate the number of times to allow not found instead of the default 20.
		// Using math.Floor, the calculation will err on the lower side of the creation timeout, so as to return before the overall create timeout occurs.
		if policy.IssuerParameters != nil && policy.IssuerParameters.Name != nil && *policy.IssuerParameters.Name != "Self" {
			stateConf.PollInterval = 30 * time.Second
			stateConf.NotFoundChecks = int(math.Floor(float64(stateConf.Timeout) / float64(stateConf.PollInterval)))
		}

		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("Error waiting for Certificate %q in Vault %q to become available: %s", name, keyVaultBaseUrl, err)
		}
	}

	resp, err := client.GetCertificate(ctx, keyVaultBaseUrl, name, "")
	if err != nil {
		return err
	}

	d.SetId(*resp.ID)

	return resourceArmKeyVaultCertificateRead(d, meta)
}

func keyVaultCertificateCreationRefreshFunc(ctx context.Context, client *keyvault.BaseClient, keyVaultBaseUrl string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.GetCertificate(ctx, keyVaultBaseUrl, name, "")
		if err != nil {
			return nil, "", fmt.Errorf("Error issuing read request in keyVaultCertificateCreationRefreshFunc for Certificate %q in Vault %q: %s", name, keyVaultBaseUrl, err)
		}

		if res.Policy != nil &&
			res.Policy.IssuerParameters != nil &&
			res.Policy.IssuerParameters.Name != nil &&
			strings.EqualFold(*(res.Policy.IssuerParameters.Name), "unknown") {
			return res, "Ready", nil
		}

		if res.Sid == nil || *res.Sid == "" {
			return nil, "Provisioning", nil
		}

		return res, "Ready", nil
	}
}

func resourceArmKeyVaultCertificateRead(d *schema.ResourceData, meta interface{}) error {
	keyVaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseKeyVaultChildID(d.Id())
	if err != nil {
		return err
	}

	keyVaultId, err := azure.GetKeyVaultIDFromBaseUrl(ctx, keyVaultClient, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("Error retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultId == nil {
		log.Printf("[DEBUG] Unable to determine the Resource ID for the Key Vault at URL %q - removing from state!", id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	ok, err := azure.KeyVaultExists(ctx, keyVaultClient, *keyVaultId)
	if err != nil {
		return fmt.Errorf("Error checking if key vault %q for Certificate %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Certificate %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	cert, err := client.GetCertificate(ctx, id.KeyVaultBaseUrl, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(cert.Response) {
			log.Printf("[DEBUG] Certificate %q was not found in Key Vault at URI %q - removing from state", id.Name, id.KeyVaultBaseUrl)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Key Vault Certificate: %+v", err)
	}

	d.Set("name", id.Name)

	certificatePolicy := flattenKeyVaultCertificatePolicy(cert.Policy)
	if err := d.Set("certificate_policy", certificatePolicy); err != nil {
		return fmt.Errorf("Error setting Key Vault Certificate Policy: %+v", err)
	}

	if err := d.Set("certificate_attribute", flattenKeyVaultCertificateAttribute(cert.Attributes)); err != nil {
		return fmt.Errorf("setting Key Vault Certificate Attributes: %+v", err)
	}

	// Computed
	d.Set("version", id.Version)
	d.Set("secret_id", cert.Sid)

	certificateData := ""
	if contents := cert.Cer; contents != nil {
		certificateData = strings.ToUpper(hex.EncodeToString(*contents))
	}
	d.Set("certificate_data", certificateData)

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

func resourceArmKeyVaultCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	keyVaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseKeyVaultChildID(d.Id())
	if err != nil {
		return err
	}

	keyVaultId, err := azure.GetKeyVaultIDFromBaseUrl(ctx, keyVaultClient, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("Error retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultId == nil {
		return fmt.Errorf("Unable to determine the Resource ID for the Key Vault at URL %q", id.KeyVaultBaseUrl)
	}

	ok, err := azure.KeyVaultExists(ctx, keyVaultClient, *keyVaultId)
	if err != nil {
		return fmt.Errorf("Error checking if key vault %q for Certificate %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Certificate %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	resp, err := client.DeleteCertificate(ctx, id.KeyVaultBaseUrl, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}

		return fmt.Errorf("Error deleting Certificate %q from Key Vault: %+v", id.Name, err)
	}

	return nil
}

func expandKeyVaultCertificatePolicy(d *schema.ResourceData) keyvault.CertificatePolicy {
	policies := d.Get("certificate_policy").([]interface{})
	policyRaw := policies[0].(map[string]interface{})
	policy := keyvault.CertificatePolicy{}

	issuers := policyRaw["issuer_parameters"].([]interface{})
	issuer := issuers[0].(map[string]interface{})
	policy.IssuerParameters = &keyvault.IssuerParameters{
		Name: utils.String(issuer["name"].(string)),
	}

	properties := policyRaw["key_properties"].([]interface{})
	props := properties[0].(map[string]interface{})
	policy.KeyProperties = &keyvault.KeyProperties{
		Exportable: utils.Bool(props["exportable"].(bool)),
		KeySize:    utils.Int32(int32(props["key_size"].(int))),
		KeyType:    utils.String(props["key_type"].(string)),
		ReuseKey:   utils.Bool(props["reuse_key"].(bool)),
	}

	lifetimeActions := make([]keyvault.LifetimeAction, 0)
	actions := policyRaw["lifetime_action"].([]interface{})
	for _, v := range actions {
		action := v.(map[string]interface{})
		lifetimeAction := keyvault.LifetimeAction{}

		if v, ok := action["action"]; ok {
			as := v.([]interface{})
			a := as[0].(map[string]interface{})
			lifetimeAction.Action = &keyvault.Action{
				ActionType: keyvault.ActionType(a["action_type"].(string)),
			}
		}

		if v, ok := action["trigger"]; ok {
			triggers := v.([]interface{})
			trigger := triggers[0].(map[string]interface{})
			lifetimeAction.Trigger = &keyvault.Trigger{}

			d := trigger["days_before_expiry"].(int)
			if d > 0 {
				lifetimeAction.Trigger.DaysBeforeExpiry = utils.Int32(int32(d))
			}

			p := trigger["lifetime_percentage"].(int)
			if p > 0 {
				lifetimeAction.Trigger.LifetimePercentage = utils.Int32(int32(p))
			}
		}

		lifetimeActions = append(lifetimeActions, lifetimeAction)
	}
	policy.LifetimeActions = &lifetimeActions

	secrets := policyRaw["secret_properties"].([]interface{})
	secret := secrets[0].(map[string]interface{})
	policy.SecretProperties = &keyvault.SecretProperties{
		ContentType: utils.String(secret["content_type"].(string)),
	}

	certificateProperties := policyRaw["x509_certificate_properties"].([]interface{})
	for _, v := range certificateProperties {
		cert := v.(map[string]interface{})

		ekus := cert["extended_key_usage"].([]interface{})
		extendedKeyUsage := utils.ExpandStringSlice(ekus)

		keyUsage := make([]keyvault.KeyUsageType, 0)
		keys := cert["key_usage"].([]interface{})
		for _, key := range keys {
			keyUsage = append(keyUsage, keyvault.KeyUsageType(key.(string)))
		}

		subjectAlternativeNames := &keyvault.SubjectAlternativeNames{}
		if v, ok := cert["subject_alternative_names"]; ok {
			if sans := v.([]interface{}); len(sans) > 0 {
				if sans[0] != nil {
					san := sans[0].(map[string]interface{})

					emails := san["emails"].(*schema.Set).List()
					if len(emails) > 0 {
						subjectAlternativeNames.Emails = utils.ExpandStringSlice(emails)
					}

					dnsNames := san["dns_names"].(*schema.Set).List()
					if len(dnsNames) > 0 {
						subjectAlternativeNames.DNSNames = utils.ExpandStringSlice(dnsNames)
					}

					upns := san["upns"].(*schema.Set).List()
					if len(upns) > 0 {
						subjectAlternativeNames.Upns = utils.ExpandStringSlice(upns)
					}
				}
			}
		}

		policy.X509CertificateProperties = &keyvault.X509CertificateProperties{
			ValidityInMonths:        utils.Int32(int32(cert["validity_in_months"].(int))),
			Subject:                 utils.String(cert["subject"].(string)),
			KeyUsage:                &keyUsage,
			Ekus:                    extendedKeyUsage,
			SubjectAlternativeNames: subjectAlternativeNames,
		}
	}

	return policy
}

func flattenKeyVaultCertificatePolicy(input *keyvault.CertificatePolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	policy := make(map[string]interface{})

	if params := input.IssuerParameters; params != nil {
		issuerParams := make(map[string]interface{})
		issuerParams["name"] = *params.Name
		policy["issuer_parameters"] = []interface{}{issuerParams}
	}

	// key properties
	if props := input.KeyProperties; props != nil {
		keyProps := make(map[string]interface{})
		keyProps["exportable"] = *props.Exportable
		keyProps["key_size"] = int(*props.KeySize)
		keyProps["key_type"] = *props.KeyType
		keyProps["reuse_key"] = *props.ReuseKey

		policy["key_properties"] = []interface{}{keyProps}
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
		keyProps := make(map[string]interface{})
		keyProps["content_type"] = *props.ContentType

		policy["secret_properties"] = []interface{}{keyProps}
	}

	// x509 Certificate Properties
	if props := input.X509CertificateProperties; props != nil {
		certProps := make(map[string]interface{})

		usages := make([]string, 0)
		for _, usage := range *props.KeyUsage {
			usages = append(usages, string(usage))
		}

		sanOutputs := make([]interface{}, 0)
		if san := props.SubjectAlternativeNames; san != nil {
			sanOutput := make(map[string]interface{})
			if emails := san.Emails; emails != nil {
				sanOutput["emails"] = set.FromStringSlice(*emails)
			}
			if dnsNames := san.DNSNames; dnsNames != nil {
				sanOutput["dns_names"] = set.FromStringSlice(*dnsNames)
			}
			if upns := san.Upns; upns != nil {
				sanOutput["upns"] = set.FromStringSlice(*upns)
			}

			sanOutputs = append(sanOutputs, sanOutput)
		}

		certProps["key_usage"] = usages
		certProps["subject"] = *props.Subject
		certProps["validity_in_months"] = int(*props.ValidityInMonths)
		if props.Ekus != nil {
			certProps["extended_key_usage"] = props.Ekus
		}
		certProps["subject_alternative_names"] = sanOutputs
		policy["x509_certificate_properties"] = []interface{}{certProps}
	}

	return []interface{}{policy}
}

func flattenKeyVaultCertificateAttribute(input *keyvault.CertificateAttributes) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	enabled := false
	created := ""
	expires := ""
	notBefore := ""
	updated := ""
	if input.Enabled != nil {
		enabled = *input.Enabled
	}
	if input.Created != nil {
		created = time.Time(*input.Created).Format(time.RFC3339)
	}
	if input.Expires != nil {
		expires = time.Time(*input.Expires).Format(time.RFC3339)
	}
	if input.NotBefore != nil {
		notBefore = time.Time(*input.NotBefore).Format(time.RFC3339)
	}
	if input.Updated != nil {
		updated = time.Time(*input.Updated).Format(time.RFC3339)
	}
	return []interface{}{
		map[string]interface{}{
			"created":        created,
			"enabled":        enabled,
			"expires":        expires,
			"not_before":     notBefore,
			"recovery_level": string(input.RecoveryLevel),
			"updated":        updated,
		},
	}
}

type KeyVaultCertificateImportParameters struct {
	CertificateData     string
	CertificatePassword string
}

func expandKeyVaultCertificate(v interface{}) KeyVaultCertificateImportParameters {
	certs := v.([]interface{})
	cert := certs[0].(map[string]interface{})

	return KeyVaultCertificateImportParameters{
		CertificateData:     cert["contents"].(string),
		CertificatePassword: cert["password"].(string),
	}
}
