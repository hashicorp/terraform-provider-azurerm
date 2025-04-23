// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SpringCloudConfigurationServiceModel struct {
	Name                 string                       `tfschema:"name"`
	SpringCloudServiceId string                       `tfschema:"spring_cloud_service_id"`
	Generation           string                       `tfschema:"generation"`
	RefreshInterval      int64                        `tfschema:"refresh_interval_in_seconds"`
	Repository           []SpringCloudRepositoryModel `tfschema:"repository"`
}

type SpringCloudRepositoryModel struct {
	Name                  string   `tfschema:"name"`
	Label                 string   `tfschema:"label"`
	Patterns              []string `tfschema:"patterns"`
	Uri                   string   `tfschema:"uri"`
	CaCertificateId       string   `tfschema:"ca_certificate_id"`
	HostKey               string   `tfschema:"host_key"`
	HostKeyAlgorithm      string   `tfschema:"host_key_algorithm"`
	Password              string   `tfschema:"password"`
	PrivateKey            string   `tfschema:"private_key"`
	SearchPaths           []string `tfschema:"search_paths"`
	StrictHostKeyChecking bool     `tfschema:"strict_host_key_checking"`
	Username              string   `tfschema:"username"`
}

type SpringCloudConfigurationServiceResource struct{}

func (s SpringCloudConfigurationServiceResource) DeprecationMessage() string {
	return features.DeprecatedInFivePointOh("Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_configuration_service` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.")
}

var (
	_ sdk.ResourceWithUpdate                      = SpringCloudConfigurationServiceResource{}
	_ sdk.ResourceWithStateMigration              = SpringCloudConfigurationServiceResource{}
	_ sdk.ResourceWithDeprecationAndNoReplacement = SpringCloudConfigurationServiceResource{}
)

func (s SpringCloudConfigurationServiceResource) ResourceType() string {
	return "azurerm_spring_cloud_configuration_service"
}

func (s SpringCloudConfigurationServiceResource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"default",
			}, false),
		},

		"spring_cloud_service_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SpringCloudServiceID,
		},

		"generation": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(appplatform.ConfigurationServiceGenerationGenOne),
				string(appplatform.ConfigurationServiceGenerationGenTwo),
			}, false),
		},

		"refresh_interval_in_seconds": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"repository": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"label": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"patterns": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"uri": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"ca_certificate_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.SpringCloudCertificateID,
					},

					"host_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"host_key_algorithm": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"password": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"private_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"search_paths": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"strict_host_key_checking": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"username": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}

func (s SpringCloudConfigurationServiceResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (s SpringCloudConfigurationServiceResource) ModelObject() interface{} {
	return &SpringCloudConfigurationServiceModel{}
}

func (s SpringCloudConfigurationServiceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.SpringCloudConfigurationServiceID
}

func (s SpringCloudConfigurationServiceResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.SpringCloudConfigurationServiceV0ToV1{},
		},
	}
}

func (s SpringCloudConfigurationServiceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudConfigurationServiceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AppPlatform.AppPlatformClient
			springId, err := commonids.ParseSpringCloudServiceID(model.SpringCloudServiceId)
			if err != nil {
				return err
			}
			id := appplatform.NewConfigurationServiceID(springId.SubscriptionId, springId.ResourceGroupName, springId.ServiceName, model.Name)

			existing, err := client.ConfigurationServicesGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(s.ResourceType(), id)
			}

			configurationServiceResource := appplatform.ConfigurationServiceResource{
				Properties: &appplatform.ConfigurationServiceProperties{
					Generation: pointer.To(appplatform.ConfigurationServiceGeneration(model.Generation)),
					Settings: &appplatform.ConfigurationServiceSettings{
						GitProperty: &appplatform.ConfigurationServiceGitProperty{
							Repositories: expandConfigurationServiceConfigurationServiceGitRepositoryArray(model.Repository),
						},
						RefreshIntervalInSeconds: pointer.To(model.RefreshInterval),
					},
				},
			}
			err = client.ConfigurationServicesCreateOrUpdateThenPoll(ctx, id, configurationServiceResource)
			if err != nil {
				return fmt.Errorf("creating/updating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (s SpringCloudConfigurationServiceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudConfigurationServiceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			springId, err := commonids.ParseSpringCloudServiceID(model.SpringCloudServiceId)
			if err != nil {
				return err
			}
			id := appplatform.NewConfigurationServiceID(springId.SubscriptionId, springId.ResourceGroupName, springId.ServiceName, model.Name)

			client := metadata.Client.AppPlatform.AppPlatformClient
			existing, err := client.ConfigurationServicesGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}
			if existing.Model == nil || existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			properties := existing.Model.Properties
			if metadata.ResourceData.HasChange("generation") {
				properties.Generation = pointer.To(appplatform.ConfigurationServiceGeneration(model.Generation))
			}

			if metadata.ResourceData.HasChange("repository") {
				properties.Settings.GitProperty.Repositories = expandConfigurationServiceConfigurationServiceGitRepositoryArray(model.Repository)
			}

			if metadata.ResourceData.HasChange("refresh_interval_in_seconds") {
				properties.Settings.RefreshIntervalInSeconds = pointer.To(model.RefreshInterval)
			}

			configurationServiceResource := appplatform.ConfigurationServiceResource{
				Properties: properties,
			}
			err = client.ConfigurationServicesCreateOrUpdateThenPoll(ctx, id, configurationServiceResource)
			if err != nil {
				return fmt.Errorf("creating/updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (s SpringCloudConfigurationServiceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.AppPlatformClient

			id, err := appplatform.ParseConfigurationServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ConfigurationServicesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			springId := commonids.NewSpringCloudServiceID(id.SubscriptionId, id.ResourceGroupName, id.SpringName)

			var model SpringCloudConfigurationServiceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := SpringCloudConfigurationServiceModel{
				Name:                 id.ConfigurationServiceName,
				SpringCloudServiceId: springId.ID(),
			}

			if resp.Model != nil {
				if props := resp.Model.Properties; props != nil {
					state.Generation = string(pointer.From(props.Generation))
					if props.Settings != nil && props.Settings.GitProperty != nil {
						state.Repository = flattenConfigurationServiceConfigurationServiceGitRepositoryArray(props.Settings.GitProperty.Repositories, model.Repository)
						state.RefreshInterval = pointer.From(props.Settings.RefreshIntervalInSeconds)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (s SpringCloudConfigurationServiceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.AppPlatformClient

			id, err := appplatform.ParseConfigurationServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			err = client.ConfigurationServicesDeleteThenPoll(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandConfigurationServiceConfigurationServiceGitRepositoryArray(input []SpringCloudRepositoryModel) *[]appplatform.ConfigurationServiceGitRepository {
	if len(input) == 0 {
		return nil
	}
	results := make([]appplatform.ConfigurationServiceGitRepository, 0)
	for _, v := range input {
		repo := appplatform.ConfigurationServiceGitRepository{
			Name:                  v.Name,
			Patterns:              v.Patterns,
			Uri:                   v.Uri,
			Label:                 v.Label,
			SearchPaths:           pointer.To(v.SearchPaths),
			Username:              pointer.To(v.Username),
			Password:              pointer.To(v.Password),
			HostKey:               pointer.To(v.HostKey),
			HostKeyAlgorithm:      pointer.To(v.HostKeyAlgorithm),
			PrivateKey:            pointer.To(v.PrivateKey),
			StrictHostKeyChecking: pointer.To(v.StrictHostKeyChecking),
		}
		if v.CaCertificateId != "" {
			repo.CaCertResourceId = pointer.To(v.CaCertificateId)
		}
		results = append(results, repo)
	}
	return &results
}

func flattenConfigurationServiceConfigurationServiceGitRepositoryArray(input *[]appplatform.ConfigurationServiceGitRepository, old []SpringCloudRepositoryModel) []SpringCloudRepositoryModel {
	results := make([]SpringCloudRepositoryModel, 0)
	if input == nil {
		return results
	}

	oldItems := make(map[string]SpringCloudRepositoryModel)
	for _, v := range old {
		oldItems[v.Name] = v
	}

	for _, item := range *input {
		var strictHostKeyChecking bool
		if item.StrictHostKeyChecking != nil {
			strictHostKeyChecking = *item.StrictHostKeyChecking
		}

		var hostKey string
		var hostKeyAlgorithm string
		var privateKey string
		var username string
		var password string
		if oldItem, ok := oldItems[item.Name]; ok {
			hostKey = oldItem.HostKey
			hostKeyAlgorithm = oldItem.HostKeyAlgorithm
			privateKey = oldItem.PrivateKey
			username = oldItem.Username
			password = oldItem.Password
		}

		var caCertificateId string
		if item.CaCertResourceId != nil {
			certificatedId, err := appplatform.ParseCertificateIDInsensitively(*item.CaCertResourceId)
			if err == nil {
				caCertificateId = certificatedId.ID()
			}
		}
		results = append(results, SpringCloudRepositoryModel{
			Name:                  item.Name,
			Label:                 item.Label,
			Patterns:              item.Patterns,
			Uri:                   item.Uri,
			CaCertificateId:       caCertificateId,
			HostKey:               hostKey,
			HostKeyAlgorithm:      hostKeyAlgorithm,
			Password:              password,
			PrivateKey:            privateKey,
			SearchPaths:           pointer.From(item.SearchPaths),
			StrictHostKeyChecking: strictHostKeyChecking,
			Username:              username,
		})
	}
	return results
}
