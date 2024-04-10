// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
)

func resourceKeyVaultCertificate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		// TODO: support Updating additional properties once we have more information about what can be updated
		Create: resourceKeyVaultCertificateCreate,
		Read:   resourceKeyVaultCertificateRead,
		Delete: resourceKeyVaultCertificateDelete,
		Update: resourceKeyVaultCertificateUpdate,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.ParseNestedItemID(id)
			return err
		}, nestedItemResourceImporter),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			// TODO: Change this back to 5min, once https://github.com/hashicorp/terraform-provider-azurerm/issues/11059 is addressed.
			Read:   pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: keyVaultValidate.NestedItemName,
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
				Computed: true,
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
										Type:     pluginsdk.TypeString,
										Optional: true,
										Computed: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(keyvault.JSONWebKeyCurveNameP256),
											string(keyvault.JSONWebKeyCurveNameP256K),
											string(keyvault.JSONWebKeyCurveNameP384),
											string(keyvault.JSONWebKeyCurveNameP521),
										}, false),
									},
									"exportable": {
										Type:     pluginsdk.TypeBool,
										Required: true,
									},
									"key_size": {
										Type:     pluginsdk.TypeInt,
										Optional: true,
										Computed: true,
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
											string(keyvault.JSONWebKeyTypeEC),
											string(keyvault.JSONWebKeyTypeECHSM),
											string(keyvault.JSONWebKeyTypeRSA),
											string(keyvault.JSONWebKeyTypeRSAHSM),
											string(keyvault.JSONWebKeyTypeOct),
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
													Type:     pluginsdk.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(keyvault.CertificatePolicyActionAutoRenew),
														string(keyvault.CertificatePolicyActionEmailContacts),
													}, false),
												},
											},
										},
									},
									//lintignore:XS003
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
							Computed: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"extended_key_usage": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
									"key_usage": {
										Type:     pluginsdk.TypeSet,
										Required: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												string(keyvault.KeyUsageTypeCRLSign),
												string(keyvault.KeyUsageTypeDataEncipherment),
												string(keyvault.KeyUsageTypeDecipherOnly),
												string(keyvault.KeyUsageTypeDigitalSignature),
												string(keyvault.KeyUsageTypeEncipherOnly),
												string(keyvault.KeyUsageTypeKeyAgreement),
												string(keyvault.KeyUsageTypeKeyCertSign),
												string(keyvault.KeyUsageTypeKeyEncipherment),
												string(keyvault.KeyUsageTypeNonRepudiation),
											}, false),
										},
									},
									"subject": {
										Type:     pluginsdk.TypeString,
										Required: true,
									},
									"subject_alternative_names": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Computed: true,
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

			"tags": tags.Schema(),
		},
	}
}

func createCertificate(d *pluginsdk.ResourceData, meta interface{}) (keyvault.CertificateBundle, error) {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return keyvault.CertificateBundle{}, err
	}

	keyVaultBaseUrl, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return keyvault.CertificateBundle{}, fmt.Errorf("looking up Base URI for Certificate %q in %s: %+v", name, *keyVaultId, err)
	}

	t := d.Get("tags").(map[string]interface{})

	policy, err := expandKeyVaultCertificatePolicy(d)
	if err != nil {
		return keyvault.CertificateBundle{}, fmt.Errorf("expanding certificate policy: %s", err)
	}

	parameters := keyvault.CertificateCreateParameters{
		CertificatePolicy: policy,
		Tags:              tags.Expand(t),
	}

	result, err := client.CreateCertificate(ctx, *keyVaultBaseUrl, name, parameters)
	if err != nil {
		return keyvault.CertificateBundle{
			Response: result.Response,
		}, err
	}

	log.Printf("[DEBUG] Waiting for Key Vault Certificate %q in Vault %q to be provisioned", name, *keyVaultBaseUrl)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Provisioning"},
		Target:     []string{"Ready"},
		Refresh:    keyVaultCertificateCreationRefreshFunc(ctx, client, *keyVaultBaseUrl, name),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutCreate),
	}
	// It has been observed that at least one certificate issuer responds to a request with manual processing by issuer staff. SLA's may differ among issuers.
	// The total create timeout duration is divided by a modified poll interval of 30s to calculate the number of times to allow not found instead of the default 20.
	// Using math.Floor, the calculation will err on the lower side of the creation timeout, so as to return before the overall create timeout occurs.
	if policy != nil && policy.IssuerParameters != nil && policy.IssuerParameters.Name != nil && *policy.IssuerParameters.Name != "Self" {
		stateConf.PollInterval = 30 * time.Second
		stateConf.NotFoundChecks = int(math.Floor(float64(stateConf.Timeout) / float64(stateConf.PollInterval)))
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return keyvault.CertificateBundle{}, fmt.Errorf("waiting for Certificate %q in Vault %q to become available: %s", name, *keyVaultBaseUrl, err)
	}
	return client.GetCertificate(ctx, *keyVaultBaseUrl, name, "")
}

func resourceKeyVaultCertificateCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	keyVaultBaseUrl, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("looking up Base URI for Certificate %q in %s: %+v", name, *keyVaultId, err)
	}

	existing, err := client.GetCertificate(ctx, *keyVaultBaseUrl, name, "")
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Certificate %q in %s: %s", name, *keyVaultBaseUrl, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_key_vault_certificate", *existing.ID)
	}

	t := d.Get("tags").(map[string]interface{})
	policy, err := expandKeyVaultCertificatePolicy(d)
	if err != nil {
		return fmt.Errorf("expanding certificate policy: %s", err)
	}

	var newCert keyvault.CertificateBundle
	if v, ok := d.GetOk("certificate"); ok {
		// Import
		certificate := expandKeyVaultCertificate(v)
		importParameters := keyvault.CertificateImportParameters{
			Base64EncodedCertificate: utils.String(certificate.CertificateData),
			Password:                 utils.String(certificate.CertificatePassword),
			CertificatePolicy:        policy,
			Tags:                     tags.Expand(t),
		}
		newCert, err = client.ImportCertificate(ctx, *keyVaultBaseUrl, name, importParameters)
		if err != nil {
			if meta.(*clients.Client).Features.KeyVault.RecoverSoftDeletedCerts && utils.ResponseWasConflict(newCert.Response) {
				if err = recoverDeletedCertificate(ctx, d, meta, *keyVaultBaseUrl, name); err != nil {
					return fmt.Errorf("recover deleted certificate: %+v", err)
				}
				newCert, err = client.ImportCertificate(ctx, *keyVaultBaseUrl, name, importParameters)
				if err != nil {
					return fmt.Errorf("update recovered certificate: %+v", err)
				}
			} else {
				return err
			}
		}
	} else {
		// Generate new
		newCert, err = createCertificate(d, meta)
		if err != nil {
			if meta.(*clients.Client).Features.KeyVault.RecoverSoftDeletedCerts && utils.ResponseWasConflict(newCert.Response) {
				if err = recoverDeletedCertificate(ctx, d, meta, *keyVaultBaseUrl, name); err != nil {
					return fmt.Errorf("recover deleted certificate: %+v", err)
				}
				// after we recovered the existing certificate we still have to apply our changes
				newCert, err = createCertificate(d, meta)
				if err != nil {
					return fmt.Errorf("update recovered certificate: %+v", err)
				}
			} else {
				return err
			}
		}
	}

	certificateId, err := parse.ParseNestedItemID(*newCert.ID)
	if err != nil {
		return err
	}
	d.SetId(certificateId.ID())

	return resourceKeyVaultCertificateRead(d, meta)
}

func recoverDeletedCertificate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}, keyVaultBaseUrl string, name string) error {
	client := meta.(*clients.Client).KeyVault.ManagementClient
	recoveredCertificate, err := client.RecoverDeletedCertificate(ctx, keyVaultBaseUrl, name)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Recovering Secret %q with ID: %q", name, *recoveredCertificate.ID)
	if certificate := recoveredCertificate.ID; certificate != nil {
		stateConf := &pluginsdk.StateChangeConf{
			Pending:                   []string{"pending"},
			Target:                    []string{"available"},
			Refresh:                   keyVaultChildItemRefreshFunc(*certificate),
			Delay:                     30 * time.Second,
			PollInterval:              10 * time.Second,
			ContinuousTargetOccurence: 10,
			Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for Key Vault Secret %q to become available: %s", name, err)
		}
		log.Printf("[DEBUG] Secret %q recovered with ID: %q", name, *recoveredCertificate.ID)
	}
	return nil
}

func resourceKeyVaultCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseNestedItemID(d.Id())
	if err != nil {
		return err
	}

	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	meta.(*clients.Client).KeyVault.AddToCache(*keyVaultId, id.KeyVaultBaseUrl)

	if d.HasChange("certificate") {
		if v, ok := d.GetOk("certificate"); ok {
			// Import new version of certificate
			certificate := expandKeyVaultCertificate(v)
			importParameters := keyvault.CertificateImportParameters{
				Base64EncodedCertificate: utils.String(certificate.CertificateData),
				Password:                 utils.String(certificate.CertificatePassword),
			}
			resp, err := client.ImportCertificate(ctx, id.KeyVaultBaseUrl, id.Name, importParameters)
			if err != nil {
				return err
			}

			if resp.ID == nil {
				return fmt.Errorf("error: Certificate %q in Vault %q get nil ID from server", id.Name, id.KeyVaultBaseUrl)
			}

			certificateId, err := parse.ParseNestedItemID(*resp.ID)
			if err != nil {
				return err
			}
			d.SetId(certificateId.ID())
		}
	}

	// update lifetime_action only should not recreate a certificate
	var lifeTimeOld, lifeTimeNew interface{}
	var policyOld, policyNew map[string]interface{}

	policyOldRaw, policyNewRaw := d.GetChange("certificate_policy")
	policyOldList, policyNewList := policyOldRaw.([]interface{}), policyNewRaw.([]interface{})

	if len(policyOldList) > 0 {
		policyOld = policyOldList[0].(map[string]interface{})
		lifeTimeOld = policyOld["lifetime_action"]
		delete(policyOld, "lifetime_action")
	}
	if len(policyNewList) > 0 {
		policyNew = policyNewList[0].(map[string]interface{})
		lifeTimeNew = policyNew["lifetime_action"]
		delete(policyNew, "lifetime_action")
	}

	// do not recreate cerfiticate when only lifetime_action changes
	if !cmp.Equal(policyNewList, policyOldList) {
		policyNew["lifetime_action"] = lifeTimeNew
		newCert, err := createCertificate(d, meta)
		if err != nil {
			return err
		}
		certificateId, err := parse.ParseNestedItemID(*newCert.ID)
		if err != nil {
			return err
		}
		d.SetId(certificateId.ID())
	}

	if updateLifetime := !cmp.Equal(lifeTimeOld, lifeTimeNew); d.HasChange("tags") || updateLifetime {
		patch := keyvault.CertificateUpdateParameters{}
		if d.HasChange("tags") {
			if t, ok := d.GetOk("tags"); ok {
				patch.Tags = tags.Expand(t.(map[string]interface{}))
			}
		}

		if updateLifetime {
			patch.CertificatePolicy = &keyvault.CertificatePolicy{
				LifetimeActions: expandKeyVaultCertificatePolicyLifetimeAction(lifeTimeNew),
			}
		}

		if _, err = client.UpdateCertificate(ctx, id.KeyVaultBaseUrl, id.Name, "", patch); err != nil {
			return err
		}
	}
	return resourceKeyVaultCertificateRead(d, meta)
}

func keyVaultCertificateCreationRefreshFunc(ctx context.Context, client *keyvault.BaseClient, keyVaultBaseUrl string, name string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		operation, err := client.GetCertificateOperation(ctx, keyVaultBaseUrl, name)
		if err != nil {
			return nil, "", fmt.Errorf("failed to read CertificateOperation in keyVaultCertificateCreationRefreshFunc for Certificate %q in Vault %q: %s", name, keyVaultBaseUrl, err)
		}
		if operation.Status == nil {
			return nil, "", fmt.Errorf("missing status in certificate operation")
		}

		if strings.EqualFold(*operation.Status, "inProgress") {
			if issuer := operation.IssuerParameters; issuer != nil {
				if strings.EqualFold(pointer.From(issuer.Name), "unknown") {
					return operation, "Ready", nil
				}
			}

			return operation, "Provisioning", nil
		}

		if strings.EqualFold(*operation.Status, "completed") {
			return operation, "Ready", nil
		}

		return nil, "", fmt.Errorf("certifcate creation faild in state '%s'", *operation.Status)
	}
}

