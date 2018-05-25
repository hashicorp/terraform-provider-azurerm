package azurerm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKeyVaultCertificate() *schema.Resource {
	return &schema.Resource{
		// TODO: support Updating once we have more information about what can be updated
		Create: resourceArmKeyVaultCertificateCreate,
		Read:   resourceArmKeyVaultCertificateRead,
		Delete: resourceArmKeyVaultCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateKeyVaultChildName,
			},

			"vault_uri": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"certificate": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"contents": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"password": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
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
										ValidateFunc: validateIntInSlice([]int{
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
													ConflictsWith: []string{
														"certificate_policy.0.lifetime_action.0.trigger.0.lifetime_percentage",
													},
												},
												"lifetime_percentage": {
													Type:     schema.TypeInt,
													Optional: true,
													ForceNew: true,
													ConflictsWith: []string{
														"certificate_policy.0.lifetime_action.0.trigger.0.days_before_expiry",
													},
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

			"tags": tagsSchema(),
		},
	}
}

func resourceArmKeyVaultCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).keyVaultManagementClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	keyVaultBaseUrl := d.Get("vault_uri").(string)
	tags := d.Get("tags").(map[string]interface{})

	policy := expandKeyVaultCertificatePolicy(d)

	if v, ok := d.GetOk("certificate"); ok {
		// Import
		certificate := expandKeyVaultCertificate(v)
		importParameters := keyvault.CertificateImportParameters{
			Base64EncodedCertificate: utils.String(certificate.CertificateData),
			Password:                 utils.String(certificate.CertificatePassword),
			CertificatePolicy:        &policy,
			Tags:                     expandTags(tags),
		}
		_, err := client.ImportCertificate(ctx, keyVaultBaseUrl, name, importParameters)
		if err != nil {
			return err
		}
	} else {
		// Generate new
		parameters := keyvault.CertificateCreateParameters{
			CertificatePolicy: &policy,
			Tags:              expandTags(tags),
		}
		_, err := client.CreateCertificate(ctx, keyVaultBaseUrl, name, parameters)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Waiting for Key Vault Certificate %q in Vault %q to be provisioned", name, keyVaultBaseUrl)
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"Provisioning"},
			Target:     []string{"Ready"},
			Refresh:    keyVaultCertificateCreationRefreshFunc(ctx, client, keyVaultBaseUrl, name),
			Timeout:    60 * time.Minute,
			MinTimeout: 15 * time.Second,
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

func keyVaultCertificateCreationRefreshFunc(ctx context.Context, client keyvault.BaseClient, keyVaultBaseUrl string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.GetCertificate(ctx, keyVaultBaseUrl, name, "")
		if err != nil {
			return nil, "", fmt.Errorf("Error issuing read request in keyVaultCertificateCreationRefreshFunc for Certificate %q in Vault %q: %s", name, keyVaultBaseUrl, err)
		}

		if res.Sid == nil || *res.Sid == "" {
			return nil, "Provisioning", nil
		}

		return res, "Ready", nil
	}
}

func resourceArmKeyVaultCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).keyVaultManagementClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseKeyVaultChildID(d.Id())
	if err != nil {
		return err
	}

	cert, err := client.GetCertificate(ctx, id.KeyVaultBaseUrl, id.Name, "")

	if err != nil {
		if utils.ResponseWasNotFound(cert.Response) {
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", id.Name)
	d.Set("vault_uri", id.KeyVaultBaseUrl)

	certificatePolicy := flattenKeyVaultCertificatePolicy(cert.Policy)
	if err := d.Set("certificate_policy", certificatePolicy); err != nil {
		return fmt.Errorf("Error flattening Key Vault Certificate Policy: %+v", err)
	}

	// Computed
	d.Set("version", id.Version)
	d.Set("secret_id", cert.Sid)

	if contents := cert.Cer; contents != nil {
		d.Set("certificate_data", string(*contents))
	}
	flattenAndSetTags(d, cert.Tags)

	return nil
}

func resourceArmKeyVaultCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).keyVaultManagementClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseKeyVaultChildID(d.Id())
	if err != nil {
		return err
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

		keyUsage := make([]keyvault.KeyUsageType, 0)
		keys := cert["key_usage"].([]interface{})
		for _, key := range keys {
			keyUsage = append(keyUsage, keyvault.KeyUsageType(key.(string)))
		}

		policy.X509CertificateProperties = &keyvault.X509CertificateProperties{
			ValidityInMonths: utils.Int32(int32(cert["validity_in_months"].(int))),
			Subject:          utils.String(cert["subject"].(string)),
			KeyUsage:         &keyUsage,
		}
	}

	return policy
}

func flattenKeyVaultCertificatePolicy(input *keyvault.CertificatePolicy) []interface{} {
	policy := make(map[string]interface{}, 0)

	if params := input.IssuerParameters; params != nil {
		issuerParams := make(map[string]interface{}, 0)
		issuerParams["name"] = *params.Name
		policy["issuer_parameters"] = []interface{}{issuerParams}
	}

	// key properties
	if props := input.KeyProperties; props != nil {
		keyProps := make(map[string]interface{}, 0)
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
			lifetimeAction := make(map[string]interface{}, 0)

			actionOutput := make(map[string]interface{}, 0)
			if act := action.Action; act != nil {
				actionOutput["action_type"] = string(act.ActionType)
			}
			lifetimeAction["action"] = []interface{}{actionOutput}

			triggerOutput := make(map[string]interface{}, 0)
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
		keyProps := make(map[string]interface{}, 0)
		keyProps["content_type"] = *props.ContentType

		policy["secret_properties"] = []interface{}{keyProps}
	}

	// x509 Certificate Properties
	if props := input.X509CertificateProperties; props != nil {
		certProps := make(map[string]interface{}, 0)

		usages := make([]string, 0)
		for _, usage := range *props.KeyUsage {
			usages = append(usages, string(usage))
		}

		certProps["key_usage"] = usages
		certProps["subject"] = *props.Subject
		certProps["validity_in_months"] = int(*props.ValidityInMonths)

		policy["x509_certificate_properties"] = []interface{}{certProps}
	}

	return []interface{}{policy}
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
