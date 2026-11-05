// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_run_command

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-03-01/virtualmachineruncommands"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Source struct {
	MediaLink                string `tfschema:"media_link"`
	DefaultConfigurationLink string `tfschema:"default_configuration_link"`
}

type Properties struct {
	Output Output `tfschema:"output"`
}

type Output struct {
	AccessSAS string `tfschema:"accessSAS"`
}

var (
	_ sdk.Resource           = Resource{}
	_ sdk.ResourceWithUpdate = Resource{}
)

type Resource struct{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return virtualmachineruncommands.ValidateVirtualMachineRunCommandID
}

func (r Resource) ResourceType() string {
	return "azurerm_virtual_machine_run_command"
}
