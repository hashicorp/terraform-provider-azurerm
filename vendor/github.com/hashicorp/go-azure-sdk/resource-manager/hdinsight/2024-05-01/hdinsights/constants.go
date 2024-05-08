package hdinsights

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Action string

const (
	ActionCANCEL          Action = "CANCEL"
	ActionDELETE          Action = "DELETE"
	ActionLASTSTATEUPDATE Action = "LAST_STATE_UPDATE"
	ActionLISTSAVEPOINT   Action = "LIST_SAVEPOINT"
	ActionNEW             Action = "NEW"
	ActionRELAUNCH        Action = "RE_LAUNCH"
	ActionSAVEPOINT       Action = "SAVEPOINT"
	ActionSTART           Action = "START"
	ActionSTATELESSUPDATE Action = "STATELESS_UPDATE"
	ActionSTOP            Action = "STOP"
	ActionUPDATE          Action = "UPDATE"
)

func PossibleValuesForAction() []string {
	return []string{
		string(ActionCANCEL),
		string(ActionDELETE),
		string(ActionLASTSTATEUPDATE),
		string(ActionLISTSAVEPOINT),
		string(ActionNEW),
		string(ActionRELAUNCH),
		string(ActionSAVEPOINT),
		string(ActionSTART),
		string(ActionSTATELESSUPDATE),
		string(ActionSTOP),
		string(ActionUPDATE),
	}
}

