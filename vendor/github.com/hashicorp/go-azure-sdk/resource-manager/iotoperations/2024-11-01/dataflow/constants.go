package dataflow

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataflowMappingType string

const (
	DataflowMappingTypeBuiltInFunction DataflowMappingType = "BuiltInFunction"
	DataflowMappingTypeCompute         DataflowMappingType = "Compute"
	DataflowMappingTypeNewProperties   DataflowMappingType = "NewProperties"
	DataflowMappingTypePassThrough     DataflowMappingType = "PassThrough"
	DataflowMappingTypeRename          DataflowMappingType = "Rename"
)

func PossibleValuesForDataflowMappingType() []string {
	return []string{
		string(DataflowMappingTypeBuiltInFunction),
		string(DataflowMappingTypeCompute),
		string(DataflowMappingTypeNewProperties),
		string(DataflowMappingTypePassThrough),
		string(DataflowMappingTypeRename),
	}
}

func (s *DataflowMappingType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataflowMappingType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataflowMappingType(input string) (*DataflowMappingType, error) {
	vals := map[string]DataflowMappingType{
		"builtinfunction": DataflowMappingTypeBuiltInFunction,
		"compute":         DataflowMappingTypeCompute,
		"newproperties":   DataflowMappingTypeNewProperties,
		"passthrough":     DataflowMappingTypePassThrough,
		"rename":          DataflowMappingTypeRename,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataflowMappingType(input)
	return &out, nil
}

type ExtendedLocationType string

const (
	ExtendedLocationTypeCustomLocation ExtendedLocationType = "CustomLocation"
)

func PossibleValuesForExtendedLocationType() []string {
	return []string{
		string(ExtendedLocationTypeCustomLocation),
	}
}

func (s *ExtendedLocationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExtendedLocationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExtendedLocationType(input string) (*ExtendedLocationType, error) {
	vals := map[string]ExtendedLocationType{
		"customlocation": ExtendedLocationTypeCustomLocation,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExtendedLocationType(input)
	return &out, nil
}

type FilterType string

const (
	FilterTypeFilter FilterType = "Filter"
)

func PossibleValuesForFilterType() []string {
	return []string{
		string(FilterTypeFilter),
	}
}

func (s *FilterType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFilterType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFilterType(input string) (*FilterType, error) {
	vals := map[string]FilterType{
		"filter": FilterTypeFilter,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FilterType(input)
	return &out, nil
}

type OperationType string

const (
	OperationTypeBuiltInTransformation OperationType = "BuiltInTransformation"
	OperationTypeDestination           OperationType = "Destination"
	OperationTypeSource                OperationType = "Source"
)

func PossibleValuesForOperationType() []string {
	return []string{
		string(OperationTypeBuiltInTransformation),
		string(OperationTypeDestination),
		string(OperationTypeSource),
	}
}

func (s *OperationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperationType(input string) (*OperationType, error) {
	vals := map[string]OperationType{
		"builtintransformation": OperationTypeBuiltInTransformation,
		"destination":           OperationTypeDestination,
		"source":                OperationTypeSource,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperationType(input)
	return &out, nil
}

type OperationalMode string

const (
	OperationalModeDisabled OperationalMode = "Disabled"
	OperationalModeEnabled  OperationalMode = "Enabled"
)

func PossibleValuesForOperationalMode() []string {
	return []string{
		string(OperationalModeDisabled),
		string(OperationalModeEnabled),
	}
}

func (s *OperationalMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperationalMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperationalMode(input string) (*OperationalMode, error) {
	vals := map[string]OperationalMode{
		"disabled": OperationalModeDisabled,
		"enabled":  OperationalModeEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperationalMode(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateProvisioning),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":     ProvisioningStateAccepted,
		"canceled":     ProvisioningStateCanceled,
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"provisioning": ProvisioningStateProvisioning,
		"succeeded":    ProvisioningStateSucceeded,
		"updating":     ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type SourceSerializationFormat string

const (
	SourceSerializationFormatJson SourceSerializationFormat = "Json"
)

func PossibleValuesForSourceSerializationFormat() []string {
	return []string{
		string(SourceSerializationFormatJson),
	}
}

func (s *SourceSerializationFormat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSourceSerializationFormat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSourceSerializationFormat(input string) (*SourceSerializationFormat, error) {
	vals := map[string]SourceSerializationFormat{
		"json": SourceSerializationFormatJson,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SourceSerializationFormat(input)
	return &out, nil
}

type TransformationSerializationFormat string

const (
	TransformationSerializationFormatDelta   TransformationSerializationFormat = "Delta"
	TransformationSerializationFormatJson    TransformationSerializationFormat = "Json"
	TransformationSerializationFormatParquet TransformationSerializationFormat = "Parquet"
)

func PossibleValuesForTransformationSerializationFormat() []string {
	return []string{
		string(TransformationSerializationFormatDelta),
		string(TransformationSerializationFormatJson),
		string(TransformationSerializationFormatParquet),
	}
}

func (s *TransformationSerializationFormat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTransformationSerializationFormat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTransformationSerializationFormat(input string) (*TransformationSerializationFormat, error) {
	vals := map[string]TransformationSerializationFormat{
		"delta":   TransformationSerializationFormatDelta,
		"json":    TransformationSerializationFormatJson,
		"parquet": TransformationSerializationFormatParquet,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TransformationSerializationFormat(input)
	return &out, nil
}
