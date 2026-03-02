package localusers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListLocalUserIncludeParam string

const (
	ListLocalUserIncludeParamNfsvThree ListLocalUserIncludeParam = "nfsv3"
)

func PossibleValuesForListLocalUserIncludeParam() []string {
	return []string{
		string(ListLocalUserIncludeParamNfsvThree),
	}
}

func (s *ListLocalUserIncludeParam) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseListLocalUserIncludeParam(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseListLocalUserIncludeParam(input string) (*ListLocalUserIncludeParam, error) {
	vals := map[string]ListLocalUserIncludeParam{
		"nfsv3": ListLocalUserIncludeParamNfsvThree,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ListLocalUserIncludeParam(input)
	return &out, nil
}