func (s *Action) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAction(input string) (*Action, error) {
	vals := map[string]Action{
		"cancel":            ActionCANCEL,
		"delete":            ActionDELETE,
		"last_state_update": ActionLASTSTATEUPDATE,
		"list_savepoint":    ActionLISTSAVEPOINT,
		"new":               ActionNEW,
		"re_launch":         ActionRELAUNCH,
		"savepoint":         ActionSAVEPOINT,
		"start":             ActionSTART,
		"stateless_update":  ActionSTATELESSUPDATE,
		"stop":              ActionSTOP,
		"update":            ActionUPDATE,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Action(input)
	return &out, nil
}

type AutoscaleType string

const (
	AutoscaleTypeLoadBased     AutoscaleType = "LoadBased"
	AutoscaleTypeScheduleBased AutoscaleType = "ScheduleBased"
)

func PossibleValuesForAutoscaleType() []string {
	return []string{
		string(AutoscaleTypeLoadBased),
		string(AutoscaleTypeScheduleBased),
	}
}

func (s *AutoscaleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutoscaleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutoscaleType(input string) (*AutoscaleType, error) {
	vals := map[string]AutoscaleType{
		"loadbased":     AutoscaleTypeLoadBased,
		"schedulebased": AutoscaleTypeScheduleBased,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutoscaleType(input)
	return &out, nil
}

type Category string

const (
	CategoryCustom     Category = "custom"
	CategoryPredefined Category = "predefined"
)

func PossibleValuesForCategory() []string {
	return []string{
		string(CategoryCustom),
		string(CategoryPredefined),
	}
}

func (s *Category) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCategory(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCategory(input string) (*Category, error) {
	vals := map[string]Category{
		"custom":     CategoryCustom,
		"predefined": CategoryPredefined,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Category(input)
	return &out, nil
}

type ClusterAvailableUpgradeType string

const (
	ClusterAvailableUpgradeTypeAKSPatchUpgrade     ClusterAvailableUpgradeType = "AKSPatchUpgrade"
	ClusterAvailableUpgradeTypeHotfixUpgrade       ClusterAvailableUpgradeType = "HotfixUpgrade"
	ClusterAvailableUpgradeTypePatchVersionUpgrade ClusterAvailableUpgradeType = "PatchVersionUpgrade"
)

func PossibleValuesForClusterAvailableUpgradeType() []string {
	return []string{
		string(ClusterAvailableUpgradeTypeAKSPatchUpgrade),
		string(ClusterAvailableUpgradeTypeHotfixUpgrade),
		string(ClusterAvailableUpgradeTypePatchVersionUpgrade),
	}
}

func (s *ClusterAvailableUpgradeType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterAvailableUpgradeType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusterAvailableUpgradeType(input string) (*ClusterAvailableUpgradeType, error) {
	vals := map[string]ClusterAvailableUpgradeType{
		"akspatchupgrade":     ClusterAvailableUpgradeTypeAKSPatchUpgrade,
		"hotfixupgrade":       ClusterAvailableUpgradeTypeHotfixUpgrade,
		"patchversionupgrade": ClusterAvailableUpgradeTypePatchVersionUpgrade,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterAvailableUpgradeType(input)
	return &out, nil
}

type ClusterPoolAvailableUpgradeType string

const (
	ClusterPoolAvailableUpgradeTypeAKSPatchUpgrade ClusterPoolAvailableUpgradeType = "AKSPatchUpgrade"
	ClusterPoolAvailableUpgradeTypeNodeOsUpgrade   ClusterPoolAvailableUpgradeType = "NodeOsUpgrade"
)

func PossibleValuesForClusterPoolAvailableUpgradeType() []string {
	return []string{
		string(ClusterPoolAvailableUpgradeTypeAKSPatchUpgrade),
		string(ClusterPoolAvailableUpgradeTypeNodeOsUpgrade),
	}
}

func (s *ClusterPoolAvailableUpgradeType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterPoolAvailableUpgradeType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusterPoolAvailableUpgradeType(input string) (*ClusterPoolAvailableUpgradeType, error) {
	vals := map[string]ClusterPoolAvailableUpgradeType{
		"akspatchupgrade": ClusterPoolAvailableUpgradeTypeAKSPatchUpgrade,
		"nodeosupgrade":   ClusterPoolAvailableUpgradeTypeNodeOsUpgrade,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterPoolAvailableUpgradeType(input)
	return &out, nil
}

type ClusterPoolUpgradeHistoryType string

const (
	ClusterPoolUpgradeHistoryTypeAKSPatchUpgrade ClusterPoolUpgradeHistoryType = "AKSPatchUpgrade"
	ClusterPoolUpgradeHistoryTypeNodeOsUpgrade   ClusterPoolUpgradeHistoryType = "NodeOsUpgrade"
)

func PossibleValuesForClusterPoolUpgradeHistoryType() []string {
	return []string{
		string(ClusterPoolUpgradeHistoryTypeAKSPatchUpgrade),
		string(ClusterPoolUpgradeHistoryTypeNodeOsUpgrade),
	}
}

func (s *ClusterPoolUpgradeHistoryType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterPoolUpgradeHistoryType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusterPoolUpgradeHistoryType(input string) (*ClusterPoolUpgradeHistoryType, error) {
	vals := map[string]ClusterPoolUpgradeHistoryType{
		"akspatchupgrade": ClusterPoolUpgradeHistoryTypeAKSPatchUpgrade,
		"nodeosupgrade":   ClusterPoolUpgradeHistoryTypeNodeOsUpgrade,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterPoolUpgradeHistoryType(input)
	return &out, nil
}

type ClusterPoolUpgradeHistoryUpgradeResultType string

const (
	ClusterPoolUpgradeHistoryUpgradeResultTypeFailed  ClusterPoolUpgradeHistoryUpgradeResultType = "Failed"
	ClusterPoolUpgradeHistoryUpgradeResultTypeSucceed ClusterPoolUpgradeHistoryUpgradeResultType = "Succeed"
)

func PossibleValuesForClusterPoolUpgradeHistoryUpgradeResultType() []string {
	return []string{
		string(ClusterPoolUpgradeHistoryUpgradeResultTypeFailed),
		string(ClusterPoolUpgradeHistoryUpgradeResultTypeSucceed),
	}
}

func (s *ClusterPoolUpgradeHistoryUpgradeResultType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterPoolUpgradeHistoryUpgradeResultType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusterPoolUpgradeHistoryUpgradeResultType(input string) (*ClusterPoolUpgradeHistoryUpgradeResultType, error) {
	vals := map[string]ClusterPoolUpgradeHistoryUpgradeResultType{
		"failed":  ClusterPoolUpgradeHistoryUpgradeResultTypeFailed,
		"succeed": ClusterPoolUpgradeHistoryUpgradeResultTypeSucceed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterPoolUpgradeHistoryUpgradeResultType(input)
	return &out, nil
}

type ClusterPoolUpgradeType string

const (
	ClusterPoolUpgradeTypeAKSPatchUpgrade ClusterPoolUpgradeType = "AKSPatchUpgrade"
	ClusterPoolUpgradeTypeNodeOsUpgrade   ClusterPoolUpgradeType = "NodeOsUpgrade"
)

func PossibleValuesForClusterPoolUpgradeType() []string {
	return []string{
		string(ClusterPoolUpgradeTypeAKSPatchUpgrade),
		string(ClusterPoolUpgradeTypeNodeOsUpgrade),
	}
}

func (s *ClusterPoolUpgradeType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterPoolUpgradeType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusterPoolUpgradeType(input string) (*ClusterPoolUpgradeType, error) {
	vals := map[string]ClusterPoolUpgradeType{
		"akspatchupgrade": ClusterPoolUpgradeTypeAKSPatchUpgrade,
		"nodeosupgrade":   ClusterPoolUpgradeTypeNodeOsUpgrade,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterPoolUpgradeType(input)
	return &out, nil
}

type ClusterUpgradeHistorySeverityType string

const (
	ClusterUpgradeHistorySeverityTypeCritical ClusterUpgradeHistorySeverityType = "critical"
	ClusterUpgradeHistorySeverityTypeHigh     ClusterUpgradeHistorySeverityType = "high"
	ClusterUpgradeHistorySeverityTypeLow      ClusterUpgradeHistorySeverityType = "low"
	ClusterUpgradeHistorySeverityTypeMedium   ClusterUpgradeHistorySeverityType = "medium"
)

func PossibleValuesForClusterUpgradeHistorySeverityType() []string {
	return []string{
		string(ClusterUpgradeHistorySeverityTypeCritical),
		string(ClusterUpgradeHistorySeverityTypeHigh),
		string(ClusterUpgradeHistorySeverityTypeLow),
		string(ClusterUpgradeHistorySeverityTypeMedium),
	}
}

func (s *ClusterUpgradeHistorySeverityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterUpgradeHistorySeverityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusterUpgradeHistorySeverityType(input string) (*ClusterUpgradeHistorySeverityType, error) {
	vals := map[string]ClusterUpgradeHistorySeverityType{
		"critical": ClusterUpgradeHistorySeverityTypeCritical,
		"high":     ClusterUpgradeHistorySeverityTypeHigh,
		"low":      ClusterUpgradeHistorySeverityTypeLow,
		"medium":   ClusterUpgradeHistorySeverityTypeMedium,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterUpgradeHistorySeverityType(input)
	return &out, nil
}

type ClusterUpgradeHistoryType string

const (
	ClusterUpgradeHistoryTypeAKSPatchUpgrade             ClusterUpgradeHistoryType = "AKSPatchUpgrade"
	ClusterUpgradeHistoryTypeHotfixUpgrade               ClusterUpgradeHistoryType = "HotfixUpgrade"
	ClusterUpgradeHistoryTypeHotfixUpgradeRollback       ClusterUpgradeHistoryType = "HotfixUpgradeRollback"
	ClusterUpgradeHistoryTypePatchVersionUpgrade         ClusterUpgradeHistoryType = "PatchVersionUpgrade"
	ClusterUpgradeHistoryTypePatchVersionUpgradeRollback ClusterUpgradeHistoryType = "PatchVersionUpgradeRollback"
)

func PossibleValuesForClusterUpgradeHistoryType() []string {
	return []string{
		string(ClusterUpgradeHistoryTypeAKSPatchUpgrade),
		string(ClusterUpgradeHistoryTypeHotfixUpgrade),
		string(ClusterUpgradeHistoryTypeHotfixUpgradeRollback),
		string(ClusterUpgradeHistoryTypePatchVersionUpgrade),
		string(ClusterUpgradeHistoryTypePatchVersionUpgradeRollback),
	}
}

func (s *ClusterUpgradeHistoryType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterUpgradeHistoryType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusterUpgradeHistoryType(input string) (*ClusterUpgradeHistoryType, error) {
	vals := map[string]ClusterUpgradeHistoryType{
		"akspatchupgrade":             ClusterUpgradeHistoryTypeAKSPatchUpgrade,
		"hotfixupgrade":               ClusterUpgradeHistoryTypeHotfixUpgrade,
		"hotfixupgraderollback":       ClusterUpgradeHistoryTypeHotfixUpgradeRollback,
		"patchversionupgrade":         ClusterUpgradeHistoryTypePatchVersionUpgrade,
		"patchversionupgraderollback": ClusterUpgradeHistoryTypePatchVersionUpgradeRollback,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterUpgradeHistoryType(input)
	return &out, nil
}

type ClusterUpgradeHistoryUpgradeResultType string

const (
	ClusterUpgradeHistoryUpgradeResultTypeFailed  ClusterUpgradeHistoryUpgradeResultType = "Failed"
	ClusterUpgradeHistoryUpgradeResultTypeSucceed ClusterUpgradeHistoryUpgradeResultType = "Succeed"
)

func PossibleValuesForClusterUpgradeHistoryUpgradeResultType() []string {
	return []string{
		string(ClusterUpgradeHistoryUpgradeResultTypeFailed),
		string(ClusterUpgradeHistoryUpgradeResultTypeSucceed),
	}
}

func (s *ClusterUpgradeHistoryUpgradeResultType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterUpgradeHistoryUpgradeResultType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusterUpgradeHistoryUpgradeResultType(input string) (*ClusterUpgradeHistoryUpgradeResultType, error) {
	vals := map[string]ClusterUpgradeHistoryUpgradeResultType{
		"failed":  ClusterUpgradeHistoryUpgradeResultTypeFailed,
		"succeed": ClusterUpgradeHistoryUpgradeResultTypeSucceed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterUpgradeHistoryUpgradeResultType(input)
	return &out, nil
}

type ClusterUpgradeType string

const (
	ClusterUpgradeTypeAKSPatchUpgrade     ClusterUpgradeType = "AKSPatchUpgrade"
	ClusterUpgradeTypeHotfixUpgrade       ClusterUpgradeType = "HotfixUpgrade"
	ClusterUpgradeTypePatchVersionUpgrade ClusterUpgradeType = "PatchVersionUpgrade"
)

func PossibleValuesForClusterUpgradeType() []string {
	return []string{
		string(ClusterUpgradeTypeAKSPatchUpgrade),
		string(ClusterUpgradeTypeHotfixUpgrade),
		string(ClusterUpgradeTypePatchVersionUpgrade),
	}
}

func (s *ClusterUpgradeType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterUpgradeType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusterUpgradeType(input string) (*ClusterUpgradeType, error) {
	vals := map[string]ClusterUpgradeType{
		"akspatchupgrade":     ClusterUpgradeTypeAKSPatchUpgrade,
		"hotfixupgrade":       ClusterUpgradeTypeHotfixUpgrade,
		"patchversionupgrade": ClusterUpgradeTypePatchVersionUpgrade,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterUpgradeType(input)
	return &out, nil
}

type ComparisonOperator string

const (
	ComparisonOperatorGreaterThan        ComparisonOperator = "greaterThan"
	ComparisonOperatorGreaterThanOrEqual ComparisonOperator = "greaterThanOrEqual"
	ComparisonOperatorLessThan           ComparisonOperator = "lessThan"
	ComparisonOperatorLessThanOrEqual    ComparisonOperator = "lessThanOrEqual"
)

func PossibleValuesForComparisonOperator() []string {
	return []string{
		string(ComparisonOperatorGreaterThan),
		string(ComparisonOperatorGreaterThanOrEqual),
		string(ComparisonOperatorLessThan),
		string(ComparisonOperatorLessThanOrEqual),
	}
}

func (s *ComparisonOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseComparisonOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseComparisonOperator(input string) (*ComparisonOperator, error) {
	vals := map[string]ComparisonOperator{
		"greaterthan":        ComparisonOperatorGreaterThan,
		"greaterthanorequal": ComparisonOperatorGreaterThanOrEqual,
		"lessthan":           ComparisonOperatorLessThan,
		"lessthanorequal":    ComparisonOperatorLessThanOrEqual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ComparisonOperator(input)
	return &out, nil
}

type ContentEncoding string

const (
	ContentEncodingBaseSixFour ContentEncoding = "Base64"
	ContentEncodingNone        ContentEncoding = "None"
)

func PossibleValuesForContentEncoding() []string {
	return []string{
		string(ContentEncodingBaseSixFour),
		string(ContentEncodingNone),
	}
}

func (s *ContentEncoding) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContentEncoding(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContentEncoding(input string) (*ContentEncoding, error) {
	vals := map[string]ContentEncoding{
		"base64": ContentEncodingBaseSixFour,
		"none":   ContentEncodingNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContentEncoding(input)
	return &out, nil
}

type CurrentClusterAksVersionStatus string

const (
	CurrentClusterAksVersionStatusDeprecated CurrentClusterAksVersionStatus = "Deprecated"
	CurrentClusterAksVersionStatusSupported  CurrentClusterAksVersionStatus = "Supported"
)

func PossibleValuesForCurrentClusterAksVersionStatus() []string {
	return []string{
		string(CurrentClusterAksVersionStatusDeprecated),
		string(CurrentClusterAksVersionStatusSupported),
	}
}

func (s *CurrentClusterAksVersionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCurrentClusterAksVersionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCurrentClusterAksVersionStatus(input string) (*CurrentClusterAksVersionStatus, error) {
	vals := map[string]CurrentClusterAksVersionStatus{
		"deprecated": CurrentClusterAksVersionStatusDeprecated,
		"supported":  CurrentClusterAksVersionStatusSupported,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CurrentClusterAksVersionStatus(input)
	return &out, nil
}

type CurrentClusterPoolAksVersionStatus string

const (
	CurrentClusterPoolAksVersionStatusDeprecated CurrentClusterPoolAksVersionStatus = "Deprecated"
	CurrentClusterPoolAksVersionStatusSupported  CurrentClusterPoolAksVersionStatus = "Supported"
)

func PossibleValuesForCurrentClusterPoolAksVersionStatus() []string {
	return []string{
		string(CurrentClusterPoolAksVersionStatusDeprecated),
		string(CurrentClusterPoolAksVersionStatusSupported),
	}
}

func (s *CurrentClusterPoolAksVersionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCurrentClusterPoolAksVersionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCurrentClusterPoolAksVersionStatus(input string) (*CurrentClusterPoolAksVersionStatus, error) {
	vals := map[string]CurrentClusterPoolAksVersionStatus{
		"deprecated": CurrentClusterPoolAksVersionStatusDeprecated,
		"supported":  CurrentClusterPoolAksVersionStatusSupported,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CurrentClusterPoolAksVersionStatus(input)
	return &out, nil
}

type DataDiskType string

const (
	DataDiskTypePremiumSSDLRS     DataDiskType = "Premium_SSD_LRS"
	DataDiskTypePremiumSSDVTwoLRS DataDiskType = "Premium_SSD_v2_LRS"
	DataDiskTypePremiumSSDZRS     DataDiskType = "Premium_SSD_ZRS"
	DataDiskTypeStandardHDDLRS    DataDiskType = "Standard_HDD_LRS"
	DataDiskTypeStandardSSDLRS    DataDiskType = "Standard_SSD_LRS"
	DataDiskTypeStandardSSDZRS    DataDiskType = "Standard_SSD_ZRS"
)

func PossibleValuesForDataDiskType() []string {
	return []string{
		string(DataDiskTypePremiumSSDLRS),
		string(DataDiskTypePremiumSSDVTwoLRS),
		string(DataDiskTypePremiumSSDZRS),
		string(DataDiskTypeStandardHDDLRS),
		string(DataDiskTypeStandardSSDLRS),
		string(DataDiskTypeStandardSSDZRS),
	}
}

func (s *DataDiskType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataDiskType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataDiskType(input string) (*DataDiskType, error) {
	vals := map[string]DataDiskType{
		"premium_ssd_lrs":    DataDiskTypePremiumSSDLRS,
		"premium_ssd_v2_lrs": DataDiskTypePremiumSSDVTwoLRS,
		"premium_ssd_zrs":    DataDiskTypePremiumSSDZRS,
		"standard_hdd_lrs":   DataDiskTypeStandardHDDLRS,
		"standard_ssd_lrs":   DataDiskTypeStandardSSDLRS,
		"standard_ssd_zrs":   DataDiskTypeStandardSSDZRS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataDiskType(input)
	return &out, nil
}

type DbConnectionAuthenticationMode string

const (
	DbConnectionAuthenticationModeIdentityAuth DbConnectionAuthenticationMode = "IdentityAuth"
	DbConnectionAuthenticationModeSqlAuth      DbConnectionAuthenticationMode = "SqlAuth"
)

func PossibleValuesForDbConnectionAuthenticationMode() []string {
	return []string{
		string(DbConnectionAuthenticationModeIdentityAuth),
		string(DbConnectionAuthenticationModeSqlAuth),
	}
}

func (s *DbConnectionAuthenticationMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDbConnectionAuthenticationMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDbConnectionAuthenticationMode(input string) (*DbConnectionAuthenticationMode, error) {
	vals := map[string]DbConnectionAuthenticationMode{
		"identityauth": DbConnectionAuthenticationModeIdentityAuth,
		"sqlauth":      DbConnectionAuthenticationModeSqlAuth,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DbConnectionAuthenticationMode(input)
	return &out, nil
}

type DeploymentMode string

const (
	DeploymentModeApplication DeploymentMode = "Application"
	DeploymentModeSession     DeploymentMode = "Session"
)

func PossibleValuesForDeploymentMode() []string {
	return []string{
		string(DeploymentModeApplication),
		string(DeploymentModeSession),
	}
}

func (s *DeploymentMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeploymentMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeploymentMode(input string) (*DeploymentMode, error) {
	vals := map[string]DeploymentMode{
		"application": DeploymentModeApplication,
		"session":     DeploymentModeSession,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeploymentMode(input)
	return &out, nil
}

type JobType string

const (
	JobTypeFlinkJob JobType = "FlinkJob"
)

func PossibleValuesForJobType() []string {
	return []string{
		string(JobTypeFlinkJob),
	}
}

func (s *JobType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJobType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJobType(input string) (*JobType, error) {
	vals := map[string]JobType{
		"flinkjob": JobTypeFlinkJob,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobType(input)
	return &out, nil
}

type KeyVaultObjectType string

const (
	KeyVaultObjectTypeCertificate KeyVaultObjectType = "Certificate"
	KeyVaultObjectTypeKey         KeyVaultObjectType = "Key"
	KeyVaultObjectTypeSecret      KeyVaultObjectType = "Secret"
)

func PossibleValuesForKeyVaultObjectType() []string {
	return []string{
		string(KeyVaultObjectTypeCertificate),
		string(KeyVaultObjectTypeKey),
		string(KeyVaultObjectTypeSecret),
	}
}

func (s *KeyVaultObjectType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeyVaultObjectType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeyVaultObjectType(input string) (*KeyVaultObjectType, error) {
	vals := map[string]KeyVaultObjectType{
		"certificate": KeyVaultObjectTypeCertificate,
		"key":         KeyVaultObjectTypeKey,
		"secret":      KeyVaultObjectTypeSecret,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyVaultObjectType(input)
	return &out, nil
}

type LibraryManagementAction string

const (
	LibraryManagementActionInstall   LibraryManagementAction = "Install"
	LibraryManagementActionUninstall LibraryManagementAction = "Uninstall"
)

func PossibleValuesForLibraryManagementAction() []string {
	return []string{
		string(LibraryManagementActionInstall),
		string(LibraryManagementActionUninstall),
	}
}

func (s *LibraryManagementAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLibraryManagementAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLibraryManagementAction(input string) (*LibraryManagementAction, error) {
	vals := map[string]LibraryManagementAction{
		"install":   LibraryManagementActionInstall,
		"uninstall": LibraryManagementActionUninstall,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LibraryManagementAction(input)
	return &out, nil
}

type ManagedIdentityType string

const (
	ManagedIdentityTypeCluster  ManagedIdentityType = "cluster"
	ManagedIdentityTypeInternal ManagedIdentityType = "internal"
	ManagedIdentityTypeUser     ManagedIdentityType = "user"
)

func PossibleValuesForManagedIdentityType() []string {
	return []string{
		string(ManagedIdentityTypeCluster),
		string(ManagedIdentityTypeInternal),
		string(ManagedIdentityTypeUser),
	}
}

func (s *ManagedIdentityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedIdentityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedIdentityType(input string) (*ManagedIdentityType, error) {
	vals := map[string]ManagedIdentityType{
		"cluster":  ManagedIdentityTypeCluster,
		"internal": ManagedIdentityTypeInternal,
		"user":     ManagedIdentityTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedIdentityType(input)
	return &out, nil
}

type MetastoreDbConnectionAuthenticationMode string

const (
	MetastoreDbConnectionAuthenticationModeIdentityAuth MetastoreDbConnectionAuthenticationMode = "IdentityAuth"
	MetastoreDbConnectionAuthenticationModeSqlAuth      MetastoreDbConnectionAuthenticationMode = "SqlAuth"
)

func PossibleValuesForMetastoreDbConnectionAuthenticationMode() []string {
	return []string{
		string(MetastoreDbConnectionAuthenticationModeIdentityAuth),
		string(MetastoreDbConnectionAuthenticationModeSqlAuth),
	}
}

func (s *MetastoreDbConnectionAuthenticationMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMetastoreDbConnectionAuthenticationMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMetastoreDbConnectionAuthenticationMode(input string) (*MetastoreDbConnectionAuthenticationMode, error) {
	vals := map[string]MetastoreDbConnectionAuthenticationMode{
		"identityauth": MetastoreDbConnectionAuthenticationModeIdentityAuth,
		"sqlauth":      MetastoreDbConnectionAuthenticationModeSqlAuth,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MetastoreDbConnectionAuthenticationMode(input)
	return &out, nil
}

type OutboundType string

const (
	OutboundTypeLoadBalancer       OutboundType = "loadBalancer"
	OutboundTypeUserDefinedRouting OutboundType = "userDefinedRouting"
)

func PossibleValuesForOutboundType() []string {
	return []string{
		string(OutboundTypeLoadBalancer),
		string(OutboundTypeUserDefinedRouting),
	}
}

func (s *OutboundType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOutboundType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOutboundType(input string) (*OutboundType, error) {
	vals := map[string]OutboundType{
		"loadbalancer":       OutboundTypeLoadBalancer,
		"userdefinedrouting": OutboundTypeUserDefinedRouting,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OutboundType(input)
	return &out, nil
}

type ProvisioningStatus string

const (
	ProvisioningStatusAccepted  ProvisioningStatus = "Accepted"
	ProvisioningStatusCanceled  ProvisioningStatus = "Canceled"
	ProvisioningStatusFailed    ProvisioningStatus = "Failed"
	ProvisioningStatusSucceeded ProvisioningStatus = "Succeeded"
)

func PossibleValuesForProvisioningStatus() []string {
	return []string{
		string(ProvisioningStatusAccepted),
		string(ProvisioningStatusCanceled),
		string(ProvisioningStatusFailed),
		string(ProvisioningStatusSucceeded),
	}
}

func (s *ProvisioningStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningStatus(input string) (*ProvisioningStatus, error) {
	vals := map[string]ProvisioningStatus{
		"accepted":  ProvisioningStatusAccepted,
		"canceled":  ProvisioningStatusCanceled,
		"failed":    ProvisioningStatusFailed,
		"succeeded": ProvisioningStatusSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningStatus(input)
	return &out, nil
}

type RangerUsersyncMode string

const (
	RangerUsersyncModeAutomatic RangerUsersyncMode = "automatic"
	RangerUsersyncModeStatic    RangerUsersyncMode = "static"
)

func PossibleValuesForRangerUsersyncMode() []string {
	return []string{
		string(RangerUsersyncModeAutomatic),
		string(RangerUsersyncModeStatic),
	}
}

func (s *RangerUsersyncMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRangerUsersyncMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRangerUsersyncMode(input string) (*RangerUsersyncMode, error) {
	vals := map[string]RangerUsersyncMode{
		"automatic": RangerUsersyncModeAutomatic,
		"static":    RangerUsersyncModeStatic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RangerUsersyncMode(input)
	return &out, nil
}

type ScaleActionType string

const (
	ScaleActionTypeScaledown ScaleActionType = "scaledown"
	ScaleActionTypeScaleup   ScaleActionType = "scaleup"
)

func PossibleValuesForScaleActionType() []string {
	return []string{
		string(ScaleActionTypeScaledown),
		string(ScaleActionTypeScaleup),
	}
}

func (s *ScaleActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScaleActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScaleActionType(input string) (*ScaleActionType, error) {
	vals := map[string]ScaleActionType{
		"scaledown": ScaleActionTypeScaledown,
		"scaleup":   ScaleActionTypeScaleup,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScaleActionType(input)
	return &out, nil
}

type ScheduleDay string

const (
	ScheduleDayFriday    ScheduleDay = "Friday"
	ScheduleDayMonday    ScheduleDay = "Monday"
	ScheduleDaySaturday  ScheduleDay = "Saturday"
	ScheduleDaySunday    ScheduleDay = "Sunday"
	ScheduleDayThursday  ScheduleDay = "Thursday"
	ScheduleDayTuesday   ScheduleDay = "Tuesday"
	ScheduleDayWednesday ScheduleDay = "Wednesday"
)

func PossibleValuesForScheduleDay() []string {
	return []string{
		string(ScheduleDayFriday),
		string(ScheduleDayMonday),
		string(ScheduleDaySaturday),
		string(ScheduleDaySunday),
		string(ScheduleDayThursday),
		string(ScheduleDayTuesday),
		string(ScheduleDayWednesday),
	}
}

func (s *ScheduleDay) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScheduleDay(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScheduleDay(input string) (*ScheduleDay, error) {
	vals := map[string]ScheduleDay{
		"friday":    ScheduleDayFriday,
		"monday":    ScheduleDayMonday,
		"saturday":  ScheduleDaySaturday,
		"sunday":    ScheduleDaySunday,
		"thursday":  ScheduleDayThursday,
		"tuesday":   ScheduleDayTuesday,
		"wednesday": ScheduleDayWednesday,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScheduleDay(input)
	return &out, nil
}

type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
)

func PossibleValuesForSeverity() []string {
	return []string{
		string(SeverityCritical),
		string(SeverityHigh),
		string(SeverityLow),
		string(SeverityMedium),
	}
}

func (s *Severity) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSeverity(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSeverity(input string) (*Severity, error) {
	vals := map[string]Severity{
		"critical": SeverityCritical,
		"high":     SeverityHigh,
		"low":      SeverityLow,
		"medium":   SeverityMedium,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Severity(input)
	return &out, nil
}

type Status string

const (
	StatusINSTALLED       Status = "INSTALLED"
	StatusINSTALLFAILED   Status = "INSTALL_FAILED"
	StatusINSTALLING      Status = "INSTALLING"
	StatusUNINSTALLFAILED Status = "UNINSTALL_FAILED"
	StatusUNINSTALLING    Status = "UNINSTALLING"
)

func PossibleValuesForStatus() []string {
	return []string{
		string(StatusINSTALLED),
		string(StatusINSTALLFAILED),
		string(StatusINSTALLING),
		string(StatusUNINSTALLFAILED),
		string(StatusUNINSTALLING),
	}
}

func (s *Status) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStatus(input string) (*Status, error) {
	vals := map[string]Status{
		"installed":        StatusINSTALLED,
		"install_failed":   StatusINSTALLFAILED,
		"installing":       StatusINSTALLING,
		"uninstall_failed": StatusUNINSTALLFAILED,
		"uninstalling":     StatusUNINSTALLING,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Status(input)
	return &out, nil
}

type Type string

const (
	TypeMaven Type = "maven"
	TypePypi  Type = "pypi"
)

func PossibleValuesForType() []string {
	return []string{
		string(TypeMaven),
		string(TypePypi),
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
		"maven": TypeMaven,
		"pypi":  TypePypi,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Type(input)
	return &out, nil
}

type UpgradeMode string

const (
	UpgradeModeLASTSTATEUPDATE UpgradeMode = "LAST_STATE_UPDATE"
	UpgradeModeSTATELESSUPDATE UpgradeMode = "STATELESS_UPDATE"
	UpgradeModeUPDATE          UpgradeMode = "UPDATE"
)

func PossibleValuesForUpgradeMode() []string {
	return []string{
		string(UpgradeModeLASTSTATEUPDATE),
		string(UpgradeModeSTATELESSUPDATE),
		string(UpgradeModeUPDATE),
	}
}

func (s *UpgradeMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUpgradeMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUpgradeMode(input string) (*UpgradeMode, error) {
	vals := map[string]UpgradeMode{
		"last_state_update": UpgradeModeLASTSTATEUPDATE,
		"stateless_update":  UpgradeModeSTATELESSUPDATE,
		"update":            UpgradeModeUPDATE,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UpgradeMode(input)
	return &out, nil
}
