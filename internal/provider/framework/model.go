// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package framework

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProviderModel struct {
	SubscriptionId                 types.String `tfsdk:"subscription_id"`
	ClientId                       types.String `tfsdk:"client_id"`
	ClientIdFilePath               types.String `tfsdk:"client_id_file_path"`
	TenantId                       types.String `tfsdk:"tenant_id"`
	AuxiliaryTenantIds             types.List   `tfsdk:"auxiliary_tenant_ids"`
	Environment                    types.String `tfsdk:"environment"`
	MetaDataHost                   types.String `tfsdk:"metadata_host"`
	ClientCertificate              types.String `tfsdk:"client_certificate"`
	ClientCertificatePath          types.String `tfsdk:"client_certificate_path"`
	ClientCertificatePassword      types.String `tfsdk:"client_certificate_password"`
	ClientSecret                   types.String `tfsdk:"client_secret"`
	ClientSecretFilePath           types.String `tfsdk:"client_secret_file_path"`
	ADOPipelineServiceConnectionID types.String `tfsdk:"ado_pipeline_service_connection_id"`
	OIDCRequestToken               types.String `tfsdk:"oidc_request_token"`
	OIDCRequestURL                 types.String `tfsdk:"oidc_request_url"`
	OIDCToken                      types.String `tfsdk:"oidc_token"`
	OIDCTokenFilePath              types.String `tfsdk:"oidc_token_file_path"`
	UseOIDC                        types.Bool   `tfsdk:"use_oidc"`
	UseMSI                         types.Bool   `tfsdk:"use_msi"`
	MSIEndpoint                    types.String `tfsdk:"msi_endpoint"`
	UseCLI                         types.Bool   `tfsdk:"use_cli"`
	UseAKSWorkloadIdentity         types.Bool   `tfsdk:"use_aks_workload_identity"`
	PartnerId                      types.String `tfsdk:"partner_id"`
	DisableCorrelationRequestId    types.Bool   `tfsdk:"disable_correlation_request_id"`
	DisableTerraformPartnerId      types.Bool   `tfsdk:"disable_terraform_partner_id"`
	StorageUseAzureAD              types.Bool   `tfsdk:"storage_use_azuread"`
	Features                       types.List   `tfsdk:"features"`
	SkipProviderRegistration       types.Bool   `tfsdk:"skip_provider_registration"` // TODO - Remove in 5.0
	ResourceProviderRegistrations  types.String `tfsdk:"resource_provider_registrations"`
	ResourceProvidersToRegister    types.List   `tfsdk:"resource_providers_to_register"`
}

type Features struct {
	APIManagement            types.List `tfsdk:"api_management"`
	AppConfiguration         types.List `tfsdk:"app_configuration"`
	ApplicationInsights      types.List `tfsdk:"application_insights"`
	CognitiveAccount         types.List `tfsdk:"cognitive_account"`
	KeyVault                 types.List `tfsdk:"key_vault"`
	LogAnalyticsWorkspace    types.List `tfsdk:"log_analytics_workspace"`
	TemplateDeployment       types.List `tfsdk:"template_deployment"`
	VirtualMachine           types.List `tfsdk:"virtual_machine"`
	VirtualMachineScaleSet   types.List `tfsdk:"virtual_machine_scale_set"`
	ResourceGroup            types.List `tfsdk:"resource_group"`
	ManagedDisk              types.List `tfsdk:"managed_disk"`
	Storage                  types.List `tfsdk:"storage"`
	Subscription             types.List `tfsdk:"subscription"`
	PostgresqlFlexibleServer types.List `tfsdk:"postgresql_flexible_server"`
	MachineLearning          types.List `tfsdk:"machine_learning"`
	RecoveryService          types.List `tfsdk:"recovery_service"`
	RecoveryServicesVaults   types.List `tfsdk:"recovery_services_vaults"`
	NetApp                   types.List `tfsdk:"netapp"`
	DatabricksWorkspace      types.List `tfsdk:"databricks_workspace"`
}

// FeaturesAttributes and the other block attribute vars are required for unit testing on the Load func
// New features blocks and attributes must be added here and to unit tests.
var FeaturesAttributes = map[string]attr.Type{
	"api_management":             types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(APIManagementAttributes)),
	"app_configuration":          types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(AppConfigurationAttributes)),
	"application_insights":       types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(ApplicationInsightsAttributes)),
	"cognitive_account":          types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(CognitiveAccountAttributes)),
	"key_vault":                  types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(KeyVaultAttributes)),
	"log_analytics_workspace":    types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(LogAnalyticsWorkspaceAttributes)),
	"template_deployment":        types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(TemplateDeploymentAttributes)),
	"virtual_machine":            types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(VirtualMachineAttributes)),
	"virtual_machine_scale_set":  types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(VirtualMachineScaleSetAttributes)),
	"resource_group":             types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(ResourceGroupAttributes)),
	"managed_disk":               types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(ManagedDiskAttributes)),
	"storage":                    types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(StorageAttributes)),
	"subscription":               types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(SubscriptionAttributes)),
	"postgresql_flexible_server": types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(PostgresqlFlexibleServerAttributes)),
	"machine_learning":           types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(MachineLearningAttributes)),
	"recovery_service":           types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(RecoveryServiceAttributes)),
	"recovery_services_vaults":   types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(RecoveryServiceVaultsAttributes)),
	"netapp":                     types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(NetAppAttributes)),
	"databricks_workspace":       types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(DatabricksWorkspaceAttributes)),
}

