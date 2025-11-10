// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SpringCloudCustomizedAcceleratorModel struct {
	Name                     string               `tfschema:"name"`
	SpringCloudAcceleratorId string               `tfschema:"spring_cloud_accelerator_id"`
	AcceleratorTags          []string             `tfschema:"accelerator_tags"`
	AcceleratorType          string               `tfschema:"accelerator_type"`
	Description              string               `tfschema:"description"`
	DisplayName              string               `tfschema:"display_name"`
	GitRepository            []GitRepositoryModel `tfschema:"git_repository"`
	IconURL                  string               `tfschema:"icon_url"`
}

type GitRepositoryModel struct {
	BasicAuth         []BasicAuthModel `tfschema:"basic_auth"`
	SshAuth           []SshAuthModel   `tfschema:"ssh_auth"`
	Branch            string           `tfschema:"branch"`
	CaCertificateId   string           `tfschema:"ca_certificate_id"`
	Commit            string           `tfschema:"commit"`
	GitTag            string           `tfschema:"git_tag"`
	IntervalInSeconds int64            `tfschema:"interval_in_seconds"`
	Url               string           `tfschema:"url"`
	Path              string           `tfschema:"path"`
}

type BasicAuthModel struct {
	Username string `tfschema:"username"`
	Password string `tfschema:"password"`
}

type SshAuthModel struct {
	PrivateKey          string `tfschema:"private_key"`
	HostKey             string `tfschema:"host_key"`
	PrivateKeyAlgorithm string `tfschema:"host_key_algorithm"`
}

type SpringCloudCustomizedAcceleratorResource struct{}

func (s SpringCloudCustomizedAcceleratorResource) DeprecationMessage() string {
	return features.DeprecatedInFivePointOh("Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_customized_accelerator` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.")
}

var (
	_ sdk.ResourceWithUpdate                      = SpringCloudCustomizedAcceleratorResource{}
	_ sdk.ResourceWithStateMigration              = SpringCloudCustomizedAcceleratorResource{}
	_ sdk.ResourceWithDeprecationAndNoReplacement = SpringCloudCustomizedAcceleratorResource{}
)

func (s SpringCloudCustomizedAcceleratorResource) ResourceType() string {
	return "azurerm_spring_cloud_customized_accelerator"
}

func (s SpringCloudCustomizedAcceleratorResource) ModelObject() interface{} {
	return &SpringCloudCustomizedAcceleratorModel{}
}

func (s SpringCloudCustomizedAcceleratorResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return appplatform.ValidateCustomizedAcceleratorID
}

func (s SpringCloudCustomizedAcceleratorResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.SpringCloudCustomizedAcceleratorV0ToV1{},
		},
	}
}

