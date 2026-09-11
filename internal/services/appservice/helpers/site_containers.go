// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
)

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

// ValidateSiteContainerEnvironmentVariables returns an error if the same `name` is used by more than one
// `environment_variable` block. Because `environment_variable` is a Set hashed on all of its attributes, two
// blocks that share a `name` but differ on `app_setting_name` hash to distinct elements and would otherwise be
// accepted, producing an ambiguous mapping for a single container environment variable name.
func ValidateSiteContainerEnvironmentVariables(input []SiteContainerEnvironmentVariable) error {
	seen := make(map[string]struct{}, len(input))
	for _, env := range input {
		if _, ok := seen[env.Name]; ok {
			return fmt.Errorf("`environment_variable` blocks must have unique `name` values, but `name` %q is defined more than once", env.Name)
		}
		seen[env.Name] = struct{}{}
	}

	return nil
}

func ExpandSiteContainerEnvironmentVariables(input []SiteContainerEnvironmentVariable) *[]webapps.EnvironmentVariable {
	envs := make([]webapps.EnvironmentVariable, 0, len(input))
	for _, env := range input {
		envs = append(envs, webapps.EnvironmentVariable{
			Name:  env.Name,
			Value: env.AppSettingName,
		})
	}

	return &envs
}

func ExpandSiteContainerVolumeMounts(input []SiteContainerVolumeMount) *[]webapps.VolumeMount {
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

func FlattenSiteContainerEnvironmentVariables(input *[]webapps.EnvironmentVariable) []SiteContainerEnvironmentVariable {
	envs := make([]SiteContainerEnvironmentVariable, 0)
	if input == nil {
		return envs
	}

	for _, env := range *input {
		envs = append(envs, SiteContainerEnvironmentVariable{
			Name:           env.Name,
			AppSettingName: env.Value,
		})
	}

	return envs
}

func FlattenSiteContainerVolumeMounts(input *[]webapps.VolumeMount) []SiteContainerVolumeMount {
	mounts := make([]SiteContainerVolumeMount, 0)
	if input == nil {
		return mounts
	}

	for _, m := range *input {
		mounts = append(mounts, SiteContainerVolumeMount{
			ContainerMountPath: m.ContainerMountPath,
			Data:               pointer.From(m.Data),
			ReadOnly:           pointer.From(m.ReadOnly),
			VolumeSubPath:      m.VolumeSubPath,
		})
	}

	return mounts
}
