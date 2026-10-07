// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sdk

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
	*s = parseResourceValidationType(decoded)
	return nil
}

func parseResourceValidationType(input string) ResourceValidationType {
	vals := map[string]ResourceValidationType{
		"armfull":    ResourceValidationTypeArmFull,
		"armpartial": ResourceValidationTypeArmPartial,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return v
	}

	// otherwise presume it's an undefined value and best-effort it
	return ResourceValidationType(input)
}
