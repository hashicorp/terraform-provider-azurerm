package nginx

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2022-08-01/nginxconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2022-08-01/nginxdeployment"
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
		Content:     pointer.FromString(c.Content),
		VirtualPath: pointer.FromString(c.VirtualPath),
	}
}

type ProtectedFile struct {
	Content     string `tfschema:"content"`
	VirtualPath string `tfschema:"virtual_path"`
}

func (c ProtectedFile) toSDKModel() nginxconfiguration.NginxConfigurationFile {
	return nginxconfiguration.NginxConfigurationFile{
		Content:     pointer.FromString(c.Content),
		VirtualPath: pointer.FromString(c.VirtualPath),
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
	var files []nginxconfiguration.NginxConfigurationFile
	for _, file := range c.ConfigFile {
		files = append(files, file.toSDKModel())
	}
	return &files
}

func (c ConfigurationModel) toSDKProtectedFiles() *[]nginxconfiguration.NginxConfigurationFile {
	if len(c.ProtectedFile) == 0 {
		return nil
	}
	var files []nginxconfiguration.NginxConfigurationFile
	for _, file := range c.ProtectedFile {
		files = append(files, file.toSDKModel())
	}
	return &files
}

// ToSDKModel used in both Create and Update
func (c ConfigurationModel) ToSDKModel() nginxconfiguration.NginxConfiguration {
	req := nginxconfiguration.NginxConfiguration{
		Name: pointer.FromString(defaultConfigurationName),
		Properties: &nginxconfiguration.NginxConfigurationProperties{
			RootFile: pointer.FromString(c.RootFile),
		},
	}

	req.Properties.Files = c.toSDKFiles()
	req.Properties.ProtectedFiles = c.toSDKProtectedFiles()

	if c.PackageData != "" {
		req.Properties.Package = &nginxconfiguration.NginxConfigurationPackage{
			Data: pointer.FromString(c.PackageData),
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
			Type:     pluginsdk.TypeSet,
			Required: true,
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
			Type:     pluginsdk.TypeSet,
			Optional: true,
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
				},
			},
		},

		"package_data": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
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
			future, err := client.ConfigurationsCreateOrUpdate(ctx, id, req)
			if err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			if err := future.Poller.PollUntilDone(); err != nil {
				return fmt.Errorf("waiting for creation of %s: %v", id, err)
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
				return err
			}

			if result.Model == nil {
				return fmt.Errorf("retrieving %s got nil model", id)
			}

			var output ConfigurationModel
			deployID := nginxdeployment.NewNginxDeploymentID(id.SubscriptionId, id.ResourceGroupName, id.NginxDeploymentName)
			output.NginxDeploymentId = deployID.ID()

			if prop := result.Model.Properties; prop != nil {
				output.RootFile = pointer.ToString(prop.RootFile)

				if prop.Package != nil && prop.Package.Data != nil {
					output.PackageData = pointer.ToString(prop.Package.Data)
				}

				if files := prop.Files; files != nil {
					for _, file := range *files {
						output.ConfigFile = append(output.ConfigFile, ConfigFile{
							Content:     pointer.ToString(file.Content),
							VirtualPath: pointer.ToString(file.VirtualPath),
						})
					}
				}

				if files := prop.ProtectedFiles; files != nil {
					for _, file := range *files {
						output.ProtectedFile = append(output.ProtectedFile, ProtectedFile{
							Content:     pointer.ToString(file.Content),
							VirtualPath: pointer.ToString(file.VirtualPath),
						})
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

			upd := nginxconfiguration.NginxConfiguration{
				Name:       pointer.FromString(defaultConfigurationName),
				Properties: &nginxconfiguration.NginxConfigurationProperties{},
			}

			// root file is required in update
			upd.Properties.RootFile = pointer.FromString(model.RootFile)

			if meta.ResourceData.HasChange("config_file") {
				upd.Properties.Files = model.toSDKFiles()
			}

			if meta.ResourceData.HasChange("protected_file") {
				upd.Properties.ProtectedFiles = model.toSDKProtectedFiles()
			}

			if meta.ResourceData.HasChange("package_data") {
				upd.Properties.Package = &nginxconfiguration.NginxConfigurationPackage{
					Data: pointer.FromString(model.PackageData),
				}
			}

			result, err := client.ConfigurationsCreateOrUpdate(ctx, *id, upd)
			if err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			if err := result.Poller.PollUntilDone(); err != nil {
				return fmt.Errorf("waiting update %s: %v", *id, err)
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
			result, err := client.ConfigurationsDelete(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}
			if err := result.Poller.PollUntilDone(); err != nil {
				return fmt.Errorf("waiting deleting %s: %v", *id, err)
			}
			return nil
		},
	}
}

func (m ConfigurationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return nginxconfiguration.ValidateConfigurationID
}
