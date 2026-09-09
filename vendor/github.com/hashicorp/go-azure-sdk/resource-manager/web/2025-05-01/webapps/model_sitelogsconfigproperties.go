package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SiteLogsConfigProperties struct {
	ApplicationLogs       *ApplicationLogsConfig `json:"applicationLogs,omitempty"`
	DetailedErrorMessages *EnabledConfig         `json:"detailedErrorMessages,omitempty"`
	FailedRequestsTracing *EnabledConfig         `json:"failedRequestsTracing,omitempty"`
	HTTPLogs              *HTTPLogsConfig        `json:"httpLogs,omitempty"`
}
