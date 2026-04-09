package storageinsights

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageInsightState string

const (
	StorageInsightStateERROR StorageInsightState = "ERROR"
	StorageInsightStateOK    StorageInsightState = "OK"
)

func PossibleValuesForStorageInsightState() []string {
	return []string{
		string(StorageInsightStateERROR),
		string(StorageInsightStateOK),
	}
}

func (s *StorageInsightState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageInsightState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageInsightState(input string) (*StorageInsightState, error) {
	vals := map[string]StorageInsightState{
		"error": StorageInsightStateERROR,
		"ok":    StorageInsightStateOK,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageInsightState(input)
	return &out, nil
}
