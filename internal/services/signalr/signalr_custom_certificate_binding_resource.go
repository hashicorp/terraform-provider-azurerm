package signalr

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2023-02-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CustomCertBindingSignalrModel struct {
	Name               string `tfschema:"name"`
	SignalRServiceId   string `tfschema:"signalr_service_id"`
	CustomCertId       string `tfschema:"custom_certificate_id"`
	CertificateVersion string `tfschema:"certificate_version"`
}

type CustomCertBindingSignalrResource struct{}

var _ sdk.Resource = CustomCertBindingSignalrResource{}

func (r CustomCertBindingSignalrResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"signalr_service_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: signalr.ValidateSignalRID,
		},

		"custom_certificate_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.Any(
				keyVaultValidate.NestedItemId,
				keyVaultValidate.NestedItemIdWithOptionalVersion,
			),
		},

		"certificate_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r CustomCertBindingSignalrResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CustomCertBindingSignalrResource) ModelObject() interface{} {
	return &CustomCertBindingSignalrModel{}
}

func (r CustomCertBindingSignalrResource) ResourceType() string {
	return "azurerm_signalr_custom_certificate_binding"
}

func (r CustomCertBindingSignalrResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var signalRCustomCertBinding CustomCertBindingSignalrModel
			if err := metadata.Decode(&signalRCustomCertBinding); err != nil {
				return err
			}
			client := metadata.Client.SignalR.SignalRClient

			signalRServiceId, err := signalr.ParseSignalRID(metadata.ResourceData.Get("signalr_service_id").(string))
			if err != nil {
				return fmt.Errorf("parsing signalr service id error: %+v", err)
			}

			keyVaultCertificateId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(metadata.ResourceData.Get("custom_certificate_id").(string))
			if err != nil {
				return fmt.Errorf("parsing custom certificate id error: %+v", err)
			}

			keyVaultUri := keyVaultCertificateId.KeyVaultBaseUrl
			keyVaultSecretName := keyVaultCertificateId.Name

			id := signalr.NewCustomCertificateID(signalRServiceId.SubscriptionId, signalRServiceId.ResourceGroupName, signalRServiceId.SignalRName, metadata.ResourceData.Get("name").(string))

			existing, err := client.CustomCertificatesGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing SignalR service custom cert binding error %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			customCert := signalr.CustomCertificate{
				Properties: signalr.CustomCertificateProperties{
					KeyVaultBaseUri:    keyVaultUri,
					KeyVaultSecretName: keyVaultSecretName,
				},
			}

			if certVersion := keyVaultCertificateId.Version; certVersion != "" {
				if signalRCustomCertBinding.CertificateVersion != "" && certVersion != signalRCustomCertBinding.CertificateVersion {
					return fmt.Errorf("certificate version in cert id is different from `certificate_version`")
				}
				customCert.Properties.KeyVaultSecretVersion = utils.String(certVersion)
			}

			if _, err := client.CustomCertificatesCreateOrUpdate(ctx, id, customCert); err != nil {
				return fmt.Errorf("creating signalR custom certificate binding: %s: %+v", id, err)
			}

			stateConf := &pluginsdk.StateChangeConf{
				Pending:    []string{string(signalr.ProvisioningStateCreating), string(signalr.ProvisioningStateUpdating), string(signalr.ProvisioningStateFailed)},
				Target:     []string{string(signalr.ProvisioningStateSucceeded)},
				Refresh:    signalrCustomCertBindingStateRefreshFunc(ctx, client, id),
				MinTimeout: 15 * time.Second,
				Timeout:    30 * time.Minute,
			}

			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r CustomCertBindingSignalrResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SignalR.SignalRClient
			keyVaultClient := metadata.Client.KeyVault
			resourcesClient := metadata.Client.Resource
			id, err := signalr.ParseCustomCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.CustomCertificatesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading SignalR custom certificate binding %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: got nil model", *id)
			}

			vaultBasedUri := resp.Model.Properties.KeyVaultBaseUri
			certName := resp.Model.Properties.KeyVaultSecretName

			keyVaultIdRaw, err := keyVaultClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, vaultBasedUri)
			if err != nil {
				return fmt.Errorf("getting key vault base uri from %s: %+v", id, err)
			}
			vaultId, err := keyVaultParse.VaultID(*keyVaultIdRaw)
			if err != nil {
				return fmt.Errorf("parsing key vault %s: %+v", vaultId, err)
			}

			signalrServiceId := signalr.NewSignalRID(id.SubscriptionId, id.ResourceGroupName, id.SignalRName).ID()

			certVersion := ""
			if resp.Model.Properties.KeyVaultSecretVersion != nil {
				certVersion = *resp.Model.Properties.KeyVaultSecretVersion
			}
			nestedItem, err := keyVaultParse.NewNestedItemID(vaultBasedUri, "certificates", certName, certVersion)
			if err != nil {
				return err
			}

			certId := nestedItem.ID()

			state := CustomCertBindingSignalrModel{
				Name:               id.CustomCertificateName,
				CustomCertId:       certId,
				SignalRServiceId:   signalrServiceId,
				CertificateVersion: utils.NormalizeNilableString(resp.Model.Properties.KeyVaultSecretVersion),
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CustomCertBindingSignalrResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SignalR.SignalRClient

			id, err := signalr.ParseCustomCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var signalRCustomCertBinding CustomCertBindingSignalrModel
			if err := metadata.Decode(&signalRCustomCertBinding); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.CustomCertificatesGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			props := existing.Model
			if props == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if metadata.ResourceData.HasChange("certificate_version") {
				existing.Model.Properties.KeyVaultSecretVersion = utils.String(signalRCustomCertBinding.CertificateVersion)
			}

			if err := client.CustomCertificatesCreateOrUpdateThenPoll(ctx, *id, *props); err != nil {
				return fmt.Errorf("creating signalR custom certificate binding: %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r CustomCertBindingSignalrResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SignalR.SignalRClient

			id, err := signalr.ParseCustomCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			if _, err := client.CustomCertificatesDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r CustomCertBindingSignalrResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return signalr.ValidateCustomCertificateID
}

func signalrCustomCertBindingStateRefreshFunc(ctx context.Context, client *signalr.SignalRClient, id signalr.CustomCertificateId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.CustomCertificatesGet(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id, err)
		}

		if model := resp.Model; model != nil {
			if model.Properties.ProvisioningState != nil {
				return resp, string(*model.Properties.ProvisioningState), nil
			}
		}

		return nil, "", fmt.Errorf("error fetching the custom certificate provisioing state")
	}
}
