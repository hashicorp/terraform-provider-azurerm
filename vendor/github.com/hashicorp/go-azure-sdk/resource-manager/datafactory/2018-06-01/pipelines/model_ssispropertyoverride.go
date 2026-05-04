package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SSISPropertyOverride struct {
	IsSensitive *bool       `json:"isSensitive,omitempty"`
	Value       interface{} `json:"value"`
}
