// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/managedhsms"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	kv74 "github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func resourceArmKeyVaultManagedHardwareSecurityModuleRead(d *pluginsdk.ResourceData, meta any) error {
	hsmClient := meta.(*clients.Client).ManagedHSMs.ManagedHsmClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managedhsms.ParseManagedHSMID(d.Id())
	if err != nil {
		return err
	}

	resp, err := hsmClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[ERROR] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ManagedHSMName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			d.Set("tenant_id", pointer.From(props.TenantId))
			d.Set("admin_object_ids", pluginsdk.FlattenSlice(props.InitialAdminObjectIds))
			d.Set("hsm_uri", props.HsmUri)
			d.Set("soft_delete_retention_days", props.SoftDeleteRetentionInDays)
			d.Set("purge_protection_enabled", props.EnablePurgeProtection)

			publicAccessEnabled := true
			if props.PublicNetworkAccess != nil && *props.PublicNetworkAccess != managedhsms.PublicNetworkAccessEnabled {
				publicAccessEnabled = false
			}
			d.Set("public_network_access_enabled", publicAccessEnabled)

			if err := d.Set("network_acls", flattenMHSMNetworkAcls(props.NetworkAcls)); err != nil {
				return fmt.Errorf("setting `network_acls`: %+v", err)
			}
		}

		skuName := ""
		if sku := model.Sku; sku != nil {
			skuName = string(sku.Name)
		}
		d.Set("sku_name", skuName)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func flattenMHSMNetworkAcls(acl *managedhsms.MHSMNetworkRuleSet) []any {
	bypass := string(managedhsms.NetworkRuleBypassOptionsAzureServices)
	defaultAction := string(managedhsms.NetworkRuleActionAllow)

	if acl != nil {
		if acl.Bypass != nil {
			bypass = string(*acl.Bypass)
		}
		if acl.DefaultAction != nil {
			defaultAction = string(*acl.DefaultAction)
		}
	}

	return []any{
		map[string]any{
			"bypass":         bypass,
			"default_action": defaultAction,
		},
	}
}

func securityDomainDownload(ctx context.Context, sdClient *kv74.HSMSecurityDomainClient, keyClient kv74.BaseClient, vaultBaseUrl string, certIds []any, quorum int) (encDataStr string, err error) {
	var param kv74.CertificateInfoObject

	param.Required = pointer.To(int32(quorum))
	certs := make([]kv74.SecurityDomainJSONWebKey, 0, len(certIds))
	for _, certID := range certIds {
		certIDStr, ok := certID.(string)
		if !ok {
			continue
		}

		keyID, err := keyvault.ParseNestedItemID(certIDStr, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeCertificate)
		if err != nil {
			return "", fmt.Errorf("parsing %q: %+v", certIDStr, err)
		}
		certRes, err := keyClient.GetCertificate(ctx, keyID.KeyVaultBaseURL, keyID.Name, keyID.Version)
		if err != nil {
			return "", fmt.Errorf("retrieving key %s: %v", certID, err)
		}
		if certRes.Cer == nil {
			return "", fmt.Errorf("got nil key for %s", certID)
		}
		cert := kv74.SecurityDomainJSONWebKey{
			Kty:    pointer.To("RSA"),
			KeyOps: &[]string{""},
			Alg:    pointer.To("RSA-OAEP-256"),
		}
		if certRes.Policy != nil && certRes.Policy.KeyProperties != nil {
			cert.Kty = pointer.To(string(certRes.Policy.KeyProperties.KeyType))
		}
		x5c := ""
		if contents := certRes.Cer; contents != nil {
			x5c = base64.StdEncoding.EncodeToString(*contents)
		}
		cert.X5c = &[]string{x5c}

		sum256 := sha256.Sum256([]byte(x5c))
		s256Dst := make([]byte, base64.StdEncoding.EncodedLen(len(sum256)))
		base64.URLEncoding.Encode(s256Dst, sum256[:])
		cert.X5tS256 = pointer.To(string(s256Dst))
		certs = append(certs, cert)
	}
	param.Certificates = &certs

	future, err := sdClient.Download(ctx, vaultBaseUrl, param)
	if err != nil {
		return "", fmt.Errorf("downloading for %s: %v", vaultBaseUrl, err)
	}

	originResponse := future.Response()
	data, err := io.ReadAll(originResponse.Body)
	if err != nil {
		return "", err
	}
	var encData struct {
		Value string `json:"value"`
	}

	if err = json.Unmarshal(data, &encData); err != nil {
		return "", fmt.Errorf("unmarshal EncData: %v", err)
	}

	pollerType := custompollers.NewHSMDownloadPoller(sdClient, vaultBaseUrl)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return "", fmt.Errorf("waiting for security domain to download: %+v", err)
	}

	return encData.Value, err
}

func keyVaultHSMCustomizeDiff(_ context.Context, d *pluginsdk.ResourceDiff, _ any) error {
	if oldVal, newVal := d.GetChange("security_domain_key_vault_certificate_ids"); len(oldVal.([]any)) != 0 && len(newVal.([]any)) == 0 {
		if err := d.ForceNew("security_domain_key_vault_certificate_ids"); err != nil {
			return err
		}
	}

	if oldVal, newVal := d.GetChange("security_domain_quorum"); oldVal.(int) != 0 && newVal.(int) == 0 {
		if err := d.ForceNew("security_domain_quorum"); err != nil {
			return err
		}
	}

	return nil
}
