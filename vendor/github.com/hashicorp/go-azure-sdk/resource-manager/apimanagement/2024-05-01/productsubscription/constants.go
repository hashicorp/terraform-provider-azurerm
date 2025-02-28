package productsubscription

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubscriptionState string

const (
	SubscriptionStateActive    SubscriptionState = "active"
	SubscriptionStateCancelled SubscriptionState = "cancelled"
	SubscriptionStateExpired   SubscriptionState = "expired"
	SubscriptionStateRejected  SubscriptionState = "rejected"
	SubscriptionStateSubmitted SubscriptionState = "submitted"
	SubscriptionStateSuspended SubscriptionState = "suspended"
)

func PossibleValuesForSubscriptionState() []string {
	return []string{
		string(SubscriptionStateActive),
		string(SubscriptionStateCancelled),
		string(SubscriptionStateExpired),
		string(SubscriptionStateRejected),
		string(SubscriptionStateSubmitted),
		string(SubscriptionStateSuspended),
	}
}

func (s *SubscriptionState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSubscriptionState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSubscriptionState(input string) (*SubscriptionState, error) {
	vals := map[string]SubscriptionState{
		"active":    SubscriptionStateActive,
		"cancelled": SubscriptionStateCancelled,
		"expired":   SubscriptionStateExpired,
		"rejected":  SubscriptionStateRejected,
		"submitted": SubscriptionStateSubmitted,
		"suspended": SubscriptionStateSuspended,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SubscriptionState(input)
	return &out, nil
}
