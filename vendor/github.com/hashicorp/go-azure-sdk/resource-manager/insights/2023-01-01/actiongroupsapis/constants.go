package actiongroupsapis

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReceiverStatus string

const (
	ReceiverStatusDisabled     ReceiverStatus = "Disabled"
	ReceiverStatusEnabled      ReceiverStatus = "Enabled"
	ReceiverStatusNotSpecified ReceiverStatus = "NotSpecified"
)

func PossibleValuesForReceiverStatus() []string {
	return []string{
		string(ReceiverStatusDisabled),
		string(ReceiverStatusEnabled),
		string(ReceiverStatusNotSpecified),
	}
}

func (s *ReceiverStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReceiverStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReceiverStatus(input string) (*ReceiverStatus, error) {
	vals := map[string]ReceiverStatus{
		"disabled":     ReceiverStatusDisabled,
		"enabled":      ReceiverStatusEnabled,
		"notspecified": ReceiverStatusNotSpecified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReceiverStatus(input)
	return &out, nil
}
