// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_global_vm_shutdown_schedule

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/globalschedules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDevTestGlobalVMShutdownScheduleCreateUpdate,
		Read:   resourceDevTestGlobalVMShutdownScheduleRead,
		Update: resourceDevTestGlobalVMShutdownScheduleCreateUpdate,
		Delete: resourceDevTestGlobalVMShutdownScheduleDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := globalschedules.ParseScheduleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: devTestGlobalVmShutdownScheduleSchema(),
	}
}
