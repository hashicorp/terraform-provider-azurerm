package apikey

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatadogApiKeyListResponse struct {
	NextLink *string          `json:"nextLink,omitempty"`
	Value    *[]DatadogApiKey `json:"value,omitempty"`
}
