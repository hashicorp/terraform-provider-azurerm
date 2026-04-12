package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebHookActivityTypeProperties struct {
	Authentication         *WebActivityAuthentication `json:"authentication,omitempty"`
	Body                   *interface{}               `json:"body,omitempty"`
	Headers                *map[string]interface{}    `json:"headers,omitempty"`
	Method                 WebHookActivityMethod      `json:"method"`
	ReportStatusOnCallBack *interface{}               `json:"reportStatusOnCallBack,omitempty"`
	Timeout                *string                    `json:"timeout,omitempty"`
	Url                    interface{}                `json:"url"`
}
