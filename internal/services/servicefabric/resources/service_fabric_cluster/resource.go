// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package service_fabric_cluster

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabric/2021-06-01/cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceServiceFabricClusterCreateUpdate,
		Read:   resourceServiceFabricClusterRead,
		Update: resourceServiceFabricClusterCreateUpdate,
		Delete: resourceServiceFabricClusterDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := cluster.ParseClusterID(id)
			return err
		}),

		Schema: serviceFabricClusterSchema(),
	}
}
