package resourcedetails

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceStatus string

const (
	ResourceStatusAllocated      ResourceStatus = "Allocated"
	ResourceStatusLeased         ResourceStatus = "Leased"
	ResourceStatusNotReady       ResourceStatus = "NotReady"
	ResourceStatusPendingReimage ResourceStatus = "PendingReimage"
	ResourceStatusPendingReturn  ResourceStatus = "PendingReturn"
	ResourceStatusProvisioning   ResourceStatus = "Provisioning"
	ResourceStatusReady          ResourceStatus = "Ready"
	ResourceStatusReimaging      ResourceStatus = "Reimaging"
	ResourceStatusReturned       ResourceStatus = "Returned"
	ResourceStatusStarting       ResourceStatus = "Starting"
	ResourceStatusUpdating       ResourceStatus = "Updating"
)

func PossibleValuesForResourceStatus() []string {
	return []string{
		string(ResourceStatusAllocated),
		string(ResourceStatusLeased),
		string(ResourceStatusNotReady),
		string(ResourceStatusPendingReimage),
		string(ResourceStatusPendingReturn),
		string(ResourceStatusProvisioning),
		string(ResourceStatusReady),
		string(ResourceStatusReimaging),
		string(ResourceStatusReturned),
		string(ResourceStatusStarting),
		string(ResourceStatusUpdating),
	}
}

func (s *ResourceStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceStatus(input string) (*ResourceStatus, error) {
	vals := map[string]ResourceStatus{
		"allocated":      ResourceStatusAllocated,
		"leased":         ResourceStatusLeased,
		"notready":       ResourceStatusNotReady,
		"pendingreimage": ResourceStatusPendingReimage,
		"pendingreturn":  ResourceStatusPendingReturn,
		"provisioning":   ResourceStatusProvisioning,
		"ready":          ResourceStatusReady,
		"reimaging":      ResourceStatusReimaging,
		"returned":       ResourceStatusReturned,
		"starting":       ResourceStatusStarting,
		"updating":       ResourceStatusUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceStatus(input)
	return &out, nil
}
