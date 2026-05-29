package analyticsitemsapis

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ItemScope string

const (
	ItemScopeShared ItemScope = "shared"
	ItemScopeUser   ItemScope = "user"
)

func PossibleValuesForItemScope() []string {
	return []string{
		string(ItemScopeShared),
		string(ItemScopeUser),
	}
}

func (s *ItemScope) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseItemScope(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseItemScope(input string) (*ItemScope, error) {
	vals := map[string]ItemScope{
		"shared": ItemScopeShared,
		"user":   ItemScopeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ItemScope(input)
	return &out, nil
}

type ItemType string

const (
	ItemTypeFunction ItemType = "function"
	ItemTypeNone     ItemType = "none"
	ItemTypeQuery    ItemType = "query"
	ItemTypeRecent   ItemType = "recent"
)

func PossibleValuesForItemType() []string {
	return []string{
		string(ItemTypeFunction),
		string(ItemTypeNone),
		string(ItemTypeQuery),
		string(ItemTypeRecent),
	}
}

func (s *ItemType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseItemType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseItemType(input string) (*ItemType, error) {
	vals := map[string]ItemType{
		"function": ItemTypeFunction,
		"none":     ItemTypeNone,
		"query":    ItemTypeQuery,
		"recent":   ItemTypeRecent,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ItemType(input)
	return &out, nil
}

type ItemTypeParameter string

const (
	ItemTypeParameterFolder   ItemTypeParameter = "folder"
	ItemTypeParameterFunction ItemTypeParameter = "function"
	ItemTypeParameterNone     ItemTypeParameter = "none"
	ItemTypeParameterQuery    ItemTypeParameter = "query"
	ItemTypeParameterRecent   ItemTypeParameter = "recent"
)

func PossibleValuesForItemTypeParameter() []string {
	return []string{
		string(ItemTypeParameterFolder),
		string(ItemTypeParameterFunction),
		string(ItemTypeParameterNone),
		string(ItemTypeParameterQuery),
		string(ItemTypeParameterRecent),
	}
}

func (s *ItemTypeParameter) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseItemTypeParameter(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseItemTypeParameter(input string) (*ItemTypeParameter, error) {
	vals := map[string]ItemTypeParameter{
		"folder":   ItemTypeParameterFolder,
		"function": ItemTypeParameterFunction,
		"none":     ItemTypeParameterNone,
		"query":    ItemTypeParameterQuery,
		"recent":   ItemTypeParameterRecent,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ItemTypeParameter(input)
	return &out, nil
}
