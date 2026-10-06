package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TriggeredWebJobProperties struct {
	Error                  *string                 `json:"error,omitempty"`
	ExtraInfoUrl           *string                 `json:"extra_info_url,omitempty"`
	HistoryUrl             *string                 `json:"history_url,omitempty"`
	LatestRun              *TriggeredJobRun        `json:"latest_run,omitempty"`
	PublicNetworkAccess    *string                 `json:"publicNetworkAccess,omitempty"`
	RunCommand             *string                 `json:"run_command,omitempty"`
	SchedulerLogsUrl       *string                 `json:"scheduler_logs_url,omitempty"`
	Settings               *map[string]interface{} `json:"settings,omitempty"`
	StorageAccountRequired *bool                   `json:"storageAccountRequired,omitempty"`
	Url                    *string                 `json:"url,omitempty"`
	UsingSdk               *bool                   `json:"using_sdk,omitempty"`
	WebJobType             *WebJobType             `json:"web_job_type,omitempty"`
}
