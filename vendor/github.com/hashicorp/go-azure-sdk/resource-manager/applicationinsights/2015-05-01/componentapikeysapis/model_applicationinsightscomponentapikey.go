package componentapikeysapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationInsightsComponentAPIKey struct {
	ApiKey                *string   `json:"apiKey,omitempty"`
	CreatedDate           *string   `json:"createdDate,omitempty"`
	Id                    *string   `json:"id,omitempty"`
	LinkedReadProperties  *[]string `json:"linkedReadProperties,omitempty"`
	LinkedWriteProperties *[]string `json:"linkedWriteProperties,omitempty"`
	Name                  *string   `json:"name,omitempty"`
}
