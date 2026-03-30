package restorables

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiType string

const (
	ApiTypeCassandra   ApiType = "Cassandra"
	ApiTypeGremlin     ApiType = "Gremlin"
	ApiTypeGremlinVTwo ApiType = "GremlinV2"
	ApiTypeMongoDB     ApiType = "MongoDB"
	ApiTypeSql         ApiType = "Sql"
	ApiTypeTable       ApiType = "Table"
)

func PossibleValuesForApiType() []string {
	return []string{
		string(ApiTypeCassandra),
		string(ApiTypeGremlin),
		string(ApiTypeGremlinVTwo),
		string(ApiTypeMongoDB),
		string(ApiTypeSql),
		string(ApiTypeTable),
	}
}

func (s *ApiType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApiType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApiType(input string) (*ApiType, error) {
	vals := map[string]ApiType{
		"cassandra": ApiTypeCassandra,
		"gremlin":   ApiTypeGremlin,
		"gremlinv2": ApiTypeGremlinVTwo,
		"mongodb":   ApiTypeMongoDB,
		"sql":       ApiTypeSql,
		"table":     ApiTypeTable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApiType(input)
	return &out, nil
}

type CompositePathSortOrder string

const (
	CompositePathSortOrderAscending  CompositePathSortOrder = "ascending"
	CompositePathSortOrderDescending CompositePathSortOrder = "descending"
)

func PossibleValuesForCompositePathSortOrder() []string {
	return []string{
		string(CompositePathSortOrderAscending),
		string(CompositePathSortOrderDescending),
	}
}

func (s *CompositePathSortOrder) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCompositePathSortOrder(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCompositePathSortOrder(input string) (*CompositePathSortOrder, error) {
	vals := map[string]CompositePathSortOrder{
		"ascending":  CompositePathSortOrderAscending,
		"descending": CompositePathSortOrderDescending,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CompositePathSortOrder(input)
	return &out, nil
}

type ConflictResolutionMode string

const (
	ConflictResolutionModeCustom         ConflictResolutionMode = "Custom"
	ConflictResolutionModeLastWriterWins ConflictResolutionMode = "LastWriterWins"
)

func PossibleValuesForConflictResolutionMode() []string {
	return []string{
		string(ConflictResolutionModeCustom),
		string(ConflictResolutionModeLastWriterWins),
	}
}

func (s *ConflictResolutionMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConflictResolutionMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConflictResolutionMode(input string) (*ConflictResolutionMode, error) {
	vals := map[string]ConflictResolutionMode{
		"custom":         ConflictResolutionModeCustom,
		"lastwriterwins": ConflictResolutionModeLastWriterWins,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConflictResolutionMode(input)
	return &out, nil
}

type CreateMode string

const (
	CreateModeDefault CreateMode = "Default"
	CreateModeRestore CreateMode = "Restore"
)

func PossibleValuesForCreateMode() []string {
	return []string{
		string(CreateModeDefault),
		string(CreateModeRestore),
	}
}

func (s *CreateMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCreateMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCreateMode(input string) (*CreateMode, error) {
	vals := map[string]CreateMode{
		"default": CreateModeDefault,
		"restore": CreateModeRestore,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CreateMode(input)
	return &out, nil
}

type DataType string

const (
	DataTypeLineString   DataType = "LineString"
	DataTypeMultiPolygon DataType = "MultiPolygon"
	DataTypeNumber       DataType = "Number"
	DataTypePoint        DataType = "Point"
	DataTypePolygon      DataType = "Polygon"
	DataTypeString       DataType = "String"
)

func PossibleValuesForDataType() []string {
	return []string{
		string(DataTypeLineString),
		string(DataTypeMultiPolygon),
		string(DataTypeNumber),
		string(DataTypePoint),
		string(DataTypePolygon),
		string(DataTypeString),
	}
}

func (s *DataType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataType(input string) (*DataType, error) {
	vals := map[string]DataType{
		"linestring":   DataTypeLineString,
		"multipolygon": DataTypeMultiPolygon,
		"number":       DataTypeNumber,
		"point":        DataTypePoint,
		"polygon":      DataTypePolygon,
		"string":       DataTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataType(input)
	return &out, nil
}

type IndexKind string

const (
	IndexKindHash    IndexKind = "Hash"
	IndexKindRange   IndexKind = "Range"
	IndexKindSpatial IndexKind = "Spatial"
)

func PossibleValuesForIndexKind() []string {
	return []string{
		string(IndexKindHash),
		string(IndexKindRange),
		string(IndexKindSpatial),
	}
}

func (s *IndexKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIndexKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIndexKind(input string) (*IndexKind, error) {
	vals := map[string]IndexKind{
		"hash":    IndexKindHash,
		"range":   IndexKindRange,
		"spatial": IndexKindSpatial,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IndexKind(input)
	return &out, nil
}

type IndexingMode string

const (
	IndexingModeConsistent IndexingMode = "consistent"
	IndexingModeLazy       IndexingMode = "lazy"
	IndexingModeNone       IndexingMode = "none"
)

func PossibleValuesForIndexingMode() []string {
	return []string{
		string(IndexingModeConsistent),
		string(IndexingModeLazy),
		string(IndexingModeNone),
	}
}

func (s *IndexingMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIndexingMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIndexingMode(input string) (*IndexingMode, error) {
	vals := map[string]IndexingMode{
		"consistent": IndexingModeConsistent,
		"lazy":       IndexingModeLazy,
		"none":       IndexingModeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IndexingMode(input)
	return &out, nil
}

type OperationType string

const (
	OperationTypeCreate          OperationType = "Create"
	OperationTypeDelete          OperationType = "Delete"
	OperationTypeRecreate        OperationType = "Recreate"
	OperationTypeReplace         OperationType = "Replace"
	OperationTypeSystemOperation OperationType = "SystemOperation"
)

func PossibleValuesForOperationType() []string {
	return []string{
		string(OperationTypeCreate),
		string(OperationTypeDelete),
		string(OperationTypeRecreate),
		string(OperationTypeReplace),
		string(OperationTypeSystemOperation),
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
		"create":          OperationTypeCreate,
		"delete":          OperationTypeDelete,
		"recreate":        OperationTypeRecreate,
		"replace":         OperationTypeReplace,
		"systemoperation": OperationTypeSystemOperation,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperationType(input)
	return &out, nil
}

type PartitionKind string

const (
	PartitionKindHash      PartitionKind = "Hash"
	PartitionKindMultiHash PartitionKind = "MultiHash"
	PartitionKindRange     PartitionKind = "Range"
)

func PossibleValuesForPartitionKind() []string {
	return []string{
		string(PartitionKindHash),
		string(PartitionKindMultiHash),
		string(PartitionKindRange),
	}
}

func (s *PartitionKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePartitionKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePartitionKind(input string) (*PartitionKind, error) {
	vals := map[string]PartitionKind{
		"hash":      PartitionKindHash,
		"multihash": PartitionKindMultiHash,
		"range":     PartitionKindRange,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PartitionKind(input)
	return &out, nil
}

type SpatialType string

const (
	SpatialTypeLineString   SpatialType = "LineString"
	SpatialTypeMultiPolygon SpatialType = "MultiPolygon"
	SpatialTypePoint        SpatialType = "Point"
	SpatialTypePolygon      SpatialType = "Polygon"
)

func PossibleValuesForSpatialType() []string {
	return []string{
		string(SpatialTypeLineString),
		string(SpatialTypeMultiPolygon),
		string(SpatialTypePoint),
		string(SpatialTypePolygon),
	}
}

func (s *SpatialType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSpatialType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSpatialType(input string) (*SpatialType, error) {
	vals := map[string]SpatialType{
		"linestring":   SpatialTypeLineString,
		"multipolygon": SpatialTypeMultiPolygon,
		"point":        SpatialTypePoint,
		"polygon":      SpatialTypePolygon,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SpatialType(input)
	return &out, nil
}
