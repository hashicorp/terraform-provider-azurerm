package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContinuousWebJobProperties struct {
	DetailedStatus *string                 `json:"detailed_status,omitempty"`
	Error          *string                 `json:"error,omitempty"`
	ExtraInfoUrl   *string                 `json:"extra_info_url,omitempty"`
	LogUrl         *string                 `json:"log_url,omitempty"`
	RunCommand     *string                 `json:"run_command,omitempty"`
	Settings       *map[string]interface{} `json:"settings,omitempty"`
	Status         *ContinuousWebJobStatus `json:"status,omitempty"`
	Url            *string                 `json:"url,omitempty"`
	UsingSdk       *bool                   `json:"using_sdk,omitempty"`
	WebJobType     *WebJobType             `json:"web_job_type,omitempty"`
}
