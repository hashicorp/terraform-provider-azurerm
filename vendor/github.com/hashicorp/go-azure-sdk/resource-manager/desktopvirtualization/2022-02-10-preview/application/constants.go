package application

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CommandLineSetting string

const (
	CommandLineSettingAllow      CommandLineSetting = "Allow"
	CommandLineSettingDoNotAllow CommandLineSetting = "DoNotAllow"
	CommandLineSettingRequire    CommandLineSetting = "Require"
)

func PossibleValuesForCommandLineSetting() []string {
	return []string{
		string(CommandLineSettingAllow),
		string(CommandLineSettingDoNotAllow),
		string(CommandLineSettingRequire),
	}
}

func (s *CommandLineSetting) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCommandLineSetting(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCommandLineSetting(input string) (*CommandLineSetting, error) {
	vals := map[string]CommandLineSetting{
		"allow":      CommandLineSettingAllow,
		"donotallow": CommandLineSettingDoNotAllow,
		"require":    CommandLineSettingRequire,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CommandLineSetting(input)
	return &out, nil
}

type RemoteApplicationType string

const (
	RemoteApplicationTypeInBuilt         RemoteApplicationType = "InBuilt"
	RemoteApplicationTypeMsixApplication RemoteApplicationType = "MsixApplication"
)

func PossibleValuesForRemoteApplicationType() []string {
	return []string{
		string(RemoteApplicationTypeInBuilt),
		string(RemoteApplicationTypeMsixApplication),
	}
}

func (s *RemoteApplicationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRemoteApplicationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRemoteApplicationType(input string) (*RemoteApplicationType, error) {
	vals := map[string]RemoteApplicationType{
		"inbuilt":         RemoteApplicationTypeInBuilt,
		"msixapplication": RemoteApplicationTypeMsixApplication,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RemoteApplicationType(input)
	return &out, nil
}
