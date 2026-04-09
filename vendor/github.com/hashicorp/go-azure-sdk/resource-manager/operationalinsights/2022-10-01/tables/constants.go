package tables

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ColumnDataTypeHintEnum string

const (
	ColumnDataTypeHintEnumArmPath ColumnDataTypeHintEnum = "armPath"
	ColumnDataTypeHintEnumGuid    ColumnDataTypeHintEnum = "guid"
	ColumnDataTypeHintEnumIP      ColumnDataTypeHintEnum = "ip"
	ColumnDataTypeHintEnumUri     ColumnDataTypeHintEnum = "uri"
)

func PossibleValuesForColumnDataTypeHintEnum() []string {
	return []string{
		string(ColumnDataTypeHintEnumArmPath),
		string(ColumnDataTypeHintEnumGuid),
		string(ColumnDataTypeHintEnumIP),
		string(ColumnDataTypeHintEnumUri),
	}
}

func (s *ColumnDataTypeHintEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseColumnDataTypeHintEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseColumnDataTypeHintEnum(input string) (*ColumnDataTypeHintEnum, error) {
	vals := map[string]ColumnDataTypeHintEnum{
		"armpath": ColumnDataTypeHintEnumArmPath,
		"guid":    ColumnDataTypeHintEnumGuid,
		"ip":      ColumnDataTypeHintEnumIP,
		"uri":     ColumnDataTypeHintEnumUri,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ColumnDataTypeHintEnum(input)
	return &out, nil
}

type ColumnTypeEnum string

const (
	ColumnTypeEnumBoolean  ColumnTypeEnum = "boolean"
	ColumnTypeEnumDateTime ColumnTypeEnum = "dateTime"
	ColumnTypeEnumDynamic  ColumnTypeEnum = "dynamic"
	ColumnTypeEnumGuid     ColumnTypeEnum = "guid"
	ColumnTypeEnumInt      ColumnTypeEnum = "int"
	ColumnTypeEnumLong     ColumnTypeEnum = "long"
	ColumnTypeEnumReal     ColumnTypeEnum = "real"
	ColumnTypeEnumString   ColumnTypeEnum = "string"
)

func PossibleValuesForColumnTypeEnum() []string {
	return []string{
		string(ColumnTypeEnumBoolean),
		string(ColumnTypeEnumDateTime),
		string(ColumnTypeEnumDynamic),
		string(ColumnTypeEnumGuid),
		string(ColumnTypeEnumInt),
		string(ColumnTypeEnumLong),
		string(ColumnTypeEnumReal),
		string(ColumnTypeEnumString),
	}
}

func (s *ColumnTypeEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseColumnTypeEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseColumnTypeEnum(input string) (*ColumnTypeEnum, error) {
	vals := map[string]ColumnTypeEnum{
		"boolean":  ColumnTypeEnumBoolean,
		"datetime": ColumnTypeEnumDateTime,
		"dynamic":  ColumnTypeEnumDynamic,
		"guid":     ColumnTypeEnumGuid,
		"int":      ColumnTypeEnumInt,
		"long":     ColumnTypeEnumLong,
		"real":     ColumnTypeEnumReal,
		"string":   ColumnTypeEnumString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ColumnTypeEnum(input)
	return &out, nil
}

type ProvisioningStateEnum string

const (
	ProvisioningStateEnumDeleting   ProvisioningStateEnum = "Deleting"
	ProvisioningStateEnumInProgress ProvisioningStateEnum = "InProgress"
	ProvisioningStateEnumSucceeded  ProvisioningStateEnum = "Succeeded"
	ProvisioningStateEnumUpdating   ProvisioningStateEnum = "Updating"
)

func PossibleValuesForProvisioningStateEnum() []string {
	return []string{
		string(ProvisioningStateEnumDeleting),
		string(ProvisioningStateEnumInProgress),
		string(ProvisioningStateEnumSucceeded),
		string(ProvisioningStateEnumUpdating),
	}
}

func (s *ProvisioningStateEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningStateEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningStateEnum(input string) (*ProvisioningStateEnum, error) {
	vals := map[string]ProvisioningStateEnum{
		"deleting":   ProvisioningStateEnumDeleting,
		"inprogress": ProvisioningStateEnumInProgress,
		"succeeded":  ProvisioningStateEnumSucceeded,
		"updating":   ProvisioningStateEnumUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningStateEnum(input)
	return &out, nil
}

type SourceEnum string

const (
	SourceEnumCustomer  SourceEnum = "customer"
	SourceEnumMicrosoft SourceEnum = "microsoft"
)

func PossibleValuesForSourceEnum() []string {
	return []string{
		string(SourceEnumCustomer),
		string(SourceEnumMicrosoft),
	}
}

func (s *SourceEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSourceEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSourceEnum(input string) (*SourceEnum, error) {
	vals := map[string]SourceEnum{
		"customer":  SourceEnumCustomer,
		"microsoft": SourceEnumMicrosoft,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SourceEnum(input)
	return &out, nil
}

type TablePlanEnum string

const (
	TablePlanEnumAnalytics TablePlanEnum = "Analytics"
	TablePlanEnumBasic     TablePlanEnum = "Basic"
)

func PossibleValuesForTablePlanEnum() []string {
	return []string{
		string(TablePlanEnumAnalytics),
		string(TablePlanEnumBasic),
	}
}

func (s *TablePlanEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTablePlanEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTablePlanEnum(input string) (*TablePlanEnum, error) {
	vals := map[string]TablePlanEnum{
		"analytics": TablePlanEnumAnalytics,
		"basic":     TablePlanEnumBasic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TablePlanEnum(input)
	return &out, nil
}

type TableSubTypeEnum string

const (
	TableSubTypeEnumAny                     TableSubTypeEnum = "Any"
	TableSubTypeEnumClassic                 TableSubTypeEnum = "Classic"
	TableSubTypeEnumDataCollectionRuleBased TableSubTypeEnum = "DataCollectionRuleBased"
)

func PossibleValuesForTableSubTypeEnum() []string {
	return []string{
		string(TableSubTypeEnumAny),
		string(TableSubTypeEnumClassic),
		string(TableSubTypeEnumDataCollectionRuleBased),
	}
}

func (s *TableSubTypeEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTableSubTypeEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTableSubTypeEnum(input string) (*TableSubTypeEnum, error) {
	vals := map[string]TableSubTypeEnum{
		"any":                     TableSubTypeEnumAny,
		"classic":                 TableSubTypeEnumClassic,
		"datacollectionrulebased": TableSubTypeEnumDataCollectionRuleBased,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TableSubTypeEnum(input)
	return &out, nil
}

type TableTypeEnum string

const (
	TableTypeEnumCustomLog     TableTypeEnum = "CustomLog"
	TableTypeEnumMicrosoft     TableTypeEnum = "Microsoft"
	TableTypeEnumRestoredLogs  TableTypeEnum = "RestoredLogs"
	TableTypeEnumSearchResults TableTypeEnum = "SearchResults"
)

func PossibleValuesForTableTypeEnum() []string {
	return []string{
		string(TableTypeEnumCustomLog),
		string(TableTypeEnumMicrosoft),
		string(TableTypeEnumRestoredLogs),
		string(TableTypeEnumSearchResults),
	}
}

func (s *TableTypeEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTableTypeEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTableTypeEnum(input string) (*TableTypeEnum, error) {
	vals := map[string]TableTypeEnum{
		"customlog":     TableTypeEnumCustomLog,
		"microsoft":     TableTypeEnumMicrosoft,
		"restoredlogs":  TableTypeEnumRestoredLogs,
		"searchresults": TableTypeEnumSearchResults,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TableTypeEnum(input)
	return &out, nil
}
