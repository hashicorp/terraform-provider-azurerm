package deletedstorage

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletionRecoveryLevel string

const (
	DeletionRecoveryLevelCustomizedRecoverable                              DeletionRecoveryLevel = "CustomizedRecoverable"
	DeletionRecoveryLevelCustomizedRecoverablePositiveProtectedSubscription DeletionRecoveryLevel = "CustomizedRecoverable+ProtectedSubscription"
	DeletionRecoveryLevelCustomizedRecoverablePositivePurgeable             DeletionRecoveryLevel = "CustomizedRecoverable+Purgeable"
	DeletionRecoveryLevelPurgeable                                          DeletionRecoveryLevel = "Purgeable"
	DeletionRecoveryLevelRecoverable                                        DeletionRecoveryLevel = "Recoverable"
	DeletionRecoveryLevelRecoverablePositiveProtectedSubscription           DeletionRecoveryLevel = "Recoverable+ProtectedSubscription"
	DeletionRecoveryLevelRecoverablePositivePurgeable                       DeletionRecoveryLevel = "Recoverable+Purgeable"
)

func PossibleValuesForDeletionRecoveryLevel() []string {
	return []string{
		string(DeletionRecoveryLevelCustomizedRecoverable),
		string(DeletionRecoveryLevelCustomizedRecoverablePositiveProtectedSubscription),
		string(DeletionRecoveryLevelCustomizedRecoverablePositivePurgeable),
		string(DeletionRecoveryLevelPurgeable),
		string(DeletionRecoveryLevelRecoverable),
		string(DeletionRecoveryLevelRecoverablePositiveProtectedSubscription),
		string(DeletionRecoveryLevelRecoverablePositivePurgeable),
	}
}

func (s *DeletionRecoveryLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeletionRecoveryLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeletionRecoveryLevel(input string) (*DeletionRecoveryLevel, error) {
	vals := map[string]DeletionRecoveryLevel{
		"customizedrecoverable":                       DeletionRecoveryLevelCustomizedRecoverable,
		"customizedrecoverable+protectedsubscription": DeletionRecoveryLevelCustomizedRecoverablePositiveProtectedSubscription,
		"customizedrecoverable+purgeable":             DeletionRecoveryLevelCustomizedRecoverablePositivePurgeable,
		"purgeable":                                   DeletionRecoveryLevelPurgeable,
		"recoverable":                                 DeletionRecoveryLevelRecoverable,
		"recoverable+protectedsubscription":           DeletionRecoveryLevelRecoverablePositiveProtectedSubscription,
		"recoverable+purgeable":                       DeletionRecoveryLevelRecoverablePositivePurgeable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeletionRecoveryLevel(input)
	return &out, nil
}

type SasTokenType string

const (
	SasTokenTypeAccount SasTokenType = "account"
	SasTokenTypeService SasTokenType = "service"
)

func PossibleValuesForSasTokenType() []string {
	return []string{
		string(SasTokenTypeAccount),
		string(SasTokenTypeService),
	}
}

func (s *SasTokenType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSasTokenType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSasTokenType(input string) (*SasTokenType, error) {
	vals := map[string]SasTokenType{
		"account": SasTokenTypeAccount,
		"service": SasTokenTypeService,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SasTokenType(input)
	return &out, nil
}
