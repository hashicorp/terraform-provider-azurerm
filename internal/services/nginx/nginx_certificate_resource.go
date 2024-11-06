// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nginx

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-06-01-preview/nginxcertificate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-06-01-preview/nginxdeployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	keyvaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
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
	return map[string]*pluginsdk.Schema{
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
			ValidateFunc: nginxdeployment.ValidateNginxDeploymentID,
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
			ValidateFunc: keyvaultValidate.NestedItemIdWithOptionalVersion,
		},
	}
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
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Nginx.NginxCertificate

			var model CertificateModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			deployID, _ := nginxdeployment.ParseNginxDeploymentID(model.NginxDeploymentId)

			subscriptionID := meta.Client.Account.SubscriptionId
			id := nginxcertificate.NewCertificateID(subscriptionID, deployID.ResourceGroupName, deployID.NginxDeploymentName, model.Name)
			existing, err := client.CertificatesGet(ctx, id)
			if !response.WasNotFound(existing.HttpResponse) {
				if err != nil {
					return fmt.Errorf("retreiving %s: %v", id, err)
				}
				return meta.ResourceRequiresImport(m.ResourceType(), id)
			}

			req := nginxcertificate.NginxCertificate{
				Properties: &nginxcertificate.NginxCertificateProperties{
					CertificateVirtualPath: pointer.FromString(model.CertificateVirtualPath),
					KeyVaultSecretId:       pointer.FromString(model.KeyVaultSecretId),
					KeyVirtualPath:         pointer.FromString(model.KeyVirtualPath),
				},
			}

			err = client.CertificatesCreateOrUpdateThenPoll(ctx, id, req)
			if err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			meta.SetID(id)
			return nil
		},
	}
}

func (m CertificateResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {

			client := meta.Client.Nginx.NginxCertificate
			id, err := nginxcertificate.ParseCertificateID(meta.ResourceData.Id())
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
				upd.Properties.KeyVirtualPath = pointer.FromString(model.KeyVirtualPath)
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
			id, err := nginxcertificate.ParseCertificateID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			client := meta.Client.Nginx.NginxCertificate
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

			output.Name = pointer.ToString(result.Model.Name)
			output.NginxDeploymentId = nginxdeployment.NewNginxDeploymentID(id.SubscriptionId, id.ResourceGroupName, id.NginxDeploymentName).ID()
			prop := result.Model.Properties
			output.KeyVirtualPath = pointer.ToString(prop.KeyVirtualPath)
			output.KeyVaultSecretId = pointer.ToString(prop.KeyVaultSecretId)
			output.CertificateVirtualPath = pointer.ToString(prop.CertificateVirtualPath)
			return meta.Encode(&output)
		},
	}
}

func (m CertificateResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := nginxcertificate.ParseCertificateID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			meta.Logger.Infof("deleting %s", id)
			client := meta.Client.Nginx.NginxCertificate

			if err := client.CertificatesDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}

			return nil
		},
	}
}

func (m CertificateResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return nginxcertificate.ValidateCertificateID
}
