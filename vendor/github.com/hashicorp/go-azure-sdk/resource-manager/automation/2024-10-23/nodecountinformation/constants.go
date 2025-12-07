package nodecountinformation

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CountType string

const (
	CountTypeNodeconfiguration CountType = "nodeconfiguration"
	CountTypeStatus            CountType = "status"
)

func PossibleValuesForCountType() []string {
	return []string{
		string(CountTypeNodeconfiguration),
		string(CountTypeStatus),
	}
}

func (s *CountType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCountType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCountType(input string) (*CountType, error) {
	vals := map[string]CountType{
		"nodeconfiguration": CountTypeNodeconfiguration,
		"status":            CountTypeStatus,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CountType(input)
	return &out, nil
}
