package flexcomponents

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HardwareType string

const (
	HardwareTypeCELL    HardwareType = "CELL"
	HardwareTypeCOMPUTE HardwareType = "COMPUTE"
)

func PossibleValuesForHardwareType() []string {
	return []string{
		string(HardwareTypeCELL),
		string(HardwareTypeCOMPUTE),
	}
}

func (s *HardwareType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHardwareType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHardwareType(input string) (*HardwareType, error) {
	vals := map[string]HardwareType{
		"cell":    HardwareTypeCELL,
		"compute": HardwareTypeCOMPUTE,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HardwareType(input)
	return &out, nil
}

type SystemShapes string

const (
	SystemShapesExaDbXS              SystemShapes = "ExaDbXS"
	SystemShapesExadataPointXNineM   SystemShapes = "Exadata.X9M"
	SystemShapesExadataPointXOneOneM SystemShapes = "Exadata.X11M"
)

func PossibleValuesForSystemShapes() []string {
	return []string{
		string(SystemShapesExaDbXS),
		string(SystemShapesExadataPointXNineM),
		string(SystemShapesExadataPointXOneOneM),
	}
}

func (s *SystemShapes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSystemShapes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSystemShapes(input string) (*SystemShapes, error) {
	vals := map[string]SystemShapes{
		"exadbxs":      SystemShapesExaDbXS,
		"exadata.x9m":  SystemShapesExadataPointXNineM,
		"exadata.x11m": SystemShapesExadataPointXOneOneM,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SystemShapes(input)
	return &out, nil
}
