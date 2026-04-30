// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nginx

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2025-11-01/nginxcertificates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2025-11-01/nginxdeployments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type CertificateModel struct {
	Name                   string `tfschema:"name"`
	NginxDeploymentId      string `tfschema:"nginx_deployment_id"`
	KeyVirtualPath         string `tfschema:"key_virtual_path"`
	CertificateVirtualPath string `tfschema:"certificate_virtual_path"`
	KeyVaultSecretId       string `tfschema:"key_vault_secret_id"`
}

type CertificateResource struct{}

var _ sdk.ResourceWithUpdate = (*CertificateResource)(nil)

func (m CertificateResource) Arguments() map[string]*pluginsdk.Schema {
	args := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"nginx_deployment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: nginxdeployments.ValidateNginxDeploymentID,
		},

		"key_virtual_path": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"certificate_virtual_path": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"key_vault_secret_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeSecret),
		},
	}

	if !features.FivePointOh() {
		args["key_vault_secret_id"].ValidateFunc = keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeAny)
	}

	return args
}

func (m CertificateResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (m CertificateResource) ModelObject() interface{} {
	return &CertificateModel{}
}

func (m CertificateResource) ResourceType() string {
	return "azurerm_nginx_certificate"
}

func (m CertificateResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Nginx.NginxCertificates

			var model CertificateModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			deployID, _ := nginxdeployments.ParseNginxDeploymentID(model.NginxDeploymentId)

			subscriptionID := metadata.Client.Account.SubscriptionId
			id := nginxcertificates.NewCertificateID(subscriptionID, deployID.ResourceGroupName, deployID.NginxDeploymentName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.CertificatesGet(ctx, id)
				if !response.WasNotFound(existing.HttpResponse) {
					if err != nil {
						return fmt.Errorf("retreiving %s: %v", id, err)
					}
					return metadata.ResourceRequiresImport(m.ResourceType(), id)
				}
			}

			req := nginxcertificates.NginxCertificate{
				Properties: &nginxcertificates.NginxCertificateProperties{
					CertificateVirtualPath: pointer.To(model.CertificateVirtualPath),
					KeyVaultSecretId:       pointer.To(model.KeyVaultSecretId),
					KeyVirtualPath:         pointer.To(model.KeyVirtualPath),
				},
			}

			if err := client.CertificatesCreateOrUpdateCallbackThenPoll(ctx, id, req, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (m CertificateResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Nginx.NginxCertificates
			id, err := nginxcertificates.ParseCertificateID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model CertificateModel
			if err = meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			// retrieve from GET
			existing, err := client.CertificatesGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving exists when updating: +%v", *id)
			}
			if existing.Model == nil && existing.Model.Properties == nil {
				return fmt.Errorf("retrieving as nil when updating for %v", *id)
			}

			// have to pass all existing properties to update
			upd := existing.Model
			if meta.ResourceData.HasChange("key_virtual_path") {
				upd.Properties.KeyVirtualPath = pointer.To(model.KeyVirtualPath)
			}

			if meta.ResourceData.HasChange("certificate_virtual_path") {
				upd.Properties.CertificateVirtualPath = pointer.To(model.CertificateVirtualPath)
			}

			if meta.ResourceData.HasChange("key_vault_secret_id") {
				upd.Properties.KeyVaultSecretId = pointer.To(model.KeyVaultSecretId)
			}

			err = client.CertificatesCreateOrUpdateThenPoll(ctx, *id, *upd)
			if err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}
			return nil
		},
	}
}

func (m CertificateResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := nginxcertificates.ParseCertificateID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			client := meta.Client.Nginx.NginxCertificates
			result, err := client.CertificatesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return meta.MarkAsGone(id)
				}
				return err
			}

			if result.Model == nil {
				return fmt.Errorf("retrieving %s got nil model", id)
			}
			var output CertificateModel

			output.Name = pointer.From(result.Model.Name)
			output.NginxDeploymentId = nginxdeployments.NewNginxDeploymentID(id.SubscriptionId, id.ResourceGroupName, id.NginxDeploymentName).ID()
			prop := result.Model.Properties
			output.KeyVirtualPath = pointer.From(prop.KeyVirtualPath)
			output.KeyVaultSecretId = pointer.From(prop.KeyVaultSecretId)
			output.CertificateVirtualPath = pointer.From(prop.CertificateVirtualPath)
			return meta.Encode(&output)
		},
	}
}

func (m CertificateResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := nginxcertificates.ParseCertificateID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			meta.Logger.Infof("deleting %s", id)
			client := meta.Client.Nginx.NginxCertificates

			if err := client.CertificatesDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}

			return nil
		},
	}
}

func (m CertificateResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return nginxcertificates.ValidateCertificateID
}