type APIManagement struct {
	PurgeSoftDeleteOnDestroy types.Bool `tfsdk:"purge_soft_delete_on_destroy"`
	RecoverSoftDeleted       types.Bool `tfsdk:"recover_soft_deleted"`
}

var APIManagementAttributes = map[string]attr.Type{
	"purge_soft_delete_on_destroy": types.BoolType,
	"recover_soft_deleted":         types.BoolType,
}

type AppConfiguration struct {
	PurgeSoftDeleteOnDestroy types.Bool `tfsdk:"purge_soft_delete_on_destroy"`
	RecoverSoftDeleted       types.Bool `tfsdk:"recover_soft_deleted"`
}

var AppConfigurationAttributes = map[string]attr.Type{
	"purge_soft_delete_on_destroy": types.BoolType,
	"recover_soft_deleted":         types.BoolType,
}

type ApplicationInsights struct {
	DisableGeneratedRule types.Bool `tfsdk:"disable_generated_rule"`
}

var ApplicationInsightsAttributes = map[string]attr.Type{
	"disable_generated_rule": types.BoolType,
}

type CognitiveAccount struct {
	PurgeSoftDeleteOnDestroy types.Bool `tfsdk:"purge_soft_delete_on_destroy"`
}

var CognitiveAccountAttributes = map[string]attr.Type{
	"purge_soft_delete_on_destroy": types.BoolType,
}

type KeyVault struct {
	PurgeSoftDeleteOnDestroy                             types.Bool `tfsdk:"purge_soft_delete_on_destroy"`
	PurgeSoftDeletedCertificatesOnDestroy                types.Bool `tfsdk:"purge_soft_deleted_certificates_on_destroy"`
	PurgeSoftDeletedKeysOnDestroy                        types.Bool `tfsdk:"purge_soft_deleted_keys_on_destroy"`
	PurgeSoftDeletedSecretsOnDestroy                     types.Bool `tfsdk:"purge_soft_deleted_secrets_on_destroy"`
	PurgeSoftDeletedHardwareSecurityModulesOnDestroy     types.Bool `tfsdk:"purge_soft_deleted_hardware_security_modules_on_destroy"`
	PurgeSoftDeletedHardwareSecurityModulesKeysOnDestroy types.Bool `tfsdk:"purge_soft_deleted_hardware_security_module_keys_on_destroy"`
	RecoverSoftDeletedCertificates                       types.Bool `tfsdk:"recover_soft_deleted_certificates"`
	RecoverSoftDeletedKeyVaults                          types.Bool `tfsdk:"recover_soft_deleted_key_vaults"`
	RecoverSoftDeletedKeys                               types.Bool `tfsdk:"recover_soft_deleted_keys"`
	RecoverSoftDeletedSecrets                            types.Bool `tfsdk:"recover_soft_deleted_secrets"`
	RecoverSoftDeletedHSMKeys                            types.Bool `tfsdk:"recover_soft_deleted_hardware_security_module_keys"`
}

var KeyVaultAttributes = map[string]attr.Type{
	"purge_soft_delete_on_destroy":                                types.BoolType,
	"purge_soft_deleted_certificates_on_destroy":                  types.BoolType,
	"purge_soft_deleted_keys_on_destroy":                          types.BoolType,
	"purge_soft_deleted_secrets_on_destroy":                       types.BoolType,
	"purge_soft_deleted_hardware_security_modules_on_destroy":     types.BoolType,
	"purge_soft_deleted_hardware_security_module_keys_on_destroy": types.BoolType,
	"recover_soft_deleted_certificates":                           types.BoolType,
	"recover_soft_deleted_key_vaults":                             types.BoolType,
	"recover_soft_deleted_keys":                                   types.BoolType,
	"recover_soft_deleted_secrets":                                types.BoolType,
	"recover_soft_deleted_hardware_security_module_keys":          types.BoolType,
}

type LogAnalyticsWorkspace struct {
	PermanentlyDeleteOnDestroy types.Bool `tfsdk:"permanently_delete_on_destroy"`
}

var LogAnalyticsWorkspaceAttributes = map[string]attr.Type{
	"permanently_delete_on_destroy": types.BoolType,
}

type TemplateDeployment struct {
	DeleteNestedItemsDuringDeletion types.Bool `tfsdk:"delete_nested_items_during_deletion"`
}

var TemplateDeploymentAttributes = map[string]attr.Type{
	"delete_nested_items_during_deletion": types.BoolType,
}

