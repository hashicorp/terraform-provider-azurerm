// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_protection_container_mapping

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainermappings"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSiteRecoveryContainerMappingCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	fabricName := d.Get("recovery_fabric_name").(string)
	policyId := d.Get("recovery_replication_policy_id").(string)
	protectionContainerName := d.Get("recovery_source_protection_container_name").(string)
	targetContainerId := d.Get("recovery_target_protection_container_id").(string)
	name := d.Get("name").(string)

	client := meta.(*clients.Client).RecoveryServices.ContainerMappingClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := replicationprotectioncontainermappings.NewReplicationProtectionContainerMappingID(subscriptionId, resGroup, vaultName, fabricName, protectionContainerName, name)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing site recovery protection container mapping %s (fabric %s, container %s): %+v", name, fabricName, protectionContainerName, err)
			}
		}

		if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
			return tf.ImportAsExistsError("azurerm_site_recovery_protection_container_mapping", *existing.Model.Id)
		}
	}

	parameters := replicationprotectioncontainermappings.CreateProtectionContainerMappingInput{
		Properties: &replicationprotectioncontainermappings.CreateProtectionContainerMappingInputProperties{
			TargetProtectionContainerId: &targetContainerId,
			PolicyId:                    &policyId,
			ProviderSpecificInput:       replicationprotectioncontainermappings.A2AContainerMappingInput{},
		},
	}

	autoUpdateEnabledValue, automationAccountArmId, authType := expandAutoUpdateSettings(d.Get("automatic_update").([]any))

	if autoUpdateEnabledValue == replicationprotectioncontainermappings.AgentAutoUpdateStatusEnabled {
		parameters.Properties.ProviderSpecificInput = replicationprotectioncontainermappings.A2AContainerMappingInput{
			AgentAutoUpdateStatus:               &autoUpdateEnabledValue,
			AutomationAccountArmId:              automationAccountArmId,
			AutomationAccountAuthenticationType: authType,
		}
	}

	if err := client.CreateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating site recovery protection container mapping %s (vault %s): %+v", name, vaultName, err)
	}
	d.SetId(id.ID())

	return resourceSiteRecoveryContainerMappingRead(d, meta)
}

func expandAutoUpdateSettings(input []any) (enabled replicationprotectioncontainermappings.AgentAutoUpdateStatus, automationAccountId *string, authType *replicationprotectioncontainermappings.AutomationAccountAuthenticationType) {
	if len(input) == 0 {
		return replicationprotectioncontainermappings.AgentAutoUpdateStatusDisabled, nil, nil
	}

	v := input[0].(map[string]any)

	var accountIdOutput *string
	if accountId := v["automation_account_id"].(string); accountId != "" {
		accountIdOutput = &accountId
	}

	return replicationprotectioncontainermappings.AgentAutoUpdateStatusEnabled, accountIdOutput, pointer.ToEnum[replicationprotectioncontainermappings.AutomationAccountAuthenticationType](v["authentication_type"].(string))
}
