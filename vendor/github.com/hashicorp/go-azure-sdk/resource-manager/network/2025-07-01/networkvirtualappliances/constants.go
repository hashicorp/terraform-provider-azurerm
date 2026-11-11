package networkvirtualappliances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NicTypeInRequest string

const (
	NicTypeInRequestPrivateNic NicTypeInRequest = "PrivateNic"
	NicTypeInRequestPublicNic  NicTypeInRequest = "PublicNic"
)

func PossibleValuesForNicTypeInRequest() []string {
	return []string{
		string(NicTypeInRequestPrivateNic),
		string(NicTypeInRequestPublicNic),
	}
}

func (s *NicTypeInRequest) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNicTypeInRequest(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNicTypeInRequest(input string) (*NicTypeInRequest, error) {
	vals := map[string]NicTypeInRequest{
		"privatenic": NicTypeInRequestPrivateNic,
		"publicnic":  NicTypeInRequestPublicNic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NicTypeInRequest(input)
	return &out, nil
}

type NicTypeInResponse string

const (
	NicTypeInResponseAdditionalNic NicTypeInResponse = "AdditionalNic"
	NicTypeInResponsePrivateNic    NicTypeInResponse = "PrivateNic"
	NicTypeInResponsePublicNic     NicTypeInResponse = "PublicNic"
)

func PossibleValuesForNicTypeInResponse() []string {
	return []string{
		string(NicTypeInResponseAdditionalNic),
		string(NicTypeInResponsePrivateNic),
		string(NicTypeInResponsePublicNic),
	}
}

func (s *NicTypeInResponse) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNicTypeInResponse(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNicTypeInResponse(input string) (*NicTypeInResponse, error) {
	vals := map[string]NicTypeInResponse{
		"additionalnic": NicTypeInResponseAdditionalNic,
		"privatenic":    NicTypeInResponsePrivateNic,
		"publicnic":     NicTypeInResponsePublicNic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NicTypeInResponse(input)
	return &out, nil
}

type NvaNicType string

const (
	NvaNicTypeAdditionalPrivateNic NvaNicType = "AdditionalPrivateNic"
	NvaNicTypeAdditionalPublicNic  NvaNicType = "AdditionalPublicNic"
	NvaNicTypePrivateNic           NvaNicType = "PrivateNic"
	NvaNicTypePublicNic            NvaNicType = "PublicNic"
)

func PossibleValuesForNvaNicType() []string {
	return []string{
		string(NvaNicTypeAdditionalPrivateNic),
		string(NvaNicTypeAdditionalPublicNic),
		string(NvaNicTypePrivateNic),
		string(NvaNicTypePublicNic),
	}
}

func (s *NvaNicType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNvaNicType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNvaNicType(input string) (*NvaNicType, error) {
	vals := map[string]NvaNicType{
		"additionalprivatenic": NvaNicTypeAdditionalPrivateNic,
		"additionalpublicnic":  NvaNicTypeAdditionalPublicNic,
		"privatenic":           NvaNicTypePrivateNic,
		"publicnic":            NvaNicTypePublicNic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NvaNicType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"canceled":  ProvisioningStateCanceled,
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
