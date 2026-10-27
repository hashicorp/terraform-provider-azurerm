// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_virtual_machine

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -parent-id "virtual_machine_id"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlVirtualMachineCreateUpdate,
		Read:   resourceMsSqlVirtualMachineRead,
		Update: resourceMsSqlVirtualMachineCreateUpdate,
		Delete: resourceMsSqlVirtualMachineDelete,

		CustomizeDiff: pluginsdk.CustomizeDiffShim(resourceMsSqlVirtualMachineCustomDiff),

		Importer: pluginsdk.ImporterValidatingIdentity(&sqlvirtualmachines.SqlVirtualMachineId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&sqlvirtualmachines.SqlVirtualMachineId{}),
		},
		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: mssqlVirtualMachineSchema(),
	}
}
