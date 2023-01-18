package replicationrecoveryplans

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type A2ARpRecoveryPointType string

const (
	A2ARpRecoveryPointTypeLatest                      A2ARpRecoveryPointType = "Latest"
	A2ARpRecoveryPointTypeLatestApplicationConsistent A2ARpRecoveryPointType = "LatestApplicationConsistent"
	A2ARpRecoveryPointTypeLatestCrashConsistent       A2ARpRecoveryPointType = "LatestCrashConsistent"
	A2ARpRecoveryPointTypeLatestProcessed             A2ARpRecoveryPointType = "LatestProcessed"
)

func PossibleValuesForA2ARpRecoveryPointType() []string {
	return []string{
		string(A2ARpRecoveryPointTypeLatest),
		string(A2ARpRecoveryPointTypeLatestApplicationConsistent),
		string(A2ARpRecoveryPointTypeLatestCrashConsistent),
		string(A2ARpRecoveryPointTypeLatestProcessed),
	}
}

func parseA2ARpRecoveryPointType(input string) (*A2ARpRecoveryPointType, error) {
	vals := map[string]A2ARpRecoveryPointType{
		"latest":                      A2ARpRecoveryPointTypeLatest,
		"latestapplicationconsistent": A2ARpRecoveryPointTypeLatestApplicationConsistent,
		"latestcrashconsistent":       A2ARpRecoveryPointTypeLatestCrashConsistent,
		"latestprocessed":             A2ARpRecoveryPointTypeLatestProcessed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := A2ARpRecoveryPointType(input)
	return &out, nil
}

type AlternateLocationRecoveryOption string

const (
	AlternateLocationRecoveryOptionCreateVMIfNotFound AlternateLocationRecoveryOption = "CreateVmIfNotFound"
	AlternateLocationRecoveryOptionNoAction           AlternateLocationRecoveryOption = "NoAction"
)

func PossibleValuesForAlternateLocationRecoveryOption() []string {
	return []string{
		string(AlternateLocationRecoveryOptionCreateVMIfNotFound),
		string(AlternateLocationRecoveryOptionNoAction),
	}
}

func parseAlternateLocationRecoveryOption(input string) (*AlternateLocationRecoveryOption, error) {
	vals := map[string]AlternateLocationRecoveryOption{
		"createvmifnotfound": AlternateLocationRecoveryOptionCreateVMIfNotFound,
		"noaction":           AlternateLocationRecoveryOptionNoAction,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AlternateLocationRecoveryOption(input)
	return &out, nil
}

type DataSyncStatus string

const (
	DataSyncStatusForDownTime        DataSyncStatus = "ForDownTime"
	DataSyncStatusForSynchronization DataSyncStatus = "ForSynchronization"
)

func PossibleValuesForDataSyncStatus() []string {
	return []string{
		string(DataSyncStatusForDownTime),
		string(DataSyncStatusForSynchronization),
	}
}

func parseDataSyncStatus(input string) (*DataSyncStatus, error) {
	vals := map[string]DataSyncStatus{
		"fordowntime":        DataSyncStatusForDownTime,
		"forsynchronization": DataSyncStatusForSynchronization,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataSyncStatus(input)
	return &out, nil
}

type FailoverDeploymentModel string

const (
	FailoverDeploymentModelClassic         FailoverDeploymentModel = "Classic"
	FailoverDeploymentModelNotApplicable   FailoverDeploymentModel = "NotApplicable"
	FailoverDeploymentModelResourceManager FailoverDeploymentModel = "ResourceManager"
)

func PossibleValuesForFailoverDeploymentModel() []string {
	return []string{
		string(FailoverDeploymentModelClassic),
		string(FailoverDeploymentModelNotApplicable),
		string(FailoverDeploymentModelResourceManager),
	}
}

func parseFailoverDeploymentModel(input string) (*FailoverDeploymentModel, error) {
	vals := map[string]FailoverDeploymentModel{
		"classic":         FailoverDeploymentModelClassic,
		"notapplicable":   FailoverDeploymentModelNotApplicable,
		"resourcemanager": FailoverDeploymentModelResourceManager,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FailoverDeploymentModel(input)
	return &out, nil
}

type HyperVReplicaAzureRpRecoveryPointType string

const (
	HyperVReplicaAzureRpRecoveryPointTypeLatest                      HyperVReplicaAzureRpRecoveryPointType = "Latest"
	HyperVReplicaAzureRpRecoveryPointTypeLatestApplicationConsistent HyperVReplicaAzureRpRecoveryPointType = "LatestApplicationConsistent"
	HyperVReplicaAzureRpRecoveryPointTypeLatestProcessed             HyperVReplicaAzureRpRecoveryPointType = "LatestProcessed"
)

func PossibleValuesForHyperVReplicaAzureRpRecoveryPointType() []string {
	return []string{
		string(HyperVReplicaAzureRpRecoveryPointTypeLatest),
		string(HyperVReplicaAzureRpRecoveryPointTypeLatestApplicationConsistent),
		string(HyperVReplicaAzureRpRecoveryPointTypeLatestProcessed),
	}
}

func parseHyperVReplicaAzureRpRecoveryPointType(input string) (*HyperVReplicaAzureRpRecoveryPointType, error) {
	vals := map[string]HyperVReplicaAzureRpRecoveryPointType{
		"latest":                      HyperVReplicaAzureRpRecoveryPointTypeLatest,
		"latestapplicationconsistent": HyperVReplicaAzureRpRecoveryPointTypeLatestApplicationConsistent,
		"latestprocessed":             HyperVReplicaAzureRpRecoveryPointTypeLatestProcessed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HyperVReplicaAzureRpRecoveryPointType(input)
	return &out, nil
}

type InMageRcmFailbackRecoveryPointType string

const (
	InMageRcmFailbackRecoveryPointTypeApplicationConsistent InMageRcmFailbackRecoveryPointType = "ApplicationConsistent"
	InMageRcmFailbackRecoveryPointTypeCrashConsistent       InMageRcmFailbackRecoveryPointType = "CrashConsistent"
)

func PossibleValuesForInMageRcmFailbackRecoveryPointType() []string {
	return []string{
		string(InMageRcmFailbackRecoveryPointTypeApplicationConsistent),
		string(InMageRcmFailbackRecoveryPointTypeCrashConsistent),
	}
}

func parseInMageRcmFailbackRecoveryPointType(input string) (*InMageRcmFailbackRecoveryPointType, error) {
	vals := map[string]InMageRcmFailbackRecoveryPointType{
		"applicationconsistent": InMageRcmFailbackRecoveryPointTypeApplicationConsistent,
		"crashconsistent":       InMageRcmFailbackRecoveryPointTypeCrashConsistent,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InMageRcmFailbackRecoveryPointType(input)
	return &out, nil
}

type InMageV2RpRecoveryPointType string

const (
	InMageV2RpRecoveryPointTypeLatest                      InMageV2RpRecoveryPointType = "Latest"
	InMageV2RpRecoveryPointTypeLatestApplicationConsistent InMageV2RpRecoveryPointType = "LatestApplicationConsistent"
	InMageV2RpRecoveryPointTypeLatestCrashConsistent       InMageV2RpRecoveryPointType = "LatestCrashConsistent"
	InMageV2RpRecoveryPointTypeLatestProcessed             InMageV2RpRecoveryPointType = "LatestProcessed"
)

func PossibleValuesForInMageV2RpRecoveryPointType() []string {
	return []string{
		string(InMageV2RpRecoveryPointTypeLatest),
		string(InMageV2RpRecoveryPointTypeLatestApplicationConsistent),
		string(InMageV2RpRecoveryPointTypeLatestCrashConsistent),
		string(InMageV2RpRecoveryPointTypeLatestProcessed),
	}
}

func parseInMageV2RpRecoveryPointType(input string) (*InMageV2RpRecoveryPointType, error) {
	vals := map[string]InMageV2RpRecoveryPointType{
		"latest":                      InMageV2RpRecoveryPointTypeLatest,
		"latestapplicationconsistent": InMageV2RpRecoveryPointTypeLatestApplicationConsistent,
		"latestcrashconsistent":       InMageV2RpRecoveryPointTypeLatestCrashConsistent,
		"latestprocessed":             InMageV2RpRecoveryPointTypeLatestProcessed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InMageV2RpRecoveryPointType(input)
	return &out, nil
}

type MultiVMSyncPointOption string

const (
	MultiVMSyncPointOptionUseMultiVMSyncRecoveryPoint MultiVMSyncPointOption = "UseMultiVmSyncRecoveryPoint"
	MultiVMSyncPointOptionUsePerVMRecoveryPoint       MultiVMSyncPointOption = "UsePerVmRecoveryPoint"
)

func PossibleValuesForMultiVMSyncPointOption() []string {
	return []string{
		string(MultiVMSyncPointOptionUseMultiVMSyncRecoveryPoint),
		string(MultiVMSyncPointOptionUsePerVMRecoveryPoint),
	}
}

func parseMultiVMSyncPointOption(input string) (*MultiVMSyncPointOption, error) {
	vals := map[string]MultiVMSyncPointOption{
		"usemultivmsyncrecoverypoint": MultiVMSyncPointOptionUseMultiVMSyncRecoveryPoint,
		"usepervmrecoverypoint":       MultiVMSyncPointOptionUsePerVMRecoveryPoint,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MultiVMSyncPointOption(input)
	return &out, nil
}

type PossibleOperationsDirections string

const (
	PossibleOperationsDirectionsPrimaryToRecovery PossibleOperationsDirections = "PrimaryToRecovery"
	PossibleOperationsDirectionsRecoveryToPrimary PossibleOperationsDirections = "RecoveryToPrimary"
)

func PossibleValuesForPossibleOperationsDirections() []string {
	return []string{
		string(PossibleOperationsDirectionsPrimaryToRecovery),
		string(PossibleOperationsDirectionsRecoveryToPrimary),
	}
}

func parsePossibleOperationsDirections(input string) (*PossibleOperationsDirections, error) {
	vals := map[string]PossibleOperationsDirections{
		"primarytorecovery": PossibleOperationsDirectionsPrimaryToRecovery,
		"recoverytoprimary": PossibleOperationsDirectionsRecoveryToPrimary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PossibleOperationsDirections(input)
	return &out, nil
}

type RecoveryPlanActionLocation string

const (
	RecoveryPlanActionLocationPrimary  RecoveryPlanActionLocation = "Primary"
	RecoveryPlanActionLocationRecovery RecoveryPlanActionLocation = "Recovery"
)

func PossibleValuesForRecoveryPlanActionLocation() []string {
	return []string{
		string(RecoveryPlanActionLocationPrimary),
		string(RecoveryPlanActionLocationRecovery),
	}
}

func parseRecoveryPlanActionLocation(input string) (*RecoveryPlanActionLocation, error) {
	vals := map[string]RecoveryPlanActionLocation{
		"primary":  RecoveryPlanActionLocationPrimary,
		"recovery": RecoveryPlanActionLocationRecovery,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RecoveryPlanActionLocation(input)
	return &out, nil
}

type RecoveryPlanGroupType string

const (
	RecoveryPlanGroupTypeBoot     RecoveryPlanGroupType = "Boot"
	RecoveryPlanGroupTypeFailover RecoveryPlanGroupType = "Failover"
	RecoveryPlanGroupTypeShutdown RecoveryPlanGroupType = "Shutdown"
)

func PossibleValuesForRecoveryPlanGroupType() []string {
	return []string{
		string(RecoveryPlanGroupTypeBoot),
		string(RecoveryPlanGroupTypeFailover),
		string(RecoveryPlanGroupTypeShutdown),
	}
}

func parseRecoveryPlanGroupType(input string) (*RecoveryPlanGroupType, error) {
	vals := map[string]RecoveryPlanGroupType{
		"boot":     RecoveryPlanGroupTypeBoot,
		"failover": RecoveryPlanGroupTypeFailover,
		"shutdown": RecoveryPlanGroupTypeShutdown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RecoveryPlanGroupType(input)
	return &out, nil
}

type RecoveryPlanPointType string

const (
	RecoveryPlanPointTypeLatest                      RecoveryPlanPointType = "Latest"
	RecoveryPlanPointTypeLatestApplicationConsistent RecoveryPlanPointType = "LatestApplicationConsistent"
	RecoveryPlanPointTypeLatestCrashConsistent       RecoveryPlanPointType = "LatestCrashConsistent"
	RecoveryPlanPointTypeLatestProcessed             RecoveryPlanPointType = "LatestProcessed"
)

func PossibleValuesForRecoveryPlanPointType() []string {
	return []string{
		string(RecoveryPlanPointTypeLatest),
		string(RecoveryPlanPointTypeLatestApplicationConsistent),
		string(RecoveryPlanPointTypeLatestCrashConsistent),
		string(RecoveryPlanPointTypeLatestProcessed),
	}
}

func parseRecoveryPlanPointType(input string) (*RecoveryPlanPointType, error) {
	vals := map[string]RecoveryPlanPointType{
		"latest":                      RecoveryPlanPointTypeLatest,
		"latestapplicationconsistent": RecoveryPlanPointTypeLatestApplicationConsistent,
		"latestcrashconsistent":       RecoveryPlanPointTypeLatestCrashConsistent,
		"latestprocessed":             RecoveryPlanPointTypeLatestProcessed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RecoveryPlanPointType(input)
	return &out, nil
}

type ReplicationProtectedItemOperation string

const (
	ReplicationProtectedItemOperationCancelFailover      ReplicationProtectedItemOperation = "CancelFailover"
	ReplicationProtectedItemOperationChangePit           ReplicationProtectedItemOperation = "ChangePit"
	ReplicationProtectedItemOperationCommit              ReplicationProtectedItemOperation = "Commit"
	ReplicationProtectedItemOperationCompleteMigration   ReplicationProtectedItemOperation = "CompleteMigration"
	ReplicationProtectedItemOperationDisableProtection   ReplicationProtectedItemOperation = "DisableProtection"
	ReplicationProtectedItemOperationFailback            ReplicationProtectedItemOperation = "Failback"
	ReplicationProtectedItemOperationFinalizeFailback    ReplicationProtectedItemOperation = "FinalizeFailback"
	ReplicationProtectedItemOperationPlannedFailover     ReplicationProtectedItemOperation = "PlannedFailover"
	ReplicationProtectedItemOperationRepairReplication   ReplicationProtectedItemOperation = "RepairReplication"
	ReplicationProtectedItemOperationReverseReplicate    ReplicationProtectedItemOperation = "ReverseReplicate"
	ReplicationProtectedItemOperationSwitchProtection    ReplicationProtectedItemOperation = "SwitchProtection"
	ReplicationProtectedItemOperationTestFailover        ReplicationProtectedItemOperation = "TestFailover"
	ReplicationProtectedItemOperationTestFailoverCleanup ReplicationProtectedItemOperation = "TestFailoverCleanup"
	ReplicationProtectedItemOperationUnplannedFailover   ReplicationProtectedItemOperation = "UnplannedFailover"
)

func PossibleValuesForReplicationProtectedItemOperation() []string {
	return []string{
		string(ReplicationProtectedItemOperationCancelFailover),
		string(ReplicationProtectedItemOperationChangePit),
		string(ReplicationProtectedItemOperationCommit),
		string(ReplicationProtectedItemOperationCompleteMigration),
		string(ReplicationProtectedItemOperationDisableProtection),
		string(ReplicationProtectedItemOperationFailback),
		string(ReplicationProtectedItemOperationFinalizeFailback),
		string(ReplicationProtectedItemOperationPlannedFailover),
		string(ReplicationProtectedItemOperationRepairReplication),
		string(ReplicationProtectedItemOperationReverseReplicate),
		string(ReplicationProtectedItemOperationSwitchProtection),
		string(ReplicationProtectedItemOperationTestFailover),
		string(ReplicationProtectedItemOperationTestFailoverCleanup),
		string(ReplicationProtectedItemOperationUnplannedFailover),
	}
}

func parseReplicationProtectedItemOperation(input string) (*ReplicationProtectedItemOperation, error) {
	vals := map[string]ReplicationProtectedItemOperation{
		"cancelfailover":      ReplicationProtectedItemOperationCancelFailover,
		"changepit":           ReplicationProtectedItemOperationChangePit,
		"commit":              ReplicationProtectedItemOperationCommit,
		"completemigration":   ReplicationProtectedItemOperationCompleteMigration,
		"disableprotection":   ReplicationProtectedItemOperationDisableProtection,
		"failback":            ReplicationProtectedItemOperationFailback,
		"finalizefailback":    ReplicationProtectedItemOperationFinalizeFailback,
		"plannedfailover":     ReplicationProtectedItemOperationPlannedFailover,
		"repairreplication":   ReplicationProtectedItemOperationRepairReplication,
		"reversereplicate":    ReplicationProtectedItemOperationReverseReplicate,
		"switchprotection":    ReplicationProtectedItemOperationSwitchProtection,
		"testfailover":        ReplicationProtectedItemOperationTestFailover,
		"testfailovercleanup": ReplicationProtectedItemOperationTestFailoverCleanup,
		"unplannedfailover":   ReplicationProtectedItemOperationUnplannedFailover,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReplicationProtectedItemOperation(input)
	return &out, nil
}

type RpInMageRecoveryPointType string

const (
	RpInMageRecoveryPointTypeCustom     RpInMageRecoveryPointType = "Custom"
	RpInMageRecoveryPointTypeLatestTag  RpInMageRecoveryPointType = "LatestTag"
	RpInMageRecoveryPointTypeLatestTime RpInMageRecoveryPointType = "LatestTime"
)

func PossibleValuesForRpInMageRecoveryPointType() []string {
	return []string{
		string(RpInMageRecoveryPointTypeCustom),
		string(RpInMageRecoveryPointTypeLatestTag),
		string(RpInMageRecoveryPointTypeLatestTime),
	}
}

func parseRpInMageRecoveryPointType(input string) (*RpInMageRecoveryPointType, error) {
	vals := map[string]RpInMageRecoveryPointType{
		"custom":     RpInMageRecoveryPointTypeCustom,
		"latesttag":  RpInMageRecoveryPointTypeLatestTag,
		"latesttime": RpInMageRecoveryPointTypeLatestTime,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RpInMageRecoveryPointType(input)
	return &out, nil
}

type SourceSiteOperations string

const (
	SourceSiteOperationsNotRequired SourceSiteOperations = "NotRequired"
	SourceSiteOperationsRequired    SourceSiteOperations = "Required"
)

func PossibleValuesForSourceSiteOperations() []string {
	return []string{
		string(SourceSiteOperationsNotRequired),
		string(SourceSiteOperationsRequired),
	}
}

func parseSourceSiteOperations(input string) (*SourceSiteOperations, error) {
	vals := map[string]SourceSiteOperations{
		"notrequired": SourceSiteOperationsNotRequired,
		"required":    SourceSiteOperationsRequired,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SourceSiteOperations(input)
	return &out, nil
}
