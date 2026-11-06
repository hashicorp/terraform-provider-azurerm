package goldengatedeployments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureResourceProvisioningState string

const (
	AzureResourceProvisioningStateCanceled     AzureResourceProvisioningState = "Canceled"
	AzureResourceProvisioningStateFailed       AzureResourceProvisioningState = "Failed"
	AzureResourceProvisioningStateProvisioning AzureResourceProvisioningState = "Provisioning"
	AzureResourceProvisioningStateSucceeded    AzureResourceProvisioningState = "Succeeded"
)

func PossibleValuesForAzureResourceProvisioningState() []string {
	return []string{
		string(AzureResourceProvisioningStateCanceled),
		string(AzureResourceProvisioningStateFailed),
		string(AzureResourceProvisioningStateProvisioning),
		string(AzureResourceProvisioningStateSucceeded),
	}
}

func (s *AzureResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureResourceProvisioningState(input string) (*AzureResourceProvisioningState, error) {
	vals := map[string]AzureResourceProvisioningState{
		"canceled":     AzureResourceProvisioningStateCanceled,
		"failed":       AzureResourceProvisioningStateFailed,
		"provisioning": AzureResourceProvisioningStateProvisioning,
		"succeeded":    AzureResourceProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureResourceProvisioningState(input)
	return &out, nil
}

type CategoryType string

const (
	CategoryTypeDataReplication CategoryType = "DataReplication"
	CategoryTypeDataTransforms  CategoryType = "DataTransforms"
	CategoryTypeStreamAnalytics CategoryType = "StreamAnalytics"
)

func PossibleValuesForCategoryType() []string {
	return []string{
		string(CategoryTypeDataReplication),
		string(CategoryTypeDataTransforms),
		string(CategoryTypeStreamAnalytics),
	}
}

func (s *CategoryType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCategoryType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCategoryType(input string) (*CategoryType, error) {
	vals := map[string]CategoryType{
		"datareplication": CategoryTypeDataReplication,
		"datatransforms":  CategoryTypeDataTransforms,
		"streamanalytics": CategoryTypeStreamAnalytics,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CategoryType(input)
	return &out, nil
}

type CredentialType string

const (
	CredentialTypeGoldenGate CredentialType = "GoldenGate"
	CredentialTypeIAM        CredentialType = "IAM"
)

func PossibleValuesForCredentialType() []string {
	return []string{
		string(CredentialTypeGoldenGate),
		string(CredentialTypeIAM),
	}
}

func (s *CredentialType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCredentialType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCredentialType(input string) (*CredentialType, error) {
	vals := map[string]CredentialType{
		"goldengate": CredentialTypeGoldenGate,
		"iam":        CredentialTypeIAM,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CredentialType(input)
	return &out, nil
}

type DayOfWeekName string

const (
	DayOfWeekNameFriday    DayOfWeekName = "Friday"
	DayOfWeekNameMonday    DayOfWeekName = "Monday"
	DayOfWeekNameSaturday  DayOfWeekName = "Saturday"
	DayOfWeekNameSunday    DayOfWeekName = "Sunday"
	DayOfWeekNameThursday  DayOfWeekName = "Thursday"
	DayOfWeekNameTuesday   DayOfWeekName = "Tuesday"
	DayOfWeekNameWednesday DayOfWeekName = "Wednesday"
)

func PossibleValuesForDayOfWeekName() []string {
	return []string{
		string(DayOfWeekNameFriday),
		string(DayOfWeekNameMonday),
		string(DayOfWeekNameSaturday),
		string(DayOfWeekNameSunday),
		string(DayOfWeekNameThursday),
		string(DayOfWeekNameTuesday),
		string(DayOfWeekNameWednesday),
	}
}

func (s *DayOfWeekName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDayOfWeekName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDayOfWeekName(input string) (*DayOfWeekName, error) {
	vals := map[string]DayOfWeekName{
		"friday":    DayOfWeekNameFriday,
		"monday":    DayOfWeekNameMonday,
		"saturday":  DayOfWeekNameSaturday,
		"sunday":    DayOfWeekNameSunday,
		"thursday":  DayOfWeekNameThursday,
		"tuesday":   DayOfWeekNameTuesday,
		"wednesday": DayOfWeekNameWednesday,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DayOfWeekName(input)
	return &out, nil
}

type DeploymentLifecycleState string

const (
	DeploymentLifecycleStateActive         DeploymentLifecycleState = "Active"
	DeploymentLifecycleStateCanceled       DeploymentLifecycleState = "Canceled"
	DeploymentLifecycleStateCanceling      DeploymentLifecycleState = "Canceling"
	DeploymentLifecycleStateCreating       DeploymentLifecycleState = "Creating"
	DeploymentLifecycleStateDeleted        DeploymentLifecycleState = "Deleted"
	DeploymentLifecycleStateDeleting       DeploymentLifecycleState = "Deleting"
	DeploymentLifecycleStateFailed         DeploymentLifecycleState = "Failed"
	DeploymentLifecycleStateInActive       DeploymentLifecycleState = "InActive"
	DeploymentLifecycleStateInProgress     DeploymentLifecycleState = "In Progress"
	DeploymentLifecycleStateNeedsAttention DeploymentLifecycleState = "Needs Attention"
	DeploymentLifecycleStateSucceeded      DeploymentLifecycleState = "Succeeded"
	DeploymentLifecycleStateUpdating       DeploymentLifecycleState = "Updating"
	DeploymentLifecycleStateWaiting        DeploymentLifecycleState = "Waiting"
)

func PossibleValuesForDeploymentLifecycleState() []string {
	return []string{
		string(DeploymentLifecycleStateActive),
		string(DeploymentLifecycleStateCanceled),
		string(DeploymentLifecycleStateCanceling),
		string(DeploymentLifecycleStateCreating),
		string(DeploymentLifecycleStateDeleted),
		string(DeploymentLifecycleStateDeleting),
		string(DeploymentLifecycleStateFailed),
		string(DeploymentLifecycleStateInActive),
		string(DeploymentLifecycleStateInProgress),
		string(DeploymentLifecycleStateNeedsAttention),
		string(DeploymentLifecycleStateSucceeded),
		string(DeploymentLifecycleStateUpdating),
		string(DeploymentLifecycleStateWaiting),
	}
}

func (s *DeploymentLifecycleState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeploymentLifecycleState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeploymentLifecycleState(input string) (*DeploymentLifecycleState, error) {
	vals := map[string]DeploymentLifecycleState{
		"active":          DeploymentLifecycleStateActive,
		"canceled":        DeploymentLifecycleStateCanceled,
		"canceling":       DeploymentLifecycleStateCanceling,
		"creating":        DeploymentLifecycleStateCreating,
		"deleted":         DeploymentLifecycleStateDeleted,
		"deleting":        DeploymentLifecycleStateDeleting,
		"failed":          DeploymentLifecycleStateFailed,
		"inactive":        DeploymentLifecycleStateInActive,
		"in progress":     DeploymentLifecycleStateInProgress,
		"needs attention": DeploymentLifecycleStateNeedsAttention,
		"succeeded":       DeploymentLifecycleStateSucceeded,
		"updating":        DeploymentLifecycleStateUpdating,
		"waiting":         DeploymentLifecycleStateWaiting,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeploymentLifecycleState(input)
	return &out, nil
}

type DeploymentType string

const (
	DeploymentTypeBigData                    DeploymentType = "BigData"
	DeploymentTypeDATABASEDBTwoI             DeploymentType = "DATABASE_DB2I"
	DeploymentTypeDataTransforms             DeploymentType = "DataTransforms"
	DeploymentTypeDatabaseDBTwoZOS           DeploymentType = "DatabaseDB2ZOS"
	DeploymentTypeDatabaseMicrosoftSQLServer DeploymentType = "DatabaseMicrosoftSQLServer"
	DeploymentTypeDatabaseMySQL              DeploymentType = "DatabaseMySQL"
	DeploymentTypeDatabaseOracle             DeploymentType = "DatabaseOracle"
	DeploymentTypeDatabasePostGreSQL         DeploymentType = "DatabasePostGreSQL"
	DeploymentTypeGGSA                       DeploymentType = "GGSA"
	DeploymentTypeOgg                        DeploymentType = "Ogg"
)

func PossibleValuesForDeploymentType() []string {
	return []string{
		string(DeploymentTypeBigData),
		string(DeploymentTypeDATABASEDBTwoI),
		string(DeploymentTypeDataTransforms),
		string(DeploymentTypeDatabaseDBTwoZOS),
		string(DeploymentTypeDatabaseMicrosoftSQLServer),
		string(DeploymentTypeDatabaseMySQL),
		string(DeploymentTypeDatabaseOracle),
		string(DeploymentTypeDatabasePostGreSQL),
		string(DeploymentTypeGGSA),
		string(DeploymentTypeOgg),
	}
}

func (s *DeploymentType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeploymentType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeploymentType(input string) (*DeploymentType, error) {
	vals := map[string]DeploymentType{
		"bigdata":                    DeploymentTypeBigData,
		"database_db2i":              DeploymentTypeDATABASEDBTwoI,
		"datatransforms":             DeploymentTypeDataTransforms,
		"databasedb2zos":             DeploymentTypeDatabaseDBTwoZOS,
		"databasemicrosoftsqlserver": DeploymentTypeDatabaseMicrosoftSQLServer,
		"databasemysql":              DeploymentTypeDatabaseMySQL,
		"databaseoracle":             DeploymentTypeDatabaseOracle,
		"databasepostgresql":         DeploymentTypeDatabasePostGreSQL,
		"ggsa":                       DeploymentTypeGGSA,
		"ogg":                        DeploymentTypeOgg,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeploymentType(input)
	return &out, nil
}

type FrequencyType string

const (
	FrequencyTypeDaily   FrequencyType = "Daily"
	FrequencyTypeMonthly FrequencyType = "Monthly"
	FrequencyTypeWeekly  FrequencyType = "Weekly"
)

func PossibleValuesForFrequencyType() []string {
	return []string{
		string(FrequencyTypeDaily),
		string(FrequencyTypeMonthly),
		string(FrequencyTypeWeekly),
	}
}

func (s *FrequencyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFrequencyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFrequencyType(input string) (*FrequencyType, error) {
	vals := map[string]FrequencyType{
		"daily":   FrequencyTypeDaily,
		"monthly": FrequencyTypeMonthly,
		"weekly":  FrequencyTypeWeekly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FrequencyType(input)
	return &out, nil
}

type GoldenGateConnectionAssignmentLifecycleState string

const (
	GoldenGateConnectionAssignmentLifecycleStateACTIVE   GoldenGateConnectionAssignmentLifecycleState = "ACTIVE"
	GoldenGateConnectionAssignmentLifecycleStateCREATING GoldenGateConnectionAssignmentLifecycleState = "CREATING"
	GoldenGateConnectionAssignmentLifecycleStateDELETED  GoldenGateConnectionAssignmentLifecycleState = "DELETED"
	GoldenGateConnectionAssignmentLifecycleStateDELETING GoldenGateConnectionAssignmentLifecycleState = "DELETING"
	GoldenGateConnectionAssignmentLifecycleStateFAILED   GoldenGateConnectionAssignmentLifecycleState = "FAILED"
	GoldenGateConnectionAssignmentLifecycleStateUPDATING GoldenGateConnectionAssignmentLifecycleState = "UPDATING"
)

func PossibleValuesForGoldenGateConnectionAssignmentLifecycleState() []string {
	return []string{
		string(GoldenGateConnectionAssignmentLifecycleStateACTIVE),
		string(GoldenGateConnectionAssignmentLifecycleStateCREATING),
		string(GoldenGateConnectionAssignmentLifecycleStateDELETED),
		string(GoldenGateConnectionAssignmentLifecycleStateDELETING),
		string(GoldenGateConnectionAssignmentLifecycleStateFAILED),
		string(GoldenGateConnectionAssignmentLifecycleStateUPDATING),
	}
}

func (s *GoldenGateConnectionAssignmentLifecycleState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGoldenGateConnectionAssignmentLifecycleState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGoldenGateConnectionAssignmentLifecycleState(input string) (*GoldenGateConnectionAssignmentLifecycleState, error) {
	vals := map[string]GoldenGateConnectionAssignmentLifecycleState{
		"active":   GoldenGateConnectionAssignmentLifecycleStateACTIVE,
		"creating": GoldenGateConnectionAssignmentLifecycleStateCREATING,
		"deleted":  GoldenGateConnectionAssignmentLifecycleStateDELETED,
		"deleting": GoldenGateConnectionAssignmentLifecycleStateDELETING,
		"failed":   GoldenGateConnectionAssignmentLifecycleStateFAILED,
		"updating": GoldenGateConnectionAssignmentLifecycleStateUPDATING,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GoldenGateConnectionAssignmentLifecycleState(input)
	return &out, nil
}

type LicenseModel string

const (
	LicenseModelBringYourOwnLicense LicenseModel = "BringYourOwnLicense"
	LicenseModelLicenseIncluded     LicenseModel = "LicenseIncluded"
)

func PossibleValuesForLicenseModel() []string {
	return []string{
		string(LicenseModelBringYourOwnLicense),
		string(LicenseModelLicenseIncluded),
	}
}

func (s *LicenseModel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLicenseModel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLicenseModel(input string) (*LicenseModel, error) {
	vals := map[string]LicenseModel{
		"bringyourownlicense": LicenseModelBringYourOwnLicense,
		"licenseincluded":     LicenseModelLicenseIncluded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LicenseModel(input)
	return &out, nil
}

type SetupType string

const (
	SetupTypeDevelopmentOrTesting SetupType = "DevelopmentOrTesting"
	SetupTypeProduction           SetupType = "Production"
)

func PossibleValuesForSetupType() []string {
	return []string{
		string(SetupTypeDevelopmentOrTesting),
		string(SetupTypeProduction),
	}
}

func (s *SetupType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSetupType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSetupType(input string) (*SetupType, error) {
	vals := map[string]SetupType{
		"developmentortesting": SetupTypeDevelopmentOrTesting,
		"production":           SetupTypeProduction,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SetupType(input)
	return &out, nil
}
