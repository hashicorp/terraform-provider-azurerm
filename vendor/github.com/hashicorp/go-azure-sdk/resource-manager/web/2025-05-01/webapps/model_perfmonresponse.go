package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PerfMonResponse struct {
	Code    *string     `json:"code,omitempty"`
	Data    *PerfMonSet `json:"data,omitempty"`
	Message *string     `json:"message,omitempty"`
}
