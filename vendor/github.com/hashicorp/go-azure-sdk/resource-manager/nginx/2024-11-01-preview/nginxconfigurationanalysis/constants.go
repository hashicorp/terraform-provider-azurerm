package nginxconfigurationanalysis

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Level string

const (
	LevelInfo    Level = "Info"
	LevelWarning Level = "Warning"
)

func PossibleValuesForLevel() []string {
	return []string{
		string(LevelInfo),
		string(LevelWarning),
	}
}

func (s *Level) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLevel(input string) (*Level, error) {
	vals := map[string]Level{
		"info":    LevelInfo,
		"warning": LevelWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Level(input)
	return &out, nil
}
