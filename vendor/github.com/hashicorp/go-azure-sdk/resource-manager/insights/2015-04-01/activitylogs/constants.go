package activitylogs

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventLevel string

const (
	EventLevelCritical      EventLevel = "Critical"
	EventLevelError         EventLevel = "Error"
	EventLevelInformational EventLevel = "Informational"
	EventLevelVerbose       EventLevel = "Verbose"
	EventLevelWarning       EventLevel = "Warning"
)

func PossibleValuesForEventLevel() []string {
	return []string{
		string(EventLevelCritical),
		string(EventLevelError),
		string(EventLevelInformational),
		string(EventLevelVerbose),
		string(EventLevelWarning),
	}
}

func parseEventLevel(input string) (*EventLevel, error) {
	vals := map[string]EventLevel{
		"critical":      EventLevelCritical,
		"error":         EventLevelError,
		"informational": EventLevelInformational,
		"verbose":       EventLevelVerbose,
		"warning":       EventLevelWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EventLevel(input)
	return &out, nil
}
