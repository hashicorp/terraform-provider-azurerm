package logger

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoggerType string

const (
	LoggerTypeApplicationInsights LoggerType = "applicationInsights"
	LoggerTypeAzureEventHub       LoggerType = "azureEventHub"
	LoggerTypeAzureMonitor        LoggerType = "azureMonitor"
)

func PossibleValuesForLoggerType() []string {
	return []string{
		string(LoggerTypeApplicationInsights),
		string(LoggerTypeAzureEventHub),
		string(LoggerTypeAzureMonitor),
	}
}

func (s *LoggerType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLoggerType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLoggerType(input string) (*LoggerType, error) {
	vals := map[string]LoggerType{
		"applicationinsights": LoggerTypeApplicationInsights,
		"azureeventhub":       LoggerTypeAzureEventHub,
		"azuremonitor":        LoggerTypeAzureMonitor,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LoggerType(input)
	return &out, nil
}
