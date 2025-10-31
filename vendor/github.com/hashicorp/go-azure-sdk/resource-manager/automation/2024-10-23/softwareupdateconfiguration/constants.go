package softwareupdateconfiguration

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinuxUpdateClasses string

const (
	LinuxUpdateClassesCritical     LinuxUpdateClasses = "Critical"
	LinuxUpdateClassesOther        LinuxUpdateClasses = "Other"
	LinuxUpdateClassesSecurity     LinuxUpdateClasses = "Security"
	LinuxUpdateClassesUnclassified LinuxUpdateClasses = "Unclassified"
)

func PossibleValuesForLinuxUpdateClasses() []string {
	return []string{
		string(LinuxUpdateClassesCritical),
		string(LinuxUpdateClassesOther),
		string(LinuxUpdateClassesSecurity),
		string(LinuxUpdateClassesUnclassified),
	}
}

func (s *LinuxUpdateClasses) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLinuxUpdateClasses(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLinuxUpdateClasses(input string) (*LinuxUpdateClasses, error) {
	vals := map[string]LinuxUpdateClasses{
		"critical":     LinuxUpdateClassesCritical,
		"other":        LinuxUpdateClassesOther,
		"security":     LinuxUpdateClassesSecurity,
		"unclassified": LinuxUpdateClassesUnclassified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LinuxUpdateClasses(input)
	return &out, nil
}

type OperatingSystemType string

const (
	OperatingSystemTypeLinux   OperatingSystemType = "Linux"
	OperatingSystemTypeWindows OperatingSystemType = "Windows"
)

func PossibleValuesForOperatingSystemType() []string {
	return []string{
		string(OperatingSystemTypeLinux),
		string(OperatingSystemTypeWindows),
	}
}

func (s *OperatingSystemType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperatingSystemType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperatingSystemType(input string) (*OperatingSystemType, error) {
	vals := map[string]OperatingSystemType{
		"linux":   OperatingSystemTypeLinux,
		"windows": OperatingSystemTypeWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperatingSystemType(input)
	return &out, nil
}

type ScheduleDay string

const (
	ScheduleDayFriday    ScheduleDay = "Friday"
	ScheduleDayMonday    ScheduleDay = "Monday"
	ScheduleDaySaturday  ScheduleDay = "Saturday"
	ScheduleDaySunday    ScheduleDay = "Sunday"
	ScheduleDayThursday  ScheduleDay = "Thursday"
	ScheduleDayTuesday   ScheduleDay = "Tuesday"
	ScheduleDayWednesday ScheduleDay = "Wednesday"
)

func PossibleValuesForScheduleDay() []string {
	return []string{
		string(ScheduleDayFriday),
		string(ScheduleDayMonday),
		string(ScheduleDaySaturday),
		string(ScheduleDaySunday),
		string(ScheduleDayThursday),
		string(ScheduleDayTuesday),
		string(ScheduleDayWednesday),
	}
}

func (s *ScheduleDay) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScheduleDay(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScheduleDay(input string) (*ScheduleDay, error) {
	vals := map[string]ScheduleDay{
		"friday":    ScheduleDayFriday,
		"monday":    ScheduleDayMonday,
		"saturday":  ScheduleDaySaturday,
		"sunday":    ScheduleDaySunday,
		"thursday":  ScheduleDayThursday,
		"tuesday":   ScheduleDayTuesday,
		"wednesday": ScheduleDayWednesday,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScheduleDay(input)
	return &out, nil
}

type ScheduleFrequency string

const (
	ScheduleFrequencyDay     ScheduleFrequency = "Day"
	ScheduleFrequencyHour    ScheduleFrequency = "Hour"
	ScheduleFrequencyMinute  ScheduleFrequency = "Minute"
	ScheduleFrequencyMonth   ScheduleFrequency = "Month"
	ScheduleFrequencyOneTime ScheduleFrequency = "OneTime"
	ScheduleFrequencyWeek    ScheduleFrequency = "Week"
)

func PossibleValuesForScheduleFrequency() []string {
	return []string{
		string(ScheduleFrequencyDay),
		string(ScheduleFrequencyHour),
		string(ScheduleFrequencyMinute),
		string(ScheduleFrequencyMonth),
		string(ScheduleFrequencyOneTime),
		string(ScheduleFrequencyWeek),
	}
}

func (s *ScheduleFrequency) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScheduleFrequency(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScheduleFrequency(input string) (*ScheduleFrequency, error) {
	vals := map[string]ScheduleFrequency{
		"day":     ScheduleFrequencyDay,
		"hour":    ScheduleFrequencyHour,
		"minute":  ScheduleFrequencyMinute,
		"month":   ScheduleFrequencyMonth,
		"onetime": ScheduleFrequencyOneTime,
		"week":    ScheduleFrequencyWeek,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScheduleFrequency(input)
	return &out, nil
}

type TagOperators string

const (
	TagOperatorsAll TagOperators = "All"
	TagOperatorsAny TagOperators = "Any"
)

func PossibleValuesForTagOperators() []string {
	return []string{
		string(TagOperatorsAll),
		string(TagOperatorsAny),
	}
}

func (s *TagOperators) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTagOperators(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTagOperators(input string) (*TagOperators, error) {
	vals := map[string]TagOperators{
		"all": TagOperatorsAll,
		"any": TagOperatorsAny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TagOperators(input)
	return &out, nil
}

type WindowsUpdateClasses string

const (
	WindowsUpdateClassesCritical     WindowsUpdateClasses = "Critical"
	WindowsUpdateClassesDefinition   WindowsUpdateClasses = "Definition"
	WindowsUpdateClassesFeaturePack  WindowsUpdateClasses = "FeaturePack"
	WindowsUpdateClassesSecurity     WindowsUpdateClasses = "Security"
	WindowsUpdateClassesServicePack  WindowsUpdateClasses = "ServicePack"
	WindowsUpdateClassesTools        WindowsUpdateClasses = "Tools"
	WindowsUpdateClassesUnclassified WindowsUpdateClasses = "Unclassified"
	WindowsUpdateClassesUpdateRollup WindowsUpdateClasses = "UpdateRollup"
	WindowsUpdateClassesUpdates      WindowsUpdateClasses = "Updates"
)

func PossibleValuesForWindowsUpdateClasses() []string {
	return []string{
		string(WindowsUpdateClassesCritical),
		string(WindowsUpdateClassesDefinition),
		string(WindowsUpdateClassesFeaturePack),
		string(WindowsUpdateClassesSecurity),
		string(WindowsUpdateClassesServicePack),
		string(WindowsUpdateClassesTools),
		string(WindowsUpdateClassesUnclassified),
		string(WindowsUpdateClassesUpdateRollup),
		string(WindowsUpdateClassesUpdates),
	}
}

func (s *WindowsUpdateClasses) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWindowsUpdateClasses(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWindowsUpdateClasses(input string) (*WindowsUpdateClasses, error) {
	vals := map[string]WindowsUpdateClasses{
		"critical":     WindowsUpdateClassesCritical,
		"definition":   WindowsUpdateClassesDefinition,
		"featurepack":  WindowsUpdateClassesFeaturePack,
		"security":     WindowsUpdateClassesSecurity,
		"servicepack":  WindowsUpdateClassesServicePack,
		"tools":        WindowsUpdateClassesTools,
		"unclassified": WindowsUpdateClassesUnclassified,
		"updaterollup": WindowsUpdateClassesUpdateRollup,
		"updates":      WindowsUpdateClassesUpdates,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WindowsUpdateClasses(input)
	return &out, nil
}
