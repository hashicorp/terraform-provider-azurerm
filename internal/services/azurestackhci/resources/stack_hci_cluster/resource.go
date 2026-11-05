// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_cluster

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmStackHCIClusterCreate,
		Read:   resourceArmStackHCIClusterRead,
		Update: resourceArmStackHCIClusterUpdate,
		Delete: resourceArmStackHCIClusterDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := clusters.ParseClusterID(id)
			return err
		}),

		Schema: stackHciClusterSchema(),
	}
}
