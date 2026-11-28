// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_protection_container

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages an Azure Site Recovery protection container. Protection containers serve as containers for replicated VMs and belong to a single region / recovery fabric. Protection containers can contain more than one replicated VM. To replicate a VM, a container must exist in both the source and target Azure regions.",
		Create:      resourceSiteRecoveryProtectionContainerCreate,
		Read:        resourceSiteRecoveryProtectionContainerRead,
		Delete:      resourceSiteRecoveryProtectionContainerDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := replicationprotectioncontainers.ParseReplicationProtectionContainerID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: siteRecoveryProtectionContainerSchema(),
	}
}
