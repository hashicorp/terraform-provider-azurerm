// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/certificates"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/deletedcertificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultCertificateRead(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := keyvault.ParseNestedItemID(d.Id(), keyvault.VersionTypeVersioned, keyvault.NestedItemTypeCertificate)
	if err != nil {
		return err
	}

	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, id.KeyVaultBaseURL)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseURL, err)
	}
	if keyVaultIdRaw == nil {
		log.Printf("[DEBUG] Unable to determine the Resource ID for the Key Vault at URL %q - removing from state!", id.KeyVaultBaseURL)
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

	client := keyVaultsClient.DataPlaneKeyVaultClient.Certificates.Clone(id.KeyVaultBaseURL)
	resp, err := client.GetCertificate(ctx, certificates.NewCertificateversionID(id.KeyVaultBaseURL, id.Name, ""))
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Certificate %q was not found in Key Vault at URI %q - removing from state", id.Name, id.KeyVaultBaseURL)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Key Vault Certificate: %+v", err)
	}

	cert := resp.Model
	if cert == nil {
		return fmt.Errorf("retrieving Certificate %q in Key Vault at URI %q: model was nil", id.Name, id.KeyVaultBaseURL)
	}

	var certBytes []byte
	if cert.Cer != nil {
		certBytes, err = base64.StdEncoding.DecodeString(*cert.Cer)
		if err != nil {
			return fmt.Errorf("decoding certificate data for Certificate %q in Key Vault at URI %q: %+v", id.Name, id.KeyVaultBaseURL, err)
		}
	}

	d.Set("name", id.Name)

	if err := d.Set("certificate_policy", flattenKeyVaultCertificatePolicy(cert.Policy, certBytes)); err != nil {
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
		secretId, err := keyvault.ParseNestedItemID(*cert.Sid, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeSecret)
		if err != nil {
			return err
		}
		d.Set("versionless_secret_id", secretId.VersionlessID())
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

	return tags.FlattenAndSet(d, cert.Tags)
}

func flattenKeyVaultCertificatePolicy(input *certificates.CertificatePolicy, certData []byte) []any {
	if input == nil {
		return []any{}
	}

	policy := make(map[string]any)

	if params := input.Issuer; params != nil {
		issuerParams := make(map[string]any)
		issuerParams["name"] = *params.Name
		policy["issuer_parameters"] = []any{issuerParams}
	}

	// key properties
	if props := input.KeyProps; props != nil {
		keyProps := make(map[string]any)
		keyProps["curve"] = pointer.FromEnum(props.Crv)
		keyProps["exportable"] = *props.Exportable
		keyProps["key_size"] = int(*props.KeySize)
		keyProps["key_type"] = pointer.FromEnum(props.Kty)
		keyProps["reuse_key"] = *props.ReuseKey

		policy["key_properties"] = []any{keyProps}
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
		keyProps := make(map[string]any)
		keyProps["content_type"] = *props.ContentType

		policy["secret_properties"] = []any{keyProps}
	}

	// x509 Certificate Properties
	if props := input.X509Props; props != nil {
		certProps := make(map[string]any)

		usages := make([]string, 0)
		for _, usage := range *props.KeyUsage {
			usages = append(usages, string(usage))
		}

		sanOutputs := make([]any, 0)
		if san := props.Sans; san != nil {
			sanOutput := make(map[string]any)
			if emails := san.Emails; emails != nil {
				sanOutput["emails"] = set.FromStringSlice(*emails)
			}
			if dnsNames := san.DnsNames; dnsNames != nil {
				sanOutput["dns_names"] = set.FromStringSlice(*dnsNames)
			}
			if upns := san.Upns; upns != nil {
				sanOutput["upns"] = set.FromStringSlice(*upns)
			}

			sanOutputs = append(sanOutputs, sanOutput)
		} else if len(certData) > 0 {
			sanOutput := make(map[string]any)
			cert, err := x509.ParseCertificate(certData)
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
		certProps["validity_in_months"] = int(*props.ValidityMonths)
		if props.Ekus != nil {
			certProps["extended_key_usage"] = props.Ekus
		}
		certProps["subject_alternative_names"] = sanOutputs
		policy["x509_certificate_properties"] = []any{certProps}
	}

	return []any{policy}
}

func flattenKeyVaultCertificateAttribute(input *certificates.CertificateAttributes) []any {
	if input == nil {
		return []any{}
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
		created = time.Unix(*input.Created, 0).UTC().Format(time.RFC3339)
	}
	if input.Exp != nil {
		expires = time.Unix(*input.Exp, 0).UTC().Format(time.RFC3339)
	}
	if input.Nbf != nil {
		notBefore = time.Unix(*input.Nbf, 0).UTC().Format(time.RFC3339)
	}
	if input.Updated != nil {
		updated = time.Unix(*input.Updated, 0).UTC().Format(time.RFC3339)
	}
	return []any{
		map[string]any{
			"created":        created,
			"enabled":        enabled,
			"expires":        expires,
			"not_before":     notBefore,
			"recovery_level": pointer.FromEnum(input.RecoveryLevel),
			"updated":        updated,
		},
	}
}

func recoverDeletedCertificate(ctx context.Context, d *pluginsdk.ResourceData, meta any, keyVaultBaseUrl string, name string) error {
	client := meta.(*clients.Client).KeyVault.DataPlaneKeyVaultClient.DeletedCertificates.Clone(keyVaultBaseUrl)
	recoveredCertificate, err := client.RecoverDeletedCertificate(ctx, deletedcertificates.NewDeletedcertificateID(keyVaultBaseUrl, name))
	if err != nil {
		return err
	}
	if recoveredCertificate.Model == nil || recoveredCertificate.Model.Id == nil {
		return fmt.Errorf("recovering Certificate %q in Key Vault at URI %q: unable to determine certificate ID", name, keyVaultBaseUrl)
	}
	certificate := *recoveredCertificate.Model.Id
	log.Printf("[DEBUG] Recovering Secret %q with ID: %q", name, certificate)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"pending"},
		Target:                    []string{"available"},
		Refresh:                   helpers.KeyVaultChildItemRefreshFunc(ctx, certificate),
		Delay:                     30 * time.Second,
		PollInterval:              10 * time.Second,
		ContinuousTargetOccurence: 10,
		Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Key Vault Secret %q to become available: %s", name, err)
	}
	log.Printf("[DEBUG] Secret %q recovered with ID: %q", name, certificate)
	return nil
}

