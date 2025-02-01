// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nginx

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxdeployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

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

type ConfigurationModel struct {
	NginxDeploymentId string          `tfschema:"nginx_deployment_id"`
	ConfigFile        []ConfigFile    `tfschema:"config_file"`
	ProtectedFile     []ProtectedFile `tfschema:"protected_file"`
	PackageData       string          `tfschema:"package_data"`
	RootFile          string          `tfschema:"root_file"`
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

type ConfigurationResource struct{}

var _ sdk.Resource = (*ConfigurationResource)(nil)

func (m ConfigurationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"nginx_deployment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: nginxdeployment.ValidateNginxDeploymentID,
		},

		"config_file": {
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			AtLeastOneOf: []string{"config_file", "package_data"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"content": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsBase64,
					},

					"virtual_path": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"protected_file": {
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			RequiredWith: []string{"config_file"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"content": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsBase64,
					},

					"virtual_path": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"content_hash": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"package_data": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			AtLeastOneOf:  []string{"config_file", "package_data"},
			ConflictsWith: []string{"protected_file", "config_file"},
		},

		"root_file": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (m ConfigurationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (m ConfigurationResource) ModelObject() interface{} {
	return &ConfigurationModel{}
}

func (m ConfigurationResource) ResourceType() string {
	return "azurerm_nginx_configuration"
}

func (m ConfigurationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Nginx.NginxConfiguration

			var model ConfigurationModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			deployID, err := nginxdeployment.ParseNginxDeploymentID(model.NginxDeploymentId)
			if err != nil {
				return err
			}

			subscriptionID := meta.Client.Account.SubscriptionId
			id := nginxconfiguration.NewConfigurationID(subscriptionID, deployID.ResourceGroupName, deployID.NginxDeploymentName, defaultConfigurationName)

			existing, err := client.ConfigurationsGet(ctx, id)
			if !response.WasNotFound(existing.HttpResponse) {
				if err != nil {
					return fmt.Errorf("retreiving %s: %v", id, err)
				}
				return meta.ResourceRequiresImport(m.ResourceType(), id)
			}

			req := model.ToSDKModel()

			if err := client.ConfigurationsCreateOrUpdateThenPoll(ctx, id, req); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			meta.SetID(id)
			return nil
		},
	}
}

func (m ConfigurationResource) Read() sdk.ResourceFunc {
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
				output.RootFile = pointer.ToString(prop.RootFile)

				if prop.Package != nil && prop.Package.Data != nil {
					output.PackageData = pointer.ToString(prop.Package.Data)
				}

				if files := prop.Files; files != nil {
					configs := []ConfigFile{}
					for _, file := range *files {
						if pointer.From(file.Content) != "" {
							configs = append(configs, ConfigFile{
								Content:     pointer.ToString(file.Content),
								VirtualPath: pointer.ToString(file.VirtualPath),
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
							VirtualPath: pointer.ToString(file.VirtualPath),
							ContentHash: pointer.ToString(file.ContentHash),
						}
						// GET returns protected files without content, so fill in from state
						for _, protectedFile := range output.ProtectedFile {
							if protectedFile.VirtualPath == pointer.ToString(file.VirtualPath) {
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

func (m ConfigurationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.Nginx.NginxConfiguration
			id, err := nginxconfiguration.ParseConfigurationID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ConfigurationModel
			if err = meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			// retrieve from GET
			existing, err := client.ConfigurationsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving exists: +%v", *id)
			}
			if existing.Model == nil && existing.Model.Properties == nil {
				return fmt.Errorf("retrieving as nil for %v", *id)
			}

			// full update - fill in the existing fields from the API and then patch it
			upd := nginxconfiguration.NginxConfigurationRequest{
				Name: pointer.To(defaultConfigurationName),
				Properties: &nginxconfiguration.NginxConfigurationRequestProperties{
					RootFile: existing.Model.Properties.RootFile,
					Files:    existing.Model.Properties.Files,
					Package:  existing.Model.Properties.Package,
				},
			}

			if existing.Model.Properties.ProtectedFiles != nil {
				var pfs []nginxconfiguration.NginxConfigurationProtectedFileRequest
				for _, f := range *existing.Model.Properties.ProtectedFiles {
					pfs = append(pfs, nginxconfiguration.NginxConfigurationProtectedFileRequest{
						VirtualPath: f.VirtualPath,
					})
				}
				upd.Properties.ProtectedFiles = pointer.To(pfs)
			}

			if meta.ResourceData.HasChange("root_file") {
				upd.Properties.RootFile = pointer.To(model.RootFile)
			}

			if meta.ResourceData.HasChange("config_file") {
				upd.Properties.Files = model.toSDKFiles()
			}

			if meta.ResourceData.HasChange("protected_file") {
				upd.Properties.ProtectedFiles = model.toSDKProtectedFiles()
			}

			if meta.ResourceData.HasChange("package_data") {
				upd.Properties.Package = &nginxconfiguration.NginxConfigurationPackage{
					Data: pointer.To(model.PackageData),
				}
			}

			if err := client.ConfigurationsCreateOrUpdateThenPoll(ctx, *id, upd); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}

func (m ConfigurationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := nginxconfiguration.ParseConfigurationID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			meta.Logger.Infof("deleting %s", id)
			client := meta.Client.Nginx.NginxConfiguration

			if err := client.ConfigurationsDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}

			return nil
		},
	}
}

func (m ConfigurationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return nginxconfiguration.ValidateConfigurationID
}
