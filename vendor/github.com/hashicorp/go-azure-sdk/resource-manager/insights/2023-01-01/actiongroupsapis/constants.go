package actiongroupsapis

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReceiverStatus string

const (
	ReceiverStatusDisabled     ReceiverStatus = "Disabled"
	ReceiverStatusEnabled      ReceiverStatus = "Enabled"
	ReceiverStatusNotSpecified ReceiverStatus = "NotSpecified"
)

func PossibleValuesForReceiverStatus() []string {
	return []string{
		string(ReceiverStatusDisabled),
		string(ReceiverStatusEnabled),
		string(ReceiverStatusNotSpecified),
	}
}

func parseReceiverStatus(input string) (*ReceiverStatus, error) {
	vals := map[string]ReceiverStatus{
		"disabled":     ReceiverStatusDisabled,
		"enabled":      ReceiverStatusEnabled,
		"notspecified": ReceiverStatusNotSpecified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReceiverStatus(input)
	return &out, nil
}
