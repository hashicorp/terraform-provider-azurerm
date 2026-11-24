// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_dapr_component

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
)

type ContainerAppEnvironmentDaprComponentModel struct {
	Name                 string                 `tfschema:"name"`
	ManagedEnvironmentId string                 `tfschema:"container_app_environment_id"`
	ComponentType        string                 `tfschema:"component_type"`
	Version              string                 `tfschema:"version"`
	IgnoreErrors         bool                   `tfschema:"ignore_errors"`
	InitTimeout          string                 `tfschema:"init_timeout"`
	Secrets              []helpers.DaprSecret   `tfschema:"secret"`
	Scopes               []string               `tfschema:"scopes"`
	Metadata             []helpers.DaprMetadata `tfschema:"metadata"`
}

func (r Resource) ModelObject() any {
	return &ContainerAppEnvironmentDaprComponentModel{}
}
