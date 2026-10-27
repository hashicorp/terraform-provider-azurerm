// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package qumulo_file_system

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/qumulostorage/2026-04-16/filesystems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.ResourceWithUpdate = Resource{}

type Resource struct{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return filesystems.ValidateFileSystemID
}

func (r Resource) ResourceType() string {
	return "azurerm_qumulo_file_system"
}
