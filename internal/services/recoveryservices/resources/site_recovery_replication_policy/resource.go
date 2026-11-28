// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_replication_policy

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages an Azure Site Recovery replication policy within a recovery vault. Replication policies define the frequency at which recovery points are created and how long they are stored.",
		Create:      resourceSiteRecoveryReplicationPolicyCreate,
		Read:        resourceSiteRecoveryReplicationPolicyRead,
		Update:      resourceSiteRecoveryReplicationPolicyUpdate,
		Delete:      resourceSiteRecoveryReplicationPolicyDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := replicationpolicies.ParseReplicationPolicyID(id)
			return err
		}),
		CustomizeDiff: resourceSiteRecoveryReplicationPolicyCustomDiff,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: siteRecoveryReplicationPolicySchema(),
	}
}
