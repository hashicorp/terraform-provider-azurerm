// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_replicated_vm

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSiteRecoveryReplicatedItemCreate,
		Read:   resourceSiteRecoveryReplicatedItemRead,
		Update: resourceSiteRecoveryReplicatedItemUpdate,
		Delete: resourceSiteRecoveryReplicatedItemDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := replicationprotecteditems.ParseReplicationProtectedItemID(id)
			return err
		}),
		CustomizeDiff: resourceSiteRecoveryReplicatedVMCustomizeDiff,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(180 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(80 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(80 * time.Minute),
		},

		Schema: siteRecoveryReplicatedVmSchema(),
	}
}
