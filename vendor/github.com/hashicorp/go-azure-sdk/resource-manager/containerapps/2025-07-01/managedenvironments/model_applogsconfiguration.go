package managedenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppLogsConfiguration struct {
	Destination               *string                    `json:"destination,omitempty"`
	LogAnalyticsConfiguration *LogAnalyticsConfiguration `json:"logAnalyticsConfiguration,omitempty"`
}
