package networksecurityperimeterassociableresourcetypes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NspReadinessState string

const (
	NspReadinessStateGA         NspReadinessState = "GA"
	NspReadinessStateNotReady   NspReadinessState = "NotReady"
	NspReadinessStateOnboarding NspReadinessState = "Onboarding"
	NspReadinessStatePreview    NspReadinessState = "Preview"
)

func PossibleValuesForNspReadinessState() []string {
	return []string{
		string(NspReadinessStateGA),
		string(NspReadinessStateNotReady),
		string(NspReadinessStateOnboarding),
		string(NspReadinessStatePreview),
	}
}

func (s *NspReadinessState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNspReadinessState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNspReadinessState(input string) (*NspReadinessState, error) {
	vals := map[string]NspReadinessState{
		"ga":         NspReadinessStateGA,
		"notready":   NspReadinessStateNotReady,
		"onboarding": NspReadinessStateOnboarding,
		"preview":    NspReadinessStatePreview,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NspReadinessState(input)
	return &out, nil
}
