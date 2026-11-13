package apiversions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LifecycleStage string

const (
	LifecycleStageDeprecated  LifecycleStage = "deprecated"
	LifecycleStageDesign      LifecycleStage = "design"
	LifecycleStageDevelopment LifecycleStage = "development"
	LifecycleStagePreview     LifecycleStage = "preview"
	LifecycleStageProduction  LifecycleStage = "production"
	LifecycleStageRetired     LifecycleStage = "retired"
	LifecycleStageTesting     LifecycleStage = "testing"
)

func PossibleValuesForLifecycleStage() []string {
	return []string{
		string(LifecycleStageDeprecated),
		string(LifecycleStageDesign),
		string(LifecycleStageDevelopment),
		string(LifecycleStagePreview),
		string(LifecycleStageProduction),
		string(LifecycleStageRetired),
		string(LifecycleStageTesting),
	}
}

func (s *LifecycleStage) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLifecycleStage(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLifecycleStage(input string) (*LifecycleStage, error) {
	vals := map[string]LifecycleStage{
		"deprecated":  LifecycleStageDeprecated,
		"design":      LifecycleStageDesign,
		"development": LifecycleStageDevelopment,
		"preview":     LifecycleStagePreview,
		"production":  LifecycleStageProduction,
		"retired":     LifecycleStageRetired,
		"testing":     LifecycleStageTesting,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LifecycleStage(input)
	return &out, nil
}
