package synonymmaps

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Prefer string

const (
	PreferReturnRepresentation Prefer = "return=representation"
)

func PossibleValuesForPrefer() []string {
	return []string{
		string(PreferReturnRepresentation),
	}
}

func (s *Prefer) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrefer(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrefer(input string) (*Prefer, error) {
	vals := map[string]Prefer{
		"return=representation": PreferReturnRepresentation,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Prefer(input)
	return &out, nil
}

type SynonymMapFormat string

const (
	SynonymMapFormatSolr SynonymMapFormat = "solr"
)

func PossibleValuesForSynonymMapFormat() []string {
	return []string{
		string(SynonymMapFormatSolr),
	}
}

func (s *SynonymMapFormat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSynonymMapFormat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSynonymMapFormat(input string) (*SynonymMapFormat, error) {
	vals := map[string]SynonymMapFormat{
		"solr": SynonymMapFormatSolr,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SynonymMapFormat(input)
	return &out, nil
}
