package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StopOnNoConnectConfiguration struct {
	GracePeriodMinutes *int64                       `json:"gracePeriodMinutes,omitempty"`
	Status             *StopOnNoConnectEnableStatus `json:"status,omitempty"`
}
