// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
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
	Primary                     bool                               `tfschema:"primary"`
	TargetPort                  int64                              `tfschema:"target_port"`
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
				"primary": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
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
// `hasApplicationStack` reports whether `site_config.0.application_stack` is configured
// on the parent resource. The check is passed in as a bool because the application stack
// struct types differ across the four web app resource variants.
func ValidateSiteContainers(containers []SiteContainer, hasApplicationStack bool) error {
	if len(containers) == 0 {
		return nil
	}

	if hasApplicationStack {
		return errors.New("`site_container` cannot be used when `site_config.0.application_stack` is specified")
	}

	primaryCount := 0
	names := make(map[string]struct{}, len(containers))
	for _, container := range containers {
		if container.Primary {
			primaryCount++
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
	if primaryCount != 1 {
		return errors.New("exactly one `site_container` must have `primary` set to `true`")
	}

	return nil
}

// ReconcileSiteContainers diffs the desired set of `site_container` blocks against what
// currently exists on the Web App (or Slot) and issues create-or-update / delete calls
// to converge. The three closures abstract over the web-app vs slot SDK methods so the
// caller only supplies the resource-shaped wrappers.
func ReconcileSiteContainers(
	ctx context.Context,
	resourceLabel string,
	containers []SiteContainer,
	listExisting func(ctx context.Context) ([]webapps.SiteContainer, error),
	createOrUpdate func(ctx context.Context, name string, container webapps.SiteContainer) error,
	deleteContainer func(ctx context.Context, name string) error,
) error {
	expanded, err := ExpandSiteContainers(containers)
	if err != nil {
		return fmt.Errorf("expanding `site_container`: %+v", err)
	}

	desired := make(map[string]webapps.SiteContainer, len(expanded))
	for _, container := range expanded {
		name := pointer.From(container.Name)
		if name == "" {
			return errors.New("`site_container` entries must include `name`")
		}
		desired[name] = container
	}

	existing, err := listExisting(ctx)
	if err != nil {
		return fmt.Errorf("listing Site Containers for %s: %+v", resourceLabel, err)
	}

	names := make([]string, 0, len(desired))
	for name := range desired {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		if err := createOrUpdate(ctx, name, desired[name]); err != nil {
			return fmt.Errorf("creating or updating Site Container `%s` for %s: %+v", name, resourceLabel, err)
		}
	}

	for _, existingContainer := range existing {
		name := pointer.From(existingContainer.Name)
		if name == "" {
			continue
		}
		if _, keep := desired[name]; keep {
			continue
		}
		if err := deleteContainer(ctx, name); err != nil {
			return fmt.Errorf("deleting Site Container `%s` for %s: %+v", name, resourceLabel, err)
		}
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
			IsMain:               container.Primary,
			VolumeMounts:         expandSiteContainerVolumeMounts(container.VolumeMounts),
		}

		if container.TargetPort != 0 {
			props.TargetPort = pointer.To(strconv.FormatInt(container.TargetPort, 10))
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

// FlattenSiteContainers converts the Azure response into the schema model. The order of
// `site_container`, `environment_variable`, and `volume_mount` is preserved from `existing`
// (the previously stored state) to avoid spurious diffs on a TypeList; any items returned
// by Azure that are not present in the existing state are appended in name/path order.
func FlattenSiteContainers(input []webapps.SiteContainer, existing []SiteContainer) ([]SiteContainer, map[string]struct{}) {
	missingSecrets := make(map[string]struct{})
	flattenedByName := make(map[string]SiteContainer, len(input))
	for _, container := range input {
		props := container.Properties
		if props == nil {
			continue
		}

		targetPort := int64(0)
		if props.TargetPort != nil {
			if parsed, err := strconv.ParseInt(*props.TargetPort, 10, 64); err == nil {
				targetPort = parsed
			}
		}

		name := pointer.From(container.Name)
		flattened := SiteContainer{
			Name:       name,
			Image:      props.Image,
			Primary:    props.IsMain,
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
		} else if name != "" {
			missingSecrets[name] = struct{}{}
		}

		existingForContainer := lookupSiteContainer(existing, name)
		flattened.EnvironmentVariables = flattenSiteContainerEnvVars(props.EnvironmentVariables, existingForContainer.EnvironmentVariables)
		flattened.VolumeMounts = flattenSiteContainerVolumeMounts(props.VolumeMounts, existingForContainer.VolumeMounts)

		flattenedByName[name] = flattened
	}

	result := make([]SiteContainer, 0, len(flattenedByName))
	seen := make(map[string]struct{}, len(flattenedByName))
	for _, prior := range existing {
		if container, ok := flattenedByName[prior.Name]; ok {
			result = append(result, container)
			seen[prior.Name] = struct{}{}
		}
	}

	remaining := make([]string, 0, len(flattenedByName)-len(seen))
	for name := range flattenedByName {
		if _, ok := seen[name]; ok {
			continue
		}
		remaining = append(remaining, name)
	}
	sort.Strings(remaining)
	for _, name := range remaining {
		result = append(result, flattenedByName[name])
	}

	return result, missingSecrets
}

func lookupSiteContainer(existing []SiteContainer, name string) SiteContainer {
	for _, e := range existing {
		if e.Name == name {
			return e
		}
	}
	return SiteContainer{}
}

func flattenSiteContainerEnvVars(input *[]webapps.EnvironmentVariable, existing []SiteContainerEnvironmentVariable) []SiteContainerEnvironmentVariable {
	envs := []SiteContainerEnvironmentVariable{}
	if input == nil || len(*input) == 0 {
		return envs
	}

	byName := make(map[string]SiteContainerEnvironmentVariable, len(*input))
	for _, env := range *input {
		byName[env.Name] = SiteContainerEnvironmentVariable{
			Name:           env.Name,
			AppSettingName: env.Value,
		}
	}

	seen := make(map[string]struct{}, len(byName))
	for _, prior := range existing {
		if env, ok := byName[prior.Name]; ok {
			envs = append(envs, env)
			seen[prior.Name] = struct{}{}
		}
	}

	remaining := make([]string, 0, len(byName)-len(seen))
	for name := range byName {
		if _, ok := seen[name]; ok {
			continue
		}
		remaining = append(remaining, name)
	}
	sort.Strings(remaining)
	for _, name := range remaining {
		envs = append(envs, byName[name])
	}

	return envs
}

func flattenSiteContainerVolumeMounts(input *[]webapps.VolumeMount, existing []SiteContainerVolumeMount) []SiteContainerVolumeMount {
	mounts := []SiteContainerVolumeMount{}
	if input == nil || len(*input) == 0 {
		return mounts
	}

	type mountKey struct {
		path, sub string
	}
	byKey := make(map[mountKey]SiteContainerVolumeMount, len(*input))
	keys := make([]mountKey, 0, len(*input))
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
		key := mountKey{path: flattened.ContainerMountPath, sub: flattened.VolumeSubPath}
		byKey[key] = flattened
		keys = append(keys, key)
	}

	seen := make(map[mountKey]struct{}, len(byKey))
	for _, prior := range existing {
		key := mountKey{path: prior.ContainerMountPath, sub: prior.VolumeSubPath}
		if mount, ok := byKey[key]; ok {
			mounts = append(mounts, mount)
			seen[key] = struct{}{}
		}
	}

	sort.Slice(keys, func(i, j int) bool {
		if keys[i].path == keys[j].path {
			return keys[i].sub < keys[j].sub
		}
		return keys[i].path < keys[j].path
	})
	for _, key := range keys {
		if _, ok := seen[key]; ok {
			continue
		}
		mounts = append(mounts, byKey[key])
		seen[key] = struct{}{}
	}

	return mounts
}
