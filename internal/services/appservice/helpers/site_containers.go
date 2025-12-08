// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"fmt"
	"sort"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// SiteContainer models the Terraform schema for linux web app sidecar containers.
type SiteContainer struct {
	Name                        string                             `tfschema:"name"`
	Image                       string                             `tfschema:"image"`
	IsMain                      bool                               `tfschema:"is_main"`
	TargetPort                  string                             `tfschema:"target_port"`
	AuthType                    string                             `tfschema:"auth_type"`
	StartUpCommand              string                             `tfschema:"startup_command"`
	UserManagedIdentityClientID string                             `tfschema:"user_managed_identity_client_id"`
	Username                    string                             `tfschema:"username"`
	PasswordSecret              string                             `tfschema:"password_secret"`
	EnvironmentVariables        []SiteContainerEnvironmentVariable `tfschema:"environment_variable"`
	VolumeMounts                []SiteContainerVolumeMount         `tfschema:"volume_mount"`
}

// SiteContainerEnvironmentVariable captures container environment variables.
type SiteContainerEnvironmentVariable struct {
	Name  string `tfschema:"name"`
	Value string `tfschema:"value"`
}

// SiteContainerVolumeMount represents a single volume mount entry.
type SiteContainerVolumeMount struct {
	ContainerMountPath string `tfschema:"container_mount_path"`
	Data               string `tfschema:"data"`
	ReadOnly           bool   `tfschema:"read_only"`
	VolumeSubPath      string `tfschema:"volume_sub_path"`
}

func SiteContainerSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"image": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"is_main": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
				"target_port": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"auth_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      string(webapps.AuthTypeAnonymous),
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForAuthType(), false),
				},
				"startup_command": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
				"user_managed_identity_client_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsUUID,
				},
				"username": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
				"password_secret": {
					Type:      pluginsdk.TypeString,
					Optional:  true,
					Sensitive: true,
				},
				"environment_variable": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"name": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"value": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
				"volume_mount": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"container_mount_path": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"data": {
								Type:     pluginsdk.TypeString,
								Optional: true,
							},
							"read_only": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
							},
							"volume_sub_path": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},
		},
	}
}

func ExpandSiteContainers(input []SiteContainer) ([]webapps.SiteContainer, error) {
	if len(input) == 0 {
		return nil, nil
	}

	result := make([]webapps.SiteContainer, 0, len(input))
	for _, container := range input {
		if container.Name == "" {
			return nil, fmt.Errorf("`name` must be specified for `site_container`")
		}
		if container.Image == "" {
			return nil, fmt.Errorf("`image` must be specified for `site_container` %q", container.Name)
		}
		if container.TargetPort == "" {
			return nil, fmt.Errorf("`target_port` must be specified for `site_container` %q", container.Name)
		}

		authType := webapps.AuthType(container.AuthType)
		props := &webapps.SiteContainerProperties{
			AuthType:                    pointer.To(authType),
			EnvironmentVariables:        expandSiteContainerEnvVars(container.EnvironmentVariables),
			Image:                       container.Image,
			IsMain:                      container.IsMain,
			PasswordSecret:              pointer.To(container.PasswordSecret),
			StartUpCommand:              pointer.To(container.StartUpCommand),
			TargetPort:                  pointer.To(container.TargetPort),
			UserManagedIdentityClientId: pointer.To(container.UserManagedIdentityClientID),
			UserName:                    pointer.To(container.Username),
			VolumeMounts:                expandSiteContainerVolumeMounts(container.VolumeMounts),
		}

		if container.PasswordSecret == "" {
			props.PasswordSecret = nil
		}
		if container.StartUpCommand == "" {
			props.StartUpCommand = nil
		}
		if container.UserManagedIdentityClientID == "" {
			props.UserManagedIdentityClientId = nil
		}
		if container.Username == "" {
			props.UserName = nil
		}
		if props.EnvironmentVariables == nil {
			props.EnvironmentVariables = &[]webapps.EnvironmentVariable{}
		}
		if props.VolumeMounts == nil {
			props.VolumeMounts = &[]webapps.VolumeMount{}
		}

		result = append(result, webapps.SiteContainer{
			Name:       pointer.To(container.Name),
			Properties: props,
		})
	}

	return result, nil
}

