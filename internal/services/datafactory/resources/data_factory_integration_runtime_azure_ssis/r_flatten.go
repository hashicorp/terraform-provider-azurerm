// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_integration_runtime_azure_ssis

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenDataFactoryIntegrationRuntimeAzureSsisVnetIntegration(vnetProperties *integrationruntimes.IntegrationRuntimeVNetProperties) []any {
	if vnetProperties == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"vnet_id":     pointer.From(vnetProperties.VNetId),
			"subnet_id":   pointer.From(vnetProperties.SubnetId),
			"subnet_name": pointer.From(vnetProperties.Subnet),
			"public_ips":  pluginsdk.FlattenSlice(vnetProperties.PublicIPs),
		},
	}
}

func flattenDataFactoryIntegrationRuntimeAzureSsisCatalogInfo(ssisProperties *integrationruntimes.IntegrationRuntimeSsisCatalogInfo, d *pluginsdk.ResourceData) []any {
	if ssisProperties == nil {
		return []any{}
	}

	var administratorPassword string

	var pricingTier, elasticPoolName string
	elasticPoolName, elasticPoolNameMatched := parseDataFactoryIntegrationRuntimeElasticPool(pointer.FromEnum(ssisProperties.CatalogPricingTier))
	if !elasticPoolNameMatched {
		pricingTier = pointer.FromEnum(ssisProperties.CatalogPricingTier)
	}

	// read back
	if adminPassword, ok := d.GetOk("catalog_info.0.administrator_password"); ok {
		administratorPassword = adminPassword.(string)
	}

	return []any{
		map[string]any{
			"server_endpoint":        pointer.From(ssisProperties.CatalogServerEndpoint),
			"pricing_tier":           pricingTier,
			"elastic_pool_name":      elasticPoolName,
			"administrator_login":    pointer.From(ssisProperties.CatalogAdminUserName),
			"administrator_password": administratorPassword,
			"dual_standby_pair_name": pointer.From(ssisProperties.DualStandbyPairName),
		},
	}
}

func flattenDataFactoryIntegrationRuntimeAzureSsisProxy(input *integrationruntimes.IntegrationRuntimeDataProxyProperties) []any {
	if input == nil {
		return []any{}
	}

	var selfHostedIntegrationRuntimeName, stagingStorageLinkedServiceName string
	if input.ConnectVia != nil {
		selfHostedIntegrationRuntimeName = pointer.From(input.ConnectVia.ReferenceName)
	}
	if input.StagingLinkedService != nil {
		stagingStorageLinkedServiceName = pointer.From(input.StagingLinkedService.ReferenceName)
	}
	return []any{
		map[string]any{
			"path":                                 pointer.From(input.Path),
			"self_hosted_integration_runtime_name": selfHostedIntegrationRuntimeName,
			"staging_storage_linked_service_name":  stagingStorageLinkedServiceName,
		},
	}
}

func flattenDataFactoryIntegrationRuntimeUserAssignedCredential(credentialProperties *integrationruntimes.CredentialReference) *string {
	if credentialProperties == nil {
		return nil
	}

	return &credentialProperties.ReferenceName
}

func flattenDataFactoryIntegrationRuntimeAzureSsisCustomSetupScript(customSetupScriptProperties *integrationruntimes.IntegrationRuntimeCustomSetupScriptProperties, d *pluginsdk.ResourceData) []any {
	if customSetupScriptProperties == nil {
		return []any{}
	}

	customSetupScript := map[string]string{
		"blob_container_uri": pointer.From(customSetupScriptProperties.BlobContainerUri),
	}

	if sasToken, ok := d.GetOk("custom_setup_script.0.sas_token"); ok {
		customSetupScript["sas_token"] = sasToken.(string)
	}

	return []any{customSetupScript}
}

func flattenDataFactoryIntegrationRuntimeAzureSsisPackageStore(input *[]integrationruntimes.PackageStore) []any {
	if input == nil {
		return []any{}
	}

	result := make([]any, 0)
	for _, item := range *input {
		result = append(result, map[string]any{
			"name":                item.Name,
			"linked_service_name": pointer.From(item.PackageStoreLinkedService.ReferenceName),
		})
	}
	return result
}

