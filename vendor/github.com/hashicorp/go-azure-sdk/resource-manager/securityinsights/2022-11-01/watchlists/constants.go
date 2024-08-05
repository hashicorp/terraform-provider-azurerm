package watchlists

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Source string

const (
	SourceLocalFile     Source = "Local file"
	SourceRemoteStorage Source = "Remote storage"
)

func PossibleValuesForSource() []string {
	return []string{
		string(SourceLocalFile),
		string(SourceRemoteStorage),
	}
}

func (s *Source) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSource(input string) (*Source, error) {
	vals := map[string]Source{
		"local file":     SourceLocalFile,
		"remote storage": SourceRemoteStorage,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Source(input)
	return &out, nil
}
