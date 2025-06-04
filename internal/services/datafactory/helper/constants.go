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

type SecretType string

const (
	SecretTypeAzureKeyVaultSecret SecretType = "AzureKeyVaultSecret"
	SecretTypeSecureString        SecretType = "SecureString"
)
