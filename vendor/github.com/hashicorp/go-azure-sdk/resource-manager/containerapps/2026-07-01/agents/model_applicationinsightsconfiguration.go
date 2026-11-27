package agents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationInsightsConfiguration struct {
	AppId            *string `json:"appId,omitempty"`
	ConnectionString *string `json:"connectionString,omitempty"`
}