func keyVaultCertificateCreationRefreshFunc(ctx context.Context, client *certificates.CertificatesClient, keyVaultBaseUrl string, name string) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		resp, err := client.GetCertificateOperation(ctx, certificates.NewCertificateID(keyVaultBaseUrl, name))
		if err != nil {
			return nil, "", fmt.Errorf("failed to read CertificateOperation in keyVaultCertificateCreationRefreshFunc for Certificate %q in Vault %q: %s", name, keyVaultBaseUrl, err)
		}
		operation := resp.Model
		if operation == nil || operation.Status == nil {
			return nil, "", fmt.Errorf("missing status in certificate operation")
		}

		if strings.EqualFold(*operation.Status, "inProgress") {
			if issuer := operation.Issuer; issuer != nil {
				if strings.EqualFold(pointer.From(issuer.Name), "unknown") {
					return operation, "Ready", nil
				}
			}

			return operation, "Provisioning", nil
		}

		if strings.EqualFold(*operation.Status, "completed") {
			return operation, "Ready", nil
		}

		return nil, "", fmt.Errorf("certificate creation failed in state '%s'", *operation.Status)
	}
}

var _ helpers.DeleteAndPurgeNestedItem = deleteAndPurgeCertificate{}

type deleteAndPurgeCertificate struct {
	certificatesClient        *certificates.CertificatesClient
	deletedCertificatesClient *deletedcertificates.DeletedCertificatesClient
	keyVaultUri               string
	name                      string
}

func (d deleteAndPurgeCertificate) DeleteNestedItem(ctx context.Context) (autorest.Response, error) {
	resp, err := d.certificatesClient.DeleteCertificate(ctx, certificates.NewCertificateID(d.keyVaultUri, d.name))
	return autorest.Response{Response: resp.HttpResponse}, err
}

func (d deleteAndPurgeCertificate) NestedItemHasBeenDeleted(ctx context.Context) (autorest.Response, error) {
	resp, err := d.certificatesClient.GetCertificate(ctx, certificates.NewCertificateversionID(d.keyVaultUri, d.name, ""))
	return autorest.Response{Response: resp.HttpResponse}, err
}

func (d deleteAndPurgeCertificate) PurgeNestedItem(ctx context.Context) (autorest.Response, error) {
	resp, err := d.deletedCertificatesClient.PurgeDeletedCertificate(ctx, deletedcertificates.NewDeletedcertificateID(d.keyVaultUri, d.name))
	return autorest.Response{Response: resp.HttpResponse}, err
}

func (d deleteAndPurgeCertificate) NestedItemHasBeenPurged(ctx context.Context) (autorest.Response, error) {
	resp, err := d.deletedCertificatesClient.GetDeletedCertificate(ctx, deletedcertificates.NewDeletedcertificateID(d.keyVaultUri, d.name))
	return autorest.Response{Response: resp.HttpResponse}, err
}

type KeyVaultCertificateImportParameters struct {
	CertificateData     string
	CertificatePassword string
}
