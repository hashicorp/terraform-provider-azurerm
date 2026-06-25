// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name linux_web_app_site_container -service-package-name appservice -properties "name" -compare-values "subscription_id:linux_web_app_id,resource_group_name:linux_web_app_id,site_name:linux_web_app_id"

type LinuxWebAppSiteContainerResource struct{}

type LinuxWebAppSiteContainerModel struct {
	Name                        string                                     `tfschema:"name"`
	LinuxWebAppId               string                                     `tfschema:"linux_web_app_id"`
	Image                       string                                     `tfschema:"image"`
	TargetPort                  int64                                      `tfschema:"target_port"`
	Primary                     bool                                       `tfschema:"primary"`
	AuthenticationType          string                                     `tfschema:"authentication_type"`
	StartUpCommand              string                                     `tfschema:"startup_command"`
	UserManagedIdentityClientID string                                     `tfschema:"user_managed_identity_client_id"`
	Username                    string                                     `tfschema:"username"`
	PasswordSecret              string                                     `tfschema:"password_secret"`
	EnvironmentVariables        []helpers.SiteContainerEnvironmentVariable `tfschema:"environment_variable"`
	VolumeMounts                []helpers.SiteContainerVolumeMount         `tfschema:"volume_mount"`
}

var (
	_ sdk.ResourceWithUpdate   = LinuxWebAppSiteContainerResource{}
	_ sdk.ResourceWithIdentity = LinuxWebAppSiteContainerResource{}
)

func (r LinuxWebAppSiteContainerResource) Identity() resourceids.ResourceId {
	return &webapps.SitecontainerId{}
}

func (r LinuxWebAppSiteContainerResource) ModelObject() interface{} {
	return &LinuxWebAppSiteContainerModel{}
}

func (r LinuxWebAppSiteContainerResource) ResourceType() string {
	return "azurerm_linux_web_app_site_container"
}

func (r LinuxWebAppSiteContainerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return webapps.ValidateSitecontainerID
}

func (r LinuxWebAppSiteContainerResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?$`),
				"`name` must start and end with an alphanumeric character and may contain hyphens",
			),
		},

		"linux_web_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateWebAppID,
		},

		"image": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"target_port": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IsPortNumber,
		},

		"authentication_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(webapps.AuthTypeAnonymous),
			ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForAuthType(), false),
		},

		"environment_variable": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			MinItems: 1,
			Set: func(v interface{}) int {
				m := v.(map[string]interface{})
				return pluginsdk.HashString(m["name"].(string))
			},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"app_setting_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"password_secret": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"primary": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"startup_command": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"user_managed_identity_client_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
		},

		"username": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"volume_mount": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"container_mount_path": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"data": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"read_only": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"volume_sub_path": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}

func (r LinuxWebAppSiteContainerResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LinuxWebAppSiteContainerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			var model LinuxWebAppSiteContainerModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			if err := model.validateAuthentication(); err != nil {
				return err
			}

			appId, err := commonids.ParseWebAppID(model.LinuxWebAppId)
			if err != nil {
				return err
			}

			id := webapps.NewSitecontainerID(appId.SubscriptionId, appId.ResourceGroupName, appId.SiteName, model.Name)

			existing, err := client.GetSiteContainer(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(r.ResourceType(), id.ID())
			}

			if _, err := client.CreateOrUpdateSiteContainer(ctx, id, webapps.SiteContainer{
				Name:       pointer.To(model.Name),
				Properties: model.expandProperties(),
			}); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r LinuxWebAppSiteContainerResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := webapps.ParseSitecontainerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.GetSiteContainer(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			var config LinuxWebAppSiteContainerModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

			container := webapps.SiteContainer{}
			if existing.Model != nil {
				container = *existing.Model
			}

			state := flattenLinuxWebAppSiteContainer(*id, container, config.PasswordSecret)

			return metadata.Encode(&state)
		},
	}
}

func (r LinuxWebAppSiteContainerResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := webapps.ParseSitecontainerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model LinuxWebAppSiteContainerModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			if err := model.validateAuthentication(); err != nil {
				return err
			}

			if _, err := client.CreateOrUpdateSiteContainer(ctx, *id, webapps.SiteContainer{
				Name:       pointer.To(id.SitecontainerName),
				Properties: model.expandProperties(),
			}); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r LinuxWebAppSiteContainerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := webapps.ParseSitecontainerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.DeleteSiteContainer(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (m LinuxWebAppSiteContainerModel) expandProperties() *webapps.SiteContainerProperties {
	props := &webapps.SiteContainerProperties{
		AuthType:             pointer.ToEnum[webapps.AuthType](m.AuthenticationType),
		EnvironmentVariables: helpers.ExpandSiteContainerEnvironmentVariables(m.EnvironmentVariables),
		Image:                m.Image,
		IsMain:               m.Primary,
		TargetPort:           pointer.To(strconv.FormatInt(m.TargetPort, 10)),
		VolumeMounts:         helpers.ExpandSiteContainerVolumeMounts(m.VolumeMounts),
	}

	if m.PasswordSecret != "" {
		props.PasswordSecret = pointer.To(m.PasswordSecret)
	}
	if m.StartUpCommand != "" {
		props.StartUpCommand = pointer.To(m.StartUpCommand)
	}
	if m.UserManagedIdentityClientID != "" {
		props.UserManagedIdentityClientId = pointer.To(m.UserManagedIdentityClientID)
	}
	if m.Username != "" {
		props.UserName = pointer.To(m.Username)
	}

	return props
}

func (m LinuxWebAppSiteContainerModel) validateAuthentication() error {
	switch webapps.AuthType(m.AuthenticationType) {
	case webapps.AuthTypeUserCredentials:
		if m.Username == "" || m.PasswordSecret == "" {
			return fmt.Errorf("`username` and `password_secret` must be set when `authentication_type` is `%s`", webapps.AuthTypeUserCredentials)
		}
	case webapps.AuthTypeUserAssigned:
		if m.UserManagedIdentityClientID == "" {
			return fmt.Errorf("`user_managed_identity_client_id` must be set when `authentication_type` is `%s`", webapps.AuthTypeUserAssigned)
		}
	}

	return nil
}

func flattenLinuxWebAppSiteContainer(id webapps.SitecontainerId, container webapps.SiteContainer, priorSecret string) LinuxWebAppSiteContainerModel {
	state := LinuxWebAppSiteContainerModel{
		Name:           id.SitecontainerName,
		LinuxWebAppId:  commonids.NewAppServiceID(id.SubscriptionId, id.ResourceGroupName, id.SiteName).ID(),
		PasswordSecret: priorSecret,
	}

	if props := container.Properties; props != nil {
		state.Image = props.Image
		state.Primary = props.IsMain
		state.AuthenticationType = string(webapps.AuthTypeAnonymous)
		if props.AuthType != nil {
			state.AuthenticationType = string(*props.AuthType)
		}
		if props.TargetPort != nil {
			if parsed, err := strconv.ParseInt(*props.TargetPort, 10, 64); err == nil {
				state.TargetPort = parsed
			}
		}
		state.StartUpCommand = pointer.From(props.StartUpCommand)
		state.UserManagedIdentityClientID = pointer.From(props.UserManagedIdentityClientId)
		state.Username = pointer.From(props.UserName)
		if props.PasswordSecret != nil {
			state.PasswordSecret = pointer.From(props.PasswordSecret)
		}
		state.EnvironmentVariables = helpers.FlattenSiteContainerEnvironmentVariables(props.EnvironmentVariables)
		state.VolumeMounts = helpers.FlattenSiteContainerVolumeMounts(props.VolumeMounts)
	}

	return state
}
