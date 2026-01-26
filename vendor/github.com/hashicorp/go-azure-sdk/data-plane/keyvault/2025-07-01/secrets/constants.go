package secrets

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentType string

const (
	ContentTypeApplicationXNegativepemNegativefile ContentType = "application/x-pem-file"
	ContentTypeApplicationXNegativepkcsOneTwo      ContentType = "application/x-pkcs12"
)

func PossibleValuesForContentType() []string {
	return []string{
		string(ContentTypeApplicationXNegativepemNegativefile),
		string(ContentTypeApplicationXNegativepkcsOneTwo),
	}
}

func (s *ContentType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContentType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContentType(input string) (*ContentType, error) {
	vals := map[string]ContentType{
		"application/x-pem-file": ContentTypeApplicationXNegativepemNegativefile,
		"application/x-pkcs12":   ContentTypeApplicationXNegativepkcsOneTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContentType(input)
	return &out, nil
}

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
