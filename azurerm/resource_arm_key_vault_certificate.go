package azurerm

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/dataplane/keyvault"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKeyVaultCertificate() *schema.Resource {
	return &schema.Resource{
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
							Type: schema.TypeList,
							// in the case of Import
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
					},
				},
			},

			// Computed
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmKeyVaultCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).keyVaultManagementClient

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
		_, err := client.ImportCertificate(keyVaultBaseUrl, name, importParameters)
		if err != nil {
			return err
		}
	} else {
		// Generate new
		parameters := keyvault.CertificateCreateParameters{
			CertificatePolicy: &policy,
			Tags:              expandTags(tags),
		}
		_, err := client.CreateCertificate(keyVaultBaseUrl, name, parameters)
		if err != nil {
			return err
		}
	}

	resp, err := client.GetCertificate(keyVaultBaseUrl, name, "")
	if err != nil {
		return err
	}

	d.SetId(*resp.ID)

	return resourceArmKeyVaultCertificateRead(d, meta)
}

func resourceArmKeyVaultCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).keyVaultManagementClient

	id, err := parseKeyVaultCertificateID(d.Id())
	if err != nil {
		return err
	}

	cert, err := client.GetCertificate(id.KeyVaultBaseUrl, id.Name, "")

	if err != nil {
		if utils.ResponseWasNotFound(cert.Response) {
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", id.Name)
	d.Set("vault_url", id.KeyVaultBaseUrl)

	certificatePolicy := flattenKeyVaultCertificatePolicy(cert.Policy)
	if err := d.Set("certificate_policy", certificatePolicy); err != nil {
		return fmt.Errorf("Error flattening Key Vault Certificate Policy: %+v", err)
	}

	// Computed
	d.Set("version", id.Version)
	flattenAndSetTags(d, cert.Tags)

	return nil
}

func resourceArmKeyVaultCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).keyVaultManagementClient

	id, err := parseKeyVaultCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.DeleteCertificate(id.KeyVaultBaseUrl, id.Name)
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

	/*
		TODO: is this needed?
		policy.X509CertificateProperties = &keyvault.X509CertificateProperties{
			ValidityInMonths: utils.Int32(12),
			Subject:          utils.String(""),
		}
	*/

	return policy
}

func flattenKeyVaultCertificatePolicy(input *keyvault.CertificatePolicy) []interface{} {
	policy := make(map[string]interface{}, 0)

	if params := input.IssuerParameters; params != nil {
		issuerParams := make(map[string]interface{}, 0)
		issuerParams["name"] = *params.Name
		// TODO: add this field
		//issuerParams["certificate_type"] = params.CertificateType
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
		}
	}
	policy["lifetime_action"] = lifetimeActions

	// secret properties
	if props := input.SecretProperties; props != nil {
		keyProps := make(map[string]interface{}, 0)
		keyProps["content_type"] = *props.ContentType

		policy["secret_properties"] = []interface{}{keyProps}
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

func parseKeyVaultCertificateID(id string) (*KeyVaultCertificate, error) {
	// example: https://tharvey-keyvault.vault.azure.net/certificates/bird/fdf067c93bbb4b22bff4d8b7a9a56217
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse Azure KeyVault Certificate Id: %q", err)
	}

	path := idURL.Path

	path = strings.TrimSpace(path)
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	components := strings.Split(path, "/")

	if len(components) != 3 {
		return nil, fmt.Errorf("Azure KeyVault Certificate Id should have 3 segments, got %d: %q", len(components), path)
	}

	key := KeyVaultCertificate{
		KeyVaultBaseUrl: fmt.Sprintf("%s://%s/", idURL.Scheme, idURL.Host),
		Name:            components[1],
		Version:         components[2],
	}

	return &key, nil
}

type KeyVaultCertificate struct {
	KeyVaultBaseUrl string
	Name            string
	Version         string
}
