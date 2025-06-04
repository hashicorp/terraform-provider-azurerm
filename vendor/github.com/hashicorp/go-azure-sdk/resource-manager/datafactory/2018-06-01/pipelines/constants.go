package pipelines

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActivityOnInactiveMarkAs string

const (
	ActivityOnInactiveMarkAsFailed    ActivityOnInactiveMarkAs = "Failed"
	ActivityOnInactiveMarkAsSkipped   ActivityOnInactiveMarkAs = "Skipped"
	ActivityOnInactiveMarkAsSucceeded ActivityOnInactiveMarkAs = "Succeeded"
)

func PossibleValuesForActivityOnInactiveMarkAs() []string {
	return []string{
		string(ActivityOnInactiveMarkAsFailed),
		string(ActivityOnInactiveMarkAsSkipped),
		string(ActivityOnInactiveMarkAsSucceeded),
	}
}

func (s *ActivityOnInactiveMarkAs) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseActivityOnInactiveMarkAs(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseActivityOnInactiveMarkAs(input string) (*ActivityOnInactiveMarkAs, error) {
	vals := map[string]ActivityOnInactiveMarkAs{
		"failed":    ActivityOnInactiveMarkAsFailed,
		"skipped":   ActivityOnInactiveMarkAsSkipped,
		"succeeded": ActivityOnInactiveMarkAsSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ActivityOnInactiveMarkAs(input)
	return &out, nil
}

type ActivityState string

const (
	ActivityStateActive   ActivityState = "Active"
	ActivityStateInactive ActivityState = "Inactive"
)

func PossibleValuesForActivityState() []string {
	return []string{
		string(ActivityStateActive),
		string(ActivityStateInactive),
	}
}

func (s *ActivityState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseActivityState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseActivityState(input string) (*ActivityState, error) {
	vals := map[string]ActivityState{
		"active":   ActivityStateActive,
		"inactive": ActivityStateInactive,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ActivityState(input)
	return &out, nil
}

type AzureFunctionActivityMethod string

const (
	AzureFunctionActivityMethodDELETE  AzureFunctionActivityMethod = "DELETE"
	AzureFunctionActivityMethodGET     AzureFunctionActivityMethod = "GET"
	AzureFunctionActivityMethodHEAD    AzureFunctionActivityMethod = "HEAD"
	AzureFunctionActivityMethodOPTIONS AzureFunctionActivityMethod = "OPTIONS"
	AzureFunctionActivityMethodPOST    AzureFunctionActivityMethod = "POST"
	AzureFunctionActivityMethodPUT     AzureFunctionActivityMethod = "PUT"
	AzureFunctionActivityMethodTRACE   AzureFunctionActivityMethod = "TRACE"
)

func PossibleValuesForAzureFunctionActivityMethod() []string {
	return []string{
		string(AzureFunctionActivityMethodDELETE),
		string(AzureFunctionActivityMethodGET),
		string(AzureFunctionActivityMethodHEAD),
		string(AzureFunctionActivityMethodOPTIONS),
		string(AzureFunctionActivityMethodPOST),
		string(AzureFunctionActivityMethodPUT),
		string(AzureFunctionActivityMethodTRACE),
	}
}

func (s *AzureFunctionActivityMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureFunctionActivityMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureFunctionActivityMethod(input string) (*AzureFunctionActivityMethod, error) {
	vals := map[string]AzureFunctionActivityMethod{
		"delete":  AzureFunctionActivityMethodDELETE,
		"get":     AzureFunctionActivityMethodGET,
		"head":    AzureFunctionActivityMethodHEAD,
		"options": AzureFunctionActivityMethodOPTIONS,
		"post":    AzureFunctionActivityMethodPOST,
		"put":     AzureFunctionActivityMethodPUT,
		"trace":   AzureFunctionActivityMethodTRACE,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureFunctionActivityMethod(input)
	return &out, nil
}

type AzurePostgreSqlWriteMethodEnum string

const (
	AzurePostgreSqlWriteMethodEnumBulkInsert  AzurePostgreSqlWriteMethodEnum = "BulkInsert"
	AzurePostgreSqlWriteMethodEnumCopyCommand AzurePostgreSqlWriteMethodEnum = "CopyCommand"
	AzurePostgreSqlWriteMethodEnumUpsert      AzurePostgreSqlWriteMethodEnum = "Upsert"
)

func PossibleValuesForAzurePostgreSqlWriteMethodEnum() []string {
	return []string{
		string(AzurePostgreSqlWriteMethodEnumBulkInsert),
		string(AzurePostgreSqlWriteMethodEnumCopyCommand),
		string(AzurePostgreSqlWriteMethodEnumUpsert),
	}
}

func (s *AzurePostgreSqlWriteMethodEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzurePostgreSqlWriteMethodEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzurePostgreSqlWriteMethodEnum(input string) (*AzurePostgreSqlWriteMethodEnum, error) {
	vals := map[string]AzurePostgreSqlWriteMethodEnum{
		"bulkinsert":  AzurePostgreSqlWriteMethodEnumBulkInsert,
		"copycommand": AzurePostgreSqlWriteMethodEnumCopyCommand,
		"upsert":      AzurePostgreSqlWriteMethodEnumUpsert,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzurePostgreSqlWriteMethodEnum(input)
	return &out, nil
}

type AzureSearchIndexWriteBehaviorType string

const (
	AzureSearchIndexWriteBehaviorTypeMerge  AzureSearchIndexWriteBehaviorType = "Merge"
	AzureSearchIndexWriteBehaviorTypeUpload AzureSearchIndexWriteBehaviorType = "Upload"
)

func PossibleValuesForAzureSearchIndexWriteBehaviorType() []string {
	return []string{
		string(AzureSearchIndexWriteBehaviorTypeMerge),
		string(AzureSearchIndexWriteBehaviorTypeUpload),
	}
}

func (s *AzureSearchIndexWriteBehaviorType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureSearchIndexWriteBehaviorType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureSearchIndexWriteBehaviorType(input string) (*AzureSearchIndexWriteBehaviorType, error) {
	vals := map[string]AzureSearchIndexWriteBehaviorType{
		"merge":  AzureSearchIndexWriteBehaviorTypeMerge,
		"upload": AzureSearchIndexWriteBehaviorTypeUpload,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureSearchIndexWriteBehaviorType(input)
	return &out, nil
}

type BigDataPoolReferenceType string

const (
	BigDataPoolReferenceTypeBigDataPoolReference BigDataPoolReferenceType = "BigDataPoolReference"
)

func PossibleValuesForBigDataPoolReferenceType() []string {
	return []string{
		string(BigDataPoolReferenceTypeBigDataPoolReference),
	}
}

func (s *BigDataPoolReferenceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBigDataPoolReferenceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBigDataPoolReferenceType(input string) (*BigDataPoolReferenceType, error) {
	vals := map[string]BigDataPoolReferenceType{
		"bigdatapoolreference": BigDataPoolReferenceTypeBigDataPoolReference,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BigDataPoolReferenceType(input)
	return &out, nil
}

type CassandraSourceReadConsistencyLevels string

const (
	CassandraSourceReadConsistencyLevelsALL         CassandraSourceReadConsistencyLevels = "ALL"
	CassandraSourceReadConsistencyLevelsEACHQUORUM  CassandraSourceReadConsistencyLevels = "EACH_QUORUM"
	CassandraSourceReadConsistencyLevelsLOCALONE    CassandraSourceReadConsistencyLevels = "LOCAL_ONE"
	CassandraSourceReadConsistencyLevelsLOCALQUORUM CassandraSourceReadConsistencyLevels = "LOCAL_QUORUM"
	CassandraSourceReadConsistencyLevelsLOCALSERIAL CassandraSourceReadConsistencyLevels = "LOCAL_SERIAL"
	CassandraSourceReadConsistencyLevelsONE         CassandraSourceReadConsistencyLevels = "ONE"
	CassandraSourceReadConsistencyLevelsQUORUM      CassandraSourceReadConsistencyLevels = "QUORUM"
	CassandraSourceReadConsistencyLevelsSERIAL      CassandraSourceReadConsistencyLevels = "SERIAL"
	CassandraSourceReadConsistencyLevelsTHREE       CassandraSourceReadConsistencyLevels = "THREE"
	CassandraSourceReadConsistencyLevelsTWO         CassandraSourceReadConsistencyLevels = "TWO"
)

func PossibleValuesForCassandraSourceReadConsistencyLevels() []string {
	return []string{
		string(CassandraSourceReadConsistencyLevelsALL),
		string(CassandraSourceReadConsistencyLevelsEACHQUORUM),
		string(CassandraSourceReadConsistencyLevelsLOCALONE),
		string(CassandraSourceReadConsistencyLevelsLOCALQUORUM),
		string(CassandraSourceReadConsistencyLevelsLOCALSERIAL),
		string(CassandraSourceReadConsistencyLevelsONE),
		string(CassandraSourceReadConsistencyLevelsQUORUM),
		string(CassandraSourceReadConsistencyLevelsSERIAL),
		string(CassandraSourceReadConsistencyLevelsTHREE),
		string(CassandraSourceReadConsistencyLevelsTWO),
	}
}

func (s *CassandraSourceReadConsistencyLevels) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCassandraSourceReadConsistencyLevels(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCassandraSourceReadConsistencyLevels(input string) (*CassandraSourceReadConsistencyLevels, error) {
	vals := map[string]CassandraSourceReadConsistencyLevels{
		"all":          CassandraSourceReadConsistencyLevelsALL,
		"each_quorum":  CassandraSourceReadConsistencyLevelsEACHQUORUM,
		"local_one":    CassandraSourceReadConsistencyLevelsLOCALONE,
		"local_quorum": CassandraSourceReadConsistencyLevelsLOCALQUORUM,
		"local_serial": CassandraSourceReadConsistencyLevelsLOCALSERIAL,
		"one":          CassandraSourceReadConsistencyLevelsONE,
		"quorum":       CassandraSourceReadConsistencyLevelsQUORUM,
		"serial":       CassandraSourceReadConsistencyLevelsSERIAL,
		"three":        CassandraSourceReadConsistencyLevelsTHREE,
		"two":          CassandraSourceReadConsistencyLevelsTWO,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CassandraSourceReadConsistencyLevels(input)
	return &out, nil
}

type ConfigurationType string

const (
	ConfigurationTypeArtifact   ConfigurationType = "Artifact"
	ConfigurationTypeCustomized ConfigurationType = "Customized"
	ConfigurationTypeDefault    ConfigurationType = "Default"
)

func PossibleValuesForConfigurationType() []string {
	return []string{
		string(ConfigurationTypeArtifact),
		string(ConfigurationTypeCustomized),
		string(ConfigurationTypeDefault),
	}
}

func (s *ConfigurationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConfigurationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConfigurationType(input string) (*ConfigurationType, error) {
	vals := map[string]ConfigurationType{
		"artifact":   ConfigurationTypeArtifact,
		"customized": ConfigurationTypeCustomized,
		"default":    ConfigurationTypeDefault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConfigurationType(input)
	return &out, nil
}

type CredentialReferenceType string

const (
	CredentialReferenceTypeCredentialReference CredentialReferenceType = "CredentialReference"
)

func PossibleValuesForCredentialReferenceType() []string {
	return []string{
		string(CredentialReferenceTypeCredentialReference),
	}
}

func (s *CredentialReferenceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCredentialReferenceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCredentialReferenceType(input string) (*CredentialReferenceType, error) {
	vals := map[string]CredentialReferenceType{
		"credentialreference": CredentialReferenceTypeCredentialReference,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CredentialReferenceType(input)
	return &out, nil
}

type DataFlowReferenceType string

const (
	DataFlowReferenceTypeDataFlowReference DataFlowReferenceType = "DataFlowReference"
)

func PossibleValuesForDataFlowReferenceType() []string {
	return []string{
		string(DataFlowReferenceTypeDataFlowReference),
	}
}

func (s *DataFlowReferenceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataFlowReferenceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataFlowReferenceType(input string) (*DataFlowReferenceType, error) {
	vals := map[string]DataFlowReferenceType{
		"dataflowreference": DataFlowReferenceTypeDataFlowReference,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataFlowReferenceType(input)
	return &out, nil
}

type DatasetReferenceType string

const (
	DatasetReferenceTypeDatasetReference DatasetReferenceType = "DatasetReference"
)

func PossibleValuesForDatasetReferenceType() []string {
	return []string{
		string(DatasetReferenceTypeDatasetReference),
	}
}

func (s *DatasetReferenceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDatasetReferenceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDatasetReferenceType(input string) (*DatasetReferenceType, error) {
	vals := map[string]DatasetReferenceType{
		"datasetreference": DatasetReferenceTypeDatasetReference,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatasetReferenceType(input)
	return &out, nil
}

type DependencyCondition string

const (
	DependencyConditionCompleted DependencyCondition = "Completed"
	DependencyConditionFailed    DependencyCondition = "Failed"
	DependencyConditionSkipped   DependencyCondition = "Skipped"
	DependencyConditionSucceeded DependencyCondition = "Succeeded"
)

func PossibleValuesForDependencyCondition() []string {
	return []string{
		string(DependencyConditionCompleted),
		string(DependencyConditionFailed),
		string(DependencyConditionSkipped),
		string(DependencyConditionSucceeded),
	}
}

func (s *DependencyCondition) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDependencyCondition(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDependencyCondition(input string) (*DependencyCondition, error) {
	vals := map[string]DependencyCondition{
		"completed": DependencyConditionCompleted,
		"failed":    DependencyConditionFailed,
		"skipped":   DependencyConditionSkipped,
		"succeeded": DependencyConditionSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DependencyCondition(input)
	return &out, nil
}

type DynamicsSinkWriteBehavior string

const (
	DynamicsSinkWriteBehaviorUpsert DynamicsSinkWriteBehavior = "Upsert"
)

func PossibleValuesForDynamicsSinkWriteBehavior() []string {
	return []string{
		string(DynamicsSinkWriteBehaviorUpsert),
	}
}

func (s *DynamicsSinkWriteBehavior) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDynamicsSinkWriteBehavior(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDynamicsSinkWriteBehavior(input string) (*DynamicsSinkWriteBehavior, error) {
	vals := map[string]DynamicsSinkWriteBehavior{
		"upsert": DynamicsSinkWriteBehaviorUpsert,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DynamicsSinkWriteBehavior(input)
	return &out, nil
}

type ExpressionType string

const (
	ExpressionTypeExpression ExpressionType = "Expression"
)

func PossibleValuesForExpressionType() []string {
	return []string{
		string(ExpressionTypeExpression),
	}
}

func (s *ExpressionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExpressionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExpressionType(input string) (*ExpressionType, error) {
	vals := map[string]ExpressionType{
		"expression": ExpressionTypeExpression,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExpressionType(input)
	return &out, nil
}

type ExpressionV2Type string

const (
	ExpressionV2TypeBinary   ExpressionV2Type = "Binary"
	ExpressionV2TypeConstant ExpressionV2Type = "Constant"
	ExpressionV2TypeField    ExpressionV2Type = "Field"
	ExpressionV2TypeNAry     ExpressionV2Type = "NAry"
	ExpressionV2TypeUnary    ExpressionV2Type = "Unary"
)

func PossibleValuesForExpressionV2Type() []string {
	return []string{
		string(ExpressionV2TypeBinary),
		string(ExpressionV2TypeConstant),
		string(ExpressionV2TypeField),
		string(ExpressionV2TypeNAry),
		string(ExpressionV2TypeUnary),
	}
}

func (s *ExpressionV2Type) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExpressionV2Type(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExpressionV2Type(input string) (*ExpressionV2Type, error) {
	vals := map[string]ExpressionV2Type{
		"binary":   ExpressionV2TypeBinary,
		"constant": ExpressionV2TypeConstant,
		"field":    ExpressionV2TypeField,
		"nary":     ExpressionV2TypeNAry,
		"unary":    ExpressionV2TypeUnary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExpressionV2Type(input)
	return &out, nil
}

type HDInsightActivityDebugInfoOption string

const (
	HDInsightActivityDebugInfoOptionAlways  HDInsightActivityDebugInfoOption = "Always"
	HDInsightActivityDebugInfoOptionFailure HDInsightActivityDebugInfoOption = "Failure"
	HDInsightActivityDebugInfoOptionNone    HDInsightActivityDebugInfoOption = "None"
)

func PossibleValuesForHDInsightActivityDebugInfoOption() []string {
	return []string{
		string(HDInsightActivityDebugInfoOptionAlways),
		string(HDInsightActivityDebugInfoOptionFailure),
		string(HDInsightActivityDebugInfoOptionNone),
	}
}

func (s *HDInsightActivityDebugInfoOption) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHDInsightActivityDebugInfoOption(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHDInsightActivityDebugInfoOption(input string) (*HDInsightActivityDebugInfoOption, error) {
	vals := map[string]HDInsightActivityDebugInfoOption{
		"always":  HDInsightActivityDebugInfoOptionAlways,
		"failure": HDInsightActivityDebugInfoOptionFailure,
		"none":    HDInsightActivityDebugInfoOptionNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HDInsightActivityDebugInfoOption(input)
	return &out, nil
}

type IntegrationRuntimeReferenceType string

const (
	IntegrationRuntimeReferenceTypeIntegrationRuntimeReference IntegrationRuntimeReferenceType = "IntegrationRuntimeReference"
)

func PossibleValuesForIntegrationRuntimeReferenceType() []string {
	return []string{
		string(IntegrationRuntimeReferenceTypeIntegrationRuntimeReference),
	}
}

func (s *IntegrationRuntimeReferenceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIntegrationRuntimeReferenceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIntegrationRuntimeReferenceType(input string) (*IntegrationRuntimeReferenceType, error) {
	vals := map[string]IntegrationRuntimeReferenceType{
		"integrationruntimereference": IntegrationRuntimeReferenceTypeIntegrationRuntimeReference,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IntegrationRuntimeReferenceType(input)
	return &out, nil
}

type NotebookParameterType string

const (
	NotebookParameterTypeBool   NotebookParameterType = "bool"
	NotebookParameterTypeFloat  NotebookParameterType = "float"
	NotebookParameterTypeInt    NotebookParameterType = "int"
	NotebookParameterTypeString NotebookParameterType = "string"
)

func PossibleValuesForNotebookParameterType() []string {
	return []string{
		string(NotebookParameterTypeBool),
		string(NotebookParameterTypeFloat),
		string(NotebookParameterTypeInt),
		string(NotebookParameterTypeString),
	}
}

func (s *NotebookParameterType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNotebookParameterType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNotebookParameterType(input string) (*NotebookParameterType, error) {
	vals := map[string]NotebookParameterType{
		"bool":   NotebookParameterTypeBool,
		"float":  NotebookParameterTypeFloat,
		"int":    NotebookParameterTypeInt,
		"string": NotebookParameterTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NotebookParameterType(input)
	return &out, nil
}

type NotebookReferenceType string

const (
	NotebookReferenceTypeNotebookReference NotebookReferenceType = "NotebookReference"
)

func PossibleValuesForNotebookReferenceType() []string {
	return []string{
		string(NotebookReferenceTypeNotebookReference),
	}
}

func (s *NotebookReferenceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNotebookReferenceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNotebookReferenceType(input string) (*NotebookReferenceType, error) {
	vals := map[string]NotebookReferenceType{
		"notebookreference": NotebookReferenceTypeNotebookReference,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NotebookReferenceType(input)
	return &out, nil
}

type ParameterType string

const (
	ParameterTypeArray        ParameterType = "Array"
	ParameterTypeBool         ParameterType = "Bool"
	ParameterTypeFloat        ParameterType = "Float"
	ParameterTypeInt          ParameterType = "Int"
	ParameterTypeObject       ParameterType = "Object"
	ParameterTypeSecureString ParameterType = "SecureString"
	ParameterTypeString       ParameterType = "String"
)

func PossibleValuesForParameterType() []string {
	return []string{
		string(ParameterTypeArray),
		string(ParameterTypeBool),
		string(ParameterTypeFloat),
		string(ParameterTypeInt),
		string(ParameterTypeObject),
		string(ParameterTypeSecureString),
		string(ParameterTypeString),
	}
}

func (s *ParameterType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseParameterType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseParameterType(input string) (*ParameterType, error) {
	vals := map[string]ParameterType{
		"array":        ParameterTypeArray,
		"bool":         ParameterTypeBool,
		"float":        ParameterTypeFloat,
		"int":          ParameterTypeInt,
		"object":       ParameterTypeObject,
		"securestring": ParameterTypeSecureString,
		"string":       ParameterTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ParameterType(input)
	return &out, nil
}

type PipelineReferenceType string

const (
	PipelineReferenceTypePipelineReference PipelineReferenceType = "PipelineReference"
)

func PossibleValuesForPipelineReferenceType() []string {
	return []string{
		string(PipelineReferenceTypePipelineReference),
	}
}

func (s *PipelineReferenceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePipelineReferenceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePipelineReferenceType(input string) (*PipelineReferenceType, error) {
	vals := map[string]PipelineReferenceType{
		"pipelinereference": PipelineReferenceTypePipelineReference,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PipelineReferenceType(input)
	return &out, nil
}

type PolybaseSettingsRejectType string

const (
	PolybaseSettingsRejectTypePercentage PolybaseSettingsRejectType = "percentage"
	PolybaseSettingsRejectTypeValue      PolybaseSettingsRejectType = "value"
)

func PossibleValuesForPolybaseSettingsRejectType() []string {
	return []string{
		string(PolybaseSettingsRejectTypePercentage),
		string(PolybaseSettingsRejectTypeValue),
	}
}

func (s *PolybaseSettingsRejectType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePolybaseSettingsRejectType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePolybaseSettingsRejectType(input string) (*PolybaseSettingsRejectType, error) {
	vals := map[string]PolybaseSettingsRejectType{
		"percentage": PolybaseSettingsRejectTypePercentage,
		"value":      PolybaseSettingsRejectTypeValue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PolybaseSettingsRejectType(input)
	return &out, nil
}

type SalesforceSinkWriteBehavior string

const (
	SalesforceSinkWriteBehaviorInsert SalesforceSinkWriteBehavior = "Insert"
	SalesforceSinkWriteBehaviorUpsert SalesforceSinkWriteBehavior = "Upsert"
)

func PossibleValuesForSalesforceSinkWriteBehavior() []string {
	return []string{
		string(SalesforceSinkWriteBehaviorInsert),
		string(SalesforceSinkWriteBehaviorUpsert),
	}
}

func (s *SalesforceSinkWriteBehavior) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSalesforceSinkWriteBehavior(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSalesforceSinkWriteBehavior(input string) (*SalesforceSinkWriteBehavior, error) {
	vals := map[string]SalesforceSinkWriteBehavior{
		"insert": SalesforceSinkWriteBehaviorInsert,
		"upsert": SalesforceSinkWriteBehaviorUpsert,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SalesforceSinkWriteBehavior(input)
	return &out, nil
}

type SalesforceV2SinkWriteBehavior string

const (
	SalesforceV2SinkWriteBehaviorInsert SalesforceV2SinkWriteBehavior = "Insert"
	SalesforceV2SinkWriteBehaviorUpsert SalesforceV2SinkWriteBehavior = "Upsert"
)

func PossibleValuesForSalesforceV2SinkWriteBehavior() []string {
	return []string{
		string(SalesforceV2SinkWriteBehaviorInsert),
		string(SalesforceV2SinkWriteBehaviorUpsert),
	}
}

func (s *SalesforceV2SinkWriteBehavior) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSalesforceV2SinkWriteBehavior(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSalesforceV2SinkWriteBehavior(input string) (*SalesforceV2SinkWriteBehavior, error) {
	vals := map[string]SalesforceV2SinkWriteBehavior{
		"insert": SalesforceV2SinkWriteBehaviorInsert,
		"upsert": SalesforceV2SinkWriteBehaviorUpsert,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SalesforceV2SinkWriteBehavior(input)
	return &out, nil
}

type SapCloudForCustomerSinkWriteBehavior string

const (
	SapCloudForCustomerSinkWriteBehaviorInsert SapCloudForCustomerSinkWriteBehavior = "Insert"
	SapCloudForCustomerSinkWriteBehaviorUpdate SapCloudForCustomerSinkWriteBehavior = "Update"
)

func PossibleValuesForSapCloudForCustomerSinkWriteBehavior() []string {
	return []string{
		string(SapCloudForCustomerSinkWriteBehaviorInsert),
		string(SapCloudForCustomerSinkWriteBehaviorUpdate),
	}
}

func (s *SapCloudForCustomerSinkWriteBehavior) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSapCloudForCustomerSinkWriteBehavior(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSapCloudForCustomerSinkWriteBehavior(input string) (*SapCloudForCustomerSinkWriteBehavior, error) {
	vals := map[string]SapCloudForCustomerSinkWriteBehavior{
		"insert": SapCloudForCustomerSinkWriteBehaviorInsert,
		"update": SapCloudForCustomerSinkWriteBehaviorUpdate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SapCloudForCustomerSinkWriteBehavior(input)
	return &out, nil
}

type ScriptActivityLogDestination string

const (
	ScriptActivityLogDestinationActivityOutput ScriptActivityLogDestination = "ActivityOutput"
	ScriptActivityLogDestinationExternalStore  ScriptActivityLogDestination = "ExternalStore"
)

func PossibleValuesForScriptActivityLogDestination() []string {
	return []string{
		string(ScriptActivityLogDestinationActivityOutput),
		string(ScriptActivityLogDestinationExternalStore),
	}
}

func (s *ScriptActivityLogDestination) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScriptActivityLogDestination(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScriptActivityLogDestination(input string) (*ScriptActivityLogDestination, error) {
	vals := map[string]ScriptActivityLogDestination{
		"activityoutput": ScriptActivityLogDestinationActivityOutput,
		"externalstore":  ScriptActivityLogDestinationExternalStore,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScriptActivityLogDestination(input)
	return &out, nil
}

type ScriptActivityParameterDirection string

const (
	ScriptActivityParameterDirectionInput       ScriptActivityParameterDirection = "Input"
	ScriptActivityParameterDirectionInputOutput ScriptActivityParameterDirection = "InputOutput"
	ScriptActivityParameterDirectionOutput      ScriptActivityParameterDirection = "Output"
)

func PossibleValuesForScriptActivityParameterDirection() []string {
	return []string{
		string(ScriptActivityParameterDirectionInput),
		string(ScriptActivityParameterDirectionInputOutput),
		string(ScriptActivityParameterDirectionOutput),
	}
}

func (s *ScriptActivityParameterDirection) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScriptActivityParameterDirection(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScriptActivityParameterDirection(input string) (*ScriptActivityParameterDirection, error) {
	vals := map[string]ScriptActivityParameterDirection{
		"input":       ScriptActivityParameterDirectionInput,
		"inputoutput": ScriptActivityParameterDirectionInputOutput,
		"output":      ScriptActivityParameterDirectionOutput,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScriptActivityParameterDirection(input)
	return &out, nil
}

type ScriptActivityParameterType string

const (
	ScriptActivityParameterTypeBoolean        ScriptActivityParameterType = "Boolean"
	ScriptActivityParameterTypeDateTime       ScriptActivityParameterType = "DateTime"
	ScriptActivityParameterTypeDateTimeOffset ScriptActivityParameterType = "DateTimeOffset"
	ScriptActivityParameterTypeDecimal        ScriptActivityParameterType = "Decimal"
	ScriptActivityParameterTypeDouble         ScriptActivityParameterType = "Double"
	ScriptActivityParameterTypeGuid           ScriptActivityParameterType = "Guid"
	ScriptActivityParameterTypeIntOneSix      ScriptActivityParameterType = "Int16"
	ScriptActivityParameterTypeIntSixFour     ScriptActivityParameterType = "Int64"
	ScriptActivityParameterTypeIntThreeTwo    ScriptActivityParameterType = "Int32"
	ScriptActivityParameterTypeSingle         ScriptActivityParameterType = "Single"
	ScriptActivityParameterTypeString         ScriptActivityParameterType = "String"
	ScriptActivityParameterTypeTimespan       ScriptActivityParameterType = "Timespan"
)

func PossibleValuesForScriptActivityParameterType() []string {
	return []string{
		string(ScriptActivityParameterTypeBoolean),
		string(ScriptActivityParameterTypeDateTime),
		string(ScriptActivityParameterTypeDateTimeOffset),
		string(ScriptActivityParameterTypeDecimal),
		string(ScriptActivityParameterTypeDouble),
		string(ScriptActivityParameterTypeGuid),
		string(ScriptActivityParameterTypeIntOneSix),
		string(ScriptActivityParameterTypeIntSixFour),
		string(ScriptActivityParameterTypeIntThreeTwo),
		string(ScriptActivityParameterTypeSingle),
		string(ScriptActivityParameterTypeString),
		string(ScriptActivityParameterTypeTimespan),
	}
}

func (s *ScriptActivityParameterType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScriptActivityParameterType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScriptActivityParameterType(input string) (*ScriptActivityParameterType, error) {
	vals := map[string]ScriptActivityParameterType{
		"boolean":        ScriptActivityParameterTypeBoolean,
		"datetime":       ScriptActivityParameterTypeDateTime,
		"datetimeoffset": ScriptActivityParameterTypeDateTimeOffset,
		"decimal":        ScriptActivityParameterTypeDecimal,
		"double":         ScriptActivityParameterTypeDouble,
		"guid":           ScriptActivityParameterTypeGuid,
		"int16":          ScriptActivityParameterTypeIntOneSix,
		"int64":          ScriptActivityParameterTypeIntSixFour,
		"int32":          ScriptActivityParameterTypeIntThreeTwo,
		"single":         ScriptActivityParameterTypeSingle,
		"string":         ScriptActivityParameterTypeString,
		"timespan":       ScriptActivityParameterTypeTimespan,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScriptActivityParameterType(input)
	return &out, nil
}

type SparkConfigurationReferenceType string

const (
	SparkConfigurationReferenceTypeSparkConfigurationReference SparkConfigurationReferenceType = "SparkConfigurationReference"
)

func PossibleValuesForSparkConfigurationReferenceType() []string {
	return []string{
		string(SparkConfigurationReferenceTypeSparkConfigurationReference),
	}
}

func (s *SparkConfigurationReferenceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSparkConfigurationReferenceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSparkConfigurationReferenceType(input string) (*SparkConfigurationReferenceType, error) {
	vals := map[string]SparkConfigurationReferenceType{
		"sparkconfigurationreference": SparkConfigurationReferenceTypeSparkConfigurationReference,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SparkConfigurationReferenceType(input)
	return &out, nil
}

type SparkJobReferenceType string

const (
	SparkJobReferenceTypeSparkJobDefinitionReference SparkJobReferenceType = "SparkJobDefinitionReference"
)

func PossibleValuesForSparkJobReferenceType() []string {
	return []string{
		string(SparkJobReferenceTypeSparkJobDefinitionReference),
	}
}

func (s *SparkJobReferenceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSparkJobReferenceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSparkJobReferenceType(input string) (*SparkJobReferenceType, error) {
	vals := map[string]SparkJobReferenceType{
		"sparkjobdefinitionreference": SparkJobReferenceTypeSparkJobDefinitionReference,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SparkJobReferenceType(input)
	return &out, nil
}

type SsisLogLocationType string

const (
	SsisLogLocationTypeFile SsisLogLocationType = "File"
)

func PossibleValuesForSsisLogLocationType() []string {
	return []string{
		string(SsisLogLocationTypeFile),
	}
}

func (s *SsisLogLocationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSsisLogLocationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSsisLogLocationType(input string) (*SsisLogLocationType, error) {
	vals := map[string]SsisLogLocationType{
		"file": SsisLogLocationTypeFile,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SsisLogLocationType(input)
	return &out, nil
}

type SsisPackageLocationType string

const (
	SsisPackageLocationTypeFile          SsisPackageLocationType = "File"
	SsisPackageLocationTypeInlinePackage SsisPackageLocationType = "InlinePackage"
	SsisPackageLocationTypePackageStore  SsisPackageLocationType = "PackageStore"
	SsisPackageLocationTypeSSISDB        SsisPackageLocationType = "SSISDB"
)

func PossibleValuesForSsisPackageLocationType() []string {
	return []string{
		string(SsisPackageLocationTypeFile),
		string(SsisPackageLocationTypeInlinePackage),
		string(SsisPackageLocationTypePackageStore),
		string(SsisPackageLocationTypeSSISDB),
	}
}

func (s *SsisPackageLocationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSsisPackageLocationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSsisPackageLocationType(input string) (*SsisPackageLocationType, error) {
	vals := map[string]SsisPackageLocationType{
		"file":          SsisPackageLocationTypeFile,
		"inlinepackage": SsisPackageLocationTypeInlinePackage,
		"packagestore":  SsisPackageLocationTypePackageStore,
		"ssisdb":        SsisPackageLocationTypeSSISDB,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SsisPackageLocationType(input)
	return &out, nil
}

type Type string

const (
	TypeLinkedServiceReference Type = "LinkedServiceReference"
)

func PossibleValuesForType() []string {
	return []string{
		string(TypeLinkedServiceReference),
	}
}

func (s *Type) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseType(input string) (*Type, error) {
	vals := map[string]Type{
		"linkedservicereference": TypeLinkedServiceReference,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Type(input)
	return &out, nil
}

type VariableType string

const (
	VariableTypeArray  VariableType = "Array"
	VariableTypeBool   VariableType = "Bool"
	VariableTypeString VariableType = "String"
)

func PossibleValuesForVariableType() []string {
	return []string{
		string(VariableTypeArray),
		string(VariableTypeBool),
		string(VariableTypeString),
	}
}

func (s *VariableType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVariableType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVariableType(input string) (*VariableType, error) {
	vals := map[string]VariableType{
		"array":  VariableTypeArray,
		"bool":   VariableTypeBool,
		"string": VariableTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VariableType(input)
	return &out, nil
}

type WebActivityMethod string

const (
	WebActivityMethodDELETE WebActivityMethod = "DELETE"
	WebActivityMethodGET    WebActivityMethod = "GET"
	WebActivityMethodPOST   WebActivityMethod = "POST"
	WebActivityMethodPUT    WebActivityMethod = "PUT"
)

func PossibleValuesForWebActivityMethod() []string {
	return []string{
		string(WebActivityMethodDELETE),
		string(WebActivityMethodGET),
		string(WebActivityMethodPOST),
		string(WebActivityMethodPUT),
	}
}

func (s *WebActivityMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWebActivityMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWebActivityMethod(input string) (*WebActivityMethod, error) {
	vals := map[string]WebActivityMethod{
		"delete": WebActivityMethodDELETE,
		"get":    WebActivityMethodGET,
		"post":   WebActivityMethodPOST,
		"put":    WebActivityMethodPUT,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WebActivityMethod(input)
	return &out, nil
}

type WebHookActivityMethod string

const (
	WebHookActivityMethodPOST WebHookActivityMethod = "POST"
)

func PossibleValuesForWebHookActivityMethod() []string {
	return []string{
		string(WebHookActivityMethodPOST),
	}
}

func (s *WebHookActivityMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWebHookActivityMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWebHookActivityMethod(input string) (*WebHookActivityMethod, error) {
	vals := map[string]WebHookActivityMethod{
		"post": WebHookActivityMethodPOST,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WebHookActivityMethod(input)
	return &out, nil
}
