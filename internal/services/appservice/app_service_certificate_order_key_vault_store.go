package appservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/appservicecertificateorders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type CertificateOrderCertificateResource struct{}

type CertificateOrderCertificateModel struct {
	Name               string `tfschema:"name"`
	CertificateOrderId string `tfschema:"certificate_order_id"`
	Location           string `tfschema:"location"`
	KeyVaultId         string `tfschema:"key_vault_id"`
	KeyVaultSecretName string `tfschema:"key_vault_secret_name"`
}

func (r CertificateOrderCertificateResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"certificate_order_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CertificateOrderID,
		},

		"key_vault_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateKeyVaultID,
			// TODO -- remove when issue https://github.com/Azure/azure-rest-api-specs/issues/28498 is addressed
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"key_vault_secret_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: keyVaultValidate.NestedItemName,
		},
	}
}

func (r CertificateOrderCertificateResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r CertificateOrderCertificateResource) ModelObject() interface{} {
	return &CertificateOrderCertificateModel{}
}

func (r CertificateOrderCertificateResource) ResourceType() string {
	return "azurerm_app_service_certificate_order_key_vault_store"
}

func (r CertificateOrderCertificateResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var certificateOrderCertificate CertificateOrderCertificateModel
			if err := metadata.Decode(&certificateOrderCertificate); err != nil {
				return err
			}

			client := metadata.Client.AppService.AppServiceCertificatesOrderClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			certificateOrderId, err := appservicecertificateorders.ParseCertificateOrderID(certificateOrderCertificate.CertificateOrderId)
			if err != nil {
				return err
			}
			id := appservicecertificateorders.NewCertificateID(subscriptionId, certificateOrderId.ResourceGroupName, certificateOrderId.CertificateOrderName, certificateOrderCertificate.Name)

			keyVaultId, err := commonids.ParseKeyVaultID(certificateOrderCertificate.KeyVaultId)
			if err != nil {
				return err
			}

			existing, err := client.GetCertificate(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("retrieving %s: %v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			certOrderCertificate := appservicecertificateorders.AppServiceCertificateResource{
				Name: pointer.To(certificateOrderCertificate.Name),
				Properties: &appservicecertificateorders.AppServiceCertificate{
					KeyVaultId:         pointer.To(keyVaultId.ID()),
					KeyVaultSecretName: pointer.To(certificateOrderCertificate.KeyVaultSecretName),
				},
			}

			if err := client.CreateOrUpdateCertificateThenPoll(ctx, id, certOrderCertificate); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r CertificateOrderCertificateResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.AppServiceCertificatesOrderClient
			id, err := appservicecertificateorders.ParseCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			certificateOrderCertificate, err := client.GetCertificate(ctx, *id)
			if err != nil {
				if response.WasNotFound(certificateOrderCertificate.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := CertificateOrderCertificateModel{
				Name: id.CertificateName,
			}

			certificateOrderId := appservicecertificateorders.NewCertificateOrderID(id.SubscriptionId, id.ResourceGroupName, id.CertificateOrderName)
			state.CertificateOrderId = certificateOrderId.ID()

			// we need to parse the key vault id insensitively as the resource group part was changed https://github.com/Azure/azure-rest-api-specs/issues/new?assignees=&labels=bug&projects=&template=02_bug.yml&title=%5BBUG%5D
			if model := certificateOrderCertificate.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				if props := model.Properties; props != nil {
					if props.KeyVaultId != nil {
						keyVaultId, err := commonids.ParseKeyVaultIDInsensitively(*props.KeyVaultId)
						if err != nil {
							return err
						}
						state.KeyVaultId = keyVaultId.ID()
					}
					state.KeyVaultSecretName = pointer.From(props.KeyVaultSecretName)
				}
			}
			if err := metadata.Encode(&state); err != nil {
				return fmt.Errorf("encoding: %+v", err)
			}

			return nil
		},
	}
}

func (r CertificateOrderCertificateResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := appservicecertificateorders.ParseCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.AppService.AppServiceCertificatesOrderClient

			if _, err := client.DeleteCertificate(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r CertificateOrderCertificateResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := appservicecertificateorders.ParseCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.AppService.AppServiceCertificatesOrderClient

			var state CertificateOrderCertificateModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.GetCertificate(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := *existing.Model

			if metadata.ResourceData.HasChange("key_vault_id") {
				model.Properties.KeyVaultId = pointer.To(state.KeyVaultId)
			}

			if metadata.ResourceData.HasChange("key_vault_secret_name") {
				model.Properties.KeyVaultSecretName = pointer.To(state.KeyVaultSecretName)
			}

			if err := client.CreateOrUpdateCertificateThenPoll(ctx, *id, model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r CertificateOrderCertificateResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return appservicecertificateorders.ValidateCertificateID
}
