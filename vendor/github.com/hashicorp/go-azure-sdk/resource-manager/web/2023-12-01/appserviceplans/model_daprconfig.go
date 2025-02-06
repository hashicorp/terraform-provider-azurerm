package appserviceplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DaprConfig struct {
	AppId              *string       `json:"appId,omitempty"`
	AppPort            *int64        `json:"appPort,omitempty"`
	EnableApiLogging   *bool         `json:"enableApiLogging,omitempty"`
	Enabled            *bool         `json:"enabled,omitempty"`
	HTTPMaxRequestSize *int64        `json:"httpMaxRequestSize,omitempty"`
	HTTPReadBufferSize *int64        `json:"httpReadBufferSize,omitempty"`
	LogLevel           *DaprLogLevel `json:"logLevel,omitempty"`
}
