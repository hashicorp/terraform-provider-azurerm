package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Cron struct {
	Expression *string `json:"expression,omitempty"`
	StartTime  *string `json:"startTime,omitempty"`
	TimeZone   *string `json:"timeZone,omitempty"`
}
