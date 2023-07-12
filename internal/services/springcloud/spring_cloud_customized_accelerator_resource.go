// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

type SpringCloudCustomizedAcceleratorModel struct {
	Name                     string               `tfschema:"name"`
	SpringCloudAcceleratorId string               `tfschema:"spring_cloud_accelerator_id"`
	AcceleratorTags          []string             `tfschema:"accelerator_tags"`
	Description              string               `tfschema:"description"`
	DisplayName              string               `tfschema:"display_name"`
	GitRepository            []GitRepositoryModel `tfschema:"git_repository"`
	IconUrl                  string               `tfschema:"icon_url"`
}

type GitRepositoryModel struct {
	BasicAuth         []BasicAuthModel `tfschema:"basic_auth"`
	SshAuth           []SshAuthModel   `tfschema:"ssh_auth"`
	Branch            string           `tfschema:"branch"`
	CaCertificateId   string           `tfschema:"ca_certificate_id"`
	Commit            string           `tfschema:"commit"`
	GitTag            string           `tfschema:"git_tag"`
	IntervalInSeconds int              `tfschema:"interval_in_seconds"`
	Url               string           `tfschema:"url"`
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

var _ sdk.ResourceWithUpdate = SpringCloudCustomizedAcceleratorResource{}

func (s SpringCloudCustomizedAcceleratorResource) ResourceType() string {
	return "azurerm_spring_cloud_customized_accelerator"
}

func (s SpringCloudCustomizedAcceleratorResource) ModelObject() interface{} {
	return &SpringCloudCustomizedAcceleratorModel{}
}

func (s SpringCloudCustomizedAcceleratorResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.SpringCloudCustomizedAcceleratorID
}

func (s SpringCloudCustomizedAcceleratorResource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"spring_cloud_accelerator_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SpringCloudAcceleratorID,
		},

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

			client := metadata.Client.AppPlatform.CustomizedAcceleratorClient
			springAcceleratorId, err := parse.SpringCloudAcceleratorID(model.SpringCloudAcceleratorId)
			if err != nil {
				return fmt.Errorf("parsing spring service ID: %+v", err)
			}
			id := parse.NewSpringCloudCustomizedAcceleratorID(springAcceleratorId.SubscriptionId, springAcceleratorId.ResourceGroup, springAcceleratorId.SpringName, springAcceleratorId.ApplicationAcceleratorName, model.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.ApplicationAcceleratorName, id.CustomizedAcceleratorName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(s.ResourceType(), id)
			}

			CustomizedAcceleratorResource := appplatform.CustomizedAcceleratorResource{
				Properties: &appplatform.CustomizedAcceleratorProperties{
					DisplayName:     utils.String(model.DisplayName),
					Description:     utils.String(model.Description),
					IconURL:         utils.String(model.IconUrl),
					AcceleratorTags: utils.ToPtr(model.AcceleratorTags),
					GitRepository:   expandSpringCloudCustomizedAcceleratorGitRepository(model.GitRepository),
				},
			}
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.ApplicationAcceleratorName, id.CustomizedAcceleratorName, CustomizedAcceleratorResource)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
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
			client := metadata.Client.AppPlatform.CustomizedAcceleratorClient

			id, err := parse.SpringCloudCustomizedAcceleratorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SpringCloudCustomizedAcceleratorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.ApplicationAcceleratorName, id.CustomizedAcceleratorName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Properties
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("git_repository") {
				properties.GitRepository = expandSpringCloudCustomizedAcceleratorGitRepository(model.GitRepository)
			}

			if metadata.ResourceData.HasChange("accelerator_tags") {
				properties.AcceleratorTags = &model.AcceleratorTags
			}

			if metadata.ResourceData.HasChange("description") {
				properties.Description = &model.Description
			}

			if metadata.ResourceData.HasChange("display_name") {
				properties.DisplayName = &model.DisplayName
			}

			if metadata.ResourceData.HasChange("icon_url") {
				properties.IconURL = &model.IconUrl
			}

			CustomizedAcceleratorResource := appplatform.CustomizedAcceleratorResource{
				Properties: properties,
			}
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.ApplicationAcceleratorName, id.CustomizedAcceleratorName, CustomizedAcceleratorResource)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (s SpringCloudCustomizedAcceleratorResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.CustomizedAcceleratorClient

			id, err := parse.SpringCloudCustomizedAcceleratorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.ApplicationAcceleratorName, id.CustomizedAcceleratorName)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			state := SpringCloudCustomizedAcceleratorModel{
				Name:                     id.CustomizedAcceleratorName,
				SpringCloudAcceleratorId: parse.NewSpringCloudAcceleratorID(id.SubscriptionId, id.ResourceGroup, id.SpringName, id.ApplicationAcceleratorName).ID(),
			}

			if props := resp.Properties; props != nil {
				if props.AcceleratorTags != nil {
					state.AcceleratorTags = *props.AcceleratorTags
				}
				if props.Description != nil {
					state.Description = *props.Description
				}
				if props.DisplayName != nil {
					state.DisplayName = *props.DisplayName
				}
				if props.GitRepository != nil {
					var model SpringCloudCustomizedAcceleratorModel
					if err := metadata.Decode(&model); err != nil {
						return fmt.Errorf("decoding: %+v", err)
					}
					state.GitRepository = flattenSpringCloudCustomizedAcceleratorGitRepository(model.GitRepository, props.GitRepository)
				}
				if props.IconURL != nil {
					state.IconUrl = *props.IconURL
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
			client := metadata.Client.AppPlatform.CustomizedAcceleratorClient

			id, err := parse.SpringCloudCustomizedAcceleratorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.SpringName, id.ApplicationAcceleratorName, id.CustomizedAcceleratorName)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				if !response.WasNotFound(future.Response()) {
					return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

func expandSpringCloudCustomizedAcceleratorGitRepository(repository []GitRepositoryModel) *appplatform.AcceleratorGitRepository {
	if len(repository) == 0 {
		return nil
	}
	repo := repository[0]
	var authSetting appplatform.BasicAcceleratorAuthSetting
	var caCertResourceID *string
	if repo.CaCertificateId != "" {
		caCertResourceID = utils.String(repo.CaCertificateId)
	}
	authSetting = appplatform.AcceleratorPublicSetting{
		CaCertResourceID: caCertResourceID,
	}
	if len(repo.BasicAuth) != 0 {
		basicAuth := repo.BasicAuth[0]
		authSetting = appplatform.AcceleratorBasicAuthSetting{
			Username:         utils.String(basicAuth.Username),
			Password:         utils.String(basicAuth.Password),
			CaCertResourceID: caCertResourceID,
		}
	}
	if len(repo.SshAuth) != 0 {
		sshAuth := repo.SshAuth[0]
		authSetting = appplatform.AcceleratorSSHSetting{
			HostKey:          utils.String(sshAuth.HostKey),
			HostKeyAlgorithm: utils.String(sshAuth.PrivateKeyAlgorithm),
			PrivateKey:       utils.String(sshAuth.PrivateKey),
		}
	}
	res := &appplatform.AcceleratorGitRepository{
		URL:         utils.String(repo.Url),
		Branch:      utils.String(repo.Branch),
		Commit:      utils.String(repo.Commit),
		GitTag:      utils.String(repo.GitTag),
		AuthSetting: authSetting,
	}
	if repo.IntervalInSeconds != 0 {
		res.IntervalInSeconds = utils.Int32(int32(repo.IntervalInSeconds))
	}
	return res
}

func flattenSpringCloudCustomizedAcceleratorGitRepository(state []GitRepositoryModel, input *appplatform.AcceleratorGitRepository) []GitRepositoryModel {
	if input == nil {
		return []GitRepositoryModel{}
	}

	basicAuth := make([]BasicAuthModel, 0)

	caCertificateId := ""
	if publicAuthSetting, ok := input.AuthSetting.AsAcceleratorPublicSetting(); ok && publicAuthSetting != nil && publicAuthSetting.CaCertResourceID != nil {
		certificatedId, err := parse.SpringCloudCertificateIDInsensitively(*publicAuthSetting.CaCertResourceID)
		if err == nil {
			caCertificateId = certificatedId.ID()
		}
	}
	if basicAuthSetting, ok := input.AuthSetting.AsAcceleratorBasicAuthSetting(); ok && basicAuthSetting != nil {
		if basicAuthSetting.CaCertResourceID != nil {
			certificatedId, err := parse.SpringCloudCertificateIDInsensitively(*basicAuthSetting.CaCertResourceID)
			if err == nil {
				caCertificateId = certificatedId.ID()
			}
		}
		var basicAuthState BasicAuthModel
		if len(state) != 0 && len(state[0].BasicAuth) != 0 {
			basicAuthState = state[0].BasicAuth[0]
		}
		basicAuth = append(basicAuth, BasicAuthModel{
			Username: *basicAuthSetting.Username,
			Password: basicAuthState.Password,
		})
	}

	sshAuth := make([]SshAuthModel, 0)
	if sshAuthSetting, ok := input.AuthSetting.AsAcceleratorSSHSetting(); ok && sshAuthSetting != nil {
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

	intervalInSeconds := 0
	if input.IntervalInSeconds != nil {
		intervalInSeconds = int(*input.IntervalInSeconds)
	}

	url := ""
	if input.URL != nil {
		url = *input.URL
	}

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
		},
	}
}
