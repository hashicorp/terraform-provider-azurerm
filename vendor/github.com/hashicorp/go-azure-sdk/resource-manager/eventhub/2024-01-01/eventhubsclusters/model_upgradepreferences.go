package eventhubsclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpgradePreferences struct {
	StartDayOfWeek *StartDayOfWeek `json:"startDayOfWeek,omitempty"`
	StartHourOfDay *int64          `json:"startHourOfDay,omitempty"`
}