func resourceKeyVaultCertificateRead(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseNestedItemID(d.Id())
	if err != nil {
		return err
	}

	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultIdRaw == nil {
		log.Printf("[DEBUG] Unable to determine the Resource ID for the Key Vault at URL %q - removing from state!", id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("checking if %s for Certificate %q exists: %v", *keyVaultId, id.Name, err)
	}
	if !ok {
		log.Printf("[DEBUG] Certificate %q was not found in %s - removing from state", id.Name, *keyVaultId)
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

		return fmt.Errorf("reading Key Vault Certificate: %+v", err)
	}

	d.Set("name", id.Name)

	certificatePolicy := flattenKeyVaultCertificatePolicy(cert.Policy, cert.Cer)
	if err := d.Set("certificate_policy", certificatePolicy); err != nil {
		return fmt.Errorf("setting Key Vault Certificate Policy: %+v", err)
	}

	if err := d.Set("certificate_attribute", flattenKeyVaultCertificateAttribute(cert.Attributes)); err != nil {
		return fmt.Errorf("setting Key Vault Certificate Attributes: %+v", err)
	}

	// Computed
	d.Set("version", id.Version)
	d.Set("secret_id", cert.Sid)
	d.Set("versionless_id", id.VersionlessID())

	d.Set("resource_manager_id", parse.NewCertificateID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, id.Name, id.Version).ID())
	d.Set("resource_manager_versionless_id", parse.NewCertificateVersionlessID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, id.Name).ID())

	if cert.Sid != nil {
		secretId, err := parse.ParseNestedItemID(*cert.Sid)
		if err != nil {
			return err
		}
		d.Set("versionless_secret_id", secretId.VersionlessID())
	}

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

func resourceKeyVaultCertificateDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseNestedItemID(d.Id())
	if err != nil {
		return err
	}

	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultIdRaw == nil {
		return fmt.Errorf("Unable to determine the Resource ID for the Key Vault at URL %q", id.KeyVaultBaseUrl)
	}

	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	kv, err := keyVaultsClient.VaultsClient.Get(ctx, *keyVaultId)
	if err != nil {
		if response.WasNotFound(kv.HttpResponse) {
			log.Printf("[DEBUG] Certificate %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("checking if key vault %q for Certificate %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}

	shouldPurge := meta.(*clients.Client).Features.KeyVault.PurgeSoftDeletedCertsOnDestroy
	if shouldPurge && kv.Model != nil && utils.NormaliseNilableBool(kv.Model.Properties.EnablePurgeProtection) {
		log.Printf("[DEBUG] cannot purge certificate %q because vault %q has purge protection enabled", id.Name, keyVaultId.String())
		shouldPurge = false
	}

	description := fmt.Sprintf("Certificate %q (Key Vault %q)", id.Name, id.KeyVaultBaseUrl)
	deleter := deleteAndPurgeCertificate{
		client:      client,
		keyVaultUri: id.KeyVaultBaseUrl,
		name:        id.Name,
	}
	if err := deleteAndOptionallyPurge(ctx, description, shouldPurge, deleter); err != nil {
		return err
	}

	return nil
}

var _ deleteAndPurgeNestedItem = deleteAndPurgeCertificate{}

type deleteAndPurgeCertificate struct {
	client      *keyvault.BaseClient
	keyVaultUri string
	name        string
}

func (d deleteAndPurgeCertificate) DeleteNestedItem(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.DeleteCertificate(ctx, d.keyVaultUri, d.name)
	return resp.Response, err
}

func (d deleteAndPurgeCertificate) NestedItemHasBeenDeleted(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.GetCertificate(ctx, d.keyVaultUri, d.name, "")
	return resp.Response, err
}

func (d deleteAndPurgeCertificate) PurgeNestedItem(ctx context.Context) (autorest.Response, error) {
	return d.client.PurgeDeletedCertificate(ctx, d.keyVaultUri, d.name)
}

func (d deleteAndPurgeCertificate) NestedItemHasBeenPurged(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.GetDeletedCertificate(ctx, d.keyVaultUri, d.name)
	return resp.Response, err
}

func expandKeyVaultCertificatePolicy(d *pluginsdk.ResourceData) (*keyvault.CertificatePolicy, error) {
	policies := d.Get("certificate_policy").([]interface{})
	if len(policies) == 0 || policies[0] == nil {
		return nil, nil
	}

	policyRaw := policies[0].(map[string]interface{})
	policy := keyvault.CertificatePolicy{}

	issuers := policyRaw["issuer_parameters"].([]interface{})
	issuer := issuers[0].(map[string]interface{})
	policy.IssuerParameters = &keyvault.IssuerParameters{
		Name: utils.String(issuer["name"].(string)),
	}

	properties := policyRaw["key_properties"].([]interface{})
	props := properties[0].(map[string]interface{})

	curve := props["curve"].(string)
	keyType := props["key_type"].(string)
	keySize := props["key_size"].(int)

	if keyType == string(keyvault.JSONWebKeyTypeEC) || keyType == string(keyvault.JSONWebKeyTypeECHSM) {
		if curve == "" {
			return nil, fmt.Errorf("`curve` is required when creating an EC key")
		}
		// determine key_size if not specified
		if keySize == 0 {
			switch curve {
			case string(keyvault.JSONWebKeyCurveNameP256), string(keyvault.JSONWebKeyCurveNameP256K):
				keySize = 256
			case string(keyvault.JSONWebKeyCurveNameP384):
				keySize = 384
			case string(keyvault.JSONWebKeyCurveNameP521):
				keySize = 521
			}
		}
	} else if keyType == string(keyvault.JSONWebKeyTypeRSA) || keyType == string(keyvault.JSONWebKeyTypeRSAHSM) {
		if keySize == 0 {
			return nil, fmt.Errorf("`key_size` is required when creating an RSA key")
		}
	}

	policy.KeyProperties = &keyvault.KeyProperties{
		Curve:      keyvault.JSONWebKeyCurveName(curve),
		Exportable: utils.Bool(props["exportable"].(bool)),
		KeySize:    utils.Int32(int32(keySize)),
		KeyType:    keyvault.JSONWebKeyType(keyType),
		ReuseKey:   utils.Bool(props["reuse_key"].(bool)),
	}

	policy.LifetimeActions = expandKeyVaultCertificatePolicyLifetimeAction(policyRaw["lifetime_action"])

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
		keys := cert["key_usage"].(*pluginsdk.Set).List()
		for _, key := range keys {
			keyUsage = append(keyUsage, keyvault.KeyUsageType(key.(string)))
		}

		subjectAlternativeNames := &keyvault.SubjectAlternativeNames{}
		if v, ok := cert["subject_alternative_names"]; ok {
			if sans := v.([]interface{}); len(sans) > 0 {
				if sans[0] != nil {
					san := sans[0].(map[string]interface{})

					emails := san["emails"].(*pluginsdk.Set).List()
					if len(emails) > 0 {
						subjectAlternativeNames.Emails = utils.ExpandStringSlice(emails)
					}

					dnsNames := san["dns_names"].(*pluginsdk.Set).List()
					if len(dnsNames) > 0 {
						subjectAlternativeNames.DNSNames = utils.ExpandStringSlice(dnsNames)
					}

					upns := san["upns"].(*pluginsdk.Set).List()
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

	return &policy, nil
}

func expandKeyVaultCertificatePolicyLifetimeAction(actions interface{}) *[]keyvault.LifetimeAction {
	lifetimeActions := make([]keyvault.LifetimeAction, 0)
	if actions == nil {
		return &lifetimeActions
	}

	for _, v := range actions.([]interface{}) {
		action := v.(map[string]interface{})
		lifetimeAction := keyvault.LifetimeAction{}

		if v, ok := action["action"]; ok {
			as := v.([]interface{})
			a := as[0].(map[string]interface{})
			lifetimeAction.Action = &keyvault.Action{
				ActionType: keyvault.CertificatePolicyAction(a["action_type"].(string)),
			}
		}

		if v, ok := action["trigger"]; ok {
			triggers := v.([]interface{})
			if triggers[0] != nil {
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
		}

		lifetimeActions = append(lifetimeActions, lifetimeAction)
	}
	return &lifetimeActions
}

func flattenKeyVaultCertificatePolicy(input *keyvault.CertificatePolicy, certData *[]byte) []interface{} {
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
		keyProps["curve"] = string(props.Curve)
		keyProps["exportable"] = *props.Exportable
		keyProps["key_size"] = int(*props.KeySize)
		keyProps["key_type"] = string(props.KeyType)
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
		} else if certData != nil && len(*certData) > 0 {
			sanOutput := make(map[string]interface{})
			cert, err := x509.ParseCertificate(*certData)
			if err != nil {
				log.Printf("[DEBUG] Unable to read certificate data: %v", err)
			} else {
				sanOutput["emails"] = set.FromStringSlice(cert.EmailAddresses)
				sanOutput["dns_names"] = set.FromStringSlice(cert.DNSNames)
				sanOutput["upns"] = set.FromStringSlice([]string{})
				sanOutputs = append(sanOutputs, sanOutput)
			}
		}

		certProps["key_usage"] = usages
		certProps["subject"] = ""
		if props.Subject != nil {
			certProps["subject"] = *props.Subject
		}
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
