// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_sql_pool

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseSqlPoolUpdate(d *pluginsdk.ResourceData, meta any) error {
	sqlClient := meta.(*clients.Client).Synapse.SqlPoolClient
	geoBackUpClient := meta.(*clients.Client).Synapse.SqlPoolGeoBackupPoliciesClient
	sqlPTDEClient := meta.(*clients.Client).Synapse.SqlPoolTransparentDataEncryptionClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlPoolID(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("data_encrypted") {
		status := synapse.TransparentDataEncryptionStatusDisabled
		if d.Get("data_encrypted").(bool) {
			status = synapse.TransparentDataEncryptionStatusEnabled
		}

		parameter := synapse.TransparentDataEncryption{
			TransparentDataEncryptionProperties: &synapse.TransparentDataEncryptionProperties{
				Status: status,
			},
		}

		if _, err := sqlPTDEClient.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, parameter); err != nil {
			return fmt.Errorf("updating `data_encrypted`: %+v", err)
		}
	}

	if d.HasChange("geo_backup_policy_enabled") {
		state := synapse.GeoBackupPolicyStateEnabled
		if !d.Get("geo_backup_policy_enabled").(bool) {
			state = synapse.GeoBackupPolicyStateDisabled
		}

		geoBackupParams := synapse.GeoBackupPolicy{
			GeoBackupPolicyProperties: &synapse.GeoBackupPolicyProperties{
				State: state,
			},
		}

		if _, err := geoBackUpClient.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, geoBackupParams); err != nil {
			return fmt.Errorf("updating `geo_backup_policy_enabled`: %+v", err)
		}
	}

	if d.HasChanges("sku_name", "tags") {
		sqlPoolInfo := synapse.SQLPoolPatchInfo{
			Sku: &synapse.Sku{
				Name: pointer.To(d.Get("sku_name").(string)),
			},
			Tags: tags.Expand(d.Get("tags").(map[string]any)),
		}

		if _, err := sqlClient.Update(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, sqlPoolInfo); err != nil {
			return fmt.Errorf("updating %s: %+v", *id, err)
		}

		// wait for sku scale completion
		if d.HasChange("sku_name") {
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}

			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{
					"Scaling",
				},
				Target: []string{
					"Online",
				},
				Refresh:                   synapseSqlPoolScaleStateRefreshFunc(ctx, sqlClient, id.ResourceGroup, id.WorkspaceName, id.Name),
				MinTimeout:                5 * time.Second,
				ContinuousTargetOccurence: 3,
				Timeout:                   time.Until(deadline),
			}

			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for scaling of %s: %+v", *id, err)
			}
		}
	}

	return resourceSynapseSqlPoolRead(d, meta)
}
