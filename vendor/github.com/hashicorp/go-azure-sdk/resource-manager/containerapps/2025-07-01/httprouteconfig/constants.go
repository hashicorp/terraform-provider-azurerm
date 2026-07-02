package httprouteconfig

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BindingType string

const (
	BindingTypeAuto       BindingType = "Auto"
	BindingTypeDisabled   BindingType = "Disabled"
	BindingTypeSniEnabled BindingType = "SniEnabled"
)

func PossibleValuesForBindingType() []string {
	return []string{
		string(BindingTypeAuto),
		string(BindingTypeDisabled),
		string(BindingTypeSniEnabled),
	}
}

func (s *BindingType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBindingType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBindingType(input string) (*BindingType, error) {
	vals := map[string]BindingType{
		"auto":       BindingTypeAuto,
		"disabled":   BindingTypeDisabled,
		"snienabled": BindingTypeSniEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BindingType(input)
	return &out, nil
}

type HTTPRouteProvisioningState string

const (
	HTTPRouteProvisioningStateCanceled  HTTPRouteProvisioningState = "Canceled"
	HTTPRouteProvisioningStateDeleting  HTTPRouteProvisioningState = "Deleting"
	HTTPRouteProvisioningStateFailed    HTTPRouteProvisioningState = "Failed"
	HTTPRouteProvisioningStatePending   HTTPRouteProvisioningState = "Pending"
	HTTPRouteProvisioningStateSucceeded HTTPRouteProvisioningState = "Succeeded"
	HTTPRouteProvisioningStateUpdating  HTTPRouteProvisioningState = "Updating"
	HTTPRouteProvisioningStateWaiting   HTTPRouteProvisioningState = "Waiting"
)

func PossibleValuesForHTTPRouteProvisioningState() []string {
	return []string{
		string(HTTPRouteProvisioningStateCanceled),
		string(HTTPRouteProvisioningStateDeleting),
		string(HTTPRouteProvisioningStateFailed),
		string(HTTPRouteProvisioningStatePending),
		string(HTTPRouteProvisioningStateSucceeded),
		string(HTTPRouteProvisioningStateUpdating),
		string(HTTPRouteProvisioningStateWaiting),
	}
}

func (s *HTTPRouteProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHTTPRouteProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHTTPRouteProvisioningState(input string) (*HTTPRouteProvisioningState, error) {
	vals := map[string]HTTPRouteProvisioningState{
		"canceled":  HTTPRouteProvisioningStateCanceled,
		"deleting":  HTTPRouteProvisioningStateDeleting,
		"failed":    HTTPRouteProvisioningStateFailed,
		"pending":   HTTPRouteProvisioningStatePending,
		"succeeded": HTTPRouteProvisioningStateSucceeded,
		"updating":  HTTPRouteProvisioningStateUpdating,
		"waiting":   HTTPRouteProvisioningStateWaiting,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HTTPRouteProvisioningState(input)
	return &out, nil
}
