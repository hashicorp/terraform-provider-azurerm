package user

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppType string

const (
	AppTypeDeveloperPortal AppType = "developerPortal"
	AppTypePortal          AppType = "portal"
)

func PossibleValuesForAppType() []string {
	return []string{
		string(AppTypeDeveloperPortal),
		string(AppTypePortal),
	}
}

func (s *AppType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAppType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAppType(input string) (*AppType, error) {
	vals := map[string]AppType{
		"developerportal": AppTypeDeveloperPortal,
		"portal":          AppTypePortal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AppType(input)
	return &out, nil
}

type Confirmation string

const (
	ConfirmationInvite Confirmation = "invite"
	ConfirmationSignup Confirmation = "signup"
)

func PossibleValuesForConfirmation() []string {
	return []string{
		string(ConfirmationInvite),
		string(ConfirmationSignup),
	}
}

func (s *Confirmation) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConfirmation(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConfirmation(input string) (*Confirmation, error) {
	vals := map[string]Confirmation{
		"invite": ConfirmationInvite,
		"signup": ConfirmationSignup,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Confirmation(input)
	return &out, nil
}

type GroupType string

const (
	GroupTypeCustom   GroupType = "custom"
	GroupTypeExternal GroupType = "external"
	GroupTypeSystem   GroupType = "system"
)

func PossibleValuesForGroupType() []string {
	return []string{
		string(GroupTypeCustom),
		string(GroupTypeExternal),
		string(GroupTypeSystem),
	}
}

func (s *GroupType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGroupType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGroupType(input string) (*GroupType, error) {
	vals := map[string]GroupType{
		"custom":   GroupTypeCustom,
		"external": GroupTypeExternal,
		"system":   GroupTypeSystem,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GroupType(input)
	return &out, nil
}

type UserState string

const (
	UserStateActive  UserState = "active"
	UserStateBlocked UserState = "blocked"
	UserStateDeleted UserState = "deleted"
	UserStatePending UserState = "pending"
)

func PossibleValuesForUserState() []string {
	return []string{
		string(UserStateActive),
		string(UserStateBlocked),
		string(UserStateDeleted),
		string(UserStatePending),
	}
}

func (s *UserState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUserState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUserState(input string) (*UserState, error) {
	vals := map[string]UserState{
		"active":  UserStateActive,
		"blocked": UserStateBlocked,
		"deleted": UserStateDeleted,
		"pending": UserStatePending,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UserState(input)
	return &out, nil
}
