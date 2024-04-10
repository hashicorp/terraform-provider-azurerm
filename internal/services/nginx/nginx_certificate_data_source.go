// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nginx

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-01-01-preview/nginxcertificate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-01-01-preview/nginxdeployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type CertificateDataSourceModel struct {
	Name                   string `tfschema:"name"`
	NginxDeploymentId      string `tfschema:"nginx_deployment_id"`
	KeyVirtualPath         string `tfschema:"key_virtual_path"`
	CertificateVirtualPath string `tfschema:"certificate_virtual_path"`
	KeyVaultSecretId       string `tfschema:"key_vault_secret_id"`
}

type CertificateDataSource struct{}

var _ sdk.DataSource = CertificateDataSource{}

func (m CertificateDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"nginx_deployment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: nginxdeployment.ValidateNginxDeploymentID,
		},
	}
}

func (m CertificateDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"key_virtual_path": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"certificate_virtual_path": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"key_vault_secret_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (m CertificateDataSource) ModelObject() interface{} {
	return &CertificateDataSourceModel{}
}

func (m CertificateDataSource) ResourceType() string {
	return "azurerm_nginx_certificate"
}

func (m CertificateDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Nginx.NginxCertificate
			var model CertificateDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}
			deploymentId, err := nginxdeployment.ParseNginxDeploymentID(model.NginxDeploymentId)
			if err != nil {
				return err
			}
			id := nginxcertificate.NewCertificateID(
				deploymentId.SubscriptionId,
				deploymentId.ResourceGroupName,
				deploymentId.NginxDeploymentName,
				model.Name,
			)
			result, err := client.CertificatesGet(ctx, id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			output := CertificateDataSourceModel{
				Name:              id.CertificateName,
				NginxDeploymentId: deploymentId.ID(),
			}

			if model := result.Model; model != nil {
				prop := result.Model.Properties
				output.KeyVirtualPath = pointer.From(prop.KeyVirtualPath)
				output.KeyVaultSecretId = pointer.From(prop.KeyVaultSecretId)
				output.CertificateVirtualPath = pointer.From(prop.CertificateVirtualPath)
			}

			metadata.SetID(id)
			return metadata.Encode(&output)
		},
	}
}
