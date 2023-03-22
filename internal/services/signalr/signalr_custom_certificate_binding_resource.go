package signalr

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2023-02-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"time"
)

type SignalRCustomCertBindingResource struct{}

type SignalRCustomCertBindingModel struct {
	Name             string `tfschema:"name"`
	SignalRServiceId string `tfschema:"signalR_service_id"`
	CustomCertId     string `tfschema:"custom_cert_id"`
}

func (r SignalRCustomCertBindingResource) Arguments() map[string]*pluginsdk.Schema {
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: keyVaultValidate.NestedItemId,
		},
	}
}

func (r SignalRCustomCertBindingResource) ModelObject() interface{} {
	return &SignalRCustomCertBindingModel{}
}

func (r SignalRCustomCertBindingResource) ResourceType() string {
	return "azurerm_signalr_custom_certificate_binding"
}

func (r SignalRCustomCertBindingResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var signalRCustomCertBinding SignalRCustomCertBindingModel
			if err := metadata.Decode(&signalRCustomCertBinding); err != nil {
				return err
			}
			client := metadata.Client.SignalR.SignalRClient

			signalRServiceId, err := signalr.ParseSignalRID(metadata.ResourceData.Get("signalr_service_id").(string))
			if err != nil {
				return fmt.Errorf("parsing signalr service id error: %+v", err)
			}

			keyVaultCertificateId, err := keyVaultParse.ParseNestedItemID(metadata.ResourceData.Get("custom_certificate_id").(string))
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

			props := signalr.CustomCertificateProperties{
				KeyVaultUri: keyVaultCertificateId.KeyVaultBaseUrl,
			}
		},
	}
}
