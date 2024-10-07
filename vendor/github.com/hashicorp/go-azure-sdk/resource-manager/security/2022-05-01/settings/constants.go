package settings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SettingKind string

const (
	SettingKindAlertSuppressionSetting SettingKind = "AlertSuppressionSetting"
	SettingKindAlertSyncSettings       SettingKind = "AlertSyncSettings"
	SettingKindDataExportSettings      SettingKind = "DataExportSettings"
)

func PossibleValuesForSettingKind() []string {
	return []string{
		string(SettingKindAlertSuppressionSetting),
		string(SettingKindAlertSyncSettings),
		string(SettingKindDataExportSettings),
	}
}

func (s *SettingKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSettingKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSettingKind(input string) (*SettingKind, error) {
	vals := map[string]SettingKind{
		"alertsuppressionsetting": SettingKindAlertSuppressionSetting,
		"alertsyncsettings":       SettingKindAlertSyncSettings,
		"dataexportsettings":      SettingKindDataExportSettings,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SettingKind(input)
	return &out, nil
}

type SettingName string

const (
	SettingNameMCAS                           SettingName = "MCAS"
	SettingNameSentinel                       SettingName = "Sentinel"
	SettingNameWDATP                          SettingName = "WDATP"
	SettingNameWDATPEXCLUDELINUXPUBLICPREVIEW SettingName = "WDATP_EXCLUDE_LINUX_PUBLIC_PREVIEW"
	SettingNameWDATPUNIFIEDSOLUTION           SettingName = "WDATP_UNIFIED_SOLUTION"
)

func PossibleValuesForSettingName() []string {
	return []string{
		string(SettingNameMCAS),
		string(SettingNameSentinel),
		string(SettingNameWDATP),
		string(SettingNameWDATPEXCLUDELINUXPUBLICPREVIEW),
		string(SettingNameWDATPUNIFIEDSOLUTION),
	}
}

func (s *SettingName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSettingName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSettingName(input string) (*SettingName, error) {
	vals := map[string]SettingName{
		"mcas":                               SettingNameMCAS,
		"sentinel":                           SettingNameSentinel,
		"wdatp":                              SettingNameWDATP,
		"wdatp_exclude_linux_public_preview": SettingNameWDATPEXCLUDELINUXPUBLICPREVIEW,
		"wdatp_unified_solution":             SettingNameWDATPUNIFIEDSOLUTION,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SettingName(input)
	return &out, nil
}
