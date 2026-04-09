package globalschedules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RetargetScheduleProperties struct {
	CurrentResourceId *string `json:"currentResourceId,omitempty"`
	TargetResourceId  *string `json:"targetResourceId,omitempty"`
}
