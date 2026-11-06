package subscriptionquotaallocation

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RequestState string

const (
	RequestStateAccepted   RequestState = "Accepted"
	RequestStateCanceled   RequestState = "Canceled"
	RequestStateCreated    RequestState = "Created"
	RequestStateFailed     RequestState = "Failed"
	RequestStateInProgress RequestState = "InProgress"
	RequestStateInvalid    RequestState = "Invalid"
	RequestStateSucceeded  RequestState = "Succeeded"
)

func PossibleValuesForRequestState() []string {
	return []string{
		string(RequestStateAccepted),
		string(RequestStateCanceled),
		string(RequestStateCreated),
		string(RequestStateFailed),
		string(RequestStateInProgress),
		string(RequestStateInvalid),
		string(RequestStateSucceeded),
	}
}

func (s *RequestState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRequestState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRequestState(input string) (*RequestState, error) {
	vals := map[string]RequestState{
		"accepted":   RequestStateAccepted,
		"canceled":   RequestStateCanceled,
		"created":    RequestStateCreated,
		"failed":     RequestStateFailed,
		"inprogress": RequestStateInProgress,
		"invalid":    RequestStateInvalid,
		"succeeded":  RequestStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RequestState(input)
	return &out, nil
}
