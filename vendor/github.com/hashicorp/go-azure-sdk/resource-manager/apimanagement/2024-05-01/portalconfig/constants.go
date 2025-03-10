package portalconfig

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PortalSettingsCspMode string

const (
	PortalSettingsCspModeDisabled   PortalSettingsCspMode = "disabled"
	PortalSettingsCspModeEnabled    PortalSettingsCspMode = "enabled"
	PortalSettingsCspModeReportOnly PortalSettingsCspMode = "reportOnly"
)

func PossibleValuesForPortalSettingsCspMode() []string {
	return []string{
		string(PortalSettingsCspModeDisabled),
		string(PortalSettingsCspModeEnabled),
		string(PortalSettingsCspModeReportOnly),
	}
}

func (s *PortalSettingsCspMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePortalSettingsCspMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePortalSettingsCspMode(input string) (*PortalSettingsCspMode, error) {
	vals := map[string]PortalSettingsCspMode{
		"disabled":   PortalSettingsCspModeDisabled,
		"enabled":    PortalSettingsCspModeEnabled,
		"reportonly": PortalSettingsCspModeReportOnly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PortalSettingsCspMode(input)
	return &out, nil
}
