package configurations

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationSource string

const (
	ConfigurationSourceSystemNegativedefault ConfigurationSource = "system-default"
	ConfigurationSourceUserNegativeoverride  ConfigurationSource = "user-override"
)

func PossibleValuesForConfigurationSource() []string {
	return []string{
		string(ConfigurationSourceSystemNegativedefault),
		string(ConfigurationSourceUserNegativeoverride),
	}
}

func (s *ConfigurationSource) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConfigurationSource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConfigurationSource(input string) (*ConfigurationSource, error) {
	vals := map[string]ConfigurationSource{
		"system-default": ConfigurationSourceSystemNegativedefault,
		"user-override":  ConfigurationSourceUserNegativeoverride,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConfigurationSource(input)
	return &out, nil
}

type IsConfigPendingRestart string

const (
	IsConfigPendingRestartFalse IsConfigPendingRestart = "False"
	IsConfigPendingRestartTrue  IsConfigPendingRestart = "True"
)

func PossibleValuesForIsConfigPendingRestart() []string {
	return []string{
		string(IsConfigPendingRestartFalse),
		string(IsConfigPendingRestartTrue),
	}
}

func (s *IsConfigPendingRestart) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIsConfigPendingRestart(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIsConfigPendingRestart(input string) (*IsConfigPendingRestart, error) {
	vals := map[string]IsConfigPendingRestart{
		"false": IsConfigPendingRestartFalse,
		"true":  IsConfigPendingRestartTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IsConfigPendingRestart(input)
	return &out, nil
}

type IsDynamicConfig string

const (
	IsDynamicConfigFalse IsDynamicConfig = "False"
	IsDynamicConfigTrue  IsDynamicConfig = "True"
)

func PossibleValuesForIsDynamicConfig() []string {
	return []string{
		string(IsDynamicConfigFalse),
		string(IsDynamicConfigTrue),
	}
}

func (s *IsDynamicConfig) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIsDynamicConfig(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIsDynamicConfig(input string) (*IsDynamicConfig, error) {
	vals := map[string]IsDynamicConfig{
		"false": IsDynamicConfigFalse,
		"true":  IsDynamicConfigTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IsDynamicConfig(input)
	return &out, nil
}

type IsReadOnly string

const (
	IsReadOnlyFalse IsReadOnly = "False"
	IsReadOnlyTrue  IsReadOnly = "True"
)

func PossibleValuesForIsReadOnly() []string {
	return []string{
		string(IsReadOnlyFalse),
		string(IsReadOnlyTrue),
	}
}

func (s *IsReadOnly) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIsReadOnly(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIsReadOnly(input string) (*IsReadOnly, error) {
	vals := map[string]IsReadOnly{
		"false": IsReadOnlyFalse,
		"true":  IsReadOnlyTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IsReadOnly(input)
	return &out, nil
}
