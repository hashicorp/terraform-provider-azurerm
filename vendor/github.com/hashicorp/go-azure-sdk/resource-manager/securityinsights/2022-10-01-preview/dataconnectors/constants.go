package dataconnectors

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailabilityStatus int64

const (
	AvailabilityStatusOne AvailabilityStatus = 1
)

func PossibleValuesForAvailabilityStatus() []int64 {
	return []int64{
		int64(AvailabilityStatusOne),
	}
}

func parseAvailabilityStatus(input int64) (*AvailabilityStatus, error) {
	vals := map[int64]AvailabilityStatus{
		1: AvailabilityStatusOne,
	}
	if v, ok := vals[input]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AvailabilityStatus(input)
	return &out, nil
}

type ConnectivityType string

const (
	ConnectivityTypeIsConnectedQuery ConnectivityType = "IsConnectedQuery"
)

func PossibleValuesForConnectivityType() []string {
	return []string{
		string(ConnectivityTypeIsConnectedQuery),
	}
}

func parseConnectivityType(input string) (*ConnectivityType, error) {
	vals := map[string]ConnectivityType{
		"isconnectedquery": ConnectivityTypeIsConnectedQuery,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectivityType(input)
	return &out, nil
}

type DataConnectorKind string

const (
	DataConnectorKindAPIPolling                                DataConnectorKind = "APIPolling"
	DataConnectorKindAmazonWebServicesCloudTrail               DataConnectorKind = "AmazonWebServicesCloudTrail"
	DataConnectorKindAmazonWebServicesSThree                   DataConnectorKind = "AmazonWebServicesS3"
	DataConnectorKindAzureActiveDirectory                      DataConnectorKind = "AzureActiveDirectory"
	DataConnectorKindAzureAdvancedThreatProtection             DataConnectorKind = "AzureAdvancedThreatProtection"
	DataConnectorKindAzureSecurityCenter                       DataConnectorKind = "AzureSecurityCenter"
	DataConnectorKindDynamicsThreeSixFive                      DataConnectorKind = "Dynamics365"
	DataConnectorKindGenericUI                                 DataConnectorKind = "GenericUI"
	DataConnectorKindIOT                                       DataConnectorKind = "IOT"
	DataConnectorKindMicrosoftCloudAppSecurity                 DataConnectorKind = "MicrosoftCloudAppSecurity"
	DataConnectorKindMicrosoftDefenderAdvancedThreatProtection DataConnectorKind = "MicrosoftDefenderAdvancedThreatProtection"
	DataConnectorKindMicrosoftThreatIntelligence               DataConnectorKind = "MicrosoftThreatIntelligence"
	DataConnectorKindMicrosoftThreatProtection                 DataConnectorKind = "MicrosoftThreatProtection"
	DataConnectorKindOfficeATP                                 DataConnectorKind = "OfficeATP"
	DataConnectorKindOfficeIRM                                 DataConnectorKind = "OfficeIRM"
	DataConnectorKindOfficePowerBI                             DataConnectorKind = "OfficePowerBI"
	DataConnectorKindOfficeThreeSixFive                        DataConnectorKind = "Office365"
	DataConnectorKindOfficeThreeSixFiveProject                 DataConnectorKind = "Office365Project"
	DataConnectorKindThreatIntelligence                        DataConnectorKind = "ThreatIntelligence"
	DataConnectorKindThreatIntelligenceTaxii                   DataConnectorKind = "ThreatIntelligenceTaxii"
)

func PossibleValuesForDataConnectorKind() []string {
	return []string{
		string(DataConnectorKindAPIPolling),
		string(DataConnectorKindAmazonWebServicesCloudTrail),
		string(DataConnectorKindAmazonWebServicesSThree),
		string(DataConnectorKindAzureActiveDirectory),
		string(DataConnectorKindAzureAdvancedThreatProtection),
		string(DataConnectorKindAzureSecurityCenter),
		string(DataConnectorKindDynamicsThreeSixFive),
		string(DataConnectorKindGenericUI),
		string(DataConnectorKindIOT),
		string(DataConnectorKindMicrosoftCloudAppSecurity),
		string(DataConnectorKindMicrosoftDefenderAdvancedThreatProtection),
		string(DataConnectorKindMicrosoftThreatIntelligence),
		string(DataConnectorKindMicrosoftThreatProtection),
		string(DataConnectorKindOfficeATP),
		string(DataConnectorKindOfficeIRM),
		string(DataConnectorKindOfficePowerBI),
		string(DataConnectorKindOfficeThreeSixFive),
		string(DataConnectorKindOfficeThreeSixFiveProject),
		string(DataConnectorKindThreatIntelligence),
		string(DataConnectorKindThreatIntelligenceTaxii),
	}
}

func parseDataConnectorKind(input string) (*DataConnectorKind, error) {
	vals := map[string]DataConnectorKind{
		"apipolling":                    DataConnectorKindAPIPolling,
		"amazonwebservicescloudtrail":   DataConnectorKindAmazonWebServicesCloudTrail,
		"amazonwebservicess3":           DataConnectorKindAmazonWebServicesSThree,
		"azureactivedirectory":          DataConnectorKindAzureActiveDirectory,
		"azureadvancedthreatprotection": DataConnectorKindAzureAdvancedThreatProtection,
		"azuresecuritycenter":           DataConnectorKindAzureSecurityCenter,
		"dynamics365":                   DataConnectorKindDynamicsThreeSixFive,
		"genericui":                     DataConnectorKindGenericUI,
		"iot":                           DataConnectorKindIOT,
		"microsoftcloudappsecurity":     DataConnectorKindMicrosoftCloudAppSecurity,
		"microsoftdefenderadvancedthreatprotection": DataConnectorKindMicrosoftDefenderAdvancedThreatProtection,
		"microsoftthreatintelligence":               DataConnectorKindMicrosoftThreatIntelligence,
		"microsoftthreatprotection":                 DataConnectorKindMicrosoftThreatProtection,
		"officeatp":                                 DataConnectorKindOfficeATP,
		"officeirm":                                 DataConnectorKindOfficeIRM,
		"officepowerbi":                             DataConnectorKindOfficePowerBI,
		"office365":                                 DataConnectorKindOfficeThreeSixFive,
		"office365project":                          DataConnectorKindOfficeThreeSixFiveProject,
		"threatintelligence":                        DataConnectorKindThreatIntelligence,
		"threatintelligencetaxii":                   DataConnectorKindThreatIntelligenceTaxii,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataConnectorKind(input)
	return &out, nil
}

type DataTypeState string

const (
	DataTypeStateDisabled DataTypeState = "Disabled"
	DataTypeStateEnabled  DataTypeState = "Enabled"
)

func PossibleValuesForDataTypeState() []string {
	return []string{
		string(DataTypeStateDisabled),
		string(DataTypeStateEnabled),
	}
}

func parseDataTypeState(input string) (*DataTypeState, error) {
	vals := map[string]DataTypeState{
		"disabled": DataTypeStateDisabled,
		"enabled":  DataTypeStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataTypeState(input)
	return &out, nil
}

type PermissionProviderScope string

const (
	PermissionProviderScopeResourceGroup PermissionProviderScope = "ResourceGroup"
	PermissionProviderScopeSubscription  PermissionProviderScope = "Subscription"
	PermissionProviderScopeWorkspace     PermissionProviderScope = "Workspace"
)

func PossibleValuesForPermissionProviderScope() []string {
	return []string{
		string(PermissionProviderScopeResourceGroup),
		string(PermissionProviderScopeSubscription),
		string(PermissionProviderScopeWorkspace),
	}
}

func parsePermissionProviderScope(input string) (*PermissionProviderScope, error) {
	vals := map[string]PermissionProviderScope{
		"resourcegroup": PermissionProviderScopeResourceGroup,
		"subscription":  PermissionProviderScopeSubscription,
		"workspace":     PermissionProviderScopeWorkspace,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PermissionProviderScope(input)
	return &out, nil
}

type PollingFrequency string

const (
	PollingFrequencyOnceADay    PollingFrequency = "OnceADay"
	PollingFrequencyOnceAMinute PollingFrequency = "OnceAMinute"
	PollingFrequencyOnceAnHour  PollingFrequency = "OnceAnHour"
)

func PossibleValuesForPollingFrequency() []string {
	return []string{
		string(PollingFrequencyOnceADay),
		string(PollingFrequencyOnceAMinute),
		string(PollingFrequencyOnceAnHour),
	}
}

func parsePollingFrequency(input string) (*PollingFrequency, error) {
	vals := map[string]PollingFrequency{
		"onceaday":    PollingFrequencyOnceADay,
		"onceaminute": PollingFrequencyOnceAMinute,
		"onceanhour":  PollingFrequencyOnceAnHour,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PollingFrequency(input)
	return &out, nil
}

type ProviderName string

const (
	ProviderNameMicrosoftPointAuthorizationPolicyAssignments           ProviderName = "Microsoft.Authorization/policyAssignments"
	ProviderNameMicrosoftPointOperationalInsightsSolutions             ProviderName = "Microsoft.OperationalInsights/solutions"
	ProviderNameMicrosoftPointOperationalInsightsWorkspaces            ProviderName = "Microsoft.OperationalInsights/workspaces"
	ProviderNameMicrosoftPointOperationalInsightsWorkspacesDatasources ProviderName = "Microsoft.OperationalInsights/workspaces/datasources"
	ProviderNameMicrosoftPointOperationalInsightsWorkspacesSharedKeys  ProviderName = "Microsoft.OperationalInsights/workspaces/sharedKeys"
	ProviderNameMicrosoftPointaadiamDiagnosticSettings                 ProviderName = "microsoft.aadiam/diagnosticSettings"
)

func PossibleValuesForProviderName() []string {
	return []string{
		string(ProviderNameMicrosoftPointAuthorizationPolicyAssignments),
		string(ProviderNameMicrosoftPointOperationalInsightsSolutions),
		string(ProviderNameMicrosoftPointOperationalInsightsWorkspaces),
		string(ProviderNameMicrosoftPointOperationalInsightsWorkspacesDatasources),
		string(ProviderNameMicrosoftPointOperationalInsightsWorkspacesSharedKeys),
		string(ProviderNameMicrosoftPointaadiamDiagnosticSettings),
	}
}

func parseProviderName(input string) (*ProviderName, error) {
	vals := map[string]ProviderName{
		"microsoft.authorization/policyassignments":            ProviderNameMicrosoftPointAuthorizationPolicyAssignments,
		"microsoft.operationalinsights/solutions":              ProviderNameMicrosoftPointOperationalInsightsSolutions,
		"microsoft.operationalinsights/workspaces":             ProviderNameMicrosoftPointOperationalInsightsWorkspaces,
		"microsoft.operationalinsights/workspaces/datasources": ProviderNameMicrosoftPointOperationalInsightsWorkspacesDatasources,
		"microsoft.operationalinsights/workspaces/sharedkeys":  ProviderNameMicrosoftPointOperationalInsightsWorkspacesSharedKeys,
		"microsoft.aadiam/diagnosticsettings":                  ProviderNameMicrosoftPointaadiamDiagnosticSettings,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProviderName(input)
	return &out, nil
}

type SettingType string

const (
	SettingTypeCopyableLabel         SettingType = "CopyableLabel"
	SettingTypeInfoMessage           SettingType = "InfoMessage"
	SettingTypeInstructionStepsGroup SettingType = "InstructionStepsGroup"
)

func PossibleValuesForSettingType() []string {
	return []string{
		string(SettingTypeCopyableLabel),
		string(SettingTypeInfoMessage),
		string(SettingTypeInstructionStepsGroup),
	}
}

func parseSettingType(input string) (*SettingType, error) {
	vals := map[string]SettingType{
		"copyablelabel":         SettingTypeCopyableLabel,
		"infomessage":           SettingTypeInfoMessage,
		"instructionstepsgroup": SettingTypeInstructionStepsGroup,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SettingType(input)
	return &out, nil
}
