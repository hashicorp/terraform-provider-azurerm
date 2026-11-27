// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_sql_pool_extended_auditing_policy

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseSqlPoolExtendedAuditingPolicyCreateUpdate,
		Read:   resourceSynapseSqlPoolExtendedAuditingPolicyRead,
		Update: resourceSynapseSqlPoolExtendedAuditingPolicyCreateUpdate,
		Delete: resourceSynapseSqlPoolExtendedAuditingPolicyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SqlPoolExtendedAuditingPolicyID(id)
			return err
		}),

		Schema: synapseSqlPoolExtendedAuditingPolicySchema(),
	}
}
