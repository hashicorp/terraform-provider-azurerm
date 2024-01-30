package quotas

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceState string

const (
	ResourceStateCanceled  ResourceState = "Canceled"
	ResourceStateDeleted   ResourceState = "Deleted"
	ResourceStateFailed    ResourceState = "Failed"
	ResourceStateSucceeded ResourceState = "Succeeded"
)

func PossibleValuesForResourceState() []string {
	return []string{
		string(ResourceStateCanceled),
		string(ResourceStateDeleted),
		string(ResourceStateFailed),
		string(ResourceStateSucceeded),
	}
}

func (s *ResourceState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceState(input string) (*ResourceState, error) {
	vals := map[string]ResourceState{
		"canceled":  ResourceStateCanceled,
		"deleted":   ResourceStateDeleted,
		"failed":    ResourceStateFailed,
		"succeeded": ResourceStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceState(input)
	return &out, nil
}
