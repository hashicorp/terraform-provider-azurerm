package servicelinker

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessKeyPermissions string

const (
	AccessKeyPermissionsListen AccessKeyPermissions = "Listen"
	AccessKeyPermissionsManage AccessKeyPermissions = "Manage"
	AccessKeyPermissionsRead   AccessKeyPermissions = "Read"
	AccessKeyPermissionsSend   AccessKeyPermissions = "Send"
	AccessKeyPermissionsWrite  AccessKeyPermissions = "Write"
)

func PossibleValuesForAccessKeyPermissions() []string {
	return []string{
		string(AccessKeyPermissionsListen),
		string(AccessKeyPermissionsManage),
		string(AccessKeyPermissionsRead),
		string(AccessKeyPermissionsSend),
		string(AccessKeyPermissionsWrite),
	}
}

func (s *AccessKeyPermissions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccessKeyPermissions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccessKeyPermissions(input string) (*AccessKeyPermissions, error) {
	vals := map[string]AccessKeyPermissions{
		"listen": AccessKeyPermissionsListen,
		"manage": AccessKeyPermissionsManage,
		"read":   AccessKeyPermissionsRead,
		"send":   AccessKeyPermissionsSend,
		"write":  AccessKeyPermissionsWrite,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessKeyPermissions(input)
	return &out, nil
}

type ActionType string

const (
	ActionTypeEnable ActionType = "enable"
	ActionTypeOptOut ActionType = "optOut"
)

func PossibleValuesForActionType() []string {
	return []string{
		string(ActionTypeEnable),
		string(ActionTypeOptOut),
	}
}

func (s *ActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseActionType(input string) (*ActionType, error) {
	vals := map[string]ActionType{
		"enable": ActionTypeEnable,
		"optout": ActionTypeOptOut,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ActionType(input)
	return &out, nil
}

type AllowType string

const (
	AllowTypeFalse AllowType = "false"
	AllowTypeTrue  AllowType = "true"
)

func PossibleValuesForAllowType() []string {
	return []string{
		string(AllowTypeFalse),
		string(AllowTypeTrue),
	}
}

func (s *AllowType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAllowType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAllowType(input string) (*AllowType, error) {
	vals := map[string]AllowType{
		"false": AllowTypeFalse,
		"true":  AllowTypeTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AllowType(input)
	return &out, nil
}

type AuthMode string

const (
	AuthModeOptInAllAuth  AuthMode = "optInAllAuth"
	AuthModeOptOutAllAuth AuthMode = "optOutAllAuth"
)

func PossibleValuesForAuthMode() []string {
	return []string{
		string(AuthModeOptInAllAuth),
		string(AuthModeOptOutAllAuth),
	}
}

func (s *AuthMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAuthMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAuthMode(input string) (*AuthMode, error) {
	vals := map[string]AuthMode{
		"optinallauth":  AuthModeOptInAllAuth,
		"optoutallauth": AuthModeOptOutAllAuth,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuthMode(input)
	return &out, nil
}

type AuthType string

const (
	AuthTypeAccessKey                   AuthType = "accessKey"
	AuthTypeEasyAuthMicrosoftEntraID    AuthType = "easyAuthMicrosoftEntraID"
	AuthTypeSecret                      AuthType = "secret"
	AuthTypeServicePrincipalCertificate AuthType = "servicePrincipalCertificate"
	AuthTypeServicePrincipalSecret      AuthType = "servicePrincipalSecret"
	AuthTypeSystemAssignedIdentity      AuthType = "systemAssignedIdentity"
	AuthTypeUserAccount                 AuthType = "userAccount"
	AuthTypeUserAssignedIdentity        AuthType = "userAssignedIdentity"
)

func PossibleValuesForAuthType() []string {
	return []string{
		string(AuthTypeAccessKey),
		string(AuthTypeEasyAuthMicrosoftEntraID),
		string(AuthTypeSecret),
		string(AuthTypeServicePrincipalCertificate),
		string(AuthTypeServicePrincipalSecret),
		string(AuthTypeSystemAssignedIdentity),
		string(AuthTypeUserAccount),
		string(AuthTypeUserAssignedIdentity),
	}
}

func (s *AuthType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAuthType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAuthType(input string) (*AuthType, error) {
	vals := map[string]AuthType{
		"accesskey":                   AuthTypeAccessKey,
		"easyauthmicrosoftentraid":    AuthTypeEasyAuthMicrosoftEntraID,
		"secret":                      AuthTypeSecret,
		"serviceprincipalcertificate": AuthTypeServicePrincipalCertificate,
		"serviceprincipalsecret":      AuthTypeServicePrincipalSecret,
		"systemassignedidentity":      AuthTypeSystemAssignedIdentity,
		"useraccount":                 AuthTypeUserAccount,
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

func (s *AzureResourceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureResourceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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
	ClientTypeDapr                    ClientType = "dapr"
	ClientTypeDjango                  ClientType = "django"
	ClientTypeDotnet                  ClientType = "dotnet"
	ClientTypeGo                      ClientType = "go"
	ClientTypeJava                    ClientType = "java"
	ClientTypeJmsNegativespringBoot   ClientType = "jms-springBoot"
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
		string(ClientTypeDapr),
		string(ClientTypeDjango),
		string(ClientTypeDotnet),
		string(ClientTypeGo),
		string(ClientTypeJava),
		string(ClientTypeJmsNegativespringBoot),
		string(ClientTypeKafkaNegativespringBoot),
		string(ClientTypeNodejs),
		string(ClientTypeNone),
		string(ClientTypePhp),
		string(ClientTypePython),
		string(ClientTypeRuby),
		string(ClientTypeSpringBoot),
	}
}

func (s *ClientType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClientType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClientType(input string) (*ClientType, error) {
	vals := map[string]ClientType{
		"dapr":             ClientTypeDapr,
		"django":           ClientTypeDjango,
		"dotnet":           ClientTypeDotnet,
		"go":               ClientTypeGo,
		"java":             ClientTypeJava,
		"jms-springboot":   ClientTypeJmsNegativespringBoot,
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

type DaprBindingComponentDirection string

const (
	DaprBindingComponentDirectionInput  DaprBindingComponentDirection = "input"
	DaprBindingComponentDirectionOutput DaprBindingComponentDirection = "output"
)

func PossibleValuesForDaprBindingComponentDirection() []string {
	return []string{
		string(DaprBindingComponentDirectionInput),
		string(DaprBindingComponentDirectionOutput),
	}
}

func (s *DaprBindingComponentDirection) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDaprBindingComponentDirection(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDaprBindingComponentDirection(input string) (*DaprBindingComponentDirection, error) {
	vals := map[string]DaprBindingComponentDirection{
		"input":  DaprBindingComponentDirectionInput,
		"output": DaprBindingComponentDirectionOutput,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DaprBindingComponentDirection(input)
	return &out, nil
}

type DaprMetadataRequired string

const (
	DaprMetadataRequiredFalse DaprMetadataRequired = "false"
	DaprMetadataRequiredTrue  DaprMetadataRequired = "true"
)

func PossibleValuesForDaprMetadataRequired() []string {
	return []string{
		string(DaprMetadataRequiredFalse),
		string(DaprMetadataRequiredTrue),
	}
}

func (s *DaprMetadataRequired) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDaprMetadataRequired(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDaprMetadataRequired(input string) (*DaprMetadataRequired, error) {
	vals := map[string]DaprMetadataRequired{
		"false": DaprMetadataRequiredFalse,
		"true":  DaprMetadataRequiredTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DaprMetadataRequired(input)
	return &out, nil
}

type DeleteOrUpdateBehavior string

const (
	DeleteOrUpdateBehaviorDefault       DeleteOrUpdateBehavior = "Default"
	DeleteOrUpdateBehaviorForcedCleanup DeleteOrUpdateBehavior = "ForcedCleanup"
)

func PossibleValuesForDeleteOrUpdateBehavior() []string {
	return []string{
		string(DeleteOrUpdateBehaviorDefault),
		string(DeleteOrUpdateBehaviorForcedCleanup),
	}
}

func (s *DeleteOrUpdateBehavior) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeleteOrUpdateBehavior(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeleteOrUpdateBehavior(input string) (*DeleteOrUpdateBehavior, error) {
	vals := map[string]DeleteOrUpdateBehavior{
		"default":       DeleteOrUpdateBehaviorDefault,
		"forcedcleanup": DeleteOrUpdateBehaviorForcedCleanup,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeleteOrUpdateBehavior(input)
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

func (s *SecretType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecretType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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
	TargetServiceTypeSelfHostedServer         TargetServiceType = "SelfHostedServer"
)

func PossibleValuesForTargetServiceType() []string {
	return []string{
		string(TargetServiceTypeAzureResource),
		string(TargetServiceTypeConfluentBootstrapServer),
		string(TargetServiceTypeConfluentSchemaRegistry),
		string(TargetServiceTypeSelfHostedServer),
	}
}

func (s *TargetServiceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTargetServiceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTargetServiceType(input string) (*TargetServiceType, error) {
	vals := map[string]TargetServiceType{
		"azureresource":            TargetServiceTypeAzureResource,
		"confluentbootstrapserver": TargetServiceTypeConfluentBootstrapServer,
		"confluentschemaregistry":  TargetServiceTypeConfluentSchemaRegistry,
		"selfhostedserver":         TargetServiceTypeSelfHostedServer,
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

func (s *VNetSolutionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVNetSolutionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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
