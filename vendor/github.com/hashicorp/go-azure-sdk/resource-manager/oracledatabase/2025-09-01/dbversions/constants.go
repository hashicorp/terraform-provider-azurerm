package dbversions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BaseDbSystemShapes string

const (
	BaseDbSystemShapesVMPointStandardPointxEightSix BaseDbSystemShapes = "VM.Standard.x86"
)

func PossibleValuesForBaseDbSystemShapes() []string {
	return []string{
		string(BaseDbSystemShapesVMPointStandardPointxEightSix),
	}
}

func (s *BaseDbSystemShapes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBaseDbSystemShapes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBaseDbSystemShapes(input string) (*BaseDbSystemShapes, error) {
	vals := map[string]BaseDbSystemShapes{
		"vm.standard.x86": BaseDbSystemShapesVMPointStandardPointxEightSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BaseDbSystemShapes(input)
	return &out, nil
}

type ShapeFamilyType string

const (
	ShapeFamilyTypeEXADATA        ShapeFamilyType = "EXADATA"
	ShapeFamilyTypeEXADBXS        ShapeFamilyType = "EXADB_XS"
	ShapeFamilyTypeSINGLENODE     ShapeFamilyType = "SINGLENODE"
	ShapeFamilyTypeVirtualMachine ShapeFamilyType = "VIRTUALMACHINE"
)

func PossibleValuesForShapeFamilyType() []string {
	return []string{
		string(ShapeFamilyTypeEXADATA),
		string(ShapeFamilyTypeEXADBXS),
		string(ShapeFamilyTypeSINGLENODE),
		string(ShapeFamilyTypeVirtualMachine),
	}
}

func (s *ShapeFamilyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseShapeFamilyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseShapeFamilyType(input string) (*ShapeFamilyType, error) {
	vals := map[string]ShapeFamilyType{
		"exadata":        ShapeFamilyTypeEXADATA,
		"exadb_xs":       ShapeFamilyTypeEXADBXS,
		"singlenode":     ShapeFamilyTypeSINGLENODE,
		"virtualmachine": ShapeFamilyTypeVirtualMachine,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ShapeFamilyType(input)
	return &out, nil
}

type StorageManagementType string

const (
	StorageManagementTypeLVM StorageManagementType = "LVM"
)

func PossibleValuesForStorageManagementType() []string {
	return []string{
		string(StorageManagementTypeLVM),
	}
}

func (s *StorageManagementType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageManagementType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageManagementType(input string) (*StorageManagementType, error) {
	vals := map[string]StorageManagementType{
		"lvm": StorageManagementTypeLVM,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageManagementType(input)
	return &out, nil
}