func flattenDataFactoryIntegrationRuntimeAzureSsisExpressCustomSetUp(input *[]integrationruntimes.CustomSetupBase, d *pluginsdk.ResourceData) []any {
	if input == nil {
		return []any{}
	}

	// retrieve old state
	oldState := make(map[string]any)
	if arr := d.Get("express_custom_setup").([]any); len(arr) > 0 {
		oldState = arr[0].(map[string]any)
	}
	oldComponents := make([]any, 0)
	if rawComponent, ok := oldState["component"]; ok {
		if v := rawComponent.([]any); len(v) > 0 {
			oldComponents = v
		}
	}
	oldCmdKey := make([]any, 0)
	if rawCmdKey, ok := oldState["command_key"]; ok {
		if v := rawCmdKey.([]any); len(v) > 0 {
			oldCmdKey = v
		}
	}

	env := make(map[string]any)
	powershellVersion := ""
	components := make([]any, 0)
	cmdkeys := make([]any, 0)
	for _, item := range *input {
		switch v := item.(type) {
		case integrationruntimes.AzPowerShellSetup:
			powershellVersion = v.TypeProperties.Version
		case integrationruntimes.ComponentSetup:
			name := v.TypeProperties.ComponentName

			var keyVaultLicense *integrationruntimes.AzureKeyVaultSecretReference
			if license, ok := v.TypeProperties.LicenseKey.(integrationruntimes.AzureKeyVaultSecretReference); ok {
				keyVaultLicense = &license
			}

			components = append(components, map[string]any{
				"name":              name,
				"key_vault_license": flattenDataFactoryIntegrationRuntimeAzureSsisKeyVaultSecretReference(keyVaultLicense),
				"license": readBackSensitiveValue(oldComponents, "license", map[string]string{
					"name": name,
				}),
			})
		case integrationruntimes.EnvironmentVariableSetup:
			env[v.TypeProperties.VariableName] = v.TypeProperties.VariableValue
		case integrationruntimes.CmdkeySetup:
			var name, userName string
			if v.TypeProperties.TargetName != nil {
				if v, ok := v.TypeProperties.TargetName.(string); ok {
					name = v
				}
			}
			if v.TypeProperties.UserName != nil {
				if v, ok := v.TypeProperties.UserName.(string); ok {
					userName = v
				}
			}
			var keyVaultPassword *integrationruntimes.AzureKeyVaultSecretReference
			if v.TypeProperties.Password != nil {
				if reference, ok := v.TypeProperties.Password.(integrationruntimes.AzureKeyVaultSecretReference); ok {
					keyVaultPassword = &reference
				}
			}
			cmdkeys = append(cmdkeys, map[string]any{
				"target_name": name,
				"user_name":   userName,
				"password": readBackSensitiveValue(oldCmdKey, "password", map[string]string{
					"target_name": name,
					"user_name":   userName,
				}),
				"key_vault_password": flattenDataFactoryIntegrationRuntimeAzureSsisKeyVaultSecretReference(keyVaultPassword),
			})
		}
	}

	return []any{
		map[string]any{
			"environment":        env,
			"powershell_version": powershellVersion,
			"component":          components,
			"command_key":        cmdkeys,
		},
	}
}

func flattenDataFactoryIntegrationRuntimeAzureSsisKeyVaultSecretReference(input *integrationruntimes.AzureKeyVaultSecretReference) []any {
	if input == nil {
		return []any{}
	}
	var secretName, secretVersion string
	if input.SecretName != nil {
		if v, ok := input.SecretName.(string); ok {
			secretName = v
		}
	}
	if input.SecretVersion != nil {
		if v, ok := (*input.SecretVersion).(string); ok {
			secretVersion = v
		}
	}
	return []any{
		map[string]any{
			"linked_service_name": input.Store.ReferenceName,
			"parameters":          pointer.From(input.Store.Parameters),
			"secret_name":         secretName,
			"secret_version":      secretVersion,
		},
	}
}

func flattenDataFactoryIntegrationRuntimeCustomerVnetIntegration(input *integrationruntimes.IntegrationRuntimeCustomerVirtualNetwork) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"subnet_id": pointer.From(input.SubnetId),
		},
	}
}

func flattenDataFactoryIntegrationRuntimeAzureSsisPipelineExternalComputeScaleProperties(input *integrationruntimes.PipelineExternalComputeScaleProperties) []any {
	if input == nil {
		return []any{}
	}
	return []any{
		map[string]any{
			"number_of_external_nodes": pointer.From(input.NumberOfPipelineNodes),
			"number_of_pipeline_nodes": pointer.From(input.NumberOfPipelineNodes),
			"time_to_live":             pointer.From(input.TimeToLive),
		},
	}
}

func flattenDataFactoryIntegrationRuntimeAzureSsisCopyComputeScale(input *integrationruntimes.CopyComputeScaleProperties) []any {
	if input == nil {
		return []any{}
	}
	return []any{
		map[string]any{
			"data_integration_unit": pointer.From(input.DataIntegrationUnit),
			"time_to_live":          pointer.From(input.TimeToLive),
		},
	}
}
