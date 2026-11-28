// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_fabric

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationfabrics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages an Azure Site Recovery Replication Fabric within a Recovery Services Vault. Only Azure fabrics are supported at this time. Replication Fabrics serve as a container within an Azure region for other Site Recovery resources such as protection containers, protected items, network mappings.",
		Create:      resourceSiteRecoveryFabricCreate,
		Read:        resourceSiteRecoveryFabricRead,
		Delete:      resourceSiteRecoveryFabricDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := replicationfabrics.ParseReplicationFabricID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(45 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: siteRecoveryFabricSchema(),
	}
}
