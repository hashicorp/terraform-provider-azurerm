package dataflows

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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
