// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// MainSiteContainerName is the fixed Site Container name used for the single main Site Container that is
// managed inline via `site_config.0.main_site_container`. Additional (non-main) Site Containers are managed
// out of band via the standalone `azurerm_linux_web_app_site_container` resource.
const MainSiteContainerName = "main"

type SiteContainerEnvironmentVariable struct {
	Name  string `tfschema:"name"`
	Value string `tfschema:"value"`
}

type SiteContainerVolumeMount struct {
	ContainerMountPath string `tfschema:"container_mount_path"`
	VolumeSubPath      string `tfschema:"volume_sub_path"`
	ReadOnly           bool   `tfschema:"read_only"`
}

type MainSiteContainer struct {
	Image                string                             `tfschema:"image"`
	TargetPort           string                             `tfschema:"target_port"`
	EnvironmentVariables []SiteContainerEnvironmentVariable `tfschema:"environment_variable"`
	VolumeMounts         []SiteContainerVolumeMount         `tfschema:"volume_mount"`
}

func MainSiteContainerSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"image": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The image to use for the main Site Container, e.g. `mcr.microsoft.com/appsvc/staticsite:latest`.",
				},

				"target_port": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(`^([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`),
						"target_port must be a valid TCP port number between 1 and 65535",
					),
					Description: "The port the main Site Container listens on.",
				},

				"environment_variable": siteContainerEnvironmentVariableSchema(),

				"volume_mount": siteContainerVolumeMountSchema(),
			},
		},
	}
}

func MainSiteContainerSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"image": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"target_port": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"environment_variable": siteContainerEnvironmentVariableSchemaComputed(),

				"volume_mount": siteContainerVolumeMountSchemaComputed(),
			},
		},
	}
}

func siteContainerEnvironmentVariableSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The name of the environment variable.",
				},

				"value": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The value of the environment variable.",
				},
			},
		},
	}
}

func siteContainerEnvironmentVariableSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"value": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func siteContainerVolumeMountSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"container_mount_path": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The path within the container at which the volume should be mounted.",
				},

				"volume_sub_path": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The path within the volume to mount.",
				},

				"read_only": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should this volume mount be read-only?",
				},
			},
		},
	}
}

func siteContainerVolumeMountSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"container_mount_path": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"volume_sub_path": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"read_only": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
			},
		},
	}
}

// ExpandMainSiteContainer expands the `main_site_container` block into the `SiteContainerProperties` used to
// create or update the main Site Container via the dedicated Site Containers API. Returns nil when no
// `main_site_container` block is configured.
func ExpandMainSiteContainer(input []MainSiteContainer) *webapps.SiteContainerProperties {
	if len(input) == 0 {
		return nil
	}

	container := input[0]

	props := &webapps.SiteContainerProperties{
		Image:  container.Image,
		IsMain: true,
	}

	if container.TargetPort != "" {
		props.TargetPort = pointer.To(container.TargetPort)
	}

	if envVars := expandSiteContainerEnvironmentVariables(container.EnvironmentVariables); len(envVars) > 0 {
		props.EnvironmentVariables = &envVars
	}

	if volumeMounts := expandSiteContainerVolumeMounts(container.VolumeMounts); len(volumeMounts) > 0 {
		props.VolumeMounts = &volumeMounts
	}

	return props
}

// FlattenMainSiteContainer flattens the main Site Container, as returned by the Site Containers API, into the
// `main_site_container` block. Returns an empty slice when the given Site Container is nil, i.e. when the Linux
// Web App is not currently running in Site Containers mode.
func FlattenMainSiteContainer(input *webapps.SiteContainer) []MainSiteContainer {
	if input == nil || input.Properties == nil {
		return []MainSiteContainer{}
	}

	props := *input.Properties

	return []MainSiteContainer{
		{
			Image:                props.Image,
			TargetPort:           pointer.From(props.TargetPort),
			EnvironmentVariables: flattenSiteContainerEnvironmentVariables(props.EnvironmentVariables),
			VolumeMounts:         flattenSiteContainerVolumeMounts(props.VolumeMounts),
		},
	}
}

// FindMainSiteContainer lists the Site Containers configured for the given Web App and returns the one marked as
// the main Site Container, if any. Non-main Site Containers are managed by the standalone
// `azurerm_linux_web_app_site_container` resource and are intentionally left untouched.
func FindMainSiteContainer(ctx context.Context, client *webapps.WebAppsClient, id commonids.AppServiceId) (*webapps.SiteContainer, error) {
	containers, err := client.ListSiteContainersComplete(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("listing Site Containers for %s: %+v", id, err)
	}

	for _, container := range containers.Items {
		if container.Properties != nil && container.Properties.IsMain {
			return pointer.To(container), nil
		}
	}

	return nil, nil
}

func expandSiteContainerEnvironmentVariables(input []SiteContainerEnvironmentVariable) []webapps.EnvironmentVariable {
	result := make([]webapps.EnvironmentVariable, 0, len(input))
	for _, v := range input {
		result = append(result, webapps.EnvironmentVariable{
			Name:  v.Name,
			Value: v.Value,
		})
	}

	return result
}

func flattenSiteContainerEnvironmentVariables(input *[]webapps.EnvironmentVariable) []SiteContainerEnvironmentVariable {
	if input == nil {
		return []SiteContainerEnvironmentVariable{}
	}

	result := make([]SiteContainerEnvironmentVariable, 0, len(*input))
	for _, v := range *input {
		result = append(result, SiteContainerEnvironmentVariable{
			Name:  v.Name,
			Value: v.Value,
		})
	}

	return result
}

func expandSiteContainerVolumeMounts(input []SiteContainerVolumeMount) []webapps.VolumeMount {
	result := make([]webapps.VolumeMount, 0, len(input))
	for _, v := range input {
		result = append(result, webapps.VolumeMount{
			ContainerMountPath: v.ContainerMountPath,
			VolumeSubPath:      v.VolumeSubPath,
			ReadOnly:           pointer.To(v.ReadOnly),
		})
	}

	return result
}

func flattenSiteContainerVolumeMounts(input *[]webapps.VolumeMount) []SiteContainerVolumeMount {
	if input == nil {
		return []SiteContainerVolumeMount{}
	}

	result := make([]SiteContainerVolumeMount, 0, len(*input))
	for _, v := range *input {
		result = append(result, SiteContainerVolumeMount{
			ContainerMountPath: v.ContainerMountPath,
			VolumeSubPath:      v.VolumeSubPath,
			ReadOnly:           pointer.From(v.ReadOnly),
		})
	}

	return result
}
