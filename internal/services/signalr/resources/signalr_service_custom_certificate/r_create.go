// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_service_custom_certificate

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			keyVaultCertificateId, err := keyvault.ParseNestedItemID(metadata.ResourceData.Get("custom_certificate_id").(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeCertificate)
			if err != nil {
				return fmt.Errorf("parsing custom certificate id error: %+v", err)
			}

			keyVaultUri := keyVaultCertificateId.KeyVaultBaseURL
			keyVaultSecretName := keyVaultCertificateId.Name

			id := signalr.NewCustomCertificateID(signalRServiceId.SubscriptionId, signalRServiceId.ResourceGroupName, signalRServiceId.SignalRName, metadata.ResourceData.Get("name").(string))

			locks.ByID(signalRServiceId.ID())
			defer locks.UnlockByID(signalRServiceId.ID())

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.CustomCertificatesGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing SignalR service custom cert error %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			customCert := signalr.CustomCertificate{
				Properties: signalr.CustomCertificateProperties{
					KeyVaultBaseUri:    keyVaultUri,
					KeyVaultSecretName: keyVaultSecretName,
				},
			}

			if certVersion := keyVaultCertificateId.Version; certVersion != "" {
				if customCertSignalrService.CertificateVersion != "" && certVersion != customCertSignalrService.CertificateVersion {
					return errors.New("certificate version in cert id is different from `certificate_version`")
				}
				customCert.Properties.KeyVaultSecretVersion = pointer.To(certVersion)
			}

			if err := client.CustomCertificatesCreateOrUpdateCallbackThenPoll(ctx, id, customCert, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}
