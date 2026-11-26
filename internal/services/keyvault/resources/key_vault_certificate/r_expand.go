// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func expandKeyVaultCertificatePolicy(d *pluginsdk.ResourceData) (*certificates.CertificatePolicy, error) {
	policies := d.Get("certificate_policy").([]any)
	if len(policies) == 0 || policies[0] == nil {
		return nil, nil
	}

	policyRaw := policies[0].(map[string]any)
	policy := certificates.CertificatePolicy{}

	issuers := policyRaw["issuer_parameters"].([]any)
	issuer := issuers[0].(map[string]any)
	policy.Issuer = &certificates.IssuerParameters{
		Name: pointer.To(issuer["name"].(string)),
	}

	properties := policyRaw["key_properties"].([]any)
	props := properties[0].(map[string]any)

	curve := props["curve"].(string)
	keyType := props["key_type"].(string)
	keySize := props["key_size"].(int)

	if keyType == string(certificates.JsonWebKeyTypeEC) || keyType == string(certificates.JsonWebKeyTypeECNegativeHSM) {
		if curve == "" {
			return nil, fmt.Errorf("`curve` is required when creating an EC key")
		}
		// determine key_size if not specified
		if keySize == 0 {
			switch curve {
			case string(certificates.JsonWebKeyCurveNamePNegativeTwoFiveSix), string(certificates.JsonWebKeyCurveNamePNegativeTwoFiveSixK):
				keySize = 256
			case string(certificates.JsonWebKeyCurveNamePNegativeThreeEightFour):
				keySize = 384
			case string(certificates.JsonWebKeyCurveNamePNegativeFiveTwoOne):
				keySize = 521
			}
		}
	} else if keyType == string(certificates.JsonWebKeyTypeRSA) || keyType == string(certificates.JsonWebKeyTypeRSANegativeHSM) {
		if keySize == 0 {
			return nil, fmt.Errorf("`key_size` is required when creating an RSA key")
		}
	}

	policy.KeyProps = &certificates.KeyProperties{
		Exportable: pointer.To(props["exportable"].(bool)),
		KeySize:    pointer.To(int64(keySize)),
		Kty:        pointer.ToEnum[certificates.JsonWebKeyType](keyType),
		ReuseKey:   pointer.To(props["reuse_key"].(bool)),
	}
	if curve != "" {
		policy.KeyProps.Crv = pointer.ToEnum[certificates.JsonWebKeyCurveName](curve)
	}

	policy.LifetimeActions = expandKeyVaultCertificatePolicyLifetimeAction(policyRaw["lifetime_action"])

	secrets := policyRaw["secret_properties"].([]any)
	secret := secrets[0].(map[string]any)
	policy.SecretProps = &certificates.SecretProperties{
		ContentType: pointer.To(secret["content_type"].(string)),
	}

	certificateProperties := policyRaw["x509_certificate_properties"].([]any)
	for _, v := range certificateProperties {
		cert := v.(map[string]any)

		ekus := cert["extended_key_usage"].([]any)
		extendedKeyUsage := pluginsdk.ExpandStringSlice(ekus)

		keyUsage := make([]certificates.KeyUsageType, 0)
		keys := cert["key_usage"].(*pluginsdk.Set).List()
		for _, key := range keys {
			keyUsage = append(keyUsage, certificates.KeyUsageType(key.(string)))
		}

		subjectAlternativeNames := &certificates.SubjectAlternativeNames{}
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
						subjectAlternativeNames.DnsNames = pluginsdk.ExpandStringSlice(dnsNames)
					}

					upns := san["upns"].(*pluginsdk.Set).List()
					if len(upns) > 0 {
						subjectAlternativeNames.Upns = pluginsdk.ExpandStringSlice(upns)
					}
				}
			}
		}

		policy.X509Props = &certificates.X509CertificateProperties{
			ValidityMonths: pointer.To(int64(cert["validity_in_months"].(int))),
			Subject:        pointer.To(cert["subject"].(string)),
			KeyUsage:       &keyUsage,
			Ekus:           extendedKeyUsage,
			Sans:           subjectAlternativeNames,
		}
	}

	return &policy, nil
}

func expandKeyVaultCertificatePolicyLifetimeAction(actions any) *[]certificates.LifetimeAction {
	lifetimeActions := make([]certificates.LifetimeAction, 0)
	if actions == nil {
		return &lifetimeActions
	}

	for _, v := range actions.([]any) {
		action := v.(map[string]any)
		lifetimeAction := certificates.LifetimeAction{}

		if v, ok := action["action"]; ok {
			as := v.([]any)
			a := as[0].(map[string]any)
			lifetimeAction.Action = &certificates.Action{
				ActionType: pointer.ToEnum[certificates.CertificatePolicyAction](a["action_type"].(string)),
			}
		}

		if v, ok := action["trigger"]; ok {
			triggers := v.([]any)
			if triggers[0] != nil {
				trigger := triggers[0].(map[string]any)
				lifetimeAction.Trigger = &certificates.Trigger{}

				d := trigger["days_before_expiry"].(int)
				if d > 0 {
					lifetimeAction.Trigger.DaysBeforeExpiry = pointer.To(int64(d))
				}

				p := trigger["lifetime_percentage"].(int)
				if p > 0 {
					lifetimeAction.Trigger.LifetimePercentage = pointer.To(int64(p))
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

func createCertificate(d *pluginsdk.ResourceData, meta any) (*certificates.CertificateBundle, *http.Response, error) {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return nil, nil, err
	}

	keyVaultBaseUrl, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return nil, nil, fmt.Errorf("looking up Base URI for Certificate %q in %s: %+v", name, *keyVaultId, err)
	}

	client := keyVaultsClient.DataPlaneKeyVaultClient.Certificates.Clone(*keyVaultBaseUrl)

	t := d.Get("tags").(map[string]any)

	policy, err := expandKeyVaultCertificatePolicy(d)
	if err != nil {
		return nil, nil, fmt.Errorf("expanding certificate policy: %s", err)
	}

	parameters := certificates.CertificateCreateParameters{
		Policy: policy,
		Tags:   tags.Expand(t),
	}

	result, err := client.CreateCertificate(ctx, certificates.NewCertificateID(*keyVaultBaseUrl, name), parameters)
	if err != nil {
		return nil, result.HttpResponse, err
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
	if policy != nil && policy.Issuer != nil && policy.Issuer.Name != nil && *policy.Issuer.Name != "Self" {
		stateConf.PollInterval = 30 * time.Second
		stateConf.NotFoundChecks = int(math.Floor(float64(stateConf.Timeout) / float64(stateConf.PollInterval)))
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return nil, nil, fmt.Errorf("waiting for Certificate %q in Vault %q to become available: %s", name, *keyVaultBaseUrl, err)
	}

	resp, err := client.GetCertificate(ctx, certificates.NewCertificateversionID(*keyVaultBaseUrl, name, ""))
	return resp.Model, resp.HttpResponse, err
}
