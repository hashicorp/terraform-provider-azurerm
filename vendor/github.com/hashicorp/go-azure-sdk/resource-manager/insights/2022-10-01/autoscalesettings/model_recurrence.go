package autoscalesettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Recurrence struct {
	Frequency RecurrenceFrequency `json:"frequency"`
	Schedule  RecurrentSchedule   `json:"schedule"`
}
