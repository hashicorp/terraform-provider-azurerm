package blobauditing

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobAuditingPolicyState string

const (
	BlobAuditingPolicyStateDisabled BlobAuditingPolicyState = "Disabled"
	BlobAuditingPolicyStateEnabled  BlobAuditingPolicyState = "Enabled"
)

func PossibleValuesForBlobAuditingPolicyState() []string {
	return []string{
		string(BlobAuditingPolicyStateDisabled),
		string(BlobAuditingPolicyStateEnabled),
	}
}

func (s *BlobAuditingPolicyState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBlobAuditingPolicyState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBlobAuditingPolicyState(input string) (*BlobAuditingPolicyState, error) {
	vals := map[string]BlobAuditingPolicyState{
		"disabled": BlobAuditingPolicyStateDisabled,
		"enabled":  BlobAuditingPolicyStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BlobAuditingPolicyState(input)
	return &out, nil
}
