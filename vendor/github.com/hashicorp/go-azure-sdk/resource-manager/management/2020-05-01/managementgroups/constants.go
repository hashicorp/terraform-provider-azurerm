package managementgroups

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Expand string

const (
	ExpandChildren Expand = "children"
	ExpandPath     Expand = "path"
)

func PossibleValuesForExpand() []string {
	return []string{
		string(ExpandChildren),
		string(ExpandPath),
	}
}

func (s *Expand) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExpand(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExpand(input string) (*Expand, error) {
	vals := map[string]Expand{
		"children": ExpandChildren,
		"path":     ExpandPath,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Expand(input)
	return &out, nil
}

type ManagementGroupChildType string

const (
	ManagementGroupChildTypeMicrosoftPointManagementManagementGroups ManagementGroupChildType = "Microsoft.Management/managementGroups"
	ManagementGroupChildTypeSubscriptions                            ManagementGroupChildType = "/subscriptions"
)

func PossibleValuesForManagementGroupChildType() []string {
	return []string{
		string(ManagementGroupChildTypeMicrosoftPointManagementManagementGroups),
		string(ManagementGroupChildTypeSubscriptions),
	}
}

func (s *ManagementGroupChildType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagementGroupChildType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagementGroupChildType(input string) (*ManagementGroupChildType, error) {
	vals := map[string]ManagementGroupChildType{
		"microsoft.management/managementgroups": ManagementGroupChildTypeMicrosoftPointManagementManagementGroups,
		"/subscriptions":                        ManagementGroupChildTypeSubscriptions,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagementGroupChildType(input)
	return &out, nil
}
