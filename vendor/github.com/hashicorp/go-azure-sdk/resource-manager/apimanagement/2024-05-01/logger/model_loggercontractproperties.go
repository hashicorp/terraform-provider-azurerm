package logger

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoggerContractProperties struct {
	Credentials *map[string]string `json:"credentials,omitempty"`
	Description *string            `json:"description,omitempty"`
	IsBuffered  *bool              `json:"isBuffered,omitempty"`
	LoggerType  LoggerType         `json:"loggerType"`
	ResourceId  *string            `json:"resourceId,omitempty"`
}