func (s SpringCloudCustomizedAcceleratorResource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"spring_cloud_accelerator_id": commonschema.ResourceIDReferenceRequiredForceNew(&appplatform.ApplicationAcceleratorId{}),

		"git_repository": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"url": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"basic_auth": {
						Type:          pluginsdk.TypeList,
						Optional:      true,
						ForceNew:      true,
						MaxItems:      1,
						ConflictsWith: []string{"git_repository.0.ssh_auth"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"username": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"password": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"ssh_auth": {
						Type:          pluginsdk.TypeList,
						Optional:      true,
						ForceNew:      true,
						MaxItems:      1,
						ConflictsWith: []string{"git_repository.0.basic_auth"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"private_key": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"host_key": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"host_key_algorithm": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"branch": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ExactlyOneOf: []string{"git_repository.0.branch", "git_repository.0.commit", "git_repository.0.git_tag"},
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"ca_certificate_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.SpringCloudCertificateID,
					},

					"commit": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ExactlyOneOf: []string{"git_repository.0.branch", "git_repository.0.commit", "git_repository.0.git_tag"},
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"git_tag": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ExactlyOneOf: []string{"git_repository.0.branch", "git_repository.0.commit", "git_repository.0.git_tag"},
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"interval_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(10),
					},

					"path": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"accelerator_tags": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"accelerator_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  appplatform.CustomizedAcceleratorTypeAccelerator,
			ValidateFunc: validation.StringInSlice([]string{
				string(appplatform.CustomizedAcceleratorTypeAccelerator),
				string(appplatform.CustomizedAcceleratorTypeFragment),
			}, false),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"icon_url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (s SpringCloudCustomizedAcceleratorResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (s SpringCloudCustomizedAcceleratorResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudCustomizedAcceleratorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AppPlatform.AppPlatformClient
			springAcceleratorId, err := appplatform.ParseApplicationAcceleratorID(model.SpringCloudAcceleratorId)
			if err != nil {
				return fmt.Errorf("parsing spring service ID: %+v", err)
			}
			id := appplatform.NewCustomizedAcceleratorID(springAcceleratorId.SubscriptionId, springAcceleratorId.ResourceGroupName, springAcceleratorId.SpringName, springAcceleratorId.ApplicationAcceleratorName, model.Name)

			existing, err := client.CustomizedAcceleratorsGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(s.ResourceType(), id)
			}

			CustomizedAcceleratorResource := appplatform.CustomizedAcceleratorResource{
				Properties: &appplatform.CustomizedAcceleratorProperties{
					AcceleratorType: pointer.To(appplatform.CustomizedAcceleratorType(model.AcceleratorType)),
					DisplayName:     pointer.To(model.DisplayName),
					Description:     pointer.To(model.Description),
					IconURL:         pointer.To(model.IconURL),
					AcceleratorTags: pointer.To(model.AcceleratorTags),
					GitRepository:   expandSpringCloudCustomizedAcceleratorGitRepository(model.GitRepository),
				},
			}
			err = client.CustomizedAcceleratorsCreateOrUpdateThenPoll(ctx, id, CustomizedAcceleratorResource)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (s SpringCloudCustomizedAcceleratorResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.AppPlatformClient

			id, err := appplatform.ParseCustomizedAcceleratorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SpringCloudCustomizedAcceleratorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.CustomizedAcceleratorsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model.Properties
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("git_repository") {
				properties.GitRepository = expandSpringCloudCustomizedAcceleratorGitRepository(model.GitRepository)
			}

			if metadata.ResourceData.HasChange("accelerator_tags") {
				properties.AcceleratorTags = &model.AcceleratorTags
			}

			if metadata.ResourceData.HasChange("accelerator_type") {
				properties.AcceleratorType = pointer.To(appplatform.CustomizedAcceleratorType(model.AcceleratorType))
			}

			if metadata.ResourceData.HasChange("description") {
				properties.Description = &model.Description
			}

			if metadata.ResourceData.HasChange("display_name") {
				properties.DisplayName = &model.DisplayName
			}

			if metadata.ResourceData.HasChange("icon_url") {
				properties.IconURL = &model.IconURL
			}

			CustomizedAcceleratorResource := appplatform.CustomizedAcceleratorResource{
				Properties: properties,
			}
			err = client.CustomizedAcceleratorsCreateOrUpdateThenPoll(ctx, *id, CustomizedAcceleratorResource)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (s SpringCloudCustomizedAcceleratorResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.AppPlatformClient

			id, err := appplatform.ParseCustomizedAcceleratorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.CustomizedAcceleratorsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			state := SpringCloudCustomizedAcceleratorModel{
				Name:                     id.CustomizedAcceleratorName,
				SpringCloudAcceleratorId: appplatform.NewApplicationAcceleratorID(id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.ApplicationAcceleratorName).ID(),
			}

			if props := resp.Model.Properties; props != nil {
				if props.AcceleratorTags != nil {
					state.AcceleratorTags = *props.AcceleratorTags
				}
				if props.AcceleratorType != nil {
					state.AcceleratorType = string(*props.AcceleratorType)
				}
				if props.Description != nil {
					state.Description = *props.Description
				}
				if props.DisplayName != nil {
					state.DisplayName = *props.DisplayName
				}

				var model SpringCloudCustomizedAcceleratorModel
				if err := metadata.Decode(&model); err != nil {
					return fmt.Errorf("decoding: %+v", err)
				}
				state.GitRepository = flattenSpringCloudCustomizedAcceleratorGitRepository(model.GitRepository, props.GitRepository)

				if props.IconURL != nil {
					state.IconURL = *props.IconURL
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (s SpringCloudCustomizedAcceleratorResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.AppPlatformClient

			id, err := appplatform.ParseCustomizedAcceleratorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			err = client.CustomizedAcceleratorsDeleteThenPoll(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandSpringCloudCustomizedAcceleratorGitRepository(repository []GitRepositoryModel) appplatform.AcceleratorGitRepository {
	if len(repository) == 0 {
		return appplatform.AcceleratorGitRepository{}
	}
	repo := repository[0]
	var authSetting appplatform.AcceleratorAuthSetting
	var caCertResourceID *string
	if repo.CaCertificateId != "" {
		caCertResourceID = pointer.To(repo.CaCertificateId)
	}
	authSetting = appplatform.AcceleratorPublicSetting{
		CaCertResourceId: caCertResourceID,
	}
	if len(repo.BasicAuth) != 0 {
		basicAuth := repo.BasicAuth[0]
		authSetting = appplatform.AcceleratorBasicAuthSetting{
			Username:         basicAuth.Username,
			Password:         pointer.To(basicAuth.Password),
			CaCertResourceId: caCertResourceID,
		}
	}
	if len(repo.SshAuth) != 0 {
		sshAuth := repo.SshAuth[0]
		authSetting = appplatform.AcceleratorSshSetting{
			HostKey:          pointer.To(sshAuth.HostKey),
			HostKeyAlgorithm: pointer.To(sshAuth.PrivateKeyAlgorithm),
			PrivateKey:       pointer.To(sshAuth.PrivateKey),
		}
	}
	res := appplatform.AcceleratorGitRepository{
		Url:         repo.Url,
		Branch:      pointer.To(repo.Branch),
		Commit:      pointer.To(repo.Commit),
		GitTag:      pointer.To(repo.GitTag),
		AuthSetting: authSetting,
		SubPath:     pointer.To(repo.Path),
	}
	if repo.IntervalInSeconds != 0 {
		res.IntervalInSeconds = pointer.To(repo.IntervalInSeconds)
	}
	return res
}

func flattenSpringCloudCustomizedAcceleratorGitRepository(state []GitRepositoryModel, input appplatform.AcceleratorGitRepository) []GitRepositoryModel {
	basicAuth := make([]BasicAuthModel, 0)

	caCertificateId := ""
	if publicAuthSetting, ok := input.AuthSetting.(appplatform.AcceleratorPublicSetting); ok && publicAuthSetting.CaCertResourceId != nil {
		certificatedId, err := parse.SpringCloudCertificateIDInsensitively(*publicAuthSetting.CaCertResourceId)
		if err == nil {
			caCertificateId = certificatedId.ID()
		}
	}
	if basicAuthSetting, ok := input.AuthSetting.(appplatform.AcceleratorBasicAuthSetting); ok {
		if basicAuthSetting.CaCertResourceId != nil {
			certificatedId, err := parse.SpringCloudCertificateIDInsensitively(*basicAuthSetting.CaCertResourceId)
			if err == nil {
				caCertificateId = certificatedId.ID()
			}
		}
		var basicAuthState BasicAuthModel
		if len(state) != 0 && len(state[0].BasicAuth) != 0 {
			basicAuthState = state[0].BasicAuth[0]
		}
		basicAuth = append(basicAuth, BasicAuthModel{
			Username: basicAuthSetting.Username,
			Password: basicAuthState.Password,
		})
	}

	sshAuth := make([]SshAuthModel, 0)
	if _, ok := input.AuthSetting.(appplatform.AcceleratorSshSetting); ok {
		var sshAuthState SshAuthModel
		if len(state) != 0 && len(state[0].SshAuth) != 0 {
			sshAuthState = state[0].SshAuth[0]
		}
		sshAuth = append(sshAuth, sshAuthState)
	}

	branch := ""
	if input.Branch != nil {
		branch = *input.Branch
	}

	commit := ""
	if input.Commit != nil {
		commit = *input.Commit
	}

	gitTag := ""
	if input.GitTag != nil {
		gitTag = *input.GitTag
	}

	var intervalInSeconds int64
	if input.IntervalInSeconds != nil {
		intervalInSeconds = *input.IntervalInSeconds
	}

	subPath := ""
	if input.SubPath != nil {
		subPath = *input.SubPath
	}

	url := input.Url

	return []GitRepositoryModel{
		{
			BasicAuth:         basicAuth,
			SshAuth:           sshAuth,
			Branch:            branch,
			CaCertificateId:   caCertificateId,
			Commit:            commit,
			GitTag:            gitTag,
			IntervalInSeconds: intervalInSeconds,
			Url:               url,
			Path:              subPath,
		},
	}
}
