// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_integration_runtime_azure_ssis

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandDataFactoryIntegrationRuntimeAzureSsisVirtualNetwork(input []any) *integrationruntimes.IntegrationRuntimeVNetProperties {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)
	result := &integrationruntimes.IntegrationRuntimeVNetProperties{}

	if vnetID := v["vnet_id"].(string); vnetID != "" {
		result.VNetId = pointer.To(vnetID)
		result.Subnet = pointer.To(v["subnet_name"].(string))
	}

	if subnetID := v["subnet_id"].(string); subnetID != "" {
		result.SubnetId = pointer.To(subnetID)
	}

	if publicIPs := v["public_ips"].([]any); len(publicIPs) > 0 {
		result.PublicIPs = pluginsdk.ExpandStringSlice(publicIPs)
	}

	return result
}

func expandDataFactoryIntegrationRuntimeAzureSsisCopyComputeScale(input []any) *integrationruntimes.CopyComputeScaleProperties {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)
	result := &integrationruntimes.CopyComputeScaleProperties{}

	if diUnit := v["data_integration_unit"].(int); diUnit != 0 {
		result.DataIntegrationUnit = pointer.To(int64(diUnit))
	}

	if ttl := v["time_to_live"].(int); ttl != 0 {
		result.TimeToLive = pointer.To(int64(ttl))
	}

	return result
}

func expandDataFactoryIntegrationRuntimeAzureSsisPipelineExternalComputeScale(input []any) *integrationruntimes.PipelineExternalComputeScaleProperties {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)
	result := &integrationruntimes.PipelineExternalComputeScaleProperties{}

	if numExternalNodes := v["number_of_external_nodes"].(int); numExternalNodes != 0 {
		result.NumberOfExternalNodes = pointer.To(int64(numExternalNodes))
	}

	if numPipelineNodes := v["number_of_pipeline_nodes"].(int); numPipelineNodes != 0 {
		result.NumberOfPipelineNodes = pointer.To(int64(numPipelineNodes))
	}

	if timeToLive := v["time_to_live"].(int); timeToLive != 0 {
		result.TimeToLive = pointer.To(int64(timeToLive))
	}

	return result
}

func expandDataFactoryIntegrationRuntimeAzureSsisCredential(input string) *integrationruntimes.CredentialReference {
	if input == "" {
		return nil
	}

	return &integrationruntimes.CredentialReference{
		ReferenceName: input,
		Type:          integrationruntimes.CredentialReferenceTypeCredentialReference,
	}
}

func expandDataFactoryIntegrationRuntimeAzureSsisCatalogInfo(input []any) *integrationruntimes.IntegrationRuntimeSsisCatalogInfo {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)
	result := &integrationruntimes.IntegrationRuntimeSsisCatalogInfo{
		CatalogServerEndpoint: pointer.To(v["server_endpoint"].(string)),
		CatalogPricingTier:    pointer.ToEnum[integrationruntimes.IntegrationRuntimeSsisCatalogPricingTier](v["pricing_tier"].(string)),
	}

	if epName := v["elastic_pool_name"].(string); epName != "" {
		result.CatalogPricingTier = pointer.ToEnum[integrationruntimes.IntegrationRuntimeSsisCatalogPricingTier](formatDataFactoryIntegrationRuntimeElasticPool(epName))
	}

	if adminUsername := v["administrator_login"].(string); adminUsername != "" {
		result.CatalogAdminUserName = pointer.To(adminUsername)
	}

	if adminPassword := v["administrator_password"].(string); adminPassword != "" {
		result.CatalogAdminPassword = &integrationruntimes.SecureString{
			Value: adminPassword,
			Type:  string(helpers.SecretTypeSecureString),
		}
	}

	if dspName := v["dual_standby_pair_name"].(string); dspName != "" {
		result.DualStandbyPairName = pointer.To(dspName)
	}

	return result
}

func expandDataFactoryIntegrationRuntimeAzureSsisCustomSetupScript(input []any) *integrationruntimes.IntegrationRuntimeCustomSetupScriptProperties {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)
	return &integrationruntimes.IntegrationRuntimeCustomSetupScriptProperties{
		BlobContainerUri: pointer.To(v["blob_container_uri"].(string)),
		SasToken: &integrationruntimes.SecureString{
			Value: v["sas_token"].(string),
			Type:  string(helpers.SecretTypeSecureString),
		},
	}
}

func expandDataFactoryIntegrationRuntimeAzureSsisProxy(input []any) *integrationruntimes.IntegrationRuntimeDataProxyProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	raw := input[0].(map[string]any)

	result := &integrationruntimes.IntegrationRuntimeDataProxyProperties{
		ConnectVia: &integrationruntimes.EntityReference{
			Type:          pointer.To(integrationruntimes.IntegrationRuntimeEntityReferenceTypeIntegrationRuntimeReference),
			ReferenceName: pointer.To(raw["self_hosted_integration_runtime_name"].(string)),
		},
		StagingLinkedService: &integrationruntimes.EntityReference{
			Type:          pointer.To(integrationruntimes.IntegrationRuntimeEntityReferenceTypeLinkedServiceReference),
			ReferenceName: pointer.To(raw["staging_storage_linked_service_name"].(string)),
		},
	}
	if path := raw["path"].(string); len(path) > 0 {
		result.Path = pointer.To(path)
	}
	return result
}

