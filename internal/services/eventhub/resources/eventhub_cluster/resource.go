// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub_cluster

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubsclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceEventHubClusterCreateUpdate,
		Read:   resourceEventHubClusterRead,
		Update: resourceEventHubClusterCreateUpdate,
		Delete: resourceEventHubClusterDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := eventhubsclusters.ParseClusterID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			// You can't delete a cluster until at least 4 hours have passed from the initial creation.
			Delete: pluginsdk.DefaultTimeout(300 * time.Minute),
		},

		Schema: eventhubClusterSchema(),
	}
}
