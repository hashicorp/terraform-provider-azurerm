// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_virtual_machine

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachinegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlVirtualMachineCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.VirtualMachinesClient
	vmclient := meta.(*clients.Client).Compute.VirtualMachinesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vmId, err := virtualmachines.ParseVirtualMachineID(d.Get("virtual_machine_id").(string))
	if err != nil {
		return err
	}
	id := sqlvirtualmachines.NewSqlVirtualMachineID(vmId.SubscriptionId, vmId.ResourceGroupName, vmId.VirtualMachineName)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id, sqlvirtualmachines.GetOperationOptions{Expand: pointer.To("*")})
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for present of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_mssql_virtual_machine", id.ID())
			}
		}
	}

	// get location from vm
	respvm, err := vmclient.Get(ctx, *vmId, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", vmId, err)
	}

	if respvm.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", vmId)
	}
	if respvm.Model.Location == "" {
		return fmt.Errorf("retrieving %s: `location` is empty", vmId)
	}
	var sqlVmGroupId string
	if sqlVmGroupId = d.Get("sql_virtual_machine_group_id").(string); sqlVmGroupId != "" {
		parsedVmGroupId, err := sqlvirtualmachinegroups.ParseSqlVirtualMachineGroupIDInsensitively(sqlVmGroupId)
		if err != nil {
			return err
		}
		sqlVmGroupId = parsedVmGroupId.ID()
	}

	sqlInstance, err := expandSqlVirtualMachineSQLInstance(d.Get("sql_instance").([]any))
	if err != nil {
		return fmt.Errorf("expanding `sql_instance`: %+v", err)
	}

	connectivityType := sqlvirtualmachines.ConnectivityType(d.Get("sql_connectivity_type").(string))
	sqlServerLicenseType := sqlvirtualmachines.SqlServerLicenseType(d.Get("sql_license_type").(string))
	autoBackupSettings, err := expandSqlVirtualMachineAutoBackupSettings(d.Get("auto_backup").([]any))
	if err != nil {
		return fmt.Errorf("expanding `auto_backup`: %+v", err)
	}

	parameters := sqlvirtualmachines.SqlVirtualMachine{
		Location: respvm.Model.Location,
		Properties: &sqlvirtualmachines.SqlVirtualMachineProperties{
			AutoBackupSettings:               autoBackupSettings,
			AutoPatchingSettings:             expandSqlVirtualMachineAutoPatchingSettings(d.Get("auto_patching").([]any)),
			AssessmentSettings:               expandSqlVirtualMachineAssessmentSettings(d.Get("assessment").([]any)),
			KeyVaultCredentialSettings:       expandSqlVirtualMachineKeyVaultCredential(d.Get("key_vault_credential").([]any)),
			WsfcDomainCredentials:            expandSqlVirtualMachineWsfcDomainCredentials(d.Get("wsfc_domain_credential").([]any)),
			SqlVirtualMachineGroupResourceId: pointer.To(sqlVmGroupId),
			ServerConfigurationsManagementSettings: &sqlvirtualmachines.ServerConfigurationsManagementSettings{
				AdditionalFeaturesServerConfigurations: &sqlvirtualmachines.AdditionalFeaturesServerConfigurations{
					IsRServicesEnabled: pointer.To(d.Get("r_services_enabled").(bool)),
				},
				SqlConnectivityUpdateSettings: &sqlvirtualmachines.SqlConnectivityUpdateSettings{
					ConnectivityType:      &connectivityType,
					Port:                  pointer.To(int64(d.Get("sql_connectivity_port").(int))),
					SqlAuthUpdatePassword: pointer.To(d.Get("sql_connectivity_update_password").(string)),
					SqlAuthUpdateUserName: pointer.To(d.Get("sql_connectivity_update_username").(string)),
				},
				SqlInstanceSettings: sqlInstance,
			},
			SqlManagement:                pointer.To(sqlvirtualmachines.SqlManagementModeFull),
			SqlServerLicenseType:         &sqlServerLicenseType,
			StorageConfigurationSettings: expandSqlVirtualMachineStorageConfigurationSettings(d.Get("storage_configuration").([]any)),
			VirtualMachineResourceId:     pointer.To(d.Get("virtual_machine_id").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		d.SetId(id.ID())
		if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
			return err
		}
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	// Wait for the auto backup settings to take effect
	// See: https://github.com/Azure/azure-rest-api-specs/issues/12818
	if autoBackup := d.Get("auto_backup"); (d.IsNewResource() && len(autoBackup.([]any)) > 0) || (!d.IsNewResource() && d.HasChange("auto_backup")) {
		log.Printf("[DEBUG] Waiting for SQL Virtual Machine %q AutoBackupSettings to take effect", d.Id())
		stateConf := &pluginsdk.StateChangeConf{
			Pending:                   []string{"Retry", "Pending"},
			Target:                    []string{"Updated"},
			Refresh:                   resourceMsSqlVirtualMachineAutoBackupSettingsRefreshFunc(ctx, client, d),
			MinTimeout:                1 * time.Minute,
			ContinuousTargetOccurence: 2,
		}

		if d.IsNewResource() {
			stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
		} else {
			stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for SQL Virtual Machine %q AutoBackupSettings to take effect: %+v", d.Id(), err)
		}
	}

	// Wait for the auto patching settings to take effect
	// See: https://github.com/Azure/azure-rest-api-specs/issues/12818
	if autoPatching := d.Get("auto_patching"); (d.IsNewResource() && len(autoPatching.([]any)) > 0) || (!d.IsNewResource() && d.HasChange("auto_patching")) {
		log.Printf("[DEBUG] Waiting for SQL Virtual Machine %q AutoPatchingSettings to take effect", d.Id())
		stateConf := &pluginsdk.StateChangeConf{
			Pending:                   []string{"Retry", "Pending"},
			Target:                    []string{"Updated"},
			Refresh:                   resourceMsSqlVirtualMachineAutoPatchingSettingsRefreshFunc(ctx, client, d),
			MinTimeout:                1 * time.Minute,
			ContinuousTargetOccurence: 2,
		}

		if d.IsNewResource() {
			stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
		} else {
			stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for SQL Virtual Machine %q AutoPatchingSettings to take effect: %+v", d.Id(), err)
		}
	}

	return resourceMsSqlVirtualMachineRead(d, meta)
}
