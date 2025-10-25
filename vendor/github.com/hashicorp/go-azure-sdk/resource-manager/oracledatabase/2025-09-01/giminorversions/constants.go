package giminorversions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ShapeFamily string

const (
	ShapeFamilyEXADATA ShapeFamily = "EXADATA"
	ShapeFamilyEXADBXS ShapeFamily = "EXADB_XS"
)

func PossibleValuesForShapeFamily() []string {
	return []string{
		string(ShapeFamilyEXADATA),
		string(ShapeFamilyEXADBXS),
	}
}

func (s *ShapeFamily) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseShapeFamily(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseShapeFamily(input string) (*ShapeFamily, error) {
	vals := map[string]ShapeFamily{
		"exadata":  ShapeFamilyEXADATA,
		"exadb_xs": ShapeFamilyEXADBXS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ShapeFamily(input)
	return &out, nil
}
