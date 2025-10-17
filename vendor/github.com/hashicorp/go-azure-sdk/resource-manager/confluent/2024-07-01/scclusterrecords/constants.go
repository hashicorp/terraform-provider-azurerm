package scclusterrecords

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Package string

const (
	PackageADVANCED   Package = "ADVANCED"
	PackageESSENTIALS Package = "ESSENTIALS"
)

func PossibleValuesForPackage() []string {
	return []string{
		string(PackageADVANCED),
		string(PackageESSENTIALS),
	}
}

func (s *Package) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePackage(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePackage(input string) (*Package, error) {
	vals := map[string]Package{
		"advanced":   PackageADVANCED,
		"essentials": PackageESSENTIALS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Package(input)
	return &out, nil
}
