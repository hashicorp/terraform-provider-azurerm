package deploymentscripts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzurePowerShellScriptProperties struct {
	Arguments              *string                      `json:"arguments,omitempty"`
	AzPowerShellVersion    string                       `json:"azPowerShellVersion"`
	CleanupPreference      *CleanupOptions              `json:"cleanupPreference,omitempty"`
	ContainerSettings      *ContainerConfiguration      `json:"containerSettings,omitempty"`
	EnvironmentVariables   *[]EnvironmentVariable       `json:"environmentVariables,omitempty"`
	ForceUpdateTag         *string                      `json:"forceUpdateTag,omitempty"`
	Outputs                *map[string]interface{}      `json:"outputs,omitempty"`
	PrimaryScriptUri       *string                      `json:"primaryScriptUri,omitempty"`
	ProvisioningState      *ScriptProvisioningState     `json:"provisioningState,omitempty"`
	RetentionInterval      string                       `json:"retentionInterval"`
	ScriptContent          *string                      `json:"scriptContent,omitempty"`
	Status                 *ScriptStatus                `json:"status,omitempty"`
	StorageAccountSettings *StorageAccountConfiguration `json:"storageAccountSettings,omitempty"`
	SupportingScriptUris   *[]string                    `json:"supportingScriptUris,omitempty"`
	Timeout                *string                      `json:"timeout,omitempty"`
}
