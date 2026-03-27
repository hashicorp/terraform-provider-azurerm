package dscconfiguration

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentSourceType string

const (
	ContentSourceTypeEmbeddedContent ContentSourceType = "embeddedContent"
	ContentSourceTypeUri             ContentSourceType = "uri"
)

func PossibleValuesForContentSourceType() []string {
	return []string{
		string(ContentSourceTypeEmbeddedContent),
		string(ContentSourceTypeUri),
	}
}

func (s *ContentSourceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContentSourceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContentSourceType(input string) (*ContentSourceType, error) {
	vals := map[string]ContentSourceType{
		"embeddedcontent": ContentSourceTypeEmbeddedContent,
		"uri":             ContentSourceTypeUri,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContentSourceType(input)
	return &out, nil
}

type DscConfigurationProvisioningState string

const (
	DscConfigurationProvisioningStateSucceeded DscConfigurationProvisioningState = "Succeeded"
)

func PossibleValuesForDscConfigurationProvisioningState() []string {
	return []string{
		string(DscConfigurationProvisioningStateSucceeded),
	}
}

func (s *DscConfigurationProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDscConfigurationProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDscConfigurationProvisioningState(input string) (*DscConfigurationProvisioningState, error) {
	vals := map[string]DscConfigurationProvisioningState{
		"succeeded": DscConfigurationProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DscConfigurationProvisioningState(input)
	return &out, nil
}

type DscConfigurationState string

const (
	DscConfigurationStateEdit      DscConfigurationState = "Edit"
	DscConfigurationStateNew       DscConfigurationState = "New"
	DscConfigurationStatePublished DscConfigurationState = "Published"
)

func PossibleValuesForDscConfigurationState() []string {
	return []string{
		string(DscConfigurationStateEdit),
		string(DscConfigurationStateNew),
		string(DscConfigurationStatePublished),
	}
}

func (s *DscConfigurationState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDscConfigurationState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDscConfigurationState(input string) (*DscConfigurationState, error) {
	vals := map[string]DscConfigurationState{
		"edit":      DscConfigurationStateEdit,
		"new":       DscConfigurationStateNew,
		"published": DscConfigurationStatePublished,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DscConfigurationState(input)
	return &out, nil
}
