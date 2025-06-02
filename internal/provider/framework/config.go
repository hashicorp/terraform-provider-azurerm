// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package framework

import (
	"context"
	"os"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	providerfeatures "github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceproviders"
)

type ProviderConfig struct {
	clientBuilder clients.ClientBuilder
	Client        *clients.Client
}

// Load handles the heavy lifting of configuring the provider and handling defaults
func (p *ProviderConfig) Load(ctx context.Context, data *ProviderModel, tfVersion string, diags *diag.Diagnostics) {
	env := &environments.Environment{}
	var err error

	subscriptionId := getEnvStringOrDefault(data.SubscriptionId, "ARM_SUBSCRIPTION_ID", "")
	if subscriptionId == "" {
		diags.Append(diag.NewErrorDiagnostic("Configuring subscription", "`subscription_id` is a required provider property when performing a plan/apply operation"))
		return
	}

	if metadataHost := getEnvStringOrDefault(data.MetaDataHost, "ARM_METADATA_HOSTNAME", ""); metadataHost != "" {
		env, err = environments.FromEndpoint(ctx, metadataHost)
		if err != nil {
			diags.Append(diag.NewErrorDiagnostic("Configuring metadata host", err.Error()))
			return
		}
	} else {
		if v := getEnvStringOrDefault(data.Environment, "ARM_ENVIRONMENT", "public"); v != "" {
			env, err = environments.FromName(v)
			if err != nil {
				diags.Append(diag.NewErrorDiagnostic("Configuring metadata host", err.Error()))
				return
			}
		}
	}

	var clientCertificateData []byte
	if encodedCert := getEnvStringOrDefault(data.ClientCertificate, "ARM_CLIENT_CERTIFICATE", ""); encodedCert != "" {
		clientCertificateData, err = decodeCertificate(encodedCert)
		if err != nil {
			diags.Append(diag.NewErrorDiagnostic("decoding client certificate", err.Error()))
			if diags.HasError() {
				return
			}
		}
	}

	clientSecret, err := getClientSecret(data)
	if err != nil {
		diags.Append(diag.NewErrorDiagnostic("configuring client secret", err.Error()))
		if diags.HasError() {
			return
		}
	}

	oidcToken, err := getOidcToken(data)
	if err != nil {
		diags.Append(diag.NewErrorDiagnostic("", err.Error()))
		return
	}

	enableOIDC := getEnvBoolIfValueAbsent(data.UseOIDC, "ARM_USE_OIDC") || getEnvBoolIfValueAbsent(data.UseAKSWorkloadIdentity, "ARM_USE_AKS_WORKLOAD_IDENTITY")
	auxTenants := getEnvListOfStringsIfAbsent(data.AuxiliaryTenantIds, "ARM_AUXILIARY_TENANT_IDS", ";")

	oidcReqURL := getEnvStringsOrDefault(data.OIDCRequestURL, []string{"ARM_OIDC_REQUEST_URL", "ACTIONS_ID_TOKEN_REQUEST_URL", "SYSTEM_OIDCREQUESTURI"}, "")
	oidcReqToken := getEnvStringsOrDefault(data.OIDCRequestToken, []string{"ARM_OIDC_REQUEST_TOKEN", "ACTIONS_ID_TOKEN_REQUEST_TOKEN", "SYSTEM_ACCESSTOKEN"}, "")

	// ARM_OIDC_AZURE_SERVICE_CONNECTION_ID is to be compatible with `azapi` provider.
	adoPipelineServiceConnectionID := getEnvStringsOrDefault(data.ADOPipelineServiceConnectionID, []string{"ARM_ADO_PIPELINE_SERVICE_CONNECTION_ID", "ARM_OIDC_AZURE_SERVICE_CONNECTION_ID"}, "")

	authConfig := &auth.Credentials{
		Environment:        *env,
		ClientID:           getEnvStringIfValueAbsent(data.ClientId, "ARM_CLIENT_ID"),
		TenantID:           getEnvStringIfValueAbsent(data.TenantId, "ARM_TENANT_ID"),
		AuxiliaryTenantIDs: auxTenants,

		ClientCertificateData:     clientCertificateData,
		ClientCertificatePath:     getEnvStringOrDefault(data.ClientCertificatePath, "ARM_CLIENT_CERTIFICATE_PATH", ""),
		ClientCertificatePassword: getEnvStringOrDefault(data.ClientCertificatePassword, "ARM_CLIENT_CERTIFICATE_PASSWORD", ""),
		ClientSecret:              *clientSecret,

		OIDCAssertionToken:    *oidcToken,
		OIDCTokenRequestURL:   oidcReqURL,
		OIDCTokenRequestToken: oidcReqToken,

		ADOPipelineServiceConnectionID: adoPipelineServiceConnectionID,

		CustomManagedIdentityEndpoint: getEnvStringOrDefault(data.MSIEndpoint, "ARM_MSI_ENDPOINT", ""),

		AzureCliSubscriptionIDHint: subscriptionId,

		EnableAuthenticatingUsingClientCertificate: true,
		EnableAuthenticatingUsingClientSecret:      true,
		EnableAuthenticationUsingOIDC:              enableOIDC,
		EnableAuthenticationUsingGitHubOIDC:        enableOIDC,
		EnableAuthenticationUsingADOPipelineOIDC:   enableOIDC,
		EnableAuthenticatingUsingAzureCLI:          getEnvBoolOrDefault(data.UseCLI, "ARM_USE_CLI", true),
		EnableAuthenticatingUsingManagedIdentity:   getEnvBoolOrDefault(data.UseMSI, "ARM_USE_MSI", false),
	}

	p.clientBuilder.SubscriptionID = getEnvStringIfValueAbsent(data.SubscriptionId, "ARM_SUBSCRIPTION_ID")

	partnerId := getEnvStringIfValueAbsent(data.PartnerId, "ARM_PARTNER_ID")
	if _, errs := provider.ValidatePartnerID(partnerId, "ARM_PARTNER_ID"); len(errs) > 0 {
		diags.Append(diag.NewErrorDiagnostic("validating ARM_PARTNER_ID", errs[0].Error()))
		return
	}
	p.clientBuilder.PartnerID = partnerId
	p.clientBuilder.DisableCorrelationRequestID = getEnvBoolOrDefault(data.DisableCorrelationRequestId, "ARM_DISABLE_CORRELATION_REQUEST_ID", false)
	p.clientBuilder.DisableTerraformPartnerID = getEnvBoolOrDefault(data.DisableTerraformPartnerId, "ARM_DISABLE_TERRAFORM_PARTNER_ID", false)
	p.clientBuilder.StorageUseAzureAD = getEnvBoolOrDefault(data.StorageUseAzureAD, "ARM_STORAGE_USE_AZUREAD", false)

	f := providerfeatures.UserFeatures{}

	// features is required, but we'll play safe here
	if !data.Features.IsNull() && !data.Features.IsUnknown() {
		var featuresList []Features
		d := data.Features.ElementsAs(ctx, &featuresList, true)
		diags.Append(d...)
		if diags.HasError() {
			return
		}

		features := featuresList[0]

		if !features.APIManagement.IsNull() && !features.APIManagement.IsUnknown() {
			var feature []APIManagement
			d := features.APIManagement.ElementsAs(ctx, &feature, true)
			diags.Append(d...)
			if diags.HasError() {
				return
			}

			f.ApiManagement.PurgeSoftDeleteOnDestroy = true
			if !feature[0].PurgeSoftDeleteOnDestroy.IsNull() && !feature[0].PurgeSoftDeleteOnDestroy.IsUnknown() {
				f.ApiManagement.PurgeSoftDeleteOnDestroy = feature[0].PurgeSoftDeleteOnDestroy.ValueBool()
			}

			f.ApiManagement.RecoverSoftDeleted = true
			if !feature[0].RecoverSoftDeleted.IsNull() && !feature[0].RecoverSoftDeleted.IsUnknown() {
				f.ApiManagement.RecoverSoftDeleted = feature[0].RecoverSoftDeleted.ValueBool()
			}
		} else {
			f.ApiManagement.PurgeSoftDeleteOnDestroy = true
			f.ApiManagement.RecoverSoftDeleted = true
		}

		if !features.AppConfiguration.IsNull() && !features.AppConfiguration.IsUnknown() {
			var feature []AppConfiguration
			d := features.AppConfiguration.ElementsAs(ctx, &feature, true)
			diags.Append(d...)
			if diags.HasError() {
				return
			}

			f.AppConfiguration.PurgeSoftDeleteOnDestroy = true
			if !feature[0].PurgeSoftDeleteOnDestroy.IsNull() && !feature[0].PurgeSoftDeleteOnDestroy.IsUnknown() {
				f.AppConfiguration.PurgeSoftDeleteOnDestroy = feature[0].PurgeSoftDeleteOnDestroy.ValueBool()
			}

			f.AppConfiguration.RecoverSoftDeleted = true
			if !feature[0].RecoverSoftDeleted.IsNull() && !feature[0].RecoverSoftDeleted.IsUnknown() {
				f.AppConfiguration.RecoverSoftDeleted = feature[0].RecoverSoftDeleted.ValueBool()
			}
		} else {
			f.AppConfiguration.PurgeSoftDeleteOnDestroy = true
			f.AppConfiguration.RecoverSoftDeleted = true
		}

		if !features.ApplicationInsights.IsNull() && !features.ApplicationInsights.IsUnknown() {
			var feature []ApplicationInsights
			d := features.ApplicationInsights.ElementsAs(ctx, &feature, true)
			diags.Append(d...)
			if diags.HasError() {
				return
			}

			f.ApplicationInsights.DisableGeneratedRule = false
			if !feature[0].DisableGeneratedRule.IsNull() && !feature[0].DisableGeneratedRule.IsUnknown() {
				f.ApplicationInsights.DisableGeneratedRule = feature[0].DisableGeneratedRule.ValueBool()
			}
		}

		if !features.CognitiveAccount.IsNull() && !features.CognitiveAccount.IsUnknown() {
			var feature []CognitiveAccount
			d := features.CognitiveAccount.ElementsAs(ctx, &feature, true)
			diags.Append(d...)
			if diags.HasError() {
				return
			}

			f.CognitiveAccount.PurgeSoftDeleteOnDestroy = true
			if !feature[0].PurgeSoftDeleteOnDestroy.IsNull() && !feature[0].PurgeSoftDeleteOnDestroy.IsUnknown() {
				f.CognitiveAccount.PurgeSoftDeleteOnDestroy = feature[0].PurgeSoftDeleteOnDestroy.ValueBool()
			}
		} else {
			f.CognitiveAccount.PurgeSoftDeleteOnDestroy = true
		}

		if !features.KeyVault.IsNull() && !features.KeyVault.IsUnknown() {
			var feature []KeyVault
			d := features.KeyVault.ElementsAs(ctx, &feature, true)
			diags.Append(d...)
			if diags.HasError() {
				return
			}

			f.KeyVault.PurgeSoftDeleteOnDestroy = true
			if !feature[0].PurgeSoftDeleteOnDestroy.IsNull() && !feature[0].PurgeSoftDeleteOnDestroy.IsUnknown() {
				f.KeyVault.PurgeSoftDeleteOnDestroy = feature[0].PurgeSoftDeleteOnDestroy.ValueBool()
			}

			f.KeyVault.PurgeSoftDeletedCertsOnDestroy = true
			if !feature[0].PurgeSoftDeletedCertificatesOnDestroy.IsNull() && !feature[0].PurgeSoftDeletedCertificatesOnDestroy.IsUnknown() {
				f.KeyVault.PurgeSoftDeletedCertsOnDestroy = feature[0].PurgeSoftDeletedCertificatesOnDestroy.ValueBool()
			}

			f.KeyVault.PurgeSoftDeletedKeysOnDestroy = true
			if !feature[0].PurgeSoftDeletedKeysOnDestroy.IsNull() && !feature[0].PurgeSoftDeletedKeysOnDestroy.IsUnknown() {
				f.KeyVault.PurgeSoftDeletedKeysOnDestroy = feature[0].PurgeSoftDeletedKeysOnDestroy.ValueBool()
			}

			f.KeyVault.PurgeSoftDeletedSecretsOnDestroy = true
			if !feature[0].PurgeSoftDeletedKeysOnDestroy.IsNull() && !feature[0].PurgeSoftDeletedKeysOnDestroy.IsUnknown() {
				f.KeyVault.PurgeSoftDeletedKeysOnDestroy = feature[0].PurgeSoftDeletedKeysOnDestroy.ValueBool()
			}

			f.KeyVault.PurgeSoftDeletedHSMsOnDestroy = true
			if !feature[0].PurgeSoftDeletedHardwareSecurityModulesOnDestroy.IsNull() && !feature[0].PurgeSoftDeletedHardwareSecurityModulesOnDestroy.IsUnknown() {
				f.KeyVault.PurgeSoftDeletedHSMsOnDestroy = feature[0].PurgeSoftDeletedHardwareSecurityModulesOnDestroy.ValueBool()
			}

			f.KeyVault.PurgeSoftDeletedHSMKeysOnDestroy = true
			if !feature[0].PurgeSoftDeletedHardwareSecurityModulesKeysOnDestroy.IsNull() && !feature[0].PurgeSoftDeletedHardwareSecurityModulesKeysOnDestroy.IsUnknown() {
				f.KeyVault.PurgeSoftDeletedHSMKeysOnDestroy = feature[0].PurgeSoftDeletedHardwareSecurityModulesKeysOnDestroy.ValueBool()
			}

			f.KeyVault.RecoverSoftDeletedCerts = true
			if !feature[0].RecoverSoftDeletedCertificates.IsNull() && !feature[0].RecoverSoftDeletedCertificates.IsUnknown() {
				f.KeyVault.RecoverSoftDeletedCerts = feature[0].RecoverSoftDeletedCertificates.ValueBool()
			}

			f.KeyVault.RecoverSoftDeletedKeyVaults = true
			if !feature[0].RecoverSoftDeletedKeyVaults.IsNull() && !feature[0].RecoverSoftDeletedKeyVaults.IsUnknown() {
				f.KeyVault.RecoverSoftDeletedKeyVaults = feature[0].RecoverSoftDeletedKeyVaults.ValueBool()
			}

			f.KeyVault.RecoverSoftDeletedKeys = true
			if !feature[0].RecoverSoftDeletedKeys.IsNull() && !feature[0].RecoverSoftDeletedKeys.IsUnknown() {
				f.KeyVault.RecoverSoftDeletedKeys = feature[0].RecoverSoftDeletedKeys.ValueBool()
			}

			f.KeyVault.RecoverSoftDeletedSecrets = true
			if !feature[0].RecoverSoftDeletedSecrets.IsNull() && !feature[0].RecoverSoftDeletedSecrets.IsUnknown() {
				f.KeyVault.RecoverSoftDeletedSecrets = feature[0].RecoverSoftDeletedSecrets.ValueBool()
			}

			f.KeyVault.RecoverSoftDeletedHSMKeys = true
			if !feature[0].RecoverSoftDeletedHSMKeys.IsNull() && !feature[0].RecoverSoftDeletedHSMKeys.IsUnknown() {
				f.KeyVault.RecoverSoftDeletedHSMKeys = feature[0].RecoverSoftDeletedHSMKeys.ValueBool()
			}
		} else {
			f.KeyVault.PurgeSoftDeleteOnDestroy = true
			f.KeyVault.PurgeSoftDeletedCertsOnDestroy = true
			f.KeyVault.PurgeSoftDeletedKeysOnDestroy = true
			f.KeyVault.PurgeSoftDeletedSecretsOnDestroy = true
			f.KeyVault.PurgeSoftDeletedHSMsOnDestroy = true
			f.KeyVault.PurgeSoftDeletedHSMKeysOnDestroy = true
			f.KeyVault.RecoverSoftDeletedCerts = true
			f.KeyVault.RecoverSoftDeletedKeyVaults = true
			f.KeyVault.RecoverSoftDeletedKeys = true
			f.KeyVault.RecoverSoftDeletedSecrets = true
			f.KeyVault.RecoverSoftDeletedHSMKeys = true
		}

		if !features.LogAnalyticsWorkspace.IsNull() && !features.LogAnalyticsWorkspace.IsUnknown() {
			var feature []LogAnalyticsWorkspace
			d := features.LogAnalyticsWorkspace.ElementsAs(ctx, &feature, true)
			diags.Append(d...)
			if diags.HasError() {
				return
			}

			f.LogAnalyticsWorkspace.PermanentlyDeleteOnDestroy = false
			if !feature[0].PermanentlyDeleteOnDestroy.IsNull() && !feature[0].PermanentlyDeleteOnDestroy.IsUnknown() {
				f.LogAnalyticsWorkspace.PermanentlyDeleteOnDestroy = feature[0].PermanentlyDeleteOnDestroy.ValueBool()
			}
		} else {
			f.LogAnalyticsWorkspace.PermanentlyDeleteOnDestroy = false
		}

		if !features.TemplateDeployment.IsNull() && !features.TemplateDeployment.IsUnknown() {
			var feature []TemplateDeployment
			d := features.TemplateDeployment.ElementsAs(ctx, &feature, true)
			diags.Append(d...)
			if diags.HasError() {
				return
			}

			f.TemplateDeployment.DeleteNestedItemsDuringDeletion = false
			if !feature[0].DeleteNestedItemsDuringDeletion.IsNull() && !feature[0].DeleteNestedItemsDuringDeletion.IsUnknown() {
				f.TemplateDeployment.DeleteNestedItemsDuringDeletion = feature[0].DeleteNestedItemsDuringDeletion.ValueBool()
			}
		} else {
			f.TemplateDeployment.DeleteNestedItemsDuringDeletion = false
		}

		if !features.VirtualMachine.IsNull() && !features.VirtualMachine.IsUnknown() {
			var feature []VirtualMachine
			d := features.VirtualMachine.ElementsAs(ctx, &feature, true)
			diags.Append(d...)
			if diags.HasError() {
				return
			}

			f.VirtualMachine.DeleteOSDiskOnDeletion = false
			if !feature[0].DeleteOsDiskOnDeletion.IsNull() && !feature[0].DeleteOsDiskOnDeletion.IsUnknown() {
				f.VirtualMachine.DeleteOSDiskOnDeletion = feature[0].DeleteOsDiskOnDeletion.ValueBool()
			}

			f.VirtualMachine.SkipShutdownAndForceDelete = false
			if !feature[0].SkipShutdownAndForceDelete.IsNull() && !feature[0].SkipShutdownAndForceDelete.IsUnknown() {
				f.VirtualMachine.SkipShutdownAndForceDelete = feature[0].SkipShutdownAndForceDelete.ValueBool()
			}
		} else {
			f.VirtualMachine.DeleteOSDiskOnDeletion = false
			f.VirtualMachine.SkipShutdownAndForceDelete = false
		}

		if !features.VirtualMachineScaleSet.IsNull() && !features.VirtualMachineScaleSet.IsUnknown() {
			var feature []VirtualMachineScaleSet
			d := features.VirtualMachineScaleSet.ElementsAs(ctx, &feature, true)
			diags.Append(d...)
			if diags.HasError() {
				return
			}

			f.VirtualMachineScaleSet.ForceDelete = false
			if !feature[0].ForceDelete.IsNull() && !feature[0].ForceDelete.IsUnknown() {
				f.VirtualMachineScaleSet.ForceDelete = feature[0].ForceDelete.ValueBool()
			}

			f.VirtualMachineScaleSet.ReimageOnManualUpgrade = true
			if !feature[0].ReimageOnManualUpgrade.IsNull() && !feature[0].ReimageOnManualUpgrade.IsUnknown() {
				f.VirtualMachineScaleSet.ReimageOnManualUpgrade = feature[0].ReimageOnManualUpgrade.ValueBool()
			}

			f.VirtualMachineScaleSet.RollInstancesWhenRequired = true
			if !feature[0].RollInstancesWhenRequired.IsNull() && !feature[0].RollInstancesWhenRequired.IsUnknown() {
				f.VirtualMachineScaleSet.RollInstancesWhenRequired = feature[0].RollInstancesWhenRequired.ValueBool()
			}

			f.VirtualMachineScaleSet.ScaleToZeroOnDelete = false
			if !feature[0].ScaleToZeroBeforeDeletion.IsNull() && !feature[0].ScaleToZeroBeforeDeletion.IsUnknown() {
				f.VirtualMachineScaleSet.ScaleToZeroOnDelete = feature[0].ScaleToZeroBeforeDeletion.ValueBool()
			}
		} else {
			f.VirtualMachineScaleSet.ForceDelete = false
			f.VirtualMachineScaleSet.ReimageOnManualUpgrade = true
			f.VirtualMachineScaleSet.RollInstancesWhenRequired = true
			f.VirtualMachineScaleSet.ScaleToZeroOnDelete = false
		}

		if !features.ResourceGroup.IsNull() && !features.ResourceGroup.IsUnknown() {
			var feature []ResourceGroup
			d := features.ResourceGroup.ElementsAs(ctx, &feature, true)
			diags.Append(d...)
			if diags.HasError() {
				return
			}

			f.ResourceGroup.PreventDeletionIfContainsResources = os.Getenv("TF_ACC") == ""
			if !feature[0].PreventDeletionIfContainsResources.IsNull() && !feature[0].PreventDeletionIfContainsResources.IsUnknown() {
				f.ResourceGroup.PreventDeletionIfContainsResources = feature[0].PreventDeletionIfContainsResources.ValueBool()
			}
		} else {
			f.ResourceGroup.PreventDeletionIfContainsResources = os.Getenv("TF_ACC") == ""
		}

		if !features.ManagedDisk.IsNull() && !features.ManagedDisk.IsUnknown() {
			var feature []ManagedDisk
			d := features.ManagedDisk.ElementsAs(ctx, &feature, true)
			diags.Append(d...)
			if diags.HasError() {
				return
			}

			f.ManagedDisk.ExpandWithoutDowntime = true
			if !feature[0].ExpandWithoutDowntime.IsNull() && !feature[0].ExpandWithoutDowntime.IsUnknown() {
				f.ManagedDisk.ExpandWithoutDowntime = feature[0].ExpandWithoutDowntime.ValueBool()
			}
		} else {
			f.ManagedDisk.ExpandWithoutDowntime = true
		}

		if !features.Storage.IsNull() && !features.Storage.IsUnknown() {
			var feature []Storage
			d := features.Storage.ElementsAs(ctx, &feature, true)
			diags.Append(d...)
			if diags.HasError() {
				return
			}
			f.Storage.DataPlaneAvailable = true
			if !feature[0].DataPlaneAvailable.IsNull() && !feature[0].DataPlaneAvailable.IsUnknown() {
				f.Storage.DataPlaneAvailable = feature[0].DataPlaneAvailable.ValueBool()
			}
		}

		if !features.Subscription.IsNull() && !features.Subscription.IsUnknown() {
			var feature []Subscription
			d := features.Subscription.ElementsAs(ctx, &feature, true)
			diags.Append(d...)
			if diags.HasError() {
				return
			}

			f.Subscription.PreventCancellationOnDestroy = false
			if !feature[0].PreventCancellationOnDestroy.IsNull() && !feature[0].PreventCancellationOnDestroy.IsUnknown() {
				f.Subscription.PreventCancellationOnDestroy = feature[0].PreventCancellationOnDestroy.ValueBool()
			}
		} else {
			f.Subscription.PreventCancellationOnDestroy = false
		}

		if !features.PostgresqlFlexibleServer.IsNull() && !features.PostgresqlFlexibleServer.IsUnknown() {
			var feature []PostgresqlFlexibleServer
			d := features.PostgresqlFlexibleServer.ElementsAs(ctx, &feature, true)
			diags.Append(d...)
			if diags.HasError() {
				return
			}

			f.PostgresqlFlexibleServer.RestartServerOnConfigurationValueChange = true
			if !feature[0].RestartServerOnConfigurationValueChange.IsNull() && !feature[0].RestartServerOnConfigurationValueChange.IsUnknown() {
				f.PostgresqlFlexibleServer.RestartServerOnConfigurationValueChange = feature[0].RestartServerOnConfigurationValueChange.ValueBool()
			}
		} else {
			f.PostgresqlFlexibleServer.RestartServerOnConfigurationValueChange = true
		}

		if !features.RecoveryService.IsNull() && !features.RecoveryService.IsUnknown() {
			var feature []RecoveryService
			d := features.RecoveryService.ElementsAs(ctx, &feature, true)
			diags.Append(d...)
			if diags.HasError() {
				return
			}

			f.RecoveryService.VMBackupStopProtectionAndRetainDataOnDestroy = false
			if !feature[0].VMBackupStopProtectionAndRetainDataOnDestroy.IsNull() && !feature[0].VMBackupStopProtectionAndRetainDataOnDestroy.IsUnknown() {
				f.RecoveryService.VMBackupStopProtectionAndRetainDataOnDestroy = feature[0].VMBackupStopProtectionAndRetainDataOnDestroy.ValueBool()
			}

			f.RecoveryService.VMBackupSuspendProtectionAndRetainDataOnDestroy = false
			if !feature[0].VMBackupSuspendProtectionAndRetainDataOnDestroy.IsNull() && !feature[0].VMBackupSuspendProtectionAndRetainDataOnDestroy.IsUnknown() {
				f.RecoveryService.VMBackupSuspendProtectionAndRetainDataOnDestroy = feature[0].VMBackupSuspendProtectionAndRetainDataOnDestroy.ValueBool()
			}

			f.RecoveryService.PurgeProtectedItemsFromVaultOnDestroy = false
			if !feature[0].PurgeProtectedItemsFromVaultOnDestroy.IsNull() && !feature[0].PurgeProtectedItemsFromVaultOnDestroy.IsUnknown() {
				f.RecoveryService.PurgeProtectedItemsFromVaultOnDestroy = feature[0].PurgeProtectedItemsFromVaultOnDestroy.ValueBool()
			}
		} else {
			f.RecoveryService.VMBackupStopProtectionAndRetainDataOnDestroy = false
			f.RecoveryService.PurgeProtectedItemsFromVaultOnDestroy = false
		}

		if !features.NetApp.IsNull() && !features.NetApp.IsUnknown() {
			var feature []NetApp
			d := features.NetApp.ElementsAs(ctx, &feature, true)
			diags.Append(d...)
			if diags.HasError() {
				return
			}

			f.NetApp.DeleteBackupsOnBackupVaultDestroy = false
			if !feature[0].DeleteBackupsOnBackupVaultDestroy.IsNull() && !feature[0].DeleteBackupsOnBackupVaultDestroy.IsUnknown() {
				f.NetApp.DeleteBackupsOnBackupVaultDestroy = feature[0].DeleteBackupsOnBackupVaultDestroy.ValueBool()
			}

			f.NetApp.PreventVolumeDestruction = true
			if !feature[0].PreventVolumeDestruction.IsNull() && !feature[0].PreventVolumeDestruction.IsUnknown() {
				f.NetApp.PreventVolumeDestruction = feature[0].PreventVolumeDestruction.ValueBool()
			}
		} else {
			f.NetApp.DeleteBackupsOnBackupVaultDestroy = false
			f.NetApp.PreventVolumeDestruction = true
		}

		if !features.DatabricksWorkspace.IsNull() && !features.DatabricksWorkspace.IsUnknown() {
			var feature []DatabricksWorkspace
			d := features.DatabricksWorkspace.ElementsAs(ctx, &feature, true)
			diags.Append(d...)
			if diags.HasError() {
				return
			}

			f.DatabricksWorkspace.ForceDelete = false
			if !feature[0].ForceDelete.IsNull() && !feature[0].ForceDelete.IsUnknown() {
				f.DatabricksWorkspace.ForceDelete = feature[0].ForceDelete.ValueBool()
			}
		} else {
			f.DatabricksWorkspace.ForceDelete = false
		}
	}

	p.clientBuilder.Features = f
	p.clientBuilder.AuthConfig = authConfig
	p.clientBuilder.CustomCorrelationRequestID = os.Getenv("ARM_CORRELATION_REQUEST_ID")
	p.clientBuilder.TerraformVersion = tfVersion

	client, err := clients.Build(ctx, p.clientBuilder)
	if err != nil {
		diags.Append(diag.NewErrorDiagnostic("building client", err.Error()))
		return
	}

	if diags.HasError() {
		return
	}

	client.StopContext = ctx

	resourceProviderRegistrationSet := getEnvStringOrDefault(data.ResourceProviderRegistrations, "ARM_RESOURCE_PROVIDER_REGISTRATIONS", resourceproviders.ProviderRegistrationsCore)
	if !providerfeatures.FivePointOh() {
		resourceProviderRegistrationSet = getEnvStringOrDefault(data.ResourceProviderRegistrations, "ARM_RESOURCE_PROVIDER_REGISTRATIONS", resourceproviders.ProviderRegistrationsLegacy)
	}

	if !data.SkipProviderRegistration.ValueBool() {
		if resourceProviderRegistrationSet != resourceproviders.ProviderRegistrationsLegacy {
			diags.Append(diag.NewErrorDiagnostic("resource provider registration misconfiguration", "provider property `skip_provider_registration` cannot be set at the same time as `resource_provider_registrations`, please remove `skip_provider_registration` from your configuration or unset the `ARM_SKIP_PROVIDER_REGISTRATION` environment variable"))
		}

		resourceProviderRegistrationSet = resourceproviders.ProviderRegistrationsNone
	}

	requiredResourceProviders, err := resourceproviders.GetResourceProvidersSet(resourceProviderRegistrationSet)
	if err != nil {
		diags.Append(diag.NewErrorDiagnostic("building resource providers", err.Error()))
		return
	}

	additionalResourceProvidersToRegister := make([]string, 0)
	if !data.ResourceProvidersToRegister.IsNull() {
		data.ResourceProvidersToRegister.ElementsAs(ctx, &additionalResourceProvidersToRegister, false)
		if len(additionalResourceProvidersToRegister) > 0 {
			additionalProviders := make(resourceproviders.ResourceProviders)
			for _, rp := range additionalResourceProvidersToRegister {
				additionalProviders.Add(rp)
			}
		}
	}

	subId := commonids.NewSubscriptionID(client.Account.SubscriptionId)
	ctx2, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	if err = resourceproviders.EnsureRegistered(ctx2, client.Resource.ResourceProvidersClient, subId, requiredResourceProviders); err != nil {
		diags.AddError("registering resource providers", err.Error())
		return
	}

	p.Client = client
}
