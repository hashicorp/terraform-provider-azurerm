package managedhsm

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-07-01/managedhsms"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	keyvault2 "github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

type KeyVaultMHSMSecurityDomainResource struct{}

var _ sdk.ResourceWithUpdate = KeyVaultMHSMSecurityDomainResource{}

func (r KeyVaultMHSMSecurityDomainResource) ModelObject() interface{} {
	return &KeyVaultMHSMSecurityDomainResourceSchema{}
}

type KeyVaultMHSMSecurityDomainResourceSchema struct {
	ManagedHSMID   string   `tfschema:"managed_hsm_id"`
	CertificateIds []string `tfschema:"security_domain_key_vault_certificate_ids"`
	Quorum         int      `tfschema:"security_domain_quorum"`
	EncryptedData  string   `tfschema:"security_domain_encrypted_data"`
}

func (r KeyVaultMHSMSecurityDomainResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return managedhsms.ValidateManagedHSMID
}

func (r KeyVaultMHSMSecurityDomainResource) ResourceType() string {
	return "azurerm_key_vault_managed_hardware_security_module_security_domain"
}

func (r KeyVaultMHSMSecurityDomainResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"managed_hsm_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedhsms.ValidateManagedHSMID,
		},

		"security_domain_key_vault_certificate_ids": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 3,
			MaxItems: 10,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeCertificate),
			},
		},

		"security_domain_quorum": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(2, 10),
		},
	}
}

func (r KeyVaultMHSMSecurityDomainResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"security_domain_encrypted_data": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},
	}
}

func (r KeyVaultMHSMSecurityDomainResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			kvClient := metadata.Client.ManagedHSMs
			hsmClient := kvClient.ManagedHsmClient
			keyVaultClient := metadata.Client.KeyVault.ManagementClient

			var config KeyVaultMHSMSecurityDomainResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			managedHsmId, err := managedhsms.ParseManagedHSMID(config.ManagedHSMID)
			if err != nil {
				return err
			}

			resp, err := hsmClient.Get(ctx, *managedHsmId)
			if err != nil || resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.HsmUri == nil {
				return fmt.Errorf("got nil HSMUri for %s: %+v", managedHsmId, err)
			}

			certs := make([]interface{}, len(config.CertificateIds))
			for i, v := range config.CertificateIds {
				certs[i] = v
			}

			encData, err := securityDomainDownload(ctx, kvClient.DataPlaneSecurityDomainsClient, *keyVaultClient, *resp.Model.Properties.HsmUri, certs, config.Quorum)
			if err != nil {
				return fmt.Errorf("downloading security domain for %q: %+v", managedHsmId, err)
			}

			config.EncryptedData = encData

			metadata.SetID(managedHsmId)
			return metadata.Encode(&config)
		},
	}
}

func (r KeyVaultMHSMSecurityDomainResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			hsmClient := metadata.Client.ManagedHSMs.ManagedHsmClient

			managedHsmId, err := managedhsms.ParseManagedHSMID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := hsmClient.Get(ctx, *managedHsmId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*managedHsmId)
				}
				return fmt.Errorf("retrieving %s: %+v", managedHsmId, err)
			}

			return nil
		},
	}
}

// Update doesn't make any changes to the resource, it's only used for re-downloading the security domain with updated
// encryption conditions. e.g. Compromised keys, key rotation, etc.
func (r KeyVaultMHSMSecurityDomainResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			kvClient := metadata.Client.ManagedHSMs
			hsmClient := kvClient.ManagedHsmClient
			keyVaultClient := metadata.Client.KeyVault.ManagementClient

			var config KeyVaultMHSMSecurityDomainResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			managedHsmId, err := managedhsms.ParseManagedHSMID(config.ManagedHSMID)
			if err != nil {
				return err
			}

			resp, err := hsmClient.Get(ctx, *managedHsmId)
			if err != nil || resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.HsmUri == nil {
				return fmt.Errorf("got nil HSMUri for %s: %+v", managedHsmId, err)
			}

			certs := make([]interface{}, len(config.CertificateIds))
			for i, v := range config.CertificateIds {
				certs[i] = v
			}

			encData, err := securityDomainDownload(ctx, kvClient.DataPlaneSecurityDomainsClient, *keyVaultClient, *resp.Model.Properties.HsmUri, certs, config.Quorum)
			if err != nil {
				return fmt.Errorf("downloading security domain for %q: %+v", managedHsmId, err)
			}

			config.EncryptedData = encData
			return metadata.Encode(&config)
		},
	}
}

// Delete is not possible here as there's no _actual_ resource to manage, this resource facilitates a download of a security domain for HSM recovery purposes.
func (r KeyVaultMHSMSecurityDomainResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// No-op for delete, the state is simply removed.
			return nil
		},
	}
}

func securityDomainDownload(ctx context.Context, sdClient *keyvault2.HSMSecurityDomainClient, keyClient keyvault2.BaseClient, vaultBaseUrl string, certIds []interface{}, quorum int) (encDataStr string, err error) {
	var param keyvault2.CertificateInfoObject

	param.Required = pointer.To(int32(quorum))
	certs := make([]keyvault2.SecurityDomainJSONWebKey, 0, len(certIds))
	for _, certID := range certIds {
		certIDStr, ok := certID.(string)
		if !ok {
			continue
		}

		nestedItemType := keyvault.NestedItemTypeCertificate
		if !features.FivePointOh() {
			nestedItemType = keyvault.NestedItemTypeAny
		}

		keyID, err := keyvault.ParseNestedItemID(certIDStr, keyvault.VersionTypeVersioned, nestedItemType)
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
		cert := keyvault2.SecurityDomainJSONWebKey{
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

	err = json.Unmarshal(data, &encData)
	if err != nil {
		return "", fmt.Errorf("unmarshal EncData: %v", err)
	}

	pollerType := custompollers.NewHSMDownloadPoller(sdClient, vaultBaseUrl)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return "", fmt.Errorf("waiting for security domain to download: %+v", err)
	}

	return encData.Value, err
}
