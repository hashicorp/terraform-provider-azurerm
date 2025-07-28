// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helper

type AuthorizationType string

const (
	AuthorizationTypeKey                          AuthorizationType = "Key"
	AuthorizationTypeLinkedIntegrationRuntimeType AuthorizationType = "LinkedIntegrationRuntimeType"
	AuthorizationTypeRBAC                         AuthorizationType = "RBAC"
)

type CustomSetupType string

const (
	CustomSetupTypeAzPowerShellSetup        CustomSetupType = "AzPowerShellSetup"
	CustomSetupTypeCmdkeySetup              CustomSetupType = "CmdkeySetup"
	CustomSetupTypeComponentSetup           CustomSetupType = "ComponentSetup"
	CustomSetupTypeEnvironmentVariableSetup CustomSetupType = "EnvironmentVariableSetup"
)

type DataFlowType string

const (
	DataFlowTypeFlowlet           = "Flowlet"
	DataFlowTypeMappingDataFlow   = "MappingDataFlow"
	DataFlowTypeWranglingDataFlow = "WranglingDataFlow"
)

type SecretType string

const (
	SecretTypeAzureKeyVaultSecret SecretType = "AzureKeyVaultSecret"
	SecretTypeSecureString        SecretType = "SecureString"
)
