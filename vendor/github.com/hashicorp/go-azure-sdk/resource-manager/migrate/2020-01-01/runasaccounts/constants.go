package runasaccounts

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CredentialType string

const (
	CredentialTypeHyperVFabric  CredentialType = "HyperVFabric"
	CredentialTypeLinuxGuest    CredentialType = "LinuxGuest"
	CredentialTypeLinuxServer   CredentialType = "LinuxServer"
	CredentialTypeVMwareFabric  CredentialType = "VMwareFabric"
	CredentialTypeWindowsGuest  CredentialType = "WindowsGuest"
	CredentialTypeWindowsServer CredentialType = "WindowsServer"
)

func PossibleValuesForCredentialType() []string {
	return []string{
		string(CredentialTypeHyperVFabric),
		string(CredentialTypeLinuxGuest),
		string(CredentialTypeLinuxServer),
		string(CredentialTypeVMwareFabric),
		string(CredentialTypeWindowsGuest),
		string(CredentialTypeWindowsServer),
	}
}

func (s *CredentialType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCredentialType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCredentialType(input string) (*CredentialType, error) {
	vals := map[string]CredentialType{
		"hypervfabric":  CredentialTypeHyperVFabric,
		"linuxguest":    CredentialTypeLinuxGuest,
		"linuxserver":   CredentialTypeLinuxServer,
		"vmwarefabric":  CredentialTypeVMwareFabric,
		"windowsguest":  CredentialTypeWindowsGuest,
		"windowsserver": CredentialTypeWindowsServer,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CredentialType(input)
	return &out, nil
}
