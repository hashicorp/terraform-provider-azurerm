package webtestsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebTestProperties struct {
	Configuration      *WebTestPropertiesConfiguration   `json:"Configuration,omitempty"`
	Description        *string                           `json:"Description,omitempty"`
	Enabled            *bool                             `json:"Enabled,omitempty"`
	Frequency          *int64                            `json:"Frequency,omitempty"`
	Kind               WebTestKind                       `json:"Kind"`
	Locations          []WebTestGeolocation              `json:"Locations"`
	Name               string                            `json:"Name"`
	ProvisioningState  *string                           `json:"provisioningState,omitempty"`
	Request            *WebTestPropertiesRequest         `json:"Request,omitempty"`
	RetryEnabled       *bool                             `json:"RetryEnabled,omitempty"`
	SyntheticMonitorId string                            `json:"SyntheticMonitorId"`
	Timeout            *int64                            `json:"Timeout,omitempty"`
	ValidationRules    *WebTestPropertiesValidationRules `json:"ValidationRules,omitempty"`
}
