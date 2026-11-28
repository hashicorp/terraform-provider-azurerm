// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vmware_cluster

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVmwareClusterCreate,
		Read:   resourceVmwareClusterRead,
		Update: resourceVmwareClusterUpdate,
		Delete: resourceVmwareClusterDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(5 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(5 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(5 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := clusters.ParseClusterID(id)
			return err
		}),

		Schema: vmwareClusterSchema(),
	}
}
