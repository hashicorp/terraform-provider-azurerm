// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package signalr

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2023-02-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CustomCertSignalrServiceResourceModel struct {
	Name               string `tfschema:"name"`
	SignalRServiceId   string `tfschema:"signalr_service_id"`
	CustomCertId       string `tfschema:"custom_certificate_id"`
	CertificateVersion string `tfschema:"certificate_version"`
}

type CustomCertSignalrServiceResource struct{}

var _ sdk.Resource = CustomCertSignalrServiceResource{}

func (r CustomCertSignalrServiceResource) Arguments() map[string]*pluginsdk.Schema {
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
	}
}

func (r CustomCertSignalrServiceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"certificate_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r CustomCertSignalrServiceResource) ModelObject() interface{} {
	return &CustomCertSignalrServiceResourceModel{}
}

func (r CustomCertSignalrServiceResource) ResourceType() string {
	return "azurerm_signalr_service_custom_certificate"
}

func (r CustomCertSignalrServiceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var customCertSignalrService CustomCertSignalrServiceResourceModel
			if err := metadata.Decode(&customCertSignalrService); err != nil {
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

			locks.ByID(signalRServiceId.ID())
			defer locks.UnlockByID(signalRServiceId.ID())

			existing, err := client.CustomCertificatesGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing SignalR service custom cert error %s: %+v", id, err)
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
				if customCertSignalrService.CertificateVersion != "" && certVersion != customCertSignalrService.CertificateVersion {
					return fmt.Errorf("certificate version in cert id is different from `certificate_version`")
				}
				customCert.Properties.KeyVaultSecretVersion = utils.String(certVersion)
			}

			if err := client.CustomCertificatesCreateOrUpdateThenPoll(ctx, id, customCert); err != nil {
				return fmt.Errorf("creating signalR custom certificate: %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r CustomCertSignalrServiceResource) Read() sdk.ResourceFunc {
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
				return fmt.Errorf("reading SignalR custom certificate %s: %+v", id, err)
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
			vaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
			if err != nil {
				return fmt.Errorf("parsing key vault %s: %+v", vaultId, err)
			}

			signalrServiceId := signalr.NewSignalRID(id.SubscriptionId, id.ResourceGroupName, id.SignalRName).ID()

			certVersion := ""
			if resp.Model.Properties.KeyVaultSecretVersion != nil {
				certVersion = *resp.Model.Properties.KeyVaultSecretVersion
			}
			nestedItem, err := keyVaultParse.NewNestedItemID(vaultBasedUri, keyVaultParse.NestedItemTypeCertificate, certName, certVersion)
			if err != nil {
				return err
			}

			certId := nestedItem.ID()

			state := CustomCertSignalrServiceResourceModel{
				Name:               id.CustomCertificateName,
				CustomCertId:       certId,
				SignalRServiceId:   signalrServiceId,
				CertificateVersion: utils.NormalizeNilableString(resp.Model.Properties.KeyVaultSecretVersion),
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CustomCertSignalrServiceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SignalR.SignalRClient

			id, err := signalr.ParseCustomCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			signalrId := signalr.NewSignalRID(id.SubscriptionId, id.ResourceGroupName, id.SignalRName)

			locks.ByID(signalrId.ID())
			defer locks.UnlockByID(signalrId.ID())

			if _, err := client.CustomCertificatesDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}
			stateConf := &pluginsdk.StateChangeConf{
				Pending:                   []string{"Exists"},
				Target:                    []string{"NotFound"},
				Refresh:                   signalrServiceCustomCertificateDeleteRefreshFunc(ctx, client, *id),
				Timeout:                   time.Until(deadline),
				PollInterval:              10 * time.Second,
				ContinuousTargetOccurence: 20,
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be fully deleted: %+v", *id, err)
			}
			return nil
		},
	}
}

func (r CustomCertSignalrServiceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return signalr.ValidateCustomCertificateID
}

func signalrServiceCustomCertificateDeleteRefreshFunc(ctx context.Context, client *signalr.SignalRClient, id signalr.CustomCertificateId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.CustomCertificatesGet(ctx, id)
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return "NotFound", "NotFound", nil
			}

			return nil, "", fmt.Errorf("checking if %s has been deleted: %+v", id, err)
		}

		return res, "Exists", nil
	}
}