type VirtualMachine struct {
	DeleteOsDiskOnDeletion           types.Bool `tfsdk:"delete_os_disk_on_deletion"`
	GracefulShutdown                 types.Bool `tfsdk:"graceful_shutdown"` // TODO: Remove in 5.0 - Currently not possible to deprecate feature block struct items via feature flagging. Feature made redundant/ineffective by a breaking API change.
	SkipShutdownAndForceDelete       types.Bool `tfsdk:"skip_shutdown_and_force_delete"`
	DetachImplicitDataDiskOnDeletion types.Bool `tfsdk:"detach_implicit_data_disk_on_deletion"`
}

var VirtualMachineAttributes = map[string]attr.Type{
	"delete_os_disk_on_deletion":            types.BoolType,
	"detach_implicit_data_disk_on_deletion": types.BoolType,
	"graceful_shutdown":                     types.BoolType, // TODO: Remove in 5.0 - Currently not possible to deprecate feature block struct items via feature flagging. Feature made redundant/ineffective by a breaking API change.
	"skip_shutdown_and_force_delete":        types.BoolType,
}

type VirtualMachineScaleSet struct {
	ForceDelete               types.Bool `tfsdk:"force_delete"`
	ReimageOnManualUpgrade    types.Bool `tfsdk:"reimage_on_manual_upgrade"`
	RollInstancesWhenRequired types.Bool `tfsdk:"roll_instances_when_required"`
	ScaleToZeroBeforeDeletion types.Bool `tfsdk:"scale_to_zero_before_deletion"`
}

var VirtualMachineScaleSetAttributes = map[string]attr.Type{
	"force_delete":                  types.BoolType,
	"reimage_on_manual_upgrade":     types.BoolType,
	"roll_instances_when_required":  types.BoolType,
	"scale_to_zero_before_deletion": types.BoolType,
}

type ResourceGroup struct {
	PreventDeletionIfContainsResources types.Bool `tfsdk:"prevent_deletion_if_contains_resources"`
}

var ResourceGroupAttributes = map[string]attr.Type{
	"prevent_deletion_if_contains_resources": types.BoolType,
}

type ManagedDisk struct {
	ExpandWithoutDowntime types.Bool `tfsdk:"expand_without_downtime"`
}

var ManagedDiskAttributes = map[string]attr.Type{
	"expand_without_downtime": types.BoolType,
}

type Storage struct {
	DataPlaneAvailable types.Bool `tfsdk:"data_plane_available"`
}

var StorageAttributes = map[string]attr.Type{
	"data_plane_available": types.BoolType,
}

type Subscription struct {
	PreventCancellationOnDestroy types.Bool `tfsdk:"prevent_cancellation_on_destroy"`
}

var SubscriptionAttributes = map[string]attr.Type{
	"prevent_cancellation_on_destroy": types.BoolType,
}

type PostgresqlFlexibleServer struct {
	RestartServerOnConfigurationValueChange types.Bool `tfsdk:"restart_server_on_configuration_value_change"`
}

var PostgresqlFlexibleServerAttributes = map[string]attr.Type{
	"restart_server_on_configuration_value_change": types.BoolType,
}

type MachineLearning struct {
	PurgeSoftDeletedWorkspaceOnDestroy types.Bool `tfsdk:"purge_soft_deleted_workspace_on_destroy"`
}

var MachineLearningAttributes = map[string]attr.Type{
	"purge_soft_deleted_workspace_on_destroy": types.BoolType,
}

type RecoveryService struct {
	VMBackupStopProtectionAndRetainDataOnDestroy    types.Bool `tfsdk:"vm_backup_stop_protection_and_retain_data_on_destroy"`
	VMBackupSuspendProtectionAndRetainDataOnDestroy types.Bool `tfsdk:"vm_backup_suspend_protection_and_retain_data_on_destroy"`
	PurgeProtectedItemsFromVaultOnDestroy           types.Bool `tfsdk:"purge_protected_items_from_vault_on_destroy"`
}

var RecoveryServiceAttributes = map[string]attr.Type{
	"vm_backup_stop_protection_and_retain_data_on_destroy":    types.BoolType,
	"vm_backup_suspend_protection_and_retain_data_on_destroy": types.BoolType,
	"purge_protected_items_from_vault_on_destroy":             types.BoolType,
}

type RecoveryServiceVaults struct {
	RecoverSoftDeletedBackupProtectedVm types.Bool `tfsdk:"recover_soft_deleted_backup_protected_vm"`
}

var RecoveryServiceVaultsAttributes = map[string]attr.Type{
	"recover_soft_deleted_backup_protected_vm": types.BoolType,
}

type NetApp struct {
	DeleteBackupsOnBackupVaultDestroy types.Bool `tfsdk:"delete_backups_on_backup_vault_destroy"`
	PreventVolumeDestruction          types.Bool `tfsdk:"prevent_volume_destruction"`
}

var NetAppAttributes = map[string]attr.Type{
	"delete_backups_on_backup_vault_destroy": types.BoolType,
	"prevent_volume_destruction":             types.BoolType,
}

type DatabricksWorkspace struct {
	ForceDelete types.Bool `tfsdk:"force_delete"`
}

var DatabricksWorkspaceAttributes = map[string]attr.Type{
	"force_delete": types.BoolType,
}
