// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceBackupPolicyVm() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceBackupPolicyVmRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: dataSourceBackupPolicyVmSchema(),
	}
}

func dataSourceBackupPolicyVmRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectionPoliciesClient
	subscriptionid := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := protectionpolicies.NewBackupPolicyID(subscriptionid, d.Get("resource_group_name").(string), d.Get("recovery_vault_name").(string), d.Get("name").(string))

	log.Printf("[DEBUG] Reading %s", id)

	protectionPolicy, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(protectionPolicy.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	if protectionPolicy.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}

	d.SetId(id.ID())

	return nil
}

func dataSourceBackupPolicyVmSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"recovery_vault_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.RecoveryServicesVaultName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}
