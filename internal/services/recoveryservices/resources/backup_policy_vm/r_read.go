// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_policy_vm

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBackupProtectionPolicyVMRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectionPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := protectionpolicies.ParseBackupPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	return resourceBackupProtectionPolicyVMFlatten(d, id, resp.Model)
}

func resourceBackupProtectionPolicyVMWaitForDeletion(ctx context.Context, client *protectionpolicies.ProtectionPoliciesClient, id protectionpolicies.BackupPolicyId, d *pluginsdk.ResourceData) error {
	state := &pluginsdk.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"Found"},
		Target:     []string{"NotFound"},
		Refresh:    resourceBackupProtectionPolicyVMRefreshFunc(ctx, client, id),
		Timeout:    d.Timeout(pluginsdk.TimeoutDelete),
	}

	if _, err := state.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to provision: %+v", id, err)
	}

	return nil
}

func resourceBackupProtectionPolicyVMRefreshFunc(ctx context.Context, client *protectionpolicies.ProtectionPoliciesClient, id protectionpolicies.BackupPolicyId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return resp, "NotFound", nil
			}

			return resp, "Error", fmt.Errorf("making Read request on %s: %+v", id, err)
		}

		return resp, "Found", nil
	}
}
