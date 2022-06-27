package tenantconfiguration

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationName string

const (
	ConfigurationNameDefault ConfigurationName = "default"
)

func PossibleValuesForConfigurationName() []string {
	return []string{
		string(ConfigurationNameDefault),
	}
}

func parseConfigurationName(input string) (*ConfigurationName, error) {
	vals := map[string]ConfigurationName{
		"default": ConfigurationNameDefault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConfigurationName(input)
	return &out, nil
}