func expandDataFactoryIntegrationRuntimeAzureSsisExpressCustomSetUp(input []any) *[]integrationruntimes.CustomSetupBase {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	raw := input[0].(map[string]any)

	result := make([]integrationruntimes.CustomSetupBase, 0)
	if env := raw["environment"].(map[string]any); len(env) > 0 {
		for k, v := range env {
			result = append(result, &integrationruntimes.EnvironmentVariableSetup{
				Type: string(helpers.CustomSetupTypeEnvironmentVariableSetup),
				TypeProperties: integrationruntimes.EnvironmentVariableSetupTypeProperties{
					VariableName:  k,
					VariableValue: v.(string),
				},
			})
		}
	}

	if powershellVersion := raw["powershell_version"].(string); powershellVersion != "" {
		result = append(result, &integrationruntimes.AzPowerShellSetup{
			Type: string(helpers.CustomSetupTypeAzPowerShellSetup),
			TypeProperties: integrationruntimes.AzPowerShellSetupTypeProperties{
				Version: powershellVersion,
			},
		})
	}

	if components := raw["component"].([]any); len(components) > 0 {
		for _, item := range components {
			raw := item.(map[string]any)

			var license integrationruntimes.SecretBase
			if v := raw["license"].(string); v != "" {
				license = &integrationruntimes.SecureString{
					Type:  string(helpers.SecretTypeSecureString),
					Value: v,
				}
			} else {
				license = expandDataFactoryIntegrationRuntimeAzureSsisKeyVaultSecretReference(raw["key_vault_license"].([]any))
			}

			result = append(result, &integrationruntimes.ComponentSetup{
				Type: string(helpers.CustomSetupTypeComponentSetup),
				TypeProperties: integrationruntimes.LicensedComponentSetupTypeProperties{
					ComponentName: raw["name"].(string),
					LicenseKey:    license,
				},
			})
		}
	}
	if cmdKeys := raw["command_key"].([]any); len(cmdKeys) > 0 {
		for _, item := range cmdKeys {
			raw := item.(map[string]any)

			var password integrationruntimes.SecretBase
			if v := raw["password"].(string); v != "" {
				password = &integrationruntimes.SecureString{
					Type:  string(helpers.SecretTypeSecureString),
					Value: v,
				}
			} else {
				password = expandDataFactoryIntegrationRuntimeAzureSsisKeyVaultSecretReference(raw["key_vault_password"].([]any))
			}

			result = append(result, &integrationruntimes.CmdkeySetup{
				Type: string(helpers.CustomSetupTypeCmdkeySetup),
				TypeProperties: integrationruntimes.CmdkeySetupTypeProperties{
					TargetName: pointer.To(raw["target_name"].(string)),
					UserName:   pointer.To(raw["user_name"].(string)),
					Password:   password,
				},
			})
		}
	}

	return &result
}

func expandDataFactoryIntegrationRuntimeAzureSsisPackageStore(input []any) *[]integrationruntimes.PackageStore {
	if len(input) == 0 {
		return nil
	}

	result := make([]integrationruntimes.PackageStore, 0)
	for _, item := range input {
		raw := item.(map[string]any)
		result = append(result, integrationruntimes.PackageStore{
			Name: raw["name"].(string),
			PackageStoreLinkedService: integrationruntimes.EntityReference{
				Type:          pointer.To(integrationruntimes.IntegrationRuntimeEntityReferenceTypeLinkedServiceReference),
				ReferenceName: pointer.To(raw["linked_service_name"].(string)),
			},
		})
	}
	return &result
}

func expandDataFactoryIntegrationRuntimeAzureSsisKeyVaultSecretReference(input []any) *integrationruntimes.AzureKeyVaultSecretReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)
	reference := &integrationruntimes.AzureKeyVaultSecretReference{
		SecretName: raw["secret_name"].(string),
		Store: integrationruntimes.LinkedServiceReference{
			Type:          integrationruntimes.TypeLinkedServiceReference,
			ReferenceName: raw["linked_service_name"].(string),
		},
		Type: string(helpers.SecretTypeAzureKeyVaultSecret),
	}
	if v := raw["secret_version"].(string); v != "" {
		reference.SecretVersion = pointer.To(raw["secret_version"])
	}
	if v := raw["parameters"].(map[string]any); len(v) > 0 {
		reference.Store.Parameters = &v
	}
	return reference
}

func expandDataFactoryIntegrationRuntimeCustomerVirtualNetwork(input []any) *integrationruntimes.IntegrationRuntimeCustomerVirtualNetwork {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	raw := input[0].(map[string]any)
	return &integrationruntimes.IntegrationRuntimeCustomerVirtualNetwork{
		SubnetId: pointer.To(raw["subnet_id"].(string)),
	}
}
