// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nginx_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxdeployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := nginxconfiguration.ParseConfigurationID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			client := meta.Client.Nginx.NginxConfiguration
			result, err := client.ConfigurationsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return meta.MarkAsGone(id)
				}
				return err
			}

			if result.Model == nil {
				return fmt.Errorf("retrieving %s got nil model", id)
			}

			var output ConfigurationModel
			// protected files content field not return by API so decode from state
			if err := meta.Decode(&output); err != nil {
				return err
			}

			deployID := nginxdeployment.NewNginxDeploymentID(id.SubscriptionId, id.ResourceGroupName, id.NginxDeploymentName)
			output.NginxDeploymentId = deployID.ID()

			if prop := result.Model.Properties; prop != nil {
				output.RootFile = pointer.From(prop.RootFile)

				if prop.Package != nil && prop.Package.Data != nil {
					output.PackageData = pointer.From(prop.Package.Data)
				}

				if files := prop.Files; files != nil {
					configs := []ConfigFile{}
					for _, file := range *files {
						if pointer.From(file.Content) != "" {
							configs = append(configs, ConfigFile{
								Content:     pointer.From(file.Content),
								VirtualPath: pointer.From(file.VirtualPath),
							})
						}
					}
					if len(configs) > 0 {
						output.ConfigFile = configs
					}
				}

				if files := prop.ProtectedFiles; files != nil {
					configs := []ProtectedFile{}
					for _, file := range *files {
						config := ProtectedFile{
							VirtualPath: pointer.From(file.VirtualPath),
							ContentHash: pointer.From(file.ContentHash),
						}
						// GET returns protected files without content, so fill in from state
						for _, protectedFile := range output.ProtectedFile {
							if protectedFile.VirtualPath == pointer.From(file.VirtualPath) {
								config.Content = protectedFile.Content
								break
							}
						}
						configs = append(configs, config)
					}
					if len(configs) > 0 {
						output.ProtectedFile = configs
					}
				}
			}

			return meta.Encode(&output)
		},
	}
}

const defaultConfigurationName = "default"

type ConfigFile struct {
	Content     string `tfschema:"content"`
	VirtualPath string `tfschema:"virtual_path"`
}

func (c ConfigFile) toSDKModel() nginxconfiguration.NginxConfigurationFile {
	return nginxconfiguration.NginxConfigurationFile{
		Content:     pointer.To(c.Content),
		VirtualPath: pointer.To(c.VirtualPath),
	}
}

type ProtectedFile struct {
	Content     string `tfschema:"content"`
	VirtualPath string `tfschema:"virtual_path"`
	ContentHash string `tfschema:"content_hash"`
}

func (c ProtectedFile) toSDKModel() nginxconfiguration.NginxConfigurationProtectedFileRequest {
	return nginxconfiguration.NginxConfigurationProtectedFileRequest{
		Content:     pointer.To(c.Content),
		VirtualPath: pointer.To(c.VirtualPath),
	}
}

func (c ConfigurationModel) toSDKFiles() *[]nginxconfiguration.NginxConfigurationFile {
	files := make([]nginxconfiguration.NginxConfigurationFile, 0, len(c.ConfigFile))
	for _, file := range c.ConfigFile {
		files = append(files, file.toSDKModel())
	}
	return &files
}

func (c ConfigurationModel) toSDKProtectedFiles() *[]nginxconfiguration.NginxConfigurationProtectedFileRequest {
	if len(c.ProtectedFile) == 0 {
		return nil
	}
	files := []nginxconfiguration.NginxConfigurationProtectedFileRequest{}
	for _, file := range c.ProtectedFile {
		files = append(files, file.toSDKModel())
	}
	return &files
}

// ToSDKModel used in both Create and Update
func (c ConfigurationModel) ToSDKModel() nginxconfiguration.NginxConfigurationRequest {
	req := nginxconfiguration.NginxConfigurationRequest{
		Name: pointer.To(defaultConfigurationName),
		Properties: &nginxconfiguration.NginxConfigurationRequestProperties{
			RootFile: pointer.To(c.RootFile),
		},
	}

	req.Properties.Files = c.toSDKFiles()
	req.Properties.ProtectedFiles = c.toSDKProtectedFiles()

	if c.PackageData != "" {
		req.Properties.Package = &nginxconfiguration.NginxConfigurationPackage{
			Data: pointer.To(c.PackageData),
		}
	}

	return req
}
