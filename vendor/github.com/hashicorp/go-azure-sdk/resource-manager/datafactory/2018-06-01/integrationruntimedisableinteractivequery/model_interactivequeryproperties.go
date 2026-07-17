package integrationruntimedisableinteractivequery

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InteractiveQueryProperties struct {
	AutoTerminationMinutes *int64                       `json:"autoTerminationMinutes,omitempty"`
	Status                 *InteractiveCapabilityStatus `json:"status,omitempty"`
}
