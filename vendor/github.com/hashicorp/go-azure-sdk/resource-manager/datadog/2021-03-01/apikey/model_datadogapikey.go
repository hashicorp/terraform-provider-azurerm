package apikey

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatadogApiKey struct {
	Created   *string `json:"created,omitempty"`
	CreatedBy *string `json:"createdBy,omitempty"`
	Key       string  `json:"key"`
	Name      *string `json:"name,omitempty"`
}
