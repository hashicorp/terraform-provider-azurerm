// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nginx

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-06-01-preview/nginxconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-06-01-preview/nginxdeployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ConfigurationDataSourceModel struct {
	NginxDeploymentId string          `tfschema:"nginx_deployment_id"`
	ConfigFile        []ConfigFile    `tfschema:"config_file"`
	ProtectedFile     []ProtectedFile `tfschema:"protected_file"`
	PackageData       string          `tfschema:"package_data"`
	RootFile          string          `tfschema:"root_file"`
}

type ConfigurationDataSource struct{}

var _ sdk.DataSource = ConfigurationDataSource{}

func (m ConfigurationDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"nginx_deployment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: nginxdeployment.ValidateNginxDeploymentID,
		},
	}
}

func (m ConfigurationDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"config_file": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"content": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"virtual_path": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"protected_file": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"content": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},

					"virtual_path": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"package_data": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"root_file": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (m ConfigurationDataSource) ModelObject() interface{} {
	return &ConfigurationDataSourceModel{}
}

func (m ConfigurationDataSource) ResourceType() string {
	return "azurerm_nginx_configuration"
}

func (m ConfigurationDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Nginx.NginxConfiguration
			var model ConfigurationDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}
			deploymentId, err := nginxdeployment.ParseNginxDeploymentID(model.NginxDeploymentId)
			if err != nil {
				return err
			}
			id := nginxconfiguration.NewConfigurationID(
				deploymentId.SubscriptionId,
				deploymentId.ResourceGroupName,
				deploymentId.NginxDeploymentName,
				defaultConfigurationName,
			)
			result, err := client.ConfigurationsGet(ctx, id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			output := ConfigurationDataSourceModel{
				NginxDeploymentId: deploymentId.ID(),
			}

			if model := result.Model; model != nil {
				prop := result.Model.Properties
				output.RootFile = pointer.From(prop.RootFile)

				if prop.Package != nil && prop.Package.Data != nil {
					output.PackageData = pointer.From(prop.Package.Data)
				}

				if files := prop.Files; files != nil {
					for _, file := range *files {
						output.ConfigFile = append(output.ConfigFile, ConfigFile{
							Content:     pointer.From(file.Content),
							VirtualPath: pointer.From(file.VirtualPath),
						})
					}
				}

				if files := prop.ProtectedFiles; files != nil {
					for _, file := range *files {
						output.ProtectedFile = append(output.ProtectedFile, ProtectedFile{
							Content:     pointer.From(file.Content),
							VirtualPath: pointer.From(file.VirtualPath),
						})
					}
				}
			}

			metadata.SetID(id)
			return metadata.Encode(&output)
		},
	}
}
