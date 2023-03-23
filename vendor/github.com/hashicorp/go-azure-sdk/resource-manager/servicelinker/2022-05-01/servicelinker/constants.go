package servicelinker

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthType string

const (
	AuthTypeSecret                      AuthType = "secret"
	AuthTypeServicePrincipalCertificate AuthType = "servicePrincipalCertificate"
	AuthTypeServicePrincipalSecret      AuthType = "servicePrincipalSecret"
	AuthTypeSystemAssignedIdentity      AuthType = "systemAssignedIdentity"
	AuthTypeUserAssignedIdentity        AuthType = "userAssignedIdentity"
)

func PossibleValuesForAuthType() []string {
	return []string{
		string(AuthTypeSecret),
		string(AuthTypeServicePrincipalCertificate),
		string(AuthTypeServicePrincipalSecret),
		string(AuthTypeSystemAssignedIdentity),
		string(AuthTypeUserAssignedIdentity),
	}
}

func parseAuthType(input string) (*AuthType, error) {
	vals := map[string]AuthType{
		"secret":                      AuthTypeSecret,
		"serviceprincipalcertificate": AuthTypeServicePrincipalCertificate,
		"serviceprincipalsecret":      AuthTypeServicePrincipalSecret,
		"systemassignedidentity":      AuthTypeSystemAssignedIdentity,
		"userassignedidentity":        AuthTypeUserAssignedIdentity,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuthType(input)
	return &out, nil
}

type AzureResourceType string

const (
	AzureResourceTypeKeyVault AzureResourceType = "KeyVault"
)

func PossibleValuesForAzureResourceType() []string {
	return []string{
		string(AzureResourceTypeKeyVault),
	}
}

func parseAzureResourceType(input string) (*AzureResourceType, error) {
	vals := map[string]AzureResourceType{
		"keyvault": AzureResourceTypeKeyVault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureResourceType(input)
	return &out, nil
}

type ClientType string

const (
	ClientTypeDjango                  ClientType = "django"
	ClientTypeDotnet                  ClientType = "dotnet"
	ClientTypeGo                      ClientType = "go"
	ClientTypeJava                    ClientType = "java"
	ClientTypeKafkaNegativespringBoot ClientType = "kafka-springBoot"
	ClientTypeNodejs                  ClientType = "nodejs"
	ClientTypeNone                    ClientType = "none"
	ClientTypePhp                     ClientType = "php"
	ClientTypePython                  ClientType = "python"
	ClientTypeRuby                    ClientType = "ruby"
	ClientTypeSpringBoot              ClientType = "springBoot"
)

func PossibleValuesForClientType() []string {
	return []string{
		string(ClientTypeDjango),
		string(ClientTypeDotnet),
		string(ClientTypeGo),
		string(ClientTypeJava),
		string(ClientTypeKafkaNegativespringBoot),
		string(ClientTypeNodejs),
		string(ClientTypeNone),
		string(ClientTypePhp),
		string(ClientTypePython),
		string(ClientTypeRuby),
		string(ClientTypeSpringBoot),
	}
}

func parseClientType(input string) (*ClientType, error) {
	vals := map[string]ClientType{
		"django":           ClientTypeDjango,
		"dotnet":           ClientTypeDotnet,
		"go":               ClientTypeGo,
		"java":             ClientTypeJava,
		"kafka-springboot": ClientTypeKafkaNegativespringBoot,
		"nodejs":           ClientTypeNodejs,
		"none":             ClientTypeNone,
		"php":              ClientTypePhp,
		"python":           ClientTypePython,
		"ruby":             ClientTypeRuby,
		"springboot":       ClientTypeSpringBoot,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClientType(input)
	return &out, nil
}

type SecretType string

const (
	SecretTypeKeyVaultSecretReference SecretType = "keyVaultSecretReference"
	SecretTypeKeyVaultSecretUri       SecretType = "keyVaultSecretUri"
	SecretTypeRawValue                SecretType = "rawValue"
)

func PossibleValuesForSecretType() []string {
	return []string{
		string(SecretTypeKeyVaultSecretReference),
		string(SecretTypeKeyVaultSecretUri),
		string(SecretTypeRawValue),
	}
}

func parseSecretType(input string) (*SecretType, error) {
	vals := map[string]SecretType{
		"keyvaultsecretreference": SecretTypeKeyVaultSecretReference,
		"keyvaultsecreturi":       SecretTypeKeyVaultSecretUri,
		"rawvalue":                SecretTypeRawValue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecretType(input)
	return &out, nil
}

type TargetServiceType string

const (
	TargetServiceTypeAzureResource            TargetServiceType = "AzureResource"
	TargetServiceTypeConfluentBootstrapServer TargetServiceType = "ConfluentBootstrapServer"
	TargetServiceTypeConfluentSchemaRegistry  TargetServiceType = "ConfluentSchemaRegistry"
)

func PossibleValuesForTargetServiceType() []string {
	return []string{
		string(TargetServiceTypeAzureResource),
		string(TargetServiceTypeConfluentBootstrapServer),
		string(TargetServiceTypeConfluentSchemaRegistry),
	}
}

func parseTargetServiceType(input string) (*TargetServiceType, error) {
	vals := map[string]TargetServiceType{
		"azureresource":            TargetServiceTypeAzureResource,
		"confluentbootstrapserver": TargetServiceTypeConfluentBootstrapServer,
		"confluentschemaregistry":  TargetServiceTypeConfluentSchemaRegistry,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TargetServiceType(input)
	return &out, nil
}

type VNetSolutionType string

const (
	VNetSolutionTypePrivateLink     VNetSolutionType = "privateLink"
	VNetSolutionTypeServiceEndpoint VNetSolutionType = "serviceEndpoint"
)

func PossibleValuesForVNetSolutionType() []string {
	return []string{
		string(VNetSolutionTypePrivateLink),
		string(VNetSolutionTypeServiceEndpoint),
	}
}

func parseVNetSolutionType(input string) (*VNetSolutionType, error) {
	vals := map[string]VNetSolutionType{
		"privatelink":     VNetSolutionTypePrivateLink,
		"serviceendpoint": VNetSolutionTypeServiceEndpoint,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VNetSolutionType(input)
	return &out, nil
}