func expandSiteContainerEnvVars(input []SiteContainerEnvironmentVariable) *[]webapps.EnvironmentVariable {
	if len(input) == 0 {
		return nil
	}

	envs := make([]webapps.EnvironmentVariable, 0, len(input))
	for _, env := range input {
		envs = append(envs, webapps.EnvironmentVariable{
			Name:  env.Name,
			Value: env.Value,
		})
	}

	return &envs
}

func expandSiteContainerVolumeMounts(input []SiteContainerVolumeMount) *[]webapps.VolumeMount {
	if len(input) == 0 {
		return nil
	}

	mounts := make([]webapps.VolumeMount, 0, len(input))
	for _, m := range input {
		mount := webapps.VolumeMount{
			ContainerMountPath: m.ContainerMountPath,
			VolumeSubPath:      m.VolumeSubPath,
		}
		if m.Data != "" {
			mount.Data = pointer.To(m.Data)
		}
		if m.ReadOnly {
			mount.ReadOnly = pointer.To(true)
		}
		mounts = append(mounts, mount)
	}

	return &mounts
}

func FlattenSiteContainers(input []webapps.SiteContainer) []SiteContainer {
	if len(input) == 0 {
		return nil
	}

	result := make([]SiteContainer, 0, len(input))
	for _, container := range input {
		props := container.Properties
		if props == nil {
			continue
		}

		flattened := SiteContainer{
			Name:       pointer.From(container.Name),
			Image:      props.Image,
			IsMain:     props.IsMain,
			TargetPort: pointer.From(props.TargetPort),
			AuthType: func() string {
				if props.AuthType == nil {
					return string(webapps.AuthTypeAnonymous)
				}
				return string(*props.AuthType)
			}(),
			StartUpCommand:              pointer.From(props.StartUpCommand),
			UserManagedIdentityClientID: pointer.From(props.UserManagedIdentityClientId),
			Username:                    pointer.From(props.UserName),
			PasswordSecret:              pointer.From(props.PasswordSecret),
		}

		if props.EnvironmentVariables != nil {
			flattened.EnvironmentVariables = flattenSiteContainerEnvVars(*props.EnvironmentVariables)
		}

		if props.VolumeMounts != nil {
			flattened.VolumeMounts = flattenSiteContainerVolumeMounts(*props.VolumeMounts)
		}

		result = append(result, flattened)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

func flattenSiteContainerEnvVars(input []webapps.EnvironmentVariable) []SiteContainerEnvironmentVariable {
	if len(input) == 0 {
		return nil
	}

	envs := make([]SiteContainerEnvironmentVariable, 0, len(input))
	for _, env := range input {
		envs = append(envs, SiteContainerEnvironmentVariable{
			Name:  env.Name,
			Value: env.Value,
		})
	}

	sort.Slice(envs, func(i, j int) bool {
		return envs[i].Name < envs[j].Name
	})

	return envs
}

func flattenSiteContainerVolumeMounts(input []webapps.VolumeMount) []SiteContainerVolumeMount {
	if len(input) == 0 {
		return nil
	}

	mounts := make([]SiteContainerVolumeMount, 0, len(input))
	for _, mount := range input {
		flattened := SiteContainerVolumeMount{
			ContainerMountPath: mount.ContainerMountPath,
			VolumeSubPath:      mount.VolumeSubPath,
		}
		if mount.Data != nil {
			flattened.Data = pointer.From(mount.Data)
		}
		if mount.ReadOnly != nil {
			flattened.ReadOnly = pointer.From(mount.ReadOnly)
		}
		mounts = append(mounts, flattened)
	}

	sort.Slice(mounts, func(i, j int) bool {
		if mounts[i].ContainerMountPath == mounts[j].ContainerMountPath {
			return mounts[i].VolumeSubPath < mounts[j].VolumeSubPath
		}
		return mounts[i].ContainerMountPath < mounts[j].ContainerMountPath
	})

	return mounts
}
