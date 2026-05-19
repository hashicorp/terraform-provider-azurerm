// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SiteContainer struct {
	Name                        string                             `tfschema:"name"`
	Image                       string                             `tfschema:"image"`
	Main                        bool                               `tfschema:"main"`
	TargetPort                  int                                `tfschema:"target_port"`
	AuthenticationType          string                             `tfschema:"authentication_type"`
	StartUpCommand              string                             `tfschema:"startup_command"`
	UserManagedIdentityClientID string                             `tfschema:"user_managed_identity_client_id"`
	Username                    string                             `tfschema:"username"`
	PasswordSecret              string                             `tfschema:"password_secret"`
	EnvironmentVariables        []SiteContainerEnvironmentVariable `tfschema:"environment_variable"`
	VolumeMounts                []SiteContainerVolumeMount         `tfschema:"volume_mount"`
}

type SiteContainerEnvironmentVariable struct {
	Name           string `tfschema:"name"`
	AppSettingName string `tfschema:"app_setting_name"`
}

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
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?$`),
						"name must start and end with an alphanumeric character and may contain hyphens",
					),
				},
				"image": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"main": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
				"target_port": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IsPortNumber,
				},
				"authentication_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      string(webapps.AuthTypeAnonymous),
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForAuthType(), false),
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
				"password_secret": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
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
							"app_setting_name": {
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
			},
		},
	}
}

// ValidateSiteContainers performs cross-field validation for the `site_container` block.
// Callers pass whether the parent resource has a `site_config` block and whether
// `site_config.0.application_stack` is set, since those struct types differ across
// the four web app resource variants.
func ValidateSiteContainers(containers []SiteContainer, hasSiteConfig bool, hasApplicationStack bool) error {
	if len(containers) == 0 {
		return nil
	}

	if !hasSiteConfig {
		return errors.New("`site_config` must be configured when using `site_container`")
	}
	if hasApplicationStack {
		return errors.New("`site_container` cannot be used when `site_config.0.application_stack` is specified")
	}

	mainCount := 0
	names := make(map[string]struct{}, len(containers))
	for _, container := range containers {
		if container.Main {
			mainCount++
		}
		if _, exists := names[container.Name]; exists {
			return fmt.Errorf("`site_container` names must be unique; duplicate `%s` found", container.Name)
		}
		names[container.Name] = struct{}{}

		switch webapps.AuthType(container.AuthenticationType) {
		case webapps.AuthTypeUserCredentials:
			if container.Username == "" || container.PasswordSecret == "" {
				return fmt.Errorf("`username` and `password_secret` must be set for `site_container` %q when `authentication_type` is `UserCredentials`", container.Name)
			}
		case webapps.AuthTypeUserAssigned:
			if container.UserManagedIdentityClientID == "" {
				return fmt.Errorf("`user_managed_identity_client_id` must be set for `site_container` %q when `authentication_type` is `UserAssigned`", container.Name)
			}
		}
	}
	if mainCount != 1 {
		return errors.New("exactly one `site_container` must have `main` set to `true`")
	}

	return nil
}

func ExpandSiteContainers(input []SiteContainer) ([]webapps.SiteContainer, error) {
	if len(input) == 0 {
		return nil, nil
	}

	result := make([]webapps.SiteContainer, 0, len(input))
	for _, container := range input {
		if container.Name == "" {
			return nil, errors.New("`name` must be specified for `site_container`")
		}
		if container.Image == "" {
			return nil, fmt.Errorf("`image` must be specified for `site_container` %q", container.Name)
		}

		authType := webapps.AuthType(container.AuthenticationType)
		props := &webapps.SiteContainerProperties{
			AuthType:             pointer.To(authType),
			EnvironmentVariables: expandSiteContainerEnvVars(container.EnvironmentVariables),
			Image:                container.Image,
			IsMain:               container.Main,
			VolumeMounts:         expandSiteContainerVolumeMounts(container.VolumeMounts),
		}

		if container.TargetPort != 0 {
			props.TargetPort = pointer.To(strconv.Itoa(container.TargetPort))
		}
		if container.PasswordSecret != "" {
			props.PasswordSecret = pointer.To(container.PasswordSecret)
		}
		if container.StartUpCommand != "" {
			props.StartUpCommand = pointer.To(container.StartUpCommand)
		}
		if container.UserManagedIdentityClientID != "" {
			props.UserManagedIdentityClientId = pointer.To(container.UserManagedIdentityClientID)
		}
		if container.Username != "" {
			props.UserName = pointer.To(container.Username)
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
			Value: env.AppSettingName,
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

func FlattenSiteContainers(input []webapps.SiteContainer) ([]SiteContainer, map[string]struct{}) {
	result := make([]SiteContainer, 0, len(input))
	missingSecrets := make(map[string]struct{})
	for _, container := range input {
		props := container.Properties
		if props == nil {
			continue
		}

		targetPort := 0
		if props.TargetPort != nil {
			if parsed, err := strconv.Atoi(*props.TargetPort); err == nil {
				targetPort = parsed
			}
		}

		flattened := SiteContainer{
			Name:       pointer.From(container.Name),
			Image:      props.Image,
			Main:       props.IsMain,
			TargetPort: targetPort,
			AuthenticationType: func() string {
				if props.AuthType == nil {
					return string(webapps.AuthTypeAnonymous)
				}
				return string(*props.AuthType)
			}(),
			StartUpCommand:              pointer.From(props.StartUpCommand),
			UserManagedIdentityClientID: pointer.From(props.UserManagedIdentityClientId),
			Username:                    pointer.From(props.UserName),
		}

		if props.PasswordSecret != nil {
			flattened.PasswordSecret = pointer.From(props.PasswordSecret)
		} else if flattened.Name != "" {
			missingSecrets[flattened.Name] = struct{}{}
		}

		flattened.EnvironmentVariables = flattenSiteContainerEnvVars(props.EnvironmentVariables)
		flattened.VolumeMounts = flattenSiteContainerVolumeMounts(props.VolumeMounts)

		result = append(result, flattened)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, missingSecrets
}

func flattenSiteContainerEnvVars(input *[]webapps.EnvironmentVariable) []SiteContainerEnvironmentVariable {
	envs := []SiteContainerEnvironmentVariable{}
	if input == nil || len(*input) == 0 {
		return envs
	}

	for _, env := range *input {
		envs = append(envs, SiteContainerEnvironmentVariable{
			Name:           env.Name,
			AppSettingName: env.Value,
		})
	}

	sort.Slice(envs, func(i, j int) bool {
		return envs[i].Name < envs[j].Name
	})

	return envs
}

func flattenSiteContainerVolumeMounts(input *[]webapps.VolumeMount) []SiteContainerVolumeMount {
	mounts := []SiteContainerVolumeMount{}
	if input == nil || len(*input) == 0 {
		return mounts
	}

	for _, mount := range *input {
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
