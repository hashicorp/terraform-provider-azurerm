// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	kv "github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func expandKeyVaultCertificatePolicy(d *pluginsdk.ResourceData) (*kv.CertificatePolicy, error) {
	policies := d.Get("certificate_policy").([]any)
	if len(policies) == 0 || policies[0] == nil {
		return nil, nil
	}

	policyRaw := policies[0].(map[string]any)
	policy := kv.CertificatePolicy{}

	issuers := policyRaw["issuer_parameters"].([]any)
	issuer := issuers[0].(map[string]any)
	policy.IssuerParameters = &kv.IssuerParameters{
		Name: pointer.To(issuer["name"].(string)),
	}

	properties := policyRaw["key_properties"].([]any)
	props := properties[0].(map[string]any)

	curve := props["curve"].(string)
	keyType := props["key_type"].(string)
	keySize := props["key_size"].(int)

	if keyType == string(kv.JSONWebKeyTypeEC) || keyType == string(kv.JSONWebKeyTypeECHSM) {
		if curve == "" {
			return nil, fmt.Errorf("`curve` is required when creating an EC key")
		}
		// determine key_size if not specified
		if keySize == 0 {
			switch curve {
			case string(kv.JSONWebKeyCurveNameP256), string(kv.JSONWebKeyCurveNameP256K):
				keySize = 256
			case string(kv.JSONWebKeyCurveNameP384):
				keySize = 384
			case string(kv.JSONWebKeyCurveNameP521):
				keySize = 521
			}
		}
	} else if keyType == string(kv.JSONWebKeyTypeRSA) || keyType == string(kv.JSONWebKeyTypeRSAHSM) {
		if keySize == 0 {
			return nil, fmt.Errorf("`key_size` is required when creating an RSA key")
		}
	}

	policy.KeyProperties = &kv.KeyProperties{
		Curve:      kv.JSONWebKeyCurveName(curve),
		Exportable: pointer.To(props["exportable"].(bool)),
		KeySize:    pointer.To(int32(keySize)),
		KeyType:    kv.JSONWebKeyType(keyType),
		ReuseKey:   pointer.To(props["reuse_key"].(bool)),
	}

	policy.LifetimeActions = expandKeyVaultCertificatePolicyLifetimeAction(policyRaw["lifetime_action"])

	secrets := policyRaw["secret_properties"].([]any)
	secret := secrets[0].(map[string]any)
	policy.SecretProperties = &kv.SecretProperties{
		ContentType: pointer.To(secret["content_type"].(string)),
	}

	certificateProperties := policyRaw["x509_certificate_properties"].([]any)
	for _, v := range certificateProperties {
		cert := v.(map[string]any)

		ekus := cert["extended_key_usage"].([]any)
		extendedKeyUsage := pluginsdk.ExpandStringSlice(ekus)

		keyUsage := make([]kv.KeyUsageType, 0)
		keys := cert["key_usage"].(*pluginsdk.Set).List()
		for _, key := range keys {
			keyUsage = append(keyUsage, kv.KeyUsageType(key.(string)))
		}

		subjectAlternativeNames := &kv.SubjectAlternativeNames{}
		if v, ok := cert["subject_alternative_names"]; ok {
			if sans := v.([]any); len(sans) > 0 {
				if sans[0] != nil {
					san := sans[0].(map[string]any)

					emails := san["emails"].(*pluginsdk.Set).List()
					if len(emails) > 0 {
						subjectAlternativeNames.Emails = pluginsdk.ExpandStringSlice(emails)
					}

					dnsNames := san["dns_names"].(*pluginsdk.Set).List()
					if len(dnsNames) > 0 {
						subjectAlternativeNames.DNSNames = pluginsdk.ExpandStringSlice(dnsNames)
					}

					upns := san["upns"].(*pluginsdk.Set).List()
					if len(upns) > 0 {
						subjectAlternativeNames.Upns = pluginsdk.ExpandStringSlice(upns)
					}
				}
			}
		}

		policy.X509CertificateProperties = &kv.X509CertificateProperties{
			ValidityInMonths:        pointer.To(int32(cert["validity_in_months"].(int))),
			Subject:                 pointer.To(cert["subject"].(string)),
			KeyUsage:                &keyUsage,
			Ekus:                    extendedKeyUsage,
			SubjectAlternativeNames: subjectAlternativeNames,
		}
	}

	return &policy, nil
}

func expandKeyVaultCertificatePolicyLifetimeAction(actions any) *[]kv.LifetimeAction {
	lifetimeActions := make([]kv.LifetimeAction, 0)
	if actions == nil {
		return &lifetimeActions
	}

	for _, v := range actions.([]any) {
		action := v.(map[string]any)
		lifetimeAction := kv.LifetimeAction{}

		if v, ok := action["action"]; ok {
			as := v.([]any)
			a := as[0].(map[string]any)
			lifetimeAction.Action = &kv.Action{
				ActionType: kv.CertificatePolicyAction(a["action_type"].(string)),
			}
		}

		if v, ok := action["trigger"]; ok {
			triggers := v.([]any)
			if triggers[0] != nil {
				trigger := triggers[0].(map[string]any)
				lifetimeAction.Trigger = &kv.Trigger{}

				d := trigger["days_before_expiry"].(int)
				if d > 0 {
					lifetimeAction.Trigger.DaysBeforeExpiry = pointer.To(int32(d))
				}

				p := trigger["lifetime_percentage"].(int)
				if p > 0 {
					lifetimeAction.Trigger.LifetimePercentage = pointer.To(int32(p))
				}
			}
		}

		lifetimeActions = append(lifetimeActions, lifetimeAction)
	}
	return &lifetimeActions
}

func expandKeyVaultCertificate(v any) KeyVaultCertificateImportParameters {
	certs := v.([]any)
	cert := certs[0].(map[string]any)

	return KeyVaultCertificateImportParameters{
		CertificateData:     cert["contents"].(string),
		CertificatePassword: cert["password"].(string),
	}
}

func createCertificate(d *pluginsdk.ResourceData, meta any) (kv.CertificateBundle, error) {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return kv.CertificateBundle{}, err
	}

	keyVaultBaseUrl, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return kv.CertificateBundle{}, fmt.Errorf("looking up Base URI for Certificate %q in %s: %+v", name, *keyVaultId, err)
	}

	t := d.Get("tags").(map[string]any)

	policy, err := expandKeyVaultCertificatePolicy(d)
	if err != nil {
		return kv.CertificateBundle{}, fmt.Errorf("expanding certificate policy: %s", err)
	}

	parameters := kv.CertificateCreateParameters{
		CertificatePolicy: policy,
		Tags:              tags.Expand(t),
	}

	result, err := client.CreateCertificate(ctx, *keyVaultBaseUrl, name, parameters)
	if err != nil {
		return kv.CertificateBundle{
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
		return kv.CertificateBundle{}, fmt.Errorf("waiting for Certificate %q in Vault %q to become available: %s", name, *keyVaultBaseUrl, err)
	}
	return client.GetCertificate(ctx, *keyVaultBaseUrl, name, "")
}
