// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vmware_private_cloud

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/privateclouds"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVmwarePrivateCloudCreate,
		Read:   resourceVmwarePrivateCloudRead,
		Update: resourceVmwarePrivateCloudUpdate,
		Delete: resourceVmwarePrivateCloudDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(10 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(10 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(10 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := privateclouds.ParsePrivateCloudID(id)
			return err
		}),

		Schema: vmwarePrivateCloudSchema(),
	}
}
