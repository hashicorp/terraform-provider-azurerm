package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActivityPolicy struct {
	Retry                  *interface{} `json:"retry,omitempty"`
	RetryIntervalInSeconds *int64       `json:"retryIntervalInSeconds,omitempty"`
	SecureInput            *bool        `json:"secureInput,omitempty"`
	SecureOutput           *bool        `json:"secureOutput,omitempty"`
	Timeout                *interface{} `json:"timeout,omitempty"`
}
