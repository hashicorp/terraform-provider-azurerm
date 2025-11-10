// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package signalr

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeCertificate),
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

			keyVaultCertificateId, err := keyvault.ParseNestedItemID(metadata.ResourceData.Get("custom_certificate_id").(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeCertificate)
			if err != nil {
				return fmt.Errorf("parsing custom certificate id error: %+v", err)
			}

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
					KeyVaultBaseUri:    keyVaultCertificateId.KeyVaultBaseURL,
					KeyVaultSecretName: keyVaultCertificateId.Name,
				},
			}

			if certVersion := keyVaultCertificateId.Version; certVersion != nil {
				if customCertSignalrService.CertificateVersion != "" && *certVersion != customCertSignalrService.CertificateVersion {
					return fmt.Errorf("certificate version in cert id is different from `certificate_version`") // TODO: consider deprecating `certificate_version` and enforce versioned key ID input? or if version is optional grab it from the cert ID, otherwise nil?
				}
				customCert.Properties.KeyVaultSecretVersion = certVersion
			}

			if err := client.CustomCertificatesCreateOrUpdateThenPoll(ctx, id, customCert); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
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
			id, err := signalr.ParseCustomCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.CustomCertificatesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving SignalR custom certificate %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			props := resp.Model.Properties

			signalrServiceId := signalr.NewSignalRID(id.SubscriptionId, id.ResourceGroupName, id.SignalRName).ID()

			nestedItem, err := keyvault.NewNestedItemID(props.KeyVaultBaseUri, keyvault.NestedItemTypeCertificate, props.KeyVaultSecretName, props.KeyVaultSecretVersion)
			if err != nil {
				return err
			}

			state := CustomCertSignalrServiceResourceModel{
				Name:               id.CustomCertificateName,
				CustomCertId:       nestedItem.ID(),
				SignalRServiceId:   signalrServiceId,
				CertificateVersion: pointer.From(props.KeyVaultSecretVersion),
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
