package sshpublickeys

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SshEncryptionTypes string

const (
	SshEncryptionTypesEdTwoFiveFiveOneNine SshEncryptionTypes = "Ed25519"
	SshEncryptionTypesRSA                  SshEncryptionTypes = "RSA"
)

func PossibleValuesForSshEncryptionTypes() []string {
	return []string{
		string(SshEncryptionTypesEdTwoFiveFiveOneNine),
		string(SshEncryptionTypesRSA),
	}
}

func (s *SshEncryptionTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSshEncryptionTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSshEncryptionTypes(input string) (*SshEncryptionTypes, error) {
	vals := map[string]SshEncryptionTypes{
		"ed25519": SshEncryptionTypesEdTwoFiveFiveOneNine,
		"rsa":     SshEncryptionTypesRSA,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SshEncryptionTypes(input)
	return &out, nil
}
