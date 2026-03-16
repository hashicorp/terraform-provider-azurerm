package resourcevalidationclient

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceValidationType string

const (
	ResourceValidationTypeArmFull    ResourceValidationType = "ArmFull"
	ResourceValidationTypeArmPartial ResourceValidationType = "ArmPartial"
)

func PossibleValuesForResourceValidationType() []string {
	return []string{
		string(ResourceValidationTypeArmFull),
		string(ResourceValidationTypeArmPartial),
	}
}

func (s *ResourceValidationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceValidationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceValidationType(input string) (*ResourceValidationType, error) {
	vals := map[string]ResourceValidationType{
		"armfull":    ResourceValidationTypeArmFull,
		"armpartial": ResourceValidationTypeArmPartial,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceValidationType(input)
	return &out, nil
}
